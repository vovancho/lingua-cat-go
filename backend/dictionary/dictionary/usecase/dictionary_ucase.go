package usecase

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/dictionary/domain"
	"strings"
	"time"
)

var (
	DictsRandomCountError = errors.New("DICTIONARY_RANDOM_COUNT_INVALID")
)

type Timeout time.Duration

func NewDictionaryUseCase(dr domain.DictionaryRepository, v *validator.Validate, timeout Timeout) domain.DictionaryUseCase {
	return &dictionaryUseCase{
		dictionaryRepo: dr,
		validate:       v,
		contextTimeout: time.Duration(timeout),
	}
}

type dictionaryUseCase struct {
	dictionaryRepo domain.DictionaryRepository
	validate       *validator.Validate
	contextTimeout time.Duration
}

func (d dictionaryUseCase) GetByIDs(ctx context.Context, ids []domain.DictionaryID) ([]domain.Dictionary, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()

	dictionaries, err := d.dictionaryRepo.GetByIDs(ctx, ids)
	if err != nil {
		// Если это таймаут — не затираем ошибку
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		return nil, domain.DictNotFoundError
	}

	return dictionaries, nil
}

func (d dictionaryUseCase) GetRandomDictionaries(ctx context.Context, lang domain.DictionaryLang, limit uint8) ([]domain.Dictionary, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()

	if limit < 4 || limit > 8 {
		err := DictsRandomCountError
		return nil, err
	}

	if !lang.IsValid() {
		err := domain.DictLangInvalidError
		return nil, err
	}

	dicts, err := d.dictionaryRepo.GetRandomDictionaries(ctx, lang, limit)
	if err != nil {
		// Если это таймаут — не затираем ошибку
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		return nil, domain.DictsNotFoundError
	}
	if len(dicts) != int(limit) {
		return nil, domain.DictsNotFoundError
	}

	return dicts, nil
}

func (d dictionaryUseCase) Store(ctx context.Context, dict *domain.Dictionary) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()

	dict.Name = strings.ToLower(strings.TrimSpace(dict.Name))

	if len(dict.Translations) == 0 {
		return domain.DictTranslationRequiredError
	}

	for i, t := range dict.Translations {
		dict.Translations[i].Dictionary.Name = strings.ToLower(strings.TrimSpace(dict.Translations[i].Dictionary.Name))

		if t.Dictionary.Lang == dict.Lang {
			return domain.DictTranslationLangInvalidError
		}
	}

	if err := d.validate.Struct(dict); err != nil {
		return err
	}

	isExists, err := d.dictionaryRepo.IsExistsByNameAndLang(ctx, dict.Name, dict.Lang)
	if err != nil {
		return err
	}
	if isExists {
		return domain.DictExistsError
	}

	for _, t := range dict.Translations {
		isExists, err = d.dictionaryRepo.IsExistsByNameAndLang(ctx, t.Dictionary.Name, t.Dictionary.Lang)
		if err != nil {
			return err
		}
		if isExists {
			return domain.DictExistsError
		}
	}

	if err = d.dictionaryRepo.Store(ctx, dict); err != nil {
		return err
	}

	return nil
}

func (d dictionaryUseCase) ChangeName(ctx context.Context, id domain.DictionaryID, name string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()

	dictionaries, err := d.dictionaryRepo.GetByIDs(ctx, []domain.DictionaryID{id})
	if err != nil {
		// Если это таймаут — не затираем ошибку
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}
		return domain.DictNotFoundError
	}
	if len(dictionaries) == 0 {
		return domain.DictNotFoundError
	}
	dict := dictionaries[0]

	newDictName := strings.ToLower(strings.TrimSpace(name))
	if newDictName == dict.Name {
		return nil
	}
	dict.Name = newDictName

	if err = d.validate.StructPartial(dict, "Name"); err != nil {
		return err
	}

	isExists, err := d.dictionaryRepo.IsExistsByNameAndLang(ctx, dict.Name, dict.Lang)
	if err != nil {
		return err
	}
	if isExists {
		return domain.DictExistsError
	}

	if err = d.dictionaryRepo.ChangeName(ctx, id, dict.Name); err != nil {
		return err
	}

	return nil
}

func (d dictionaryUseCase) Delete(ctx context.Context, id domain.DictionaryID) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()

	// Получаем словарь по ID с использованием GetByIDs
	dictionaries, err := d.dictionaryRepo.GetByIDs(ctx, []domain.DictionaryID{id})
	if err != nil {
		// Если это таймаут — не затираем ошибку
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}
		return domain.DictNotFoundError
	}

	if len(dictionaries) == 0 {
		return domain.DictNotFoundError
	}

	if err := d.dictionaryRepo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

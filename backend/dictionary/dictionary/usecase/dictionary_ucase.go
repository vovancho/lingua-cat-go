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

func NewDictionaryUseCase(dr domain.DictionaryRepository, v *validator.Validate, timeout time.Duration) domain.DictionaryUseCase {
	return &dictionaryUseCase{
		dictionaryRepo: dr,
		validate:       v,
		contextTimeout: timeout,
	}
}

type dictionaryUseCase struct {
	dictionaryRepo domain.DictionaryRepository
	validate       *validator.Validate
	contextTimeout time.Duration
}

func (d dictionaryUseCase) GetByID(ctx context.Context, id domain.DictionaryID) (*domain.Dictionary, error) {
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()

	dictionary, err := d.dictionaryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.DictNotFoundError
	}

	return dictionary, nil
}

func (d dictionaryUseCase) GetRandomDictionaries(ctx context.Context, lang domain.DictionaryLang, count uint8) ([]domain.Dictionary, error) {
	if count < 4 || count > 8 {
		err := DictsRandomCountError
		return nil, err
	}

	if !lang.IsValid() {
		err := domain.DictLangInvalidError
		return nil, err
	}

	dicts, err := d.dictionaryRepo.GetRandomDictionaries(ctx, lang, count)
	if err != nil {
		return nil, domain.DictsNotFoundError
	}
	if len(dicts) != int(count) {
		return nil, domain.DictsNotFoundError
	}

	return dicts, nil
}

func (d dictionaryUseCase) Store(ctx context.Context, dict *domain.Dictionary) (err error) {
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

	if err = d.validate.Struct(dict); err != nil {
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

	err = d.dictionaryRepo.Store(ctx, dict)

	return
}

func (d dictionaryUseCase) ChangeName(ctx context.Context, id domain.DictionaryID, name string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()

	// имя в нижний регистр
	// проверить что сущность существует
	// изменить имя сущности и валидировать ее
	// проверить новое слово на существование (уникальность)

	err = d.dictionaryRepo.ChangeName(ctx, id, name)

	return
}

func (d dictionaryUseCase) Delete(ctx context.Context, id domain.DictionaryID) (err error) {
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()

	// проверить что сущность существует

	err = d.dictionaryRepo.Delete(ctx, id)

	return
}

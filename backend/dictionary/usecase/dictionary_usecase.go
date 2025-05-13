package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/dictionary/domain"
)

func NewDictionaryUseCase(repo domain.DictionaryRepository, validator *validator.Validate) domain.DictionaryUseCase {
	return &dictionaryUseCase{
		dictionaryRepo: repo,
		validator:      validator,
	}
}

type dictionaryUseCase struct {
	dictionaryRepo domain.DictionaryRepository
	validator      *validator.Validate
}

func (uc dictionaryUseCase) GetByIDs(ctx context.Context, ids []domain.DictionaryID) ([]domain.Dictionary, error) {
	dictionaries, err := uc.dictionaryRepo.GetByIDs(ctx, ids)
	if err != nil {
		return nil, domain.DictNotFoundError
	}

	return dictionaries, nil
}

func (uc dictionaryUseCase) GetRandomDictionaries(ctx context.Context, lang domain.DictionaryLang, limit uint8) ([]domain.Dictionary, error) {
	if limit < 4 || limit > 8 {
		return nil, domain.DictsRandomCountError
	}

	if !lang.IsValid() {
		return nil, domain.DictLangInvalidError
	}

	dicts, err := uc.dictionaryRepo.GetRandomDictionaries(ctx, lang, limit)
	if err != nil {
		return nil, domain.DictsNotFoundError
	}

	return dicts, nil
}

func (uc dictionaryUseCase) Store(ctx context.Context, dict *domain.Dictionary) error {
	normalizeDictionary(dict)

	if len(dict.Translations) == 0 {
		return domain.DictTranslationRequiredError
	}

	for _, t := range dict.Translations {
		if t.Dictionary.Lang == dict.Lang {
			return domain.DictTranslationLangInvalidError
		}
	}

	if err := uc.validator.Struct(dict); err != nil {
		return err
	}

	if err := uc.isDictExists(ctx, dict.Name, dict.Lang); err != nil {
		return err
	}

	for _, t := range dict.Translations {
		if err := uc.isDictExists(ctx, t.Dictionary.Name, t.Dictionary.Lang); err != nil {
			return err
		}
	}

	if err := uc.dictionaryRepo.Store(ctx, dict); err != nil {
		return err
	}

	return nil
}

func (uc dictionaryUseCase) ChangeName(ctx context.Context, id domain.DictionaryID, name string) error {
	dictionaries, err := uc.dictionaryRepo.GetByIDs(ctx, []domain.DictionaryID{id})
	if err != nil {
		return domain.DictNotFoundError
	}

	if len(dictionaries) == 0 {
		return domain.DictNotFoundError
	}

	dict := dictionaries[0]

	newDictName := strings.ToLower(strings.TrimSpace(name))
	// Ничего не делаем, если имя не изменилось
	if newDictName == dict.Name {
		return nil
	}
	dict.Name = newDictName

	if err := uc.validator.StructPartial(dict, "Name"); err != nil {
		return err
	}

	if err := uc.isDictExists(ctx, dict.Name, dict.Lang); err != nil {
		return err
	}

	if err := uc.dictionaryRepo.ChangeName(ctx, id, dict.Name); err != nil {
		return err
	}

	return nil
}

func (uc dictionaryUseCase) Delete(ctx context.Context, id domain.DictionaryID) error {
	// Получаем словарь по ID с использованием GetByIDs
	dictionaries, err := uc.dictionaryRepo.GetByIDs(ctx, []domain.DictionaryID{id})
	if err != nil {
		return domain.DictNotFoundError
	}

	if len(dictionaries) == 0 {
		return domain.DictNotFoundError
	}

	if err := uc.dictionaryRepo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (uc *dictionaryUseCase) isDictExists(ctx context.Context, name string, lang domain.DictionaryLang) error {
	exists, err := uc.dictionaryRepo.IsExistsByNameAndLang(ctx, name, lang)
	if err != nil {
		return fmt.Errorf("failed to check existence: %w", err)
	}
	if exists {
		return domain.DictExistsError
	}

	return nil
}

func normalizeDictionary(dict *domain.Dictionary) {
	dict.Name = strings.ToLower(strings.TrimSpace(dict.Name))

	for i, t := range dict.Translations {
		t.Dictionary.Name = strings.ToLower(strings.TrimSpace(t.Dictionary.Name))
		dict.Translations[i] = t
	}
}

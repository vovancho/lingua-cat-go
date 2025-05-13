package usecase

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
)

func NewDictionaryUseCase(repo domain.DictionaryRepository, validator *validator.Validate) domain.DictionaryUseCase {
	return &dictionaryUseCase{
		dictionaryRepo: repo,
		validate:       validator,
	}
}

type dictionaryUseCase struct {
	dictionaryRepo domain.DictionaryRepository
	validate       *validator.Validate
}

func (d dictionaryUseCase) GetRandomDictionaries(ctx context.Context, lang domain.DictionaryLang, limit uint8) ([]domain.Dictionary, error) {
	if limit < 4 || limit > 8 {
		return nil, domain.DictionariesLimitError
	}

	dicts, err := d.dictionaryRepo.GetRandomDictionaries(ctx, lang, limit)
	if err != nil {
		return nil, domain.DictionariesNotFoundError
	}

	return dicts, nil
}

func (d dictionaryUseCase) GetDictionariesByIds(ctx context.Context, dictIds []domain.DictionaryID) ([]domain.Dictionary, error) {
	dicts, err := d.dictionaryRepo.GetDictionariesByIds(ctx, dictIds)
	if err != nil {
		return nil, domain.DictionariesNotFoundError
	}

	return dicts, nil
}

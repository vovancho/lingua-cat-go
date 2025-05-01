package usecase

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	"time"
)

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

func (d dictionaryUseCase) GetRandomDictionaries(ctx context.Context, lang domain.DictionaryLang, limit uint8) ([]domain.Dictionary, error) {
	//TODO implement me
	panic("implement me")
}

func (d dictionaryUseCase) GetDictionaryByIds(ctx context.Context, dictIds []domain.DictionaryID) ([]domain.Dictionary, error) {
	//TODO implement me
	panic("implement me")
}

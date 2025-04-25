package usecase

import (
	"context"
	"github.com/vovancho/lingua-cat-go/dictionary/domain"
	"time"
)

func NewDictionaryUseCase(dr domain.DictionaryRepository, timeout time.Duration) domain.DictionaryUseCase {
	return &dictionaryUseCase{
		dictionaryRepo: dr,
		contextTimeout: timeout,
	}
}

type dictionaryUseCase struct {
	dictionaryRepo domain.DictionaryRepository
	contextTimeout time.Duration
}

func (d dictionaryUseCase) GetByID(ctx context.Context, id uint64) (domain.Dictionary, error) {
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()

	dictionary, err := d.dictionaryRepo.GetByID(ctx, id)

	if err != nil {
		return dictionary, err
	}

	return dictionary, nil
}

func (d dictionaryUseCase) Store(ctx context.Context, dict *domain.Dictionary) (err error) {
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()

	err = d.dictionaryRepo.Store(ctx, dict)

	return
}

func (d dictionaryUseCase) ChangeName(ctx context.Context, id uint64, name string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()

	err = d.dictionaryRepo.ChangeName(ctx, id, name)

	return
}

func (d dictionaryUseCase) Delete(ctx context.Context, id uint64) (err error) {
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()

	err = d.dictionaryRepo.Delete(ctx, id)

	return
}

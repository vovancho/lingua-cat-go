package usecase

import (
	"context"
	"errors"
	"github.com/vovancho/lingua-cat-go/dictionary/domain"
	"time"
)

var (
	DictsRandomCountError = errors.New("DICTIONARY_RANDOM_COUNT_INVALID")
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

	// валидировать сущность (+проверить что переводы другого языка)
	// проверить слово на существование (уникальность)

	err = d.dictionaryRepo.Store(ctx, dict)

	return
}

func (d dictionaryUseCase) ChangeName(ctx context.Context, id domain.DictionaryID, name string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, d.contextTimeout)
	defer cancel()

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

package grpc

import (
	"context"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
)

type postgresDictionaryRepository struct {
}

func NewPostgresDictionaryRepository() domain.DictionaryRepository {
	return &postgresDictionaryRepository{}
}

func (p postgresDictionaryRepository) GetRandomDictionaries(ctx context.Context, lang domain.DictionaryLang, limit uint8) ([]domain.Dictionary, error) {
	//TODO implement me
	panic("implement me")
}

func (p postgresDictionaryRepository) GetDictionaryByIds(ctx context.Context, dictIds []domain.DictionaryID) ([]domain.Dictionary, error) {
	//TODO implement me
	panic("implement me")
}

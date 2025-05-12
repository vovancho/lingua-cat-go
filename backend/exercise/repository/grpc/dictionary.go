package grpc

import (
	"context"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	"github.com/vovancho/lingua-cat-go/exercise/repository/grpc/gen"
	"github.com/vovancho/lingua-cat-go/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type dictionaryRepository struct {
	client dictionary.DictionaryServiceClient
	auth   *auth.AuthService
}

func NewDictionaryRepository(conn *grpc.ClientConn, auth *auth.AuthService) domain.DictionaryRepository {
	return &dictionaryRepository{
		client: dictionary.NewDictionaryServiceClient(conn),
		auth:   auth,
	}
}

func (r dictionaryRepository) GetRandomDictionaries(ctx context.Context, lang domain.DictionaryLang, limit uint8) ([]domain.Dictionary, error) {
	token, err := r.auth.GetJWTToken(ctx)
	if err != nil {
		return nil, err
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)

	resp, err := r.client.GetRandomDictionaries(ctx, &dictionary.GetRandomDictionariesRequest{
		Lang:  string(lang),
		Limit: int32(limit),
	})
	if err != nil {
		return nil, err
	}

	dictionaries := r.newDictionariesByGrpcResponse(resp.Dictionaries)

	return dictionaries, nil
}

func (r dictionaryRepository) GetDictionariesByIds(ctx context.Context, dictIds []domain.DictionaryID) ([]domain.Dictionary, error) {
	token, err := r.auth.GetJWTToken(ctx)
	if err != nil {
		return nil, err
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)

	intIDs := make([]int64, len(dictIds))
	for i, id := range dictIds {
		intIDs[i] = int64(id)
	}

	resp, err := r.client.GetDictionaries(ctx, &dictionary.GetDictionariesRequest{
		Ids: intIDs,
	})
	if err != nil {
		return nil, err
	}

	dictionaries := r.newDictionariesByGrpcResponse(resp.Dictionaries)

	return dictionaries, nil
}

func (r dictionaryRepository) newDictionariesByGrpcResponse(dicts []*dictionary.Dictionary) []domain.Dictionary {
	var dictionaries []domain.Dictionary
	for _, dt := range dicts {
		dict := domain.Dictionary{
			ID:   domain.DictionaryID(dt.Id),
			Lang: domain.DictionaryLang(dt.Lang),
			Name: dt.Name,
			Type: domain.DictionaryType(dt.Type),
		}

		for _, t := range dt.Translations {
			var dSentences []domain.Sentence
			for _, ts := range t.Sentences {
				dSentences = append(dSentences, domain.Sentence{
					TextRU: ts.TextRu,
					TextEN: ts.TextEn,
				})
			}

			dict.Translations = append(dict.Translations, domain.Translation{
				Dictionary: domain.Dictionary{
					ID:        domain.DictionaryID(t.Id),
					Lang:      domain.DictionaryLang(t.Lang),
					Name:      t.Name,
					Type:      domain.DictionaryType(t.Type),
					Sentences: dSentences,
				},
			})
		}

		for _, s := range dt.Sentences {
			dict.Sentences = append(dict.Sentences, domain.Sentence{
				TextRU: s.TextRu,
				TextEN: s.TextEn,
			})
		}

		dictionaries = append(dictionaries, dict)
	}

	return dictionaries
}

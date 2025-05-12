package grpc

import (
	"context"
	"fmt"

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
	ctx, err := r.withAuthContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get auth context: %w", err)
	}

	req := &dictionary.GetRandomDictionariesRequest{Lang: string(lang), Limit: int32(limit)}
	resp, err := r.client.GetRandomDictionaries(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("grpc call GetRandomDictionaries: %w", err)
	}

	dictionaries := r.newDictionariesByGrpcResponse(resp.Dictionaries)

	return dictionaries, nil
}

func (r dictionaryRepository) GetDictionariesByIds(ctx context.Context, dictIds []domain.DictionaryID) ([]domain.Dictionary, error) {
	ctx, err := r.withAuthContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get auth context: %w", err)
	}

	intIDs := make([]int64, len(dictIds))
	for i, id := range dictIds {
		intIDs[i] = int64(id)
	}

	req := &dictionary.GetDictionariesRequest{Ids: intIDs}
	resp, err := r.client.GetDictionaries(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("grpc call GetDictionaries: %w", err)
	}

	dictionaries := r.newDictionariesByGrpcResponse(resp.Dictionaries)

	return dictionaries, nil
}

func (r dictionaryRepository) newDictionariesByGrpcResponse(dicts []*dictionary.Dictionary) []domain.Dictionary {
	var dictionaries []domain.Dictionary
	for _, dt := range dicts {
		dict := domain.Dictionary{
			ID:           domain.DictionaryID(dt.Id),
			Lang:         domain.DictionaryLang(dt.Lang),
			Name:         dt.Name,
			Type:         domain.DictionaryType(dt.Type),
			Translations: r.newTranslations(dt.Translations),
			Sentences:    r.newSentences(dt.Sentences),
		}

		dictionaries = append(dictionaries, dict)
	}

	return dictionaries
}

func (r dictionaryRepository) newTranslations(translations []*dictionary.Translation) []domain.Translation {
	result := make([]domain.Translation, 0, len(translations))
	for _, t := range translations {
		result = append(result, domain.Translation{
			Dictionary: domain.Dictionary{
				ID:        domain.DictionaryID(t.Id),
				Lang:      domain.DictionaryLang(t.Lang),
				Name:      t.Name,
				Type:      domain.DictionaryType(t.Type),
				Sentences: r.newSentences(t.Sentences),
			},
		})
	}

	return result
}

func (r dictionaryRepository) newSentences(sentences []*dictionary.Sentence) []domain.Sentence {
	result := make([]domain.Sentence, 0, len(sentences))
	for _, s := range sentences {
		result = append(result, domain.Sentence{
			TextRU: s.TextRu,
			TextEN: s.TextEn,
		})
	}

	return result
}

func (r dictionaryRepository) withAuthContext(ctx context.Context) (context.Context, error) {
	token, err := r.auth.GetJWTToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("get JWT token: %w", err)
	}

	return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token), nil
}

package grpc

import (
	"context"

	pb "github.com/vovancho/lingua-cat-go/dictionary/delivery/grpc/gen"

	"github.com/vovancho/lingua-cat-go/dictionary/domain"
)

type DictionaryHandler struct {
	dictionaryUseCase domain.DictionaryUseCase
}

func NewDictionaryHandler(dictionaryUseCase domain.DictionaryUseCase) *DictionaryHandler {
	return &DictionaryHandler{
		dictionaryUseCase: dictionaryUseCase,
	}
}

func (h *DictionaryHandler) GetRandomDictionaries(ctx context.Context, req *pb.GetRandomDictionariesRequest) (*pb.GetRandomDictionariesResponse, error) {
	lang := domain.DictionaryLang(req.GetLang())
	limit := uint8(req.GetLimit())

	dictionaries, err := h.dictionaryUseCase.GetRandomDictionaries(ctx, lang, limit)
	if err != nil {
		return nil, err
	}

	protoDictionaries := buildProtoDictionaries(dictionaries)

	return &pb.GetRandomDictionariesResponse{Dictionaries: protoDictionaries}, nil
}

func (h *DictionaryHandler) GetDictionaries(ctx context.Context, req *pb.GetDictionariesRequest) (*pb.GetDictionariesResponse, error) {
	rawIds := req.GetIds()
	dictionaryIds := make([]domain.DictionaryID, len(rawIds))
	for i, id := range rawIds {
		dictionaryIds[i] = domain.DictionaryID(id)
	}

	dictionaries, err := h.dictionaryUseCase.GetByIDs(ctx, dictionaryIds)
	if err != nil {
		return nil, err
	}

	protoDictionaries := buildProtoDictionaries(dictionaries)

	return &pb.GetDictionariesResponse{Dictionaries: protoDictionaries}, nil
}

// Вспомогательная функция для сборки protobuf-словарей из доменных
func buildProtoDictionaries(dictionaries []domain.Dictionary) []*pb.Dictionary {
	protoDictionaries := make([]*pb.Dictionary, len(dictionaries))
	for i, dict := range dictionaries {
		translations := make([]*pb.Translation, len(dict.Translations))
		for j, t := range dict.Translations {
			tSentences := make([]*pb.Sentence, len(t.Dictionary.Sentences))
			for k, ts := range t.Dictionary.Sentences {
				tSentences[k] = &pb.Sentence{
					TextRu: ts.TextRU,
					TextEn: ts.TextEN,
				}
			}

			translations[j] = &pb.Translation{
				Id:        int64(t.Dictionary.ID),
				Lang:      string(t.Dictionary.Lang),
				Name:      t.Dictionary.Name,
				Type:      int32(t.Dictionary.Type),
				Sentences: tSentences,
			}
		}

		sentences := make([]*pb.Sentence, len(dict.Sentences))
		for j, s := range dict.Sentences {
			sentences[j] = &pb.Sentence{
				TextRu: s.TextRU,
				TextEn: s.TextEN,
			}
		}

		protoDictionaries[i] = &pb.Dictionary{
			Id:           int64(dict.ID),
			Lang:         string(dict.Lang),
			Name:         dict.Name,
			Type:         int32(dict.Type),
			Translations: translations,
			Sentences:    sentences,
		}
	}

	return protoDictionaries
}

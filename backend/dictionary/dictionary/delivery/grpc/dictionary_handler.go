package grpc

import (
	"context"
	"github.com/go-playground/validator/v10"
	pb "github.com/vovancho/lingua-cat-go/dictionary/dictionary/delivery/grpc/gen"
	"github.com/vovancho/lingua-cat-go/dictionary/domain"
)

type DictionaryHandler struct {
	DUseCase domain.DictionaryUseCase
	validate *validator.Validate
}

func NewDictionaryHandler(v *validator.Validate, d domain.DictionaryUseCase) *DictionaryHandler {
	return &DictionaryHandler{
		DUseCase: d,
		validate: v,
	}
}

func (d *DictionaryHandler) GetRandomDictionaries(ctx context.Context, req *pb.GetRandomDictionariesRequest) (*pb.GetRandomDictionariesResponse, error) {
	lang := domain.DictionaryLang(req.GetLang())
	limit := uint8(req.GetLimit())

	dictionaries, err := d.DUseCase.GetRandomDictionaries(ctx, lang, limit)
	if err != nil {
		return nil, err
	}

	protoDictionaries := make([]*pb.Dictionary, len(dictionaries))
	for i, dict := range dictionaries {
		translations := make([]*pb.Translation, len(dict.Translations))
		for j, t := range dict.Translations {
			translations[j] = &pb.Translation{
				Lang: string(t.Dictionary.Lang),
				Name: t.Dictionary.Name,
				Type: int32(t.Dictionary.Type),
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

	return &pb.GetRandomDictionariesResponse{
		Dictionaries: protoDictionaries,
	}, nil
}

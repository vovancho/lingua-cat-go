package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/dictionary/domain"
	_internalError "github.com/vovancho/lingua-cat-go/dictionary/internal/error"
	"github.com/vovancho/lingua-cat-go/dictionary/internal/request"
	"github.com/vovancho/lingua-cat-go/dictionary/internal/response"
	"net/http"
)

type DictionaryStoreRequest struct {
	Lang         string `json:"lang" validate:"required,len=2"`
	Name         string `json:"name" validate:"required,min=2"`
	Type         uint16 `json:"type" validate:"required,oneof=1 2 3"`
	Translations []struct {
		Lang string `json:"lang" validate:"required,len=2"`
		Name string `json:"name" validate:"required,min=2"`
		Type uint16 `json:"type" validate:"required,oneof=1 2 3"`
	} `json:"translations"`
	Sentences []struct {
		TextRU string `json:"text_ru" validate:"required,min=5"`
		TextEN string `json:"text_en" validate:"required,min=5"`
	} `json:"sentences"`
}

type DictionaryChangeNameRequest struct {
	Name string `json:"name" validate:"required,min=2"`
}

type DictionaryHandler struct {
	DUseCase domain.DictionaryUseCase
	validate *validator.Validate
}

func NewDictionaryHandler(router *http.ServeMux, v *validator.Validate, d domain.DictionaryUseCase) {
	handler := &DictionaryHandler{
		DUseCase: d,
		validate: v,
	}

	router.HandleFunc("GET /dictionary/{id}", request.WithID(handler.GetByID))
	router.HandleFunc("POST /dictionary", handler.Store)
	router.HandleFunc("POST /dictionary/{id}/name", request.WithID(handler.ChangeName))
	router.HandleFunc("DELETE /dictionary/{id}", request.WithID(handler.Delete))
}

func (d *DictionaryHandler) GetByID(w http.ResponseWriter, r *http.Request, id uint64) {
	dictionary, err := d.DUseCase.GetByID(r.Context(), domain.DictionaryID(id))
	if err != nil {
		response.Error(err, r)
		return
	}

	response.JSON(w, http.StatusOK, response.APIResponse{
		Message: "Dictionary got successfully",
		Data:    dictionary,
	})
}

func (d *DictionaryHandler) Store(w http.ResponseWriter, r *http.Request) {
	var requestBody DictionaryStoreRequest
	if err := d.validateRequest(r, &requestBody); err != nil {
		response.Error(err, r)
		return
	}

	dictionary := newDictionaryByRequest(requestBody)
	if err := d.DUseCase.Store(r.Context(), &dictionary); err != nil {
		response.Error(err, r)
		return
	}

	response.JSON(w, http.StatusCreated, response.APIResponse{
		Message: "Dictionary created successfully",
		Data:    dictionary,
	})
}

func (d *DictionaryHandler) ChangeName(w http.ResponseWriter, r *http.Request, id uint64) {
	var requestBody DictionaryChangeNameRequest
	if err := d.validateRequest(r, &requestBody); err != nil {
		response.Error(err, r)
		return
	}

	if err := d.DUseCase.ChangeName(r.Context(), domain.DictionaryID(id), requestBody.Name); err != nil {
		response.Error(err, r)
		return
	}

	response.JSON(w, http.StatusOK, response.APIResponse{
		Message: "Dictionary name changed successfully",
	})
}

func (d *DictionaryHandler) Delete(w http.ResponseWriter, r *http.Request, id uint64) {
	if err := d.DUseCase.Delete(r.Context(), domain.DictionaryID(id)); err != nil {
		response.Error(err, r)
		return
	}

	response.JSON(w, http.StatusOK, response.APIResponse{
		Message: "Dictionary deleted successfully",
	})
}

func (d *DictionaryHandler) validateRequest(r *http.Request, req any) error {
	if err := request.FromJSON(r, req); err != nil {
		return _internalError.NewAppError(http.StatusBadRequest, "Некорректный синтаксис JSON", _internalError.InvalidDecodeJsonError)
	}

	if err := d.validate.Struct(req); err != nil {
		return err
	}

	return nil
}

func newDictionaryByRequest(req DictionaryStoreRequest) domain.Dictionary {
	dictionary := domain.Dictionary{
		Name: req.Name,
		Type: req.Type,
		Lang: req.Lang,
	}

	for _, s := range req.Sentences {
		sentence := domain.Sentence{
			TextRU: s.TextRU, // Русский перевод
			TextEN: s.TextEN, // Английский текст
		}
		dictionary.Sentences = append(dictionary.Sentences, sentence)
	}

	for _, t := range req.Translations {
		transDict := domain.Dictionary{
			Name: t.Name,
			Type: t.Type,
			Lang: t.Lang,
		}
		translation := domain.Translation{
			Dictionary: transDict,
		}
		dictionary.Translations = append(dictionary.Translations, translation)
	}

	return dictionary
}

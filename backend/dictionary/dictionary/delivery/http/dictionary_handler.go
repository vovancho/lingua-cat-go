package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/dictionary/domain"
	"log/slog"
	"net/http"
	"strconv"
)

type APIResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type ValidationError struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

func (ve ValidationError) Error() string {
	return ve.Message
}

func newValidationError(err error) ValidationError {
	ve := ValidationError{
		Message: "Validation error",
		Errors:  make(map[string]string),
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range validationErrors {
			field := fe.StructField()
			tag := fe.Tag()

			// Кастомные сообщения для разных тегов
			switch tag {
			case "required":
				ve.Errors[field] = "This field is required"
			case "min":
				ve.Errors[field] = fmt.Sprintf("Minimum length is %s", fe.Param())
			case "len":
				ve.Errors[field] = fmt.Sprintf("Must be exactly %s characters", fe.Param())
			default:
				ve.Errors[field] = fmt.Sprintf("Field validation failed: %s", tag)
			}
		}
	}

	return ve
}

type DictionaryHandler struct {
	DUseCase  domain.DictionaryUseCase
	validator *validator.Validate
}

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

func NewDictionaryHandler(router *http.ServeMux, v *validator.Validate, d domain.DictionaryUseCase) {
	handler := &DictionaryHandler{
		DUseCase:  d,
		validator: v,
	}

	router.HandleFunc("GET /dictionary/{id}", handler.GetByID)
	router.HandleFunc("POST /dictionary", handler.Store)
	router.HandleFunc("POST /dictionary/{id}/name", handler.ChangeName)
	router.HandleFunc("DELETE /dictionary/{id}", handler.Delete)
}

func (d *DictionaryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dictionary, err := d.DUseCase.GetByID(r.Context(), domain.DictionaryID(id))
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, `{"message":"Failed to get dictionary"}`, http.StatusInternalServerError)

		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Message: "Dictionary got successfully",
		Data:    dictionary,
	})
}

func (d *DictionaryHandler) Store(w http.ResponseWriter, r *http.Request) {
	var requestBody DictionaryStoreRequest

	if err := d.validateRequest(r, &requestBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	dictionary := mapStoreRequestToDomain(requestBody)

	if err := d.DUseCase.Store(r.Context(), &dictionary); err != nil {
		slog.Error(err.Error())
		respondWithError(w, http.StatusInternalServerError, "Failed to store dictionary", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, APIResponse{
		Message: "Dictionary created successfully",
		Data:    dictionary,
	})
}

func (d *DictionaryHandler) ChangeName(w http.ResponseWriter, r *http.Request) {
	var requestBody DictionaryChangeNameRequest

	if err := d.validateRequest(r, &requestBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	idString := r.PathValue("id")
	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := d.DUseCase.ChangeName(r.Context(), domain.DictionaryID(id), requestBody.Name); err != nil {
		slog.Error(err.Error())
		respondWithError(w, http.StatusInternalServerError, "Failed to change dictionary name", err)
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Message: "Dictionary name changed successfully",
	})
}

func (d *DictionaryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := d.DUseCase.Delete(r.Context(), domain.DictionaryID(id)); err != nil {
		slog.Error(err.Error())
		http.Error(w, `{"message":"Failed to delete dictionary"}`, http.StatusInternalServerError)

		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Message: "Dictionary deleted successfully",
	})
}

func respondWithError(w http.ResponseWriter, code int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := map[string]interface{}{
		"message": message,
	}

	if err != nil {
		response["error"] = err.Error()
	}

	json.NewEncoder(w).Encode(response)
}

func (d *DictionaryHandler) validateRequest(r *http.Request, req any) error {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}
	defer r.Body.Close()

	if err := d.validator.Struct(req); err != nil {
		return newValidationError(err)
	}

	return nil
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, `{"message":"Failed to encode response"}`, http.StatusInternalServerError)
	}
}

func mapStoreRequestToDomain(req DictionaryStoreRequest) domain.Dictionary {
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

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

type ResponseError struct {
	Message string `json:"message"`
}

type ValidationError struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
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
	DUseCase domain.DictionaryUseCase
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

func NewDictionaryHandler(router *http.ServeMux, d domain.DictionaryUseCase) {
	handler := &DictionaryHandler{
		DUseCase: d,
	}

	router.HandleFunc("GET /dictionary/{id}", handler.GetByID)
	router.HandleFunc("POST /dictionary", handler.Store)
	router.HandleFunc("POST /dictionary/{id}/name", handler.ChangeName)
	router.HandleFunc("DELETE /dictionary/{id}", handler.Delete)
	router.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(map[string]any{
			"message": "success",
			"host":    r.Host,
			"path":    r.URL.Path,
		}); err != nil {
			http.Error(w, `{"message":"Failed to encode response"}`, http.StatusInternalServerError)

			return
		}
	})
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]any{
		"message":    "Dictionary got successfully",
		"dictionary": dictionary,
	}); err != nil {
		http.Error(w, `{"message":"Failed to encode response"}`, http.StatusInternalServerError)

		return
	}
}

func (d *DictionaryHandler) Store(w http.ResponseWriter, r *http.Request) {
	var requestBody DictionaryStoreRequest

	// Декодируем JSON из тела запроса в структуру
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		http.Error(w, `{"message":"Invalid JSON format"}`, http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		respondWithValidationError(w, newValidationError(err))
		return
	}

	dictionary := domain.Dictionary{
		Name: requestBody.Name,
		Type: requestBody.Type,
		Lang: requestBody.Lang,
	}

	for _, s := range requestBody.Sentences {
		sentence := domain.Sentence{
			TextRU: s.TextRU, // Русский перевод
			TextEN: s.TextEN, // Английский текст
		}
		dictionary.Sentences = append(dictionary.Sentences, sentence)
	}

	for _, t := range requestBody.Translations {
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

	w.Header().Set("Content-Type", "application/json")
	if err := d.DUseCase.Store(r.Context(), &dictionary); err != nil {
		slog.Error(err.Error())
		http.Error(w, `{"message":"Failed to store dictionary"}`, http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(map[string]any{
		"message":    "Dictionary created successfully",
		"dictionary": dictionary,
	}); err != nil {
		http.Error(w, `{"message":"Failed to encode response"}`, http.StatusInternalServerError)

		return
	}
}

func (d *DictionaryHandler) ChangeName(w http.ResponseWriter, r *http.Request) {
	var requestBody DictionaryChangeNameRequest

	// Декодируем JSON из тела запроса в структуру
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		http.Error(w, `{"message":"Invalid JSON format"}`, http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		respondWithValidationError(w, newValidationError(err))
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
		http.Error(w, `{"message":"Failed to change dictionary name"}`, http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]any{
		"message": "Dictionary name changed successfully",
	}); err != nil {
		http.Error(w, `{"message":"Failed to encode response"}`, http.StatusInternalServerError)

		return
	}
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]any{
		"message": "Dictionary deleted successfully",
	}); err != nil {
		http.Error(w, `{"message":"Failed to encode response"}`, http.StatusInternalServerError)

		return
	}
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

func respondWithValidationError(w http.ResponseWriter, ve ValidationError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(w).Encode(ve)
}

package http

import (
	"encoding/json"
	"github.com/vovancho/lingua-cat-go/dictionary/domain"
	"log/slog"
	"net/http"
)

type ResponseError struct {
	Message string `json:"message"`
}

type DictionaryHandler struct {
	DUseCase domain.DictionaryUseCase
}

type DictionaryStoreRequest struct {
	Lang         string `json:"lang"`
	Name         string `json:"name"`
	Type         uint16 `json:"type"`
	Translations []struct {
		Lang string `json:"lang"`
		Name string `json:"name"`
		Type uint16 `json:"type"`
	} `json:"translations"`
	Sentences []struct {
		TextRU string `json:"text_ru"`
		TextEN string `json:"text_en"`
	} `json:"sentences"`
}

func NewDictionaryHandler(router *http.ServeMux, d domain.DictionaryUseCase) {
	handler := &DictionaryHandler{
		DUseCase: d,
	}

	router.HandleFunc("POST /dictionary", handler.Store)
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

func (d *DictionaryHandler) Store(w http.ResponseWriter, r *http.Request) {
	var requestBody DictionaryStoreRequest

	// Декодируем JSON из тела запроса в структуру
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		http.Error(w, `{"message":"Invalid JSON format"}`, http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

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

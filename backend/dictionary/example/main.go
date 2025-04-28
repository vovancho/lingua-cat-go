package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	rutranslations "github.com/go-playground/validator/v10/translations/ru"
	"net/http"
	"strconv"
)

type APIResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

type BodyRequest struct {
	Name string `json:"name" validate:"required,min=5"`
	Age  uint8  `json:"age" validate:"required,min=18,max=100"`
}

type ErrorKey struct{}

type ValidationErrorItem struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func main() {
	validate := validator.New()

	// Инициализация переводчика
	ruLocale := ru.New()
	uni := ut.New(ruLocale, ruLocale)

	var ok bool
	trans, ok := uni.GetTranslator("ru")
	if !ok {
		panic("не удалось получить переводчик для ru")
	}

	if err := rutranslations.RegisterDefaultTranslations(validate, trans); err != nil {
		panic("не удалось зарегистрировать русские переводы: " + err.Error())
	}

	router := http.NewServeMux()

	router.HandleFunc("POST /example1/{id}",
		func(w http.ResponseWriter, r *http.Request) {
			idString := r.PathValue("id")
			id, err := strconv.ParseUint(idString, 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			fmt.Println("Route ID:", id)

			var requestBody BodyRequest

			if err := validateRequest(r, validate, &requestBody); err != nil {
				ctx := context.WithValue(r.Context(), ErrorKey{}, err)
				*r = *r.WithContext(ctx)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			respondWithJSON(w, http.StatusOK, APIResponse{
				Message: "Response successfully",
			})
		},
	)

	server := http.Server{
		Addr:    ":80",
		Handler: errorHandleMiddleware(router, trans),
	}
	fmt.Println("Test Server is listening on port 80")
	server.ListenAndServe()
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, `{"message":"Failed to encode response"}`, http.StatusInternalServerError)
	}
}

func validateRequest(r *http.Request, v *validator.Validate, req any) error {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}
	defer r.Body.Close()

	if err := v.Struct(req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			return errs
		}
		return err
	}

	return nil
}

func errorHandleMiddleware(next http.Handler, trans ut.Translator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Вызываем следующий обработчик в цепочке
		next.ServeHTTP(w, r)

		// Проверяем, есть ли ошибка в контексте
		if err, ok := r.Context().Value(ErrorKey{}).(error); ok && err != nil {
			switch e := err.(type) {
			case validator.ValidationErrors:
				var details []ValidationErrorItem
				for _, ei := range e {
					details = append(details, ValidationErrorItem{
						Field:   ei.Field(),
						Message: ei.Translate(trans),
					})
				}

				respondWithJSON(w, http.StatusBadRequest, APIResponse{
					Message: "Ошибка валидации",
					Data:    details,
				})
			case error:
				fmt.Println("Error:", e)
				respondWithJSON(w, http.StatusInternalServerError, APIResponse{
					Message: "Внутренняя ошибка сервера",
					Error:   e.Error(),
				})
			}
		}
	})
}

package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	_internalError "github.com/vovancho/lingua-cat-go/pkg/error"
)

type Responder interface {
	Error(w http.ResponseWriter, err error)
	Success(w http.ResponseWriter, status int, data any)
}

type responder struct {
	trans ut.Translator
}

func NewResponder(trans ut.Translator) Responder {
	return &responder{trans: trans}
}

func (r *responder) Success(w http.ResponseWriter, status int, data any) {
	JSON(w, status, APIResponse{
		Data: data,
	})
}

func (r *responder) Error(w http.ResponseWriter, err error) {
	HandleError(w, err, r.trans)
}

type APIResponse struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func (r APIResponse) isEmpty() bool {
	return r.Message == "" && r.Data == nil && r.Error == ""
}

func JSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if apiResp, ok := payload.(APIResponse); ok && apiResp.isEmpty() {
		return
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, `{"message":"Failed to encode response"}`, http.StatusInternalServerError)
	}
}

func HandleError(w http.ResponseWriter, err error, trans ut.Translator) {
	switch e := err.(type) {
	case validator.ValidationErrors:
		var details []ValidationErrorItem
		for _, ei := range e {
			var message string
			if trans != nil {
				message = ei.Translate(trans)
			} else {
				message = ei.Error()
			}

			details = append(details, ValidationErrorItem{
				Field:   formatFieldName(ei.Namespace()),
				Message: message,
			})
		}

		JSON(w, http.StatusBadRequest, APIResponse{
			Message: "Ошибка валидации",
			Data:    details,
		})
	case _internalError.AppErrorInterface:
		var unwrapped string
		if errUnwrapped := e.Unwrap(); errUnwrapped != nil {
			unwrapped = errUnwrapped.Error()
		}

		JSON(w, e.StatusCode(), APIResponse{
			Message: e.Error(),
			Error:   unwrapped,
		})
	default:
		slog.Error("unhandled error", "error", err)

		JSON(w, http.StatusInternalServerError, APIResponse{
			Message: "Внутренняя ошибка сервера",
			Error:   err.Error(),
		})
	}
}

func formatFieldName(namespace string) string {
	parts := strings.SplitN(namespace, ".", 2)
	if len(parts) < 2 {
		return namespace // Если нет точек, возвращаем как есть
	}

	return parts[1]
}

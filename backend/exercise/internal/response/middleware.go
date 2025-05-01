package response

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	_internalError "github.com/vovancho/lingua-cat-go/exercise/internal/error"
	"log/slog"
	"net/http"
	"strings"
)

type ValidationErrorItem struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorKey struct{}

func ErrorMiddleware(next http.Handler, trans ut.Translator) http.Handler {
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
						Field:   formatFieldName(ei.Namespace()),
						Message: ei.Translate(trans),
					})
				}

				JSON(w, http.StatusBadRequest, APIResponse{
					Message: "Ошибка валидации",
					Data:    details,
				})
			case _internalError.AppErrorInterface:
				JSON(w, e.StatusCode(), APIResponse{
					Message: e.Error(),
					Error:   e.Unwrap().Error(),
				})
			case error:
				slog.Error("unhandled error", "error", e)
				JSON(w, http.StatusInternalServerError, APIResponse{
					Message: "Внутренняя ошибка сервера",
					Error:   e.Error(),
				})
			}
		}
	})
}

func formatFieldName(namespace string) string {
	parts := strings.SplitN(namespace, ".", 2)
	if len(parts) < 2 {
		return namespace // Если нет точек, возвращаем как есть
	}
	return parts[1]
}

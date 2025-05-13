package response

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

type ValidationErrorItem struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorKey struct{}

func ErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recovered", "error", rec, "stack", string(debug.Stack()))

				JSON(w, http.StatusInternalServerError, APIResponse{
					Message: "Внутренняя ошибка сервера",
					Error:   "panic: см. логи сервера",
				})
			}
		}()

		next.ServeHTTP(w, r)
	})
}

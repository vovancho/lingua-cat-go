package response

import (
	"context"
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func JSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, `{"message":"Failed to encode response"}`, http.StatusInternalServerError)
	}
}

func Error(err error, r *http.Request) {
	ctx := context.WithValue(r.Context(), ErrorKey{}, err)
	*r = *r.WithContext(ctx)
}

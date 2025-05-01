package response

import (
	"context"
	"encoding/json"
	"net/http"
)

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

func Error(err error, r *http.Request) {
	ctx := context.WithValue(r.Context(), ErrorKey{}, err)
	*r = *r.WithContext(ctx)
}

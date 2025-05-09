package request

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	_internalError "github.com/vovancho/lingua-cat-go/analytics/internal/error"
	"github.com/vovancho/lingua-cat-go/analytics/internal/response"
	"github.com/vovancho/lingua-cat-go/pkg/auth"
	"net/http"
)

func FromJSON(r *http.Request, dto any) error {
	if err := json.NewDecoder(r.Body).Decode(dto); err != nil {
		return err
	}
	defer r.Body.Close()

	return nil
}

type HandlerFuncWithID func(w http.ResponseWriter, r *http.Request, id auth.UserID)

func WithID(h HandlerFuncWithID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")

		id, err := uuid.Parse(idString)
		if err != nil {
			messageErr := fmt.Sprintf("Парсинг \"%s\": некорректный формат UUID", idString)
			appErr := _internalError.NewAppError(http.StatusBadRequest, messageErr, _internalError.InvalidPathParamError)
			response.Error(appErr, r)
			return
		}

		h(w, r, auth.UserID(id))
	}
}

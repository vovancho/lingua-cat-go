package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/vovancho/lingua-cat-go/pkg/auth"
	_internalError "github.com/vovancho/lingua-cat-go/pkg/error"
	"github.com/vovancho/lingua-cat-go/pkg/response"
)

func FromJSON(r *http.Request, dto any) error {
	if err := json.NewDecoder(r.Body).Decode(dto); err != nil {
		return err
	}
	defer r.Body.Close()

	return nil
}

type HandlerFuncWithID func(w http.ResponseWriter, r *http.Request, id uint64)
type HandlerFuncWithUserID func(w http.ResponseWriter, r *http.Request, id auth.UserID)

func WithID(h HandlerFuncWithID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			messageErr := fmt.Sprintf("Парсинг \"%s\": некорректный формат ID", idString)
			appErr := _internalError.NewAppError(http.StatusBadRequest, messageErr, _internalError.InvalidPathParamError)
			response.Error(appErr, r)
			return
		}

		h(w, r, id)
	}
}

func WithUserID(h HandlerFuncWithUserID) http.HandlerFunc {
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

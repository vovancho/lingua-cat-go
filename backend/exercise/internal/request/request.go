package request

import (
	"encoding/json"
	"fmt"
	_internalError "github.com/vovancho/lingua-cat-go/exercise/internal/error"
	"github.com/vovancho/lingua-cat-go/exercise/internal/response"
	"net/http"
	"strconv"
)

func FromJSON(r *http.Request, dto any) error {
	if err := json.NewDecoder(r.Body).Decode(dto); err != nil {
		return err
	}
	defer r.Body.Close()

	return nil
}

type HandlerFuncWithID func(w http.ResponseWriter, r *http.Request, id uint64)

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

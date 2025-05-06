package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/analytics/domain"
	"github.com/vovancho/lingua-cat-go/analytics/internal/auth"
	_internalError "github.com/vovancho/lingua-cat-go/analytics/internal/error"
	"github.com/vovancho/lingua-cat-go/analytics/internal/request"
	"github.com/vovancho/lingua-cat-go/analytics/internal/response"
	"net/http"
)

type ExerciseCompleteHandler struct {
	ECUseCase domain.ExerciseCompleteUseCase
	validate  *validator.Validate
	auth      *auth.AuthService
}

func NewExerciseCompleteHandler(router *http.ServeMux, v *validator.Validate, auth *auth.AuthService, ec domain.ExerciseCompleteUseCase) {
	handler := &ExerciseCompleteHandler{
		ECUseCase: ec,
		validate:  v,
		auth:      auth,
	}

	router.HandleFunc("GET /v1/analytics/user/{id}", request.WithID(handler.GetByUserID))
}

func (ec *ExerciseCompleteHandler) GetByUserID(w http.ResponseWriter, r *http.Request, id auth.UserID) {
	exerciseCompleteList, err := ec.ECUseCase.GetByUserID(r.Context(), id)
	if err != nil {
		appError := _internalError.NewAppError(http.StatusNotFound, "Аналитика не найдена", err)
		response.Error(appError, r)
		return
	}

	response.JSON(w, http.StatusOK, response.APIResponse{
		Data: map[string]any{
			"analytics": exerciseCompleteList,
		},
	})
}

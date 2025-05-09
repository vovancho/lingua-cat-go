package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/analytics/domain"
	_internalError "github.com/vovancho/lingua-cat-go/analytics/internal/error"
	"github.com/vovancho/lingua-cat-go/analytics/internal/request"
	"github.com/vovancho/lingua-cat-go/analytics/internal/response"
	"github.com/vovancho/lingua-cat-go/pkg/auth"
	"net/http"
)

type AnalyticsData struct {
	Analytics []domain.ExerciseComplete `json:"analytics"`
}

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

// GetByUserID godoc
// @Summary Получить аналитику по пользователю
// @Description Получает аналитику завершенных упражнений для указанного пользователя.
// @Security BearerAuth
// @Tags Analytics
// @Param id path string true "ID пользователя" format(uuid)
// @Success 200 {object} response.APIResponse{data=AnalyticsData} "Аналитика найдена"
// @Failure 404 {object} response.APIResponse "Аналитика не найдена"
// @Router /v1/analytics/user/{id} [get]
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

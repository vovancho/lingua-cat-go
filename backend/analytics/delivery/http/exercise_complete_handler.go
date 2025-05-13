package http

import (
	"net/http"

	"github.com/vovancho/lingua-cat-go/analytics/domain"
	"github.com/vovancho/lingua-cat-go/pkg/auth"
	_internalError "github.com/vovancho/lingua-cat-go/pkg/error"
	"github.com/vovancho/lingua-cat-go/pkg/request"
	"github.com/vovancho/lingua-cat-go/pkg/response"
)

type AnalyticsData struct {
	Analytics []domain.ExerciseComplete `json:"analytics"`
}

type exerciseCompleteHandler struct {
	responder               response.Responder
	exerciseCompleteUseCase domain.ExerciseCompleteUseCase
	auth                    *auth.AuthService
}

func NewExerciseCompleteHandler(
	router *http.ServeMux,
	responder response.Responder,
	exerciseCompleteUseCase domain.ExerciseCompleteUseCase,
	auth *auth.AuthService,
) {
	handler := &exerciseCompleteHandler{
		responder:               responder,
		exerciseCompleteUseCase: exerciseCompleteUseCase,
		auth:                    auth,
	}

	router.HandleFunc("GET /v1/analytics/user/{id}", request.WithUserID(handler.GetByUserID))
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
func (h *exerciseCompleteHandler) GetByUserID(w http.ResponseWriter, r *http.Request, id auth.UserID) {
	exerciseCompleteList, err := h.exerciseCompleteUseCase.GetItemsByUserID(r.Context(), id)
	if err != nil {
		appError := _internalError.NewAppError(http.StatusNotFound, "Аналитика не найдена", err)
		h.responder.Error(w, appError)

		return
	}

	h.responder.Success(w, http.StatusOK, map[string]any{
		"analytics": exerciseCompleteList,
	})
}

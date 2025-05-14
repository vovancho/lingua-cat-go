package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	"github.com/vovancho/lingua-cat-go/pkg/auth"
	_internalError "github.com/vovancho/lingua-cat-go/pkg/error"
	"github.com/vovancho/lingua-cat-go/pkg/request"
	"github.com/vovancho/lingua-cat-go/pkg/response"
)

type ExerciseStoreRequest struct {
	Lang       domain.ExerciseLang `json:"lang" validate:"required,valid_exercise_lang"`
	TaskAmount uint16              `json:"task_amount" validate:"min=1,max=100"`
}

type ExerciseData struct {
	Exercise domain.Exercise `json:"exercise"`
}

type exerciseHandler struct {
	responder       response.Responder
	exerciseUseCase domain.ExerciseUseCase
	validator       *validator.Validate
	auth            *auth.AuthService
}

func NewExerciseHandler(router *http.ServeMux, responder response.Responder, exerciseUseCase domain.ExerciseUseCase, validator *validator.Validate, auth *auth.AuthService) {
	handler := &exerciseHandler{
		responder:       responder,
		exerciseUseCase: exerciseUseCase,
		validator:       validator,
		auth:            auth,
	}

	router.HandleFunc("GET /v1/exercise/{id}", request.WithID(handler.GetByID))
	router.HandleFunc("POST /v1/exercise", handler.Store)
}

// GetByID godoc
// @Summary Получить упражнение по ID
// @Description Получает упражнение по указанному идентификатору.
// @Security BearerAuth
// @Tags Exercise
// @Param id path uint64 true "ID упражнения"
// @Success 200 {object} response.APIResponse{data=ExerciseData} "Упражнение найдено"
// @Failure 404 {object} response.APIResponse "Упражнение не найдено"
// @Router /v1/exercise/{id} [get]
func (h *exerciseHandler) GetByID(w http.ResponseWriter, r *http.Request, id uint64) {
	exercise, err := h.exerciseUseCase.GetByID(r.Context(), domain.ExerciseID(id))
	if err != nil {
		appError := _internalError.NewAppError(http.StatusNotFound, "Упражнение не найдено", err)
		h.responder.Error(w, appError)

		return
	}

	h.responder.Success(w, http.StatusOK, map[string]any{
		"exercise": exercise,
	})
}

// Store godoc
// @Summary Создать новое упражнение
// @Description Создает новое упражнение с предоставленными данными.
// @Security BearerAuth
// @Tags Exercise
// @Param exercise body ExerciseStoreRequest true "Данные упражнения"
// @Success 201 {object} response.APIResponse{data=ExerciseData} "Упражнение создано"
// @Failure 400 {object} response.APIResponse "Некорректный запрос"
// @Failure 401 {object} response.APIResponse "Неавторизованный доступ"
// @Router /v1/exercise [post]
func (h *exerciseHandler) Store(w http.ResponseWriter, r *http.Request) {
	userID, err := h.auth.GetUserID(r.Context())
	if err != nil {
		err = _internalError.NewAppError(http.StatusUnauthorized, "Не удалось получить userID", err)
		h.responder.Error(w, err)

		return
	}

	var requestBody ExerciseStoreRequest
	if err := h.validateRequest(r, &requestBody); err != nil {
		h.responder.Error(w, err)

		return
	}

	exercise := domain.Exercise{
		UserID:     *userID,
		Lang:       requestBody.Lang,
		TaskAmount: requestBody.TaskAmount,
	}

	if err = h.exerciseUseCase.Store(r.Context(), &exercise); err != nil {
		err = _internalError.NewAppError(http.StatusBadRequest, "Ошибка сохранения упражнения", err)
		h.responder.Error(w, err)

		return
	}

	h.responder.Success(w, http.StatusCreated, map[string]any{
		"exercise": exercise,
	})
}

func (h *exerciseHandler) validateRequest(r *http.Request, req any) error {
	if err := request.FromJSON(r, req); err != nil {
		return _internalError.NewAppError(http.StatusBadRequest, "Некорректный синтаксис JSON", _internalError.InvalidDecodeJsonError)
	}

	if err := h.validator.Struct(req); err != nil {
		return err
	}

	return nil
}

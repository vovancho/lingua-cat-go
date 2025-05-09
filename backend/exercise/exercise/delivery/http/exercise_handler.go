package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	"github.com/vovancho/lingua-cat-go/exercise/internal/auth"
	_internalError "github.com/vovancho/lingua-cat-go/exercise/internal/error"
	"github.com/vovancho/lingua-cat-go/exercise/internal/request"
	"github.com/vovancho/lingua-cat-go/exercise/internal/response"
	"net/http"
)

type ExerciseStoreRequest struct {
	Lang       domain.ExerciseLang `json:"lang" validate:"required,valid_exercise_lang"`
	TaskAmount uint16              `json:"task_amount" validate:"required,min=1,max=100"`
}

type ExerciseData struct {
	Exercise domain.Exercise `json:"exercise"`
}

type ExerciseHandler struct {
	EUseCase domain.ExerciseUseCase
	validate *validator.Validate
	auth     *auth.AuthService
}

func NewExerciseHandler(router *http.ServeMux, v *validator.Validate, auth *auth.AuthService, e domain.ExerciseUseCase) {
	handler := &ExerciseHandler{
		EUseCase: e,
		validate: v,
		auth:     auth,
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
func (d *ExerciseHandler) GetByID(w http.ResponseWriter, r *http.Request, id uint64) {
	exercise, err := d.EUseCase.GetByID(r.Context(), domain.ExerciseID(id))
	if err != nil {
		appError := _internalError.NewAppError(http.StatusNotFound, "Упражнение не найдено", err)
		response.Error(appError, r)
		return
	}

	response.JSON(w, http.StatusOK, response.APIResponse{
		Data: map[string]any{
			"exercise": exercise,
		},
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
func (e *ExerciseHandler) Store(w http.ResponseWriter, r *http.Request) {
	userID, err := e.auth.GetUserID(r.Context())
	if err != nil {
		err = _internalError.NewAppError(http.StatusUnauthorized, "Не удалось получить userID", err)
		response.Error(err, r)
		return
	}

	var requestBody ExerciseStoreRequest
	if err := e.validateRequest(r, &requestBody); err != nil {
		response.Error(err, r)
		return
	}

	exercise := domain.Exercise{
		UserID:     *userID,
		Lang:       requestBody.Lang,
		TaskAmount: requestBody.TaskAmount,
	}

	if err = e.EUseCase.Store(r.Context(), &exercise); err != nil {
		err = _internalError.NewAppError(http.StatusBadRequest, "Ошибка сохранения упражнения", err)
		response.Error(err, r)
		return
	}

	response.JSON(w, http.StatusCreated, response.APIResponse{
		Data: map[string]any{
			"exercise": exercise,
		},
	})
}

func (e *ExerciseHandler) validateRequest(r *http.Request, req any) error {
	if err := request.FromJSON(r, req); err != nil {
		return _internalError.NewAppError(http.StatusBadRequest, "Некорректный синтаксис JSON", _internalError.InvalidDecodeJsonError)
	}

	if err := e.validate.Struct(req); err != nil {
		return err
	}

	return nil
}

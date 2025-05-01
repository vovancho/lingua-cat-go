package http

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	"github.com/vovancho/lingua-cat-go/exercise/internal/auth"
	_internalError "github.com/vovancho/lingua-cat-go/exercise/internal/error"
	"github.com/vovancho/lingua-cat-go/exercise/internal/request"
	"github.com/vovancho/lingua-cat-go/exercise/internal/response"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	TUseCase domain.TaskUseCase
	EUseCase domain.ExerciseUseCase
	validate *validator.Validate
	auth     *auth.AuthService
}

func NewTaskHandler(
	router *http.ServeMux,
	v *validator.Validate,
	auth *auth.AuthService,
	t domain.TaskUseCase,
	e domain.ExerciseUseCase,
) {
	handler := &TaskHandler{
		TUseCase: t,
		EUseCase: e,
		validate: v,
		auth:     auth,
	}

	router.HandleFunc("GET /v1/exercise/{id}/task/{taskId}", request.WithID(withTaskID(handler.GetByID)))
	router.HandleFunc("POST /v1/exercise/{id}/task", request.WithID(handler.Create))
	router.HandleFunc("POST /v1/exercise/{id}/task/{taskId}/word-selected", request.WithID(withTaskID(handler.SelectWord)))
}

func (t *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request, exerciseID domain.ExerciseID, taskID *domain.TaskID) {
	task, err := t.TUseCase.GetByID(r.Context(), *taskID)
	if err != nil || task.Exercise.ID != exerciseID {
		appError := _internalError.NewAppError(http.StatusNotFound, "Задача не найдена", err)
		response.Error(appError, r)
		return
	}

	response.JSON(w, http.StatusOK, response.APIResponse{
		Data: map[string]any{
			"task": task,
		},
	})
}

func (t *TaskHandler) Create(w http.ResponseWriter, r *http.Request, id uint64) {
	exerciseID := domain.ExerciseID(id)
	userID, err := t.auth.GetUserID(r.Context())
	if err != nil {
		err = _internalError.NewAppError(http.StatusUnauthorized, "Не удалось получить userID", err)
		response.Error(err, r)
		return
	}

	if ok, err := t.EUseCase.IsExerciseOwner(r.Context(), exerciseID, *userID); !ok {
		err = _internalError.NewAppError(http.StatusForbidden, "Только автор упражнения может получить задачу", err)
		response.Error(err, r)
		return
	}

	task, err := t.TUseCase.Create(r.Context(), exerciseID)
	if err != nil {
		err = _internalError.NewAppError(http.StatusBadRequest, "Ошибка генерации задачи", err)
		response.Error(err, r)
		return
	}

	response.JSON(w, http.StatusCreated, response.APIResponse{
		Data: map[string]any{
			"task": task,
		},
	})
}

func (d *TaskHandler) SelectWord(w http.ResponseWriter, r *http.Request, exerciseID domain.ExerciseID, taskID *domain.TaskID) {
	response.JSON(w, http.StatusCreated, response.APIResponse{
		Data: map[string]any{
			"task": "task",
		},
	})
}

type HandlerFuncWithExerciseIDAndTaskID func(w http.ResponseWriter, r *http.Request, exerciseID domain.ExerciseID, taskID *domain.TaskID)

func withTaskID(h HandlerFuncWithExerciseIDAndTaskID) request.HandlerFuncWithID {
	return func(w http.ResponseWriter, r *http.Request, id uint64) {
		exerciseID := domain.ExerciseID(id)

		idString := r.PathValue("taskId")
		var taskID *domain.TaskID
		if idString != "" {
			id, err := strconv.ParseUint(idString, 10, 64)
			if err != nil {
				messageErr := fmt.Sprintf("Парсинг \"%s\": некорректный формат ID задачи", idString)
				appErr := _internalError.NewAppError(http.StatusBadRequest, messageErr, _internalError.InvalidPathParamError)
				response.Error(appErr, r)
				return
			}
			t := domain.TaskID(id)
			taskID = &t
		}

		h(w, r, exerciseID, taskID)
	}
}

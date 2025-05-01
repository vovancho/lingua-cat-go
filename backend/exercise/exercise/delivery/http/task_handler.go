package http

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	_internalError "github.com/vovancho/lingua-cat-go/exercise/internal/error"
	"github.com/vovancho/lingua-cat-go/exercise/internal/request"
	"github.com/vovancho/lingua-cat-go/exercise/internal/response"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	TUseCase domain.TaskUseCase
	validate *validator.Validate
}

func NewTaskHandler(router *http.ServeMux, v *validator.Validate, t domain.TaskUseCase) {
	handler := &TaskHandler{
		TUseCase: t,
		validate: v,
	}

	router.HandleFunc("GET /v1/exercise/{id}/task/{taskId}", request.WithID(withTaskID(handler.GetByID)))
	router.HandleFunc("POST /v1/exercise/{id}/task", request.WithID(handler.Store))
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

func (d *TaskHandler) Store(w http.ResponseWriter, r *http.Request, id uint64) {
	exerciseID := domain.ExerciseID(id)

	response.JSON(w, http.StatusCreated, response.APIResponse{
		Data: map[string]any{
			"task": exerciseID,
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

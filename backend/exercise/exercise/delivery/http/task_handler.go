package http

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	_internalError "github.com/vovancho/lingua-cat-go/exercise/internal/error"
	"github.com/vovancho/lingua-cat-go/exercise/internal/request"
	"github.com/vovancho/lingua-cat-go/exercise/internal/response"
	"github.com/vovancho/lingua-cat-go/pkg/auth"
	"net/http"
	"strconv"
)

type TaskWordSelectRequest struct {
	WordSelect uint64 `json:"word_select" validate:"required,gt=0"`
}

type TaskData struct {
	Task domain.Task `json:"task"`
}

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

// GetByID godoc
// @Summary Получить задачу по ID
// @Description Получает задачу по идентификатору упражнения и идентификатору задачи.
// @Security BearerAuth
// @Tags Task
// @Param id path uint64 true "ID упражнения"
// @Param taskId path uint64 true "ID задачи"
// @Success 200 {object} response.APIResponse{data=TaskData} "Задача найдена"
// @Failure 400 {object} response.APIResponse "Некорректный формат ID задачи"
// @Failure 404 {object} response.APIResponse "Задача не найдена или не принадлежит упражнению"
// @Router /v1/exercise/{id}/task/{taskId} [get]
func (t *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request, exerciseID domain.ExerciseID, taskID *domain.TaskID) {
	task, err := t.TUseCase.GetByID(r.Context(), *taskID)
	if err != nil {
		appError := _internalError.NewAppError(http.StatusNotFound, "Задача не найдена", err)
		response.Error(appError, r)
		return
	}
	if task.Exercise.ID != exerciseID {
		appError := _internalError.NewAppError(http.StatusNotFound, "Задача не найдена", domain.TaskNotFoundError)
		response.Error(appError, r)
		return
	}

	response.JSON(w, http.StatusOK, response.APIResponse{
		Data: map[string]any{
			"task": task,
		},
	})
}

// Create godoc
// @Summary Создать новую задачу
// @Description Создает новую задачу для указанного упражнения. Требуется, чтобы пользователь был автором упражнения
// @Security BearerAuth
// @Tags Task
// @Param id path uint64 true "ID упражнения"
// @Success 201 {object} response.APIResponse{data=TaskData} "Задача создана"
// @Failure 400 {object} response.APIResponse "Ошибка генерации задачи"
// @Failure 401 {object} response.APIResponse "Неавторизованный доступ"
// @Failure 403 {object} response.APIResponse "Только автор упражнения может создать задачу"
// @Router /v1/exercise/{id}/task [post]
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

// SelectWord godoc
// @Summary Выбрать слово задачи
// @Description Выбирает слово задачи в указанном упражнении. Требуется, чтобы пользователь был автором упражнения и задачи
// @Security BearerAuth
// @Tags Task
// @Param id path uint64 true "ID упражнения"
// @Param taskId path uint64 true "ID задачи"
// @Param word_select body TaskWordSelectRequest true "ID выбранного слова"
// @Success 201 {object} response.APIResponse{data=TaskData} "Слово выбрано, задача обновлена"
// @Failure 400 {object} response.APIResponse "Некорректный запрос"
// @Failure 401 {object} response.APIResponse "Неавторизованный доступ"
// @Failure 403 {object} response.APIResponse "Только автор упражнения или задачи может выбрать слово"
// @Router /v1/exercise/{id}/task/{taskId}/word-selected [post]
func (t *TaskHandler) SelectWord(w http.ResponseWriter, r *http.Request, exerciseID domain.ExerciseID, taskID *domain.TaskID) {
	var requestBody TaskWordSelectRequest
	if err := t.validateRequest(r, &requestBody); err != nil {
		response.Error(err, r)
		return
	}
	dictionaryID := domain.DictionaryID(requestBody.WordSelect)

	userID, err := t.auth.GetUserID(r.Context())
	if err != nil {
		err = _internalError.NewAppError(http.StatusUnauthorized, "Не удалось получить userID", err)
		response.Error(err, r)
		return
	}

	if ok, err := t.EUseCase.IsExerciseOwner(r.Context(), exerciseID, *userID); !ok {
		err = _internalError.NewAppError(http.StatusForbidden, "Только автор упражнения может выбрать слово", err)
		response.Error(err, r)
		return
	}

	if ok, err := t.TUseCase.IsTaskOwnerExercise(r.Context(), exerciseID, *taskID); !ok {
		err = _internalError.NewAppError(http.StatusForbidden, "Только автор задания может выбрать слово", err)
		response.Error(err, r)
		return
	}

	task, err := t.TUseCase.SelectWord(r.Context(), exerciseID, *taskID, dictionaryID)
	if err != nil {
		err = _internalError.NewAppError(http.StatusBadRequest, "Ошибка выбора слова", err)
		response.Error(err, r)
		return
	}

	response.JSON(w, http.StatusCreated, response.APIResponse{
		Data: map[string]any{
			"task": task,
		},
	})
}

func (e *TaskHandler) validateRequest(r *http.Request, req any) error {
	if err := request.FromJSON(r, req); err != nil {
		return _internalError.NewAppError(http.StatusBadRequest, "Некорректный синтаксис JSON", _internalError.InvalidDecodeJsonError)
	}

	if err := e.validate.Struct(req); err != nil {
		return err
	}

	return nil
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

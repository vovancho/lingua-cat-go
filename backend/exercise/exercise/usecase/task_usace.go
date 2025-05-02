package usecase

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	"math/rand"
	"time"
)

func NewTaskUseCase(
	eUseCase domain.ExerciseUseCase,
	dUseCase domain.DictionaryUseCase,
	tr domain.TaskRepository,
	v *validator.Validate,
	timeout Timeout,
) domain.TaskUseCase {
	return &taskUseCase{
		eUseCase:       eUseCase,
		dUseCase:       dUseCase,
		taskRepo:       tr,
		validate:       v,
		contextTimeout: time.Duration(timeout),
	}
}

type taskUseCase struct {
	eUseCase       domain.ExerciseUseCase
	dUseCase       domain.DictionaryUseCase
	taskRepo       domain.TaskRepository
	validate       *validator.Validate
	contextTimeout time.Duration
}

func (t taskUseCase) GetByID(ctx context.Context, id domain.TaskID) (*domain.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (t taskUseCase) Create(ctx context.Context, exerciseID domain.ExerciseID) (*domain.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	// получить сущность упражнения
	exercise, err := t.eUseCase.GetByID(ctx, exerciseID)
	if err != nil {
		return nil, err
	}

	// проверить возможность добавления новой задачи
	isCreateTaskAllowed := exercise.SelectedCounter < exercise.TaskAmount && (exercise.ProcessedCounter == 0 || exercise.SelectedCounter == exercise.ProcessedCounter)
	if !isCreateTaskAllowed {
		return nil, domain.NewTaskNotAllowedError
	}

	// получить словари в dictionaryService по языку упражнения
	lang := domain.DictionaryLang(exercise.Lang)
	limit := uint8(4)
	dictionaries, err := t.dUseCase.GetRandomDictionaries(ctx, lang, limit)
	if err != nil {
		return nil, err
	}

	if len(dictionaries) == 0 {
		return nil, domain.DictionariesNotFoundError
	}

	// проверить, что словари правильного языка
	for _, dict := range dictionaries {
		if dict.Lang != lang {
			return nil, domain.DictionaryLangIncorrectError
		}
	}

	// получение случайного словаря - как корректного в задаче
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := r.Intn(len(dictionaries))
	randomDictionary := dictionaries[randomIndex]

	// собрать готовую сущность задачи со словарями и упражнением
	task := domain.Task{
		Words:           dictionaries,
		WordIDCorrected: randomDictionary.ID,
		Exercise:        *exercise,
	}

	// сохранить сущность задачи + увеличить счетчик обработанных задач
	if err := t.taskRepo.Store(ctx, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (t taskUseCase) IsTaskOwnerExercise(ctx context.Context, exerciseID domain.ExerciseID, taskID domain.TaskID) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	ok, err := t.taskRepo.IsTaskOwnerExercise(ctx, exerciseID, taskID)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (t taskUseCase) SelectWord(ctx context.Context, exerciseID domain.ExerciseID, taskId domain.TaskID, dictId domain.DictionaryID) (*domain.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	// получить задачу со словарями
	// проверить, что задача принадлежит упражнению
	// проверить, что dictId есть в Words
	// принять выбранное слово SetWordSelected
	// проверить, что exercise.taskAmount == exercise.processedCounter (упражнение завершено)
	// если упражнение завершено, получить потраченное время (created_at - updated_at) и отправить сообщение в kafka

	//TODO implement me
	panic("implement me")
}

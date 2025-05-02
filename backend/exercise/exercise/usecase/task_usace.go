package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	"math/rand"
	"slices"
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
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	task, wordIDs, wordCorrectID, wordSelectedID, err := t.taskRepo.GetByID(ctx, id)
	if err != nil {
		// Если это таймаут — не затираем ошибку
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		fmt.Println(err)

		return nil, domain.TaskNotFoundError
	}

	// получить словари в dictionaryService по ID
	dictionaries, err := t.dUseCase.GetDictionariesByIds(ctx, wordIDs)
	if err != nil {
		return nil, err
	}

	// если какие-то словари были удалены, то задача невалидна
	if len(dictionaries) != len(wordIDs) {
		return nil, domain.TaskNotFoundError
	}

	var (
		wordCorrect  *domain.Dictionary
		wordSelected *domain.Dictionary
	)

	for _, dict := range dictionaries {
		if dict.ID == wordCorrectID {
			wordCorrect = &dict
		}
		if dict.ID == wordSelectedID {
			// так как wordSelected может быть nil, проверим, что ID задан
			wordSelected = &dict
		}
	}

	if wordCorrect == nil {
		return nil, fmt.Errorf("не найден словарь с ID wordCorrect = %d", wordCorrectID)
	}

	task.Words = dictionaries
	task.WordCorrect = *wordCorrect
	task.WordSelected = wordSelected

	return task, nil
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
		if exercise.SelectedCounter == exercise.TaskAmount {
			return nil, domain.ExerciseCompletedError
		}

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
		Words:       dictionaries,
		WordCorrect: randomDictionary,
		Exercise:    *exercise,
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
	task, err := t.GetByID(ctx, taskId)
	if err != nil {
		return nil, domain.TaskNotFoundError
	}

	// проверить, что задача принадлежит упражнению
	if exerciseID != task.Exercise.ID {
		return nil, domain.TaskNotFoundError
	}

	// проверить, что dictId есть в Words
	if !slices.ContainsFunc(task.Words, func(d domain.Dictionary) bool {
		return d.ID == dictId
	}) {
		return nil, domain.DictionaryNotFoundError
	}

	// принять выбранное слово SetWordSelected
	if err := t.taskRepo.SetWordSelected(ctx, task, dictId); err != nil {
		return nil, err
	}

	// проверить, что exercise.taskAmount == exercise.processedCounter (упражнение завершено)
	if task.Exercise.TaskAmount == task.Exercise.ProcessedCounter {
		// если упражнение завершено, получить потраченное время (created_at - updated_at) и отправить сообщение в kafka
		spentTime := task.Exercise.UpdatedAt.Sub(task.Exercise.CreatedAt)

		fmt.Println("Spent time: ", spentTime.String())
	}

	return task, nil
}

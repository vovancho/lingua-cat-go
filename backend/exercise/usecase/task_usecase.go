package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	_watermillSql "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	"github.com/vovancho/lingua-cat-go/pkg/eventpublisher"
)

const dictionaryLimit uint8 = 4

type ExerciseCompletedPublisherInterface interface {
	eventpublisher.PublisherInterface
}

func NewTaskUseCase(
	exerciseUseCase domain.ExerciseUseCase,
	dictionaryUseCase domain.DictionaryUseCase,
	repo domain.TaskRepository,
	validator *validator.Validate,
	publisher ExerciseCompletedPublisherInterface,
) domain.TaskUseCase {
	return &taskUseCase{
		exerciseUseCase:   exerciseUseCase,
		dictionaryUseCase: dictionaryUseCase,
		taskRepo:          repo,
		validate:          validator,
		publisher:         publisher,
	}
}

type taskUseCase struct {
	exerciseUseCase   domain.ExerciseUseCase
	dictionaryUseCase domain.DictionaryUseCase
	taskRepo          domain.TaskRepository
	validate          *validator.Validate
	publisher         ExerciseCompletedPublisherInterface
}

func (uc taskUseCase) GetByID(ctx context.Context, id domain.TaskID) (*domain.Task, error) {
	taskWithDetails, err := uc.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.TaskNotFoundError
	}

	// получить словари в dictionaryService по ID
	dictionaries, err := uc.dictionaryUseCase.GetDictionariesByIds(ctx, taskWithDetails.WordIDs)
	if err != nil {
		return nil, err
	}

	// если какие-то словари были удалены, то задача невалидна
	if len(dictionaries) != len(taskWithDetails.WordIDs) {
		return nil, domain.TaskNotFoundError
	}

	var wordCorrect *domain.Dictionary
	var wordSelected *domain.Dictionary

	for _, dict := range dictionaries {
		if dict.ID == taskWithDetails.WordCorrectID {
			wordCorrect = &dict
		}
		// так как WordSelectedID может быть nil, проверим, что он задан
		if dict.ID == taskWithDetails.WordSelectedID {
			wordSelected = &dict
		}
	}

	if wordCorrect == nil {
		return nil, fmt.Errorf("не найден словарь с ID wordCorrect = %d", taskWithDetails.WordCorrectID)
	}

	taskWithDetails.Task.Words = dictionaries
	taskWithDetails.Task.WordCorrect = *wordCorrect
	taskWithDetails.Task.WordSelected = wordSelected

	return taskWithDetails.Task, nil
}

func (uc taskUseCase) Create(ctx context.Context, exerciseID domain.ExerciseID) (*domain.Task, error) {
	// получить сущность упражнения
	exercise, err := uc.exerciseUseCase.GetByID(ctx, exerciseID)
	if err != nil {
		return nil, err
	}

	// проверить возможность добавления новой задачи
	if !uc.isCreateTaskAllowed(exercise) {
		if exercise.SelectedCounter == exercise.TaskAmount {
			return nil, domain.ExerciseCompletedError
		}

		return nil, domain.NewTaskNotAllowedError
	}

	// получить словари в dictionaryService по языку упражнения
	lang := domain.DictionaryLang(exercise.Lang)
	dictionaries, err := uc.dictionaryUseCase.GetRandomDictionaries(ctx, lang, dictionaryLimit)
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
	randomDictionary := uc.getRandomDictionary(dictionaries)

	// собрать готовую сущность задачи со словарями и упражнением
	task := domain.Task{
		Words:       dictionaries,
		WordCorrect: randomDictionary,
		Exercise:    *exercise,
	}

	// сохранить сущность задачи + увеличить счетчик обработанных задач
	if err := uc.taskRepo.Store(ctx, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (uc taskUseCase) IsTaskOwnerExercise(ctx context.Context, exerciseID domain.ExerciseID, taskID domain.TaskID) (bool, error) {
	isOwner, err := uc.taskRepo.IsTaskOwnerExercise(ctx, exerciseID, taskID)
	if err != nil {
		return false, domain.ExerciseIsNotOwnerOfTask
	}

	return isOwner, nil
}

func (uc taskUseCase) SelectWord(ctx context.Context, exerciseID domain.ExerciseID, taskId domain.TaskID, dictId domain.DictionaryID) (*domain.Task, error) {
	// получить задачу со словарями
	task, err := uc.GetByID(ctx, taskId)
	if err != nil {
		return nil, domain.TaskNotFoundError
	}

	// проверить, что задача принадлежит упражнению
	if exerciseID != task.Exercise.ID {
		return nil, domain.ExerciseIsNotOwnerOfTask
	}

	// проверить, что dictId есть в Words
	for _, dict := range task.Words {
		if dict.ID == dictId {
			break
		}

		return nil, domain.DictionaryNotFoundError
	}

	// публикуем событие об выполненном упражнении
	afterWordSetCallback := func(tx _watermillSql.ContextExecutor, t domain.Task) error {
		if t.Exercise.TaskAmount == t.Exercise.SelectedCounter {
			if err := uc.publishExerciseCompletedEvent(ctx, tx, t); err != nil {
				return err
			}
		}

		return nil
	}

	// принять выбранное слово SetWordSelected
	if err := uc.taskRepo.SetWordSelected(ctx, task, dictId, afterWordSetCallback); err != nil {
		return nil, err
	}

	return task, nil
}

func (uc taskUseCase) isCreateTaskAllowed(exercise *domain.Exercise) bool {
	return exercise.SelectedCounter < exercise.TaskAmount &&
		(exercise.ProcessedCounter == 0 || exercise.SelectedCounter == exercise.ProcessedCounter)
}

func (uc taskUseCase) getRandomDictionary(dictionaries []domain.Dictionary) domain.Dictionary {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := r.Intn(len(dictionaries))

	return dictionaries[randomIndex]
}

func (uc taskUseCase) publishExerciseCompletedEvent(ctx context.Context, tx _watermillSql.ContextExecutor, task domain.Task) error {
	spentTime := task.Exercise.UpdatedAt.Sub(task.Exercise.CreatedAt)

	exerciseCompletedEvent := domain.ExerciseCompletedEvent{
		UserID:              task.Exercise.UserID,
		ExerciseID:          task.Exercise.ID,
		ExerciseLang:        task.Exercise.Lang,
		SpentTime:           int64(spentTime.Seconds()),
		WordsCount:          task.Exercise.TaskAmount,
		WordsCorrectedCount: task.Exercise.CorrectedCounter,
	}

	if err := uc.publisher.Publish(ctx, tx, exerciseCompletedEvent); err != nil {
		return fmt.Errorf("publish message: %w", err)
	}

	return nil
}

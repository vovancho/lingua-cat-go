package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	_watermillSql "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type ExerciseCompletedTopic string

func NewTaskUseCase(
	eUseCase domain.ExerciseUseCase,
	dUseCase domain.DictionaryUseCase,
	tr domain.TaskRepository,
	v *validator.Validate,
	timeout Timeout,
	ect ExerciseCompletedTopic,
) domain.TaskUseCase {
	return &taskUseCase{
		eUseCase:               eUseCase,
		dUseCase:               dUseCase,
		taskRepo:               tr,
		validate:               v,
		contextTimeout:         time.Duration(timeout),
		exerciseCompletedTopic: string(ect),
	}
}

type taskUseCase struct {
	eUseCase               domain.ExerciseUseCase
	dUseCase               domain.DictionaryUseCase
	taskRepo               domain.TaskRepository
	validate               *validator.Validate
	contextTimeout         time.Duration
	exerciseCompletedTopic string
}

func (t taskUseCase) GetByID(ctx context.Context, id domain.TaskID) (*domain.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	taskWithDetails, err := t.taskRepo.GetByID(ctx, id)
	if err != nil {
		// Если это таймаут — не затираем ошибку
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		fmt.Println(err)

		return nil, domain.TaskNotFoundError
	}

	// получить словари в dictionaryService по ID
	dictionaries, err := t.dUseCase.GetDictionariesByIds(ctx, taskWithDetails.WordIDs)
	if err != nil {
		return nil, err
	}

	// если какие-то словари были удалены, то задача невалидна
	if len(dictionaries) != len(taskWithDetails.WordIDs) {
		return nil, domain.TaskNotFoundError
	}

	var (
		wordCorrect  *domain.Dictionary
		wordSelected *domain.Dictionary
	)

	for _, dict := range dictionaries {
		if dict.ID == taskWithDetails.WordCorrectID {
			wordCorrect = &dict
		}
		if dict.ID == taskWithDetails.WordSelectedID {
			// так как wordSelected может быть nil, проверим, что ID задан
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

func (tuc taskUseCase) SelectWord(ctx context.Context, exerciseID domain.ExerciseID, taskId domain.TaskID, dictId domain.DictionaryID) (*domain.Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, tuc.contextTimeout)
	defer cancel()

	// получить задачу со словарями
	task, err := tuc.GetByID(ctx, taskId)
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

	afterWordSetCallback := func(ce _watermillSql.ContextExecutor, t domain.Task) error {
		if t.Exercise.TaskAmount == t.Exercise.SelectedCounter {
			spentTime := t.Exercise.UpdatedAt.Sub(t.Exercise.CreatedAt)

			exerciseCompletedEvent := domain.ExerciseCompletedEvent{
				UserID:              t.Exercise.UserID,
				ExerciseID:          t.Exercise.ID,
				ExerciseLang:        t.Exercise.Lang,
				SpentTime:           int64(spentTime.Seconds()),
				WordsCount:          t.Exercise.TaskAmount,
				WordsCorrectedCount: t.Exercise.CorrectedCounter,
			}

			payload, err := json.Marshal(exerciseCompletedEvent)
			if err != nil {
				return fmt.Errorf("marshal payload: %w", err)
			}

			msg := message.NewMessage(watermill.NewUUID(), payload)

			// Inject tracing context
			propagator := otel.GetTextMapPropagator()
			carrier := propagation.MapCarrier{}
			propagator.Inject(ctx, carrier)
			for k, v := range carrier {
				msg.Metadata[k] = v
			}

			txPublisher, err := sql.NewPublisher(
				ce,
				sql.PublisherConfig{
					SchemaAdapter: sql.DefaultPostgreSQLSchema{},
				},
				nil, // logger should be injected if available
			)
			if err != nil {
				return fmt.Errorf("create tx publisher: %w", err)
			}

			err = txPublisher.Publish(tuc.exerciseCompletedTopic, msg)
			if err != nil {
				return fmt.Errorf("publish message: %w", err)
			}
		}
		return nil
	}

	// принять выбранное слово SetWordSelected
	if err := tuc.taskRepo.SetWordSelected(ctx, task, dictId, afterWordSetCallback); err != nil {
		return nil, err
	}

	return task, nil
}

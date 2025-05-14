package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	"github.com/vovancho/lingua-cat-go/pkg/auth"
)

// MockExerciseUseCase implements domain.ExerciseUseCase for testing
type MockExerciseUseCase struct {
	mock.Mock
}

func (m *MockExerciseUseCase) GetByID(ctx context.Context, id domain.ExerciseID) (*domain.Exercise, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Exercise), args.Error(1)
}

func (m *MockExerciseUseCase) IsExerciseOwner(ctx context.Context, exerciseID domain.ExerciseID, userID auth.UserID) (bool, error) {
	args := m.Called(ctx, exerciseID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockExerciseUseCase) Store(ctx context.Context, exercise *domain.Exercise) error {
	args := m.Called(ctx, exercise)
	return args.Error(0)
}

// MockDictionaryUseCase implements domain.DictionaryUseCase for testing
type MockDictionaryUseCase struct {
	mock.Mock
}

func (m *MockDictionaryUseCase) GetRandomDictionaries(ctx context.Context, lang domain.DictionaryLang, limit uint8) ([]domain.Dictionary, error) {
	args := m.Called(ctx, lang, limit)
	return args.Get(0).([]domain.Dictionary), args.Error(1)
}

func (m *MockDictionaryUseCase) GetDictionariesByIds(ctx context.Context, dictIds []domain.DictionaryID) ([]domain.Dictionary, error) {
	args := m.Called(ctx, dictIds)
	return args.Get(0).([]domain.Dictionary), args.Error(1)
}

// MockTaskRepository implements domain.TaskRepository for testing
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) GetByID(ctx context.Context, id domain.TaskID) (*domain.TaskWithDetails, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.TaskWithDetails), args.Error(1)
}

func (m *MockTaskRepository) IsTaskOwnerExercise(ctx context.Context, exerciseID domain.ExerciseID, taskID domain.TaskID) (bool, error) {
	args := m.Called(ctx, exerciseID, taskID)
	return args.Bool(0), args.Error(1)
}

func (m *MockTaskRepository) Store(ctx context.Context, task *domain.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) SetWordSelected(ctx context.Context, task *domain.Task, dictId domain.DictionaryID, afterWordSetCallback func(sql.ContextExecutor, domain.Task) error) error {
	args := m.Called(ctx, task, dictId, afterWordSetCallback)
	return args.Error(0)
}

// MockPublisher implements eventpublisher.PublisherInterface for testing
type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) Publish(ctx context.Context, tx sql.ContextExecutor, message interface{}) error {
	args := m.Called(ctx, tx, message)
	return args.Error(0)
}

func TestTaskUseCase_GetByID(t *testing.T) {
	ctx := context.Background()
	taskID := domain.TaskID(1)
	userID, _ := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	exerciseID := domain.ExerciseID(1)
	dictIDs := []domain.DictionaryID{1, 2, 3, 4}
	wordCorrectID := domain.DictionaryID(1)
	dictionaries := []domain.Dictionary{
		{ID: 1, Lang: domain.EnDictionary, Name: "word1", Type: domain.SimpleDictionary},
		{ID: 2, Lang: domain.EnDictionary, Name: "word2", Type: domain.SimpleDictionary},
		{ID: 3, Lang: domain.EnDictionary, Name: "word3", Type: domain.SimpleDictionary},
		{ID: 4, Lang: domain.EnDictionary, Name: "word4", Type: domain.SimpleDictionary},
	}
	taskWithDetails := &domain.TaskWithDetails{
		Task: &domain.Task{
			ID: taskID,
			Exercise: domain.Exercise{
				ID:               exerciseID,
				UserID:           auth.UserID(userID),
				Lang:             domain.EnExercise,
				TaskAmount:       10,
				ProcessedCounter: 0,
				SelectedCounter:  0,
				CorrectedCounter: 0,
			},
		},
		WordIDs:       dictIDs,
		WordCorrectID: wordCorrectID,
	}

	t.Run("Success", func(t *testing.T) {
		uc, _, mockDictUC, mockTaskRepo, _ := newTaskUseCaseWithMocks()

		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(taskWithDetails, nil).Once()
		mockDictUC.On("GetDictionariesByIds", mock.Anything, dictIDs).Return(dictionaries, nil).Once()

		result, err := uc.GetByID(ctx, taskID)
		assert.NoError(t, err)
		assert.Equal(t, dictionaries, result.Words)
		assert.Equal(t, dictionaries[0], result.WordCorrect)
		assert.Nil(t, result.WordSelected)
		mockTaskRepo.AssertExpectations(t)
		mockDictUC.AssertExpectations(t)
	})

	t.Run("TaskNotFound", func(t *testing.T) {
		uc, _, mockDictUC, mockTaskRepo, _ := newTaskUseCaseWithMocks()

		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return((*domain.TaskWithDetails)(nil), errors.New("not found")).Once()

		result, err := uc.GetByID(ctx, taskID)
		assert.ErrorIs(t, err, domain.TaskNotFoundError)
		assert.Nil(t, result)
		mockTaskRepo.AssertExpectations(t)
		mockDictUC.AssertNotCalled(t, "GetDictionariesByIds")
	})

	t.Run("DictionariesNotFound", func(t *testing.T) {
		uc, _, mockDictUC, mockTaskRepo, _ := newTaskUseCaseWithMocks()

		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(taskWithDetails, nil).Once()
		mockDictUC.On("GetDictionariesByIds", mock.Anything, dictIDs).Return([]domain.Dictionary{}, domain.DictionariesNotFoundError).Once()

		result, err := uc.GetByID(ctx, taskID)
		assert.ErrorIs(t, err, domain.DictionariesNotFoundError)
		assert.Nil(t, result)
		mockTaskRepo.AssertExpectations(t)
		mockDictUC.AssertExpectations(t)
	})

	t.Run("WordCorrectNotFound", func(t *testing.T) {
		uc, _, mockDictUC, mockTaskRepo, _ := newTaskUseCaseWithMocks()

		invalidTask := &domain.TaskWithDetails{
			Task:          taskWithDetails.Task,
			WordIDs:       dictIDs,
			WordCorrectID: domain.DictionaryID(999), // ID не существует
		}
		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(invalidTask, nil).Once()
		mockDictUC.On("GetDictionariesByIds", mock.Anything, dictIDs).Return(dictionaries, nil).Once()

		result, err := uc.GetByID(ctx, taskID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "не найден словарь с ID wordCorrect")
		assert.Nil(t, result)
		mockTaskRepo.AssertExpectations(t)
		mockDictUC.AssertExpectations(t)
	})

	t.Run("PartialDictionaries", func(t *testing.T) {
		uc, _, mockDictUC, mockTaskRepo, _ := newTaskUseCaseWithMocks()

		partialDictionaries := []domain.Dictionary{
			{ID: 1, Lang: domain.EnDictionary, Name: "word1", Type: domain.SimpleDictionary},
		}
		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(taskWithDetails, nil).Once()
		mockDictUC.On("GetDictionariesByIds", mock.Anything, dictIDs).Return(partialDictionaries, nil).Once()

		result, err := uc.GetByID(ctx, taskID)
		assert.ErrorIs(t, err, domain.TaskNotFoundError)
		assert.Nil(t, result)
		mockTaskRepo.AssertExpectations(t)
		mockDictUC.AssertExpectations(t)
	})

	t.Run("WordSelected", func(t *testing.T) {
		uc, _, mockDictUC, mockTaskRepo, _ := newTaskUseCaseWithMocks()

		taskWithSelected := &domain.TaskWithDetails{
			Task:           taskWithDetails.Task,
			WordIDs:        dictIDs,
			WordCorrectID:  wordCorrectID,
			WordSelectedID: domain.DictionaryID(2),
		}
		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(taskWithSelected, nil).Once()
		mockDictUC.On("GetDictionariesByIds", mock.Anything, dictIDs).Return(dictionaries, nil).Once()

		result, err := uc.GetByID(ctx, taskID)
		assert.NoError(t, err)
		assert.Equal(t, dictionaries, result.Words)
		assert.Equal(t, dictionaries[0], result.WordCorrect)
		assert.Equal(t, &dictionaries[1], result.WordSelected)
		mockTaskRepo.AssertExpectations(t)
		mockDictUC.AssertExpectations(t)
	})
}

func TestTaskUseCase_Create(t *testing.T) {
	ctx := context.Background()
	exerciseID := domain.ExerciseID(1)
	userID, _ := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	dictionaries := []domain.Dictionary{
		{ID: 1, Lang: domain.EnDictionary, Name: "word1", Type: domain.SimpleDictionary},
		{ID: 2, Lang: domain.EnDictionary, Name: "word2", Type: domain.SimpleDictionary},
		{ID: 3, Lang: domain.EnDictionary, Name: "word3", Type: domain.SimpleDictionary},
		{ID: 4, Lang: domain.EnDictionary, Name: "word4", Type: domain.SimpleDictionary},
	}
	exercise := &domain.Exercise{
		ID:               exerciseID,
		UserID:           auth.UserID(userID),
		Lang:             domain.EnExercise,
		TaskAmount:       10,
		ProcessedCounter: 0,
		SelectedCounter:  0,
		CorrectedCounter: 0,
	}

	t.Run("Success", func(t *testing.T) {
		uc, mockExerciseUC, mockDictUC, mockTaskRepo, _ := newTaskUseCaseWithMocks()

		mockExerciseUC.On("GetByID", mock.Anything, exerciseID).Return(exercise, nil).Once()
		mockDictUC.On("GetRandomDictionaries", mock.Anything, domain.DictionaryLang(exercise.Lang), uint8(4)).Return(dictionaries, nil).Once()
		mockTaskRepo.On("Store", mock.Anything, mock.AnythingOfType("*domain.Task")).Return(nil).Run(func(args mock.Arguments) {
			task := args.Get(1).(*domain.Task)
			assert.Equal(t, dictionaries, task.Words)
			assert.Contains(t, dictionaries, task.WordCorrect)
			assert.Equal(t, exercise, &task.Exercise)
		}).Once()

		result, err := uc.Create(ctx, exerciseID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, dictionaries, result.Words)
		assert.Contains(t, dictionaries, result.WordCorrect)
		assert.Nil(t, result.WordSelected)
		mockExerciseUC.AssertExpectations(t)
		mockDictUC.AssertExpectations(t)
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("ExerciseNotFound", func(t *testing.T) {
		uc, mockExerciseUC, mockDictUC, mockTaskRepo, _ := newTaskUseCaseWithMocks()

		mockExerciseUC.On("GetByID", mock.Anything, exerciseID).Return((*domain.Exercise)(nil), domain.ExerciseNotFoundError).Once()

		result, err := uc.Create(ctx, exerciseID)
		assert.ErrorIs(t, err, domain.ExerciseNotFoundError)
		assert.Nil(t, result)
		mockExerciseUC.AssertExpectations(t)
		mockDictUC.AssertNotCalled(t, "GetRandomDictionaries")
		mockTaskRepo.AssertNotCalled(t, "Store")
	})

	t.Run("ExerciseCompleted", func(t *testing.T) {
		uc, mockExerciseUC, mockDictUC, mockTaskRepo, _ := newTaskUseCaseWithMocks()

		completedExercise := &domain.Exercise{
			ID:               exerciseID,
			UserID:           auth.UserID(userID),
			Lang:             domain.EnExercise,
			TaskAmount:       10,
			ProcessedCounter: 10,
			SelectedCounter:  10,
			CorrectedCounter: 8,
		}
		mockExerciseUC.On("GetByID", mock.Anything, exerciseID).Return(completedExercise, nil).Once()

		result, err := uc.Create(ctx, exerciseID)
		assert.ErrorIs(t, err, domain.ExerciseCompletedError)
		assert.Nil(t, result)
		mockExerciseUC.AssertExpectations(t)
		mockDictUC.AssertNotCalled(t, "GetRandomDictionaries")
		mockTaskRepo.AssertNotCalled(t, "Store")
	})

	t.Run("NewTaskNotAllowed", func(t *testing.T) {
		uc, mockExerciseUC, mockDictUC, mockTaskRepo, _ := newTaskUseCaseWithMocks()

		invalidExercise := &domain.Exercise{
			ID:               exerciseID,
			UserID:           auth.UserID(userID),
			Lang:             domain.EnExercise,
			TaskAmount:       10,
			ProcessedCounter: 5,
			SelectedCounter:  4,
			CorrectedCounter: 4,
		}
		mockExerciseUC.On("GetByID", mock.Anything, exerciseID).Return(invalidExercise, nil).Once()

		result, err := uc.Create(ctx, exerciseID)
		assert.ErrorIs(t, err, domain.NewTaskNotAllowedError)
		assert.Nil(t, result)
		mockExerciseUC.AssertExpectations(t)
		mockDictUC.AssertNotCalled(t, "GetRandomDictionaries")
		mockTaskRepo.AssertNotCalled(t, "Store")
	})

	t.Run("DictionariesNotFound", func(t *testing.T) {
		uc, mockExerciseUC, mockDictUC, mockTaskRepo, _ := newTaskUseCaseWithMocks()

		mockExerciseUC.On("GetByID", mock.Anything, exerciseID).Return(exercise, nil).Once()
		mockDictUC.On("GetRandomDictionaries", mock.Anything, domain.DictionaryLang(exercise.Lang), uint8(4)).Return([]domain.Dictionary{}, domain.DictionariesNotFoundError).Once()

		result, err := uc.Create(ctx, exerciseID)
		assert.ErrorIs(t, err, domain.DictionariesNotFoundError)
		assert.Nil(t, result)
		mockExerciseUC.AssertExpectations(t)
		mockDictUC.AssertExpectations(t)
		mockTaskRepo.AssertNotCalled(t, "Store")
	})

	t.Run("DictionaryLangIncorrect", func(t *testing.T) {
		uc, mockExerciseUC, mockDictUC, mockTaskRepo, _ := newTaskUseCaseWithMocks()

		invalidDictionaries := []domain.Dictionary{
			{ID: 1, Lang: domain.RuDictionary, Name: "word1", Type: domain.SimpleDictionary},
		}
		mockExerciseUC.On("GetByID", mock.Anything, exerciseID).Return(exercise, nil).Once()
		mockDictUC.On("GetRandomDictionaries", mock.Anything, domain.DictionaryLang(exercise.Lang), uint8(4)).Return(invalidDictionaries, nil).Once()

		result, err := uc.Create(ctx, exerciseID)
		assert.ErrorIs(t, err, domain.DictionaryLangIncorrectError)
		assert.Nil(t, result)
		mockExerciseUC.AssertExpectations(t)
		mockDictUC.AssertExpectations(t)
		mockTaskRepo.AssertNotCalled(t, "Store")
	})

	t.Run("RepositoryError", func(t *testing.T) {
		uc, mockExerciseUC, mockDictUC, mockTaskRepo, _ := newTaskUseCaseWithMocks()

		mockExerciseUC.On("GetByID", mock.Anything, exerciseID).Return(exercise, nil).Once()
		mockDictUC.On("GetRandomDictionaries", mock.Anything, domain.DictionaryLang(exercise.Lang), uint8(4)).Return(dictionaries, nil).Once()
		mockTaskRepo.On("Store", mock.Anything, mock.AnythingOfType("*domain.Task")).Return(errors.New("repo error")).Once()

		result, err := uc.Create(ctx, exerciseID)
		assert.Error(t, err)
		assert.EqualError(t, err, "repo error")
		assert.Nil(t, result)
		mockExerciseUC.AssertExpectations(t)
		mockDictUC.AssertExpectations(t)
		mockTaskRepo.AssertExpectations(t)
	})
}

func TestTaskUseCase_IsTaskOwnerExercise(t *testing.T) {
	ctx := context.Background()
	exerciseID := domain.ExerciseID(1)
	taskID := domain.TaskID(1)

	t.Run("Success_Owner", func(t *testing.T) {
		uc, _, _, mockTaskRepo, _ := newTaskUseCaseWithMocks()

		mockTaskRepo.On("IsTaskOwnerExercise", mock.Anything, exerciseID, taskID).Return(true, nil).Once()

		isOwner, err := uc.IsTaskOwnerExercise(ctx, exerciseID, taskID)
		assert.NoError(t, err)
		assert.True(t, isOwner)
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("Success_NotOwner", func(t *testing.T) {
		uc, _, _, mockTaskRepo, _ := newTaskUseCaseWithMocks()

		mockTaskRepo.On("IsTaskOwnerExercise", mock.Anything, exerciseID, taskID).Return(false, nil).Once()

		isOwner, err := uc.IsTaskOwnerExercise(ctx, exerciseID, taskID)
		assert.NoError(t, err)
		assert.False(t, isOwner)
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		uc, _, _, mockTaskRepo, _ := newTaskUseCaseWithMocks()

		mockTaskRepo.On("IsTaskOwnerExercise", mock.Anything, exerciseID, taskID).Return(false, errors.New("repo error")).Once()

		isOwner, err := uc.IsTaskOwnerExercise(ctx, exerciseID, taskID)
		assert.ErrorIs(t, err, domain.ExerciseIsNotOwnerOfTask)
		assert.False(t, isOwner)
		mockTaskRepo.AssertExpectations(t)
	})
}

func TestTaskUseCase_SelectWord(t *testing.T) {
	ctx := context.Background()
	exerciseID := domain.ExerciseID(1)
	taskID := domain.TaskID(1)
	userID, _ := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	dictID := domain.DictionaryID(1)
	dictionaries := []domain.Dictionary{
		{ID: 1, Lang: domain.EnDictionary, Name: "word1", Type: domain.SimpleDictionary},
		{ID: 2, Lang: domain.EnDictionary, Name: "word2", Type: domain.SimpleDictionary},
		{ID: 3, Lang: domain.EnDictionary, Name: "word3", Type: domain.SimpleDictionary},
		{ID: 4, Lang: domain.EnDictionary, Name: "word4", Type: domain.SimpleDictionary},
	}
	task := &domain.Task{
		ID:          taskID,
		Words:       dictionaries,
		WordCorrect: dictionaries[0],
		Exercise: domain.Exercise{
			ID:               exerciseID,
			UserID:           auth.UserID(userID),
			Lang:             domain.EnExercise,
			TaskAmount:       10,
			ProcessedCounter: 10,
			SelectedCounter:  10,
			CorrectedCounter: 8,
		},
	}
	taskWithDetails := &domain.TaskWithDetails{
		Task:          task,
		WordIDs:       []domain.DictionaryID{1, 2, 3, 4},
		WordCorrectID: dictID,
	}

	t.Run("Success", func(t *testing.T) {
		uc, _, mockDictUC, mockTaskRepo, mockPublisher := newTaskUseCaseWithMocks()

		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(taskWithDetails, nil).Once()
		mockDictUC.On("GetDictionariesByIds", mock.Anything, taskWithDetails.WordIDs).Return(dictionaries, nil).Once()
		mockTaskRepo.On("SetWordSelected", mock.Anything, task, dictID, mock.AnythingOfType("func(sql.ContextExecutor, domain.Task) error")).Return(nil).Run(func(args mock.Arguments) {
			callback := args.Get(3).(func(sql.ContextExecutor, domain.Task) error)
			err := callback(nil, *task)
			assert.NoError(t, err)
		}).Once()
		mockPublisher.On("Publish", mock.Anything, mock.Anything, mock.AnythingOfType("domain.ExerciseCompletedEvent")).Return(nil).Once()

		result, err := uc.SelectWord(ctx, exerciseID, taskID, dictID)
		assert.NoError(t, err)
		assert.Equal(t, task, result)
		mockTaskRepo.AssertExpectations(t)
		mockDictUC.AssertExpectations(t)
		mockPublisher.AssertExpectations(t)
	})

	t.Run("TaskNotFound", func(t *testing.T) {
		uc, _, mockDictUC, mockTaskRepo, _ := newTaskUseCaseWithMocks()

		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return((*domain.TaskWithDetails)(nil), errors.New("not found")).Once()

		result, err := uc.SelectWord(ctx, exerciseID, taskID, dictID)
		assert.ErrorIs(t, err, domain.TaskNotFoundError)
		assert.Nil(t, result)
		mockTaskRepo.AssertExpectations(t)
		mockDictUC.AssertNotCalled(t, "GetDictionariesByIds")
		mockTaskRepo.AssertNotCalled(t, "SetWordSelected")
	})

	t.Run("WordAlreadySelected", func(t *testing.T) {
		uc, _, mockDictUC, mockTaskRepo, _ := newTaskUseCaseWithMocks()

		taskWithSelected := &domain.TaskWithDetails{
			Task: &domain.Task{
				ID:           taskID,
				Words:        dictionaries,
				WordCorrect:  dictionaries[0],
				WordSelected: &dictionaries[1],
				Exercise:     task.Exercise,
			},
			WordIDs:        []domain.DictionaryID{1, 2, 3, 4},
			WordCorrectID:  dictID,
			WordSelectedID: domain.DictionaryID(2),
		}
		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(taskWithSelected, nil).Once()
		mockDictUC.On("GetDictionariesByIds", mock.Anything, taskWithSelected.WordIDs).Return(dictionaries, nil).Once()

		result, err := uc.SelectWord(ctx, exerciseID, taskID, dictID)
		assert.ErrorIs(t, err, domain.TaskWordAlreadySelectedError)
		assert.Nil(t, result)
		mockTaskRepo.AssertExpectations(t)
		mockDictUC.AssertExpectations(t)
		mockTaskRepo.AssertNotCalled(t, "SetWordSelected")
	})

	t.Run("ExerciseNotOwner", func(t *testing.T) {
		uc, _, mockDictUC, mockTaskRepo, _ := newTaskUseCaseWithMocks()

		wrongExerciseID := domain.ExerciseID(2)
		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(taskWithDetails, nil).Once()
		mockDictUC.On("GetDictionariesByIds", mock.Anything, taskWithDetails.WordIDs).Return(dictionaries, nil).Once()

		result, err := uc.SelectWord(ctx, wrongExerciseID, taskID, dictID)
		assert.ErrorIs(t, err, domain.ExerciseIsNotOwnerOfTask)
		assert.Nil(t, result)
		mockTaskRepo.AssertExpectations(t)
		mockDictUC.AssertExpectations(t)
		mockTaskRepo.AssertNotCalled(t, "SetWordSelected")
	})

	t.Run("DictionaryNotFound", func(t *testing.T) {
		uc, _, mockDictUC, mockTaskRepo, _ := newTaskUseCaseWithMocks()

		invalidDictID := domain.DictionaryID(999)
		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(taskWithDetails, nil).Once()
		mockDictUC.On("GetDictionariesByIds", mock.Anything, taskWithDetails.WordIDs).Return(dictionaries, nil).Once()

		result, err := uc.SelectWord(ctx, exerciseID, taskID, invalidDictID)
		assert.ErrorIs(t, err, domain.DictionaryNotFoundError)
		assert.Nil(t, result)
		mockTaskRepo.AssertExpectations(t)
		mockDictUC.AssertExpectations(t)
		mockTaskRepo.AssertNotCalled(t, "SetWordSelected")
	})

	t.Run("RepositoryError", func(t *testing.T) {
		uc, _, mockDictUC, mockTaskRepo, _ := newTaskUseCaseWithMocks()

		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(taskWithDetails, nil).Once()
		mockDictUC.On("GetDictionariesByIds", mock.Anything, taskWithDetails.WordIDs).Return(dictionaries, nil).Once()
		mockTaskRepo.On("SetWordSelected", mock.Anything, task, dictID, mock.AnythingOfType("func(sql.ContextExecutor, domain.Task) error")).Return(errors.New("repo error")).Once()

		result, err := uc.SelectWord(ctx, exerciseID, taskID, dictID)
		assert.Error(t, err)
		assert.EqualError(t, err, "repo error")
		assert.Nil(t, result)
		mockTaskRepo.AssertExpectations(t)
		mockDictUC.AssertExpectations(t)
	})

	t.Run("PublishError", func(t *testing.T) {
		uc, _, mockDictUC, mockTaskRepo, mockPublisher := newTaskUseCaseWithMocks()

		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(taskWithDetails, nil).Once()
		mockDictUC.On("GetDictionariesByIds", mock.Anything, taskWithDetails.WordIDs).Return(dictionaries, nil).Once()
		mockTaskRepo.On("SetWordSelected", mock.Anything, task, dictID, mock.AnythingOfType("func(sql.ContextExecutor, domain.Task) error")).Run(func(args mock.Arguments) {
			callback := args.Get(3).(func(sql.ContextExecutor, domain.Task) error)
			err := callback(nil, *task)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "publish message")
		}).Return(errors.New("execute afterWordSetCallback: publish error")).Once()
		mockPublisher.On("Publish", mock.Anything, mock.Anything, mock.AnythingOfType("domain.ExerciseCompletedEvent")).Return(errors.New("publish error")).Once()

		result, err := uc.SelectWord(ctx, exerciseID, taskID, dictID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "execute afterWordSetCallback: publish error")
		assert.Nil(t, result)
		mockTaskRepo.AssertExpectations(t)
		mockDictUC.AssertExpectations(t)
		mockPublisher.AssertExpectations(t)
	})
}

func newTaskUseCaseWithMocks() (domain.TaskUseCase, *MockExerciseUseCase, *MockDictionaryUseCase, *MockTaskRepository, *MockPublisher) {
	exerciseUC := new(MockExerciseUseCase)
	dictUC := new(MockDictionaryUseCase)
	taskRepo := new(MockTaskRepository)
	publisher := new(MockPublisher)
	uc := NewTaskUseCase(exerciseUC, dictUC, taskRepo, publisher)
	return uc, exerciseUC, dictUC, taskRepo, publisher
}

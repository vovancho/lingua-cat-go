package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vovancho/lingua-cat-go/analytics/domain"
	"github.com/vovancho/lingua-cat-go/pkg/auth"
)

// MockExerciseCompleteRepository implements domain.ExerciseCompleteRepository for testing
type MockExerciseCompleteRepository struct {
	mock.Mock
}

func (m *MockExerciseCompleteRepository) GetItemsByUserID(ctx context.Context, userId auth.UserID) ([]domain.ExerciseComplete, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).([]domain.ExerciseComplete), args.Error(1)
}

func (m *MockExerciseCompleteRepository) Store(ctx context.Context, ec *domain.ExerciseComplete) error {
	args := m.Called(ctx, ec)
	return args.Error(0)
}

func TestExerciseCompleteUseCase_GetItemsByUserID(t *testing.T) {
	ctx := context.Background()
	userID, _ := uuid.Parse("123e4567-e89b-12d3-a456-426614174000") // Валидный UUID
	expectedList := []domain.ExerciseComplete{
		{
			UserID:              auth.UserID(userID),
			UserName:            "testuser",
			ExerciseID:          domain.ExerciseID(1),
			ExerciseLang:        domain.EnExercise,
			SpentTime:           1000,
			WordsCount:          10,
			WordsCorrectedCount: 8,
			EventTime:           time.Now(),
		},
		{
			UserID:              auth.UserID(userID),
			UserName:            "testuser",
			ExerciseID:          domain.ExerciseID(2),
			ExerciseLang:        domain.RuExercise,
			SpentTime:           1500,
			WordsCount:          15,
			WordsCorrectedCount: 12,
			EventTime:           time.Now(),
		},
	}

	t.Run("Success", func(t *testing.T) {
		uc, mockRepo := newExerciseCompleteUseCaseWithMock()

		mockRepo.On("GetItemsByUserID", mock.Anything, auth.UserID(userID)).Return(expectedList, nil).Once()

		result, err := uc.GetItemsByUserID(ctx, auth.UserID(userID))
		assert.NoError(t, err)
		assert.Equal(t, expectedList, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		uc, mockRepo := newExerciseCompleteUseCaseWithMock()

		mockRepo.On("GetItemsByUserID", mock.Anything, auth.UserID(userID)).Return([]domain.ExerciseComplete{}, errors.New("not found")).Once()

		result, err := uc.GetItemsByUserID(ctx, auth.UserID(userID))
		assert.ErrorIs(t, err, domain.ExerciseCompleteListNotFoundError)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestExerciseCompleteUseCase_Store(t *testing.T) {
	ctx := context.Background()
	userID, _ := uuid.Parse("123e4567-e89b-12d3-a456-426614174000") // Валидный UUID
	exerciseComplete := &domain.ExerciseComplete{
		UserID:              auth.UserID(userID),
		UserName:            "testuser",
		ExerciseID:          domain.ExerciseID(1),
		ExerciseLang:        domain.EnExercise,
		SpentTime:           1000,
		WordsCount:          10,
		WordsCorrectedCount: 8,
		EventTime:           time.Now(),
	}

	t.Run("Success", func(t *testing.T) {
		uc, mockRepo := newExerciseCompleteUseCaseWithMock()

		mockRepo.On("Store", mock.Anything, exerciseComplete).Return(nil).Once()

		err := uc.Store(ctx, exerciseComplete)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidLang", func(t *testing.T) {
		uc, mockRepo := newExerciseCompleteUseCaseWithMock()

		invalidExerciseComplete := &domain.ExerciseComplete{
			UserID:              auth.UserID(userID),
			UserName:            "testuser",
			ExerciseID:          domain.ExerciseID(1),
			ExerciseLang:        domain.ExerciseLang("invalid"),
			SpentTime:           1000,
			WordsCount:          10,
			WordsCorrectedCount: 8,
			EventTime:           time.Now(),
		}

		err := uc.Store(ctx, invalidExerciseComplete)
		assert.Error(t, err)
		assert.IsType(t, validator.ValidationErrors{}, err)
		validationErrors := err.(validator.ValidationErrors)
		assert.Equal(t, "ExerciseLang", validationErrors[0].Field())
		assert.Equal(t, "valid_exercise_lang", validationErrors[0].Tag())
		mockRepo.AssertNotCalled(t, "Store")
	})

	t.Run("MissingUserID", func(t *testing.T) {
		uc, mockRepo := newExerciseCompleteUseCaseWithMock()

		invalidExerciseComplete := &domain.ExerciseComplete{
			UserName:            "testuser",
			ExerciseID:          domain.ExerciseID(1),
			ExerciseLang:        domain.EnExercise,
			SpentTime:           1000,
			WordsCount:          10,
			WordsCorrectedCount: 8,
			EventTime:           time.Now(),
		}

		err := uc.Store(ctx, invalidExerciseComplete)
		assert.Error(t, err)
		assert.IsType(t, validator.ValidationErrors{}, err)
		validationErrors := err.(validator.ValidationErrors)
		assert.Equal(t, "UserID", validationErrors[0].Field())
		assert.Equal(t, "required", validationErrors[0].Tag())
		mockRepo.AssertNotCalled(t, "Store")
	})

	t.Run("MissingUserName", func(t *testing.T) {
		uc, mockRepo := newExerciseCompleteUseCaseWithMock()

		invalidExerciseComplete := &domain.ExerciseComplete{
			UserID:              auth.UserID(userID),
			ExerciseID:          domain.ExerciseID(1),
			ExerciseLang:        domain.EnExercise,
			SpentTime:           1000,
			WordsCount:          10,
			WordsCorrectedCount: 8,
			EventTime:           time.Now(),
		}

		err := uc.Store(ctx, invalidExerciseComplete)
		assert.Error(t, err)
		assert.IsType(t, validator.ValidationErrors{}, err)
		validationErrors := err.(validator.ValidationErrors)
		assert.Equal(t, "UserName", validationErrors[0].Field())
		assert.Equal(t, "required", validationErrors[0].Tag())
		mockRepo.AssertNotCalled(t, "Store")
	})

	t.Run("MissingExerciseID", func(t *testing.T) {
		uc, mockRepo := newExerciseCompleteUseCaseWithMock()

		invalidExerciseComplete := &domain.ExerciseComplete{
			UserID:              auth.UserID(userID),
			UserName:            "testuser",
			ExerciseLang:        domain.EnExercise,
			SpentTime:           1000,
			WordsCount:          10,
			WordsCorrectedCount: 8,
			EventTime:           time.Now(),
		}

		err := uc.Store(ctx, invalidExerciseComplete)
		assert.Error(t, err)
		assert.IsType(t, validator.ValidationErrors{}, err)
		validationErrors := err.(validator.ValidationErrors)
		assert.Equal(t, "ExerciseID", validationErrors[0].Field())
		assert.Equal(t, "required", validationErrors[0].Tag())
		mockRepo.AssertNotCalled(t, "Store")
	})

	t.Run("MissingExerciseLang", func(t *testing.T) {
		uc, mockRepo := newExerciseCompleteUseCaseWithMock()

		invalidExerciseComplete := &domain.ExerciseComplete{
			UserID:              auth.UserID(userID),
			UserName:            "testuser",
			ExerciseID:          domain.ExerciseID(1),
			SpentTime:           1000,
			WordsCount:          10,
			WordsCorrectedCount: 8,
			EventTime:           time.Now(),
		}

		err := uc.Store(ctx, invalidExerciseComplete)
		assert.Error(t, err)
		assert.IsType(t, validator.ValidationErrors{}, err)
		validationErrors := err.(validator.ValidationErrors)
		assert.Equal(t, "ExerciseLang", validationErrors[0].Field())
		assert.Equal(t, "required", validationErrors[0].Tag())
		mockRepo.AssertNotCalled(t, "Store")
	})

	t.Run("InvalidWordsCount_Zero", func(t *testing.T) {
		uc, mockRepo := newExerciseCompleteUseCaseWithMock()

		invalidExerciseComplete := &domain.ExerciseComplete{
			UserID:              auth.UserID(userID),
			UserName:            "testuser",
			ExerciseID:          domain.ExerciseID(1),
			ExerciseLang:        domain.EnExercise,
			SpentTime:           1000,
			WordsCount:          0,
			WordsCorrectedCount: 8,
			EventTime:           time.Now(),
		}

		err := uc.Store(ctx, invalidExerciseComplete)
		assert.Error(t, err)
		assert.IsType(t, validator.ValidationErrors{}, err)
		validationErrors := err.(validator.ValidationErrors)
		assert.Equal(t, "WordsCount", validationErrors[0].Field())
		assert.Equal(t, "min", validationErrors[0].Tag())
		mockRepo.AssertNotCalled(t, "Store")
	})

	t.Run("MissingWordsCorrectedCount", func(t *testing.T) {
		uc, mockRepo := newExerciseCompleteUseCaseWithMock()

		invalidExerciseComplete := &domain.ExerciseComplete{
			UserID:              auth.UserID(userID),
			UserName:            "testuser",
			ExerciseID:          domain.ExerciseID(1),
			ExerciseLang:        domain.EnExercise,
			SpentTime:           1000,
			WordsCount:          10,
			WordsCorrectedCount: 0,
			EventTime:           time.Now(),
		}

		err := uc.Store(ctx, invalidExerciseComplete)
		assert.Error(t, err)
		assert.IsType(t, validator.ValidationErrors{}, err)
		validationErrors := err.(validator.ValidationErrors)
		assert.Equal(t, "WordsCorrectedCount", validationErrors[0].Field())
		assert.Equal(t, "required", validationErrors[0].Tag())
		mockRepo.AssertNotCalled(t, "Store")
	})

	t.Run("RepositoryError", func(t *testing.T) {
		uc, mockRepo := newExerciseCompleteUseCaseWithMock()

		mockRepo.On("Store", mock.Anything, exerciseComplete).Return(errors.New("repo error")).Once()

		err := uc.Store(ctx, exerciseComplete)
		assert.Error(t, err)
		assert.EqualError(t, err, "repo error")
		mockRepo.AssertExpectations(t)
	})
}

func newValidator() *validator.Validate {
	v := validator.New()
	if err := domain.RegisterAll(v, nil); err != nil {
		panic(err)
	}
	return v
}

func newExerciseCompleteUseCaseWithMock() (domain.ExerciseCompleteUseCase, *MockExerciseCompleteRepository) {
	repo := new(MockExerciseCompleteRepository)
	uc := NewExerciseCompleteUseCase(repo, newValidator())
	return uc, repo
}

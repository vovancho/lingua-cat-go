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
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	"github.com/vovancho/lingua-cat-go/pkg/auth"
)

// MockExerciseRepository implements domain.ExerciseRepository for testing
type MockExerciseRepository struct {
	mock.Mock
}

func (m *MockExerciseRepository) GetByID(ctx context.Context, id domain.ExerciseID) (*domain.Exercise, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Exercise), args.Error(1)
}

func (m *MockExerciseRepository) IsExerciseOwner(ctx context.Context, exerciseID domain.ExerciseID, userID auth.UserID) (bool, error) {
	args := m.Called(ctx, exerciseID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockExerciseRepository) Store(ctx context.Context, exercise *domain.Exercise) error {
	args := m.Called(ctx, exercise)
	return args.Error(0)
}

func TestExerciseUseCase_GetByID(t *testing.T) {
	ctx := context.Background()
	exerciseID := domain.ExerciseID(1)
	userID, _ := uuid.Parse("123e4567-e89b-12d3-a456-426614174000") // Валидный UUID
	expectedExercise := &domain.Exercise{
		ID:               exerciseID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		UserID:           auth.UserID(userID),
		Lang:             domain.EnExercise,
		TaskAmount:       10,
		ProcessedCounter: 5,
		SelectedCounter:  4,
		CorrectedCounter: 3,
	}

	t.Run("Success", func(t *testing.T) {
		uc, mockRepo := newExerciseUseCaseWithMock()

		mockRepo.On("GetByID", mock.Anything, exerciseID).Return(expectedExercise, nil).Once()

		result, err := uc.GetByID(ctx, exerciseID)
		assert.NoError(t, err)
		assert.Equal(t, expectedExercise, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		uc, mockRepo := newExerciseUseCaseWithMock()

		mockRepo.On("GetByID", mock.Anything, exerciseID).Return((*domain.Exercise)(nil), errors.New("not found")).Once()

		result, err := uc.GetByID(ctx, exerciseID)
		assert.ErrorIs(t, err, domain.ExerciseNotFoundError)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestExerciseUseCase_IsExerciseOwner(t *testing.T) {
	ctx := context.Background()
	exerciseID := domain.ExerciseID(1)
	userID, _ := uuid.Parse("123e4567-e89b-12d3-a456-426614174000") // Валидный UUID

	t.Run("Success_Owner", func(t *testing.T) {
		uc, mockRepo := newExerciseUseCaseWithMock()

		mockRepo.On("IsExerciseOwner", mock.Anything, exerciseID, auth.UserID(userID)).Return(true, nil).Once()

		isOwner, err := uc.IsExerciseOwner(ctx, exerciseID, auth.UserID(userID))
		assert.NoError(t, err)
		assert.True(t, isOwner)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success_NotOwner", func(t *testing.T) {
		uc, mockRepo := newExerciseUseCaseWithMock()

		mockRepo.On("IsExerciseOwner", mock.Anything, exerciseID, auth.UserID(userID)).Return(false, nil).Once()

		isOwner, err := uc.IsExerciseOwner(ctx, exerciseID, auth.UserID(userID))
		assert.NoError(t, err)
		assert.False(t, isOwner)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		uc, mockRepo := newExerciseUseCaseWithMock()

		mockRepo.On("IsExerciseOwner", mock.Anything, exerciseID, auth.UserID(userID)).Return(false, errors.New("repo error")).Once()

		isOwner, err := uc.IsExerciseOwner(ctx, exerciseID, auth.UserID(userID))
		assert.ErrorIs(t, err, domain.UserIsNotOwnerOfExercise)
		assert.False(t, isOwner)
		mockRepo.AssertExpectations(t)
	})
}

func TestExerciseUseCase_Store(t *testing.T) {
	ctx := context.Background()
	userID, _ := uuid.Parse("123e4567-e89b-12d3-a456-426614174000") // Валидный UUID
	exercise := &domain.Exercise{
		UserID:           auth.UserID(userID),
		Lang:             domain.EnExercise,
		TaskAmount:       10,
		ProcessedCounter: 0,
		SelectedCounter:  0,
		CorrectedCounter: 0,
	}

	t.Run("Success", func(t *testing.T) {
		uc, mockRepo := newExerciseUseCaseWithMock()

		mockRepo.On("Store", mock.Anything, exercise).Return(nil).Once()

		err := uc.Store(ctx, exercise)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidLang", func(t *testing.T) {
		uc, mockRepo := newExerciseUseCaseWithMock()

		invalidExercise := &domain.Exercise{
			UserID:           auth.UserID(userID),
			Lang:             domain.ExerciseLang("invalid"),
			TaskAmount:       10,
			ProcessedCounter: 0,
			SelectedCounter:  0,
			CorrectedCounter: 0,
		}

		err := uc.Store(ctx, invalidExercise)
		assert.Error(t, err)
		assert.IsType(t, validator.ValidationErrors{}, err)
		validationErrors := err.(validator.ValidationErrors)
		assert.Equal(t, "Lang", validationErrors[0].Field())
		assert.Equal(t, "valid_exercise_lang", validationErrors[0].Tag())
		mockRepo.AssertNotCalled(t, "Store")
	})

	t.Run("MissingUserID", func(t *testing.T) {
		uc, mockRepo := newExerciseUseCaseWithMock()

		invalidExercise := &domain.Exercise{
			Lang:             domain.EnExercise,
			TaskAmount:       10,
			ProcessedCounter: 0,
			SelectedCounter:  0,
			CorrectedCounter: 0,
		}

		err := uc.Store(ctx, invalidExercise)
		assert.Error(t, err)
		assert.IsType(t, validator.ValidationErrors{}, err)
		validationErrors := err.(validator.ValidationErrors)
		assert.Equal(t, "UserID", validationErrors[0].Field())
		assert.Equal(t, "required", validationErrors[0].Tag())
		mockRepo.AssertNotCalled(t, "Store")
	})

	t.Run("InvalidTaskAmount_Zero", func(t *testing.T) {
		uc, mockRepo := newExerciseUseCaseWithMock()

		invalidExercise := &domain.Exercise{
			UserID:           auth.UserID(userID),
			Lang:             domain.EnExercise,
			TaskAmount:       0,
			ProcessedCounter: 0,
			SelectedCounter:  0,
			CorrectedCounter: 0,
		}

		err := uc.Store(ctx, invalidExercise)
		assert.Error(t, err)
		assert.IsType(t, validator.ValidationErrors{}, err)
		validationErrors := err.(validator.ValidationErrors)
		assert.Equal(t, "TaskAmount", validationErrors[0].Field())
		assert.Equal(t, "min", validationErrors[0].Tag())
		mockRepo.AssertNotCalled(t, "Store")
	})

	t.Run("InvalidTaskAmount_TooHigh", func(t *testing.T) {
		uc, mockRepo := newExerciseUseCaseWithMock()

		invalidExercise := &domain.Exercise{
			UserID:           auth.UserID(userID),
			Lang:             domain.EnExercise,
			TaskAmount:       101,
			ProcessedCounter: 0,
			SelectedCounter:  0,
			CorrectedCounter: 0,
		}

		err := uc.Store(ctx, invalidExercise)
		assert.Error(t, err)
		assert.IsType(t, validator.ValidationErrors{}, err)
		validationErrors := err.(validator.ValidationErrors)
		assert.Equal(t, "TaskAmount", validationErrors[0].Field())
		assert.Equal(t, "max", validationErrors[0].Tag())
		mockRepo.AssertNotCalled(t, "Store")
	})

	t.Run("RepositoryError", func(t *testing.T) {
		uc, mockRepo := newExerciseUseCaseWithMock()

		mockRepo.On("Store", mock.Anything, exercise).Return(errors.New("repo error")).Once()

		err := uc.Store(ctx, exercise)
		assert.Error(t, err)
		assert.EqualError(t, err, "repo error")
		mockRepo.AssertExpectations(t)
	})
}

func newExerciseValidator() *validator.Validate {
	v := validator.New()
	if err := domain.RegisterAll(v, nil); err != nil {
		panic(err)
	}
	return v
}

func newExerciseUseCaseWithMock() (domain.ExerciseUseCase, *MockExerciseRepository) {
	repo := new(MockExerciseRepository)
	uc := NewExerciseUseCase(repo, newExerciseValidator())
	return uc, repo
}

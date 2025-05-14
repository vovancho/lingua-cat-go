package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vovancho/lingua-cat-go/analytics/domain"
	"github.com/vovancho/lingua-cat-go/pkg/auth"
)

// MockUserRepository implements domain.UserRepository for testing
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(ctx context.Context, userID auth.UserID) (*domain.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestUserUseCase_GetByID(t *testing.T) {
	ctx := context.Background()
	userID, _ := uuid.Parse("123e4567-e89b-12d3-a456-426614174000") // Валидный UUID
	expectedUser := &domain.User{
		ID:       auth.UserID(userID),
		Username: "testuser",
	}

	t.Run("Success", func(t *testing.T) {
		uc, mockRepo := newUserUseCaseWithMock()

		mockRepo.On("GetByID", mock.Anything, auth.UserID(userID)).Return(expectedUser, nil).Once()

		result, err := uc.GetByID(ctx, auth.UserID(userID))
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		uc, mockRepo := newUserUseCaseWithMock()

		mockRepo.On("GetByID", mock.Anything, auth.UserID(userID)).Return((*domain.User)(nil), errors.New("not found")).Once()

		result, err := uc.GetByID(ctx, auth.UserID(userID))
		assert.ErrorIs(t, err, domain.UserNotFoundError)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func newUserUseCaseWithMock() (domain.UserUseCase, *MockUserRepository) {
	repo := new(MockUserRepository)
	uc := NewUserUseCase(repo)
	return uc, repo
}

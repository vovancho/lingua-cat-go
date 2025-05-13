package usecase

import (
	"context"

	"github.com/vovancho/lingua-cat-go/analytics/domain"
	"github.com/vovancho/lingua-cat-go/pkg/auth"
)

func NewUserUseCase(repo domain.UserRepository) domain.UserUseCase {
	return &userUseCase{
		userRepo: repo,
	}
}

type userUseCase struct {
	userRepo domain.UserRepository
}

func (uc userUseCase) GetByID(ctx context.Context, userId auth.UserID) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, domain.UserNotFoundError
	}

	return user, nil
}

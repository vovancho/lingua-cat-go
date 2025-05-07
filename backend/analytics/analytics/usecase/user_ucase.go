package usecase

import (
	"context"
	"errors"
	"github.com/vovancho/lingua-cat-go/analytics/domain"
	"github.com/vovancho/lingua-cat-go/analytics/internal/auth"
	"time"
)

func NewUserUseCase(ecr domain.UserRepository, timeout Timeout) domain.UserUseCase {
	return &userUseCase{
		userRepo:       ecr,
		contextTimeout: time.Duration(timeout),
	}
}

type userUseCase struct {
	userRepo       domain.UserRepository
	contextTimeout time.Duration
}

func (u userUseCase) GetByID(ctx context.Context, userId auth.UserID) (*domain.User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.GetByID(ctx, userId)
	if err != nil {
		// Если это таймаут — не затираем ошибку
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		return nil, err
		//return nil, domain.UserNotFoundError
	}

	return user, nil
}

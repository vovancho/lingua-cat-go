package usecase

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/analytics/domain"
	"github.com/vovancho/lingua-cat-go/analytics/internal/auth"
	"time"
)

type Timeout time.Duration

func NewExerciseCompleteUseCase(ecr domain.ExerciseCompleteRepository, v *validator.Validate, timeout Timeout) domain.ExerciseCompleteUseCase {
	return &exerciseCompleteUseCase{
		exerciseCompleteRepo: ecr,
		validate:             v,
		contextTimeout:       time.Duration(timeout),
	}
}

type exerciseCompleteUseCase struct {
	exerciseCompleteRepo domain.ExerciseCompleteRepository
	validate             *validator.Validate
	contextTimeout       time.Duration
}

func (ecu exerciseCompleteUseCase) GetByUserID(ctx context.Context, userId auth.UserID) ([]domain.ExerciseComplete, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, ecu.contextTimeout)
	defer cancel()

	ecuList, err := ecu.exerciseCompleteRepo.GetByUserID(ctx, userId)
	if err != nil {
		// Если это таймаут — не затираем ошибку
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		return nil, domain.ExerciseCompleteListNotFoundError
	}

	return ecuList, nil
}

func (ecu exerciseCompleteUseCase) Store(ctx context.Context, ec *domain.ExerciseComplete) error {
	//TODO implement me
	panic("implement me")
}

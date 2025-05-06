package usecase

import (
	"context"
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

func (e exerciseCompleteUseCase) GetByUserID(ctx context.Context, userId auth.UserID) ([]domain.ExerciseComplete, error) {
	//TODO implement me
	panic("implement me")
}

func (e exerciseCompleteUseCase) Store(ctx context.Context, ec *domain.ExerciseComplete) error {
	//TODO implement me
	panic("implement me")
}

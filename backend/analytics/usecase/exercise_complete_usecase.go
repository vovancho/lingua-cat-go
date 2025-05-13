package usecase

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/analytics/domain"
	"github.com/vovancho/lingua-cat-go/pkg/auth"
)

func NewExerciseCompleteUseCase(repo domain.ExerciseCompleteRepository, validator *validator.Validate) domain.ExerciseCompleteUseCase {
	return &exerciseCompleteUseCase{
		exerciseCompleteRepo: repo,
		validator:            validator,
	}
}

type exerciseCompleteUseCase struct {
	exerciseCompleteRepo domain.ExerciseCompleteRepository
	validator            *validator.Validate
}

func (uc exerciseCompleteUseCase) GetItemsByUserID(ctx context.Context, userId auth.UserID) ([]domain.ExerciseComplete, error) {
	ecuList, err := uc.exerciseCompleteRepo.GetItemsByUserID(ctx, userId)
	if err != nil {
		return nil, domain.ExerciseCompleteListNotFoundError
	}

	return ecuList, nil
}

func (uc exerciseCompleteUseCase) Store(ctx context.Context, ec *domain.ExerciseComplete) error {
	if err := uc.validator.Struct(ec); err != nil {
		return err
	}

	if err := uc.exerciseCompleteRepo.Store(ctx, ec); err != nil {
		return err
	}

	return nil
}

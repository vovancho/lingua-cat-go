package usecase

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	"github.com/vovancho/lingua-cat-go/pkg/auth"
)

func NewExerciseUseCase(repo domain.ExerciseRepository, validator *validator.Validate) domain.ExerciseUseCase {
	return &exerciseUseCase{
		exerciseRepo: repo,
		validate:     validator,
	}
}

type exerciseUseCase struct {
	exerciseRepo   domain.ExerciseRepository
	validate       *validator.Validate
	contextTimeout time.Duration
}

func (e exerciseUseCase) GetByID(ctx context.Context, id domain.ExerciseID) (*domain.Exercise, error) {
	exercise, err := e.exerciseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ExerciseNotFoundError
	}

	return exercise, nil
}

func (e exerciseUseCase) IsExerciseOwner(ctx context.Context, exerciseID domain.ExerciseID, userID auth.UserID) (bool, error) {
	isOwner, err := e.exerciseRepo.IsExerciseOwner(ctx, exerciseID, userID)
	if err != nil {
		return false, domain.UserIsNotOwnerOfExercise
	}

	return isOwner, nil
}

func (e exerciseUseCase) Store(ctx context.Context, exercise *domain.Exercise) error {
	if err := e.validate.Struct(exercise); err != nil {
		return err
	}

	if err := e.exerciseRepo.Store(ctx, exercise); err != nil {
		return err
	}

	return nil
}

package usecase

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	"time"
)

func NewExerciseUseCase(er domain.ExerciseRepository, v *validator.Validate, timeout Timeout) domain.ExerciseUseCase {
	return &exerciseUseCase{
		exerciseRepo:   er,
		validate:       v,
		contextTimeout: time.Duration(timeout),
	}
}

type exerciseUseCase struct {
	exerciseRepo   domain.ExerciseRepository
	validate       *validator.Validate
	contextTimeout time.Duration
}

func (e exerciseUseCase) GetByID(ctx context.Context, id domain.ExerciseID) (*domain.Exercise, error) {
	//TODO implement me
	panic("implement me")
}

func (e exerciseUseCase) Store(ctx context.Context, exercise *domain.Exercise) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	if err := e.validate.Struct(exercise); err != nil {
		fmt.Println(exercise)
		return err
	}

	if err := e.exerciseRepo.Store(ctx, exercise); err != nil {
		return err
	}

	return nil
}

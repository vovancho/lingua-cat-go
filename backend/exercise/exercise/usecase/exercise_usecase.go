package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	"github.com/vovancho/lingua-cat-go/exercise/internal/auth"
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
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	exercise, err := e.exerciseRepo.GetByID(ctx, id)
	if err != nil {
		// Если это таймаут — не затираем ошибку
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}

		return nil, domain.ExerciseNotFoundError
	}

	return exercise, nil
}

func (e exerciseUseCase) IsExerciseOwner(ctx context.Context, exerciseID domain.ExerciseID, userID auth.UserID) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	ok, err := e.exerciseRepo.IsExerciseOwner(ctx, exerciseID, userID)
	if err != nil {
		return false, err
	}

	return ok, nil
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

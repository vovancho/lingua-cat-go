package usecase

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	"time"
)

func NewTaskUseCase(tr domain.TaskRepository, dr domain.DictionaryRepository, v *validator.Validate, timeout Timeout) domain.TaskUseCase {
	return &taskUseCase{
		taskRepo:       tr,
		dictRepo:       dr,
		validate:       v,
		contextTimeout: time.Duration(timeout),
	}
}

type taskUseCase struct {
	taskRepo       domain.TaskRepository
	dictRepo       domain.DictionaryRepository
	validate       *validator.Validate
	contextTimeout time.Duration
}

func (t taskUseCase) GetByID(ctx context.Context, id domain.TaskID) (*domain.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (t taskUseCase) Create(ctx context.Context, exerciseID domain.ExerciseID) (*domain.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (t taskUseCase) SelectWord(ctx context.Context, exerciseID domain.ExerciseID, taskId domain.TaskID, dictId domain.DictionaryID) error {
	//TODO implement me
	panic("implement me")
}

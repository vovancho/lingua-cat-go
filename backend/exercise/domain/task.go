package domain

import "context"

type TaskID uint64

type Task struct {
	ID              TaskID
	Words           []Dictionary
	WordIDCorrected DictionaryID
	WordIDSelected  *DictionaryID
	Exercise        Exercise
}

type TaskUseCase interface {
	GetByID(ctx context.Context, id TaskID) (*Task, error)
	Create(ctx context.Context, exerciseID ExerciseID) (*Task, error)
	SelectWord(ctx context.Context, exerciseID ExerciseID, taskId TaskID, dictId DictionaryID) error
}

type TaskRepository interface {
	GetByID(ctx context.Context, id TaskID) (*Task, error)
	Store(ctx context.Context, task *Task) error
	SetWordSelected(ctx context.Context, taskId TaskID, dictId DictionaryID) error
}

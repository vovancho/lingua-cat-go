package domain

import (
	"context"
)

type TaskID uint64

type Task struct {
	ID           TaskID       `json:"id" db:"id"`
	Words        []Dictionary `json:"words" db:"-"`
	WordCorrect  Dictionary   `json:"word_correct" db:"-"`
	WordSelected *Dictionary  `json:"word_selected" db:"-"`
	Exercise     Exercise     `json:"exercise" db:"exercise"`
}

type TaskUseCase interface {
	GetByID(ctx context.Context, id TaskID) (*Task, error)
	IsTaskOwnerExercise(ctx context.Context, exerciseID ExerciseID, taskID TaskID) (bool, error)
	Create(ctx context.Context, exerciseID ExerciseID) (*Task, error)
	SelectWord(ctx context.Context, exerciseID ExerciseID, taskId TaskID, dictId DictionaryID) (*Task, error)
}

type TaskRepository interface {
	GetByID(ctx context.Context, id TaskID) (*Task, []DictionaryID, DictionaryID, DictionaryID, error)
	IsTaskOwnerExercise(ctx context.Context, exerciseID ExerciseID, taskID TaskID) (bool, error)
	Store(ctx context.Context, task *Task) error
	SetWordSelected(ctx context.Context, task *Task, dictId DictionaryID) error
}

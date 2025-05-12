package domain

import (
	"context"

	"github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
)

type TaskID uint64

type Task struct {
	ID           TaskID       `json:"id" db:"id"`
	Words        []Dictionary `json:"words" db:"-"`
	WordCorrect  Dictionary   `json:"word_correct" db:"-"`
	WordSelected *Dictionary  `json:"word_selected" db:"-"`
	Exercise     Exercise     `json:"exercise" db:"exercise"`
}

type TaskWithDetails struct {
	Task           *Task
	WordIDs        []DictionaryID
	WordCorrectID  DictionaryID
	WordSelectedID DictionaryID
}

type TaskUseCase interface {
	GetByID(ctx context.Context, id TaskID) (*Task, error)
	IsTaskOwnerExercise(ctx context.Context, exerciseID ExerciseID, taskID TaskID) (bool, error)
	Create(ctx context.Context, exerciseID ExerciseID) (*Task, error)
	SelectWord(ctx context.Context, exerciseID ExerciseID, taskId TaskID, dictId DictionaryID) (*Task, error)
}

type TaskRepository interface {
	GetByID(ctx context.Context, id TaskID) (*TaskWithDetails, error)
	IsTaskOwnerExercise(ctx context.Context, exerciseID ExerciseID, taskID TaskID) (bool, error)
	Store(ctx context.Context, task *Task) error
	SetWordSelected(ctx context.Context, task *Task, dictId DictionaryID, afterWordSetCallback func(ce sql.ContextExecutor, t Task) error) error
}

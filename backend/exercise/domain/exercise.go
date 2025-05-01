package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type ExerciseID uint64
type UserID uuid.UUID
type ExerciseLang string

const (
	RuExercise ExerciseLang = "ru"
	EnExercise ExerciseLang = "en"
)

type Exercise struct {
	ID               ExerciseID
	CreatedAt        time.Time
	UpdatedAt        time.Time
	UserID           UserID
	ExerciseLang     ExerciseLang
	TaskAmount       uint16
	ProcessedCounter uint16
	SelectedCounter  uint16
	CorrectedCounter uint16
}

type ExerciseUseCase interface {
	GetByID(ctx context.Context, id ExerciseID) (*Exercise, error)
	Create(ctx context.Context, userId UserID, exerciseLang ExerciseLang, taskAmount uint16) (Exercise, error)
}

type ExerciseRepository interface {
	GetByID(ctx context.Context, id ExerciseID) (*Exercise, error)
	Store(ctx context.Context, exercise *Exercise) error
}

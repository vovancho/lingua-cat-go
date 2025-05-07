package domain

import (
	"context"
	"github.com/vovancho/lingua-cat-go/exercise/internal/auth"
	"time"
)

type ExerciseID uint64
type ExerciseLang string

const (
	RuExercise ExerciseLang = "ru"
	EnExercise ExerciseLang = "en"
)

func (l ExerciseLang) IsValid() bool {
	return l == RuExercise || l == EnExercise
}

type Exercise struct {
	ID               ExerciseID   `json:"id" db:"id"`
	CreatedAt        time.Time    `json:"-" db:"created_at"`
	UpdatedAt        time.Time    `json:"-" db:"updated_at"`
	UserID           auth.UserID  `json:"user_id" db:"user_id" validate:"required"`
	Lang             ExerciseLang `json:"lang" db:"lang" validate:"required,valid_exercise_lang"`
	TaskAmount       uint16       `json:"task_amount" db:"task_amount" validate:"required,min=1,max=100"`
	ProcessedCounter uint16       `json:"processed_counter" db:"processed_counter"`
	SelectedCounter  uint16       `json:"selected_counter" db:"selected_counter"`
	CorrectedCounter uint16       `json:"corrected_counter" db:"corrected_counter"`
}

type ExerciseCompletedEvent struct {
	UserID              auth.UserID  `json:"user_id"`
	ExerciseID          ExerciseID   `json:"exercise_id"`
	ExerciseLang        ExerciseLang `json:"exercise_lang"`
	SpentTime           int64        `json:"spent_time"` // в миллисекундах
	WordsCount          uint16       `json:"words_count"`
	WordsCorrectedCount uint16       `json:"words_corrected_count"`
}

type ExerciseUseCase interface {
	GetByID(ctx context.Context, id ExerciseID) (*Exercise, error)
	IsExerciseOwner(ctx context.Context, exerciseID ExerciseID, userID auth.UserID) (bool, error)
	Store(ctx context.Context, exercise *Exercise) error
}

type ExerciseRepository interface {
	GetByID(ctx context.Context, id ExerciseID) (*Exercise, error)
	IsExerciseOwner(ctx context.Context, exerciseID ExerciseID, userID auth.UserID) (bool, error)
	Store(ctx context.Context, exercise *Exercise) error
}

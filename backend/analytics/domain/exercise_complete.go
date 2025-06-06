package domain

import (
	"context"
	"time"

	"github.com/vovancho/lingua-cat-go/pkg/auth"
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

type ExerciseComplete struct {
	UserID              auth.UserID  `json:"user_id" db:"user_id" validate:"required"`
	UserName            string       `json:"user_name" db:"user_name" validate:"required"`
	ExerciseID          ExerciseID   `json:"exercise_id" db:"exercise_id" validate:"required"`
	ExerciseLang        ExerciseLang `json:"exercise_lang" db:"exercise_lang" validate:"required,valid_exercise_lang"`
	SpentTime           uint64       `json:"spent_time" db:"spent_time"`
	WordsCount          uint16       `json:"words_count" db:"words_count" validate:"min=1"`
	WordsCorrectedCount uint16       `json:"words_corrected_count" db:"words_corrected_count" validate:"required"`
	EventTime           time.Time    `json:"event_time" db:"event_time"`
}

type ExerciseCompleteUseCase interface {
	GetItemsByUserID(ctx context.Context, userId auth.UserID) ([]ExerciseComplete, error)
	Store(ctx context.Context, ec *ExerciseComplete) error
}

type ExerciseCompleteRepository interface {
	GetItemsByUserID(ctx context.Context, userId auth.UserID) ([]ExerciseComplete, error)
	Store(ctx context.Context, ec *ExerciseComplete) error
}

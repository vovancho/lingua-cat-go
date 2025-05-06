package domain

import (
	"context"
	"github.com/vovancho/lingua-cat-go/analytics/internal/auth"
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

type ExerciseComplete struct {
	UserID              auth.UserID  `json:"user_id"`
	UserName            string       `json:"user_name"`
	ExerciseID          ExerciseID   `json:"exercise_id"`
	ExerciseLang        ExerciseLang `json:"exercise_lang"`
	SpentTime           time.Time    `json:"spent_time"`
	WordsCount          uint16       `json:"words_count"`
	WordsCorrectedCount uint16       `json:"words_corrected_count"`
	EventTime           time.Time    `json:"event_time"`
}

type ExerciseCompleteUseCase interface {
	GetByUserID(ctx context.Context, userId auth.UserID) ([]ExerciseComplete, error)
	Store(ctx context.Context, ec *ExerciseComplete) error
}

type ExerciseCompleteRepository interface {
	GetByUserID(ctx context.Context, userId auth.UserID) ([]ExerciseComplete, error)
	Store(ctx context.Context, ec *ExerciseComplete) error
}

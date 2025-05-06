package clickhouse

import (
	"context"
	"github.com/vovancho/lingua-cat-go/analytics/domain"
	"github.com/vovancho/lingua-cat-go/analytics/internal/auth"
	"github.com/vovancho/lingua-cat-go/analytics/internal/db"
)

type clickhouseExerciseCompleteRepository struct {
	Conn db.DB
}

func NewClickhouseExerciseCompleteRepository(conn db.DB) domain.ExerciseCompleteRepository {
	return &clickhouseExerciseCompleteRepository{conn}
}

func (c clickhouseExerciseCompleteRepository) GetByUserID(ctx context.Context, userId auth.UserID) ([]domain.ExerciseComplete, error) {
	//TODO implement me
	panic("implement me")
}

func (c clickhouseExerciseCompleteRepository) Store(ctx context.Context, ec *domain.ExerciseComplete) error {
	//TODO implement me
	panic("implement me")
}

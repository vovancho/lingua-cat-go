package postgres

import (
	"context"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	"github.com/vovancho/lingua-cat-go/exercise/internal/db"
)

type postgresExerciseRepository struct {
	Conn db.DB
}

func NewPostgresExerciseRepository(conn db.DB) domain.ExerciseRepository {
	return &postgresExerciseRepository{conn}
}

func (p postgresExerciseRepository) GetByID(ctx context.Context, id domain.ExerciseID) (*domain.Exercise, error) {
	//TODO implement me
	panic("implement me")
}

func (p postgresExerciseRepository) Store(ctx context.Context, exercise *domain.Exercise) error {
	//TODO implement me
	panic("implement me")
}

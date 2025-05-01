package postgres

import (
	"context"
	"fmt"
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
	const query = `
		INSERT INTO exercise (user_id, lang, task_amount)
		VALUES (:user_id, :lang, :task_amount)
		RETURNING id, created_at, updated_at, user_id, lang, task_amount, processed_counter, selected_counter, corrected_counter`

	rows, err := p.Conn.NamedQueryContext(ctx, query, exercise)
	if err != nil {
		return fmt.Errorf("insert exercise: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(exercise); err != nil {
			return fmt.Errorf("scan exercise: %w", err)
		}
		return nil
	}

	return fmt.Errorf("no rows returned after insert")
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	"github.com/vovancho/lingua-cat-go/pkg/auth"
)

type repository struct {
	conn *sqlx.DB
}

func NewExerciseRepository(conn *sqlx.DB) domain.ExerciseRepository {
	return &repository{conn}
}

func (r repository) GetByID(ctx context.Context, id domain.ExerciseID) (*domain.Exercise, error) {
	const query = `
		SELECT id, created_at, updated_at, user_id, lang, task_amount, processed_counter, selected_counter, corrected_counter
		FROM exercise WHERE id = $1`

	var exercise domain.Exercise
	if err := r.conn.GetContext(ctx, &exercise, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("exercise not found: %w", err)
		}

		return nil, fmt.Errorf("get exercise: %w", err)
	}

	return &exercise, nil
}

func (r repository) IsExerciseOwner(ctx context.Context, exerciseID domain.ExerciseID, userID auth.UserID) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM exercise WHERE id = $1 AND user_id = $2)`

	var exists bool
	if err := r.conn.GetContext(ctx, &exists, query, exerciseID, userID); err != nil {
		return false, fmt.Errorf("check exercise owner: %w", err)
	}

	return exists, nil
}

func (r repository) Store(ctx context.Context, exercise *domain.Exercise) error {
	const query = `
        INSERT INTO exercise (user_id, lang, task_amount)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, updated_at, user_id, lang, task_amount, processed_counter, selected_counter, corrected_counter`

	if err := r.conn.GetContext(ctx, exercise, query, exercise.UserID, exercise.Lang, exercise.TaskAmount); err != nil {
		return fmt.Errorf("insert exercise: %w", err)
	}

	return nil
}

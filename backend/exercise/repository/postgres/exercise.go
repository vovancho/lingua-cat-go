package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	"github.com/vovancho/lingua-cat-go/pkg/auth"
)

type exerciseRepository struct {
	conn *sqlx.DB
}

func NewExerciseRepository(conn *sqlx.DB) domain.ExerciseRepository {
	return &exerciseRepository{conn}
}

func (r exerciseRepository) GetByID(ctx context.Context, id domain.ExerciseID) (*domain.Exercise, error) {
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

func (r exerciseRepository) IsExerciseOwner(ctx context.Context, exerciseID domain.ExerciseID, userID auth.UserID) (bool, error) {
	const query = `SELECT id FROM exercise WHERE id = $1 AND user_id = $2`

	var id domain.DictionaryID
	err := r.conn.GetContext(ctx, &id, query, exerciseID, userID)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r exerciseRepository) Store(ctx context.Context, exercise *domain.Exercise) error {
	const query = `
		INSERT INTO exercise (user_id, lang, task_amount)
		VALUES (:user_id, :lang, :task_amount)
		RETURNING id, created_at, updated_at, user_id, lang, task_amount, processed_counter, selected_counter, corrected_counter`

	rows, err := r.conn.NamedQueryContext(ctx, query, exercise)
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

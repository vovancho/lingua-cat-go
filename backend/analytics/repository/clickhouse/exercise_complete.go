package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/vovancho/lingua-cat-go/analytics/domain"
	"github.com/vovancho/lingua-cat-go/pkg/auth"
)

type exerciseCompleteRepository struct {
	conn *sqlx.DB
}

func NewExerciseCompleteRepository(conn *sqlx.DB) domain.ExerciseCompleteRepository {
	return &exerciseCompleteRepository{conn}
}

func (r exerciseCompleteRepository) GetByUserID(ctx context.Context, userId auth.UserID) ([]domain.ExerciseComplete, error) {
	const query = `SELECT user_id, user_name, exercise_id, exercise_lang, spent_time, words_count, words_corrected_count, event_time
		FROM analytics.exercise_complete
		WHERE user_id = $1
		ORDER BY event_time DESC`

	var ecList []domain.ExerciseComplete

	if err := r.conn.SelectContext(ctx, &ecList, query, userId); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("analytics not found: %w", err)
		}
		return nil, fmt.Errorf("get analytics: %w", err)
	}

	return ecList, nil
}

func (r exerciseCompleteRepository) Store(ctx context.Context, ec *domain.ExerciseComplete) error {
	const query = `
		INSERT INTO analytics.exercise_complete 
		    (user_id, user_name, exercise_id, exercise_lang, spent_time, words_count, words_corrected_count)
		VALUES 
		    (:user_id, :user_name, :exercise_id, :exercise_lang, :spent_time, :words_count, :words_corrected_count)`

	_, err := r.conn.NamedExecContext(ctx, query, ec)
	if err != nil {
		return fmt.Errorf("insert exercise_complete: %w", err)
	}
	return nil
}

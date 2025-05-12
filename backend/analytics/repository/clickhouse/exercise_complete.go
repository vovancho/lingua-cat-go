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

func (r exerciseCompleteRepository) GetItemsByUserID(ctx context.Context, userId auth.UserID) ([]domain.ExerciseComplete, error) {
	const query = `SELECT user_id, user_name, exercise_id, exercise_lang, spent_time, words_count, words_corrected_count, event_time
		FROM analytics.exercise_complete
		WHERE user_id = $1
		ORDER BY event_time DESC`

	var ecList []domain.ExerciseComplete
	if err := r.conn.SelectContext(ctx, &ecList, query, userId); err != nil {
		if err == sql.ErrNoRows {
			return []domain.ExerciseComplete{}, nil
		}

		return nil, fmt.Errorf("get analytics: %w", err)
	}

	return ecList, nil
}

func (r exerciseCompleteRepository) Store(ctx context.Context, ec *domain.ExerciseComplete) error {
	const query = `
		INSERT INTO analytics.exercise_complete (user_id, user_name, exercise_id, exercise_lang, spent_time, words_count, words_corrected_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.conn.ExecContext(ctx, query, ec.UserID, ec.UserName, ec.ExerciseID, ec.ExerciseLang, ec.SpentTime, ec.WordsCount, ec.WordsCorrectedCount)
	if err != nil {
		return fmt.Errorf("insert exercise_complete: %w", err)
	}

	return nil
}

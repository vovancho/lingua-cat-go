package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	"github.com/vovancho/lingua-cat-go/exercise/internal/db"
)

type postgresTaskRepository struct {
	Conn db.DB
}

func NewPostgresTaskRepository(conn db.DB) domain.TaskRepository {
	return &postgresTaskRepository{conn}
}

func (p postgresTaskRepository) GetByID(ctx context.Context, id domain.TaskID) (*domain.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (p postgresTaskRepository) Store(ctx context.Context, task *domain.Task) error {
	return p.withTransaction(ctx, func(tx *sqlx.Tx) error {
		// Вставка задачи
		taskID, err := p.insertTask(ctx, tx, task)
		if err != nil {
			return err
		}
		task.ID = taskID

		// Инкремент счетчика получения задачи
		if err := p.incrementProcessedCounter(ctx, tx, task.Exercise.ID); err != nil {
			return err
		}

		return nil
	})
}

func (p postgresTaskRepository) SetWordSelected(ctx context.Context, taskId domain.TaskID, dictId domain.DictionaryID) error {
	//TODO implement me
	panic("implement me")
}

// withTransaction выполняет callback в контексте транзакции
func (p postgresTaskRepository) withTransaction(ctx context.Context, callback func(*sqlx.Tx) error) error {
	tx, err := p.Conn.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if err = callback(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction: %v; original error: %w", rollbackErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (p postgresTaskRepository) insertTask(ctx context.Context, tx *sqlx.Tx, t *domain.Task) (domain.TaskID, error) {
	var id domain.TaskID
	const query = `INSERT INTO task (exercise_id, words, word_correct) VALUES (:exercise_id, :words, :word_correct) RETURNING id`
	nstmt, err := tx.PrepareNamed(query)
	if err != nil {
		return 0, fmt.Errorf("prepare named query: %w", err)
	}
	err = nstmt.QueryRowxContext(ctx, t).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert task: %w", err)
	}
	return id, nil
}

func (p postgresTaskRepository) incrementProcessedCounter(ctx context.Context, tx *sqlx.Tx, exerciseId domain.ExerciseID) error {
	const query = `UPDATE exercise SET processed_counter = processed_counter + 1 WHERE id = :id AND processed_counter < task_amount`
	res, err := tx.NamedExecContext(ctx, query, map[string]any{
		"id": exerciseId,
	})
	if err != nil {
		return fmt.Errorf("increment processed counter: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if (rowsAffected == 0) || err != nil {
		return fmt.Errorf("increment processed counter not affected")
	}

	return nil
}

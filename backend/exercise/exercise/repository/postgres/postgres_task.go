package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

func (p postgresTaskRepository) IsTaskOwnerExercise(ctx context.Context, exerciseID domain.ExerciseID, taskID domain.TaskID) (bool, error) {
	const query = `SELECT id FROM task WHERE id = $1 AND exercise_id = $2`

	var id domain.TaskID
	err := p.Conn.GetContext(ctx, &id, query, taskID, exerciseID)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
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
		counter, err := p.incrementProcessedCounter(ctx, tx, task.Exercise.ID)
		if err != nil {
			return err
		}
		task.Exercise.ProcessedCounter = counter

		return nil
	})
}

func (p postgresTaskRepository) SetWordSelected(ctx context.Context, task *domain.Task, dictId domain.DictionaryID) error {

	// установить task.word_selected = dictId
	// сделать инкремент exercise.selected_counter
	// сделать инкремент exercise.corrected_counter, если task.word_selected == task.word_correct

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

	wordIDs := make([]int64, len(t.Words))
	for i, w := range t.Words {
		wordIDs[i] = int64(w.ID)
	}

	dto := map[string]any{
		"exercise_id":  t.Exercise.ID,
		"words":        pq.Array(wordIDs),
		"word_correct": t.WordIDCorrected,
	}

	const query = `
		INSERT INTO task (exercise_id, words, word_correct)
		VALUES (:exercise_id, :words, :word_correct)
		RETURNING id`

	nstmt, err := tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("prepare named query: %w", err)
	}
	defer nstmt.Close()

	err = nstmt.QueryRowxContext(ctx, dto).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert task: %w", err)
	}

	return id, nil
}

func (p postgresTaskRepository) incrementProcessedCounter(ctx context.Context, tx *sqlx.Tx, exerciseId domain.ExerciseID) (uint16, error) {
	const query = `
		UPDATE exercise
		SET processed_counter = processed_counter + 1
		WHERE id = :id AND processed_counter < task_amount
		RETURNING processed_counter`

	stmt, err := tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("prepare named statement: %w", err)
	}
	defer stmt.Close()

	var counter uint16
	err = stmt.QueryRowxContext(ctx, map[string]any{"id": exerciseId}).Scan(&counter)
	if err != nil {
		return 0, fmt.Errorf("increment processed counter: %w", err)
	}

	return counter, nil
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	"github.com/vovancho/lingua-cat-go/exercise/internal/db"
	"slices"
)

type postgresTaskRepository struct {
	Conn db.DB
}

func NewPostgresTaskRepository(conn db.DB) domain.TaskRepository {
	return &postgresTaskRepository{conn}
}

func (p postgresTaskRepository) GetByID(ctx context.Context, id domain.TaskID) (
	*domain.Task,
	[]domain.DictionaryID,
	domain.DictionaryID,
	domain.DictionaryID,
	error,
) {
	const query = `
		SELECT t.id,
		       e.id AS "exercise.id",
		       e.created_at AS "exercise.created_at",
		       e.updated_at AS "exercise.updated_at",
		       e.user_id AS "exercise.user_id",
		       e.lang AS "exercise.lang",
		       e.task_amount AS "exercise.task_amount",
		       e.processed_counter AS "exercise.processed_counter",
		       e.selected_counter AS "exercise.selected_counter",
		       e.corrected_counter AS "exercise.corrected_counter",
		       t.word_correct,
		       t.word_selected,
		       t.words AS word_ids
		FROM task t
		INNER JOIN exercise e ON e.id = t.exercise_id
		WHERE t.id = $1`

	var raw struct {
		ID             domain.TaskID        `db:"id"`
		Exercise       domain.Exercise      `db:"exercise"`
		WordIDs        pq.Int64Array        `db:"word_ids"`
		WordCorrectID  domain.DictionaryID  `db:"word_correct"`
		WordSelectedID *domain.DictionaryID `db:"word_selected"`
	}

	if err := p.Conn.GetContext(ctx, &raw, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, 0, 0, fmt.Errorf("task not found: %w", err)
		}
		return nil, nil, 0, 0, fmt.Errorf("get task by id: %w", err)
	}

	wordIDs := make([]domain.DictionaryID, len(raw.WordIDs))
	for i, id := range raw.WordIDs {
		wordIDs[i] = domain.DictionaryID(id)
	}

	// Для удобства возвращаем значение даже если WordSelectedID == nil
	var wordSelectedID domain.DictionaryID
	if raw.WordSelectedID != nil {
		wordSelectedID = *raw.WordSelectedID
	}

	task := &domain.Task{
		ID:       raw.ID,
		Exercise: raw.Exercise,
		// Words, WordCorrect, WordSelected будут заполняться на уровне usecase, когда подгружаются словари по ID
	}

	return task, wordIDs, raw.WordCorrectID, wordSelectedID, nil
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
	found := slices.IndexFunc(task.Words, func(w domain.Dictionary) bool {
		return w.ID == dictId
	})
	if found == -1 {
		return domain.DictionaryNotFoundError
	}
	dict := task.Words[found]
	task.WordSelected = &dict

	return p.withTransaction(ctx, func(tx *sqlx.Tx) error {
		// обновить поле word_selected в таблице task
		_, err := tx.ExecContext(ctx, `UPDATE task SET word_selected = $1 WHERE id = $2`, dictId, task.ID)
		if err != nil {
			return fmt.Errorf("update task.word_selected: %w", err)
		}

		// инкрементировать selected_counter в exercise
		err = tx.GetContext(
			ctx, &task.Exercise.SelectedCounter,
			`UPDATE exercise SET selected_counter = selected_counter + 1, updated_at = NOW() WHERE id = $1 RETURNING selected_counter`,
			task.Exercise.ID,
		)
		if err != nil {
			return fmt.Errorf("increment selected_counter: %w", err)
		}

		// если выбрано правильно — инкрементировать corrected_counter
		if dictId == task.WordCorrect.ID {
			err = tx.GetContext(
				ctx,
				&task.Exercise.CorrectedCounter,
				`UPDATE exercise SET corrected_counter = corrected_counter + 1 WHERE id = $1 RETURNING corrected_counter`,
				task.Exercise.ID,
			)
			if err != nil {
				return fmt.Errorf("increment corrected_counter: %w", err)
			}
		}

		return nil
	})
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
		"word_correct": t.WordCorrect.ID,
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

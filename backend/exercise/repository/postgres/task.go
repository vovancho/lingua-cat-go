package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"slices"

	_watermillSql "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
	"github.com/vovancho/lingua-cat-go/pkg/txmanager"
)

type taskRepository struct {
	conn *sqlx.DB
	tx   *txmanager.Manager
}

func NewTaskRepository(conn *sqlx.DB, tx *txmanager.Manager) domain.TaskRepository {
	return &taskRepository{conn, tx}
}

func (r taskRepository) GetByID(ctx context.Context, id domain.TaskID) (*domain.TaskWithDetails, error) {
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

	if err := r.conn.GetContext(ctx, &raw, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found: %w", err)
		}

		return nil, fmt.Errorf("get task by id: %w", err)
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
	}

	return &domain.TaskWithDetails{
		Task:           task,
		WordIDs:        wordIDs,
		WordCorrectID:  raw.WordCorrectID,
		WordSelectedID: wordSelectedID,
	}, nil
}

func (r taskRepository) IsTaskOwnerExercise(ctx context.Context, exerciseID domain.ExerciseID, taskID domain.TaskID) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM task WHERE id = $1 AND exercise_id = $2)`

	var exists bool
	if err := r.conn.GetContext(ctx, &exists, query, taskID, exerciseID); err != nil {
		return false, fmt.Errorf("check task owner: %w", err)
	}

	return exists, nil
}

func (r taskRepository) Store(ctx context.Context, task *domain.Task) error {
	return r.tx.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		// Вставка задачи
		taskID, err := r.insertTask(ctx, tx, task)
		if err != nil {
			return fmt.Errorf("insert task: %w", err)
		}
		task.ID = taskID

		// Инкремент счетчика получения задачи
		counter, err := r.incrementProcessedCounter(ctx, tx, task.Exercise.ID)
		if err != nil {
			return fmt.Errorf("increment processed counter for task: %w", err)
		}
		task.Exercise.ProcessedCounter = counter

		return nil
	})
}

func (r taskRepository) SetWordSelected(
	ctx context.Context,
	task *domain.Task,
	dictId domain.DictionaryID,
	afterWordSetCallback func(ce _watermillSql.ContextExecutor, t domain.Task) error,
) error {
	found := slices.IndexFunc(task.Words, func(w domain.Dictionary) bool {
		return w.ID == dictId
	})

	if found == -1 {
		return domain.DictionaryNotFoundError
	}

	dict := task.Words[found]
	task.WordSelected = &dict

	return r.tx.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		// Обновление выбранного слова в задаче
		if err := r.updateTaskWordSelected(ctx, tx, task.ID, dictId); err != nil {
			return err
		}

		// Увеличение счетчика выбранных задач
		if err := r.incrementSelectedCounter(ctx, tx, task); err != nil {
			return err
		}

		// Если слово выбрано правильно, увеличить счетчик правильных ответов
		if dictId == task.WordCorrect.ID {
			if err := r.incrementCorrectedCounter(ctx, tx, task); err != nil {
				return err
			}
		}

		// Выполнение колбека после выбора слова, если он передан
		if afterWordSetCallback != nil {
			if err := afterWordSetCallback(tx, *task); err != nil {
				return fmt.Errorf("execute afterWordSetCallback: %w", err)
			}
		}

		return nil
	})
}

func (r taskRepository) insertTask(ctx context.Context, tx *sqlx.Tx, t *domain.Task) (domain.TaskID, error) {
	wordIDs := make([]int64, len(t.Words))
	for i, w := range t.Words {
		wordIDs[i] = int64(w.ID)
	}

	const query = `
		INSERT INTO task (exercise_id, words, word_correct)
		VALUES ($1, $2, $3)
		RETURNING id`

	var id domain.TaskID
	if err := tx.GetContext(ctx, &id, query, t.Exercise.ID, pq.Array(wordIDs), t.WordCorrect.ID); err != nil {
		return 0, fmt.Errorf("insert task: %w", err)
	}

	return id, nil
}

func (r taskRepository) incrementProcessedCounter(ctx context.Context, tx *sqlx.Tx, exerciseId domain.ExerciseID) (uint16, error) {
	const query = `
		UPDATE exercise
		SET processed_counter = processed_counter + 1
		WHERE id = $1 AND processed_counter < task_amount
		RETURNING processed_counter`

	var counter uint16
	if err := tx.GetContext(ctx, &counter, query, exerciseId); err != nil {
		return 0, fmt.Errorf("increment processed counter: %w", err)
	}

	return counter, nil
}

func (r taskRepository) updateTaskWordSelected(ctx context.Context, tx *sqlx.Tx, taskID domain.TaskID, dictID domain.DictionaryID) error {
	const query = `UPDATE task SET word_selected = $1 WHERE id = $2`
	if _, err := tx.ExecContext(ctx, query, dictID, taskID); err != nil {
		return fmt.Errorf("update task.word_selected: %w", err)
	}

	return nil
}

func (r taskRepository) incrementSelectedCounter(ctx context.Context, tx *sqlx.Tx, task *domain.Task) error {
	const query = `UPDATE exercise 
                   SET selected_counter = selected_counter + 1, updated_at = NOW() 
                   WHERE id = $1 
                   RETURNING selected_counter, updated_at`

	if err := tx.GetContext(ctx, &task.Exercise, query, task.Exercise.ID); err != nil {
		return fmt.Errorf("increment selected_counter: %w", err)
	}

	return nil
}

func (r taskRepository) incrementCorrectedCounter(ctx context.Context, tx *sqlx.Tx, task *domain.Task) error {
	const query = `UPDATE exercise SET corrected_counter = corrected_counter + 1 WHERE id = $1 RETURNING corrected_counter`
	if err := tx.GetContext(ctx, &task.Exercise.CorrectedCounter, query, task.Exercise.ID); err != nil {
		return fmt.Errorf("increment corrected_counter: %w", err)
	}

	return nil
}

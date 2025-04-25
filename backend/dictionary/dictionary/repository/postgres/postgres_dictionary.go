package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/vovancho/lingua-cat-go/dictionary/domain"
	"log/slog"
)

const (
	queryInsertDictionary         = `INSERT INTO dictionary (lang, name, type) VALUES ($1, $2, $3) RETURNING id`
	queryInsertTranslation        = `INSERT INTO translation (dictionary_id, translation_id) VALUES ($1, $2)`
	queryInsertSentence           = `INSERT INTO sentence (text_ru, text_en) VALUES ($1, $2) RETURNING id`
	queryInsertDictionarySentence = `INSERT INTO dictionary_sentence (dictionary_id, sentence_id) VALUES ($1, $2)`
	queryUpdateDictionaryName     = `UPDATE dictionary set name = $1 WHERE id = $2 AND deleted_at IS NULL`
)

type postgresDictionaryRepository struct {
	Conn *sql.DB
}

func NewPostgresDictionaryRepository(conn *sql.DB) domain.DictionaryRepository {
	return &postgresDictionaryRepository{conn}
}

func (p postgresDictionaryRepository) GetByID(ctx context.Context, id uint64) (res domain.Dictionary, err error) {
	query := `SELECT id, name, type, lang, created_at FROM dictionary WHERE id = ? AND deleted_at IS NULL`

	list, err := p.fetch(ctx, query, id)
	if err != nil {
		return domain.Dictionary{}, err
	}

	if len(list) == 0 {
		return res, domain.ErrNotFound
	}

	return
}

func (p postgresDictionaryRepository) Store(ctx context.Context, d *domain.Dictionary) (err error) {
	tx, err := p.Conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("failed to rollback transaction: %v; original error: %w", rollbackErr, err)
			}
			return
		}
		// Если ошибок нет, коммитим транзакцию
		if commitErr := tx.Commit(); commitErr != nil {
			err = fmt.Errorf("failed to commit transaction: %w", commitErr)
		}
	}()

	// Вставка основного словаря
	dictID, err := p.insertDictionary(ctx, tx, d)
	if err != nil {
		return err
	}
	d.ID = uint64(dictID)

	// Вставка переводов
	for i := range d.Translations {
		transDictID, err := p.insertDictionary(ctx, tx, &d.Translations[i].Dictionary)
		if err != nil {
			return err
		}
		d.Translations[i].Dictionary.ID = uint64(transDictID)

		// Вставка связей переводов (в обе стороны)
		if err := p.insertTranslation(ctx, tx, d.ID, d.Translations[i].Dictionary.ID); err != nil {
			return err
		}
		if err := p.insertTranslation(ctx, tx, d.Translations[i].Dictionary.ID, d.ID); err != nil {
			return err
		}
	}

	// Вставка предложений
	for i := range d.Sentences {
		sentenceID, err := p.insertSentence(ctx, tx, &d.Sentences[i])
		if err != nil {
			return err
		}

		// Вставка связей для основного словаря
		if err := p.insertDictionarySentence(ctx, tx, d.ID, sentenceID); err != nil {
			return err
		}

		// Вставка связей для переводов
		for _, trans := range d.Translations {
			if err := p.insertDictionarySentence(ctx, tx, trans.Dictionary.ID, sentenceID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p postgresDictionaryRepository) ChangeName(ctx context.Context, id uint64, name string) (err error) {
	// Обновление текста словаря
	if err := p.updateDictionaryName(ctx, p.Conn, id, name); err != nil {
		return err
	}

	return nil
}

func (p postgresDictionaryRepository) Delete(ctx context.Context, id uint64) (err error) {
	query := `UPDATE dictionary set deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`

	stmt, err := p.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAfected != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", rowsAfected)
		return
	}

	return
}

func (p postgresDictionaryRepository) fetch(ctx context.Context, query string, args ...any) (result []domain.Dictionary, err error) {
	rows, err := p.Conn.QueryContext(ctx, query, args...)

	if err != nil {
		slog.Error(err.Error())

		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error(err.Error())
		}
	}()

	result = make([]domain.Dictionary, 0)
	for rows.Next() {
		t := domain.Dictionary{}
		err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Type,
			&t.Lang,
			&t.CreatedAt,
			&t.DeletedAt,
		)

		if err != nil {
			slog.Error(err.Error())

			return nil, err
		}
	}

	return result, nil
}

func (p postgresDictionaryRepository) insertDictionary(ctx context.Context, tx *sql.Tx, d *domain.Dictionary) (int64, error) {
	var id int64
	err := tx.QueryRowContext(ctx, queryInsertDictionary, d.Lang, d.Name, d.Type).Scan(&id)
	if err != nil {
		slog.Error("failed to insert dictionary", "error", err)
		return 0, fmt.Errorf("insert dictionary: %w", err)
	}
	return id, nil
}

func (p postgresDictionaryRepository) insertTranslation(ctx context.Context, tx *sql.Tx, dictID, transID uint64) error {
	_, err := tx.ExecContext(ctx, queryInsertTranslation, dictID, transID)
	if err != nil {
		slog.Error("failed to insert translation", "error", err)
		return fmt.Errorf("insert translation: %w", err)
	}
	return nil
}

func (p postgresDictionaryRepository) insertSentence(ctx context.Context, tx *sql.Tx, s *domain.Sentence) (int64, error) {
	var id int64
	err := tx.QueryRowContext(ctx, queryInsertSentence, s.TextRU, s.TextEN).Scan(&id)
	if err != nil {
		slog.Error("failed to insert sentence", "error", err)
		return 0, fmt.Errorf("insert sentence: %w", err)
	}
	return id, nil
}

func (p postgresDictionaryRepository) insertDictionarySentence(ctx context.Context, tx *sql.Tx, dictID uint64, sentenceID int64) error {
	_, err := tx.ExecContext(ctx, queryInsertDictionarySentence, dictID, sentenceID)
	if err != nil {
		slog.Error("failed to insert dictionary_sentence", "error", err)
		return fmt.Errorf("insert dictionary_sentence: %w", err)
	}
	return nil
}

func (p postgresDictionaryRepository) updateDictionaryName(ctx context.Context, db *sql.DB, dictID uint64, dictName string) error {
	_, err := db.ExecContext(ctx, queryUpdateDictionaryName, dictName, dictID)
	if err != nil {
		slog.Error("failed to update dictionary name", "error", err)
		return fmt.Errorf("update dictionary name: %w", err)
	}
	return nil
}

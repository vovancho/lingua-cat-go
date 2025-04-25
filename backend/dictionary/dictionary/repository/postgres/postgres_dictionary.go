package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/vovancho/lingua-cat-go/dictionary/domain"
	"log/slog"
	"strings"
	"time"
)

const (
	queryInsertDictionary                        = `INSERT INTO dictionary (lang, name, type) VALUES ($1, $2, $3) RETURNING id`
	queryInsertTranslation                       = `INSERT INTO translation (dictionary_id, translation_id) VALUES ($1, $2)`
	queryInsertSentence                          = `INSERT INTO sentence (text_ru, text_en) VALUES ($1, $2) RETURNING id`
	queryInsertDictionarySentence                = `INSERT INTO dictionary_sentence (dictionary_id, sentence_id) VALUES ($1, $2)`
	queryUpdateDictionaryName                    = `UPDATE dictionary set name = $1 WHERE id = $2 AND deleted_at IS NULL`
	queryGetDictionaryById                       = `SELECT id, name, type, lang, deleted_at FROM dictionary WHERE id = $1 AND deleted_at IS NULL`
	queryGetTranslationsByDictionaryId           = `SELECT d.id, d.name, d.type, d.lang, d.deleted_at FROM dictionary d INNER JOIN translation t ON t.translation_id = d.id WHERE t.dictionary_id = $1 AND d.deleted_at IS NULL AND t.deleted_at IS NULL`
	queryGetSentencesByDictionaryIds             = `SELECT s.id, s.text_ru, s.text_en, ds.dictionary_id FROM sentence s INNER JOIN dictionary_sentence ds ON ds.sentence_id = s.id WHERE ds.dictionary_id IN (%s) AND s.deleted_at IS NULL ORDER BY ds.dictionary_id, s.id`
	queryGetDictionaryIdById                     = `SELECT id FROM dictionary WHERE id = $1 AND deleted_at IS NULL`
	queryGetIdsOfTranslationDictionariesToDelete = `SELECT d.id FROM dictionary d INNER JOIN translation t ON d.id = t.translation_id WHERE t.dictionary_id = $1 AND d.deleted_at IS NULL AND NOT EXISTS ( SELECT 1 FROM translation t2 WHERE t2.translation_id = d.id AND t2.dictionary_id <> $1 AND t2.deleted_at IS NULL)`
	queryDeleteTranslations                      = `UPDATE translation SET deleted_at = $1 WHERE dictionary_id IN (%s)`
	queryDeleteSentences                         = `UPDATE sentence s SET deleted_at = $1 FROM dictionary_sentence ds WHERE ds.sentence_id = s.id AND ds.dictionary_id IN (%s)`
	queryDeleteDictionaries                      = `UPDATE dictionary SET deleted_at = $1 WHERE id IN (%s)`
)

type postgresDictionaryRepository struct {
	Conn *sql.DB
}

func NewPostgresDictionaryRepository(conn *sql.DB) domain.DictionaryRepository {
	return &postgresDictionaryRepository{conn}
}

// GetByID возвращает словарь по ID с переводами и предложениями
func (p postgresDictionaryRepository) GetByID(ctx context.Context, id uint64) (dict domain.Dictionary, err error) {
	// Получаем основной словарь
	dict, err = p.getDictionaryByID(ctx, id)
	if err != nil {
		return domain.Dictionary{}, err
	}

	// Получаем переводы
	translations, err := p.getTranslationsByDictionaryID(ctx, id)
	if err != nil {
		return domain.Dictionary{}, err
	}
	dict.Translations = translations

	// Собираем ID словарей (основной + переводы)
	dictIDs := []uint64{dict.ID}
	for _, t := range translations {
		dictIDs = append(dictIDs, t.Dictionary.ID)
	}

	// Получаем предложения для всех словарей
	sentencesMap, err := p.getSentencesByDictionaryIDs(ctx, dictIDs)
	if err != nil {
		return domain.Dictionary{}, err
	}

	// Распределяем предложения по словарям
	dict.Sentences = sentencesMap[dict.ID]
	for i := range dict.Translations {
		dict.Translations[i].Dictionary.Sentences = sentencesMap[dict.Translations[i].Dictionary.ID]
	}

	return dict, nil
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

		if commitErr := tx.Commit(); commitErr != nil {
			err = fmt.Errorf("failed to commit transaction: %w", commitErr)
		}
	}()

	// Проверка существования словаря
	var dictID uint64
	err = tx.QueryRowContext(ctx, queryGetDictionaryIdById, id).Scan(&dictID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("dictionary with id %d not found", id)
	}
	if err != nil {
		slog.Error("failed to check dictionary existence", "error", err)
		return fmt.Errorf("check dictionary existence: %w", err)
	}

	// Получение переводов удаляемого словаря
	rows, err := tx.QueryContext(ctx, queryGetIdsOfTranslationDictionariesToDelete, id)
	if err != nil {
		slog.Error("failed to get translations", "error", err)
		return fmt.Errorf("get translations: %w", err)
	}
	defer rows.Close()

	// Собираем ID словарей (основной + переводы)
	dictIDs := []uint64{id}
	for rows.Next() {
		var transDictID uint64
		if err := rows.Scan(&transDictID); err != nil {
			slog.Error("failed to scan translation id", "error", err)
			return fmt.Errorf("scan translation id: %w", err)
		}
		dictIDs = append(dictIDs, transDictID)
	}

	placeholders := make([]string, len(dictIDs))
	args := make([]any, len(dictIDs)+1)
	args[0] = time.Now()
	for i, dictID2 := range dictIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args[i+1] = dictID2
	}

	// Удаление переводов (основной + переводы)
	query := fmt.Sprintf(queryDeleteTranslations, strings.Join(placeholders, ","))

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		slog.Error("failed to delete translation", "error", err)
		return fmt.Errorf("delete translation: %w", err)
	}

	// Удаление предложений (основной + переводы)
	query2 := fmt.Sprintf(queryDeleteSentences, strings.Join(placeholders, ","))

	_, err = tx.ExecContext(ctx, query2, args...)
	if err != nil {
		slog.Error("failed to delete sentences", "error", err)
		return fmt.Errorf("delete sentences: %w", err)
	}

	// Удаление словарей (основной + переводы)
	query3 := fmt.Sprintf(queryDeleteDictionaries, strings.Join(placeholders, ","))

	_, err = tx.ExecContext(ctx, query3, args...)
	if err != nil {
		slog.Error("failed to delete dictionary", "error", err)
		return fmt.Errorf("delete dictionary: %w", err)
	}

	return nil
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

// getDictionaryByID получает словарь по ID
func (p postgresDictionaryRepository) getDictionaryByID(ctx context.Context, id uint64) (domain.Dictionary, error) {
	row := p.Conn.QueryRowContext(ctx, queryGetDictionaryById, id)
	var d domain.Dictionary
	err := row.Scan(&d.ID, &d.Name, &d.Type, &d.Lang, &d.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Dictionary{}, fmt.Errorf("dictionary not found: %w", err)
		}
		slog.Error("failed to scan dictionary", "error", err)
		return domain.Dictionary{}, fmt.Errorf("scan dictionary: %w", err)
	}
	return d, nil
}

// getTranslationsByDictionaryID получает переводы для словаря
func (p postgresDictionaryRepository) getTranslationsByDictionaryID(ctx context.Context, dictionaryID uint64) ([]domain.Translation, error) {
	rows, err := p.Conn.QueryContext(ctx, queryGetTranslationsByDictionaryId, dictionaryID)
	if err != nil {
		slog.Error("failed to query translations", "error", err)
		return nil, fmt.Errorf("query translations: %w", err)
	}
	defer rows.Close()

	var translations []domain.Translation
	for rows.Next() {
		var d domain.Dictionary
		if err := rows.Scan(&d.ID, &d.Name, &d.Type, &d.Lang, &d.DeletedAt); err != nil {
			slog.Error("failed to scan translation dictionary", "error", err)
			return nil, fmt.Errorf("scan translation dictionary: %w", err)
		}
		translations = append(translations, domain.Translation{Dictionary: d})
	}
	if err := rows.Err(); err != nil {
		slog.Error("error iterating translations rows", "error", err)
		return nil, fmt.Errorf("iterate translations rows: %w", err)
	}
	return translations, nil
}

// getSentencesByDictionaryIDs получает предложения для списка ID словарей
func (p postgresDictionaryRepository) getSentencesByDictionaryIDs(ctx context.Context, dictIDs []uint64) (map[uint64][]domain.Sentence, error) {
	if len(dictIDs) == 0 {
		return make(map[uint64][]domain.Sentence), nil
	}

	placeholders := make([]string, len(dictIDs))
	args := make([]any, len(dictIDs))
	for i, id := range dictIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	query := fmt.Sprintf(queryGetSentencesByDictionaryIds, strings.Join(placeholders, ","))

	rows, err := p.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		slog.Error("failed to query sentences", "error", err)
		return nil, fmt.Errorf("query sentences: %w", err)
	}
	defer rows.Close()

	sentencesMap := make(map[uint64][]domain.Sentence)
	for rows.Next() {
		var s domain.Sentence
		var dictID uint64
		if err := rows.Scan(&s.ID, &s.TextRU, &s.TextEN, &dictID); err != nil {
			slog.Error("failed to scan sentence", "error", err)
			return nil, fmt.Errorf("scan sentence: %w", err)
		}
		sentencesMap[dictID] = append(sentencesMap[dictID], s)
	}
	if err := rows.Err(); err != nil {
		slog.Error("error iterating sentences rows", "error", err)
		return nil, fmt.Errorf("iterate sentences rows: %w", err)
	}
	return sentencesMap, nil
}

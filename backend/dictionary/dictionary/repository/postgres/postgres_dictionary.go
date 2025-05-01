package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/vovancho/lingua-cat-go/dictionary/internal/db"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vovancho/lingua-cat-go/dictionary/domain"
)

type postgresDictionaryRepository struct {
	Conn db.DB
}

func NewPostgresDictionaryRepository(conn db.DB) domain.DictionaryRepository {
	return &postgresDictionaryRepository{conn}
}

// GetByID возвращает словарь по ID с переводами и предложениями
func (p postgresDictionaryRepository) GetByID(ctx context.Context, id domain.DictionaryID) (*domain.Dictionary, error) {
	// Получаем основной словарь
	dict, err := p.getDictionaryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Получаем переводы
	translationsMap, err := p.getTranslationsByDictionariesIDs(ctx, []domain.DictionaryID{id})
	if err != nil {
		return nil, err
	}
	if len(translationsMap) == 0 {
		return nil, fmt.Errorf("translations not found")
	}
	dict.Translations = translationsMap[dict.ID]

	// Собираем ID словарей (основной + переводы)
	dictIDs := []domain.DictionaryID{dict.ID}
	for transDictId := range translationsMap {
		dictIDs = append(dictIDs, transDictId)
	}

	// Получаем предложения для всех словарей
	sentencesMap, err := p.getSentencesByDictionaryIDs(ctx, dictIDs)
	if err != nil {
		return nil, err
	}

	// Распределяем предложения по словарям
	if len(sentencesMap) > 0 {
		dict.Sentences = sentencesMap[dict.ID]
		for i := range dict.Translations {
			dict.Translations[i].Dictionary.Sentences = sentencesMap[dict.Translations[i].Dictionary.ID]
		}
	}

	return dict, nil
}

func (p postgresDictionaryRepository) IsExistsByNameAndLang(ctx context.Context, name string, lang domain.DictionaryLang) (bool, error) {
	const query = `SELECT id FROM dictionary WHERE name = $1 AND lang = $2 AND deleted_at IS NULL`

	var id domain.DictionaryID
	err := p.Conn.GetContext(ctx, &id, query, name, lang)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetRandomDictionaries возвращает случайный набор словарей по определенному языку с переводами и предложениями
func (p postgresDictionaryRepository) GetRandomDictionaries(ctx context.Context, lang domain.DictionaryLang, limit uint8) ([]domain.Dictionary, error) {
	// Получаем набор случайных словарей
	dicts, err := p.getDictionariesByLangAndRandomIDs(ctx, lang, limit)
	if err != nil {
		return nil, err
	}

	// Собираем ID словарей (основные)
	dictIDs := []domain.DictionaryID{}
	for _, dict := range dicts {
		dictIDs = append(dictIDs, dict.ID)
	}

	// Получаем переводы
	translationsMap, err := p.getTranslationsByDictionariesIDs(ctx, dictIDs)
	if err != nil {
		return nil, err
	}
	if len(translationsMap) == 0 {
		return nil, fmt.Errorf("translations not found")
	}
	for i := range dicts {
		dicts[i].Translations = translationsMap[dicts[i].ID]
	}

	// Собираем ID словарей (переводы)
	for transDictId := range translationsMap {
		dictIDs = append(dictIDs, transDictId)
	}

	// Получаем предложения для всех словарей
	sentencesMap, err := p.getSentencesByDictionaryIDs(ctx, dictIDs)
	if err != nil {
		return nil, err
	}

	// Распределяем предложения по словарям
	if len(sentencesMap) > 0 {
		for i := range dicts {
			dicts[i].Sentences = sentencesMap[dicts[i].ID]
			for j := range dicts[i].Translations {
				dicts[i].Translations[j].Dictionary.Sentences = sentencesMap[dicts[i].Translations[j].Dictionary.ID]
			}
		}
	}

	return dicts, nil
}

func (p postgresDictionaryRepository) Store(ctx context.Context, d *domain.Dictionary) error {
	return p.withTransaction(ctx, func(tx *sqlx.Tx) error {
		// Вставка основного словаря
		dictID, err := p.insertDictionary(ctx, tx, d)
		if err != nil {
			return err
		}
		d.ID = dictID

		// Вставка переводов
		for i := range d.Translations {
			transDictID, err := p.insertDictionary(ctx, tx, &d.Translations[i].Dictionary)
			if err != nil {
				return err
			}
			d.Translations[i].Dictionary.ID = transDictID

			// Вставка связей переводов (в обе стороны)
			if err = p.insertTranslation(ctx, tx, d.ID, d.Translations[i].Dictionary.ID); err != nil {
				return err
			}
			if err = p.insertTranslation(ctx, tx, d.Translations[i].Dictionary.ID, d.ID); err != nil {
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
			if err = p.insertDictionarySentence(ctx, tx, d.ID, sentenceID); err != nil {
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
	})
}

func (p postgresDictionaryRepository) ChangeName(ctx context.Context, id domain.DictionaryID, name string) error {
	const query = `UPDATE dictionary SET name = :name WHERE id = :id AND deleted_at IS NULL`
	res, err := p.Conn.NamedExecContext(ctx, query, map[string]any{
		"id":   id,
		"name": name,
	})
	if err != nil {
		return fmt.Errorf("update dictionary name: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if (rowsAffected == 0) || err != nil {
		return fmt.Errorf("update dictionary name not affected")
	}

	return nil
}

func (p postgresDictionaryRepository) Delete(ctx context.Context, id domain.DictionaryID) error {
	return p.withTransaction(ctx, func(tx *sqlx.Tx) error {
		// Проверка существования словаря
		if err := p.checkDictionaryExists(ctx, tx, id); err != nil {
			return err
		}

		// Получение ID словарей для удаления
		dictIDs, err := p.getDictionariesToDelete(ctx, tx, id)
		if err != nil {
			return err
		}

		// Удаление переводов
		if err = p.deleteTranslations(ctx, tx, dictIDs); err != nil {
			return err
		}

		// Удаление предложений
		if err = p.deleteSentences(ctx, tx, dictIDs); err != nil {
			return err
		}

		// Удаление словарей
		if err = p.deleteDictionaries(ctx, tx, dictIDs); err != nil {
			return err
		}

		return nil
	})
}

func (p postgresDictionaryRepository) insertDictionary(ctx context.Context, tx *sqlx.Tx, d *domain.Dictionary) (domain.DictionaryID, error) {
	var id domain.DictionaryID
	const query = `INSERT INTO dictionary (lang, name, type) VALUES (:lang, :name, :type) RETURNING id`
	nstmt, err := tx.PrepareNamed(query)
	if err != nil {
		return 0, fmt.Errorf("prepare named query: %w", err)
	}
	err = nstmt.QueryRowxContext(ctx, d).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert dictionary: %w", err)
	}
	return id, nil
}

func (p postgresDictionaryRepository) insertTranslation(ctx context.Context, tx *sqlx.Tx, dictID, transID domain.DictionaryID) error {
	const query = `INSERT INTO translation (dictionary_id, translation_id) VALUES (:dictionary_id, :translation_id)`
	_, err := tx.NamedExecContext(ctx, query, map[string]any{
		"dictionary_id":  dictID,
		"translation_id": transID,
	})
	if err != nil {
		return fmt.Errorf("insert translation: %w", err)
	}
	return nil
}

func (p postgresDictionaryRepository) insertSentence(ctx context.Context, tx *sqlx.Tx, s *domain.Sentence) (int64, error) {
	var id int64
	const query = `INSERT INTO sentence (text_ru, text_en) VALUES (:text_ru, :text_en) RETURNING id`
	nstmt, err := tx.PrepareNamed(query)
	if err != nil {
		return 0, fmt.Errorf("prepare named query: %w", err)
	}
	err = nstmt.QueryRowxContext(ctx, s).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert sentence: %w", err)
	}
	return id, nil
}

func (p postgresDictionaryRepository) insertDictionarySentence(ctx context.Context, tx *sqlx.Tx, dictID domain.DictionaryID, sentenceID int64) error {
	const query = `INSERT INTO dictionary_sentence (dictionary_id, sentence_id) VALUES (:dictionary_id, :sentence_id)`
	_, err := tx.NamedExecContext(ctx, query, map[string]any{
		"dictionary_id": dictID,
		"sentence_id":   sentenceID,
	})
	if err != nil {
		return fmt.Errorf("insert dictionary_sentence: %w", err)
	}
	return nil
}

// getSentencesByDictionaryIDs получает предложения для списка ID словарей
func (p postgresDictionaryRepository) getSentencesByDictionaryIDs(ctx context.Context, dictIDs []domain.DictionaryID) (map[domain.DictionaryID][]domain.Sentence, error) {
	if len(dictIDs) == 0 {
		return make(map[domain.DictionaryID][]domain.Sentence), nil
	}

	var query = `
		SELECT s.id,
		       s.text_ru,
		       s.text_en,
		       ds.dictionary_id
		FROM sentence s
		INNER JOIN dictionary_sentence ds ON ds.sentence_id = s.id
		WHERE ds.dictionary_id IN (?)
		  AND s.deleted_at IS NULL
		ORDER BY ds.dictionary_id, s.id`
	query, args, err := sqlx.In(query, dictIDs)
	if err != nil {
		return nil, fmt.Errorf("prepare IN query: %w", err)
	}
	query = p.Conn.Rebind(query)

	var results []struct {
		domain.Sentence
		DictionaryID domain.DictionaryID `db:"dictionary_id"`
	}
	if err := p.Conn.SelectContext(ctx, &results, query, args...); err != nil {
		return nil, fmt.Errorf("query sentences: %w", err)
	}

	sentencesMap := make(map[domain.DictionaryID][]domain.Sentence)
	for _, r := range results {
		sentencesMap[r.DictionaryID] = append(sentencesMap[r.DictionaryID], r.Sentence)
	}
	return sentencesMap, nil
}

// checkDictionaryExists проверяет существование словаря
func (p postgresDictionaryRepository) checkDictionaryExists(ctx context.Context, tx *sqlx.Tx, id domain.DictionaryID) error {
	var dictID domain.DictionaryID
	const query = `SELECT id FROM dictionary WHERE id = $1 AND deleted_at IS NULL`
	err := tx.GetContext(ctx, &dictID, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("dictionary with id %d not found", id)
		}
		return fmt.Errorf("check dictionary existence: %w", err)
	}
	return nil
}

// getDictionariesToDelete получает ID словарей для удаления
func (p postgresDictionaryRepository) getDictionariesToDelete(ctx context.Context, tx *sqlx.Tx, id domain.DictionaryID) ([]domain.DictionaryID, error) {
	var dictIDs []domain.DictionaryID
	const query = `
		SELECT d.id
		FROM dictionary d
		INNER JOIN translation t ON d.id = t.translation_id
		WHERE t.dictionary_id = $1
		  AND d.deleted_at IS NULL
		  AND NOT EXISTS (
			SELECT 1
			FROM translation t2
			WHERE t2.translation_id = d.id
			  AND t2.dictionary_id <> $1
			  AND t2.deleted_at IS NULL
		  )`

	if err := tx.GetContext(ctx, &dictIDs, query, id); err != nil {
		return nil, fmt.Errorf("get translations: %w", err)
	}
	dictIDs = append(dictIDs, id) // Добавляем основной словарь
	return dictIDs, nil
}

// deleteTranslations удаляет переводы
func (p postgresDictionaryRepository) deleteTranslations(ctx context.Context, tx *sqlx.Tx, dictIDs []domain.DictionaryID) error {
	if len(dictIDs) == 0 {
		return nil
	}
	var query = `UPDATE translation SET deleted_at = ? WHERE dictionary_id IN (?)`
	query, args, err := sqlx.In(query, time.Now(), dictIDs)
	if err != nil {
		return fmt.Errorf("prepare IN query: %w", err)
	}
	query = tx.Rebind(query)
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete translations: %w", err)
	}
	return nil
}

// deleteSentences удаляет предложения
func (p postgresDictionaryRepository) deleteSentences(ctx context.Context, tx *sqlx.Tx, dictIDs []domain.DictionaryID) error {
	if len(dictIDs) == 0 {
		return nil
	}
	var query = `UPDATE sentence s SET deleted_at = ? FROM dictionary_sentence ds WHERE ds.sentence_id = s.id AND ds.dictionary_id IN (?)`
	query, args, err := sqlx.In(query, time.Now(), dictIDs)
	if err != nil {
		return fmt.Errorf("prepare IN query: %w", err)
	}
	query = tx.Rebind(query)
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete sentences: %w", err)
	}
	return nil
}

// deleteDictionaries удаляет словари
func (p postgresDictionaryRepository) deleteDictionaries(ctx context.Context, tx *sqlx.Tx, dictIDs []domain.DictionaryID) error {
	if len(dictIDs) == 0 {
		return nil
	}
	var query = `UPDATE dictionary SET deleted_at = ? WHERE id IN (?)`
	query, args, err := sqlx.In(query, time.Now(), dictIDs)
	if err != nil {
		return fmt.Errorf("prepare IN query: %w", err)
	}
	query = tx.Rebind(query)
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete dictionaries: %w", err)
	}
	return nil
}

// withTransaction выполняет callback в контексте транзакции
func (p postgresDictionaryRepository) withTransaction(ctx context.Context, callback func(*sqlx.Tx) error) error {
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

func (p postgresDictionaryRepository) getDictionaryByID(ctx context.Context, id domain.DictionaryID) (*domain.Dictionary, error) {
	const query = `SELECT id, name, type, lang, deleted_at FROM dictionary WHERE id = $1 AND deleted_at IS NULL`
	var dict domain.Dictionary
	if err := p.Conn.GetContext(ctx, &dict, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("dictionary not found: %w", err)
		}
		return nil, fmt.Errorf("get dictionary: %w", err)
	}
	return &dict, nil
}

func (p postgresDictionaryRepository) getDictionariesByLangAndRandomIDs(ctx context.Context, lang domain.DictionaryLang, count uint8) ([]domain.Dictionary, error) {
	const query = `SELECT id, name, type, lang, deleted_at FROM dictionary WHERE lang = $1 AND deleted_at IS NULL ORDER BY RANDOM() LIMIT $2`
	var dicts []domain.Dictionary

	if err := p.Conn.SelectContext(ctx, &dicts, query, lang, count); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("dictionaries not found: %w", err)
		}
		return nil, fmt.Errorf("get random dictionaries: %w", err)
	}

	return dicts, nil
}

func (p postgresDictionaryRepository) getTranslationsByDictionariesIDs(ctx context.Context, dictIDs []domain.DictionaryID) (map[domain.DictionaryID][]domain.Translation, error) {
	var query = `
		SELECT t.id,
		       t.deleted_at,
		       d.id   AS "dictionary.id",
		       d.name AS "dictionary.name",
		       d.type AS "dictionary.type",
		       d.lang AS "dictionary.lang",
		       d.deleted_at AS "dictionary.deleted_at",
		       t.dictionary_id
		FROM dictionary d
		INNER JOIN translation t ON t.translation_id = d.id
		WHERE t.dictionary_id IN (?)
		  AND d.deleted_at IS NULL
		  AND t.deleted_at IS NULL`

	translationsMap := make(map[domain.DictionaryID][]domain.Translation)
	if len(dictIDs) == 0 {
		return translationsMap, nil
	}

	query, args, err := sqlx.In(query, dictIDs)
	if err != nil {
		return nil, fmt.Errorf("prepare IN query: %w", err)
	}
	query = p.Conn.Rebind(query)

	var results []struct {
		domain.Translation
		DictionaryID domain.DictionaryID `db:"dictionary_id"`
	}
	if err := p.Conn.SelectContext(ctx, &results, query, args...); err != nil {
		return nil, fmt.Errorf("query translations: %w", err)
	}

	for _, r := range results {
		translationsMap[r.DictionaryID] = append(translationsMap[r.DictionaryID], r.Translation)
	}
	return translationsMap, nil
}

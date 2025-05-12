package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vovancho/lingua-cat-go/dictionary/domain"
	"github.com/vovancho/lingua-cat-go/pkg/txmanager"
)

type repository struct {
	conn *sqlx.DB
	tx   *txmanager.Manager
}

func NewDictionaryRepository(conn *sqlx.DB, tx *txmanager.Manager) domain.DictionaryRepository {
	return &repository{conn, tx}
}

// GetByIDs возвращает словари по множеству ID с переводами и предложениями
func (r repository) GetByIDs(ctx context.Context, ids []domain.DictionaryID) ([]domain.Dictionary, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	// Получаем основные словари
	dictionaries, err := r.getDictionariesByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	// Составим карту словарей для быстрого доступа
	dictMap := make(map[domain.DictionaryID]*domain.Dictionary, len(dictionaries))
	var allDictIDs []domain.DictionaryID
	for _, dict := range dictionaries {
		dictMap[dict.ID] = dict
		allDictIDs = append(allDictIDs, dict.ID)
	}

	// Получаем переводы для всех словарей
	translationsMap, err := r.getTranslationsByDictionariesIDs(ctx, allDictIDs)
	if err != nil {
		return nil, err
	}

	// Добавляем переводы в словари и собираем все ID (основные и переводов)
	for dictID, translations := range translationsMap {
		if dict, ok := dictMap[dictID]; ok {
			dict.Translations = translations
			for _, t := range translations {
				allDictIDs = append(allDictIDs, t.Dictionary.ID)
			}
		}
	}

	// Получаем предложения
	sentencesMap, err := r.getSentencesByDictionaryIDs(ctx, allDictIDs)
	if err != nil {
		return nil, err
	}

	// Добавляем предложения в основные словари и их переводы
	for _, dict := range dictMap {
		dict.Sentences = sentencesMap[dict.ID]
		for i := range dict.Translations {
			dict.Translations[i].Dictionary.Sentences = sentencesMap[dict.Translations[i].Dictionary.ID]
		}
	}

	// Преобразуем карту обратно в слайс
	result := make([]domain.Dictionary, 0, len(dictMap))
	for _, dict := range dictMap {
		result = append(result, *dict)
	}

	return result, nil
}

func (r repository) IsExistsByNameAndLang(ctx context.Context, name string, lang domain.DictionaryLang) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM dictionary WHERE name = $1 AND lang = $2 AND deleted_at IS NULL)`

	var exists bool
	if err := r.conn.GetContext(ctx, &exists, query, name, lang); err != nil {
		return false, fmt.Errorf("check dictionary existence: %w", err)
	}

	return exists, nil
}

// GetRandomDictionaries возвращает случайный набор словарей по определенному языку с переводами и предложениями
func (r repository) GetRandomDictionaries(ctx context.Context, lang domain.DictionaryLang, limit uint8) ([]domain.Dictionary, error) {
	// Получаем набор случайных словарей
	dicts, err := r.getDictionariesByLangAndRandomIDs(ctx, lang, limit)
	if err != nil {
		return nil, err
	}

	// Собираем ID словарей (основные)
	dictIDs := []domain.DictionaryID{}
	for _, dict := range dicts {
		dictIDs = append(dictIDs, dict.ID)
	}

	// Получаем переводы
	translationsMap, err := r.getTranslationsByDictionariesIDs(ctx, dictIDs)
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
	sentencesMap, err := r.getSentencesByDictionaryIDs(ctx, dictIDs)
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

func (r repository) Store(ctx context.Context, d *domain.Dictionary) error {
	return r.tx.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		// Вставка основного словаря
		dictID, err := r.insertDictionary(ctx, tx, d)
		if err != nil {
			return err
		}
		d.ID = dictID

		// Вставка переводов
		for i := range d.Translations {
			transDictID, err := r.insertDictionary(ctx, tx, &d.Translations[i].Dictionary)
			if err != nil {
				return err
			}
			d.Translations[i].Dictionary.ID = transDictID

			// Вставка связей переводов (в обе стороны)
			if err = r.insertTranslation(ctx, tx, d.ID, d.Translations[i].Dictionary.ID); err != nil {
				return err
			}
			if err = r.insertTranslation(ctx, tx, d.Translations[i].Dictionary.ID, d.ID); err != nil {
				return err
			}
		}

		// Вставка предложений
		for i := range d.Sentences {
			sentenceID, err := r.insertSentence(ctx, tx, &d.Sentences[i])
			if err != nil {
				return err
			}

			// Вставка связей для основного словаря
			if err = r.insertDictionarySentence(ctx, tx, d.ID, sentenceID); err != nil {
				return err
			}

			// Вставка связей для переводов
			for _, trans := range d.Translations {
				if err := r.insertDictionarySentence(ctx, tx, trans.Dictionary.ID, sentenceID); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r repository) ChangeName(ctx context.Context, id domain.DictionaryID, name string) error {
	const query = `UPDATE dictionary SET name = :name WHERE id = :id AND deleted_at IS NULL`
	res, err := r.conn.NamedExecContext(ctx, query, map[string]any{
		"id":   id,
		"name": name,
	})
	if err != nil {
		return fmt.Errorf("update dictionary name: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("update dictionary name not affected")
	}

	return nil
}

func (r repository) Delete(ctx context.Context, id domain.DictionaryID) error {
	return r.tx.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		// Проверка существования словаря
		if err := r.checkDictionaryExists(ctx, tx, id); err != nil {
			return err
		}

		// Получение ID словарей для удаления
		dictIDs, err := r.getDictionariesToDelete(ctx, tx, id)
		if err != nil {
			return err
		}

		// Удаление переводов
		if err = r.deleteTranslations(ctx, tx, dictIDs); err != nil {
			return err
		}

		// Удаление предложений
		if err = r.deleteSentences(ctx, tx, dictIDs); err != nil {
			return err
		}

		// Удаление словарей
		if err = r.deleteDictionaries(ctx, tx, dictIDs); err != nil {
			return err
		}

		return nil
	})
}

func (r repository) insertDictionary(ctx context.Context, tx *sqlx.Tx, d *domain.Dictionary) (domain.DictionaryID, error) {
	const query = `INSERT INTO dictionary (lang, name, type) VALUES (:lang, :name, :type) RETURNING id`
	nstmt, err := tx.PrepareNamed(query)
	if err != nil {
		return 0, fmt.Errorf("prepare named query: %w", err)
	}

	var id domain.DictionaryID
	if err = nstmt.QueryRowxContext(ctx, d).Scan(&id); err != nil {
		return 0, fmt.Errorf("insert dictionary: %w", err)
	}

	return id, nil
}

func (r repository) insertTranslation(ctx context.Context, tx *sqlx.Tx, dictID, transID domain.DictionaryID) error {
	arg := map[string]any{
		"dictionary_id":  dictID,
		"translation_id": transID,
	}

	const query = `INSERT INTO translation (dictionary_id, translation_id) VALUES (:dictionary_id, :translation_id)`
	if _, err := tx.NamedExecContext(ctx, query, arg); err != nil {
		return fmt.Errorf("insert translation: %w", err)
	}

	return nil
}

func (r repository) insertSentence(ctx context.Context, tx *sqlx.Tx, s *domain.Sentence) (int64, error) {
	const query = `INSERT INTO sentence (text_ru, text_en) VALUES (:text_ru, :text_en) RETURNING id`
	nstmt, err := tx.PrepareNamed(query)
	if err != nil {
		return 0, fmt.Errorf("prepare named query: %w", err)
	}

	var id int64
	if err = nstmt.QueryRowxContext(ctx, s).Scan(&id); err != nil {
		return 0, fmt.Errorf("insert sentence: %w", err)
	}

	return id, nil
}

func (r repository) insertDictionarySentence(ctx context.Context, tx *sqlx.Tx, dictID domain.DictionaryID, sentenceID int64) error {
	arg := map[string]any{
		"dictionary_id": dictID,
		"sentence_id":   sentenceID,
	}

	const query = `INSERT INTO dictionary_sentence (dictionary_id, sentence_id) VALUES (:dictionary_id, :sentence_id)`
	if _, err := tx.NamedExecContext(ctx, query, arg); err != nil {
		return fmt.Errorf("insert dictionary_sentence: %w", err)
	}

	return nil
}

// getSentencesByDictionaryIDs получает предложения для списка ID словарей
func (r repository) getSentencesByDictionaryIDs(ctx context.Context, dictIDs []domain.DictionaryID) (map[domain.DictionaryID][]domain.Sentence, error) {
	if len(dictIDs) == 0 {
		return make(map[domain.DictionaryID][]domain.Sentence), nil
	}

	query := `
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
	query = r.conn.Rebind(query)

	var results []struct {
		domain.Sentence
		DictionaryID domain.DictionaryID `db:"dictionary_id"`
	}
	if err := r.conn.SelectContext(ctx, &results, query, args...); err != nil {
		return nil, fmt.Errorf("query sentences: %w", err)
	}

	sentencesMap := make(map[domain.DictionaryID][]domain.Sentence)
	for _, r := range results {
		sentencesMap[r.DictionaryID] = append(sentencesMap[r.DictionaryID], r.Sentence)
	}

	return sentencesMap, nil
}

// checkDictionaryExists проверяет существование словаря
func (r repository) checkDictionaryExists(ctx context.Context, tx *sqlx.Tx, id domain.DictionaryID) error {
	const query = `SELECT id FROM dictionary WHERE id = $1 AND deleted_at IS NULL`

	var dictID domain.DictionaryID
	if err := tx.GetContext(ctx, &dictID, query, id); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("dictionary with id %d not found", id)
		}

		return fmt.Errorf("check dictionary existence: %w", err)
	}

	return nil
}

// getDictionariesToDelete получает ID словарей для удаления
func (r repository) getDictionariesToDelete(ctx context.Context, tx *sqlx.Tx, id domain.DictionaryID) ([]domain.DictionaryID, error) {
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

	var dictIDs []domain.DictionaryID
	if err := tx.GetContext(ctx, &dictIDs, query, id); err != nil {
		return nil, fmt.Errorf("get translations: %w", err)
	}

	dictIDs = append(dictIDs, id) // Добавляем основной словарь

	return dictIDs, nil
}

// deleteTranslations удаляет переводы
func (r repository) deleteTranslations(ctx context.Context, tx *sqlx.Tx, dictIDs []domain.DictionaryID) error {
	if len(dictIDs) == 0 {
		return nil
	}

	query := `UPDATE translation SET deleted_at = ? WHERE dictionary_id IN (?)`
	query, args, err := sqlx.In(query, time.Now(), dictIDs)
	if err != nil {
		return fmt.Errorf("prepare IN query: %w", err)
	}
	query = tx.Rebind(query)

	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("delete translations: %w", err)
	}

	return nil
}

// deleteSentences удаляет предложения
func (r repository) deleteSentences(ctx context.Context, tx *sqlx.Tx, dictIDs []domain.DictionaryID) error {
	if len(dictIDs) == 0 {
		return nil
	}

	query := `UPDATE sentence s SET deleted_at = ? FROM dictionary_sentence ds WHERE ds.sentence_id = s.id AND ds.dictionary_id IN (?)`
	query, args, err := sqlx.In(query, time.Now(), dictIDs)
	if err != nil {
		return fmt.Errorf("prepare IN query: %w", err)
	}
	query = tx.Rebind(query)

	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("delete sentences: %w", err)
	}

	return nil
}

// deleteDictionaries удаляет словари
func (r repository) deleteDictionaries(ctx context.Context, tx *sqlx.Tx, dictIDs []domain.DictionaryID) error {
	if len(dictIDs) == 0 {
		return nil
	}

	query := `UPDATE dictionary SET deleted_at = ? WHERE id IN (?)`
	query, args, err := sqlx.In(query, time.Now(), dictIDs)
	if err != nil {
		return fmt.Errorf("prepare IN query: %w", err)
	}
	query = tx.Rebind(query)

	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("delete dictionaries: %w", err)
	}

	return nil
}

func (r repository) getDictionariesByIDs(ctx context.Context, ids []domain.DictionaryID) ([]*domain.Dictionary, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("dictionaries not found")
	}

	query := `
		SELECT id, name, type, lang, deleted_at
		FROM dictionary
		WHERE id IN (?) AND deleted_at IS NULL
	`

	query, args, err := sqlx.In(query, ids)
	if err != nil {
		return nil, fmt.Errorf("prepare IN query: %w", err)
	}
	query = r.conn.Rebind(query)

	var dictionaries []*domain.Dictionary
	if err := r.conn.SelectContext(ctx, &dictionaries, query, args...); err != nil {
		return nil, fmt.Errorf("get dictionaries: %w", err)
	}

	return dictionaries, nil
}

func (r repository) getDictionariesByLangAndRandomIDs(ctx context.Context, lang domain.DictionaryLang, count uint8) ([]domain.Dictionary, error) {
	const query = `SELECT id, name, type, lang, deleted_at FROM dictionary WHERE lang = $1 AND deleted_at IS NULL ORDER BY RANDOM() LIMIT $2`

	var dicts []domain.Dictionary
	if err := r.conn.SelectContext(ctx, &dicts, query, lang, count); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("dictionaries not found: %w", err)
		}
		return nil, fmt.Errorf("get random dictionaries: %w", err)
	}

	return dicts, nil
}

func (r repository) getTranslationsByDictionariesIDs(ctx context.Context, dictIDs []domain.DictionaryID) (map[domain.DictionaryID][]domain.Translation, error) {
	query := `
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
	query = r.conn.Rebind(query)

	var results []struct {
		domain.Translation
		DictionaryID domain.DictionaryID `db:"dictionary_id"`
	}
	if err := r.conn.SelectContext(ctx, &results, query, args...); err != nil {
		return nil, fmt.Errorf("query translations: %w", err)
	}

	for _, r := range results {
		translationsMap[r.DictionaryID] = append(translationsMap[r.DictionaryID], r.Translation)
	}

	return translationsMap, nil
}

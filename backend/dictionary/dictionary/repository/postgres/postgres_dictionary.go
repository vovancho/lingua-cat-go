package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/vovancho/lingua-cat-go/dictionary/domain"
	"log/slog"
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
	query := `INSERT INTO dictionary (lang, name, type) VALUES ($1, $2, $3) RETURNING id`
	stmt, err := p.Conn.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("Failed to prepare statement", "error", err)

		return
	}
	defer stmt.Close()

	var lastID int64
	err = stmt.QueryRowContext(ctx, d.Lang, d.Name, d.Type).Scan(&lastID)
	if err != nil {
		slog.Error("Failed to execute query", "error", err)

		return err
	}

	d.ID = uint64(lastID)

	return nil
}

func (p postgresDictionaryRepository) ChangeName(ctx context.Context, id uint64, name string) (err error) {
	query := `UPDATE dictionary set name = ? WHERE id = ? AND deleted_at IS NULL`

	stmt, err := p.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, name, id)
	if err != nil {
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return
	}

	if affect != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", affect)

		return
	}

	return
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

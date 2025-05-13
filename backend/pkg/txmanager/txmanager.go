package txmanager

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Manager struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Manager {
	return &Manager{db: db}
}

func (m *Manager) WithTransaction(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := m.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction rollback failed: %w", rbErr)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}

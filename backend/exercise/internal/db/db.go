package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error)
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	Rebind(query string) string
}

type DSN string

func NewDB(dsn DSN) (*sqlx.DB, error) {
	dbConn, err := sqlx.Open("postgres", string(dsn))
	if err != nil {
		return nil, fmt.Errorf("failed to configure database connection: %w", err)
	}

	if err = dbConn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database server: %w", err)
	}

	return dbConn, nil
}

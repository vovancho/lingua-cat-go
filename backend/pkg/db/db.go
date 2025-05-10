package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	//_ "github.com/lib/pq"
)

type DB interface {
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
	NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)
	NamedQueryContext(ctx context.Context, query string, arg any) (*sqlx.Rows, error)
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	Rebind(query string) string
}

type DriverName string
type DSN string

func NewDB(driverName DriverName, dsn DSN) (*sqlx.DB, error) {
	dbConn, err := sqlx.Open(string(driverName), string(dsn))
	if err != nil {
		return nil, fmt.Errorf("failed to configure database connection: %w", err)
	}

	if err = dbConn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database server: %w", err)
	}

	return dbConn, nil
}

package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDB(dsn string) (*sqlx.DB, error) {
	dbConn, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to configure database connection: %w", err)
	}

	if err = dbConn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database server: %w", err)
	}

	return dbConn, nil
}

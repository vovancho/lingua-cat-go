package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

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

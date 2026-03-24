package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Init opens and verifies a PostgreSQL connection pool.
// Returns *sql.DB which is safe for concurrent use and manages the pool internally.
func Init(dsn string) (*sql.DB, error) {
	database, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open failed: %w", err)
	}

	if err := database.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	// Connection pool settings
	database.SetMaxOpenConns(25)
	database.SetMaxIdleConns(10)

	return database, nil
}

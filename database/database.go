package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func Init() (*sql.DB, error) {
	// Open a connection to database
	connStr := "host=database user=username dbname=mailchump sslmode=disable password=password"
	db, err := sql.Open("database", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open a DB connection: %w", err)
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	return db, nil
}

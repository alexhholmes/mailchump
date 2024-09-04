package pgdb

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Client struct {
	db *sql.DB
}

func NewClient() (client *Client, err error) {
	// Open a connection to pgdb
	connStr := "host=postgres user=username dbname=mailchump sslmode=disable password=password"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return client, fmt.Errorf("failed to open a db connection: %w", err)
	}

	// Verify the connection
	if err = db.Ping(); err != nil {
		return client, fmt.Errorf("failed to ping db: %w", err)
	}

	return &Client{db: db}, nil
}

func (c *Client) Close() error {
	return c.db.Close()
}

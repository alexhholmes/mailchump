package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"mailchump/api"
	"mailchump/pgdb"
)

// init initializes the pgdb tables and adds test data. This is for local testing purposes
// only. Do not set the environment variable `INIT_DB` in non-local environments.
func init() {
	env := strings.ToLower(os.Getenv("ENVIRONMENT"))
	if env == "local" && os.Getenv("INIT_DB") != "" {
		// Initialize the pgdb tables and add test data
		db, err := pgdb.Init()
		if err != nil {
			log.Fatalf("Could not make db connection for dev environment initialization: %s", err.Error())
		}
		defer func(db *sql.DB) {
			_ = db.Close()
		}(db)

		// Read the tables.sql file
		query, err := os.ReadFile("tables.sql")
		if err != nil {
			log.Fatalf("Could not read migration file: %s", err.Error())
		}

		// Transaction to migrate entire tables.sql file the fresh database
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("Could not start migration transaction: %s", err.Error())
		}

		// Execute the tables.sql file statements to create the empty tables
		_, err = tx.Exec(string(query))
		if err != nil {
			roll := tx.Rollback()
			if roll != nil {
				log.Fatalf("Could not rollback migration transaction: %s", err.Error())
			}
			log.Fatalf("Could not execute migration statement: %s", err.Error())
		}

		// Commit the transaction
		err = tx.Commit()
		if err != nil {
			log.Fatalf("Could not commit migration transaction: %s", err.Error())
		}
	}
}

func main() {
	env := strings.ToLower(os.Getenv("ENVIRONMENT"))

	// Lets us initialize the pgdb and add test data without running the server for local
	// development.
	if env == "local" && os.Getenv("INIT_ONLY") != "" {
		return
	}

	// pprof web server. See: https://golang.org/pkg/net/http/pprof/
	if env == "local" || env == "dev" {
		go func() {
			pprof := "0.0.0.0:6060"
			slog.Info("server is listening", "pprof", pprof)
			log.Fatal(http.ListenAndServe(pprof, nil))
		}()
	}

	err := api.Run()
	if err != nil {
		os.Exit(1)
	}
}

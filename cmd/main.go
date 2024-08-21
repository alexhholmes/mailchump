package main

import (
	"database/sql"
	"log"
	"log/slog"
	"mailchump/pgdb"
	"net/http"
	"os"
	"strings"

	"mailchump/api"
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

		query := `
CREATE TABLE IF NOT EXISTS subscriptions (
id UUID PRIMARY KEY,
email VARCHAR(255) NOT NULL UNIQUE,
"from" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
until TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`
		_, err = db.Exec(query)
		if err != nil {
			log.Fatalf("Could not execute DB initialization query: %s", err.Error())
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

	// TODO implement os signal handling
	err := api.Run()
	if err != nil {
		os.Exit(1)
	}
}

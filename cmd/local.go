package main

import (
	"database/sql"
	"log"
	"log/slog"
	"mailchump/api"
	"mailchump/database"
	"net/http"
	_ "net/http/pprof"
	"os"
)

// init initializes the database tables and adds test data. This is for local testing
// purposes only. In a production environment, use `main.go` with a proper standalone
// database.
func init() {
	// Initialize the database tables and add test data
	db, err := database.Init()
	if err != nil {
		log.Fatalf("Could not make db connection for dev environment initialization: %s", err.Error())
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	query := `CREATE TABLE IF NOT EXISTS subscriptions (
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

func main() {
	// Lets us initialize the database and add test data without running the server.
	// Useful for GoLand's live SQL testing feature.
	if os.Getenv("INIT_ONLY") != "" {
		return
	}

	// pprof web server. See: https://golang.org/pkg/net/http/pprof/
	go func() {
		pprof := "0.0.0.0:6060"
		slog.Info("server is listening", "pprof", pprof)
		log.Fatal(http.ListenAndServe(pprof, nil))
	}()

	err := api.Run()
	if err != nil {
		os.Exit(1)
	}
}

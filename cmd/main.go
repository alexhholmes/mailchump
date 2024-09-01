package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"mailchump/pkg/api"
	"mailchump/pkg/pgdb"
)

// init initializes the pgdb tables and adds test data. This is for local testing purposes
// only. Do not set the environment variable `INIT_DB` in non-local environments.
func init() {
	env := strings.ToLower(os.Getenv("ENVIRONMENT"))
	if env == "local" && os.Getenv("INIT_DB") != "" {
		pgdb.InitializeLocalDB()
	}
	slog.Info("Database migration complete")

	// Lets us initialize the database and add test data without running the server; for
	// local development.
	if env == "local" && os.Getenv("INIT_ONLY") != "" {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

		<-signalChan
		os.Exit(0)
	}
}

func main() {
	env := strings.ToLower(os.Getenv("ENVIRONMENT"))

	// pprof web server. See: https://golang.org/pkg/net/http/pprof/
	if env == "local" || env == "dev" {
		go func() {
			pprof := "0.0.0.0:6060"
			slog.Info("Server is listening (pprof)", "address", pprof)
			log.Fatal(http.ListenAndServe(pprof, nil))
		}()
	}

	err := api.Run()
	if err != nil {
		os.Exit(1)
	}
}

package main

import (
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
		pgdb.InitializeLocalDB()
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

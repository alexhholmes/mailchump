package api

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"mailchump/api/healthcheck"
	"mailchump/api/newsletters"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"mailchump/gen"
	"mailchump/pgdb"
)

func Run() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	env := strings.ToLower(os.Getenv("ENVIRONMENT"))
	if env == "local" || env == "dev" {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	slog.SetDefault(logger)

	server, err := newServer()
	if err != nil {
		slog.Error("server fatal startup error", "error", err)
		return err
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			log.Fatalf("failed to close DB connection: %s", err.Error())
		}
	}(server.db)

	r := http.NewServeMux()

	// Get an `http.Handler` that we can use
	h := gen.HandlerFromMux(server, r)

	s := &http.Server{
		Handler: h,
		Addr:    "0.0.0.0:8080",
	}

	slog.Info("server is listening", "address", s.Addr)
	err = s.ListenAndServe()
	if err != nil {
		slog.Error(err.Error(), "error", err)
		slog.Error("server fatal runtime error", "error", err)
		return err
	}

	return nil
}

type server struct {
	db *sql.DB
	newsletters.NewsletterHandler
	healthcheck.HealthCheckHandler
}

func newServer() (server, error) {
	db, err := pgdb.Init()
	if err != nil {
		return server{}, fmt.Errorf("failed to open a DB connection: %w", err)
	}

	return server{
		db: db,
	}, nil
}

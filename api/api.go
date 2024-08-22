package api

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"mailchump/api/healthcheck"
	"mailchump/api/newsletters"
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

	server, closer, err := NewHandler()
	if err != nil {
		slog.Error("server fatal startup error", "error", err)
		return err
	}
	defer closer()

	// Get an `http.Handler` that we can use
	r := http.NewServeMux()
	h := gen.HandlerFromMux(&server, r)
	s := &http.Server{
		Handler: h,
		Addr:    "0.0.0.0:8080",
	}

	slog.Info("server is listening", "address", s.Addr)
	err = s.ListenAndServe()
	if err != nil {
		slog.Error("server fatal runtime error", "error", err)
		return err
	}

	return nil
}

type Handler struct {
	db *sql.DB
	newsletters.NewsletterHandler
	healthcheck.HealthCheckHandler
}

func NewHandler() (Handler, func(), error) {
	db, err := pgdb.Init()
	if err != nil {
		return Handler{}, nil, fmt.Errorf("failed to open a DB connection: %w", err)
	}

	return Handler{db: db}, func() {
		if db != nil {
			err = db.Close()
			if err != nil {
				log.Fatalf("failed to close DB connection: %s", err.Error())
			}
		}
	}, nil
}

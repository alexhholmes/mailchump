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
	"mailchump/gen"
	"mailchump/postgres"
)

func Run() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if strings.ToLower(os.Getenv("ENVIRONMENT")) == "dev" {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	slog.SetDefault(logger)

	server, err := newServer()
	if err != nil {
		slog.Error("Server fatal startup error", "error", err)
		return err
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(server.db)

	r := http.NewServeMux()

	// Get an `http.Handler` that we can use
	h := gen.HandlerFromMux(server, r)

	s := &http.Server{
		Handler: h,
		Addr:    "0.0.0.0:8080",
	}

	slog.Info("Server is listening", "address", s.Addr)
	err = s.ListenAndServe()
	if err != nil {
		slog.Error(err.Error(), "error", err)
		slog.Error("Server fatal runtime error", "error", err)
		return err
	}

	return nil
}

type Server struct {
	db *sql.DB
}

func newServer() (Server, error) {
	db, err := postgres.Init()
	if err != nil {
		return Server{}, fmt.Errorf("failed to open a DB connection: %w", err)
	}

	return Server{
		db: db,
	}, nil
}

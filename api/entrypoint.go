package api

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"mailchump/gen"
)

func Run() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if strings.ToLower(os.Getenv("ENVIRONMENT")) == "dev" {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	slog.SetDefault(logger)

	server, err := NewServer()
	if err != nil {
		slog.Error("server fatal startup error", "error", err)
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
		slog.Error("server fatal runtime error", "error", err)
		return err
	}

	return nil
}

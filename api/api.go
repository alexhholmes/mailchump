package api

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"mailchump/api/gen"
	"mailchump/api/healthcheck"
	"mailchump/api/newsletters"
	"mailchump/middleware"
	"mailchump/pgdb"
)

// Run starts the server, initializing the logger and a handler instance that will be
// used by the code-generated router.
func Run() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	env := strings.ToLower(os.Getenv("ENVIRONMENT"))
	if env == "local" || env == "dev" {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	slog.SetDefault(logger)

	server, closer, err := NewHandler()
	if err != nil {
		slog.Error("Server fatal startup error", "error", err)
		return err
	}
	defer closer()

	// Get an `http.Handler` that we can use
	r := http.NewServeMux()
	h := gen.HandlerWithOptions(
		&server, gen.StdHTTPServerOptions{
			BaseRouter: r,
			Middlewares: []gen.MiddlewareFunc{
				middleware.RecoveryMiddleware,
				middleware.LogRequestMiddleware,
				middleware.CreateAuthMiddleware(server.db),
			},
		},
	)
	s := &http.Server{
		Handler: h,
		Addr:    "0.0.0.0:8080",
	}

	// Used receive shutdown signal from SIGINT and SIGTERM
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("Server is listening", "address", s.Addr)
		err = s.ListenAndServe()
		if err != nil {
			log.Fatalf("Server fatal runtime error: %s", err.Error())
		}
	}()

	// Handle graceful shutdown
	sig := <-signalChan
	slog.Info("Server received shutdown signal", "signal", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = s.Shutdown(ctx); err != nil {
		log.Fatalf("Server failed to shutdown: %s", err.Error())
	}
	slog.Info("Server shutdown successfully")

	return nil
}

// Check that the Handler implements the generated API interface
var _ gen.ServerInterface = (*Handler)(nil)

// Handler is a composition of the endpoint handlers. This allows the individuals handlers
// to share the same resources, such as the database connection.
type Handler struct {
	db *sql.DB
	// cache *redis.Client
	newsletters.NewsletterHandler
	healthcheck.HealthHandler
}

// NewHandler creates a new Handler instance, initializing the database connection.
func NewHandler() (h Handler, close func(), err error) {
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

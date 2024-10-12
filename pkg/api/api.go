package api

import (
	"context"
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

	"mailchump/pkg/api/gen"
	"mailchump/pkg/api/healthcheck"
	"mailchump/pkg/api/newsletters"
	"mailchump/pkg/middleware"
	"mailchump/pkg/pgdb"
)

// Run starts the server, initializing the logger and a handler instance that will be
// used by the code-generated router.
func Run() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	env := strings.ToLower(os.Getenv("ENVIRONMENT"))
	if env == "local" || env == "dev" {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}
	slog.SetDefault(logger)

	server, closer, err := NewHandler()
	if err != nil {
		slog.Error("Server fatal startup error", "error", err)
		return err
	}
	defer closer()

	// Base router for page serving
	baseRouter := http.NewServeMux()
	baseRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, "templates/index.html")
	})

	// Create a sub-router with the generated OpenAPI spec. Register the API
	// routes, and strip the `/api` prefix since we don't specify it in the API
	// spec.
	h := gen.HandlerWithOptions(
		&server, gen.StdHTTPServerOptions{
			BaseRouter: http.NewServeMux(),
			Middlewares: []gen.MiddlewareFunc{
				middleware.RecoveryMiddleware,
				middleware.LogRequestMiddleware,
				middleware.CreateAuthMiddleware(),
			},
		},
	)
	baseRouter.Handle("/api/", http.StripPrefix("/api", h))
	baseRouter.Handle("/healthcheck", h)

	s := &http.Server{
		Handler: baseRouter,
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
	newsletters.NewsletterHandler
	healthcheck.HealthHandler
}

// NewHandler creates a new Handler instance, initializing the database connection.
func NewHandler() (h Handler, close func(), err error) {
	db, err := pgdb.NewClient()
	if err != nil {
		return Handler{}, nil, fmt.Errorf("failed to open a db connection: %w", err)
	}

	return Handler{
			NewsletterHandler: newsletters.NewsletterHandler{DB: db},
		}, func() {
			err = db.Close()
			if err != nil {
				log.Fatalf("failed to close db connection: %s", err.Error())
			}
		}, nil
}

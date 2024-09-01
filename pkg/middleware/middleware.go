package middleware

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"mailchump/pkg/api/util"
)

func CreateAuthMiddleware(db *sql.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// TODO authentication middleware for the api with jwt token from cookie added to context
				// TODO add user id to context as "user" key
				// TODO if the user is not authenticated, direct to login

				env := strings.ToLower(os.Getenv("ENV"))
				if env == "local" || env == "dev" {
					ctx := context.WithValue(r.Context(), util.ContextUser, "00000000-0000-0000-0000-000000000000")
					next.ServeHTTP(w, r.WithContext(ctx))
				}

				next.ServeHTTP(w, r)
			},
		)
	}
}

func HeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			contentType := strings.ToLower(r.Header.Get("Content-Type"))
			if contentType != "" && !strings.HasPrefix(contentType, "application/json") {
				http.Error(w, "Content-Type header must be application/json", http.StatusBadRequest)
				return
			}

			acceptType := r.Header.Get("Accept")
			if acceptType != "" && !strings.HasPrefix(acceptType, "application/json") {
				http.Error(w, "Accept-Encoding header must be application/json", http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")

			next.ServeHTTP(w, r)
		},
	)
}

func LogRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			slogWithReq := slog.Default().With(
				"method", r.Method,
				"path", r.URL.Path,
				"remote", r.RemoteAddr,
				"host", r.Host,
				"agent", r.UserAgent(),
				"tls", r.TLS,
				"proto", r.Proto,
				"user", func() string {
					if c, err := r.Cookie("user"); err == nil {
						return c.Value
					}
					return ""
				}(),
			)
			ctx := context.WithValue(r.Context(), util.ContextLogger, slogWithReq)

			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}

// RecoveryMiddleware recovers from panics, logs the error, and sends a 500 status code.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log, ok := r.Context().Value(util.ContextLogger).(*slog.Logger)
					if !ok {
						log = slog.Default()
					}

					log.Info("Recovered from panic: %v", err)
					http.Error(
						w,
						http.StatusText(http.StatusInternalServerError),
						http.StatusInternalServerError,
					)
				}
			}()
			next.ServeHTTP(w, r)
		},
	)
}

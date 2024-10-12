package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"mailchump/pkg/api/util"
)

func CreateAuthMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// TODO authentication middleware for the api with jwt token from cookie added to context
				// TODO add user id to context as "user" key
				// TODO if the user is not authenticated, direct to login

				env := strings.ToLower(os.Getenv("ENV"))
				if env == "local" || env == "dev" {
					ctx := context.WithValue(r.Context(), util.ContextUser, util.Key("00000000-0000-0000-0000-000000000000"))
					next.ServeHTTP(w, r.WithContext(ctx))
				}

				next.ServeHTTP(w, r)
			},
		)
	}
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
			slogWithReq.Info("Request")

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

					log.Error("Recovered from panic", "error", err)
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

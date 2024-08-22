package middleware

import (
	"net/http"

	"mailchump/api/gen"
)

// TODO authentication middleware for the api with jwt token from cookie added to context
// TODO add user id to context as "user" key

func AuthMiddleware() gen.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO extract jwt token from cookie
			// TODO validate jwt token
			// TODO add user id to context as "user" key
			next.ServeHTTP(w, r)
		}
	}
}
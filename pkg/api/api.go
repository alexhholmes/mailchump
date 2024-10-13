package api

import (
	_ "github.com/lib/pq"

	"mailchump/pkg/api/gen"
	"mailchump/pkg/api/healthcheck"
	"mailchump/pkg/api/newsletters"
	"mailchump/pkg/pgdb"
)

// Check that the Handler implements the generated API interface
var _ gen.ServerInterface = (*Handler)(nil)

// Handler is a composition of the endpoint handlers. This allows the individual
// handlers to share the same resources, such as the database connection.
type Handler struct {
	newsletters.NewsletterHandler
	healthcheck.HealthHandler
}

// NewHandler creates a new Handler instance, initializing the database
// connection.
func NewHandler(db *pgdb.Client) (h Handler) {
	return Handler{
		NewsletterHandler: newsletters.NewsletterHandler{DB: db},
	}
}

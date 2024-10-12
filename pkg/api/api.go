package api

import (
	"fmt"
	"log"

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

package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

type Server struct {
	db *sql.DB
}

func NewServer() (Server, error) {
	// Open a connection to postgres
	connStr := "user=username dbname=mailchump sslmode=disable password=password"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return Server{}, fmt.Errorf("failed to open a DB connection: %w", err)
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		return Server{}, fmt.Errorf("failed to ping DB: %w", err)
	}

	return Server{
		db: db,
	}, nil
}

// GetHealthcheck returns HTTP status 200.
// GET /ping
func (Server) GetHealthcheck(w http.ResponseWriter, r *http.Request) {
	resp := Health{
		Status: "pong",
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// PostSubscribe returns HTTP status 200.
// POST /subscribe
func (Server) PostSubscribe(w http.ResponseWriter, r *http.Request) {
	var req Subscription
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(BadRequest{Message: err.Error()})
		return
	}

	resp := SubscriptionResponse{
		Status: "subscribed",
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// PostUnsubscribe returns HTTP status 200.
// POST /unsubscribe
func (Server) PostUnsubscribe(w http.ResponseWriter, r *http.Request) {
	var req Subscription
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(BadRequest{Message: err.Error()})
		return
	}

	resp := SubscriptionResponse{
		Status: "unsubscribed",
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"mailchump/postgres"
	"net/http"

	_ "github.com/lib/pq"
	"mailchump/gen"
	"mailchump/model"
)

type Server struct {
	db *sql.DB
}

func NewServer() (Server, error) {
	db, err := postgres.Init()
	if err != nil {
		return Server{}, fmt.Errorf("failed to open a DB connection: %w", err)
	}

	return Server{
		db: db,
	}, nil
}

// PostSubscribe returns HTTP status 200.
// POST /subscribe
func (s Server) PostSubscribe(w http.ResponseWriter, r *http.Request) {
	var req gen.SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(gen.BadRequest{Message: err.Error()})
		return
	}

	var subscription model.Subscription
	subscription.FromReq(req)
	if err := subscription.Validate(); err != nil {
		slog.Info("Invalid subscription", "error", err)
	}
	if err := subscription.Create(s.db); errors.Is(err, model.ErrCreateSubscriptionAlreadyExists) {
		slog.Info("Subscription already exists", "error", err)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(gen.BadRequest{Message: err.Error()})
		return
	} else if err != nil {
		slog.Warn("Create subscription", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(gen.InternalServerError{Message: "Internal Server Error"})
		return
	}

	resp := gen.SubscriptionResponse{
		Status: "Subscribed",
		Email:  &subscription.Email,
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

// PostUnsubscribe returns HTTP status 200.
// POST /unsubscribe
func (s Server) PostUnsubscribe(w http.ResponseWriter, r *http.Request) {
	var req gen.SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(gen.BadRequest{Message: err.Error()})
		return
	}

	var subscription model.Subscription
	subscription.FromReq(req)
	if err := subscription.Validate(); err != nil {
		slog.Info("Invalid subscription", "error", err)
	}
	if err := subscription.Remove(s.db); errors.Is(err, model.ErrRemoveSubscriptionDoesNotExist) {
		slog.Info("Subscription does not exist", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(gen.BadRequest{Message: err.Error()})
		return
	} else if err != nil {
		slog.Warn("Remove subscription", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(gen.InternalServerError{Message: "Internal Server Error"})
		return
	}

	resp := gen.SubscriptionResponse{
		Status: "Unsubscribed",
		Email:  &req.Email,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// GetHealthcheck returns HTTP status 200.
// GET /ping
func (s Server) GetHealthcheck(w http.ResponseWriter, r *http.Request) {
	resp := gen.Health{
		Status: "OK",
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

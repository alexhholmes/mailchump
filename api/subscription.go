package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"mailchump/gen"
	"mailchump/model"
)

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

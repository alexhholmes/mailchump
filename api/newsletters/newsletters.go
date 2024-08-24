package newsletters

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"mailchump/api/gen"
	"mailchump/api/util"
	"mailchump/model"
)

type NewsletterHandler struct {
	db *sql.DB
}

func (h *NewsletterHandler) GetNewsletters(w http.ResponseWriter, r *http.Request) {
	newsletters := model.Newsletters{}
	err := newsletters.GetAll(r.Context(), h.db)
	if err != nil {
		slog.Warn("Failed to get newsletters", "error", err)
		http.Error(w, util.ErrInternalServerError.Error(), http.StatusInternalServerError)
		return
	}

	user := r.Context().Value(util.ContextUser).(util.Key)
	response := gen.AllNewsletterResponse{
		Count:       len(newsletters),
		Newsletters: newsletters.ToResponse(user),
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *NewsletterHandler) CreateNewsletter(w http.ResponseWriter, r *http.Request) {
	var newsletter model.Newsletter
	err := json.NewDecoder(r.Body).Decode(&newsletter)
	if err != nil {
		slog.Warn("Failed to decode request body", "error", err)
		http.Error(w, util.ErrMalformedRequest.Error(), http.StatusBadRequest)
		return
	}

	if err = newsletter.Create(r.Context(), h.db); err != nil {
		if errors.Is(err, ErrNewsletterAlreadyExists) {
			slog.Info("Create newsletter; newsletter already exists", "error", err)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		slog.Warn("Failed to create newsletter", "error", err)
		http.Error(w, util.ErrInternalServerError.Error(), http.StatusInternalServerError)
		return
	}

	user := r.Context().Value(util.ContextUser).(util.Key)
	response := newsletter.ToResponse(user)
	slog.Info("Create newsletter")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *NewsletterHandler) DeleteNewsletterById(w http.ResponseWriter, r *http.Request, id string) {
	//TODO implement me
	panic("implement me")
}

func (h *NewsletterHandler) GetNewsletterById(w http.ResponseWriter, r *http.Request, id string) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		slog.Info("Failed to parse id", "error", err)
		http.Error(w, util.ErrInvalidUUID.Error(), http.StatusBadRequest)
		return
	}

	newsletter := model.Newsletter{Id: parsed}
	err = newsletter.Get(r.Context(), h.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Info("Newsletter not found", "error", err)
			http.Error(w, ErrNewsletterNotFound.Error(), http.StatusNotFound)
			return
		}
		slog.Warn("Failed to get newsletter", "error", err)
		http.Error(w, util.ErrInternalServerError.Error(), http.StatusInternalServerError)
	}

	user := r.Context().Value(util.ContextUser).(util.Key)
	response := newsletter.ToResponse(user)
	slog.Info("Get newsletter by id", "owner", response.Owner)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *NewsletterHandler) HideNewsletter(w http.ResponseWriter, r *http.Request, id string) {
	//TODO implement me
	panic("implement me")
}

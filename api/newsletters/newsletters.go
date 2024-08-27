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
	log := r.Context().Value(util.ContextLogger).(slog.Logger)

	newsletters := model.Newsletters{}
	err := newsletters.GetAll(r.Context(), h.db)
	if err != nil {
		log.Warn("Failed to get newsletters", "error", err)
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
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
	log := r.Context().Value(util.ContextLogger).(slog.Logger)

	// TODO use newsletter request from gen
	newsletter := model.Newsletter{}
	err := json.NewDecoder(r.Body).Decode(&newsletter)
	if err != nil {
		log.Warn("Failed to decode request body", "error", err)
		http.Error(w, util.ErrMalformedRequest.Error(), http.StatusBadRequest)
		return
	}

	user := r.Context().Value(util.ContextUser)
	newsletter.OwnerID, err = uuid.Parse(user.(string))
	if err != nil {
		log.Error("Failed to parse user id", "error", err)
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	if err = newsletter.Create(r.Context(), h.db); err != nil {
		if errors.Is(err, model.ErrAlreadyExists) {
			log.Info("Create newsletter; newsletter already exists", "error", err)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		log.Warn("Failed to create newsletter", "error", err)
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	response := newsletter.ToResponse(user.(util.Key))
	log.Info(
		"Create newsletter",
		"id", newsletter.Id,
		"owner", newsletter.OwnerID,
		"title", newsletter.Title,
	)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *NewsletterHandler) DeleteNewsletterById(
	w http.ResponseWriter,
	r *http.Request,
	id string,
) {
	log := r.Context().Value(util.ContextLogger).(slog.Logger)

	parsed, err := uuid.Parse(id)
	if err != nil {
		log.Info("Failed to parse id", "error", err)
		http.Error(w, util.ErrInvalidUUID.Error(), http.StatusBadRequest)
		return
	}
	newsletter := model.Newsletter{Id: parsed}

	err = newsletter.Delete(r.Context(), h.db)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			log.Info("Newsletter not found", "error", err)
			http.Error(w, ErrNewsletterNotFound.Error(), http.StatusNotFound)
			return
		}
		log.Warn("Failed to delete newsletter", "error", err)
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}

	log.Info("Delete newsletter", "id", id)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(
		gen.StatusResponse{
			Status: "newsletter deleted",
		},
	)
}

func (h *NewsletterHandler) GetNewsletterById(
	w http.ResponseWriter,
	r *http.Request,
	id string,
) {
	log := r.Context().Value(util.ContextLogger).(slog.Logger)

	parsed, err := uuid.Parse(id)
	if err != nil {
		log.Info("Failed to parse id", "error", err)
		http.Error(w, util.ErrInvalidUUID.Error(), http.StatusBadRequest)
		return
	}
	newsletter := model.Newsletter{Id: parsed}

	err = newsletter.Get(r.Context(), h.db)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			log.Info("Newsletter not found", "error", err)
			http.Error(w, ErrNewsletterNotFound.Error(), http.StatusNotFound)
			return
		}
		log.Warn("Failed to get newsletter", "error", err)
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}

	user := r.Context().Value(util.ContextUser).(util.Key)
	response := newsletter.ToResponse(user)
	log.Info("Get newsletter by id", "owner", response.Owner)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *NewsletterHandler) HideNewsletter(
	w http.ResponseWriter,
	r *http.Request,
	id string,
) {
	log := r.Context().Value(util.ContextLogger).(slog.Logger)

	parsed, err := uuid.Parse(id)
	if err != nil {
		log.Info("Failed to parse id", "error", err)
		http.Error(w, util.ErrInvalidUUID.Error(), http.StatusBadRequest)
		return
	}
	newsletter := model.Newsletter{Id: parsed}

	if ok, err := newsletter.IsOwner(r.Context(), h.db); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			log.Info("Newsletter not found", "error", err)
			http.Error(w, ErrNewsletterNotFound.Error(), http.StatusNotFound)
			return
		}
		log.Warn("Failed to check if user is owner", "error", err)
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	} else if !ok {
		log.Info(
			"User is not the owner of the newsletter",
			"newsletter", id,
			"user", r.Context().Value(util.ContextUser),
			"owner", newsletter.OwnerID,
		)
		http.Error(w, util.ErrForbidden.Error(), http.StatusForbidden)
		return
	}

	// Check that the user is the newsletter owner
	err = newsletter.GetOwnerID(r.Context(), h.db)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			log.Info("Newsletter not found", "error", err)
			http.Error(w, ErrNewsletterNotFound.Error(), http.StatusNotFound)
			return
		}
		log.Warn("Failed to get newsletter", "error", err)
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}
	user := r.Context().Value(util.ContextUser).(string)
	if user != newsletter.OwnerID.String() {
		log.Info("User is not the owner of the newsletter",
			"user", user,
			"owner", newsletter.OwnerID,
		)
		http.Error(w, util.ErrForbidden.Error(), http.StatusForbidden)
	}

	err = newsletter.Hide(r.Context(), h.db)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			log.Info("Newsletter not found", "error", err)
			http.Error(w, ErrNewsletterNotFound.Error(), http.StatusNotFound)
			return
		}
		log.Warn("Failed to hide newsletter", "error", err)
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}

	log.Info("Hide newsletter",
		"id", id,
		"user", r.Context().Value(util.ContextUser),
		"hidden", newsletter.Hidden,
	)

	status := "newsletter hidden"
	if !newsletter.Hidden {
		status = "newsletter unhidden"
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(
		gen.StatusResponse{
			Status: status,
		},
	)
}

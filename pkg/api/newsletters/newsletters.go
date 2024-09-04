package newsletters

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"mailchump/pkg/api/gen"
	"mailchump/pkg/api/util"
	"mailchump/pkg/model"
	"mailchump/pkg/pgdb"
)

// NewsletterStore is an interface for the database operations required by the
// newsletter handlers. This is for mocking the database in tests.
type NewsletterStore interface {
	GetAllNewsletters(ctx context.Context) (model.Newsletters, error)
	GetNewsletterById(ctx context.Context, id string) (model.Newsletter, error)
	GetNewsletterOwnerID(ctx context.Context, id string) (uuid.UUID, error)
	DeleteNewsletter(ctx context.Context, id string) error
	HideNewsletter(ctx context.Context, id string) (isHidden bool, err error)
}

var _ NewsletterStore = &pgdb.Client{}

type NewsletterHandler struct {
	DB NewsletterStore
}

func (h *NewsletterHandler) GetNewsletters(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value(util.ContextLogger).(*slog.Logger)

	newsletters, err := h.DB.GetAllNewsletters(r.Context())
	if err != nil {
		log.Warn("Failed to get newsletters", "error", err)
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}

	response := gen.AllNewsletterResponse{
		Newsletters: []gen.NewsletterResponse{},
	}
	user, ok := r.Context().Value(util.ContextUser).(util.Key)
	if !ok {
		user = util.Key(uuid.Nil.String())
	}
	response.Newsletters = newsletters.ToResponse(user)
	response.Count = len(response.Newsletters)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *NewsletterHandler) DeleteNewsletterById(
	w http.ResponseWriter,
	r *http.Request,
	id string,

) {
	log := r.Context().Value(util.ContextLogger).(slog.Logger)

	_, err := uuid.Parse(id)
	if err != nil {
		log.Info("Failed to parse id", "error", err)
		http.Error(w, util.ErrInvalidUUID.Error(), http.StatusBadRequest)
		return
	}

	err = h.DB.DeleteNewsletter(r.Context(), id)
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

	_, err := uuid.Parse(id)
	if err != nil {
		log.Info("Failed to parse id", "error", err)
		http.Error(w, util.ErrInvalidUUID.Error(), http.StatusBadRequest)
		return
	}

	newsletter, err := h.DB.GetNewsletterById(r.Context(), id)
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

	_, err := uuid.Parse(id)
	if err != nil {
		log.Info("Failed to parse id", "error", err)
		http.Error(w, util.ErrInvalidUUID.Error(), http.StatusBadRequest)
		return
	}

	// Check that the user is the newsletter owner
	owner, err := h.DB.GetNewsletterOwnerID(r.Context(), id)
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
	if user != owner.String() {
		log.Info("User is not the owner of the newsletter",
			"user", user,
			"owner", owner.String(),
		)
		http.Error(w, util.ErrForbidden.Error(), http.StatusForbidden)
	}

	isHidden, err := h.DB.HideNewsletter(r.Context(), id)
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
		"hidden", isHidden,
	)

	status := "newsletter hidden"
	if !isHidden {
		status = "newsletter visible"
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(
		gen.StatusResponse{
			Status: status,
		},
	)
}

package healthcheck

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"mailchump/pkg/api/gen"
	"mailchump/pkg/api/util"
)

type HealthHandler struct{}

// GetHealthcheck returns HTTP status 200.
// GET /ping
func (h *HealthHandler) GetHealthcheck(w http.ResponseWriter, r *http.Request) {
	r.Context().Value(util.ContextLogger).(*slog.Logger).Info("Healthcheck")
	resp := gen.HealthCheck{
		Status: "ok",
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

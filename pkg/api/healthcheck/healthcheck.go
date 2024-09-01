package healthcheck

import (
	"encoding/json"
	"net/http"

	"mailchump/pkg/api/gen"
)

type HealthHandler struct{}

// GetHealthcheck returns HTTP status 200.
// GET /ping
func (h *HealthHandler) GetHealthcheck(w http.ResponseWriter, _ *http.Request) {
	resp := gen.HealthCheck{
		Status: "ok",
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

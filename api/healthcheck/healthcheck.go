package healthcheck

import (
	"encoding/json"
	"net/http"

	"mailchump/gen"
)

type HealthCheckHandler struct{}

// GetHealthcheck returns HTTP status 200.
// GET /ping
func (h HealthCheckHandler) GetHealthcheck(w http.ResponseWriter, r *http.Request) {
	resp := gen.HealthCheck{
		Status: "OK",
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

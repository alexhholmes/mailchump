package api

import (
	"encoding/json"
	"net/http"

	"mailchump/gen"
)

// GetHealthcheck returns HTTP status 200.
// GET /ping
func (s Server) GetHealthcheck(w http.ResponseWriter, r *http.Request) {
	resp := gen.Health{
		Status: "OK",
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

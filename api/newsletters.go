package api

import (
	"encoding/json"
	"mailchump/gen"
	"net/http"
)

// GetNewsletters returns a list of all newsletters.
func (s server) GetNewsletters(w http.ResponseWriter, r *http.Request) {
	// TODO implement me

	resp := struct {
		Count       int                      `json:"count"`
		Newsletters []gen.NewsletterResponse `json:"newsletters"`
	}{
		Count:       0,
		Newsletters: []gen.NewsletterResponse{},
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (s server) PostNewsletters(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (s server) DeleteNewslettersId(w http.ResponseWriter, r *http.Request, id string) {
	//TODO implement me
	panic("implement me")
}

func (s server) GetNewslettersId(w http.ResponseWriter, r *http.Request, id string) {
	//TODO implement me
	panic("implement me")
}

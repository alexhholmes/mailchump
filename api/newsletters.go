package api

import (
	"encoding/json"
	"net/http"

	"mailchump/gen"
	"mailchump/model"
)

func (s server) GetAllNewsletters(w http.ResponseWriter, r *http.Request) {
	newsletters := model.Newsletters{}
	err := newsletters.GetAllNewsletters(r.Context(), s.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(gen.InternalServerError{Message: "failed to retrieve newsletters"})
		return
	}

	resp := struct {
		Count       int                      `json:"count"`
		Newsletters []gen.NewsletterResponse `json:"newsletters"`
	}{
		Count:       len(newsletters),
		Newsletters: nil,
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
	//TODO implement me
	panic("implement me")
}

func (s server) CreateNewsletter(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (s server) DeleteNewsletterById(w http.ResponseWriter, r *http.Request, id string) {
	//TODO implement me
	panic("implement me")
}

func (s server) GetNewsletterById(w http.ResponseWriter, r *http.Request, id string) {
	//TODO implement me
	panic("implement me")
}

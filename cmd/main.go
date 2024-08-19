package main

import (
	"log"
	"net/http"

	"mailchump/api"
)

func main() {
	server, err := api.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	r := http.NewServeMux()

	// get an `http.Handler` that we can use
	h := api.HandlerFromMux(server, r)

	s := &http.Server{
		Handler: h,
		Addr:    "0.0.0.0:8080",
	}

	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())
}

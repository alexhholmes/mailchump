package main

import (
	"log"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"

	"mailchump/api"
)

func main() {
	// Lets us initialize the pgdb and add test data without running the server.
	// Useful for GoLand's live SQL testing feature.
	if os.Getenv("INIT_ONLY") != "" {
		return
	}

	// pprof web server. See: https://golang.org/pkg/net/http/pprof/
	go func() {
		pprof := "0.0.0.0:6060"
		slog.Info("server is listening", "pprof", pprof)
		log.Fatal(http.ListenAndServe(pprof, nil))
	}()

	err := api.Run()
	if err != nil {
		os.Exit(1)
	}
}

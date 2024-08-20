package main

import (
	"os"

	"mailchump/api"
)

func main() {
	// TODO implement os signal handling
	err := api.Run()
	if err != nil {
		os.Exit(1)
	}
}

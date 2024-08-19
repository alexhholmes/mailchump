package main

import (
	"mailchump/api"
	"os"
)

func main() {
	// TODO implement os signal handling
	err := api.Run()
	if err != nil {
		os.Exit(1)
	}
	return
}

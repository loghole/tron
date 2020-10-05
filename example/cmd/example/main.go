package main

import (
	_ "example/config"
	v1 "example/internal/app/controllers/v1"
	log "log"

	tron "github.com/loghole/tron"
)

func main() {
	app, err := tron.New()
	if err != nil {
		log.Fatalf("can't create app: %s", err)
	}

	// Init all ..

	var (
		strings = v1.NewStrings()
	)

	app.Run(strings)

	// Stop all...
}

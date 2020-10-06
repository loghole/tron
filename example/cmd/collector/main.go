package main

import (
	log "log"

	config "github.com/loghole/collector/config"
	v1 "github.com/loghole/collector/internal/app/controllers/v1"
	tron "github.com/loghole/tron"
)

func main() {
	app, err := tron.New()
	if err != nil {
		log.Fatalf("can't create app: %s", err)
	}

	log.Println(config.GetExampleValue())

	// Init all ..

	var (
		strings = v1.NewStrings()
	)

	app.Run(strings)

	// Stop all...
}

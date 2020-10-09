package main

import (
	log "log"

	config "github.com/loghole/example/config"
	v1 "github.com/loghole/example/internal/app/controllers/v1"
	tron "github.com/loghole/tron"
)

func main() {
	app, err := tron.New(tron.AddLogCaller())
	if err != nil {
		log.Fatalf("can't create app: %s", err)
	}

	app.Logger().Info(config.GetExampleValue())

	// Init all ..

	var (
		strings = v1.NewStrings()
	)

	if err := app.WithRunOptions().Run(strings); err != nil {
		app.Logger().Fatalf("can't run app: %v", err)
	}

	// Stop all...
}

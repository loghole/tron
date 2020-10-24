package main

import (
	config "example/config"
	stringsV1 "example/internal/app/controllers/strings/v1"
	log "log"

	tron "github.com/loghole/tron"
)

func main() {
	app, err := tron.New(tron.AddLogCaller())
	if err != nil {
		log.Fatalf("can't create app: %s", err)
	}

	defer app.Close()

	app.Logger().Info(config.GetExampleValue())

	// Init all ..

	var (
		stringsV1Handler = stringsV1.NewStrings()
	)

	if err := app.WithRunOptions().Run(stringsV1Handler); err != nil {
		app.Logger().Fatalf("can't run app: %v", err)
	}

	// Stop all...
}

package main

import (
	"log"

	"github.com/loghole/tron/example/config"

	"github.com/loghole/tron"
)

func main() {
	app, err := tron.New(tron.AddLogCaller())
	if err != nil {
		log.Fatalf("can't create app: %s", err)
	}

	defer app.Close()

	app.Logger().Info(config.GetExampleValue())

	// Init handlers
	var ()

	if err := app.WithRunOptions().Run(); err != nil {
		app.Logger().Fatalf("can't run app: %v", err)
	}
}

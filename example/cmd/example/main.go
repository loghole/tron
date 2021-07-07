package main

import (
	"log"

	"github.com/loghole/tron/example/config"
	stringsV1 "github.com/loghole/tron/example/internal/app/api/strings/v1"

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
	var (
		stringsV1Impl = stringsV1.NewImplementation()
	)

	if err := app.WithRunOptions().Run(stringsV1Impl); err != nil {
		app.Logger().Fatalf("can't run app: %v", err)
	}
}

package main

import (
	"log"

	"github.com/loghole/tron"
	"github.com/loghole/tron/transport"

	stringsV1 "github.com/loghole/tron/example/internal/app/api/strings/v1"
)

func main() {
	app, err := tron.New(tron.AddLogCaller(), tron.WithRealtimeConfig())
	if err != nil {
		log.Fatalf("can't create app: %v", err)
	}

	defer app.Close()

	handlers := []transport.Service{
		stringsV1.NewImplementation(),
	}

	if err := app.Run(handlers...); err != nil {
		app.Logger().Fatalf("can't run app: %v", err)
	}
}

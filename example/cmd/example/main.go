package main

import (
	"log"
	"os"

	"github.com/loghole/tron"
	"github.com/loghole/tron/transport"

	stringsV1 "github.com/loghole/tron/example/internal/app/api/strings/v1"
)

func main() {
	app, err := tron.New(
		tron.AddLogCaller(),
		tron.WithPublicHTTP(8080),
		tron.WithPublicGRPC(8081),
		tron.WithAdminHTTP(8082),
		tron.WithLoggerLevel(os.Getenv("LOGGER_LEVEL")),
	)
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

	if err := app.Wait(); err != nil {
		app.Logger().Errorf("wait: %v", err)
	}
}

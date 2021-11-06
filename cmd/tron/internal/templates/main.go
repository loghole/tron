package templates

const MainFile = `
package main

import (
	"log"

	"{{ .Module }}/config"

	"github.com/loghole/tron"
	"github.com/loghole/tron/transport"
)

func main() {
	app, err := tron.New(tron.AddLogCaller(), tron.WithRealtimeConfig())
	if err != nil {
		log.Fatalf("can't create app: %v", err)
	}

	defer app.Close()

	handlers := []transport.Service{
		// TODO: init handlers.
	}

	if err := app.Run(handlers...); err != nil {
		app.Logger().Fatalf("can't run app: %v", err)
	}
}
`
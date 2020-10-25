VERSION  := $(shell git describe --tags --always)
LDFLAGS  := -X 'github.com/loghole/tron/cmd/tron/internal/version.CliVersion=$(VERSION)'
TRON_BIN := $(GOPATH)/bin/tron

build:
	go build -o $(TRON_BIN) -ldflags "$(LDFLAGS)" *.go

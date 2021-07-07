VERSION  := $(shell git describe --tags --always)
LDFLAGS  := -X 'github.com/loghole/tron/cmd/tron/internal/version.CliVersion=$(VERSION)'

.PHONY: build
build:
	go build -o $(GOPATH)/bin/tron -ldflags "$(LDFLAGS)" *.go

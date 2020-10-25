VERSION  := $(shell git describe --tags --always)
LDFLAGS  := -X 'github.com/loghole/tron/cmd/tron/internal/version.CliVersion=$(VERSION)'

build:
	go build -o $(GOPATH)/bin/tron -ldflags "$(LDFLAGS)" *.go

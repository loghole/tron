VERSION  := $(shell git describe --tags --always | perl -pe 's/cmd\/tron\///')
LDFLAGS  := -X 'github.com/loghole/tron/cmd/tron/internal/command.cliVersion=$(VERSION)'
TRON_BIN := $(GOPATH)/bin/tron

build:
	go build -o $(TRON_BIN) -ldflags "$(LDFLAGS)" *.go

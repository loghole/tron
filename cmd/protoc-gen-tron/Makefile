VERSION  := $(shell git describe --tags --always)
LDFLAGS  := -X 'github.com/loghole/tron/cmd/protoc-gen-tron/internal/version.CliVersion=$(VERSION)'

build:
	go build -o protoc-gen-tron -ldflags "$(LDFLAGS)" *.go

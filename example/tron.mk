# Code generated by tron v0.17.1-rc1.1. DO NOT EDIT.
# You can extend or override anything in ./Makefile

LOCAL_BIN:=$(CURDIR)/bin

DOCKERFILE   = .deploy/docker/Dockerfile
DOCKER_IMAGE = loghole/tron/example

VERSION  := $(shell git describe --tags --always)
GIT_HASH := $(shell git rev-parse HEAD 2> /dev/null)
BUILD_TS := $(shell date +%FT%T%z)

LDFLAGS:=-X 'github.com/loghole/tron/internal/app.ServiceName=example' \
		 -X 'github.com/loghole/tron/internal/app.AppName=github.com/loghole/tron/example' \
		 -X 'github.com/loghole/tron/internal/app.GitHash=$(GIT_HASH)' \
		 -X 'github.com/loghole/tron/internal/app.Version=$(VERSION)' \
		 -X 'github.com/loghole/tron/internal/app.BuildAt=$(BUILD_TS)'

# generate code from proto
.PHONY: generate
generate:
	tron generate -v --proto=api

.PHONY: generate-config
generate-config:
	tron generate --config -v

.PHONY: gotest
gotest:
	go test -race -v -cover -coverprofile coverage.out ./...

.PHONY: lint
lint:
	golangci-lint run -v

.PHONY: docker-image
docker-image:
	docker build \
	-f $(DOCKERFILE) \
	-t $(DOCKER_IMAGE):latest \
	-t $(DOCKER_IMAGE):$(VERSION) \
	.

.PHONY: run-local
run-local:
	go run -ldflags "$(LDFLAGS)" cmd/example/main.go --local-config-enabled

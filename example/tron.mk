# Code generated by tron v0.3.0-4-g4cb4b95. DO NOT EDIT.
# You can extend or override anything in ./Makefile

LOCAL_BIN:=$(CURDIR)/bin

DOCKERFILE   = .deploy/docker/Dockerfile
DOCKER_IMAGE = example

VERSION  := $(shell git describe --exact-match --abbrev=0 --tags 2> /dev/null)
GIT_HASH := $(shell git rev-parse HEAD 2> /dev/null)
BUILD_TS := $(shell date +%FT%T%z)

LDFLAGS:=-X 'github.com/loghole/tron/internal/app.ServiceName=example' \
		 -X 'github.com/loghole/tron/internal/app.AppName=example' \
		 -X 'github.com/loghole/tron/internal/app.GitHash=$(GIT_HASH)' \
		 -X 'github.com/loghole/tron/internal/app.Version=$(VERSION)' \
		 -X 'github.com/loghole/tron/internal/app.BuildAt=$(BUILD_TS)'

.PHONY: .generate
.generate:
	tron generate --proto=api -v

# generate code from proto and config
.PHONY: generate
generate: .pb-deps .generate

# generate code from proto but without downloading proto deps
.PHONY: fast-generate
fast-generate: .generate

.PHONY: generate-config
generate-config:
	tron generate --config -v

# install proto dependencies
.PHONY: .pb-deps
.pb-deps:
	$(info #Installing proto dependencies...)
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

gotest:
	go test -race -v -cover -coverprofile coverage.out ./...

lint:
	golangci-lint run -v

docker-image:
	docker build \
	-f $(DOCKERFILE) \
	-t $(DOCKER_IMAGE) \
	-t $(DOCKER_IMAGE):$(VERSION) \
	.

run-local:
	go run -ldflags "$(LDFLAGS)" cmd/example/main.go --local-config-enabled

# Code generated by tron v0.26.0-rc1.1. DO NOT EDIT.
# You can extend or override anything in ./Makefile

LOCAL_BIN:=$(CURDIR)/bin

DOCKERFILE   = build/Dockerfile
DOCKER_IMAGE = loghole/tron/example

VERSION  := $(shell git describe --tags --always)
GIT_HASH := $$(git rev-parse HEAD)
BUILD_TS := $(shell date +%FT%T%z)

LDFLAGS:=-X 'github.com/loghole/tron/internal/app.ServiceName=example' \
		 -X 'github.com/loghole/tron/internal/app.AppName=github.com/loghole/tron/example' \
		 -X 'github.com/loghole/tron/internal/app.GitHash=$(GIT_HASH)' \
		 -X 'github.com/loghole/tron/internal/app.Version=$(VERSION)' \
		 -X 'github.com/loghole/tron/internal/app.BuildAt=$(BUILD_TS)'

DOCKER_COMPOSE_RUN ?= docker-compose

.PHONY: default
default: docker-compose docker-volumes tidy generate ## Init docker volumes, download deps and generate code

# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_\/-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: generate
generate: ## Generate code from proto
	tron generate -v --proto=api

.PHONY: tidy
tidy: ## Run go mod tidy
	${DOCKER_COMPOSE_RUN} run --rm --no-deps app /bin/sh -c "go mod tidy"

.PHONY: test
test: ## Run tests
	${DOCKER_COMPOSE_RUN} run --rm app /bin/sh -c "go test -race -v -cover -coverprofile coverage.out ./..."

.PHONY: lint
lint: ## Run linter
	${DOCKER_COMPOSE_RUN} run --rm linter /bin/sh -c "golangci-lint run ./... -c .golangci.yaml -v"

.PHONY: docker-image
docker-image: ## Create docker image
	docker build \
	-f $(DOCKERFILE) \
	-t $(DOCKER_IMAGE):latest \
	-t $(DOCKER_IMAGE):$(VERSION) \
	.

.PHONY: docker-volumes
docker-volumes: ## Create docker cache volumes
	docker volume create go-mod-cache
	docker volume create go-build-cache
	docker volume create go-lint-cache

.PHONY: docker-compose
docker-compose: ## Generate local docker-compose.override file
	test -s docker-compose.override.yaml || cp docker-compose.override.example.yaml docker-compose.override.yaml

.PHONY: docker-run
docker-run:
	${DOCKER_COMPOSE_RUN} run --rm --service-ports app /bin/sh -c "go run -ldflags \"$(LDFLAGS)\" cmd/example/main.go"

.PHONY: docker-down
docker-down: ## Down app containers
	${DOCKER_COMPOSE_RUN} down --volumes --remove-orphans

.PHONY: docker-rebuild
docker-rebuild: ## Rebuild development container
	${DOCKER_COMPOSE_RUN} build app

%:
	@true

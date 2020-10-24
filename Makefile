GO_TEST_PACKAGES = $(shell go list ./... | egrep -v '(pkg|cmd)')

gomod:
	go mod download

test:
	go test -race -v -cover -coverprofile coverage.out $(GO_TEST_PACKAGES)

lint:
	golangci-lint run -v
	cd cmd/tron && golangci-lint run -v

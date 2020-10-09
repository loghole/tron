GO_TEST_PACKAGES = $(shell go list ./... | egrep -v '(pkg|cmd)')

gomod:
	go mod download

gotest:
	go test -race -v -cover -coverprofile coverage.out $(GO_TEST_PACKAGES)

golint: .cmd_lint
	golangci-lint run -v

.cmd_lint:
	cd cmd/tron && golangci-lint run -v
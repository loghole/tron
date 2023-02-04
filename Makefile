GO_TEST_PACKAGES = $(shell go list ./... | egrep -v '(pkg|cmd)')

.PHONY: gomod
gomod:
	go mod download

.PHONY: test
test:
	go test -race -v -cover -coverprofile coverage.out $(GO_TEST_PACKAGES)

.PHONY: lint
lint:
	golangci-lint run -v
	cd cmd/tron && golangci-lint run -v --path-prefix=cmd/tron
	cd cmd/protoc-gen-tron && golangci-lint run -v --path-prefix=cmd/protoc-gen-tron

.PHONY: update-swagger
update-swagger:
	rm -fr /tmp/swagger-ui
	git clone https://github.com/swagger-api/swagger-ui.git /tmp/swagger-ui
	cd /tmp/swagger-ui; \
		mkdir ./html; \
		cat ./dist/index.html | perl -pe 's/https?:\/\/petstore.swagger.io\/v2\///g' > ./html/index.html; \
		cp ./dist/oauth2-redirect.html ./html; \
		cp ./dist/*.js ./html; \
		cp ./dist/*.css ./html; \
		cp ./dist/*.png ./html

	cp -R /tmp/swagger-ui/html ./internal/admin/
	rm -rf /tmp/swagger-ui

.PHONY: tidy
tidy:
	go mod tidy -compat=1.17
	cd cmd/tron && go mod tidy -compat=1.17
	cd cmd/protoc-gen-tron && go mod tidy -compat=1.17

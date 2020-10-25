GO_TEST_PACKAGES = $(shell go list ./... | egrep -v '(pkg|cmd)')

gomod:
	go mod download

test:
	go test -race -v -cover -coverprofile coverage.out $(GO_TEST_PACKAGES)

lint:
	golangci-lint run -v
	cd cmd/tron && golangci-lint run -v

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

	go-bindata -fs -o internal/admin/bindata.go -pkg "admin" -prefix "/tmp/swagger-ui/html" /tmp/swagger-ui/html/...
	rm -rf /tmp/swagger-ui

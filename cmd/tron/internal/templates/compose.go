package templates

const DefaultDockerCompose = `version: '3'
services:
  app:
    image: "{{ .ImageName }}:${GIT_COMMIT:-1}"
    build:
      context: .
      dockerfile: build/Dockerfile.dev
    container_name: "{{ .ImageName }}_${GIT_COMMIT:-1}-${BUILD_NUMBER:-1}"
    volumes:
      - ./:/src
      - go-mod-cache:/go/pkg
      - go-build-cache:/root/.cache/go-build
    working_dir: /src

  linter:
    image: golangci/golangci-lint:v1.51
    volumes:
      - ./:/src
      - go-mod-cache:/go/pkg
      - go-build-cache:/root/.cache/go-build
      - go-lint-cache:/root/.cache/golangci-lint
    working_dir: /src

volumes:
  go-mod-cache:
    external: true
  go-build-cache:
    external: true
  go-lint-cache:
    external: true
`

const DefaultOverrideDockerCompose = `version: '3'
services:
  app:
    environment:
      LOGGER_LEVEL: debug
    ports:
      - "127.0.0.1:8080:8080"
      - "127.0.0.1:8081:8081"
      - "127.0.0.1:8082:8082"
`

package templates

const DefaultDockerfileTemplate = `# BUILD BACKEND
FROM golang:{{ .GoVersion }}-alpine as builder

RUN apk add --no-cache ca-certificates tzdata git

# Create appuser.
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735RUN
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /src

COPY go.mod ./go.sum ./
RUN go mod download && go mod verify

COPY . .

ARG SERVICE_NAME={{ .Name }}
ARG APP_NAME={{ .Module }}
ARG GIT_HASH=unknown
ARG VERSION=unknown
ARG BUILD_TS=unknown

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s \
    -X github.com/loghole/tron/internal/app.ServiceName=$SERVICE_NAME \
    -X github.com/loghole/tron/internal/app.AppName=$APP_NAME \
    -X github.com/loghole/tron/internal/app.GitHash=${GIT_HASH} \
    -X github.com/loghole/tron/internal/app.Version=${VERSION} \
    -X github.com/loghole/tron/internal/app.BuildAt=${BUILD_TS} \
    " -o /app cmd/$SERVICE_NAME/*.go

# BUILD FINAL CONTAINER
FROM scratch as final
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /app /app

USER appuser:appuser

ENTRYPOINT ["/app"]
`

const DefaultDockerfileDev = `# Dockerfile for development
FROM golang:{{ .GoVersion }}-alpine

RUN apk add git make

ENV WORKDIR=/src
RUN mkdir -p ${WORKDIR}
WORKDIR ${WORKDIR}
`

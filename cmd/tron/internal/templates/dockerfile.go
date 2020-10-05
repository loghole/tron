package templates

const DefaultDockerfileTemplate = `# BUILD BACKEND
FROM golang:1.15-alpine as builder

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

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s \
    -X github.com/loghole/tron/internal/app.ServiceName=$SERVICE_NAME \
    -X github.com/loghole/tron/internal/app.AppName=$APP_NAME \
    -X github.com/loghole/tron/internal/app.GitHash=$(git rev-parse HEAD) \
    -X github.com/loghole/tron/internal/app.Version=$(git describe --tags --always) \
    -X github.com/loghole/tron/internal/app.BuildAt=$(date --utc +%FT%TZ) \
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

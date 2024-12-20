# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.23.3 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod verify
RUN go mod download

COPY .  ./

# Build main app
RUN CGO_ENABLED=0 GOOS=linux go build -o /sso ./cmd/rssaggregator

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /sso /sso
COPY --from=build-stage /app/config /config

USER nonroot:nonroot

ENTRYPOINT ["/sso", "--config","./config/local.yml"]
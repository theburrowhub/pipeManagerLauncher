# Description: Dockerfile to build the applications in the pipelineManager repository.

# Builder stage
FROM golang:1.22 AS builder

# Build arguments
ARG APP_VERSION=0.0.1

ENV APP_VERSION=${APP_VERSION}
ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /go/src/github.com/sergiotejon/pipeManager

COPY .. .

RUN go mod download

RUN go build \
      -a -installsuffix cgo \
      -ldflags "-X github.com/sergiotejon/pipeManager/internal/pkg/version.Version=${APP_VERSION}" \
      -o dashboard \
      cmd/dashboard/main.go

# Final stage
FROM alpine:3.20.3

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /go/src/github.com/sergiotejon/pipeManager/dashboard ./dashboard

ENTRYPOINT ["dashboard"]
CMD ["-c", "/etc/pipe-manager/config.yaml"]
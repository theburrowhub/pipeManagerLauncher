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
      -o pipeline-converter \
      cmd/pipeline-converter/main.go

# Final stage
FROM alpine:3.20.3

RUN apk --no-cache add \
    ca-certificates \
    openssh-client \
    git \
    bash

RUN echo 'eval $(ssh-agent -s)' >> /root/.bashrc

WORKDIR /app

COPY --from=builder /go/src/github.com/sergiotejon/pipeManager/pipeline-converter ./pipeline-converter

# Set the entrypoint to start the SSH agent
ENTRYPOINT ["/bin/bash", "-c", "source /root/.bashrc && exec \"$@\"", "--"]
CMD ["/app/pipeline-converter", "-c", "/etc/pipe-manager/config.yaml"]
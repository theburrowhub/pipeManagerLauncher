# Description: Dockerfile to build the applications in the pipelineManager repository.

# Builder stage
FROM golang:1.22 AS builder

# Build arguments
ARG APP_NAME=app
ARG APP_VERSION=0.0.1
ARG USER_NAME=sergiotejon
ARG REPORT_NAME=pipeManager

ENV APP_NAME=${APP_NAME}
ENV APP_VERSION=${APP_VERSION}
ENV USER_NAME=${USER_NAME}
ENV REPORT_NAME=${REPORT_NAME}

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /go/src/github.com/$USER_NAME/$REPORT_NAME

COPY . .

RUN go mod download

RUN go build \
      -a -installsuffix cgo \
      -ldflags "-X github.com/$USER_NAME/$REPORT_NAME/internal/pkg/version.Version=${APP_VERSION}" \
      -o ${APP_NAME} \
      cmd/${APP_NAME}/main.go

# Final stage
FROM alpine:3.20.3

# Build arguments
ARG APP_NAME=app
ARG APP_VERSION=0.0.1
ARG USER_NAME=sergiotejon
ARG REPORT_NAME=pipeManager

ENV APP_NAME=${APP_NAME}
ENV APP_VERSION=${APP_VERSION}
ENV USER_NAME=${USER_NAME}
ENV REPORT_NAME=${REPORT_NAME}

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /go/src/github.com/${USER_NAME}/${REPORT_NAME}/${APP_NAME} ./app

ENTRYPOINT ["./app"]
CMD ["-c", "./config/pipe-manager.conf", "-l", "80"]
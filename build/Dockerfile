# Description: Dockerfile to build the applications in the pipelineManager repository.

# Builder stage
FROM golang:1.22 AS builder

# Build arguments
ARG APP_NAME=dashBoard
ARG USER_NAME=sergiotejon
ARG REPORT_NAME=pipeManager

ENV APP_NAME=${APP_NAME}
ENV USER_NAME=${USER_NAME}
ENV REPORT_NAME=${REPORT_NAME}

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /go/src/github.com/$USER_NAME/$REPORT_NAME

COPY . .

RUN go mod download

RUN go build \
      -a -installsuffix cgo \
      -o ${APP_NAME} \
      cmd/${APP_NAME}/main.go

# Final stage
FROM alpine:latest

# Build arguments
ARG APP_NAME=dashBoard
ARG USER_NAME=sergiotejon
ARG REPORT_NAME=pipeManager

ENV APP_NAME=${APP_NAME}
ENV USER_NAME=${USER_NAME}
ENV REPORT_NAME=${REPORT_NAME}

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /go/src/github.com/${USER_NAME}/${REPORT_NAME}/${APP_NAME} .

CMD ./${APP_NAME}
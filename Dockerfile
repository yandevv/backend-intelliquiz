FROM golang:1.24.6-alpine3.22 AS base-image

RUN apk add --no-cache \
    libwebp-dev \
    gcc \
    musl-dev \
    pkgconfig

ENV CGO_ENABLED=1
ENV GOOS=linux

RUN go install github.com/swaggo/swag/cmd/swag@latest

FROM base-image AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY ./src ./src

WORKDIR /build/src

RUN swag init --parseDependency --parseInternal

RUN go build -o main

EXPOSE 8080

CMD ["/build/src/main", "--migrate=true", "--fresh=true"]
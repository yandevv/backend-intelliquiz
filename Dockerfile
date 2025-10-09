FROM golang:1.24.6-alpine3.22 AS base

RUN apk add --no-cache \
    libwebp-dev \
    gcc \
    musl-dev \
    pkgconfig

ENV CGO_ENABLED=1
ENV GOOS=linux

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY ./src ./src

WORKDIR /build/src

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init --parseDependency --parseInternal

RUN go build -o main

EXPOSE 8080

CMD ["/build/src/main", "--migrate=true", "--fresh=true"]
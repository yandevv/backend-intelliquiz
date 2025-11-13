FROM golang:1.25.3-alpine3.22 AS base-image

RUN apk add --no-cache \
    libwebp-dev \
    gcc \
    musl-dev \
    pkgconfig

ENV CGO_ENABLED=1
ENV GOOS=linux

RUN go install github.com/air-verse/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest

FROM base-image AS deps

WORKDIR /intelliquiz

COPY go.mod go.sum ./

RUN go mod download

COPY ./src ./src

WORKDIR /intelliquiz/src

RUN swag init --parseDependency --parseInternal

FROM deps AS build

RUN go build -o main

EXPOSE 8080

CMD ["/intelliquiz/src/main"]

FROM deps AS development

WORKDIR /intelliquiz/src

EXPOSE 8080

CMD ["air"]
# syntax=docker/dockerfile:1
FROM golang:tip-alpine3.24

ARG BOT_STATUS
ARG BOT_KEY
ENV BOT_STATUS=$BOT_STATUS
ENV BOT_KEY=$BOT_KEY

WORKDIR /app

COPY go.mod go.sum ./

RUN apk update \
    && apk upgrade \
    && apk add --no-cache sqlite \
    && go mod download

COPY . .

RUN go mod tidy \
    && go build -o main .

CMD ["./main"]
# syntax=docker/dockerfile:1
FROM golang:tip-alpine3.24

ARG BOT_STATUS
ENV BOT_STATUS=$BOT_STATUS

WORKDIR /app

COPY go.mod go.sum ./

RUN apk update \
    && apk upgrade \
    && apk add --no-cache sqlite \
    && go mod download

COPY . .

RUN --mount=type=secret,id=botkey,env=BOT_KEY ... \
&& go mod tidy \
&& go build -o main . \

CMD ["./main"]
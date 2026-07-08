# syntax=docker/dockerfile:1

FROM golang:tip-alpine3.24

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go mod tidy

RUN go build -o main .

CMD ["./main"]
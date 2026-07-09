# syntax=docker/dockerfile:1
FROM golang:1.25-alpine AS build

WORKDIR /src

# gcc/musl-dev are required to compile the cgo-based go-sqlite3 driver
RUN apk add --no-cache build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o /gofox .

FROM alpine:3.22

# BOT_KEY and BOT_STATUS are injected at runtime by the ECS task
# definition — never bake them into the image.
RUN adduser -D -H gofox

WORKDIR /app
COPY --from=build /gofox ./gofox

# /app must stay writable: the bot creates gofox.db in its working directory
RUN chown gofox:gofox /app
USER gofox

CMD ["./gofox"]

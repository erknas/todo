FROM golang:alpine AS builder

RUN apk add --no-cache \
     gcc \
     musl-dev

WORKDIR /app

COPY . /app

RUN CGO_ENABLED=1 go build -o todoapp cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/todoapp /app/todoapp
COPY web /app/web
COPY .env /app/.env

ARG TODO_PORT
EXPOSE $TODO_PORT

CMD ["./todoapp"]
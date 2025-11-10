FROM golang:1.25.1-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /coopera-bot ./cmd/main.go

FROM alpine:3.21

COPY --from=builder /coopera-bot /coopera-bot

ENTRYPOINT ["/coopera-bot"]

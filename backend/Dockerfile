FROM golang:1.25.1-alpine3.21 AS builder

WORKDIR /app

ARG VERSION="(devel)"
ARG COMMIT="none"

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="\
    -X 'main.version=${VERSION}' \
    -X 'main.commit=${COMMIT}'" \
    -o /coopera-backend ./cmd/main.go

FROM alpine:3.23

WORKDIR /app

COPY --from=builder /coopera-backend /coopera-backend
COPY migrations ./migrations

EXPOSE 8080

ENTRYPOINT ["/coopera-backend"]

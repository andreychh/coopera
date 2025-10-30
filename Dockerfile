FROM golang:1.25.1-alpine3.21 AS builder

WORKDIR /app

# Arguments that will be passed from the CI script
ARG VERSION="(devel)"
ARG COMMIT="none"

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="\
    -X 'main.version=${VERSION}' \
    -X 'main.commit=${COMMIT}'" \
    -o /coopera .

FROM alpine:3.22

COPY --from=builder /coopera /coopera

EXPOSE 8080

ENTRYPOINT ["/coopera"]

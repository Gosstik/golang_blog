FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN GOTOOLCHAIN=auto go mod download

COPY tests/loadtest/ tests/loadtest/

RUN GOTOOLCHAIN=auto go build -o /loadtest ./tests/loadtest/main.go

FROM alpine:3.20

COPY --from=builder /loadtest /loadtest

CMD ["/loadtest"]

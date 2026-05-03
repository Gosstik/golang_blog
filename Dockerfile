# Stage 1: Build
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN GOTOOLCHAIN=auto go mod download

COPY api/ api/
COPY internal/ internal/
RUN GOTOOLCHAIN=auto go build -o /blog ./internal/main.go

# Stage 2: Runtime
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /blog /app/blog
COPY api/json/ /app/api/json/
COPY contrib/swagger-ui/ /app/contrib/swagger-ui/

EXPOSE 50051 8090

CMD ["/app/blog"]

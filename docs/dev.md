# DEV

Initiate project:

```bash
go mod init github.com/Gosstik/golang_blog

go env -w GOTOOLCHAIN=auto # allows not to write GOTOOLCHAIN=auto before 'go install'
# GOTOOLCHAIN=auto go install <package>

# proto generation
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
# add to ~/.bashrc
export PATH=$(go env GOPATH)/bin:$PATH

go build
```

Add external library:

```bash
git submodule add https://github.com/googleapis/googleapis contrib/googleapis
```

Swagger UI was taken from here: <https://github.com/swagger-api/swagger-ui/tree/master/dist>.

## Run

```bash
docker compose up --build -d

# Seed db with some data (for smoke test)
docker compose exec -T postgres psql -U blog -d blog < postgresql/seeds/001_mock_data.sql

### Stop and clean

# Remove only containers
docker compose down
# Additionally remove volumes
docker compose down -v
# Also remove built blog image
docker compose down --rmi local
```

## Unit tests

```bash
go test -v ./internal/handlers/posts_service/...
```

## Smoke tests

```bash
# List posts on empty db:
curl -s -H "X-User-Uuid: 11111111-1111-1111-1111-111111111111" \
  "http://localhost:8090/v1/posts?limit=10"

# Create a post:
curl -s -X POST \
  -H "X-User-Uuid: 11111111-1111-1111-1111-111111111111" \
  -H "Content-Type: application/json" \
  -d '{"contentText":"Hello from the blog!"}' \
  "http://localhost:8090/v1/posts"

# List posts again:
curl -s -H "X-User-Uuid: 11111111-1111-1111-1111-111111111111" \
  "http://localhost:8090/v1/posts?limit=10" | python3 -m json.tool

# Like a post:
curl -s -X POST \
  -H "X-User-Uuid: 11111111-1111-1111-1111-111111111111" \
  "http://localhost:8090/v1/posts/<post-uuid>/like"

# Check prometheus metrics:
curl -s http://localhost:8090/metrics | grep grpc_server

# Swagger UI:
# Open http://localhost:8090/swagger-ui/ in browser
```

## Load tests

!!! To avoid port conflicts (because of include in docker-compose.yml) stop simple service running before loadtest.

```bash
docker compose -f tests/loadtest/docker-compose.yml up --build loadtest

# Stop and clean
docker compose -f tests/loadtest/docker-compose.yml down -v
# Also remove built images
docker compose -f tests/loadtest/docker-compose.yml down --rmi local
```

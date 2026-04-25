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

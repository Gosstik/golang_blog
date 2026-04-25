package posts_service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	api_proto "github.com/Gosstik/golang_blog/api/proto"
)

const (
	// UserUuidHeader is the metadata key used to pass user UUID from HTTP header.
	UserUuidHeader = "x-user-uuid"
)

// Handler implements api_proto.PostsServiceServer with hardcoded mock responses.
type Handler struct {
	api_proto.UnimplementedPostsServiceServer
}

func NewHandler() *Handler {
	return &Handler{}
}

// getUserUuid extracts user UUID from gRPC incoming metadata (forwarded from
// X-User-Uuid HTTP header by grpc-gateway).
func getUserUuid(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}

	values := md.Get(UserUuidHeader)
	if len(values) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing X-User-Uuid header")
	}

	return values[0], nil
}

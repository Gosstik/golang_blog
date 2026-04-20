package posts_service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api_proto "github.com/Gosstik/golang_blog/api/proto"
)

func (h *Handler) PostV1Delete(ctx context.Context, req *api_proto.V1PostsDeleteRequest) (*api_proto.V1PostsDeleteResponse, error) {
	_, err := getUserUuid(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetPostUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "post_uuid is required")
	}

	return &api_proto.V1PostsDeleteResponse{}, nil
}

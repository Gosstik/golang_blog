package posts_service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api_proto "github.com/Gosstik/golang_blog/api/proto"
)

func (h *Handler) PostV1LikesListUsers(ctx context.Context, req *api_proto.V1PostsLikesListUsersRequest) (*api_proto.V1PostsLikesListUsersResponse, error) {
	_, err := getUserUuid(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetPostUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "post_uuid is required")
	}

	limit := int(req.GetLimit())
	if limit <= 0 {
		return nil, status.Error(codes.InvalidArgument, "limit must be greater than 0")
	}

	users := mockUsers
	if limit > len(users) {
		limit = len(users)
	}
	users = users[:limit]

	return &api_proto.V1PostsLikesListUsersResponse{
		Users:  users,
		Cursor: nil,
	}, nil
}

package posts_service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api_proto "github.com/Gosstik/golang_blog/api/proto"
)

func (h *Handler) PostsV1List(ctx context.Context, req *api_proto.V1PostsListRequest) (*api_proto.V1PostsListResponse, error) {
	_, err := getUserUuid(ctx)
	if err != nil {
		return nil, err
	}

	limit := int(req.GetLimit())
	if limit <= 0 {
		return nil, status.Error(codes.InvalidArgument, "limit must be greater than 0")
	}

	posts := mockPosts
	if limit > len(posts) {
		limit = len(posts)
	}
	posts = posts[:limit]

	return &api_proto.V1PostsListResponse{
		Posts:  posts,
		Cursor: nil,
	}, nil
}

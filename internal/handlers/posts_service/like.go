package posts_service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api_proto "github.com/Gosstik/golang_blog/api/proto"
)

func (h *Handler) PostV1Like(ctx context.Context, req *api_proto.V1PostsLikeRequest) (*api_proto.V1PostsLikeResponse, error) {
	userUuid, err := getUserUuid(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetPostUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "post_uuid is required")
	}

	if _, err := h.postRepo.GetByUUID(ctx, req.GetPostUuid()); err != nil {
		return nil, status.Errorf(codes.NotFound, "post not found: %v", err)
	}

	count, err := h.likesRepo.Like(ctx, req.GetPostUuid(), userUuid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to like post: %v", err)
	}

	return &api_proto.V1PostsLikeResponse{
		LikesCount: count,
	}, nil
}

package posts_service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api_proto "github.com/Gosstik/golang_blog/api/proto"
)

func (h *Handler) PostV1Delete(ctx context.Context, req *api_proto.V1PostsDeleteRequest) (*api_proto.V1PostsDeleteResponse, error) {
	userUuid, err := getUserUuid(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetPostUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "post_uuid is required")
	}

	existing, err := h.postRepo.GetByUUID(ctx, req.GetPostUuid())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "post not found: %v", err)
	}

	if existing.AuthorUUID != userUuid {
		return nil, status.Error(codes.PermissionDenied, "you can only delete your own posts")
	}

	if err := h.postRepo.Delete(ctx, req.GetPostUuid()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete post: %v", err)
	}

	_ = h.cacheRepo.InvalidatePostsListCache(ctx)

	return &api_proto.V1PostsDeleteResponse{}, nil
}

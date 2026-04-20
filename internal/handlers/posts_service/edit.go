package posts_service

import (
	"context"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api_proto "github.com/Gosstik/golang_blog/api/proto"
	"github.com/Gosstik/golang_blog/internal/entities"
)

func (h *Handler) PostV1Edit(ctx context.Context, req *api_proto.V1PostsEditRequest) (*api_proto.V1PostsEditResponse, error) {
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
		return nil, status.Error(codes.PermissionDenied, "you can only edit your own posts")
	}

	imageUrls := req.GetContentImageUrls()
	if imageUrls == nil {
		imageUrls = []string{}
	}

	post := &entities.BlogPost{
		UUID:             req.GetPostUuid(),
		ContentText:      req.GetContentText(),
		ContentImageUrls: pq.StringArray(imageUrls),
	}

	if err := h.postRepo.Update(ctx, post); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update post: %v", err)
	}

	updated, err := h.postRepo.GetByUUID(ctx, req.GetPostUuid())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to reload post: %v", err)
	}

	likesCount, _ := h.likesRepo.GetLikesCount(ctx, updated.UUID)
	likedByMe, _ := h.likesRepo.IsLikedByUser(ctx, updated.UUID, userUuid)

	_ = h.cacheRepo.InvalidatePostsListCache(ctx)

	return &api_proto.V1PostsEditResponse{
		Post: postToProto(updated, likesCount, likedByMe),
	}, nil
}

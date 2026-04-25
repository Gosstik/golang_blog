package posts_service

import (
	"context"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api_proto "github.com/Gosstik/golang_blog/api/proto"
	"github.com/Gosstik/golang_blog/internal/entities"
)

func (h *Handler) PostV1Create(ctx context.Context, req *api_proto.V1PostsCreateRequest) (*api_proto.V1PostsCreateResponse, error) {
	userUuid, err := getUserUuid(ctx)
	if err != nil {
		return nil, err
	}

	_, err = h.userRepo.GetByUUID(ctx, userUuid)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	imageUrls := req.GetContentImageUrls()
	if imageUrls == nil {
		imageUrls = []string{}
	}

	post := &entities.BlogPost{
		AuthorUUID:       userUuid,
		ContentText:      req.GetContentText(),
		ContentImageUrls: pq.StringArray(imageUrls),
	}

	if err := h.postRepo.Create(ctx, post); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create post: %v", err)
	}

	// Reload with author.
	created, err := h.postRepo.GetByUUID(ctx, post.UUID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to reload post: %v", err)
	}

	// Invalidate main page cache.
	_ = h.cacheRepo.InvalidatePostsListCache(ctx)

	return &api_proto.V1PostsCreateResponse{
		Post: postToProto(created, 0, false),
	}, nil
}

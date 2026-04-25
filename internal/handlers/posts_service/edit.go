package posts_service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	api_proto "github.com/Gosstik/golang_blog/api/proto"
)

func (h *Handler) PostV1Edit(ctx context.Context, req *api_proto.V1PostsEditRequest) (*api_proto.V1PostsEditResponse, error) {
	userUuid, err := getUserUuid(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetPostUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "post_uuid is required")
	}

	author, err := findUserByUuid(userUuid)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	post := &api_proto.BlogPost{
		Uuid:             req.GetPostUuid(),
		Author:           author,
		CreatedAt:        mockPosts[0].CreatedAt,
		UpdatedAt:        timestamppb.Now(),
		LikesCount:       5,
		ContentText:      req.GetContentText(),
		ContentImageUrls: req.GetContentImageUrls(),
		LikedByMe:        false,
	}

	return &api_proto.V1PostsEditResponse{
		Post: post,
	}, nil
}

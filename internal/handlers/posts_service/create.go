package posts_service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	api_proto "github.com/Gosstik/golang_blog/api/proto"
)

func (h *Handler) PostV1Create(ctx context.Context, req *api_proto.V1PostsCreateRequest) (*api_proto.V1PostsCreateResponse, error) {
	userUuid, err := getUserUuid(ctx)
	if err != nil {
		return nil, err
	}

	author, err := findUserByUuid(userUuid)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	post := &api_proto.BlogPost{
		Uuid:             "post-uuid-new",
		Author:           author,
		CreatedAt:        timestamppb.Now(),
		LikesCount:       0,
		ContentText:      req.GetContentText(),
		ContentImageUrls: req.GetContentImageUrls(),
		LikedByMe:        false,
	}

	return &api_proto.V1PostsCreateResponse{
		Post: post,
	}, nil
}

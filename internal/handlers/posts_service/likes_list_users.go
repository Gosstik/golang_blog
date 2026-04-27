package posts_service

import (
	"context"
	"strconv"

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

	var redisCursor int64
	if c := req.GetCursor(); c != nil && c.GetValue() != "" {
		redisCursor, err = strconv.ParseInt(c.GetValue(), 10, 64)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid cursor: %v", err)
		}
	}

	userUUIDs, nextCursor, err := h.likesRepo.GetLikedUserUUIDs(ctx, req.GetPostUuid(), limit, redisCursor)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get liked users: %v", err)
	}

	if len(userUUIDs) == 0 {
		return &api_proto.V1PostsLikesListUsersResponse{
			Users:  nil,
			Cursor: nil,
		}, nil
	}

	users, err := h.userRepo.GetByUUIDs(ctx, userUUIDs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get users: %v", err)
	}

	pbUsers := make([]*api_proto.User, len(users))
	for i, u := range users {
		pbUsers[i] = userToProto(&u)
	}

	var pbCursor *api_proto.Cursor
	if nextCursor > 0 {
		pbCursor = &api_proto.Cursor{Value: strconv.FormatInt(nextCursor, 10)}
	}

	return &api_proto.V1PostsLikesListUsersResponse{
		Users:  pbUsers,
		Cursor: pbCursor,
	}, nil
}

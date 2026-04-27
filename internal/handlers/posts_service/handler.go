package posts_service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	api_proto "github.com/Gosstik/golang_blog/api/proto"
	"github.com/Gosstik/golang_blog/internal/repositories"
)

const (
	UserUuidHeader = "x-user-uuid"
)

type Handler struct {
	api_proto.UnimplementedPostsServiceServer

	postRepo  repositories.PostRepository
	userRepo  repositories.UserRepository
	likesRepo repositories.LikesRepository
	cacheRepo repositories.CacheRepository
}

func NewHandler(
	postRepo repositories.PostRepository,
	userRepo repositories.UserRepository,
	likesRepo repositories.LikesRepository,
	cacheRepo repositories.CacheRepository,
) *Handler {
	return &Handler{
		postRepo:  postRepo,
		userRepo:  userRepo,
		likesRepo: likesRepo,
		cacheRepo: cacheRepo,
	}
}

func getUserUuid(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}

	values := md.Get(UserUuidHeader)
	if len(values) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing X-User-Uuid header")
	}

	return values[0], nil
}

package posts_service

import (
	"context"
	"encoding/json"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	api_proto "github.com/Gosstik/golang_blog/api/proto"
	"github.com/Gosstik/golang_blog/internal/entities"
)

const mainPageLimit = 10
const mainPageCacheTTL = 30 * time.Second

// cachedPost is a lightweight struct for caching posts without likes data.
type cachedPost struct {
	UUID             string    `json:"uuid"`
	AuthorUUID       string    `json:"author_uuid"`
	AuthorNickname   string    `json:"author_nickname"`
	AuthorName       string    `json:"author_name"`
	AuthorSurname    string    `json:"author_surname"`
	AuthorAvatarUrl  string    `json:"author_avatar_url"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
	ContentText      string    `json:"content_text"`
	ContentImageUrls []string  `json:"content_image_urls"`
}

func (h *Handler) PostsV1List(ctx context.Context, req *api_proto.V1PostsListRequest) (*api_proto.V1PostsListResponse, error) {
	userUuid, err := getUserUuid(ctx)
	if err != nil {
		return nil, err
	}

	limit := int(req.GetLimit())
	if limit <= 0 {
		return nil, status.Error(codes.InvalidArgument, "limit must be greater than 0")
	}

	isMainPage := limit == mainPageLimit && req.GetCursor() == nil

	var cached []cachedPost

	// Try cache for main page.
	if isMainPage {
		data, err := h.cacheRepo.GetPostsListCache(ctx)
		if err == nil && data != nil {
			if json.Unmarshal(data, &cached) == nil {
				return h.enrichCachedPosts(ctx, cached, userUuid)
			}
		}
	}

	// Cache miss or not main page — query Postgres.
	var cursorTime *time.Time
	if c := req.GetCursor(); c != nil && c.GetValue() != "" {
		t, err := time.Parse(time.RFC3339Nano, c.GetValue())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid cursor: %v", err)
		}
		cursorTime = &t
	}

	posts, err := h.postRepo.List(ctx, limit+1, cursorTime)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list posts: %v", err)
	}

	var nextCursor *api_proto.Cursor
	if len(posts) > limit {
		last := posts[limit-1]
		nextCursor = &api_proto.Cursor{Value: last.CreatedAt.Format(time.RFC3339Nano)}
		posts = posts[:limit]
	}

	// Collect post UUIDs for batch Redis calls.
	postUUIDs := make([]string, len(posts))
	for i, p := range posts {
		postUUIDs[i] = p.UUID
	}

	likesCounts, err := h.likesRepo.GetLikesCountBatch(ctx, postUUIDs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get likes counts: %v", err)
	}

	likedByMe, err := h.likesRepo.IsLikedByUserBatch(ctx, postUUIDs, userUuid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get liked status: %v", err)
	}

	// Build response.
	pbPosts := make([]*api_proto.BlogPost, len(posts))
	for i, p := range posts {
		pbPosts[i] = postToProto(&p, likesCounts[p.UUID], likedByMe[p.UUID])
	}

	// Cache main page result (without likes data).
	if isMainPage && nextCursor == nil {
		h.cacheMainPage(ctx, posts)
	}

	return &api_proto.V1PostsListResponse{
		Posts:  pbPosts,
		Cursor: nextCursor,
	}, nil
}

func (h *Handler) enrichCachedPosts(ctx context.Context, cached []cachedPost, userUuid string) (*api_proto.V1PostsListResponse, error) {
	postUUIDs := make([]string, len(cached))
	for i, c := range cached {
		postUUIDs[i] = c.UUID
	}

	likesCounts, err := h.likesRepo.GetLikesCountBatch(ctx, postUUIDs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get likes counts: %v", err)
	}

	likedByMe, err := h.likesRepo.IsLikedByUserBatch(ctx, postUUIDs, userUuid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get liked status: %v", err)
	}

	pbPosts := make([]*api_proto.BlogPost, len(cached))
	for i, c := range cached {
		pb := &api_proto.BlogPost{
			Uuid: c.UUID,
			Author: &api_proto.User{
				Uuid:      c.AuthorUUID,
				Nickname:  c.AuthorNickname,
				Name:      c.AuthorName,
				Surname:   c.AuthorSurname,
				AvatarUrl: c.AuthorAvatarUrl,
			},
			CreatedAt:        timestamppb.New(c.CreatedAt),
			LikesCount:       likesCounts[c.UUID],
			ContentText:      c.ContentText,
			ContentImageUrls: c.ContentImageUrls,
			LikedByMe:        likedByMe[c.UUID],
		}
		if c.UpdatedAt != nil {
			pb.UpdatedAt = timestamppb.New(*c.UpdatedAt)
		}
		pbPosts[i] = pb
	}

	return &api_proto.V1PostsListResponse{
		Posts:  pbPosts,
		Cursor: nil,
	}, nil
}

func (h *Handler) cacheMainPage(ctx context.Context, posts []entities.BlogPost) {
	cached := make([]cachedPost, len(posts))
	for i, p := range posts {
		cached[i] = cachedPost{
			UUID:             p.UUID,
			AuthorUUID:       p.Author.UUID,
			AuthorNickname:   p.Author.Nickname,
			AuthorName:       p.Author.Name,
			AuthorSurname:    p.Author.Surname,
			AuthorAvatarUrl:  p.Author.AvatarUrl,
			CreatedAt:        p.CreatedAt,
			UpdatedAt:        p.UpdatedAt,
			ContentText:      p.ContentText,
			ContentImageUrls: []string(p.ContentImageUrls),
		}
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return
	}
	_ = h.cacheRepo.SetPostsListCache(ctx, data, mainPageCacheTTL)
}

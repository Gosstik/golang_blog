package posts_service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	api_proto "github.com/Gosstik/golang_blog/api/proto"
	"github.com/Gosstik/golang_blog/internal/entities"
	"github.com/lib/pq"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ---------------------------------------------------------------------------
// Repositories' mocks
// ---------------------------------------------------------------------------

type mockPostRepo struct {
	listFn     func(ctx context.Context, limit int, cursor *time.Time) ([]entities.BlogPost, error)
	listCalled bool
}

func (m *mockPostRepo) List(ctx context.Context, limit int, cursor *time.Time) ([]entities.BlogPost, error) {
	m.listCalled = true
	return m.listFn(ctx, limit, cursor)
}
func (m *mockPostRepo) GetByUUID(ctx context.Context, uuid string) (*entities.BlogPost, error) {
	return nil, nil
}
func (m *mockPostRepo) Create(ctx context.Context, post *entities.BlogPost) error { return nil }
func (m *mockPostRepo) Update(ctx context.Context, post *entities.BlogPost) error { return nil }
func (m *mockPostRepo) Delete(ctx context.Context, uuid string) error             { return nil }

type mockUserRepo struct{}

func (m *mockUserRepo) GetByUUID(ctx context.Context, uuid string) (*entities.User, error) {
	return &entities.User{UUID: uuid, Nickname: "test"}, nil
}
func (m *mockUserRepo) GetByUUIDs(ctx context.Context, uuids []string) ([]entities.User, error) {
	return nil, nil
}

type mockLikesRepo struct {
	counts map[string]int64
	liked  map[string]bool
}

func (m *mockLikesRepo) Like(ctx context.Context, postUUID, userUUID string) (int64, error) {
	return 0, nil
}
func (m *mockLikesRepo) Unlike(ctx context.Context, postUUID, userUUID string) (int64, error) {
	return 0, nil
}
func (m *mockLikesRepo) GetLikesCount(ctx context.Context, postUUID string) (int64, error) {
	return m.counts[postUUID], nil
}
func (m *mockLikesRepo) IsLikedByUser(ctx context.Context, postUUID, userUUID string) (bool, error) {
	return m.liked[postUUID], nil
}
func (m *mockLikesRepo) GetLikedUserUUIDs(ctx context.Context, postUUID string, limit int, cursor int64) ([]string, int64, error) {
	return nil, 0, nil
}
func (m *mockLikesRepo) GetLikesCountBatch(ctx context.Context, postUUIDs []string) (map[string]int64, error) {
	return m.counts, nil
}
func (m *mockLikesRepo) IsLikedByUserBatch(ctx context.Context, postUUIDs []string, userUUID string) (map[string]bool, error) {
	return m.liked, nil
}

type mockCacheRepo struct {
	data      []byte
	getCalled bool
	setCalled bool
}

func (m *mockCacheRepo) GetPostsListCache(ctx context.Context) ([]byte, error) {
	m.getCalled = true
	if m.data != nil {
		return m.data, nil
	}
	return nil, nil
}
func (m *mockCacheRepo) SetPostsListCache(ctx context.Context, data []byte, ttl time.Duration) error {
	m.setCalled = true
	return nil
}
func (m *mockCacheRepo) InvalidatePostsListCache(ctx context.Context) error { return nil }

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func ctxWithUser(userUUID string) context.Context {
	md := metadata.New(map[string]string{UserUuidHeader: userUUID})
	return metadata.NewIncomingContext(context.Background(), md)
}

func assertPostsEqual(t *testing.T, expected, actual []*api_proto.BlogPost) {
	t.Helper()
	if diff := cmp.Diff(expected, actual, protocmp.Transform()); diff != "" {
		t.Errorf("posts mismatch (-expected +actual):\n%s", diff)
	}
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestPostsV1List_DBPath_NonDefaultParams(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Microsecond)

	testPosts := []entities.BlogPost{
		{
			UUID:             "post-uuid-1",
			AuthorUUID:       "user-1",
			Author:           entities.User{UUID: "user-1", Nickname: "alice", Name: "Alice", Surname: "A"},
			CreatedAt:        now,
			ContentText:      "First post",
			ContentImageUrls: pq.StringArray{"https://img.example.com/1.png"},
		},
		{
			UUID:             "post-uuid-2",
			AuthorUUID:       "user-2",
			Author:           entities.User{UUID: "user-2", Nickname: "bob", Name: "Bob", Surname: "B"},
			CreatedAt:        now.Add(-time.Minute),
			ContentText:      "Second post",
			ContentImageUrls: pq.StringArray{},
		},
	}

	postRepo := &mockPostRepo{
		listFn: func(ctx context.Context, limit int, cursor *time.Time) ([]entities.BlogPost, error) {
			if limit != 6 {
				t.Errorf("expected limit=6 (5+1), got %d", limit)
			}
			return testPosts, nil
		},
	}

	likesRepo := &mockLikesRepo{
		counts: map[string]int64{"post-uuid-1": 42, "post-uuid-2": 7},
		liked:  map[string]bool{"post-uuid-1": true, "post-uuid-2": false},
	}

	cacheRepo := &mockCacheRepo{data: nil} // cache miss

	handler := NewHandler(postRepo, &mockUserRepo{}, likesRepo, cacheRepo)

	// Non-default params (cursor + limit=5)
	cursorTime := now.Add(-30 * time.Second).Format(time.RFC3339Nano)
	req := &api_proto.V1PostsListRequest{
		Limit:  5,
		Cursor: &api_proto.Cursor{Value: cursorTime},
	}

	resp, err := handler.PostsV1List(ctxWithUser("user-1"), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !postRepo.listCalled {
		t.Error("expected PostRepository.List to be called")
	}

	if cacheRepo.getCalled {
		t.Error("cache should NOT be checked for non-default params (isMainPage=false)")
	}

	expectedPosts := []*api_proto.BlogPost{
		{
			Uuid:             "post-uuid-1",
			Author:           &api_proto.User{Uuid: "user-1", Nickname: "alice", Name: "Alice", Surname: "A"},
			CreatedAt:        timestamppb.New(now),
			LikesCount:       42,
			ContentText:      "First post",
			ContentImageUrls: []string{"https://img.example.com/1.png"},
			LikedByMe:        true,
		},
		{
			Uuid:             "post-uuid-2",
			Author:           &api_proto.User{Uuid: "user-2", Nickname: "bob", Name: "Bob", Surname: "B"},
			CreatedAt:        timestamppb.New(now.Add(-time.Minute)),
			LikesCount:       7,
			ContentText:      "Second post",
			ContentImageUrls: []string{},
			LikedByMe:        false,
		},
	}

	assertPostsEqual(t, expectedPosts, resp.Posts)

	if resp.Cursor != nil {
		t.Error("expected no cursor when posts < limit")
	}
}

func TestPostsV1List_CachePath_DefaultParams(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Microsecond)

	cached := []cachedPost{
		{
			UUID:            "cached-post-1",
			AuthorUUID:      "user-1",
			AuthorNickname:  "alice",
			AuthorName:      "Alice",
			AuthorSurname:   "A",
			AuthorAvatarUrl: "https://example.com/alice.png",
			CreatedAt:       now,
			ContentText:     "Cached post 1",
		},
		{
			UUID:            "cached-post-2",
			AuthorUUID:      "user-2",
			AuthorNickname:  "bob",
			AuthorName:      "Bob",
			AuthorSurname:   "B",
			AuthorAvatarUrl: "https://example.com/bob.png",
			CreatedAt:       now.Add(-time.Minute),
			ContentText:     "Cached post 2",
		},
	}

	cacheData, err := json.Marshal(cached)
	if err != nil {
		t.Fatalf("failed to marshal cache data: %v", err)
	}

	postRepo := &mockPostRepo{
		listFn: func(ctx context.Context, limit int, cursor *time.Time) ([]entities.BlogPost, error) {
			t.Fatal("PostRepository.List should NOT be called when cache hit")
			return nil, nil
		},
	}

	likesRepo := &mockLikesRepo{
		counts: map[string]int64{"cached-post-1": 100, "cached-post-2": 0},
		liked:  map[string]bool{"cached-post-1": false, "cached-post-2": true},
	}

	cacheRepo := &mockCacheRepo{data: cacheData}

	handler := NewHandler(postRepo, &mockUserRepo{}, likesRepo, cacheRepo)

	// Default params (no cursor + limit=10)
	req := &api_proto.V1PostsListRequest{
		Limit: mainPageLimit,
	}

	resp, err := handler.PostsV1List(ctxWithUser("user-1"), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if postRepo.listCalled {
		t.Error("PostRepository.List should NOT be called on cache hit")
	}

	if !cacheRepo.getCalled {
		t.Error("expected cache to be checked for default params")
	}

	expectedPosts := []*api_proto.BlogPost{
		{
			Uuid:      "cached-post-1",
			Author:    &api_proto.User{Uuid: "user-1", Nickname: "alice", Name: "Alice", Surname: "A", AvatarUrl: "https://example.com/alice.png"},
			CreatedAt: timestamppb.New(now),
			LikesCount:       100,
			ContentText:      "Cached post 1",
			ContentImageUrls: nil,
			LikedByMe:        false,
		},
		{
			Uuid:      "cached-post-2",
			Author:    &api_proto.User{Uuid: "user-2", Nickname: "bob", Name: "Bob", Surname: "B", AvatarUrl: "https://example.com/bob.png"},
			CreatedAt: timestamppb.New(now.Add(-time.Minute)),
			LikesCount:       0,
			ContentText:      "Cached post 2",
			ContentImageUrls: nil,
			LikedByMe:        true,
		},
	}

	assertPostsEqual(t, expectedPosts, resp.Posts)

	if resp.Cursor != nil {
		t.Error("expected no cursor from cache path")
	}
}

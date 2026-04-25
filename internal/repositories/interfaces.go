package repositories

import (
	"context"
	"time"

	"github.com/Gosstik/golang_blog/internal/entities"
)

// PostRepository manages blog posts in PostgreSQL.
type PostRepository interface {
	List(ctx context.Context, limit int, cursorCreatedAt *time.Time) ([]entities.BlogPost, error)
	GetByUUID(ctx context.Context, uuid string) (*entities.BlogPost, error)
	Create(ctx context.Context, post *entities.BlogPost) error
	Update(ctx context.Context, post *entities.BlogPost) error
	Delete(ctx context.Context, uuid string) error
}

// UserRepository manages users in PostgreSQL.
type UserRepository interface {
	GetByUUID(ctx context.Context, uuid string) (*entities.User, error)
	GetByUUIDs(ctx context.Context, uuids []string) ([]entities.User, error)
}

// LikesRepository manages likes in Redis using Sets.
// Key per post: "post:{postUUID}:likes" — members are user UUIDs.
type LikesRepository interface {
	Like(ctx context.Context, postUUID, userUUID string) (int64, error)
	Unlike(ctx context.Context, postUUID, userUUID string) (int64, error)
	GetLikesCount(ctx context.Context, postUUID string) (int64, error)
	IsLikedByUser(ctx context.Context, postUUID, userUUID string) (bool, error)
	GetLikedUserUUIDs(ctx context.Context, postUUID string, limit int, cursor int64) ([]string, int64, error)
	// Batch methods for efficient main page cache enrichment (pipeline).
	GetLikesCountBatch(ctx context.Context, postUUIDs []string) (map[string]int64, error)
	IsLikedByUserBatch(ctx context.Context, postUUIDs []string, userUUID string) (map[string]bool, error)
}

// CacheRepository manages cached responses in Redis.
type CacheRepository interface {
	GetPostsListCache(ctx context.Context) ([]byte, error)
	SetPostsListCache(ctx context.Context, data []byte, ttl time.Duration) error
	InvalidatePostsListCache(ctx context.Context) error
}

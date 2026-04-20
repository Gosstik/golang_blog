package entities

import (
	"time"

	"github.com/google/uuid"
)

type PostUuid struct {
	uuid.UUID
}

func NewPostUuid(uuid uuid.UUID) PostUuid {
	return PostUuid{
		UUID: uuid,
	}
}

type BlogPost struct {
	PostUuid         PostUuid   `db:"post_uuid"`
	AuthorUuid       UserUuid   `db:"author_uuid"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        *time.Time `db:"updated_at,omitempty"`
	LikesCount       int64      `db:"likes_count"`
	ContentText      string     `db:"content_text,omitempty"`
	ContentImageUrls []string   `db:"content_image_urls,omitempty"`
}

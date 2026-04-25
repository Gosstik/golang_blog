package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/Gosstik/golang_blog/internal/entities"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) List(ctx context.Context, limit int, cursorCreatedAt *time.Time) ([]entities.BlogPost, error) {
	var posts []entities.BlogPost
	q := r.db.WithContext(ctx).
		Preload("Author").
		Order("created_at DESC").
		Limit(limit)

	if cursorCreatedAt != nil {
		q = q.Where("created_at < ?", cursorCreatedAt)
	}

	if err := q.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *PostRepository) GetByUUID(ctx context.Context, uuid string) (*entities.BlogPost, error) {
	var post entities.BlogPost
	if err := r.db.WithContext(ctx).Preload("Author").Where("uuid = ?", uuid).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) Create(ctx context.Context, post *entities.BlogPost) error {
	return r.db.WithContext(ctx).Omit("UpdatedAt").Create(post).Error
}

func (r *PostRepository) Update(ctx context.Context, post *entities.BlogPost) error {
	now := time.Now()
	post.UpdatedAt = &now
	return r.db.WithContext(ctx).
		Model(post).
		Where("uuid = ?", post.UUID).
		Updates(map[string]interface{}{
			"content_text":       post.ContentText,
			"content_image_urls": post.ContentImageUrls,
			"updated_at":         post.UpdatedAt,
		}).Error
}

func (r *PostRepository) Delete(ctx context.Context, uuid string) error {
	return r.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&entities.BlogPost{}).Error
}

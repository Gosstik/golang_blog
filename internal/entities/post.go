package entities

import (
	"time"

	"github.com/lib/pq"
)

type BlogPost struct {
	UUID             string         `gorm:"column:uuid;type:uuid;primaryKey;default:gen_random_uuid()"`
	AuthorUUID       string         `gorm:"column:author_uuid;type:uuid;not null;index"`
	Author           User           `gorm:"foreignKey:AuthorUUID;references:UUID"`
	CreatedAt        time.Time      `gorm:"column:created_at"`
	UpdatedAt        *time.Time     `gorm:"column:updated_at"`
	ContentText      string         `gorm:"column:content_text"`
	ContentImageUrls pq.StringArray `gorm:"column:content_image_urls;type:text[]"`
}

func (BlogPost) TableName() string { return "posts" }

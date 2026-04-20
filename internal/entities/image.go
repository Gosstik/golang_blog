package entities

type Image struct {
	UUID        string `gorm:"column:uuid;type:uuid;primaryKey;default:gen_random_uuid()"`
	Url         string `gorm:"column:url;not null"`
	ContentType string `gorm:"column:content_type;not null"`
}

func (Image) TableName() string { return "images" }

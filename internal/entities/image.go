package entities

import (
	"github.com/google/uuid"
)

type ImageUuid struct {
	uuid.UUID
}

func NewImageUuid(uuid uuid.UUID) ImageUuid {
	return ImageUuid{
		UUID: uuid,
	}
}

type Image struct {
	ImageUuid   ImageUuid `db:"image_uuid"`
	Url         string    `db:"url"`
	ContentType string    `db:"content_type"`
}

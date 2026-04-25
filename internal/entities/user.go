package entities

import (
	"github.com/google/uuid"
)

type UserUuid struct {
	uuid.UUID
}

func NewUserUuid(uuid uuid.UUID) UserUuid {
	return UserUuid{
		UUID: uuid,
	}
}

type User struct {
	UserUuid  UserUuid `db:"user_uuid"`
	Nickname  string   `db:"nickname"`
	Name      string   `db:"name"`
	Surname   string   `db:"surname"`
	AvatarUrl string   `db:"avatar_url,omitempty"`
}

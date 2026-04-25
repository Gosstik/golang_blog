package postgres

import (
	"context"

	"gorm.io/gorm"

	"github.com/Gosstik/golang_blog/internal/entities"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByUUID(ctx context.Context, uuid string) (*entities.User, error) {
	var user entities.User
	if err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUUIDs(ctx context.Context, uuids []string) ([]entities.User, error) {
	var users []entities.User
	if err := r.db.WithContext(ctx).Where("uuid IN ?", uuids).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

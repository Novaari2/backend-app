package users

import (
	"auth-app/internal/entity"
	"context"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) Save(ctx context.Context, user *entity.User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	return err
}

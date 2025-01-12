package users

import (
	"auth-app/internal/entity"
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Save(ctx context.Context, user *entity.User) error
	FindByNik(ctx context.Context, nik string) (*entity.User, error)
}

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

func (r *repository) FindByNik(ctx context.Context, nik string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("nik = ?", nik).First(&user).Error
	return &user, err
}

package users

import (
	"auth-app/internal/entity"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (r *MockRepository) Save(ctx context.Context, user *entity.User) error {
	args := r.Called(ctx, user)
	return args.Error(0)
}

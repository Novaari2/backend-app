package mocks

import (
	"auth-app/internal/entity"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Save(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) FindByNik(ctx context.Context, nik string) (*entity.User, error) {
	args := m.Called(ctx, nik)
	return args.Get(0).(*entity.User), args.Error(1)
}

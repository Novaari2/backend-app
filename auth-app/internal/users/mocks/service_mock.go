package mocks

import (
	"auth-app/internal/users"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) RegisterUser(ctx context.Context, req *users.RegisterUserRequest) (*users.RegisterUserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) != nil {
		return args.Get(0).(*users.RegisterUserResponse), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockService) LoginUser(ctx context.Context, req *users.LoginUserRequest) (*users.LoginUserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) != nil {
		return args.Get(0).(*users.LoginUserResponse), args.Error(1)
	}

	return nil, args.Error(1)
}

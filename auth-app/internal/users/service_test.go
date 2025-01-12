package users

import (
	"auth-app/internal/entity"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUserSuccess(t *testing.T) {
	ctx := context.Background()

	mockRepo := new(MockRepository)

	mockPasswordGenerator := func(length int) string { return "mockpass" }

	service := NewService(mockRepo, mockPasswordGenerator)

	req := &RegisterUserRequest{
		Nik:  "1234567890123456",
		Role: "admin",
	}

	expectedUser := &entity.User{
		Nik:      req.Nik,
		Role:     req.Role,
		Password: "hashedpassword",
	}

	mockRepo.On("Save", ctx, mock.MatchedBy(func(user *entity.User) bool {
		return user.Nik == expectedUser.Nik &&
			user.Role == expectedUser.Role
	})).Return(nil)

	resp, err := service.RegisterUser(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, req.Nik, resp.Nik)
	assert.Equal(t, req.Role, resp.Role)
	assert.Equal(t, "mockpass", resp.Password)

	mockRepo.AssertExpectations(t)
}

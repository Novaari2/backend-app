package users_test

import (
	"auth-app/internal/entity"
	"auth-app/internal/users"
	"auth-app/internal/users/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUser(t *testing.T) {
	ctx := context.Background()
	mockPasswordGenerator := func(length int) string { return "mockpass" }

	tests := []struct {
		name         string
		mockRepoFunc func(mockRepo *mocks.MockRepository)
		req          *users.RegisterUserRequest
		wantErr      string
		wantResp     *users.RegisterUserResponse
	}{
		{
			name: "Success",
			mockRepoFunc: func(mockRepo *mocks.MockRepository) {
				mockRepo.On("Save", ctx, mock.MatchedBy(func(user *entity.User) bool {
					return user.Nik == "1234567890123456" && user.Role == "admin"
				})).Return(nil)
			},
			req: &users.RegisterUserRequest{
				Nik:  "1234567890123456",
				Role: "admin",
			},
			wantErr: "",
			wantResp: &users.RegisterUserResponse{
				Nik:      "1234567890123456",
				Role:     "admin",
				Password: "mockpass",
			},
		},
		{
			name: "Failed to save user",
			mockRepoFunc: func(mockRepo *mocks.MockRepository) {
				mockRepo.On("Save", ctx, mock.MatchedBy(func(user *entity.User) bool {
					return user.Nik == "1234567890123456" && user.Role == "admin"
				})).Return(errors.New("Failed to save user"))
			},
			req: &users.RegisterUserRequest{
				Nik:  "1234567890123456",
				Role: "admin",
			},
			wantErr:  "Failed to save user",
			wantResp: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockRepository)
			tt.mockRepoFunc(mockRepo)

			service := users.NewService(mockRepo, mockPasswordGenerator, nil)

			resp, err := service.RegisterUser(ctx, tt.req)

			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResp, resp)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestLoginUser(t *testing.T) {
	ctx := context.Background()
	const staticHash = "$2a$14$kcuProf0TwyYPC9YXua8JeYAAXN.UFg1ba8CbumfPagrBkY37xeUC"

	tests := []struct {
		name          string
		mockRepoFunc  func(mockRepo *mocks.MockRepository)
		mockTokenFunc func(mockTokenGen *mocks.MockTokenGenerator)
		req           *users.LoginUserRequest
		expectedResp  *users.LoginUserResponse
		expectedError string
	}{
		{
			name: "Login success",
			mockRepoFunc: func(mockRepo *mocks.MockRepository) {
				mockRepo.On("FindByNik", ctx, "1234567890123456").Return(&entity.User{
					ID:       1,
					Nik:      "1234567890123456",
					Role:     "admin",
					Password: staticHash,
				}, nil)
			},
			mockTokenFunc: func(mockTokenGen *mocks.MockTokenGenerator) {
				mockTokenGen.On("GenerateJWT", "1234567890123456", staticHash).Return("mocked.jwt.token", nil)
			},
			req: &users.LoginUserRequest{
				Nik:      "1234567890123456",
				Password: "correctpassword",
			},
			expectedResp: &users.LoginUserResponse{
				ID:    1,
				Nik:   "1234567890123456",
				Role:  "admin",
				Token: "mocked.jwt.token",
			},
			expectedError: "",
		},
		{
			name: "Invalid password",
			mockRepoFunc: func(mockRepo *mocks.MockRepository) {
				mockRepo.On("FindByNik", ctx, "1234567890123456").Return(&entity.User{
					ID:       1,
					Nik:      "1234567890123456",
					Role:     "admin",
					Password: staticHash,
				}, nil)
			},
			mockTokenFunc: func(mockTokenGen *mocks.MockTokenGenerator) {
				mockTokenGen.On("GenerateJWT", mock.Anything, mock.Anything).Maybe()
			},
			req: &users.LoginUserRequest{
				Nik:      "1234567890123456",
				Password: "wrongpassword",
			},
			expectedResp:  nil,
			expectedError: "invalid password",
		},
		{
			name: "User not found",
			mockRepoFunc: func(mockRepo *mocks.MockRepository) {
				mockRepo.On("FindByNik", ctx, "notexistnik").Return(&entity.User{}, errors.New("user not found"))
			},
			mockTokenFunc: func(mockTokenGen *mocks.MockTokenGenerator) {
				mockTokenGen.On("GenerateJWT", mock.Anything, mock.Anything).Maybe()
			},
			req: &users.LoginUserRequest{
				Nik:      "notexistnik",
				Password: "password",
			},
			expectedResp:  nil,
			expectedError: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockRepository)
			mockTokenGen := new(mocks.MockTokenGenerator)

			tt.mockRepoFunc(mockRepo)
			tt.mockTokenFunc(mockTokenGen)

			service := users.NewService(mockRepo, nil, mockTokenGen)

			resp, err := service.LoginUser(ctx, tt.req)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp.ID, resp.ID)
				assert.Equal(t, tt.expectedResp.Nik, resp.Nik)
				assert.Equal(t, tt.expectedResp.Role, resp.Role)
				assert.Equal(t, tt.expectedResp.Token, resp.Token)
			}

			mockRepo.AssertExpectations(t)
			mockTokenGen.AssertExpectations(t)
		})
	}
}

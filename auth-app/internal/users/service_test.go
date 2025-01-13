package users_test

import (
	"auth-app/internal/entity"
	"auth-app/internal/users"
	"auth-app/internal/users/mocks"
	"auth-app/internal/utils"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	staticHash  = "$2a$14$kcuProf0TwyYPC9YXua8JeYAAXN.UFg1ba8CbumfPagrBkY37xeUC"
	staticToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuaWsiOiIzMzc0MDcyMDEyOTgwMDA1IiwicGFzc3dvcmQiOiJtb2NrcGFzcyIsImV4cCI6MTczNjgzNTY4MiwiaWF0IjoxNzM2NzQ5MjgyfQ.S0bYsFRE0hexrSLqyw5FevRQnioRB0XT8-LKliL07-c"
	staticPass  = "$2a$14$VE311U4UriU6bcASrNRi/.E3LAJywZRfWXt.mv9Ysnfrz2dlDgnSe"
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

func TestClaimToken(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		mockBehavior  func(mockTokenGen *mocks.MockTokenGenerator)
		req           string
		expectedResp  *users.PayloadResponse
		expectedError string
	}{
		{
			name: "Claim success",
			mockBehavior: func(mockTokenGen *mocks.MockTokenGenerator) {
				mockTokenGen.On("ValidateJwt", staticToken).Return(&utils.JWTClaims{
					Nik:      "3374072012980005",
					Password: staticPass,
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Unix(1736835682, 0)),
					},
				}, nil)
			},
			req: staticToken,
			expectedResp: &users.PayloadResponse{
				Nik:      "3374072012980005",
				Password: staticPass,
				Exp:      int64(1736835682),
			},
			expectedError: "",
		},
		{
			name: "Claim failed",
			mockBehavior: func(mockTokenGen *mocks.MockTokenGenerator) {
				mockTokenGen.On("ValidateJwt", "mocked.jwt.token").Return(&utils.JWTClaims{}, errors.New("invalid token"))
			},
			req:           "mocked.jwt.token",
			expectedResp:  nil,
			expectedError: "invalid token",
		},
		{
			name: "Empty token",
			mockBehavior: func(mockTokenGen *mocks.MockTokenGenerator) {
				mockTokenGen.On("ValidateJwt", "").Return(&utils.JWTClaims{}, errors.New("empty token"))
			},
			req:           "",
			expectedResp:  nil,
			expectedError: "empty token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTokenGen := new(mocks.MockTokenGenerator)

			tt.mockBehavior(mockTokenGen)

			service := users.NewService(nil, nil, mockTokenGen)

			resp, err := service.ValidateToken(ctx, tt.req)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp.Nik, resp.Nik)
				assert.Equal(t, tt.expectedResp.Password, resp.Password)
				assert.Equal(t, tt.expectedResp.Exp, resp.Exp)
			}

			mockTokenGen.AssertExpectations(t)
		})
	}
}

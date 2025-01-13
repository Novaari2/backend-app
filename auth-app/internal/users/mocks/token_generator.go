package mocks

import (
	"auth-app/internal/utils"

	"github.com/stretchr/testify/mock"
)

type MockTokenGenerator struct {
	mock.Mock
}

func (m *MockTokenGenerator) GenerateJWT(nik string, password string) (string, error) {
	args := m.Called(nik, password)
	return args.String(0), args.Error(1)
}

func (m *MockTokenGenerator) ValidateJwt(tokenString string) (*utils.JWTClaims, error) {
	args := m.Called(tokenString)
	if claims, ok := args.Get(0).(*utils.JWTClaims); ok {
		return claims, args.Error(1)
	}

	return nil, args.Error(1)
}

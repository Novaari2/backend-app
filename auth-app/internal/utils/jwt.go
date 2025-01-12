package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Nik      string `json:"nik"`
	Password string `json:"password"`
	jwt.RegisteredClaims
}

type TokenGenerator interface {
	GenerateJWT(nik, password string) (string, error)
	// ValidateJwt(tokenString string) (*JWTClaims, error)
}

type DefaultTokenGenerator struct {
	secretKey string
}

func NewTokenGenerator(secretKey string) TokenGenerator {
	return &DefaultTokenGenerator{
		secretKey: secretKey,
	}
}

func (t *DefaultTokenGenerator) GenerateJWT(nik string, password string) (string, error) {
	claims := JWTClaims{
		Nik:      nik,
		Password: password,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(t.secretKey))
}

// func (t *DefaultTokenGenerator) ValidateJwt(tokenString string) (*JWTClaims, error) {
// 	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(t.secretKey), nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	claims, ok := token.Claims.(*JWTClaims)
// 	if !ok || !token.Valid {
// 		return nil, fmt.Errorf("invalid token")
// 	}

// 	return claims, nil
// }

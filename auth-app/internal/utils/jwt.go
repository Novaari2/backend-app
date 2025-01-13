package utils

import (
	"errors"
	"fmt"
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
	ValidateJwt(tokenString string) (*JWTClaims, error)
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

func (t *DefaultTokenGenerator) ValidateJwt(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(t.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		if claims.ExpiresAt.Before(time.Now()) {
			return nil, errors.New("token has expired")
		}
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}

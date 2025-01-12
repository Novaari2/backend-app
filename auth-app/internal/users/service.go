package users

import (
	"auth-app/internal/entity"
	"auth-app/internal/utils"
	"context"

	"github.com/rs/zerolog/log"
)

type PasswordGenerator func(int) string

type Repository interface {
	Save(ctx context.Context, user *entity.User) error
}

type Service struct {
	repo             Repository
	generatePassword PasswordGenerator
}

func NewService(repo Repository, generatePassword PasswordGenerator) *Service {
	return &Service{
		repo:             repo,
		generatePassword: generatePassword,
	}
}

func (s *Service) RegisterUser(ctx context.Context, req *RegisterUserRequest) (RegisterUserResponse, error) {
	pass := s.generatePassword(6)

	hashPass, err := utils.HashPassword(pass)
	if err != nil {
		log.Err(err).Msg("Failed to hash password")
		return RegisterUserResponse{}, err
	}

	user := &entity.User{
		Nik:      req.Nik,
		Role:     req.Role,
		Password: hashPass,
	}

	err = s.repo.Save(ctx, user)
	if err != nil {
		log.Err(err).Msg("Failed to save user")
		return RegisterUserResponse{}, err
	}

	return RegisterUserResponse{
		Nik:      req.Nik,
		Role:     req.Role,
		Password: pass,
	}, nil
}

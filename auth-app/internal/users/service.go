package users

import (
	"auth-app/internal/entity"
	"auth-app/internal/utils"
	"context"

	"github.com/rs/zerolog/log"
)

type Service interface {
	RegisterUser(ctx context.Context, req *RegisterUserRequest) (*RegisterUserResponse, error)
	LoginUser(ctx context.Context, req *LoginUserRequest) (*LoginUserResponse, error)
}

type PasswordGenerator func(int) string

type service struct {
	repo             Repository
	generatePassword PasswordGenerator
	tokenGenerator   utils.TokenGenerator
}

func NewService(repo Repository, generatePassword PasswordGenerator, tokenGenerator utils.TokenGenerator) *service {
	return &service{
		repo:             repo,
		generatePassword: generatePassword,
		tokenGenerator:   tokenGenerator,
	}
}

func (s *service) RegisterUser(ctx context.Context, req *RegisterUserRequest) (*RegisterUserResponse, error) {
	pass := s.generatePassword(6)

	hashPass, err := utils.HashPassword(pass)
	if err != nil {
		log.Err(err).Msg("Failed to hash password")
		return &RegisterUserResponse{}, err
	}

	user := &entity.User{
		Nik:      req.Nik,
		Role:     req.Role,
		Password: hashPass,
	}

	err = s.repo.Save(ctx, user)
	if err != nil {
		log.Err(err).Msg("Failed to save user")
		return &RegisterUserResponse{}, err
	}

	return &RegisterUserResponse{
		Nik:      req.Nik,
		Role:     req.Role,
		Password: pass,
	}, nil
}

func (s *service) LoginUser(ctx context.Context, req *LoginUserRequest) (*LoginUserResponse, error) {
	user, err := s.repo.FindByNik(ctx, req.Nik)
	if err != nil {
		log.Err(err).Msg("Failed to find user")
		return &LoginUserResponse{}, err
	}

	err = utils.CheckPasswordHash(req.Password, user.Password)
	if err != nil {
		log.Err(err).Msg("Password does not match")
		return &LoginUserResponse{}, err
	}

	token, err := s.tokenGenerator.GenerateJWT(user.Nik, user.Password)
	if err != nil {
		log.Err(err).Msg("Failed to generate token")
		return &LoginUserResponse{}, err
	}

	return &LoginUserResponse{
		ID:    user.ID,
		Nik:   req.Nik,
		Role:  user.Role,
		Token: token,
	}, nil
}

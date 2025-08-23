package services

import (
	"errors"
	"strings"

	"github.com/vinhnglx/go-rest-api-template/internal/models"
	"github.com/vinhnglx/go-rest-api-template/internal/repositories"
	"github.com/vinhnglx/go-rest-api-template/internal/util"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(req models.RegisterRequest) (*models.UserResponse, error) {
	if s.userRepo.EmailExists(req.Email) {
		return nil, errors.New("email already exists")
	}

	// TODO: Can add the password strength validation here later

	hashedPassword, err := util.HashPassword(req.Password)

	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &models.User{
		Name:         strings.TrimSpace(req.Name),
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash: hashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &models.UserResponse{
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *AuthService) Login(req models.LoginRequest) (*models.AuthResponse, error) {
	emailExists := s.userRepo.EmailExists(req.Email)

	if !emailExists {
		return nil, errors.New("email does not exist")
	}

	user, err := s.userRepo.GetByEmail(strings.ToLower(strings.TrimSpace(req.Email)))

	if err != nil {
		return nil, errors.New("failed to retrieve user")
	}

	if !util.CheckPassword(user.PasswordHash, req.Password) {
		return nil, errors.New("failed to retrieve user")
	}

	token, err := util.GenerateJWT(user.ID, user.Email)

	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &models.AuthResponse{
		Token: token,
	}, nil
}

package service

import (
	"context"
	"errors"

	"github.com/hamwiwatsapon/go-ticket-booking/internal/domain"
	"github.com/hamwiwatsapon/go-ticket-booking/internal/dto"
	"github.com/hamwiwatsapon/go-ticket-booking/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, createDTO *dto.CreateUserDTO) (*dto.UserResponseDTO, error) {
	// Check if user already exists
	existingUser, err := s.repo.FindByEmail(ctx, createDTO.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user domain model
	user := &domain.User{
		Email:     createDTO.Email,
		Password:  string(hashedPassword),
		FirstName: createDTO.FirstName,
		LastName:  createDTO.LastName,
		Role:      createDTO.Role,
	}

	// Save to repository
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Convert to response DTO
	return &dto.UserResponseDTO{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *UserService) AuthenticateUser(ctx context.Context, email, password string) (*dto.UserResponseDTO, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &dto.UserResponseDTO{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

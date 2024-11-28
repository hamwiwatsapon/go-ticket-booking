package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/hamwiwatsapon/go-ticket-booking/internal/domain"
	"github.com/hamwiwatsapon/go-ticket-booking/internal/dto"
	"github.com/hamwiwatsapon/go-ticket-booking/internal/repository"
	"github.com/hamwiwatsapon/go-ticket-booking/pkg/utils"
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

func (s *UserService) Login(ctx context.Context, email, password string) (*dto.AuthResponse, error) {
	// Find user by email
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT tokens
	tokens, err := utils.GenerateTokenPair(fmt.Sprint(user.ID), user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	// Prepare response
	return &dto.AuthResponse{
		User: dto.UserResponseDTO{
			ID:        fmt.Sprint(user.ID),
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
		},
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *UserService) RefreshAccessToken(ctx context.Context, refreshToken string) (*dto.AuthResponse, error) {
	// Validate and generate new tokens
	newTokens, err := utils.RefreshAccessToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Optionally fetch user details if needed
	// In this example, we're not fetching full user details to keep it simple
	return &dto.AuthResponse{
		AccessToken:  newTokens.AccessToken,
		RefreshToken: newTokens.RefreshToken,
	}, nil
}

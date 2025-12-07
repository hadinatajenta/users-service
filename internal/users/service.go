package users

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateUser(ctx context.Context, input CreateUserInput) (UserResponse, error)
	GetUser(ctx context.Context, id int64) (UserResponse, error)
	ListUsers(ctx context.Context) ([]UserResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateUser(ctx context.Context, input CreateUserInput) (UserResponse, error) {
	if err := validateCreateInput(input); err != nil {
		return UserResponse{}, err
	}

	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		return UserResponse{}, fmt.Errorf("hash password: %w", err)
	}

	user := User{
		Name:     strings.TrimSpace(input.Name),
		Email:    strings.TrimSpace(input.Email),
		Password: hashedPassword,
	}

	created, err := s.repo.Create(ctx, user)
	if err != nil {
		return UserResponse{}, err
	}

	return toOutput(created), nil
}

func (s *service) GetUser(ctx context.Context, id int64) (UserResponse, error) {
	if id <= 0 {
		return UserResponse{}, fmt.Errorf("user id must be positive")
	}
	found, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return UserResponse{}, err
	}
	return toOutput(found), nil
}

func (s *service) ListUsers(ctx context.Context) ([]UserResponse, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]UserResponse, 0, len(users))
	for _, u := range users {
		out = append(out, toOutput(u))
	}
	return out, nil
}

func validateCreateInput(input CreateUserInput) error {
	if strings.TrimSpace(input.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if strings.TrimSpace(input.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if strings.TrimSpace(input.Password) == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func toOutput(u User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

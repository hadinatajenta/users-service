package service

import (
	"context"
	"fmt"
	"strings"

	"users-service/internal/entity"
	"users-service/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

type CreateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s *UserService) CreateUser(ctx context.Context, input CreateUserInput) (entity.User, error) {
	if err := validateCreateInput(input); err != nil {
		return entity.User{}, err
	}

	user := entity.User{
		Name:  strings.TrimSpace(input.Name),
		Email: strings.TrimSpace(input.Email),
	}

	return s.repo.Create(ctx, user)
}

func (s *UserService) GetUser(ctx context.Context, id int64) (entity.User, error) {
	if id <= 0 {
		return entity.User{}, fmt.Errorf("user id must be positive")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context) ([]entity.User, error) {
	return s.repo.List(ctx)
}

func validateCreateInput(input CreateUserInput) error {
	if strings.TrimSpace(input.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if strings.TrimSpace(input.Email) == "" {
		return fmt.Errorf("email is required")
	}
	return nil
}

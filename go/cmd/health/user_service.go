package main

import (
	"context"
	"errors"
	"strings"
)

var ErrNotFound = errors.New("user not found")

type UserRepository interface {
	Create(ctx context.Context, name, email string) (User, error)
	List(ctx context.Context) ([]User, error)
	Get(ctx context.Context, id int) (User, error)
	Update(ctx context.Context, id int, name, email string) (User, error)
	Delete(ctx context.Context, id int) error
}

type UserService struct {
	repository UserRepository
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewUserService(repository UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func validateCreateUser(name, email string) error {
	if strings.TrimSpace(name) == "" {
		return &ValidationError{
			Message: "name is required",
		}
	}

	if strings.TrimSpace(email) == "" {
		return &ValidationError{
			Message: "email is required",
		}
	}

	parts := strings.Split(email, "@")

	if len(parts) != 2 ||
		strings.TrimSpace(parts[0]) == "" ||
		strings.TrimSpace(parts[1]) == "" {
		return &ValidationError{
			Message: "invalid email",
		}
	}

	return nil
}

func (s *UserService) Create(ctx context.Context, name, email string) (User, error) {
	if err := validateCreateUser(name, email); err != nil {
		return User{}, err
	}

	return s.repository.Create(ctx, name, email)
}

func (s *UserService) List(ctx context.Context) ([]User, error) {
	return s.repository.List(ctx)
}

func (s *UserService) Get(ctx context.Context, id int) (User, error) {
	return s.repository.Get(ctx, id)
}

func (s *UserService) Update(ctx context.Context, id int, name, email string) (User, error) {
	return s.repository.Update(ctx, id, name, email)
}

func (s *UserService) Delete(ctx context.Context, id int) error {
	return s.repository.Delete(ctx, id)
}

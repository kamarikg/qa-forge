package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockUserRepository struct {
	createCalled bool
	listCalled   bool
	getCalled    bool
	updateCalled bool
	deleteCalled bool

	user User
}

func (m *mockUserRepository) Create(
	ctx context.Context,
	name, email string,
) (User, error) {
	m.createCalled = true

	m.user = User{
		ID:    1,
		Name:  name,
		Email: email,
	}

	return m.user, nil
}

func (m *mockUserRepository) List(
	ctx context.Context,
) ([]User, error) {
	m.listCalled = true

	return []User{m.user}, nil
}

func (m *mockUserRepository) Get(
	ctx context.Context,
	id int,
) (User, error) {
	m.getCalled = true

	if id != m.user.ID {
		return User{}, ErrNotFound
	}

	return m.user, nil
}

func (m *mockUserRepository) Update(
	ctx context.Context,
	id int,
	name string,
	email string,
) (User, error) {
	m.updateCalled = true

	if id != m.user.ID {
		return User{}, ErrNotFound
	}

	m.user = User{
		ID:    id,
		Name:  name,
		Email: email,
	}

	return m.user, nil
}

func (m *mockUserRepository) Delete(
	ctx context.Context,
	id int,
) error {
	m.deleteCalled = true

	if id != m.user.ID {
		return ErrNotFound
	}

	m.user = User{}

	return nil
}

func TestUserServiceCreate(t *testing.T) {
	repository := &mockUserRepository{}
	service := NewUserService(repository)

	user, err := service.Create(
		context.Background(),
		"Alice",
		"alice@example.com",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !repository.createCalled {
		t.Fatal("expected Create to be called")
	}

	if user.Name != "Alice" {
		t.Fatalf(
			"expected name %q, got %q",
			"Alice",
			user.Email,
		)
	}

	if user.Email != "alice@example.com" {
		t.Fatalf(
			"expected email %q, got %q",
			"alice@example.com",
			user.Email,
		)
	}
}

func TestUserServiceGet(t *testing.T) {
	repository := &mockUserRepository{
		user: User{
			ID:    1,
			Name:  "Alice",
			Email: "alice@ecample.com",
		},
	}

	service := NewUserService(repository)

	user, err := service.Get(
		context.Background(),
		1,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !repository.getCalled {
		t.Fatal("expected Get to be called")
	}

	if user.Name != "Alice" {
		t.Fatalf(
			"expected name %q, got %q",
			"Alice",
			user.Name,
		)
	}
}

func TestUserServiceDelete(t *testing.T) {
	repository := &mockUserRepository{
		user: User{
			ID:    1,
			Name:  "Alice",
			Email: "alice@example.com",
		},
	}

	service := NewUserService(repository)

	err := service.Delete(
		context.Background(),
		1,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !repository.deleteCalled {
		t.Fatal("expected Delete to be called")
	}
}

func TestUserServiceCreateRejectsEmptyEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
	}{
		{
			name:  "empty email",
			email: "",
		},
		{
			name:  "whitespace email",
			email: "   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &mockUserRepository{}
			service := NewUserService(repository)

			_, err := service.Create(
				context.Background(),
				"Alice",
				tt.email,
			)

			if err == nil {
				t.Fatalf("expected error for invalid email")
			}

			if repository.createCalled {
				t.Fatalf("repository Create must not be called")
			}
		})
	}
}

func TestUserServiceCreateRejectsInvalidName(t *testing.T) {
	tests := []struct {
		name      string
		inputName string
	}{
		{
			name:      "empty name",
			inputName: "",
		},
		{
			name:      "whitespace name",
			inputName: "   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &mockUserRepository{}
			service := NewUserService(repository)

			_, err := service.Create(
				context.Background(),
				tt.inputName,
				"alice@example.com",
			)

			if err == nil {
				t.Fatalf("expected error for invalid name")
			}

			if repository.createCalled {
				t.Fatalf("repository Create must not be called")
			}
		})
	}
}

func TestUserServiceCreateRetunrsValidationError(t *testing.T) {
	repository := &mockUserRepository{}
	service := NewUserService(repository)

	_, err := service.Create(
		context.Background(),
		"Alice",
		"",
	)

	if err == nil {
		t.Fatalf("expected ValidationError, got %T", err)
	}

	var validationErr *ValidationError

	if !errors.As(err, &validationErr) {
		t.Fatalf(
			"ecxpected ValidationError, got %T",
			err,
		)
	}

	if validationErr.Message != "email is required" {
		t.Fatalf(
			"expected message %q, got %q",
			"email is reauired",
			validationErr.Message,
		)
	}

	if repository.createCalled {
		t.Fatalf("repository Create must not be called")
	}
}

func TestCreateUserHandlerValidationError(t *testing.T) {
	repository := &mockUserRepository{}
	service := NewUserService(repository)
	handler := createUserHandler(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/users/",
		strings.NewReader(`{
			"name": "Alice",
			"email": "   "
		}`),
	)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}

	if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf(
			"expected Content-Type application/json, got %q",
			contentType,
		)
	}

	var response ErrorResponse

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf(
			"failed to decode error responce: %v",
			err,
		)
	}

	if response.Error != "email is required" {
		t.Fatalf(
			"expected error %q, got %q",
			"email is required",
			response.Error,
		)
	}

	if repository.createCalled {
		t.Fatalf("repository Create must not be called")
	}
}

func TestUserServiceCreateRejectsInvalidEmail(t *testing.T) {
	tests := []struct {
		name       string
		inputEmail string
	}{
		{
			name:       "missing at sign",
			inputEmail: "alice.example.com",
		},
		{
			name:       "missing domain",
			inputEmail: "alice@",
		},
		{
			name:       "missing local part",
			inputEmail: "@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &mockUserRepository{}
			service := NewUserService(repository)

			_, err := service.Create(
				context.Background(),
				"Alice",
				tt.inputEmail,
			)

			if err == nil {
				t.Fatalf("expected error for invalid email")
			}

			if repository.createCalled {
				t.Fatalf("repository Create must not be called")
			}
		})
	}
}

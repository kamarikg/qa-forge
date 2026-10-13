package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		wantStatus int
	}{
		{
			name:       "GET returns 200",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
		},
		{
			name:       "POST returns 405",
			method:     http.MethodPost,
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "PUT returns 405",
			method:     http.MethodPut,
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "DELETE returns 405",
			method:     http.MethodDelete,
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				tt.method,
				"/health",
				nil,
			)

			rec := httptest.NewRecorder()

			healthHandler(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf(
					"expected status %d, got %d",
					tt.wantStatus,
					rec.Code,
				)
			}
		})
	}
}

func TestVersionHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()

	versionHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	expectedContentType := "application/json"
	actualContentType := rec.Header().Get("Content-Type")

	if actualContentType != expectedContentType {
		t.Fatalf(
			"expected Content-Type %q, got %q",
			expectedContentType,
			actualContentType,
		)
	}

	expectedBody := "{\"version\":\"0.1.0\"}\n"

	if rec.Body.String() != expectedBody {
		t.Fatalf(
			"expected body %q, got %q",
			expectedBody,
			rec.Body.String(),
		)
	}
}

func TestCreateUserHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "valid user",
			body:       `{"name":"Alice","email":"alice@example.com"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "missing name",
			body:       `{"email":"alice@example.com"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing email",
			body:       `{"name":"Alice"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty name",
			body:       `{"name":"","email":"alice@example.com"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty email",
			body:       `{"name":"Alice","email":""}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			body:       `{"name":`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty body",
			body:       "",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewUserStore()
			service := NewUserService(store)

			req := httptest.NewRequest(
				http.MethodPost,
				"/users",
				strings.NewReader(tt.body),
			)

			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			createUserHandler(service)(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf(
					"expected status %d, got %d",
					tt.wantStatus,
					rec.Code,
				)
			}
		})
	}
}

func TestListUsersHandler(t *testing.T) {
	store := NewUserStore()
	service := NewUserService(store)

	store.users[1] = User{
		ID:    1,
		Name:  "Kamarik",
		Email: "kamarik@example.com",
	}

	store.users[2] = User{
		ID:    2,
		Name:  "Alex",
		Email: "alex@example.com",
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/users",
		nil,
	)

	rec := httptest.NewRecorder()

	listUsersHandler(service)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expexted status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	expectadContentType := "application/json"
	actualContentType := rec.Header().Get("Content-Type")

	if actualContentType != expectadContentType {
		t.Fatalf(
			"expected Content-Type %q, got %q",
			expectadContentType,
			actualContentType,
		)
	}

	var users []User

	if err := json.NewDecoder(rec.Body).Decode(&users); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(users) != 2 {
		t.Fatalf(
			"expected 2 users, got %d",
			len(users),
		)
	}
}

func TestGetUserHandler(t *testing.T) {
	store := NewUserStore()
	service := NewUserService(store)

	createdUser := User{
		ID:    store.nextUserID,
		Name:  "Kamarik",
		Email: "kamarik@example.com",
	}

	store.users[createdUser.ID] = createdUser
	store.nextUserID++

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	rec := httptest.NewRecorder()

	getUserHandler(service)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf(
			"expacted Content-Type application/json, got %q",
			contentType,
		)
	}

	var user User

	if err := json.NewDecoder(rec.Body).Decode(&user); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if user.ID != createdUser.ID {
		t.Fatalf(
			"expacted user Id %d, got %d",
			createdUser.ID,
			user.ID,
		)
	}

	if user.Name != createdUser.Name {
		t.Fatalf(
			"expected name %q, got %q",
			createdUser.Name,
			user.Name,
		)
	}

	if user.Email != createdUser.Email {
		t.Fatalf(
			"expected email %q, got %q",
			createdUser.Email,
			user.Email,
		)
	}

}

func TestUpdateUserHandler(t *testing.T) {
	store := NewUserStore()
	service := NewUserService(store)

	store.users[1] = User{
		ID:    1,
		Name:  "Kamarik",
		Email: "kamarik@example.com",
	}

	requestBody := `{
		"name":  "Updated Kamarik",
		"email": "updated@example.com"
	}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/users/1",
		strings.NewReader(requestBody),
	)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	updateUserHandler(service)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expacted status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Type") != "application/json" {
		t.Fatalf(
			"expacted Content-Type application/json, got %q",
			rec.Header().Get("Content-Type"),
		)
	}

	var user User

	if err := json.NewDecoder(rec.Body).Decode(&user); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if user.ID != 1 {
		t.Fatalf("expected ID 1, got %d", user.ID)
	}

	if user.Name != "Updated Kamarik" {
		t.Fatalf("expected update name, got %q", user.Name)
	}

	if user.Email != "updated@example.com" {
		t.Fatalf("expected updated email, got %q", user.Email)
	}
}

func TestDeleteUserHandler(t *testing.T) {
	store := NewUserStore()
	service := NewUserService(store)

	store.users[1] = User{
		ID:    1,
		Name:  "Kamarik",
		Email: "kamarik@example.com",
	}

	req := httptest.NewRequest(
		http.MethodDelete,
		"/users/1",
		nil,
	)

	rec := httptest.NewRecorder()

	deleteUserHandler(service)(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNoContent,
			rec.Code,
		)
	}

	if _, ok := store.users[1]; ok {
		t.Fatalf("expected user to be deleted")
	}
}

package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
	})
}

func createUserHandler(service *UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request CreateUserRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			writeJSONError(
				w,
				http.StatusBadRequest,
				"invalid JSON",
			)
			return
		}

		user, err := service.Create(
			r.Context(),
			request.Name,
			request.Email,
		)
		if err != nil {
			writeJSONError(
				w,
				http.StatusBadRequest,
				err.Error(),
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_ = json.NewEncoder(w).Encode(user)
	}
}

func listUsersHandler(service *UserService) http.HandlerFunc {
	return listUserHandler(service)
}

func listUserHandler(service *UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		users, err := service.List(r.Context())
		if err != nil {
			writeJSONError(
				w,
				http.StatusInternalServerError,
				"internal server error",
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_ = json.NewEncoder(w).Encode(users)
	}
}

func getUserHandler(service *UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "merhod not allowed", http.StatusMethodNotAllowed)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/users/")

		userID, err := strconv.Atoi(path)
		if err != nil || userID <= 0 {
			writeJSONError(
				w,
				http.StatusBadRequest,
				"invalid user id",
			)
			return
		}

		user, err := service.Get(r.Context(), userID)
		if err != nil {
			writeJSONError(
				w,
				http.StatusNotFound,
				"user not found",
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_ = json.NewEncoder(w).Encode(user)
	}
}

func updateUserHandler(service *UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/users/")

		userID, err := strconv.Atoi(path)
		if err != nil || userID <= 0 {
			writeJSONError(
				w,
				http.StatusBadRequest,
				"invalid user id",
			)
			return
		}

		var request CreateUserRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			writeJSONError(
				w,
				http.StatusBadRequest,
				"invalid request body",
			)
			return
		}

		request.Name = strings.TrimSpace(request.Name)
		request.Email = strings.TrimSpace(request.Email)

		if request.Name == "" {
			writeJSONError(
				w,
				http.StatusBadRequest,
				"name is required",
			)
			return
		}

		if request.Email == "" {
			writeJSONError(
				w,
				http.StatusBadRequest,
				"email is required",
			)
			return
		}

		user, err := service.Update(
			r.Context(),
			userID,
			request.Name,
			request.Email,
		)
		if err != nil {
			writeJSONError(
				w,
				http.StatusNotFound,
				"user not found",
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_ = json.NewEncoder(w).Encode(user)
	}
}

func deleteUserHandler(service *UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/users/")

		userID, err := strconv.Atoi(path)
		if err != nil || userID <= 0 {
			writeJSONError(
				w,
				http.StatusBadRequest,
				"invalid user id",
			)
			return
		}

		if err := service.Delete(r.Context(), userID); err != nil {
			writeJSONError(
				w,
				http.StatusNotFound,
				"user not found",
			)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

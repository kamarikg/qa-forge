package main

import (
	"context"
	"errors"
	"testing"
)

func TestUserStoreCreate(t *testing.T) {
	store := NewUserStore()
	ctx := context.Background()

	user1, err := store.Create(
		ctx,
		"Alice",
		"alice@example.com",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	user2, err := store.Create(
		ctx,
		"Bob",
		"bob@example.com",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user1.ID != 1 {
		t.Fatalf("expected firts ID 1, got %d", user1.ID)
	}

	if user2.ID != 2 {
		t.Fatalf("expected second ID 2, got %d", user2.ID)
	}

	if user1.Name != "Alice" {
		t.Fatalf("expected name Alice, got %q", user1.Name)
	}

	if user2.Email != "bob@example.com" {
		t.Fatalf(
			"expected email bob@example.com, got %q",
			user2.Email,
		)
	}
}

func TestUserStoreGet(t *testing.T) {
	store := NewUserStore()
	ctx := context.Background()

	created, err := store.Create(
		ctx,
		"Alice",
		"alice@ecample.com",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		name     string
		id       int
		wantErr  error
		wantUser User
	}{
		{
			name:     "existing user",
			id:       created.ID,
			wantErr:  nil,
			wantUser: created,
		},
		{
			name:    "missing user",
			id:      999,
			wantErr: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := store.Get(ctx, tt.id)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"expected error %v, got %v",
					tt.wantErr,
					err,
				)
			}

			if tt.wantErr == nil && user != tt.wantUser {
				t.Fatalf(
					"expected user %+v, got %+v",
					tt.wantUser,
					user,
				)
			}
		})
	}
}

func TestUserStoreList(t *testing.T) {
	store := NewUserStore()
	ctx := context.Background()

	_, err := store.Create(
		ctx,
		"Alice",
		"alice@example.com",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = store.Create(
		ctx,
		"Bob",
		"bob@example.com",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	users, err := store.List(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
}

func TestUserStoreUpdate(t *testing.T) {
	store := NewUserStore()
	ctx := context.Background()

	created, err := store.Create(
		ctx,
		"Alice",
		"alice@example.com",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, err := store.Update(
		ctx,
		created.ID,
		"Updated Alice",
		"updated@example.com",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updated.ID != created.ID {
		t.Fatalf(
			"expected ID %d, got %d",
			created.ID,
			updated.ID,
		)
	}

	if updated.Name != "Updated Alice" {
		t.Fatalf(
			"expected name %q, got %q",
			"Updated Alice",
			updated.Name,
		)
	}

	if updated.Email != "updated@example.com" {
		t.Fatalf(
			"expected email %q, got %q",
			"updated@example.com",
			updated.Email,
		)
	}

	_, err = store.Update(
		ctx,
		999,
		"Nobody",
		"nobody@example.com",
	)

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf(
			"expected ErrNotFound, got %v",
			err,
		)
	}
}

func TestUserStoreDelete(t *testing.T) {
	store := NewUserStore()
	ctx := context.Background()

	created, err := store.Create(
		ctx,
		"Alice",
		"alice@example.com",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := store.Delete(ctx, created.ID); err != nil {
		t.Fatalf("expected delete to succeed, got %v", err)
	}

	if _, err := store.Get(ctx, created.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf(
			"expected deleted user to return ErrNotFound, got %v",
			err,
		)
	}

	if err := store.Delete(ctx, 999); !errors.Is(err, ErrNotFound) {
		t.Fatalf(
			"expected delete of missing user to return ErrNotFound, got %v",
			err,
		)
	}
}

package main

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repositoryFactory func(t *testing.T) UserRepository

func registerUserCleanup(t *testing.T, repository UserRepository, id int) {
	t.Helper()

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := repository.Delete(ctx, id)
		if err != nil && !errors.Is(err, ErrNotFound) {
			t.Errorf("cleanup Delete(%d) error = %v", id, err)
		}
	})
}

func TestUserRepositoryContract(t *testing.T) {
	factories := map[string]repositoryFactory{
		"UserStore": func(t *testing.T) UserRepository {
			return NewUserStore()
		},

		"PostgreSQLRepository": func(t *testing.T) UserRepository {
			dsn := os.Getenv("QA_FORGE_TEST_DATABASE_URL")
			if dsn == "" {
				t.Fatal("QA_FORGE_TEST_DATABASE_URL must be set for PostgreSQL contract tests")
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			pool, err := pgxpool.New(ctx, dsn)
			if err != nil {
				t.Fatalf("connect to PostgreSQL: %v", err)
			}

			if err := pool.Ping(ctx); err != nil {
				pool.Close()
				t.Fatalf("ping PostgreSQL: %v", err)
			}

			t.Cleanup(pool.Close)

			return NewPostgreSQLRepository(pool)
		},
	}

	for name, factory := range factories {
		t.Run(name, func(t *testing.T) {
			t.Run("create_and_get", func(t *testing.T) {
				repository := factory(t)
				ctx := context.Background()

				created, err := repository.Create(
					ctx,
					"Kamarik",
					"kamarik@example.com",
				)
				if err != nil {
					t.Fatalf("Create() error = %v", err)
				}
				registerUserCleanup(t, repository, created.ID)

				if created.ID <= 0 {
					t.Fatalf("expected positive user ID, got %d", created.ID)
				}

				if created.Name != "Kamarik" {
					t.Fatalf("expected name Kamarik, got %q", created.Name)
				}

				if created.Email != "kamarik@example.com" {
					t.Fatalf(
						"expected email kamarik@example.com, got %q",
						created.Email,
					)
				}

				got, err := repository.Get(ctx, created.ID)
				if err != nil {
					t.Fatalf("Get() error = %v", err)
				}

				if got != created {
					t.Fatalf("Get() = %+v, want %+v", got, created)
				}
			})

			t.Run("list", func(t *testing.T) {
				repository := factory(t)
				ctx := context.Background()

				alice, err := repository.Create(
					ctx,
					"Alice",
					"alice@example.com",
				)
				if err != nil {
					t.Fatalf("Create(Alice) error = %v", err)
				}
				registerUserCleanup(t, repository, alice.ID)

				bob, err := repository.Create(
					ctx,
					"Bob",
					"bob@example.com",
				)
				if err != nil {
					t.Fatalf("Create(Bob) error = %v", err)
				}
				registerUserCleanup(t, repository, bob.ID)

				users, err := repository.List(ctx)
				if err != nil {
					t.Fatalf("List() error = %v", err)
				}

				gotByID := make(map[int]User, len(users))
				for _, user := range users {
					gotByID[user.ID] = user
				}

				if got, ok := gotByID[alice.ID]; !ok {
					t.Errorf("List() does not contain Alice with ID %d", alice.ID)
				} else if got != alice {
					t.Errorf("List() Alice = %+v, want %+v", got, alice)
				}

				if got, ok := gotByID[bob.ID]; !ok {
					t.Errorf("List() does not contain Bob with ID %d", bob.ID)
				} else if got != bob {
					t.Errorf("List() Bob = %+v, want %+v", got, bob)
				}
			})

			t.Run("update", func(t *testing.T) {
				repository := factory(t)
				ctx := context.Background()

				created, err := repository.Create(
					ctx,
					"Old Name",
					"old@example.com",
				)
				if err != nil {
					t.Fatalf("Create() error = %v", err)
				}
				registerUserCleanup(t, repository, created.ID)

				updated, err := repository.Update(
					ctx,
					created.ID,
					"New Name",
					"new@example.com",
				)
				if err != nil {
					t.Fatalf("Update() error = %v", err)
				}

				if updated.ID != created.ID {
					t.Fatalf(
						"Update() changed ID: got %d, want %d",
						updated.ID,
						created.ID,
					)
				}

				if updated.Name != "New Name" {
					t.Errorf("Update() name = %q, want %q", updated.Name, "New Name")
				}

				if updated.Email != "new@example.com" {
					t.Errorf("Update() email = %q, want %q", updated.Email, "new@example.com")
				}

				got, err := repository.Get(ctx, created.ID)
				if err != nil {
					t.Fatalf("Get() after Update() error = %v", err)
				}
				if got != updated {
					t.Errorf("Get() after Update() = %+v, want %+v", got, updated)
				}
			})

			t.Run("delete", func(t *testing.T) {
				repository := factory(t)
				ctx := context.Background()

				created, err := repository.Create(
					ctx,
					"Delete Me",
					"delete@example.com",
				)
				if err != nil {
					t.Fatalf("Create() error = %v", err)
				}
				registerUserCleanup(t, repository, created.ID)

				if err := repository.Delete(ctx, created.ID); err != nil {
					t.Fatalf("Delete() error = %v", err)
				}

				_, err = repository.Get(ctx, created.ID)
				if !errors.Is(err, ErrNotFound) {
					t.Fatalf("Get() after Delete() error = %v, want ErrNotFound", err)
				}
			})

			t.Run("missing_entities", func(t *testing.T) {
				repository := factory(t)
				ctx := context.Background()

				missing, err := repository.Create(
					ctx,
					"Temporary User",
					"temporary@example.com",
				)
				if err != nil {
					t.Fatalf("Create() error = %v", err)
				}

				missingID := missing.ID

				if err := repository.Delete(ctx, missingID); err != nil {
					t.Fatalf("Delete() before missing-entity checks error = %v", err)
				}

				_, err = repository.Get(ctx, missingID)
				if !errors.Is(err, ErrNotFound) {
					t.Fatalf("Get() error = %v, want ErrNotFound", err)
				}

				_, err = repository.Update(
					ctx,
					missingID,
					"Nobody",
					"nobody@example.com",
				)
				if !errors.Is(err, ErrNotFound) {
					t.Fatalf("Update() error = %v, want ErrNotFound", err)
				}

				err = repository.Delete(ctx, missingID)
				if !errors.Is(err, ErrNotFound) {
					t.Fatalf("Delete() error = %v, want ErrNotFound", err)
				}
			})

			t.Run("delete_does_not_affect_other_users", func(t *testing.T) {
				repository := factory(t)
				ctx := context.Background()

				first, err := repository.Create(
					ctx,
					"First",
					"first@example.com",
				)
				if err != nil {
					t.Fatalf("Create(first) error = %v", err)
				}
				registerUserCleanup(t, repository, first.ID)

				second, err := repository.Create(
					ctx,
					"Second",
					"second@example.com",
				)
				if err != nil {
					t.Fatalf("Create(second) error = %v", err)
				}
				registerUserCleanup(t, repository, second.ID)

				if err := repository.Delete(ctx, first.ID); err != nil {
					t.Fatalf("Delete(first) error = %v", err)
				}

				got, err := repository.Get(ctx, second.ID)
				if err != nil {
					t.Fatalf("Get(second) error = %v", err)
				}

				if got != second {
					t.Fatalf("Get(second) = %+v, want %+v", got, second)
				}
			})
		})
	}
}

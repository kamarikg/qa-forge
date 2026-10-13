package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgreSQLRepository struct {
	db *pgxpool.Pool
}

func NewPostgreSQLRepository(db *pgxpool.Pool) *PostgreSQLRepository {
	return &PostgreSQLRepository{
		db: db,
	}
}

func (r *PostgreSQLRepository) Create(
	ctx context.Context,
	name, email string,
) (User, error) {
	var user User

	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id, name, email
		`,
		name,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)

	if err != nil {
		return User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (r *PostgreSQLRepository) List(
	ctx context.Context,
) ([]User, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT id, name, email
		FROM users
		ORDER BY id
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	users := make([]User, 0)

	for rows.Next() {
		var user User

		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
		); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return users, nil
}

func (r *PostgreSQLRepository) Get(
	ctx context.Context,
	id int,
) (User, error) {
	var user User

	err := r.db.QueryRow(
		ctx,
		`
		SELECT id, name, email
		FROM users
		WHERE id = $1
		`,
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}

	if err != nil {
		return User{}, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}

func (r *PostgreSQLRepository) Update(
	ctx context.Context,
	id int,
	name, email string,
) (User, error) {
	var user User

	err := r.db.QueryRow(
		ctx,
		`
		UPDATE users
		SET name = $2, email = $3
		WHERE id = $1
		RETURNING id, name, email
		`,
		id,
		name,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}

	if err != nil {
		return User{}, fmt.Errorf("update user: %w", err)
	}

	return user, nil
}

func (r *PostgreSQLRepository) Delete(
	ctx context.Context,
	id int,
) error {
	result, err := r.db.Exec(
		ctx,
		`
		DELETE FROM users
		WHERE id = $1
		`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

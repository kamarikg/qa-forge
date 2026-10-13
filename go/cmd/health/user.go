package main

import "context"

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserStore struct {
	users      map[int]User
	nextUserID int
}

func NewUserStore() *UserStore {
	return &UserStore{
		users:      make(map[int]User),
		nextUserID: 1,
	}
}

func (s *UserStore) Create(ctx context.Context, name, email string) (User, error) {
	user := User{
		ID:    s.nextUserID,
		Name:  name,
		Email: email,
	}

	s.users[user.ID] = user
	s.nextUserID++

	return user, nil
}

func (s *UserStore) List(ctx context.Context) ([]User, error) {
	users := make([]User, 0, len(s.users))

	for _, user := range s.users {
		users = append(users, user)
	}

	return users, nil
}

func (s *UserStore) Get(ctx context.Context, id int) (User, error) {
	user, ok := s.users[id]
	if !ok {
		return User{}, ErrNotFound
	}

	return user, nil
}

func (s *UserStore) Update(ctx context.Context, id int, name, email string) (User, error) {
	if _, ok := s.users[id]; !ok {
		return User{}, ErrNotFound
	}

	user := User{
		ID:    id,
		Name:  name,
		Email: email,
	}

	s.users[id] = user

	return user, nil
}

func (s *UserStore) Delete(ctx context.Context, id int) error {
	if _, ok := s.users[id]; !ok {
		return ErrNotFound
	}

	delete(s.users, id)

	return nil
}

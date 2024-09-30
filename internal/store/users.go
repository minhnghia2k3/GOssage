package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type IUsers interface {
	Create(context.Context, *User) error
	GetByID(context.Context, int64) (*User, error)
}

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserStorage struct {
	db *sql.DB
}

func (s *UserStorage) Create(ctx context.Context, users *User) error {
	query := `
	INSERT INTO users (username, email, password)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at
`
	err := s.db.QueryRowContext(
		ctx,
		query,
		users.Username,
		users.Email,
		users.Password,
	).Scan(
		&users.ID,
		&users.CreatedAt,
		&users.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStorage) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
	SELECT id, username, email, created_at, updated_at
	FROM users
	WHERE id = $1
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var user User
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

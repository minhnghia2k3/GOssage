package storage

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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

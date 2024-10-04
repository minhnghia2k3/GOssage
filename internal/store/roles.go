package store

import (
	"context"
	"database/sql"
	"errors"
)

type IRoles interface {
	GetByName(ctx context.Context, name string) (*Role, error)
}

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Level       int64  `json:"level"`
	Description string `json:"description"`
}

type RoleStorage struct {
	db *sql.DB
}

func (s *RoleStorage) GetByName(ctx context.Context, name string) (*Role, error) {
	query := `
	SELECT id, name, level, description FROM roles
	WHERE name = $1;
`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var role Role
	err := s.db.QueryRowContext(ctx, query, name).Scan(
		&role.ID,
		&role.Name,
		&role.Level,
		&role.Description,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &role, nil
}

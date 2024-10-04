package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	ErrFollowSelf        = errors.New("cannot follow yourself")
	QueryTimeOutDuration = 5 * time.Second
)

type Storage struct {
	Posts     IPosts
	Users     IUsers
	Followers IFollower
	Comments  IComments
	Roles     IRoles
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStorage{db: db},
		Users:     &UserStorage{db: db},
		Followers: &FollowerStorage{db: db},
		Comments:  &CommentStorage{db: db},
		Roles:     &RoleStorage{db: db},
	}
}

func withTx(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err = fn(tx); err != nil {
		// Rollback in case anything fails in the closure
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

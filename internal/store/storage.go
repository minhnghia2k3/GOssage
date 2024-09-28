package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	QueryTimeOutDuration = 5 * time.Second
)

type Storage struct {
	Posts IPosts
	Users interface {
		Create(context.Context, *User) error
	}
	Comments IComments
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStorage{db: db},
		Users:    &UserStorage{db: db},
		Comments: &CommentStorage{db: db},
	}
}

package store

import (
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
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStorage{db: db},
		Users:     &UserStorage{db: db},
		Followers: &FollowerStorage{db: db},
		Comments:  &CommentStorage{db: db},
	}
}

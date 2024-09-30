package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"time"
)

type IFollower interface {
	Follow(ctx context.Context, userID, followerID int64) error
	Unfollow(ctx context.Context, userID, followerID int64) error
}

type Follower struct {
	UserID     int64     `json:"user_id"`
	FollowerID int64     `json:"follower_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type FollowerStorage struct {
	db *sql.DB
}

func (s *FollowerStorage) Follow(ctx context.Context, userID, followerID int64) error {
	query := `
	INSERT INTO followers (user_id, follower_id)
	VALUES ($1, $2)
`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, followerID)

	if err != nil {
		var pqError *pq.Error
		switch errors.As(err, &pqError) {
		case pqError.Code == "23505":
			return ErrConflict
		case pqError.Code == "23514":
			return ErrFollowSelf
		default:
			return err
		}
	}

	return nil
}

func (s *FollowerStorage) Unfollow(ctx context.Context, userID, followerID int64) error {
	query := `
	DELETE FROM followers
	WHERE user_id = $1 AND follower_id = $2
`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, userID, followerID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return ErrNotFound
	}

	return err
}

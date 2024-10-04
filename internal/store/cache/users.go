package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/minhnghia2k3/GOssage/internal/store"
	"github.com/redis/go-redis/v9"
	"time"
)

type IUsers interface {
	Get(ctx context.Context, userID int64) (*store.User, error)
	Set(ctx context.Context, user *store.User) error
}

type UserStorage struct {
	rdb *redis.Client
}

func (s *UserStorage) Get(ctx context.Context, userID int64) (*store.User, error) {
	key := fmt.Sprintf("user-%d", userID)

	val, err := s.rdb.Get(ctx, key).Result()

	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			return nil, nil
		default:
			return nil, err
		}
	}

	// Unmarshal data JSON => Go struct
	var user store.User
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

const UserExpTime = time.Minute

func (s *UserStorage) Set(ctx context.Context, user *store.User) error {
	key := fmt.Sprintf("user-%d", user.ID)

	// Marshal data GO struct => JSON
	v, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.rdb.Set(ctx, key, v, UserExpTime).Err()
}

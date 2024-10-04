package cache

import "github.com/redis/go-redis/v9"

type Storage struct {
	Users IUsers
}

func NewRedisStorage(rdb *redis.Client) *Storage {
	return &Storage{
		Users: &UserStorage{rdb: rdb},
	}
}

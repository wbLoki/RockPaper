package db

import "github.com/redis/go-redis/v9"

func NewRedisStorage(cfg redis.Options) *redis.Client {
	rdb := redis.NewClient(&cfg)
	return rdb
}

// &redis.Options{
// 	Addr:     "localhost:6379",
// 	Password: "password",
// 	DB:       0,
// }

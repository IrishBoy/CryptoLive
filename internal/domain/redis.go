package domain

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// InitializeRedis initializes connection with Redis
func InitializeRedis(addr, password string, userName string, db int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: userName,
		Password: password, // no password set if password is empty
		DB:       db,       // use default DB
	})

	// Ping to check if connection is successful
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

// func NewRedisClient(userName string, password string) *redis.Client {

// 	return redis.NewClient(&redis.Options{
// 		Addr:     "redis-14349.c304.europe-west1-2.gce.cloud.redislabs.com:14349",
// 		Password: password,
// 		Username: userName,
// 		DB:       0,
// 	})
// }

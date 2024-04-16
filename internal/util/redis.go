package util

import (
	"fmt"
	"github.com/go-redis/redis/v8"
)

func NewRedis(env *Env) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", env.RedisHost, env.RedisPort),
		Password: env.RedisPassword,
		DB:       0,
	})
}

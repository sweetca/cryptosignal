package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

func NewClient(settings *Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", settings.RedisHost, settings.RedisPort),
		Password:     "",
		DB:           0,
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		PoolTimeout:  time.Second,
		PoolSize:     settings.RedisPoolSize,
	})

	return client
}

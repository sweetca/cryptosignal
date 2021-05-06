package redis

import (
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"log"
)

const (
	ChannelWriteStatDone = "write_stat_done"
)

type Config struct {
	RedisPort     int    `envconfig:"redis_port" default:"6379"`
	RedisHost     string `envconfig:"redis_host" default:"127.0.0.1"`
	RedisPoolSize int    `envconfig:"redis_pool_size" default:"10"`
}

func Init() *Config {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		log.Fatal("redis config issue", zap.Error(err))
	}

	return &c
}

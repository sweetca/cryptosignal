package config

import (
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"log"
)

type Config struct {
	AppPort int `envconfig:"data_maker_app_port" default:"8082"`

	RedisPort     int    `envconfig:"redis_port" default:"6379"`
	RedisHost     string `envconfig:"redis_host" default:"127.0.0.1"`
	RedisPoolSize int    `envconfig:"redis_pool_size" default:"10"`
}

func Init() *Config {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		log.Fatal("data science config issue", zap.Error(err))
	}

	return &c
}

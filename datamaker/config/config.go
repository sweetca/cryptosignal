package config

import (
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"log"
)

type Config struct {
	AppPort              int      `envconfig:"data_maker_app_port"`
	CryptoList           []string `envconfig:"crypto_list"`
	CoinMarketCapAPI     string   `envconfig:"data_maker_coinmarketcap_api"`
	CoinMarketCapKey     string   `envconfig:"data_maker_coinmarketcap_key"`
	CoinMarketCapConvert string   `envconfig:"data_maker_coinmarketcap_convert"`

	RedisPort     int    `envconfig:"redis_port"`
	RedisHost     string `envconfig:"redis_host"`
	RedisPoolSize int    `envconfig:"redis_pool_size"`
}

func Init() *Config {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		log.Fatal("data maker config issue", zap.Error(err))
	}

	return &c
}

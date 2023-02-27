package config

import "github.com/caarlos0/env/v7"

type Config struct {
	Env  string `env:"TODO_ENV" envDefault:"dev"`
	Port int    `env:"PORT" envDefault:"80"`
}

func NewConfig() *Config {
	return &Config{}
}

func New() (*Config, error) {
	cfg := NewConfig()

	// 環境変数をパースする
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

package config

import "github.com/caarlos0/env"

type Config struct {
	Address string `env:"ADDRESS" envDefault:"localhost:8080"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	return cfg, err
}

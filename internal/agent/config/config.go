package config

import "github.com/caarlos0/env"

type Config struct {
	Address        string `env:"ADDRESS" envDefault:"localhost:8080"`
	PullInterval   int    `env:"POLL_INTERVAL" envDefault:"2"`
	ReportInterval int    `env:"REPORT_INTERVAL" envDefault:"10"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	return cfg, err
}

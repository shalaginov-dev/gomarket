package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port      string `envconfig:"PORT" default:"8080"`
	Env       string `envconfig:"ENV" default:"development"`
	DBDSN     string `envconfig:"DB_DSN"`
	JWTSecret string `envconfig:"JWT_SECRET"`
	JWTExpiry int    `envconfig:"JWT_EXPIRY" default:"15"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

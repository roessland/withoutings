package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DatabaseURL          string `envconfig:"DATABASE_URL"`
	SessionSecret        []byte `envconfig:"SESSION_SECRET"`
	WithingsClientID     string `envconfig:"WITHINGS_CLIENT_ID"`
	WithingsClientSecret string `envconfig:"WITHINGS_CLIENT_SECRET"`
	WithingsRedirectURL  string `envconfig:"WITHINGS_REDIRECT_URL"`
}

func LoadFromEnv() (*Config, error) {
	c := Config{}
	err := envconfig.Process("WOT", &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

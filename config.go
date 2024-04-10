package main

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	PostgresDSN        string `env:"POSTGRES_DSN,required"`
	HTTPServerAddr     string `env:"HTTP_SERVER_ADDR,default=:8000"`
	StravaClientID     string `env:"STRAVA_CLIENT_ID,required"`
	StravaClientSecret string `env:"STRAVA_CLIENT_SECRET,required"`
	StravaRedirectURL  string `env:"STRAVA_REDIRECT_URL,default=http://localhost:8000/strava/callback"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

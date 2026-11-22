package server

import (
	"fmt"
	"net/url"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr            string        `envconfig:"ADDR" required:"true"`
	ShutdownTimeout time.Duration `envconfig:"TIMEOUT" default:"30s"`
	FrontendBaseURL string        `envconfig:"FRONTEND_BASE_URL" required:"true"`
}

func NewConfig() (Config, error) {
	var cfg Config

	if err := envconfig.Process("HTTP", &cfg); err != nil {
		return Config{}, fmt.Errorf("envconfig process: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, fmt.Errorf("validate: %w", err)
	}

	return cfg, nil
}

func (c Config) Validate() error {
	u, err := url.Parse(c.FrontendBaseURL)
	if err != nil {
		return fmt.Errorf("parse FRONTEND_BASE_URL: %w", err)
	}
	if u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("invalid FRONTEND_BASE_URL: %q", c.FrontendBaseURL)
	}
	return nil
}

func NewConfigMust() Config {
	cfg, err := NewConfig()
	if err != nil {
		panic(err)
	}

	return cfg
}

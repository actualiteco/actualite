package smtpmailhogclient

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Host    string        `envconfig:"HOST" required:"true"`
	Port    int           `envconfig:"PORT" default:"1025"`
	From    string        `envconfig:"FROM" default:"no-reply@actualite.co"`
	Timeout time.Duration `envconfig:"TIMEOUT" default:"10s"`
}

func newConfig() (config, error) {
	var cfg config
	if err := envconfig.Process("SMTP_MAILHOG", &cfg); err != nil {
		return config{}, fmt.Errorf("load smtp mailhog config: %w", err)
	}
	return cfg, nil
}

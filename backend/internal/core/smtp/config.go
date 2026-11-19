package smtpclient

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Host     string        `envconfig:"HOST" required:"true"`
	Port     int           `envconfig:"PORT" required:"true"`
	Username string        `envconfig:"USERNAME" required:"true"`
	Password string        `envconfig:"PASSWORD" required:"true"`
	From     string        `envconfig:"FROM" default:"actualiteco@gmail.com"`
	Timeout  time.Duration `envconfig:"TIMEOUT" required:"true"`
}

func NewConfigMust() config {
	cfg, err := newConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}

func newConfig() (config, error) {
	var cfg config
	if err := envconfig.Process("SMTP", &cfg); err != nil {
		return config{}, fmt.Errorf("load smtp config: %w", err)
	}

	return cfg, nil
}

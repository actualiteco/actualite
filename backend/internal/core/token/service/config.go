package tokenservice

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Secret     []byte        `envconfig:"SECRET"`
	AccessTTL  time.Duration `envconfig:"ACCESS_TTL"`
	RefreshTTL time.Duration `envconfig:"REFRESH_TTL"`
}

func NewConfig() (config, error) {
	var cfg config
	if err := envconfig.Process("TOKEN", &cfg); err != nil {
		return config{}, fmt.Errorf("token service config: %w", err)
	}

	return cfg, nil
}

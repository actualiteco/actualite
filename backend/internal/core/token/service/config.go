package tokenservice

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Secret           []byte        `envconfig:"SECRET" required:"true"`
	AccessTTL        time.Duration `envconfig:"ACCESS_TTL" required:"true"`
	RefreshTTL       time.Duration `envconfig:"REFRESH_TTL" required:"true"`
	PasswordResetTTL time.Duration `envconfig:"PASSWORD_RESET_TTL" required:"true"`
}

func NewConfig() (config, error) {
	var cfg config
	if err := envconfig.Process("TOKEN", &cfg); err != nil {
		return config{}, fmt.Errorf("token service config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return config{}, fmt.Errorf("validate token service config: %w", err)
	}

	return cfg, nil
}

func (c *config) Validate() error {
	if len(c.Secret) == 0 {
		return fmt.Errorf("TOKEN_SECRET is empty")
	}

	if c.AccessTTL <= 0 {
		return fmt.Errorf("TOKEN_ACCESS_TTL=%v must be greater than 0", c.AccessTTL)
	}

	if c.RefreshTTL <= 0 {
		return fmt.Errorf("TOKEN_REFRESH_TTL=%v must be greater than 0", c.RefreshTTL)
	}

	if c.RefreshTTL < c.AccessTTL {
		return fmt.Errorf("TOKEN_REFRESH_TTL=%v must be greater than TOKEN_ACCESS_TTL=%v", c.RefreshTTL, c.AccessTTL)
	}

	return nil
}

package storage

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	MasterURL    string        `envconfig:"MASTER_URL" required:"true"`
	PublicURL    string        `envconfig:"PUBLIC_URL" required:"true"`
	VolumeServer string        `envconfig:"VOLUME_SERVER" required:"true"`
	Timeout      time.Duration `envconfig:"TIMEOUT" default:"10s"`
}

func NewConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("SEAWEEDFS", &cfg); err != nil {
		return Config{}, fmt.Errorf("envconfig process: %w", err)
	}

	return cfg, nil
}

func NewConfigMust() Config {
	cfg, err := NewConfig()
	if err != nil {
		panic(err)
	}

	return cfg
}

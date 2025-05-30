package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

const (
	defaultOrderSizeLimit = 1000000
	defaultCORSOrigin     = "*"
)

// Config holds the configuration for the application.
type Config struct {
	Port           string `mapstructure:"port"`
	CORSOrigin     string `mapstructure:"cors_origin"`
	OrderSizeLimit int    `mapstructure:"order_size_limit"`
}

// New returns a new, populated instance of Config.
func New() (Config, error) {
	v := viper.New()

	v.SetDefault("order_size_limit", defaultOrderSizeLimit)
	v.SetDefault("cors_origin", defaultCORSOrigin)
	v.AutomaticEnv()

	err := v.BindEnv("port", "PORT")
	if err != nil {
		return Config{}, fmt.Errorf("failed to bind env var for port: %w", err)
	}
	err = v.BindEnv("cors_origin", "CORS_ORIGIN")
	if err != nil {
		return Config{}, fmt.Errorf("failed to bind env var for cors_origin: %w", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if config.Port == "" {
		return Config{}, errors.New("port must be set")
	}

	return config, nil
}

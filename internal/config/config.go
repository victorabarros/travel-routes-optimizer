package config

import (
	"github.com/spf13/viper"
)

// Config summarises all environment variables.
type Config struct {
	Port     int    `mapstructure:"port"`
	LogLevel string `mapstructure:"log_level"`
	Server   string `mapstructure:"server"`
}

var (
	cfg *Config
)

// Load return all environment variables loaded.
func Load() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}
	viper.SetDefault("LOG_LEVEL", "INFO")
	viper.SetDefault("PORT", 8092)
	viper.SetDefault("SERVER", "challenge-bexs-server")

	viper.AutomaticEnv()
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

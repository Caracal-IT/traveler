package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds application configuration.
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Log    LogConfig    `mapstructure:"log"`
}

// ServerConfig holds server-specific configuration.
type ServerConfig struct {
	Port int `mapstructure:"port"`
}

// LogConfig holds logging configuration.
type LogConfig struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"` // Optional: if empty, logs only to stdout
}

// Load reads configuration from a YAML file.
func Load(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// Set defaults
	v.SetDefault("server.port", 8080)
	v.SetDefault("log.level", "info")
	v.SetDefault("log.file", "") // Empty means stdout only

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// LoadOrDefault attempts to load configuration from the given path,
// falling back to defaults if the file doesn't exist.
func LoadOrDefault(configPath string) *Config {
	cfg, err := Load(configPath)
	if err != nil {
		// Return default config
		return &Config{
			Server: ServerConfig{Port: 8080},
			Log:    LogConfig{Level: "info"},
		}
	}
	return cfg
}

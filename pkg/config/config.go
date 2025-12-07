package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds application configuration.
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Log    LogConfig    `mapstructure:"log"`
	Auth   AuthConfig   `mapstructure:"auth"`
}

// ServerConfig holds server-specific configuration.
type ServerConfig struct {
	Port int `mapstructure:"port"`
}

// LogConfig holds logging configuration.
type LogConfig struct {
	Level         string           `mapstructure:"level"`
	File          string           `mapstructure:"file"` // Optional: if empty, logs only to stdout
	Elasticsearch ElasticLogConfig `mapstructure:"elasticsearch"`
}

// ElasticLogConfig controls optional shipping of logs to Elasticsearch.
type ElasticLogConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	URL     string `mapstructure:"url"`   // e.g. http://localhost:9200
	Index   string `mapstructure:"index"` // e.g. traveler-logs
	Buffer  int    `mapstructure:"buffer"`
	Workers int    `mapstructure:"workers"`
}

// AuthConfig holds authentication settings (Keycloak/OpenID Connect).
type AuthConfig struct {
	// Issuer is the base issuer URL of the realm, e.g. http://localhost:8081/realms/traveler-dev
	Issuer string `mapstructure:"issuer"`
	// Audience is the expected audience/client_id in tokens, e.g. traveler-app
	Audience string `mapstructure:"audience"`
	// JWKSURL optionally overrides the JWKS endpoint URL used to validate tokens.
	// If empty, it will be derived from Issuer as: <issuer>/protocol/openid-connect/certs
	JWKSURL string `mapstructure:"jwks_url"`
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
	// Elastic logging defaults (disabled)
	v.SetDefault("log.elasticsearch.enabled", false)
	v.SetDefault("log.elasticsearch.url", "http://localhost:9200")
	v.SetDefault("log.elasticsearch.index", "traveler-logs")
	v.SetDefault("log.elasticsearch.buffer", 1024)
	v.SetDefault("log.elasticsearch.workers", 1)
	// Reasonable dev defaults for local Keycloak in docker
	v.SetDefault("auth.issuer", "http://localhost:8081/realms/traveler-dev")
	v.SetDefault("auth.audience", "traveler-app")

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

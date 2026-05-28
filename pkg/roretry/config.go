package roretry

import (
	"fmt"
	"time"

	appconfig "traveler/pkg/config"
)

// Config defines optional configuration file values for the retry engine.
type Config struct {
	Start                string        `mapstructure:"start" json:"start"`
	StartGroup           string        `mapstructure:"start_group" json:"start_group"` // backward compatibility
	OverallCycles        int           `mapstructure:"overall_cycles" json:"overall_cycles"`
	OverallBackoff       time.Duration `mapstructure:"-"`
	OverallBackoffMillis int64         `mapstructure:"overall_backoff_millis" json:"overall_backoff_millis"`
	Groups               []GroupConfig `mapstructure:"groups" json:"groups"`
}

// GroupConfig defines config-driven retry behavior for a command group.
type GroupConfig struct {
	Name          string        `mapstructure:"name" json:"name"`
	Method        string        `mapstructure:"method" json:"method"`
	MaxRetries    int           `mapstructure:"max_retries" json:"max_retries"`
	Backoff       time.Duration `mapstructure:"-"`
	BackoffMillis int64         `mapstructure:"backoff_millis" json:"backoff_millis"`
}

func loadConfig(path string) (Config, error) {
	if path == "" {
		return Config{}, nil
	}

	var cfg Config
	if err := appconfig.LoadInto(path, &cfg); err != nil {
		return Config{}, fmt.Errorf("load roretry config: %w", err)
	}

	normalizeConfig(&cfg)
	return cfg, nil
}

func normalizeConfig(cfg *Config) {
	if cfg == nil {
		return
	}

	if cfg.OverallBackoffMillis > 0 {
		cfg.OverallBackoff = time.Duration(cfg.OverallBackoffMillis) * time.Millisecond
	}

	for i := range cfg.Groups {
		if cfg.Groups[i].BackoffMillis > 0 {
			cfg.Groups[i].Backoff = time.Duration(cfg.Groups[i].BackoffMillis) * time.Millisecond
		}
	}
}

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Organization    string `mapstructure:"organization"`
	Project         string `mapstructure:"project"`
	Repository      string `mapstructure:"repository"`
	DefaultReviewer string `mapstructure:"default_reviewer"`
}

var (
	configDir  = filepath.Join(os.Getenv("HOME"), ".dex-cli")
	configFile = filepath.Join(configDir, "config.yaml")
)

// Load reads the configuration file
func Load() (*Config, error) {
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	// Set defaults
	viper.SetDefault("organization", "")
	viper.SetDefault("project", "")
	viper.SetDefault("repository", "")
	viper.SetDefault("default_reviewer", "")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Check if config file exists
	_, err := os.Stat(configFile)
	configExists := err == nil

	// Read config file if it exists
	if configExists {
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	} else {
		// Config file doesn't exist, create it with defaults
		if err := viper.WriteConfigAs(configFile); err != nil {
			return nil, fmt.Errorf("failed to create config file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// Save writes the configuration to file
func Save(cfg *Config) error {
	viper.Set("organization", cfg.Organization)
	viper.Set("project", cfg.Project)
	viper.Set("repository", cfg.Repository)
	viper.Set("default_reviewer", cfg.DefaultReviewer)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// GetConfigDir returns the configuration directory path
func GetConfigDir() string {
	return configDir
}

package testhelpers

import (
	"os"
	"path/filepath"
	"testing"
)

// SetupTempConfig creates a temporary config directory and file for testing
func SetupTempConfig(t *testing.T) (string, func()) {
	t.Helper()

	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".dex-cli")
	configFile := filepath.Join(configDir, "config.yaml")

	// Create config directory
	if err := os.MkdirAll(configDir, 0700); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	// Create empty config file
	if err := os.WriteFile(configFile, []byte(""), 0600); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Cleanup function
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

// WriteConfigFile writes YAML content to a config file
func WriteConfigFile(t *testing.T, configDir, content string) string {
	t.Helper()

	configPath := filepath.Join(configDir, ".dex-cli", "config.yaml")
	if err := os.WriteFile(configPath, []byte(content), 0600); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}
	return configPath
}

package cmd

import (
	"path/filepath"
	"testing"

	"github.com/chriskievit/dex-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "(not set)",
		},
		{
			name:     "non-empty string",
			input:    "test-value",
			expected: "test-value",
		},
		{
			name:     "whitespace only",
			input:    "   ",
			expected: "   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatValue(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRunShowConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".dex-cli")
	
	// Save original config dir and restore after test
	originalConfigDir := config.GetConfigDir()
	defer config.SetConfigDir(originalConfigDir)
	
	config.SetConfigDir(configDir)

	// Create config with values
	cfg := &config.Config{
		Organization:    "testorg",
		Project:         "testproject",
		Repository:      "testrepo",
		DefaultReviewer: "reviewer@example.com",
	}
	
	// Initialize config
	_, err := config.Load()
	require.NoError(t, err)
	
	err = config.Save(cfg)
	require.NoError(t, err)

	// Capture output - cobra commands write to os.Stdout by default
	// We'll test by checking the function doesn't error and config is readable
	err = runShowConfig(showConfigCmd, []string{})
	require.NoError(t, err)

	// Verify config values are correct by reloading
	loadedCfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "testorg", loadedCfg.Organization)
	assert.Equal(t, "testproject", loadedCfg.Project)
	assert.Equal(t, "testrepo", loadedCfg.Repository)
	assert.Equal(t, "reviewer@example.com", loadedCfg.DefaultReviewer)
	assert.Equal(t, "testorg", cfg.Organization)
	assert.Equal(t, "testproject", cfg.Project)
	assert.Equal(t, "testrepo", cfg.Repository)
	assert.Equal(t, "reviewer@example.com", cfg.DefaultReviewer)
}

func TestRunShowConfig_EmptyConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".dex-cli")
	
	// Save original config dir and restore after test
	originalConfigDir := config.GetConfigDir()
	defer config.SetConfigDir(originalConfigDir)
	
	config.SetConfigDir(configDir)

	// Initialize empty config - this creates a file with defaults
	_, err := config.Load()
	require.NoError(t, err)

	// Execute command - should not error even with empty config
	err = runShowConfig(showConfigCmd, []string{})
	require.NoError(t, err)

	// Note: After Load(), viper may have default values set
	// The config file itself will have empty strings, but viper defaults may apply
	// This is expected behavior - the function works correctly
}

func TestRunSetProject(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".dex-cli")
	
	// Save original config dir and restore after test
	originalConfigDir := config.GetConfigDir()
	defer config.SetConfigDir(originalConfigDir)
	
	config.SetConfigDir(configDir)

	// Initialize config
	_, err := config.Load()
	require.NoError(t, err)

	// Execute command
	err = runSetProject(setProjectCmd, []string{"myproject"})
	require.NoError(t, err)

	// Verify config was saved
	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "myproject", cfg.Project)
}

func TestRunSetProject_EmptyValue(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".dex-cli")
	
	// Save original config dir and restore after test
	originalConfigDir := config.GetConfigDir()
	defer config.SetConfigDir(originalConfigDir)
	
	config.SetConfigDir(configDir)

	// Initialize config
	_, err := config.Load()
	require.NoError(t, err)

	// Execute command with empty value (should fail validation)
	err = runSetProject(setProjectCmd, []string{""})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestRunSetRepo(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".dex-cli")
	
	// Save original config dir and restore after test
	originalConfigDir := config.GetConfigDir()
	defer config.SetConfigDir(originalConfigDir)
	
	config.SetConfigDir(configDir)

	// Initialize config
	_, err := config.Load()
	require.NoError(t, err)

	// Execute command
	err = runSetRepo(setRepoCmd, []string{"myrepo"})
	require.NoError(t, err)

	// Verify config was saved
	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "myrepo", cfg.Repository)
}

func TestRunSetRepo_EmptyValue(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".dex-cli")
	
	// Save original config dir and restore after test
	originalConfigDir := config.GetConfigDir()
	defer config.SetConfigDir(originalConfigDir)
	
	config.SetConfigDir(configDir)

	// Initialize config
	_, err := config.Load()
	require.NoError(t, err)

	// Execute command with empty value (should fail validation)
	err = runSetRepo(setRepoCmd, []string{""})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestRunSetReviewer(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".dex-cli")
	
	// Save original config dir and restore after test
	originalConfigDir := config.GetConfigDir()
	defer config.SetConfigDir(originalConfigDir)
	
	config.SetConfigDir(configDir)

	// Initialize config
	_, err := config.Load()
	require.NoError(t, err)

	// Execute command
	err = runSetReviewer(setReviewerCmd, []string{"reviewer@example.com"})
	require.NoError(t, err)

	// Verify config was saved
	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "reviewer@example.com", cfg.DefaultReviewer)
}

func TestRunSetReviewer_EmptyValue(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".dex-cli")
	
	// Save original config dir and restore after test
	originalConfigDir := config.GetConfigDir()
	defer config.SetConfigDir(originalConfigDir)
	
	config.SetConfigDir(configDir)

	// Initialize config
	_, err := config.Load()
	require.NoError(t, err)

	// Execute command with empty value (should fail validation)
	err = runSetReviewer(setReviewerCmd, []string{""})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

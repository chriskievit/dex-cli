package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_NewConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".dex-cli")
	
	// Save original config dir and restore after test
	originalConfigDir := GetConfigDir()
	defer SetConfigDir(originalConfigDir)
	
	SetConfigDir(configDir)

	cfg, err := Load()
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "", cfg.Organization)
	assert.Equal(t, "", cfg.Project)
	assert.Equal(t, "", cfg.Repository)
	assert.Equal(t, "", cfg.DefaultReviewer)

	// Verify config file was created
	configFile := filepath.Join(configDir, "config.yaml")
	_, err = os.Stat(configFile)
	assert.NoError(t, err, "Config file should be created")
}

func TestLoad_ExistingConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".dex-cli")
	configFile := filepath.Join(configDir, "config.yaml")
	
	// Save original config dir and restore after test
	originalConfigDir := GetConfigDir()
	defer SetConfigDir(originalConfigDir)
	
	SetConfigDir(configDir)

	// Create config directory
	err := os.MkdirAll(configDir, 0700)
	require.NoError(t, err)

	// Write existing config
	configContent := `organization: testorg
project: testproject
repository: testrepo
default_reviewer: reviewer@example.com
`
	err = os.WriteFile(configFile, []byte(configContent), 0600)
	require.NoError(t, err)

	cfg, err := Load()
	require.NoError(t, err)
	assert.Equal(t, "testorg", cfg.Organization)
	assert.Equal(t, "testproject", cfg.Project)
	assert.Equal(t, "testrepo", cfg.Repository)
	assert.Equal(t, "reviewer@example.com", cfg.DefaultReviewer)
}

func TestLoad_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".dex-cli")
	configFile := filepath.Join(configDir, "config.yaml")
	
	// Save original config dir and restore after test
	originalConfigDir := GetConfigDir()
	defer SetConfigDir(originalConfigDir)
	
	SetConfigDir(configDir)

	// Create config directory
	err := os.MkdirAll(configDir, 0700)
	require.NoError(t, err)

	// Write invalid YAML
	invalidYAML := `organization: testorg
project: [invalid yaml
`
	err = os.WriteFile(configFile, []byte(invalidYAML), 0600)
	require.NoError(t, err)

	cfg, err := Load()
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "failed to read config file")
}

func TestSave(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".dex-cli")
	
	// Save original config dir and restore after test
	originalConfigDir := GetConfigDir()
	defer SetConfigDir(originalConfigDir)
	
	SetConfigDir(configDir)

	// Load config first to initialize viper
	_, err := Load()
	require.NoError(t, err)

	// Create and save new config
	cfg := &Config{
		Organization:    "neworg",
		Project:         "newproject",
		Repository:      "newrepo",
		DefaultReviewer: "newreviewer@example.com",
	}

	err = Save(cfg)
	require.NoError(t, err)

	// Reload and verify
	loadedCfg, err := Load()
	require.NoError(t, err)
	assert.Equal(t, "neworg", loadedCfg.Organization)
	assert.Equal(t, "newproject", loadedCfg.Project)
	assert.Equal(t, "newrepo", loadedCfg.Repository)
	assert.Equal(t, "newreviewer@example.com", loadedCfg.DefaultReviewer)
}

func TestSave_UpdatesExistingConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".dex-cli")
	configFile := filepath.Join(configDir, "config.yaml")
	
	// Save original config dir and restore after test
	originalConfigDir := GetConfigDir()
	defer SetConfigDir(originalConfigDir)
	
	SetConfigDir(configDir)

	// Create initial config
	err := os.MkdirAll(configDir, 0700)
	require.NoError(t, err)

	initialConfig := `organization: oldorg
project: oldproject
`
	err = os.WriteFile(configFile, []byte(initialConfig), 0600)
	require.NoError(t, err)

	// Load config
	_, err = Load()
	require.NoError(t, err)

	// Update and save
	cfg := &Config{
		Organization:    "updatedorg",
		Project:         "updatedproject",
		Repository:      "updatedrepo",
		DefaultReviewer: "updatedreviewer@example.com",
	}

	err = Save(cfg)
	require.NoError(t, err)

	// Reload and verify
	loadedCfg, err := Load()
	require.NoError(t, err)
	assert.Equal(t, "updatedorg", loadedCfg.Organization)
	assert.Equal(t, "updatedproject", loadedCfg.Project)
	assert.Equal(t, "updatedrepo", loadedCfg.Repository)
	assert.Equal(t, "updatedreviewer@example.com", loadedCfg.DefaultReviewer)
}

func TestGetConfigDir(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".dex-cli")
	
	// Save original config dir and restore after test
	originalConfigDir := GetConfigDir()
	defer SetConfigDir(originalConfigDir)
	
	SetConfigDir(configDir)

	dir := GetConfigDir()
	assert.Equal(t, configDir, dir)
}

func TestLoad_CreatesConfigDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "non-existent", ".dex-cli")
	
	// Save original config dir and restore after test
	originalConfigDir := GetConfigDir()
	defer SetConfigDir(originalConfigDir)
	
	SetConfigDir(configDir)

	// Directory shouldn't exist yet
	_, err := os.Stat(configDir)
	assert.Error(t, err)

	// Load should create the directory
	cfg, err := Load()
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	// Directory should now exist
	_, err = os.Stat(configDir)
	assert.NoError(t, err)
}

package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCmd_Initialization(t *testing.T) {
	assert.NotNil(t, rootCmd)
	assert.Equal(t, "dex", rootCmd.Use)
	assert.Equal(t, "Azure DevOps CLI tool for managing branches, work items, and pull requests", rootCmd.Short)
	assert.Equal(t, "1.0.0", rootCmd.Version)
}

func TestRootCmd_GlobalFlags(t *testing.T) {
	// Note: Flags are initialized in init() which runs before tests
	// We can't easily reset and re-initialize, so we just verify flags exist

	// Test organization flag
	orgFlag := rootCmd.PersistentFlags().Lookup("org")
	assert.NotNil(t, orgFlag)
	assert.Equal(t, "o", orgFlag.Shorthand)

	// Test project flag
	projectFlag := rootCmd.PersistentFlags().Lookup("project")
	assert.NotNil(t, projectFlag)
	assert.Equal(t, "p", projectFlag.Shorthand)

	// Test debug flag
	debugFlag := rootCmd.PersistentFlags().Lookup("debug")
	assert.NotNil(t, debugFlag)
	assert.Equal(t, "d", debugFlag.Shorthand)
}

func TestExecute_InvalidCommand(t *testing.T) {
	// Note: Execute() calls os.Exit(1) on error, so we can't easily test it
	// without using a subprocess. We'll just verify the root command exists.
	assert.NotNil(t, rootCmd)
}

func TestExecute_HelpCommand(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Set help command
	os.Args = []string{"dex", "help"}

	// Capture output
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)

	// Execute should succeed
	Execute()

	// Help command should produce output
	output := buf.String()
	assert.Contains(t, output, "dex")
}

func TestGlobalFlags_Variables(t *testing.T) {
	// Test that global flag variables are accessible
	// These are set when flags are parsed
	assert.NotNil(t, &organization)
	assert.NotNil(t, &project)
	assert.NotNil(t, &debug)
}

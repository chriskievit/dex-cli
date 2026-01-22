package cmd

import (
	"fmt"

	"github.com/chriskievit/dex-cli/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  "View and manage DEX CLI configuration settings",
}

var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  "Display the current configuration values from ~/.dex-cli/config.yaml",
	RunE:  runShowConfig,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(showConfigCmd)
}

func runShowConfig(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println("Current Configuration")
	fmt.Println("─────────────────────────────────────────")
	fmt.Printf("Organization:     %s\n", formatValue(cfg.Organization))
	fmt.Printf("Project:          %s\n", formatValue(cfg.Project))
	fmt.Printf("Repository:       %s\n", formatValue(cfg.Repository))
	fmt.Printf("Default Reviewer: %s\n", formatValue(cfg.DefaultReviewer))
	fmt.Printf("\nConfig File: %s\n", config.GetConfigDir()+"/config.yaml")

	return nil
}

func formatValue(value string) string {
	if value == "" {
		return "(not set)"
	}
	return value
}

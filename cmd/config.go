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

var setConfigCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration values",
	Long:  "Set configuration values in ~/.dex-cli/config.yaml",
}

var setProjectCmd = &cobra.Command{
	Use:   "project [value]",
	Short: "Set the project configuration value",
	Long:  "Set the Azure DevOps project name in the configuration",
	Args:  cobra.ExactArgs(1),
	RunE:  runSetProject,
}

var setRepoCmd = &cobra.Command{
	Use:   "repo [value]",
	Short: "Set the repository configuration value",
	Long:  "Set the Azure DevOps repository name in the configuration",
	Args:  cobra.ExactArgs(1),
	RunE:  runSetRepo,
}

var setReviewerCmd = &cobra.Command{
	Use:   "reviewer [value]",
	Short: "Set the default reviewer configuration value",
	Long:  "Set the default reviewer for pull requests in the configuration",
	Args:  cobra.ExactArgs(1),
	RunE:  runSetReviewer,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(showConfigCmd)
	configCmd.AddCommand(setConfigCmd)
	setConfigCmd.AddCommand(setProjectCmd)
	setConfigCmd.AddCommand(setRepoCmd)
	setConfigCmd.AddCommand(setReviewerCmd)
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

func runSetProject(cmd *cobra.Command, args []string) error {
	value := args[0]
	if value == "" {
		return fmt.Errorf("project value cannot be empty")
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	cfg.Project = value

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Project set to: %s\n", value)
	return nil
}

func runSetRepo(cmd *cobra.Command, args []string) error {
	value := args[0]
	if value == "" {
		return fmt.Errorf("repository value cannot be empty")
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	cfg.Repository = value

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Repository set to: %s\n", value)
	return nil
}

func runSetReviewer(cmd *cobra.Command, args []string) error {
	value := args[0]
	if value == "" {
		return fmt.Errorf("reviewer value cannot be empty")
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	cfg.DefaultReviewer = value

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Default reviewer set to: %s\n", value)
	return nil
}

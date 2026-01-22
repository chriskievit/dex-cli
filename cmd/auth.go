package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/chriskievit/dex-cli/internal/auth"
	"github.com/chriskievit/dex-cli/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
	Long:  "Manage Azure DevOps authentication credentials stored securely in your system keychain",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Azure DevOps",
	Long:  "Store your Azure DevOps Personal Access Token (PAT) securely in the system keychain",
	RunE:  runLogin,
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from Azure DevOps",
	Long:  "Remove your Azure DevOps credentials from the system keychain",
	RunE:  runLogout,
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check authentication status",
	Long:  "Check if you are authenticated with Azure DevOps",
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(statusCmd)
}

func runLogin(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	// Get organization
	org := organization
	if org == "" {
		fmt.Print("Azure DevOps Organization: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read organization: %w", err)
		}
		org = strings.TrimSpace(input)
	}

	if org == "" {
		return fmt.Errorf("organization is required")
	}

	// Get PAT (hidden input)
	fmt.Print("Personal Access Token (PAT): ")
	tokenBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // New line after password input
	if err != nil {
		return fmt.Errorf("failed to read token: %w", err)
	}

	token := strings.TrimSpace(string(tokenBytes))
	if token == "" {
		return fmt.Errorf("token is required")
	}

	// Store token in keychain
	if err := auth.StoreToken(org, token); err != nil {
		return fmt.Errorf("failed to store credentials: %w", err)
	}

	// Update config with organization
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	cfg.Organization = org
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("✓ Successfully authenticated with organization: %s\n", org)
	fmt.Println("Your credentials are stored securely in the system keychain")

	return nil
}

func runLogout(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	org := organization
	if org == "" {
		org = cfg.Organization
	}

	if org == "" {
		return fmt.Errorf("no organization configured. Use --org flag or login first")
	}

	if err := auth.DeleteToken(org); err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	fmt.Printf("✓ Successfully logged out from organization: %s\n", org)

	return nil
}

func runStatus(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	org := organization
	if org == "" {
		org = cfg.Organization
	}

	if org == "" {
		fmt.Println("✗ Not authenticated")
		fmt.Println("Run 'dex-cli auth login' to authenticate")
		return nil
	}

	_, err = auth.GetToken(org, debug)
	if err != nil {
		fmt.Printf("✗ Not authenticated with organization: %s\n", org)
		fmt.Println("Run 'dex-cli auth login' to authenticate")
		return nil
	}

	fmt.Printf("✓ Authenticated with organization: %s\n", org)
	if cfg.Project != "" {
		fmt.Printf("  Default project: %s\n", cfg.Project)
	}
	if cfg.Repository != "" {
		fmt.Printf("  Default repository: %s\n", cfg.Repository)
	}

	return nil
}

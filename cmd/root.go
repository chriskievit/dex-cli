package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	organization string
	project      string
	debug        bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "dex",
	Short: "Azure DevOps CLI tool for managing branches, work items, and pull requests",
	Long: `DEX CLI is a secure command-line tool for Azure DevOps that helps you:
  - Create Git branches linked to work items
  - Create pull requests with automatic work item linking
  - View work item details
  
All credentials are stored securely in your system's keychain.`,
	Version: "1.0.0",
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&organization, "org", "o", "", "Azure DevOps organization")
	rootCmd.PersistentFlags().StringVarP(&project, "project", "p", "", "Azure DevOps project")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug output")
}

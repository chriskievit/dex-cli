package cmd

import (
	"fmt"

	"strconv"

	"github.com/chriskievit/dex-cli/internal/auth"
	"github.com/chriskievit/dex-cli/internal/azdo"
	"github.com/chriskievit/dex-cli/internal/config"
	"github.com/spf13/cobra"
)

var workitemCmd = &cobra.Command{
	Use:   "workitem",
	Short: "Manage work items",
	Long:  "View and manage Azure DevOps work items",
}

var showWorkitemCmd = &cobra.Command{
	Use:   "show <work-item-id>",
	Short: "Show work item details",
	Long:  "Display detailed information about an Azure DevOps work item",
	Args:  cobra.ExactArgs(1),
	RunE:  runShowWorkitem,
}

func init() {
	rootCmd.AddCommand(workitemCmd)
	workitemCmd.AddCommand(showWorkitemCmd)
}

func runShowWorkitem(cmd *cobra.Command, args []string) error {
	workItemIDStr := args[0]

	// Parse work item ID
	workItemID, err := strconv.Atoi(workItemIDStr)
	if err != nil {
		return fmt.Errorf("invalid work item ID: %s", workItemIDStr)
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	org := organization
	if org == "" {
		org = cfg.Organization
	}
	if org == "" {
		return fmt.Errorf("organization not configured. Use --org flag or run 'dex-cli auth login'")
	}

	// Get authentication token
	token, err := auth.GetToken(org, debug)
	if err != nil {
		return err
	}

	// Create Azure DevOps client
	client := azdo.NewClient(org, token, debug)

	// Fetch work item
	workItem, err := client.GetWorkItem(workItemID)
	if err != nil {
		return fmt.Errorf("failed to fetch work item: %w", err)
	}

	// Display work item details
	fmt.Printf("Work Item #%d\n", workItem.ID)
	fmt.Printf("─────────────────────────────────────────\n")
	fmt.Printf("Title:       %s\n", workItem.GetTitle())
	fmt.Printf("Type:        %s\n", workItem.GetWorkItemType())
	fmt.Printf("State:       %s\n", workItem.GetState())
	fmt.Printf("Assigned To: %s\n", workItem.GetAssignedTo())
	fmt.Printf("URL:         https://dev.azure.com/%s/_workitems/edit/%d\n", org, workItem.ID)

	return nil
}

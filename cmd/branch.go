package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/chriskievit/dex-cli/internal/auth"
	"github.com/chriskievit/dex-cli/internal/azdo"
	"github.com/chriskievit/dex-cli/internal/config"
	"github.com/chriskievit/dex-cli/internal/git"
	"github.com/spf13/cobra"
)

var (
	fromBranch string
)

var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "Manage Git branches",
	Long:  "Create and manage Git branches linked to Azure DevOps work items",
}

var createBranchCmd = &cobra.Command{
	Use:   "create <work-item-id> <description>",
	Short: "Create a new branch linked to a work item",
	Long: `Create a new Git branch following the naming convention: {work-item-type}/{work-item-id}/{description}

Example:
  dex-cli branch create 12345 add-login-feature
  dex-cli branch create 12345 fix-bug --from develop

The work item type will be fetched from Azure DevOps automatically.`,
	Args: cobra.ExactArgs(2),
	RunE: runCreateBranch,
}

func init() {
	rootCmd.AddCommand(branchCmd)
	branchCmd.AddCommand(createBranchCmd)

	createBranchCmd.Flags().StringVarP(&fromBranch, "from", "f", "", "Base branch to create from (defaults to main/master)")
}

func runCreateBranch(cmd *cobra.Command, args []string) error {
	workItemIDStr := args[0]
	description := args[1]

	// Parse work item ID
	workItemID, err := strconv.Atoi(workItemIDStr)
	if err != nil {
		return fmt.Errorf("invalid work item ID: %s", workItemIDStr)
	}

	// Validate description format (lowercase, hyphens only)
	if !isValidDescription(description) {
		return fmt.Errorf("description must be lowercase with hyphens only (e.g., add-login-feature)")
	}

	// Check if we're in a Git repository
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	if !git.IsGitRepository(cwd) {
		return fmt.Errorf("not a git repository. Please run this command from within a git repository")
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

	// Fetch work item to get the type
	if debug {
		fmt.Printf("Fetching work item %d...\n", workItemID)
	}

	workItem, err := client.GetWorkItem(workItemID)
	if err != nil {
		return fmt.Errorf("failed to fetch work item: %w", err)
	}

	workItemType := workItem.GetWorkItemType()
	workItemTitle := workItem.GetTitle()

	if debug {
		fmt.Printf("Work Item: %s #%d - %s\n", workItemType, workItemID, workItemTitle)
	}

	// Generate branch name: {type}/{id}/{description}
	branchName := fmt.Sprintf("%s/%d/%s", workItemType, workItemID, description)

	// Check if branch already exists
	exists, err := git.BranchExists(cwd, branchName)
	if err != nil {
		return fmt.Errorf("failed to check if branch exists: %w", err)
	}
	if exists {
		return fmt.Errorf("branch '%s' already exists", branchName)
	}

	// Determine base branch
	baseBranch := fromBranch
	if baseBranch == "" {
		baseBranch, err = git.GetDefaultBranch(cwd)
		if err != nil {
			return fmt.Errorf("failed to determine default branch: %w", err)
		}
	}

	// Create the branch
	fmt.Printf("Creating branch: %s (from %s)\n", branchName, baseBranch)
	if err := git.CreateBranch(cwd, branchName, baseBranch); err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	fmt.Printf("✓ Successfully created branch: %s\n", branchName)
	fmt.Printf("  Work Item: #%d - %s\n", workItemID, workItemTitle)
	fmt.Printf("  Type: %s\n", workItemType)

	return nil
}

func isValidDescription(desc string) bool {
	// Description should be lowercase letters, numbers, and hyphens only
	matched, _ := regexp.MatchString(`^[a-z0-9]+(-[a-z0-9]+)*$`, desc)
	return matched
}

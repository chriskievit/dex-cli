package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/chriskievit/dex-cli/internal/auth"
	"github.com/chriskievit/dex-cli/internal/azdo"
	"github.com/chriskievit/dex-cli/internal/config"
	"github.com/chriskievit/dex-cli/internal/git"
	"github.com/spf13/cobra"
)

var workitemCmd = &cobra.Command{
	Use:   "workitem",
	Short: "Manage work items",
	Long:  "View and manage Azure DevOps work items",
}

var (
	startBaseBranch string
)

var showWorkitemCmd = &cobra.Command{
	Use:   "show <work-item-id>",
	Short: "Show work item details",
	Long:  "Display detailed information about an Azure DevOps work item",
	Args:  cobra.ExactArgs(1),
	RunE:  runShowWorkitem,
}

var startWorkitemCmd = &cobra.Command{
	Use:   "start <work-item-id>",
	Short: "Start work on a work item",
	Long: `Start work on a work item by creating a new branch and linking it to the work item.

This command will:
  1. Verify the work item exists
  2. Optionally checkout a base branch (if --from is specified)
  3. Create a new branch from the current branch following the naming convention
  4. Push the branch to remote
  5. Link the branch to the work item via commit message

Example:
  dex workitem start 12345
  dex workitem start 12345 --from develop`,
	Args: cobra.ExactArgs(1),
	RunE: runStartWorkitem,
}

func init() {
	rootCmd.AddCommand(workitemCmd)
	workitemCmd.AddCommand(showWorkitemCmd)
	workitemCmd.AddCommand(startWorkitemCmd)

	startWorkitemCmd.Flags().StringVarP(&startBaseBranch, "from", "f", "", "Base branch to checkout before creating new branch")
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
		return fmt.Errorf("organization not configured. Use --org flag or run 'dex auth login'")
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

func runStartWorkitem(cmd *cobra.Command, args []string) error {
	// Check if we're in a Git repository first (fail fast)
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	if !git.IsGitRepository(cwd) {
		return fmt.Errorf("not a git repository. Please run this command from within a git repository")
	}

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
		return fmt.Errorf("organization not configured. Use --org flag or run 'dex auth login'")
	}

	// Get authentication token
	token, err := auth.GetToken(org, debug)
	if err != nil {
		return err
	}

	// Create Azure DevOps client
	client := azdo.NewClient(org, token, debug)

	// Step 1: Verify work item exists
	fmt.Printf("Step 1: Verifying work item #%d exists...\n", workItemID)
	workItem, err := client.GetWorkItem(workItemID)
	if err != nil {
		return fmt.Errorf("failed to fetch work item: %w", err)
	}
	fmt.Printf("✓ Work item found: %s #%d - %s\n", workItem.GetWorkItemType(), workItemID, workItem.GetTitle())

	// Step 2: Optionally checkout base branch
	currentBranch, err := git.GetCurrentBranch(cwd)
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	if startBaseBranch != "" {
		fmt.Printf("Step 2: Checking out base branch '%s'...\n", startBaseBranch)
		if err := git.CheckoutBranch(cwd, startBaseBranch); err != nil {
			return fmt.Errorf("failed to checkout base branch: %w", err)
		}
		fmt.Printf("✓ Checked out branch: %s\n", startBaseBranch)
		currentBranch = startBaseBranch
	} else {
		fmt.Printf("Step 2: Using current branch '%s' as base\n", currentBranch)
	}

	// Step 3: Generate branch name from work item
	workItemType := workItem.GetWorkItemType()
	workItemTitle := workItem.GetTitle()
	description := generateBranchDescription(workItemTitle)
	branchName := fmt.Sprintf("%s/%d/%s", workItemType, workItemID, description)

	// Check if branch already exists
	exists, err := git.BranchExists(cwd, branchName)
	if err != nil {
		return fmt.Errorf("failed to check if branch exists: %w", err)
	}
	if exists {
		return fmt.Errorf("branch '%s' already exists", branchName)
	}

	// Step 3: Create the branch
	fmt.Printf("Step 3: Creating branch '%s' from '%s'...\n", branchName, currentBranch)
	if err := git.CreateBranch(cwd, branchName, currentBranch); err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}
	fmt.Printf("✓ Branch created: %s\n", branchName)

	// Step 4: Create commit with work item reference to link it
	fmt.Printf("Step 4: Linking branch to work item #%d...\n", workItemID)
	commitMessage := fmt.Sprintf("Start work on #%d: %s", workItemID, workItemTitle)
	if err := git.CreateCommit(cwd, commitMessage); err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}
	fmt.Printf("✓ Commit created with work item reference\n")

	// Step 5: Push the branch
	fmt.Printf("Step 5: Pushing branch to remote...\n")
	if err := git.PushBranch(cwd, branchName); err != nil {
		return fmt.Errorf("failed to push branch: %w", err)
	}
	fmt.Printf("✓ Branch pushed to remote\n")

	fmt.Printf("\n✓ Successfully started work on work item #%d\n", workItemID)
	fmt.Printf("  Branch: %s\n", branchName)
	fmt.Printf("  Work Item: %s #%d - %s\n", workItemType, workItemID, workItemTitle)

	return nil
}

// generateBranchDescription converts a work item title to a valid branch description
// Converts to lowercase, replaces spaces/special chars with hyphens, removes invalid chars
func generateBranchDescription(title string) string {
	// Convert to lowercase
	desc := strings.ToLower(title)

	// Replace spaces and common separators with hyphens
	desc = regexp.MustCompile(`[\s_]+`).ReplaceAllString(desc, "-")

	// Remove all characters except lowercase letters, numbers, and hyphens
	desc = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(desc, "")

	// Remove multiple consecutive hyphens
	desc = regexp.MustCompile(`-+`).ReplaceAllString(desc, "-")

	// Remove leading/trailing hyphens
	desc = strings.Trim(desc, "-")

	// Limit length to 50 characters
	if len(desc) > 50 {
		desc = desc[:50]
		desc = strings.Trim(desc, "-")
	}

	// If empty after sanitization, use a default
	if desc == "" {
		desc = "work-item"
	}

	return desc
}

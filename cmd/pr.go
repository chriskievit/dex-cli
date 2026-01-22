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
	sourceBranch string
	targetBranch string
	prTitle      string
	prDesc       string
	workItemID   int
	isDraft      bool
)

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Manage pull requests",
	Long:  "Create and manage Azure DevOps pull requests",
}

var createPRCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new pull request",
	Long: `Create a new pull request in Azure DevOps.

The source branch defaults to your current Git branch.
Work item ID will be automatically extracted from the branch name if it follows the naming convention.

Example:
  dex-cli pr create --target main --title "Add login feature"
  dex-cli pr create --source feature/123/login --target main --title "Add login" --workitem 123`,
	RunE: runCreatePR,
}

func init() {
	rootCmd.AddCommand(prCmd)
	prCmd.AddCommand(createPRCmd)

	createPRCmd.Flags().StringVarP(&sourceBranch, "source", "s", "", "Source branch (defaults to current branch)")
	createPRCmd.Flags().StringVarP(&targetBranch, "target", "t", "", "Target branch (required)")
	createPRCmd.Flags().StringVar(&prTitle, "title", "", "Pull request title (required)")
	createPRCmd.Flags().StringVar(&prDesc, "description", "", "Pull request description")
	createPRCmd.Flags().IntVarP(&workItemID, "workitem", "w", 0, "Work item ID to link (auto-detected from branch name)")
	createPRCmd.Flags().BoolVar(&isDraft, "draft", false, "Create as draft pull request")

	createPRCmd.MarkFlagRequired("target")
	createPRCmd.MarkFlagRequired("title")
}

func runCreatePR(cmd *cobra.Command, args []string) error {
	// Check if we're in a Git repository
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	if !git.IsGitRepository(cwd) {
		return fmt.Errorf("not a git repository. Please run this command from within a git repository")
	}

	// Determine source branch
	source := sourceBranch
	if source == "" {
		source, err = git.GetCurrentBranch(cwd)
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}
		if debug {
			fmt.Printf("Using current branch as source: %s\n", source)
		}
	}

	// Validate source != target
	if source == targetBranch {
		return fmt.Errorf("source branch cannot be the same as target branch: %s", source)
	}

	// Extract work item ID from branch name if not provided
	wiID := workItemID
	if wiID == 0 {
		wiID = extractWorkItemFromBranch(source)
		if wiID > 0 && debug {
			fmt.Printf("Extracted work item ID from branch name: %d\n", wiID)
		}
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

	proj := project
	if proj == "" {
		proj = cfg.Project
	}
	if proj == "" {
		return fmt.Errorf("project not configured. Use --project flag or set in config")
	}

	repo := cfg.Repository
	if repo == "" {
		return fmt.Errorf("repository not configured. Please set in config file at %s", config.GetConfigDir())
	}

	// Get authentication token
	token, err := auth.GetToken(org, debug)
	if err != nil {
		return err
	}

	// Create Azure DevOps client
	client := azdo.NewClient(org, token, debug)

	// Get repository information
	if debug {
		fmt.Printf("Getting repository information for: %s\n", repo)
	}

	repository, err := client.GetRepository(proj, repo)
	if err != nil {
		return fmt.Errorf("failed to get repository: %w", err)
	}

	// Prepare PR request
	prRequest := &azdo.CreatePRRequest{
		SourceRefName: azdo.FormatRefName(source),
		TargetRefName: azdo.FormatRefName(targetBranch),
		Title:         prTitle,
		Description:   prDesc,
		IsDraft:       isDraft,
	}

	// Add work item link if available
	if wiID > 0 {
		prRequest.WorkItemRefs = []map[string]interface{}{
			{
				"id": strconv.Itoa(wiID),
			},
		}
	}

	// Create pull request
	fmt.Printf("Creating pull request...\n")
	fmt.Printf("  Source: %s\n", source)
	fmt.Printf("  Target: %s\n", targetBranch)
	fmt.Printf("  Title: %s\n", prTitle)
	if wiID > 0 {
		fmt.Printf("  Work Item: #%d\n", wiID)
	}

	pr, err := client.CreatePullRequest(proj, repository.ID, prRequest)
	if err != nil {
		return fmt.Errorf("failed to create pull request: %w", err)
	}

	fmt.Printf("\n✓ Successfully created pull request #%d\n", pr.PullRequestID)
	fmt.Printf("  URL: https://dev.azure.com/%s/%s/_git/%s/pullrequest/%d\n",
		org, proj, repo, pr.PullRequestID)

	return nil
}

// extractWorkItemFromBranch attempts to extract work item ID from branch name
// Expected format: {type}/{id}/{description}
func extractWorkItemFromBranch(branchName string) int {
	// Pattern: type/12345/description
	re := regexp.MustCompile(`^[a-z-]+/(\d+)/`)
	matches := re.FindStringSubmatch(branchName)
	if len(matches) > 1 {
		id, err := strconv.Atoi(matches[1])
		if err == nil {
			return id
		}
	}
	return 0
}

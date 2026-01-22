package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// IsGitRepository checks if the current directory is a git repository
func IsGitRepository(dir string) bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = dir
	err := cmd.Run()
	return err == nil
}

// GetCurrentBranch returns the name of the current branch
func GetCurrentBranch(dir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetDefaultBranch returns the default branch (main or master)
func GetDefaultBranch(dir string) (string, error) {
	// Try to get the default branch from remote
	cmd := exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err == nil {
		// Extract branch name from refs/remotes/origin/main
		branch := strings.TrimSpace(string(output))
		parts := strings.Split(branch, "/")
		if len(parts) > 0 {
			return parts[len(parts)-1], nil
		}
	}

	// Fallback: check if main exists
	cmd = exec.Command("git", "rev-parse", "--verify", "main")
	cmd.Dir = dir
	if err := cmd.Run(); err == nil {
		return "main", nil
	}

	// Fallback: check if master exists
	cmd = exec.Command("git", "rev-parse", "--verify", "master")
	cmd.Dir = dir
	if err := cmd.Run(); err == nil {
		return "master", nil
	}

	return "", fmt.Errorf("could not determine default branch")
}

// CreateBranch creates a new branch from the specified base branch
func CreateBranch(dir, branchName, baseBranch string) error {
	// First, ensure we're on the base branch or it exists
	cmd := exec.Command("git", "checkout", baseBranch)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to checkout base branch %s: %w", baseBranch, err)
	}

	// Create and checkout the new branch
	cmd = exec.Command("git", "checkout", "-b", branchName)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create branch %s: %w", branchName, err)
	}

	return nil
}

// GetRemoteURL returns the remote URL for origin
func GetRemoteURL(dir string) (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get remote URL: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// BranchExists checks if a branch exists locally
func BranchExists(dir, branchName string) (bool, error) {
	cmd := exec.Command("git", "rev-parse", "--verify", branchName)
	cmd.Dir = dir
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if branch exists: %w", err)
	}
	return true, nil
}

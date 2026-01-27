package testhelpers

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// SetupTempGitRepo creates a temporary git repository for testing
func SetupTempGitRepo(t *testing.T) (string, func()) {
	t.Helper()

	tmpDir := t.TempDir()

	// Initialize git repository with explicit initial branch
	cmd := exec.Command("git", "init", "--initial-branch=main")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		// Fallback to regular git init if --initial-branch is not supported
		cmd = exec.Command("git", "init")
		cmd.Dir = tmpDir
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to initialize git repository: %v", err)
		}
		// Rename branch to main if it was created as master
		cmd = exec.Command("git", "branch", "-m", "master", "main")
		cmd.Dir = tmpDir
		cmd.Run() // Ignore error if branch is already main
	}

	// Configure git user (required for commits)
	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.name: %v", err)
	}

	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to configure git user.email: %v", err)
	}

	// Create initial commit on main branch
	readmePath := filepath.Join(tmpDir, "README.md")
	if err := os.WriteFile(readmePath, []byte("# Test Repo\n"), 0644); err != nil {
		t.Fatalf("Failed to create README: %v", err)
	}

	cmd = exec.Command("git", "add", "README.md")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to add README: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	// Cleanup function
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

// CreateBranch creates a branch in the git repository
func CreateBranch(t *testing.T, repoDir, branchName string) {
	t.Helper()

	cmd := exec.Command("git", "checkout", "-b", branchName)
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create branch %s: %v", branchName, err)
	}
}

// CheckoutBranch checks out a branch in the git repository
func CheckoutBranch(t *testing.T, repoDir, branchName string) {
	t.Helper()

	cmd := exec.Command("git", "checkout", branchName)
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to checkout branch %s: %v", branchName, err)
	}
}

// CreateCommit creates a commit in the git repository
func CreateCommit(t *testing.T, repoDir, message string) {
	t.Helper()

	cmd := exec.Command("git", "commit", "--allow-empty", "-m", message)
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}
}

// GetCurrentBranch returns the current branch name
func GetCurrentBranch(t *testing.T, repoDir string) string {
	t.Helper()

	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = repoDir
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to get current branch: %v", err)
	}
	return string(output[:len(output)-1]) // Remove newline
}

package git

import (
	"os/exec"
	"testing"

	"github.com/chriskievit/dex-cli/internal/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsGitRepository(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T) string
		expected bool
	}{
		{
			name: "valid git repository",
			setup: func(t *testing.T) string {
				repoDir, _ := testhelpers.SetupTempGitRepo(t)
				return repoDir
			},
			expected: true,
		},
		{
			name: "not a git repository",
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			expected: false,
		},
		{
			name: "non-existent directory",
			setup: func(t *testing.T) string {
				return "/non/existent/path"
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setup(t)
			result := IsGitRepository(dir)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetCurrentBranch(t *testing.T) {
	repoDir, cleanup := testhelpers.SetupTempGitRepo(t)
	defer cleanup()

	branch, err := GetCurrentBranch(repoDir)
	require.NoError(t, err)
	// Git version determines default branch name (main or master)
	assert.Contains(t, []string{"main", "master"}, branch)
}

func TestGetCurrentBranch_NotGitRepo(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := GetCurrentBranch(tmpDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get current branch")
}

func TestGetDefaultBranch_Main(t *testing.T) {
	repoDir, cleanup := testhelpers.SetupTempGitRepo(t)
	defer cleanup()

	// Get current branch (could be main or master)
	currentBranch := testhelpers.GetCurrentBranch(t, repoDir)

	branch, err := GetDefaultBranch(repoDir)
	require.NoError(t, err)
	// Should return the current default branch
	assert.Contains(t, []string{"main", "master"}, branch)
	// Should match current branch if it's a default branch
	if currentBranch == "main" || currentBranch == "master" {
		assert.Equal(t, currentBranch, branch)
	}
}

func TestGetDefaultBranch_Master(t *testing.T) {
	repoDir, cleanup := testhelpers.SetupTempGitRepo(t)
	defer cleanup()

	// Get current branch name
	currentBranch := testhelpers.GetCurrentBranch(t, repoDir)

	// Rename current branch to master if it's not already
	if currentBranch != "master" {
		cmd := exec.Command("git", "branch", "-m", currentBranch, "master")
		cmd.Dir = repoDir
		require.NoError(t, cmd.Run())
	}

	cmd := exec.Command("git", "checkout", "master")
	cmd.Dir = repoDir
	require.NoError(t, cmd.Run())

	branch, err := GetDefaultBranch(repoDir)
	require.NoError(t, err)
	assert.Equal(t, "master", branch)
}

func TestGetDefaultBranch_NoDefaultBranch(t *testing.T) {
	repoDir, cleanup := testhelpers.SetupTempGitRepo(t)
	defer cleanup()

	// Get current branch name
	currentBranch := testhelpers.GetCurrentBranch(t, repoDir)

	// Create a different branch
	cmd := exec.Command("git", "checkout", "-b", "develop")
	cmd.Dir = repoDir
	require.NoError(t, cmd.Run())

	// Delete the default branch (main or master)
	cmd = exec.Command("git", "branch", "-D", currentBranch)
	cmd.Dir = repoDir
	cmd.Run() // Ignore error if already deleted

	// Try to get default branch - should fail
	_, err := GetDefaultBranch(repoDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not determine default branch")
}

func TestCreateBranch(t *testing.T) {
	repoDir, cleanup := testhelpers.SetupTempGitRepo(t)
	defer cleanup()

	// Get current branch (main or master)
	currentBranch := testhelpers.GetCurrentBranch(t, repoDir)

	branchName := "feature/test-branch"
	err := CreateBranch(repoDir, branchName, currentBranch)
	require.NoError(t, err)

	// Verify branch was created
	cmd := exec.Command("git", "rev-parse", "--verify", branchName)
	cmd.Dir = repoDir
	err = cmd.Run()
	assert.NoError(t, err, "Branch should exist")
}

func TestCreateBranch_InvalidBaseBranch(t *testing.T) {
	repoDir, cleanup := testhelpers.SetupTempGitRepo(t)
	defer cleanup()

	branchName := "feature/test-branch"
	err := CreateBranch(repoDir, branchName, "non-existent-base")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to checkout base branch")
}

func TestBranchExists(t *testing.T) {
	repoDir, cleanup := testhelpers.SetupTempGitRepo(t)
	defer cleanup()

	// Create a branch
	testhelpers.CreateBranch(t, repoDir, "test-branch")

	tests := []struct {
		name     string
		branch   string
		expected bool
	}{
		{
			name:     "branch exists",
			branch:   "test-branch",
			expected: true,
		},
		{
			name:     "branch does not exist",
			branch:   "non-existent-branch",
			expected: false,
		},
		{
			name:     "default branch exists",
			branch:   testhelpers.GetCurrentBranch(t, repoDir),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := BranchExists(repoDir, tt.branch)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, exists)
		})
	}
}

func TestBranchExists_NotGitRepo(t *testing.T) {
	tmpDir := t.TempDir()

	// BranchExists may return false without error for non-git directories
	// depending on how git rev-parse behaves
	exists, err := BranchExists(tmpDir, "main")
	// Either should return false, or an error
	if err != nil {
		assert.Contains(t, err.Error(), "failed to check if branch exists")
	} else {
		assert.False(t, exists)
	}
}

func TestCheckoutBranch(t *testing.T) {
	repoDir, cleanup := testhelpers.SetupTempGitRepo(t)
	defer cleanup()

	// Create a branch
	testhelpers.CreateBranch(t, repoDir, "test-branch")

	// Checkout the branch
	err := CheckoutBranch(repoDir, "test-branch")
	require.NoError(t, err)

	// Verify we're on the branch
	branch := testhelpers.GetCurrentBranch(t, repoDir)
	assert.Equal(t, "test-branch", branch)
}

func TestCheckoutBranch_NonExistentBranch(t *testing.T) {
	repoDir, cleanup := testhelpers.SetupTempGitRepo(t)
	defer cleanup()

	err := CheckoutBranch(repoDir, "non-existent-branch")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to checkout branch")
}

func TestCreateCommit(t *testing.T) {
	repoDir, cleanup := testhelpers.SetupTempGitRepo(t)
	defer cleanup()

	message := "Test commit message"
	err := CreateCommit(repoDir, message)
	require.NoError(t, err)

	// Verify commit was created
	cmd := exec.Command("git", "log", "-1", "--pretty=%s")
	cmd.Dir = repoDir
	output, err := cmd.Output()
	require.NoError(t, err)
	assert.Contains(t, string(output), message)
}

func TestCreateCommit_NotGitRepo(t *testing.T) {
	tmpDir := t.TempDir()

	err := CreateCommit(tmpDir, "Test message")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create commit")
}

func TestGetRemoteURL(t *testing.T) {
	repoDir, cleanup := testhelpers.SetupTempGitRepo(t)
	defer cleanup()

	// Add a remote
	cmd := exec.Command("git", "remote", "add", "origin", "https://github.com/test/repo.git")
	cmd.Dir = repoDir
	require.NoError(t, cmd.Run())

	url, err := GetRemoteURL(repoDir)
	require.NoError(t, err)
	assert.Equal(t, "https://github.com/test/repo.git", url)
}

func TestGetRemoteURL_NoRemote(t *testing.T) {
	repoDir, cleanup := testhelpers.SetupTempGitRepo(t)
	defer cleanup()

	_, err := GetRemoteURL(repoDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get remote URL")
}

func TestPushBranch(t *testing.T) {
	// This test requires a remote repository, which is complex to set up
	// For now, we'll test that it fails gracefully when there's no remote
	repoDir, cleanup := testhelpers.SetupTempGitRepo(t)
	defer cleanup()

	// Create a branch
	testhelpers.CreateBranch(t, repoDir, "test-branch")

	// Push should fail without a remote
	err := PushBranch(repoDir, "test-branch")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to push branch")
}

func TestCreateBranch_FromCurrentBranch(t *testing.T) {
	repoDir, cleanup := testhelpers.SetupTempGitRepo(t)
	defer cleanup()

	// Create a feature branch
	testhelpers.CreateBranch(t, repoDir, "feature/base")

	// Create another branch from the feature branch
	branchName := "feature/child"
	err := CreateBranch(repoDir, branchName, "feature/base")
	require.NoError(t, err)

	// Verify branch was created
	cmd := exec.Command("git", "rev-parse", "--verify", branchName)
	cmd.Dir = repoDir
	err = cmd.Run()
	assert.NoError(t, err, "Branch should exist")

	// Verify it's based on feature/base
	cmd = exec.Command("git", "log", "--oneline", "feature/base.."+branchName)
	cmd.Dir = repoDir
	output, err := cmd.Output()
	require.NoError(t, err)
	// If output is empty, it means the branches point to the same commit
	assert.Empty(t, output, "New branch should be based on feature/base")
}

func TestGetDefaultBranch_RemoteHEAD(t *testing.T) {
	repoDir, cleanup := testhelpers.SetupTempGitRepo(t)
	defer cleanup()

	// Set up remote tracking
	cmd := exec.Command("git", "branch", "--set-upstream-to=origin/main", "main")
	cmd.Dir = repoDir
	cmd.Run() // May fail if remote doesn't exist, that's okay

	// Set symbolic ref for remote HEAD
	cmd = exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main")
	cmd.Dir = repoDir
	cmd.Run() // May fail, that's okay - we'll fall back to main/master check

	// Should still work with fallback
	branch, err := GetDefaultBranch(repoDir)
	require.NoError(t, err)
	assert.Contains(t, []string{"main", "master"}, branch)
}

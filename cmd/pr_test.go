package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractWorkItemFromBranch(t *testing.T) {
	tests := []struct {
		name     string
		branch   string
		expected int
	}{
		{
			name:     "valid format with user-story",
			branch:   "user-story/12345/add-login-feature",
			expected: 12345,
		},
		{
			name:     "valid format with bug",
			branch:   "bug/67890/fix-crash",
			expected: 67890,
		},
		{
			name:     "valid format with task",
			branch:   "task/11111/implement-feature",
			expected: 11111,
		},
		{
			name:     "valid format with feature",
			branch:   "feature/22222/new-ui",
			expected: 22222,
		},
		{
			name:     "valid format with epic",
			branch:   "epic/33333/major-update",
			expected: 33333,
		},
		{
			name:     "valid format with single digit ID",
			branch:   "bug/1/fix",
			expected: 1,
		},
		{
			name:     "valid format with long ID",
			branch:   "user-story/1234567890/long-description",
			expected: 1234567890,
		},
		{
			name:     "valid format with hyphenated type",
			branch:   "user-story/12345/test-branch",
			expected: 12345,
		},
		{
			name:     "invalid no work item ID",
			branch:   "feature/add-login",
			expected: 0,
		},
		{
			name:     "invalid wrong format",
			branch:   "feature12345add-login",
			expected: 0,
		},
		{
			name:     "invalid empty string",
			branch:   "",
			expected: 0,
		},
		{
			name:     "invalid just type",
			branch:   "feature",
			expected: 0,
		},
		{
			name:     "invalid missing description",
			branch:   "feature/12345",
			expected: 0, // Regex requires trailing slash
		},
		{
			name:     "invalid non-numeric ID",
			branch:   "feature/abc123/description",
			expected: 0,
		},
		{
			name:     "invalid ID with letters",
			branch:   "feature/123abc/description",
			expected: 0,
		},
		{
			name:     "valid with refs/heads prefix",
			branch:   "refs/heads/user-story/12345/add-login",
			expected: 0, // Pattern doesn't match refs/heads prefix
		},
		{
			name:     "valid format with multiple segments",
			branch:   "user-story/12345/add-login-feature-with-oauth",
			expected: 12345,
		},
		{
			name:     "invalid starts with number",
			branch:   "12345/feature/add-login",
			expected: 0,
		},
		{
			name:     "invalid with uppercase type",
			branch:   "Feature/12345/add-login",
			expected: 0, // Regex only matches lowercase
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractWorkItemFromBranch(tt.branch)
			assert.Equal(t, tt.expected, result, "Branch: %q", tt.branch)
		})
	}
}

func TestLoadPRTemplate(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(t *testing.T) string
		expectedContent string
		expectedError  bool
	}{
		{
			name: "template in .azuredevops directory",
			setup: func(t *testing.T) string {
				tmpDir := t.TempDir()
				azureDevOpsDir := filepath.Join(tmpDir, ".azuredevops")
				err := os.MkdirAll(azureDevOpsDir, 0755)
				require.NoError(t, err)
				
				templatePath := filepath.Join(azureDevOpsDir, "pull_request_template.md")
				content := "# PR Template\n\nThis is a test template."
				err = os.WriteFile(templatePath, []byte(content), 0644)
				require.NoError(t, err)
				
				return tmpDir
			},
			expectedContent: "# PR Template\n\nThis is a test template.",
			expectedError:   false,
		},
		{
			name: "template in .github directory",
			setup: func(t *testing.T) string {
				tmpDir := t.TempDir()
				githubDir := filepath.Join(tmpDir, ".github")
				err := os.MkdirAll(githubDir, 0755)
				require.NoError(t, err)
				
				templatePath := filepath.Join(githubDir, "pull_request_template.md")
				content := "# GitHub Style Template\n\nDescription here."
				err = os.WriteFile(templatePath, []byte(content), 0644)
				require.NoError(t, err)
				
				return tmpDir
			},
			expectedContent: "# GitHub Style Template\n\nDescription here.",
			expectedError:   false,
		},
		{
			name: "template in repository root",
			setup: func(t *testing.T) string {
				tmpDir := t.TempDir()
				
				templatePath := filepath.Join(tmpDir, "pull_request_template.md")
				content := "# Root Template\n\nRoot level template."
				err := os.WriteFile(templatePath, []byte(content), 0644)
				require.NoError(t, err)
				
				return tmpDir
			},
			expectedContent: "# Root Template\n\nRoot level template.",
			expectedError:   false,
		},
		{
			name: "uppercase template in .azuredevops",
			setup: func(t *testing.T) string {
				tmpDir := t.TempDir()
				azureDevOpsDir := filepath.Join(tmpDir, ".azuredevops")
				err := os.MkdirAll(azureDevOpsDir, 0755)
				require.NoError(t, err)
				
				templatePath := filepath.Join(azureDevOpsDir, "PULL_REQUEST_TEMPLATE.md")
				content := "# Uppercase Template\n\nUppercase variant."
				err = os.WriteFile(templatePath, []byte(content), 0644)
				require.NoError(t, err)
				
				return tmpDir
			},
			expectedContent: "# Uppercase Template\n\nUppercase variant.",
			expectedError:   false,
		},
		{
			name: "prefers .azuredevops over .github",
			setup: func(t *testing.T) string {
				tmpDir := t.TempDir()
				
				// Create both templates
				azureDevOpsDir := filepath.Join(tmpDir, ".azuredevops")
				err := os.MkdirAll(azureDevOpsDir, 0755)
				require.NoError(t, err)
				azureTemplate := filepath.Join(azureDevOpsDir, "pull_request_template.md")
				err = os.WriteFile(azureTemplate, []byte("# Azure DevOps Template"), 0644)
				require.NoError(t, err)
				
				githubDir := filepath.Join(tmpDir, ".github")
				err = os.MkdirAll(githubDir, 0755)
				require.NoError(t, err)
				githubTemplate := filepath.Join(githubDir, "pull_request_template.md")
				err = os.WriteFile(githubTemplate, []byte("# GitHub Template"), 0644)
				require.NoError(t, err)
				
				return tmpDir
			},
			expectedContent: "# Azure DevOps Template",
			expectedError:   false,
		},
		{
			name: "prefers .github over root",
			setup: func(t *testing.T) string {
				tmpDir := t.TempDir()
				
				// Create .github template
				githubDir := filepath.Join(tmpDir, ".github")
				err := os.MkdirAll(githubDir, 0755)
				require.NoError(t, err)
				githubTemplate := filepath.Join(githubDir, "pull_request_template.md")
				err = os.WriteFile(githubTemplate, []byte("# GitHub Template"), 0644)
				require.NoError(t, err)
				
				// Create root template
				rootTemplate := filepath.Join(tmpDir, "pull_request_template.md")
				err = os.WriteFile(rootTemplate, []byte("# Root Template"), 0644)
				require.NoError(t, err)
				
				return tmpDir
			},
			expectedContent: "# GitHub Template",
			expectedError:   false,
		},
		{
			name: "no template found",
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			expectedContent: "",
			expectedError:   true,
		},
		{
			name: "template with special characters",
			setup: func(t *testing.T) string {
				tmpDir := t.TempDir()
				azureDevOpsDir := filepath.Join(tmpDir, ".azuredevops")
				err := os.MkdirAll(azureDevOpsDir, 0755)
				require.NoError(t, err)
				
				templatePath := filepath.Join(azureDevOpsDir, "pull_request_template.md")
				content := "# Template\n\n- Item 1\n- Item 2\n\n**Bold text** and *italic*.\n\n```\ncode block\n```"
				err = os.WriteFile(templatePath, []byte(content), 0644)
				require.NoError(t, err)
				
				return tmpDir
			},
			expectedContent: "# Template\n\n- Item 1\n- Item 2\n\n**Bold text** and *italic*.\n\n```\ncode block\n```",
			expectedError:   false,
		},
		{
			name: "template with unicode characters",
			setup: func(t *testing.T) string {
				tmpDir := t.TempDir()
				azureDevOpsDir := filepath.Join(tmpDir, ".azuredevops")
				err := os.MkdirAll(azureDevOpsDir, 0755)
				require.NoError(t, err)
				
				templatePath := filepath.Join(azureDevOpsDir, "pull_request_template.md")
				content := "# Template with Unicode\n\nCafé résumé naïve"
				err = os.WriteFile(templatePath, []byte(content), 0644)
				require.NoError(t, err)
				
				return tmpDir
			},
			expectedContent: "# Template with Unicode\n\nCafé résumé naïve",
			expectedError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoDir := tt.setup(t)
			
			content, err := loadPRTemplate(repoDir)
			
			if tt.expectedError {
				assert.Error(t, err)
				assert.Empty(t, content)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedContent, content)
			}
		})
	}
}

package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

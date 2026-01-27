package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateBranchDescription(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple title",
			input:    "Add login feature",
			expected: "add-login-feature",
		},
		{
			name:     "title with underscores",
			input:    "Add_login_feature",
			expected: "add-login-feature",
		},
		{
			name:     "title with multiple spaces",
			input:    "Add  login   feature",
			expected: "add-login-feature",
		},
		{
			name:     "title with special characters",
			input:    "Add login feature!",
			expected: "add-login-feature",
		},
		{
			name:     "title with numbers",
			input:    "Fix bug #123",
			expected: "fix-bug-123",
		},
		{
			name:     "title with mixed case",
			input:    "Implement New Authentication System",
			expected: "implement-new-authentication-system",
		},
		{
			name:     "title with punctuation",
			input:    "Fix bug: Crash on startup",
			expected: "fix-bug-crash-on-startup",
		},
		{
			name:     "title with parentheses",
			input:    "Add feature (WIP)",
			expected: "add-feature-wip",
		},
		{
			name:     "title with brackets",
			input:    "Fix [Critical] bug",
			expected: "fix-critical-bug",
		},
		{
			name:     "title with quotes",
			input:    `Add "new" feature`,
			expected: "add-new-feature",
		},
		{
			name:     "title with slashes",
			input:    "Add feature/login",
			expected: "add-featurelogin",
		},
		{
			name:     "title with backslashes",
			input:    "Add feature\\login",
			expected: "add-featurelogin",
		},
		{
			name:     "title with multiple consecutive hyphens in result",
			input:    "Add---login---feature",
			expected: "add-login-feature",
		},
		{
			name:     "title starting with special characters",
			input:    "!!!Add login feature",
			expected: "add-login-feature",
		},
		{
			name:     "title ending with special characters",
			input:    "Add login feature!!!",
			expected: "add-login-feature",
		},
		{
			name:     "title with leading/trailing spaces",
			input:    "  Add login feature  ",
			expected: "add-login-feature",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "work-item",
		},
		{
			name:     "only special characters",
			input:    "!!!@@@###",
			expected: "work-item",
		},
		{
			name:     "only spaces",
			input:    "     ",
			expected: "work-item",
		},
		{
			name:     "long title over 50 characters",
			input:    "This is a very long work item title that exceeds fifty characters in length",
			expected: "this-is-a-very-long-work-item-title-that-exceeds-f",
		},
		{
			name:     "long title exactly 50 characters",
			input:    "This is a work item title that is exactly fifty",
			expected: "this-is-a-work-item-title-that-is-exactly-fifty",
		},
		{
			name:     "long title with special characters",
			input:    "This is a very long work item title!!! that exceeds fifty characters",
			expected: "this-is-a-very-long-work-item-title-that-exceeds-f",
		},
		{
			name:     "title with unicode characters",
			input:    "Add café feature",
			expected: "add-caf-feature",
		},
		{
			name:     "title with emoji",
			input:    "Add 🚀 feature",
			expected: "add-feature",
		},
		{
			name:     "title with tabs",
			input:    "Add\tlogin\tfeature",
			expected: "add-login-feature",
		},
		{
			name:     "title with newlines",
			input:    "Add\nlogin\nfeature",
			expected: "add-login-feature",
		},
		{
			name:     "title with only numbers",
			input:    "12345",
			expected: "12345",
		},
		{
			name:     "title with hyphens",
			input:    "Add-login-feature",
			expected: "add-login-feature",
		},
		{
			name:     "title with multiple words and numbers",
			input:    "Fix bug 123 in module 456",
			expected: "fix-bug-123-in-module-456",
		},
		{
			name:     "title that becomes empty after sanitization",
			input:    "!!!",
			expected: "work-item",
		},
		{
			name:     "title with leading hyphen after sanitization",
			input:    "-Add login",
			expected: "add-login",
		},
		{
			name:     "title with trailing hyphen after sanitization",
			input:    "Add login-",
			expected: "add-login",
		},
		{
			name:     "title truncated at hyphen boundary",
			input:    "This is a very long work item title that exceeds fifty characters and should be truncated",
			expected: "this-is-a-very-long-work-item-title-that-exceeds-f",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateBranchDescription(tt.input)
			assert.Equal(t, tt.expected, result, "Input: %q", tt.input)
		})
	}
}

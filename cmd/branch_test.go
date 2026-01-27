package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidDescription(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid lowercase with hyphens",
			input:    "add-login-feature",
			expected: true,
		},
		{
			name:     "valid single word",
			input:    "feature",
			expected: true,
		},
		{
			name:     "valid with numbers",
			input:    "fix-bug-123",
			expected: true,
		},
		{
			name:     "valid all numbers",
			input:    "12345",
			expected: true,
		},
		{
			name:     "valid single character",
			input:    "a",
			expected: true,
		},
		{
			name:     "invalid uppercase letters",
			input:    "Add-Login-Feature",
			expected: false,
		},
		{
			name:     "invalid with spaces",
			input:    "add login feature",
			expected: false,
		},
		{
			name:     "invalid with underscores",
			input:    "add_login_feature",
			expected: false,
		},
		{
			name:     "invalid with special characters",
			input:    "add-login-feature!",
			expected: false,
		},
		{
			name:     "invalid empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "invalid starts with hyphen",
			input:    "-add-login-feature",
			expected: false,
		},
		{
			name:     "invalid ends with hyphen",
			input:    "add-login-feature-",
			expected: false,
		},
		{
			name:     "invalid multiple consecutive hyphens",
			input:    "add--login--feature",
			expected: false,
		},
		{
			name:     "valid with single hyphens",
			input:    "add-login-feature",
			expected: true,
		},
		{
			name:     "invalid with dots",
			input:    "add.login.feature",
			expected: false,
		},
		{
			name:     "invalid with slashes",
			input:    "add/login/feature",
			expected: false,
		},
		{
			name:     "valid long description",
			input:    "implement-new-authentication-system-with-oauth",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidDescription(tt.input)
			assert.Equal(t, tt.expected, result, "Input: %q", tt.input)
		})
	}
}

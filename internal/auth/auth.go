package auth

import (
	"fmt"
	"strings"

	"github.com/zalando/go-keyring"
)

const (
	serviceName = "dex-cli"
)

// StoreToken securely stores the PAT in the system keychain
func StoreToken(organization, token string) error {
	if organization == "" {
		return fmt.Errorf("organization cannot be empty")
	}
	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	// Validate token format (basic check)
	if len(token) < 20 {
		return fmt.Errorf("invalid token format")
	}

	// Use organization URL as the username for keyring
	username := normalizeOrganization(organization)

	if err := keyring.Set(serviceName, username, token); err != nil {
		return fmt.Errorf("failed to store token in keychain: %w", err)
	}

	return nil
}

// GetToken retrieves the PAT from the system keychain
func GetToken(organization string) (string, error) {
	if organization == "" {
		return "", fmt.Errorf("organization cannot be empty")
	}

	username := normalizeOrganization(organization)

	token, err := keyring.Get(serviceName, username)
	if err != nil {
		if err == keyring.ErrNotFound {
			return "", fmt.Errorf("no credentials found for organization %s. Please run 'dex-cli auth login' first", organization)
		}
		return "", fmt.Errorf("failed to retrieve token from keychain: %w", err)
	}

	return token, nil
}

// DeleteToken removes the PAT from the system keychain
func DeleteToken(organization string) error {
	if organization == "" {
		return fmt.Errorf("organization cannot be empty")
	}

	username := normalizeOrganization(organization)

	if err := keyring.Delete(serviceName, username); err != nil {
		if err == keyring.ErrNotFound {
			return fmt.Errorf("no credentials found for organization %s", organization)
		}
		return fmt.Errorf("failed to delete token from keychain: %w", err)
	}

	return nil
}

// normalizeOrganization ensures consistent organization format
func normalizeOrganization(org string) string {
	// Remove any protocol prefix
	org = strings.TrimPrefix(org, "https://")
	org = strings.TrimPrefix(org, "http://")

	// Remove trailing slashes
	org = strings.TrimSuffix(org, "/")

	// If it's just the org name, return as-is
	// If it's a full URL, extract the org name
	if strings.Contains(org, "dev.azure.com/") {
		parts := strings.Split(org, "dev.azure.com/")
		if len(parts) > 1 {
			org = strings.TrimSuffix(parts[1], "/")
		}
	}

	return org
}

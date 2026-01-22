package azdo

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	apiVersion = "7.0"
	baseURL    = "https://dev.azure.com"
)

// Client represents an Azure DevOps API client
type Client struct {
	organization string
	token        string
	httpClient   *http.Client
	debug        bool
}

// NewClient creates a new Azure DevOps API client
func NewClient(organization, token string, debug bool) *Client {
	// Normalize organization name - extract just the org name from URL if needed
	org := normalizeOrganization(organization)

	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}

	c := &Client{
		organization: org,
		token:        token,
		httpClient:   client,
		debug:        debug,
	}

	c.debugLog("[DEBUG] Organization (input): %q\n", organization)
	c.debugLog("[DEBUG] Organization (normalized): %q\n", org)

	return c
}

// normalizeOrganization extracts the organization name from a full URL or returns it as-is
func normalizeOrganization(org string) string {
	// Remove any protocol prefix
	org = strings.TrimPrefix(org, "https://")
	org = strings.TrimPrefix(org, "http://")

	// Remove trailing slashes
	org = strings.TrimSuffix(org, "/")

	// If it's a full URL, extract the org name
	if strings.Contains(org, "dev.azure.com/") {
		parts := strings.Split(org, "dev.azure.com/")
		if len(parts) > 1 {
			org = strings.TrimSuffix(parts[1], "/")
			// Remove any path after the org name
			if idx := strings.Index(org, "/"); idx != -1 {
				org = org[:idx]
			}
		}
	}

	return org
}

// debugLog conditionally prints debug messages if debug is enabled
func (c *Client) debugLog(format string, args ...interface{}) {
	if c.debug {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}

// doRequest performs an HTTP request with authentication
func (c *Client) doRequest(method, url string, body interface{}) ([]byte, error) {
	c.debugLog("[DEBUG] HTTP Request: %s %s\n", method, url)

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authentication header (Basic Auth with PAT)
	// Format: "{username}:{PAT}" base64 encoded
	// Username is empty for Azure DevOps PAT authentication
	authString := fmt.Sprintf(":%s", c.token)
	authEncoded := base64.StdEncoding.EncodeToString([]byte(authString))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authEncoded))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	c.debugLog("[DEBUG] Authorization header set (Basic Auth with empty username): %s\n", authEncoded)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	c.debugLog("[DEBUG] HTTP Response Status: %d %s\n", resp.StatusCode, resp.Status)
	c.debugLog("[DEBUG] Response Headers: %v\n", resp.Header)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	c.debugLog("[DEBUG] Response Body (first 500 chars): %s\n", truncateString(string(respBody), 500))
	if len(respBody) > 500 {
		c.debugLog("[DEBUG] Response Body length: %d bytes\n", len(respBody))
	}

	// Accept 2xx status codes (including 203 Non-Authoritative Information)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		c.debugLog("[DEBUG] Full Response Body: %s\n", string(respBody))
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// buildURL constructs the full API URL with proper path encoding
func (c *Client) buildURL(project, path string) string {
	orgEncoded := url.PathEscape(c.organization)
	if project != "" {
		projectEncoded := url.PathEscape(project)
		return fmt.Sprintf("%s/%s/%s/_apis/%s?api-version=%s", baseURL, orgEncoded, projectEncoded, path, apiVersion)
	}
	return fmt.Sprintf("%s/%s/_apis/%s?api-version=%s", baseURL, orgEncoded, path, apiVersion)
}

// truncateString truncates a string to a maximum length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

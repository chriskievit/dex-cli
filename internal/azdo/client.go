package azdo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
}

// NewClient creates a new Azure DevOps API client
func NewClient(organization, token string) *Client {
	return &Client{
		organization: organization,
		token:        token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequest performs an HTTP request with authentication
func (c *Client) doRequest(method, url string, body interface{}) ([]byte, error) {
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
	req.SetBasicAuth("", c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// buildURL constructs the full API URL
func (c *Client) buildURL(project, path string) string {
	if project != "" {
		return fmt.Sprintf("%s/%s/%s/_apis/%s?api-version=%s", baseURL, c.organization, project, path, apiVersion)
	}
	return fmt.Sprintf("%s/%s/_apis/%s?api-version=%s", baseURL, c.organization, path, apiVersion)
}

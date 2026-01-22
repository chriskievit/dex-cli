package azdo

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// Repository represents an Azure DevOps Git repository
type Repository struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// PullRequest represents an Azure DevOps pull request
type PullRequest struct {
	PullRequestID int    `json:"pullRequestId"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	SourceRefName string `json:"sourceRefName"`
	TargetRefName string `json:"targetRefName"`
	Status        string `json:"status"`
}

// CreatePRRequest represents the request body for creating a pull request
type CreatePRRequest struct {
	SourceRefName string                   `json:"sourceRefName"`
	TargetRefName string                   `json:"targetRefName"`
	Title         string                   `json:"title"`
	Description   string                   `json:"description,omitempty"`
	WorkItemRefs  []map[string]interface{} `json:"workItemRefs,omitempty"`
	Reviewers     []map[string]interface{} `json:"reviewers,omitempty"`
	IsDraft       bool                     `json:"isDraft,omitempty"`
}

// GetRepository retrieves repository information by name
func (c *Client) GetRepository(project, repoName string) (*Repository, error) {
	repoNameEncoded := url.PathEscape(repoName)
	apiURL := c.buildURL(project, fmt.Sprintf("git/repositories/%s", repoNameEncoded))

	respBody, err := c.doRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}

	var repo Repository
	if err := json.Unmarshal(respBody, &repo); err != nil {
		return nil, fmt.Errorf("failed to parse repository response: %w", err)
	}

	return &repo, nil
}

// CreatePullRequest creates a new pull request
func (c *Client) CreatePullRequest(project, repoID string, req *CreatePRRequest) (*PullRequest, error) {
	repoIDEncoded := url.PathEscape(repoID)
	apiURL := c.buildURL(project, fmt.Sprintf("git/repositories/%s/pullrequests", repoIDEncoded))

	respBody, err := c.doRequest("POST", apiURL, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create pull request: %w", err)
	}

	var pr PullRequest
	if err := json.Unmarshal(respBody, &pr); err != nil {
		return nil, fmt.Errorf("failed to parse pull request response: %w", err)
	}

	return &pr, nil
}

// FormatRefName formats a branch name to a full ref name
func FormatRefName(branchName string) string {
	if !hasPrefix(branchName, "refs/") {
		return "refs/heads/" + branchName
	}
	return branchName
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

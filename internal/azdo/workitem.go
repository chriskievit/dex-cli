package azdo

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// WorkItem represents an Azure DevOps work item
type WorkItem struct {
	ID     int                    `json:"id"`
	Fields map[string]interface{} `json:"fields"`
}

// GetWorkItem retrieves a work item by ID
func (c *Client) GetWorkItem(id int) (*WorkItem, error) {
	// Properly encode the organization name in the URL path
	orgEncoded := url.PathEscape(c.organization)
	apiURL := fmt.Sprintf("%s/%s/_apis/wit/workitems/%d?api-version=%s", baseURL, orgEncoded, id, apiVersion)

	c.debugLog("[DEBUG] Organization (original): %q\n", c.organization)
	c.debugLog("[DEBUG] Organization (encoded): %q\n", orgEncoded)
	c.debugLog("[DEBUG] Constructed URL: %s\n", apiURL)

	respBody, err := c.doRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get work item: %w", err)
	}

	c.debugLog("[DEBUG] Attempting to unmarshal JSON response (length: %d bytes)\n", len(respBody))
	
	var workItem WorkItem
	if err := json.Unmarshal(respBody, &workItem); err != nil {
		c.debugLog("[DEBUG] JSON Unmarshal error: %v\n", err)
		c.debugLog("[DEBUG] Response body that failed to parse: %s\n", string(respBody))
		return nil, fmt.Errorf("failed to parse work item response: %w", err)
	}
	
	c.debugLog("[DEBUG] Successfully unmarshaled work item: ID=%d\n", workItem.ID)

	return &workItem, nil
}

// GetWorkItemType returns the work item type (User Story, Bug, Task, etc.)
func (wi *WorkItem) GetWorkItemType() string {
	if workItemType, ok := wi.Fields["System.WorkItemType"].(string); ok {
		return normalizeWorkItemType(workItemType)
	}
	return "unknown"
}

// GetTitle returns the work item title
func (wi *WorkItem) GetTitle() string {
	if title, ok := wi.Fields["System.Title"].(string); ok {
		return title
	}
	return ""
}

// GetState returns the work item state
func (wi *WorkItem) GetState() string {
	if state, ok := wi.Fields["System.State"].(string); ok {
		return state
	}
	return ""
}

// GetAssignedTo returns the work item assignee
func (wi *WorkItem) GetAssignedTo() string {
	if assignedTo, ok := wi.Fields["System.AssignedTo"].(map[string]interface{}); ok {
		if displayName, ok := assignedTo["displayName"].(string); ok {
			return displayName
		}
	}
	return "Unassigned"
}

// normalizeWorkItemType converts work item type to a normalized format for branch naming
func normalizeWorkItemType(workItemType string) string {
	// Convert to lowercase and replace spaces with hyphens
	normalized := strings.ToLower(workItemType)
	normalized = strings.ReplaceAll(normalized, " ", "-")

	// Common mappings
	switch normalized {
	case "user-story":
		return "user-story"
	case "bug":
		return "bug"
	case "task":
		return "task"
	case "feature":
		return "feature"
	case "epic":
		return "epic"
	default:
		return normalized
	}
}

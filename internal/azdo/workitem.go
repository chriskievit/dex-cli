package azdo

import (
	"encoding/json"
	"fmt"
	"strings"
)

// WorkItem represents an Azure DevOps work item
type WorkItem struct {
	ID     int                    `json:"id"`
	Fields map[string]interface{} `json:"fields"`
}

// GetWorkItem retrieves a work item by ID
func (c *Client) GetWorkItem(id int) (*WorkItem, error) {
	url := fmt.Sprintf("%s/%s/_apis/wit/workitems/%d?api-version=%s", baseURL, c.organization, id, apiVersion)

	respBody, err := c.doRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get work item: %w", err)
	}

	var workItem WorkItem
	if err := json.Unmarshal(respBody, &workItem); err != nil {
		return nil, fmt.Errorf("failed to parse work item response: %w", err)
	}

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

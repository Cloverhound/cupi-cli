package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// UserTemplate represents a CUC user template
type UserTemplate struct {
	ObjectId    string `json:"ObjectId"`
	Alias       string `json:"Alias"`
	DisplayName string `json:"DisplayName"`
}

type userTemplatesResponse struct {
	Total     string         `json:"@total"`
	Templates OneOrMany[UserTemplate] `json:"UserTemplate"`
}

// ListUserTemplates returns CUC user templates.
func ListUserTemplates(host string, port int, user, pass string, query string, rowsPerPage int) ([]UserTemplate, error) {
	path := "/usertemplates"
	params := url.Values{}
	if query != "" {
		params.Set("query", query)
	}
	if rowsPerPage > 0 {
		params.Set("rowsPerPage", fmt.Sprintf("%d", rowsPerPage))
	}
	if len(params) > 0 {
		path = path + "?" + params.Encode()
	}

	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list user templates: %w", err)
	}

	var resp userTemplatesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse user templates response: %w", err)
	}

	return resp.Templates, nil
}

// GetUserTemplate retrieves a user template by alias or ObjectId.
func GetUserTemplate(host string, port int, user, pass string, aliasOrID string) (*UserTemplate, error) {
	if isUUID(aliasOrID) {
		body, err := Get(host, port, user, pass, "/usertemplates/"+aliasOrID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user template: %w", err)
		}
		var t UserTemplate
		if err := json.Unmarshal(body, &t); err != nil {
			return nil, fmt.Errorf("failed to parse user template: %w", err)
		}
		return &t, nil
	}

	q := fmt.Sprintf("(alias is %s)", aliasOrID)
	templates, err := ListUserTemplates(host, port, user, pass, q, 1)
	if err != nil {
		return nil, err
	}
	if len(templates) == 0 {
		return nil, fmt.Errorf("user template '%s' not found", aliasOrID)
	}
	return &templates[0], nil
}

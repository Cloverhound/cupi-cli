package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// CallHandlerTemplate represents a call handler template
type CallHandlerTemplate struct {
	ObjectId    string `json:"ObjectId"`
	DisplayName string `json:"DisplayName"`
	IsDefault   string `json:"IsDefault"`
}

type callHandlerTemplatesResponse struct {
	Total string                 `json:"@total"`
	Items []CallHandlerTemplate  `json:"CallhandlerTemplate"`
}

// ListCallHandlerTemplates returns call handler templates
func ListCallHandlerTemplates(host string, port int, user, pass string, query string, rowsPerPage int) ([]CallHandlerTemplate, error) {
	path := "/callhandlertemplates"
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
		return nil, fmt.Errorf("failed to list call handler templates: %w", err)
	}

	var resp callHandlerTemplatesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse call handler templates response: %w", err)
	}

	return resp.Items, nil
}

// GetCallHandlerTemplate retrieves a call handler template by name or ObjectId
func GetCallHandlerTemplate(host string, port int, user, pass, nameOrID string) (*CallHandlerTemplate, error) {
	if isUUID(nameOrID) {
		return getCallHandlerTemplateByID(host, port, user, pass, nameOrID)
	}
	return getCallHandlerTemplateByName(host, port, user, pass, nameOrID)
}

func getCallHandlerTemplateByID(host string, port int, user, pass, objectID string) (*CallHandlerTemplate, error) {
	body, err := Get(host, port, user, pass, "/callhandlertemplates/"+objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get call handler template: %w", err)
	}

	var cht CallHandlerTemplate
	if err := json.Unmarshal(body, &cht); err != nil {
		return nil, fmt.Errorf("failed to parse call handler template: %w", err)
	}

	return &cht, nil
}

func getCallHandlerTemplateByName(host string, port int, user, pass, name string) (*CallHandlerTemplate, error) {
	q := fmt.Sprintf("(displayname is %s)", name)
	templates, err := ListCallHandlerTemplates(host, port, user, pass, q, 1)
	if err != nil {
		return nil, err
	}
	if len(templates) == 0 {
		return nil, fmt.Errorf("call handler template '%s' not found", name)
	}
	return &templates[0], nil
}

// CreateCallHandlerTemplate creates a new call handler template
func CreateCallHandlerTemplate(host string, port int, user, pass string, fields map[string]interface{}) error {
	_, err := Post(host, port, user, pass, "/callhandlertemplates", fields)
	if err != nil {
		return fmt.Errorf("failed to create call handler template: %w", err)
	}
	return nil
}

// UpdateCallHandlerTemplate updates a call handler template by name or ObjectId
func UpdateCallHandlerTemplate(host string, port int, user, pass, nameOrID string, fields map[string]interface{}) error {
	cht, err := GetCallHandlerTemplate(host, port, user, pass, nameOrID)
	if err != nil {
		return err
	}
	return Put(host, port, user, pass, "/callhandlertemplates/"+cht.ObjectId, fields)
}

// DeleteCallHandlerTemplate deletes a call handler template by name or ObjectId
func DeleteCallHandlerTemplate(host string, port int, user, pass, nameOrID string) error {
	cht, err := GetCallHandlerTemplate(host, port, user, pass, nameOrID)
	if err != nil {
		return err
	}
	return Delete(host, port, user, pass, "/callhandlertemplates/"+cht.ObjectId)
}

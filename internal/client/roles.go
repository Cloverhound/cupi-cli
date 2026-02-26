package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// CustomRole represents a custom role
type CustomRole struct {
	ObjectId    string `json:"ObjectId"`
	DisplayName string `json:"DisplayName"`
	Description string `json:"Description"`
}

type customRolesResponse struct {
	Total string       `json:"@total"`
	Items []CustomRole `json:"CustomRole"`
}

// ListCustomRoles returns all custom roles
func ListCustomRoles(host string, port int, user, pass string, query string, rowsPerPage int) ([]CustomRole, error) {
	path := "/customroles"
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
		return nil, fmt.Errorf("failed to list custom roles: %w", err)
	}

	var resp customRolesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse custom roles response: %w", err)
	}

	return resp.Items, nil
}

// GetCustomRole retrieves a custom role by name or ObjectId
func GetCustomRole(host string, port int, user, pass, nameOrID string) (*CustomRole, error) {
	if isUUID(nameOrID) {
		return getCustomRoleByID(host, port, user, pass, nameOrID)
	}
	return getCustomRoleByName(host, port, user, pass, nameOrID)
}

func getCustomRoleByID(host string, port int, user, pass, objectID string) (*CustomRole, error) {
	body, err := Get(host, port, user, pass, "/customroles/"+objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get custom role: %w", err)
	}

	var cr CustomRole
	if err := json.Unmarshal(body, &cr); err != nil {
		return nil, fmt.Errorf("failed to parse custom role: %w", err)
	}

	return &cr, nil
}

func getCustomRoleByName(host string, port int, user, pass, name string) (*CustomRole, error) {
	q := fmt.Sprintf("(displayname is %s)", name)
	roles, err := ListCustomRoles(host, port, user, pass, q, 1)
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return nil, fmt.Errorf("custom role '%s' not found", name)
	}
	return &roles[0], nil
}

// CreateCustomRole creates a new custom role
func CreateCustomRole(host string, port int, user, pass string, fields map[string]interface{}) error {
	_, err := Post(host, port, user, pass, "/customroles", fields)
	if err != nil {
		return fmt.Errorf("failed to create custom role: %w", err)
	}
	return nil
}

// UpdateCustomRole updates a custom role by name or ObjectId
func UpdateCustomRole(host string, port int, user, pass, nameOrID string, fields map[string]interface{}) error {
	cr, err := GetCustomRole(host, port, user, pass, nameOrID)
	if err != nil {
		return err
	}
	return Put(host, port, user, pass, "/customroles/"+cr.ObjectId, fields)
}

// DeleteCustomRole deletes a custom role by name or ObjectId
func DeleteCustomRole(host string, port int, user, pass, nameOrID string) error {
	cr, err := GetCustomRole(host, port, user, pass, nameOrID)
	if err != nil {
		return err
	}
	return Delete(host, port, user, pass, "/customroles/"+cr.ObjectId)
}

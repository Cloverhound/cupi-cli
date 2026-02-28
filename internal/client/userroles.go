package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// UserRole represents a role assigned to a user
type UserRole struct {
	ObjectId     string `json:"ObjectId"`
	RoleObjectId string `json:"RoleObjectId"`
	UserObjectId string `json:"UserObjectId"`
}

type userRolesResponse struct {
	Total string     `json:"@total"`
	Items OneOrMany[UserRole] `json:"UserRole"`
}

// Role represents a role in the system
type Role struct {
	ObjectId    string `json:"ObjectId"`
	DisplayName string `json:"DisplayName"`
	RoleType    string `json:"RoleType"`
}

type rolesResponse struct {
	Total string `json:"@total"`
	Items OneOrMany[Role] `json:"Role"`
}

// ListUserRoles returns all roles assigned to a user
func ListUserRoles(host string, port int, user, pass, userObjectId string) ([]UserRole, error) {
	path := fmt.Sprintf("/users/%s/userroles", url.PathEscape(userObjectId))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list user roles: %w", err)
	}

	var resp userRolesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse user roles response: %w", err)
	}

	return resp.Items, nil
}

// AddUserRole assigns a role to a user
func AddUserRole(host string, port int, user, pass, userObjectId, roleObjectId string) error {
	path := fmt.Sprintf("/users/%s/userroles", url.PathEscape(userObjectId))
	fields := map[string]interface{}{
		"RoleObjectId": roleObjectId,
	}
	_, err := Post(host, port, user, pass, path, fields)
	if err != nil {
		return fmt.Errorf("failed to add user role: %w", err)
	}
	return nil
}

// RemoveUserRole removes a role from a user
func RemoveUserRole(host string, port int, user, pass, userObjectId, roleObjectId string) error {
	path := fmt.Sprintf("/users/%s/userroles/%s", url.PathEscape(userObjectId), url.PathEscape(roleObjectId))
	if err := Delete(host, port, user, pass, path); err != nil {
		return fmt.Errorf("failed to remove user role: %w", err)
	}
	return nil
}

// ListRoles returns all available roles in the system
func ListRoles(host string, port int, user, pass string) ([]Role, error) {
	body, err := Get(host, port, user, pass, "/roles")
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}

	var resp rolesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse roles response: %w", err)
	}

	return resp.Items, nil
}

package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// User represents a CUC mailbox user
type User struct {
	ObjectId    string `json:"ObjectId"`
	Alias       string `json:"Alias"`
	DisplayName string `json:"DisplayName"`
	DtmfAccessId string `json:"DtmfAccessId"`
	FirstName   string `json:"FirstName"`
	LastName    string `json:"LastName"`
	Department  string `json:"Department"`
	Title       string `json:"Title"`
	IsTemplate  string `json:"IsTemplate"`
}

type usersListResponse struct {
	Total string `json:"@total"`
	Users []User `json:"User"`
}

// ListUsers returns CUC users, optionally with a search query and row limit.
func ListUsers(host string, port int, user, pass string, query string, rowsPerPage int) ([]User, error) {
	path := "/users"
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
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	var resp usersListResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse users response: %w", err)
	}

	return resp.Users, nil
}

// GetUser retrieves a user by alias or ObjectId.
// If input looks like a UUID, GETs directly; otherwise searches by alias.
func GetUser(host string, port int, user, pass string, aliasOrID string) (*User, error) {
	if isUUID(aliasOrID) {
		return getUserByID(host, port, user, pass, aliasOrID)
	}
	return getUserByAlias(host, port, user, pass, aliasOrID)
}

func getUserByID(host string, port int, user, pass, objectID string) (*User, error) {
	body, err := Get(host, port, user, pass, "/users/"+objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	var u User
	if err := json.Unmarshal(body, &u); err != nil {
		return nil, fmt.Errorf("failed to parse user: %w", err)
	}
	return &u, nil
}

func getUserByAlias(host string, port int, user, pass, alias string) (*User, error) {
	q := fmt.Sprintf("(alias is %s)", alias)
	users, err := ListUsers(host, port, user, pass, q, 1)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("user '%s' not found", alias)
	}
	return &users[0], nil
}

// CreateUser creates a new CUC user. templateAlias is required.
func CreateUser(host string, port int, user, pass string, templateAlias string, fields map[string]interface{}) (*User, error) {
	path := fmt.Sprintf("/users?templateAlias=%s", url.QueryEscape(templateAlias))
	body, err := Post(host, port, user, pass, path, fields)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	var u User
	if err := json.Unmarshal(body, &u); err != nil {
		// Some versions return just the ObjectId URI; try to extract
		if len(body) > 0 {
			return &User{Alias: fmt.Sprintf("%v", fields["Alias"])}, nil
		}
		return nil, fmt.Errorf("failed to parse create response: %w", err)
	}
	return &u, nil
}

// UpdateUser updates a CUC user by alias or ObjectId.
func UpdateUser(host string, port int, user, pass string, aliasOrID string, fields map[string]interface{}) error {
	u, err := GetUser(host, port, user, pass, aliasOrID)
	if err != nil {
		return err
	}
	return Put(host, port, user, pass, "/users/"+u.ObjectId, fields)
}

// DeleteUser deletes a CUC user by alias or ObjectId.
func DeleteUser(host string, port int, user, pass string, aliasOrID string) error {
	u, err := GetUser(host, port, user, pass, aliasOrID)
	if err != nil {
		return err
	}
	return Delete(host, port, user, pass, "/users/"+u.ObjectId)
}

// isUUID returns true if the string looks like a UUID (contains hyphens and is 36 chars).
func isUUID(s string) bool {
	return len(s) == 36 && strings.Count(s, "-") == 4
}

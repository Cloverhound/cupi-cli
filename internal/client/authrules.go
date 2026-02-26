package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// AuthRule represents an authentication rule
type AuthRule struct {
	ObjectId         string `json:"ObjectId"`
	DisplayName      string `json:"DisplayName"`
	MaxLogonAttempts string `json:"MaxLogonAttempts"`
	LockoutDuration  string `json:"LockoutDuration"`
}

type authRulesResponse struct {
	Total string     `json:"@total"`
	Items []AuthRule `json:"AuthenticationRule"`
}

// ListAuthRules returns all authentication rules
func ListAuthRules(host string, port int, user, pass string, query string, rowsPerPage int) ([]AuthRule, error) {
	path := "/authenticationrules"
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
		return nil, fmt.Errorf("failed to list authentication rules: %w", err)
	}

	var resp authRulesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse authentication rules response: %w", err)
	}

	return resp.Items, nil
}

// GetAuthRule retrieves an authentication rule by name or ObjectId
func GetAuthRule(host string, port int, user, pass, nameOrID string) (*AuthRule, error) {
	if isUUID(nameOrID) {
		return getAuthRuleByID(host, port, user, pass, nameOrID)
	}
	return getAuthRuleByName(host, port, user, pass, nameOrID)
}

func getAuthRuleByID(host string, port int, user, pass, objectID string) (*AuthRule, error) {
	body, err := Get(host, port, user, pass, "/authenticationrules/"+objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get authentication rule: %w", err)
	}

	var ar AuthRule
	if err := json.Unmarshal(body, &ar); err != nil {
		return nil, fmt.Errorf("failed to parse authentication rule: %w", err)
	}

	return &ar, nil
}

func getAuthRuleByName(host string, port int, user, pass, name string) (*AuthRule, error) {
	q := fmt.Sprintf("(displayname is %s)", name)
	rules, err := ListAuthRules(host, port, user, pass, q, 1)
	if err != nil {
		return nil, err
	}
	if len(rules) == 0 {
		return nil, fmt.Errorf("authentication rule '%s' not found", name)
	}
	return &rules[0], nil
}

// CreateAuthRule creates a new authentication rule
func CreateAuthRule(host string, port int, user, pass string, fields map[string]interface{}) error {
	_, err := Post(host, port, user, pass, "/authenticationrules", fields)
	if err != nil {
		return fmt.Errorf("failed to create authentication rule: %w", err)
	}
	return nil
}

// UpdateAuthRule updates an authentication rule by name or ObjectId
func UpdateAuthRule(host string, port int, user, pass, nameOrID string, fields map[string]interface{}) error {
	ar, err := GetAuthRule(host, port, user, pass, nameOrID)
	if err != nil {
		return err
	}
	return Put(host, port, user, pass, "/authenticationrules/"+ar.ObjectId, fields)
}

// DeleteAuthRule deletes an authentication rule by name or ObjectId
func DeleteAuthRule(host string, port int, user, pass, nameOrID string) error {
	ar, err := GetAuthRule(host, port, user, pass, nameOrID)
	if err != nil {
		return err
	}
	return Delete(host, port, user, pass, "/authenticationrules/"+ar.ObjectId)
}

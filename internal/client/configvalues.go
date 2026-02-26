package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// ConfigValue represents a system configuration value
type ConfigValue struct {
	FullName    string `json:"FullName"`
	Value       string `json:"Value"`
	Description string `json:"Description"`
	Type        string `json:"Type"`
}

type configValuesResponse struct {
	Total string       `json:"@total"`
	Items []ConfigValue `json:"ConfigurationValue"`
}

// ListConfigValues returns configuration values
func ListConfigValues(host string, port int, user, pass string, query string, rowsPerPage int) ([]ConfigValue, error) {
	path := "/configurationvalues"
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
		return nil, fmt.Errorf("failed to list configuration values: %w", err)
	}

	var resp configValuesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse configuration values response: %w", err)
	}

	return resp.Items, nil
}

// GetConfigValue retrieves a configuration value by full name
func GetConfigValue(host string, port int, user, pass, fullName string) (*ConfigValue, error) {
	path := fmt.Sprintf("/configurationvalues/%s", url.PathEscape(fullName))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get configuration value: %w", err)
	}

	var cv ConfigValue
	if err := json.Unmarshal(body, &cv); err != nil {
		return nil, fmt.Errorf("failed to parse configuration value: %w", err)
	}

	return &cv, nil
}

// UpdateConfigValue updates a configuration value
func UpdateConfigValue(host string, port int, user, pass, fullName, value string) error {
	path := fmt.Sprintf("/configurationvalues/%s", url.PathEscape(fullName))
	fields := map[string]interface{}{
		"Value": value,
	}
	if err := Put(host, port, user, pass, path, fields); err != nil {
		return fmt.Errorf("failed to update configuration value: %w", err)
	}
	return nil
}

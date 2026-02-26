package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// AlternateName represents an alternate name for a user
type AlternateName struct {
	ObjectId           string `json:"ObjectId"`
	FirstName          string `json:"FirstName"`
	LastName           string `json:"LastName"`
	GlobalUserObjectId string `json:"GlobalUserObjectId"`
}

type alternateNamesResponse struct {
	Total string         `json:"@total"`
	Items []AlternateName `json:"AlternateName"`
}

// ListAlternateNames returns alternate names for a user
func ListAlternateNames(host string, port int, user, pass, userObjectId string) ([]AlternateName, error) {
	q := fmt.Sprintf("(GlobalUserObjectId is %s)", userObjectId)
	params := url.Values{}
	params.Set("query", q)
	path := "/alternatenames?" + params.Encode()

	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list alternate names: %w", err)
	}

	var resp alternateNamesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse alternate names response: %w", err)
	}

	return resp.Items, nil
}

// GetAlternateName retrieves a specific alternate name
func GetAlternateName(host string, port int, user, pass, objectId string) (*AlternateName, error) {
	path := fmt.Sprintf("/alternatenames/%s", url.PathEscape(objectId))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get alternate name: %w", err)
	}

	var an AlternateName
	if err := json.Unmarshal(body, &an); err != nil {
		return nil, fmt.Errorf("failed to parse alternate name: %w", err)
	}

	return &an, nil
}

// CreateAlternateName creates a new alternate name
func CreateAlternateName(host string, port int, user, pass string, fields map[string]interface{}) error {
	_, err := Post(host, port, user, pass, "/alternatenames", fields)
	if err != nil {
		return fmt.Errorf("failed to create alternate name: %w", err)
	}
	return nil
}

// UpdateAlternateName updates an alternate name
func UpdateAlternateName(host string, port int, user, pass, objectId string, fields map[string]interface{}) error {
	path := fmt.Sprintf("/alternatenames/%s", url.PathEscape(objectId))
	if err := Put(host, port, user, pass, path, fields); err != nil {
		return fmt.Errorf("failed to update alternate name: %w", err)
	}
	return nil
}

// DeleteAlternateName deletes an alternate name
func DeleteAlternateName(host string, port int, user, pass, objectId string) error {
	path := fmt.Sprintf("/alternatenames/%s", url.PathEscape(objectId))
	if err := Delete(host, port, user, pass, path); err != nil {
		return fmt.Errorf("failed to delete alternate name: %w", err)
	}
	return nil
}

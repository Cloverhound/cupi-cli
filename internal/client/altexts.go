package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// AlternateExtension represents a user's alternate extension
type AlternateExtension struct {
	ObjectId          string `json:"ObjectId"`
	DtmfAccessId      string `json:"DtmfAccessId"`
	IdIndex           int    `json:"IdIndex"`
	PartitionObjectId string `json:"PartitionObjectId"`
}

type alternateExtensionsResponse struct {
	Total string                 `json:"@total"`
	Items OneOrMany[AlternateExtension] `json:"AlternateExtension"`
}

// ListAlternateExtensions returns alternate extensions for a user
func ListAlternateExtensions(host string, port int, user, pass, userObjectId string) ([]AlternateExtension, error) {
	path := fmt.Sprintf("/users/%s/alternateextensions", url.PathEscape(userObjectId))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list alternate extensions: %w", err)
	}

	var resp alternateExtensionsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse alternate extensions response: %w", err)
	}

	return resp.Items, nil
}

// GetAlternateExtension retrieves a specific alternate extension
func GetAlternateExtension(host string, port int, user, pass, userObjectId, objectId string) (*AlternateExtension, error) {
	path := fmt.Sprintf("/users/%s/alternateextensions/%s", url.PathEscape(userObjectId), url.PathEscape(objectId))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get alternate extension: %w", err)
	}

	var ae AlternateExtension
	if err := json.Unmarshal(body, &ae); err != nil {
		return nil, fmt.Errorf("failed to parse alternate extension: %w", err)
	}

	return &ae, nil
}

// CreateAlternateExtension creates a new alternate extension
func CreateAlternateExtension(host string, port int, user, pass, userObjectId string, fields map[string]interface{}) error {
	path := fmt.Sprintf("/users/%s/alternateextensions", url.PathEscape(userObjectId))
	_, err := Post(host, port, user, pass, path, fields)
	if err != nil {
		return fmt.Errorf("failed to create alternate extension: %w", err)
	}
	return nil
}

// UpdateAlternateExtension updates an alternate extension
func UpdateAlternateExtension(host string, port int, user, pass, userObjectId, objectId string, fields map[string]interface{}) error {
	path := fmt.Sprintf("/users/%s/alternateextensions/%s", url.PathEscape(userObjectId), url.PathEscape(objectId))
	if err := Put(host, port, user, pass, path, fields); err != nil {
		return fmt.Errorf("failed to update alternate extension: %w", err)
	}
	return nil
}

// DeleteAlternateExtension deletes an alternate extension
func DeleteAlternateExtension(host string, port int, user, pass, userObjectId, objectId string) error {
	path := fmt.Sprintf("/users/%s/alternateextensions/%s", url.PathEscape(userObjectId), url.PathEscape(objectId))
	if err := Delete(host, port, user, pass, path); err != nil {
		return fmt.Errorf("failed to delete alternate extension: %w", err)
	}
	return nil
}

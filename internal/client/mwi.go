package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// MWI represents a Message Waiting Indicator for a user
type MWI struct {
	ObjectId            string `json:"ObjectId"`
	DisplayName         string `json:"DisplayName"`
	Active              string `json:"Active"`
	MWIExtension        string `json:"MWIExtension"`
	MediaSwitchObjectId string `json:"MediaSwitchObjectId"`
}

type mwisResponse struct {
	Total string `json:"@total"`
	Items OneOrMany[MWI] `json:"MWI"`
}

// ListMWIs returns all MWIs for a user
func ListMWIs(host string, port int, user, pass, userObjectId string) ([]MWI, error) {
	path := fmt.Sprintf("/users/%s/mwis", url.PathEscape(userObjectId))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list MWIs: %w", err)
	}

	var resp mwisResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse MWIs response: %w", err)
	}

	return resp.Items, nil
}

// GetMWI retrieves a specific MWI
func GetMWI(host string, port int, user, pass, userObjectId, mwiObjectId string) (*MWI, error) {
	path := fmt.Sprintf("/users/%s/mwis/%s", url.PathEscape(userObjectId), url.PathEscape(mwiObjectId))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get MWI: %w", err)
	}

	var m MWI
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf("failed to parse MWI: %w", err)
	}

	return &m, nil
}

// CreateMWI creates a new MWI for a user
func CreateMWI(host string, port int, user, pass, userObjectId string, fields map[string]interface{}) error {
	path := fmt.Sprintf("/users/%s/mwis", url.PathEscape(userObjectId))
	_, err := Post(host, port, user, pass, path, fields)
	if err != nil {
		return fmt.Errorf("failed to create MWI: %w", err)
	}
	return nil
}

// UpdateMWI updates a MWI
func UpdateMWI(host string, port int, user, pass, userObjectId, mwiObjectId string, fields map[string]interface{}) error {
	path := fmt.Sprintf("/users/%s/mwis/%s", url.PathEscape(userObjectId), url.PathEscape(mwiObjectId))
	if err := Put(host, port, user, pass, path, fields); err != nil {
		return fmt.Errorf("failed to update MWI: %w", err)
	}
	return nil
}

// DeleteMWI deletes a MWI
func DeleteMWI(host string, port int, user, pass, userObjectId, mwiObjectId string) error {
	path := fmt.Sprintf("/users/%s/mwis/%s", url.PathEscape(userObjectId), url.PathEscape(mwiObjectId))
	if err := Delete(host, port, user, pass, path); err != nil {
		return fmt.Errorf("failed to delete MWI: %w", err)
	}
	return nil
}

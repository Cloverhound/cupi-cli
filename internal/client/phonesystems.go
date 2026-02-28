package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// PhoneSystem represents a phone system (call manager)
type PhoneSystem struct {
	ObjectId        string `json:"ObjectId"`
	DisplayName     string `json:"DisplayName"`
	CallManagerType string `json:"CallManagerType"`
}

type phoneSystemsResponse struct {
	Total string       `json:"@total"`
	Items OneOrMany[PhoneSystem] `json:"PhoneSystem"`
}

// AXLServer represents an AXL server in a phone system
type AXLServer struct {
	ObjectId   string `json:"ObjectId"`
	ServerName string `json:"ServerName"`
	Port       string `json:"Port"`
}

type axlServersResponse struct {
	Total string      `json:"@total"`
	Items OneOrMany[AXLServer] `json:"AXLServer"`
}

// ListPhoneSystems returns all phone systems
func ListPhoneSystems(host string, port int, user, pass string, query string, rowsPerPage int) ([]PhoneSystem, error) {
	path := "/phonesystems"
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
		return nil, fmt.Errorf("failed to list phone systems: %w", err)
	}

	var resp phoneSystemsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse phone systems response: %w", err)
	}

	return resp.Items, nil
}

// GetPhoneSystem retrieves a phone system by name or ObjectId
func GetPhoneSystem(host string, port int, user, pass, nameOrID string) (*PhoneSystem, error) {
	if isUUID(nameOrID) {
		return getPhoneSystemByID(host, port, user, pass, nameOrID)
	}
	return getPhoneSystemByName(host, port, user, pass, nameOrID)
}

func getPhoneSystemByID(host string, port int, user, pass, objectID string) (*PhoneSystem, error) {
	body, err := Get(host, port, user, pass, "/phonesystems/"+objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get phone system: %w", err)
	}

	var ps PhoneSystem
	if err := json.Unmarshal(body, &ps); err != nil {
		return nil, fmt.Errorf("failed to parse phone system: %w", err)
	}

	return &ps, nil
}

func getPhoneSystemByName(host string, port int, user, pass, name string) (*PhoneSystem, error) {
	q := fmt.Sprintf("(displayname is %s)", name)
	systems, err := ListPhoneSystems(host, port, user, pass, q, 1)
	if err != nil {
		return nil, err
	}
	if len(systems) == 0 {
		return nil, fmt.Errorf("phone system '%s' not found", name)
	}
	return &systems[0], nil
}

// CreatePhoneSystem creates a new phone system
func CreatePhoneSystem(host string, port int, user, pass string, fields map[string]interface{}) error {
	_, err := Post(host, port, user, pass, "/phonesystems", fields)
	if err != nil {
		return fmt.Errorf("failed to create phone system: %w", err)
	}
	return nil
}

// UpdatePhoneSystem updates a phone system by name or ObjectId
func UpdatePhoneSystem(host string, port int, user, pass, nameOrID string, fields map[string]interface{}) error {
	ps, err := GetPhoneSystem(host, port, user, pass, nameOrID)
	if err != nil {
		return err
	}
	return Put(host, port, user, pass, "/phonesystems/"+ps.ObjectId, fields)
}

// ListAXLServers returns all AXL servers for a phone system
func ListAXLServers(host string, port int, user, pass, phoneSystemObjectId string) ([]AXLServer, error) {
	path := fmt.Sprintf("/phonesystems/%s/axlservers", url.PathEscape(phoneSystemObjectId))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list AXL servers: %w", err)
	}

	var resp axlServersResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse AXL servers response: %w", err)
	}

	return resp.Items, nil
}

// CreateAXLServer creates a new AXL server for a phone system
func CreateAXLServer(host string, port int, user, pass, phoneSystemObjectId string, fields map[string]interface{}) error {
	path := fmt.Sprintf("/phonesystems/%s/axlservers", url.PathEscape(phoneSystemObjectId))
	_, err := Post(host, port, user, pass, path, fields)
	if err != nil {
		return fmt.Errorf("failed to create AXL server: %w", err)
	}
	return nil
}

// UpdateAXLServer updates an AXL server
func UpdateAXLServer(host string, port int, user, pass, phoneSystemObjectId, axlServerObjectId string, fields map[string]interface{}) error {
	path := fmt.Sprintf("/phonesystems/%s/axlservers/%s", url.PathEscape(phoneSystemObjectId), url.PathEscape(axlServerObjectId))
	if err := Put(host, port, user, pass, path, fields); err != nil {
		return fmt.Errorf("failed to update AXL server: %w", err)
	}
	return nil
}

// DeleteAXLServer deletes an AXL server
func DeleteAXLServer(host string, port int, user, pass, phoneSystemObjectId, axlServerObjectId string) error {
	path := fmt.Sprintf("/phonesystems/%s/axlservers/%s", url.PathEscape(phoneSystemObjectId), url.PathEscape(axlServerObjectId))
	if err := Delete(host, port, user, pass, path); err != nil {
		return fmt.Errorf("failed to delete AXL server: %w", err)
	}
	return nil
}

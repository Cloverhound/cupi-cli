package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// DirectoryHandler represents a directory handler
type DirectoryHandler struct {
	ObjectId    string `json:"ObjectId"`
	DisplayName string `json:"DisplayName"`
	DtmfAccessId string `json:"DtmfAccessId"`
}

type directoryHandlersResponse struct {
	Total string              `json:"@total"`
	Items []DirectoryHandler  `json:"DirectoryHandler"`
}

// ListDirectoryHandlers returns directory handlers
func ListDirectoryHandlers(host string, port int, user, pass string, query string, rowsPerPage int) ([]DirectoryHandler, error) {
	path := "/handlers/directoryhandlers"
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
		return nil, fmt.Errorf("failed to list directory handlers: %w", err)
	}

	var resp directoryHandlersResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse directory handlers response: %w", err)
	}

	return resp.Items, nil
}

// GetDirectoryHandler retrieves a directory handler by name or ObjectId
func GetDirectoryHandler(host string, port int, user, pass, nameOrID string) (*DirectoryHandler, error) {
	if isUUID(nameOrID) {
		return getDirectoryHandlerByID(host, port, user, pass, nameOrID)
	}
	return getDirectoryHandlerByName(host, port, user, pass, nameOrID)
}

func getDirectoryHandlerByID(host string, port int, user, pass, objectID string) (*DirectoryHandler, error) {
	body, err := Get(host, port, user, pass, "/handlers/directoryhandlers/"+objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get directory handler: %w", err)
	}

	var dh DirectoryHandler
	if err := json.Unmarshal(body, &dh); err != nil {
		return nil, fmt.Errorf("failed to parse directory handler: %w", err)
	}

	return &dh, nil
}

func getDirectoryHandlerByName(host string, port int, user, pass, name string) (*DirectoryHandler, error) {
	q := fmt.Sprintf("(displayname is %s)", name)
	handlers, err := ListDirectoryHandlers(host, port, user, pass, q, 1)
	if err != nil {
		return nil, err
	}
	if len(handlers) == 0 {
		return nil, fmt.Errorf("directory handler '%s' not found", name)
	}
	return &handlers[0], nil
}

// CreateDirectoryHandler creates a new directory handler
func CreateDirectoryHandler(host string, port int, user, pass string, fields map[string]interface{}) error {
	_, err := Post(host, port, user, pass, "/handlers/directoryhandlers", fields)
	if err != nil {
		return fmt.Errorf("failed to create directory handler: %w", err)
	}
	return nil
}

// UpdateDirectoryHandler updates a directory handler by name or ObjectId
func UpdateDirectoryHandler(host string, port int, user, pass, nameOrID string, fields map[string]interface{}) error {
	dh, err := GetDirectoryHandler(host, port, user, pass, nameOrID)
	if err != nil {
		return err
	}
	return Put(host, port, user, pass, "/handlers/directoryhandlers/"+dh.ObjectId, fields)
}

// DeleteDirectoryHandler deletes a directory handler by name or ObjectId
func DeleteDirectoryHandler(host string, port int, user, pass, nameOrID string) error {
	dh, err := GetDirectoryHandler(host, port, user, pass, nameOrID)
	if err != nil {
		return err
	}
	return Delete(host, port, user, pass, "/handlers/directoryhandlers/"+dh.ObjectId)
}

package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// MenuEntry represents a menu entry in a call handler
type MenuEntry struct {
	TouchtoneKey              string `json:"TouchtoneKey"`
	Action                    string `json:"Action"`
	TargetConversation        string `json:"TargetConversation"`
	TargetHandlerObjectId     string `json:"TargetHandlerObjectId"`
}

type menuEntriesResponse struct {
	Total string      `json:"@total"`
	Items []MenuEntry `json:"MenuItem"`
}

// ListMenuEntries returns all menu entries for a call handler
func ListMenuEntries(host string, port int, user, pass, handlerObjectId string) ([]MenuEntry, error) {
	path := fmt.Sprintf("/handlers/callhandlers/%s/menuentries", url.PathEscape(handlerObjectId))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list menu entries: %w", err)
	}

	var resp menuEntriesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse menu entries response: %w", err)
	}

	return resp.Items, nil
}

// GetMenuEntry retrieves a specific menu entry by key
func GetMenuEntry(host string, port int, user, pass, handlerObjectId, key string) (*MenuEntry, error) {
	path := fmt.Sprintf("/handlers/callhandlers/%s/menuentries/%s", url.PathEscape(handlerObjectId), url.PathEscape(key))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu entry: %w", err)
	}

	var m MenuEntry
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, fmt.Errorf("failed to parse menu entry: %w", err)
	}

	return &m, nil
}

// UpdateMenuEntry updates a menu entry
func UpdateMenuEntry(host string, port int, user, pass, handlerObjectId, key string, fields map[string]interface{}) error {
	path := fmt.Sprintf("/handlers/callhandlers/%s/menuentries/%s", url.PathEscape(handlerObjectId), url.PathEscape(key))
	if err := Put(host, port, user, pass, path, fields); err != nil {
		return fmt.Errorf("failed to update menu entry: %w", err)
	}
	return nil
}

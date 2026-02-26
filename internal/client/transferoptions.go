package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// TransferOption represents transfer options for a call handler
type TransferOption struct {
	TransferType string `json:"TransferType"`
	Action       string `json:"Action"`
	Extension    string `json:"Extension"`
	RingCount    string `json:"RingCount"`
	Enabled      string `json:"Enabled"`
}

type transferOptionsResponse struct {
	Total string           `json:"@total"`
	Items []TransferOption `json:"TransferOption"`
}

// ListTransferOptions returns all transfer options for a call handler
func ListTransferOptions(host string, port int, user, pass, handlerObjectId string) ([]TransferOption, error) {
	path := fmt.Sprintf("/handlers/callhandlers/%s/transferoptions", url.PathEscape(handlerObjectId))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list transfer options: %w", err)
	}

	var resp transferOptionsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse transfer options response: %w", err)
	}

	return resp.Items, nil
}

// GetTransferOption retrieves a specific transfer option by type
func GetTransferOption(host string, port int, user, pass, handlerObjectId, transferType string) (*TransferOption, error) {
	path := fmt.Sprintf("/handlers/callhandlers/%s/transferoptions/%s", url.PathEscape(handlerObjectId), url.PathEscape(transferType))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get transfer option: %w", err)
	}

	var t TransferOption
	if err := json.Unmarshal(body, &t); err != nil {
		return nil, fmt.Errorf("failed to parse transfer option: %w", err)
	}

	return &t, nil
}

// UpdateTransferOption updates a transfer option
func UpdateTransferOption(host string, port int, user, pass, handlerObjectId, transferType string, fields map[string]interface{}) error {
	path := fmt.Sprintf("/handlers/callhandlers/%s/transferoptions/%s", url.PathEscape(handlerObjectId), url.PathEscape(transferType))
	if err := Put(host, port, user, pass, path, fields); err != nil {
		return fmt.Errorf("failed to update transfer option: %w", err)
	}
	return nil
}

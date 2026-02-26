package client

import (
	"encoding/json"
	"fmt"
)

// SMTPServerConfig represents SMTP server configuration
type SMTPServerConfig struct {
	ObjectId string `json:"ObjectId"`
	SmartHost string `json:"SmartHost"`
	Port     string `json:"Port"`
	UseSsl   string `json:"UseSsl"`
}

type smtpServerConfigsResponse struct {
	Total string               `json:"@total"`
	Items []SMTPServerConfig   `json:"SMTPServerConfig"`
}

// SMTPClientConfig represents SMTP client configuration
type SMTPClientConfig struct {
	ObjectId   string `json:"ObjectId"`
	ServerName string `json:"ServerName"`
	Port       string `json:"Port"`
	UseAuth    string `json:"UseAuth"`
}

type smtpClientConfigsResponse struct {
	Total string               `json:"@total"`
	Items []SMTPClientConfig   `json:"SMTPClientConfig"`
}

// GetSMTPServerConfig retrieves SMTP server configuration
func GetSMTPServerConfig(host string, port int, user, pass string) (*SMTPServerConfig, error) {
	body, err := Get(host, port, user, pass, "/smtpserver/serverconfigs")
	if err != nil {
		return nil, fmt.Errorf("failed to get SMTP server config: %w", err)
	}

	var resp smtpServerConfigsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse SMTP server config response: %w", err)
	}

	if len(resp.Items) == 0 {
		return nil, fmt.Errorf("no SMTP server config found")
	}

	return &resp.Items[0], nil
}

// UpdateSMTPServerConfig updates SMTP server configuration
func UpdateSMTPServerConfig(host string, port int, user, pass string, fields map[string]interface{}) error {
	// Get current config to find the ObjectId (usually empty or ignored in PUT)
	if err := Put(host, port, user, pass, "/smtpserver/serverconfigs", fields); err != nil {
		return fmt.Errorf("failed to update SMTP server config: %w", err)
	}
	return nil
}

// GetSMTPClientConfig retrieves SMTP client configuration
func GetSMTPClientConfig(host string, port int, user, pass string) (*SMTPClientConfig, error) {
	body, err := Get(host, port, user, pass, "/smtpclient/clientconfigs")
	if err != nil {
		return nil, fmt.Errorf("failed to get SMTP client config: %w", err)
	}

	var resp smtpClientConfigsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse SMTP client config response: %w", err)
	}

	if len(resp.Items) == 0 {
		return nil, fmt.Errorf("no SMTP client config found")
	}

	return &resp.Items[0], nil
}

// UpdateSMTPClientConfig updates SMTP client configuration
func UpdateSMTPClientConfig(host string, port int, user, pass string, fields map[string]interface{}) error {
	if err := Put(host, port, user, pass, "/smtpclient/clientconfigs", fields); err != nil {
		return fmt.Errorf("failed to update SMTP client config: %w", err)
	}
	return nil
}

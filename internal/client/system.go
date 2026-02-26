package client

import (
	"encoding/json"
	"fmt"
)

// SystemInfo represents CUC system information
type SystemInfo struct {
	DisplayName        string `json:"DisplayName"`
	Version            string `json:"Version"`
	SerialNumber       string `json:"SerialNumber"`
	IpAddress          string `json:"IpAddress"`
	Hostname           string `json:"Hostname"`
	DomainName         string `json:"DomainName"`
	SmtpSmartHost      string `json:"SmtpSmartHost"`
	MaxMailboxSize     string `json:"MaxMailboxSize"`
}

// GetSystemInfo retrieves CUC system information.
func GetSystemInfo(host string, port int, user, pass string) (*SystemInfo, error) {
	body, err := Get(host, port, user, pass, "/systeminfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get system info: %w", err)
	}

	var info SystemInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("failed to parse system info: %w", err)
	}

	return &info, nil
}

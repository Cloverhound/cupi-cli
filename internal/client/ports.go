package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// Port represents a system port
type Port struct {
	ObjectId               string `json:"ObjectId"`
	DisplayName            string `json:"DisplayName"`
	Enabled                string `json:"Enabled"`
	AnswerCalls            string `json:"AnswerCalls"`
	PortNumber             string `json:"PortNumber"`
	MediaRemoteServiceIP   string `json:"MediaRemoteServiceIP"`
}

type portsResponse struct {
	Total string `json:"@total"`
	Items OneOrMany[Port] `json:"Port"`
}

// ListPorts returns all ports
func ListPorts(host string, port int, user, pass string, query string, rowsPerPage int) ([]Port, error) {
	path := "/ports"
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
		return nil, fmt.Errorf("failed to list ports: %w", err)
	}

	var resp portsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse ports response: %w", err)
	}

	return resp.Items, nil
}

// GetPort retrieves a port by name or ObjectId
func GetPort(host string, port int, user, pass, nameOrID string) (*Port, error) {
	if isUUID(nameOrID) {
		return getPortByID(host, port, user, pass, nameOrID)
	}
	return getPortByName(host, port, user, pass, nameOrID)
}

func getPortByID(host string, port int, user, pass, objectID string) (*Port, error) {
	body, err := Get(host, port, user, pass, "/ports/"+objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get port: %w", err)
	}

	var p Port
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, fmt.Errorf("failed to parse port: %w", err)
	}

	return &p, nil
}

func getPortByName(host string, port int, user, pass, name string) (*Port, error) {
	q := fmt.Sprintf("(displayname is %s)", name)
	ports, err := ListPorts(host, port, user, pass, q, 1)
	if err != nil {
		return nil, err
	}
	if len(ports) == 0 {
		return nil, fmt.Errorf("port '%s' not found", name)
	}
	return &ports[0], nil
}

// CreatePort creates a new port
func CreatePort(host string, port int, user, pass string, fields map[string]interface{}) error {
	_, err := Post(host, port, user, pass, "/ports", fields)
	if err != nil {
		return fmt.Errorf("failed to create port: %w", err)
	}
	return nil
}

// UpdatePort updates a port by name or ObjectId
func UpdatePort(host string, port int, user, pass, nameOrID string, fields map[string]interface{}) error {
	p, err := GetPort(host, port, user, pass, nameOrID)
	if err != nil {
		return err
	}
	return Put(host, port, user, pass, "/ports/"+p.ObjectId, fields)
}

// DeletePort deletes a port by name or ObjectId
func DeletePort(host string, port int, user, pass, nameOrID string) error {
	p, err := GetPort(host, port, user, pass, nameOrID)
	if err != nil {
		return err
	}
	return Delete(host, port, user, pass, "/ports/"+p.ObjectId)
}

package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// PortGroup represents a port group
type PortGroup struct {
	ObjectId    string `json:"ObjectId"`
	DisplayName string `json:"DisplayName"`
	MWIOnCode   string `json:"MWIOnCode"`
	MWIOffCode  string `json:"MWIOffCode"`
}

type portGroupsResponse struct {
	Total string      `json:"@total"`
	Items []PortGroup `json:"PortGroup"`
}

// ListPortGroups returns all port groups
func ListPortGroups(host string, port int, user, pass string, query string, rowsPerPage int) ([]PortGroup, error) {
	path := "/portgroups"
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
		return nil, fmt.Errorf("failed to list port groups: %w", err)
	}

	var resp portGroupsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse port groups response: %w", err)
	}

	return resp.Items, nil
}

// GetPortGroup retrieves a port group by name or ObjectId
func GetPortGroup(host string, port int, user, pass, nameOrID string) (*PortGroup, error) {
	if isUUID(nameOrID) {
		return getPortGroupByID(host, port, user, pass, nameOrID)
	}
	return getPortGroupByName(host, port, user, pass, nameOrID)
}

func getPortGroupByID(host string, port int, user, pass, objectID string) (*PortGroup, error) {
	body, err := Get(host, port, user, pass, "/portgroups/"+objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get port group: %w", err)
	}

	var pg PortGroup
	if err := json.Unmarshal(body, &pg); err != nil {
		return nil, fmt.Errorf("failed to parse port group: %w", err)
	}

	return &pg, nil
}

func getPortGroupByName(host string, port int, user, pass, name string) (*PortGroup, error) {
	q := fmt.Sprintf("(displayname is %s)", name)
	groups, err := ListPortGroups(host, port, user, pass, q, 1)
	if err != nil {
		return nil, err
	}
	if len(groups) == 0 {
		return nil, fmt.Errorf("port group '%s' not found", name)
	}
	return &groups[0], nil
}

// UpdatePortGroup updates a port group by name or ObjectId
func UpdatePortGroup(host string, port int, user, pass, nameOrID string, fields map[string]interface{}) error {
	pg, err := GetPortGroup(host, port, user, pass, nameOrID)
	if err != nil {
		return err
	}
	return Put(host, port, user, pass, "/portgroups/"+pg.ObjectId, fields)
}

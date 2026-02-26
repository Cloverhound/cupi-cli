package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// COS represents a CUC Class of Service
type COS struct {
	ObjectId    string `json:"ObjectId"`
	DisplayName string `json:"DisplayName"`
}

type cosListResponse struct {
	Total string `json:"@total"`
	COSes []COS  `json:"Cos"`
}

// ListCOS returns CUC classes of service.
func ListCOS(host string, port int, user, pass string, query string, rowsPerPage int) ([]COS, error) {
	path := "/coses"
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
		return nil, fmt.Errorf("failed to list classes of service: %w", err)
	}

	var resp cosListResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse COS response: %w", err)
	}

	return resp.COSes, nil
}

// GetCOS retrieves a COS by name or ObjectId.
func GetCOS(host string, port int, user, pass string, nameOrID string) (*COS, error) {
	if isUUID(nameOrID) {
		return getCOSByID(host, port, user, pass, nameOrID)
	}
	return getCOSByName(host, port, user, pass, nameOrID)
}

func getCOSByID(host string, port int, user, pass, objectID string) (*COS, error) {
	body, err := Get(host, port, user, pass, "/coses/"+objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get COS: %w", err)
	}
	var c COS
	if err := json.Unmarshal(body, &c); err != nil {
		return nil, fmt.Errorf("failed to parse COS: %w", err)
	}
	return &c, nil
}

func getCOSByName(host string, port int, user, pass, name string) (*COS, error) {
	q := fmt.Sprintf("(displayname is %s)", name)
	coses, err := ListCOS(host, port, user, pass, q, 1)
	if err != nil {
		return nil, err
	}
	if len(coses) == 0 {
		return nil, fmt.Errorf("class of service '%s' not found", name)
	}
	return &coses[0], nil
}

// UpdateCOS updates a COS by name or ObjectId.
func UpdateCOS(host string, port int, user, pass string, nameOrID string, fields map[string]interface{}) error {
	c, err := GetCOS(host, port, user, pass, nameOrID)
	if err != nil {
		return err
	}
	return Put(host, port, user, pass, "/coses/"+c.ObjectId, fields)
}

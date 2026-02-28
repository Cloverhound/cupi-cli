package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// SearchSpace represents a search space
type SearchSpace struct {
	ObjectId    string `json:"ObjectId"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

type searchSpacesResponse struct {
	Total string       `json:"@total"`
	Items OneOrMany[SearchSpace] `json:"SearchSpace"`
}

// SearchSpaceMember represents a partition member in a search space
type SearchSpaceMember struct {
	ObjectId          string `json:"ObjectId"`
	PartitionObjectId string `json:"PartitionObjectId"`
	Index             int    `json:"Index"`
}

type searchSpaceMembersResponse struct {
	Total string              `json:"@total"`
	Items OneOrMany[SearchSpaceMember] `json:"SearchSpaceMember"`
}

// ListSearchSpaces returns all search spaces
func ListSearchSpaces(host string, port int, user, pass string, query string, rowsPerPage int) ([]SearchSpace, error) {
	path := "/searchspaces"
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
		return nil, fmt.Errorf("failed to list search spaces: %w", err)
	}

	var resp searchSpacesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse search spaces response: %w", err)
	}

	return resp.Items, nil
}

// GetSearchSpace retrieves a search space by name or ObjectId
func GetSearchSpace(host string, port int, user, pass, nameOrID string) (*SearchSpace, error) {
	if isUUID(nameOrID) {
		return getSearchSpaceByID(host, port, user, pass, nameOrID)
	}
	return getSearchSpaceByName(host, port, user, pass, nameOrID)
}

func getSearchSpaceByID(host string, port int, user, pass, objectID string) (*SearchSpace, error) {
	body, err := Get(host, port, user, pass, "/searchspaces/"+objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get search space: %w", err)
	}

	var ss SearchSpace
	if err := json.Unmarshal(body, &ss); err != nil {
		return nil, fmt.Errorf("failed to parse search space: %w", err)
	}

	return &ss, nil
}

func getSearchSpaceByName(host string, port int, user, pass, name string) (*SearchSpace, error) {
	q := fmt.Sprintf("(name is %s)", name)
	spaces, err := ListSearchSpaces(host, port, user, pass, q, 1)
	if err != nil {
		return nil, err
	}
	if len(spaces) == 0 {
		return nil, fmt.Errorf("search space '%s' not found", name)
	}
	return &spaces[0], nil
}

// CreateSearchSpace creates a new search space
func CreateSearchSpace(host string, port int, user, pass string, fields map[string]interface{}) error {
	_, err := Post(host, port, user, pass, "/searchspaces", fields)
	if err != nil {
		return fmt.Errorf("failed to create search space: %w", err)
	}
	return nil
}

// UpdateSearchSpace updates a search space by name or ObjectId
func UpdateSearchSpace(host string, port int, user, pass, nameOrID string, fields map[string]interface{}) error {
	ss, err := GetSearchSpace(host, port, user, pass, nameOrID)
	if err != nil {
		return err
	}
	return Put(host, port, user, pass, "/searchspaces/"+ss.ObjectId, fields)
}

// DeleteSearchSpace deletes a search space by name or ObjectId
func DeleteSearchSpace(host string, port int, user, pass, nameOrID string) error {
	ss, err := GetSearchSpace(host, port, user, pass, nameOrID)
	if err != nil {
		return err
	}
	return Delete(host, port, user, pass, "/searchspaces/"+ss.ObjectId)
}

// ListSearchSpaceMembers returns members of a search space
func ListSearchSpaceMembers(host string, port int, user, pass, ssObjectId string) ([]SearchSpaceMember, error) {
	path := fmt.Sprintf("/searchspaces/%s/searchspacemembers", url.PathEscape(ssObjectId))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list search space members: %w", err)
	}

	var resp searchSpaceMembersResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse search space members response: %w", err)
	}

	return resp.Items, nil
}

// AddSearchSpaceMember adds a partition member to a search space
func AddSearchSpaceMember(host string, port int, user, pass, ssObjectId, partitionObjectId string) error {
	path := fmt.Sprintf("/searchspaces/%s/searchspacemembers", url.PathEscape(ssObjectId))
	fields := map[string]interface{}{
		"PartitionObjectId": partitionObjectId,
	}
	_, err := Post(host, port, user, pass, path, fields)
	if err != nil {
		return fmt.Errorf("failed to add search space member: %w", err)
	}
	return nil
}

// RemoveSearchSpaceMember removes a member from a search space
func RemoveSearchSpaceMember(host string, port int, user, pass, ssObjectId, memberObjectId string) error {
	path := fmt.Sprintf("/searchspaces/%s/searchspacemembers/%s", url.PathEscape(ssObjectId), url.PathEscape(memberObjectId))
	if err := Delete(host, port, user, pass, path); err != nil {
		return fmt.Errorf("failed to remove search space member: %w", err)
	}
	return nil
}

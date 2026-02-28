package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// PrivateList represents a user's private distribution list
type PrivateList struct {
	ObjectId    string `json:"ObjectId"`
	DisplayName string `json:"DisplayName"`
	NumericId   string `json:"NumericId"`
}

type privateListsResponse struct {
	Total string        `json:"@total"`
	Items OneOrMany[PrivateList] `json:"PrivateList"`
}

// PrivateListMember represents a member of a private list
type PrivateListMember struct {
	ObjectId       string `json:"ObjectId"`
	MemberObjectId string `json:"MemberObjectId"`
	MemberType     string `json:"MemberType"`
}

type privateListMembersResponse struct {
	Total string               `json:"@total"`
	Items OneOrMany[PrivateListMember] `json:"PrivateListMember"`
}

// ListPrivateLists returns all private lists for a user
func ListPrivateLists(host string, port int, user, pass, userObjectId string) ([]PrivateList, error) {
	path := fmt.Sprintf("/users/%s/privatelists", url.PathEscape(userObjectId))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list private lists: %w", err)
	}

	var resp privateListsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse private lists response: %w", err)
	}

	return resp.Items, nil
}

// GetPrivateList retrieves a specific private list
func GetPrivateList(host string, port int, user, pass, userObjectId, listObjectId string) (*PrivateList, error) {
	path := fmt.Sprintf("/users/%s/privatelists/%s", url.PathEscape(userObjectId), url.PathEscape(listObjectId))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get private list: %w", err)
	}

	var pl PrivateList
	if err := json.Unmarshal(body, &pl); err != nil {
		return nil, fmt.Errorf("failed to parse private list: %w", err)
	}

	return &pl, nil
}

// CreatePrivateList creates a new private list for a user
func CreatePrivateList(host string, port int, user, pass, userObjectId string, fields map[string]interface{}) error {
	path := fmt.Sprintf("/users/%s/privatelists", url.PathEscape(userObjectId))
	_, err := Post(host, port, user, pass, path, fields)
	if err != nil {
		return fmt.Errorf("failed to create private list: %w", err)
	}
	return nil
}

// UpdatePrivateList updates a private list
func UpdatePrivateList(host string, port int, user, pass, userObjectId, listObjectId string, fields map[string]interface{}) error {
	path := fmt.Sprintf("/users/%s/privatelists/%s", url.PathEscape(userObjectId), url.PathEscape(listObjectId))
	if err := Put(host, port, user, pass, path, fields); err != nil {
		return fmt.Errorf("failed to update private list: %w", err)
	}
	return nil
}

// DeletePrivateList deletes a private list
func DeletePrivateList(host string, port int, user, pass, userObjectId, listObjectId string) error {
	path := fmt.Sprintf("/users/%s/privatelists/%s", url.PathEscape(userObjectId), url.PathEscape(listObjectId))
	if err := Delete(host, port, user, pass, path); err != nil {
		return fmt.Errorf("failed to delete private list: %w", err)
	}
	return nil
}

// ListPrivateListMembers returns members of a private list
func ListPrivateListMembers(host string, port int, user, pass, userObjectId, listObjectId string) ([]PrivateListMember, error) {
	path := fmt.Sprintf("/users/%s/privatelists/%s/privatelistmembers", url.PathEscape(userObjectId), url.PathEscape(listObjectId))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list private list members: %w", err)
	}

	var resp privateListMembersResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse private list members response: %w", err)
	}

	return resp.Items, nil
}

// AddPrivateListMember adds a member to a private list
func AddPrivateListMember(host string, port int, user, pass, userObjectId, listObjectId, memberObjectId string) error {
	path := fmt.Sprintf("/users/%s/privatelists/%s/privatelistmembers", url.PathEscape(userObjectId), url.PathEscape(listObjectId))
	fields := map[string]interface{}{
		"MemberObjectId": memberObjectId,
	}
	_, err := Post(host, port, user, pass, path, fields)
	if err != nil {
		return fmt.Errorf("failed to add private list member: %w", err)
	}
	return nil
}

// RemovePrivateListMember removes a member from a private list
func RemovePrivateListMember(host string, port int, user, pass, userObjectId, listObjectId, memberObjectId string) error {
	path := fmt.Sprintf("/users/%s/privatelists/%s/privatelistmembers/%s", url.PathEscape(userObjectId), url.PathEscape(listObjectId), url.PathEscape(memberObjectId))
	if err := Delete(host, port, user, pass, path); err != nil {
		return fmt.Errorf("failed to remove private list member: %w", err)
	}
	return nil
}

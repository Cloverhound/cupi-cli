package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// DistList represents a CUC distribution list
type DistList struct {
	ObjectId    string `json:"ObjectId"`
	Alias       string `json:"Alias"`
	DisplayName string `json:"DisplayName"`
	DtmfAccessId string `json:"DtmfAccessId"`
}

type distListsResponse struct {
	Total     string     `json:"@total"`
	DistLists []DistList `json:"DistributionList"`
}

// DistListMember represents a member of a distribution list
type DistListMember struct {
	ObjectId        string `json:"ObjectId"`
	MemberObjectId  string `json:"MemberObjectId"`
	MemberType      string `json:"MemberType"`
}

type distListMembersResponse struct {
	Total   string           `json:"@total"`
	Members []DistListMember `json:"DistributionListMember"`
}

// ListDistLists returns CUC distribution lists.
func ListDistLists(host string, port int, user, pass string, query string, rowsPerPage int) ([]DistList, error) {
	path := "/distributionlists"
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
		return nil, fmt.Errorf("failed to list distribution lists: %w", err)
	}

	var resp distListsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse distribution lists response: %w", err)
	}

	return resp.DistLists, nil
}

// GetDistList retrieves a distribution list by alias or ObjectId.
func GetDistList(host string, port int, user, pass string, aliasOrID string) (*DistList, error) {
	if isUUID(aliasOrID) {
		return getDistListByID(host, port, user, pass, aliasOrID)
	}
	return getDistListByAlias(host, port, user, pass, aliasOrID)
}

func getDistListByID(host string, port int, user, pass, objectID string) (*DistList, error) {
	body, err := Get(host, port, user, pass, "/distributionlists/"+objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get distribution list: %w", err)
	}
	var dl DistList
	if err := json.Unmarshal(body, &dl); err != nil {
		return nil, fmt.Errorf("failed to parse distribution list: %w", err)
	}
	return &dl, nil
}

func getDistListByAlias(host string, port int, user, pass, alias string) (*DistList, error) {
	q := fmt.Sprintf("(alias is %s)", alias)
	lists, err := ListDistLists(host, port, user, pass, q, 1)
	if err != nil {
		return nil, err
	}
	if len(lists) == 0 {
		return nil, fmt.Errorf("distribution list '%s' not found", alias)
	}
	return &lists[0], nil
}

// CreateDistList creates a new distribution list.
func CreateDistList(host string, port int, user, pass string, fields map[string]interface{}) (*DistList, error) {
	body, err := Post(host, port, user, pass, "/distributionlists", fields)
	if err != nil {
		return nil, fmt.Errorf("failed to create distribution list: %w", err)
	}
	var dl DistList
	if err := json.Unmarshal(body, &dl); err != nil {
		return &DistList{}, nil
	}
	return &dl, nil
}

// UpdateDistList updates a distribution list by alias or ObjectId.
func UpdateDistList(host string, port int, user, pass string, aliasOrID string, fields map[string]interface{}) error {
	dl, err := GetDistList(host, port, user, pass, aliasOrID)
	if err != nil {
		return err
	}
	return Put(host, port, user, pass, "/distributionlists/"+dl.ObjectId, fields)
}

// DeleteDistList deletes a distribution list by alias or ObjectId.
func DeleteDistList(host string, port int, user, pass string, aliasOrID string) error {
	dl, err := GetDistList(host, port, user, pass, aliasOrID)
	if err != nil {
		return err
	}
	return Delete(host, port, user, pass, "/distributionlists/"+dl.ObjectId)
}

// ListDistListMembers returns members of a distribution list.
func ListDistListMembers(host string, port int, user, pass string, aliasOrID string) ([]DistListMember, error) {
	dl, err := GetDistList(host, port, user, pass, aliasOrID)
	if err != nil {
		return nil, err
	}

	body, err := Get(host, port, user, pass, "/distributionlists/"+dl.ObjectId+"/distributionlistmembers")
	if err != nil {
		return nil, fmt.Errorf("failed to list members: %w", err)
	}

	var resp distListMembersResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse members response: %w", err)
	}

	return resp.Members, nil
}

// AddDistListMember adds a member to a distribution list.
func AddDistListMember(host string, port int, user, pass string, listAliasOrID, memberObjectID string) error {
	dl, err := GetDistList(host, port, user, pass, listAliasOrID)
	if err != nil {
		return err
	}
	fields := map[string]interface{}{
		"MemberObjectId": memberObjectID,
	}
	_, err = Post(host, port, user, pass, "/distributionlists/"+dl.ObjectId+"/distributionlistmembers", fields)
	return err
}

// RemoveDistListMember removes a member from a distribution list.
func RemoveDistListMember(host string, port int, user, pass string, listAliasOrID, memberObjectID string) error {
	dl, err := GetDistList(host, port, user, pass, listAliasOrID)
	if err != nil {
		return err
	}
	return Delete(host, port, user, pass, "/distributionlists/"+dl.ObjectId+"/distributionlistmembers/"+memberObjectID)
}

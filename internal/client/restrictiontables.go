package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// RestrictionTable represents a restriction table
type RestrictionTable struct {
	ObjectId    string `json:"ObjectId"`
	DisplayName string `json:"DisplayName"`
}

type restrictionTablesResponse struct {
	Total string             `json:"@total"`
	Items OneOrMany[RestrictionTable] `json:"RestrictionTable"`
}

// RestrictionPattern represents a pattern in a restriction table
type RestrictionPattern struct {
	ObjectId      string `json:"ObjectId"`
	NumberPattern string `json:"NumberPattern"`
	Blocked       string `json:"Blocked"`
	Index         int    `json:"Index"`
}

type restrictionPatternsResponse struct {
	Total string               `json:"@total"`
	Items OneOrMany[RestrictionPattern] `json:"RestrictionPattern"`
}

// ListRestrictionTables returns all restriction tables
func ListRestrictionTables(host string, port int, user, pass string, query string, rowsPerPage int) ([]RestrictionTable, error) {
	path := "/restrictiontables"
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
		return nil, fmt.Errorf("failed to list restriction tables: %w", err)
	}

	var resp restrictionTablesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse restriction tables response: %w", err)
	}

	return resp.Items, nil
}

// GetRestrictionTable retrieves a restriction table by name or ObjectId
func GetRestrictionTable(host string, port int, user, pass, nameOrID string) (*RestrictionTable, error) {
	if isUUID(nameOrID) {
		return getRestrictionTableByID(host, port, user, pass, nameOrID)
	}
	return getRestrictionTableByName(host, port, user, pass, nameOrID)
}

func getRestrictionTableByID(host string, port int, user, pass, objectID string) (*RestrictionTable, error) {
	body, err := Get(host, port, user, pass, "/restrictiontables/"+objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get restriction table: %w", err)
	}

	var rt RestrictionTable
	if err := json.Unmarshal(body, &rt); err != nil {
		return nil, fmt.Errorf("failed to parse restriction table: %w", err)
	}

	return &rt, nil
}

func getRestrictionTableByName(host string, port int, user, pass, name string) (*RestrictionTable, error) {
	q := fmt.Sprintf("(displayname is %s)", name)
	tables, err := ListRestrictionTables(host, port, user, pass, q, 1)
	if err != nil {
		return nil, err
	}
	if len(tables) == 0 {
		return nil, fmt.Errorf("restriction table '%s' not found", name)
	}
	return &tables[0], nil
}

// CreateRestrictionTable creates a new restriction table
func CreateRestrictionTable(host string, port int, user, pass string, fields map[string]interface{}) error {
	_, err := Post(host, port, user, pass, "/restrictiontables", fields)
	if err != nil {
		return fmt.Errorf("failed to create restriction table: %w", err)
	}
	return nil
}

// UpdateRestrictionTable updates a restriction table by name or ObjectId
func UpdateRestrictionTable(host string, port int, user, pass, nameOrID string, fields map[string]interface{}) error {
	rt, err := GetRestrictionTable(host, port, user, pass, nameOrID)
	if err != nil {
		return err
	}
	return Put(host, port, user, pass, "/restrictiontables/"+rt.ObjectId, fields)
}

// DeleteRestrictionTable deletes a restriction table by name or ObjectId
func DeleteRestrictionTable(host string, port int, user, pass, nameOrID string) error {
	rt, err := GetRestrictionTable(host, port, user, pass, nameOrID)
	if err != nil {
		return err
	}
	return Delete(host, port, user, pass, "/restrictiontables/"+rt.ObjectId)
}

// ListRestrictionPatterns returns all patterns in a restriction table
func ListRestrictionPatterns(host string, port int, user, pass, tableObjectId string) ([]RestrictionPattern, error) {
	path := fmt.Sprintf("/restrictiontables/%s/restrictionpatterns", url.PathEscape(tableObjectId))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list restriction patterns: %w", err)
	}

	var resp restrictionPatternsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse restriction patterns response: %w", err)
	}

	return resp.Items, nil
}

// GetRestrictionPattern retrieves a specific restriction pattern
func GetRestrictionPattern(host string, port int, user, pass, tableObjectId, patternObjectId string) (*RestrictionPattern, error) {
	path := fmt.Sprintf("/restrictiontables/%s/restrictionpatterns/%s", url.PathEscape(tableObjectId), url.PathEscape(patternObjectId))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get restriction pattern: %w", err)
	}

	var rp RestrictionPattern
	if err := json.Unmarshal(body, &rp); err != nil {
		return nil, fmt.Errorf("failed to parse restriction pattern: %w", err)
	}

	return &rp, nil
}

// CreateRestrictionPattern creates a new restriction pattern
func CreateRestrictionPattern(host string, port int, user, pass, tableObjectId string, fields map[string]interface{}) error {
	path := fmt.Sprintf("/restrictiontables/%s/restrictionpatterns", url.PathEscape(tableObjectId))
	_, err := Post(host, port, user, pass, path, fields)
	if err != nil {
		return fmt.Errorf("failed to create restriction pattern: %w", err)
	}
	return nil
}

// UpdateRestrictionPattern updates a restriction pattern
func UpdateRestrictionPattern(host string, port int, user, pass, tableObjectId, patternObjectId string, fields map[string]interface{}) error {
	path := fmt.Sprintf("/restrictiontables/%s/restrictionpatterns/%s", url.PathEscape(tableObjectId), url.PathEscape(patternObjectId))
	if err := Put(host, port, user, pass, path, fields); err != nil {
		return fmt.Errorf("failed to update restriction pattern: %w", err)
	}
	return nil
}

// DeleteRestrictionPattern deletes a restriction pattern
func DeleteRestrictionPattern(host string, port int, user, pass, tableObjectId, patternObjectId string) error {
	path := fmt.Sprintf("/restrictiontables/%s/restrictionpatterns/%s", url.PathEscape(tableObjectId), url.PathEscape(patternObjectId))
	if err := Delete(host, port, user, pass, path); err != nil {
		return fmt.Errorf("failed to delete restriction pattern: %w", err)
	}
	return nil
}

package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// Partition represents a dial plan partition
type Partition struct {
	ObjectId    string `json:"ObjectId"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

type partitionsResponse struct {
	Total string      `json:"@total"`
	Items OneOrMany[Partition] `json:"Partition"`
}

// ListPartitions returns all partitions
func ListPartitions(host string, port int, user, pass string, query string, rowsPerPage int) ([]Partition, error) {
	path := "/partitions"
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
		return nil, fmt.Errorf("failed to list partitions: %w", err)
	}

	var resp partitionsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse partitions response: %w", err)
	}

	return resp.Items, nil
}

// GetPartition retrieves a partition by name or ObjectId
func GetPartition(host string, port int, user, pass, nameOrID string) (*Partition, error) {
	if isUUID(nameOrID) {
		return getPartitionByID(host, port, user, pass, nameOrID)
	}
	return getPartitionByName(host, port, user, pass, nameOrID)
}

func getPartitionByID(host string, port int, user, pass, objectID string) (*Partition, error) {
	body, err := Get(host, port, user, pass, "/partitions/"+objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get partition: %w", err)
	}

	var p Partition
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, fmt.Errorf("failed to parse partition: %w", err)
	}

	return &p, nil
}

func getPartitionByName(host string, port int, user, pass, name string) (*Partition, error) {
	q := fmt.Sprintf("(name is %s)", name)
	partitions, err := ListPartitions(host, port, user, pass, q, 1)
	if err != nil {
		return nil, err
	}
	if len(partitions) == 0 {
		return nil, fmt.Errorf("partition '%s' not found", name)
	}
	return &partitions[0], nil
}

// CreatePartition creates a new partition
func CreatePartition(host string, port int, user, pass string, fields map[string]interface{}) error {
	_, err := Post(host, port, user, pass, "/partitions", fields)
	if err != nil {
		return fmt.Errorf("failed to create partition: %w", err)
	}
	return nil
}

// UpdatePartition updates a partition by name or ObjectId
func UpdatePartition(host string, port int, user, pass, nameOrID string, fields map[string]interface{}) error {
	p, err := GetPartition(host, port, user, pass, nameOrID)
	if err != nil {
		return err
	}
	return Put(host, port, user, pass, "/partitions/"+p.ObjectId, fields)
}

// DeletePartition deletes a partition by name or ObjectId
func DeletePartition(host string, port int, user, pass, nameOrID string) error {
	p, err := GetPartition(host, port, user, pass, nameOrID)
	if err != nil {
		return err
	}
	return Delete(host, port, user, pass, "/partitions/"+p.ObjectId)
}

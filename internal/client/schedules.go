package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// Schedule represents a CUC schedule
type Schedule struct {
	ObjectId    string `json:"ObjectId"`
	DisplayName string `json:"DisplayName"`
	IsHoliday   string `json:"IsHoliday"`
}

type schedulesResponse struct {
	Total     string     `json:"@total"`
	Schedules []Schedule `json:"Schedule"`
}

// ListSchedules returns CUC schedules.
func ListSchedules(host string, port int, user, pass string, query string, rowsPerPage int) ([]Schedule, error) {
	path := "/schedules"
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
		return nil, fmt.Errorf("failed to list schedules: %w", err)
	}

	var resp schedulesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse schedules response: %w", err)
	}

	return resp.Schedules, nil
}

// GetSchedule retrieves a schedule by name or ObjectId.
func GetSchedule(host string, port int, user, pass string, nameOrID string) (*Schedule, error) {
	if isUUID(nameOrID) {
		body, err := Get(host, port, user, pass, "/schedules/"+nameOrID)
		if err != nil {
			return nil, fmt.Errorf("failed to get schedule: %w", err)
		}
		var s Schedule
		if err := json.Unmarshal(body, &s); err != nil {
			return nil, fmt.Errorf("failed to parse schedule: %w", err)
		}
		return &s, nil
	}

	q := fmt.Sprintf("(displayname is %s)", nameOrID)
	schedules, err := ListSchedules(host, port, user, pass, q, 1)
	if err != nil {
		return nil, err
	}
	if len(schedules) == 0 {
		return nil, fmt.Errorf("schedule '%s' not found", nameOrID)
	}
	return &schedules[0], nil
}

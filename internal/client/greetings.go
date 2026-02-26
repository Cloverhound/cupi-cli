package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// Greeting represents a call handler greeting
type Greeting struct {
	GreetingType      string `json:"GreetingType"`
	Enabled           string `json:"Enabled"`
	PlayWhat          string `json:"PlayWhat"`
	TimeExpiresSetFor string `json:"TimeExpiresSetFor"`
}

type greetingsResponse struct {
	Total string     `json:"@total"`
	Items []Greeting `json:"Greeting"`
}

// ListGreetings returns all greetings for a call handler
func ListGreetings(host string, port int, user, pass, handlerObjectId string) ([]Greeting, error) {
	path := fmt.Sprintf("/handlers/callhandlers/%s/greetings", url.PathEscape(handlerObjectId))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list greetings: %w", err)
	}

	var resp greetingsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse greetings response: %w", err)
	}

	return resp.Items, nil
}

// GetGreeting retrieves a specific greeting by type
func GetGreeting(host string, port int, user, pass, handlerObjectId, greetingType string) (*Greeting, error) {
	path := fmt.Sprintf("/handlers/callhandlers/%s/greetings/%s", url.PathEscape(handlerObjectId), url.PathEscape(greetingType))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get greeting: %w", err)
	}

	var g Greeting
	if err := json.Unmarshal(body, &g); err != nil {
		return nil, fmt.Errorf("failed to parse greeting: %w", err)
	}

	return &g, nil
}

// UpdateGreeting updates a call handler greeting
func UpdateGreeting(host string, port int, user, pass, handlerObjectId, greetingType string, fields map[string]interface{}) error {
	path := fmt.Sprintf("/handlers/callhandlers/%s/greetings/%s", url.PathEscape(handlerObjectId), url.PathEscape(greetingType))
	if err := Put(host, port, user, pass, path, fields); err != nil {
		return fmt.Errorf("failed to update greeting: %w", err)
	}
	return nil
}

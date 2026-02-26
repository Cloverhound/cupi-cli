package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// CallHandler represents a CUC call handler
type CallHandler struct {
	ObjectId    string `json:"ObjectId"`
	DisplayName string `json:"DisplayName"`
	DtmfAccessId string `json:"DtmfAccessId"`
	IsPrimary   string `json:"IsPrimary"`
}

type callHandlersResponse struct {
	Total    string        `json:"@total"`
	Handlers []CallHandler `json:"Callhandler"`
}

// ListCallHandlers returns CUC call handlers.
func ListCallHandlers(host string, port int, user, pass string, query string, rowsPerPage int) ([]CallHandler, error) {
	path := "/handlers/callhandlers"
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
		return nil, fmt.Errorf("failed to list call handlers: %w", err)
	}

	var resp callHandlersResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse call handlers response: %w", err)
	}

	return resp.Handlers, nil
}

// GetCallHandler retrieves a call handler by DisplayName or ObjectId.
func GetCallHandler(host string, port int, user, pass string, nameOrID string) (*CallHandler, error) {
	if isUUID(nameOrID) {
		return getCallHandlerByID(host, port, user, pass, nameOrID)
	}
	return getCallHandlerByName(host, port, user, pass, nameOrID)
}

func getCallHandlerByID(host string, port int, user, pass, objectID string) (*CallHandler, error) {
	body, err := Get(host, port, user, pass, "/handlers/callhandlers/"+objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get call handler: %w", err)
	}
	var h CallHandler
	if err := json.Unmarshal(body, &h); err != nil {
		return nil, fmt.Errorf("failed to parse call handler: %w", err)
	}
	return &h, nil
}

func getCallHandlerByName(host string, port int, user, pass, name string) (*CallHandler, error) {
	q := fmt.Sprintf("(displayname is %s)", name)
	handlers, err := ListCallHandlers(host, port, user, pass, q, 1)
	if err != nil {
		return nil, err
	}
	if len(handlers) == 0 {
		return nil, fmt.Errorf("call handler '%s' not found", name)
	}
	return &handlers[0], nil
}

// CreateCallHandler creates a new call handler. templateObjectID is required.
func CreateCallHandler(host string, port int, user, pass string, templateObjectID string, fields map[string]interface{}) (*CallHandler, error) {
	path := fmt.Sprintf("/handlers/callhandlers?templateObjectId=%s", url.QueryEscape(templateObjectID))
	body, err := Post(host, port, user, pass, path, fields)
	if err != nil {
		return nil, fmt.Errorf("failed to create call handler: %w", err)
	}
	var h CallHandler
	if err := json.Unmarshal(body, &h); err != nil {
		return &CallHandler{}, nil
	}
	return &h, nil
}

// UpdateCallHandler updates a call handler by name or ObjectId.
func UpdateCallHandler(host string, port int, user, pass string, nameOrID string, fields map[string]interface{}) error {
	h, err := GetCallHandler(host, port, user, pass, nameOrID)
	if err != nil {
		return err
	}
	return Put(host, port, user, pass, "/handlers/callhandlers/"+h.ObjectId, fields)
}

// DeleteCallHandler deletes a call handler by name or ObjectId.
func DeleteCallHandler(host string, port int, user, pass string, nameOrID string) error {
	h, err := GetCallHandler(host, port, user, pass, nameOrID)
	if err != nil {
		return err
	}
	return Delete(host, port, user, pass, "/handlers/callhandlers/"+h.ObjectId)
}

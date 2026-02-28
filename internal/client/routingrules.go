package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// RoutingRule represents a dial plan routing rule
type RoutingRule struct {
	ObjectId    string `json:"ObjectId"`
	DisplayName string `json:"DisplayName"`
	Type        string `json:"Type"`
	RouteAction string `json:"RouteAction"`
	Enabled     string `json:"Enabled"`
	Index       int    `json:"Index"`
}

type routingRulesResponse struct {
	Total string        `json:"@total"`
	Items OneOrMany[RoutingRule] `json:"RoutingRule"`
}

// RoutingRuleCondition represents a condition in a routing rule
type RoutingRuleCondition struct {
	ObjectId     string `json:"ObjectId"`
	OperatorType string `json:"OperatorType"`
	Parameter    string `json:"Parameter"`
	OperandTwo   string `json:"OperandTwo"`
}

type routingRuleConditionsResponse struct {
	Total string                  `json:"@total"`
	Items OneOrMany[RoutingRuleCondition] `json:"RoutingRuleCondition"`
}

// ListRoutingRules returns all routing rules
func ListRoutingRules(host string, port int, user, pass string, query string, rowsPerPage int) ([]RoutingRule, error) {
	path := "/routingrules"
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
		return nil, fmt.Errorf("failed to list routing rules: %w", err)
	}

	var resp routingRulesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse routing rules response: %w", err)
	}

	return resp.Items, nil
}

// GetRoutingRule retrieves a routing rule by name or ObjectId
func GetRoutingRule(host string, port int, user, pass, nameOrID string) (*RoutingRule, error) {
	if isUUID(nameOrID) {
		return getRoutingRuleByID(host, port, user, pass, nameOrID)
	}
	return getRoutingRuleByName(host, port, user, pass, nameOrID)
}

func getRoutingRuleByID(host string, port int, user, pass, objectID string) (*RoutingRule, error) {
	body, err := Get(host, port, user, pass, "/routingrules/"+objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get routing rule: %w", err)
	}

	var rr RoutingRule
	if err := json.Unmarshal(body, &rr); err != nil {
		return nil, fmt.Errorf("failed to parse routing rule: %w", err)
	}

	return &rr, nil
}

func getRoutingRuleByName(host string, port int, user, pass, name string) (*RoutingRule, error) {
	q := fmt.Sprintf("(displayname is %s)", name)
	rules, err := ListRoutingRules(host, port, user, pass, q, 1)
	if err != nil {
		return nil, err
	}
	if len(rules) == 0 {
		return nil, fmt.Errorf("routing rule '%s' not found", name)
	}
	return &rules[0], nil
}

// CreateRoutingRule creates a new routing rule
func CreateRoutingRule(host string, port int, user, pass string, fields map[string]interface{}) error {
	_, err := Post(host, port, user, pass, "/routingrules", fields)
	if err != nil {
		return fmt.Errorf("failed to create routing rule: %w", err)
	}
	return nil
}

// UpdateRoutingRule updates a routing rule by name or ObjectId
func UpdateRoutingRule(host string, port int, user, pass, nameOrID string, fields map[string]interface{}) error {
	rr, err := GetRoutingRule(host, port, user, pass, nameOrID)
	if err != nil {
		return err
	}
	return Put(host, port, user, pass, "/routingrules/"+rr.ObjectId, fields)
}

// DeleteRoutingRule deletes a routing rule by name or ObjectId
func DeleteRoutingRule(host string, port int, user, pass, nameOrID string) error {
	rr, err := GetRoutingRule(host, port, user, pass, nameOrID)
	if err != nil {
		return err
	}
	return Delete(host, port, user, pass, "/routingrules/"+rr.ObjectId)
}

// ListRoutingRuleConditions returns conditions for a routing rule
func ListRoutingRuleConditions(host string, port int, user, pass, ruleObjectId string) ([]RoutingRuleCondition, error) {
	path := fmt.Sprintf("/routingrules/%s/routingruleconditions", url.PathEscape(ruleObjectId))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list routing rule conditions: %w", err)
	}

	var resp routingRuleConditionsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse routing rule conditions response: %w", err)
	}

	return resp.Items, nil
}

// GetRoutingRuleCondition retrieves a specific routing rule condition
func GetRoutingRuleCondition(host string, port int, user, pass, ruleObjectId, conditionObjectId string) (*RoutingRuleCondition, error) {
	path := fmt.Sprintf("/routingrules/%s/routingruleconditions/%s", url.PathEscape(ruleObjectId), url.PathEscape(conditionObjectId))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get routing rule condition: %w", err)
	}

	var c RoutingRuleCondition
	if err := json.Unmarshal(body, &c); err != nil {
		return nil, fmt.Errorf("failed to parse routing rule condition: %w", err)
	}

	return &c, nil
}

// CreateRoutingRuleCondition creates a new condition for a routing rule
func CreateRoutingRuleCondition(host string, port int, user, pass, ruleObjectId string, fields map[string]interface{}) error {
	path := fmt.Sprintf("/routingrules/%s/routingruleconditions", url.PathEscape(ruleObjectId))
	_, err := Post(host, port, user, pass, path, fields)
	if err != nil {
		return fmt.Errorf("failed to create routing rule condition: %w", err)
	}
	return nil
}

// UpdateRoutingRuleCondition updates a routing rule condition
func UpdateRoutingRuleCondition(host string, port int, user, pass, ruleObjectId, conditionObjectId string, fields map[string]interface{}) error {
	path := fmt.Sprintf("/routingrules/%s/routingruleconditions/%s", url.PathEscape(ruleObjectId), url.PathEscape(conditionObjectId))
	if err := Put(host, port, user, pass, path, fields); err != nil {
		return fmt.Errorf("failed to update routing rule condition: %w", err)
	}
	return nil
}

// DeleteRoutingRuleCondition deletes a routing rule condition
func DeleteRoutingRuleCondition(host string, port int, user, pass, ruleObjectId, conditionObjectId string) error {
	path := fmt.Sprintf("/routingrules/%s/routingruleconditions/%s", url.PathEscape(ruleObjectId), url.PathEscape(conditionObjectId))
	if err := Delete(host, port, user, pass, path); err != nil {
		return fmt.Errorf("failed to delete routing rule condition: %w", err)
	}
	return nil
}

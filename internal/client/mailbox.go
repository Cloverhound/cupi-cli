package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// MailboxAttributes represents mailbox configuration and statistics
type MailboxAttributes struct {
	ObjectId       string `json:"ObjectId"`
	IsReceiveEnabled string `json:"IsReceiveEnabled"`
	ReceiveQuota   string `json:"ReceiveQuota"`
	SendQuota      string `json:"SendQuota"`
	QuotaWarning   string `json:"QuotaWarning"`
	MessageCount   string `json:"MessageCount"`
	TotalByteSize  string `json:"TotalByteSize"`
}

// GetMailboxAttributes retrieves mailbox attributes for a user
func GetMailboxAttributes(host string, port int, user, pass, userObjectId string) (*MailboxAttributes, error) {
	path := fmt.Sprintf("/users/%s/mailboxattributes", url.PathEscape(userObjectId))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get mailbox attributes: %w", err)
	}

	var ma MailboxAttributes
	if err := json.Unmarshal(body, &ma); err != nil {
		return nil, fmt.Errorf("failed to parse mailbox attributes: %w", err)
	}

	return &ma, nil
}

// UpdateMailboxAttributes updates mailbox attributes for a user
func UpdateMailboxAttributes(host string, port int, user, pass, userObjectId string, fields map[string]interface{}) error {
	path := fmt.Sprintf("/users/%s/mailboxattributes", url.PathEscape(userObjectId))
	if err := Put(host, port, user, pass, path, fields); err != nil {
		return fmt.Errorf("failed to update mailbox attributes: %w", err)
	}
	return nil
}

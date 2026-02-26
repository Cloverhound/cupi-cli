package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// Credential represents user credentials (PIN or password)
type Credential struct {
	ObjectId      string `json:"ObjectId"`
	CredentialType string `json:"CredentialType"`
	Locked        string `json:"Locked"`
	HackCount     string `json:"HackCount"`
	IsMustChange  string `json:"IsMustChange"`
	DoesntExpire  string `json:"DoesntExpire"`
	NeverExpires  string `json:"NeverExpires"`
}

// GetCredential retrieves credential info for a user (PIN or password)
func GetCredential(host string, port int, user, pass, userObjectId, credType string) (*Credential, error) {
	path := fmt.Sprintf("/users/%s/credential/%s", url.PathEscape(userObjectId), url.PathEscape(credType))
	body, err := Get(host, port, user, pass, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get credential: %w", err)
	}

	var c Credential
	if err := json.Unmarshal(body, &c); err != nil {
		return nil, fmt.Errorf("failed to parse credential: %w", err)
	}

	return &c, nil
}

// UpdateCredential updates credential fields
func UpdateCredential(host string, port int, user, pass, userObjectId, credType string, fields map[string]interface{}) error {
	path := fmt.Sprintf("/users/%s/credential/%s", url.PathEscape(userObjectId), url.PathEscape(credType))
	if err := Put(host, port, user, pass, path, fields); err != nil {
		return fmt.Errorf("failed to update credential: %w", err)
	}
	return nil
}

// UnlockCredential unlocks a credential by clearing lock and hack count
func UnlockCredential(host string, port int, user, pass, userObjectId, credType string) error {
	fields := map[string]interface{}{
		"Locked":    "false",
		"HackCount": "0",
	}
	return UpdateCredential(host, port, user, pass, userObjectId, credType, fields)
}

// SetCredential sets a new credential value (PIN or password)
func SetCredential(host string, port int, user, pass, userObjectId, credType, newValue string) error {
	fields := map[string]interface{}{
		"Credentials": newValue,
	}
	return UpdateCredential(host, port, user, pass, userObjectId, credType, fields)
}

package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zalando/go-keyring"
)

const (
	CredTypeCUPI        = "cupi"
	CredTypeApplication = "application"
	CredTypePlatform    = "platform"

	keyringSvc = "cupi-cli"
)

// KeyringKey returns the account name used in the OS keystore.
func KeyringKey(host, credType string) string {
	return fmt.Sprintf("%s:%s", host, credType)
}

// StorePassword saves a password in the OS keystore.
func StorePassword(host, credType, password string) error {
	if err := keyring.Set(keyringSvc, KeyringKey(host, credType), password); err != nil {
		return fmt.Errorf("failed to store credentials in OS keystore: %w", err)
	}
	return nil
}

// GetPassword retrieves a password from the OS keystore.
// If not found, falls back to legacy JSON file and migrates automatically.
func GetPassword(host, credType string) (string, error) {
	key := KeyringKey(host, credType)

	password, err := keyring.Get(keyringSvc, key)
	if err == nil {
		return password, nil
	}

	// Not in keystore — check legacy JSON file and migrate if found
	if legacy, legacyErr := getLegacyPassword(host, credType); legacyErr == nil {
		if setErr := keyring.Set(keyringSvc, key, legacy); setErr == nil {
			deleteLegacyPassword(host, credType)
			return legacy, nil
		}
		return legacy, nil
	}

	return "", fmt.Errorf("credentials not found: run 'cupi auth login' to authenticate")
}

// DeletePassword removes a password from the OS keystore.
func DeletePassword(host, credType string) error {
	err := keyring.Delete(keyringSvc, KeyringKey(host, credType))
	deleteLegacyPassword(host, credType)
	if err != nil && err != keyring.ErrNotFound {
		return fmt.Errorf("failed to delete credentials from OS keystore: %w", err)
	}
	return nil
}

// DeleteAllPasswords removes all credential types for the given host.
func DeleteAllPasswords(host string) error {
	for _, credType := range []string{CredTypeCUPI, CredTypeApplication, CredTypePlatform} {
		_ = DeletePassword(host, credType)
	}
	return nil
}

// --- Legacy JSON file migration helpers ---

func legacyStorePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cupi-cli", ".credentials")
}

func loadLegacyStore() (map[string]string, error) {
	data, err := os.ReadFile(legacyStorePath())
	if err != nil {
		return nil, err
	}
	var store map[string]string
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	return store, nil
}

func getLegacyPassword(host, credType string) (string, error) {
	store, err := loadLegacyStore()
	if err != nil {
		return "", err
	}
	key := fmt.Sprintf("cupi-cli:%s:%s", host, credType)
	password, ok := store[key]
	if !ok {
		return "", fmt.Errorf("not found")
	}
	return password, nil
}

func deleteLegacyPassword(host, credType string) {
	store, err := loadLegacyStore()
	if err != nil {
		return
	}
	key := fmt.Sprintf("cupi-cli:%s:%s", host, credType)
	delete(store, key)
	if len(store) == 0 {
		os.Remove(legacyStorePath())
		return
	}
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(legacyStorePath(), append(data, '\n'), 0600)
}

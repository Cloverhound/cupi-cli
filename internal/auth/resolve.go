package auth

import (
	"fmt"

	"github.com/Cloverhound/cupi-cli/internal/appconfig"
)

// ResolveCreds returns (username, password, error) for the given server and credType.
// Looks up username from config, password from keyring.
func ResolveCreds(server *appconfig.ServerConfig, credType string) (string, string, error) {
	cred, ok := server.Credentials[credType]
	if !ok {
		return "", "", fmt.Errorf("credential type '%s' not configured for server", credType)
	}

	username := cred.Username
	if username == "" {
		return "", "", fmt.Errorf("username not configured for credential type '%s'", credType)
	}

	password, err := GetPassword(server.Host, credType)
	if err != nil {
		return "", "", fmt.Errorf("failed to retrieve password from keyring: %w", err)
	}

	return username, password, nil
}

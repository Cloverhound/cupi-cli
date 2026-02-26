package discovery

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Cloverhound/cupi-cli/internal/client"
)

// TestCUPIAuth verifies credentials by GETting /vmrest/users?rowsPerPage=0
// Returns (version, error). Version may be empty if not in response.
func TestCUPIAuth(host string, port int, user, pass string) (string, error) {
	if port == 0 {
		port = 443
	}
	url := fmt.Sprintf("https://%s:%d/vmrest/users?rowsPerPage=0", host, port)

	httpClient := client.NewHTTPClient()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(user, pass)
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to connect to %s: %w", host, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("authentication failed: invalid credentials (HTTP 401)")
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("connectivity check failed: HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	// Try to extract @version if present
	var versionResp struct {
		Version string `json:"@version"`
	}
	if err := json.Unmarshal(body, &versionResp); err == nil && versionResp.Version != "" {
		return versionResp.Version, nil
	}

	return "", nil
}

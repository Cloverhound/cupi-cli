package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

// ErrDryRun is returned when a mutating request is intercepted by dry-run mode.
var ErrDryRun = errors.New("dry-run: request not sent")

// cupiBase returns the base URL for CUPI REST API
func cupiBase(host string, port int) string {
	if port == 0 {
		port = 443
	}
	return fmt.Sprintf("https://%s:%d/vmrest", host, port)
}

// Request makes an authenticated HTTP request to the CUPI REST API.
// method: GET, POST, PUT, DELETE
// path: e.g. "/users" or "/users/{objectId}"
// body: optional JSON body (for POST/PUT); pass nil for GET/DELETE
func Request(host string, port int, user, pass, method, path string, body interface{}) ([]byte, int, error) {
	base := cupiBase(host, port)
	url := base + path

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	// Dry-run intercept for mutating operations
	if os.Getenv("CUPI_DRY_RUN") != "" && (method == "POST" || method == "PUT" || method == "DELETE") {
		fmt.Fprintf(os.Stderr, "[DRY RUN] %s %s\n", method, url)
		if body != nil {
			b, _ := json.MarshalIndent(body, "", "  ")
			fmt.Fprintf(os.Stderr, "[DRY RUN] Body: %s\n", string(b))
		}
		return nil, 0, ErrDryRun
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(user, pass)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if os.Getenv("CUPI_DEBUG") != "" {
		fmt.Fprintf(os.Stderr, "=== CUPI Request ===\n%s %s\n", method, url)
		if body != nil {
			b, _ := json.MarshalIndent(body, "", "  ")
			fmt.Fprintf(os.Stderr, "Body: %s\n", string(b))
		}
	}

	httpClient := NewHTTPClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response: %w", err)
	}

	if os.Getenv("CUPI_DEBUG") != "" {
		fmt.Fprintf(os.Stderr, "=== CUPI Response HTTP %d ===\n%s\n", resp.StatusCode, string(respBody))
	}

	return respBody, resp.StatusCode, nil
}

// Get makes an authenticated GET request to the CUPI REST API.
func Get(host string, port int, user, pass, path string) ([]byte, error) {
	body, status, err := Request(host, port, user, pass, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, cupiError(status, body)
	}
	return body, nil
}

// Post makes an authenticated POST request to the CUPI REST API.
func Post(host string, port int, user, pass, path string, body interface{}) ([]byte, error) {
	respBody, status, err := Request(host, port, user, pass, "POST", path, body)
	if err != nil {
		return nil, err
	}
	if status != http.StatusCreated && status != http.StatusOK {
		return nil, cupiError(status, respBody)
	}
	return respBody, nil
}

// Put makes an authenticated PUT request to the CUPI REST API.
func Put(host string, port int, user, pass, path string, body interface{}) error {
	respBody, status, err := Request(host, port, user, pass, "PUT", path, body)
	if err != nil {
		return err
	}
	if status != http.StatusNoContent && status != http.StatusOK {
		return cupiError(status, respBody)
	}
	return nil
}

// Delete makes an authenticated DELETE request to the CUPI REST API.
func Delete(host string, port int, user, pass, path string) error {
	respBody, status, err := Request(host, port, user, pass, "DELETE", path, nil)
	if err != nil {
		return err
	}
	if status != http.StatusNoContent && status != http.StatusOK {
		return cupiError(status, respBody)
	}
	return nil
}

// cupiError parses a CUPI error response and returns a formatted error.
// CUPI error format: {"errors":{"code":"...","error":[...]}}
func cupiError(status int, body []byte) error {
	var errResp struct {
		Errors struct {
			Code  string   `json:"code"`
			Error []string `json:"error"`
		} `json:"errors"`
	}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &errResp); err == nil && len(errResp.Errors.Error) > 0 {
			return fmt.Errorf("HTTP %d: %s", status, errResp.Errors.Error[0])
		}
	}
	msg := string(body)
	if len(msg) > 200 {
		msg = msg[:200] + "..."
	}
	return fmt.Errorf("HTTP %d: %s", status, msg)
}

// Integration tests for cupi-cli.
// These tests require a live Cisco Unity Connection server.
//
// Required environment variables:
//
//	CUPI_TEST_HOST   CUC hostname or IP
//	CUPI_TEST_USER   CUPI admin username
//	CUPI_TEST_PASS   CUPI admin password
//
// Optional:
//
//	CUPI_TEST_PORT   Port (default: 443)
//
// Run:
//
//	CUPI_TEST_HOST=cuc.example.com CUPI_TEST_USER=admin CUPI_TEST_PASS=secret \
//	  go test ./tests/ -v -timeout 120s
package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/Cloverhound/cupi-cli/internal/client"
)

type testEnv struct {
	Host string
	Port int
	User string
	Pass string
}

// getTestEnv reads credentials from env vars and skips the test if not set.
func getTestEnv(t *testing.T) testEnv {
	t.Helper()
	host := os.Getenv("CUPI_TEST_HOST")
	user := os.Getenv("CUPI_TEST_USER")
	pass := os.Getenv("CUPI_TEST_PASS")
	portStr := os.Getenv("CUPI_TEST_PORT")

	if host == "" || user == "" || pass == "" {
		t.Skip("set CUPI_TEST_HOST, CUPI_TEST_USER, CUPI_TEST_PASS to run integration tests")
	}

	port := 443
	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}
	return testEnv{Host: host, Port: port, User: user, Pass: pass}
}

// uniqueSuffix returns a short timestamp-based suffix for unique test object names.
func uniqueSuffix() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()%1_000_000_000)
}

// getFirstCallHandlerTemplateID retrieves the ObjectId of the first available
// call handler template. Skips the test if none are found.
func getFirstCallHandlerTemplateID(t *testing.T, env testEnv) string {
	t.Helper()
	body, err := client.Get(env.Host, env.Port, env.User, env.Pass, "/callhandlertemplates?rowsPerPage=1")
	if err != nil {
		t.Skipf("could not list call handler templates: %v", err)
	}
	var resp struct {
		Templates []struct {
			ObjectId string `json:"ObjectId"`
		} `json:"CallhandlerTemplate"`
	}
	if err := json.Unmarshal(body, &resp); err != nil || len(resp.Templates) == 0 {
		t.Skip("no call handler templates found on this server")
	}
	return resp.Templates[0].ObjectId
}

// verifyConnectivity does a lightweight auth check — fails fast with a clear message.
func verifyConnectivity(t *testing.T, env testEnv) {
	t.Helper()
	_, err := client.Get(env.Host, env.Port, env.User, env.Pass, "/users?rowsPerPage=0")
	if err != nil {
		t.Fatalf("connectivity check failed (check host/credentials): %v", err)
	}
}

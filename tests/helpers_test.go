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
// If CUPI_TEST_SUFFIX is set, that value is used instead — enabling a two-phase
// run where Add/Modify and Delete are invoked separately but share the same alias.
func uniqueSuffix() string {
	if s := os.Getenv("CUPI_TEST_SUFFIX"); s != "" {
		return s
	}
	return fmt.Sprintf("%d", time.Now().UnixNano()%1_000_000_000)
}

// getFirstCallHandlerTemplateID retrieves the ObjectId of the first available
// call handler template. Skips the test if none are found.
func getFirstCallHandlerTemplateID(t *testing.T, env testEnv) string {
	t.Helper()
	templates, err := client.ListCallHandlerTemplates(env.Host, env.Port, env.User, env.Pass, "", 1)
	if err != nil {
		t.Skipf("could not list call handler templates: %v", err)
	}
	if len(templates) == 0 {
		t.Skip("no call handler templates found on this server")
	}
	return templates[0].ObjectId
}

// getDistListMemberObjectID returns the ObjectId of a system distribution list
// (e.g. allvoicemailusers) that can be used as a member of the test list.
// Using a distlist as a member avoids user-license constraints. Skips if none found.
func getDistListMemberObjectID(t *testing.T, env testEnv) string {
	t.Helper()
	lists, err := client.ListDistLists(env.Host, env.Port, env.User, env.Pass, "", 10)
	if err != nil {
		t.Skipf("could not list distribution lists: %v", err)
	}
	// Skip the test list itself (prefix cupi-tl-) and return the first system list.
	for _, dl := range lists {
		if len(dl.Alias) > 0 && dl.Alias[:7] != "cupi-tl" {
			return dl.ObjectId
		}
	}
	t.Skip("no system distribution list available to use as distlist member")
	return ""
}

// verifyConnectivity does a lightweight auth check — fails fast with a clear message.
func verifyConnectivity(t *testing.T, env testEnv) {
	t.Helper()
	_, err := client.Get(env.Host, env.Port, env.User, env.Pass, "/users?rowsPerPage=0")
	if err != nil {
		t.Fatalf("connectivity check failed (check host/credentials): %v", err)
	}
}

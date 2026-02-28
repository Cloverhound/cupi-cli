package tests

import (
	"strings"
	"testing"

	"github.com/Cloverhound/cupi-cli/internal/client"
)

// TestUsersAddModifyDelete creates a mailbox user, updates it, and deletes it.
func TestUsersAddModifyDelete(t *testing.T) {
	env := getTestEnv(t)
	verifyConnectivity(t, env)

	alias := "cupi-test-" + uniqueSuffix()
	t.Logf("test user alias: %s", alias)

	// --- Add ---
	t.Run("Add", func(t *testing.T) {
		fields := map[string]interface{}{
			"Alias":        alias,
			"DtmfAccessId": "79999",
			"FirstName":    "CUPITest",
			"LastName":     "User",
			"DisplayName":  "CUPITest User",
		}
		u, err := client.CreateUser(env.Host, env.Port, env.User, env.Pass, "voicemailusertemplate", fields)
		if err != nil {
			t.Fatalf("CreateUser: %v", err)
		}
		t.Logf("created user ObjectId: %s alias: %s", u.ObjectId, u.Alias)

		// Verify it exists
		got, err := client.GetUser(env.Host, env.Port, env.User, env.Pass, alias)
		if err != nil {
			t.Fatalf("GetUser after create: %v", err)
		}
		if got.Alias != alias {
			t.Errorf("alias mismatch: want %q got %q", alias, got.Alias)
		}
	})

	// --- Modify ---
	t.Run("Modify", func(t *testing.T) {
		updates := map[string]interface{}{
			"DisplayName": "CUPITest User Updated",
			"Department":  "TestDept",
			"Title":       "TestTitle",
		}
		if err := client.UpdateUser(env.Host, env.Port, env.User, env.Pass, alias, updates); err != nil {
			t.Fatalf("UpdateUser: %v", err)
		}

		got, err := client.GetUser(env.Host, env.Port, env.User, env.Pass, alias)
		if err != nil {
			t.Fatalf("GetUser after update: %v", err)
		}
		if got.DisplayName != "CUPITest User Updated" {
			t.Errorf("DisplayName: want %q got %q", "CUPITest User Updated", got.DisplayName)
		}
		if got.Department != "TestDept" {
			t.Errorf("Department: want %q got %q", "TestDept", got.Department)
		}
		// Title is accepted by the PUT but may be silently ignored on some CUC
		// versions/configurations (e.g. LDAP-sync or restricted COS). Log only.
		if got.Title != "TestTitle" {
			t.Logf("Title not persisted (server may ignore this field): got %q", got.Title)
		}
	})

	// --- Delete ---
	t.Run("Delete", func(t *testing.T) {
		if err := client.DeleteUser(env.Host, env.Port, env.User, env.Pass, alias); err != nil {
			t.Fatalf("DeleteUser: %v", err)
		}

		// Verify deletion
		_, err := client.GetUser(env.Host, env.Port, env.User, env.Pass, alias)
		if err == nil {
			t.Errorf("expected error after deletion but GetUser succeeded")
		}
		t.Logf("confirmed deleted: %v", err)
	})
}

// TestUsersList verifies listing returns results and supports a query filter.
func TestUsersList(t *testing.T) {
	env := getTestEnv(t)
	verifyConnectivity(t, env)

	t.Run("ListAll", func(t *testing.T) {
		users, err := client.ListUsers(env.Host, env.Port, env.User, env.Pass, "", 10)
		if err != nil {
			t.Fatalf("ListUsers: %v", err)
		}
		if len(users) == 0 {
			t.Fatal("expected at least one user, got none")
		}
		t.Logf("ListUsers returned %d users (max 10)", len(users))
	})

	t.Run("ListWithQuery", func(t *testing.T) {
		// Query for users whose alias starts with a common prefix — may return 0 but must not error
		users, err := client.ListUsers(env.Host, env.Port, env.User, env.Pass, "(alias startswith a)", 5)
		if err != nil {
			t.Fatalf("ListUsers with query: %v", err)
		}
		t.Logf("query '(alias startswith a)' returned %d users", len(users))
		for _, u := range users {
			if !strings.HasPrefix(strings.ToLower(u.Alias), "a") {
				t.Errorf("result alias %q doesn't match query prefix 'a'", u.Alias)
			}
		}
	})
}

// TestUsersGetByAlias and GetByObjectId round-trip.
func TestUsersGetRoundtrip(t *testing.T) {
	env := getTestEnv(t)
	verifyConnectivity(t, env)

	// List to get a known user
	users, err := client.ListUsers(env.Host, env.Port, env.User, env.Pass, "", 1)
	if err != nil || len(users) == 0 {
		t.Skip("no users available for round-trip test")
	}
	first := users[0]

	t.Run("GetByAlias", func(t *testing.T) {
		got, err := client.GetUser(env.Host, env.Port, env.User, env.Pass, first.Alias)
		if err != nil {
			t.Fatalf("GetUser by alias: %v", err)
		}
		if got.ObjectId != first.ObjectId {
			t.Errorf("ObjectId mismatch: want %q got %q", first.ObjectId, got.ObjectId)
		}
	})

	t.Run("GetByObjectId", func(t *testing.T) {
		got, err := client.GetUser(env.Host, env.Port, env.User, env.Pass, first.ObjectId)
		if err != nil {
			t.Fatalf("GetUser by ObjectId: %v", err)
		}
		if got.Alias != first.Alias {
			t.Errorf("Alias mismatch: want %q got %q", first.Alias, got.Alias)
		}
	})
}

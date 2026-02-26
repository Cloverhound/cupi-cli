package tests

import (
	"testing"

	"github.com/Cloverhound/cupi-cli/internal/client"
)

// TestCOSListAndGet verifies that listing and fetching COS objects works.
// COS cannot be created or deleted via CUPI — only listed, fetched, and updated.
func TestCOSListAndGet(t *testing.T) {
	env := getTestEnv(t)
	verifyConnectivity(t, env)

	t.Run("List", func(t *testing.T) {
		coses, err := client.ListCOS(env.Host, env.Port, env.User, env.Pass, "", 10)
		if err != nil {
			t.Fatalf("ListCOS: %v", err)
		}
		if len(coses) == 0 {
			t.Fatal("expected at least one COS, got none")
		}
		t.Logf("ListCOS returned %d entries (max 10)", len(coses))
		for _, c := range coses {
			t.Logf("  COS: %s (%s)", c.DisplayName, c.ObjectId)
		}
	})

	t.Run("GetByName", func(t *testing.T) {
		coses, err := client.ListCOS(env.Host, env.Port, env.User, env.Pass, "", 1)
		if err != nil || len(coses) == 0 {
			t.Skip("no COS available")
		}
		first := coses[0]

		got, err := client.GetCOS(env.Host, env.Port, env.User, env.Pass, first.DisplayName)
		if err != nil {
			t.Fatalf("GetCOS by name: %v", err)
		}
		if got.ObjectId != first.ObjectId {
			t.Errorf("ObjectId mismatch: want %q got %q", first.ObjectId, got.ObjectId)
		}
	})

	t.Run("GetByObjectId", func(t *testing.T) {
		coses, err := client.ListCOS(env.Host, env.Port, env.User, env.Pass, "", 1)
		if err != nil || len(coses) == 0 {
			t.Skip("no COS available")
		}
		first := coses[0]

		got, err := client.GetCOS(env.Host, env.Port, env.User, env.Pass, first.ObjectId)
		if err != nil {
			t.Fatalf("GetCOS by ObjectId: %v", err)
		}
		if got.DisplayName != first.DisplayName {
			t.Errorf("DisplayName mismatch: want %q got %q", first.DisplayName, got.DisplayName)
		}
	})
}

// TestCOSUpdate updates the display name of the first COS then restores it.
func TestCOSUpdate(t *testing.T) {
	env := getTestEnv(t)
	verifyConnectivity(t, env)

	coses, err := client.ListCOS(env.Host, env.Port, env.User, env.Pass, "", 1)
	if err != nil || len(coses) == 0 {
		t.Skip("no COS available to update")
	}
	target := coses[0]
	originalName := target.DisplayName
	testName := originalName + " [cupi-test]"

	t.Logf("updating COS %q (ObjectId: %s)", originalName, target.ObjectId)

	// --- Modify ---
	t.Run("Modify", func(t *testing.T) {
		updates := map[string]interface{}{
			"DisplayName": testName,
		}
		if err := client.UpdateCOS(env.Host, env.Port, env.User, env.Pass, target.ObjectId, updates); err != nil {
			t.Fatalf("UpdateCOS (set test name): %v", err)
		}

		got, err := client.GetCOS(env.Host, env.Port, env.User, env.Pass, target.ObjectId)
		if err != nil {
			t.Fatalf("GetCOS after update: %v", err)
		}
		if got.DisplayName != testName {
			t.Errorf("DisplayName: want %q got %q", testName, got.DisplayName)
		}
	})

	// --- Restore ---
	t.Run("Restore", func(t *testing.T) {
		restore := map[string]interface{}{
			"DisplayName": originalName,
		}
		if err := client.UpdateCOS(env.Host, env.Port, env.User, env.Pass, target.ObjectId, restore); err != nil {
			t.Fatalf("UpdateCOS (restore original name): %v", err)
		}

		got, err := client.GetCOS(env.Host, env.Port, env.User, env.Pass, target.ObjectId)
		if err != nil {
			t.Fatalf("GetCOS after restore: %v", err)
		}
		if got.DisplayName != originalName {
			t.Errorf("DisplayName after restore: want %q got %q", originalName, got.DisplayName)
		}
		t.Logf("COS display name restored to %q", originalName)
	})
}

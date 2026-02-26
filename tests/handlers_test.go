package tests

import (
	"testing"

	"github.com/Cloverhound/cupi-cli/internal/client"
)

// TestHandlersAddModifyDelete creates a call handler, updates its DTMF access ID,
// then deletes it. Requires a call handler template to exist on the server.
func TestHandlersAddModifyDelete(t *testing.T) {
	env := getTestEnv(t)
	verifyConnectivity(t, env)

	templateID := getFirstCallHandlerTemplateID(t, env)
	t.Logf("using call handler template ObjectId: %s", templateID)

	displayName := "CUPI Test Handler " + uniqueSuffix()
	t.Logf("test handler display name: %s", displayName)

	// --- Add ---
	t.Run("Add", func(t *testing.T) {
		fields := map[string]interface{}{
			"DisplayName": displayName,
		}
		h, err := client.CreateCallHandler(env.Host, env.Port, env.User, env.Pass, templateID, fields)
		if err != nil {
			t.Fatalf("CreateCallHandler: %v", err)
		}
		t.Logf("created handler ObjectId: %s", h.ObjectId)

		got, err := client.GetCallHandler(env.Host, env.Port, env.User, env.Pass, displayName)
		if err != nil {
			t.Fatalf("GetCallHandler after create: %v", err)
		}
		if got.DisplayName != displayName {
			t.Errorf("DisplayName mismatch: want %q got %q", displayName, got.DisplayName)
		}
	})

	// --- Modify ---
	t.Run("Modify", func(t *testing.T) {
		updates := map[string]interface{}{
			"DtmfAccessId": "79998",
		}
		if err := client.UpdateCallHandler(env.Host, env.Port, env.User, env.Pass, displayName, updates); err != nil {
			t.Fatalf("UpdateCallHandler: %v", err)
		}

		got, err := client.GetCallHandler(env.Host, env.Port, env.User, env.Pass, displayName)
		if err != nil {
			t.Fatalf("GetCallHandler after update: %v", err)
		}
		if got.DtmfAccessId != "79998" {
			t.Errorf("DtmfAccessId: want %q got %q", "79998", got.DtmfAccessId)
		}
	})

	// --- Delete ---
	t.Run("Delete", func(t *testing.T) {
		if err := client.DeleteCallHandler(env.Host, env.Port, env.User, env.Pass, displayName); err != nil {
			t.Fatalf("DeleteCallHandler: %v", err)
		}

		_, err := client.GetCallHandler(env.Host, env.Port, env.User, env.Pass, displayName)
		if err == nil {
			t.Errorf("expected error after deletion but GetCallHandler succeeded")
		}
		t.Logf("confirmed deleted: %v", err)
	})
}

// TestHandlersList verifies listing returns at least one handler (Opening Greeting always exists).
func TestHandlersList(t *testing.T) {
	env := getTestEnv(t)
	verifyConnectivity(t, env)

	handlers, err := client.ListCallHandlers(env.Host, env.Port, env.User, env.Pass, "", 10)
	if err != nil {
		t.Fatalf("ListCallHandlers: %v", err)
	}
	if len(handlers) == 0 {
		t.Fatal("expected at least one call handler (Opening Greeting), got none")
	}
	t.Logf("ListCallHandlers returned %d handlers (max 10)", len(handlers))
}

// TestHandlersGetRoundtrip verifies get-by-name and get-by-ObjectId return consistent data.
func TestHandlersGetRoundtrip(t *testing.T) {
	env := getTestEnv(t)
	verifyConnectivity(t, env)

	handlers, err := client.ListCallHandlers(env.Host, env.Port, env.User, env.Pass, "", 1)
	if err != nil || len(handlers) == 0 {
		t.Skip("no call handlers available for round-trip test")
	}
	first := handlers[0]

	t.Run("GetByName", func(t *testing.T) {
		got, err := client.GetCallHandler(env.Host, env.Port, env.User, env.Pass, first.DisplayName)
		if err != nil {
			t.Fatalf("GetCallHandler by name: %v", err)
		}
		if got.ObjectId != first.ObjectId {
			t.Errorf("ObjectId mismatch: want %q got %q", first.ObjectId, got.ObjectId)
		}
	})

	t.Run("GetByObjectId", func(t *testing.T) {
		got, err := client.GetCallHandler(env.Host, env.Port, env.User, env.Pass, first.ObjectId)
		if err != nil {
			t.Fatalf("GetCallHandler by ObjectId: %v", err)
		}
		if got.DisplayName != first.DisplayName {
			t.Errorf("DisplayName mismatch: want %q got %q", first.DisplayName, got.DisplayName)
		}
	})
}

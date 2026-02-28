package tests

import (
	"strings"
	"testing"

	"github.com/Cloverhound/cupi-cli/internal/client"
)

// TestDistListsAddModifyDelete creates a distribution list, updates it,
// adds a user member, removes that member, then deletes the list.
func TestDistListsAddModifyDelete(t *testing.T) {
	env := getTestEnv(t)
	verifyConnectivity(t, env)

	alias := "cupi-tl-" + uniqueSuffix()
	t.Logf("test distlist alias: %s", alias)

	memberObjectId := getDistListMemberObjectID(t, env)
	t.Logf("member ObjectId (first regular user): %s", memberObjectId)
	memberAdded := false

	// --- Add ---
	t.Run("Add", func(t *testing.T) {
		fields := map[string]interface{}{
			"Alias":       alias,
			"DisplayName": "CUPI Test List",
		}
		dl, err := client.CreateDistList(env.Host, env.Port, env.User, env.Pass, fields)
		if err != nil {
			t.Fatalf("CreateDistList: %v", err)
		}
		t.Logf("created distlist ObjectId: %s", dl.ObjectId)

		got, err := client.GetDistList(env.Host, env.Port, env.User, env.Pass, alias)
		if err != nil {
			t.Fatalf("GetDistList after create: %v", err)
		}
		if got.Alias != alias {
			t.Errorf("alias mismatch: want %q got %q", alias, got.Alias)
		}
	})

	// --- Modify ---
	t.Run("Modify", func(t *testing.T) {
		updates := map[string]interface{}{
			"DisplayName": "CUPI Test List Updated",
		}
		if err := client.UpdateDistList(env.Host, env.Port, env.User, env.Pass, alias, updates); err != nil {
			t.Fatalf("UpdateDistList: %v", err)
		}

		got, err := client.GetDistList(env.Host, env.Port, env.User, env.Pass, alias)
		if err != nil {
			t.Fatalf("GetDistList after update: %v", err)
		}
		if got.DisplayName != "CUPI Test List Updated" {
			t.Errorf("DisplayName: want %q got %q", "CUPI Test List Updated", got.DisplayName)
		}
	})

	// --- Members: Add ---
	t.Run("MembersAdd", func(t *testing.T) {
		if err := client.AddDistListMember(env.Host, env.Port, env.User, env.Pass, alias, memberObjectId); err != nil {
			// CUC may reject member adds with a DB-level FK constraint during license
			// revalidation (ccnullfkfilter). Skip rather than fail in that case.
			if strings.Contains(err.Error(), "ccnullfkfilter") {
				t.Skipf("server rejected member add (CUC license revalidation in progress): %v", err)
			}
			t.Fatalf("AddDistListMember: %v", err)
		}

		members, err := client.ListDistListMembers(env.Host, env.Port, env.User, env.Pass, alias)
		if err != nil {
			t.Fatalf("ListDistListMembers after add: %v", err)
		}
		found := false
		for _, m := range members {
			if m.MemberObjectId == memberObjectId {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("member %s not found in list after add", memberObjectId)
		}
		memberAdded = true
		t.Logf("member count: %d", len(members))
	})

	// --- Members: Remove ---
	t.Run("MembersRemove", func(t *testing.T) {
		if !memberAdded {
			t.Skip("skipping remove: member was not successfully added")
		}
		if err := client.RemoveDistListMember(env.Host, env.Port, env.User, env.Pass, alias, memberObjectId); err != nil {
			t.Fatalf("RemoveDistListMember: %v", err)
		}

		members, err := client.ListDistListMembers(env.Host, env.Port, env.User, env.Pass, alias)
		if err != nil {
			t.Fatalf("ListDistListMembers after remove: %v", err)
		}
		for _, m := range members {
			if m.MemberObjectId == memberObjectId {
				t.Errorf("member %s still present after remove", memberObjectId)
			}
		}
	})

	// --- Delete ---
	t.Run("Delete", func(t *testing.T) {
		if err := client.DeleteDistList(env.Host, env.Port, env.User, env.Pass, alias); err != nil {
			t.Fatalf("DeleteDistList: %v", err)
		}

		_, err := client.GetDistList(env.Host, env.Port, env.User, env.Pass, alias)
		if err == nil {
			t.Errorf("expected error after deletion but GetDistList succeeded")
		}
		t.Logf("confirmed deleted: %v", err)
	})
}

// TestDistListsList verifies listing returns results.
func TestDistListsList(t *testing.T) {
	env := getTestEnv(t)
	verifyConnectivity(t, env)

	lists, err := client.ListDistLists(env.Host, env.Port, env.User, env.Pass, "", 10)
	if err != nil {
		t.Fatalf("ListDistLists: %v", err)
	}
	t.Logf("ListDistLists returned %d lists (max 10)", len(lists))
}

// TestDistListsGetRoundtrip verifies get-by-alias and get-by-ObjectId return consistent data.
func TestDistListsGetRoundtrip(t *testing.T) {
	env := getTestEnv(t)
	verifyConnectivity(t, env)

	lists, err := client.ListDistLists(env.Host, env.Port, env.User, env.Pass, "", 1)
	if err != nil || len(lists) == 0 {
		t.Skip("no distribution lists available for round-trip test")
	}
	first := lists[0]

	t.Run("GetByAlias", func(t *testing.T) {
		got, err := client.GetDistList(env.Host, env.Port, env.User, env.Pass, first.Alias)
		if err != nil {
			t.Fatalf("GetDistList by alias: %v", err)
		}
		if got.ObjectId != first.ObjectId {
			t.Errorf("ObjectId mismatch: want %q got %q", first.ObjectId, got.ObjectId)
		}
	})

	t.Run("GetByObjectId", func(t *testing.T) {
		got, err := client.GetDistList(env.Host, env.Port, env.User, env.Pass, first.ObjectId)
		if err != nil {
			t.Fatalf("GetDistList by ObjectId: %v", err)
		}
		if got.Alias != first.Alias {
			t.Errorf("Alias mismatch: want %q got %q", first.Alias, got.Alias)
		}
	})
}

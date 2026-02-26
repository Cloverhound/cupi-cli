package tests

import (
	"testing"

	"github.com/Cloverhound/cupi-cli/internal/client"
)

// TestSystemInfo verifies the system info endpoint returns data.
func TestSystemInfo(t *testing.T) {
	env := getTestEnv(t)
	verifyConnectivity(t, env)

	info, err := client.GetSystemInfo(env.Host, env.Port, env.User, env.Pass)
	if err != nil {
		t.Fatalf("GetSystemInfo: %v", err)
	}
	if info.Hostname == "" && info.Version == "" && info.DisplayName == "" {
		t.Error("GetSystemInfo returned empty fields")
	}
	t.Logf("system: hostname=%q version=%q displayName=%q",
		info.Hostname, info.Version, info.DisplayName)
}

// TestTemplatesList verifies user templates can be listed.
func TestTemplatesList(t *testing.T) {
	env := getTestEnv(t)
	verifyConnectivity(t, env)

	tmpls, err := client.ListUserTemplates(env.Host, env.Port, env.User, env.Pass, "", 10)
	if err != nil {
		t.Fatalf("ListUserTemplates: %v", err)
	}
	if len(tmpls) == 0 {
		t.Fatal("expected at least one user template, got none")
	}
	t.Logf("ListUserTemplates returned %d templates (max 10)", len(tmpls))
	for _, tmpl := range tmpls {
		t.Logf("  template: %s (%s)", tmpl.DisplayName, tmpl.ObjectId)
	}
}

// TestSchedulesList verifies schedules can be listed.
func TestSchedulesList(t *testing.T) {
	env := getTestEnv(t)
	verifyConnectivity(t, env)

	scheds, err := client.ListSchedules(env.Host, env.Port, env.User, env.Pass, "", 10)
	if err != nil {
		t.Fatalf("ListSchedules: %v", err)
	}
	t.Logf("ListSchedules returned %d schedules (max 10)", len(scheds))
}

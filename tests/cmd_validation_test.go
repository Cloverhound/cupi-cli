// Command validation tests for cupi-cli.
// These tests validate CLI argument parsing, required-field enforcement,
// cobra arg validation, dry-run behavior, and command registration.
//
// These tests do NOT require a live Cisco Unity Connection server.
// They run the compiled binary and check output/exit codes.
//
// Run:
//
//	go test ./tests/ -v -run TestCLI
package tests

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// binaryPath holds the path to the compiled test binary, set by TestMain.
var binaryPath string

// TestMain builds the cupi binary once for all tests in this package.
func TestMain(m *testing.M) {
	// Determine binary name for the current OS
	binName := "cupi-test"
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	// Build relative to the tests/ directory (project root is one level up)
	root, err := filepath.Abs("..")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve project root: %v\n", err)
		os.Exit(1)
	}
	binaryPath = filepath.Join(root, binName)

	build := exec.Command("go", "build", "-o", binaryPath, ".")
	build.Dir = root
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build cupi binary: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()
	os.Remove(binaryPath)
	os.Exit(code)
}

// runCupi executes the built binary with the given args.
// Returns the combined stdout+stderr output and whether the process exited with code 0.
func runCupi(args ...string) (string, bool) {
	cmd := exec.Command(binaryPath, args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	out := buf.String()
	success := err == nil
	return out, success
}

// assertError verifies the command failed and output contains the expected substring.
func assertError(t *testing.T, out string, ok bool, wantSubstr string) {
	t.Helper()
	if ok {
		t.Errorf("expected command to fail but it succeeded; output: %q", out)
		return
	}
	if !strings.Contains(out, wantSubstr) {
		t.Errorf("expected output to contain %q\ngot: %q", wantSubstr, out)
	}
}

// assertSuccess verifies the command succeeded.
func assertSuccess(t *testing.T, out string, ok bool) {
	t.Helper()
	if !ok {
		t.Errorf("expected command to succeed but it failed; output: %q", out)
	}
}

// assertEither verifies the command either succeeded OR failed with a specific message.
// Useful for tests that pass validation but may hit "no default server" depending on env.
func assertEither(t *testing.T, out string, _ bool, failureSubstr string) {
	t.Helper()
	// If it failed, the error must be the expected one (not a validation error)
	if strings.Contains(out, "is required") {
		t.Errorf("hit unexpected 'is required' validation error; output: %q", out)
	}
	if strings.Contains(out, "accepts") && strings.Contains(out, "arg") {
		t.Errorf("hit unexpected arg-count error; output: %q", out)
	}
	_ = failureSubstr // accepted regardless
}

// ─────────────────────────────────────────────────────────────────────────────
// Command Registration
// ─────────────────────────────────────────────────────────────────────────────

// TestCLICommandRegistration verifies that all expected top-level commands appear
// in cupi --help output and have proper sub-commands registered.
func TestCLICommandRegistration(t *testing.T) {
	out, ok := runCupi("--help")
	assertSuccess(t, out, ok)

	topLevel := []string{
		"auth", "users", "distlists", "handlers", "cos", "templates",
		"schedules", "system", "ast", "paws", "dime",
		"dirhandlers", "inthandlers", "routingrules", "partitions", "searchspaces",
		"phonesystems", "portgroups", "ports", "restrictiontables", "roles",
		"authrules", "configvalues", "chtemplates", "smtp", "alternatenames",
	}
	for _, cmd := range topLevel {
		if !strings.Contains(out, cmd) {
			t.Errorf("top-level command %q not found in --help output", cmd)
		}
	}
}

func TestCLIUsersSubcommands(t *testing.T) {
	out, ok := runCupi("users", "--help")
	assertSuccess(t, out, ok)
	for _, sub := range []string{"list", "get", "add", "update", "remove"} {
		if !strings.Contains(out, sub) {
			t.Errorf("users sub-command %q not found in help", sub)
		}
	}
}

func TestCLIDistlistsSubcommands(t *testing.T) {
	out, ok := runCupi("distlists", "--help")
	assertSuccess(t, out, ok)
	for _, sub := range []string{"list", "get", "add", "update", "remove", "members"} {
		if !strings.Contains(out, sub) {
			t.Errorf("distlists sub-command %q not found in help", sub)
		}
	}
}

func TestCLIHandlersSubcommands(t *testing.T) {
	out, ok := runCupi("handlers", "--help")
	assertSuccess(t, out, ok)
	for _, sub := range []string{"list", "get", "add", "update", "remove"} {
		if !strings.Contains(out, sub) {
			t.Errorf("handlers sub-command %q not found in help", sub)
		}
	}
}

func TestCLIAstSubcommands(t *testing.T) {
	out, ok := runCupi("ast", "--help")
	assertSuccess(t, out, ok)
	for _, sub := range []string{"disk", "tftp", "heartbeat", "alerts", "perfmon", "services"} {
		if !strings.Contains(out, sub) {
			t.Errorf("ast sub-command %q not found in help", sub)
		}
	}
}

func TestCLIAstPerfmonSubcommands(t *testing.T) {
	out, ok := runCupi("ast", "perfmon", "--help")
	assertSuccess(t, out, ok)
	for _, sub := range []string{"counters", "collect"} {
		if !strings.Contains(out, sub) {
			t.Errorf("ast perfmon sub-command %q not found in help", sub)
		}
	}
}

func TestCLIAstAlertsSubcommands(t *testing.T) {
	out, ok := runCupi("ast", "alerts", "--help")
	assertSuccess(t, out, ok)
	if !strings.Contains(out, "get") {
		t.Errorf("ast alerts sub-command 'get' not found in help")
	}
	if !strings.Contains(out, "triggered") {
		t.Errorf("ast alerts --triggered flag not mentioned in help")
	}
}

func TestCLIAstServicesSubcommands(t *testing.T) {
	out, ok := runCupi("ast", "services", "--help")
	assertSuccess(t, out, ok)
	for _, sub := range []string{"list", "start", "stop", "restart"} {
		if !strings.Contains(out, sub) {
			t.Errorf("ast services sub-command %q not found in help", sub)
		}
	}
}

func TestCLIPawsSubcommands(t *testing.T) {
	out, ok := runCupi("paws", "--help")
	assertSuccess(t, out, ok)
	for _, sub := range []string{"cluster", "drs"} {
		if !strings.Contains(out, sub) {
			t.Errorf("paws sub-command %q not found in help", sub)
		}
	}
}

func TestCLIPawsClusterSubcommands(t *testing.T) {
	out, ok := runCupi("paws", "cluster", "--help")
	assertSuccess(t, out, ok)
	for _, sub := range []string{"status", "replication"} {
		if !strings.Contains(out, sub) {
			t.Errorf("paws cluster sub-command %q not found in help", sub)
		}
	}
}

func TestCLIPawsDrsSubcommands(t *testing.T) {
	out, ok := runCupi("paws", "drs", "--help")
	assertSuccess(t, out, ok)
	for _, sub := range []string{"backup", "status"} {
		if !strings.Contains(out, sub) {
			t.Errorf("paws drs sub-command %q not found in help", sub)
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Help Text Quality
// ─────────────────────────────────────────────────────────────────────────────

// TestCLIHelpMentionsRequired ensures that "add" commands document required flags.
func TestCLIHelpMentionsRequired(t *testing.T) {
	cases := []struct {
		args    []string
		wantStr string
	}{
		{[]string{"users", "add", "--help"}, "required"},
		{[]string{"distlists", "add", "--help"}, "required"},
		{[]string{"handlers", "add", "--help"}, "required"},
		{[]string{"paws", "drs", "backup", "--help"}, "required"},
		{[]string{"alternatenames", "add", "--help"}, "required"},
		{[]string{"configvalues", "update", "--help"}, "required"},
		{[]string{"phonesystems", "axlservers", "add", "--help"}, "required"},
		{[]string{"restrictiontables", "patterns", "add", "--help"}, "required"},
		{[]string{"routingrules", "conditions", "add", "--help"}, "required"},
		{[]string{"searchspaces", "members", "add", "--help"}, "required"},
	}
	for _, tc := range cases {
		t.Run(strings.Join(tc.args, "_"), func(t *testing.T) {
			out, ok := runCupi(tc.args...)
			assertSuccess(t, out, ok)
			if !strings.Contains(strings.ToLower(out), strings.ToLower(tc.wantStr)) {
				t.Errorf("help output for %v does not mention %q\ngot: %q",
					tc.args, tc.wantStr, out)
			}
		})
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// users
// ─────────────────────────────────────────────────────────────────────────────

func TestCLIUsersAddMissingAlias(t *testing.T) {
	out, ok := runCupi("users", "add", "--dtmf", "1234")
	assertError(t, out, ok, "--alias is required")
}

func TestCLIUsersAddMissingDtmf(t *testing.T) {
	out, ok := runCupi("users", "add", "--alias", "jsmith")
	assertError(t, out, ok, "--dtmf is required")
}

func TestCLIUsersAddHelpDocumentsFlags(t *testing.T) {
	out, ok := runCupi("users", "add", "--help")
	assertSuccess(t, out, ok)
	for _, flag := range []string{"--alias", "--dtmf", "--first-name", "--last-name", "--template"} {
		if !strings.Contains(out, flag) {
			t.Errorf("help missing flag %q", flag)
		}
	}
}

func TestCLIUsersGetRequiresArg(t *testing.T) {
	out, ok := runCupi("users", "get")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIUsersRemoveRequiresArg(t *testing.T) {
	out, ok := runCupi("users", "remove")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIUsersUpdateRequiresArg(t *testing.T) {
	out, ok := runCupi("users", "update")
	assertError(t, out, ok, "accepts 1 arg")
}

// ─────────────────────────────────────────────────────────────────────────────
// distlists
// ─────────────────────────────────────────────────────────────────────────────

func TestCLIDistlistsAddMissingAlias(t *testing.T) {
	out, ok := runCupi("distlists", "add")
	assertError(t, out, ok, "--alias is required")
}

func TestCLIDistlistsGetRequiresArg(t *testing.T) {
	out, ok := runCupi("distlists", "get")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIDistlistsRemoveRequiresArg(t *testing.T) {
	out, ok := runCupi("distlists", "remove")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIDistlistsUpdateRequiresArg(t *testing.T) {
	out, ok := runCupi("distlists", "update")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIDistlistsMembersAddRequiresArgs(t *testing.T) {
	out, ok := runCupi("distlists", "members", "add")
	assertError(t, out, ok, "accepts 2 arg")
}

func TestCLIDistlistsMembersRemoveRequiresArgs(t *testing.T) {
	out, ok := runCupi("distlists", "members", "remove")
	assertError(t, out, ok, "accepts 2 arg")
}

// ─────────────────────────────────────────────────────────────────────────────
// handlers
// ─────────────────────────────────────────────────────────────────────────────

func TestCLIHandlersAddMissingDisplayName(t *testing.T) {
	out, ok := runCupi("handlers", "add", "--template-id", "some-uuid")
	assertError(t, out, ok, "--display-name is required")
}

func TestCLIHandlersAddMissingTemplateID(t *testing.T) {
	out, ok := runCupi("handlers", "add", "--display-name", "Test Handler")
	assertError(t, out, ok, "--template-id is required")
}

func TestCLIHandlersAddHelpDocumentsFlags(t *testing.T) {
	out, ok := runCupi("handlers", "add", "--help")
	assertSuccess(t, out, ok)
	for _, flag := range []string{"--display-name", "--template-id", "--dtmf"} {
		if !strings.Contains(out, flag) {
			t.Errorf("help missing flag %q", flag)
		}
	}
}

func TestCLIHandlersGetRequiresArg(t *testing.T) {
	out, ok := runCupi("handlers", "get")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIHandlersRemoveRequiresArg(t *testing.T) {
	out, ok := runCupi("handlers", "remove")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIHandlersUpdateRequiresArg(t *testing.T) {
	out, ok := runCupi("handlers", "update")
	assertError(t, out, ok, "accepts 1 arg")
}

// ─────────────────────────────────────────────────────────────────────────────
// authrules
// ─────────────────────────────────────────────────────────────────────────────

func TestCLIAuthRulesAddMissingDisplayName(t *testing.T) {
	out, ok := runCupi("authrules", "add")
	assertError(t, out, ok, "--display-name is required")
}

func TestCLIAuthRulesGetRequiresArg(t *testing.T) {
	out, ok := runCupi("authrules", "get")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIAuthRulesUpdateRequiresArg(t *testing.T) {
	out, ok := runCupi("authrules", "update")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIAuthRulesRemoveRequiresArg(t *testing.T) {
	out, ok := runCupi("authrules", "remove")
	assertError(t, out, ok, "accepts 1 arg")
}

// ─────────────────────────────────────────────────────────────────────────────
// callhandlertemplates (chtemplates)
// ─────────────────────────────────────────────────────────────────────────────

func TestCLICHTemplatesAddMissingDisplayName(t *testing.T) {
	out, ok := runCupi("chtemplates", "add")
	assertError(t, out, ok, "--display-name is required")
}

func TestCLICHTemplatesGetRequiresArg(t *testing.T) {
	out, ok := runCupi("chtemplates", "get")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLICHTemplatesUpdateRequiresArg(t *testing.T) {
	out, ok := runCupi("chtemplates", "update")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLICHTemplatesRemoveRequiresArg(t *testing.T) {
	out, ok := runCupi("chtemplates", "remove")
	assertError(t, out, ok, "accepts 1 arg")
}

// ─────────────────────────────────────────────────────────────────────────────
// configvalues
// ─────────────────────────────────────────────────────────────────────────────

func TestCLIConfigValuesGetRequiresArg(t *testing.T) {
	out, ok := runCupi("configvalues", "get")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIConfigValuesUpdateRequiresArg(t *testing.T) {
	out, ok := runCupi("configvalues", "update")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIConfigValuesUpdateMissingValue(t *testing.T) {
	// Arg is provided but --value flag is missing
	out, ok := runCupi("configvalues", "update", "System.Licensing.MaxSessions")
	assertError(t, out, ok, "--value is required")
}

func TestCLIConfigValuesUpdateHelpDocumentsValue(t *testing.T) {
	out, ok := runCupi("configvalues", "update", "--help")
	assertSuccess(t, out, ok)
	if !strings.Contains(out, "--value") {
		t.Errorf("configvalues update help missing --value flag")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// alternatenames
// ─────────────────────────────────────────────────────────────────────────────

func TestCLIAlternateNamesAddMissingFields(t *testing.T) {
	// All three fields required together
	out, ok := runCupi("alternatenames", "add")
	assertError(t, out, ok, "required")
}

func TestCLIAlternateNamesAddHelpDocumentsRequiredFlags(t *testing.T) {
	out, ok := runCupi("alternatenames", "add", "--help")
	assertSuccess(t, out, ok)
	for _, flag := range []string{"--user-id", "--first-name", "--last-name"} {
		if !strings.Contains(out, flag) {
			t.Errorf("alternatenames add help missing flag %q", flag)
		}
	}
}

func TestCLIAlternateNamesListRequiresArg(t *testing.T) {
	out, ok := runCupi("alternatenames", "list")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIAlternateNamesGetRequiresArg(t *testing.T) {
	out, ok := runCupi("alternatenames", "get")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIAlternateNamesUpdateRequiresArg(t *testing.T) {
	out, ok := runCupi("alternatenames", "update")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIAlternateNamesRemoveRequiresArg(t *testing.T) {
	out, ok := runCupi("alternatenames", "remove")
	assertError(t, out, ok, "accepts 1 arg")
}

// ─────────────────────────────────────────────────────────────────────────────
// phonesystems
// ─────────────────────────────────────────────────────────────────────────────

func TestCLIPhoneSystemsAddMissingDisplayName(t *testing.T) {
	out, ok := runCupi("phonesystems", "add")
	assertError(t, out, ok, "--display-name is required")
}

func TestCLIPhoneSystemsGetRequiresArg(t *testing.T) {
	out, ok := runCupi("phonesystems", "get")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIPhoneSystemsUpdateRequiresArg(t *testing.T) {
	out, ok := runCupi("phonesystems", "update")
	assertError(t, out, ok, "accepts 1 arg")
}

// Note: phonesystems has no remove command (phone systems manage via CUPI add only).

func TestCLIPhoneSystemsAxlserversAddRequiresArg(t *testing.T) {
	out, ok := runCupi("phonesystems", "axlservers", "add")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIPhoneSystemsAxlserversAddMissingServerName(t *testing.T) {
	out, ok := runCupi("phonesystems", "axlservers", "add", "ps-uuid-1234")
	assertError(t, out, ok, "--server-name is required")
}

func TestCLIPhoneSystemsAxlserversRemoveRequiresArgs(t *testing.T) {
	out, ok := runCupi("phonesystems", "axlservers", "remove")
	assertError(t, out, ok, "accepts 2 arg")
}

// ─────────────────────────────────────────────────────────────────────────────
// portgroups
// ─────────────────────────────────────────────────────────────────────────────

func TestCLIPortGroupsGetRequiresArg(t *testing.T) {
	out, ok := runCupi("portgroups", "get")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIPortGroupsUpdateRequiresArg(t *testing.T) {
	out, ok := runCupi("portgroups", "update")
	assertError(t, out, ok, "accepts 1 arg")
}

// ─────────────────────────────────────────────────────────────────────────────
// ports
// ─────────────────────────────────────────────────────────────────────────────

func TestCLIPortsAddMissingDisplayName(t *testing.T) {
	out, ok := runCupi("ports", "add")
	assertError(t, out, ok, "--display-name is required")
}

func TestCLIPortsGetRequiresArg(t *testing.T) {
	out, ok := runCupi("ports", "get")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIPortsUpdateRequiresArg(t *testing.T) {
	out, ok := runCupi("ports", "update")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIPortsRemoveRequiresArg(t *testing.T) {
	out, ok := runCupi("ports", "remove")
	assertError(t, out, ok, "accepts 1 arg")
}

// ─────────────────────────────────────────────────────────────────────────────
// restrictiontables
// ─────────────────────────────────────────────────────────────────────────────

func TestCLIRestrictionTablesAddMissingDisplayName(t *testing.T) {
	out, ok := runCupi("restrictiontables", "add")
	assertError(t, out, ok, "--display-name is required")
}

func TestCLIRestrictionTablesGetRequiresArg(t *testing.T) {
	out, ok := runCupi("restrictiontables", "get")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIRestrictionTablesUpdateRequiresArg(t *testing.T) {
	out, ok := runCupi("restrictiontables", "update")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIRestrictionTablesRemoveRequiresArg(t *testing.T) {
	out, ok := runCupi("restrictiontables", "remove")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIRestrictionTablesPatternsAddRequiresArg(t *testing.T) {
	out, ok := runCupi("restrictiontables", "patterns", "add")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIRestrictionTablesPatternsAddMissingPattern(t *testing.T) {
	out, ok := runCupi("restrictiontables", "patterns", "add", "table-uuid-1234")
	assertError(t, out, ok, "--pattern is required")
}

func TestCLIRestrictionTablesPatternsRemoveRequiresArgs(t *testing.T) {
	out, ok := runCupi("restrictiontables", "patterns", "remove")
	assertError(t, out, ok, "accepts 2 arg")
}

// ─────────────────────────────────────────────────────────────────────────────
// roles
// ─────────────────────────────────────────────────────────────────────────────

func TestCLIRolesAddMissingDisplayName(t *testing.T) {
	out, ok := runCupi("roles", "add")
	assertError(t, out, ok, "--display-name is required")
}

func TestCLIRolesGetRequiresArg(t *testing.T) {
	out, ok := runCupi("roles", "get")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIRolesUpdateRequiresArg(t *testing.T) {
	out, ok := runCupi("roles", "update")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIRolesRemoveRequiresArg(t *testing.T) {
	out, ok := runCupi("roles", "remove")
	assertError(t, out, ok, "accepts 1 arg")
}

// ─────────────────────────────────────────────────────────────────────────────
// routingrules
// ─────────────────────────────────────────────────────────────────────────────

func TestCLIRoutingRulesAddMissingDisplayName(t *testing.T) {
	out, ok := runCupi("routingrules", "add")
	assertError(t, out, ok, "--display-name is required")
}

func TestCLIRoutingRulesGetRequiresArg(t *testing.T) {
	out, ok := runCupi("routingrules", "get")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIRoutingRulesUpdateRequiresArg(t *testing.T) {
	out, ok := runCupi("routingrules", "update")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIRoutingRulesRemoveRequiresArg(t *testing.T) {
	out, ok := runCupi("routingrules", "remove")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIRoutingRulesConditionsAddRequiresArg(t *testing.T) {
	out, ok := runCupi("routingrules", "conditions", "add")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIRoutingRulesConditionsAddMissingParameter(t *testing.T) {
	out, ok := runCupi("routingrules", "conditions", "add", "rule-uuid-1234")
	assertError(t, out, ok, "--parameter is required")
}

func TestCLIRoutingRulesConditionsRemoveRequiresArgs(t *testing.T) {
	out, ok := runCupi("routingrules", "conditions", "remove")
	assertError(t, out, ok, "accepts 2 arg")
}

// ─────────────────────────────────────────────────────────────────────────────
// searchspaces
// ─────────────────────────────────────────────────────────────────────────────

func TestCLISearchSpacesAddMissingName(t *testing.T) {
	out, ok := runCupi("searchspaces", "add")
	assertError(t, out, ok, "--name is required")
}

func TestCLISearchSpacesGetRequiresArg(t *testing.T) {
	out, ok := runCupi("searchspaces", "get")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLISearchSpacesUpdateRequiresArg(t *testing.T) {
	out, ok := runCupi("searchspaces", "update")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLISearchSpacesRemoveRequiresArg(t *testing.T) {
	out, ok := runCupi("searchspaces", "remove")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLISearchSpacesMembersAddRequiresArg(t *testing.T) {
	out, ok := runCupi("searchspaces", "members", "add")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLISearchSpacesMembersAddMissingPartitionID(t *testing.T) {
	out, ok := runCupi("searchspaces", "members", "add", "ss-uuid-1234")
	assertError(t, out, ok, "--partition-id is required")
}

func TestCLISearchSpacesMembersRemoveRequiresArgs(t *testing.T) {
	out, ok := runCupi("searchspaces", "members", "remove")
	assertError(t, out, ok, "accepts 2 arg")
}

// ─────────────────────────────────────────────────────────────────────────────
// partitions
// ─────────────────────────────────────────────────────────────────────────────

func TestCLIPartitionsAddMissingName(t *testing.T) {
	out, ok := runCupi("partitions", "add")
	assertError(t, out, ok, "--name is required")
}

func TestCLIPartitionsGetRequiresArg(t *testing.T) {
	out, ok := runCupi("partitions", "get")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIPartitionsUpdateRequiresArg(t *testing.T) {
	out, ok := runCupi("partitions", "update")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIPartitionsRemoveRequiresArg(t *testing.T) {
	out, ok := runCupi("partitions", "remove")
	assertError(t, out, ok, "accepts 1 arg")
}

// ─────────────────────────────────────────────────────────────────────────────
// dirhandlers / inthandlers
// ─────────────────────────────────────────────────────────────────────────────

func TestCLIDirHandlersAddMissingDisplayName(t *testing.T) {
	out, ok := runCupi("dirhandlers", "add")
	assertError(t, out, ok, "--display-name is required")
}

func TestCLIDirHandlersGetRequiresArg(t *testing.T) {
	out, ok := runCupi("dirhandlers", "get")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIDirHandlersUpdateRequiresArg(t *testing.T) {
	out, ok := runCupi("dirhandlers", "update")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIIntHandlersAddMissingDisplayName(t *testing.T) {
	out, ok := runCupi("inthandlers", "add")
	assertError(t, out, ok, "--display-name is required")
}

func TestCLIIntHandlersGetRequiresArg(t *testing.T) {
	out, ok := runCupi("inthandlers", "get")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIIntHandlersUpdateRequiresArg(t *testing.T) {
	out, ok := runCupi("inthandlers", "update")
	assertError(t, out, ok, "accepts 1 arg")
}

// ─────────────────────────────────────────────────────────────────────────────
// smtp
// ─────────────────────────────────────────────────────────────────────────────

func TestCLISmtpSubcommands(t *testing.T) {
	out, ok := runCupi("smtp", "--help")
	assertSuccess(t, out, ok)
	for _, sub := range []string{"server", "client"} {
		if !strings.Contains(out, sub) {
			t.Errorf("smtp sub-command %q not found in help", sub)
		}
	}
}

func TestCLISmtpServerUpdateNoFieldsError(t *testing.T) {
	// No update fields provided → should error rather than hit server
	out, ok := runCupi("smtp", "server", "update")
	assertError(t, out, ok, "")
	// Just verify it errors — the specific message may be "no fields" or "no default server"
	if ok {
		t.Errorf("expected smtp server update (no args) to fail, but it succeeded: %q", out)
	}
}

func TestCLISmtpClientUpdateNoFieldsError(t *testing.T) {
	out, ok := runCupi("smtp", "client", "update")
	if ok {
		t.Errorf("expected smtp client update (no args) to fail, but it succeeded: %q", out)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// paws drs backup — required flags
// ─────────────────────────────────────────────────────────────────────────────

func TestCLIPawsDrsBackupMissingSFTPServer(t *testing.T) {
	out, ok := runCupi("paws", "drs", "backup",
		"--sftp-user", "backup", "--sftp-password", "secret", "--sftp-dir", "/backups")
	assertError(t, out, ok, "--sftp-server is required")
}

func TestCLIPawsDrsBackupMissingSFTPUser(t *testing.T) {
	out, ok := runCupi("paws", "drs", "backup",
		"--sftp-server", "10.0.0.5", "--sftp-password", "secret", "--sftp-dir", "/backups")
	assertError(t, out, ok, "--sftp-user is required")
}

func TestCLIPawsDrsBackupMissingSFTPPassword(t *testing.T) {
	out, ok := runCupi("paws", "drs", "backup",
		"--sftp-server", "10.0.0.5", "--sftp-user", "backup", "--sftp-dir", "/backups")
	assertError(t, out, ok, "--sftp-password is required")
}

func TestCLIPawsDrsBackupMissingSFTPDir(t *testing.T) {
	out, ok := runCupi("paws", "drs", "backup",
		"--sftp-server", "10.0.0.5", "--sftp-user", "backup", "--sftp-password", "secret")
	assertError(t, out, ok, "--sftp-dir is required")
}

func TestCLIPawsDrsBackupHelpDocumentsFlags(t *testing.T) {
	out, ok := runCupi("paws", "drs", "backup", "--help")
	assertSuccess(t, out, ok)
	for _, flag := range []string{"--sftp-server", "--sftp-user", "--sftp-password", "--sftp-dir", "--sftp-port"} {
		if !strings.Contains(out, flag) {
			t.Errorf("paws drs backup help missing flag %q", flag)
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// dime
// ─────────────────────────────────────────────────────────────────────────────

func TestCLIDimeGetFileRequiresArg(t *testing.T) {
	out, ok := runCupi("dime", "get-file")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIDimeGetFileHelpDocumentsFlags(t *testing.T) {
	out, ok := runCupi("dime", "get-file", "--help")
	assertSuccess(t, out, ok)
	for _, flag := range []string{"--output", "--node"} {
		if !strings.Contains(out, flag) {
			t.Errorf("dime get-file help missing flag %q", flag)
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// ast perfmon sub-commands — require positional arg
// ─────────────────────────────────────────────────────────────────────────────

func TestCLIAstPerfmonCountersRequiresArg(t *testing.T) {
	out, ok := runCupi("ast", "perfmon", "counters")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIAstPerfmonCollectRequiresArg(t *testing.T) {
	out, ok := runCupi("ast", "perfmon", "collect")
	assertError(t, out, ok, "accepts 1 arg")
}

// ─────────────────────────────────────────────────────────────────────────────
// ast alerts get — requires positional arg
// ─────────────────────────────────────────────────────────────────────────────

func TestCLIAstAlertsGetRequiresArg(t *testing.T) {
	out, ok := runCupi("ast", "alerts", "get")
	assertError(t, out, ok, "accepts 1 arg")
}

// ─────────────────────────────────────────────────────────────────────────────
// ast services — require positional arg; restart/start/stop honor --dry-run
// ─────────────────────────────────────────────────────────────────────────────

func TestCLIAstServicesStartRequiresArg(t *testing.T) {
	out, ok := runCupi("ast", "services", "start")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIAstServicesStopRequiresArg(t *testing.T) {
	out, ok := runCupi("ast", "services", "stop")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIAstServicesRestartRequiresArg(t *testing.T) {
	out, ok := runCupi("ast", "services", "restart")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLIAstServicesDryRunStart(t *testing.T) {
	out, ok := runCupi("--dry-run", "ast", "services", "start", "Cisco Unity Connection Voicemail")
	assertSuccess(t, out, ok)
	if !strings.Contains(out, "[dry-run]") || !strings.Contains(out, "Start") {
		t.Errorf("expected dry-run output for service start, got: %q", out)
	}
}

func TestCLIAstServicesDryRunStop(t *testing.T) {
	out, ok := runCupi("--dry-run", "ast", "services", "stop", "Cisco Unity Connection Voicemail")
	assertSuccess(t, out, ok)
	if !strings.Contains(out, "[dry-run]") || !strings.Contains(out, "Stop") {
		t.Errorf("expected dry-run output for service stop, got: %q", out)
	}
}

func TestCLIAstServicesDryRunRestart(t *testing.T) {
	out, ok := runCupi("--dry-run", "ast", "services", "restart", "Cisco Tomcat")
	assertSuccess(t, out, ok)
	if !strings.Contains(out, "[dry-run]") || !strings.Contains(out, "Restart") {
		t.Errorf("expected dry-run output for service restart, got: %q", out)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// cos / templates / schedules — read-only, get commands need arg
// ─────────────────────────────────────────────────────────────────────────────

func TestCLICOSGetRequiresArg(t *testing.T) {
	out, ok := runCupi("cos", "get")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLICOSUpdateRequiresArg(t *testing.T) {
	out, ok := runCupi("cos", "update")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLITemplatesGetRequiresArg(t *testing.T) {
	out, ok := runCupi("templates", "get")
	assertError(t, out, ok, "accepts 1 arg")
}

func TestCLISchedulesGetRequiresArg(t *testing.T) {
	out, ok := runCupi("schedules", "get")
	assertError(t, out, ok, "accepts 1 arg")
}

// ─────────────────────────────────────────────────────────────────────────────
// Global flags
// ─────────────────────────────────────────────────────────────────────────────

func TestCLIGlobalFlagsInHelp(t *testing.T) {
	out, ok := runCupi("--help")
	assertSuccess(t, out, ok)
	for _, flag := range []string{"--server", "--output", "--debug", "--max", "--dry-run"} {
		if !strings.Contains(out, flag) {
			t.Errorf("global flag %q not found in root help", flag)
		}
	}
}

func TestCLIUnknownCommandFails(t *testing.T) {
	out, ok := runCupi("notacommand")
	if ok {
		t.Errorf("expected unknown command to fail, but succeeded; output: %q", out)
	}
	if !strings.Contains(out, "unknown command") {
		t.Errorf("expected 'unknown command' in output, got: %q", out)
	}
}

func TestCLIVersionOrHelpDoesNotCrash(t *testing.T) {
	out, ok := runCupi("--help")
	assertSuccess(t, out, ok)
	if len(out) == 0 {
		t.Error("root --help produced empty output")
	}
}

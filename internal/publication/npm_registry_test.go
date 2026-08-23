package publication

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// stubTool writes an executable script named tool onto PATH (via a fresh
// scratch directory prepended for the test) that prints version to stdout.
func stubTool(t *testing.T, dir, tool, version string) {
	t.Helper()
	script := "#!/bin/sh\necho '" + version + "'\n"
	path := filepath.Join(dir, tool)
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
}

// withStubbedTools prepends a scratch PATH directory containing stub npm and
// node binaries reporting the given versions, restoring PATH on cleanup.
func withStubbedTools(t *testing.T, npmVersion, nodeVersion string) {
	t.Helper()
	dir := t.TempDir()
	stubTool(t, dir, "npm", npmVersion)
	stubTool(t, dir, "node", "v"+nodeVersion)
	original := os.Getenv("PATH")
	if err := os.Setenv("PATH", dir+string(os.PathListSeparator)+original); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Setenv("PATH", original) })
}

// TestNPMCLIRegistryRefusesStagedOpsBelowToolFloor asserts the staged-flow
// precondition belongs only to the public npm adapter. An npm/node pair
// below the documented floor must fail StageSubmit and Approve. This must
// happen before either ever shells out to `npm publish` or `npm ...` for
// the real operation.
func TestNPMCLIRegistryRefusesStagedOpsBelowToolFloor(t *testing.T) {
	withStubbedTools(t, "11.14.0", "22.14.0")
	registry := NewNPMCLIRegistry("")
	if _, err := registry.StageSubmit(context.Background(), "redbench", "1.0.0", []byte("x")); err == nil {
		t.Fatal("expected StageSubmit to refuse an npm version below the staged-flow floor")
	} else if !strings.Contains(err.Error(), "npm") {
		t.Fatalf("expected the tool-floor error to mention npm, got: %v", err)
	}
	if err := registry.Approve(context.Background(), "stage-1"); err == nil {
		t.Fatal("expected Approve to refuse an npm version below the staged-flow floor")
	}
}

// TestNPMCLIRegistryAcceptsStagedOpsAtToolFloor confirms a tool pair at or
// above the floor passes the precondition. Only then does it hit the
// construct-only "not implemented" stub, which is expected for this adapter.
func TestNPMCLIRegistryAcceptsStagedOpsAtToolFloor(t *testing.T) {
	withStubbedTools(t, "11.15.0", "22.14.0")
	registry := NewNPMCLIRegistry("")
	_, err := registry.StageSubmit(context.Background(), "redbench", "1.0.0", []byte("x"))
	if err == nil || strings.Contains(err.Error(), "npm 11.15") || strings.Contains(err.Error(), "node 22.14") {
		t.Fatalf("expected the tool-floor check to pass and fall through to the not-implemented stub, got: %v", err)
	}
	if !strings.Contains(err.Error(), "not implemented") {
		t.Fatalf("expected the construct-only not-implemented error, got: %v", err)
	}
}

// TestFixtureRegistryStagedOpsNeverCheckToolVersions is the row-1 precondition
// isolation. The fixture adapter must never require a locally installed npm
// or node version; that requirement belongs only to the public npm adapter.
// A PATH with no npm/node binaries at all must not affect the fixture's
// StageSubmit, since it never shells either tool.
func TestFixtureRegistryStagedOpsNeverCheckToolVersions(t *testing.T) {
	empty := t.TempDir()
	original := os.Getenv("PATH")
	if err := os.Setenv("PATH", empty); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Setenv("PATH", original) })

	// No HTTP server is reachable at this bogus base, so StageSubmit fails
	// on the network call. It must fail there, never on a missing npm/node
	// tool: that proves the fixture adapter carries no such precondition.
	registry := NewFixtureRegistry("http://127.0.0.1:1")
	_, err := registry.StageSubmit(context.Background(), "redbench", "1.0.0", []byte("x"))
	if err == nil {
		t.Fatal("expected a connection error against an unreachable base")
	}
	if strings.Contains(err.Error(), "npm") || strings.Contains(err.Error(), "node") {
		t.Fatalf("fixture adapter must never require npm/node tooling, got: %v", err)
	}
}

// stubNPMView writes an npm stub onto a fresh PATH directory that responds to
// `npm view ...` as directed. exitCode != 0 simulates the real npm CLI's E404
// on a missing version (Integrity's absent case). exitCode == 0 prints
// stdout, simulating `npm view ... dist.integrity --json` output for a
// version that exists.
func stubNPMView(t *testing.T, stdout string, exitCode int) {
	t.Helper()
	dir := t.TempDir()
	script := "#!/bin/sh\n"
	if exitCode != 0 {
		script += "exit " + strconv.Itoa(exitCode) + "\n"
	} else {
		script += "cat <<'EOF'\n" + stdout + "\nEOF\n"
	}
	if err := os.WriteFile(filepath.Join(dir, "npm"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	original := os.Getenv("PATH")
	if err := os.Setenv("PATH", dir+string(os.PathListSeparator)+original); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Setenv("PATH", original) })
}

// TestNPMCLIRegistryIntegrityAbsentVersion is the CLI adapter's absent-version
// case: npm view exiting non-zero (E404) must classify as not-live with no
// error, never a mismatch.
func TestNPMCLIRegistryIntegrityAbsentVersion(t *testing.T) {
	stubNPMView(t, "", 1)
	registry := NewNPMCLIRegistry("")
	integrity, live, err := registry.Integrity(context.Background(), "redbench", "1.0.0")
	if err != nil {
		t.Fatalf("expected no error for an absent version, got: %v", err)
	}
	if live {
		t.Fatal("expected an absent version to classify as not live")
	}
	if integrity != "" {
		t.Fatalf("expected empty integrity for an absent version, got: %q", integrity)
	}
}

// TestNPMCLIRegistryIntegrityPresentVersion is the CLI adapter's ordinary
// live case: a version with a real integrity value must classify as live
// with that value.
func TestNPMCLIRegistryIntegrityPresentVersion(t *testing.T) {
	stubNPMView(t, `"sha512-abc123"`, 0)
	registry := NewNPMCLIRegistry("")
	integrity, live, err := registry.Integrity(context.Background(), "redbench", "1.0.0")
	if err != nil {
		t.Fatalf("expected no error for a present version, got: %v", err)
	}
	if !live {
		t.Fatal("expected a present version with integrity to classify as live")
	}
	if integrity != "sha512-abc123" {
		t.Fatalf("integrity = %q, want sha512-abc123", integrity)
	}
}

// TestNPMCLIRegistryIntegrityPresentEmptyIsError is the malformed-integrity
// edge: the version exists (npm view succeeded) but reports an empty or
// whitespace integrity — a malformed or hostile registry response. Integrity
// must fail closed with an attributed error here, never misread this as
// live=true-with-empty-integrity or as live=false-not-yet-published.
func TestNPMCLIRegistryIntegrityPresentEmptyIsError(t *testing.T) {
	for name, stdout := range map[string]string{
		"empty string":      `""`,
		"whitespace string": `"   "`,
	} {
		t.Run(name, func(t *testing.T) {
			stubNPMView(t, stdout, 0)
			registry := NewNPMCLIRegistry("")
			integrity, live, err := registry.Integrity(context.Background(), "redbench", "1.0.0")
			if err == nil {
				t.Fatal("expected an error for a present version with empty integrity")
			}
			if !strings.Contains(err.Error(), "no integrity value") {
				t.Fatalf("error did not attribute the empty-integrity response: %v", err)
			}
			if live {
				t.Fatal("expected a present-but-empty-integrity version not to classify as live")
			}
			if integrity != "" {
				t.Fatalf("expected empty integrity on error, got: %q", integrity)
			}
		})
	}
}

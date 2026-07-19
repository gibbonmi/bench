package publication

import (
	"context"
	"os"
	"path/filepath"
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
// precondition belongs only to the public npm adapter: an npm/node pair below
// the documented floor must fail StageSubmit and Approve before ever shelling
// out to `npm publish` or `npm ...` for the real operation.
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
// above the floor passes the precondition (and only then hits the
// construct-only "not implemented" stub, which is expected for this adapter).
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
// isolation: the fixture adapter must never require a locally installed npm
// or node version — that requirement belongs only to the public npm adapter.
// A PATH with no npm/node binaries at all must not affect the fixture's
// StageSubmit, since it never shells either tool.
func TestFixtureRegistryStagedOpsNeverCheckToolVersions(t *testing.T) {
	empty := t.TempDir()
	original := os.Getenv("PATH")
	if err := os.Setenv("PATH", empty); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Setenv("PATH", original) })

	// No HTTP server is reachable at this bogus base, so StageSubmit will
	// fail on the network call — but it must fail there, never on a missing
	// npm/node tool, proving the fixture adapter carries no such precondition.
	registry := NewFixtureRegistry("http://127.0.0.1:1")
	_, err := registry.StageSubmit(context.Background(), "redbench", "1.0.0", []byte("x"))
	if err == nil {
		t.Fatal("expected a connection error against an unreachable base")
	}
	if strings.Contains(err.Error(), "npm") || strings.Contains(err.Error(), "node") {
		t.Fatalf("fixture adapter must never require npm/node tooling, got: %v", err)
	}
}

package spec

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
)

// TestRetirePreflightsEveryPlannedRemovalBeforeDeleting pins the retirement plan's
// atomicity: a later blocked target refuses before the earlier review pickup moves.
func TestRetirePreflightsEveryPlannedRemovalBeforeDeleting(t *testing.T) {
	root := retireRepo(t, "s", "Status: implemented\nRoadmap: FT7\n", nil)
	folder := filepath.Join(root, "specs", "s")
	pickup := filepath.Join(root, "reviews", "s.md")
	if err := os.MkdirAll(filepath.Dir(pickup), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(pickup, []byte("pickup\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ticket := filepath.Join(folder, "tickets", "01.md")
	if err := os.MkdirAll(filepath.Dir(ticket), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(ticket, []byte("ticket\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(folder, 0o755) })
	if err := os.Chmod(folder, 0o500); err != nil {
		capability.Capability(t, capability.Privilege, fmt.Sprintf("cannot restrict spec folder: %v", err))
	}
	if probe, err := os.CreateTemp(folder, "retire-preflight-probe-"); err == nil {
		_ = probe.Close()
		_ = os.Remove(probe.Name())
		capability.Capability(t, capability.Privilege, "read-only spec folder remains writable by this user")
	}

	out, code := runRetire(t, root, "s")
	if code != 1 {
		t.Fatalf("code = %d, want 1; out = %q", code, out)
	}
	for _, path := range []string{"specs/s/spec.md", "specs/s/tickets"} {
		if !strings.Contains(out, path) {
			t.Errorf("refusal = %q, want blocked path %q", out, path)
		}
	}
	if strings.Index(out, "specs/s/spec.md") > strings.Index(out, "specs/s/tickets") {
		t.Errorf("refusal paths are not deterministic: %q", out)
	}
	if _, err := os.Stat(pickup); err != nil {
		t.Errorf("earlier pickup was removed before refusal: %v", err)
	}
	for _, path := range []string{filepath.Join(folder, "spec.md"), filepath.Join(folder, "tickets")} {
		if _, err := os.Stat(path); err != nil {
			t.Errorf("planned target was removed before refusal: %s: %v", path, err)
		}
	}
}

// TestRetirePreflightsNestedRecursiveRemovalBeforeDeleting keeps RemoveAll from
// discovering a blocked descendant only after it has retired the review pickup.
func TestRetirePreflightsNestedRecursiveRemovalBeforeDeleting(t *testing.T) {
	root := retireRepo(t, "s", "Status: implemented\nRoadmap: FT7\n", nil)
	folder := filepath.Join(root, "specs", "s")
	pickup := filepath.Join(root, "reviews", "s.md")
	if err := os.MkdirAll(filepath.Dir(pickup), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(pickup, []byte("pickup\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(folder, "tickets", "nested")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nested, "01.md"), []byte("ticket\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(nested, 0o755) })
	if err := os.Chmod(nested, 0o500); err != nil {
		capability.Capability(t, capability.Privilege, fmt.Sprintf("cannot restrict nested ticket directory: %v", err))
	}
	if probe, err := os.CreateTemp(nested, "retire-preflight-probe-"); err == nil {
		_ = probe.Close()
		_ = os.Remove(probe.Name())
		capability.Capability(t, capability.Privilege, "read-only nested ticket directory remains writable by this user")
	}

	out, code := runRetire(t, root, "s")
	if code != 1 {
		t.Fatalf("code = %d, want 1; out = %q", code, out)
	}
	if !strings.Contains(out, "specs/s/tickets") {
		t.Errorf("refusal = %q, want recursive blocked path", out)
	}
	for _, path := range []string{pickup, filepath.Join(folder, "spec.md"), filepath.Join(folder, "tickets"), folder} {
		if _, err := os.Stat(path); err != nil {
			t.Errorf("planned target was removed before refusal: %s: %v", path, err)
		}
	}
}

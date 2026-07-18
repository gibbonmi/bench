package conformance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// This contract exercises diagnostic aggregation separately from workflow inspection.
func TestReleasePreflightDiagnosticsDistinguishAndAggregate(t *testing.T) {
	missing := t.TempDir()
	if got := strings.Join(checkReleasePreflight(missing), "\n"); !strings.Contains(got, "release preflight registry is absent") || !strings.Contains(got, "release requirement registry is absent") {
		t.Fatalf("absent diagnostics = %q", got)
	}
	malformed := t.TempDir()
	if err := os.MkdirAll(filepath.Join(malformed, "internal", "releaseevidence"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(malformed, "internal", "releaseevidence", "registry.json"), []byte("{\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(malformed, "internal", "releaseevidence", "requirements.json"), []byte("{\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got := strings.Join(checkReleasePreflight(malformed), "\n")
	if !strings.Contains(got, "release preflight registry is malformed") || !strings.Contains(got, "release requirement registry is malformed") {
		t.Fatalf("malformed diagnostics = %q", got)
	}
}

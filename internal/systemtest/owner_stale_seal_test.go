//go:build system

package systemtest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/freshness"
)

// TestSystemOwnerRefusesAStaleSelectedSeal is BF8. The owner's refusal fires inside
// TestMain, so this row runs the suite binary again as a child under a selected
// executable whose seal records a source digest the kit never produced. The child reads
// its verdict from the exit code and the setup line, which no in-process assertion can
// reach.
func TestSystemOwnerRefusesAStaleSelectedSeal(t *testing.T) {
	stale := staleSealedCopy(t, owner.selected.path)
	overrides := []string{"BENCH_RUN_BINARY=" + stale, "BENCH_KIT=" + owner.kit}

	result := owner.runAt(owner.root, overrides, os.Args[0], "-test.run=^$")
	if result.code == 0 {
		t.Fatalf("suite setup under a stale seal exited 0, want a refusal: %q", result.stderr)
	}
	if !strings.Contains(result.stderr, "system owner:") {
		t.Fatalf("suite setup stderr = %q, want the owner's setup refusal", result.stderr)
	}
	if !strings.Contains(result.stderr, freshness.RebuildAction(owner.kit)) {
		t.Fatalf("suite setup stderr = %q, want the rebuild action %q", result.stderr, freshness.RebuildAction(owner.kit))
	}
}

// staleSealedCopy publishes a copy of executable beside a seal whose source digest names
// a tree no kit produced. The copy keeps an intact seal pair, so only a source-digest
// check refuses it.
func staleSealedCopy(t *testing.T, executable string) string {
	t.Helper()
	binary, err := os.ReadFile(executable)
	if err != nil {
		t.Fatal(err)
	}
	sealed, err := os.ReadFile(executable + ".seal")
	if err != nil {
		t.Fatal(err)
	}
	sources, _, err := freshness.SealDigests(executable)
	if err != nil {
		t.Fatal(err)
	}
	copied := filepath.Join(t.TempDir(), "bench")
	if err := os.WriteFile(copied, binary, 0o755); err != nil {
		t.Fatal(err)
	}
	stale := strings.Repeat("a", len(sources))
	if stale == sources {
		stale = strings.Repeat("b", len(sources))
	}
	if err := os.WriteFile(copied+".seal", []byte(strings.Replace(string(sealed), sources, stale, 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	return copied
}

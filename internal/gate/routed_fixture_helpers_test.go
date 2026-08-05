package gate

// Shared execution-observability helpers: a repository fixture observable at both
// execution layers, and accessors over the durable markers it writes. The failure class
// these observe is a verdict that credits work nobody graded — which gate ran, which
// phases ran, what the slot records — rather than return values alone.

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/canary"
)

// routedKitFixture builds a repository observable at both execution layers: the
// resolved gate script appends to .git/full-runs, and each declared phase appends its
// name to .git/phase-runs. A full run therefore leaves a gate marker and no phase
// markers; a narrowed run leaves phase markers and no gate marker — the two cannot be
// confused. The phase names carry the declaration's meaning, so the fixture guards
// both memberships before relying on them.
func routedKitFixture(t *testing.T) string {
	t.Helper()
	scope := ReducedScope()
	if !scope.Excludable(canary.PhaseTest) || scope.Excludable(conformancePhaseName) {
		t.Fatal("fixture phase names no longer match the declaration; repoint the fixture")
	}
	if !scope.Member("ROADMAP.md") || !scope.Member("capture/learnings.md") {
		t.Fatal("fixture capture paths are no longer declared; repoint the fixture")
	}
	root := t.TempDir()
	// The fixture claims the kit-root identity for itself, exactly as the contract
	// fixtures do, so these tests hold whether or not an enclosing gate run exported
	// the real kit's BENCH_KIT into the test environment.
	t.Setenv("BENCH_KIT", root)
	gitRun(t, root, "init", "-q")
	// The gate script routes through the gate-phases hand-off exactly as the kit's own
	// entry does — a stand-in binary keeps the exec inert, and the marker line above it
	// is the layer this fixture observes.
	writeGateTestFile(t, root, ".bench/gate.sh",
		"#!/usr/bin/env bash\necho full >> .git/full-runs\nexec true gate-phases \"$PWD\"\n", 0o755)
	writeGateTestFile(t, root, ".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":[],"paths":[],"tools":[]}`, 0o644)
	writeGateTestFile(t, root, ".bench/phase-conformance.sh", "echo conformance >> .git/phase-runs\n", 0o644)
	writeGateTestFile(t, root, ".bench/phase-test.sh", "echo test >> .git/phase-runs\n", 0o644)
	writeGateTestFile(t, root, canary.PhaseManifestPath, `{"phases":[`+
		`{"name":"conformance","argv":["bash",".bench/phase-conformance.sh"]},`+
		`{"name":"test","argv":["bash",".bench/phase-test.sh"]}]}`+"\n", 0o644)
	writeGateTestFile(t, root, "ROADMAP.md", "roadmap\n", 0o644)
	writeGateTestFile(t, root, "capture/learnings.md", "learnings\n", 0o644)
	writeGateTestFile(t, root, "graded.txt", "graded content\n", 0o644)
	return root
}

func markerLines(t *testing.T, root, name string) []string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, ".git", name))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	return strings.Fields(string(data))
}

func fullRunCount(t *testing.T, root string) int { return len(markerLines(t, root, "full-runs")) }

func phaseRunNames(t *testing.T, root string) []string { return markerLines(t, root, "phase-runs") }

func slotRecord(t *testing.T, root string, now time.Time) verdictRecord {
	t.Helper()
	loaded := loadVerdict(cachePath(t, root), now)
	if loaded.state != Ready {
		t.Fatalf("slot record = %s/%q, want a ready verdict", loaded.state, loaded.reason)
	}
	return loaded.record
}

func mustStrippedSubject(t *testing.T, root string) subject {
	t.Helper()
	s, err := buildStrippedSubjectForGeneration(root, mustTreeGeneration(t, root))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func commonGitDirOf(t *testing.T, root string) string {
	t.Helper()
	return gitOutput(t, root, "rev-parse", "--path-format=absolute", "--git-common-dir")
}

func mustExecuteGreen(t *testing.T, root string, engine gateEngine) {
	t.Helper()
	if got := executeWithEngine(context.Background(), root, io.Discard, io.Discard, engine); got.ActionExit != 0 {
		t.Fatalf("execution = %+v, want green", got)
	}
}

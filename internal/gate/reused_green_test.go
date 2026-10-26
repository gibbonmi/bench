package gate

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

// GR22: a primary run with no gate cache that reuses a linked-worktree green records the
// tree, green, and the source run's time in the primary gate cache.
func TestReusedGreenRecordsTheGateCache(t *testing.T) {
	primary, linked, tree := gradeKeyFixture(t)
	gradeKeyRun(t, linked, tree, reuseFreshGreen, 0)
	if err := os.Remove(outcomeCachePath(t, primary)); err != nil {
		t.Fatal(err)
	}

	gradeKeyRun(t, primary, tree, reuseFreshGreen, 0)

	requireReusedGreenCache(t, primary, linked, tree)
	requireRuns(t, primary, 1, "after a primary reuse of the linked worktree's green")
}

// GR47: a primary run whose gate cache holds a red of the same tree replaces it with the
// reused green and the source run's time.
func TestReusedGreenReplacesARedCacheOfItsTree(t *testing.T) {
	primary, linked, tree := gradeKeyFixture(t)
	red := redKeyWitness.path(t, primary)
	if err := os.WriteFile(red, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	gradeKeyRun(t, primary, tree, reuseFreshGreen, 7)
	if err := os.Remove(red); err != nil {
		t.Fatal(err)
	}
	gradeKeyRun(t, linked, tree, reuseFreshGreen, 0)

	gradeKeyRun(t, primary, tree, reuseFreshGreen, 0)

	requireReusedGreenCache(t, primary, linked, tree)
	requireRuns(t, primary, 2, "after a primary reuse over its own red")
}

// GR44: a reuse while another run holds the primary gate lock prints the reuse line and
// leaves a gate cache of another tree byte for byte.
func TestReusedGreenBesideAHeldLockWritesNothing(t *testing.T) {
	primary, linked, tree := gradeKeyFixture(t)
	gradeKeyRun(t, linked, tree, reuseFreshGreen, 0)
	before := outcomeRecord(t, primary)
	startGateLockHolder(t, primary)

	if got := gradeKeyRun(t, primary, tree, reuseFreshGreen, 0); got != reuseLine {
		t.Fatalf("primary stdout = %q, want the reuse line alone", got)
	}

	if after := outcomeRecord(t, primary); !bytes.Equal(after, before) {
		t.Fatalf("gate cache changed beside a held lock\nbefore: %s\nafter: %s", before, after)
	}
}

// requireReusedGreenCache requires that the gate cache of root holds a ready green of tree
// recorded at the time of the source checkout's green.
func requireReusedGreenCache(t *testing.T, root, source, tree string) {
	t.Helper()
	got, want := decodeGateCache(t, root), decodeGateCache(t, source)
	if got.State != Ready || got.Status != "green" || got.Tree != tree || got.RecordedAt != want.RecordedAt {
		t.Fatalf("gate cache = %+v, want a ready green of %s recorded at %s", got, tree, want.RecordedAt)
	}
}

func decodeGateCache(t *testing.T, root string) verdictRecord {
	t.Helper()
	var record verdictRecord
	if err := json.Unmarshal(outcomeRecord(t, root), &record); err != nil {
		t.Fatal(err)
	}
	return record
}

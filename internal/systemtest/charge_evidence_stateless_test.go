//go:build system

package systemtest

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/chargeevidence"
)

// storeState lists the store's objects with their sizes, excluding the two lock files every
// operation opens. A consumer cursor or reading log would appear here, whatever its name.
func (j evidenceJourney) storeState(t *testing.T) []string {
	t.Helper()
	entries, err := os.ReadDir(j.store())
	if err != nil {
		t.Fatal(err)
	}
	var state []string
	for _, entry := range entries {
		if name := entry.Name(); name == chargeevidence.OperationLockName || name == chargeevidence.WriterLockName {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			t.Fatal(err)
		}
		state = append(state, fmt.Sprintf("%s %d %s", entry.Name(), info.Size(), info.ModTime()))
	}
	sort.Strings(state)
	return state
}

// TestEvidenceStatelessProcesses is CE65: repeated and out-of-order reads in separate
// processes return the same bytes and leave the store unchanged, so no reading session or
// consumer cursor survives a read.
func TestEvidenceStatelessProcesses(t *testing.T) {
	t.Parallel()
	j := newEvidenceJourney(t)
	worktree := j.assignment(t, "stateless", "stateless\n")
	identity := identityOf(t, j.prepare(t, worktree))
	manifest := chargeevidence.Cursor{Identity: identity}.String()
	source := chargeevidence.Cursor{Identity: identity, Ordinal: 1}.String()
	before := j.storeState(t)
	// The later source page reads before the first manifest fragment, and each read repeats.
	order := []string{source, manifest, source, manifest}
	results := map[string]string{}
	for _, cursor := range order {
		result := j.read(t, worktree.path, identity, "--cursor", cursor)
		if result.code != 0 {
			t.Fatalf("read %s = (%d, %q)", cursor, result.code, result.stdout)
		}
		if seen, repeated := results[cursor]; repeated && seen != result.stdout {
			t.Fatalf("repeated read of %s returned different bytes", cursor)
		}
		results[cursor] = result.stdout
	}
	if after := j.storeState(t); strings.Join(before, "\n") != strings.Join(after, "\n") {
		t.Fatalf("reads changed the store:\nbefore:\n%s\nafter:\n%s", strings.Join(before, "\n"), strings.Join(after, "\n"))
	}
}

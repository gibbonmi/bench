// Tests for the gate row of a tip whose tree holds a reader record.
package status

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/gate/greenmarker"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

// readerLandings builds a repository whose green marker names its base commit and whose
// gate record names the base tree, with three first-parent commits after the marker. The
// commit at each index in recorded gets a reader record of its tree.
func readerLandings(t *testing.T, recorded ...int) string {
	t.Helper()
	root := initRepo(t)
	if err := os.WriteFile(filepath.Join(root, "tracked.txt"), []byte("base\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gittest.MustRun(t, root, "add", "-A")
	gittest.MustRun(t, root, "commit", "-m", "base")
	gittest.MustRun(t, root, "branch", "-M", "main")
	base := gittest.Output(t, root, "rev-parse", "HEAD")
	gittest.MustRun(t, root, "update-ref", greenmarker.Ref("main"), base)
	writeFullGateCache(t, root, gittest.Output(t, root, "rev-parse", "HEAD^{tree}"), "green")
	for landing := range 3 {
		if err := os.WriteFile(filepath.Join(root, "tracked.txt"), []byte(fmt.Sprintf("landing %d\n", landing)), 0o644); err != nil {
			t.Fatal(err)
		}
		gittest.MustRun(t, root, "commit", "-am", fmt.Sprintf("landing %d", landing))
		if slices.Contains(recorded, landing) {
			writeReaderRecord(t, root, gittest.Output(t, root, "rev-parse", "HEAD^{tree}"), base)
		}
	}
	return root
}

// writeReaderRecord installs a reader record of tree, graded against the green marker
// commit marker, in the store of root. It uses the mode and field set that the gate's
// loader requires of the reader class.
func writeReaderRecord(t *testing.T, root, tree, marker string) {
	t.Helper()
	store := filepath.Join(gittest.Output(t, root, "rev-parse", "--path-format=absolute", "--git-common-dir"), git.ReaderRecordDir)
	record, err := json.Marshal(map[string]any{
		"schema": 1, "tree": tree, "recorded_at": time.Now().UTC().Truncate(time.Second).Add(-time.Minute).Format(time.RFC3339),
		"run_binary": strings.Repeat("d", 64), "reader_set": `{"packages":["./internal/conformance"],"race":false}`,
		"destination_tree": tree, "green_marker": marker,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(store, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(store, tree), append(record, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
}

// staleRow is the row that the gate row reads for gv without its reader facts.
func staleRow(t *testing.T, gv GateInfo, root string) row {
	t.Helper()
	gv.ReaderRecorded, gv.ReaderLandings = false, 0
	rows := appendGateInfo(nil, gv, root)
	if len(rows) != 1 || !strings.HasPrefix(rows[0].detail, "stale (") {
		t.Fatalf("rows without the reader facts = %#v, want one stale row", rows)
	}
	return rows[0]
}

// TestGateRowShowsReaderGreen proves that a tip whose tree holds a reader record, with two
// reader-record commits among the three after the green marker, reads reader green with no
// command at the stale row's rank, and that no stale row prints for it.
func TestGateRowShowsReaderGreen(t *testing.T) {
	root := readerLandings(t, 1, 2)
	gv := GateVerdict(root)
	rows := appendGate(nil, root)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one gate row", rows)
	}
	got := rows[0]
	if got.signal != "gate" || got.detail != "reader green (2 of 5)" || got.action.render() != "" {
		t.Fatalf("row = %#v, want the gate row reader green (2 of 5) with no command", got)
	}
	if want := staleRow(t, gv, root).sev; got.sev != want {
		t.Fatalf("reader green rank = %d, want the stale row's rank %d", got.sev, want)
	}
	for _, r := range rows {
		if strings.HasPrefix(r.detail, "stale (") {
			t.Fatalf("rows = %#v, want no stale row for a reader-record tip", rows)
		}
	}
}

// TestGateRowKeepsTheStaleRowWithoutAReaderRecord proves that a stale tip keeps the stale row
// and its bench gate route when its own tree holds no reader record, even when an earlier
// commit after the green marker holds one. It proves the same when the tip's tree holds a
// reader record but the green marker is not an ancestor of the tip.
func TestGateRowKeepsTheStaleRowWithoutAReaderRecord(t *testing.T) {
	for _, tc := range []struct {
		name string
		repo func(t *testing.T) string
	}{
		{"no reader record on the tip", func(t *testing.T) string { return readerLandings(t, 0, 1) }},
		{"a green marker that is not an ancestor of the tip", func(t *testing.T) string {
			root := readerLandings(t, 1, 2)
			side := gittest.Output(t, root, "commit-tree", "-p", "HEAD~3", "-m", "side", "HEAD~3^{tree}")
			gittest.MustRun(t, root, "update-ref", greenmarker.Ref("main"), side)
			return root
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rows := appendGate(nil, tc.repo(t))
			if len(rows) != 1 || !strings.HasPrefix(rows[0].detail, "stale (") || rows[0].action.render() != "bench gate" {
				t.Fatalf("rows = %#v, want the stale row and its bench gate route", rows)
			}
		})
	}
}

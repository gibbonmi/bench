// Tests for the drain row's learnings component: capture/learnings.md's readability
// state and open count, split out of status_signals_test.go along this one
// responsibility (bench structure's growth budget).
package status

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/learnings"
)

// TestSignalsRendersDrainLearningsUnknownForALostDatedLine covers DL14: a journal
// holding a dated bullet is a failed read. The drain row names the journal as unknown,
// instead of the fabricated `0 open learning(s)` the decision source records.
func TestSignalsRendersDrainLearningsUnknownForALostDatedLine(t *testing.T) {
	root := initRepo(t)
	journal := learnings.JournalSchemaHeading + "\n\n- 2026-08-21 — spec anchor drift\n"
	if err := os.WriteFile(filepath.Join(root, learnings.JournalPath), []byte(journal), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")
	want := "0 idea(s), unknown (capture/learnings.md is malformed), 0 pending retro(s)"
	for _, s := range Signals(root) {
		if s.Name == "drain" {
			if s.Detail != want {
				t.Fatalf("drain detail = %q, want %q", s.Detail, want)
			}
			return
		}
	}
	t.Fatalf("Signals = %#v, want a drain row", Signals(root))
}

// TestSignalsDrainLearningsCountAgreesWithTheAXIReaderForAWriterProducedEntry closes the
// gap DL14 left open: a journal built only from entries `bench learning` itself writes,
// one of them from a title carrying a raw newline, must not make this drain row disagree
// with the learnings AXI reader (`bench learnings`) over how many entries are open. The
// two read the one file through the one parser, so a writer that stays inside its own
// reader's grammar is what keeps them agreeing.
func TestSignalsDrainLearningsCountAgreesWithTheAXIReaderForAWriterProducedEntry(t *testing.T) {
	root := initRepo(t)
	journal := learnings.JournalSchemaHeading + "\n\n" +
		learnings.FormatEntry("2026-09-11", "first", "it happened", "do the right thing", "") + "\n" +
		learnings.FormatEntry("2026-09-11", "second\nline", "it happened again", "do it again", "")
	if err := os.WriteFile(filepath.Join(root, learnings.JournalPath), []byte(journal), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")

	data, err := os.ReadFile(filepath.Join(root, learnings.JournalPath))
	if err != nil {
		t.Fatal(err)
	}
	readerRows := learnings.Rows(data)
	want := fmt.Sprintf("0 idea(s), %d open learning(s), 0 pending retro(s)", len(readerRows))
	for _, s := range Signals(root) {
		if s.Name == "drain" {
			if s.Detail != want {
				t.Fatalf("drain detail = %q, want %q (AXI reader rows=%d)", s.Detail, want, len(readerRows))
			}
			return
		}
	}
	t.Fatalf("Signals = %#v, want a drain row", Signals(root))
}

package handoffdoc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
)

// TestTwoWritersOnDistinctSectionsBothSurvive is the HS4 row. Two live phases
// close at once, each owning its own sections. Temp-and-rename alone makes each
// write atomic and still loses one: both writers read the document before either
// rename, so the later rename writes a document that never held the other's
// section.
//
// Each rewrite carries its own key, so a lost update stays lost. A writer that
// rewrote one key every round would repair the loss on its next pass, and the
// test would read green over the defect it exists to catch.
func TestTwoWritersOnDistinctSectionsBothSurvive(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "session-handoff.md")
	writers := []string{"9f2ab77", "c0ffee1"}
	const rewrites = 50

	var wait sync.WaitGroup
	errs := make([]error, len(writers))
	for i, writer := range writers {
		wait.Add(1)
		go func() {
			defer wait.Done()
			for n := 0; n < rewrites; n++ {
				key := fmt.Sprintf("%s%02d", writer, n)
				section := Section{
					Key:    key,
					Fields: []Field{{LabelRequestToken, "`" + key + "`"}},
					Next:   "`bench gate`",
					State:  "The phase is live.",
				}
				if err := WriteSection(path, section); err != nil {
					errs[i] = err
					return
				}
			}
		}()
	}
	wait.Wait()

	for i, writer := range writers {
		if errs[i] != nil {
			t.Fatalf("writer %s: %v", writer, errs[i])
		}
	}
	doc, err := Read(path)
	if err != nil {
		t.Fatalf("read the document back: %v", err)
	}
	if _, found := doc.Section(MainKey); !found {
		t.Fatal("the document lost its main section")
	}
	for _, writer := range writers {
		for n := 0; n < rewrites; n++ {
			key := fmt.Sprintf("%s%02d", writer, n)
			if _, found := doc.Section(key); !found {
				t.Fatalf("section %q is missing; the document holds %d of the %d written sections", key, len(doc.Sections)-1, len(writers)*rewrites)
			}
		}
	}
}

// TestUpdateRefusesALockAnotherWriterHolds proves the deadline. A holder that
// never releases refuses the caller by the lock path rather than hanging a phase
// close forever.
func TestUpdateRefusesALockAnotherWriterHolds(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "session-handoff.md")
	lock := LockPath(path)
	held, err := os.OpenFile(lock, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	defer held.Close()
	if err := syscall.Flock(int(held.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		t.Fatalf("hold the lock: %v", err)
	}

	err = EnsureMain(path)
	if err == nil {
		t.Fatal("Update under a held lock: want a refusal, got none")
	}
	if !strings.Contains(err.Error(), lock) {
		t.Fatalf("refusal %q does not name the lock path %q", err, lock)
	}
	if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
		t.Fatalf("the refused write left a document at %s", path)
	}
}

// TestReadTreatsAnAbsentAndAnEmptyFileAsFresh covers the empty-input edge: a first
// run in a repo needs no scaffold step ahead of it.
func TestReadTreatsAnAbsentAndAnEmptyFileAsFresh(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	for _, name := range []string{"absent.md", "empty.md"} {
		path := filepath.Join(dir, name)
		if name == "empty.md" {
			if err := os.WriteFile(path, nil, 0o644); err != nil {
				t.Fatal(err)
			}
		}
		doc, err := Read(path)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		if _, found := doc.Section(MainKey); !found || len(doc.Sections) != 1 {
			t.Fatalf("read %s returned %#v, want a fresh document holding main alone", name, doc)
		}
	}
}

// TestUpdateRewritesTheSameBytesOnARerun covers the re-run idempotency edge. A
// second close with nothing changed must not accumulate blank lines.
func TestUpdateRewritesTheSameBytesOnARerun(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "session-handoff.md")
	section := Section{Key: "9f2ab77", Next: "`bench gate`", State: "The phase is live."}
	for n := 0; n < 2; n++ {
		if err := WriteSection(path, section); err != nil {
			t.Fatalf("write %d: %v", n, err)
		}
	}
	first, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := WriteSection(path, section); err != nil {
		t.Fatalf("third write: %v", err)
	}
	second, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(first) != string(second) {
		t.Fatalf("a re-run rewrote %q, want the prior bytes %q", second, first)
	}
}

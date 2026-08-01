package retros

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
)

func TestFactsClassifiesEligibleFilesInStableOrder(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, Directory)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range map[string]string{"b.md": "second", "a.md": "first", ".hidden.md": "hidden", "note.txt": "ignored", "empty.md": ""} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	got := Facts(root)
	if got.State != bounds.StateParsed || len(got.Entries) != 3 {
		t.Fatalf("facts = %#v, want three eligible entries", got)
	}
	if got.Entries[0].Path != "capture/retros/a.md" || string(got.Entries[0].Body) != "first" || got.Entries[1].Path != "capture/retros/b.md" || got.Entries[2].State != bounds.StateEmpty {
		t.Fatalf("entries = %#v", got.Entries)
	}
}

func TestFactsKeepsDegradedAndEmptyEvidence(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, Directory)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "empty.md"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "bad.md"), []byte{0xff}, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("missing.md", filepath.Join(dir, "dangling.md")); err != nil {
		t.Fatal(err)
	}
	file, err := os.Create(filepath.Join(dir, "large.md"))
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Truncate(bounds.ControlRecordLimit + 1); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), bounds.TestDeadline(bounds.TestDeadlineFloor))
	defer cancel()
	if err := exec.CommandContext(ctx, "mkfifo", filepath.Join(dir, "wait.md")).Run(); err != nil {
		t.Fatal(err)
	}
	got := Facts(root)
	if ctx.Err() != nil {
		t.Fatalf("Facts blocked on FIFO: %v", ctx.Err())
	}
	states := map[string]bounds.FileState{}
	for _, f := range got.Entries {
		states[f.Path] = f.State
	}
	for path, want := range map[string]bounds.FileState{
		"capture/retros/empty.md":    bounds.StateEmpty,
		"capture/retros/bad.md":      bounds.StateMalformed,
		"capture/retros/dangling.md": bounds.StateUnreadable,
		"capture/retros/large.md":    bounds.StateUnreadable,
		"capture/retros/wait.md":     bounds.StateWrongType,
	} {
		if got := states[path]; got != want {
			t.Errorf("%s state = %s, want %s", path, got, want)
		}
	}
	if got.State != bounds.StateMalformed {
		t.Fatalf("aggregate state = %s, want malformed from the first stable degraded path", got.State)
	}
}

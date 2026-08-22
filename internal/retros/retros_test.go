package retros

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

// writeRetro plants one retrospective under a fresh root and returns the root.
func writeRetro(t *testing.T, name, body string) string {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, Directory)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// TestUnmarkedImprovementItemNamesPathAndLine pins RF20: the diagnostic carries the
// retro's repository-relative path and the item's line number, so the repair needs no
// search, and a marked sibling item beside it stays quiet.
func TestUnmarkedImprovementItemNamesPathAndLine(t *testing.T) {
	root := writeRetro(t, "roadmap-flow.md", strings.Join([]string{
		"## Outcome",
		"",
		"The spec landed.",
		"",
		"## Agent-experience improvements",
		"",
		"### Bench CLI",
		"",
		"- Report the board flow in the drain exit.",
		"  Feeds: FT12",
		"",
		"- Shorten the landing report.",
		"",
	}, "\n"))
	diags := ValidateImprovementMarkers(root)
	if len(diags) != 1 {
		t.Fatalf("one unmarked item beside one marked item: want one diagnostic, got %v", diags)
	}
	if !strings.Contains(diags[0], "capture/retros/roadmap-flow.md:12:") || !strings.Contains(diags[0], MissingMarkerDiagnostic) {
		t.Fatalf("diagnostic = %q, want the retro path and the item's line 12", diags[0])
	}
}

// TestAbsentAndEmptyCaptureDirectoriesStayQuiet pins RF22: a repository with no pending
// retro stays green, and absence and emptiness are the two distinct reads that reach it.
func TestAbsentAndEmptyCaptureDirectoriesStayQuiet(t *testing.T) {
	absent := t.TempDir()
	if diags := ValidateImprovementMarkers(absent); len(diags) != 0 {
		t.Fatalf("an absent capture directory: want no diagnostics, got %v", diags)
	}
	empty := t.TempDir()
	if err := os.MkdirAll(filepath.Join(empty, Directory), 0o755); err != nil {
		t.Fatal(err)
	}
	if diags := ValidateImprovementMarkers(empty); len(diags) != 0 {
		t.Fatalf("an empty capture directory: want no diagnostics, got %v", diags)
	}
}

// TestUnclassifiableRetroRedsWithItsState pins RF23: a dangling symbolic link reads as an
// absent file to a plain reader, which would count an unmarked retro compliant. The check
// reports the path and the state it was refused for instead.
func TestUnclassifiableRetroRedsWithItsState(t *testing.T) {
	root := writeRetro(t, "landed.md", "## Agent-experience improvements\n\n- Do the thing.\n  Feeds: none\n")
	if err := os.Symlink("missing.md", filepath.Join(root, Directory, "dangling.md")); err != nil {
		t.Fatal(err)
	}
	diags := ValidateImprovementMarkers(root)
	if len(diags) != 1 {
		t.Fatalf("one dangling entry beside one compliant retro: want one diagnostic, got %v", diags)
	}
	if !strings.Contains(diags[0], "capture/retros/dangling.md") || !strings.Contains(diags[0], string(bounds.StateUnreadable)) {
		t.Fatalf("diagnostic = %q, want the dangling path and its unreadable state", diags[0])
	}
}

// TestItemsOutsideTheImprovementsSectionStayQuiet pins RF24: the check grades only what
// the retro's improvements section promises, so a list item under another heading in the
// same file carries no marker duty.
func TestItemsOutsideTheImprovementsSectionStayQuiet(t *testing.T) {
	root := writeRetro(t, "landed.md", strings.Join([]string{
		"## Repair attribution",
		"",
		"- 01-report-the-board-flow | 1 round | tree-drift",
		"",
		"## Coordinator catches",
		"",
		"A delegate claimed a green gate it had not run.",
		"",
		"## Agent-experience improvements",
		"",
		"- Quote the flow report in the drain exit.",
		"  Feeds: new",
		"",
	}, "\n"))
	if diags := ValidateImprovementMarkers(root); len(diags) != 0 {
		t.Fatalf("items under other headings: want no diagnostics, got %v", diags)
	}
}

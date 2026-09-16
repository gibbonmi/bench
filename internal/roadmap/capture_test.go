package roadmap

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/learnings"
)

func TestCaptureDrainKeepsPostSnapshotLearning(t *testing.T) {
	primary := resolvedToplevel(t, newPrimaryRepo(t))
	t.Chdir(primary)
	if err := os.MkdirAll(filepath.Join(primary, "capture"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(primary, ".gitignore"), []byte("capture/IDEAS.md\ncapture/learnings.md\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var journal strings.Builder
	for i := 1; i <= 12; i++ {
		journal.WriteString(learnings.FormatEntry("2026-09-16", fmt.Sprintf("learning %d", i), "observed", "preferred", ""))
	}
	if err := os.WriteFile(filepath.Join(primary, filepath.FromSlash(learnings.JournalPath)), []byte(journal.String()), 0o644); err != nil {
		t.Fatal(err)
	}

	out, code := CaptureCommand([]string{"drain"})
	if code != 0 {
		t.Fatalf("begin drain = %q/%d", out, code)
	}
	idMatch := regexp.MustCompile(`(?m)^  ([a-z0-9-]+),sealed,`).FindStringSubmatch(out)
	if idMatch == nil {
		t.Fatalf("begin output = %q, want sealed drain id", out)
	}

	if out, code := LearningCommand([]string{"learning", "13", "--what", "arrived later", "--right", "keep it queued"}); code != 0 {
		t.Fatalf("concurrent learning = %q/%d", out, code)
	}
	snapshot, err := BuildContext(primary, true, GateCacheFact{})
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Learnings) != 12 {
		t.Fatalf("sealed context has %d learnings, want 12", len(snapshot.Learnings))
	}

	if out, code := CaptureCommand([]string{"drain", "commit", idMatch[1]}); code != 0 {
		t.Fatalf("commit drain = %q/%d", out, code)
	}
	data, err := os.ReadFile(filepath.Join(primary, filepath.FromSlash(learnings.JournalPath)))
	if err != nil {
		t.Fatal(err)
	}
	entries, malformed := learnings.Parse(data)
	if len(malformed) != 0 || len(entries) != 1 || entries[0].Title != "learning 13" {
		t.Fatalf("remaining journal = entries %#v, malformed %#v\n%s", entries, malformed, data)
	}
}

func TestCaptureCommandHelp(t *testing.T) {
	if out, code := CaptureCommand([]string{"--help"}); code != 0 || out != captureUsage+"\n" {
		t.Fatalf("help = %q/%d", out, code)
	}
}

func TestCaptureDrainAbortRestoresSealedEntriesBeforeNewEntries(t *testing.T) {
	primary := resolvedToplevel(t, newPrimaryRepo(t))
	t.Chdir(primary)
	if err := os.MkdirAll(filepath.Join(primary, "capture"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(primary, ".gitignore"), []byte("capture/IDEAS.md\ncapture/learnings.md\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(primary, filepath.FromSlash(IdeasFile)), []byte("- 2026-09-16  first\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	out, code := CaptureCommand([]string{"drain"})
	if code != 0 {
		t.Fatalf("begin drain = %q/%d", out, code)
	}
	idMatch := regexp.MustCompile(`(?m)^  ([a-z0-9-]+),sealed,`).FindStringSubmatch(out)
	if idMatch == nil {
		t.Fatalf("begin output = %q, want sealed drain id", out)
	}
	if repeated, repeatedCode := CaptureCommand([]string{"drain"}); repeatedCode != 0 || !strings.Contains(repeated, idMatch[1]) {
		t.Fatalf("repeated begin = %q/%d, want %s", repeated, repeatedCode, idMatch[1])
	}
	if out, code := IdeaCommand([]string{"second"}); code != 0 {
		t.Fatalf("concurrent idea = %q/%d", out, code)
	}
	if out, code := CaptureCommand([]string{"drain", "abort", idMatch[1]}); code != 0 {
		t.Fatalf("abort drain = %q/%d", out, code)
	}
	data, err := os.ReadFile(filepath.Join(primary, filepath.FromSlash(IdeasFile)))
	if err != nil {
		t.Fatal(err)
	}
	if first, second := strings.Index(string(data), "first"), strings.Index(string(data), "second"); first < 0 || second < 0 || first >= second {
		t.Fatalf("restored ideas = %q, want sealed entry before live entry", data)
	}
}

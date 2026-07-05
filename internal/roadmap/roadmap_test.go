package roadmap

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// newRepo initialises a git repo in a temp dir, chdirs into it (restored at test end,
// via t.Chdir), and returns its root. git.Root() resolves from the working directory.
func newRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if out, err := exec.Command("git", "-C", root, "init").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	t.Chdir(root)
	// git.Root returns the symlink-resolved toplevel; on macOS TempDir is under a
	// symlink, so resolve root the same way for comparisons/paths.
	out, err := exec.Command("git", "-C", root, "rev-parse", "--show-toplevel").Output()
	if err != nil {
		t.Fatalf("rev-parse: %v", err)
	}
	return string(out[:len(out)-1])
}

func roadmapPath(t *testing.T, root string) string {
	t.Helper()
	return filepath.Join(root, "ROADMAP.md")
}

func ideasPath(t *testing.T, root string) string {
	t.Helper()
	return filepath.Join(root, ideasFile)
}

var datedLine = regexp.MustCompile(`(?m)^- [0-9]{4}-[0-9]{2}-[0-9]{2}  ship dark mode$`)

// TestIdeaCreatesDatedLine covers idea creating the file with the dated two-space format.
func TestIdeaCreatesDatedLine(t *testing.T) {
	root := newRepo(t)
	out, code := IdeaCommand([]string{"ship", "dark", "mode"})
	if code != 0 {
		t.Fatalf("exit: got %d, want 0", code)
	}
	if out != "parked: ship dark mode\n" {
		t.Fatalf("stdout: got %q", out)
	}
	data, err := os.ReadFile(ideasPath(t, root))
	if err != nil {
		t.Fatalf("IDEAS.md not created: %v", err)
	}
	if !datedLine.Match(data) {
		t.Fatalf("line does not match dated two-space shape: %q", data)
	}
	if _, err := os.Stat(roadmapPath(t, root)); err == nil {
		t.Fatal("ROADMAP.md should not have been created")
	}
}

// TestIdeaEmptyExitsTwo covers empty and whitespace-only text: exit 2, no file created.
func TestIdeaEmptyExitsTwo(t *testing.T) {
	for _, args := range [][]string{{}, {""}, {"   ", "\t"}} {
		root := newRepo(t)
		out, code := IdeaCommand(args)
		if code != 2 {
			t.Fatalf("args %q: exit got %d, want 2", args, code)
		}
		if out != "usage: bench idea \"<text>\"\n" {
			t.Fatalf("args %q: stdout got %q", args, out)
		}
		if _, err := os.Stat(ideasPath(t, root)); err == nil {
			t.Fatalf("args %q: IDEAS.md should not have been created", args)
		}
	}
}

// TestIdeaMultiWordJoin covers unquoted multi-word args joining with single spaces.
func TestIdeaMultiWordJoin(t *testing.T) {
	root := newRepo(t)
	out, _ := IdeaCommand([]string{"capture", "all", "the", "words"})
	if out != "parked: capture all the words\n" {
		t.Fatalf("stdout: got %q", out)
	}
	data, _ := os.ReadFile(ideasPath(t, root))
	if !regexp.MustCompile(`(?m)^- [0-9-]{10}  capture all the words$`).Match(data) {
		t.Fatalf("joined line missing: %q", data)
	}
}

// TestIdeaNewlineNormalization covers appending to a hand-written last line that lacks
// a trailing newline: the new entry must not merge, so the inbox has two entries.
func TestIdeaNewlineNormalization(t *testing.T) {
	root := newRepo(t)
	if err := os.WriteFile(ideasPath(t, root), []byte("- 2026-06-01  hand added"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, code := IdeaCommand([]string{"after handedit"}); code != 0 {
		t.Fatalf("exit: got %d, want 0", code)
	}
	data, _ := os.ReadFile(ideasPath(t, root))
	if got := string(data); !strings.HasPrefix(got, "- 2026-06-01  hand added\n- ") {
		t.Fatalf("appended line did not start on a new physical line: %q", got)
	}
}

// TestRoadmapEmpty covers the empty posture for both an absent and a zero-byte file.
func TestRoadmapEmpty(t *testing.T) {
	root := newRepo(t)
	if out, code := RoadmapCommand(nil); out != "roadmap empty\n" || code != 0 {
		t.Fatalf("absent: got %q/%d", out, code)
	}
	if err := os.WriteFile(roadmapPath(t, root), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if out, code := RoadmapCommand(nil); out != "roadmap empty\n" || code != 0 {
		t.Fatalf("zero-byte: got %q/%d", out, code)
	}
}

// TestRoadmapVerbatim covers a populated file returned byte-for-byte.
func TestRoadmapVerbatim(t *testing.T) {
	root := newRepo(t)
	content := "# Roadmap\n\n- 2026-01-01  first\n- 2026-01-02  second\n"
	if err := os.WriteFile(roadmapPath(t, root), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := RoadmapCommand(nil)
	if code != 0 || out != content {
		t.Fatalf("verbatim: got %q/%d", out, code)
	}
}

// TestParkedCountMixedLines covers counting only `^- ` lines among mixed content.
func TestParkedCountMixedLines(t *testing.T) {
	root := newRepo(t)
	content := "# Roadmap\n\n- 2026-01-01  first\nnot a bullet\n-nohyphenspace\n- 2026-01-02  second\n  - indented\n"
	if err := os.WriteFile(roadmapPath(t, root), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if n := ParkedCount(root); n != 2 {
		t.Fatalf("ParkedCount got %d, want 2", n)
	}
}

// TestParkedCountAbsent covers the zero posture for a missing file.
func TestParkedCountAbsent(t *testing.T) {
	if n := ParkedCount(t.TempDir()); n != 0 {
		t.Fatalf("absent file: got %d, want 0", n)
	}
}

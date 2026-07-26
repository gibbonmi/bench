package roadmap

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
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
// An empty argument is rejected by the shared grammar and names the empty token, so it
// renders that line rather than the blank-text one the other two reach.
func TestIdeaEmptyExitsTwo(t *testing.T) {
	for _, tc := range []struct {
		args []string
		want string
	}{
		{nil, "usage: bench idea \"<text>\"\n"},
		{[]string{""}, "usage: bench idea (unknown argument: \"\")\n"},
		{[]string{"   ", "\t"}, "usage: bench idea \"<text>\"\n"},
	} {
		args := tc.args
		root := newRepo(t)
		out, code := IdeaCommand(args)
		if code != 2 {
			t.Fatalf("args %q: exit got %d, want 2", args, code)
		}
		if out != tc.want {
			t.Fatalf("args %q: stdout got %q, want %q", args, out, tc.want)
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

// TestRoadmapMissing covers the missing-roadmap posture for both an absent and a
// zero-byte file: exit 0 with a pointer to /bench-what-next, never a crash or a
// bare empty verdict.
func TestRoadmapMissing(t *testing.T) {
	root := newRepo(t)
	if out, code := RoadmapCommand(nil); out != missingRoadmap || code != 0 {
		t.Fatalf("absent: got %q/%d", out, code)
	}
	if err := os.WriteFile(roadmapPath(t, root), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if out, code := RoadmapCommand(nil); out != missingRoadmap || code != 0 {
		t.Fatalf("zero-byte: got %q/%d", out, code)
	}
	if !strings.Contains(missingRoadmap, "/bench-what-next") {
		t.Fatalf("missing-roadmap pointer does not name /bench-what-next: %q", missingRoadmap)
	}
}

// TestRoadmapVerbatim covers a populated file returned byte-for-byte before the callout.
func TestRoadmapVerbatim(t *testing.T) {
	root := newRepo(t)
	content := "# Roadmap\n\n- 2026-01-01  first\n- 2026-01-02  second\n\n## Recommended sequence\n\n1. First item - /bench-shape-idea\n2. Second item - /bench-implement-spec\n"
	if err := os.WriteFile(roadmapPath(t, root), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := RoadmapCommand(nil)
	if code != 0 || !strings.HasPrefix(out, content) {
		t.Fatalf("verbatim: got %q/%d", out, code)
	}
}

// TestRoadmapMissingSection covers a roadmap body without `## Recommended sequence`:
// exit 0, body still printed, explicit missing-section message instead of silence.
func TestRoadmapMissingSection(t *testing.T) {
	root := newRepo(t)
	content := "# Roadmap\n\n## Context\n\nNo sequence here.\n"
	if err := os.WriteFile(roadmapPath(t, root), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := RoadmapCommand(nil)
	if code != 0 {
		t.Fatalf("exit: got %d, want 0", code)
	}
	if !strings.HasPrefix(out, content) {
		t.Fatalf("body not printed verbatim: %q", out)
	}
	if !strings.Contains(out, "no ## Recommended sequence section") || !strings.Contains(out, "/bench-what-next") {
		t.Fatalf("missing-section message absent: %q", out)
	}
}

// TestRoadmapMalformedSequence covers a present section whose numbered-item count
// breaks the two-or-three contract: an explicit malformed message, not a verbatim callout.
func TestRoadmapMalformedSequence(t *testing.T) {
	cases := []struct {
		name  string
		items string
		count string
	}{
		{"no items", "", "0"},
		{"one item", "1. Only item - /bench-shape-idea\n", "1"},
		{"four items", "1. a - /x\n2. b - /x\n3. c - /x\n4. d - /x\n", "4"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := newRepo(t)
			content := "# Roadmap\n\n## Recommended sequence\n\n" + tc.items
			if err := os.WriteFile(roadmapPath(t, root), []byte(content), 0o644); err != nil {
				t.Fatal(err)
			}
			out, code := RoadmapCommand(nil)
			if code != 0 {
				t.Fatalf("exit: got %d, want 0", code)
			}
			if !strings.Contains(out, "malformed ## Recommended sequence: "+tc.count+" numbered item(s)") || !strings.Contains(out, "/bench-what-next") {
				t.Fatalf("malformed-section message absent: %q", out)
			}
			if strings.Contains(out, "## Next action") {
				t.Fatalf("malformed section still rendered as next action: %q", out)
			}
		})
	}
}

// TestRoadmapSequenceFenceBlindness covers headings inside fenced code blocks: a
// fenced `## Recommended sequence` must not start the section, and a fenced `## `
// heading inside the section must not truncate it.
func TestRoadmapSequenceFenceBlindness(t *testing.T) {
	root := newRepo(t)
	content := "# Roadmap\n\n```\n## Recommended sequence\n\n1. fake - /x\n```\n\n## Recommended sequence\n\n1. Real item - /bench-shape-idea\n\n```\n## Later\n```\n\n2. Second item - /bench-implement-spec\n\n## Later\n\nDo not include.\n"
	if err := os.WriteFile(roadmapPath(t, root), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := RoadmapCommand(nil)
	if code != 0 {
		t.Fatalf("exit: got %d, want 0", code)
	}
	callout := out[strings.LastIndex(out, "## Next action"):]
	if strings.Contains(callout, "fake - /x") {
		t.Fatalf("fenced heading hijacked the section start:\n%s", callout)
	}
	if !strings.Contains(callout, "2. Second item - /bench-implement-spec") {
		t.Fatalf("fenced heading truncated the section:\n%s", callout)
	}
	if strings.Contains(callout, "Do not include.") {
		t.Fatalf("callout included the following section:\n%s", callout)
	}
}

// TestRoadmapHeadingTrailingWhitespace covers a hand-edited heading with trailing
// spaces or tabs: it still matches instead of silently yielding no callout.
func TestRoadmapHeadingTrailingWhitespace(t *testing.T) {
	root := newRepo(t)
	content := "# Roadmap\n\n## Recommended sequence \t\n\n1. First item - /bench-shape-idea\n2. Second item - /bench-implement-spec\n"
	if err := os.WriteFile(roadmapPath(t, root), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := RoadmapCommand(nil)
	if code != 0 {
		t.Fatalf("exit: got %d, want 0", code)
	}
	if !strings.Contains(out, "## Next action") || !strings.Contains(out, "1. First item - /bench-shape-idea") {
		t.Fatalf("padded heading did not match: %q", out)
	}
}

// TestRoadmapDrainStatus covers the maintenance prompt when either drain source has rows.
func TestRoadmapDrainStatus(t *testing.T) {
	cases := []struct {
		name      string
		ideas     string
		learnings string
		want      []string
		notWanted []string
	}{
		{
			name:  "ideas only",
			ideas: "- 2026-07-05  parked\n",
			want:  []string{"ideas: 1 parked in IDEAS.md", "learnings: 0 open in .bench/learnings.md", "/bench-what-next"},
		},
		{
			name:      "learnings only",
			learnings: "## 2026-07-05 — open learning  [open]\n",
			want:      []string{"ideas: 0 parked in IDEAS.md", "learnings: 1 open in .bench/learnings.md", "/bench-what-next"},
		},
		{
			name:      "empty sources",
			ideas:     "",
			learnings: "## <date> — template  [open]\n",
			notWanted: []string{"## Drain status", "/bench-what-next"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := newRepo(t)
			if err := os.WriteFile(roadmapPath(t, root), []byte("# Roadmap\n\n## Recommended sequence\n\n1. First item - /bench-shape-idea\n2. Second item - /bench-implement-spec\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			if tc.ideas != "" {
				if err := os.WriteFile(ideasPath(t, root), []byte(tc.ideas), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			if tc.learnings != "" {
				path := filepath.Join(root, ".bench", "learnings.md")
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte(tc.learnings), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			out, code := RoadmapCommand(nil)
			if code != 0 {
				t.Fatalf("exit: got %d, want 0", code)
			}
			for _, want := range tc.want {
				if !strings.Contains(out, want) {
					t.Fatalf("missing %q in:\n%s", want, out)
				}
			}
			for _, notWanted := range tc.notWanted {
				if strings.Contains(out, notWanted) {
					t.Fatalf("unexpected %q in:\n%s", notWanted, out)
				}
			}
		})
	}
}

// TestRoadmapRecommendedSequenceCallout covers the no-drain extraction branch.
func TestRoadmapRecommendedSequenceCallout(t *testing.T) {
	root := newRepo(t)
	content := "# Roadmap\n\n## Context\n\nKeep current.\n\n## Recommended sequence\n\n1. Shape next item - /bench-shape-idea\n2. Implement next item - /bench-implement-spec\n\n## Later\n\nDo not include.\n"
	if err := os.WriteFile(roadmapPath(t, root), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := RoadmapCommand(nil)
	if code != 0 {
		t.Fatalf("exit: got %d, want 0", code)
	}
	want := "## Next action\n\n## Recommended sequence\n\n1. Shape next item - /bench-shape-idea\n2. Implement next item - /bench-implement-spec\n"
	if !strings.Contains(out, want) {
		t.Fatalf("missing recommended sequence callout %q in:\n%s", want, out)
	}
	callout := out[strings.LastIndex(out, "## Next action"):]
	if strings.Contains(callout, "Do not include.") {
		t.Fatalf("callout included the following section:\n%s", callout)
	}
}

// TestDrainCountsMixedLines covers counting only `^- ` lines among mixed inbox content.
func TestDrainCountsMixedLines(t *testing.T) {
	root := newRepo(t)
	content := "# Inbox\n\n- 2026-01-01  first\nnot a bullet\n-nohyphenspace\n- 2026-01-02  second\n  - indented\n"
	if err := os.WriteFile(ideasPath(t, root), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if ideas, _, _, _ := DrainCounts(root); ideas != 2 {
		t.Fatalf("DrainCounts ideas got %d, want 2", ideas)
	}
}

// TestDrainCountsAbsent covers the zero posture when both sources are missing.
func TestDrainCountsAbsent(t *testing.T) {
	ideas, ideasState, open, learningsState := DrainCounts(t.TempDir())
	if ideas != 0 || open != 0 {
		t.Fatalf("absent sources: got %d/%d, want 0/0", ideas, open)
	}
	if ideasState != bounds.StateParsed || learningsState != bounds.StateParsed {
		t.Fatalf("absent sources should read as the ordinary quiet state, got ideas=%s learnings=%s", ideasState, learningsState)
	}
}

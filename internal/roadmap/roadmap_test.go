package roadmap

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

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
	toplevel := string(out[:len(out)-1])
	// The capture surfaces live under a directory the repository root no longer
	// supplies for free, so a fixture root that omits it fails every write that used
	// to land beside ROADMAP.md.
	if err := os.MkdirAll(filepath.Join(toplevel, "capture"), 0o755); err != nil {
		t.Fatal(err)
	}
	return toplevel
}

func roadmapPath(t *testing.T, root string) string {
	t.Helper()
	return filepath.Join(root, "ROADMAP.md")
}

func ideasPath(t *testing.T, root string) string {
	t.Helper()
	return filepath.Join(root, IdeasFile)
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
		t.Fatalf("capture/IDEAS.md not created: %v", err)
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
			t.Fatalf("args %q: capture/IDEAS.md should not have been created", args)
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

func TestIdeaOwnedOccurrence(t *testing.T) {
	root := newRepo(t)
	if err := os.WriteFile(roadmapPath(t, root), []byte("**FT98 — active work replaces retired FT97.**\nOccurrences: baseline-01\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(ideasPath(t, root), []byte("- 2026-07-30  hand added"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := IdeaCommand([]string{"--owner", "FT98", "--incident", "new-signal", "capture", "all", "the", "words"})
	if code != 0 || out != "parked: capture all the words\n" {
		t.Fatalf("owned idea = %q/%d", out, code)
	}
	got, err := os.ReadFile(ideasPath(t, root))
	if err != nil {
		t.Fatal(err)
	}
	if !regexp.MustCompile(`(?m)^- 2026-07-30  hand added\n- [0-9-]{10}  capture all the words \[occurrence:FT98/new-signal\]$`).Match(got) {
		t.Fatalf("owned entry = %q", got)
	}
}

func TestIdeaHistoricalOwnerAbsentFromCurrentRoadmapRefuses(t *testing.T) {
	root := newRepo(t)
	if err := os.WriteFile(roadmapPath(t, root), []byte("**FT98 — retired history.**\nOccurrences: baseline-01\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, err := exec.Command("git", "-C", root, "add", RoadmapFile).CombinedOutput(); err != nil {
		t.Fatalf("stage history: %v: %s", err, out)
	}
	if out, err := exec.Command("git", "-C", root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-qm", "retire FT98").CombinedOutput(); err != nil {
		t.Fatalf("commit history: %v: %s", err, out)
	}
	if err := os.WriteFile(roadmapPath(t, root), []byte("**FT99 — current.**\nOccurrences: baseline-01\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	before := []byte("- 2026-07-30  existing idea\n")
	if err := os.WriteFile(ideasPath(t, root), before, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, code := IdeaCommand([]string{"--owner", "FT98", "--incident", "retired-signal", "refused"}); code != 1 {
		t.Fatalf("historical owner exit = %d, want 1", code)
	}
	if after, err := os.ReadFile(ideasPath(t, root)); err != nil || string(after) != string(before) {
		t.Fatalf("historical owner mutated inbox: %q, %v", after, err)
	}
}

func TestIdeaOwnedOccurrenceRefusalsPreserveInbox(t *testing.T) {
	tests := []struct {
		name    string
		roadmap string
		args    []string
		code    int
	}{
		{"missing incident", "**FT98 — active.**\nOccurrences: baseline-01\n", []string{"--owner", "FT98", "idea"}, 2},
		{"missing owner", "**FT98 — active.**\nOccurrences: baseline-01\n", []string{"--incident", "signal", "idea"}, 2},
		{"repeated owner", "**FT98 — active.**\nOccurrences: baseline-01\n", []string{"--owner", "FT98", "--owner", "FT98", "--incident", "signal", "idea"}, 2},
		{"repeated incident", "**FT98 — active.**\nOccurrences: baseline-01\n", []string{"--owner", "FT98", "--incident", "signal", "--incident", "other", "idea"}, 2},
		{"empty owner", "**FT98 — active.**\nOccurrences: baseline-01\n", []string{"--owner", "", "--incident", "signal", "idea"}, 2},
		{"empty incident", "**FT98 — active.**\nOccurrences: baseline-01\n", []string{"--owner", "FT98", "--incident", "", "idea"}, 2},
		{"flag-like owner", "**FT98 — active.**\nOccurrences: baseline-01\n", []string{"--owner", "--incident", "signal", "idea"}, 2},
		{"flag-like incident", "**FT98 — active.**\nOccurrences: baseline-01\n", []string{"--owner", "FT98", "--incident", "--owner", "idea"}, 2},
		{"malformed owner", "**FT98 — active.**\nOccurrences: baseline-01\n", []string{"--owner", "FT0", "--incident", "signal", "idea"}, 2},
		{"malformed incident", "**FT98 — active.**\nOccurrences: baseline-01\n", []string{"--owner", "FT98", "--incident", "Bad", "idea"}, 2},
		{"unknown owner", "**FT98 — active.**\nOccurrences: baseline-01\n", []string{"--owner", "FT99", "--incident", "signal", "idea"}, 1},
		{"untrusted roadmap", "**FT98 — active.**\nOccurrences: invalid_key\n", []string{"--owner", "FT98", "--incident", "signal", "idea"}, 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root := newRepo(t)
			if err := os.WriteFile(roadmapPath(t, root), []byte(tc.roadmap), 0o644); err != nil {
				t.Fatal(err)
			}
			before := []byte("- 2026-07-30  existing idea\n")
			if err := os.WriteFile(ideasPath(t, root), before, 0o644); err != nil {
				t.Fatal(err)
			}
			_, code := IdeaCommand(tc.args)
			if code != tc.code {
				t.Fatalf("exit = %d, want %d", code, tc.code)
			}
			after, err := os.ReadFile(ideasPath(t, root))
			if err != nil || string(after) != string(before) {
				t.Fatalf("inbox changed: %q, %v", after, err)
			}
		})
	}
}

// TestRoadmapMissing covers the two states a reader could collapse onto one another. An
// absent file is the maintenance-prompt posture: exit 0 with a pointer to
// /bench-what-next, never a crash or a bare empty verdict. A zero-byte file is present,
// so it takes the non-absent posture — exit 1 naming the state — and must not print the
// prompt that would tell a reader no roadmap was ever created.
func TestRoadmapMissing(t *testing.T) {
	root := newRepo(t)
	if out, code := RoadmapCommand(nil); out != missingRoadmap || code != 0 {
		t.Fatalf("absent: got %q/%d", out, code)
	}
	if err := os.WriteFile(roadmapPath(t, root), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if out, code := RoadmapCommand(nil); code != 1 || !strings.Contains(out, "error: ROADMAP.md is empty") || strings.Contains(out, "no ROADMAP.md") {
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
			want:  []string{"ideas: 1 parked in capture/IDEAS.md", "learnings: 0 open in capture/learnings.md", "/bench-what-next"},
		},
		{
			name:      "learnings only",
			learnings: "## 2026-07-05 — open learning  [open]\n",
			want:      []string{"ideas: 0 parked in capture/IDEAS.md", "learnings: 1 open in capture/learnings.md", "/bench-what-next"},
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
				path := filepath.Join(root, "capture", "learnings.md")
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

func TestRoadmapDrainStatusIncludesPendingRetros(t *testing.T) {
	root := newRepo(t)
	if err := os.WriteFile(roadmapPath(t, root), []byte("# Roadmap\n\n## Recommended sequence\n\n1. First - /bench-shape-idea\n2. Second - /bench-implement-spec\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	retro := filepath.Join(root, "capture", "retros", "done.md")
	if err := os.MkdirAll(filepath.Dir(retro), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(retro, []byte("evidence\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := RoadmapCommand(nil)
	if code != 0 || !strings.Contains(out, "retros: 1 pending in capture/retros/") {
		t.Fatalf("roadmap = %q/%d", out, code)
	}
}

func TestRoadmapDrainStatusNamesDegradedRetros(t *testing.T) {
	root := newRepo(t)
	if err := os.WriteFile(roadmapPath(t, root), []byte("# Roadmap\n\n## Recommended sequence\n\n1. First - /bench-shape-idea\n2. Second - /bench-implement-spec\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	retro := filepath.Join(root, "capture", "retros", "wait.md")
	if err := os.MkdirAll(filepath.Dir(retro), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := exec.Command("mkfifo", retro).Run(); err != nil {
		t.Fatal(err)
	}
	type result struct {
		out  string
		code int
	}
	done := make(chan result, 1)
	go func() {
		out, code := RoadmapCommand(nil)
		done <- result{out, code}
	}()
	var got result
	select {
	case got = <-done:
	case <-time.After(bounds.TestDeadline(bounds.TestDeadlineFloor)):
		t.Fatal("bench roadmap blocked on a retrospective FIFO")
	}
	out, code := got.out, got.code
	if code != 0 || !strings.Contains(out, "retros: unknown (capture/retros/ is wrong-type)") {
		t.Fatalf("roadmap = %q/%d", out, code)
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
	if drain := DrainCounts(root); drain.Ideas != 2 {
		t.Fatalf("DrainCounts ideas got %d, want 2", drain.Ideas)
	}
}

// TestDrainCountsAbsent covers the zero posture when both sources are missing.
func TestDrainCountsAbsent(t *testing.T) {
	drain := DrainCounts(t.TempDir())
	if drain.Ideas != 0 || drain.OpenLearnings != 0 || drain.Retros != 0 {
		t.Fatalf("absent sources: got %#v, want all zero", drain)
	}
	if drain.IdeasState != bounds.StateParsed || drain.LearningsState != bounds.StateParsed || drain.RetrosState != bounds.StateParsed {
		t.Fatalf("absent sources should read as the ordinary quiet state, got %#v", drain)
	}
}

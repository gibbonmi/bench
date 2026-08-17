package dashboard

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/roadmap"
	"github.com/gibbonmi/bench/internal/status"
	"github.com/gibbonmi/bench/internal/worktree"
)

var fixedClock = time.Date(2026, 7, 7, 13, 30, 0, 0, time.UTC)

// baseSnapshot is a fully-empty snapshot at the fixed clock; each test overrides the one
// field it exercises so the assertions read against a known-otherwise-empty page.
func baseSnapshot() Snapshot {
	return Snapshot{GeneratedAt: fixedClock}
}

// Row 3: the gate section shows the cached tree sha and marks a stale verdict as stale
// rather than rendering a stale green as a clean bill.
func TestRenderGateStaleShowsShaAndStale(t *testing.T) {
	s := baseSnapshot()
	s.Gate = status.GateInfo{
		Present:    true,
		Status:     "green",
		CachedTree: "cafef00dcafef00dcafef00d",
		WorkTree:   "deadbeefdeadbeefdeadbeef",
		Stale:      true,
		Timestamp:  "2026-07-07T10:00:00Z",
	}
	out := Render(s)
	if !strings.Contains(out, "cafef00d") {
		t.Errorf("gate section dropped the cached tree sha:\n%s", section(out, "Gate"))
	}
	if !strings.Contains(out, "stale") {
		t.Errorf("stale verdict not marked stale:\n%s", section(out, "Gate"))
	}
}

// Row 3 age: a past timestamp renders a humanized age against the injected clock; a
// future or unparseable one omits the age instead of rendering a negative or nonsense
// duration. Pins gateAge's guard branches so dropping one turns red here.
func TestRenderGateAge(t *testing.T) {
	cases := []struct {
		name, timestamp, want string
		absent                bool
	}{
		{name: "past renders ago", timestamp: "2026-07-07T10:30:00Z", want: "3h ago"},
		{name: "future omits age", timestamp: "2099-01-01T00:00:00Z", absent: true},
		{name: "unparseable omits age", timestamp: "not-a-timestamp", absent: true},
	}
	for _, c := range cases {
		s := baseSnapshot()
		s.Gate = status.GateInfo{Present: true, Status: "green", CachedTree: "cafef00d", WorkTree: "cafef00d", Timestamp: c.timestamp}
		out := Render(s)
		if c.absent {
			if strings.Contains(out, "ago") || strings.Contains(out, "just now") {
				t.Errorf("%s: age rendered for timestamp %q:\n%s", c.name, c.timestamp, section(out, "Gate"))
			}
			continue
		}
		if !strings.Contains(out, c.want) {
			t.Errorf("%s: gate section missing %q:\n%s", c.name, c.want, section(out, "Gate"))
		}
	}
}

// Row 3 empty edge: an absent gate cache renders a definitive empty state, not a crash.
func TestRenderGateAbsent(t *testing.T) {
	out := Render(baseSnapshot())
	if !strings.Contains(out, "No gate cache") {
		t.Errorf("absent gate did not render its empty state:\n%s", section(out, "Gate"))
	}
}

// Row 4: signals render in the ascending-severity order and membership they arrive in —
// none dropped, none re-sorted.
func TestRenderSignalsOrderedMembership(t *testing.T) {
	s := baseSnapshot()
	s.Signals = []status.Signal{
		{Severity: 0, Name: "gate", Detail: "SIG-GATE", Action: "fix before commit"},
		{Severity: 1, Name: "git", Detail: "SIG-GIT", Action: "commit on green"},
		{Severity: 4, Name: "drain", Detail: "SIG-DRAIN", Action: "/bench-what-next"},
	}
	out := Render(s)
	gate := strings.Index(out, "SIG-GATE")
	git := strings.Index(out, "SIG-GIT")
	drain := strings.Index(out, "SIG-DRAIN")
	if gate < 0 || git < 0 || drain < 0 {
		t.Fatalf("a signal was dropped (gate=%d git=%d drain=%d):\n%s", gate, git, drain, out)
	}
	if !(gate < git && git < drain) {
		t.Fatalf("signals re-ordered (gate=%d git=%d drain=%d)", gate, git, drain)
	}
}

// Row 5: the roadmap section renders the ROADMAP.md rows and the recommended sequence.
func TestRenderRoadmapRowsAndSequence(t *testing.T) {
	s := baseSnapshot()
	s.RoadmapPresent = true
	s.RoadmapText = "# Roadmap\n\nFT1 row alpha\n\n## Recommended sequence\n\n1. Shape next - /bench-shape-idea\n"
	s.Sequence = "## Recommended sequence\n\n1. Shape next - /bench-shape-idea\n"
	out := Render(s)
	if !strings.Contains(out, "FT1 row alpha") {
		t.Errorf("roadmap rows dropped:\n%s", section(out, "Roadmap"))
	}
	if !strings.Contains(out, "Recommended sequence") || !strings.Contains(out, "Shape next") {
		t.Errorf("recommended sequence dropped:\n%s", section(out, "Roadmap"))
	}
}

// Review finding P1: RoadmapText and Sequence render inside <pre>, where html/template
// already neutralizes markup, so a multi-line value must keep its real newlines rather
// than flattening to literal \n escape tokens.
func TestRenderRoadmapAndSequencePreserveNewlinesInPre(t *testing.T) {
	s := baseSnapshot()
	s.RoadmapPresent = true
	s.RoadmapText = "# Roadmap\n\nFT1 row alpha\nFT2 row beta\n"
	s.Sequence = "## Recommended sequence\n\n1. Shape next\n2. Write spec\n"
	out := Render(s)

	roadmap := section(out, "Roadmap")
	if strings.Contains(roadmap, `\n`) {
		t.Errorf("RoadmapText flattened to escaped \\n tokens instead of real newlines:\n%s", roadmap)
	}
	if !strings.Contains(out, "FT1 row alpha\nFT2 row beta") {
		t.Errorf("RoadmapText lost its real newline between rows:\n%s", roadmap)
	}
	if !strings.Contains(out, "1. Shape next\n2. Write spec") {
		t.Errorf("Sequence lost its real newline between steps:\n%s", roadmap)
	}
}

// TestGatherRenderRoadmapTextAndSequenceFromSplitTree covers story 21 (PR18) through the
// dashboard's own composition and render path: gather() calling the roadmap package's
// readers, then Render formatting the resulting Snapshot — rather than asserting
// roadmap.RoadmapText/roadmap.RecommendedSequence directly, so a dashboard-side reader
// swap or sequence-format regression reds here. Both readers already parse ROADMAP.md's
// index only, so a split tree (index plus a roadmap/ row file) must render unchanged —
// this is a regression pin, not a behavior change.
func TestGatherRenderRoadmapTextAndSequenceFromSplitTree(t *testing.T) {
	root := t.TempDir()
	const heading = "**FT7 (LOW) — x.**"
	index := "# Roadmap\n\n" + heading + "\n\n## Recommended sequence\n\n1. Shape next - /bench-shape-idea\n"
	if err := os.WriteFile(filepath.Join(root, roadmap.RoadmapFile), []byte(index), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, roadmap.RoadmapDir), 0o755); err != nil {
		t.Fatal(err)
	}
	const bodyOnlyText = "Body text that lives only in the row file."
	body := heading + "\n" + bodyOnlyText + "\n"
	if err := os.WriteFile(filepath.Join(root, roadmap.RoadmapDir, "FT7.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	snap := gather(root)
	if !snap.RoadmapPresent {
		t.Fatal("gather reported roadmap absent over a present split tree")
	}
	if snap.RoadmapText != index {
		t.Fatalf("gather RoadmapText = %q, want the index verbatim %q", snap.RoadmapText, index)
	}
	if strings.Contains(snap.RoadmapText, bodyOnlyText) {
		t.Fatalf("gather RoadmapText leaked row-file body content: %q", snap.RoadmapText)
	}
	wantSequence := "## Recommended sequence\n\n1. Shape next - /bench-shape-idea\n"
	if snap.Sequence != wantSequence {
		t.Fatalf("gather Sequence = %q, want %q", snap.Sequence, wantSequence)
	}

	out := Render(snap)
	if !strings.Contains(out, "Shape next") {
		t.Errorf("rendered page dropped the recommended sequence:\n%s", section(out, "Roadmap"))
	}
	if strings.Contains(out, bodyOnlyText) {
		t.Errorf("rendered page leaked row-file body content:\n%s", section(out, "Roadmap"))
	}
}

// Row 6: the ideas section lists each parked line.
func TestRenderIdeas(t *testing.T) {
	s := baseSnapshot()
	s.Ideas = []string{"- 2026-07-01  first idea", "- 2026-07-02  second idea"}
	out := Render(s)
	for _, want := range []string{"first idea", "second idea"} {
		if !strings.Contains(out, want) {
			t.Errorf("idea %q dropped:\n%s", want, section(out, "Ideas"))
		}
	}
}

// Row 7: the open-learnings count renders the exact integer.
func TestRenderLearningsCount(t *testing.T) {
	s := baseSnapshot()
	s.OpenLearnings = 3
	out := Render(s)
	if !strings.Contains(out, "\"count\">3<") {
		t.Errorf("open-learnings count not rendered as 3:\n%s", section(out, "learnings"))
	}
}

// Row 8: the worktree section shows each pool class.
func TestRenderWorktreeClasses(t *testing.T) {
	s := baseSnapshot()
	s.Worktrees = []worktree.Registered{
		{Path: "/repo/out", Class: worktree.ClassOutOfPool},
		{Path: "/repo/leased", Class: worktree.ClassPoolLease},
		{Path: "/repo/warm", Class: worktree.ClassPoolWarm},
	}
	out := Render(s)
	for _, class := range []string{string(worktree.ClassOutOfPool), string(worktree.ClassPoolLease), string(worktree.ClassPoolWarm)} {
		if !strings.Contains(out, class) {
			t.Errorf("worktree class %q dropped:\n%s", class, section(out, "Worktrees"))
		}
	}
}

// Row 9: the generation timestamp comes from the injected clock (deterministic renderer).
func TestRenderInjectedTimestamp(t *testing.T) {
	out := Render(baseSnapshot())
	want := fixedClock.Format(time.RFC3339)
	if !strings.Contains(out, want) {
		t.Errorf("injected timestamp %q not rendered", want)
	}
}

// Row 10: every git/file-sourced field is HTML-escaped and stripped of C0 control bytes.
func TestRenderEscapesAndSanitizes(t *testing.T) {
	s := baseSnapshot()
	s.Ideas = []string{"- danger <script>alert(1)</script>\x1b[31m\x07"}
	s.Signals = []status.Signal{{Name: "x", Detail: "<b>bold</b>", Action: "a"}}
	s.WorktreesErr = "\x1b[31m<script>alert(1)</script>\x07"
	out := Render(s)
	if !strings.Contains(out, "&lt;script&gt;") {
		t.Errorf("markup not HTML-escaped:\n%s", out)
	}
	if strings.Contains(out, "<script>") {
		t.Errorf("live <script> tag injected into the page")
	}
	if strings.Contains(out, "<b>bold</b>") {
		t.Errorf("markup in a signal detail was not escaped")
	}
	if strings.ContainsRune(out, 0x1b) || strings.ContainsRune(out, 0x07) {
		t.Errorf("raw ESC/BEL control byte passed through to the page")
	}
}

// Row 11: the page is one self-contained document — inline <style> with a dark palette, no
// external URL, no <script>.
func TestRenderSelfContained(t *testing.T) {
	out := Render(baseSnapshot())
	if !strings.Contains(out, "prefers-color-scheme: dark") {
		t.Error("no prefers-color-scheme dark rule in the style block")
	}
	for _, forbidden := range []string{"http://", "https://", "src=", "@import", "<script"} {
		if strings.Contains(out, forbidden) {
			t.Errorf("self-containment broken: page contains %q", forbidden)
		}
	}
}

// Edge of 3,5,6,8: every source absent renders a definitive empty state per section.
func TestRenderEmptyStates(t *testing.T) {
	out := Render(baseSnapshot())
	for _, want := range []string{
		"No gate cache",
		"No signals",
		"No ROADMAP.md",
		"No parked ideas",
		"No out-of-pool",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing definitive empty state %q", want)
		}
	}
}

// Edge of 10: a path with a space and a glob character renders intact and escaped.
func TestRenderPathWithSpaceAndGlob(t *testing.T) {
	s := baseSnapshot()
	s.Worktrees = []worktree.Registered{{Path: "/repo/out of *pool*", Class: worktree.ClassOutOfPool}}
	out := Render(s)
	if !strings.Contains(out, "/repo/out of *pool*") {
		t.Errorf("path with space/glob mangled:\n%s", section(out, "Worktrees"))
	}
}

// A worktree-classify failure must render as a visible error, not as the empty-pool
// message a reader would otherwise mistake for "no worktrees" — the false-empty class
// FT29 swept. It takes precedence even when Worktrees is non-empty (stale/partial data
// alongside the error is still not the truth to show).
func TestRenderWorktreeClassifyFailureIsVisible(t *testing.T) {
	s := baseSnapshot()
	s.WorktreesErr = "exit status 128: fatal: not a git repository"
	s.Worktrees = []worktree.Registered{{Path: "/repo/stale", Class: worktree.ClassOutOfPool}}
	out := Render(s)
	if !strings.Contains(out, "fatal: not a git repository") {
		t.Errorf("worktree classify failure not surfaced:\n%s", section(out, "Worktrees"))
	}
	if strings.Contains(out, "No out-of-pool, leased, or warm worktrees.") {
		t.Errorf("classify failure rendered as the empty-pool state:\n%s", section(out, "Worktrees"))
	}
	if strings.Contains(section(out, "Worktrees"), "/repo/stale") {
		t.Errorf("classify failure still rendered the (unreliable) worktree table:\n%s", section(out, "Worktrees"))
	}
}

// Row 13 (unit) + Row 2 (unit): --stdout is the only accepted flag. An unknown argument or
// a trailing token is a usage error (exit 2); --stdout renders an HTML document (exit 0).
func TestCommandArgs(t *testing.T) {
	for _, bad := range [][]string{{"--bogus"}, {"--stdout", "extra"}, {"status"}} {
		if r, c := Command(bad); c != 2 || !strings.Contains(r, "usage:") {
			t.Errorf("args %v: report %q exit %d, want usage exit 2", bad, r, c)
		}
	}
	// --stdout runs inside this repo (the test cwd) and emits the document, writing nothing.
	if r, c := Command([]string{"--stdout"}); c != 0 || !strings.HasPrefix(r, "<!DOCTYPE html>") {
		t.Errorf("--stdout: exit %d, prefix %q, want exit 0 and an HTML document", c, r[:min(20, len(r))])
	}
}

// section returns a rough slice of out around the first occurrence of marker, for readable
// failure output; it is a test aid only.
func section(out, marker string) string {
	i := strings.Index(out, marker)
	if i < 0 {
		return out
	}
	end := i + 400
	if end > len(out) {
		end = len(out)
	}
	return out[i:end]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

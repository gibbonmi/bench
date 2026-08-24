package roadmap

import (
	"fmt"
	"github.com/gibbonmi/bench/internal/usage"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/roadmap/roadmaptest"
)

// newRepo initialises a git repo in a temp dir and changes into it. t.Chdir restores the
// previous directory at the end of the test. It returns the repo's root. git.Root()
// resolves this root from the working directory.
func newRepo(t *testing.T) string {
	t.Helper()
	root := newLinkedWorktree(t, newPrimaryRepo(t))
	t.Chdir(root)
	toplevel := root
	// The capture surfaces live under a directory the repository root no longer supplies for
	// free. A fixture root that omits it fails every write that used to land beside
	// ROADMAP.md.
	if err := os.MkdirAll(filepath.Join(toplevel, "capture"), 0o755); err != nil {
		t.Fatal(err)
	}
	return toplevel
}

// newPrimaryRepo initializes a bare-history primary checkout with one commit, so a
// linked worktree can branch from it.
func newPrimaryRepo(t *testing.T) string {
	t.Helper()
	primary := t.TempDir()
	for _, args := range [][]string{
		{"init", "-q", "-b", "main"},
		{"config", "user.email", "a@b.c"},
		{"config", "user.name", "a"},
		{"commit", "-q", "--allow-empty", "-m", "root"},
	} {
		if out, err := exec.Command("git", append([]string{"-C", primary}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
	}
	return primary
}

// newLinkedWorktree adds a linked worktree of primary and returns its resolved toplevel.
// With a tracked inbox, the idea verb refuses the primary checkout, so every
// tracked-inbox write fixture lives in a linked worktree, the shape a Bench phase runs in.
func newLinkedWorktree(t *testing.T, primary string) string {
	t.Helper()
	root := filepath.Join(t.TempDir(), "linked")
	if out, err := exec.Command("git", "-C", primary, "worktree", "add", "-q", "-b", "topic", root).CombinedOutput(); err != nil {
		t.Fatalf("git worktree add: %v: %s", err, out)
	}
	return resolvedToplevel(t, root)
}

func resolvedToplevel(t *testing.T, root string) string {
	t.Helper()
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
	return filepath.Join(root, IdeasFile)
}

var datedLine = regexp.MustCompile(`(?m)^- [0-9]{4}-[0-9]{2}-[0-9]{2}  ship dark mode$`)

// TestIdeaRefusesPrimaryCheckout pins the checkout boundary for a tracked inbox: the
// primary checkout is read-only for Bench verbs, so the idea verb refuses there and
// leaves the inbox alone.
func TestIdeaRefusesPrimaryCheckout(t *testing.T) {
	primary := resolvedToplevel(t, newPrimaryRepo(t))
	t.Chdir(primary)
	out, code := IdeaCommand([]string{"ship", "dark", "mode"})
	if code != 1 || out != usage.PrimaryCheckoutRefusal()+"\n" {
		t.Fatalf("primary checkout = %q/%d, want the shared refusal on exit 1", out, code)
	}
	if _, err := os.Stat(ideasPath(t, primary)); err == nil {
		t.Fatal("refused idea still appended to the inbox")
	}
}

// TestIdeaIgnoredInboxWritesPrimary covers an ignored inbox on the primary checkout:
// the file can never dirty a landing, so the verb writes it in place and exits zero.
func TestIdeaIgnoredInboxWritesPrimary(t *testing.T) {
	primary := resolvedToplevel(t, newPrimaryRepo(t))
	if err := os.WriteFile(filepath.Join(primary, ".gitignore"), []byte("capture/IDEAS.md\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(primary)
	out, code := IdeaCommand([]string{"ship", "dark", "mode"})
	if code != 0 || out != "parked: ship dark mode\n" {
		t.Fatalf("ignored inbox on primary = %q/%d, want parked on exit 0", out, code)
	}
	data, err := os.ReadFile(ideasPath(t, primary))
	if err != nil {
		t.Fatalf("capture/IDEAS.md not created on primary: %v", err)
	}
	if !datedLine.Match(data) {
		t.Fatalf("line does not match dated two-space shape: %q", data)
	}
}

// TestIdeaIgnoredInboxRedirectsWorktreeWrite covers an ignored inbox from a linked
// worktree: the append lands in the primary checkout's copy, not the worktree copy,
// so a parked idea survives the worktree's release.
func TestIdeaIgnoredInboxRedirectsWorktreeWrite(t *testing.T) {
	primary := resolvedToplevel(t, newPrimaryRepo(t))
	if err := os.WriteFile(filepath.Join(primary, ".gitignore"), []byte("capture/IDEAS.md\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"add", ".gitignore"}, {"commit", "-q", "-m", "ignore inbox"}} {
		if out, err := exec.Command("git", append([]string{"-C", primary}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
	}
	linked := newLinkedWorktree(t, primary)
	t.Chdir(linked)
	out, code := IdeaCommand([]string{"ship", "dark", "mode"})
	if code != 0 || out != "parked: ship dark mode\n" {
		t.Fatalf("ignored inbox in worktree = %q/%d, want parked on exit 0", out, code)
	}
	if _, err := os.Stat(ideasPath(t, linked)); err == nil {
		t.Fatal("append landed in the worktree copy, not the primary checkout")
	}
	data, err := os.ReadFile(ideasPath(t, primary))
	if err != nil {
		t.Fatalf("capture/IDEAS.md not created on primary: %v", err)
	}
	if !datedLine.Match(data) {
		t.Fatalf("line does not match dated two-space shape: %q", data)
	}
}

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
// The shared grammar rejects an empty argument and names the empty token. It therefore
// renders that line rather than the blank-text line the other two cases reach.
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

// TestIdeaNewlineNormalization covers appending to a hand-written last line that lacks a
// trailing newline. The new entry must not merge with it, so the inbox ends with two
// entries.
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
	writeBoard(t, root, Row{"**FT98 — active work replaces retired FT97.**", "Next: spec\nOccurrences: baseline-01\n"})
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
	writeBoard(t, root, Row{"**FT98 — retired history.**", "Occurrences: baseline-01\n"})
	if out, err := exec.Command("git", "-C", root, "add", RoadmapFile, RoadmapDir).CombinedOutput(); err != nil {
		t.Fatalf("stage history: %v: %s", err, out)
	}
	if out, err := exec.Command("git", "-C", root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-qm", "retire FT98").CombinedOutput(); err != nil {
		t.Fatalf("commit history: %v: %s", err, out)
	}
	if err := os.Remove(filepath.Join(root, RoadmapDir, "FT98.md")); err != nil {
		t.Fatal(err)
	}
	writeBoard(t, root, Row{"**FT99 — current.**", "Occurrences: baseline-01\n"})
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
	const ledger = "Occurrences: baseline-01\n"
	tests := []struct {
		name   string
		ledger string
		args   []string
		code   int
	}{
		{"missing incident", ledger, []string{"--owner", "FT98", "idea"}, 2},
		{"missing owner", ledger, []string{"--incident", "signal", "idea"}, 2},
		{"repeated owner", ledger, []string{"--owner", "FT98", "--owner", "FT98", "--incident", "signal", "idea"}, 2},
		{"repeated incident", ledger, []string{"--owner", "FT98", "--incident", "signal", "--incident", "other", "idea"}, 2},
		{"empty owner", ledger, []string{"--owner", "", "--incident", "signal", "idea"}, 2},
		{"empty incident", ledger, []string{"--owner", "FT98", "--incident", "", "idea"}, 2},
		{"flag-like owner", ledger, []string{"--owner", "--incident", "signal", "idea"}, 2},
		{"flag-like incident", ledger, []string{"--owner", "FT98", "--incident", "--owner", "idea"}, 2},
		{"malformed owner", ledger, []string{"--owner", "FT0", "--incident", "signal", "idea"}, 2},
		{"malformed incident", ledger, []string{"--owner", "FT98", "--incident", "Bad", "idea"}, 2},
		{"unknown owner", ledger, []string{"--owner", "FT99", "--incident", "signal", "idea"}, 1},
		{"untrusted roadmap", "Occurrences: invalid_key\n", []string{"--owner", "FT98", "--incident", "signal", "idea"}, 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root := newRepo(t)
			writeBoard(t, root, Row{"**FT98 — active.**", tc.ledger})
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

func TestRoadmapBoardDocument(t *testing.T) {
	root := newRepo(t)
	board11 := make([]Row, 0, 11)
	for i := 1; i <= 11; i++ {
		board11 = append(board11, Row{fmt.Sprintf("**FT%d — row %d.**", i, i), "Next: spec\n"})
	}
	index, files := board(board11...)
	roadmaptest.WriteSplitBoard(t, root, "# Roadmap\n\n"+index+"## Recommended sequence\n\n1. First item - /bench-shape-idea\n", files)
	out, code := RoadmapCommand(nil)
	if code != 0 {
		t.Fatalf("exit = %d; stdout=%q", code, out)
	}
	document, err := axitest.DecodeDocument(out)
	if err != nil {
		t.Fatalf("stdout = %q; decode: %v", out, err)
	}
	if want := []string{"roadmap", "board", "sequence", "drain", "help"}; !reflect.DeepEqual(document.Blocks, want) {
		t.Fatalf("blocks = %q, want %q", document.Blocks, want)
	}
	roadmapRows, err := document.Rows("roadmap")
	if err != nil || len(roadmapRows) != 10 {
		t.Fatalf("roadmap rows = %v/%v, want 10", roadmapRows, err)
	}
	first := roadmapRows[0].(map[string]any)
	last := roadmapRows[9].(map[string]any)
	if first["id"] != "FT1" || last["id"] != "FT10" {
		t.Fatalf("roadmap IDs = %q through %q, want FT1 through FT10", first["id"], last["id"])
	}
	sequenceRows, err := document.Rows("sequence")
	if err != nil || len(sequenceRows) != 1 {
		t.Fatalf("sequence rows = %v/%v, want one", sequenceRows, err)
	}
	sequence := sequenceRows[0].(map[string]any)
	if fmt.Sprint(sequence["rank"]) != "1" || sequence["text"] != "First item - /bench-shape-idea" || sequence["command"] != "/bench-shape-idea" {
		t.Fatalf("sequence = %#v", sequence)
	}
	boardRows, err := document.Rows("board")
	if err != nil || len(boardRows) != 1 {
		t.Fatalf("board rows = %v/%v, want one", boardRows, err)
	}
	board := boardRows[0].(map[string]any)
	if fmt.Sprint(board["rows_shown"]) != "10" || fmt.Sprint(board["rows_total"]) != "11" || board["sequence_trusted"] != true {
		t.Fatalf("board = %#v, want shown=10 total=11 trusted=true", board)
	}
}

// TestRoadmapBoardRendersDetailFromRowFile pins PR14 (story 16). The title comes from the
// index line, while spec, occurrence_count, and occurrence_keys come from the row-file
// body. An empty ROADMAP.md still exits 1 with the record error.
func TestRoadmapBoardRendersDetailFromRowFile(t *testing.T) {
	root := newRepo(t)
	body := "Blocked until the deploy is scheduled; the spec is `specs/foo/spec.md`.\nOccurrences: alpha-1, beta-2\n"
	writeBoard(t, root, Row{"**FT1 — index title.**", body})
	out, code := RoadmapCommand(nil)
	if code != 0 {
		t.Fatalf("exit = %d; stdout=%q", code, out)
	}
	document, err := axitest.DecodeDocument(out)
	if err != nil {
		t.Fatal(err)
	}
	rows, err := document.Rows("roadmap")
	if err != nil || len(rows) != 1 {
		t.Fatalf("roadmap rows = %#v, %v", rows, err)
	}
	row := rows[0].(map[string]any)
	if row["title"] != "index title." || row["spec"] != "foo" || fmt.Sprint(row["occurrence_count"]) != "2" || row["occurrence_keys"] != "alpha-1, beta-2" {
		t.Fatalf("row = %#v", row)
	}

	root = newRepo(t)
	if err := os.WriteFile(roadmapPath(t, root), []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code = RoadmapCommand(nil)
	if code != 1 || !strings.Contains(out, "empty") {
		t.Fatalf("empty ROADMAP.md = %q/%d, want exit 1 naming empty", out, code)
	}
}

func TestRoadmapBoardRowBoundary(t *testing.T) {
	for _, count := range []int{0, 9, 10, 11} {
		t.Run(fmt.Sprintf("%d rows", count), func(t *testing.T) {
			root := newRepo(t)
			fixture := make([]Row, 0, count)
			for i := 1; i <= count; i++ {
				fixture = append(fixture, Row{fmt.Sprintf("**FT%d — row.**", i), ""})
			}
			index, files := board(fixture...)
			roadmaptest.WriteSplitBoard(t, root, "# Roadmap\n\n## Recommended sequence\n\n"+index, files)
			out, code := RoadmapCommand(nil)
			if code != 0 {
				t.Fatalf("exit = %d; stdout=%q", code, out)
			}
			document, err := axitest.DecodeDocument(out)
			if err != nil {
				t.Fatal(err)
			}
			rows, err := document.Rows("roadmap")
			if err != nil || len(rows) != min(10, count) {
				t.Fatalf("roadmap rows = %d/%v, want %d", len(rows), err, min(10, count))
			}
			board, _ := document.Rows("board")
			if len(board) != 1 || fmt.Sprint(board[0].(map[string]any)["rows_total"]) != fmt.Sprint(count) {
				t.Fatalf("board = %#v, want total=%d", board, count)
			}
		})
	}
}

func TestRoadmapBoardInputStates(t *testing.T) {
	t.Run("absent", func(t *testing.T) {
		newRepo(t)
		out, code := RoadmapCommand(nil)
		if code != 0 || !strings.Contains(out, "/bench-drain") {
			t.Fatalf("absent = %q/%d", out, code)
		}
		document, err := axitest.DecodeDocument(out)
		if err != nil {
			t.Fatal(err)
		}
		rows, _ := document.Rows("roadmap")
		sequence, _ := document.Rows("sequence")
		if len(rows) != 0 || len(sequence) != 0 {
			t.Fatalf("absent rows = %d/%d, want zero", len(rows), len(sequence))
		}
	})
	for _, tc := range []struct {
		name, content, state string
	}{
		{"empty", "", "empty"},
		{"unsupported", "not a roadmap\n", "unsupported-schema"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newRepo(t)
			if err := os.WriteFile(roadmapPath(t, root), []byte(tc.content), 0o644); err != nil {
				t.Fatal(err)
			}
			out, code := RoadmapCommand(nil)
			if code != 1 || !strings.Contains(out, tc.state) {
				t.Fatalf("%s = %q/%d", tc.name, out, code)
			}
		})
	}
	t.Run("failed read", func(t *testing.T) {
		root := newRepo(t)
		if err := os.Mkdir(roadmapPath(t, root), 0o755); err != nil {
			t.Fatal(err)
		}
		out, code := RoadmapCommand(nil)
		if code != 1 || !strings.Contains(out, "wrong-type") {
			t.Fatalf("failed read = %q/%d", out, code)
		}
	})
}

func TestRoadmapUsesOneUsageAcrossForms(t *testing.T) {
	want := "usage: bench roadmap | bench roadmap --context [--full] | bench roadmap --context --row <ID,...> | bench roadmap --flow\n"
	for _, command := range []func() (string, int){
		func() (string, int) { return RoadmapCommand([]string{"--help"}) },
		func() (string, int) {
			return ContextCommand([]string{"--help"}, func(string) GateCacheFact { return GateCacheFact{} })
		},
	} {
		out, code := command()
		if code != 0 || out != want {
			t.Fatalf("help = %q/%d, want %q/0", out, code, want)
		}
		if !strings.Contains(out, "bench roadmap |") || !strings.Contains(out, "bench roadmap --context [--full]") || !strings.Contains(out, "bench roadmap --context --row <ID,...>") || strings.Contains(out, "bench roadmap --row") {
			t.Fatalf("usage does not distinguish the three forms: %q", out)
		}
	}
}

func TestRoadmapBoardRefusesControlByte(t *testing.T) {
	root := newRepo(t)
	if err := os.WriteFile(roadmapPath(t, root), []byte("**FT1 — bad\x01 title.**\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := RoadmapCommand(nil)
	if code != 1 || !strings.Contains(out, "error: unrepresentable TOON cell") {
		t.Fatalf("control byte = %q/%d", out, code)
	}
}

func TestRoadmapBoardPendingDrainDisclosesWhatNext(t *testing.T) {
	root := newRepo(t)
	if err := os.WriteFile(roadmapPath(t, root), []byte("**FT1 — current.**\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(ideasPath(t, root), []byte("- 2026-08-13  pending idea\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	learningsPath := filepath.Join(root, "capture", "learnings.md")
	if err := os.WriteFile(learningsPath, []byte("## 2026-08-13 — pending learning  [open]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	retro := filepath.Join(root, "capture", "retros", "pending.md")
	if err := os.MkdirAll(filepath.Dir(retro), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(retro, []byte("evidence\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := RoadmapCommand(nil)
	if code != 0 {
		t.Fatalf("exit = %d; stdout=%q", code, out)
	}
	document, err := axitest.DecodeDocument(out)
	if err != nil {
		t.Fatal(err)
	}
	actions, err := document.HelpActions()
	if err != nil || len(actions) != 1 || actions[0].Cmd != "/bench-drain" {
		t.Fatalf("help = %#v/%v, want /bench-drain", actions, err)
	}
	drain, err := document.Rows("drain")
	if err != nil || len(drain) != 1 {
		t.Fatalf("drain = %#v/%v", drain, err)
	}
	facts := drain[0].(map[string]any)
	if fmt.Sprint(facts["ideas"]) != "1" || fmt.Sprint(facts["learnings"]) != "1" || fmt.Sprint(facts["retros"]) != "1" || facts["ideas_state"] != "parsed" || facts["learnings_state"] != "parsed" || facts["retros_state"] != "parsed" {
		t.Fatalf("drain facts = %#v", facts)
	}
}

func TestRecommendedSequenceGrammar(t *testing.T) {
	if got := RecommendedSequence("# Roadmap\n"); got != "" {
		t.Fatalf("missing sequence = %q", got)
	}
	if got := RecommendedSequence("## Recommended sequence\n\n1. only\n"); got != "## Recommended sequence\n\n1. only\n" {
		t.Fatalf("malformed sequence = %q", got)
	}
	fenced := "# Roadmap\n\n```\n## Recommended sequence\n1. fake\n```\n\n## Recommended sequence\n\n1. real\n\n```\n## Later\n```\n\n2. second\n\n## Later\nnope\n"
	want := "## Recommended sequence\n\n1. real\n\n```\n## Later\n```\n\n2. second\n"
	if got := RecommendedSequence(fenced); got != want {
		t.Fatalf("fenced sequence = %q, want %q", got, want)
	}
	padded := "## Recommended sequence \t\n\n1. first\n2. second\n"
	if got := RecommendedSequence(padded); got != "## Recommended sequence \t\n\n1. first\n2. second\n" {
		t.Fatalf("padded sequence = %q", got)
	}
}

func TestDrainCountsReportsDegradedRetros(t *testing.T) {
	root := newRepo(t)
	retro := filepath.Join(root, "capture", "retros", "wait.md")
	if err := os.MkdirAll(filepath.Dir(retro), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := exec.Command("mkfifo", retro).Run(); err != nil {
		t.Fatal(err)
	}
	if drain := DrainCounts(root); drain.RetrosState != bounds.StateWrongType {
		t.Fatalf("retros state = %q, want %q", drain.RetrosState, bounds.StateWrongType)
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

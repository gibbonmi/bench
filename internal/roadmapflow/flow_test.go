package roadmapflow

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/roadmap/roadmaptest"
)

// change is one planned commit against the detail directory. Every field names bare row
// IDs, so a test reads as the board movement it means rather than as paths.
type change struct {
	subject             string
	add, modify, remove []string
}

func runGit(t *testing.T, root string, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", append([]string{"-C", root}, args...)...).Output()
	if err != nil {
		t.Fatalf("git %v: %v", args, err)
	}
	return strings.TrimSpace(string(out))
}

// newFlowRepo returns a repository holding one commit that touches no detail file, so
// every fixture starts from a history whose flow is empty.
func newFlowRepo(t *testing.T) string {
	t.Helper()
	root := gittest.RepoOnBranch(t, "main")
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("fixture\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "add", "-A")
	runGit(t, root, "commit", "-q", "-m", "base")
	return root
}

// commitChange applies one planned change and returns the resulting commit identity. Each
// write carries the subject in its bytes, so a modify is always a real content change. It
// is never a no-op commit git would refuse.
func commitChange(t *testing.T, root string, c change) string {
	t.Helper()
	dir := filepath.Join(root, "roadmap")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, id := range append(append([]string{}, c.add...), c.modify...) {
		body := fmt.Sprintf("**%s — fixture.**\n%s\n", id, c.subject)
		if err := os.WriteFile(filepath.Join(dir, id+".md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, id := range c.remove {
		if err := os.Remove(filepath.Join(dir, id+".md")); err != nil {
			t.Fatal(err)
		}
	}
	runGit(t, root, "add", "-A")
	runGit(t, root, "commit", "-q", "-m", c.subject)
	return runGit(t, root, "rev-parse", "HEAD")
}

func commitChanges(t *testing.T, root string, changes ...change) []string {
	t.Helper()
	ids := make([]string, len(changes))
	for i, c := range changes {
		ids[i] = commitChange(t, root, c)
	}
	return ids
}

// writeBoard writes an index of n rows and the detail owner of each. The open mass the
// command reports is thus a real parsed board, not a count the test asserts twice.
func writeBoard(t *testing.T, root string, n int) {
	t.Helper()
	var index strings.Builder
	files := map[string]string{}
	for i := 1; i <= n; i++ {
		heading := fmt.Sprintf("**FT%d — fixture.**\n", i)
		index.WriteString(heading)
		files[fmt.Sprintf("FT%d.md", i)] = heading
	}
	roadmaptest.WriteSplitBoard(t, root, index.String(), files)
}

// runFlow calls the command with the repository as the process working directory, the
// ambient context git.Root reads.
func runFlow(t *testing.T, root string) (string, int) {
	t.Helper()
	t.Chdir(root)
	return Command([]string{"--flow"})
}

// flowCells decodes the single flow row into field-keyed cells. Reading the header's
// own field list is what keeps the assertions from depending on column order.
func flowCells(t *testing.T, out string) map[string]string {
	t.Helper()
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) < 2 {
		t.Fatalf("flow output holds no row: %q", out)
	}
	header, _, ok := strings.Cut(lines[0], "}:")
	if !ok || !strings.HasPrefix(header, "flow[1]{") {
		t.Fatalf("flow header = %q, want one flow row", lines[0])
	}
	fields := strings.Split(strings.TrimPrefix(header, "flow[1]{"), ",")
	cells := strings.Split(strings.TrimSpace(lines[1]), ",")
	if len(fields) != len(cells) {
		t.Fatalf("flow row %q does not match schema %q", lines[1], header)
	}
	decoded := map[string]string{}
	for i, field := range fields {
		decoded[field] = strings.Trim(cells[i], `"`)
	}
	return decoded
}

func requireCells(t *testing.T, cells map[string]string, want map[string]string) {
	t.Helper()
	for field, expected := range want {
		if cells[field] != expected {
			t.Fatalf("%s = %q, want %q (row %v)", field, cells[field], expected, cells)
		}
	}
}

// windowFixture is the shared history for the window's counts. It holds one drain older
// than the window, then three drains and the modify commits between them. The oldest
// drain is what an implementation that counts the whole history would wrongly include.
func windowFixture(t *testing.T, root string) (first, last string) {
	t.Helper()
	ids := commitChanges(t, root,
		change{subject: "seed the board", add: []string{"FT90", "FT91", "FT92", "FT93", "FT94"}},
		change{subject: "drain: open two", add: []string{"FT1", "FT2"}},
		change{subject: "feed three", modify: []string{"FT90", "FT91", "FT92"}},
		change{subject: "drain: open one, retire two", add: []string{"FT3"}, remove: []string{"FT90", "FT91"}},
		change{subject: "feed four", modify: []string{"FT92", "FT93", "FT94", "FT1"}},
		change{subject: "drain: open one, retire four", add: []string{"FT4"}, remove: []string{"FT92", "FT93", "FT94", "FT1"}},
		change{subject: "feed two", modify: []string{"FT2", "FT3"}},
	)
	return ids[1], ids[6]
}

// TestFlowCountsDetailFilesOverTheWindow pins RF1, RF2, and RF3.
func TestFlowCountsDetailFilesOverTheWindow(t *testing.T) {
	root := newFlowRepo(t)
	windowFixture(t, root)
	writeBoard(t, root, 7)
	out, code := runFlow(t, root)
	if code != 0 {
		t.Fatalf("exit = %d, want 0; output=%q", code, out)
	}
	requireCells(t, flowCells(t, out), map[string]string{
		"opened": "4", "fed": "9", "retired": "6", "net": "-2", "open_mass": "7", "drains": "3",
	})
}

// TestFlowNamesTheWindowBoundaries pins RF5.
func TestFlowNamesTheWindowBoundaries(t *testing.T) {
	root := newFlowRepo(t)
	ids := commitChanges(t, root,
		change{subject: "drain one", add: []string{"FT1"}},
		change{subject: "drain two", add: []string{"FT2"}},
		change{subject: "drain three", add: []string{"FT3"}},
		change{subject: "drain four", add: []string{"FT4"}},
		change{subject: "drain five", add: []string{"FT5"}},
		change{subject: "drain: this subject adds no row", modify: []string{"FT5"}},
	)
	writeBoard(t, root, 5)
	out, code := runFlow(t, root)
	if code != 0 {
		t.Fatalf("exit = %d, want 0; output=%q", code, out)
	}
	requireCells(t, flowCells(t, out), map[string]string{
		"window_from": ids[2], "window_to": ids[5], "drains": "3", "opened": "3", "fed": "1",
	})
}

// TestFlowReportsAYoungBoard pins RF6.
func TestFlowReportsAYoungBoard(t *testing.T) {
	root := newFlowRepo(t)
	ids := commitChanges(t, root,
		change{subject: "drain one", add: []string{"FT1"}},
		change{subject: "feed one", modify: []string{"FT1"}},
		change{subject: "drain two", add: []string{"FT2"}, remove: []string{"FT1"}},
	)
	writeBoard(t, root, 1)
	out, code := runFlow(t, root)
	if code != 0 {
		t.Fatalf("exit = %d, want 0; output=%q", code, out)
	}
	requireCells(t, flowCells(t, out), map[string]string{
		"drains": "2", "opened": "2", "fed": "1", "retired": "1", "net": "1",
		"window_from": ids[0], "window_to": ids[2],
	})
}

// TestFlowStatesWhetherTheTargetHolds pins RF4: a net delta of zero holds the target and
// a positive one does not.
func TestFlowStatesWhetherTheTargetHolds(t *testing.T) {
	for _, tc := range []struct {
		name       string
		changes    []change
		net, holds string
	}{
		{
			name: "net zero", net: "0", holds: "true",
			changes: []change{
				{subject: "drain one", add: []string{"FT1", "FT2"}},
				{subject: "retire both", remove: []string{"FT1", "FT2"}},
			},
		},
		{
			name: "net one", net: "1", holds: "false",
			changes: []change{
				{subject: "drain one", add: []string{"FT1"}},
				{subject: "drain two", add: []string{"FT2"}, remove: []string{"FT1"}},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := newFlowRepo(t)
			commitChanges(t, root, tc.changes...)
			writeBoard(t, root, 2)
			out, code := runFlow(t, root)
			if code != 0 {
				t.Fatalf("exit = %d, want 0; output=%q", code, out)
			}
			requireCells(t, flowCells(t, out), map[string]string{"net": tc.net, "target_met": tc.holds})
		})
	}
}

// TestFlowWithNoEventIsDefinitivelyEmpty pins RF7.
func TestFlowWithNoEventIsDefinitivelyEmpty(t *testing.T) {
	root := newFlowRepo(t)
	writeBoard(t, root, 3)
	out, code := runFlow(t, root)
	if code != 0 {
		t.Fatalf("exit = %d, want 0; output=%q", code, out)
	}
	want := "flow[0]{opened,fed,retired,net,open_mass,target_met,drains,window_from,window_to}:\nhelp[0]{cmd,why}:\n"
	if out != want {
		t.Fatalf("empty flow = %q, want %q", out, want)
	}
}

// TestFlowOutsideARepositoryRefuses pins RF8.
func TestFlowOutsideARepositoryRefuses(t *testing.T) {
	t.Chdir(t.TempDir())
	out, code := Command([]string{"--flow"})
	if code != 1 || !strings.HasPrefix(out, "error: not in a git repository — ") {
		t.Fatalf("outside a repository = %q, %d; want the structured error on exit 1", out, code)
	}
}

// TestFlowReportsADegradedBoardAsUnknown pins RF9.
func TestFlowReportsADegradedBoardAsUnknown(t *testing.T) {
	root := newFlowRepo(t)
	commitChanges(t, root,
		change{subject: "drain one", add: []string{"FT1"}},
		change{subject: "drain two", add: []string{"FT2"}},
	)
	if err := os.RemoveAll(filepath.Join(root, "roadmap")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "ROADMAP.md"), []byte("**FT1 — fixture.**\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "roadmap"), []byte("not a directory\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := runFlow(t, root)
	if code != 0 {
		t.Fatalf("exit = %d, want 0; output=%q", code, out)
	}
	requireCells(t, flowCells(t, out), map[string]string{"open_mass": "unknown", "opened": "2"})
}

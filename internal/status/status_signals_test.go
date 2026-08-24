// Tests for signal rows, their ordering, and the routes they produce in package status.
package status

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/learnings"
	"github.com/gibbonmi/bench/internal/maps"
)

func TestTimeoutGateIsDistinctHighestSeveritySignal(t *testing.T) {
	rows := appendGateInfo(nil, GateInfo{Present: true, State: "ready", Status: "timeout"}, t.TempDir())
	if len(rows) != 1 || rows[0].detail != "timeout" || rows[0].sev != 0 || rows[0].action.render() != "bench gate --fresh" {
		t.Fatalf("timeout rows = %#v", rows)
	}
}

func TestAppendSetupReportsOnlyUnadoptedRoots(t *testing.T) {
	root := t.TempDir()
	if got := appendSetup(nil, root); len(got) != 1 || got[0] != (row{0, "setup", "no .bench/", commandAction(setupAction)}) {
		t.Fatalf("unadopted setup row = %#v", got)
	}
	if err := os.Mkdir(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	if got := appendSetup(nil, root); len(got) != 0 {
		t.Fatalf("adopted setup row = %#v, want none", got)
	}
}

func TestAppendSetupStaysQuietWhenBenchPathCannotBeRead(t *testing.T) {
	root := t.TempDir()
	if err := os.Symlink(".bench", filepath.Join(root, ".bench")); err != nil {
		t.Fatal(err)
	}
	if got := appendSetup(nil, root); len(got) != 0 {
		t.Fatalf("unreadable .bench setup row = %#v, want none", got)
	}
}

func TestStagedSpecCountUsesFactsStatusReader(t *testing.T) {
	root := t.TempDir()
	write := func(path, body string) {
		t.Helper()
		full := filepath.Join(root, "specs", path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("staged/spec.md", "Status: staged\n")
	write("implemented/spec.md", "Status: implemented\n")
	write("fenced/spec.md", "```md\nStatus: staged\n```\n")
	if err := os.MkdirAll(filepath.Join(root, "specs", "missing"), 0o755); err != nil {
		t.Fatal(err)
	}
	fifo := filepath.Join(root, "specs", "fifo", "spec.md")
	if err := os.MkdirAll(filepath.Dir(fifo), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(fifo, 0o644); err != nil {
		capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
	}
	if got, _ := stagedSpecCount(root); got != 1 {
		t.Fatalf("stagedSpecCount = %d, want 1", got)
	}

	fifoRoot := t.TempDir()
	if err := syscall.Mkfifo(filepath.Join(fifoRoot, "specs"), 0o644); err != nil {
		capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
	}
	done := make(chan int, 1)
	go func() {
		count, _ := stagedSpecCount(fifoRoot)
		done <- count
	}()
	select {
	case got := <-done:
		if got != 0 {
			t.Fatalf("FIFO specs count = %d, want 0", got)
		}
	case <-time.After(time.Second):
		t.Fatal("stagedSpecCount blocked on a FIFO specs path")
	}
}

func TestAppendStagedSpecsRoutesAndOrdersBeforeRetirement(t *testing.T) {
	root := initRepo(t)
	for path, body := range map[string]string{
		"specs/staged/spec.md": "Status: staged\n",
		"specs/merged/spec.md": "Status: implemented\n",
		"reviews/staged.md":    "paired review\n",
	} {
		full := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")

	rows := appendStagedSpecs(nil, root)
	rows = appendRetirement(rows, root)
	want := []row{
		{4, "specs", "1 staged spec(s)", commandActionWithArgument(implementSpecPhaseAction, "specs/staged/spec.md")},
		{8, "specs", "1 merged spec(s) awaiting retirement", commandActionWithArgument(retireSpecAction, "<slug>")},
	}
	if !reflect.DeepEqual(rows, want) {
		t.Fatalf("staged and retirement rows = %#v, want %#v", rows, want)
	}

	second := filepath.Join(root, "specs", "second", "spec.md")
	if err := os.MkdirAll(filepath.Dir(second), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(second, []byte("Status: staged\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	rows = appendStagedSpecs(nil, root)
	if got := rows[0].action.render(); got != "/bench-implement-spec" {
		t.Fatalf("multiple staged action = %q, want bare command", got)
	}
}

func TestAppendMapsRoutesReadyOnlyWithoutUnresolvedOrInvalidMaps(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "decisions")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	ready := strings.Replace(maps.DecisionMapTemplate(), "<answer>", "Resolved.", 1)
	ready = strings.Replace(ready, "Status: shaping", "Status: ready", 1)
	if err := os.WriteFile(filepath.Join(dir, "ready.md"), []byte(ready), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := appendMaps(nil, root); !reflect.DeepEqual(got, []row{{6, "decisions", "1 ready map(s)", commandActionWithArgument(writeSpecPhaseAction, "decisions/ready.md")}}) {
		t.Fatalf("ready-only maps row = %#v", got)
	}
	if err := os.WriteFile(filepath.Join(dir, "second.md"), []byte(ready), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := appendMaps(nil, root); !reflect.DeepEqual(got, []row{{6, "decisions", "2 ready map(s)", commandAction(writeSpecPhaseAction)}}) {
		t.Fatalf("multiple-ready maps row = %#v", got)
	}
	if err := os.Remove(filepath.Join(dir, "second.md")); err != nil {
		t.Fatal(err)
	}

	shaping := strings.Replace(maps.DecisionMapTemplate(), "<answer>", "Resolved.", 1)
	if err := os.WriteFile(filepath.Join(dir, "shaping.md"), []byte(shaping), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := appendMaps(nil, root); !reflect.DeepEqual(got, []row{{6, "decisions", "1 unresolved map(s)", commandAction(shapeIdeaPhaseAction)}}) {
		t.Fatalf("ready plus shaping maps row = %#v", got)
	}

	if err := os.Remove(filepath.Join(dir, "shaping.md")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "invalid.md"), []byte("# invalid\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := appendMaps(nil, root); !reflect.DeepEqual(got, []row{{6, "decisions", "1 unresolved map(s)", commandAction(shapeIdeaPhaseAction)}}) {
		t.Fatalf("ready plus invalid maps row = %#v", got)
	}
}

func TestRenderSetupLeadsAnUnadoptedRepository(t *testing.T) {
	root := initRepo(t)
	commitFile(t, root, "tracked.txt")
	if err := os.Remove(filepath.Join(root, ".bench")); err != nil {
		t.Fatal(err)
	}
	wantBoard := "▶ bench setup  (setup)\n  setup      no .bench/                     → bench setup\n"
	if got := render(root, false); got != wantBoard {
		t.Fatalf("unadopted board = %q, want %q", got, wantBoard)
	}
	wantRoute := RouteResult{Lead: testSignal(0, "setup", "no .bench/", "bench setup")}
	if got := RouteFor(root, Signals(root), HarnessClaude); !reflect.DeepEqual(got, wantRoute) {
		t.Fatalf("unadopted route = %#v, want %#v", got, wantRoute)
	}
}

func TestSignalsOrderStagedSpecsAfterGuardsBeforeDrain(t *testing.T) {
	root := initRepo(t)
	for path, body := range map[string]string{
		".bench/lines.env":     "BENCH_CODEX_MID=test\n",
		"capture/IDEAS.md":     "- 2026-08-18  pending\n",
		"specs/staged/spec.md": "Status: staged\n",
	} {
		full := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")

	var names []string
	for _, signal := range Signals(root) {
		if signal.Name == "guards" || signal.Name == "specs" || signal.Name == "drain" {
			names = append(names, signal.Name)
		}
	}
	if !reflect.DeepEqual(names, []string{"guards", "specs", "drain"}) {
		t.Fatalf("signal order = %q, want guards, specs, drain", names)
	}
}

func TestRenderedBoardAndRouteForStagedSpecWithReviewPickup(t *testing.T) {
	root := initRepo(t)
	for path, body := range map[string]string{
		"specs/staged/spec.md": "Status: staged\n",
		"reviews/staged.md":    "pickup\n",
	} {
		full := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")

	wantBoard := "▶ /bench-implement-spec specs/staged/spec.md  (specs)\n" +
		"  specs      1 staged spec(s)               → /bench-implement-spec specs/staged/spec.md\n"
	if got := render(root, false); got != wantBoard {
		t.Fatalf("paired-review staged-spec board = %q, want %q", got, wantBoard)
	}
	wantRoute := RouteResult{Lead: testSignal(4, "specs", "1 staged spec(s)", "/bench-implement-spec specs/staged/spec.md")}
	if got := RouteFor(root, Signals(root), HarnessClaude); !reflect.DeepEqual(got, wantRoute) {
		t.Fatalf("paired-review staged-spec route = %#v, want %#v", got, wantRoute)
	}
}

func TestSignalsRenderDirtyAndUnpushedTogether(t *testing.T) {
	root := initRepo(t)
	commitFile(t, root, "tracked.txt")
	gitRun(t, root, "branch", "-M", "main")
	gitRun(t, root, "remote", "add", "origin", root)
	gitRun(t, root, "update-ref", "refs/remotes/origin/main", "HEAD")
	gitRun(t, root, "config", "branch.main.remote", "origin")
	gitRun(t, root, "config", "branch.main.merge", "refs/heads/main")
	commitFile(t, root, "ahead.txt")
	if err := os.WriteFile(filepath.Join(root, "tracked.txt"), []byte("dirty\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	want := []Signal{testSignal(1, "git", "1 dirty path, 1 unpushed commit", "/bench-final-check")}
	if got := Signals(root); !reflect.DeepEqual(got, want) {
		t.Fatalf("Signals = %#v, want %#v", got, want)
	}
}

func TestRouteForInvokesStagedSpecWhosePathContainsSpacesAheadOfDrain(t *testing.T) {
	root := initRepo(t)
	for path, body := range map[string]string{
		"capture/IDEAS.md":      "- 2026-08-18  pending\n",
		"specs/my spec/spec.md": "Status: staged\n",
	} {
		full := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")

	want := RouteResult{
		Lead: testSignal(4, "specs", "1 staged spec(s)", "/bench-implement-spec specs/my spec/spec.md"),
		RunnersUp: []Signal{
			testSignal(4, "drain", "1 idea(s), 0 open learning(s), 0 pending retro(s)", "/bench-drain"),
		},
	}
	if got := RouteFor(root, Signals(root), HarnessClaude); !reflect.DeepEqual(got, want) {
		t.Fatalf("spaced staged-spec route = %#v, want %#v", got, want)
	}
}

func TestRouteForInvokesReadyMapWhosePathContainsSpaces(t *testing.T) {
	root := initRepo(t)
	ready := strings.Replace(maps.DecisionMapTemplate(), "<answer>", "Resolved.", 1)
	ready = strings.Replace(ready, "Status: shaping", "Status: ready", 1)
	path := filepath.Join(root, "decisions", "my map.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(ready), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")

	want := RouteResult{Lead: testSignal(6, "decisions", "1 ready map(s)", "/bench-write-spec decisions/my map.md")}
	if got := RouteFor(root, Signals(root), HarnessClaude); !reflect.DeepEqual(got, want) {
		t.Fatalf("spaced ready-map route = %#v, want %#v", got, want)
	}
}

func TestGateActionNormalization(t *testing.T) {
	partial := &gate.Partition{}
	for _, tc := range []struct {
		name string
		gate GateInfo
		want string
	}{
		{"red", GateInfo{Present: true, State: string(gate.Ready), Status: "red"}, "/bench-debug"},
		{"timeout", GateInfo{Present: true, State: string(gate.Ready), Status: "timeout"}, "bench gate --fresh"},
		{"invalid", GateInfo{Present: true, State: string(gate.Invalid)}, "bench gate"},
		{"unavailable", GateInfo{Present: true, State: string(gate.Unavailable)}, "bench gate --fresh"},
		{"drifted", GateInfo{Present: true, State: string(gate.Ready), Stale: true, CachedTree: "old", WorkTree: "new"}, "bench gate"},
		{"exact-tip partial", GateInfo{Present: true, State: string(gate.Ready), Stale: true, CachedTree: "tree", WorkTree: "tree", Partition: partial}, "bench gate --fresh"},
		{"locked-pending", GateInfo{Present: true, State: string(gate.Pending), PendingStatus: "locked-pending"}, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rows := appendGateInfo(nil, tc.gate, t.TempDir())
			if len(rows) != 1 || rows[0].action.render() != tc.want {
				t.Fatalf("action = %#v, want %q", rows, tc.want)
			}
		})
	}
}

// TestSignalsRendersDrainLearningsUnknownForALostDatedLine covers DL14: a journal
// holding a dated bullet is a failed read. The drain row names the journal as unknown,
// instead of the fabricated `0 open learning(s)` the decision source records.
func TestSignalsRendersDrainLearningsUnknownForALostDatedLine(t *testing.T) {
	root := initRepo(t)
	journal := learnings.JournalSchemaHeading + "\n\n- 2026-08-21 — spec anchor drift\n"
	if err := os.WriteFile(filepath.Join(root, learnings.JournalPath), []byte(journal), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")
	want := "0 idea(s), unknown (capture/learnings.md is malformed), 0 pending retro(s)"
	for _, s := range Signals(root) {
		if s.Name == "drain" {
			if s.Detail != want {
				t.Fatalf("drain detail = %q, want %q", s.Detail, want)
			}
			return
		}
	}
	t.Fatalf("Signals = %#v, want a drain row", Signals(root))
}

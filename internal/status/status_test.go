package status

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/maps"
	"github.com/gibbonmi/bench/internal/roadmap"
	"github.com/gibbonmi/bench/internal/roadmap/roadmaptest"
)

func TestTimeoutGateIsDistinctHighestSeveritySignal(t *testing.T) {
	rows := appendGateInfo(nil, GateInfo{Present: true, State: "ready", Status: "timeout"}, t.TempDir())
	if len(rows) != 1 || rows[0].detail != "timeout" || rows[0].sev != 0 || rows[0].action != "bench gate --fresh" {
		t.Fatalf("timeout rows = %#v", rows)
	}
}

func TestAppendSetupReportsOnlyUnadoptedRoots(t *testing.T) {
	root := t.TempDir()
	if got := appendSetup(nil, root); len(got) != 1 || got[0] != (row{0, "setup", "no .bench/", "bench setup"}) {
		t.Fatalf("unadopted setup row = %#v", got)
	}
	if err := os.Mkdir(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	if got := appendSetup(nil, root); len(got) != 0 {
		t.Fatalf("adopted setup row = %#v, want none", got)
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
		{4, "specs", "1 staged spec(s)", "/bench-implement-spec specs/staged/spec.md"},
		{8, "specs", "1 merged spec(s) awaiting retirement", "bench spec retire <slug>"},
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
	if got := rows[0].action; got != "/bench-implement-spec" {
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
	if got := appendMaps(nil, root); !reflect.DeepEqual(got, []row{{6, "decisions", "1 ready map(s)", "/bench-write-spec decisions/ready.md"}}) {
		t.Fatalf("ready-only maps row = %#v", got)
	}
	if err := os.WriteFile(filepath.Join(dir, "second.md"), []byte(ready), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := appendMaps(nil, root); !reflect.DeepEqual(got, []row{{6, "decisions", "2 ready map(s)", "/bench-write-spec"}}) {
		t.Fatalf("multiple-ready maps row = %#v", got)
	}
	if err := os.Remove(filepath.Join(dir, "second.md")); err != nil {
		t.Fatal(err)
	}

	shaping := strings.Replace(maps.DecisionMapTemplate(), "<answer>", "Resolved.", 1)
	if err := os.WriteFile(filepath.Join(dir, "shaping.md"), []byte(shaping), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := appendMaps(nil, root); !reflect.DeepEqual(got, []row{{6, "decisions", "1 unresolved map(s)", "/bench-shape-idea"}}) {
		t.Fatalf("ready plus shaping maps row = %#v", got)
	}

	if err := os.Remove(filepath.Join(dir, "shaping.md")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "invalid.md"), []byte("# invalid\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := appendMaps(nil, root); !reflect.DeepEqual(got, []row{{6, "decisions", "1 unresolved map(s)", "/bench-shape-idea"}}) {
		t.Fatalf("ready plus invalid maps row = %#v", got)
	}
}

func TestRenderSetupLeadsAnUnadoptedRepository(t *testing.T) {
	root := initRepo(t)
	if err := os.Remove(filepath.Join(root, ".bench")); err != nil {
		t.Fatal(err)
	}
	if got := render(root, false); !strings.HasPrefix(got, "▶ bench setup  (setup)\n") {
		t.Fatalf("unadopted board = %q", got)
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
			if len(rows) != 1 || rows[0].action != tc.want {
				t.Fatalf("action = %#v, want %q", rows, tc.want)
			}
		})
	}
}

func TestAllProducibleBoardActionsAreInvocableOrEmpty(t *testing.T) {
	cases := []struct {
		name string
		rows func(*testing.T) []row
	}{
		{"gate locked", func(t *testing.T) []row {
			return appendGateInfo(nil, GateInfo{Present: true, State: string(gate.Pending), PendingStatus: "locked-pending"}, t.TempDir())
		}},
		{"gate interrupted", func(t *testing.T) []row {
			return appendGateInfo(nil, GateInfo{Present: true, State: string(gate.Pending)}, t.TempDir())
		}},
		{"gate invalid", func(t *testing.T) []row {
			return appendGateInfo(nil, GateInfo{Present: true, State: string(gate.Invalid)}, t.TempDir())
		}},
		{"gate unavailable", func(t *testing.T) []row {
			return appendGateInfo(nil, GateInfo{Present: true, State: string(gate.Unavailable)}, t.TempDir())
		}},
		{"gate drifted", func(t *testing.T) []row {
			return appendGateInfo(nil, GateInfo{Present: true, State: string(gate.Ready), Stale: true, CachedTree: "old", WorkTree: "new"}, t.TempDir())
		}},
		{"gate timeout", func(t *testing.T) []row {
			return appendGateInfo(nil, GateInfo{Present: true, State: string(gate.Ready), Status: "timeout"}, t.TempDir())
		}},
		{"gate red", func(t *testing.T) []row {
			return appendGateInfo(nil, GateInfo{Present: true, State: string(gate.Ready), Status: "red"}, t.TempDir())
		}},
		{"git unavailable", func(t *testing.T) []row { return appendGit(nil, t.TempDir(), Query{}) }},
		{"intent unavailable", func(t *testing.T) []row { return appendIntent(nil, t.TempDir()) }},
		{"intent live", func(t *testing.T) []row {
			root := initRepo(t)
			if err := intent.Upsert(root, intent.Entry{Key: "live", Kind: intent.KindWorktree, CreatedAt: time.Unix(1, 0), Worktree: root}); err != nil {
				t.Fatal(err)
			}
			return appendIntent(nil, root)
		}},
		{"worktree typed", func(t *testing.T) []row {
			root := initRepo(t)
			gittest.StubGit(t, root, "fail-rev-parse", filepath.Join(t.TempDir(), "argv"))
			return appendWorktree(nil, root)
		}},
		{"worktree porcelain", func(t *testing.T) []row {
			root := initRepo(t)
			gittest.StubGit(t, root, "fail-worktree", filepath.Join(t.TempDir(), "argv"))
			return appendWorktree(nil, root)
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, produced := range tc.rows(t) {
				if produced.action != "" && !IsInvocable(produced.action) {
					t.Errorf("board action %q is not invocable", produced.action)
				}
			}
		})
	}
}

// retirementCount reads specs/*/spec.md through the shared predicate. The seam is the
// directory: write spec files, then assert the count. The predicate cases (fence, CRLF,
// trailing whitespace, wrong status) are the load-bearing behaviors.
func TestRetirementCount(t *testing.T) {
	cases := []struct {
		name  string
		files map[string]string
		want  int
	}{
		{
			name:  "unfenced implemented marker is counted",
			files: map[string]string{"a/spec.md": "# spec\nStatus: implemented\n"},
			want:  1,
		},
		{
			name:  "staged status is not counted",
			files: map[string]string{"a/spec.md": "# spec\nStatus: staged\n"},
			want:  0,
		},
		{
			name:  "no status line is not counted",
			files: map[string]string{"a/spec.md": "# spec\njust prose\n"},
			want:  0,
		},
		{
			name:  "implemented inside a code fence is not counted",
			files: map[string]string{"a/spec.md": "# spec\n```\nStatus: implemented\n```\n"},
			want:  0,
		},
		{
			name:  "CRLF line endings still match",
			files: map[string]string{"a/spec.md": "# spec\r\nStatus: implemented\r\n"},
			want:  1,
		},
		{
			name:  "trailing tabs and spaces after the value still match",
			files: map[string]string{"a/spec.md": "Status: implemented\t \n"},
			want:  1,
		},
		{
			name: "multiple specs counted correctly",
			files: map[string]string{
				"a/spec.md": "Status: implemented\n",
				"b/spec.md": "Status: staged\n",
				"c/spec.md": "Status:\timplemented\n",
			},
			want: 2,
		},
		{
			name:  "hidden spec file is ignored",
			files: map[string]string{".hidden/spec.md": "Status: implemented\n"},
			want:  0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			for name, body := range tc.files {
				p := filepath.Join(root, "specs", name)
				if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			if got := retirementCount(root); got != tc.want {
				t.Errorf("retirementCount = %d, want %d", got, tc.want)
			}
		})
	}
}

// Absent specs/ directory → 0, never a panic on the missing path.
func TestRetirementCountNoSpecsDir(t *testing.T) {
	if got := retirementCount(t.TempDir()); got != 0 {
		t.Errorf("retirementCount with no specs/ = %d, want 0", got)
	}
}

func TestOrphanedPickupCount(t *testing.T) {
	root := t.TempDir()
	for name, body := range map[string]string{
		"reviews/paired.md":    "pickup\n",
		"reviews/orphaned.md":  "pickup\n",
		"specs/paired/spec.md": "Status: staged\n",
	} {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if got := orphanedPickupCount(root); got != 1 {
		t.Fatalf("orphanedPickupCount = %d, want 1", got)
	}
}

func TestRoadmapReconcileCounts(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "specs", "merged"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "specs", "merged", "spec.md"), []byte("Status: implemented\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	content := "specs/merged/spec.md specs/retired/spec.md specs/<slug>/spec.md\n```\nspecs/fenced/spec.md\n```\n"
	if err := os.WriteFile(filepath.Join(root, roadmap.RoadmapFile), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	merged, dangling, state := roadmapReconcileCounts(root)
	if merged != 1 || dangling != 1 || state.Failed() {
		t.Fatalf("roadmapReconcileCounts = (%d, %d, %s), want (1, 1, parsed)", merged, dangling, state)
	}
}

// TestRoadmapReconcileCountsFromRowFile covers story 20 (PR17): the split board keeps a
// row's spec path in roadmap/FT7.md's body, not ROADMAP.md's index line, so the reconcile
// scan must classify a spec named only there — a scan of ROADMAP.md alone counts zero.
func TestRoadmapReconcileCountsFromRowFile(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "specs", "merged"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "specs", "merged", "spec.md"), []byte("Status: implemented\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	const heading = "**FT7 (LOW) — x.**"
	body := heading + "\nMerged into specs/merged/spec.md.\n"
	roadmaptest.WriteSplitBoard(t, root, heading+"\n", map[string]string{"FT7.md": body})
	merged, dangling, state := roadmapReconcileCounts(root)
	if merged != 1 || dangling != 0 || state.Failed() {
		t.Fatalf("roadmapReconcileCounts = (%d, %d, %s), want (1, 0, parsed)", merged, dangling, state)
	}
}

// short guards the [:7] tree-prefix slice against a short or "none" hash.
func TestShort(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"", ""},
		{"none", "none"},
		{"abcdef", "abcdef"},
		{"0123456789", "0123456"},
	} {
		if got := Short(tc.in); got != tc.want {
			t.Errorf("Short(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestStaleGateDetailActionCurrentTreeNoneFailsClosed(t *testing.T) {
	detail, action := staleGateDetailAction(t.TempDir(), "0123456789abcdef", "none")
	if detail != "stale (gated tree 0123456, work tree none)" {
		t.Fatalf("detail = %q, want strong stale detail", detail)
	}
	if action != "bench gate" {
		t.Fatalf("action = %q, want bench gate", action)
	}
}

// The reduced verdict class is retired. A legacy on-disk reduced record must read as an
// invalid cache rather than as a green of any width, so the board reports it as invalid
// instead of rendering a narrowness row for evidence nothing can validate.
func TestLegacyReducedCacheReadsAsInvalid(t *testing.T) {
	root := initRepo(t)
	tree := treeOf(t, root, map[string]string{"capture/IDEAS.md": "- an idea\n"})
	gitdir := gitRun(t, root, "rev-parse", "--absolute-git-dir")
	recorded := time.Now().UTC().Truncate(time.Second).Add(-time.Minute).Format(time.RFC3339)
	record := fmt.Sprintf(`{"schema":1,"state":"ready","status":"green","tree":%q,"oracle":%q,"recorded_at":%q,"reduced":true,"phases":["conformance"],"ancestor":%q,"ancestor_recorded_at":%q}`+"\n",
		tree, strings.Repeat("0", 64), recorded, strings.Repeat("a", 40), recorded)
	if err := os.WriteFile(filepath.Join(gitdir, git.GateCacheFile), []byte(record), 0o600); err != nil {
		t.Fatal(err)
	}
	gv := GateVerdict(root)
	if !gv.Present || gv.State != "invalid" || gv.Status == "green" {
		t.Fatalf("verdict = %#v, want a legacy reduced record read as an invalid cache", gv)
	}
}

// The synthetic cache keeps exact-tip narrowness observable without establishing a
// composed green, so the status row must describe a partial verdict rather than drift.
func TestStatusRendersAPartialVerdict(t *testing.T) {
	root := initRepo(t)
	tree := treeOf(t, root, map[string]string{"f.txt": "x\n"})
	writePartialGateCache(t, root, tree, "docs", "frontend")

	gv := GateVerdict(root)
	if gv.Partition == nil || !gv.Stale || gv.CachedTree != gv.WorkTree {
		t.Fatalf("verdict = %#v, want a partial non-reusable green over the current tree", gv)
	}
	rows := appendGateInfo(nil, gv, root)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one gate row", rows)
	}
	if strings.Contains(rows[0].detail, "stale") {
		t.Errorf("detail = %q, want a partial row rather than a stale one", rows[0].detail)
	}
	if strings.Contains(rows[0].detail, "reduced") {
		t.Errorf("detail = %q, want a partial row rather than a reduced one", rows[0].detail)
	}
	for _, name := range []string{"docs", "frontend"} {
		if !strings.Contains(rows[0].detail, name) {
			t.Errorf("detail = %q, want it to name skipped component %q", rows[0].detail, name)
		}
	}
}

func TestStatusRendersCheckOnlyPartialVerdict(t *testing.T) {
	root := initRepo(t)
	tree := treeOf(t, root, map[string]string{"f.txt": "x\n"})
	writeCheckPartialGateCache(t, root, tree, "line-routing")

	gv := GateVerdict(root)
	if gv.CheckPartition == nil || !gv.Stale || gv.CachedTree != gv.WorkTree {
		t.Fatalf("verdict = %#v, want a check-only partial verdict over the current tree", gv)
	}
	rows := appendGateInfo(nil, gv, root)
	if len(rows) != 1 || !strings.Contains(rows[0].detail, "partial green") || !strings.Contains(rows[0].detail, "line-routing") {
		t.Fatalf("rows = %#v, want a check-only partial row", rows)
	}
	if strings.Contains(rows[0].detail, "stale (gated tree") {
		t.Fatalf("rows = %#v, want a partial row rather than drift", rows)
	}
}

// A partial verdict whose tree has since moved is still drift, exactly as a reduced one is:
// narrowness and staleness stay independent.
func TestPartialVerdictOnAMovedTreeIsDrift(t *testing.T) {
	root := initRepo(t)
	gated := treeOf(t, root, map[string]string{"f.txt": "x\n"})
	current := treeOf(t, root, map[string]string{"f.txt": "x // drift\n"})
	writePartialGateCache(t, root, gated, "docs")

	gv := GateVerdict(root)
	if gv.Partition == nil || gv.WorkTree != current || gv.CachedTree != gated {
		t.Fatalf("verdict = %#v, want a partial verdict over a drifted work tree", gv)
	}
	rows := appendGateInfo(nil, gv, root)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one gate row", rows)
	}
	if !strings.HasPrefix(rows[0].detail, "stale (gated tree") {
		t.Errorf("detail = %q, want the drift row rather than the partial row", rows[0].detail)
	}
	if strings.Contains(rows[0].detail, "docs") {
		t.Errorf("detail = %q, want no skipped-component name once the tree has moved", rows[0].detail)
	}
}

// The partial row's action is the operator's one lever: a fresh whole-tree run. It never
// names the bare `bench gate`, which would repeat the same partial verdict.
func TestPartialRowActionIsFresh(t *testing.T) {
	root := initRepo(t)
	writePartialGateCache(t, root, treeOf(t, root, map[string]string{"f.txt": "x\n"}), "docs")

	rows := appendGateInfo(nil, GateVerdict(root), root)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one gate row", rows)
	}
	if !strings.Contains(rows[0].action, "bench gate --fresh") {
		t.Errorf("action = %q, want the fresh whole-tree action", rows[0].action)
	}
}

// A red recorded against a tree the work tree has since left describes that run, not this
// one. The board must send the reader back to the gate rather than headline a red for work
// that is no longer in the tree — the drifted record is stale, whatever verdict it carries.
func TestDriftedRedVerdictRendersAsStaleRatherThanRed(t *testing.T) {
	root := initRepo(t)
	gated := treeOf(t, root, map[string]string{"f.txt": "x // red\n"})
	current := treeOf(t, root, map[string]string{"f.txt": "x\n"})
	writeFullGateCache(t, root, gated, "red")

	gv := GateVerdict(root)
	if gv.CachedTree != gated || gv.WorkTree != current {
		t.Fatalf("verdict = %#v, want a red recorded against a tree the work tree has left", gv)
	}
	if !gv.Stale {
		t.Fatalf("verdict = %#v, want the drifted red marked stale", gv)
	}
	rows := appendGateInfo(nil, gv, root)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one gate row", rows)
	}
	if !strings.HasPrefix(rows[0].detail, "stale (gated tree") || rows[0].action != "bench gate" {
		t.Fatalf("rows = %#v, want the drift row rather than a red one", rows)
	}
}

// writeFullGateCache installs a ready full-class record naming cachedTree with the given
// verdict, at the mode and field set the gate's loader requires of the full class.
func writeFullGateCache(t *testing.T, root, cachedTree, status string) {
	t.Helper()
	gitdir := gitRun(t, root, "rev-parse", "--absolute-git-dir")
	recorded := time.Now().UTC().Truncate(time.Second).Add(-time.Minute).Format(time.RFC3339)
	record := fmt.Sprintf(`{"schema":1,"state":"ready","status":%q,"tree":%q,"oracle":%q,"recorded_at":%q}`+"\n",
		status, cachedTree, strings.Repeat("0", 64), recorded)
	if err := os.WriteFile(filepath.Join(gitdir, git.GateCacheFile), []byte(record), 0o600); err != nil {
		t.Fatal(err)
	}
}

// writePartialGateCache installs a ready partial green naming cachedTree, skipping the
// given components, at the mode and exact field set the gate's loader requires of the
// partial class. Each skipped component carries ancestor-form evidence (an identity and
// the time it was authored), the simpler of the two forms validatePartition accepts.
func writePartialGateCache(t *testing.T, root, cachedTree string, skipped ...string) {
	t.Helper()
	gitdir := gitRun(t, root, "rev-parse", "--absolute-git-dir")
	recorded := time.Now().UTC().Truncate(time.Second).Add(-time.Minute).Format(time.RFC3339)
	authoredAt := time.Now().UTC().Truncate(time.Second).Add(-time.Hour).Format(time.RFC3339)
	identity := strings.Repeat("b", 64)

	executedJSON, err := json.Marshal([]string{"core"})
	if err != nil {
		t.Fatal(err)
	}
	skippedJSON, err := json.Marshal(skipped)
	if err != nil {
		t.Fatal(err)
	}
	evidence := make(map[string]map[string]string, len(skipped))
	for _, component := range skipped {
		evidence[component] = map[string]string{"identity": identity, "authored_at": authoredAt}
	}
	evidenceJSON, err := json.Marshal(evidence)
	if err != nil {
		t.Fatal(err)
	}

	record := fmt.Sprintf(`{"schema":1,"state":"ready","status":"green","tree":%q,"oracle":%q,"recorded_at":%q,"executed":%s,"skipped":%s,"skip_evidence":%s}`+"\n",
		cachedTree, strings.Repeat("0", 64), recorded, executedJSON, skippedJSON, evidenceJSON)
	path := filepath.Join(gitdir, git.GateCacheFile)
	if err := os.WriteFile(path, []byte(record), 0o600); err != nil {
		t.Fatal(err)
	}
}

// writeCheckPartialGateCache installs a loader-valid check-only partition so the status
// adapter must carry the gate inspection's check partition through to its public row.
func writeCheckPartialGateCache(t *testing.T, root, cachedTree, inheritedName string) {
	t.Helper()
	gitdir := gitRun(t, root, "rev-parse", "--absolute-git-dir")
	recorded := time.Now().UTC().Truncate(time.Second).Add(-time.Minute).Format(time.RFC3339)
	authoredAt := time.Now().UTC().Truncate(time.Second).Add(-time.Hour).Format(time.RFC3339)
	var executed, inherited []string
	for _, check := range registry.Checks {
		if !check.RunsAt(registry.Dev) {
			continue
		}
		if check.Name == inheritedName {
			inherited = append(inherited, check.Name)
		} else {
			executed = append(executed, check.Name)
		}
	}
	executedJSON, err := json.Marshal(executed)
	if err != nil {
		t.Fatal(err)
	}
	inheritedJSON, err := json.Marshal(inherited)
	if err != nil {
		t.Fatal(err)
	}
	evidenceJSON, err := json.Marshal(map[string]map[string]string{
		inheritedName: {"identity": strings.Repeat("b", 64), "authored_at": authoredAt},
	})
	if err != nil {
		t.Fatal(err)
	}
	record := fmt.Sprintf(`{"schema":1,"state":"ready","status":"green","tree":%q,"oracle":%q,"recorded_at":%q,"check_executed":%s,"check_inherited":%s,"check_evidence":%s}`+"\n",
		cachedTree, strings.Repeat("0", 64), recorded, executedJSON, inheritedJSON, evidenceJSON)
	if err := os.WriteFile(filepath.Join(gitdir, git.GateCacheFile), []byte(record), 0o600); err != nil {
		t.Fatal(err)
	}
}

// withFiles copies base and adds each path with throwaway content, so a case names only
// the paths it varies.
func withFiles(base map[string]string, paths ...string) map[string]string {
	out := make(map[string]string, len(base)+len(paths))
	for path, content := range base {
		out[path] = content
	}
	for _, path := range paths {
		out[path] = "drift\n"
	}
	return out
}

// treeOf materializes files as the repository's whole content and returns the tree hash
// git diff compares. The work tree is emptied first: a leftover from an earlier tree would
// join the next one and change which paths the diff reports.
func treeOf(t *testing.T, root string, files map[string]string) string {
	t.Helper()
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.Name() == ".git" {
			continue
		}
		if err := os.RemoveAll(filepath.Join(root, e.Name())); err != nil {
			t.Fatal(err)
		}
	}
	for path, content := range files {
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	gitRun(t, root, "read-tree", "--empty")
	gitRun(t, root, "add", "-A")
	return gitRun(t, root, "write-tree")
}

func initRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	gitRun(t, root, "init")
	gitRun(t, root, "config", "user.email", "t@example.com")
	gitRun(t, root, "config", "user.name", "t")
	// The capture surfaces moved under a directory the repository root no longer
	// supplies for free, so a fixture root that omits it fails every write that used
	// to land beside ROADMAP.md.
	if err := os.MkdirAll(filepath.Join(root, "capture"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, ".bench"), 0o755); err != nil {
		t.Fatal(err)
	}
	return root
}

func gitRun(t *testing.T, root string, args ...string) string {
	t.Helper()
	out, err := git.Output(append([]string{"-C", root}, args...)...)
	if err != nil {
		t.Fatalf("git %v: %v (%s)", args, err, out)
	}
	return out
}

// A clean, committed tree with no signals renders the clean message and nothing else.
func TestRenderClean(t *testing.T) {
	root := initRepo(t)
	if err := os.WriteFile(filepath.Join(root, "f.txt"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")

	out := render(root, false)
	if out != "bench: clean — nothing pending\n" {
		t.Errorf("clean render = %q", out)
	}
}

// A dirty tree leads with the git action; the capture-drain row is present but outranked.
func TestRenderDirtyLeadsGitOverDrainRow(t *testing.T) {
	root := initRepo(t)
	if err := os.WriteFile(filepath.Join(root, "f.txt"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")
	// Make the tree dirty (uncommitted change) so the git signal fires.
	if err := os.WriteFile(filepath.Join(root, "f.txt"), []byte("y\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Park an idea so the drain signal fires.
	if err := os.WriteFile(filepath.Join(root, "capture/IDEAS.md"), []byte("- 2026-07-03  an idea\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	out := render(root, false)
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if !strings.HasPrefix(lines[0], "▶ /bench-final-check  (git)") {
		t.Errorf("lead line = %q, want git action lead", lines[0])
	}
	if !strings.Contains(out, "1 idea(s), 0 open learning(s)") || !strings.Contains(out, "/bench-what-next") {
		t.Errorf("drain row missing from:\n%s", out)
	}
}

// A working roadmap alone is not pending capture: no drain row, the board stays clean.
func TestRenderWorkingRoadmapAloneIsClean(t *testing.T) {
	root := initRepo(t)
	content := "# Roadmap\n\n## Recommended sequence\n\n1. Shape next item - /bench-shape-idea\n"
	if err := os.WriteFile(filepath.Join(root, "ROADMAP.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")

	if out := render(root, false); out != "bench: clean — nothing pending\n" {
		t.Errorf("roadmap-only render = %q, want clean board", out)
	}
}

func TestAppendWorktreeIgnoresUnownedBranchPrefix(t *testing.T) {
	root := initRepo(t)
	if err := os.WriteFile(filepath.Join(root, "f.txt"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")
	gitRun(t, root, "branch", "worktree-agent-orphan")

	if rows := appendWorktree(nil, root); len(rows) != 0 {
		t.Fatalf("branch prefix created worktree status without ownership evidence: %#v", rows)
	}
}

// appendWorktree must surface a discovery failure as a visible row, never as silence that
// a reader mistakes for "no worktree signals". This filesystem refusal reaches common-dir
// resolution before porcelain, while the PATH-stub fixtures cover typed and generic routing.
func TestAppendWorktreeSurfacesClassifyFailure(t *testing.T) {
	root := initRepo(t)
	if err := os.WriteFile(filepath.Join(root, "f.txt"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")

	gitDir := filepath.Join(root, ".git")
	if err := os.Chmod(gitDir, 0o000); err != nil {
		t.Fatalf("chmod .git unreadable: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(gitDir, 0o755) })

	rows := appendWorktree(nil, root)
	if len(rows) == 0 {
		t.Fatal("appendWorktree dropped the classify failure instead of surfacing a row")
	}
	if !strings.Contains(rows[0].detail, "git common directory") || rows[0].action != "bench worktree list" {
		t.Errorf("row = %#v, want typed resolution refusal", rows[0])
	}
}

func TestAppendWorktreeKeepsTypedAndPorcelainFailureActionsDistinct(t *testing.T) {
	for _, tc := range []struct {
		mode, detail, action string
	}{
		{"fail-rev-parse", "rev-parse", "bench worktree list"},
		{"fail-worktree", "git worktree list failed", "git worktree list"},
	} {
		t.Run(tc.mode, func(t *testing.T) {
			root := initRepo(t)
			gittest.StubGit(t, root, tc.mode, filepath.Join(t.TempDir(), "argv"))
			rows := appendWorktree(nil, root)
			if len(rows) != 1 || !strings.Contains(rows[0].detail, tc.detail) || rows[0].action != tc.action {
				t.Fatalf("%s row = %#v", tc.mode, rows)
			}
		})
	}
}

func TestAppendWorktreeRendersBoundExpiryAsTypedFailure(t *testing.T) {
	restore := git.SetWorktreeListTimeoutForTest(100 * time.Millisecond)
	t.Cleanup(restore)
	root := initRepo(t)
	gittest.StubGit(t, root, "block-worktree", filepath.Join(t.TempDir(), "argv"))
	rows := appendWorktree(nil, root)
	if len(rows) != 1 || !strings.Contains(rows[0].detail, "worktree list") || rows[0].action != "bench worktree list" {
		t.Fatalf("bound row = %#v", rows)
	}
}

func TestAppendWorktreeRendersTypedAdminRefusal(t *testing.T) {
	root := initRepo(t)
	gittest.FIFOWorktreeAdmin(t, root, "typed")
	rows := appendWorktree(nil, root)
	if len(rows) != 1 || !strings.Contains(rows[0].detail, "worktrees/typed/gitdir") || !strings.Contains(rows[0].detail, "fifo") || rows[0].action != "bench worktree list" {
		t.Fatalf("typed row = %#v", rows)
	}
}

// Command rejects an unknown argument with a usage line and exit 2, prints usage on -h,
// and accepts --all as the one added token — while --all plus junk and near-misses stay
// usage errors so a typo never silently prints the default board.
func TestCommandArgs(t *testing.T) {
	if r, c := Command([]string{"--bogus"}); c != 2 || !strings.Contains(r, "usage:") {
		t.Errorf("unknown arg: report %q exit %d", r, c)
	}
	if r, c := Command([]string{"-h"}); c != 0 || !strings.Contains(r, "usage: bench status") {
		t.Errorf("help: report %q exit %d", r, c)
	}
	if r, c := Command([]string{"-h"}); !strings.Contains(r, "[--all]") {
		t.Errorf("help usage should advertise [--all], got %q exit %d", r, c)
	}
	if r, c := Command([]string{"--all"}); c != 0 {
		t.Errorf("--all should be accepted with exit 0, got report %q exit %d", r, c)
	}
	for _, bad := range [][]string{{"--all", "extra"}, {"--allx"}, {"-a"}} {
		if r, c := Command(bad); c != 2 || !strings.Contains(r, "usage:") {
			t.Errorf("args %q: report %q exit %d, want usage exit 2", bad, r, c)
		}
	}
}

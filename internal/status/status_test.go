package status

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/roadmap"
)

func TestTimeoutGateIsDistinctHighestSeveritySignal(t *testing.T) {
	rows := appendGateInfo(nil, GateInfo{Present: true, State: "ready", Status: "timeout"}, t.TempDir())
	if len(rows) != 1 || rows[0].detail != "timeout" || rows[0].sev != 0 || !strings.Contains(rows[0].action, "hang") {
		t.Fatalf("timeout rows = %#v", rows)
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
	if action != "re-run the gate" {
		t.Fatalf("action = %q, want re-run the gate", action)
	}
}

// The board's capture-only softening and the gate's reduced-run declaration are one fact,
// so the samples come from the declaration rather than from a list restated here: a path
// the declaration carries softens, one it does not falls through to the strong stale row,
// and a mixed diff fails closed. A private allowlist answering the question inside this
// package fails the moment the declaration carries a path it does not.
func TestStaleSofteningRoutesThroughDeclaration(t *testing.T) {
	scope := gate.ReducedScope()
	root := initRepo(t)
	const outside = "internal/status/status.go"
	base := map[string]string{outside: "package status\n"}
	baseTree := treeOf(t, root, base)

	var declared []string
	declared = append(declared, scope.Files()...)
	for _, dir := range scope.Directories() {
		declared = append(declared, dir+"declared-descendant.md")
	}
	for _, path := range declared {
		detail, action := staleGateDetailAction(root, baseTree, treeOf(t, root, withFiles(base, path)))
		if detail != "stale (capture-only drift)" || action != "re-run when convenient" {
			t.Errorf("drift in declared %q = (%q, %q), want the softened row", path, detail, action)
		}
	}

	undeclared := treeOf(t, root, map[string]string{outside: "package status // drift\n"})
	if detail, _ := staleGateDetailAction(root, baseTree, undeclared); !strings.HasPrefix(detail, "stale (gated tree") {
		t.Errorf("drift in undeclared %q = %q, want the strong stale row", outside, detail)
	}

	mixed := treeOf(t, root, withFiles(map[string]string{outside: "package status // drift\n"}, declared...))
	if detail, _ := staleGateDetailAction(root, baseTree, mixed); !strings.HasPrefix(detail, "stale (gated tree") {
		t.Errorf("mixed drift = %q, want the strong stale row", detail)
	}
}

// A reduced verdict graded only the phases that could observe its changeset, so it is a
// narrow green rather than drift. The row names the narrowness and never reports the tree
// against itself: the two hashes the stale row prints would be the same hash here.
func TestReducedGreenRendersItsOwnRow(t *testing.T) {
	root := initRepo(t)
	tree := treeOf(t, root, map[string]string{"capture/IDEAS.md": "- an idea\n"})
	writeReducedGateCache(t, root, tree)

	gv := GateVerdict(root)
	if !gv.Reduced || !gv.Stale || gv.CachedTree != gv.WorkTree {
		t.Fatalf("verdict = %#v, want a reduced non-reusable green over the current tree", gv)
	}
	rows := appendGateInfo(nil, gv, root)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one gate row", rows)
	}
	if strings.Contains(rows[0].detail, "stale") {
		t.Errorf("detail = %q, want a reduced row rather than a stale one", rows[0].detail)
	}
	if strings.Contains(rows[0].detail, Short(tree)) {
		t.Errorf("detail = %q, want no tree hash: the gated and work trees are the same tree", rows[0].detail)
	}
}

// The action a reduced row names has to widen the verdict. Re-running the gate over the
// same capture-only changeset records another reduced verdict and the same row, so an
// action naming it would loop forever.
func TestReducedRowActionWidensTheVerdict(t *testing.T) {
	root := initRepo(t)
	writeReducedGateCache(t, root, treeOf(t, root, map[string]string{"capture/IDEAS.md": "- an idea\n"}))

	rows := appendGateInfo(nil, GateVerdict(root), root)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one gate row", rows)
	}
	if !strings.Contains(rows[0].action, "bench gate --fresh") {
		t.Errorf("action = %q, want the escape that widens the verdict", rows[0].action)
	}
}

// Reducedness and staleness are independent: a reduced verdict whose tree no longer
// matches the work tree is genuinely stale, and undeclared drift keeps the strong row.
func TestDriftedReducedVerdictStillRendersStaleRow(t *testing.T) {
	root := initRepo(t)
	const outside = "internal/status/status.go"
	gated := treeOf(t, root, map[string]string{outside: "package status\n"})
	current := treeOf(t, root, map[string]string{outside: "package status // drift\n"})
	writeReducedGateCache(t, root, gated)

	gv := GateVerdict(root)
	if !gv.Reduced || gv.WorkTree != current {
		t.Fatalf("verdict = %#v, want a reduced verdict over a drifted work tree", gv)
	}
	rows := appendGateInfo(nil, gv, root)
	if len(rows) != 1 {
		t.Fatalf("rows = %#v, want one gate row", rows)
	}
	wantDetail := fmt.Sprintf("stale (gated tree %s, work tree %s)", Short(gated), Short(current))
	if rows[0].detail != wantDetail || rows[0].action != "re-run the gate" {
		t.Errorf("drifted reduced row = (%q, %q), want (%q, \"re-run the gate\")", rows[0].detail, rows[0].action, wantDetail)
	}
}

// writeReducedGateCache installs a ready reduced green naming cachedTree, at the mode and
// in the exact field set the gate's loader requires of the reduced class. The phase list
// comes from the declaration rather than a list restated here.
func writeReducedGateCache(t *testing.T, root, cachedTree string) {
	t.Helper()
	gitdir := gitRun(t, root, "rev-parse", "--absolute-git-dir")
	phases, err := json.Marshal(gate.ReducedScope().IncludedPhases())
	if err != nil {
		t.Fatal(err)
	}
	recorded := time.Now().UTC().Truncate(time.Second).Add(-time.Minute).Format(time.RFC3339)
	record := fmt.Sprintf(`{"schema":1,"state":"ready","status":"green","tree":%q,"oracle":%q,"recorded_at":%q,"reduced":true,"phases":%s,"ancestor":%q,"ancestor_recorded_at":%q}`+"\n",
		cachedTree, strings.Repeat("0", 64), recorded, phases, strings.Repeat("a", 40), recorded)
	path := filepath.Join(gitdir, git.GateCacheFile)
	if err := os.WriteFile(path, []byte(record), 0o600); err != nil {
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
	if !strings.HasPrefix(lines[0], "▶ commit on green  (git)") {
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

// appendWorktree must surface a `git worktree list` failure as a visible row, never as
// silence that a reader mistakes for "no worktree signals" — the false-empty class FT29
// swept. The failure is induced deterministically (the FT29 gitOpError style: break the
// git query itself, here by revoking read access to .git) rather than a PATH-shimmed git.
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
	if !strings.Contains(rows[0].detail, "worktree list failed") {
		t.Errorf("row detail = %q, want it to name the git worktree-list failure", rows[0].detail)
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

// Tests for gate-cache verdict reading, partial verdicts, and tree-drift staleness.
package status

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/git"
)

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
	if action.render() != "bench gate" {
		t.Fatalf("action = %q, want bench gate", action.render())
	}
}

// The reduced verdict class is retired. A legacy on-disk reduced record must read as an
// invalid cache rather than as a green of any width. The board reports it as invalid,
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
	if !strings.Contains(rows[0].action.render(), "bench gate --fresh") {
		t.Errorf("action = %q, want the fresh whole-tree action", rows[0].action.render())
	}
}

// A red recorded against a tree the work tree has since left describes that run, not this
// one. The board must send the reader back to the gate rather than headline a red for
// work that is no longer in the tree. The drifted record is stale, whatever verdict it
// carries.
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
	if !strings.HasPrefix(rows[0].detail, "stale (gated tree") || rows[0].action.render() != "bench gate" {
		t.Fatalf("rows = %#v, want the drift row rather than a red one", rows)
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

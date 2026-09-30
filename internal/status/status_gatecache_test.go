// Tests for gate-cache verdict reading, retired record classes, and tree-drift staleness.
package status

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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
	detail, action := staleGateDetailAction(t.TempDir(), "0123456789abcdef", "none", "")
	if detail != "stale (gated tree 0123456, work tree none)" {
		t.Fatalf("detail = %q, want strong stale detail", detail)
	}
	if action.render() != "bench gate" {
		t.Fatalf("action = %q, want bench gate", action.render())
	}
}

// The gate names why it no longer stands behind the record. Without that reason a reader
// whose two trees match sees a staleness with no cause and no way to judge the rerun.
func TestStaleGateDetailActionNamesTheGateReason(t *testing.T) {
	detail, _ := staleGateDetailAction(t.TempDir(), "0123456789abcdef", "0123456789abcdef", "evidence changed")
	if detail != "stale (gated tree 0123456, work tree 0123456; evidence changed)" {
		t.Fatalf("detail = %q, want the stale detail naming the gate reason", detail)
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

// The partial verdict classes are retired too. A legacy on-disk partial record must read
// as an invalid cache, exactly as a reduced one does, and the board sends the reader to
// the gate rather than naming components nothing can validate.
func TestLegacyPartialCacheReadsAsInvalid(t *testing.T) {
	root := initRepo(t)
	tree := treeOf(t, root, map[string]string{"f.txt": "x\n"})
	writeLegacyPartialGateCache(t, root, tree, "docs")

	gv := GateVerdict(root)
	if !gv.Present || gv.State != "invalid" || gv.Status == "green" {
		t.Fatalf("verdict = %#v, want a legacy partial record read as an invalid cache", gv)
	}
	rows := appendGateInfo(nil, gv, root)
	if len(rows) != 1 || rows[0].action.render() != "bench gate" || strings.Contains(rows[0].detail, "docs") {
		t.Fatalf("rows = %#v, want one invalid-cache row that sends the reader to bench gate", rows)
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

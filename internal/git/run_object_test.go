package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
)

// objectRepo is a repository with two commits on its default branch. The first commit
// holds a.txt, and the second commit is empty.
type objectRepo struct {
	root, commit, second, tree, blob string
}

func newObjectRepo(t *testing.T) objectRepo {
	t.Helper()
	root := newRepo(t)
	rev := func(name string) string { return strings.TrimSpace(runGit(t, root, "rev-parse", name)) }
	commit := rev("HEAD")
	runGit(t, root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "--allow-empty", "-m", "second")
	return objectRepo{root: root, commit: commit, second: rev("HEAD"), tree: rev(commit + "^{tree}"), blob: rev(commit + ":a.txt")}
}

// listing renders tree entries as one comparable string.
func listing(entries []TreeEntry, err error) (string, error) {
	lines := make([]string, len(entries))
	for i, entry := range entries {
		lines[i] = fmt.Sprintf("%s %s %s\t%s", entry.Mode, entry.Type, entry.Object, entry.Path)
	}
	return strings.Join(lines, "\n"), err
}

func requireStarts(t *testing.T, logPath string, want map[string]int) {
	t.Helper()
	for needle, count := range want {
		if got := gittest.StubGitStarts(t, logPath, needle); got != count {
			t.Fatalf("starts of %q = %d, want %d", needle, got, count)
		}
	}
}

// TestRunObjectReads proves that one run starts each pure object read once, returns a copy
// of each stored answer, and keys a resolution by its flags. Two run values share nothing.
// (Coverage rows GQ1, GQ3, GQ6, GQ7, GQ8, GQ10, GQ13, GQ14, GQ18, and GQ21.)
func TestRunObjectReads(t *testing.T) {
	f := newObjectRepo(t)
	blobRead := "cat-file blob " + f.blob
	cases := []struct {
		name   string
		read   func(run *Run, call int) (string, error)
		want   string
		starts map[string]int
	}{
		{"blob", func(run *Run, _ int) (string, error) {
			data, err := run.ReadControlBlob(f.root, f.blob)
			if err != nil || len(data) == 0 {
				return string(data), err
			}
			data[0] = 'X'
			again, err := run.ReadControlBlob(f.root, f.blob)
			return string(again), err
		}, "hello\n", map[string]int{blobRead: 1}},
		{"tree path", func(run *Run, _ int) (string, error) {
			return listing(run.ListTree(f.root, f.tree, "--", "a.txt"))
		}, "100644 blob " + f.blob + "\ta.txt", map[string]int{"ls-tree -z " + f.tree: 1}},
		{"absent path", func(run *Run, _ int) (string, error) {
			return listing(run.ListTree(f.root, f.tree, "--", "absent.txt"))
		}, "", map[string]int{"ls-tree -z " + f.tree: 1}},
		{"tree file", func(run *Run, _ int) (string, error) {
			data, err := run.ReadTreeFile(f.root, f.tree, "a.txt")
			return string(data), err
		}, "hello\n", map[string]int{"ls-tree -z " + f.tree: 1, blobRead: 1}},
		{"commit", func(run *Run, _ int) (string, error) {
			return run.ResolveCommit(f.root, f.commit, "--quiet")
		}, f.commit, map[string]int{"rev-parse --verify --quiet " + f.commit: 1}},
		{"flag sets", func(run *Run, call int) (string, error) {
			if call == 0 {
				return run.ResolveCommit(f.root, f.commit, "--quiet")
			}
			return run.ResolveCommit(f.root, f.commit)
		}, f.commit, map[string]int{f.commit + "^{commit}": 2}},
		{"ancestor", func(run *Run, _ int) (string, error) {
			yes, err := run.IsAncestor(f.root, f.commit, f.second)
			return strconv.FormatBool(yes), err
		}, "true", map[string]int{"--is-ancestor " + f.commit: 1}},
		{"non-ancestor", func(run *Run, _ int) (string, error) {
			yes, err := run.IsAncestor(f.root, f.second, f.commit)
			return strconv.FormatBool(yes), err
		}, "false", map[string]int{"--is-ancestor " + f.second: 1}},
		{"two runs", func(*Run, int) (string, error) {
			data, err := new(Run).ReadControlBlob(f.root, f.blob)
			return string(data), err
		}, "hello\n", map[string]int{blobRead: 2}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, logPath := loggedStub(t, f.root, "log")
			run := new(Run)
			for call := range 2 {
				if got, err := c.read(run, call); err != nil || got != c.want {
					t.Fatalf("read %d = %q, %v, want %q, nil", call+1, got, err, c.want)
				}
			}
			requireStarts(t, logPath, c.starts)
		})
	}
}

// TestRunObjectReadsStayLive proves that one run reads a ref, HEAD, or short-id operand
// live, and that its primary-checkout read starts the porcelain on every call.
// (Coverage rows GQ2, GQ4, GQ5, GQ9, GQ15, and GQ24.)
func TestRunObjectReadsStayLive(t *testing.T) {
	f := newObjectRepo(t)
	def, ok := ResolvedDefault(f.root)
	if !ok {
		t.Fatalf("ResolvedDefault(%s) found no default branch", f.root)
	}
	runGit(t, f.root, "branch", "side", f.commit)
	_, logPath := loggedStub(t, f.root, "log")
	run := new(Run)
	for range 2 {
		if got, err := run.PrimaryCheckout(f.root); err != nil || got != f.root {
			t.Fatalf("run.PrimaryCheckout = %q, %v, want %q", got, err, f.root)
		}
	}
	requireStarts(t, logPath, map[string]int{"--git-common-dir": 1, "worktree list": 2})

	reads := func() string {
		t.Helper()
		blob, err := run.ReadControlBlob(f.root, "HEAD:a.txt")
		if err != nil {
			t.Fatal(err)
		}
		tree, err := listing(run.ListTree(f.root, "HEAD"))
		if err != nil {
			t.Fatal(err)
		}
		tip, err := run.ResolveCommit(f.root, "refs/heads/"+def)
		if err != nil {
			t.Fatal(err)
		}
		ancestor, err := run.IsAncestor(f.root, f.second, "refs/heads/side")
		if err != nil {
			t.Fatal(err)
		}
		return strings.Join([]string{string(blob), tree, tip, strconv.FormatBool(ancestor)}, "|")
	}
	before := reads()
	if err := os.WriteFile(filepath.Join(f.root, "a.txt"), []byte("next\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, f.root, "-c", "user.email=bench@local", "-c", "user.name=bench", "commit", "-q", "-am", "next")
	runGit(t, f.root, "update-ref", "refs/heads/side", f.second)
	head := strings.TrimSpace(runGit(t, f.root, "rev-parse", "HEAD"))
	nextTree, err := listing(ListTree(f.root, head))
	if err != nil {
		t.Fatal(err)
	}
	if want, got := strings.Join([]string{"next\n", nextTree, head, "true"}, "|"), reads(); got != want || got == before {
		t.Fatalf("reads after the commit and the update-ref = %q, want %q, apart from %q", got, want, before)
	}

	short := f.tree[:7]
	for range 2 {
		if _, err := run.ListTree(f.root, short); err != nil {
			t.Fatalf("run.ListTree(%s): %v", short, err)
		}
	}
	requireStarts(t, logPath, map[string]int{"ls-tree -z " + short: 2})
}

// TestRunObjectFailures proves that a failed object read stores nothing, so each read of a
// missing full id starts git again. (Coverage rows GQ16 and GQ17.)
func TestRunObjectFailures(t *testing.T) {
	f := newObjectRepo(t)
	missing := strings.Repeat("1", 40)
	_, logPath := loggedStub(t, f.root, "log")
	run := new(Run)
	for call := range 2 {
		if _, err := run.IsAncestor(f.root, missing, f.commit); err == nil {
			t.Fatalf("check %d of the missing ancestor succeeded, want an error", call+1)
		}
		if _, err := run.ReadControlBlob(f.root, missing); err == nil {
			t.Fatalf("read %d of the missing blob succeeded, want an error", call+1)
		}
	}
	requireStarts(t, logPath, map[string]int{"--is-ancestor " + missing: 2, "cat-file blob " + missing: 2})
}

// TestRunNilCachesNothing proves that each cached form on a nil run, and each package
// read, starts its query on every call. (Coverage rows GQ19 and GQ20.)
func TestRunNilCachesNothing(t *testing.T) {
	f := newObjectRepo(t)
	var none *Run
	cases := []struct {
		name, needle string
		call         func() error
	}{
		{"CommonDir", "--git-common-dir", func() error { _, err := none.CommonDir(f.root); return err }},
		{"AdminDir", "--git-dir", func() error { _, err := none.AdminDir(f.root); return err }},
		{"RootAt", "--show-toplevel", func() error { _, err := none.RootAt(f.root); return err }},
		{"AdminPath", "--git-path index", func() error { _, err := none.AdminPath(f.root, "index"); return err }},
		{"PrimaryCheckout", "--git-common-dir", func() error { _, err := none.PrimaryCheckout(f.root); return err }},
		{"ReadControlBlob", "cat-file blob " + f.blob, func() error { _, err := none.ReadControlBlob(f.root, f.blob); return err }},
		{"ListTree", "ls-tree -z " + f.tree, func() error { _, err := none.ListTree(f.root, f.tree); return err }},
		{"ReadTreeFile", "cat-file blob " + f.blob, func() error { _, err := none.ReadTreeFile(f.root, f.tree, "a.txt"); return err }},
		{"ResolveCommit", "rev-parse --verify " + f.commit, func() error { _, err := none.ResolveCommit(f.root, f.commit); return err }},
		{"IsAncestor", "--is-ancestor " + f.commit, func() error { _, err := none.IsAncestor(f.root, f.commit, f.second); return err }},
		{"package ReadControlBlob", "cat-file blob " + f.blob, func() error { _, err := ReadControlBlob(f.root, f.blob); return err }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, logPath := loggedStub(t, f.root, "log")
			for call := range 2 {
				if err := c.call(); err != nil {
					t.Fatalf("call %d: %v", call+1, err)
				}
			}
			requireStarts(t, logPath, map[string]int{c.needle: 2})
		})
	}
}

// TestFullObjectID proves that the purity predicate accepts 40 and 64 hexadecimal digits of
// either case, and refuses every other length, a non-hex byte, and the empty value.
// (Coverage rows GQ25 and GQ47.)
func TestFullObjectID(t *testing.T) {
	sha1, sha256 := strings.Repeat("a0", 20), strings.Repeat("f9", 32)
	cases := map[string]bool{
		sha1:                    true,
		strings.ToUpper(sha1):   true,
		sha256:                  true,
		strings.ToUpper(sha256): true,
		sha1[:39]:               false,
		sha1 + "a":              false,
		sha256[:63]:             false,
		sha256 + "a":            false,
		sha1[:39] + "g":         false,
		"":                      false,
	}
	for value, want := range cases {
		if got := FullObjectID(value); got != want {
			t.Fatalf("FullObjectID(%q) = %v, want %v", value, got, want)
		}
	}
}

// TestRunAdminPathRevalidates proves that one run starts the administration-path query
// once for a root and a name, and that a read after the removal of the answer's parent
// directory refuses with no new start. (Coverage rows GQ22 and GQ23.)
func TestRunAdminPathRevalidates(t *testing.T) {
	root := newRepo(t)
	parent := filepath.Join(root, ".git", "bench-scratch")
	if err := os.Mkdir(parent, 0o755); err != nil {
		t.Fatal(err)
	}
	_, logPath := loggedStub(t, root, "log")
	run := new(Run)
	first, err := run.AdminPath(root, "bench-scratch/file")
	if err != nil {
		t.Fatalf("first read: %v", err)
	}
	if second, err := run.AdminPath(root, "bench-scratch/file"); err != nil || second != first {
		t.Fatalf("second read = %q, %v, want %q, nil", second, err, first)
	}
	if err := os.Remove(parent); err != nil {
		t.Fatal(err)
	}
	_, err = run.AdminPath(root, "bench-scratch/file")
	var resolution *ResolutionError
	if !errors.As(err, &resolution) || resolution.Path != filepath.Dir(first) {
		t.Fatalf("read after the removal = %v, want the refusal of %q", err, filepath.Dir(first))
	}
	requireStarts(t, logPath, map[string]int{"--git-path": 1})
}

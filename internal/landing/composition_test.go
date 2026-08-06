package landing

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/gate/authorization"
)

type compositionExpectation struct {
	paths        []string
	wantNames    []string
	wantContent  map[string]string
	wantMode     map[string]string
	wantIndex    map[string]string
	wantWorktree map[string]string
}

func TestLandRealGitCompositionTable(t *testing.T) {
	cases := []struct {
		name  string
		setup func(*testing.T, string) compositionExpectation
	}{
		{"addition", func(t *testing.T, root string) compositionExpectation {
			write(t, root, "added", "addition")
			return compositionResult([]string{"added"}, []string{"added", "foreign", "named"}, map[string]string{"added": "addition", "foreign": "base", "named": "base"})
		}},
		{"modification", func(t *testing.T, root string) compositionExpectation {
			write(t, root, "named", "modified")
			return compositionResult([]string{"named"}, []string{"foreign", "named"}, map[string]string{"foreign": "base", "named": "modified"})
		}},
		{"deletion", func(t *testing.T, root string) compositionExpectation {
			if err := os.Remove(filepath.Join(root, "named")); err != nil {
				t.Fatal(err)
			}
			return compositionResult([]string{"named"}, []string{"foreign"}, map[string]string{"foreign": "base"})
		}},
		{"staged-deletion", func(t *testing.T, root string) compositionExpectation {
			git(t, root, "rm", "-q", "--", "named")
			return compositionResult([]string{"named"}, []string{"foreign"}, map[string]string{"foreign": "base"})
		}},
		{"both-rename-halves", func(t *testing.T, root string) compositionExpectation {
			if err := os.Rename(filepath.Join(root, "named"), filepath.Join(root, "renamed")); err != nil {
				t.Fatal(err)
			}
			git(t, root, "add", "-A", "--", "named", "renamed")
			return compositionResult([]string{"named", "renamed"}, []string{"foreign", "renamed"}, map[string]string{"foreign": "base", "renamed": "base"})
		}},
		{"named-directory-descendants", func(t *testing.T, root string) compositionExpectation {
			write(t, root, "owned/a", "a-base")
			write(t, root, "owned/b", "b-base")
			git(t, root, "add", "--", "owned")
			git(t, root, "commit", "-qm", "directory base")
			write(t, root, "owned/a", "a-modified")
			write(t, root, "owned/c", "c-added")
			if err := os.Remove(filepath.Join(root, "owned", "b")); err != nil {
				t.Fatal(err)
			}
			return compositionResult([]string{"owned"}, []string{"foreign", "named", "owned/a", "owned/c"}, map[string]string{"foreign": "base", "named": "base", "owned/a": "a-modified", "owned/c": "c-added"})
		}},
		{"deleted-directory-descendants", func(t *testing.T, root string) compositionExpectation {
			write(t, root, "gone/a", "a")
			write(t, root, "gone/b", "b")
			git(t, root, "add", "--", "gone")
			git(t, root, "commit", "-qm", "deleted directory base")
			if err := os.RemoveAll(filepath.Join(root, "gone")); err != nil {
				t.Fatal(err)
			}
			return compositionResult([]string{"gone"}, []string{"foreign", "named"}, map[string]string{"foreign": "base", "named": "base"})
		}},
		{"duplicate-paths", func(t *testing.T, root string) compositionExpectation {
			write(t, root, "named", "once")
			return compositionResult([]string{"named", "named", "named"}, []string{"foreign", "named"}, map[string]string{"foreign": "base", "named": "once"})
		}},
		{"literal-space", func(t *testing.T, root string) compositionExpectation {
			write(t, root, "space name", "space")
			return compositionResult([]string{"space name"}, []string{"foreign", "named", "space name"}, map[string]string{"foreign": "base", "named": "base", "space name": "space"})
		}},
		{"literal-glob-characters", func(t *testing.T, root string) compositionExpectation {
			write(t, root, "literal[1]*?", "glob")
			return compositionResult([]string{"literal[1]*?"}, []string{"foreign", "literal[1]*?", "named"}, map[string]string{"foreign": "base", "literal[1]*?": "glob", "named": "base"})
		}},
		{"dash-led-name-after-separator", func(t *testing.T, root string) compositionExpectation {
			write(t, root, "-dash", "dash")
			return compositionResult([]string{"-dash"}, []string{"-dash", "foreign", "named"}, map[string]string{"-dash": "dash", "foreign": "base", "named": "base"})
		}},
		{"symlink-entry-without-target-traversal", func(t *testing.T, root string) compositionExpectation {
			write(t, root, "symlink-target/foreign", "must stay untracked")
			if err := os.Symlink("symlink-target", filepath.Join(root, "link")); err != nil {
				capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable: %v", err))
			}
			got := compositionResult([]string{"link"}, []string{"foreign", "link", "named"}, map[string]string{"foreign": "base", "link": "symlink-target", "named": "base"})
			got.wantMode = map[string]string{"link": "120000"}
			got.wantWorktree = map[string]string{"symlink-target/foreign": "must stay untracked"}
			return got
		}},
		{"gitlink-atomic-without-nested-dirt-traversal", func(t *testing.T, root string) compositionExpectation {
			nested := filepath.Join(root, "nested")
			if err := os.Mkdir(nested, 0o755); err != nil {
				t.Fatal(err)
			}
			git(t, nested, "init", "-q", "-b", "main")
			git(t, nested, "config", "user.email", "a@b.c")
			git(t, nested, "config", "user.name", "a")
			write(t, nested, "inside", "nested-base")
			git(t, nested, "add", "inside")
			git(t, nested, "commit", "-qm", "nested base")
			git(t, root, "add", "--", "nested")
			git(t, root, "commit", "-qm", "gitlink base")
			write(t, nested, "inside", "nested-next")
			git(t, nested, "add", "inside")
			git(t, nested, "commit", "-qm", "nested next")
			write(t, nested, "inside", "nested dirt")
			if err := syscall.Mkfifo(filepath.Join(nested, "nested-fifo"), 0o600); err != nil {
				capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable: %v", err))
			}
			got := compositionResult([]string{"nested"}, []string{"foreign", "named", "nested"}, map[string]string{"foreign": "base", "named": "base"})
			got.wantMode = map[string]string{"nested": "160000"}
			got.wantIndex = map[string]string{"nested": git(t, nested, "rev-parse", "HEAD")}
			got.wantWorktree = map[string]string{"nested/inside": "nested dirt"}
			return got
		}},
		{"whole-segment-prefix-sibling-exclusion", func(t *testing.T, root string) compositionExpectation {
			write(t, root, "scope/a", "scope-base")
			write(t, root, "scope-sibling/b", "sibling-base")
			git(t, root, "add", "--", "scope", "scope-sibling")
			git(t, root, "commit", "-qm", "prefix base")
			write(t, root, "scope/a", "scope-changed")
			write(t, root, "scope-sibling/b", "sibling-index")
			git(t, root, "add", "--", "scope-sibling/b")
			write(t, root, "scope-sibling/b", "sibling-worktree")
			got := compositionResult([]string{"scope"}, []string{"foreign", "named", "scope-sibling/b", "scope/a"}, map[string]string{"foreign": "base", "named": "base", "scope-sibling/b": "sibling-base", "scope/a": "scope-changed"})
			got.wantIndex = map[string]string{"scope-sibling/b": "sibling-index"}
			got.wantWorktree = map[string]string{"scope-sibling/b": "sibling-worktree"}
			return got
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := fixture(t)
			want := tc.setup(t, root)
			base := git(t, root, "rev-parse", "HEAD")
			write(t, root, "foreign", "foreign-index")
			git(t, root, "add", "--", "foreign")
			write(t, root, "foreign", "foreign-worktree")
			o := greenOwner()
			got, err := o.Land(context.Background(), Request{Root: root, Destination: "refs/heads/main", Expected: base, Message: tc.name, Paths: want.paths})
			if err != nil {
				t.Fatal(err)
			}
			assertPublishedTree(t, root, got, want.wantNames, want.wantContent, want.wantMode)
			if got := git(t, root, "show", ":foreign"); got != "foreign-index" {
				t.Fatalf("foreign index = %q, want foreign-index", got)
			}
			if got := string(mustRead(t, filepath.Join(root, "foreign"))); got != "foreign-worktree" {
				t.Fatalf("foreign worktree = %q, want foreign-worktree", got)
			}
			for path, content := range want.wantIndex {
				if mode := want.wantMode[path]; mode == "160000" {
					if got := strings.Fields(git(t, root, "ls-tree", got.Commit, "--", path))[2]; got != content {
						t.Fatalf("published gitlink %s = %q, want %q", path, got, content)
					}
					continue
				}
				if got := git(t, root, "show", ":"+path); got != content {
					t.Fatalf("index %s = %q, want %q", path, got, content)
				}
			}
			for path, content := range want.wantWorktree {
				if got := string(mustRead(t, filepath.Join(root, filepath.FromSlash(path)))); got != content {
					t.Fatalf("worktree %s = %q, want %q", path, got, content)
				}
			}
		})
	}
}

func compositionResult(paths, names []string, content map[string]string) compositionExpectation {
	return compositionExpectation{paths: paths, wantNames: names, wantContent: content}
}

func greenOwner() Owner {
	o := New()
	o.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
		return authorization.Result{Kind: authorization.Green}
	}
	return o
}

func assertPublishedTree(t *testing.T, root string, got Result, wantNames []string, wantContent, wantMode map[string]string) {
	t.Helper()
	if got.Tree != git(t, root, "rev-parse", got.Commit+"^{tree}") || git(t, root, "rev-parse", "HEAD") != got.Commit {
		t.Fatal("authorized tree was not the exact published tree")
	}
	rawNames := git(t, root, "ls-tree", "-rz", "--name-only", got.Commit)
	names := strings.Split(rawNames, "\x00")
	if names[len(names)-1] == "" {
		names = names[:len(names)-1]
	}
	sort.Strings(wantNames)
	if !reflect.DeepEqual(names, wantNames) {
		t.Fatalf("published paths = %q, want %q", names, wantNames)
	}
	for path, content := range wantContent {
		if got := git(t, root, "show", got.Commit+":"+path); got != content {
			t.Fatalf("published %s = %q, want %q", path, got, content)
		}
	}
	for path, mode := range wantMode {
		fields := strings.Fields(git(t, root, "ls-tree", got.Commit, "--", path))
		if len(fields) < 4 || fields[0] != mode {
			t.Fatalf("published mode %s = %q, want %s", path, fields, mode)
		}
	}
}

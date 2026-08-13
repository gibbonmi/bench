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

func TestComposeDivergenceProbe(t *testing.T) {
	root := fixture(t)
	base := git(t, root, "rev-parse", "HEAD")
	write(t, root, "destination-only", "destination")
	git(t, root, "add", "destination-only")
	git(t, root, "commit", "-qm", "destination")
	destination := git(t, root, "rev-parse", "HEAD")
	git(t, root, "checkout", "-qb", "source", base)
	write(t, root, "source-only", "source")
	git(t, root, "add", "source-only")
	git(t, root, "commit", "-qm", "source")
	source := git(t, root, "rev-parse", "HEAD")
	got, err := New().Compose(CompositionRequest{Root: root, Destination: destination, Source: source, ReviewBase: base})
	if err != nil || got.Base != base || got.Conflict.Kind != "" {
		t.Fatalf("compose = %+v, %v", got, err)
	}
	if git(t, root, "show", got.Tree+":destination-only") != "destination" || git(t, root, "show", got.Tree+":source-only") != "source" {
		t.Fatal("tree omitted side")
	}
}

func TestComposePartialAncestryAppliesOnlyLaterSourceCommit(t *testing.T) {
	root := fixture(t)
	base := git(t, root, "rev-parse", "HEAD")
	git(t, root, "checkout", "-qb", "source", base)
	write(t, root, "early", "early")
	git(t, root, "add", "early")
	git(t, root, "commit", "-qm", "early")
	early := git(t, root, "rev-parse", "HEAD")
	git(t, root, "checkout", "-q", "main")
	git(t, root, "merge", "-q", "--ff-only", early)
	destination := git(t, root, "rev-parse", "HEAD")
	git(t, root, "checkout", "-q", "source")
	write(t, root, "later", "later")
	git(t, root, "add", "later")
	git(t, root, "commit", "-qm", "later")
	source := git(t, root, "rev-parse", "HEAD")
	got, err := New().Compose(CompositionRequest{Root: root, Destination: destination, Source: source, ReviewBase: base})
	if err != nil || got.Base != early || got.Conflict.Kind != "" {
		t.Fatalf("compose = %+v, %v", got, err)
	}
	if git(t, root, "show", got.Tree+":early") != "early" || git(t, root, "show", got.Tree+":later") != "later" {
		t.Fatal("partial ancestry tree incomplete")
	}
}

func TestComposeClassifiesRealGitConflictsWithoutMutation(t *testing.T) {
	cases := []struct {
		name  string
		kind  string
		setup func(*testing.T, string, string) (string, string)
	}{
		{"textual", "textual", changeBoth("named", "destination", "source")},
		{"modify-delete", "modify/delete", modifyDelete},
		{"rename-rename", "rename/rename", renameRename},
		{"file-directory", "file/directory", fileDirectory},
		{"mode", "mode", modeConflict},
		{"symlink", "symlink", symlinkConflict},
		{"gitlink", "gitlink", gitlinkConflict},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := fixture(t)
			base := git(t, root, "rev-parse", "HEAD")
			destination, source := tc.setup(t, root, base)
			before := compositionState(t, root)
			got, err := New().Compose(CompositionRequest{Root: root, Destination: destination, Source: source, ReviewBase: "not-the-merge-base"})
			if err != nil {
				t.Fatal(err)
			}
			if expectedBase := git(t, root, "merge-base", destination, source); got.Base != expectedBase || got.Tree != "" || got.Conflict.Kind != tc.kind {
				t.Fatalf("compose = %+v, want base %s and %s conflict", got, expectedBase, tc.kind)
			}
			if after := compositionState(t, root); !reflect.DeepEqual(after, before) {
				t.Fatalf("prospective composition mutated repository\nbefore: %#v\nafter: %#v", before, after)
			}
		})
	}
}

func TestConflictKindRejectsEmptyMergeTreeOutput(t *testing.T) {
	if _, err := conflictKind(""); err == nil {
		t.Fatal("empty merge-tree output classified without an error")
	}
}

func TestComposeRejectsUnresolvedCommit(t *testing.T) {
	root := fixture(t)
	_, err := New().Compose(CompositionRequest{Root: root, Destination: "not-a-commit", Source: "HEAD"})
	if err == nil || !strings.Contains(err.Error(), "destination is not a commit") {
		t.Fatalf("compose error = %v", err)
	}
}

func changeBoth(path, destinationValue, sourceValue string) func(*testing.T, string, string) (string, string) {
	return func(t *testing.T, root, base string) (string, string) {
		return commitSides(t, root, base, func() { write(t, root, path, destinationValue) }, func() { write(t, root, path, sourceValue) })
	}
}

func modifyDelete(t *testing.T, root, base string) (string, string) {
	return commitSides(t, root, base, func() { write(t, root, "named", "destination") }, func() { git(t, root, "rm", "-q", "named") })
}

func renameRename(t *testing.T, root, base string) (string, string) {
	return commitSides(t, root, base, func() { git(t, root, "mv", "named", "destination-name") }, func() { git(t, root, "mv", "named", "source-name") })
}

func fileDirectory(t *testing.T, root, base string) (string, string) {
	return commitSides(t, root, base, func() { write(t, root, "clash", "file") }, func() { write(t, root, "clash/child", "child") })
}

func modeConflict(t *testing.T, root, base string) (string, string) {
	return commitSides(t, root, base, func() {
		write(t, root, "named", "destination")
		if err := os.Chmod(filepath.Join(root, "named"), 0o755); err != nil {
			t.Fatal(err)
		}
	}, func() { write(t, root, "named", "source") })
}

func symlinkConflict(t *testing.T, root, base string) (string, string) {
	return commitSides(t, root, base, func() {
		if err := os.Remove(filepath.Join(root, "named")); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink("destination-target", filepath.Join(root, "named")); err != nil {
			capability.Capability(t, capability.Symlink, err.Error())
		}
	}, func() { write(t, root, "named", "source") })
}

func gitlinkConflict(t *testing.T, root, base string) (string, string) {
	nested := filepath.Join(root, "nested")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	git(t, nested, "init", "-q", "-b", "main")
	git(t, nested, "config", "user.email", "a@b.c")
	git(t, nested, "config", "user.name", "a")
	write(t, nested, "inside", "base")
	git(t, nested, "add", "inside")
	git(t, nested, "commit", "-qm", "base")
	git(t, root, "add", "nested")
	git(t, root, "commit", "-qm", "add gitlink")
	base = git(t, root, "rev-parse", "HEAD")
	git(t, nested, "checkout", "-qb", "destination")
	write(t, nested, "inside", "destination")
	git(t, nested, "commit", "-am", "destination", "-q")
	destinationNested := git(t, nested, "rev-parse", "HEAD")
	git(t, nested, "checkout", "-qb", "source", baseForNested(t, nested))
	write(t, nested, "inside", "source")
	git(t, nested, "commit", "-am", "source", "-q")
	sourceNested := git(t, nested, "rev-parse", "HEAD")
	return commitSides(t, root, base, func() { git(t, nested, "checkout", "-q", destinationNested); git(t, root, "add", "nested") }, func() { git(t, nested, "checkout", "-q", sourceNested); git(t, root, "add", "nested") })
}

func baseForNested(t *testing.T, root string) string { return git(t, root, "rev-parse", "main") }

func commitSides(t *testing.T, root, base string, destinationChange, sourceChange func()) (string, string) {
	git(t, root, "checkout", "-qb", "destination", base)
	destinationChange()
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "destination")
	destination := git(t, root, "rev-parse", "HEAD")
	git(t, root, "checkout", "-qb", "source", base)
	sourceChange()
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "source")
	return destination, git(t, root, "rev-parse", "HEAD")
}

type compositionSnapshot struct{ refs, index, status, worktree, mergeHead string }

func compositionState(t *testing.T, root string) compositionSnapshot {
	t.Helper()
	mergeHead, err := os.ReadFile(filepath.Join(root, ".git", "MERGE_HEAD"))
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	return compositionSnapshot{
		refs:      string(gitBytes(t, root, "for-each-ref", "--format=%(refname) %(objectname)")),
		index:     string(gitBytes(t, root, "ls-files", "--stage", "-z")),
		status:    string(gitBytes(t, root, "status", "--porcelain=v2", "--untracked-files=all")),
		worktree:  string(gitBytes(t, root, "diff", "--binary", "HEAD")),
		mergeHead: string(mergeHead),
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

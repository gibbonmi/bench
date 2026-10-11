package landing

import (
	"context"
	"errors"
	"io"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate/authorization"
)

// mergeNameSetRefusal opens the refusal of a name set that omits a merged path.
const mergeNameSetRefusal = "pending merge changed paths that the commit does not name: "

// assertParents fails t unless commit records exactly the parents want, in order.
func assertParents(t *testing.T, root, commit string, want ...string) {
	t.Helper()
	line := strings.Join(append([]string{commit}, want...), " ")
	if got := git(t, root, "rev-list", "--parents", "-n", "1", commit); got != line {
		t.Fatalf("commit and parents = %q, want %q", got, line)
	}
}

// TestLandRefusesAnUnreadableMergeState proves that a directory at MERGE_HEAD refuses with
// a plain error before publication, although the name set covers every merged path.
// (Coverage row MP24.)
func TestLandRefusesAnUnreadableMergeState(t *testing.T) {
	root, base := pendingMerge(t)
	plantMergeHeadDirectory(t, root)
	_, err := greenOwner().Land(context.Background(), topicLanding(root, base, mergeNameSet...))
	if err == nil || errors.As(err, new(NamedPathError)) || !strings.HasPrefix(err.Error(), "read pending merge:") {
		t.Fatalf("Land error = %v, want a plain read pending merge refusal", err)
	}
	if head := git(t, root, "rev-parse", "HEAD"); head != base {
		t.Fatalf("HEAD moved from %s to %s", base, head)
	}
}

// TestLandRefusesAnIncompleteMergeNameSet proves that during a pending merge the name set
// must cover each merged path, the conflicted path and a deleted path included. A directory
// name covers the merged paths below it and not a sibling that shares its prefix. The
// literal spelling of a merged path with hostile bytes covers that path. The rule applies to
// an octopus merge too. A refusal changes no ref, index entry, worktree byte, or merge state.
// (Coverage rows MP35-MP38 and MP42.)
func TestLandRefusesAnIncompleteMergeNameSet(t *testing.T) {
	// hostile holds a space, glob characters, a newline, and non-ASCII bytes.
	const hostile = "a b[x]*\né.txt"
	octopus := func(t *testing.T, _ ...string) (string, string) {
		root, base, _, _ := octopusMerge(t)
		return root, base
	}
	for _, row := range []struct {
		name         string
		start        func(t *testing.T, added ...string) (root, base string)
		added, paths []string
		// omitted is each omitted path as the refusal quotes it; an empty value publishes.
		omitted string
	}{
		{"conflicted path omitted", pendingMerge, nil, []string{mergeDeletedPath}, `"conflict.txt"`},
		{"deleted path omitted", pendingMerge, nil, []string{mergeConflictPath}, `"gone.txt"`},
		{"directory covers its subtree", pendingMerge, []string{"pkg/b.go", "pkg/sub/c.go"}, []string{mergeConflictPath, mergeDeletedPath, "pkg"}, ""},
		{"directory does not cover a prefix sibling", pendingMerge, []string{"pkg/b.go", "pkgx/a.go"}, []string{mergeConflictPath, mergeDeletedPath, "pkg"}, `"pkgx/a.go"`},
		{"hostile path spelled literally", pendingMerge, []string{hostile}, []string{mergeConflictPath, mergeDeletedPath, hostile}, ""},
		{"octopus side path omitted", octopus, nil, []string{"side-a"}, `"side-b"`},
	} {
		t.Run(row.name, func(t *testing.T) {
			root, base := row.start(t, row.added...)
			before := compositionState(t, root)
			got, err := greenOwner().Land(context.Background(), topicLanding(root, base, row.paths...))
			if row.omitted == "" {
				if err != nil {
					t.Fatalf("Land error = %v, want a publication", err)
				}
				names := append([]string{"foreign", "named", mergeConflictPath}, row.added...)
				assertPublishedTree(t, root, got, names, map[string]string{mergeConflictPath: mergeResolution}, nil)
				return
			}
			if !errors.As(err, new(NamedPathError)) || err.Error() != mergeNameSetRefusal+row.omitted {
				t.Fatalf("Land error = %v, want the named-path refusal %q", err, mergeNameSetRefusal+row.omitted)
			}
			if after := compositionState(t, root); !reflect.DeepEqual(after, before) {
				t.Fatalf("refusal changed the checkout\nbefore: %#v\nafter:  %#v", before, after)
			}
		})
	}
}

// TestLandRecordsEachMergeHeadAsAParent proves that during a pending merge the published
// commit keeps the expected base as parent 1 and records each merge head after it in file
// order, for one merge head and for the two heads of an octopus merge. (Coverage rows MP20
// and MP21.)
func TestLandRecordsEachMergeHeadAsAParent(t *testing.T) {
	for _, row := range []struct {
		name  string
		merge func(*testing.T) (root, base string, heads, paths []string)
	}{
		{"one merge head", func(t *testing.T) (string, string, []string, []string) {
			root, base := pendingMerge(t)
			return root, base, []string{git(t, root, "rev-parse", mergeIncoming)}, mergeNameSet
		}},
		{"octopus merge", octopusMerge},
	} {
		t.Run(row.name, func(t *testing.T) {
			root, base, heads, paths := row.merge(t)
			got, err := greenOwner().Land(context.Background(), topicLanding(root, base, paths...))
			if err != nil {
				t.Fatalf("Land error = %v, want a publication", err)
			}
			assertParents(t, root, got.Commit, append([]string{base}, heads...)...)
		})
	}
}

// TestComposeAfterMergeCommitReadsTheMergedBase proves that after the merge commit, the
// composition of the topic tip onto the incoming commit reads the incoming commit as the
// merge base and reports no conflict. (Coverage row MP23.)
func TestComposeAfterMergeCommitReadsTheMergedBase(t *testing.T) {
	root, base := pendingMerge(t)
	got, err := greenOwner().Land(context.Background(), topicLanding(root, base, mergeNameSet...))
	if err != nil {
		t.Fatal(err)
	}
	incoming := git(t, root, "rev-parse", mergeIncoming)
	composition, err := greenOwner().Compose(CompositionRequest{Root: root, Destination: incoming, Source: got.Commit})
	if err != nil || composition.Conflict.Kind != "" || composition.Base != incoming {
		t.Fatalf("Compose = %+v, %v; want no conflict over the merge base %s", composition, err, incoming)
	}
}

// TestLandRefusesAMergeStateChangeDuringAuthorization proves that Owner.Land reads the
// pending merge again after authorization. A merge that appears, a merge that disappears,
// another merge head, or a directory at MERGE_HEAD refuses with a plain error, and no ref
// moves. The recheck compares the merge heads in order, so the same heads in another order
// also refuse. (Coverage rows MP26-MP28.)
func TestLandRefusesAMergeStateChangeDuringAuthorization(t *testing.T) {
	const changed = "pending merge changed during the commit"
	// plantHeads writes heads to MERGE_HEAD as Git does, one head on each line.
	plantHeads := func(t *testing.T, mergeHead string, heads ...string) {
		if err := os.WriteFile(mergeHead, []byte(strings.Join(heads, "\n")+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	plant := func(t *testing.T, root, mergeHead string) {
		plantHeads(t, mergeHead, git(t, root, "rev-parse", "main"))
	}
	reverse := func(t *testing.T, _, mergeHead string) {
		heads := strings.Fields(string(mustRead(t, mergeHead)))
		slices.Reverse(heads)
		plantHeads(t, mergeHead, heads...)
	}
	// Each start returns the checkout, its HEAD, and the complete name set.
	noMerge := func(t *testing.T) (string, string, []string) {
		root, base := mergeCheckout(t)
		write(t, root, "named", "changed\n")
		return root, base, []string{"named"}
	}
	conflicted := func(t *testing.T) (string, string, []string) {
		root, base := pendingMerge(t)
		return root, base, mergeNameSet
	}
	octopus := func(t *testing.T) (string, string, []string) {
		root, base, _, added := octopusMerge(t)
		return root, base, added
	}
	for _, row := range []struct {
		name  string
		start func(*testing.T) (root, base string, paths []string)
		// change runs during authorization.
		change func(t *testing.T, root, mergeHead string)
		want   string
	}{
		{"merge appears", noMerge, plant, changed},
		{"merge disappears", conflicted, func(t *testing.T, _, mergeHead string) {
			if err := os.Remove(mergeHead); err != nil {
				t.Fatal(err)
			}
		}, changed},
		{"another merge head", conflicted, plant, changed},
		{"reordered merge heads", octopus, reverse, changed},
		{"directory at MERGE_HEAD", conflicted, func(t *testing.T, root, _ string) { plantMergeHeadDirectory(t, root) }, "read pending merge:"},
	} {
		t.Run(row.name, func(t *testing.T) {
			root, base, paths := row.start(t)
			o := greenOwner()
			o.authorize = func(context.Context, string, string, io.Writer, io.Writer) authorization.Result {
				row.change(t, root, mergeStatePath(t, root, "MERGE_HEAD"))
				return authorization.Result{Kind: authorization.Green}
			}
			_, err := o.Land(context.Background(), topicLanding(root, base, paths...))
			if err == nil || errors.As(err, new(NamedPathError)) || !strings.HasPrefix(err.Error(), row.want) {
				t.Fatalf("Land error = %v, want the plain refusal %q", err, row.want)
			}
			if head := git(t, root, "rev-parse", "HEAD"); head != base {
				t.Fatalf("HEAD moved from %s to %s", base, head)
			}
		})
	}
}

// TestLandPublishesAMergeWithAnUnchangedTree proves that a pending merge waives the
// empty-diff refusal. An ours resolution of every merged path, and an ours-strategy merge
// that changes no path, each publish a merge commit whose tree is the base tree. (Coverage
// rows MP29 and MP30.)
func TestLandPublishesAMergeWithAnUnchangedTree(t *testing.T) {
	for _, row := range []struct {
		name  string
		merge func(*testing.T) (root, base string, paths []string)
	}{
		{"ours resolution", func(t *testing.T) (string, string, []string) {
			root, base := pendingMerge(t)
			for _, path := range mergeNameSet {
				write(t, root, path, string(gitBytes(t, root, "show", "HEAD:"+path)))
			}
			return root, base, mergeNameSet
		}},
		{"no merged path", func(t *testing.T) (string, string, []string) {
			root, base := mergeCheckout(t)
			git(t, root, "merge", "--no-commit", "--no-ff", "-s", "ours", mergeIncoming)
			return root, base, []string{"named"}
		}},
	} {
		t.Run(row.name, func(t *testing.T) {
			root, base, paths := row.merge(t)
			got, err := greenOwner().Land(context.Background(), topicLanding(root, base, paths...))
			if err != nil {
				t.Fatalf("Land error = %v, want a publication", err)
			}
			if tree, want := git(t, root, "rev-parse", got.Commit+"^{tree}"), git(t, root, "rev-parse", base+"^{tree}"); tree != want {
				t.Fatalf("published tree = %s, want the base tree %s", tree, want)
			}
			assertParents(t, root, got.Commit, base, git(t, root, "rev-parse", mergeIncoming))
		})
	}
}

// TestLandEndsThePendingMerge proves that the merge commit leaves no unmerged index entry
// and none of the merge-state files that the conflicted merge wrote. (Coverage rows MP45 and
// MP47-MP50.)
func TestLandEndsThePendingMerge(t *testing.T) {
	root, base := pendingMerge(t)
	for _, name := range mergeStateFiles {
		if !mergeStatePresent(t, root, name) {
			t.Fatalf("%s is absent before the merge commit", name)
		}
	}
	if _, err := greenOwner().Land(context.Background(), topicLanding(root, base, mergeNameSet...)); err != nil {
		t.Fatal(err)
	}
	if unmerged := git(t, root, "ls-files", "-u"); unmerged != "" {
		t.Errorf("unmerged entries after the merge commit = %q", unmerged)
	}
	for _, name := range mergeStateFiles {
		if mergeStatePresent(t, root, name) {
			t.Errorf("%s remains after the merge commit", name)
		}
	}
}

// TestLandReportsPublishedMergeWhenReconciliationFails proves that a reconcile failure
// after the merge commit reports the published merge commit and leaves no pending merge.
// (Coverage rows MP52 and MP53.)
func TestLandReportsPublishedMergeWhenReconciliationFails(t *testing.T) {
	root, base := pendingMerge(t)
	o := greenOwner()
	o.reconcile = func(Request, []string, composedSnapshot) error { return os.ErrPermission }
	got, err := o.Land(context.Background(), topicLanding(root, base, mergeNameSet...))
	var remainder *PublishedUnreconciledError
	if !errors.As(err, &remainder) || !errors.Is(err, os.ErrPermission) || remainder.Commit != got.Commit || git(t, root, "rev-parse", "HEAD") != got.Commit {
		t.Fatalf("Land = %+v, %v; want the remainder of the published commit", got, err)
	}
	assertParents(t, root, got.Commit, base, git(t, root, "rev-parse", mergeIncoming))
	if mergeStatePresent(t, root, "MERGE_HEAD") {
		t.Fatal("MERGE_HEAD remains after the reconcile failure")
	}
}

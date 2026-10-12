package recordcmd

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/testreport"
	"github.com/gibbonmi/bench/internal/usage"
)

var foldForm = form{name: reviewrecord.FormFold, description: "record the fold of a default-branch commit at HEAD with its resolution paths and merged test packages",
	flags: []flag{{"--source", "<commit>", false, once}}, run: fold}

func fold(f form, root, spec string, parsed usage.Result) (string, int) {
	source, out, code := commit(f, root, parsed.Flags, "--source")
	if out != "" {
		return out, code
	}
	// A later commit can change the checkout that the test package index reads, so only
	// HEAD records.
	if head, err := git.ResolveCommit(root, "HEAD"); err != nil || head != source {
		return f.refuse(fmt.Sprintf("fold source %s is not HEAD; record the fold before the next commit", source))
	}
	entry, err := reviewrecord.RecordFold(root, spec, source, func(from, incoming string) (reviewrecord.FoldDelta, error) {
		return foldDelta(root, spec, from, incoming, source)
	})
	if err != nil {
		return f.refuse(cause(spec, err))
	}
	return row("fold", []string{"id", "after", "commit", "source_digest", "resolution", "test_packages"},
		[]string{entry.ID, entry.After, entry.Commit, entry.SourceDigest, strconv.Itoa(len(entry.Resolution)), strconv.Itoa(len(entry.TestPackages))})
}

// foldDelta derives the delta of the fold commit over its parents. The automatic tree is
// the landing composition of from and incoming, which `bench worktree merge` publishes. A
// conflicted composition keeps the tree that Git writes with its conflict markers.
func foldDelta(root, spec, from, incoming, commit string) (reviewrecord.FoldDelta, error) {
	composition, err := landing.New().Compose(landing.CompositionRequest{Root: root, Destination: from, Source: incoming})
	if err != nil {
		return reviewrecord.FoldDelta{}, err
	}
	automatic := composition.Tree
	if composition.Conflict.Kind != "" {
		automatic = composition.Conflict.Tree
	}
	record, err := reviewrecord.RecordPath(spec)
	if err != nil {
		return reviewrecord.FoldDelta{}, err
	}
	resolution, err := changedPaths(root, automatic, commit, func(path string) bool { return path != record })
	if err != nil {
		return reviewrecord.FoldDelta{}, err
	}
	diff, err := git.Raw("-C", root, "diff-tree", "-r", "-p", "--binary", "--full-index", "--no-renames", automatic, commit, "--", ":(exclude,literal)"+record)
	if err != nil {
		return reviewrecord.FoldDelta{}, fmt.Errorf("fold resolution diff: %w", err)
	}
	tests, err := changedPaths(root, from, commit, func(path string) bool { return strings.HasSuffix(path, "_test.go") })
	if err != nil {
		return reviewrecord.FoldDelta{}, err
	}
	floor, err := testreport.WritesFloor(context.Background(), root, tests)
	if err != nil {
		return reviewrecord.FoldDelta{}, err
	}
	return reviewrecord.FoldDelta{Resolution: resolution, ResolutionDigest: reviewrecord.Digest(diff), TestPackages: floor.Packages}, nil
}

// changedPaths lists, in byte order, each path that differs between the trees of from and
// to and that keep selects. The NUL-framed change list keeps the exact bytes of each path.
func changedPaths(root, from, to string, keep func(string) bool) ([]string, error) {
	changes, err := git.TreeChanges(root, from, to)
	if err != nil {
		return nil, err
	}
	paths := []string{}
	for _, change := range changes {
		if keep(change.Path) {
			paths = append(paths, change.Path)
		}
	}
	slices.Sort(paths)
	return paths, nil
}

package landing

import (
	"context"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

// mergeNameSetRefusal opens the refusal of a name set that omits a merged path.
const mergeNameSetRefusal = "pending merge changed paths that the commit does not name: "

// TestLandRefusesAnUnreadableMergeState proves that a directory at MERGE_HEAD refuses with
// a plain error before publication, although the name set covers every merged path.
// (Coverage row MP24.)
func TestLandRefusesAnUnreadableMergeState(t *testing.T) {
	root, base := pendingMerge(t)
	mergeHead, err := benchgit.AdminPath(root, "MERGE_HEAD")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(mergeHead); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(mergeHead, 0o755); err != nil {
		t.Fatal(err)
	}
	_, err = greenOwner().Land(context.Background(), Request{
		Root: root, Destination: "refs/heads/topic", Expected: base, Message: "merge",
		Paths: []string{mergeConflictPath, mergeDeletedPath},
	})
	if err == nil || errors.As(err, new(NamedPathError)) || !strings.HasPrefix(err.Error(), "read pending merge:") {
		t.Fatalf("Land error = %v, want a plain read pending merge refusal", err)
	}
	if head := git(t, root, "rev-parse", "HEAD"); head != base {
		t.Fatalf("HEAD moved from %s to %s", base, head)
	}
}

// TestLandRefusesAnIncompleteMergeNameSet proves that during a pending merge the name set
// must cover each merged path, the conflicted path and a deleted path included. A directory
// name covers the merged paths below it and not a sibling that shares its prefix, and the
// literal spelling of a merged path with hostile bytes covers that path. A refusal changes
// no ref, index entry, worktree byte, or merge state. (Coverage rows MP35-MP38 and MP42.)
func TestLandRefusesAnIncompleteMergeNameSet(t *testing.T) {
	// hostile holds a space, glob characters, a newline, and non-ASCII bytes.
	const hostile = "a b[x]*\né.txt"
	for _, row := range []struct {
		name         string
		added, paths []string
		// omitted is each omitted path as the refusal quotes it; an empty value publishes.
		omitted string
	}{
		{"conflicted path omitted", nil, []string{mergeDeletedPath}, `"conflict.txt"`},
		{"deleted path omitted", nil, []string{mergeConflictPath}, `"gone.txt"`},
		{"directory covers its subtree", []string{"pkg/b.go", "pkg/sub/c.go"}, []string{mergeConflictPath, mergeDeletedPath, "pkg"}, ""},
		{"directory does not cover a prefix sibling", []string{"pkg/b.go", "pkgx/a.go"}, []string{mergeConflictPath, mergeDeletedPath, "pkg"}, `"pkgx/a.go"`},
		{"hostile path spelled literally", []string{hostile}, []string{mergeConflictPath, mergeDeletedPath, hostile}, ""},
	} {
		t.Run(row.name, func(t *testing.T) {
			root, base := pendingMerge(t, row.added...)
			before := compositionState(t, root)
			got, err := greenOwner().Land(context.Background(), Request{
				Root: root, Destination: "refs/heads/topic", Expected: base, Message: row.name, Paths: row.paths,
			})
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

package recordcmd_test

import (
	"context"
	"fmt"
	"io"
	"slices"
	"strconv"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/landing"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/testreport"
	"github.com/gibbonmi/bench/internal/toon"
)

// folding is a delegated fixture with a committed chunk 1 entry. Its default branch starts
// at the fixture plan commit, so a default-branch change can conflict with chunk 1.
func folding(t *testing.T) *recordtest.Fixture {
	t.Helper()
	f := linked(t, 2)
	f.SyncMainline()
	base, tip := advance(f, "chunk 1")
	record(t, f, "1", base, tip)
	f.Commit("retain chunk 1")
	return f
}

func foldArgs(t *testing.T, source string) []string {
	return []string{"fold", slug(t), "--source", source}
}

// published folds incoming into the fixture checkout through the merge owner that
// `bench worktree merge` runs, under a lane whose one check passes. It returns the fold.
func published(t *testing.T, f *recordtest.Fixture, incoming string) string {
	t.Helper()
	previous := f.Tip()
	fingerprint, err := landing.CheckoutFingerprint(f.Root)
	if err != nil {
		t.Fatal(err)
	}
	lane := authorization.LaneAuthority{Base: previous, Checks: []gate.Phase{{Name: "pass", Argv: []string{"true"}}}}
	result, err := landing.NewLane(lane).Merge(context.Background(), landing.MergeRequest{
		Root: f.Root, Branch: f.Git("symbolic-ref", "HEAD"), PreviousTip: previous, Incoming: incoming,
		Worktree: f.Root, Fingerprint: fingerprint, Subject: "fold the default branch", Stdout: io.Discard, Stderr: io.Discard,
	})
	if err != nil {
		t.Fatal(err)
	}
	f.Git("reset", "--merge", result.Tip)
	return result.Tip
}

// folded runs the fold form at HEAD, requires exit 0, and returns its output and the one
// fold entry of the record.
func folded(t *testing.T, f *recordtest.Fixture) (string, rr.Fold) {
	t.Helper()
	out := succeed(t, f, foldArgs(t, f.Tip()))
	folds := read(t, f).Folds
	if len(folds) != 1 {
		t.Fatalf("record folds = %+v, want one entry", folds)
	}
	return out, folds[0]
}

// TestRecordFoldWritesTheFoldEntry covers BW25 and BW26: a fold that the merge owner
// published records F1 after chunk 1 with its parents and digest, and no resolution.
func TestRecordFoldWritesTheFoldEntry(t *testing.T) {
	f := folding(t)
	from := f.Tip()
	incoming := f.Mainline(map[string]string{"main.txt": "main\n"})
	commit := published(t, f, incoming)
	out, entry := folded(t, f)
	got := []string{entry.ID, entry.After, entry.From, entry.Commit, entry.Incoming, entry.SourceDigest}
	if want := []string{"F1", "1", from, commit, incoming, sourceDigest(t, f)}; !slices.Equal(got, want) {
		t.Fatalf("fold entry id, after, from, commit, incoming, and digest = %q, want %q", got, want)
	}
	if len(entry.Resolution) != 0 {
		t.Fatalf("published fold resolution = %q, want none", entry.Resolution)
	}
	table, err := toon.Table("fold", []string{"id", "after", "commit", "source_digest", "resolution", "test_packages"},
		[][]string{{entry.ID, entry.After, entry.Commit, entry.SourceDigest, strconv.Itoa(len(entry.Resolution)), strconv.Itoa(len(entry.TestPackages))}})
	if err != nil {
		t.Fatal(err)
	}
	if out != table {
		t.Fatalf("output = %q, want %q", out, table)
	}
}

// TestRecordFoldNamesTheResolutionPaths covers BW27, BW28, and BW29: the resolution lists
// exactly the paths that differ from the automatic merge, each with its exact bytes.
func TestRecordFoldNamesTheResolutionPaths(t *testing.T) {
	for _, tc := range []struct {
		name, path string
		main       map[string]string
	}{
		{"conflict", "source.txt", map[string]string{"source.txt": "main\n", "main.txt": "main\n"}},
		{"hand edit", "notes.txt", map[string]string{"main.txt": "main\n"}},
		{"quoted path", `a "b c.txt`, map[string]string{"main.txt": "main\n"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := folding(t)
			f.HandFold(f.Mainline(tc.main), func() { f.Write(tc.path, "resolved\n") })
			if _, entry := folded(t, f); !slices.Equal(entry.Resolution, []string{tc.path}) {
				t.Fatalf("resolution = %q, want exactly %q", entry.Resolution, tc.path)
			}
		})
	}
}

// TestRecordFoldRefusesANonFold covers BW30 to BW33: each commit that is not a fold of a
// default-branch commit over the reviewed source refuses before any write.
func TestRecordFoldRefusesANonFold(t *testing.T) {
	mainline := func(f *recordtest.Fixture) string { return f.Mainline(map[string]string{"main.txt": "main\n"}) }
	for _, tc := range []struct {
		name, want string
		source     func(*recordtest.Fixture) string
	}{
		{"not HEAD", "fold source %s is not HEAD; record the fold before the next commit", func(f *recordtest.Fixture) string {
			fold := f.HandFold(mainline(f), nil)
			advance(f, "later")
			return fold
		}},
		{"one parent", "fold %s is not a two-parent merge", func(f *recordtest.Fixture) string {
			_, tip := advance(f, "one parent")
			return tip
		}},
		{"sibling tip", "fold %s does not merge a default-branch commit", func(f *recordtest.Fixture) string {
			return f.HandFold(f.Git("commit-tree", f.Tree(), "-p", f.Tip(), "-m", "sibling"), nil)
		}},
		{"unreviewed work", "fold %s does not follow the reviewed source", func(f *recordtest.Fixture) string {
			advance(f, "unreviewed")
			return f.HandFold(mainline(f), nil)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := folding(t)
			source := tc.source(f)
			refuseArgs(t, f, fmt.Sprintf(tc.want, source), foldArgs(t, source))
		})
	}
}

// TestRecordFoldDerivesMergedTestPackages covers BW41: the merged test packages are the
// floor that testreport.WritesFloor derives over the merged test file.
func TestRecordFoldDerivesMergedTestPackages(t *testing.T) {
	f := folding(t)
	test := "pkg/a/a_test.go"
	f.HandFold(f.Mainline(map[string]string{"go.mod": "module example.com/fold\n\ngo 1.20\n", "pkg/a/a.go": "package a\n", test: "package a\n"}), nil)
	floor, err := testreport.WritesFloor(context.Background(), f.Root, []string{test})
	if err != nil || len(floor.Packages) == 0 {
		t.Fatalf("floor of %s = %+v, %v; want a package", test, floor, err)
	}
	if _, entry := folded(t, f); !slices.Equal(entry.TestPackages, floor.Packages) {
		t.Fatalf("test packages = %q, want %q", entry.TestPackages, floor.Packages)
	}
}

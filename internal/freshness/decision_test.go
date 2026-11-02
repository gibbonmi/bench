package freshness

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
)

func TestSelectFreshnessFromImmutableDigests(t *testing.T) {
	digestA := strings.Repeat("a", 64)
	digestB := strings.Repeat("b", 64)
	accepted := SelectionInput{StoredSource: digestA, CurrentSource: digestA, StoredExecutable: digestB, CurrentExecutable: digestB}
	for _, tc := range []struct {
		name  string
		input SelectionInput
		want  Selection
	}{
		{name: "accepted", input: accepted, want: Selection{Accepted: true}},
		{name: "source refusal", input: SelectionInput{StoredSource: digestA, CurrentSource: digestB, StoredExecutable: digestB, CurrentExecutable: digestB}, want: Selection{Reason: "seal source digest does not match current build inputs"}},
		{name: "executable refusal", input: SelectionInput{StoredSource: digestA, CurrentSource: digestA, StoredExecutable: digestA, CurrentExecutable: digestB}, want: Selection{Reason: "seal executable digest does not match binary contents"}},
		{name: "empty", input: SelectionInput{}, want: Selection{Reason: "stored source digest is malformed"}},
		{name: "hostile", input: SelectionInput{StoredSource: "\x1b", CurrentSource: digestA, StoredExecutable: digestB, CurrentExecutable: digestB}, want: Selection{Reason: "stored source digest is malformed"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			first := Select(tc.input)
			second := Select(tc.input)
			if first != tc.want || second != first {
				t.Fatalf("Select() = %#v then %#v, want %#v", first, second, tc.want)
			}
		})
	}
}

// TestDiffChangesBuildInputs is CW01, CW02, CW03, CW41, CW42, CW49, and CW51. Each row
// commits the build fixture with two more Go files, one with a non-ASCII name, and a test
// file in its cmd/bench package, then
// commits one change and reads the diff between the two commits.
func TestDiffChangesBuildInputs(t *testing.T) {
	for _, tc := range []struct {
		name   string
		linked bool
		change []string
		write  map[string]string
		want   bool
	}{
		{name: "build script", write: map[string]string{"scripts/go-build.sh": "#!/usr/bin/env bash\nexit 0\n"}, want: true},
		{name: "markdown only", write: map[string]string{"README.md": "# fixture\n"}},
		{name: "no build-input manifest", linked: true, write: map[string]string{"scripts/go-build.sh": "#!/usr/bin/env bash\nexit 0\n"}},
		{name: "deleted input", change: []string{"rm", "-q", "cmd/bench/second.go"}, want: true},
		{name: "deleted test file", change: []string{"rm", "-q", "cmd/bench/main_test.go"}},
		{name: "input renamed to a test file", change: []string{"mv", "cmd/bench/second.go", "cmd/bench/second_test.go"}, want: true},
		{name: "deleted non-ASCII input", change: []string{"rm", "-q", "cmd/bench/zweiteé.go"}, want: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := writeBuildFixtureAt(t, gittest.RepoOnBranch(t, "main"))
			writeFixtureFiles(t, root, map[string]string{
				"cmd/bench/second.go":    "package main\n\nconst second = 2\n",
				"cmd/bench/main_test.go": "package main\n",
				"cmd/bench/zweiteé.go":   "package main\n\nconst zweite = 2\n",
			})
			if tc.linked {
				if err := os.Remove(filepath.Join(root, filepath.FromSlash(BuildInputsManifest))); err != nil {
					t.Fatal(err)
				}
			}
			gittest.MustRun(t, root, "add", "-A")
			gittest.Commit(t, root, "base")
			base := gittest.Output(t, root, "rev-parse", "HEAD")
			writeFixtureFiles(t, root, tc.write)
			if tc.change != nil {
				gittest.MustRun(t, root, tc.change...)
			}
			gittest.MustRun(t, root, "add", "-A")
			gittest.Commit(t, root, tc.name)
			if got := DiffChangesBuildInputs(root, base, "HEAD"); got != tc.want {
				t.Fatalf("DiffChangesBuildInputs(%s) = %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}

// TestRefreshBuildsOnlyAStaleExecutable is CW04, CW05, and CW08. A fresh published build
// takes no build, and a stale or absent one takes one build that publishes a valid seal.
func TestRefreshBuildsOnlyAStaleExecutable(t *testing.T) {
	for _, tc := range []struct {
		name      string
		root      func(*testing.T) string
		wantCalls int
	}{
		{name: "fresh", root: func(t *testing.T) string {
			root, _ := writePublishedFixture(t)
			return root
		}},
		{name: "stale", wantCalls: 1, root: func(t *testing.T) string {
			root, _ := writePublishedFixture(t)
			writeFixtureFiles(t, root, map[string]string{"internal/local/local.go": "package local\n\nconst Value = \"next\"\n"})
			return root
		}},
		{name: "absent", wantCalls: 1, root: writeBuildFixture},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			rebuilt, buildErr, err := Refresh(context.Background(), tc.root(t), func(_ context.Context, root, executable string) error {
				calls++
				return publishFixtureBuild(t, root, executable)
			})
			if calls != tc.wantCalls || rebuilt != (tc.wantCalls == 1) || buildErr != nil || err != nil {
				t.Fatalf("Refresh(%s) = (%v, %v, %v) after %d builds, want rebuilt %v with no errors after %d builds", tc.name, rebuilt, buildErr, err, calls, tc.wantCalls == 1, tc.wantCalls)
			}
		})
	}
}

// TestRefreshSecondReadDecides is CW06 and CW07. A build that returns nil and publishes
// nothing reports no rebuild, and a failing build's error comes back as buildErr.
func TestRefreshSecondReadDecides(t *testing.T) {
	boom := errors.New("boom")
	for _, tc := range []struct {
		name  string
		build error
	}{
		{name: "nil build that publishes nothing"},
		{name: "failing build", build: boom},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rebuilt, buildErr, err := Refresh(context.Background(), writeBuildFixture(t), func(context.Context, string, string) error {
				return tc.build
			})
			if rebuilt || !errors.Is(buildErr, tc.build) || err == nil {
				t.Fatalf("Refresh(%s) = (%v, %v, %v), want no rebuild, build error %v, and a seal refusal", tc.name, rebuilt, buildErr, err, tc.build)
			}
		})
	}
}

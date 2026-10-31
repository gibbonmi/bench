package preflight

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// TestMapSourcesHostileEntries checks the hostile tree entries of the map read and of the
// source entry: a map one byte over the read bound, a source that is a link to a regular
// file, a source that is a directory, a source that is a submodule entry, a candidate map
// stored as a link, and a tip tree that does not read are each red.
func TestMapSourcesHostileEntries(t *testing.T) {
	for _, tc := range []struct {
		name, want string
		// prefix marks a want that opens the detail, because git words the rest.
		prefix bool
		setup  func(t *testing.T)
		after  func(t *testing.T)
	}{
		{name: "over-limit map", want: "decisions/big.md: not readable at the source tip: read limit exceeded",
			setup: func(t *testing.T) {
				body := preflighttest.DecisionMap(preflighttest.PathRecord("docs/kept.md"))
				preflighttest.MustWriteFile(t, "decisions/big.md", body+strings.Repeat("\n", int(bounds.ControlRecordLimit)+1-len(body)))
			}},
		{name: "link source", want: "decisions/m.md: Sources Path docs/link.md" + missingFile,
			setup: func(t *testing.T) {
				symlink(t, "kept.md", "docs/link.md")
				preflighttest.WriteDecisionMap(t, "decisions/m.md", preflighttest.PathRecord("docs/link.md"))
			}},
		{name: "directory source", want: "decisions/m.md: Sources Path docs" + missingFile,
			setup: func(t *testing.T) {
				preflighttest.WriteDecisionMap(t, "decisions/m.md", preflighttest.PathRecord("docs"))
			}},
		{name: "submodule source", want: "decisions/m.md: Sources Path vendor/sub" + missingFile,
			setup: func(t *testing.T) {
				preflighttest.WriteDecisionMap(t, "decisions/m.md", preflighttest.PathRecord("vendor/sub"))
			},
			// The add of every change drops an index entry that the checkout lacks, so the
			// gitlink commits on its own, over an empty folder that keeps the checkout clean.
			after: func(t *testing.T) {
				if err := os.MkdirAll("vendor/sub", 0o755); err != nil {
					t.Fatal(err)
				}
				head := preflighttest.RunGit(t, "rev-parse", "HEAD")
				preflighttest.RunGit(t, "update-index", "--add", "--cacheinfo", "160000,"+head+",vendor/sub")
				preflighttest.RunGit(t, "commit", "-q", "-m", "gitlink source")
				if mode := preflighttest.RunGit(t, "ls-tree", "--format=%(objectmode)", "HEAD", "--", "vendor/sub"); mode != "160000" {
					t.Fatalf("vendor/sub tree mode = %q, want 160000", mode)
				}
			}},
		{name: "link candidate map", want: "decisions/link.md: candidate map has tree mode 120000, not a regular file",
			setup: func(t *testing.T) {
				preflighttest.WriteDecisionMap(t, "docs/map.md", preflighttest.PathRecord("docs/never.md"))
				symlink(t, "../docs/map.md", "decisions/link.md")
			}},
		{name: "unreadable tip tree", want: "decision maps not readable at the source tip: ", prefix: true,
			setup: func(t *testing.T) {
				preflighttest.WriteDecisionMap(t, "decisions/m.md", preflighttest.PathRecord("docs/kept.md"))
			},
			after: func(t *testing.T) {
				tree := preflighttest.RunGit(t, "rev-parse", "HEAD:decisions")
				if err := os.Remove(filepath.Join(".git", "objects", tree[:2], tree[2:])); err != nil {
					t.Fatal(err)
				}
			}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			slug := seedMaps(t, func(t *testing.T, _ string) { tc.setup(t) })
			if tc.after != nil {
				tc.after(t)
			}
			got := mapSourcesDetail(t, slug)
			if got != tc.want && !(tc.prefix && strings.HasPrefix(got, tc.want)) {
				t.Fatalf("map-sources detail = %q, want %q", got, tc.want)
			}
		})
	}
}

// symlink makes path a link to target and creates the parent folder of path.
func symlink(t *testing.T, target, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, path); err != nil {
		t.Fatal(err)
	}
}

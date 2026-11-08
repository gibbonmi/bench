//go:build system

package systemtest

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/prose"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
)

// markdownMerge is a copy of the kit with one assignment worktree, the merge target. The
// copy holds the kit's own build inputs, so the selected executable matches each lane
// checkout of a merge, and the built-in lane grades a Markdown fold. The target commits a
// stand-in wrapper that appends the tip of each lane pass that a merge starts to record.
type markdownMerge struct {
	repo, record string
	env          []string
	target       systemLandingWorktree
}

// newMarkdownMerge builds the fixture. The copy drops the kit's commitment policy, so the
// target needs no admission, and the target's own wrapper commit diverges it from main, so
// a fold of a later main commit is a merge.
func newMarkdownMerge(t *testing.T) markdownMerge {
	t.Helper()
	repo := copyKitRepository(t, "markdown merge [primary]-")
	systemIdentity(t, repo)
	systemGit(t, repo, "rm", "-q", ".bench/commitment.json")
	systemGit(t, repo, systemCommitArgs("-qm", "drop the commitment policy")...)
	home, err := os.MkdirTemp(owner.root, "markdown merge [home]-")
	if err != nil {
		t.Fatal(err)
	}
	f := markdownMerge{repo: repo, record: filepath.Join(home, "lane-passes")}
	f.target = systemCreateLandingWorktree(t, repo, home, "markdown-merge", "markdown-merge")
	f.env = []string{"BENCH_HOME=" + home, "BENCH_KIT=" + repo, "BENCH_RUN_BINARY=" + owner.selected.path, "BENCH_SYSTEM_ROOT=" + repo}
	// The copied wrapper is executable, and a rewrite keeps its mode.
	systemCommit(t, f.target.path, "bin/bench.sh", "#!/bin/sh\nprintf '%s\\n' \"$3\" >> "+sanitize.ShellQuote(f.record)+"\n", "stand-in wrapper")
	return f
}

// merge commits body at name on main and merges that commit into the target. It requires a
// merge of kind merge and answers the new target tip.
func (f markdownMerge) merge(t *testing.T, name, body string) string {
	t.Helper()
	systemCommit(t, f.repo, name, body, "advance "+name)
	incoming := systemGitOutput(t, f.repo, "rev-parse", "HEAD")
	merged := systemSelected(t, f.repo, f.env, "worktree", "merge", "--from", incoming, f.target.assignment)
	if merged.code != 0 || !strings.Contains(merged.stdout, "kind=merge") {
		t.Fatalf("bench worktree merge of %s = (%d, %q, %q), want a merge of kind merge", name, merged.code, merged.stdout, merged.stderr)
	}
	return systemGitOutput(t, f.target.path, "rev-parse", "HEAD")
}

// passes answers the tip of each lane pass that the stand-in wrapper recorded.
func (f markdownMerge) passes(t *testing.T) []string {
	t.Helper()
	data, err := os.ReadFile(f.record)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	return strings.Fields(string(data))
}

// TestMarkdownOnlyMergeStartsNoGrade is SP23. A merge of kind merge whose fold brings only one
// Markdown file starts no lane pass. The next merge, of a prose-policy change, starts one pass,
// and the test waits for it, so a pass of the earlier merge, which started first, is in the
// record too.
func TestMarkdownOnlyMergeStartsNoGrade(t *testing.T) {
	t.Parallel()
	f := newMarkdownMerge(t)
	markdown := f.merge(t, "notes.md", "# Notes\n\nThis file holds notes.\n")
	if got := strings.Fields(systemGitOutput(t, f.target.path, "diff", "--name-only", markdown+"^1", markdown)); !slices.Equal(got, []string{"notes.md"}) {
		t.Fatalf("merge tip changes %q, want only notes.md", got)
	}
	policy, err := os.ReadFile(filepath.Join(f.repo, filepath.FromSlash(prose.ExclusionFile)))
	if err != nil {
		t.Fatal(err)
	}
	next := f.merge(t, prose.ExclusionFile, string(policy)+"# lane fixture\n")
	processtest.WaitFor(t, "the lane pass of the next merge", func() bool { return slices.Contains(f.passes(t), next) })
	if got := f.passes(t); !slices.Equal(got, []string{next}) {
		t.Fatalf("lane passes of %q, want only the pass of the next merge %s", got, next)
	}
}

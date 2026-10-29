// The background grade that a merge starts: a merge that publishes a new tip in a kit
// worktree starts the lane pass of that tip, and a merge with nothing to publish starts none.
package worktree

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate/backgroundstate"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
)

// launchMergeSet is a kit merge set whose target holds a committed wrapper that is a stand-in
// grade of the answered launches. The commit diverges the target, so a fold of a later
// default-branch commit is a merge. The kit is the primary checkout, as the wrapper of a kit
// worktree resolves it.
func launchMergeSet(t *testing.T) (mergeSet, Creation, *processtest.Launches) {
	t.Helper()
	f := goLaneMergeFixture(t, "exit 0")
	f.kit = f.root
	target := f.created[0]
	launches := processtest.NewLaunches(t)
	// A child writes the wrapper, so no descriptor of this process holds it open for write
	// when the merge runs it.
	write := descendant(t, "sh", "-c", `mkdir -p "$(dirname "$1")" && cat >"$1" && chmod 755 "$1"`, "sh", filepath.Join(target.Path, "bin", "bench.sh"))
	write.Stdin = strings.NewReader(launches.StandIn())
	if out, err := write.CombinedOutput(); err != nil {
		t.Fatalf("write the wrapper: %v, %s", err, out)
	}
	gitRun(t, target.Path, "add", "bin/bench.sh")
	gitRun(t, target.Path, withIdentity("commit", "-qm", "wrapper")...)
	return f, target, launches
}

// mergeRunning answers the running grade of the machine queue under the home of f, or nil.
func mergeRunning(t *testing.T, f mergeSet) *backgroundstate.Slot {
	t.Helper()
	queue, err := backgroundstate.ReadQueue(f.home)
	mustNoError(t, err)
	return queue.Running
}

// TestMergeStartsABackgroundGrade is BG6. A merge of kind merge starts one lane pass, and the test
// waits for that pass. The pass names the target worktree and the new tip.
func TestMergeStartsABackgroundGrade(t *testing.T) {
	t.Parallel()
	f, target, launches := launchMergeSet(t)
	r := runVerb(t, verbMerge, mergeGoModLine(t, f, "// lane fixture"))
	if r.exit != 0 || !strings.Contains(mergedRecord(t, r.stdout), "kind=merge") {
		t.Fatalf("merge = (%d, %q, %q), want a merge of kind merge at exit 0", r.exit, r.stdout, r.stderr)
	}
	want := []string{target.Assignment.Worktree, gitOutput(t, f.root, "rev-parse", target.Assignment.Branch)}
	launches.Await(t)
	launches.End(t)
	if got := launches.Launched(t); len(got) != 1 || len(got[0].Args) != 3 || !slices.Equal(got[0].Args[1:], want) {
		t.Fatalf("launches = %+v, want one lane pass of %q", got, want)
	}
}

// TestCurrentMergeStartsNoGrade is BG10. A merge of kind current publishes nothing, so it takes
// no slot and starts no grade.
func TestCurrentMergeStartsNoGrade(t *testing.T) {
	t.Parallel()
	f, target, launches := launchMergeSet(t)
	contained := gitOutput(t, f.root, "rev-parse", "HEAD")
	r := runVerb(t, verbMerge, f.merge("--from", contained, target.Assignment.ID))
	if r.exit != 0 || !strings.Contains(mergedRecord(t, r.stdout), "kind=current") {
		t.Fatalf("merge = (%d, %q, %q), want a merge of kind current at exit 0", r.exit, r.stdout, r.stderr)
	}
	launches.End(t)
	if running, got := mergeRunning(t, f), launches.Launched(t); running != nil || len(got) != 0 {
		t.Fatalf("running grade = %+v, launches = %+v; want no grade", running, got)
	}
}

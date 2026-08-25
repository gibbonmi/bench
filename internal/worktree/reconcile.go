package worktree

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

// specbuildRefNamespace holds refs from the retired provisional spec-build lifecycle.
// The standing cleaner is its only remaining reader.
const specbuildRefNamespace = "refs/bench/specbuild/"

// lifecycleRefNamespaces bounds what the sweep may delete. Everything else under
// refs/bench/ belongs to another owner — the gate's `refs/bench/green/<branch>` verdict
// store above all. The two namespaces are enumerated rather than derived from the shared
// prefix. A ref is deleted only after its own name is checked against them.
var lifecycleRefNamespaces = []string{specbuildRefNamespace, intent.RecoveryRefNamespace}

func insideLifecycleNamespace(ref string) bool {
	for _, namespace := range lifecycleRefNamespaces {
		if strings.HasPrefix(ref, namespace) {
			return true
		}
	}
	return false
}

// sweepLifecycleRefs deletes every ref inside the lifecycle namespaces, whatever its name
// under them, and reports how many it deleted. Each ref is deleted against the object the
// listing read, so a ref something else moved between the two is refused rather than
// dropped blind. Deletions are independent: a run killed partway leaves the refs it
// already deleted deleted, and the rest exactly as they were. That is a state the next
// run finishes rather than misreads.
func sweepLifecycleRefs(j joins, root string) (int, error) {
	args := append([]string{"-C", root, "for-each-ref", "--format=%(refname) %(objectname)"}, lifecycleRefNamespaces...)
	listing, err := git.Output(args...)
	if err != nil {
		return 0, fmt.Errorf("list lifecycle refs: %w", err)
	}
	swept := 0
	for _, line := range strings.Split(listing, "\n") {
		ref, oid, ok := strings.Cut(line, " ")
		if !ok || !insideLifecycleNamespace(ref) {
			continue
		}
		if err := hit(j.cleanupBoundary, StepLifecycleSweep); err != nil {
			return swept, err
		}
		if out, err := exec.Command("git", "-C", root, "update-ref", "-d", ref, oid).CombinedOutput(); err != nil {
			return swept, fmt.Errorf("delete lifecycle ref %s: %s", ref, strings.TrimSpace(string(out)))
		}
		swept++
	}
	return swept, nil
}

// poolAssignment reports whether a ledger record is one the surviving worktree pool wrote
// and still answers for. Everything else is debris this reconcile drops. That covers a
// record in a state only the removed lifecycle produced, one naming a gone checkout, or
// one no current decoder can read.
//
// Only proven absence licenses the drop, and even then three claims outlive it. Any other
// stat error leaves the checkout's existence unknown. Dropping on unknown would delete the
// only claim on a worktree whose uncommitted work is still there. A registration git still
// holds is a claim the release path resumes from, so it outlives an absent path until the
// registration itself is pruned. An active record too young to be orphaned is a live
// session's. A drop on tree-absence alone would catch it between its `worktree add` and
// its first write.
func poolAssignment(a intent.Assignment, registered []Registered, now time.Time) bool {
	if a.State != intent.StateActive && a.State != intent.StateCleanupPending {
		return false
	}
	if _, err := os.Stat(a.Worktree); !os.IsNotExist(err) {
		return true
	}
	if isRegisteredWorktree(registered, a.Worktree) {
		return true
	}
	return a.State == intent.StateActive && !orphaned(a, now)
}

// reconcileLifecycleDebris is the standing cleaner every session start runs. It empties
// the two lifecycle ref namespaces and purges the ledger records the removed lifecycle
// left behind, reporting what each half removed. Refs go first, so a run killed between
// the halves leaves records whose refs are already gone. That is the same state a
// repository carrying only ledger debris is in, one the next run finishes from.
// The instant is the caller's explicit boundary resolution: one instant for the whole
// pass, so two records of the same age cannot straddle the staleness window and disagree.
func reconcileLifecycleDebris(j joins, root string, registered []Registered, now time.Time) (int, int, error) {
	swept, err := sweepLifecycleRefs(j, root)
	if err != nil {
		return swept, 0, err
	}
	purged, err := intent.PurgeAssignments(root, func(a intent.Assignment) bool {
		return poolAssignment(a, registered, now)
	})
	return swept, purged, err
}

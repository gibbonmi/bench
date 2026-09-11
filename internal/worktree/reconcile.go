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

// lifecycleRefNamespaces names the namespaces that the cleaner empties.
// Reset refs have a separate record-bound rule. Green verdict refs stay outside both rules.
var lifecycleRefNamespaces = []string{specbuildRefNamespace, intent.RecoveryRefNamespace}

func insideLifecycleNamespace(ref string) bool {
	for _, namespace := range lifecycleRefNamespaces {
		if strings.HasPrefix(ref, namespace) {
			return true
		}
	}
	return false
}

// sweepLifecycleRefs empties the lifecycle namespaces and deletes reset refs with no record.
// Each deletion checks the listed object, so a concurrent ref move refuses.
func sweepLifecycleRefs(j joins, root string) (int, error) {
	args := append([]string{"-C", root, "for-each-ref", "--format=%(refname) %(objectname)"}, lifecycleRefNamespaces...)
	args = append(args, intent.ResetRefNamespace)
	listing, err := git.Output(args...)
	if err != nil {
		return 0, fmt.Errorf("list lifecycle refs: %w", err)
	}
	// The ledger read comes after the listing, so a record written between the two reads
	// still protects the ref it names.
	resetPrefixes := recordedResetPrefixes(root)
	swept := 0
	for _, line := range strings.Split(listing, "\n") {
		ref, oid, ok := strings.Cut(line, " ")
		if !ok {
			continue
		}
		if !insideLifecycleNamespace(ref) {
			if !strings.HasPrefix(ref, intent.ResetRefNamespace) || resetPrefixes == nil {
				continue
			}
			parts := strings.SplitN(strings.TrimPrefix(ref, intent.ResetRefNamespace), "/", 3)
			if len(parts) == 3 && resetPrefixes[intent.ResetRefPrefix(parts[0], parts[1])] {
				continue
			}
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

func recordedResetPrefixes(root string) map[string]bool {
	// An absent ledger is not proof that every reset envelope has lost its record.
	path, err := intent.Address(root)
	if err != nil {
		return nil
	}
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() {
		return nil
	}
	assignments, err := intent.Assignments(root)
	if err != nil {
		return nil
	}
	prefixes := make(map[string]bool, len(assignments))
	for _, assignment := range assignments {
		prefixes[intent.ResetRefPrefix(assignment.OwnerID, assignment.ID)] = true
	}
	return prefixes
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

// reconcileLifecycleDebris removes obsolete refs before it purges obsolete records.
// A reset ref keeps its record's protection through this pass, even when the purge removes that record.
// The caller supplies one instant so records of the same age get the same decision.
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

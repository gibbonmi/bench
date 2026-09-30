package preflight

import (
	"errors"
	"os"

	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/usage"
)

// binarySealFacts grades root's published binary through the seal verifier,
// the one primitive that decides staleness by source digest. An absent binary
// is reported as absent rather than graded, because a linked consumer repo
// publishes none. The verifier's refusal is carried whole, so the row prints
// the reason and the rebuild sentence the verifier composed rather than a
// second copy.
//
// Only a name that resolves to nothing is an absence. The link itself is the
// artifact, so a dangling or non-regular dist/bench is present and graded,
// which is the state doctor's seal row already reds. Stat here would follow
// the link and report the one root that publishes a broken binary as the one
// root with no binary to grade.
func binarySealFacts(root string) (present bool, refusal *freshness.Refusal) {
	executable := freshness.PublishedExecutable(root)
	if _, err := os.Lstat(executable); err != nil {
		return false, nil
	}
	if err := freshness.Verify(root, executable); err != nil {
		if errors.As(err, &refusal) {
			return true, refusal
		}
		return true, &freshness.Refusal{Executable: executable, RepairRoot: root, Cause: err}
	}
	return true, nil
}

// binarySealCheck grades the published dist/bench a build is about to run.
// A stale binary answers every consumer — a hand run, a hook, the wrapper, the
// landing — from sources nobody reviewed, and no other row touches it. A root
// that publishes no binary is not applicable, so a linked consumer repo is
// reported on rather than refused. Decide calls this in build mode alone; a
// review preflight renders no such row.
//
// An assignment worktree sits under the pool path, and the hooks refuse the cd
// the verifier's rebuild sentence names. So an owned root states the reason
// alone and names the worktree build verb as its remedy. Any other root renders
// the gathered refusal whole, because its cd leaves the pool path untouched.
func binarySealCheck(f Facts) CheckResult {
	if !f.BinarySealPresent {
		return notApplicable("binary-seal")
	}
	if f.BinarySealRefusal == nil {
		return green("binary-seal")
	}
	if f.AssignmentTarget == "" {
		return red("binary-seal", f.BinarySealRefusal.Error())
	}
	row := red("binary-seal", f.BinarySealRefusal.Reason())
	row.Next = usage.WorktreeBuildFor(f.AssignmentTarget)
	return row
}

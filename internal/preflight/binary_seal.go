package preflight

import (
	"errors"
	"io/fs"
	"os"

	"github.com/gibbonmi/bench/internal/freshness"
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
// root with no binary to grade. An lstat that fails for any other cause cannot
// show that no binary exists, so the binary is present and refused with that
// cause rather than graded.
//
// The rebuild is the owner's remedy and root kind for root in the pool of home, so the
// row and the hook grade the same pool.
func binarySealFacts(home, root string) (present bool, refusal *freshness.Refusal, rebuild freshness.Rebuild) {
	rebuild = freshness.RebuildIn(home, root)
	executable := freshness.PublishedExecutable(root)
	_, err := os.Lstat(executable)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil, rebuild
	}
	if err == nil {
		err = freshness.Verify(root, executable)
	}
	if err == nil {
		return true, nil, rebuild
	}
	if errors.As(err, &refusal) {
		return true, refusal, rebuild
	}
	return true, &freshness.Refusal{Executable: executable, RepairRoot: root, Cause: err}, rebuild
}

// binarySealCheck grades the published dist/bench a build is about to run.
// A stale binary answers every consumer — a hand run, a hook, the wrapper, the
// landing — from sources nobody reviewed, and no other row touches it. A root
// that publishes no binary is not applicable, so a linked consumer repo is
// reported on rather than refused. Decide calls this in build mode alone; a
// review preflight renders no such row.
//
// The gathered rebuild owns the root kind. A worktree rebuild names the build
// verb, so the row states the reason alone and names that remedy as its next
// step. Any other root renders the gathered refusal whole, because its remedy
// runs as the refusal prints it.
func binarySealCheck(f Facts) CheckResult {
	if !f.BinarySealPresent {
		return notApplicable("binary-seal")
	}
	if f.BinarySealRefusal == nil {
		return green("binary-seal")
	}
	if !f.BinarySealRebuild.Worktree {
		return red("binary-seal", f.BinarySealRefusal.Error())
	}
	row := red("binary-seal", f.BinarySealRefusal.Reason())
	row.Next = f.BinarySealRebuild.Action
	return row
}

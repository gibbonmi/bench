package specbuild

import (
	"os"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
)

// selfSymlinkCheckout replaces the assignment's checkout with a symlink pointing at
// its own path — a loop the real classifier's Stat can never resolve, so
// ClassifyPathShape returns worktree.ShapeUnknown for it. That is the cheapest of the
// three ShapeUnknown fixtures to reach at an assignment path: the ENOTDIR site needs
// the pool parent turned into a file, which the precondition tests price out as too
// heavy for this composition.
func selfSymlinkCheckout(t *testing.T, path string) {
	t.Helper()
	replaceCheckout(t, path)
	if err := os.Symlink(path, path); err != nil {
		capability.Capability(t, capability.Symlink, "cannot create a symlink: "+err.Error())
	}
}

// TestCheckpointRefusesSelfSymlinkCheckout composes the real worktree.ClassifyPathShape
// into the precondition's ownership refusal: liveCheckout's default branch discards the
// classifier's error and treats every undecided shape as an identity fault, so an
// assignment whose checkout path cannot be resolved must refuse rather than proceed.
func TestCheckpointRefusesSelfSymlinkCheckout(t *testing.T) {
	fixture, owner := decayedFixture(t, selfSymlinkCheckout)
	if err := checkpointInvocation(t, fixture); err == nil || !strings.Contains(err.Error(), ownershipRefusal) {
		t.Fatalf("Checkpoint over a self-symlink checkout = %v, want the ownership refusal", err)
	}
	if owner.plans != 0 || owner.applies != 0 {
		t.Fatalf("self-symlink checkout reached the owner: plans=%d applies=%d", owner.plans, owner.applies)
	}
}

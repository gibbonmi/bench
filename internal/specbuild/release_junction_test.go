package specbuild

import (
	"context"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/worktree"
)

// realReleaseOwner delegates Release to the real worktree.ReleaseProvisional,
// mirroring realOwner.Create's real-delegation pattern for the release path.
type realReleaseOwner struct{ realOwner }

func (realReleaseOwner) Release(_ context.Context, root, request, path string, evidence ReleaseEvidence) error {
	return worktree.ReleaseProvisional(root, request, path, worktree.ProvisionalEvidence{
		Base: evidence.Base, CheckpointRef: evidence.CheckpointRef, Checkpoint: evidence.Checkpoint, IntegratedRef: evidence.IntegratedRef, Integrated: evidence.Integrated,
	})
}

// mismatchedReleaseOwner forwards a request that never matches the real
// assignment, so the real worktree.ReleaseProvisional refuses instead of releasing.
type mismatchedReleaseOwner struct {
	realOwner
	request string
}

func (o mismatchedReleaseOwner) Release(_ context.Context, root, _, path string, evidence ReleaseEvidence) error {
	return worktree.ReleaseProvisional(root, o.request, path, worktree.ProvisionalEvidence{
		Base: evidence.Base, CheckpointRef: evidence.CheckpointRef, Checkpoint: evidence.Checkpoint, IntegratedRef: evidence.IntegratedRef, Integrated: evidence.Integrated,
	})
}

func TestIntegrateReleasesThroughRealProvisionalRelease(t *testing.T) {
	fixture := checkpointedReleaseFixture(t)
	fixture.service.worktrees = realReleaseOwner{}
	if _, err := fixture.service.Integrate(context.Background(), "build demo", fixture.assigned.ID); err != nil {
		t.Fatalf("Integrate: %v", err)
	}
	requireReleased(t, fixture)
}

func TestIntegrateSurfacesRealProvisionalReleaseMismatch(t *testing.T) {
	fixture := checkpointedReleaseFixture(t)
	fixture.service.worktrees = mismatchedReleaseOwner{request: "an unregistered request"}
	_, err := fixture.service.Integrate(context.Background(), "build demo", fixture.assigned.ID)
	if err == nil || !strings.Contains(err.Error(), "provisional release request, assignment, or path mismatch; checkout retained") {
		t.Fatalf("Integrate error = %v", err)
	}
}

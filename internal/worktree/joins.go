package worktree

import (
	"context"
	"io"
	"os"

	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// joins is this package's injectable seam set. A verb resolves one value at its boundary
// with defaultJoins and passes it down; nothing below the boundary reads a package
// variable. One value carries every family because the families travel together: a
// landing releases its assignment, a release cleans the worktree, and that cleanup reads
// the live-binary identity. A test builds the value, replaces one field, and calls the
// internal form, so two tests never share one stub point.
//
// A field whose default has to reach another seam takes the value as its first argument.
// The default then reads the caller's joins rather than a captured copy.
type joins struct {
	landReviewed func(context.Context, landing.ReviewedRequest, landing.Admission) (landing.ReviewedResult, error)
	landStack    func(context.Context, landing.StackRequest, landing.StackAdmission) (landing.StackResult, error)
	// publicationGap runs under the publication lock, after the final admission decision
	// and before the ref update. It receives the destination root.
	publicationGap func(string)
	// reconcileCommitment runs after the publication. It releases the local claims that
	// the published delivery facts satisfy. It receives the destination store, which reads
	// through the verb's run.
	reconcileCommitment      func(commitrepo.Store) error
	releaseLandingAssignment func(joins, ambient, string, []string, io.Writer, io.Writer) int
	// cleanupBoundary is the deterministic transaction fault seam. A nil value carries no
	// fault, exactly as hit reads it, so it is also the default.
	cleanupBoundary      Fault
	cleanupLockAttempt   func(string)
	creationLockAttempt  func(string)
	claimTakeoverGap     func(string)
	claimStealGap        func(string)
	restoreClean         func(string)
	chmodPool            func(string, os.FileMode) error
	resolveRunningBinary func() (string, error)
	reauthorizeLock      func(string, string, string) error
	reauthorizeBeforeCAS func(*intent.Assignment)
	// buildSubject authors a checkout's own published executable, under the manifest
	// directory the build owner defaults to, so the manifest lands beside the wrapper. It
	// is a seam because the default reaches a build script and a Go toolchain, and a
	// fixture for that pair would make every output row wait on a real compile.
	buildSubject func(context.Context, string, string) error
	// heldInstall is the held-install predicate that the post-landing refresh reads through
	// the destination store, which reads through the verb's run. It is a seam so a landing
	// test can hold the refresh without an installed candidate.
	heldInstall func(commitrepo.Store) bool
	// publishCandidate is the install's publication owner. Its operands are the seal
	// root, the staged executable, the published executable, the manifest directory, and
	// the version.
	publishCandidate func(root, staged, executable, manifestDir, version string) error
	// doctorChild runs the doctor repair that follows an install, with argv in dir under
	// env, and answers its exit code. It is a seam because the default starts a process.
	doctorChild   func(dir string, argv, env []string, stderr io.Writer) int
	resetLayers   func(string, recoveryManifest) error
	resetEnvelope func(string, []string, resetPlan) (intent.Recovery, error)
}

// defaultJoins names the real function behind every seam. It is the one place a default
// lives, so a verb's boundary and a test's starting value cannot disagree.
func defaultJoins() joins {
	return joins{
		landReviewed: func(ctx context.Context, request landing.ReviewedRequest, admission landing.Admission) (landing.ReviewedResult, error) {
			return landing.New().LandAdmitted(ctx, request, admission)
		},
		landStack: func(ctx context.Context, request landing.StackRequest, admission landing.StackAdmission) (landing.StackResult, error) {
			return landing.New().LandStack(ctx, request, admission)
		},
		publicationGap:           func(string) {},
		reconcileCommitment:      commitrepo.Store.ReconcileDelivered,
		releaseLandingAssignment: releaseCommandWith,
		cleanupLockAttempt:       func(string) {},
		creationLockAttempt:      func(string) {},
		claimTakeoverGap:         func(string) {},
		claimStealGap:            func(string) {},
		restoreClean:             restoreCleanCheckout,
		chmodPool:                os.Chmod,
		resolveRunningBinary:     os.Executable,
		reauthorizeLock:          lockWorktree,
		resetLayers:              restoreResetLayers,
		resetEnvelope:            writeResetEnvelope,
		buildSubject:             runbinary.BuildSubject,
		heldInstall:              readHeldInstall,
		publishCandidate:         freshness.Publish,
		doctorChild:              runDoctorChild,
	}
}

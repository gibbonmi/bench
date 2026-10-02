package worktree

import (
	"context"
	"io"
	"os"

	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// joins is this package's injectable seam set. A verb resolves one value at its boundary
// with defaultJoins and passes it down; nothing below the boundary reads a package
// variable. One value carries every family because the families travel together: a
// landing releases its assignment, a release cleans the worktree, and that cleanup reads
// the ignored inventory and the live-binary identity. A test builds the value, replaces
// one field, and calls the internal form, so two tests never share one stub point.
//
// A field whose default has to reach another seam takes the value as its first argument.
// The default then reads the caller's joins rather than a captured copy.
type joins struct {
	landReviewed             func(context.Context, landing.ReviewedRequest) (landing.ReviewedResult, error)
	advanceLandingMarker     func(context.Context, string, string, string, string) error
	reconcileLanding         func(joins, string, string, string, string) error
	releaseLandingAssignment func(joins, ambient, string, []string, io.Writer, io.Writer) int
	// pruneLandedBranches retires every unclaimed assignment branch the landed
	// default branch already carries. It runs at the landing because that is the
	// one moment the proof is cheap and certain: a later commit can move or delete
	// the files a folded sibling touched, and the content proof then fails forever.
	pruneLandedBranches    func(string) (int, error)
	authorizeLandingSource func(string, string, string) (diff.SourceRange, error)
	// cleanupBoundary is the deterministic transaction fault seam. A nil value carries no
	// fault, exactly as hit reads it, so it is also the default.
	cleanupBoundary      Fault
	cleanupLockAttempt   func(string)
	creationLockAttempt  func(string)
	claimTakeoverGap     func(string)
	claimStealGap        func(string)
	restoreClean         func(string)
	chmodPool            func(string, os.FileMode) error
	ignoredLstat         func(string) (os.FileInfo, error)
	resolveRunningBinary func() (string, error)
	liveBinaryWarnings   io.Writer
	// planLandedExplicit exposes the target-path boundary so a hostile-path test can stand
	// in for the planner without the fixture the real one needs.
	planLandedExplicit   func(joins, string, string, CleanupOptions) (CleanupPlan, error)
	reauthorizeUnlock    func(string, string) error
	reauthorizeLock      func(string, string, string) error
	reauthorizeBeforeCAS func(*intent.Assignment)
	// buildSubject authors a checkout's own published executable, under the manifest
	// directory the build owner defaults to, so the manifest lands beside the wrapper. It
	// is a seam because the default reaches a build script and a Go toolchain, and a
	// fixture for that pair would make every output row wait on a real compile.
	buildSubject func(context.Context, string, string) error
	// mergeReconcile is the merge verb's publication boundary: the checkout catch-up that
	// runs after the branch ref moved. It is a seam because its failure is the one
	// outcome that reads apart from a refusal, and no fixture can make a bare reset fail.
	mergeReconcile func(string, string) error
	// resetMove attaches the assignment branch and replaces the checkout at the checkpoint.
	resetMove     func(string, string, string) error
	resetLayers   func(string, recoveryManifest) error
	resetEnvelope func(string, resetPlan) (intent.Recovery, error)
}

// defaultJoins names the real function behind every seam. It is the one place a default
// lives, so a verb's boundary and a test's starting value cannot disagree.
func defaultJoins() joins {
	return joins{
		landReviewed: func(ctx context.Context, request landing.ReviewedRequest) (landing.ReviewedResult, error) {
			return landing.New().LandReviewed(ctx, request)
		},
		advanceLandingMarker:     authorization.AdvanceMarker,
		reconcileLanding:         reconcileLandingDestination,
		releaseLandingAssignment: releaseCommandWith,
		pruneLandedBranches:      intent.PruneUnclaimedLandedBranches,
		authorizeLandingSource:   preflight.AuthorizeReviewedSource,
		cleanupLockAttempt:       func(string) {},
		creationLockAttempt:      func(string) {},
		claimTakeoverGap:         func(string) {},
		claimStealGap:            func(string) {},
		restoreClean:             restoreCleanCheckout,
		chmodPool:                os.Chmod,
		ignoredLstat:             os.Lstat,
		resolveRunningBinary:     os.Executable,
		liveBinaryWarnings:       os.Stderr,
		planLandedExplicit:       planExplicitWith,
		reauthorizeUnlock:        unlockWorktree,
		reauthorizeLock:          lockWorktree,
		mergeReconcile:           reconcileMergeCheckout,
		resetMove:                moveResetCheckout,
		resetLayers:              restoreResetLayers,
		resetEnvelope:            writeResetEnvelope,
		buildSubject:             runbinary.BuildSubject,
	}
}

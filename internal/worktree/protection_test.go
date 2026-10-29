package worktree

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// protectCheckout registers one pending user of the checkout at dir. The test process is the
// user's live actor, so the user stays unresolved until the test finishes it.
func protectCheckout(t *testing.T, dir string) (protection.Descriptor, *protection.Registration) {
	t.Helper()
	p := checkoutObserver
	d, err := Descriptor(dir)
	requireTest(t, err == nil, "describe the checkout %s: %v", dir, err)
	t.Cleanup(func() { _ = os.RemoveAll(d.Store) })
	opened, err := p.OpenResource(d)
	requireTest(t, err == nil, "open the checkout generation: %v", err)
	user, err := p.Prepare([]protection.Descriptor{opened})
	requireTest(t, err == nil, "register the checkout user: %v", err)
	return d, user
}

// TestCheckoutProtectionPredicate covers the protected-user predicate. An absent store and a
// retired generation name no user. A live pending user keeps the checkout and the sealed
// generation, and a launch after the user resolves recovers the seal. A checkout that the
// predicate cannot describe, or a store that it cannot lock, keeps the checkout too. The
// denied store keeps its directory and generation readable, because the store root is
// shared and pool reclaim reads every worktree generation in it.
func TestCheckoutProtectionPredicate(t *testing.T) {
	t.Parallel()
	run, root := new(git.Run), newWorktreeRepo(t)
	requireTest(t, checkoutUnprotected(run, root) == nil, "an absent store kept the checkout")
	err := checkoutUnprotected(run, t.TempDir())
	requireTest(t, errors.Is(err, protection.ErrRetained), "an undescribed checkout = %v, want retained", err)

	d, user := protectCheckout(t, root)
	err = checkoutUnprotected(run, root)
	requireTest(t, errors.Is(err, protection.ErrRetained) && strings.Contains(err.Error(), string(protection.Protected)), "a live user = %v, want a protected retention", err)
	lock := filepath.Join(d.Store, "lock")
	mustChmod(t, lock, 0)
	err = checkoutUnprotected(run, root)
	mustChmod(t, lock, 0o600)
	requireTest(t, errors.Is(err, protection.ErrRetained) && errors.Is(err, fs.ErrPermission), "a denied store = %v, want a retention with its cause", err)
	_, err = checkoutObserver.Reopen(d)
	requireTest(t, errors.Is(err, protection.ErrSealed), "a launch beside the live user = %v, want the sealed generation", err)

	requireTest(t, user.Finish(protection.Terminal{NoStart: true}) == nil, "finish the checkout user")
	_, err = checkoutObserver.Reopen(d)
	requireTest(t, err == nil, "a launch after the user resolved = %v, want the recovered generation", err)
	requireTest(t, checkoutUnprotected(run, root) == nil, "an unused generation kept the checkout")
	requireTest(t, checkoutUnprotected(run, root) == nil, "a retired generation kept the checkout")
}

// TestCleanupApplyRechecksProtection covers the release apply. A user of the checkout keeps
// it under the cleanup exclusion, so the release refuses and the checkout survives.
func TestCleanupApplyRechecksProtection(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "protected-release")
	protectCheckout(t, f.creation.Path)
	release := runVerb(t, verbRelease, f.call("--request", "landed-protected-release", f.creation.Path))
	requireTest(t, release.exit == 1 && strings.Contains(release.stderr, protection.ErrRetained.Error()), "protected release exit=%d stderr=%q, want a retention", release.exit, release.stderr)
	_, err := os.Stat(f.creation.Path)
	requireTest(t, err == nil, "the protected checkout is gone: %v", err)
}

// TestCleanupPlanFingerprintsProtection covers the plan evidence. A user that joins after the
// plan makes the plan's fingerprint stale, so apply refuses it before any other decision.
func TestCleanupPlanFingerprintsProtection(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "protected-plan")
	plan, err := PlanExplicit(f.root, f.creation.Path)
	mustNoError(t, err)
	protectCheckout(t, f.creation.Path)
	_, err = ApplyExplicit(f.root, f.creation.Path, plan.Fingerprint)
	requireTest(t, errors.Is(err, errStaleFingerprint), "apply after a user joined = %v, want the stale fingerprint", err)
}

// TestPoolKeepsProtectedCheckout covers pool release and reuse. A release leaves a protected
// checkout leased and unreset, and an acquire passes over it and returns its claim.
func TestPoolKeepsProtectedCheckout(t *testing.T) {
	t.Parallel()
	root, home := newWorktreeRepo(t), t.TempDir()
	wt, err := acquireAt(defaultJoins(), nil, root, "", "", home, currentTime())
	mustNoError(t, err)
	lease, err := LeaseFile(wt)
	mustNoError(t, err)
	sentinel := filepath.Join(wt, "sentinel")
	mustWrite(t, sentinel, []byte("in use\n"), 0o644)
	protectCheckout(t, wt)

	releaseWith(defaultJoins(), nil, wt)
	_, sentinelErr := os.Stat(sentinel)
	_, leaseErr := os.Stat(lease)
	requireTest(t, sentinelErr == nil && leaseErr == nil, "protected release left sentinel %v and lease %v, want both", sentinelErr, leaseErr)

	mustRemove(t, sentinel)
	mustRemove(t, lease)
	// A later instant mints a fresh name, because the first mint holds this second's name.
	next, err := acquireAt(defaultJoins(), nil, root, "", "", home, currentTime().Add(time.Second))
	mustNoError(t, err)
	_, leaseErr = os.Stat(lease)
	requireTest(t, next != wt && os.IsNotExist(leaseErr), "acquire = %s with protected lease %v, want another checkout and no claim", next, leaseErr)
}

// TestProtectedReclaim covers pool reclaim. The repository of a pool member is gone, and a
// user that its checkout registered is unresolved, so the plan retains the member's key and
// names the users. Once the user resolves, the plan reclaims the key and the apply removes it.
func TestProtectedReclaim(t *testing.T) {
	t.Parallel()
	f := newReclaimPool(t)
	source := newWorktreeRepo(t)
	key := filepath.Base(Pool(canonicalRoot(source)))
	created := mustCreate(t, source, f.home, "protected-reclaim", "protrec")
	_, user := protectCheckout(t, created.Path)
	mustNoError(t, os.RemoveAll(source))

	plan := runVerb(t, verbReclaim, f.call())
	_, verdicts := reclaimVerdicts(t, plan.stdout)
	requireTest(t, plan.exit == 0 && verdicts[key][0] == poolVerdictRetain && strings.Contains(verdicts[key][1], "users remain unresolved"),
		"plan with a live member user = %d %q/%q, want a retention that names the users", plan.exit, verdicts[key][0], verdicts[key][1])

	requireTest(t, user.Finish(protection.Terminal{NoStart: true}) == nil, "finish the member user")
	plan = runVerb(t, verbReclaim, f.call())
	fingerprint := plan.mustFingerprint(t)
	_, verdicts = reclaimVerdicts(t, plan.stdout)
	requireTest(t, verdicts[key][0] == poolVerdictReclaim, "plan after the user resolved = %q/%q, want reclaim", verdicts[key][0], verdicts[key][1])
	out := runVerb(t, verbReclaim, f.call("--apply", fingerprint))
	_, err := os.Lstat(filepath.Join(f.pool, key))
	requireTest(t, out.exit == 0 && os.IsNotExist(err), "apply code=%d out=%q, key %v, want the key removed", out.exit, out.stdout, err)
}

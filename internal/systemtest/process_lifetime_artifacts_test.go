//go:build system

package systemtest

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/gate/prospectiveartifact"
	"github.com/gibbonmi/bench/internal/gocache"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// TestProcessLifetimeArtifactsSecondProcessSweep is PL59. The landing CLI dies while its
// oracle and phase groups live on as registered bundle users. A second process's sweep must
// retain the bundle and its registration, because owner death is not user proof.
func TestProcessLifetimeArtifactsSecondProcessSweep(t *testing.T) {
	t.Parallel()
	root, home, tally, trees, ready, release := systemLandingRaceFixture(t)
	configureArtifactLandingFixture(t, root)
	private := t.TempDir()
	base := systemGitOutput(t, root, "rev-parse", "main")
	cmd, out, errOut, done := startArtifactAuthorization(t, root, home, tally, trees, ready, release, private, base, "artifact-dead", "loser.txt", "dead prospective subject")
	awaitArtifactBarrier(t, ready, done, out, errOut)
	oracle := artifactOracle(t, home)
	// The released phase lets the surviving groups end after the row has its answer.
	t.Cleanup(func() { releaseArtifactOracle(t, release, oracle) })
	bundle := oneProspectiveBundle(t, private)
	checkout := filepath.Join(bundle, prospectiveartifact.CheckoutName)
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
		t.Fatal(err)
	}
	<-done

	// A second landing sweeps before it opens its own bundle. It waits on its own barrier
	// pair, so the dead landing's phase stays alive through the sweep.
	sweepReady, sweepRelease := landingBarrier(t, t.TempDir())
	sweeper, sweepOut, sweepErr, sweepDone := startArtifactAuthorization(t, root, home, tally, trees, sweepReady, sweepRelease, private, base, "artifact-sweep", "winner.txt", "sweeping prospective subject")
	awaitArtifactBarrier(t, sweepReady, sweepDone, sweepOut, sweepErr)
	t.Cleanup(func() {
		releaseArtifactBarrier(t, sweepRelease)
		exitedWithin(t, sweeper, sweepDone, artifactBarrierWindow)
	})
	if _, err := os.Stat(checkout); err != nil {
		t.Fatalf("checkout of the dead owner with live users = %v, want retained", err)
	}
	requireCheckoutRegistered(t, root, checkout, true)
}

// TestProcessLifetimeArtifactsCompletedOwners is PL118. A green landing completes every
// resource user, so its owners release the bundle with its selected executable, the
// checkout registration, and the cache generation. A blanket retention leaves one behind.
func TestProcessLifetimeArtifactsCompletedOwners(t *testing.T) {
	t.Parallel()
	root, home, tally, trees, ready, release := systemLandingRaceFixture(t)
	configureArtifactLandingFixture(t, root)
	private := t.TempDir()
	base := systemGitOutput(t, root, "rev-parse", "main")
	cmd, out, errOut, done := startArtifactAuthorization(t, root, home, tally, trees, ready, release, private, base, "artifact-complete", "winner.txt", "complete prospective subject")
	awaitArtifactBarrier(t, ready, done, out, errOut)
	bundle := oneProspectiveBundle(t, private)
	checkout := filepath.Join(bundle, prospectiveartifact.CheckoutName)
	binary := oneOwnerBinary(t, bundle)
	requireCheckoutRegistered(t, root, checkout, true)
	releaseArtifactBarrier(t, release)
	if code := exitedWithin(t, cmd, done, artifactBarrierWindow); code != 0 {
		t.Fatalf("green landing = (%d, %q, %q), want exit 0", code, out.String(), errOut.String())
	}
	for _, path := range []string{bundle, binary} {
		if _, err := os.Lstat(path); !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("%s after the green landing = %v, want released", path, err)
		}
	}
	requireCheckoutRegistered(t, root, checkout, false)
	dir, err := gocache.Dir([]string{"HOME=" + home})
	if err != nil {
		t.Fatal(err)
	}
	store, err := gocache.ProtectionStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	// The owner retires its own generation, so recovery finds no generation to grade.
	protocol := &protection.Protocol{Validators: protection.Validators{protection.KindCache: gocache.ValidateProtection}, ReadLimit: bounds.ControlRecordLimit}
	if disposition, err := protocol.Recover(store, nil); !protection.NoGeneration(disposition, err) {
		t.Fatalf("cache protection after the green landing = %+v (%v), want a retired generation", disposition, err)
	}
}

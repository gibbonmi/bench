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
	"github.com/gibbonmi/bench/internal/worktree"
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
// checkout registration, and the cache generation. A released worktree retires the generation
// that a completed exec user left open, and a completed kit-root focused run removes its
// private test directory. A blanket retention leaves one of them behind.
func TestProcessLifetimeArtifactsCompletedOwners(t *testing.T) {
	t.Parallel()
	root, home, tally, trees, ready, release := systemLandingRaceFixture(t)
	configureArtifactLandingFixture(t, root)
	private := t.TempDir()
	base := systemGitOutput(t, root, "rev-parse", "main")
	used := systemCreateLandingWorktree(t, root, home, "artifact-used", "artifact used")
	if exec := owner.runAt(root, artifactChildEnv(root, home, tally, trees, ready, release, private), owner.selected.path, execBoundArgs(used, "true")...); exec.code != 0 {
		t.Fatalf("exec in the used worktree = (%d, %q, %q)", exec.code, exec.stdout, exec.stderr)
	}
	checkoutStore, err := worktree.Descriptor(used.path)
	if err != nil {
		t.Fatal(err)
	}
	// The release retires the generation that the completed exec user left open.
	if released := owner.runAt(root, artifactChildEnv(root, home, tally, trees, ready, release, private), owner.selected.path, "worktree", "release", "--request", used.request, used.path); released.code != 0 {
		t.Fatalf("release of the used worktree = (%d, %q, %q)", released.code, released.stdout, released.stderr)
	}
	checkouts := &protection.Protocol{Validators: protection.Validators{protection.KindWorktree: worktree.ValidateProtection}, ReadLimit: bounds.ControlRecordLimit}
	if disposition, err := checkouts.Recover(checkoutStore.Store, nil); (err != nil && !protection.NoGeneration(disposition, err)) || disposition.Resources[0].ResourceID != "" {
		t.Fatalf("worktree protection after the released exec user = %+v (%v), want a retired generation", disposition, err)
	}
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
	// A focused run at the kit root opens a private test directory, and its go test child
	// names that run's TMPDIR. The kit root builds its own private executable, as PL93 does.
	tools, mark := t.TempDir(), filepath.Join(t.TempDir(), "run")
	fakeGo(t, tools, map[string]string{goTestPattern: "printf '%s' \"$TMPDIR\" > '" + mark + "'\nprintf '%s' '" + goEvents([]string{"landingrace"}, "landingrace") + "'\n"})
	focused := systemSelected(t, root, kitRootEnv(home, root, tools, t.TempDir()), "test", "--package", "./...")
	run, err := os.ReadFile(mark)
	if focused.code != 0 || err != nil {
		t.Fatalf("focused run at the kit root = (%d, %q, %q), TMPDIR %v; want a completed private run", focused.code, focused.stdout, focused.stderr, err)
	}
	if _, err := os.Lstat(filepath.Dir(string(run))); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("private test directory %s after the completed run = %v, want released", filepath.Dir(string(run)), err)
	}
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

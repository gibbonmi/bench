package protection_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// recovered recovers d's store as its owner would from disk, and reports whether release ran.
func recovered(t *testing.T, p *protection.Protocol, d protection.Descriptor) (protection.Disposition, bool) {
	t.Helper()
	released := false
	disposition, _ := p.Recover(d.Store, func() error { released = true; return nil })
	return disposition, released
}

// requireRecoveredProtected fails unless recovery retains d for exactly one live user.
func requireRecoveredProtected(t *testing.T, p *protection.Protocol, d protection.Descriptor) {
	t.Helper()
	if disposition, released := recovered(t, p, d); disposition.State != protection.Protected || released || len(disposition.Users) != 1 {
		t.Fatalf("recovery: %+v released=%t, want one protecting user", disposition, released)
	}
}

// TestRetiredStoreHasNoGeneration holds that a retired store keeps its retired record and
// still reads as a store with no generation.
func TestRetiredStoreHasNoGeneration(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	if disposition, _ := observe(t, p, d); disposition.State != protection.Safe {
		t.Fatalf("retirement of an unused resource = %+v, want safe", disposition)
	}
	if _, err := os.Lstat(filepath.Join(d.Store, protection.RetiredFile)); err != nil {
		t.Fatalf("the retired store holds no retired record: %v", err)
	}
	if disposition, err := p.Recover(d.Store, nil); !protection.NoGeneration(disposition, err) {
		t.Fatalf("recovery of a retired store = %+v (%v), want no generation", disposition, err)
	}
}

// TestRecoveryPendingActorAbsent is PL25. A start gate opens only after the actor publishes the
// known identity, so a pending user whose actor is absent never ran and resolves as no-start. A
// pending user whose actor answers stays protected, because its PID can be a reused one.
func TestRecoveryPendingActorAbsent(t *testing.T) {
	p := newProtocol()
	absent, answering := newResource(t, p, "cache"), newResource(t, p, "cache")
	deadUser(t, absent)
	if disposition, released := recovered(t, p, absent); disposition.State != protection.Safe || !released {
		t.Fatalf("pending user of an absent actor: %+v released=%t, want no-start", disposition, released)
	}
	pendingUser(t, p, answering)
	requireRecoveredProtected(t, p, answering)
	for _, record := range userRecords(t, answering) {
		if record["identity"] != "pending" {
			t.Fatalf("the answering actor's startup record became %v", record)
		}
	}
}

// pendingUser registers one pending user of d. This test process is its live actor.
func pendingUser(t *testing.T, p *protection.Protocol, d protection.Descriptor) *protection.Registration {
	t.Helper()
	reg, err := p.Prepare([]protection.Descriptor{d})
	if err != nil {
		t.Fatal(err)
	}
	return reg
}

// deadUser records one safe user of d that keeps its record: a pending user whose actor, a
// nested CLI, exited before its start gate opened.
func deadUser(t *testing.T, d protection.Descriptor) {
	t.Helper()
	if err := processtest.Helper(t, "nested", vectorEnv(d), stageEnv+"=prepared").Run(); err != nil {
		t.Fatal(err)
	}
}

// requireOpenAfterRelease fails unless the release observation of d beside one pending user
// returns state, runs no release, and leaves admission open for the next Prepare.
func requireOpenAfterRelease(t *testing.T, p *protection.Protocol, d protection.Descriptor, state protection.State) {
	t.Helper()
	pendingUser(t, p, d)
	if disposition, released, err := observeRelease(t, p, d); err != nil || disposition.State != state || released {
		t.Fatalf("release observation = %+v (%v) released=%t, want %s with no release", disposition, err, released, state)
	}
	if _, err := p.Prepare([]protection.Descriptor{d}); err != nil {
		t.Fatalf("Prepare after the release observation = %v, want admission", err)
	}
}

// TestReleaseObservationLeavesAdmissionOpen is SP2. A release beside a protected user
// writes no record, so the generation stays open.
func TestReleaseObservationLeavesAdmissionOpen(t *testing.T) {
	p := newProtocol()
	requireOpenAfterRelease(t, p, newResource(t, p, "cache"), protection.Protected)
}

// TestReleaseObservationLeavesAdmissionOpenWhenUncertain is SP5. A probe that answers EPERM
// cannot decide the user, and the release deletes nothing, so it seals nothing either.
func TestReleaseObservationLeavesAdmissionOpenWhenUncertain(t *testing.T) {
	p := newProtocol()
	protection.SetProbe(p, func(int, syscall.Signal) error { return syscall.EPERM })
	requireOpenAfterRelease(t, p, newResource(t, p, "cache"), protection.Uncertain)
}

// TestReleaseObservationNamesTheLiveUser is SP3. The release returns the retention and the
// user that holds it.
func TestReleaseObservationNamesTheLiveUser(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	reg := pendingUser(t, p, d)
	if disposition, _, err := observeRelease(t, p, d); err != nil || disposition.State != protection.Protected || !slices.Equal(disposition.Users, reg.Disposition().Users) {
		t.Fatalf("release observation = %+v (%v), want protected by %v", disposition, err, reg.Disposition().Users)
	}
}

// TestReleaseObservationRetiresASafeGeneration is SP4. A release whose every user is safe
// runs release, removes each user, and leaves the retired record and no generation.
func TestReleaseObservationRetiresASafeGeneration(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	deadUser(t, d)
	if disposition, released, err := observeRelease(t, p, d); err != nil || disposition.State != protection.Safe || !released {
		t.Fatalf("release observation of a safe generation = %+v (%v) released=%t, want safe with a release", disposition, err, released)
	}
	if _, err := protection.ReadRetired(d.Store, readLimit); err != nil {
		t.Fatalf("the released store holds no retired record: %v", err)
	}
	if _, err := os.Lstat(filepath.Join(d.Store, protection.ResourceFile)); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("generation after the release: %v, want it absent", err)
	}
	if users := userRecords(t, d); len(users) != 0 {
		t.Fatalf("users after the release: %v, want none", users)
	}
}

// TestReleaseObservationKeepsAnEarlierSeal is SP6. A clean's seal-first observation beside a
// live user seals the generation, and a later release never opens it again.
func TestReleaseObservationKeepsAnEarlierSeal(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	pendingUser(t, p, d)
	requireRetained(t, p, d)
	if disposition, _, err := observeRelease(t, p, d); err != nil || disposition.State != protection.Protected {
		t.Fatalf("release observation after the seal = %+v (%v), want protected", disposition, err)
	}
	if _, err := p.Prepare([]protection.Descriptor{d}); !errors.Is(err, protection.ErrSealed) {
		t.Fatalf("Prepare after the release observation = %v, want %v", err, protection.ErrSealed)
	}
}

// TestReleaseObservationSealFailure is SP7. On the safe path the seal goes first, so a failed
// seal publish returns uncertain, runs no release, and removes no user.
func TestReleaseObservationSealFailure(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	deadUser(t, d)
	faults := router(map[durablefiletest.Fault]map[string]error{{Destination: "C/" + protection.ResourceFile, Ordinal: 1}: {"rename": errInjected}}, map[string]string{d.Store: "C"})
	protection.SetReplacer(p, faults.Replacer())
	if disposition, released, err := observeRelease(t, p, d); !errors.Is(err, errInjected) || disposition.State != protection.Uncertain || released {
		t.Fatalf("release observation with a failed seal = %+v (%v) released=%t, want uncertain with no release", disposition, err, released)
	}
	if users := userRecords(t, d); len(users) != 1 {
		t.Fatalf("users after the failed seal: %v, want the one user", users)
	}
}

func TestRecoveryOverlappingUsers(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	first, second := startUser(t, p, d), startUser(t, p, d)
	first.stop()
	if err := first.reg.Finish(complete); err != nil {
		t.Fatal(err)
	}
	if disposition := requireRetained(t, p, d); !slices.Equal(disposition.Users, second.reg.Disposition().Users) {
		t.Fatalf("after the first completion: %+v, want the second user to hold it", disposition)
	}
	second.stop()
	if err := second.reg.Finish(complete); err != nil {
		t.Fatal(err)
	}
	if disposition, released := observe(t, p, d); disposition.State != protection.Safe || !released {
		t.Fatalf("after both completions: %+v released=%t", disposition, released)
	}
}

func TestRecoveryMissingCompletion(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	pid := filepath.Join(t.TempDir(), "descendant")
	o := run(t, processtest.Helper(t, "nested", stageEnv+"=started", pidEnv+"="+pid), p, []protection.Descriptor{d})
	processtest.WaitPID(t, pid)
	if !o.Process.Waited || o.Process.Code != 0 || o.Protection.State != protection.Safe {
		t.Fatalf("outer launch: process=%+v protection=%+v", o.Process, o.Protection)
	}
	// The nested CLI exited zero without a terminal record while its own group lives on.
	if disposition := requireRetained(t, p, d); disposition.State != protection.Protected || len(disposition.Users) != 1 {
		t.Fatalf("missing nested completion: %+v", disposition)
	}
}

func TestRecoveryRecoverKnown(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	u := startUser(t, p, d)
	// The live owner records complete facts while its group still answers.
	if err := u.reg.Finish(complete); err != nil {
		t.Fatal(err)
	}
	requireRecoveredProtected(t, p, d)
}

func TestRecoveryNoStaleSignal(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	log := filepath.Join(t.TempDir(), "signals")
	trap := `for s in TERM INT HUP USR1 USR2; do trap "echo $s >> \"$0\"" $s; done; while :; do sleep 0.05; done`
	u := startUser(t, p, d, "sh", "-c", trap, log)
	if err := u.reg.Finish(complete); err != nil {
		t.Fatal(err)
	}
	var sent []syscall.Signal
	protection.SetProbe(p, func(target int, signal syscall.Signal) error {
		sent = append(sent, signal)
		return syscall.Kill(target, signal)
	})
	requireRecoveredProtected(t, p, d)
	signals, _ := os.ReadFile(log)
	if len(sent) == 0 || slices.ContainsFunc(sent, func(s syscall.Signal) bool { return s != 0 }) || len(signals) != 0 {
		t.Fatalf("recovery sent %v and the recorded group logged %q, want only signal zero", sent, signals)
	}
	if err := syscall.Kill(u.cmd.Process.Pid, 0); err != nil {
		t.Fatalf("the recorded group is gone after recovery: %v", err)
	}
}

func TestRecoveryOldRecord(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	u := startUser(t, p, d)
	if err := u.reg.Finish(complete); err != nil {
		t.Fatal(err)
	}
	old := time.Now().AddDate(-10, 0, 0)
	err := filepath.WalkDir(d.Store, func(path string, _ os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		return os.Chtimes(path, old, old)
	})
	if err != nil {
		t.Fatal(err)
	}
	requireRecoveredProtected(t, p, d)
}

func TestRecoveryWorkerActorDeath(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	pid := filepath.Join(t.TempDir(), "descendant")
	if err := processtest.Helper(t, "nested", vectorEnv(d), stageEnv+"=workers", pidEnv+"="+pid).Run(); err != nil {
		t.Fatal(err)
	}
	descendant := processtest.WaitPID(t, pid)
	pgid, err := syscall.Getpgid(descendant)
	if err != nil {
		t.Fatal(err)
	}
	// The actor is gone, but its pending workers' group still answers.
	requireRecoveredProtected(t, p, d)
	_ = syscall.Kill(descendant, syscall.SIGKILL)
	processtest.WaitAbsent(t, pgid)
	if disposition, released := recovered(t, p, d); disposition.State != protection.Safe || !released {
		t.Fatalf("after the group's absence: %+v released=%t", disposition, released)
	}
}

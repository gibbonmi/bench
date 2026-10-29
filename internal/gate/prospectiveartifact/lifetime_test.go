package prospectiveartifact

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// TestMain installs the prospective row of the composition root's table, so every owner in
// this package protects its bundle as a CLI owner does. The package owns the only vector, so
// an outer run's resources stay out of its launches.
func TestMain(m *testing.M) {
	if err := os.Unsetenv(protection.Env); err != nil {
		os.Exit(3)
	}
	protection.Install(protection.Protocol{Validators: protection.Validators{protection.KindProspective: ValidateProtection}, ReadLimit: bounds.ControlRecordLimit})
	processtest.Main(m, nil)
}

// testBundle plants a bundle whose record names its open protection generation.
func testBundle(t *testing.T, base string, record Record) string {
	t.Helper()
	root, err := os.MkdirTemp(base, BundlePrefix)
	if err != nil {
		t.Fatal(err)
	}
	record.Protection = protectedGeneration(t, root)
	if err := Publish(root, record); err != nil {
		t.Fatal(err)
	}
	return root
}

// protectedGeneration opens the protection generation of a planted bundle root, as an owner does.
func protectedGeneration(t *testing.T, root string) string {
	t.Helper()
	generation, err := bundleDirectory.Open(protection.Installed(), root)
	if err != nil {
		t.Fatal(err)
	}
	return generation.Resource.ResourceID
}

// TestBundleValidatorRefusesAForeignBinding pins the bundle validator that the composition root
// installs: it accepts a bundle's own descriptor and refuses one whose binding names another
// directory at the same path.
func TestBundleValidatorRefusesAForeignBinding(t *testing.T) {
	d, err := bundleDirectory.Describe(plantedBundle(t, t.TempDir()))
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateProtection(d); err != nil {
		t.Fatalf("ValidateProtection(own bundle) = %v, want nil", err)
	}
	d.Binding = "0:0"
	if err := ValidateProtection(d); err == nil || !strings.Contains(err.Error(), "binding 0:0 does not name bundle") {
		t.Fatalf("ValidateProtection(foreign binding) = %v, want a binding refusal", err)
	}
}

// materializedOwner opens one protected owner under base and materializes the repository's tree.
func materializedOwner(t *testing.T, repository, base string) *Owner {
	t.Helper()
	owner, err := (Factory{TempRoot: base}).Open(repository)
	if err != nil {
		t.Fatal(err)
	}
	if err := owner.Materialize(gittest.Output(t, repository, "write-tree")); err != nil {
		t.Fatal(err)
	}
	return owner
}

// TestDirectClose is PL58. A live group is a registered user of the bundle, so Close must
// observe it and retain the checkout and its registration.
func TestDirectClose(t *testing.T) {
	repository := fixtureRepository(t)
	owner := materializedOwner(t, repository, t.TempDir())
	registration, err := protection.Installed().Prepare(owner.Resources())
	if err != nil {
		t.Fatal(err)
	}
	// The test's own group is a live group, so it stands for the registered user.
	if err := registration.Started(os.Getpid(), syscall.Getpgrp()); err != nil {
		t.Fatal(err)
	}

	if err := owner.Close(); err == nil || !strings.Contains(err.Error(), "stays retained") {
		t.Fatalf("Close with a live group user = %v, want the retention", err)
	}
	if _, err := os.Stat(owner.Checkout()); err != nil {
		t.Fatalf("checkout after the retained Close = %v, want it in place", err)
	}
	requireCheckoutRegistered(t, repository, owner.Checkout())
}

// TestLegacyRecord is PL60. A schema-1 record names this repository and a dead owner, but it
// records no user set, so the sweep must retain the bundle and its registration.
func TestLegacyRecord(t *testing.T) {
	repository := fixtureRepository(t)
	common, err := CanonicalCommonDir(repository)
	if err != nil {
		t.Fatal(err)
	}
	base := t.TempDir()
	legacy := plantedBundle(t, base)
	if err := Publish(legacy, Record{Schema: LegacySchema, OwnerPID: 42, CommonDir: common}); err != nil {
		t.Fatal(err)
	}
	gittest.MustRun(t, repository, "worktree", "add", "-q", "--detach", filepath.Join(legacy, CheckoutName), "HEAD")

	requireSweepRetains(t, repository, base, deadProbe, legacy)
	requireCheckoutRegistered(t, repository, filepath.Join(legacy, CheckoutName))
}

// TestEmptyCurrentBundle is PL75. A current bundle with no users is removable by its own
// Close and, after its owner dies, by the next sweep.
func TestEmptyCurrentBundle(t *testing.T) {
	repository := fixtureRepository(t)
	base := t.TempDir()
	closed := materializedOwner(t, repository, base)
	if err := closed.Close(); err != nil {
		t.Fatalf("Close of an unused bundle = %v, want removal", err)
	}
	dead := materializedOwner(t, repository, base)

	sweeper, err := (Factory{TempRoot: base, Probe: deadProbe}).Open(repository)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sweeper.Close() }()
	for _, owner := range []*Owner{closed, dead} {
		if _, err := os.Stat(owner.Root()); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("unused bundle %s = %v, want absent", owner.Root(), err)
		}
		if _, err := os.Lstat(owner.generation.Resource.Store); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("protection store of %s = %v, want absent", owner.Root(), err)
		}
		requireCheckoutUnregistered(t, repository, owner.Checkout())
	}
}

// TestLiveUnusedOwner is PL121. A live owner can hold a checkout before any child registers,
// so a sweep must retain it although its generation has no users.
func TestLiveUnusedOwner(t *testing.T) {
	repository := fixtureRepository(t)
	base := t.TempDir()
	live := materializedOwner(t, repository, base)
	defer func() { _ = live.Close() }()

	requireSweepRetains(t, repository, base, func(int) error { return nil }, live.Root())
	requireCheckoutRegistered(t, repository, live.Checkout())
}

func TestOpenRefusesRecoveryWhenRegistrationRemovalFails(t *testing.T) {
	repository := fixtureRepository(t)
	common, err := benchgit.CommonDir(repository)
	if err != nil {
		t.Fatal(err)
	}
	base := t.TempDir()
	dead := testBundle(t, base, Record{Schema: RecordSchema, OwnerPID: 42, CommonDir: common})
	gittest.MustRun(t, repository, "worktree", "add", "-q", "--detach", filepath.Join(dead, CheckoutName), "HEAD")
	refusing := Factory{TempRoot: base, Probe: deadProbe, Remove: func(string, string) error { return errors.New("registration refused") }}
	planted := snapshotPath(t, base)
	registrations := worktreeRegistrations(t, repository)
	// Observation seals the dead generation before registration removal, so the first
	// refusal writes inside the protection store alone. The store resides under the store
	// root, outside the temporary root. A repeated refusal writes nothing.
	if _, err := refusing.Open(repository); err == nil {
		t.Fatal("first recovery with a registration failure = nil, want refusal")
	}
	if sealed := snapshotPath(t, base); !reflect.DeepEqual(planted, sealed) {
		t.Fatalf("temporary root after the first refusal = %v, want the planted %v", sealed, planted)
	}
	if after := worktreeRegistrations(t, repository); after != registrations {
		t.Fatalf("registrations after the first refusal =\n%s\nwant\n%s", after, registrations)
	}
	planted = snapshotPath(t, base)

	_, err = refusing.Open(repository)
	if err == nil {
		t.Fatal("recovery with a registration failure = nil, want refusal")
	}
	if _, err := os.Stat(dead); err != nil {
		t.Fatalf("dead bundle after refusal = %v, want retained", err)
	}
	if after := snapshotPath(t, base); !reflect.DeepEqual(planted, after) {
		t.Fatalf("temporary root after refusal = %v, want the planted %v", after, planted)
	}
	if after := worktreeRegistrations(t, repository); after != registrations {
		t.Fatalf("registrations after refusal =\n%s\nwant\n%s", after, registrations)
	}
}

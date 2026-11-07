package protection_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// TestGenerationRemove is the one safe removal of an owner's resource. A retained removal keeps
// every store. A safe removal deletes the resource's store and the store of each nested resource
// that the removal deleted, whether its generation is live or retired. A nested resource that
// still exists, and a resource outside the removed directory whose path shares its prefix, keep
// their stores.
func TestGenerationRemove(t *testing.T) {
	protection.SetStoreBase(t, t.TempDir())
	outer := protection.Directory{Kind: protection.KindKitTest, Noun: "run"}
	p := &protection.Protocol{Validators: protection.Validators{protection.KindKitTest: outer.Validate, selection.Kind: selection.Validate}, ReadLimit: readLimit}
	parent := t.TempDir()
	open := func(adapter protection.Directory, dir string) protection.Generation {
		t.Helper()
		if err := os.Mkdir(dir, 0o700); err != nil {
			t.Fatal(err)
		}
		g, err := adapter.Open(p, dir)
		if err != nil {
			t.Fatal(err)
		}
		return g
	}
	run := open(outer, filepath.Join(parent, "run"))
	stores := map[string]protection.Descriptor{
		"run":             run.Resource,
		"gone":            open(selection, filepath.Join(parent, "run", "gone")).Resource,
		"kept":            open(selection, filepath.Join(parent, "run", "kept")).Resource,
		"sibling":         open(selection, filepath.Join(parent, "run2")).Resource,
		"retired gone":    retire(t, p, selection, filepath.Join(parent, "run", "retired-gone")),
		"retired kept":    retire(t, p, selection, filepath.Join(parent, "run", "retired-kept")),
		"retired sibling": retire(t, p, selection, filepath.Join(parent, "run3")),
	}
	// The retired sibling's resource is gone too, so only its path outside the run keeps it.
	if err := os.RemoveAll(stores["retired sibling"].ResourcePath); err != nil {
		t.Fatal(err)
	}
	requireStores := func(when string, want map[string]bool) {
		t.Helper()
		for name, d := range stores {
			_, err := os.Stat(d.Store)
			if present := err == nil; present != want[name] {
				t.Errorf("%s: store of %s present=%v (%v), want %v", when, name, present, err, want[name])
			}
		}
	}

	if err := run.Remove("run", func() error { return errInjected }, func() error { return os.RemoveAll(run.Resource.ResourcePath) }); !errors.Is(err, errInjected) {
		t.Fatalf("failed release = %v, want the release error", err)
	}
	requireStores("after a retained removal", map[string]bool{"run": true, "gone": true, "kept": true, "sibling": true, "retired gone": true, "retired kept": true, "retired sibling": true})

	// The removal deletes one live and one retired nested resource and leaves the others, so
	// only the deleted ones' stores go with the run's own.
	if err := run.Remove("run", nil, func() error {
		return errors.Join(os.RemoveAll(stores["gone"].ResourcePath), os.RemoveAll(stores["retired gone"].ResourcePath))
	}); err != nil {
		t.Fatal(err)
	}
	requireStores("after a safe removal", map[string]bool{"kept": true, "sibling": true, "retired kept": true, "retired sibling": true})
}

// TestRetiredSweepKeepsAUser holds that a retired store whose users directory holds an entry
// stays after a sweep, although its resource is gone.
func TestRetiredSweepKeepsAUser(t *testing.T) {
	p, d, dir := goneRetired(t)
	entry := filepath.Join(d.Store, protection.UsersDir, "user")
	if err := os.Mkdir(entry, 0o700); err != nil {
		t.Fatal(err)
	}
	requireSweep(t, p, dir, d, true)
	if err := os.Remove(entry); err != nil {
		t.Fatal(err)
	}
	requireSweep(t, p, dir, d, false)
}

// TestRetiredSweepSkipsALockedStore holds that a sweep does not wait on a store lock that
// another open file holds. The store stays, and the sweep returns no error for it.
func TestRetiredSweepSkipsALockedStore(t *testing.T) {
	p, d, dir := goneRetired(t)
	holder, err := os.OpenFile(filepath.Join(d.Store, protection.LockFile), os.O_RDWR, 0)
	if err == nil {
		err = syscall.Flock(int(holder.Fd()), syscall.LOCK_EX)
	}
	if err != nil {
		t.Fatal(err)
	}
	requireSweep(t, p, dir, d, true)
	if err := holder.Close(); err != nil {
		t.Fatal(err)
	}
	requireSweep(t, p, dir, d, false)
}

// TestRetiredStoreRecordsItsPath holds that a safe retirement leaves a retired record of the
// retired generation's kind, resource path, and binding, and leaves no generation.
func TestRetiredStoreRecordsItsPath(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	if disposition, released := observe(t, p, d); disposition.State != protection.Safe || !released {
		t.Fatalf("retirement of an unused resource = %+v released=%t, want safe", disposition, released)
	}
	want := protection.Descriptor{Kind: d.Kind, ResourcePath: d.ResourcePath, Binding: d.Binding}
	if retired, err := protection.ReadRetired(d.Store, readLimit); err != nil || retired != want {
		t.Fatalf("retired record = %+v (%v), want %+v", retired, err, want)
	}
	if _, err := os.Lstat(filepath.Join(d.Store, protection.ResourceFile)); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("generation after the retirement: %v, want it absent", err)
	}
}

// TestRecoveryAbsentStore holds that an absent store names no user in the current boot.
// Observation and recovery grade the resource safe, run release, and create no store.
func TestRecoveryAbsentStore(t *testing.T) {
	p := newProtocol()
	d := newResource(t, p, "cache")
	if err := os.RemoveAll(d.Store); err != nil {
		t.Fatal(err)
	}
	for name, observe := range map[string]func(func() error) (protection.Disposition, error){
		"seal":    func(release func() error) (protection.Disposition, error) { return p.SealAndObserve(d, release) },
		"recover": func(release func() error) (protection.Disposition, error) { return p.Recover(d.Store, release) },
	} {
		released := false
		disposition, err := observe(func() error { released = true; return nil })
		if err != nil || disposition.State != protection.Safe || len(disposition.Users) != 0 || !released {
			t.Errorf("%s of an absent store = %+v (%v), released=%v; want safe with no users and a release", name, disposition, err, released)
		}
		if _, err := os.Lstat(d.Store); !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("%s of an absent store created it: %v", name, err)
		}
	}
}

// TestOpenFailureRemovesItsOwnStore holds that a failed open deletes the store that it created
// before any generation exists there, and keeps a store that it found.
func TestOpenFailureRemovesItsOwnStore(t *testing.T) {
	for name, found := range map[string]bool{"created": false, "found": true} {
		t.Run(name, func(t *testing.T) {
			p := newProtocol()
			path := filepath.Join(t.TempDir(), "cache")
			if err := os.Mkdir(path, 0o755); err != nil {
				t.Fatal(err)
			}
			fresh := candidate(path)
			if found {
				if err := os.Mkdir(fresh.Store, 0o700); err != nil {
					t.Fatal(err)
				}
			}
			faults := router(map[durablefiletest.Fault]map[string]error{{Destination: "C/resource.json", Ordinal: 1}: {"rename": errInjected}}, map[string]string{fresh.Store: "C"})
			protection.SetReplacer(p, faults.Replacer())
			if _, err := p.OpenResource(fresh); !errors.Is(err, errInjected) {
				t.Fatalf("open with a failed publication = %v, want the injected fault", err)
			}
			if _, err := os.Lstat(fresh.Store); errors.Is(err, fs.ErrNotExist) == found {
				t.Fatalf("store after the failed open: %v, want present=%v", err, found)
			}
		})
	}
}

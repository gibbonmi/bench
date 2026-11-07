package protection_test

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
)

// TestDirectoryRefusesAForeignBinding hands each directory adapter shape the descriptor of its
// own directory carrying the identity of a sibling of the same shape. Only the binding differs,
// so only the owner identity comparison refuses it.
func TestDirectoryRefusesAForeignBinding(t *testing.T) {
	for name, adapter := range map[string]protection.Directory{
		"plain":  {Kind: protection.KindKitTest, Noun: "run"},
		"prefix": {Kind: protection.KindBinary, Noun: "run", Prefix: "owned-", Class: "owned directory"},
		"lock":   {Kind: protection.KindCache, Noun: "run", Lock: "lock"},
	} {
		t.Run(name, func(t *testing.T) {
			parent := t.TempDir()
			own, other := adapterDirectory(t, adapter, parent, "own"), adapterDirectory(t, adapter, parent, "other")
			if err := adapter.Validate(own); err != nil {
				t.Fatalf("the adapter refused its own descriptor %+v: %v", own, err)
			}
			foreign := own
			foreign.Binding = other.Binding
			if err := adapter.Validate(foreign); err == nil {
				t.Fatalf("the adapter accepted the foreign binding %+v", foreign)
			}
		})
	}
}

// adapterDirectory creates one directory that adapter accepts and describes it.
func adapterDirectory(t *testing.T, adapter protection.Directory, parent, name string) protection.Descriptor {
	t.Helper()
	dir := filepath.Join(parent, adapter.Prefix+name)
	if err := os.Mkdir(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if adapter.Lock != "" {
		if err := os.WriteFile(filepath.Join(dir, adapter.Lock), nil, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	d, err := adapter.Describe(dir)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

// TestStoreRoot is PL126. The store derives from the kind, the canonical path, and the binding
// under the user's boot root, and a rewritten TMPDIR or XDG_RUNTIME_DIR leaves it unchanged, so
// a nested CLI derives the store of its inherited descriptor.
func TestStoreRoot(t *testing.T) {
	setEnv := func() {
		t.Setenv("TMPDIR", t.TempDir())
		t.Setenv("XDG_RUNTIME_DIR", t.TempDir())
	}
	setEnv()
	system := protection.SystemBase()
	base := t.TempDir()
	protection.SetStoreBase(t, base)
	adapter := protection.Directory{Kind: protection.KindBinary, Noun: "run"}
	d := adapterDirectory(t, adapter, t.TempDir(), "own")
	setEnv()
	if again := protection.SystemBase(); again != system {
		t.Fatalf("system base under a rewritten environment = %s, want %s", again, system)
	}
	again, err := adapter.Describe(d.ResourcePath)
	if err != nil || again.Store != d.Store {
		t.Fatalf("store under a rewritten environment = %s (%v), want %s", again.Store, err, d.Store)
	}
	sum := sha256.Sum256([]byte(d.ResourcePath + "\x00" + d.Binding))
	user := filepath.Join(base, fmt.Sprintf("bench-process-%d", os.Geteuid()))
	boot := filepath.Base(filepath.Dir(d.Store))
	if runtime.GOOS == "linux" {
		data, err := os.ReadFile("/proc/sys/kernel/random/boot_id")
		if err != nil {
			t.Fatal(err)
		}
		boot = strings.TrimSpace(string(data))
	}
	if want := filepath.Join(user, boot, "binary-"+hex.EncodeToString(sum[:16])); d.Store != want {
		t.Fatalf("store = %s, want %s", d.Store, want)
	}
	if err := adapter.Validate(d); err != nil {
		t.Fatalf("the adapter refused its derived store %+v: %v", d, err)
	}
	inside := d
	inside.Store = filepath.Join(d.ResourcePath, "process-users")
	if err := adapter.Validate(inside); err == nil {
		t.Fatalf("the adapter accepted a store inside its resource %+v", inside)
	}
	// A replacement directory at the same path has another binding, so it has another store.
	if err := os.Remove(d.ResourcePath); err != nil {
		t.Fatal(err)
	}
	if replaced := adapterDirectory(t, adapter, filepath.Dir(d.ResourcePath), "own"); replaced.Store == d.Store {
		t.Fatalf("a replacement directory shares the store %s", d.Store)
	}
}

// TestStoreRootRefusal is PL127. Each root directory of an opened generation becomes a link, a
// directory with a wider mode, or, for a root user, a directory that another user owns. Each
// condition refuses the next open and the next launch, so no other user can forge records.
func TestStoreRootRefusal(t *testing.T) {
	conditions := map[string]func(t *testing.T, dir string){
		"link": func(t *testing.T, dir string) {
			if err := errors.Join(os.Rename(dir, dir+".real"), os.Symlink(dir+".real", dir)); err != nil {
				t.Fatal(err)
			}
		},
		"wider mode": func(t *testing.T, dir string) {
			if err := os.Chmod(dir, 0o750); err != nil {
				t.Fatal(err)
			}
		},
		"other owner": func(t *testing.T, dir string) {
			if os.Geteuid() != 0 {
				capability.Capability(t, capability.Privilege, "only a root user can give a directory another owner")
			}
			if err := os.Chown(dir, 65534, 65534); err != nil {
				t.Fatal(err)
			}
		},
	}
	for name, condition := range conditions {
		for level, up := range map[string]int{"user root": 1, "boot root": 0} {
			t.Run(name+"/"+level, func(t *testing.T) {
				adapter, p, dir, g := openPrivateRun(t)
				root := filepath.Dir(g.Resource.Store)
				for range up {
					root = filepath.Dir(root)
				}
				condition(t, root)
				if _, err := adapter.Describe(dir); err == nil {
					t.Fatalf("the adapter described %s under a refused root", dir)
				}
				requireRefused(t, p, g.Resource)
			})
		}
	}
}

// TestObserveRecordedFailsClosed holds that ObserveRecorded keeps a resource whose users it
// cannot observe. A store whose record names the resource but cannot be read, and a store root
// that cannot be listed, each grade the resource uncertain with the error, so its owner keeps it.
func TestObserveRecordedFailsClosed(t *testing.T) {
	for name, hide := range map[string]func(t *testing.T, store string){
		"unread record": func(t *testing.T, store string) {
			if err := os.Chmod(filepath.Join(store, protection.ResourceFile), 0o644); err != nil {
				t.Fatal(err)
			}
		},
		"unlisted root": func(t *testing.T, store string) {
			root := filepath.Dir(store)
			if err := os.Chmod(root, 0o300); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = os.Chmod(root, 0o700) })
		},
	} {
		t.Run(name, func(t *testing.T) {
			adapter, _, dir, g := openPrivateRun(t)
			// The open generation's record names dir, so only the hidden condition stops a safe
			// verdict.
			hide(t, g.Resource.Store)
			d, err := adapter.ObserveRecorded(dir, readLimit)
			if err == nil || d.State != protection.Uncertain || d.Err == nil {
				t.Fatalf("observation of a store that cannot be read = %+v (%v), want uncertain with its error", d, err)
			}
		})
	}
}

// cacheKind is the directory kind that the absent-resource check tests run over.
var cacheKind = protection.Directory{Kind: protection.KindCache, Noun: "cache"}

// retiredCaches returns a protocol of the cache and selection kinds under a private store base,
// and a parent directory for the resources that it retires.
func retiredCaches(t *testing.T) (*protection.Protocol, string) {
	t.Helper()
	return privateProtocol(t, cacheKind, selection), t.TempDir()
}

// renameStore renames the store of d so that the kind of to starts its name in place of the
// kind of from. The store name and its retired record then name different kinds.
func renameStore(t *testing.T, d protection.Descriptor, from, to protection.Directory) protection.Descriptor {
	t.Helper()
	name := string(to.Kind) + strings.TrimPrefix(filepath.Base(d.Store), string(from.Kind))
	renamed := d
	renamed.Store = filepath.Join(filepath.Dir(d.Store), name)
	if err := os.Rename(d.Store, renamed.Store); err != nil {
		t.Fatal(err)
	}
	return renamed
}

// requireCheck fails unless the cache check returns the wanted error state and leaves each
// store of stores present exactly when kept names it.
func requireCheck(t *testing.T, stores map[string]protection.Descriptor, kept map[string]bool, wantErr bool) {
	t.Helper()
	if err := cacheKind.DiscardRetired(readLimit); (err != nil) != wantErr {
		t.Fatalf("cache check error = %v, want an error %t", err, wantErr)
	}
	for name, d := range stores {
		if _, err := os.Lstat(d.Store); (err == nil) != kept[name] {
			t.Errorf("store of %s after the cache check: %v, want present=%t", name, err, kept[name])
		}
	}
}

// TestDiscardRetired is RS9, RS10, and RS30. The cache check deletes the retired cache store
// whose cache is gone. It keeps the retired cache store whose cache exists, and each retired
// store whose resource is gone but whose store name or retired record names another kind.
func TestDiscardRetired(t *testing.T) {
	p, parent := retiredCaches(t)
	stores := map[string]protection.Descriptor{
		"gone cache":        goneStore(t, p, cacheKind, filepath.Join(parent, "gone")),
		"kept cache":        retire(t, p, cacheKind, filepath.Join(parent, "kept")),
		"other kind":        goneStore(t, p, selection, filepath.Join(parent, "other")),
		"other kind named":  renameStore(t, goneStore(t, p, selection, filepath.Join(parent, "named")), selection, cacheKind),
		"other kind record": renameStore(t, goneStore(t, p, cacheKind, filepath.Join(parent, "record")), cacheKind, selection),
	}
	requireCheck(t, stores, map[string]bool{"kept cache": true, "other kind": true, "other kind named": true, "other kind record": true}, false)
}

// TestDiscardRetiredFailsClosed is RS12 and RS13. A retired cache store whose record the check
// cannot read stays, and the check returns no error for it. A store root that the check cannot
// list, and a store that passes the predicate but cannot be deleted, each fail the check.
func TestDiscardRetiredFailsClosed(t *testing.T) {
	t.Run("unread record", func(t *testing.T) {
		p, parent := retiredCaches(t)
		d := goneStore(t, p, cacheKind, filepath.Join(parent, "gone"))
		if err := os.Chmod(filepath.Join(d.Store, protection.RetiredFile), 0o644); err != nil {
			t.Fatal(err)
		}
		requireCheck(t, map[string]protection.Descriptor{"unread": d}, map[string]bool{"unread": true}, false)
	})
	t.Run("unlisted root", func(t *testing.T) {
		p, parent := retiredCaches(t)
		d := goneStore(t, p, cacheKind, filepath.Join(parent, "gone"))
		root := filepath.Dir(d.Store)
		if err := os.Chmod(root, 0o300); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.Chmod(root, 0o700) })
		requireCheck(t, nil, nil, true)
	})
	t.Run("failed deletion", func(t *testing.T) {
		p, parent := retiredCaches(t)
		d := goneStore(t, p, cacheKind, filepath.Join(parent, "gone"))
		held := filepath.Join(d.Store, "held")
		if err := os.Mkdir(held, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := errors.Join(os.WriteFile(filepath.Join(held, "entry"), nil, 0o600), os.Chmod(held, 0o500)); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.Chmod(held, 0o700) })
		if probe, err := os.CreateTemp(held, "writable"); err == nil {
			_ = errors.Join(probe.Close(), os.Remove(probe.Name()))
			capability.Capability(t, capability.Privilege, "a mode 0500 directory stays writable by this user")
		}
		requireCheck(t, map[string]protection.Descriptor{"undeletable": d}, map[string]bool{"undeletable": true}, true)
	})
}

// openPrivateRun opens the generation of one run directory under a private store base, so a
// test can change its store root.
func openPrivateRun(t *testing.T) (protection.Directory, *protection.Protocol, string, protection.Generation) {
	t.Helper()
	adapter := protection.Directory{Kind: protection.KindKitTest, Noun: "run"}
	p := privateProtocol(t, adapter)
	dir := filepath.Join(t.TempDir(), "run")
	if err := os.Mkdir(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	g, err := adapter.Open(p, dir)
	if err != nil {
		t.Fatal(err)
	}
	return adapter, p, dir, g
}

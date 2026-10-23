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
				protection.SetStoreBase(t, t.TempDir())
				adapter := protection.Directory{Kind: protection.KindKitTest, Noun: "run"}
				p := &protection.Protocol{Validators: protection.Validators{protection.KindKitTest: adapter.Validate}, ReadLimit: readLimit}
				dir := filepath.Join(t.TempDir(), "run")
				if err := os.Mkdir(dir, 0o700); err != nil {
					t.Fatal(err)
				}
				g, err := adapter.Open(p, dir)
				if err != nil {
					t.Fatal(err)
				}
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

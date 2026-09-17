package chargeevidence_test

import (
	"errors"
	"math"
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	ce "github.com/gibbonmi/bench/internal/chargeevidence"
)

// storeName is the store directory name, stated independently of ce.StoreName so that a
// renamed store turns these tests red.
const storeName = "bench-charge-evidence"

func newStore(t *testing.T, options ce.StoreOptions) (*ce.Store, string) {
	t.Helper()
	common, err := os.MkdirTemp("", "ce-store")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(common) })
	return ce.OpenStore(common, options), filepath.Join(common, storeName)
}

func publish(t *testing.T, store *ce.Store, pack *ce.Pack, quota uint64) string {
	t.Helper()
	staged, err := store.Stage(pack, quota, 1)
	if err != nil {
		t.Fatalf("Stage: %v", err)
	}
	identity, err := staged.Publish(1)
	if err != nil {
		t.Fatalf("Publish: %v", err)
	}
	return identity
}

func packPath(dir, identity string) string {
	return filepath.Join(dir, strings.TrimPrefix(identity, ce.IdentityPrefix)+ce.PackSuffix)
}

func refusalClass(err error) string {
	var refusal *ce.Refusal
	if errors.As(err, &refusal) {
		return refusal.Class
	}
	return ""
}

// storeEntries lists the published and temporary pack names in dir.
func storeEntries(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	var names []string
	for _, entry := range entries {
		if name := entry.Name(); strings.HasSuffix(name, ce.PackSuffix) || strings.HasPrefix(name, ce.TempPrefix) {
			names = append(names, name)
		}
	}
	return names
}

// copyStore copies every regular file of the store at dir into a new directory at dst.
func copyStore(t *testing.T, dir, dst string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err == nil {
		err = os.Mkdir(dst, 0o700)
	}
	for _, entry := range entries {
		var data []byte
		if data, err = os.ReadFile(filepath.Join(dir, entry.Name())); err == nil {
			err = os.WriteFile(filepath.Join(dst, entry.Name()), data, 0o600)
		}
		if err != nil {
			break
		}
	}
	if err != nil {
		t.Fatal(err)
	}
}

// flipAt rewrites one byte of a published pack in place.
func flipAt(t *testing.T, path string, offset int64) {
	t.Helper()
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	b := make([]byte, 1)
	if _, err := file.ReadAt(b, offset); err != nil {
		t.Fatal(err)
	}
	b[0] ^= 0x01
	if _, err := file.WriteAt(b, offset); err != nil {
		t.Fatal(err)
	}
}

// TestEvidenceQuota is CE74 and CE76.
func TestEvidenceQuota(t *testing.T) {
	store, dir := newStore(t, ce.StoreOptions{})
	first := mustBuild(t, fixtureCandidate("# One\n"))
	identity := publish(t, store, first, ce.DefaultQuota)
	orphan := filepath.Join(dir, ce.TempPrefix+"0123456789abcdef"+ce.TempSuffix)
	if err := os.WriteFile(orphan, make([]byte, 1000), 0o600); err != nil {
		t.Fatal(err)
	}
	second := mustBuild(t, fixtureCandidate("# Two\n"))
	used := uint64(len(first.Bytes()) + 1000)
	need := used + uint64(len(second.Bytes()))
	_, err := store.Stage(second, need-1, 1)
	var capacity *ce.CapacityError
	if !errors.As(err, &capacity) || capacity.Observed != used || capacity.Required != need || capacity.Quota != need-1 {
		t.Fatalf("capacity = %+v (%v), want observed %d required %d", capacity, err, used, need)
	}
	if got := storeEntries(t, dir); len(got) != 2 {
		t.Fatalf("capacity refusal changed the store: %v", got)
	}
	if _, err := os.Stat(orphan); err != nil {
		t.Fatalf("capacity refusal removed the orphan: %v", err)
	}
	if artifact, err := store.Open(identity); err != nil {
		t.Fatalf("preserved artifact: %v", err)
	} else {
		artifact.Close()
	}
	publish(t, store, second, need)
}

// TestEvidenceAdmission is CE75 at the capacity rule: only the total quota limits a
// candidate, so an explicit quota admits an artifact larger than the default quota.
func TestEvidenceAdmission(t *testing.T) {
	const gib = uint64(1) << 30
	for _, test := range []struct {
		name                       string
		observed, candidate, quota uint64
		admitted                   bool
		required                   uint64
	}{
		{"CE75 artifact above the default quota", 0, 2 * gib, 3 * gib, true, 0},
		{"store and artifact exactly at the quota", gib, 2 * gib, 3 * gib, true, 0},
		{"one byte above the quota", gib, 2*gib + 1, 3 * gib, false, 3*gib + 1},
		{"largest representable artifact", 0, math.MaxUint64, math.MaxUint64, true, 0},
		{"unrepresentable total", math.MaxUint64, 1, math.MaxUint64, false, 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := ce.Admit(test.observed, test.candidate, test.quota)
			var capacity *ce.CapacityError
			if test.admitted != (err == nil) || !test.admitted && (!errors.As(err, &capacity) || capacity.Required != test.required ||
				capacity.Observed != test.observed || capacity.Candidate != test.candidate || capacity.Quota != test.quota) {
				t.Fatalf("Admit(%d, %d, %d) = %+v (%v)", test.observed, test.candidate, test.quota, capacity, err)
			}
		})
	}
}

// TestEvidenceQuotaOverride is CE74 and CE75 over real store bytes: a sparse temporary pack
// counts against the default quota, and a larger explicit quota admits the candidate.
func TestEvidenceQuotaOverride(t *testing.T) {
	store, dir := newStore(t, ce.StoreOptions{})
	publish(t, store, mustBuild(t, fixtureCandidate("# One\n")), ce.DefaultQuota)
	large, err := os.Create(filepath.Join(dir, ce.TempPrefix+"ffffffffffffffff"+ce.TempSuffix))
	if err != nil {
		t.Fatal(err)
	}
	if err := large.Truncate(ce.DefaultQuota); err != nil {
		t.Fatal(err)
	}
	large.Close()
	pack := mustBuild(t, fixtureCandidate("# Two\n"))
	if _, err := store.Stage(pack, ce.DefaultQuota, 1); !errors.As(err, new(*ce.CapacityError)) {
		t.Fatalf("default quota = %v, want a capacity refusal", err)
	}
	publish(t, store, pack, 2*ce.DefaultQuota)
}

// TestEvidencePublicationFailure is CE88 and the deterministic write, sync, and
// verification failures.
func TestEvidencePublicationFailure(t *testing.T) {
	for _, step := range []string{ce.StepWriteTemp, ce.StepSyncTemp, ce.StepVerifyTemp, ce.StepLink} {
		t.Run(step, func(t *testing.T) {
			fail := func(at string) error {
				if at == step {
					return errors.New("injected " + at + " failure")
				}
				return nil
			}
			store, dir := newStore(t, ce.StoreOptions{Fault: fail})
			pack := mustBuild(t, fixtureCandidate("# One\n"))
			staged, err := store.Stage(pack, ce.DefaultQuota, 1)
			if err == nil {
				_, err = staged.Publish(1)
			}
			if err == nil || !strings.Contains(err.Error(), "injected "+step) {
				t.Fatalf("%s failure = %v", step, err)
			}
			if got := storeEntries(t, dir); len(got) != 0 {
				t.Fatalf("%s failure left %v", step, got)
			}
		})
	}
	t.Run("stale attempt", func(t *testing.T) {
		store, dir := newStore(t, ce.StoreOptions{})
		staged, err := store.Stage(mustBuild(t, fixtureCandidate("# One\n")), ce.DefaultQuota, 1)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := staged.Publish(2); refusalClass(err) != ce.RefusePublish || len(storeEntries(t, dir)) != 0 {
			t.Fatalf("stale attempt = %v, entries %v", err, storeEntries(t, dir))
		}
	})
}

// TestEvidenceStoreKinds is CE91 through CE95. No kind is opened as a pack, and a FIFO is
// never read.
func TestEvidenceStoreKinds(t *testing.T) {
	for _, kind := range []string{"CE91 symlink", "CE92 FIFO", "CE93 socket", "CE94 device", "CE95 directory"} {
		t.Run(kind, func(t *testing.T) {
			store, dir := newStore(t, ce.StoreOptions{})
			decoy := publish(t, store, mustBuild(t, fixtureCandidate("# Decoy\n")), ce.DefaultQuota)
			identity := mustBuild(t, fixtureCandidate("# One\n")).Identity()
			path := packPath(dir, identity)
			var err error
			switch kind {
			case "CE91 symlink":
				err = os.Symlink(filepath.Base(packPath(dir, decoy)), path)
			case "CE92 FIFO":
				err = syscall.Mkfifo(path, 0o600)
			case "CE93 socket":
				// A relative name keeps the socket address under the platform length limit.
				t.Chdir(dir)
				var listener net.Listener
				if listener, err = net.Listen("unix", filepath.Base(path)); err == nil {
					t.Cleanup(func() { _ = listener.Close() })
				}
			case "CE94 device":
				if err = syscall.Mknod(path, syscall.S_IFCHR|0o600, 0x0103); err != nil {
					capability.Capability(t, capability.Privilege, "cannot create a character device: "+err.Error())
				}
			case "CE95 directory":
				err = os.Mkdir(path, 0o700)
			}
			if err != nil {
				t.Fatal(err)
			}
			window := bounds.TestDeadline(0)
			done := make(chan error, 1)
			go func() {
				_, openErr := store.Open(identity)
				done <- openErr
			}()
			select {
			case err := <-done:
				if refusalClass(err) != ce.RefuseUnsafe {
					t.Fatalf("%s = %v, want %s", kind, err, ce.RefuseUnsafe)
				}
			case <-time.After(window):
				t.Fatal(bounds.TestTimeoutVerdict(kind+" to refuse without blocking the reader", window))
			}
			if info, err := os.Lstat(path); err != nil || info.Mode().IsRegular() {
				t.Fatalf("%s object changed: %v", kind, err)
			}
		})
	}
}

// TestEvidenceStorePaths is CE96: an invalid identifier refuses before any path use.
func TestEvidenceStorePaths(t *testing.T) {
	store := ce.OpenStore(filepath.Join(t.TempDir(), "absent common dir"), ce.StoreOptions{})
	for _, identity := range []string{"../escape", "/absolute", "sha256:abc", "sha256:" + strings.Repeat("A", 64), "sha256:" + strings.Repeat("a", 63) + "/", "sha256:../" + strings.Repeat("a", 61)} {
		if _, err := store.Open(identity); refusalClass(err) != ce.RefuseIdentifier {
			t.Fatalf("identifier %q = %v, want %s", identity, err, ce.RefuseIdentifier)
		}
	}
}

// TestEvidenceExistingCorruption is CE97.
func TestEvidenceExistingCorruption(t *testing.T) {
	store, dir := newStore(t, ce.StoreOptions{})
	pack := mustBuild(t, fixtureCandidate("# One\n"))
	identity := publish(t, store, pack, ce.DefaultQuota)
	flipAt(t, packPath(dir, identity), int64(len(pack.Bytes())-1))
	before, err := os.ReadFile(packPath(dir, identity))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Stage(pack, ce.DefaultQuota, 1); refusalClass(err) != ce.RefuseExisting {
		t.Fatalf("corrupt existing = %v, want %s", err, ce.RefuseExisting)
	}
	if after, _ := os.ReadFile(packPath(dir, identity)); string(after) != string(before) {
		t.Fatal("preparation replaced the corrupt artifact")
	}
}

// TestEvidenceStoreDirectoryStates is CE98, CE99, CE118, and the empty store read.
func TestEvidenceStoreDirectoryStates(t *testing.T) {
	pack := mustBuild(t, fixtureCandidate("# One\n"))
	t.Run("CE118 absent read", func(t *testing.T) {
		store, dir := newStore(t, ce.StoreOptions{})
		if _, err := store.Open(pack.Identity()); refusalClass(err) != ce.RefuseAbsent {
			t.Fatalf("absent store read = %v", err)
		}
		if _, err := os.Lstat(dir); !os.IsNotExist(err) {
			t.Fatalf("read created the store: %v", err)
		}
	})
	t.Run("empty store read", func(t *testing.T) {
		store, dir := newStore(t, ce.StoreOptions{})
		if err := os.Mkdir(dir, 0o700); err != nil {
			t.Fatal(err)
		}
		if _, err := store.Open(pack.Identity()); refusalClass(err) != ce.RefuseAbsent {
			t.Fatalf("empty store read = %v, want %s", err, ce.RefuseAbsent)
		}
		if entries, err := os.ReadDir(dir); err != nil || len(entries) != 0 {
			t.Fatalf("empty store read changed the store: %v, %v", entries, err)
		}
	})
	t.Run("CE99 preparation creates", func(t *testing.T) {
		store, dir := newStore(t, ce.StoreOptions{})
		publish(t, store, pack, ce.DefaultQuota)
		if info, err := os.Lstat(dir); err != nil || !info.IsDir() || info.Mode().Perm() != 0o700 {
			t.Fatalf("created store = %v, %v", info, err)
		}
	})
	t.Run("CE98 replaced while opened", func(t *testing.T) {
		for _, operation := range []string{"read", "preparation"} {
			t.Run(operation, func(t *testing.T) {
				store, dir := newStore(t, ce.StoreOptions{})
				publish(t, store, pack, ce.DefaultQuota)
				copyStore(t, dir, dir+".copy")
				replaced := false
				racing := ce.OpenStore(filepath.Dir(dir), ce.StoreOptions{Pause: func(stage string) {
					if stage == ce.StageStoreInspected && !replaced {
						replaced = true
						if err := os.Rename(dir, dir+".old"); err != nil {
							t.Fatal(err)
						}
						if err := os.Rename(dir+".copy", dir); err != nil {
							t.Fatal(err)
						}
					}
				}})
				var err error
				if operation == "read" {
					_, err = racing.Open(pack.Identity())
				} else {
					_, err = racing.Stage(pack, ce.DefaultQuota, 1)
				}
				if !replaced || refusalClass(err) != ce.RefuseUnsafe {
					t.Fatalf("replaced store %s = %v, want %s", operation, err, ce.RefuseUnsafe)
				}
			})
		}
	})
	for _, state := range []string{"symlink", "regular file"} {
		t.Run("CE98 "+state, func(t *testing.T) {
			store, dir := newStore(t, ce.StoreOptions{})
			var err error
			if state == "symlink" {
				target := filepath.Join(filepath.Dir(dir), "elsewhere")
				if err = os.Mkdir(target, 0o700); err == nil {
					err = os.Symlink(target, dir)
				}
			} else {
				err = os.WriteFile(dir, []byte("x"), 0o600)
			}
			if err != nil {
				t.Fatal(err)
			}
			if _, err := store.Stage(pack, ce.DefaultQuota, 1); refusalClass(err) != ce.RefuseUnsafe {
				t.Fatalf("%s store preparation = %v", state, err)
			}
			if _, err := store.Open(pack.Identity()); refusalClass(err) != ce.RefuseUnsafe {
				t.Fatalf("%s store read = %v", state, err)
			}
		})
	}
}

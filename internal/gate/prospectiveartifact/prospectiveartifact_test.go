package prospectiveartifact

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"syscall"
	"testing"
	"time"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

func TestOpenPublishesPrivateOwnerRecordBeforeCheckout(t *testing.T) {
	repository := newRepository(t)
	owner, err := (Factory{TempRoot: t.TempDir()}).Open(repository)
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Lstat(filepath.Join(owner.Root(), ownerRecordName))
	if err != nil {
		t.Fatal(err)
	}
	if !info.Mode().IsRegular() || info.Mode().Perm() != 0o600 {
		t.Fatalf("owner record mode = %v, want regular 0600", info.Mode())
	}
	data, err := os.ReadFile(filepath.Join(owner.Root(), ownerRecordName))
	if err != nil {
		t.Fatal(err)
	}
	var record struct {
		Schema    int    `json:"schema"`
		OwnerPID  int    `json:"owner_pid"`
		CommonDir string `json:"common_dir"`
	}
	if err := json.Unmarshal(data, &record); err != nil {
		t.Fatal(err)
	}
	common, err := canonicalCommonDir(repository)
	if err != nil {
		t.Fatal(err)
	}
	if record.Schema != 1 || record.OwnerPID != os.Getpid() || record.CommonDir != common {
		t.Fatalf("owner record = %#v, want schema 1, pid %d, common directory %q", record, os.Getpid(), common)
	}
	if _, err := os.Lstat(owner.Checkout()); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("checkout before materialization = %v, want absent", err)
	}
}

func TestOpenRetainsForeignBundleAndRecoversDeadRecordOnlyBundle(t *testing.T) {
	repository := newRepository(t)
	common, err := benchgit.CommonDir(repository)
	if err != nil {
		t.Fatal(err)
	}
	base := t.TempDir()
	foreign := testBundle(t, base, ownerRecord{Schema: 1, OwnerPID: 41, CommonDir: filepath.Join(common, "foreign")})
	dead := testBundle(t, base, ownerRecord{Schema: 1, OwnerPID: 42, CommonDir: common})
	owner, err := (Factory{TempRoot: base, Probe: func(int) error { return syscall.ESRCH }}).Open(repository)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(foreign); err != nil {
		t.Fatalf("foreign bundle = %v, want retained", err)
	}
	if _, err := os.Stat(dead); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("dead record-only bundle = %v, want absent", err)
	}
	if err := owner.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestOpenRecoversARegisteredDeadBundleBeforeItCreatesAnother(t *testing.T) {
	repository := newRepository(t)
	common, err := benchgit.CommonDir(repository)
	if err != nil {
		t.Fatal(err)
	}
	base := t.TempDir()
	dead := testBundle(t, base, ownerRecord{Schema: 1, OwnerPID: 42, CommonDir: common})
	checkout := filepath.Join(dead, checkoutName)
	gitRun(t, repository, "worktree", "add", "-q", "--detach", checkout, "HEAD")
	run := filepath.Join(dead, "bench-run [*]", "bench")
	if err := os.MkdirAll(filepath.Dir(run), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(run, []byte("old candidate"), 0o700); err != nil {
		t.Fatal(err)
	}
	owner, err := (Factory{TempRoot: base, Probe: func(int) error { return syscall.ESRCH }}).Open(repository)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dead); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("dead bundle before new checkout = %v, want absent", err)
	}
	for _, worktree := range mustWorktrees(t, repository) {
		if filepath.Clean(worktree.Path) == filepath.Clean(checkout) {
			t.Fatalf("dead registered checkout %q survived recovery", checkout)
		}
	}
	if _, err := os.Stat(owner.Checkout()); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("new checkout before materialization = %v, want absent", err)
	}
	if err := owner.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestOpenRecoversARegisteredDeadCheckoutWithoutARunBinary(t *testing.T) {
	repository := newRepository(t)
	common, err := canonicalCommonDir(repository)
	if err != nil {
		t.Fatal(err)
	}
	base := t.TempDir()
	dead := testBundle(t, base, ownerRecord{Schema: 1, OwnerPID: 42, CommonDir: common})
	checkout := filepath.Join(dead, checkoutName)
	gitRun(t, repository, "worktree", "add", "-q", "--detach", checkout, "HEAD")

	owner, err := (Factory{TempRoot: base, Probe: func(int) error { return syscall.ESRCH }}).Open(repository)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dead); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("dead checkout-only bundle = %v, want absent", err)
	}
	requireCheckoutUnregistered(t, repository, checkout)
	if err := owner.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestOpenRecoversAStaleRegistrationWithoutACheckoutPath(t *testing.T) {
	repository := newRepository(t)
	common, err := canonicalCommonDir(repository)
	if err != nil {
		t.Fatal(err)
	}
	base := t.TempDir()
	dead := testBundle(t, base, ownerRecord{Schema: 1, OwnerPID: 42, CommonDir: common})
	checkout := filepath.Join(dead, checkoutName)
	gitRun(t, repository, "worktree", "add", "-q", "--detach", checkout, "HEAD")
	if err := os.RemoveAll(checkout); err != nil {
		t.Fatal(err)
	}

	owner, err := (Factory{TempRoot: base, Probe: func(int) error { return syscall.ESRCH }}).Open(repository)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dead); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("dead stale-registration bundle = %v, want absent", err)
	}
	requireCheckoutUnregistered(t, repository, checkout)
	if err := owner.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestOpenRefusesRecoveryWhenRegistrationRemovalFails(t *testing.T) {
	repository := newRepository(t)
	common, err := benchgit.CommonDir(repository)
	if err != nil {
		t.Fatal(err)
	}
	base := t.TempDir()
	dead := testBundle(t, base, ownerRecord{Schema: 1, OwnerPID: 42, CommonDir: common})
	gitRun(t, repository, "worktree", "add", "-q", "--detach", filepath.Join(dead, checkoutName), "HEAD")
	_, err = (Factory{
		TempRoot: base,
		Probe:    func(int) error { return syscall.ESRCH },
		Remove:   func(string, string) error { return errors.New("registration refused") },
	}).Open(repository)
	if err == nil {
		t.Fatal("recovery with a registration failure = nil, want refusal")
	}
	if _, err := os.Stat(dead); err != nil {
		t.Fatalf("dead bundle after refusal = %v, want retained", err)
	}
}

func TestCloseConfinesRemovalToItsBundle(t *testing.T) {
	repository := newRepository(t)
	parent := t.TempDir()
	base := filepath.Join(parent, "temp root [*]")
	if err := os.Mkdir(base, 0o700); err != nil {
		t.Fatal(err)
	}
	sentinel := filepath.Join(base, "keep [*] space")
	if err := os.WriteFile(sentinel, []byte("keep"), 0o600); err != nil {
		t.Fatal(err)
	}
	owner, err := (Factory{TempRoot: base}).Open(repository)
	if err != nil {
		t.Fatal(err)
	}
	root := owner.Root()
	if filepath.Dir(root) != base {
		t.Fatalf("bundle root = %q, want child of hostile temporary root %q", root, base)
	}
	if err := owner.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(root); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("closed root = %v, want absent", err)
	}
	if _, err := os.Stat(sentinel); err != nil {
		t.Fatalf("sentinel = %v, want retained", err)
	}
}

func requireCheckoutUnregistered(t *testing.T, repository, checkout string) {
	t.Helper()
	for _, worktree := range mustWorktrees(t, repository) {
		if filepath.Clean(worktree.Path) == filepath.Clean(checkout) {
			t.Fatalf("dead registered checkout %q survived recovery", checkout)
		}
	}
}

func TestCloseRemovesTheRegisteredCheckoutBeforeItsBundle(t *testing.T) {
	repository := newRepository(t)
	tree := gitOutput(t, repository, "write-tree")
	owner, err := (Factory{TempRoot: t.TempDir()}).Open(repository)
	if err != nil {
		t.Fatal(err)
	}
	if err := owner.Materialize(tree); err != nil {
		t.Fatal(err)
	}
	if err := owner.Close(); err != nil {
		t.Fatal(err)
	}
	for _, worktree := range mustWorktrees(t, repository) {
		if filepath.Clean(worktree.Path) == filepath.Clean(owner.Checkout()) {
			t.Fatalf("registered checkout %q survived bundle close", owner.Checkout())
		}
	}
	if _, err := os.Stat(owner.Root()); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("closed bundle = %v, want absent", err)
	}
}

func testBundle(t *testing.T, base string, record ownerRecord) string {
	t.Helper()
	root, err := os.MkdirTemp(base, bundlePrefix)
	if err != nil {
		t.Fatal(err)
	}
	if err := writeRecord(root, record); err != nil {
		t.Fatal(err)
	}
	return root
}

func newRepository(t *testing.T) string {
	t.Helper()
	repository := t.TempDir()
	if output, err := exec.Command("git", "-C", repository, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, output)
	}
	gitRun(t, repository, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "--allow-empty", "-qm", "base")
	return repository
}

func gitRun(t *testing.T, repository string, args ...string) {
	t.Helper()
	output, err := exec.Command("git", append([]string{"-C", repository}, args...)...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v: %s", args, err, output)
	}
}

func gitOutput(t *testing.T, repository string, args ...string) string {
	t.Helper()
	output, err := exec.Command("git", append([]string{"-C", repository}, args...)...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v: %s", args, err, output)
	}
	return string(output[:len(output)-1])
}

func mustWorktrees(t *testing.T, repository string) []benchgit.Worktree {
	t.Helper()
	worktrees, err := benchgit.Worktrees(repository)
	if err != nil {
		t.Fatal(err)
	}
	return worktrees
}

// TestSweepRetainsABundleWithAnUnsupportedRecordSchema is PAR12. The record names this
// repository and a definitely absent PID, so the unknown schema is the only fact that
// denies removal.
func TestSweepRetainsABundleWithAnUnsupportedRecordSchema(t *testing.T) {
	repository := newRepository(t)
	base := t.TempDir()
	candidate := plantedBundle(t, base)
	plantRecord(t, candidate, recordBody(t, repository, 2, 42))

	requireSweepRetains(t, repository, base, deadProbe, candidate)
}

// TestSweepRetainsABundleWithAMissingEmptyOrMalformedRecord is PAR13. Each candidate is
// otherwise removable, so a parse failure that fell through to deletion would take a
// bundle no record ever authorized.
func TestSweepRetainsABundleWithAMissingEmptyOrMalformedRecord(t *testing.T) {
	for _, row := range []struct {
		name    string
		present bool
		body    string
	}{
		{name: "missing"},
		{name: "empty", present: true},
		{name: "malformed", present: true, body: "{\"schema\":1,\"owner_pid\":42"},
	} {
		t.Run(row.name, func(t *testing.T) {
			repository := newRepository(t)
			base := t.TempDir()
			candidate := plantedBundle(t, base)
			if row.present {
				plantRecord(t, candidate, row.body)
			}
			requireSweepRetains(t, repository, base, deadProbe, candidate)
		})
	}
}

// TestSweepRetainsABundleWithANonRegularRecord is PAR14. Each candidate carries a record
// whose bytes would authorize removal, so only the file shape denies it.
func TestSweepRetainsABundleWithANonRegularRecord(t *testing.T) {
	for _, row := range []struct {
		name  string
		plant func(*testing.T, string, string)
	}{
		{name: "symbolic link", plant: func(t *testing.T, candidate, body string) {
			published := filepath.Join(candidate, "published.json")
			if err := os.WriteFile(published, []byte(body), 0o600); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(published, filepath.Join(candidate, ownerRecordName)); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "named pipe", plant: func(t *testing.T, candidate, _ string) {
			if err := syscall.Mkfifo(filepath.Join(candidate, ownerRecordName), 0o600); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "directory", plant: func(t *testing.T, candidate, _ string) {
			if err := os.Mkdir(filepath.Join(candidate, ownerRecordName), 0o600); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(row.name, func(t *testing.T) {
			repository := newRepository(t)
			base := t.TempDir()
			candidate := plantedBundle(t, base)
			row.plant(t, candidate, recordBody(t, repository, 1, 42))
			requireSweepRetains(t, repository, base, deadProbe, candidate)
		})
	}
}

// TestSweepRetainsASymbolicLinkOrSpecialFileBundleCandidate is PAR15. The link points at
// a real removable bundle outside the scanned root, so a scanner that followed the link
// would delete resources its own namespace never held.
func TestSweepRetainsASymbolicLinkOrSpecialFileBundleCandidate(t *testing.T) {
	t.Run("symbolic link", func(t *testing.T) {
		repository := newRepository(t)
		base := t.TempDir()
		target := plantedBundle(t, t.TempDir())
		plantRecord(t, target, recordBody(t, repository, 1, 42))
		link := filepath.Join(base, bundlePrefix+"link")
		if err := os.Symlink(target, link); err != nil {
			t.Fatal(err)
		}
		requireSweepRetains(t, repository, base, deadProbe, link, target)
	})
	t.Run("named pipe", func(t *testing.T) {
		repository := newRepository(t)
		base := t.TempDir()
		pipe := filepath.Join(base, bundlePrefix+"pipe")
		if err := syscall.Mkfifo(pipe, 0o600); err != nil {
			t.Fatal(err)
		}
		requireSweepRetains(t, repository, base, deadProbe, pipe)
	})
}

// TestSweepRetainsAForeignSamePrefixDirectory is PAR16. The directory carries the bundle
// prefix and arbitrary content, and its record names another repository, so the prefix
// alone must grant no authority over it.
func TestSweepRetainsAForeignSamePrefixDirectory(t *testing.T) {
	repository := newRepository(t)
	base := t.TempDir()
	foreign := plantedBundle(t, base)
	plantRecord(t, foreign, recordBody(t, newRepository(t), 1, 42))
	if err := os.MkdirAll(filepath.Join(foreign, "notes"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(foreign, "notes", "keep.txt"), []byte("foreign content\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	requireSweepRetains(t, repository, base, deadProbe, foreign)
}

// TestSweepRetainsAPermissionRefusedBundle is PAR17. A refused probe proves the owner
// exists and is unreachable, never that it is absent.
func TestSweepRetainsAPermissionRefusedBundle(t *testing.T) {
	repository := newRepository(t)
	base := t.TempDir()
	candidate := plantedBundle(t, base)
	plantRecord(t, candidate, recordBody(t, repository, 1, 42))

	requireSweepRetains(t, repository, base, func(int) error { return syscall.EPERM }, candidate)
}

// TestSweepRetainsABundleWithAnUnknownProbeFailure is PAR18. Only the absent-process
// result is death proof, so every other probe failure retains its bundle.
func TestSweepRetainsABundleWithAnUnknownProbeFailure(t *testing.T) {
	for _, row := range []struct {
		name  string
		probe func(int) error
	}{
		{name: "unclassified failure", probe: func(int) error { return errors.New("probe failed") }},
		{name: "invalid request", probe: func(int) error { return syscall.EINVAL }},
	} {
		t.Run(row.name, func(t *testing.T) {
			repository := newRepository(t)
			base := t.TempDir()
			candidate := plantedBundle(t, base)
			plantRecord(t, candidate, recordBody(t, repository, 1, 42))
			requireSweepRetains(t, repository, base, row.probe, candidate)
		})
	}
}

// TestSweepRetainsAnAnsweringPIDWithAnOldRecord is PAR19. The record is a year old, so a
// sweep that fell back on age would delete a live owner's resources.
func TestSweepRetainsAnAnsweringPIDWithAnOldRecord(t *testing.T) {
	repository := newRepository(t)
	base := t.TempDir()
	candidate := plantedBundle(t, base)
	plantRecord(t, candidate, recordBody(t, repository, 1, 42))
	stale := time.Now().Add(-365 * 24 * time.Hour)
	for _, path := range []string{filepath.Join(candidate, ownerRecordName), candidate} {
		if err := os.Chtimes(path, stale, stale); err != nil {
			t.Fatal(err)
		}
	}

	requireSweepRetains(t, repository, base, func(int) error { return nil }, candidate)
}

// TestOneSweepRemovesOnlyTheDeadBundleOfADeadAndLivePair is PAR20 at the module seam.
// Both bundles name this repository, so only the per-candidate probe result separates
// them.
func TestOneSweepRemovesOnlyTheDeadBundleOfADeadAndLivePair(t *testing.T) {
	repository := newRepository(t)
	base := t.TempDir()
	dead := plantedBundle(t, base)
	plantRecord(t, dead, recordBody(t, repository, 1, 42))
	live := plantedBundle(t, base)
	plantRecord(t, live, recordBody(t, repository, 1, 43))
	planted := snapshotPath(t, live)

	owner, err := (Factory{TempRoot: base, Probe: func(pid int) error {
		if pid == 42 {
			return syscall.ESRCH
		}
		return nil
	}}).Open(repository)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = owner.Close() }()

	if _, err := os.Stat(dead); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("dead bundle = %v, want absent", err)
	}
	if after := snapshotPath(t, live); !reflect.DeepEqual(planted, after) {
		t.Fatalf("live bundle = %v, want the planted %v", after, planted)
	}
}

// deadProbe is the one probe result that authorizes removal. A row that plants a
// removable candidate uses it, so the row's own defect is the only thing left to deny
// the deletion.
func deadProbe(int) error { return syscall.ESRCH }

// plantedBundle creates one empty same-prefix candidate under base. A row that plants an
// invalid record cannot use writeRecord, which publishes only records the sweep accepts.
func plantedBundle(t *testing.T, base string) string {
	t.Helper()
	root, err := os.MkdirTemp(base, bundlePrefix)
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func plantRecord(t *testing.T, root, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, ownerRecordName), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

func recordBody(t *testing.T, repository string, schema, pid int) string {
	t.Helper()
	common, err := canonicalCommonDir(repository)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf("{\"schema\":%d,\"owner_pid\":%d,\"common_dir\":%q}\n", schema, pid, common)
}

// requireSweepRetains reports whether one Open over base left every candidate exactly as
// its row planted it. The conservative-classification rows differ only in the hostile
// state they plant, so they ask this one question.
func requireSweepRetains(t *testing.T, repository, base string, probe func(int) error, candidates ...string) {
	t.Helper()
	planted := map[string]map[string]string{}
	for _, candidate := range candidates {
		planted[candidate] = snapshotPath(t, candidate)
	}
	owner, err := (Factory{TempRoot: base, Probe: probe}).Open(repository)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = owner.Close() }()
	for _, candidate := range candidates {
		if after := snapshotPath(t, candidate); !reflect.DeepEqual(planted[candidate], after) {
			t.Fatalf("candidate %q = %v, want the planted %v", candidate, after, planted[candidate])
		}
	}
}

// snapshotPath answers what path holds, entry by entry: the mode of every name under it,
// the bytes of every regular file, and the target of every symbolic link. It reads no
// link target and opens no special file, so a hostile candidate is described rather than
// traversed.
func snapshotPath(t *testing.T, path string) map[string]string {
	t.Helper()
	entries := map[string]string{}
	if _, err := os.Lstat(path); errors.Is(err, os.ErrNotExist) {
		return entries
	}
	parent := filepath.Dir(path)
	err := filepath.WalkDir(path, func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(parent, name)
		if err != nil {
			return err
		}
		described := info.Mode().String()
		switch {
		case info.Mode()&os.ModeSymlink != 0:
			target, err := os.Readlink(name)
			if err != nil {
				return err
			}
			described += " -> " + target
		case info.Mode().IsRegular():
			data, err := os.ReadFile(name)
			if err != nil {
				return err
			}
			described += " " + string(data)
		}
		entries[relative] = described
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return entries
}

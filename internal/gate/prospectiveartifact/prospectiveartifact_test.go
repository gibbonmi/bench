package prospectiveartifact

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"

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

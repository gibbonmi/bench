package worktree

import (
	"bytes"
	"errors"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/intent"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
)

// poolRootOperation is one verb path that selects and creates the Bench pool root before
// it uses it. Each pool-root row runs over every operation, so the pooled acquire and the
// owned create cannot diverge on how they treat the root.
type poolRootOperation struct {
	name string
	run  func(t *testing.T, j joins, root, home string) error
}

var poolRootOperations = []poolRootOperation{
	{"acquire", func(t *testing.T, j joins, root, home string) error {
		// The release returns the lease at once, so a re-entry in the same second reuses
		// the clean entry rather than colliding with its own candidate name.
		wt, err := acquireAt(j, root, "", "", home, currentTime())
		releaseWith(j, wt)
		return err
	}},
	{"create", func(t *testing.T, j joins, root, home string) error {
		request, err := randomID()
		if err != nil {
			t.Fatalf("request id: %v", err)
		}
		_, err = createAt(j, root, home, request, "pool-root-"+request[:8], nil, currentTime())
		return err
	}},
}

// poolRootFixture returns a repository, a private home, and the pool root both
// operations select for that pair.
func poolRootFixture(t *testing.T) repoPoolFixture {
	t.Helper()
	root, err := canonicalPath(newWorktreeRepo(t))
	if err != nil {
		t.Fatalf("canonical root: %v", err)
	}
	home := t.TempDir()
	return repoPoolFixture{root: root, home: home, poolRoot: poolAt(home, root)}
}

// failingPoolChmod is the chmodPool join that fails for the pool root alone.
func failingPoolChmod(pool string) func(string, os.FileMode) error {
	return func(path string, mode os.FileMode) error {
		if path == pool {
			return &os.PathError{Op: "chmod", Path: path, Err: syscall.EPERM}
		}
		return os.Chmod(path, mode)
	}
}

// recordingPoolChmod is the chmodPool join that records each path it receives and then
// applies the real mode change.
func recordingPoolChmod(seen *[]string) func(string, os.FileMode) error {
	return func(path string, mode os.FileMode) error {
		*seen = append(*seen, path)
		return os.Chmod(path, mode)
	}
}

// requirePoolRootRefusal asserts the refusal came from the branch that detail names, that
// it names the root and the wanted mode, and that the operation stopped before it put
// anything in entries or registered a checkout. An empty entries skips the content check.
func requirePoolRootRefusal(t *testing.T, err error, detail, root, pool, entries string) {
	t.Helper()
	var refused refusalError
	if !errors.As(err, &refused) {
		t.Fatalf("error = %v, want a pool-root refusal", err)
	}
	if !strings.HasPrefix(refused.detail, detail) {
		t.Fatalf("refusal detail = %q, want the %q branch", refused.detail, detail)
	}
	for _, want := range []string{pool, "0700"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("refusal %q does not name %q", err.Error(), want)
		}
	}
	if listed, readErr := os.ReadDir(entries); entries != "" && (readErr != nil || len(listed) != 0) {
		t.Fatalf("refused pool root %s holds %d entries (%v), want none", entries, len(listed), readErr)
	}
	if list := gitOutput(t, root, "worktree", "list", "--porcelain"); strings.Count(list, "worktree ") != 1 {
		t.Fatalf("a refused operation registered a checkout:\n%s", list)
	}
}

func TestPoolRootChmodFailureRefuses(t *testing.T) {
	t.Parallel()
	for _, op := range poolRootOperations {
		t.Run(op.name, func(t *testing.T) {
			t.Parallel()
			f := poolRootFixture(t)
			j := defaultJoins()
			j.chmodPool = failingPoolChmod(f.poolRoot)
			err := op.run(t, j, f.root, f.home)
			requirePoolRootRefusal(t, err, poolRootModeUnset, f.root, f.poolRoot, f.poolRoot)
			if !strings.Contains(err.Error(), syscall.EPERM.Error()) {
				t.Fatalf("refusal %q drops the join's failure", err.Error())
			}
		})
	}
}

func TestPoolRootSymlinkRefusesWithoutFollowing(t *testing.T) {
	t.Parallel()
	for _, op := range poolRootOperations {
		t.Run(op.name, func(t *testing.T) {
			t.Parallel()
			f := poolRootFixture(t)
			target := t.TempDir()
			mustChmod(t, target, 0o755)
			if err := os.MkdirAll(filepath.Dir(f.poolRoot), 0o700); err != nil {
				t.Fatalf("mkdir pool parent: %v", err)
			}
			if err := os.Symlink(target, f.poolRoot); err != nil {
				t.Fatalf("symlink pool root: %v", err)
			}
			var chmodded []string
			j := defaultJoins()
			j.chmodPool = recordingPoolChmod(&chmodded)
			requirePoolRootRefusal(t, op.run(t, j, f.root, f.home), poolRootSymlink, f.root, f.poolRoot, target)
			if info, err := os.Stat(target); err != nil || info.Mode().Perm() != 0o755 {
				t.Fatalf("symlink target mode = %v (%v), want it untouched at 0755", info, err)
			}
			if len(chmodded) != 0 {
				t.Fatalf("chmodPool reached %q through a symlinked pool root", chmodded)
			}
		})
	}
}

func TestPoolRootNonDirectoryRefuses(t *testing.T) {
	t.Parallel()
	for _, op := range poolRootOperations {
		t.Run(op.name, func(t *testing.T) {
			t.Parallel()
			f := poolRootFixture(t)
			if err := os.MkdirAll(filepath.Dir(f.poolRoot), 0o700); err != nil {
				t.Fatalf("mkdir pool parent: %v", err)
			}
			mustWrite(t, f.poolRoot, []byte("not a directory\n"), 0o644)
			requirePoolRootRefusal(t, op.run(t, defaultJoins(), f.root, f.home), poolRootNotDirectory, f.root, f.poolRoot, "")
		})
	}
}

func TestPoolRootPermissiveIsRetightened(t *testing.T) {
	t.Parallel()
	for _, op := range poolRootOperations {
		for _, failing := range []bool{false, true} {
			name := op.name + "/tightens"
			if failing {
				name = op.name + "/refuses"
			}
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				f := poolRootFixture(t)
				if err := os.MkdirAll(f.poolRoot, 0o700); err != nil {
					t.Fatalf("mkdir pool: %v", err)
				}
				mustChmod(t, f.poolRoot, 0o777)
				j := defaultJoins()
				if failing {
					j.chmodPool = failingPoolChmod(f.poolRoot)
					requirePoolRootRefusal(t, op.run(t, j, f.root, f.home), poolRootModeUnset, f.root, f.poolRoot, f.poolRoot)
					return
				}
				if err := op.run(t, j, f.root, f.home); err != nil {
					t.Fatalf("%s over a permissive pool root: %v", op.name, err)
				}
				requirePoolRootMode(t, f.poolRoot)
			})
		}
	}
}

func TestPoolRootReentryOverTightenedRootSucceeds(t *testing.T) {
	t.Parallel()
	for _, op := range poolRootOperations {
		t.Run(op.name, func(t *testing.T) {
			t.Parallel()
			f := poolRootFixture(t)
			if err := op.run(t, defaultJoins(), f.root, f.home); err != nil {
				t.Fatalf("first %s: %v", op.name, err)
			}
			before := requirePoolRootMode(t, f.poolRoot)
			if err := op.run(t, defaultJoins(), f.root, f.home); err != nil {
				t.Fatalf("re-entered %s over a tightened pool root: %v", op.name, err)
			}
			if after := requirePoolRootMode(t, f.poolRoot); !os.SameFile(before, after) {
				t.Fatal("re-entry replaced the pool root")
			}
		})
	}
}

// requirePoolRootMode asserts the pool root is a real directory at mode 0700.
func requirePoolRootMode(t *testing.T, pool string) os.FileInfo {
	t.Helper()
	info, err := os.Lstat(pool)
	if err != nil {
		t.Fatalf("lstat pool root: %v", err)
	}
	if !info.IsDir() || info.Mode().Perm() != 0o700 {
		t.Fatalf("pool root mode = %v, want a directory at 0700", info.Mode())
	}
	return info
}

func TestCommitmentTicketAssignments(t *testing.T) {
	t.Parallel()
	f := specLessLandingFixture(t, "integration")
	for _, request := range []string{"ticket-one", "ticket-two"} {
		result := runVerb(t, verbCreate, verbCall{root: f.root, home: f.home, args: []string{"--request", request, "--label", request, "--from", f.creation.Assignment.ID}})
		if result.exit != 0 {
			t.Fatalf("create = %d: %s %s", result.exit, result.stdout, result.stderr)
		}
		a, found, err := intent.FindAssignmentForRequest(f.root, request)
		if err != nil || !found {
			t.Fatalf("assignment = %v, %v", found, err)
		}
		if err := (commitrepo.Store{Root: a.Worktree}).Ready("specs/x/spec.md"); err != nil {
			t.Fatalf("sibling cannot deliver: %v", err)
		}
	}
	ledger, err := intent.Read(f.root)
	if err != nil {
		t.Fatal(err)
	}
	if ledger.Commitment == nil || len(ledger.Commitment.Bindings) != 3 || len(ledger.Commitment.Claims) != 1 {
		t.Fatalf("ticket assignments must share one outcome claim: %#v", ledger.Commitment)
	}
}

func TestCommitmentLegacyIdentity(t *testing.T) {
	t.Parallel()
	root, home := newWorktreeRepo(t), t.TempDir()
	listed := mustCreate(t, root, home, "listed", "listed")
	other := mustCreate(t, root, home, "other", "other")
	err := intent.Transact(root, intent.StrictRead, func(l intent.Ledger) (intent.Ledger, bool, error) {
		l.Commitment = &intent.CommitmentState{Continuations: []intent.LegacyContinuation{{Assignment: listed.Assignment.ID, Request: listed.Assignment.Request, Scope: []string{"owned.txt"}}}}
		return l, true, nil
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	scope, err := (commitrepo.Store{Root: listed.Path}).LegacyScope("listed")
	if err != nil || !reflect.DeepEqual(scope, []string{"owned.txt"}) {
		t.Fatalf("listed scope = %v, %v", scope, err)
	}
	address, err := intent.Address(root)
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(address)
	if err != nil {
		t.Fatal(err)
	}
	for _, request := range []string{"listed", "other"} {
		if _, err := (commitrepo.Store{Root: other.Path}).LegacyScope(request); err == nil {
			t.Fatalf("unlisted assignment borrowed continuation with %q", request)
		}
	}
	after, err := os.ReadFile(address)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatalf("refusal changed ledger: %v", err)
	}
}

package worktree

import (
	"errors"
	"os"
	"path/filepath"
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
func poolRootFixture(t *testing.T) (root, home, pool string) {
	t.Helper()
	root, err := canonicalPath(newWorktreeRepo(t))
	if err != nil {
		t.Fatalf("canonical root: %v", err)
	}
	home = t.TempDir()
	return root, home, poolAt(home, root)
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

// requirePoolRootRefusal asserts the refusal names the root and the wanted mode, and
// that the operation stopped before it put anything in the root or registered a checkout.
func requirePoolRootRefusal(t *testing.T, err error, root, pool, entries string) {
	t.Helper()
	var refused refusalError
	if !errors.As(err, &refused) {
		t.Fatalf("error = %v, want a pool-root refusal", err)
	}
	for _, want := range []string{pool, "0700"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("refusal %q does not name %q", err.Error(), want)
		}
	}
	if listed, readErr := os.ReadDir(entries); readErr != nil || len(listed) != 0 {
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
			root, home, pool := poolRootFixture(t)
			j := defaultJoins()
			j.chmodPool = failingPoolChmod(pool)
			requirePoolRootRefusal(t, op.run(t, j, root, home), root, pool, pool)
		})
	}
}

func TestPoolRootSymlinkRefusesWithoutFollowing(t *testing.T) {
	t.Parallel()
	for _, op := range poolRootOperations {
		t.Run(op.name, func(t *testing.T) {
			t.Parallel()
			root, home, pool := poolRootFixture(t)
			target := t.TempDir()
			mustChmod(t, target, 0o755)
			if err := os.MkdirAll(filepath.Dir(pool), 0o700); err != nil {
				t.Fatalf("mkdir pool parent: %v", err)
			}
			if err := os.Symlink(target, pool); err != nil {
				t.Fatalf("symlink pool root: %v", err)
			}
			requirePoolRootRefusal(t, op.run(t, defaultJoins(), root, home), root, pool, target)
			if info, err := os.Stat(target); err != nil || info.Mode().Perm() != 0o755 {
				t.Fatalf("symlink target mode = %v (%v), want it untouched at 0755", info, err)
			}
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
				root, home, pool := poolRootFixture(t)
				if err := os.MkdirAll(pool, 0o700); err != nil {
					t.Fatalf("mkdir pool: %v", err)
				}
				mustChmod(t, pool, 0o777)
				j := defaultJoins()
				if failing {
					j.chmodPool = failingPoolChmod(pool)
					requirePoolRootRefusal(t, op.run(t, j, root, home), root, pool, pool)
					return
				}
				if err := op.run(t, j, root, home); err != nil {
					t.Fatalf("%s over a permissive pool root: %v", op.name, err)
				}
				requirePoolRootMode(t, pool)
			})
		}
	}
}

func TestPoolRootReentryOverTightenedRootSucceeds(t *testing.T) {
	t.Parallel()
	for _, op := range poolRootOperations {
		t.Run(op.name, func(t *testing.T) {
			t.Parallel()
			root, home, pool := poolRootFixture(t)
			if err := op.run(t, defaultJoins(), root, home); err != nil {
				t.Fatalf("first %s: %v", op.name, err)
			}
			before := requirePoolRootMode(t, pool)
			if err := op.run(t, defaultJoins(), root, home); err != nil {
				t.Fatalf("re-entered %s over a tightened pool root: %v", op.name, err)
			}
			if after := requirePoolRootMode(t, pool); !os.SameFile(before, after) {
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

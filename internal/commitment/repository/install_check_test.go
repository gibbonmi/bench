package repository_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/brokermanifest"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/spec"
)

// installedRepository is a repository that declares Bench build inputs, whose primary broker
// manifest binds the installed executable. A linked worktree holds a build of other bytes,
// and its own wrapper directory holds a manifest that binds that build.
type installedRepository struct {
	primary, installed, worktree, worktreeBuild string
}

func newInstalledRepository(t *testing.T) installedRepository {
	t.Helper()
	primary := gittest.RepoOnBranch(t, "main")
	commitmenttest.Write(t, primary, freshness.BuildInputsManifest, "cmd\n")
	commitmenttest.Commit(t, primary, "declare the build inputs")
	worktree := filepath.Join(t.TempDir(), "worktree")
	gittest.Output(t, primary, "worktree", "add", "-q", "-b", "candidate", worktree)
	r := installedRepository{primary: primary, installed: filepath.Join(primary, "dist", "bench"), worktree: worktree, worktreeBuild: filepath.Join(worktree, "dist", "bench")}
	writeExecutable(t, r.installed, "installed candidate")
	writeExecutable(t, r.worktreeBuild, "rebuilt candidate")
	bindManifest(t, primary, r.installed)
	bindManifest(t, worktree, r.worktreeBuild)
	return r
}

func writeExecutable(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatal(err)
	}
}

// bindManifest writes the broker manifest of the checkout at root, beside its wrapper, so it
// binds executable.
func bindManifest(t *testing.T, root, executable string) {
	t.Helper()
	dir := filepath.Join(root, "bin")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, _, err := brokermanifest.Write(dir, executable, "0.0.0"); err != nil {
		t.Fatal(err)
	}
}

// TestCandidateCheck is CB26, CB27, CB28, CB60, and CB65, with one row for each other
// refusal of the check and for a repository without Bench build inputs. The declared path is
// one that no ledger type holds, so a check that tested field knowledge would refuse the
// installed candidate.
func TestCandidateCheck(t *testing.T) {
	declared := []spec.PersistedField{{Record: "intent-ledger", Path: []string{"commitment", "fieldnoledgerholds"}}}
	for _, tc := range []struct {
		name        string
		declaration []spec.PersistedField
		call        func(t *testing.T, r installedRepository) (executable, root string)
		refused     bool
	}{
		{name: "the installed candidate passes", declaration: declared, call: func(_ *testing.T, r installedRepository) (string, string) { return r.installed, r.primary }},
		{name: "another digest refuses", declaration: declared, refused: true, call: func(_ *testing.T, r installedRepository) (string, string) { return r.worktreeBuild, r.primary }},
		{name: "an undeclared spec passes any digest", call: func(_ *testing.T, r installedRepository) (string, string) { return r.worktreeBuild, r.primary }},
		{name: "a worktree root reads the primary manifest", declaration: declared, refused: true, call: func(_ *testing.T, r installedRepository) (string, string) { return r.worktreeBuild, r.worktree }},
		{name: "no primary manifest refuses", declaration: declared, refused: true, call: func(t *testing.T, r installedRepository) (string, string) {
			if err := os.Remove(filepath.Join(r.primary, "bin", brokermanifest.Name)); err != nil {
				t.Fatal(err)
			}
			return r.installed, r.primary
		}},
		{name: "a primary manifest that does not read refuses", declaration: declared, refused: true, call: func(t *testing.T, r installedRepository) (string, string) {
			writeExecutable(t, filepath.Join(r.primary, "bin", brokermanifest.Name), "not a manifest\n")
			return r.installed, r.primary
		}},
		{name: "an executable that does not digest refuses", declaration: declared, refused: true, call: func(_ *testing.T, r installedRepository) (string, string) {
			return filepath.Join(r.primary, "absent"), r.primary
		}},
		{name: "a root outside a repository refuses", declaration: declared, refused: true, call: func(t *testing.T, r installedRepository) (string, string) { return r.installed, t.TempDir() }},
		{name: "a repository without build inputs passes", declaration: declared, call: func(t *testing.T, r installedRepository) (string, string) {
			if err := os.Remove(filepath.Join(r.primary, "scripts", "go-build.inputs")); err != nil {
				t.Fatal(err)
			}
			return r.worktreeBuild, r.primary
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			executable, root := tc.call(t, newInstalledRepository(t))
			err := commitrepo.CheckCandidate(tc.declaration, executable, root)
			var raised refusalroute.Raised
			if refused := errors.As(err, &raised) && raised.Name == "candidate-not-installed"; refused != tc.refused || !refused && err != nil {
				t.Fatalf("check = %v; want a candidate-not-installed refusal %t", err, tc.refused)
			}
		})
	}
}

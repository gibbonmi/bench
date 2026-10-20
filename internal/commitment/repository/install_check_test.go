package repository_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/adopt"
	"github.com/gibbonmi/bench/internal/brokermanifest"
	"github.com/gibbonmi/bench/internal/commitment"
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
	return installAt(t, gittest.RepoOnBranch(t, "main"))
}

// installAt makes the repository of the primary checkout an installedRepository.
func installAt(t *testing.T, primary string) installedRepository {
	t.Helper()
	commitmenttest.Write(t, primary, freshness.BuildInputsManifest, "cmd\n")
	commitmenttest.Commit(t, primary, "declare the build inputs")
	worktree := filepath.Join(t.TempDir(), "worktree")
	gittest.Output(t, primary, "worktree", "add", "-q", "-b", "candidate", worktree)
	r := installedRepository{primary: primary, installed: filepath.Join(primary, "dist", "bench"), worktree: worktree, worktreeBuild: filepath.Join(worktree, "dist", "bench")}
	writeExecutable(t, r.installed, "installed candidate")
	writeExecutable(t, r.worktreeBuild, "rebuilt candidate")
	commitmenttest.BindManifest(t, primary, r.installed)
	commitmenttest.BindManifest(t, worktree, r.worktreeBuild)
	return r
}

// newOverlappingBuilds is an installedRepository whose ledger holds two open declared
// builds, one of each outcome.
func newOverlappingBuilds(t *testing.T) installedRepository {
	t.Helper()
	outcomes := []string{"A", "B"}
	root := commitmenttest.StagedInParallel(t, outcomes...)
	body := commitmenttest.StagedBody + "Persisted fields: intent-ledger:commitment.gaps\n"
	commitmenttest.EditPolicy(t, root, func(policy *commitment.Policy) {
		for i := range policy.Milestones[0].Outcomes {
			policy.Milestones[0].Outcomes[i].Deliverables[0].Source.Identity = commitment.Identity([]byte(body))
		}
	})
	for _, outcome := range outcomes {
		commitmenttest.Write(t, root, approvedSpec(t, root, outcome), body)
	}
	commitmenttest.Commit(t, root, "declare a persisted field in each spec")
	r := installAt(t, root)
	for _, outcome := range outcomes {
		startBuild(t, root, "build-"+outcome, outcome, approvedSpec(t, root, outcome))
	}
	return r
}

// wrapperAt is the wrapper of the checkout at root, which shares the directory of that
// checkout's broker manifest.
func wrapperAt(root string) string {
	return filepath.Join(filepath.Dir(brokermanifest.Path(root)), "bench.sh")
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

// TestCandidateCheck is CB26, CB27, CB28, CB60, CB65, and CB69, with one row for each other
// refusal of the check and for a repository without Bench build inputs. The declared path is
// one that no ledger type holds, so a check that tested field knowledge would refuse the
// installed candidate. One row publishes the manifest through the adoption writer beside the
// primary checkout's wrapper, so a check that reads another location refuses that row. One
// row points BENCH_WRAPPER into the worktree, whose manifest binds the running build, so a
// check that reads the manifest beside BENCH_WRAPPER passes that row. An empty face is a pass.
func TestCandidateCheck(t *testing.T) {
	declared := []spec.PersistedField{{Record: "intent-ledger", Path: []string{"commitment", "fieldnoledgerholds"}}}
	const (
		notInstalled = "candidate-not-installed"
		overlap      = "declared-builds-overlap"
	)
	for _, tc := range []struct {
		name        string
		repository  func(t *testing.T) installedRepository
		declaration []spec.PersistedField
		call        func(t *testing.T, r installedRepository) (executable, root string)
		face        string
	}{
		{name: "the installed candidate passes", declaration: declared, call: func(_ *testing.T, r installedRepository) (string, string) { return r.installed, r.primary }},
		{name: "the manifest of the production writer binds the candidate", declaration: declared, call: func(t *testing.T, r installedRepository) (string, string) {
			t.Setenv("BENCH_WRAPPER", wrapperAt(r.primary))
			if _, _, err := adopt.WriteBrokerManifest("0.0.0"); err != nil {
				t.Fatal(err)
			}
			executable, err := os.Executable()
			if err != nil {
				t.Fatal(err)
			}
			return executable, r.primary
		}},
		{name: "another digest refuses", declaration: declared, face: notInstalled, call: func(_ *testing.T, r installedRepository) (string, string) { return r.worktreeBuild, r.primary }},
		{name: "an undeclared spec passes any digest", call: func(_ *testing.T, r installedRepository) (string, string) { return r.worktreeBuild, r.primary }},
		{name: "a worktree root reads the primary manifest", declaration: declared, face: notInstalled, call: func(_ *testing.T, r installedRepository) (string, string) { return r.worktreeBuild, r.worktree }},
		{name: "a worktree wrapper does not move the manifest", declaration: declared, face: notInstalled, call: func(t *testing.T, r installedRepository) (string, string) {
			t.Setenv("BENCH_WRAPPER", wrapperAt(r.worktree))
			return r.worktreeBuild, r.worktree
		}},
		{name: "no primary manifest refuses", declaration: declared, face: notInstalled, call: func(t *testing.T, r installedRepository) (string, string) {
			if err := os.Remove(brokermanifest.Path(r.primary)); err != nil {
				t.Fatal(err)
			}
			return r.installed, r.primary
		}},
		{name: "a primary manifest that does not read refuses", declaration: declared, face: notInstalled, call: func(t *testing.T, r installedRepository) (string, string) {
			writeExecutable(t, brokermanifest.Path(r.primary), "not a manifest\n")
			return r.installed, r.primary
		}},
		{name: "an executable that does not digest refuses", declaration: declared, face: notInstalled, call: func(_ *testing.T, r installedRepository) (string, string) {
			return filepath.Join(r.primary, "absent"), r.primary
		}},
		{name: "a root outside a repository refuses", declaration: declared, face: notInstalled, call: func(t *testing.T, r installedRepository) (string, string) { return r.installed, t.TempDir() }},
		{name: "a repository without build inputs passes", declaration: declared, call: func(t *testing.T, r installedRepository) (string, string) {
			if err := os.Remove(filepath.Join(r.primary, filepath.FromSlash(freshness.BuildInputsManifest))); err != nil {
				t.Fatal(err)
			}
			return r.worktreeBuild, r.primary
		}},
		{name: "two open declared builds refuse the installed candidate", repository: newOverlappingBuilds, declaration: declared, face: overlap, call: func(_ *testing.T, r installedRepository) (string, string) {
			return r.installed, r.primary
		}},
		{name: "two open declared builds refuse before the digest test", repository: newOverlappingBuilds, declaration: declared, face: overlap, call: func(_ *testing.T, r installedRepository) (string, string) {
			return r.worktreeBuild, r.primary
		}},
		{name: "a ledger that does not read skips the count", repository: newOverlappingBuilds, declaration: declared, call: func(t *testing.T, r installedRepository) (string, string) {
			commitmenttest.PlantCommitmentField(t, r.primary, "fixture")
			return r.installed, r.primary
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repository := newInstalledRepository
			if tc.repository != nil {
				repository = tc.repository
			}
			executable, root := tc.call(t, repository(t))
			err := commitrepo.CheckCandidate(tc.declaration, executable, root)
			var raised refusalroute.Raised
			errors.As(err, &raised)
			if raised.Name != tc.face || tc.face == "" && err != nil {
				t.Fatalf("check = %v; want the refusal face %q, or a pass when it is empty", err, tc.face)
			}
			if tc.face != overlap {
				return
			}
			if route := refusalroute.Printed("", err, nil).Route; !strings.Contains(route, "bench commitment block") {
				t.Fatalf("overlap route = %q, want it to name bench commitment block", route)
			}
		})
	}
}

// TestCheckEntryReadsNoDeclarationWithoutBuildInputs proves that a repository without Bench
// build inputs keeps the declaration inert: the entry passes a spec at HEAD whose
// declaration does not parse.
func TestCheckEntryReadsNoDeclarationWithoutBuildInputs(t *testing.T) {
	const path = "specs/malformed/spec.md"
	root := gittest.RepoOnBranch(t, "main")
	commitmenttest.Write(t, root, path, "# malformed\n\nStatus: staged\nPersisted fields: gate-evidence:x\n")
	commitmenttest.Commit(t, root, "a malformed declaration")
	if err := commitrepo.CheckEntry(root, path); err != nil {
		t.Fatalf("entry without build inputs = %v, want a pass", err)
	}
}

package runtime

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

func TestRuntimeFreshInstall(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	contract.RunParallel(t, "fresh-install PATH-stripped gate resolves local CLI", testRuntimeFreshInstallGateResolves)
}

// testRuntimeFreshInstallGateResolves is the only test that installs cold with no global
// bench on PATH — the class the kit otherwise can't see. It links and inits a throwaway
// repo, then runs the scaffolded .bench/gate.sh, all under a PATH stripped to
// /usr/bin:/bin (git and bash present, no global bench). The PATH strip is load-bearing:
// a smoke that leaves a global bench on PATH false-greens against the old bare-`bench`
// scaffold, so it is asserted as part of the fixture, not left implicit.
func testRuntimeFreshInstallGateResolves(t *testing.T) {
	repo := t.TempDir()
	f := contract.NewFixtureAt(t, repo, contract.IsolatedEnv(t, repo))
	f.Git("init", "-q")

	stripped := map[string]string{"PATH": "/usr/bin:/bin"}
	contract.RunAt(t, f, repo, stripped, "bash", benchPath(t), "link").RequireExit(0)
	contract.RunAt(t, f, repo, stripped, "bash", benchPath(t), "init").RequireExit(0)

	// The fresh install must land the local CLI and the copied binary the scaffolded gate
	// resolves, so the resolution assertion below tests resolution — not mere absence.
	requireFreshFile(t, filepath.Join(repo, ".bench", "bin", "bench.sh"), "fresh link did not install .bench/bin/bench.sh")
	requireFreshFile(t, filepath.Join(repo, ".bench", "dist", "bench"), "fresh link did not copy .bench/dist/bench")

	gate := contract.RunAt(t, f, repo, stripped, "bash", filepath.Join(repo, ".bench", "gate.sh"))
	out := gate.Stdout + gate.Stderr
	// The scaffolded gate is intentionally red on its sentinel; the contract is that it
	// RESOLVES the local CLI and reaches canary, not that it goes green. A bare-`bench`
	// scaffold emits "bench: command not found" and never runs the canary under this PATH.
	contract.RequireNotContains(t, out, "command not found")
	contract.RequireNotContains(t, out, "canary sweep failed")

	// story 5: the profile templates the fresh install points at are reachable in the
	// resolved kit — the ships-regression itself is owned by the package-surface pin.
	requireFreshFile(t, filepath.Join(contract.SubjectRoot(t), "projects", "gl-axi.md"), "resolved kit dir lacks projects/ profile templates")
}

func requireFreshFile(t *testing.T, path, msg string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
}

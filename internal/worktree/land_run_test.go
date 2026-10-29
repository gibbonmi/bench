// Run-sharing test for the landing command: one landing reads the destination through the
// verb's run.
package worktree

import (
	"os"
	"path/filepath"
	"testing"
)

// landingDestinationCommonDirStarts is the observed number of destination common-directory
// starts of one landing. The readers that share the verb's run start the query once between
// them, and each reader that does not share it starts the query on its own.
const landingDestinationCommonDirStarts = 31

// TestLandingReadsTheDestinationCommonDirOnce proves that one in-process landing starts the
// destination common-directory query exactly landingDestinationCommonDirStarts times. The log
// stub on the process PATH records each start and passes it to the real git. A store or a
// helper that drops the run starts the query again, so the count rises above the pin. As
// context only: before the landing shared its run, the same landing started the query 65
// times. (Coverage row GQ34.)
func TestLandingReadsTheDestinationCommonDirOnce(t *testing.T) {
	request := "land-run-common-dir"
	f := publicLandingFixture(t, request, "", "")
	logPath, path := filepath.Join(t.TempDir(), "argv"), os.Getenv("PATH")
	journeyStubGit(t, f.root, "log", logPath)
	// The gate resolves its declared tools on PATH, so the stub goes ahead of the real PATH.
	bindEnv(t, "PATH", os.Getenv("PATH")+string(os.PathListSeparator)+path)
	if r := runVerb(t, verbLand, f.call(landArgs(request, f.base, f.tip, f.creation.Path)...)); r.exit != 0 {
		t.Fatalf("land exit = %d, want 0; stdout=%q stderr=%q", r.exit, r.stdout, r.stderr)
	}
	query := "-C " + f.root + " rev-parse --path-format=absolute --git-common-dir"
	if starts := journeyStubGitStarts(t, logPath, query); starts != landingDestinationCommonDirStarts {
		t.Fatalf("one landing started the destination common-directory query %d times, want %d", starts, landingDestinationCommonDirStarts)
	}
}

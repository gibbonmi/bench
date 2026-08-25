package worktree

// The worktree race sentinel. The race-test registry names this file's one test,
// so the gate's race phase runs two parallel journeys under -race and reds on a
// harness data race. (Coverage row WF07.)

import (
	"bytes"
	"strings"
	"testing"
)

// TestParallelJourneysShareTheHarnessSafely is WF07. Two parallel journeys each
// run a real Bench worktree creation against their own repository, and the
// harness effect log holds both journeys' starts afterwards. The child carries
// its own BENCH_HOME on Env, so the test binds no process environment and runs
// in parallel with its siblings.
func TestParallelJourneysShareTheHarnessSafely(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	journeys := []string{"race-sentinel-a", "race-sentinel-b"}
	// The inner run returns only after every parallel child journey has finished,
	// so the effect log below is complete.
	t.Run("journeys", func(t *testing.T) {
		for _, journey := range journeys {
			t.Run(journey, func(t *testing.T) {
				t.Parallel()
				root := newWorktreeRepo(t)
				var stdout, stderr bytes.Buffer
				create := descendant(t, binary, "worktree", "create", "--request", journey, "--label", journey)
				create.Dir = root
				create.Env = journeyChildEnv(t, t.TempDir())
				create.Stdout, create.Stderr = &stdout, &stderr
				if err := create.Run(); err != nil {
					t.Fatalf("worktree create exit=%d stdout=%q stderr=%q", exitCode(err), stdout.String(), stderr.String())
				}
				if !strings.Contains(stdout.String(), "worktree_create[1]") || !strings.Contains(stdout.String(), ",active") {
					t.Fatalf("worktree create stdout = %q, want one active creation receipt", stdout.String())
				}
			})
		}
	})
	log := journeyEffectLog()
	for _, journey := range journeys {
		want := "descendant " + t.Name() + "/journeys/" + journey + " " + binary
		found := false
		for _, record := range log {
			if record == want {
				found = true
			}
		}
		if !found {
			t.Fatalf("harness effect log lacks %q; parallel journeys lost an effect record", want)
		}
	}
}

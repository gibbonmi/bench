package worktree

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/capability"
)

// TestLandSuppliesTheDeclaredBenchHomeItself is the FT315 regression. A gate that
// declares BENCH_HOME must land from a shell that exports only HOME, because the
// binary derives the Bench home the way the shell wrapper does. Before the fix, the
// green gate refused as "declared environment unavailable" and named no variable. A
// declared name that nothing derives still refuses, and the refusal names it.
func TestLandSuppliesTheDeclaredBenchHomeItself(t *testing.T) {
	t.Parallel()
	binary := testRunBinary(t)
	for _, tc := range []struct {
		name, extra, refusal string
	}{
		{name: "derived"},
		{name: "unsupplied", extra: "FT315_UNSUPPLIED", refusal: "infrastructure (declared environment unavailable: FT315_UNSUPPLIED); run bench doctor"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			userHome := t.TempDir()
			home := filepath.Join(userHome, ".bench")
			request := "land-bench-home-" + tc.name
			root, creation, _, _, tally := landingFixtureAtHome(t, request, "", "", home, false)
			declared := []string{benchhome.Env}
			if tc.extra != "" {
				declared = append(declared, tc.extra)
			}
			f := landingGateFixture(t, declared...)
			seen := tally + ".bench-home"
			check := "printf '%s' \"${BENCH_HOME-unset}\" > '" + seen + "'\n[ -f owned.txt ]\nprintf g >> '" + tally + "'\n"
			f.MustWrite(t, root, "set -eu\n"+check, "set -eu\nruntime=$1\n"+check)
			gitRun(t, root, "add", ".bench/gate.sh", ".bench/gate-prospective.sh", ".bench/gate-inputs.json")
			gitRun(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "declare the bench home")
			base := gitOutput(t, root, "rev-parse", "HEAD")
			gitRun(t, creation.Path, "rebase", "main")
			refreshLandingEvidence(t, creation.Path, base)
			tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")

			var stdout, stderr bytes.Buffer
			cmd := descendant(t, binary, append([]string{"worktree", "land"}, specLessLandArgs(request, base, tip, creation.Path)...)...)
			cmd.Dir, cmd.Stdout, cmd.Stderr = root, &stdout, &stderr
			environment := os.Environ()
			for _, name := range append(declared, "HOME") {
				environment = capability.WithoutEnvironment(environment, name)
			}
			cmd.Env = append(environment, "HOME="+userHome)
			code := exitCode(cmd.Run())
			if tc.refusal != "" {
				if code != 1 || !strings.Contains(stdout.String(), tc.refusal) {
					t.Fatalf("land with an unsupplied declared name = (%d, %q, %q)", code, stdout.String(), stderr.String())
				}
				return
			}
			if code != 0 || !strings.Contains(stdout.String(), "worktree=released") {
				t.Fatalf("land without an exported BENCH_HOME = (%d, %q, %q)", code, stdout.String(), stderr.String())
			}
			if got, err := os.ReadFile(tally); err != nil || string(got) != "g" {
				t.Fatalf("gate tally = %q, %v", got, err)
			}
			if got, err := os.ReadFile(seen); err != nil || string(got) != home {
				t.Fatalf("gate saw BENCH_HOME = %q, %v; want %q", got, err, home)
			}
		})
	}
}

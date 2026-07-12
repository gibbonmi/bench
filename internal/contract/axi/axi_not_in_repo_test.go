package axi

import (
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

// notInRepoPhrase is the one rendered sentence toon.NotInRepo() emits. The Go call
// sites render it through that function; .bench/gate.sh, a shell file, cannot call
// Go and carries the literal text instead — this probe pins the literal (what a
// black-box caller actually observes on stderr) rather than importing toon.
const notInRepoPhrase = "error: not in a git repository — run inside a Bench-linked repo"

// TestOperationalNotInRepoOnePhrase sweeps every operational command that renders a
// not-in-repo message and asserts each prints the one phrase on its existing channel
// (stderr) with its existing exit code — the operational half of the one-phrase
// contract (AXI query commands keep their own pinned stdout/exit-1 posture and are
// out of this table). `bench gate pin` refuses non-TTY stdin before it reaches its
// not-in-repo branch, so a non-interactive probe can never observe that branch; it is
// pinned instead at the Go seam in internal/gate/pin_test.go
// (TestPinCommandNotInRepo), which injects isTerminal=true to bypass the unrelated
// TTY precondition.
func TestOperationalNotInRepoOnePhrase(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)

	cases := []struct {
		name     string
		wantExit int
		run      func(t *testing.T, f contract.Fixture) contract.Probe
	}{
		{
			name:     "bench gate",
			wantExit: 1,
			run: func(t *testing.T, f contract.Fixture) contract.Probe {
				return f.Bench("gate")
			},
		},
		{
			name:     "bench gate-run (plumbing)",
			wantExit: 1,
			run: func(t *testing.T, f contract.Fixture) contract.Probe {
				return f.Bench("gate-run")
			},
		},
		{
			name:     "bench gate-phases (plumbing)",
			wantExit: 3,
			run: func(t *testing.T, f contract.Fixture) contract.Probe {
				return f.Bench("gate-phases")
			},
		},
		{
			name:     ".bench/gate.sh (direct, shell literal)",
			wantExit: 3,
			run: func(t *testing.T, f contract.Fixture) contract.Probe {
				gateSh := filepath.Join(contract.SubjectRoot(t), ".bench", "gate.sh")
				return contract.RunAt(t, f, f.Root, nil, "bash", gateSh)
			},
		},
		{
			name:     "bench canary",
			wantExit: 1,
			run: func(t *testing.T, f contract.Fixture) contract.Probe {
				return f.Bench("canary")
			},
		},
		{
			name:     "bench commit",
			wantExit: 1,
			run: func(t *testing.T, f contract.Fixture) contract.Probe {
				return f.Bench("commit", "-m", "msg", "some.txt")
			},
		},
		{
			name:     "bench worktree",
			wantExit: 1,
			run: func(t *testing.T, f contract.Fixture) contract.Probe {
				return f.Bench("worktree")
			},
		},
		{
			name:     "bench worktree clean <target>",
			wantExit: 1,
			run: func(t *testing.T, f contract.Fixture) contract.Probe {
				return f.Bench("worktree", "clean", filepath.Join(f.Root, "target"))
			},
		},
		{
			name:     "bench shift",
			wantExit: 1,
			run: func(t *testing.T, f contract.Fixture) contract.Probe {
				return f.BenchEnv(map[string]string{"BENCH_AGENT": "true"}, "shift", "an objective")
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			f := contract.NewFixture(t, contract.WithNoRepo())
			out := c.run(t, f)
			out.RequireExit(c.wantExit)
			out.RequireContains(out.Stderr, notInRepoPhrase)
			if out.Stdout != "" {
				t.Fatalf("expected no stdout for an operational not-in-repo failure (stderr is the contract channel)\nstdout:\n%s", out.Stdout)
			}
		})
	}
}

// The exact-target grammar is validated before repository discovery. This keeps a
// missing target a usage error even when the caller also happens to be outside a
// repository; supplying a target reaches the operational guard pinned above.
func TestWorktreeCleanUsagePrecedesRepoLookup(t *testing.T) {
	t.Parallel()
	contract.SkipIfSubjectBenchMissing(t)
	f := contract.NewFixture(t, contract.WithNoRepo())

	out := f.Bench("worktree", "clean")

	out.RequireExit(2)
	want := "worktree_cleanup[1]{target,action,tracked,ignored,recovery,fingerprint,detail}:\n  unknown,error,unknown,unknown,none,none,\"invalid invocation; run bench worktree clean <path> [--apply <fingerprint>]\"\n"
	if out.Stdout != want || out.Stderr != "" {
		t.Fatalf("usage streams = stdout %q stderr %q", out.Stdout, out.Stderr)
	}
	out.RequireNotContains(out.Stdout+out.Stderr, notInRepoPhrase)
}

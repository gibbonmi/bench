package runtime

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/contract"
)

func TestRuntimeShiftOutcomeContracts(t *testing.T) {
	t.Parallel()
	contract.RunParallel(t, "bench shift wall deadline contract", testShiftWallDeadlineKillsAdapter)
	contract.RunParallel(t, "bench shift wall deadline mid-gate contract", testShiftWallDeadlineCancelsRunningGate)
	contract.RunParallel(t, "bench shift intent and status recovery surfacing contract", testShiftIntentAndStatusSurfaceRecovery)
	contract.RunParallel(t, "bench shift interrupt-with-mutation recovery contract", testShiftInterruptWithMutationRecovers)
	contract.RunParallel(t, "bench shift no-objective usage contract", testShiftNoObjective)
	contract.RunParallel(t, "bench shift whitespace-only objective usage contract", testShiftWhitespaceObjective)
	contract.RunParallel(t, "bench shift control-byte objective usage contract", testShiftControlByteObjective)
	contract.RunParallel(t, "bench shift invalid env usage contract", testShiftInvalidEnvUsage)
	contract.RunParallel(t, "bench shift empty env default contract", testShiftEmptyEnvDefault)
	contract.RunParallel(t, "bench shift max-wall inclusive bound contract", testShiftMaxWallInclusiveBound)
	contract.RunParallel(t, "bench shift complete result contract", testShiftCompleteResult)
	contract.RunParallel(t, "bench shift adapter-failure zero-commit contract", testShiftAdapterFailureZeroCommit)
	contract.RunParallel(t, "bench shift adapter spawn-failure contract", testShiftAdapterSpawnFailure)
	contract.RunParallel(t, "bench shift adapter-failure after-commit contract", testShiftAdapterFailureAfterCommit)
	contract.RunParallel(t, "bench shift cap-exhaustion contract", testShiftCapExhaustion)
	contract.RunParallel(t, "bench shift no-op done contract", testShiftNoOpDone)
}

// testShiftWallDeadlineKillsAdapter covers row 8: a sleeping adapter under a tiny
// BENCH_MAX_WALL is killed, its prior mutation is snapshotted, the shift exits
// incomplete/3 with a deadline detail, and the adapter does not keep running past exit.
func testShiftWallDeadlineKillsAdapter(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteExecutable("agent", "#!/usr/bin/env bash\nprintf 'partial\\n' > partial.txt\nsleep 3\nprintf 'late\\n' > late.txt\n")
	f.CommitAll("agent")
	home := t.TempDir()

	start := time.Now()
	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "2", "BENCH_MAX_WALL": "1s", "BENCH_HOME": home}, "shift", "wall deadline")
	elapsed := time.Since(start)

	probe.RequireExit(3)
	if elapsed > 15*time.Second {
		t.Fatalf("shift did not honor the wall deadline promptly: %s elapsed", elapsed)
	}
	probe.RequireContains(probe.Stdout, "deadline")
	wt := shiftWorktree(t, probe.Stdout)
	branch := shiftBranchFromStart(t, probe.Stdout)
	waitSeconds(t, 4)
	if _, err := os.Stat(filepath.Join(wt, "late.txt")); err == nil {
		t.Fatal("adapter kept running after the wall deadline and dirtied the pooled worktree")
	}
	ref := "refs/bench/recovery/" + branch
	f.Git("show-ref", "--verify", ref)
	requireNoLease(t, home)
}

// testShiftWallDeadlineCancelsRunningGate covers row 8's other half: a wall deadline
// that fires while the gate itself is running cancels the gate instead of waiting it
// out — row 8's sleeping-adapter case only exercises the wall against a hung adapter,
// so a wall that never wires cancel-gate would pass that row yet hang here.
func testShiftWallDeadlineCancelsRunningGate(t *testing.T) {
	f := shiftFixture(t, `#!/usr/bin/env bash
if [ ! -f "$BENCH_TEST_STATE/gate-ran-once" ]; then
  : > "$BENCH_TEST_STATE/gate-ran-once"
  sleep 5
  printf 'late\n' > "$BENCH_TEST_STATE/late-gate-finished"
fi
exit 0
`)
	f.WriteExecutable("agent", "#!/usr/bin/env bash\nprintf 'work\\n' > work.txt\n")
	f.CommitAll("agent")
	home := t.TempDir()
	state := t.TempDir()

	start := time.Now()
	probe := f.BenchEnv(map[string]string{
		"BENCH_TEST_STATE": state, "BENCH_AGENT": filepath.Join(f.Root, "agent"),
		"BENCH_MAX_ITERS": "2", "BENCH_MAX_WALL": "1s", "BENCH_HOME": home,
	}, "shift", "wall mid gate")
	elapsed := time.Since(start)

	probe.RequireExit(3)
	if elapsed > 15*time.Second {
		t.Fatalf("shift did not honor the wall deadline while the gate was running: %s elapsed", elapsed)
	}
	probe.RequireContains(probe.Stdout, "deadline")
	branch := shiftBranchFromStart(t, probe.Stdout)
	waitSeconds(t, 5)
	if _, err := os.Stat(filepath.Join(state, "late-gate-finished")); err == nil {
		t.Fatal("gate kept running past the wall deadline instead of being cancelled")
	}
	ref := "refs/bench/recovery/" + branch
	f.Git("show-ref", "--verify", ref)
	tree := f.Git("ls-tree", "-r", "--name-only", ref).Stdout
	if !strings.Contains(tree, "work.txt") {
		t.Fatalf("recovery snapshot did not preserve the adapter's mutation before the gate-cancelling deadline:\n%s", tree)
	}
	requireNoLease(t, home)
}

// testShiftIntentAndStatusSurfaceRecovery covers row 17: a preserving shift's outcome
// and recovery pointer land on its intent-ledger entry, and bench status renders the
// pointer, so a preserved shift is discoverable after the terminal is gone.
func testShiftIntentAndStatusSurfaceRecovery(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 1\n")
	f.WriteExecutable("agent", "#!/usr/bin/env bash\nprintf 'work\\n' > work.txt\n")
	f.CommitAll("agent")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "1", "BENCH_HOME": home}, "shift", "surface recovery")
	probe.RequireExit(1)
	branch := shiftBranchFromStart(t, probe.Stdout)
	pointer := "ref:refs/bench/recovery/" + branch

	ledgerPath := filepath.Join(gitDir(t, f), "bench-intent.json")
	data, err := os.ReadFile(ledgerPath)
	if err != nil {
		t.Fatalf("read intent ledger: %v", err)
	}
	var ledger struct {
		Entries []struct{ Objective, Outcome, Recovery string }
	}
	if err := json.Unmarshal(data, &ledger); err != nil {
		t.Fatalf("decode intent ledger: %v", err)
	}
	found := false
	for _, e := range ledger.Entries {
		if e.Objective != "surface recovery" {
			continue
		}
		found = true
		if e.Outcome != "failed" || e.Recovery != pointer {
			t.Fatalf("intent entry = %+v, want outcome=failed recovery=%s", e, pointer)
		}
	}
	if !found {
		t.Fatal("intent ledger is missing the shift entry")
	}

	status := f.BenchEnv(map[string]string{"BENCH_HOME": home}, "status")
	status.RequireExit(0)
	status.RequireContains(status.Stdout, pointer)
}

// testShiftInterruptWithMutationRecovers extends the interrupt contracts: an interrupt
// that lands after the adapter has already mutated the tree still snapshots the
// mutation (scratch excluded) before releasing the pool worktree.
func testShiftInterruptWithMutationRecovers(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteExecutable("interrupt-agent", "#!/usr/bin/env bash\nprintf 'work\\n' > work.txt\nkill -INT \"$PPID\"\nexit 130\n")
	f.CommitAll("agent")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": filepath.Join(f.Root, "interrupt-agent"), "BENCH_MAX_ITERS": "1", "BENCH_HOME": home}, "shift", "interrupt with mutation")

	probe.RequireExit(130)
	branch := shiftBranchFromStart(t, probe.Stdout)
	ref := "refs/bench/recovery/" + branch
	f.Git("show-ref", "--verify", ref)
	tree := f.Git("ls-tree", "-r", "--name-only", ref).Stdout
	if !strings.Contains(tree, "work.txt") {
		t.Fatalf("interrupted shift did not snapshot the adapter's mutation:\n%s", tree)
	}
	for _, path := range []string{".bench-notes.md", ".bench-objective"} {
		if strings.Contains(tree, path) {
			t.Fatalf("interrupted shift's snapshot rode scratch %s into the tree:\n%s", path, tree)
		}
	}
	requireNoLease(t, home)
}

// testShiftNoObjective covers row 9: an empty objective exits 2 before acquiring
// anything — the "improve the codebase" default is gone.
func testShiftNoObjective(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": "true", "BENCH_HOME": home}, "shift")

	probe.RequireExit(2)
	probe.RequireContains(probe.Stdout, "usage,\"2\",none,\"0\",\"0\",none,")
	requireNoLease(t, home)
	requireNoShiftBranch(t, f)
}

// testShiftWhitespaceObjective covers row 9's other half: a whitespace-only objective
// exits 2 before acquiring anything, same as an empty one — a bare `== ""` guard would
// let this slip through as a legitimate objective.
func testShiftWhitespaceObjective(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": "true", "BENCH_HOME": home}, "shift", "   ")

	probe.RequireExit(2)
	probe.RequireContains(probe.Stdout, "usage,\"2\",none,\"0\",\"0\",none,")
	requireNoLease(t, home)
	requireNoShiftBranch(t, f)
}

// testShiftControlByteObjective covers row 10: an objective carrying a control byte
// (ESC, 0x1b) exits 2 at entry, before it can reach the ledger or the TOON emitter.
func testShiftControlByteObjective(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": "true", "BENCH_HOME": home}, "shift", "bad\x1bobjective")

	probe.RequireExit(2)
	probe.RequireContains(probe.Stdout, "usage,\"2\",none,\"0\",\"0\",none,")
	requireNoLease(t, home)
	requireNoShiftBranch(t, f)
}

// testShiftInvalidEnvUsage covers row 11: BENCH_MAX_ITERS, BENCH_MAX_WALL set to
// values outside the accepted range each exit 2, naming the variable and range, before
// any lease is acquired.
func testShiftInvalidEnvUsage(t *testing.T) {
	const itersRange = "[1,100]"
	const wallRange = "greater than 0 and at most 24h0m0s"
	cases := []struct {
		name      string
		env       map[string]string
		want      string // substring stderr must contain: the variable name
		wantRange string // substring stderr must also contain: the accepted range
	}{
		{"zero iters", map[string]string{"BENCH_MAX_ITERS": "0"}, "BENCH_MAX_ITERS", itersRange},
		{"non-integer iters", map[string]string{"BENCH_MAX_ITERS": "abc"}, "BENCH_MAX_ITERS", itersRange},
		{"iters over max", map[string]string{"BENCH_MAX_ITERS": "101"}, "BENCH_MAX_ITERS", itersRange},
		{"zero refactor iters", map[string]string{"BENCH_REFACTOR_ITERS": "0"}, "BENCH_REFACTOR_ITERS", itersRange},
		{"wall unparseable", map[string]string{"BENCH_MAX_WALL": "abc"}, "BENCH_MAX_WALL", wallRange},
		{"wall zero", map[string]string{"BENCH_MAX_WALL": "0s"}, "BENCH_MAX_WALL", wallRange},
		{"wall negative", map[string]string{"BENCH_MAX_WALL": "-5m"}, "BENCH_MAX_WALL", wallRange},
		{"wall over bound", map[string]string{"BENCH_MAX_WALL": "48h"}, "BENCH_MAX_WALL", wallRange},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
			home := t.TempDir()
			env := map[string]string{"BENCH_AGENT": "true", "BENCH_HOME": home}
			for k, v := range tc.env {
				env[k] = v
			}

			probe := f.BenchEnv(env, "shift", "objective")

			probe.RequireExit(2)
			probe.RequireContains(probe.Stderr, tc.want)
			probe.RequireContains(probe.Stderr, tc.wantRange)
			requireNoLease(t, home)
			requireNoShiftBranch(t, f)
		})
	}
}

// testShiftEmptyEnvDefault covers row 11's other half: BENCH_MAX_ITERS explicitly set
// to the empty string behaves exactly like unset — the default cap applies and no
// usage error fires.
func testShiftEmptyEnvDefault(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	home := t.TempDir()

	// No adapter mutation ever lands, so the honest end state is no-op/4 either way —
	// the point of this row is that "" behaves as unset (no exit 2), not that the
	// default numeric cap is visible in the output.
	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": "true", "BENCH_MAX_ITERS": "", "BENCH_HOME": home}, "shift", "objective")

	probe.RequireExit(4)
}

// testShiftMaxWallInclusiveBound covers row 11's accept case: BENCH_MAX_WALL=24h sits
// exactly on the accepted upper bound, so it must run normally rather than exit 2 — an
// off-by-one on the inclusive bound would reject this legal value.
func testShiftMaxWallInclusiveBound(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	home := t.TempDir()

	// No adapter mutation ever lands, so the honest end state is no-op/4 — the point of
	// this row is that 24h is accepted (no exit 2), not the numeric cap itself.
	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": "true", "BENCH_MAX_WALL": "24h", "BENCH_HOME": home}, "shift", "objective")

	probe.RequireExit(4)
}

// testShiftCompleteResult covers rows 6 and 16: a green gate, a committing adapter, and
// a .bench/done.sh that passes once the work has landed stop the loop clean — exit 0,
// and the shift_result block carries all seven fields with outcome complete, committed
// 1, recovery none.
func testShiftCompleteResult(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteExecutable(".bench/done.sh", "#!/usr/bin/env bash\n[ -f landed.txt ]\n")
	f.WriteExecutable("agent", "#!/usr/bin/env bash\nprintf 'work\\n' > landed.txt\n")
	f.CommitAll("agent")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "3", "BENCH_HOME": home}, "shift", "land it")

	probe.RequireExit(0)
	branch := shiftBranch(t, probe.Stdout)
	requireShiftResult(t, probe.Stdout, "complete,\"0\","+branch+",\"1\",\"1\",none,")
}

// testShiftAdapterFailureZeroCommit covers rows 6 and 7: BENCH_AGENT=/bin/false with a
// green gate stops the loop as failed/1 with zero commits — today this exits 0
// "objective likely met"; the whole point of this row is that it no longer does.
func testShiftAdapterFailureZeroCommit(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": "/bin/false", "BENCH_MAX_ITERS": "2", "BENCH_HOME": home}, "shift", "objective")

	probe.RequireExit(1)
	branch := shiftBranchFromStart(t, probe.Stdout)
	requireShiftResult(t, probe.Stdout, "failed,\"1\","+branch+",\"0\",")
}

// testShiftAdapterSpawnFailure covers row 7's other half: an adapter that passes the
// preflight (it is a regular executable file) but fails to exec at iteration time — here
// because its shebang names an interpreter that does not exist — stops the loop as
// failed/1, not no-op/4. runAdapter used to swallow cmd.Start's error, so a spawn
// failure read as a clean "nothing happened" no-op; this proves it now propagates as
// real evidence of a broken adapter.
func testShiftAdapterSpawnFailure(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteExecutable("agent", "#!/no/such/interpreter\necho hi\n")
	f.CommitAll("agent")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "1", "BENCH_HOME": home}, "shift", "objective")

	probe.RequireExit(1)
	branch := shiftBranchFromStart(t, probe.Stdout)
	requireShiftResult(t, probe.Stdout, "failed,\"1\","+branch+",\"0\",")
	requireNoLease(t, home)
}

// testShiftAdapterFailureAfterCommit covers row 7: an adapter that writes a file then
// exits nonzero on every run commits its first iteration's mutation (gate green) and
// then stops the loop as incomplete/3 — folding this into failed/1 would understate the
// committed work.
func testShiftAdapterFailureAfterCommit(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteExecutable("agent", `#!/usr/bin/env bash
n=0; [ -f count ] && n="$(cat count)"
n=$((n+1)); printf '%s\n' "$n" > count
printf 'work %s\n' "$n" > "work-$n.txt"
exit 1
`)
	f.CommitAll("agent")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "3", "BENCH_HOME": home}, "shift", "objective")

	probe.RequireExit(3)
	branch := shiftBranchFromStart(t, probe.Stdout)
	requireShiftResult(t, probe.Stdout, "incomplete,\"3\","+branch+",\"1\",")
	requireEqual(t, strings.TrimSpace(f.Git("rev-list", "--count", "HEAD.."+branch).Stdout), "1", "adapter-failure-after-commit shift committed the wrong number of iterations")
}

// testShiftCapExhaustion covers row 12: an always-dirty adapter runs to the iteration
// cap, leaving the committed work on the branch and exiting incomplete/3 — today cap
// exhaustion exits 0, reading as done.
func testShiftCapExhaustion(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteExecutable("agent", `#!/usr/bin/env bash
n=0; [ -f count ] && n="$(cat count)"
n=$((n+1)); printf '%s\n' "$n" > count
printf 'work %s\n' "$n" > "work-$n.txt"
`)
	f.CommitAll("agent")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "2", "BENCH_HOME": home}, "shift", "objective")

	probe.RequireExit(3)
	branch := shiftBranch(t, probe.Stdout)
	requireShiftResult(t, probe.Stdout, "incomplete,\"3\","+branch+",\"2\",\"2\",none,")
	requireEqual(t, strings.TrimSpace(f.Git("rev-list", "--count", "HEAD.."+branch).Stdout), "2", "cap-exhaustion shift committed the wrong number of iterations")
}

// testShiftNoOpDone covers row 13: a no-op adapter with a .bench/done.sh that always
// passes commits zero iterations and exits no-op/4, with a detail naming the predicate
// — complete-on-predicate would exit 0 and hide that nothing was done.
func testShiftNoOpDone(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	f.WriteExecutable(".bench/done.sh", "#!/usr/bin/env bash\nexit 0\n")
	f.CommitAll("done.sh")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": "true", "BENCH_MAX_ITERS": "2", "BENCH_HOME": home}, "shift", "objective")

	probe.RequireExit(4)
	probe.RequireContains(probe.Stdout, "predicate")
}

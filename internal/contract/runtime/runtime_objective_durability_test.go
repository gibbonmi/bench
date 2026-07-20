package runtime

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

// Seam-C sentinel contracts: the durable state a `bench shift <objective>` leaves behind,
// observed at the built binary. Each row proves one durability claim from the objective
// handling decision — the intake cap, the mode-0600 scratch file, the ledger carrying a
// key rather than free text, the status read-back, and the reject-before-any-write
// ordering — so a regression that widened the exposure again would turn the gate red.
func TestRuntimeObjectiveDurabilityContracts(t *testing.T) {
	t.Parallel()
	contract.RunParallel(t, "over-200-rune objective is rejected at intake", testObjectiveOverCapRejected)
	contract.RunParallel(t, "objective rune-cap boundary is exact", testObjectiveCapBoundary)
	contract.RunParallel(t, "live worktree .bench-objective is mode 0600", testObjectiveScratchMode0600)
	contract.RunParallel(t, "ledger carries no objective text", testLedgerHasNoObjectiveText)
	contract.RunParallel(t, "status shows a live shift's objective from its worktree", testStatusShowsLiveShiftObjective)
	contract.RunParallel(t, "status degrades to the key when the scratch file is gone", testStatusDegradesWhenScratchAbsent)
	contract.RunParallel(t, "control-byte objective is rejected before any write", testControlByteRejectedBeforeAnyWrite)
}

// testObjectiveOverCapRejected covers row 10: an objective longer than 200 runes exits 2
// with a usage error naming the actual and maximum length, before a shift begins. Red
// before the cap existed, when an over-long objective was accepted and a shift started.
func testObjectiveOverCapRejected(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	home := t.TempDir()

	long := strings.Repeat("a", 201)
	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": "true", "BENCH_HOME": home}, "shift", long)
	probe.RequireExit(2)
	out := probe.Stdout + probe.Stderr
	for _, want := range []string{"201", "200", "runes"} {
		if !strings.Contains(out, want) {
			t.Fatalf("over-cap rejection did not name %q:\n%s", want, out)
		}
	}
	requireNoShiftBranch(t, f)
	requireNoLease(t, home)
}

// testObjectiveCapBoundary covers row 10's edge: 200 runes is accepted and 201 is
// rejected. The runes are multibyte (é, two bytes each), so a byte-counting implementation
// would reject the 200-rune objective at 400 bytes and turn the accepted case red.
func testObjectiveCapBoundary(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")

	atCap := f.BenchEnv(map[string]string{"BENCH_AGENT": "true", "BENCH_MAX_ITERS": "1", "BENCH_HOME": t.TempDir()}, "shift", strings.Repeat("é", 200))
	if atCap.ExitCode == 2 {
		t.Fatalf("a 200-rune objective was rejected as over-cap — the cap is counting bytes, not runes:\n%s%s", atCap.Stdout, atCap.Stderr)
	}
	overCap := f.BenchEnv(map[string]string{"BENCH_AGENT": "true", "BENCH_HOME": t.TempDir()}, "shift", strings.Repeat("é", 201))
	overCap.RequireExit(2)
}

// testObjectiveScratchMode0600 covers row 11: the .bench-objective file in a live shift
// worktree has mode 0600. The adapter runs with the worktree as its cwd, so it stats the
// real file while the worktree is alive and records the octal mode; the test asserts on
// that recording. Red today, when the file is written 0644.
func testObjectiveScratchMode0600(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	state := t.TempDir()
	f.WriteExecutable("agent", "#!/usr/bin/env bash\n"+
		"if stat -c '%a' .bench-objective >/dev/null 2>&1; then\n"+
		"  stat -c '%a' .bench-objective > \"$BENCH_TEST_STATE/objmode\"\n"+
		"else\n"+
		"  stat -f '%Lp' .bench-objective > \"$BENCH_TEST_STATE/objmode\"\n"+
		"fi\n")
	f.CommitAll("agent")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{
		"BENCH_AGENT": filepath.Join(f.Root, "agent"), "BENCH_MAX_ITERS": "1",
		"BENCH_HOME": home, "BENCH_TEST_STATE": state,
	}, "shift", "mode check objective")
	if probe.ExitCode == 2 {
		t.Fatalf("shift refused a valid objective:\n%s%s", probe.Stdout, probe.Stderr)
	}
	mode := strings.TrimSpace(string(mustReadRuntime(t, filepath.Join(state, "objmode"))))
	if mode != "600" {
		t.Fatalf(".bench-objective mode = %q, want 600", mode)
	}
}

// testLedgerHasNoObjectiveText covers row 12: after a shift, the intent ledger JSON
// carries no field named objective and the marker objective's text appears nowhere in the
// serialized file. Searching the whole file, not one named field, survives the field being
// renamed rather than removed. Red today, when Entry.Objective stored it verbatim.
func testLedgerHasNoObjectiveText(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	home := t.TempDir()
	const marker = "FT88durabilitymarkerobjectivexyz"

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": "true", "BENCH_MAX_ITERS": "1", "BENCH_HOME": home}, "shift", marker)
	if probe.ExitCode == 2 {
		t.Fatalf("shift refused a valid objective:\n%s%s", probe.Stdout, probe.Stderr)
	}
	ledger := string(mustReadRuntime(t, filepath.Join(gitDir(t, f), "bench-intent.json")))
	if strings.Contains(ledger, marker) {
		t.Fatalf("intent ledger carried the objective text verbatim:\n%s", ledger)
	}
	if strings.Contains(ledger, `"objective"`) {
		t.Fatalf("intent ledger still has an objective field:\n%s", ledger)
	}
}

// testStatusShowsLiveShiftObjective covers row 13: bench status shows a live interrupted
// shift's objective even though the ledger no longer stores it, by reading it back from the
// recorded worktree's .bench-objective. Red against an implementation that drops the field
// without reading the worktree file. The state is constructed directly so the "live shift
// with a surviving worktree" case is deterministic rather than timing-dependent.
func testStatusShowsLiveShiftObjective(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("tracked.txt", "base\n")
	f.CommitAll("base")

	wt := filepath.Join(t.TempDir(), "live-worktree")
	if err := os.MkdirAll(wt, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(wt, ".bench-objective"), []byte("resume this exact objective\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	ledger := `{"schema":2,"entries":[{"key":"shift-live-1","kind":"shift","created_at":"2026-07-11T00:00:00Z","worktree":"` + wt + `"}]}` + "\n"
	contract.WriteFileAbs(t, filepath.Join(gitDir(t, f), "bench-intent.json"), ledger)

	out := f.Bench("status", "--all")
	out.RequireExit(0)
	out.RequireContains(out.Stdout, "resume this exact objective")
}

// testStatusDegradesWhenScratchAbsent covers row 13's edge: when the worktree survives but
// its .bench-objective is gone — a normal end state — status renders the entry by its key
// without error rather than propagating the read failure. Red against an implementation
// that treats the missing scratch file as fatal.
func testStatusDegradesWhenScratchAbsent(t *testing.T) {
	f := contract.NewFixture(t)
	f.WriteFile("tracked.txt", "base\n")
	f.CommitAll("base")

	wt := filepath.Join(t.TempDir(), "worktree-without-scratch")
	if err := os.MkdirAll(wt, 0o755); err != nil {
		t.Fatal(err)
	}
	ledger := `{"schema":2,"entries":[{"key":"shift-degraded-1","kind":"shift","created_at":"2026-07-11T00:00:00Z","worktree":"` + wt + `"}]}` + "\n"
	contract.WriteFileAbs(t, filepath.Join(gitDir(t, f), "bench-intent.json"), ledger)

	out := f.Bench("status", "--all")
	out.RequireExit(0)
	out.RequireContains(out.Stdout, "objective=shift-degraded-1")
}

// testControlByteRejectedBeforeAnyWrite covers row 15's ordering claim: a control-byte
// objective is rejected with no intent entry written and no worktree created — validation
// runs before any durable write. Asserting the rejection alone would pass an implementation
// that wrote first and cleaned up after, so the absence of both artifacts is the assertion.
func testControlByteRejectedBeforeAnyWrite(t *testing.T) {
	f := shiftFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	home := t.TempDir()

	probe := f.BenchEnv(map[string]string{"BENCH_AGENT": "true", "BENCH_HOME": home}, "shift", "bad"+string(rune(0x1b))+"objective")
	probe.RequireExit(2)

	if _, err := os.Stat(filepath.Join(gitDir(t, f), "bench-intent.json")); err == nil {
		t.Fatal("control-byte objective wrote an intent ledger before rejection")
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat intent ledger: %v", err)
	}
	requireNoShiftBranch(t, f)
	requireNoLease(t, home)
}

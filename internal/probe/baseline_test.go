// The baseline family: every row that grades the run the probe takes over the unmutated
// tree before it writes anything. Rows DG7 through DG12, DG35, and DG36 live here, and
// PB16 rides DG36's second case, because one signal test covers both starts.

package probe

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/probe/probetest"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// interruptHelperEnv carries the fixture root, and interruptArgvEnv the JSON argv, to the
// helper half of the interrupt rows. The fact under test is what a signal does to the verb. So
// the verb runs in a process that the signal may end without the assertions.
const (
	interruptHelperEnv = "BENCH_PROBE_INTERRUPT_ROOT"
	interruptArgvEnv   = "BENCH_PROBE_INTERRUPT_ARGV"
)

// buildFailedRun is a package the run reported as failed with no failing test, which is
// what a package that does not compile prints.
const buildFailedRun = `{"Action":"fail","Package":"probefixture","Elapsed":0.01}`

// startCount answers how many `go test` children the stub recorded, which is how a row
// proves the baseline ran from outside the verb's own process.
func startCount(t *testing.T, f *fixture) int {
	t.Helper()
	body, err := os.ReadFile(f.marker)
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	return strings.Count(string(body), "test\n")
}

// requireBaselineOnly proves the baseline was the whole run: it started once, and whatever
// refused after it started no mutated run.
func requireBaselineOnly(t *testing.T, f *fixture) {
	t.Helper()
	if got := startCount(t, f); got != 1 {
		t.Fatalf("Go test starts = %d, want only the baseline's one", got)
	}
}

// DG7: the probe starts one Go test before the mutation and one after it. The stub copies
// the subject on each start, so the order is proven from outside the verb's own process.
func TestProbeRunsTheBaselineFirst(t *testing.T) {
	f := newFixture(t)
	copies := scratchDir(t)
	record := "cp " + sanitize.ShellQuote(f.subject) + " " + sanitize.ShellQuote(copies) + "/$n\n"
	installStubStarts(t, f,
		probetest.Start{Hook: record, Events: cannedPass},
		probetest.Start{Hook: record, Events: cannedFailure("caught"), Exit: 1})
	if out, code := runProbe(t, probeArgs()...); code != 0 {
		t.Fatalf("exit = %d, want 0\n%s", code, out)
	}
	if got := startCount(t, f); got != 2 {
		t.Fatalf("Go test starts = %d, want 2", got)
	}
	first, err := os.ReadFile(filepath.Join(copies, "1"))
	if err != nil || string(first) != clampSource {
		t.Fatalf("the first start's subject = (%q, %v), want the unmutated bytes", first, err)
	}
	second, err := os.ReadFile(filepath.Join(copies, "2"))
	if err != nil || strings.Contains(string(second), "n < 0") {
		t.Fatalf("the second start's subject = (%q, %v), want the mutated bytes", second, err)
	}
}

// requireBaselineRefusal drives one probe whose baseline did not pass and answers what the
// verb printed. Every non-passing kind owes the same four facts: the invalid row with the
// untouched cell, an unchanged subject, a home the probe never wrote to, and one Go start.
func requireBaselineRefusal(t *testing.T, f *fixture, cause string, failed int, first any, args ...string) string {
	t.Helper()
	out, code := runProbe(t, args...)
	requireRow(t, out, code, 1, probeRow(t, "invalid", "clamp.go", "swap", cause, failed, first, "untouched"))
	requireSubjectBytes(t, f, clampSource)
	requireHomeEmpty(t, f)
	requireBaselineOnly(t, f)
	return out
}

// DG8, DG9, DG10: a red baseline is invalid rather than a bite, it names the tests that
// were red before the probe ran, and it writes nothing.
func TestProbeRefusesARedBaseline(t *testing.T) {
	f := newFixture(t)
	installStubStarts(t, f, probetest.Start{Events: cannedFailure("already red"), Exit: 1}, probetest.Start{Events: cannedPass})
	out := requireBaselineRefusal(t, f, "baseline-failed", 1, "TestClampNegative", probeArgs()...)
	if want := selectionRow(t, "package", "./", "^TestClampNegative$", "failed", 0); !strings.Contains(out, want) {
		t.Fatalf("stdout = %q, want the selection row %q", out, want)
	}
	if !strings.Contains(out, "already red") {
		t.Fatalf("stdout = %q, want the baseline's own failures table", out)
	}
}

// DG11: a baseline that reached no verdict names its own kind, so a build failure is never
// reported as a red suite. The build-failed stream passes no run pattern, because `bench
// test` attributes a build failure under a pattern to the pattern.
func TestProbeReportsABaselineWithoutAVerdict(t *testing.T) {
	for _, tc := range []struct {
		cause, events string
		exit          int
		args          []string
	}{
		{"baseline-build-failed", buildFailedRun, 1, probeArgs("clamp.go", "--swap", "n < 0", "--with", "n > 0", "--package", "./")},
		{"baseline-no-test-run", quietRun, 0, probeArgs()},
	} {
		t.Run(tc.cause, func(t *testing.T) {
			f := newFixture(t)
			installStubStarts(t, f, probetest.Start{Events: tc.events, Exit: tc.exit}, probetest.Start{Events: cannedPass})
			requireBaselineRefusal(t, f, tc.cause, 0, nil, tc.args...)
		})
	}
}

// DG35: a baseline that prints no event at all reaches the no-packages refusal, so the
// cause names that refusal instead of collapsing onto a red suite.
func TestProbeReportsABaselineRefusal(t *testing.T) {
	f := newFixture(t)
	installStubStarts(t, f, probetest.Start{Exit: 1}, probetest.Start{Events: cannedPass})
	requireBaselineRefusal(t, f, "baseline-refused", 0, nil, probeArgs()...)
}

// PB9, DG12: --full reaches both runs, so a long diagnostic from either one is previewed
// without the flag and complete with it.
func TestProbeForwardsFullToTheFocusedRun(t *testing.T) {
	f := newFixture(t)
	long := strings.Repeat("d", bounds.PreviewRuneLimit+40)
	red := probetest.Start{Events: cannedFailure(long), Exit: 1}
	for _, tc := range []struct {
		name              string
		baseline, mutated probetest.Start
	}{
		{"mutated run", probetest.Start{Events: cannedPass}, red},
		{"baseline", red, red},
	} {
		t.Run(tc.name, func(t *testing.T) {
			installStubStarts(t, f, tc.baseline, tc.mutated)
			if previewed, _ := runProbe(t, probeArgs()...); strings.Contains(previewed, long) {
				t.Fatalf("without --full the diagnostic was not previewed:\n%s", previewed)
			}
			installStubStarts(t, f, tc.baseline, tc.mutated)
			if full, _ := runProbe(t, append(probeArgs(), "--full")...); !strings.Contains(full, long) {
				t.Fatalf("--full did not reach the run:\n%s", full)
			}
		})
	}
}

// PB16, DG36: a signal at the baseline ends the probe before any write, and a signal at the
// mutated run still puts the subject back before the verb answers, because the restore is
// deferred around that run. Either way the tree and the home come out as they went in.
func TestProbeRestoresOnInterrupt(t *testing.T) {
	if root := os.Getenv(interruptHelperEnv); root != "" {
		runInterruptHelper(t, root)
		return
	}
	sleeping := probetest.Start{Hook: "sleep 60\n", Events: cannedPass}
	for _, tc := range []struct {
		name            string
		start           int
		baseline        probetest.Start
		cause, restored string
	}{
		{"baseline", 1, sleeping, "baseline-interrupted", "untouched"},
		{"mutated run", 2, probetest.Start{Events: cannedPass}, "interrupted", "yes"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			installStubStarts(t, f, tc.baseline, sleeping)
			out := interruptedProbe(t, f, probeArgs(), func() { awaitStubStart(t, f, tc.start) })
			want := probeRow(t, "invalid", "clamp.go", "swap", tc.cause, 0, nil, tc.restored)
			if !strings.Contains(out, want) {
				t.Fatalf("helper stdout = %q, want the interrupted row %q", out, want)
			}
			requireSubjectBytes(t, f, clampSource)
			requireHomeEmpty(t, f)
		})
	}
}

// A signal during mutant 1's run stops the set after that mutant's restore.
func TestMutantSetStopsOnInterrupt(t *testing.T) {
	f := newFixture(t)
	installSet(t, f, probetest.Start{Hook: "sleep 60\n", Events: cannedPass}, bite)
	out := interruptedProbe(t, f, setArgs(mutantsFile(t, guardSwap, returnSwap)), func() { awaitStubStart(t, f, 2) })
	if want := setRows(t, interruptedRow, untouchedRow(2, "set-stopped")); !strings.HasPrefix(out, want) {
		t.Fatalf("helper stdout = %q, want %q first", out, want)
	}
	requireStarts(t, f, 2)
	requireSubjectBytes(t, f, clampSource)
	requireHomeEmpty(t, f)
}

// A signal during mutant 1's fresh build interrupts that run, and mutant 2 is never written.
// The subject sits inside the run binary's closure, so the held binary no longer validates.
func TestMutantSetStopsOnAnInterruptedBuild(t *testing.T) {
	f := newFixture(t)
	main := filepath.Join(f.root, "cmd", "bench", "main.go")
	start, err := os.ReadFile(main)
	if err != nil {
		t.Fatal(err)
	}
	parked := parkSecondBuild(t, f)
	installSet(t, f, bite, bite)
	path := mutantsFile(t, []string{"--swap", "func main() {}", "--with", "func main() { _ = 1 }"}, []string{"--swap", "package main", "--with", "package main // set"})
	out := interruptedProbe(t, f, slices.Concat([]string{"cmd/bench/main.go", "--mutants", path}, clampSelection), func() { awaitFile(t, parked) })
	if want := setRows(t, interruptedRow, untouchedRow(2, "set-stopped")); !strings.HasPrefix(out, want) {
		t.Fatalf("helper stdout = %q, want %q first", out, want)
	}
	requireStarts(t, f, 1)
	if now, err := os.ReadFile(main); err != nil || string(now) != string(start) {
		t.Fatalf("cmd/bench/main.go = (%q, %v), want the start bytes", now, err)
	}
}

// interruptedRow is the row of entry 1 whose run a signal ended after its write.
var interruptedRow = []any{1, "invalid", "swap", "interrupted", 0, nil, 0, 1, 1, nil, "yes"}

// interruptedProbe runs the verb over argv in a helper process, signals it once await returns,
// and answers what the verb printed.
func interruptedProbe(t *testing.T, f *fixture, argv []string, await func()) string {
	t.Helper()
	encoded, err := json.Marshal(argv)
	if err != nil {
		t.Fatal(err)
	}
	helper := exec.Command(os.Args[0], "-test.run=^TestProbeRestoresOnInterrupt$", "-test.timeout="+interruptDeadline().String())
	helper.Env = append(os.Environ(), interruptHelperEnv+"="+f.root, interruptArgvEnv+"="+string(encoded))
	stdout, err := helper.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	helper.Stderr = helper.Stdout
	if err := helper.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = helper.Process.Kill() })
	answered := make(chan string, 1)
	go func() {
		body, _ := io.ReadAll(stdout)
		_ = helper.Wait()
		answered <- string(body)
	}()
	await()
	if err := helper.Process.Signal(syscall.SIGINT); err != nil {
		t.Fatal(err)
	}
	select {
	case out := <-answered:
		return out
	case <-time.After(interruptDeadline()):
		t.Fatalf("the interrupted probe did not answer within %s", interruptDeadline())
	}
	return ""
}

// runInterruptHelper is the signalled half: it runs the verb for real and prints what the
// verb answered, so the parent grades production's own deferred restore.
func runInterruptHelper(t *testing.T, root string) {
	t.Helper()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	var argv []string
	if err := json.Unmarshal([]byte(os.Getenv(interruptArgvEnv)), &argv); err != nil {
		t.Fatal(err)
	}
	out, code := Command(argv)
	fmt.Print(out)
	fmt.Printf("exit=%d\n", code)
}

// awaitStubStart waits until the stub recorded the numbered `go test` child, not the `go
// list` the run binary's verification takes on the way in. The number separates the
// baseline's start from the mutated run's.
func awaitStubStart(t *testing.T, f *fixture, want int) {
	t.Helper()
	expiry := time.Now().Add(interruptDeadline())
	for time.Now().Before(expiry) {
		if startCount(t, f) >= want {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("Go test start %d did not happen within %s", want, interruptDeadline())
}

func interruptDeadline() time.Duration { return bounds.TestDeadline(runbinary.BuilderCancelGrace) }

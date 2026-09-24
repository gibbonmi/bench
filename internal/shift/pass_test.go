package shift

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/otelrecord"
)

// gateSeam is the seam of the gate span that each pass parents.
const gateSeam = "gate"

// childrenOf returns the spans of seam whose parent is parent, in start order.
func childrenOf(spans []otelrecord.Span, parent otelrecord.Span, seam string) []otelrecord.Span {
	var found []otelrecord.Span
	for _, span := range spans {
		if span.Seam == seam && span.ParentSpanID == parent.SpanID {
			found = append(found, span)
		}
	}
	sort.Slice(found, func(i, j int) bool { return found[i].Start.Before(found[j].Start) })
	return found
}

// withAgent writes an adapter script outside the repository and points BENCH_AGENT at it.
func withAgent(t *testing.T, script string) string {
	t.Helper()
	agent := filepath.Join(t.TempDir(), "agent")
	if err := os.WriteFile(agent, []byte("#!/usr/bin/env bash\n"+script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_AGENT", agent)
	return agent
}

// twoPassShift runs a two-iteration shift whose adapter changes work.txt on each run, so
// each iteration commits and the shift stops at its cap.
func twoPassShift(t *testing.T) (shift otelrecord.Span, passes, spans []otelrecord.Span, root, stdout string) {
	t.Helper()
	root = faultFixtureCore(t, "#!/usr/bin/env bash\nexit 0\n", nil)
	withAgent(t, "echo \"$RANDOM $$\" >> work.txt\n")
	t.Setenv("BENCH_MAX_ITERS", "2")
	shift, _, stdout = runRecordedShift(t, 3)
	_, spans, _ = recordedShift(t)
	return shift, childrenOf(spans, shift, iterationSeam), spans, root, stdout
}

// LE38: a two-iteration shift writes two iteration spans under the shift span.
func TestATwoIterationShiftWritesTwoPassSpans(t *testing.T) {
	_, passes, _, _, _ := twoPassShift(t)
	if len(passes) != 2 {
		t.Fatalf("the shift span has %d iteration children, want 2", len(passes))
	}
}

// LE42: each pass parents exactly one gate span with a subject and an outcome.
func TestEachPassParentsOneGateSpan(t *testing.T) {
	_, passes, spans, _, _ := twoPassShift(t)
	if len(passes) != 2 {
		t.Fatalf("the shift span has %d iteration children, want 2", len(passes))
	}
	for i, pass := range passes {
		gates := childrenOf(spans, pass, gateSeam)
		if len(gates) != 1 {
			t.Fatalf("pass %d has %d gate children, want 1", i+1, len(gates))
		}
		if gates[0].Attributes[otelrecord.AttrSubjectID] == "" || gates[0].Attributes[otelrecord.AttrOutcome] == "" {
			t.Fatalf("pass %d gate span lacks a subject or an outcome: %v", i+1, gates[0].Attributes)
		}
	}
}

// LE43: each pass that committed names the commit the branch gained in that pass.
func TestACommittedPassCarriesItsCommit(t *testing.T) {
	_, passes, _, root, stdout := twoPassShift(t)
	commits := strings.Fields(runGitOutput(t, root, "rev-list", "--reverse", "main.."+shiftBranchName(t, stdout)))
	if len(passes) != 2 || len(commits) != 2 {
		t.Fatalf("got %d passes and %d commits, want 2 and 2", len(passes), len(commits))
	}
	for i, pass := range passes {
		requireAttr(t, pass, otelrecord.AttrSubjectID, commits[i])
	}
}

// LE39: a pass whose adapter exits 3 records the exit.
func TestAPassRecordsTheAdapterExit(t *testing.T) {
	faultFixtureCore(t, "#!/usr/bin/env bash\nexit 0\n", nil)
	withAgent(t, "exit 3\n")
	shift, _, _ := runRecordedShift(t, 1)
	_, spans, _ := recordedShift(t)
	passes := childrenOf(spans, shift, iterationSeam)
	if len(passes) != 1 {
		t.Fatalf("the shift span has %d iteration children, want 1", len(passes))
	}
	requireAttr(t, passes[0], otelrecord.AttrAdapterResult, otelrecord.AdapterExited)
	requireAttr(t, passes[0], otelrecord.AttrAdapterExit, "3")
}

// LE40: a pass whose adapter cannot start records a spawn failure and no exit.
func TestAPassRecordsASpawnFailure(t *testing.T) {
	faultFixtureCore(t, "#!/usr/bin/env bash\nexit 0\n", nil)
	withAgent(t, "echo \"$RANDOM $$\" >> work.txt\nchmod -x \"$0\"\n")
	t.Setenv("BENCH_MAX_ITERS", "2")
	shift, _, _ := runRecordedShift(t, 3)
	_, spans, _ := recordedShift(t)
	passes := childrenOf(spans, shift, iterationSeam)
	if len(passes) != 2 {
		t.Fatalf("the shift span has %d iteration children, want 2", len(passes))
	}
	requireAttr(t, passes[1], otelrecord.AttrAdapterResult, otelrecord.AdapterSpawnFailed)
	requireNoAttr(t, passes[1], otelrecord.AttrAdapterExit)
}

// LE41: a refactor pass writes one refactor span under the shift span.
func TestARefactorPassWritesARefactorSpan(t *testing.T) {
	faultFixtureCore(t, "#!/usr/bin/env bash\nexit 0\n", nil)
	withAgent(t, "printf 'package work\\n\\n// a\\n// b\\n// c\\n// d\\n' > work.go\n")
	t.Setenv("BENCH_MAX_LINES", "3")
	shift, _, _ := runRecordedShift(t, 3)
	_, spans, _ := recordedShift(t)
	if refactors := childrenOf(spans, shift, refactorSeam); len(refactors) != 1 {
		t.Fatalf("the shift span has %d refactor children, want 1", len(refactors))
	}
}

// LE45: the adapter runs with the handoff of its pass span.
func TestTheAdapterReceivesThePassHandoff(t *testing.T) {
	faultFixtureCore(t, "#!/usr/bin/env bash\nexit 0\n", nil)
	handoff := filepath.Join(t.TempDir(), "handoff")
	withAgent(t, "printf '%s\\n%s\\n' \"$BENCH_OTEL_ROOT\" \"$BENCH_OTEL_TRACEPARENT\" > '"+handoff+"'\n")
	shift, _, _ := runRecordedShift(t, 4)
	_, spans, _ := recordedShift(t)
	passes := childrenOf(spans, shift, iterationSeam)
	raw, err := os.ReadFile(handoff)
	if err != nil || len(passes) != 1 {
		t.Fatalf("handoff %q (%v), %d passes, want a handoff and 1 pass", raw, err, len(passes))
	}
	root, traceparent, _ := strings.Cut(strings.TrimSpace(string(raw)), "\n")
	if repo, _ := git.Root(); root != repo {
		t.Fatalf("handoff root = %q, want the repository %q", root, repo)
	}
	if want := "00-" + passes[0].TraceID + "-" + passes[0].SpanID + "-01"; traceparent != want {
		t.Fatalf("handoff traceparent = %q, want the pass span %q", traceparent, want)
	}
}

// helperRoleEnv selects the role of a re-executed test binary.
const helperRoleEnv = "BENCH_SHIFT_HELPER_ROLE"

// TestShiftHelperProcess is the re-exec target of the interrupt tests. In the interrupt
// role it runs one shift in its working directory and exits with the shift's code, so the
// checkpoint's os.Exit ends a real process. Without the role it does nothing.
func TestShiftHelperProcess(t *testing.T) {
	if os.Getenv(helperRoleEnv) != "interrupt" {
		return
	}
	os.Exit(Loop("interrupted shift", io.Discard, io.Discard))
}

// interruptedShift runs a shift in a helper process and sends it SIGINT while its adapter
// runs. It returns the shift span, every span, and the raw record.
func interruptedShift(t *testing.T) (otelrecord.Span, []otelrecord.Span, []byte) {
	t.Helper()
	root := faultFixtureCore(t, "#!/usr/bin/env bash\nexit 0\n", nil)
	started := filepath.Join(t.TempDir(), "started")
	withAgent(t, "touch '"+started+"'\nexec sleep 30\n")
	cmd := exec.Command(os.Args[0], "-test.run=^TestShiftHelperProcess$")
	cmd.Dir = root
	cmd.Env = append(os.Environ(), helperRoleEnv+"=interrupt")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) })
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	for deadline := time.Now().Add(30 * time.Second); ; time.Sleep(20 * time.Millisecond) {
		if _, err := os.Stat(started); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("the adapter never started")
		}
	}
	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		t.Fatal(err)
	}
	select {
	case <-done:
	case <-time.After(30 * time.Second):
		t.Fatal("the interrupted shift never exited")
	}
	if code := cmd.ProcessState.ExitCode(); code != exitCodes[OutcomeInterrupted] {
		t.Fatalf("the interrupted shift exited %d, want %d", code, exitCodes[OutcomeInterrupted])
	}
	return recordedShift(t)
}

// LE37: an interrupted shift still ends its span, with work state interrupted.
func TestAnInterruptedShiftRecordsInterruptedWork(t *testing.T) {
	shift, _, _ := interruptedShift(t)
	requireAttr(t, shift, otelrecord.AttrWorkState, otelrecord.WorkInterrupted)
}

// LE44: the interrupted shift ends its open pass before it ends the shift span.
func TestAnInterruptedShiftEndsThePassFirst(t *testing.T) {
	_, _, raw := interruptedShift(t)
	passEnd, shiftEnd := -1, -1
	for i, line := range bytes.Split(bytes.TrimSpace(raw), []byte("\n")) {
		if bytes.Contains(line, []byte(startMarker)) {
			continue
		}
		switch {
		case bytes.Contains(line, []byte(seamAttr(iterationSeam))):
			passEnd = i
		case bytes.Contains(line, []byte(shiftSeamAttr)):
			shiftEnd = i
		}
	}
	if passEnd < 0 || shiftEnd < 0 || passEnd > shiftEnd {
		t.Fatalf("pass end line %d, shift end line %d, want the pass first:\n%s", passEnd, shiftEnd, raw)
	}
}

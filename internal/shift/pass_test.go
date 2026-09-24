package shift

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/otelrecord"
)

// gateSeam is the seam of the gate span that each pass parents, and greenGate is a gate
// script that always passes.
const (
	gateSeam  = "gate"
	greenGate = "#!/usr/bin/env bash\nexit 0\n"
)

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
	root = faultFixtureCore(t, greenGate, nil)
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
	faultFixtureCore(t, greenGate, nil)
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
	faultFixtureCore(t, greenGate, nil)
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
	faultFixtureCore(t, greenGate, nil)
	withAgent(t, "printf 'package work\\n\\n// a\\n// b\\n// c\\n// d\\n' > work.go\n")
	t.Setenv("BENCH_MAX_LINES", "3")
	shift, _, _ := runRecordedShift(t, 3)
	_, spans, _ := recordedShift(t)
	refactors := childrenOf(spans, shift, refactorSeam)
	if len(refactors) != 1 {
		t.Fatalf("the shift span has %d refactor children, want 1", len(refactors))
	}
	if gates := childrenOf(spans, refactors[0], gateSeam); len(gates) != 1 {
		t.Fatalf("the refactor pass has %d gate children, want 1", len(gates))
	}
	requireAttr(t, refactors[0], otelrecord.AttrAdapterResult, otelrecord.AdapterExited)
}

// LE45: the adapter runs with the handoff of its pass span.
func TestTheAdapterReceivesThePassHandoff(t *testing.T) {
	faultFixtureCore(t, greenGate, nil)
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
	root := faultFixtureCore(t, greenGate, nil)
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
	window := bounds.TestDeadline(0)
	for deadline := time.Now().Add(window); ; time.Sleep(20 * time.Millisecond) {
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
	case <-time.After(window):
		t.Fatal("the interrupted shift never exited")
	}
	if code := cmd.ProcessState.ExitCode(); code != exitCodes[OutcomeInterrupted] {
		t.Fatalf("the interrupted shift exited %d, want %d", code, exitCodes[OutcomeInterrupted])
	}
	return recordedShift(t)
}

// LE37: an interrupted shift still ends its span, with work state interrupted.
func TestAnInterruptedShiftRecordsInterruptedWork(t *testing.T) {
	shift, spans, _ := interruptedShift(t)
	requireAttr(t, shift, otelrecord.AttrWorkState, otelrecord.WorkInterrupted)
	passes := childrenOf(spans, shift, iterationSeam)
	if len(passes) != 1 {
		t.Fatalf("the shift span has %d iteration children, want 1", len(passes))
	}
	requireAttr(t, passes[0], otelrecord.AttrAdapterResult, otelrecord.AdapterExited)
	requireNoAttr(t, passes[0], otelrecord.AttrAdapterExit)
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

// memoryFiles returns the bytes of each retained memory file of the working repository.
func memoryFiles(t *testing.T) [][]byte {
	t.Helper()
	root, err := git.Root()
	if err != nil {
		t.Fatal(err)
	}
	dir := otelrecord.MemoryDir(os.Getenv("BENCH_HOME"), root)
	entries, _ := os.ReadDir(dir)
	var files [][]byte
	for _, entry := range entries {
		raw, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		files = append(files, raw)
	}
	return files
}

// memoryShift runs a shift under gate whose adapter runs script, and returns the shift
// span, the raw record, and the memory files.
func memoryShift(t *testing.T, gate, script string, code int) (otelrecord.Span, []byte, [][]byte) {
	t.Helper()
	faultFixtureCore(t, gate, nil)
	withAgent(t, script)
	shift, raw, _ := runRecordedShift(t, code)
	return shift, raw, memoryFiles(t)
}

const (
	memoryWork = "echo work >> work.txt\necho MEMMARK >> " + notesFile + "\n"
)

// LE49 and LE51: a green shift retains its notes and records their size and digest.
func TestAGreenShiftRetainsItsNotes(t *testing.T) {
	shift, _, files := memoryShift(t, greenGate, memoryWork, 3)
	if len(files) != 1 || string(files[0]) != "MEMMARK\n" {
		t.Fatalf("memory files = %q, want one file with the notes", files)
	}
	sum := sha256.Sum256(files[0])
	requireAttr(t, shift, otelrecord.AttrMemoryState, otelrecord.MemoryRetained)
	requireAttr(t, shift, otelrecord.AttrMemoryBytes, "8")
	requireAttr(t, shift, otelrecord.AttrMemoryDigest, "sha256:"+hex.EncodeToString(sum[:]))
}

// LE50: a red shift that retained its worktree retains its notes.
func TestARedShiftRetainsItsNotes(t *testing.T) {
	_, _, files := memoryShift(t, "#!/usr/bin/env bash\nexit 1\n", memoryWork, 1)
	if len(files) != 1 || string(files[0]) != "MEMMARK\n" {
		t.Fatalf("memory files = %q, want one file with the notes", files)
	}
}

// LE52: the notes text never enters the record.
func TestTheRecordHoldsNoNotesText(t *testing.T) {
	if _, raw, _ := memoryShift(t, greenGate, memoryWork, 3); bytes.Contains(raw, []byte("MEMMARK")) {
		t.Fatalf("the record holds the notes text:\n%s", raw)
	}
}

// LE53, LE54, LE88, LE55, and LE56: each notes state the adapter leaves.
func TestEachNotesStateIsRecorded(t *testing.T) {
	secret := filepath.Join(t.TempDir(), "secret")
	if err := os.WriteFile(secret, []byte("SECRETMARK\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ name, script, state string }{
		{"symlink", "rm " + notesFile + "\nln -s '" + secret + "' " + notesFile + "\n", otelrecord.MemoryRefused},
		{"fifo", "rm " + notesFile + "\nmkfifo " + notesFile + "\n", otelrecord.MemoryRefused},
		{"oversized", "head -c 3000000 /dev/zero > " + notesFile + "\n", otelrecord.MemoryRefused},
		{"deleted", "rm " + notesFile + "\n", otelrecord.MemoryAbsent},
		{"empty", "true\n", otelrecord.MemoryRetained},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// A retention that opens the FIFO blocks, and the go test deadline reds it.
			shift, _, files := memoryShift(t, greenGate, tc.script, 4)
			requireAttr(t, shift, otelrecord.AttrMemoryState, tc.state)
			if tc.state != otelrecord.MemoryRetained {
				if len(files) != 0 {
					t.Fatalf("a %s notes file wrote %d memory files", tc.name, len(files))
				}
				return
			}
			requireAttr(t, shift, otelrecord.AttrMemoryBytes, "0")
		})
	}
}

// LE95 and LE102: a failed memory write changes neither the outcome nor the exit.
func TestAFailedMemoryWriteKeepsTheOutcome(t *testing.T) {
	root := faultFixtureCore(t, greenGate, nil)
	withDonePredicate(t, root)
	withAgent(t, "echo work >> work.txt\n")
	repo, _ := git.Root()
	home := os.Getenv("BENCH_HOME")
	if err := os.MkdirAll(otelrecord.Dir(home, repo), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(t.TempDir(), otelrecord.MemoryDir(home, repo)); err != nil {
		t.Fatal(err)
	}
	shift, _, _ := runRecordedShift(t, 0)
	requireAttr(t, shift, otelrecord.AttrMemoryState, otelrecord.MemoryFailed)
	requireAttr(t, shift, otelrecord.AttrShiftOutcome, string(OutcomeComplete))
}

// LE59: a second shift starts with empty notes, and no retained memory reaches its prompt.
func TestASecondShiftStartsWithEmptyNotes(t *testing.T) {
	faultFixtureCore(t, greenGate, nil)
	seen := t.TempDir()
	withAgent(t, "cat > '"+seen+"/prompt'\ncp "+notesFile+" '"+seen+"/notes'\necho MEMMARK >> "+notesFile+"\n")
	for run := range 2 {
		if code := Loop("memory shift", io.Discard, io.Discard); code != exitCodes[OutcomeNoOp] {
			t.Fatalf("Loop = %d, want no-op", code)
		}
		if files := memoryFiles(t); run == 0 && (len(files) != 1 || !bytes.Contains(files[0], []byte("MEMMARK"))) {
			t.Fatalf("the first shift kept memory files %q, want one with MEMMARK", files)
		}
	}
	notes, _ := os.ReadFile(filepath.Join(seen, "notes"))
	prompt, _ := os.ReadFile(filepath.Join(seen, "prompt"))
	if len(notes) != 0 || len(prompt) == 0 || bytes.Contains(prompt, []byte("MEMMARK")) {
		t.Fatalf("second shift notes %q, prompt %q, want empty notes and no MEMMARK", notes, prompt)
	}
}

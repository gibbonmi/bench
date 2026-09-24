package shift

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/otelrecord"
)

// shiftSeamAttr is the encoded seam attribute of a shift span line, and startMarker is
// the encoded record-start attribute that only a start line carries.
var (
	shiftSeamAttr = seamAttr(shiftSeam)
	startMarker   = fmt.Sprintf(`{"key":%q,"value":{"stringValue":%q}}`, otelrecord.AttrRecord, otelrecord.RecordStart)
)

// seamAttr is the encoded seam attribute of a span line of seam.
func seamAttr(seam string) string {
	return fmt.Sprintf(`{"key":%q,"value":{"stringValue":%q}}`, otelrecord.AttrSeam, seam)
}

// runRecordedShift runs Loop in the fixture's repository and returns the one finished
// shift span with the raw bytes of the repository's record.
func runRecordedShift(t *testing.T, wantCode int) (otelrecord.Span, []byte, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	if code := Loop("recorded shift", &stdout, &stderr); code != wantCode {
		t.Fatalf("Loop = %d, want %d; stdout:\n%s\nstderr:\n%s", code, wantCode, stdout.String(), stderr.String())
	}
	span, _, raw := recordedShift(t)
	return span, raw, stdout.String()
}

// recordedShift reads the record of the repository at the working directory and returns
// its one finished shift span, every finished span, and the raw record bytes.
func recordedShift(t *testing.T) (otelrecord.Span, []otelrecord.Span, []byte) {
	t.Helper()
	home := os.Getenv("BENCH_HOME")
	root, err := git.Root()
	if err != nil {
		t.Fatal(err)
	}
	spans, err := otelrecord.ReadSpans(home, root)
	if err != nil {
		t.Fatalf("read the record: %v", err)
	}
	var found []otelrecord.Span
	for _, span := range spans {
		if span.Seam == shiftSeam {
			found = append(found, span)
		}
	}
	if len(found) != 1 {
		t.Fatalf("the record holds %d shift spans, want 1: %+v", len(found), spans)
	}
	raw, err := os.ReadFile(otelrecord.Path(home, root))
	if err != nil {
		t.Fatal(err)
	}
	return found[0], spans, raw
}

// withDonePredicate adds an executable .bench/done.sh that exits zero, so the first
// committed iteration meets the objective and the shift ends complete.
func withDonePredicate(t *testing.T, root string) {
	t.Helper()
	done := filepath.Join(root, ".bench", "done.sh")
	if err := os.WriteFile(done, []byte("#!/usr/bin/env bash\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	runGitOutput(t, root, "add", "-A")
	runGitOutput(t, root, "commit", "-q", "-m", "done predicate")
}

func requireAttr(t *testing.T, span otelrecord.Span, key, want string) {
	t.Helper()
	if got, ok := span.Attributes[key]; !ok || got != want {
		t.Fatalf("shift span %s = %q (present %v), want %q; attributes %v", key, got, ok, want, span.Attributes)
	}
}

func requireNoAttr(t *testing.T, span otelrecord.Span, key string) {
	t.Helper()
	if got, ok := span.Attributes[key]; ok {
		t.Fatalf("shift span carries %s = %q, want no key", key, got)
	}
}

// completeShift runs a green one-iteration shift that commits and meets its objective.
func completeShift(t *testing.T) (otelrecord.Span, []byte, string) {
	t.Helper()
	root := faultFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	withDonePredicate(t, root)
	return runRecordedShift(t, 0)
}

// retainedShift runs a shift whose first gate is red over a dirty tree, so the shift
// retains its worktree for recovery.
func retainedShift(t *testing.T) (otelrecord.Span, []byte, string) {
	t.Helper()
	faultFixture(t, "#!/usr/bin/env bash\nexit 1\n")
	return runRecordedShift(t, 1)
}

// noOpShift runs a shift whose adapter makes no change.
func noOpShift(t *testing.T) (otelrecord.Span, []byte, string) {
	t.Helper()
	t.Setenv("BENCH_AGENT", "true")
	faultFixtureNoAgentOverride(t, "#!/usr/bin/env bash\nexit 0\n")
	return runRecordedShift(t, 4)
}

// LE22: a green one-iteration shift writes one start line and one end line for its span.
func TestAGreenShiftWritesAStartAndAnEndLine(t *testing.T) {
	_, raw, _ := completeShift(t)
	var starts, ends int
	for _, line := range bytes.Split(bytes.TrimSpace(raw), []byte("\n")) {
		if !bytes.Contains(line, []byte(shiftSeamAttr)) {
			continue
		}
		if bytes.Contains(line, []byte(startMarker)) {
			starts++
		} else {
			ends++
		}
	}
	if starts != 1 || ends != 1 {
		t.Fatalf("shift lines: %d start, %d end, want 1 and 1:\n%s", starts, ends, raw)
	}
}

// LE23: a shift that exhausts its branch-name retries still ends its span.
func TestAShiftThatExhaustsItsBranchRetriesRecordsUsage(t *testing.T) {
	var taken []string
	for i := 2; i <= branchCollisionRetries; i++ {
		taken = append(taken, fmt.Sprintf("-%d", i))
	}
	shiftCollisionFixture(t, taken...)
	span, _, _ := runRecordedShift(t, 2)
	requireAttr(t, span, otelrecord.AttrShiftOutcome, string(OutcomeUsage))
	requireAttr(t, span, otelrecord.AttrWorkState, otelrecord.WorkFailed)
}

// LE24: the span carries the adapter's base name and never its directory.
func TestTheShiftSpanCarriesTheAdapterBaseName(t *testing.T) {
	faultFixtureCore(t, "#!/usr/bin/env bash\nexit 0\n", nil)
	dir := filepath.Join(t.TempDir(), "DIRMARK dir")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	agent := filepath.Join(dir, "agent")
	if err := os.WriteFile(agent, []byte("#!/usr/bin/env bash\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("BENCH_AGENT", agent)
	span, raw, _ := runRecordedShift(t, 4)
	requireAttr(t, span, otelrecord.AttrAgent, "agent")
	if bytes.Contains(raw, []byte("DIRMARK")) {
		t.Fatalf("the record holds the adapter directory:\n%s", raw)
	}
}

// LE25: a declared tier and the cap reach the span.
func TestTheShiftSpanCarriesTheTierAndTheCap(t *testing.T) {
	root := faultFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	withDonePredicate(t, root)
	t.Setenv("BENCH_MODEL", "mid")
	t.Setenv("BENCH_MAX_ITERS", "2")
	span, _, _ := runRecordedShift(t, 0)
	requireAttr(t, span, otelrecord.AttrLineTier, "mid")
	requireAttr(t, span, otelrecord.AttrShiftCap, "2")
}

// LE26: a model value outside the tier list never reaches the record.
func TestTheShiftSpanCarriesNoTierForAnOperatorModel(t *testing.T) {
	root := faultFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	withDonePredicate(t, root)
	t.Setenv("BENCH_MODEL", "custom model")
	span, raw, _ := runRecordedShift(t, 0)
	requireNoAttr(t, span, otelrecord.AttrLineTier)
	if bytes.Contains(raw, []byte("custom model")) {
		t.Fatalf("the record holds the operator model text:\n%s", raw)
	}
}

// LE27: a complete shift maps to completed work and a green outcome.
func TestACompleteShiftRecordsCompletedWork(t *testing.T) {
	span, _, _ := completeShift(t)
	requireAttr(t, span, otelrecord.AttrShiftOutcome, string(OutcomeComplete))
	requireAttr(t, span, otelrecord.AttrWorkState, otelrecord.WorkCompleted)
	requireAttr(t, span, otelrecord.AttrOutcome, otelrecord.OutcomeGreen)
}

// LE28: a shift whose first gate is red maps to failed work and a red outcome.
func TestAFailedShiftRecordsFailedWork(t *testing.T) {
	span, _, _ := retainedShift(t)
	requireAttr(t, span, otelrecord.AttrShiftOutcome, string(OutcomeFailed))
	requireAttr(t, span, otelrecord.AttrWorkState, otelrecord.WorkFailed)
	requireAttr(t, span, otelrecord.AttrOutcome, otelrecord.OutcomeRed)
}

// LE29: a shift that commits and then reaches its cap maps to completed work.
func TestAnIncompleteShiftRecordsCompletedWork(t *testing.T) {
	faultFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	span, _, _ := runRecordedShift(t, 3)
	requireAttr(t, span, otelrecord.AttrShiftOutcome, string(OutcomeIncomplete))
	requireAttr(t, span, otelrecord.AttrWorkState, otelrecord.WorkCompleted)
}

// LE30: a no-op shift maps to completed work.
func TestANoOpShiftRecordsCompletedWork(t *testing.T) {
	span, _, _ := noOpShift(t)
	requireAttr(t, span, otelrecord.AttrShiftOutcome, string(OutcomeNoOp))
	requireAttr(t, span, otelrecord.AttrWorkState, otelrecord.WorkCompleted)
}

// LE31: a shift with one commit names the branch head as its subject.
func TestACommittedShiftCarriesTheBranchHead(t *testing.T) {
	root := faultFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	withDonePredicate(t, root)
	span, _, stdout := runRecordedShift(t, 0)
	head := strings.TrimSpace(runGitOutput(t, root, "rev-parse", shiftBranchName(t, stdout)))
	requireAttr(t, span, otelrecord.AttrSubjectID, head)
}

// LE32: a no-op shift claims no commit.
func TestANoOpShiftCarriesNoSubject(t *testing.T) {
	span, _, _ := noOpShift(t)
	requireNoAttr(t, span, otelrecord.AttrSubjectID)
}

// LE33: a retained worktree reaches the span as its base name.
func TestARetainedShiftCarriesTheRecoveryBaseName(t *testing.T) {
	span, _, stdout := retainedShift(t)
	requireAttr(t, span, otelrecord.AttrRecoveryKind, "worktree")
	requireAttr(t, span, otelrecord.AttrRecoveryKey, filepath.Base(shiftWorktreePath(t, stdout)))
}

// LE34: no line of a retained shift holds the Bench home path.
func TestARetainedShiftRecordHoldsNoHomePath(t *testing.T) {
	_, raw, _ := retainedShift(t)
	if home := os.Getenv("BENCH_HOME"); bytes.Contains(raw, []byte(home)) {
		t.Fatalf("the record holds the Bench home %s:\n%s", home, raw)
	}
}

// LE35: a green shift released its worktree.
func TestAGreenShiftRecordsReleasedCleanup(t *testing.T) {
	span, _, _ := completeShift(t)
	requireAttr(t, span, otelrecord.AttrCleanup, otelrecord.CleanupReleased)
}

// LE107: a shift that left no recovery pointer records kind none and no key.
func TestAGreenShiftRecordsNoRecoveryKind(t *testing.T) {
	span, _, _ := completeShift(t)
	requireAttr(t, span, otelrecord.AttrRecoveryKind, RecoveryNone)
	requireNoAttr(t, span, otelrecord.AttrRecoveryKey)
}

// A teardown that fails before the release records no cleanup, not a release.
func TestAFailedTeardownRecordsNoCleanup(t *testing.T) {
	t.Setenv("BENCH_AGENT", "true")
	faultFixtureNoAgentOverride(t, "#!/usr/bin/env bash\nexit 0\n")
	armFault(t, func(step shiftStep) error {
		if step == stepTeardown {
			return fmt.Errorf("injected teardown failure")
		}
		return nil
	})
	span, _, _ := runRecordedShift(t, 1)
	requireAttr(t, span, otelrecord.AttrCleanup, otelrecord.CleanupNone)
}

// The span starts before the worktree acquire, so a failed acquire still ends it.
func TestAFailedAcquireStillEndsTheShiftSpan(t *testing.T) {
	faultFixture(t, "#!/usr/bin/env bash\nexit 0\n")
	home := os.Getenv("BENCH_HOME")
	if err := os.MkdirAll(home, 0o700); err != nil {
		t.Fatal(err)
	}
	// A regular file where the worktree pool directory goes makes the acquire fail.
	if err := os.WriteFile(filepath.Join(home, "worktrees"), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	span, _, _ := runRecordedShift(t, 2)
	requireAttr(t, span, otelrecord.AttrShiftOutcome, string(OutcomeUsage))
	requireAttr(t, span, otelrecord.AttrCleanup, otelrecord.CleanupNone)
	requireAttr(t, span, otelrecord.AttrRecoveryKind, RecoveryNone)
}

// LE36: a red shift that retained its worktree records the retention.
func TestARetainedShiftRecordsRetainedCleanup(t *testing.T) {
	span, _, _ := retainedShift(t)
	requireAttr(t, span, otelrecord.AttrCleanup, otelrecord.CleanupRetained)
}

// shiftIntent returns the live intent entry of key in the working repository.
func shiftIntent(t *testing.T, key string) (intent.Entry, bool) {
	t.Helper()
	root, err := git.Root()
	if err != nil {
		t.Fatal(err)
	}
	entries, err := intent.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.Key == key {
			return entry, true
		}
	}
	return intent.Entry{}, false
}

// LE61: a green shift ends its intent.
func TestAGreenShiftEndsItsIntent(t *testing.T) {
	span, _, _ := completeShift(t)
	if entry, ok := shiftIntent(t, span.Attributes[otelrecord.AttrIntentKey]); ok {
		t.Fatalf("the finished shift left a live intent entry: %+v", entry)
	}
}

// LE63: the shift records the lease line its acquire wrote.
func TestAShiftRecordsItsLease(t *testing.T) {
	faultFixtureCore(t, "#!/usr/bin/env bash\nexit 1\n", nil)
	lease := filepath.Join(t.TempDir(), "lease")
	withAgent(t, "echo work >> work.txt\ncp \"$(git rev-parse --git-path "+git.BenchLeaseFilename+")\" '"+lease+"'\n")
	span, _, _ := runRecordedShift(t, 1)
	written, err := os.ReadFile(lease)
	if err != nil {
		t.Fatal(err)
	}
	entry, ok := shiftIntent(t, span.Attributes[otelrecord.AttrIntentKey])
	if !ok || entry.Lease == "" || entry.Lease+"\n" != string(written) {
		t.Fatalf("intent lease = %q (live %v), want the acquired lease line %q", entry.Lease, ok, written)
	}
}

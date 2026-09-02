//go:build system

package systemtest

import (
	"fmt"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/otelrecord"
)

// TestOtelCrashKeepsStartedPhaseLine kills a `bench gate` run while one phase is still
// running. The processor writes each line synchronously at span start, so the record
// holds the started phase span before the kill, and the kill loses at most the line in
// flight. An end-only or buffered writer would leave the record without that phase.
//
// Row OT25.
func TestOtelCrashKeepsStartedPhaseLine(t *testing.T) {
	home := filepath.Join(t.TempDir(), "bench-home")
	scaffold := scaffoldRecordedGateRepo(t, home, `{"phases":[{"name":"slow","argv":["sleep","30"]}]}`)

	run, err := owner.startProcessGroup(scaffold.path, scaffold.environment(t, home), "bash", scaffold.wrapper, "gate", "--fresh")
	if err != nil {
		t.Fatal(err)
	}
	// The wait polls the record file itself, because the record is the only evidence the
	// phase has begun. The gate's own output is buffered until the run ends, which the
	// kill prevents.
	awaitStartedSpan(t, home, "slow")
	if err := run.kill(syscall.SIGKILL); err != nil {
		t.Fatalf("kill the gate process group: %v", err)
	}
	if !run.wait() {
		t.Fatal("the gate run exited on its own, so the record was not read after a crash")
	}

	// Every line the killed run wrote still parses, and the started phase line survives.
	var slow recordedSpan
	for _, span := range readRecordLines(t, home, true) {
		if span.Name == "slow" {
			slow = span
		}
	}
	if slow.Name == "" {
		t.Fatal("the record has no slow phase span after the kill")
	}
	if slow.Ended {
		t.Fatalf("slow phase span = %+v, want the started line of a phase the kill interrupted", slow)
	}
	if slow.Attributes[otelrecord.AttrRecord] != otelrecord.RecordStart {
		t.Fatalf("slow phase record marker = %q, want %q", slow.Attributes[otelrecord.AttrRecord], otelrecord.RecordStart)
	}
	if slow.Parent == "" {
		t.Fatalf("slow phase span = %+v, want the started root gate span as its parent", slow)
	}
}

// awaitStartedSpan blocks until the record below home holds a started span of the given
// name. It reads leniently, because the run still appends while the wait reads. The
// window contains one gate run's start-up: the run reads git worktree state before it
// enters the phase table, so that read is the longest bound the wait has to outlast.
func awaitStartedSpan(t *testing.T, home, name string) {
	t.Helper()
	window := bounds.TestDeadline(bounds.WorktreeListTimeout)
	deadline := time.Now().Add(window)
	for time.Now().Before(deadline) {
		if matches, _ := filepath.Glob(filepath.Join(home, "otel", "*", "traces.jsonl")); len(matches) == 1 {
			for _, span := range readRecordLines(t, home, false) {
				if span.Name == name {
					return
				}
			}
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatal(bounds.TestTimeoutVerdict(fmt.Sprintf("a started %s span in the record below %s", name, home), window))
}

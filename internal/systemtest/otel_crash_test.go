//go:build system

package systemtest

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	t.Parallel()
	home := filepath.Join(t.TempDir(), "bench-home")
	pidFile := filepath.Join(t.TempDir(), "slow.pid")
	// The phase end waits only until the pid file exists. So the phase writes its pid to a
	// temporary name first, and the rename publishes the whole pid at once.
	argv, err := json.Marshal([]string{"sh", "-c", `echo $$ > "$0.tmp"; mv "$0.tmp" "$0"; exec sleep 30`, pidFile})
	if err != nil {
		t.Fatal(err)
	}
	scaffold := scaffoldRecordedGateRepo(t, home, `{"phases":[{"name":"slow","argv":`+string(argv)+`}]}`)

	run, err := owner.startProcessGroup(scaffold.path, scaffold.environment(t, home), "bash", scaffold.wrapper, "gate", "--fresh")
	if err != nil {
		t.Fatal(err)
	}
	endSlow := phaseEnd(t, home, pidFile, "slow")
	// The wait polls the record file itself, because the record is the only evidence the
	// phase has begun. The gate's own output is buffered until the run ends, which the
	// kill prevents.
	awaitSpan(t, home, "slow", false)
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
	endSlow()
}

// phaseEnd answers a function that ends the phase whose pid is in pidFile. The end kills
// the phase and waits until the record below home holds the phase's end line. The gate
// starts each phase in a process group of its own, so the kill of the gate's group leaves
// the phase running. When the phase exits, its runner writes the end line below home. That
// write must come before the test removes its temporary directory, so a cleanup also runs
// the end when the test stops early. The end runs once, because a second kill can reach a
// reused pid.
func phaseEnd(t *testing.T, home, pidFile, name string) func() {
	t.Helper()
	ended := false
	end := func() {
		t.Helper()
		if ended {
			return
		}
		ended = true
		waitForFile(t, pidFile, "the "+name+" phase to write its pid file "+pidFile)
		data, err := os.ReadFile(pidFile)
		if err != nil {
			t.Fatal(err)
		}
		pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
		if err != nil {
			t.Fatalf("the %s phase pid file %s: %v", name, pidFile, err)
		}
		// A phase that already exited has written, or will write, the same end line.
		if err := syscall.Kill(pid, syscall.SIGKILL); err != nil && !errors.Is(err, syscall.ESRCH) {
			t.Fatalf("kill the %s phase %d: %v", name, pid, err)
		}
		awaitSpan(t, home, name, true)
	}
	t.Cleanup(end)
	return end
}

// awaitSpan blocks until the record below home holds a span of the given name, ended or
// started as asked. It reads leniently, because the run still appends while the wait
// reads. A wait for the started line contains one gate run's start-up. The run reads git
// worktree state before it enters the phase table, and that read is the longest bound the
// window must outlast. A wait for the end line follows the kill of the phase, and the
// runner writes that line when it reaps the phase. That wait has no bound of its own, so
// the same window contains it.
func awaitSpan(t *testing.T, home, name string, ended bool) {
	t.Helper()
	window := bounds.TestDeadline(bounds.WorktreeListTimeout)
	deadline := time.Now().Add(window)
	for time.Now().Before(deadline) {
		if matches, _ := filepath.Glob(filepath.Join(home, "otel", "*", "traces.jsonl")); len(matches) == 1 {
			for _, span := range readRecordLines(t, home, false) {
				if span.Name == name && span.Ended == ended {
					return
				}
			}
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatal(bounds.TestTimeoutVerdict(fmt.Sprintf("a %s span (ended %t) in the record below %s", name, ended, home), window))
}

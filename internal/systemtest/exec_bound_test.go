//go:build system

package systemtest

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
)

// The expectations here are authored apart from the response owner: a bounded exec
// response over 10 lines prints its first 4 lines, one spill line, and its last 5 lines.
// The spec fixes these counts, and only the owner package may read the line value from
// the policy registry.

// execBoundFixture answers a scaffolded repository and one active assignment in it.
func execBoundFixture(t *testing.T) (recordedPublicationRepo, systemLandingWorktree) {
	t.Helper()
	fixture := scaffoldRecordedPublicationRepo(t)
	return fixture, systemCreateLandingWorktree(t, fixture.root, fixture.home, "exec-bound", "exec bound")
}

func execBoundArgs(source systemLandingWorktree, child ...string) []string {
	return append([]string{"worktree", "exec", source.assignment, "--"}, child...)
}

// runBoundedExec runs the selected executable's exec verb from the primary checkout with
// input on its stdin.
func runBoundedExec(t *testing.T, fixture recordedPublicationRepo, source systemLandingWorktree, input string, child ...string) processResult {
	t.Helper()
	if err := owner.observeSelected(); err != nil {
		t.Fatal(err)
	}
	return owner.runWithInput(fixture.root, fixture.environment(fixture.home), input, owner.selected.path, execBoundArgs(source, child...)...)
}

// seqLines is what `seq first last` prints.
func seqLines(first, last int) []string {
	var lines []string
	for i := first; i <= last; i++ {
		lines = append(lines, strconv.Itoa(i))
	}
	return lines
}

// spilledPath answers the path field of a projected response's spill line, after it
// checks that the response is 10 lines of head, spill line, and tail.
func spilledPath(t *testing.T, stdout string, head, tail []string) string {
	t.Helper()
	lines := strings.Split(strings.TrimSuffix(stdout, "\n"), "\n")
	if len(lines) != 10 {
		t.Fatalf("stdout = %q, want 10 lines", stdout)
	}
	for i, want := range head {
		if lines[i] != want {
			t.Errorf("head line %d = %q, want %q", i+1, lines[i], want)
		}
	}
	for i, want := range tail {
		if got := lines[5+i]; got != want {
			t.Errorf("tail line %d = %q, want %q", i+1, got, want)
		}
	}
	_, path, found := strings.Cut(lines[4], ",path=")
	if !strings.HasPrefix(lines[4], "spilled{") || !found || !strings.HasSuffix(path, "}") {
		t.Fatalf("fifth line = %q, want the spill line", lines[4])
	}
	return strings.TrimSuffix(path, "}")
}

func readSpillFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func joinLines(lines []string) string { return strings.Join(lines, "\n") + "\n" }

// BO26: exec bounds its child's output to lines 1 to 4, the spill line, and lines 36 to
// 40, and the spill file holds the complete output.
func TestExecBoundsChildOutput(t *testing.T) {
	fixture, source := execBoundFixture(t)
	result := runBoundedExec(t, fixture, source, "", "sh", "-c", "seq 1 40")
	if result.code != 0 {
		t.Fatalf("exec = (%d, %q, %q), want exit 0", result.code, result.stdout, result.stderr)
	}
	path := spilledPath(t, result.stdout, seqLines(1, 4), seqLines(36, 40))
	if got, want := readSpillFile(t, path), joinLines(seqLines(1, 40)); got != want {
		t.Fatalf("spill file = %q, want %q", got, want)
	}
}

// BO27: the bound keeps the child's own exit code.
func TestExecKeepsChildExitUnderBound(t *testing.T) {
	fixture, source := execBoundFixture(t)
	result := runBoundedExec(t, fixture, source, "", "sh", "-c", "seq 1 40; exit 7")
	if result.code != 7 {
		t.Fatalf("exec = (%d, %q, %q), want the child's exit 7", result.code, result.stdout, result.stderr)
	}
	spilledPath(t, result.stdout, seqLines(1, 4), seqLines(37, 40))
}

// BO28: a heredoc on exec's stdin reaches the child, and the child's 12 lines spill.
func TestExecForwardsStdinUnderBound(t *testing.T) {
	fixture, source := execBoundFixture(t)
	result := runBoundedExec(t, fixture, source, "seq 1 12\n", "sh")
	if result.code != 0 {
		t.Fatalf("exec = (%d, %q, %q), want exit 0", result.code, result.stdout, result.stderr)
	}
	path := spilledPath(t, result.stdout, seqLines(1, 4), seqLines(8, 12))
	if got, want := readSpillFile(t, path), joinLines(seqLines(1, 12)); got != want {
		t.Fatalf("spill file = %q, want the child's stdin script output %q", got, want)
	}
}

// BO29: a SIGINT after the child printed 12 lines gives exit 130, and the spill file
// holds those lines.
func TestExecInterruptKeepsOutput(t *testing.T) {
	fixture, source := execBoundFixture(t)
	marker := filepath.Join(t.TempDir(), "printed")
	cmd, stdout, stderr := systemStartSelected(t, fixture.root, fixture.environment(fixture.home),
		execBoundArgs(source, "sh", "-c", "seq 1 12; : > "+marker+"; sleep 30")...)
	// The child writes the marker after its output and then idles, so this wait covers a
	// start-up handshake with no window of its own.
	window := bounds.TestDeadline(0)
	deadline := time.Now().Add(window)
	for {
		if _, err := os.Stat(marker); err == nil {
			break
		}
		if time.Now().After(deadline) {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
			t.Fatal(bounds.TestTimeoutVerdict("the exec child to print its 12 lines", window))
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
		t.Fatal(err)
	}
	if code := systemExitCode(cmd.Wait()); code != 130 {
		t.Fatalf("interrupted exec = (%d, %q, %q), want exit 130", code, stdout.String(), stderr.String())
	}
	path := spilledPath(t, stdout.String(), seqLines(1, 4), nil)
	if got, want := readSpillFile(t, path), joinLines(seqLines(1, 12)); !strings.HasPrefix(got, want) {
		t.Fatalf("spill file = %q, want the 12 printed lines %q first", got, want)
	}
}

// BO68: a child that prints 64 MiB of short lines completes, the spill file holds every
// byte, and exec's peak memory stays far below the output. An owner that kept the
// complete output would hold at least 64 MiB.
func TestExecStreamsLargeChild(t *testing.T) {
	const size = 64 << 20
	const peakLimit = 48 << 20
	fixture, source := execBoundFixture(t)
	cmd, stdout, stderr := systemStartSelected(t, fixture.root, fixture.environment(fixture.home),
		execBoundArgs(source, "sh", "-c", fmt.Sprintf("yes line | head -c %d", size))...)
	waited := make(chan error, 1)
	done := make(chan struct{})
	go func() {
		waited <- cmd.Wait()
		close(done)
	}()
	peak := peakResidentBytes(cmd.Process.Pid, done)
	if code := systemExitCode(<-waited); code != 0 {
		t.Fatalf("exec = (%d, %q, %q), want exit 0", code, stdout.String(), stderr.String())
	}
	path := spilledPath(t, stdout.String(), []string{"line", "line", "line", "line"}, []string{"line", "line", "line", "line", "line"})
	if info, err := os.Stat(path); err != nil || info.Size() != size {
		t.Fatalf("spill file = (%v, %v), want %d bytes", info, err, size)
	}
	if peak == 0 || peak >= peakLimit {
		t.Fatalf("exec peak resident set = %d bytes, want a sample below %d while it streams %d bytes", peak, peakLimit, size)
	}
}

// peakResidentBytes samples the resident set of pid until done closes, and answers the
// largest sample. The wait status is no source here: on Linux, a child that the Go
// runtime starts reports the parent's peak through its own, so it would grade this test
// process rather than exec. `ps -o rss=` reads the process's own set in KiB on Linux
// and on Darwin.
func peakResidentBytes(pid int, done <-chan struct{}) int64 {
	var peak int64
	for {
		select {
		case <-done:
			return peak
		default:
		}
		sample := owner.runAt("", nil, "ps", "-o", "rss=", "-p", strconv.Itoa(pid))
		if kib, err := strconv.ParseInt(strings.TrimSpace(sample.stdout), 10, 64); err == nil {
			peak = max(peak, kib*1024)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

// BO69: a background descendant that holds the output pipe open cannot hang exec. Exec
// returns within its wait delay with the child's exit code and the child's line.
func TestExecReturnsAtChildExit(t *testing.T) {
	fixture, source := execBoundFixture(t)
	window := bounds.TestDeadline(bounds.ExecWaitDelay)
	started := time.Now()
	result := runBoundedExec(t, fixture, source, "", "sh", "-c", "sleep 30 & echo up")
	elapsed := time.Since(started)
	if result.code != 0 || result.stdout != "up\n" {
		t.Fatalf("exec = (%d, %q, %q), want exit 0 and the line up", result.code, result.stdout, result.stderr)
	}
	if elapsed >= window {
		t.Fatal(bounds.TestTimeoutVerdict("exec to return at its child's exit", window))
	}
}

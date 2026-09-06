package probe

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/testreport"
	"github.com/gibbonmi/bench/internal/usage"
)

// interruptHelperEnv carries the fixture root to the helper half of the interrupt row.
// The fact under test is what a signal does to the verb, so the verb has to run in a
// process the signal may end without taking the assertions with it.
const interruptHelperEnv = "BENCH_PROBE_INTERRUPT_ROOT"

// PB12: the copy exists on disk, at its declared mode, while the run is in flight.
func TestProbePreservesTheSubjectBeforeTheRun(t *testing.T) {
	f := newFixture(t)
	listing := filepath.Join(scratchDir(t), "preserved")
	body := filepath.Join(filepath.Dir(listing), "preserved-body")
	// The copy is read while the run holds it, because a proven restore removes it before
	// the verb answers and the test would then read an absent file.
	hook := "find \"$BENCH_HOME/probe\" -type f -printf '%m %p\\n' > " + sanitize.ShellQuote(listing) + "\n" +
		"find \"$BENCH_HOME/probe\" -type f -exec cat {} + > " + sanitize.ShellQuote(body) + "\n"
	installStubGo(t, f, cannedFailure("caught"), 1, hook)
	if _, code := runProbe(t, "clamp.go", "--swap", "n < 0", "--with", "n > 0", "--package", "./", "--run", "^TestClampNegative$"); code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	observed, err := os.ReadFile(listing)
	if err != nil {
		t.Fatal(err)
	}
	line := strings.TrimSpace(string(observed))
	prefix := "600 " + filepath.Join(f.home, "probe") + string(filepath.Separator)
	if !strings.HasPrefix(line, prefix) || !strings.HasSuffix(line, "clamp.go") {
		t.Fatalf("preserved listing = %q, want one 0600 clamp.go below %s", line, prefix)
	}
	copied, err := os.ReadFile(body)
	if err != nil || string(copied) != clampSource {
		t.Fatalf("preserved copy = (%q, %v), want the start bytes", copied, err)
	}
}

// PB14: a restore that cannot write keeps the copy, names it, and exits 2, so a caller
// whose tree still holds the mutation has the file it needs to put back by hand.
func TestProbeReportsARestoreFailure(t *testing.T) {
	f := newFixture(t)
	installStubGo(t, f, cannedFailure("caught"), 1, "chmod 0500 .\n")
	t.Cleanup(func() { _ = os.Chmod(f.root, 0o755) })
	out, code := runProbe(t, "clamp.go", "--swap", "n < 0", "--with", "n > 0", "--package", "./", "--run", "^TestClampNegative$")
	if code != 2 {
		t.Fatalf("exit = %d, want 2\n%s", code, out)
	}
	row := probeRow(t, "restore-failed", "clamp.go", "swap", "failed", 1, "no")
	if !strings.HasPrefix(out, row) {
		t.Fatalf("stdout = %q, want %q first", out, row)
	}
	kept := preservedCopy(t, f)
	if !strings.Contains(out, "preserved[1]{path,reason}:\n  "+kept) {
		t.Fatalf("stdout = %q, want the preserved row naming %s", out, kept)
	}
	saved, err := os.ReadFile(kept)
	if err != nil || string(saved) != clampSource {
		t.Fatalf("preserved copy = (%q, %v), want the start bytes", saved, err)
	}
}

// PB43: the read-back is compared with the start bytes, not with the copy, so a copy that
// changed under the run cannot certify its own damage as a restore.
func TestProbeReportsAReadBackMismatch(t *testing.T) {
	f := newFixture(t)
	hook := "find \"$BENCH_HOME/probe\" -type f -exec truncate -s 0 {} +\n"
	installStubGo(t, f, cannedFailure("caught"), 1, hook)
	out, code := runProbe(t, "clamp.go", "--swap", "n < 0", "--with", "n > 0", "--package", "./", "--run", "^TestClampNegative$")
	if code != 2 {
		t.Fatalf("exit = %d, want 2\n%s", code, out)
	}
	row := probeRow(t, "restore-failed", "clamp.go", "swap", "failed", 1, "no")
	if !strings.HasPrefix(out, row) {
		t.Fatalf("stdout = %q, want %q first", out, row)
	}
	want := "the restored bytes differ from the bytes read at the start"
	if !strings.Contains(out, want) {
		t.Fatalf("stdout = %q, want the preserved row's reason %q", out, want)
	}
}

// PB41: the restore runs before the render, so a subject the encoder cannot carry still
// leaves a clean tree behind its refusal.
func TestProbeRestoresBeforeARenderRefusal(t *testing.T) {
	f := newFixture(t)
	name := "data\aalert.txt"
	subject := filepath.Join(f.root, name)
	writeFixtureFile(t, subject, "alpha\n", 0o644)
	installStubGo(t, f, cannedPass, 0, "")
	out, code := runProbe(t, name, "--swap", "alpha", "--with", "beta", "--package", "./", "--run", "^TestClampPositive$")
	if code != 1 {
		t.Fatalf("exit = %d, want 1\n%s", code, out)
	}
	if !strings.HasPrefix(out, "error: unrepresentable TOON cell — ") {
		t.Fatalf("stdout = %q, want the shared render error", out)
	}
	restored, err := os.ReadFile(subject)
	if err != nil || string(restored) != "alpha\n" {
		t.Fatalf("subject = (%q, %v), want the start bytes", restored, err)
	}
	requireHomeEmpty(t, f)
}

// PB16: the restore is deferred, so an interrupt during the focused run still puts the
// subject back before the verb answers.
func TestProbeRestoresOnInterrupt(t *testing.T) {
	if root := os.Getenv(interruptHelperEnv); root != "" {
		runInterruptHelper(t, root)
		return
	}
	f := newFixture(t)
	installStubGo(t, f, cannedPass, 0, "sleep 60\n")
	helper := exec.Command(os.Args[0], "-test.run=^TestProbeRestoresOnInterrupt$", "-test.timeout="+interruptDeadline().String())
	helper.Env = append(os.Environ(), interruptHelperEnv+"="+f.root)
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
	awaitStubStart(t, f)
	if err := helper.Process.Signal(syscall.SIGINT); err != nil {
		t.Fatal(err)
	}
	select {
	case out := <-answered:
		want := probeRow(t, "invalid", "clamp.go", "swap", "interrupted", 0, "yes")
		if !strings.Contains(out, want) {
			t.Fatalf("helper stdout = %q, want the interrupted row %q", out, want)
		}
	case <-time.After(interruptDeadline()):
		t.Fatalf("the interrupted probe did not answer within %s", interruptDeadline())
	}
	requireSubjectBytes(t, f, clampSource)
}

// runInterruptHelper is the signalled half: it runs the verb for real and prints what the
// verb answered, so the parent grades production's own deferred restore.
func runInterruptHelper(t *testing.T, root string) {
	t.Helper()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	out, code := Command([]string{"clamp.go", "--swap", "n < 0", "--with", "n > 0", "--package", "./", "--run", "^TestClampNegative$"})
	fmt.Print(out)
	fmt.Printf("exit=%d\n", code)
}

// awaitStubStart waits for the focused run's own `go test` child, not for the `go list`
// the run binary's verification takes on the way in.
func awaitStubStart(t *testing.T, f *fixture) {
	t.Helper()
	expiry := time.Now().Add(interruptDeadline())
	for time.Now().Before(expiry) {
		body, err := os.ReadFile(f.marker)
		if err == nil {
			for _, line := range strings.Split(string(body), "\n") {
				if line == "test" {
					return
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("the focused run's Go child did not start within %s", interruptDeadline())
}

func interruptDeadline() time.Duration { return bounds.TestDeadline(runbinary.BuilderCancelGrace) }

// PB17: an executable subject stays executable across the mutation and the restore.
func TestProbeKeepsTheSubjectMode(t *testing.T) {
	f := newFixture(t)
	script := filepath.Join(f.root, "run.sh")
	writeFixtureFile(t, script, "#!/bin/sh\necho alpha\n", 0o755)
	installStubGo(t, f, cannedPass, 0, "")
	if _, code := runProbe(t, "run.sh", "--swap", "alpha", "--with", "beta", "--package", "./", "--run", "^TestClampPositive$"); code != 1 {
		t.Fatalf("exit = %d, want 1", code)
	}
	info, err := os.Lstat(script)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o755 {
		t.Fatalf("subject mode = %v, want 0755", info.Mode().Perm())
	}
}

// PB39: a subject whose name carries a space renders through the encoder, and the
// expectation derives through the same call.
func TestProbeRendersASubjectWithASpace(t *testing.T) {
	f := newFixture(t)
	writeFixtureFile(t, filepath.Join(f.root, "my pkg", "clamp.go"), "package spaced\n\nconst alpha = 1\n", 0o644)
	installStubGo(t, f, cannedPass, 0, "")
	out, code := runProbe(t, "my pkg/clamp.go", "--swap", "alpha", "--with", "beta", "--package", "./", "--run", "^TestClampPositive$")
	if code != 1 {
		t.Fatalf("exit = %d, want 1\n%s", code, out)
	}
	if want := probeRow(t, "silent", "my pkg/clamp.go", "swap", "passed", 0, "yes"); !strings.HasPrefix(out, want) {
		t.Fatalf("stdout = %q, want %q first", out, want)
	}
}

// PB40: a subject whose last line has no newline restores byte-exact.
func TestProbeRestoresAFileWithoutTrailingNewline(t *testing.T) {
	f := newFixture(t)
	subject := filepath.Join(f.root, "notrailing.txt")
	writeFixtureFile(t, subject, "alpha", 0o644)
	installStubGo(t, f, cannedPass, 0, "")
	out, code := runProbe(t, "notrailing.txt", "--swap", "alpha", "--with", "beta", "--package", "./", "--run", "^TestClampPositive$")
	if code != 1 {
		t.Fatalf("exit = %d, want 1\n%s", code, out)
	}
	if want := probeRow(t, "silent", "notrailing.txt", "swap", "passed", 0, "yes"); !strings.HasPrefix(out, want) {
		t.Fatalf("stdout = %q, want %q first", out, want)
	}
	got, err := os.ReadFile(subject)
	if err != nil || string(got) != "alpha" {
		t.Fatalf("restored bytes = (%q, %v), want %q", got, err, "alpha")
	}
}

// PB46: the verb writes only the preserved copy and the subject, so a completed probe
// adds no entry to the administration directory, the Bench home, or the tree.
func TestProbeWritesNoRecord(t *testing.T) {
	f := newFixture(t)
	installStubGo(t, f, cannedFailure("caught"), 1, "")
	before := listRecordSurfaces(t, f)
	if _, code := runProbe(t, "clamp.go", "--swap", "n < 0", "--with", "n > 0", "--package", "./", "--run", "^TestClampNegative$"); code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	if after := listRecordSurfaces(t, f); before != after {
		t.Fatalf("record surfaces changed:\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

// listRecordSurfaces answers a stable listing of the three places a record could land.
func listRecordSurfaces(t *testing.T, f *fixture) string {
	t.Helper()
	var lines []string
	for _, dir := range []string{f.home, filepath.Join(f.root, ".git"), f.root} {
		walked := dir
		if err := filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if walked == f.root && entry.IsDir() && entry.Name() == ".git" {
				return filepath.SkipDir
			}
			lines = append(lines, path)
			return nil
		}); err != nil {
			t.Fatal(err)
		}
	}
	return strings.Join(lines, "\n")
}

// PB7: the check form builds the request `bench test` builds for the same check name, so
// the two verbs cannot disagree about what a named check selects.
func TestProbeSelectsTheCheckForm(t *testing.T) {
	f := newFixture(t)
	parsed, line, code := usage.Parse(grammar, []string{"clamp.go", "--omit", "n < 0", "--check", "line-routing"})
	if line != "" {
		t.Fatalf("parse = (%q, %d), want an accepted selection", line, code)
	}
	if got := strings.Join(selectionArgs(parsed), " "); got != "--check line-routing" {
		t.Fatalf("selection args = %q, want %q", got, "--check line-routing")
	}
	got, line, _ := testreport.Prepare(f.root, selectionArgs(parsed))
	if line != "" {
		t.Fatalf("prepare refused the probe's check form: %q", line)
	}
	want, line, _ := testreport.Prepare(f.root, []string{"--check", "line-routing"})
	if line != "" {
		t.Fatalf("prepare refused the reference check form: %q", line)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("probe request = %#v, want %#v", got, want)
	}
}

// PB38: the verdict mapping is the one place an outcome becomes a word and an exit.
func TestVerdictExitCodes(t *testing.T) {
	for _, tc := range []struct {
		kind    testreport.OutcomeKind
		verdict string
		code    int
	}{
		{testreport.OutcomeFailed, "bit", 0},
		{testreport.OutcomePassed, "silent", 1},
		{testreport.OutcomeBuildFailed, "invalid", 1},
		{testreport.OutcomeNoTestRun, "invalid", 1},
		{testreport.OutcomeRefused, "invalid", 1},
		{testreport.OutcomeInterrupted, "invalid", 1},
	} {
		t.Run(string(tc.kind), func(t *testing.T) {
			verdict, code := verdictFor(tc.kind)
			if verdict != tc.verdict || code != tc.code {
				t.Fatalf("verdictFor(%q) = (%q, %d), want (%q, %d)", tc.kind, verdict, code, tc.verdict, tc.code)
			}
		})
	}
}

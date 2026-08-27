package gate

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/gocache"
)

const (
	gateHolderRootEnv    = "BENCH_TEST_GATE_HOLDER_ROOT"
	gateHolderReadyEnv   = "BENCH_TEST_GATE_HOLDER_READY"
	gateHolderReleaseEnv = "BENCH_TEST_GATE_HOLDER_RELEASE"
)

func TestGateRunRefusesInFlightExecution(t *testing.T) {
	root := failureOutcomeFixture(t)
	outcomeWrite(t, root, ".gate-wait", "\n", 0o644)
	finished := make(chan Result, 1)
	go func() { finished <- Execute(context.Background(), root, &bytes.Buffer{}, &bytes.Buffer{}) }()
	awaitGateMarker(t, filepath.Join(root, ".gate-running"))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var stdout, stderr bytes.Buffer
	result := Execute(ctx, root, &stdout, &stderr)
	if result.ActionExit != 1 {
		t.Fatalf("contended result = %#v, stderr=%q", result, stderr.String())
	}
	if !strings.Contains(stderr.String(), "gate execution already in progress") {
		t.Fatalf("stderr = %q", stderr.String())
	}
	wantOwner := fmt.Sprintf("gate owner: pid %d (alive)", os.Getpid())
	if !strings.Contains(stderr.String(), wantOwner) {
		t.Fatalf("stderr = %q, want %q", stderr.String(), wantOwner)
	}
	if got := outcomeRuns(t, root); got != 1 {
		t.Fatalf("runs while lock is held = %d, want 1", got)
	}
	outcomeWrite(t, root, ".gate-release", "\n", 0o644)
	awaitGateResult(t, finished)
}

func TestGateRunExternalHolderDemotesReusableGreen(t *testing.T) {
	if root := os.Getenv(gateHolderRootEnv); root != "" {
		runGateLockHolder(t, root, os.Getenv(gateHolderReadyEnv), os.Getenv(gateHolderReleaseEnv))
		return
	}

	root := failureOutcomeFixture(t)
	if result := Execute(context.Background(), root, &bytes.Buffer{}, &bytes.Buffer{}); result.ActionExit != 0 {
		t.Fatalf("green result = %#v", result)
	}
	if inspection := Inspect(root); inspection.State != Ready || inspection.Status != "green" || !inspection.ReusableGreen {
		t.Fatalf("before holder inspection = %#v", inspection)
	}
	startGateLockHolder(t, root)
	var stdout, stderr bytes.Buffer
	if got := RunCommand([]string{"--fresh", root}, &stdout, &stderr); got != 1 {
		t.Fatalf("contended exit = %d, stderr=%q", got, stderr.String())
	}
	if inspection := Inspect(root); inspection.State != Pending || inspection.ReusableGreen {
		t.Fatalf("after holder inspection = %#v, want pending and non-reusable", inspection)
	}
	if got := outcomeRuns(t, root); got != 1 {
		t.Fatalf("runs while external holder is active = %d, want 1", got)
	}
}

func TestGateRunTimeoutInvalidatesOldEvidence(t *testing.T) {
	root := failureOutcomeFixture(t)
	if result := Execute(context.Background(), root, &bytes.Buffer{}, &bytes.Buffer{}); result.ActionExit != 0 {
		t.Fatalf("green result = %#v", result)
	}
	previousTimeout := gateTimeout
	gateTimeout = 100 * time.Millisecond
	t.Cleanup(func() { gateTimeout = previousTimeout })
	outcomeWrite(t, root, ".gate-sleep", "\n", 0o644)
	var stdout, stderr bytes.Buffer
	if got := RunCommand([]string{"--fresh", root}, &stdout, &stderr); got != 124 {
		t.Fatalf("timeout exit = %d, stderr=%q", got, stderr.String())
	}
	if !strings.Contains(stderr.String(), "gate: timeout") {
		t.Fatalf("stderr = %q", stderr.String())
	}
	if inspection := Inspect(root); inspection.State != Ready || inspection.Status != "timeout" || inspection.ReusableGreen {
		t.Fatalf("timeout inspection = %#v", inspection)
	}
	if err := os.Remove(filepath.Join(root, ".gate-sleep")); err != nil {
		t.Fatal(err)
	}
	if result := Execute(context.Background(), root, &bytes.Buffer{}, &bytes.Buffer{}); result.ActionExit != 0 {
		t.Fatalf("ordinary run after timeout = %#v", result)
	}
	if got := outcomeRuns(t, root); got != 3 {
		t.Fatalf("runs after green, timeout, ordinary = %d, want 3", got)
	}
}

func TestGateRunCancellationLeavesPendingForRecovery(t *testing.T) {
	root := failureOutcomeFixture(t)
	outcomeWrite(t, root, ".gate-sleep", "\n", 0o644)
	ctx, cancel := context.WithCancel(context.Background())
	finished := make(chan Result, 1)
	go func() {
		finished <- Execute(ctx, root, &bytes.Buffer{}, &bytes.Buffer{})
	}()
	awaitGateMarker(t, filepath.Join(root, ".gate-record-during"))
	cancel()

	var result Result
	select {
	case result = <-finished:
	case <-time.After(failureOutcomeDeadline()):
		t.Fatal("canceled gate did not exit")
	}
	if result.GateExit != 130 || result.ActionExit != 130 {
		t.Fatalf("canceled result = %#v, want gate and action exit 130", result)
	}
	if inspection := Inspect(root); inspection.State != Pending || inspection.Status != "" || inspection.ReusableGreen {
		t.Fatalf("canceled inspection = %#v, want pending without a terminal status", inspection)
	}

	if err := os.Remove(filepath.Join(root, ".gate-sleep")); err != nil {
		t.Fatal(err)
	}
	if recovered := Execute(context.Background(), root, &bytes.Buffer{}, &bytes.Buffer{}); recovered.ActionExit != 0 {
		t.Fatalf("recovery result = %#v", recovered)
	}
	if inspection := Inspect(root); inspection.State != Ready || inspection.Status != "green" || !inspection.ReusableGreen {
		t.Fatalf("recovery inspection = %#v, want reusable ready green", inspection)
	}
}

func TestGateRunRefusesInitialPendingPersistenceWhenCacheIsDirectory(t *testing.T) {
	root := outcomeFixture(t)
	gitdir := outcomeGit(t, root, "rev-parse", "--absolute-git-dir")
	cache := filepath.Join(gitdir, "bench-last-gate")
	if err := os.Mkdir(cache, 0o700); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if got := RunCommand([]string{"--fresh", root}, &stdout, &stderr); got != 1 {
		t.Fatalf("pending persistence refusal exit = %d, stderr=%q", got, stderr.String())
	}
	if !strings.Contains(stderr.String(), "gate pending persistence failed") {
		t.Fatalf("stderr = %q", stderr.String())
	}
	if _, err := os.Stat(filepath.Join(root, ".gate-run-count")); !os.IsNotExist(err) {
		if err == nil {
			t.Fatal("oracle ran after initial pending persistence refusal")
		}
		t.Fatalf("run counter stat = %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".gate-record-during")); !os.IsNotExist(err) {
		if err == nil {
			t.Fatal("oracle wrote a record witness after initial pending persistence refusal")
		}
		t.Fatalf("record witness stat = %v", err)
	}
}

func TestGateRunRefusesOwnerPersistenceWhenOwnerPathIsDirectory(t *testing.T) {
	root := outcomeFixture(t)
	gitdir := outcomeGit(t, root, "rev-parse", "--absolute-git-dir")
	if err := os.Mkdir(filepath.Join(gitdir, "bench-gate-owner"), 0o700); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if got := RunCommand([]string{"--fresh", root}, &stdout, &stderr); got != 1 {
		t.Fatalf("owner persistence refusal exit = %d, stderr=%q", got, stderr.String())
	}
	if got := stderr.String(); got != "gate owner persistence failed\n" {
		t.Fatalf("stderr = %q, want exact owner persistence diagnostic", got)
	}
	for _, path := range []string{".gate-run-count", ".gate-record-during"} {
		if _, err := os.Stat(filepath.Join(root, path)); !os.IsNotExist(err) {
			if err == nil {
				t.Fatalf("oracle wrote %s after owner persistence refusal", path)
			}
			t.Fatalf("oracle witness stat for %s = %v", path, err)
		}
	}
}

func TestGateRunRejectsGreenWhenEvidenceDirectoryModeChanges(t *testing.T) {
	requireDirectoryWriteDenied(t)
	root := failureOutcomeFixture(t)
	if result := Execute(context.Background(), root, &bytes.Buffer{}, &bytes.Buffer{}); result.ActionExit != 0 {
		t.Fatalf("green result = %#v", result)
	}
	gitdir := outcomeGit(t, root, "rev-parse", "--absolute-git-dir")
	restoreMode(t, filepath.Join(gitdir, "bench-gate-evidence"))
	outcomeWrite(t, root, ".gate-evidence-0500", "\n", 0o644)
	var stdout, stderr bytes.Buffer
	if got := RunCommand([]string{"--fresh", root}, &stdout, &stderr); got != 1 {
		t.Fatalf("mode-change exit = %d, stderr=%q", got, stderr.String())
	}
	if !strings.Contains(stderr.String(), "gate evidence persistence failed") {
		t.Fatalf("stderr = %q", stderr.String())
	}
	if inspection := Inspect(root); inspection.ReusableGreen {
		t.Fatalf("inspection after failed retain = %#v", inspection)
	}
}

func TestGateRunRejectsRedWhenEvidenceInvalidationFails(t *testing.T) {
	requireDirectoryWriteDenied(t)
	root := failureOutcomeFixture(t)
	if result := Execute(context.Background(), root, &bytes.Buffer{}, &bytes.Buffer{}); result.ActionExit != 0 {
		t.Fatalf("green result = %#v", result)
	}
	gitdir := outcomeGit(t, root, "rev-parse", "--absolute-git-dir")
	restoreMode(t, filepath.Join(gitdir, "bench-gate-evidence"))
	outcomeWrite(t, root, ".gate-evidence-unwritable", "\n", 0o644)
	outcomeWrite(t, root, ".gate-red", "\n", 0o644)
	var stdout, stderr bytes.Buffer
	if got := RunCommand([]string{"--fresh", root}, &stdout, &stderr); got != 1 {
		t.Fatalf("invalidation exit = %d, stderr=%q", got, stderr.String())
	}
	if !strings.Contains(stderr.String(), "gate evidence invalidation failed") {
		t.Fatalf("stderr = %q", stderr.String())
	}
	if inspection := Inspect(root); inspection.ReusableGreen {
		t.Fatalf("inspection after failed invalidation = %#v", inspection)
	}
}

func TestGateRunPreservesPendingWhenTerminalReplaceFails(t *testing.T) {
	requireDirectoryWriteDenied(t)
	root := failureOutcomeFixture(t)
	if result := Execute(context.Background(), root, &bytes.Buffer{}, &bytes.Buffer{}); result.ActionExit != 0 {
		t.Fatalf("green result = %#v", result)
	}
	gitdir := outcomeGit(t, root, "rev-parse", "--absolute-git-dir")
	restoreMode(t, gitdir)
	outcomeWrite(t, root, ".gate-gitdir-0500", "\n", 0o644)
	var stdout, stderr bytes.Buffer
	if got := RunCommand([]string{"--fresh", root}, &stdout, &stderr); got != 1 {
		t.Fatalf("terminal replace exit = %d, stderr=%q", got, stderr.String())
	}
	if !strings.Contains(stderr.String(), "gate final persistence failed") {
		t.Fatalf("stderr = %q", stderr.String())
	}
	if inspection := Inspect(root); inspection.State != Pending || inspection.ReusableGreen {
		t.Fatalf("inspection after failed terminal replace = %#v, want pending and non-reusable", inspection)
	}
}

func failureOutcomeFixture(t *testing.T) string {
	t.Helper()
	return outcomeFixture(t, `if [ -e .gate-wait ]; then
  : > .gate-running
  while [ ! -e .gate-release ]; do sleep 0.01; done
fi
if [ -e .gate-sleep ]; then sleep 5; fi
if [ -e .gate-evidence-0500 ] || [ -e .gate-evidence-unwritable ]; then chmod 500 "$gitdir/bench-gate-evidence"; fi
if [ -e .gate-gitdir-0500 ]; then chmod 500 "$gitdir"; fi
`)
}

type gateLockHolder struct {
	releasePath string
	done        <-chan processGroupResult
}

func startGateLockHolder(t *testing.T, root string) gateLockHolder {
	t.Helper()
	readyPath := filepath.Join(root, ".gate-holder-ready")
	releasePath := filepath.Join(root, ".gate-holder-release")
	cmd := &exec.Cmd{
		Path: os.Args[0],
		Args: []string{os.Args[0], "-test.run=^TestGateRunExternalHolderDemotesReusableGreen$"},
		Env: append(os.Environ(),
			gateHolderRootEnv+"="+root,
			gateHolderReadyEnv+"="+readyPath,
			gateHolderReleaseEnv+"="+releasePath,
		),
	}
	done := make(chan processGroupResult, 1)
	go func() { done <- runProcessGroupCommand(context.Background(), cmd) }()
	awaitGateMarker(t, readyPath)
	holder := gateLockHolder{releasePath: releasePath, done: done}
	t.Cleanup(func() { holder.release(t) })
	return holder
}

func (holder gateLockHolder) release(t *testing.T) {
	t.Helper()
	if _, err := os.Stat(holder.releasePath); os.IsNotExist(err) {
		if err := os.WriteFile(holder.releasePath, nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	select {
	case result := <-holder.done:
		if result.Code != 0 || result.StartErr != nil || result.Cancelled {
			t.Fatalf("holder result = %#v", result)
		}
	case <-time.After(failureOutcomeDeadline()):
		t.Fatal("external lock holder did not exit")
	}
}

func runGateLockHolder(t *testing.T, root, readyPath, releasePath string) {
	t.Helper()
	gitdir := outcomeGit(t, root, "rev-parse", "--absolute-git-dir")
	lock, err := os.OpenFile(filepath.Join(gitdir, "bench-gate.lock"), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Close()
	flock := gocache.RecordLock(syscall.F_WRLCK)
	if err := syscall.FcntlFlock(lock.Fd(), syscall.F_SETLK, &flock); err != nil {
		t.Fatal(err)
	}
	defer syscall.FcntlFlock(lock.Fd(), syscall.F_SETLK, &syscall.Flock_t{Type: syscall.F_UNLCK, Whence: int16(0), Start: 0, Len: 0})
	ownerPath := filepath.Join(gitdir, "bench-gate-owner")
	if err := os.WriteFile(ownerPath, ownerRecord(time.Now()), 0o600); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(ownerPath)
	if err := os.WriteFile(readyPath, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(releasePath); err == nil {
			return
		} else if !os.IsNotExist(err) {
			t.Fatal(err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func awaitGateMarker(t *testing.T, path string) {
	t.Helper()
	deadline := time.Now().Add(failureOutcomeDeadline())
	for {
		if _, err := os.Stat(path); err == nil {
			return
		} else if !os.IsNotExist(err) {
			t.Fatal(err)
		}
		if time.Now().After(deadline) {
			t.Fatalf("gate marker %s was not written", path)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func awaitGateResult(t *testing.T, finished <-chan Result) {
	t.Helper()
	select {
	case result := <-finished:
		if result.ActionExit != 0 {
			t.Fatalf("in-flight result = %#v", result)
		}
	case <-time.After(failureOutcomeDeadline()):
		t.Fatal("in-flight gate did not exit")
	}
}

func restoreMode(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	mode := info.Mode()
	t.Cleanup(func() {
		if err := os.Chmod(path, mode); err != nil {
			t.Errorf("restore mode %s: %v", path, err)
		}
	})
}

func failureOutcomeDeadline() time.Duration { return bounds.TestDeadline(processGroupCancelGrace) }

func requireDirectoryWriteDenied(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	restoreMode(t, dir)
	if err := os.Chmod(dir, 0o500); err != nil {
		capability.Capability(t, capability.Privilege, fmt.Sprintf("cannot remove directory write permission: %v", err))
	}
	file, err := os.Create(filepath.Join(dir, "write-probe"))
	if err == nil {
		file.Close()
		capability.Capability(t, capability.Privilege, "mode 0500 directory remains writable")
	}
}

// cacheCleanProbeEnv names the file the clean probe writes its answer to. A child with no
// entry is inert, so the probe row is a no-op inside an ordinary suite run.
const cacheCleanProbeEnv = "BENCH_TEST_CACHE_CLEAN_PROBE"

// TestCacheCleanProbe is the second process the two holder rows drive. It runs
// `bench cache clean` and records the verb's own answer, because a POSIX record lock is
// owned per process: a clean inside the holder's process could never contend with it.
func TestCacheCleanProbe(t *testing.T) {
	answerPath := os.Getenv(cacheCleanProbeEnv)
	if answerPath == "" {
		return
	}
	answer, code := gocache.Command([]string{"clean"})
	if err := os.WriteFile(answerPath, []byte(strconv.Itoa(code)+"\n"+answer), 0o600); err != nil {
		t.Fatal(err)
	}
}

// probeArgv is the child invocation that runs the clean probe: this test binary with one
// row selected.
func probeArgv(t *testing.T) []string {
	t.Helper()
	binary, err := filepath.Abs(os.Args[0])
	if err != nil {
		t.Fatal(err)
	}
	return []string{binary, "-test.run=^TestCacheCleanProbe$"}
}

// requireRefusedClean reads the probe's answer and grades it as the refusal a live holder
// produces. An unheld lock lets the clean through, and that is what a missing hold looks
// like here.
func requireRefusedClean(t *testing.T, answerPath string) {
	t.Helper()
	answer := string(outcomeRead(t, answerPath))
	code, rest, _ := strings.Cut(answer, "\n")
	if code != "1" || !strings.HasPrefix(rest, "error: cache in use — ") {
		t.Fatalf("clean beside the run = exit %s, %q; want the cache-in-use refusal at exit 1", code, rest)
	}
}

// L01: a gate run holds the shared cache lock across its phases, so `bench cache clean`
// exits 1 while the oracle is compiling. The probe runs from the gate's own child, which
// is the one point inside the run's span a second process can observe.
func TestGateRunHoldsTheCacheLockAcrossItsPhases(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	answerPath := filepath.Join(t.TempDir(), "clean-answer")
	argv := probeArgv(t)
	root := outcomeFixture(t, "HOME="+home+" "+cacheCleanProbeEnv+"="+answerPath+
		" "+argv[0]+" "+argv[1]+" >/dev/null 2>&1\n")
	// The holder derives its directory from the closure's HOME, so the closure has to
	// declare that name. A closure without it locks nothing.
	outcomeWrite(t, root, ".bench/gate-inputs.json",
		`{"schema":1,"closure":"local","environment":["HOME"],"paths":[],"tools":[]}`+"\n", 0o644)
	outcomeCommit(t, root, "declare HOME")

	var stdout, stderr bytes.Buffer
	if result := Execute(context.Background(), root, &stdout, &stderr); result.ActionExit != 0 {
		t.Fatalf("gate result = %#v, stderr=%q", result, stderr.String())
	}
	requireRefusedClean(t, answerPath)
}

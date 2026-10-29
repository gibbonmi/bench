//go:build system

package systemtest

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/gate/prospectiveartifact"
	"github.com/gibbonmi/bench/internal/gocache"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
	"github.com/gibbonmi/bench/internal/testrepo"
)

// fakeCodex is a local codex provider. At first use, it records whether a user record names
// its parent CLI. Its mode can hold it or leave a descendant in its group.
const fakeCodex = `#!/bin/sh
state=absent
if grep -rqs "\"actor_pid\":$PPID," "$FAKE_CODEX_STORE"; then state=registered; fi
printf '%s\n' "$state" > "$FAKE_CODEX_MARK.partial" && mv "$FAKE_CODEX_MARK.partial" "$FAKE_CODEX_MARK"
case "$FAKE_CODEX_MODE" in
block)
	i=0
	while [ ! -e "$FAKE_CODEX_MARK.release" ] && [ "$i" -lt 600 ]; do sleep 0.05; i=$((i + 1)); done ;;
leave)
	sleep 300 </dev/null >/dev/null 2>&1 &
	echo "$!" > "$FAKE_CODEX_MARK.descendant.partial" && mv "$FAKE_CODEX_MARK.descendant.partial" "$FAKE_CODEX_MARK.descendant" ;;
esac
printf '{"models":[{"slug":"fake-model","visibility":"list"}]}\n'
`

// nestedFixture is an enclosing cache owner whose protection reaches a nested `bench models`.
type nestedFixture struct {
	repo, bin, marks string
	protocol         *protection.Protocol
	cache            protection.Descriptor
}

func newNestedFixture(t *testing.T) nestedFixture {
	t.Helper()
	holder, err := gocache.HoldDir(filepath.Join(t.TempDir(), "go-build"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = holder.Release() })
	f := nestedFixture{repo: journeyRepo(t, "process lifetime [journey]-"), bin: t.TempDir(), marks: t.TempDir()}
	f.protocol = &protection.Protocol{Validators: protection.Validators{protection.KindCache: gocache.ValidateProtection}, ReadLimit: bounds.ControlRecordLimit}
	if f.cache, err = holder.Protect(f.protocol); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(f.bin, "codex"), []byte(fakeCodex), 0o755); err != nil {
		t.Fatal(err)
	}
	return f
}

// env is one nested CLI's environment: the fake provider first on PATH, no credentials,
// discovery online, a private TMPDIR, and the enclosing owner's descriptor vector.
func (f nestedFixture) env(t *testing.T, name, mode string) []string {
	return []string{
		"PATH=" + f.bin + string(os.PathListSeparator) + os.Getenv("PATH"),
		"TMPDIR=" + t.TempDir(),
		"OPENAI_API_KEY", "ANTHROPIC_API_KEY", "BENCH_OFFLINE",
		protection.Environ(nil, []protection.Descriptor{f.cache})[0],
		"FAKE_CODEX_STORE=" + f.cache.Store,
		"FAKE_CODEX_MARK=" + f.mark(name),
		"FAKE_CODEX_MODE=" + mode,
	}
}

func (f nestedFixture) mark(name string) string { return filepath.Join(f.marks, name) }

// requireRegistered fails unless the named provider saw its own user record at first use.
func (f nestedFixture) requireRegistered(t *testing.T, name string) {
	t.Helper()
	if state, err := os.ReadFile(f.mark(name)); err != nil || string(state) != "registered\n" {
		t.Fatalf("provider %s used the cache with protection %q (%v), want registered", name, state, err)
	}
}

// observe seals the cache protection as its cleaner would, and reports whether release ran.
func (f nestedFixture) observe() (protection.Disposition, bool) {
	released := false
	disposition, _ := f.protocol.SealAndObserve(f.cache, func() error { released = true; return nil })
	return disposition, released
}

// TestProcessLifetimeNestedRegistration has a nested CLI register its provider's separate
// group before that provider runs.
func TestProcessLifetimeNestedRegistration(t *testing.T) {
	t.Parallel()
	f := newNestedFixture(t)
	result := systemSelected(t, f.repo, f.env(t, "only", "run"), "models")
	if result.code != 0 || !strings.Contains(result.stdout, "fake-model") {
		t.Fatalf("nested models = (%d, %q, %q)", result.code, result.stdout, result.stderr)
	}
	f.requireRegistered(t, "only")
	if disposition, released := f.observe(); disposition.State != protection.Safe || !released {
		t.Fatalf("the completed nested user left %+v released=%t", disposition, released)
	}
}

// TestProcessLifetimeNestedOverlap keeps independent records for overlapping nested users.
// The second CLI exits while its provider's group lives on, and the first CLI's completion
// keeps that retained user.
func TestProcessLifetimeNestedOverlap(t *testing.T) {
	t.Parallel()
	f := newNestedFixture(t)
	blocked, stdout, stderr := systemStartSelected(t, f.repo, f.env(t, "blocked", "block"), "models")
	waitForFile(t, f.mark("blocked"), "the blocked provider's use marker")
	left := systemSelected(t, f.repo, f.env(t, "left", "leave"), "models")
	descendant := processtest.WaitPID(t, f.mark("left")+".descendant")
	pgid, err := syscall.Getpgid(descendant)
	if left.code != 0 || err != nil {
		t.Fatalf("leaving models = (%d, %q, %q), descendant group %v", left.code, left.stdout, left.stderr, err)
	}
	if err := os.WriteFile(f.mark("blocked")+".release", nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if code := exitWithin(t, blocked, bounds.ProviderTimeout); code != 0 {
		t.Fatalf("blocked models = (%d, %q, %q)", code, stdout, stderr)
	}
	f.requireRegistered(t, "blocked")
	f.requireRegistered(t, "left")
	if disposition, released := f.observe(); disposition.State != protection.Protected || len(disposition.Users) != 1 || released {
		t.Fatalf("after the first CLI completed: %+v released=%t, want the retained second user", disposition, released)
	}
	_ = syscall.Kill(descendant, syscall.SIGKILL)
	processtest.WaitAbsent(t, pgid)
	if disposition, released := f.observe(); disposition.State != protection.Safe || !released {
		t.Fatalf("after the retained group's absence: %+v released=%t", disposition, released)
	}
}

// exitWithin waits for cmd to exit inside the deadline derived from inner.
func exitWithin(t *testing.T, cmd *exec.Cmd, inner time.Duration) int {
	t.Helper()
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	return exitedWithin(t, cmd, done, inner)
}

// exitedWithin reads the exit of cmd from done, which carries the result of its one Wait,
// inside the deadline derived from inner.
func exitedWithin(t *testing.T, cmd *exec.Cmd, done <-chan error, inner time.Duration) int {
	t.Helper()
	window := bounds.TestDeadline(inner)
	select {
	case err := <-done:
		return systemExitCode(err)
	case <-time.After(window):
		_ = cmd.Process.Kill()
		t.Fatal(bounds.TestTimeoutVerdict("the CLI to exit", window))
	}
	return -1
}

// TestProcessLifetimeNestedGate is PL32. The oracle starts the nested phase runner, kills it
// once its phase group is registered, and exits zero. The phase group survives as a bundle
// user, so the landing gate must retain no green evidence and keep the bundle.
func TestProcessLifetimeNestedGate(t *testing.T) {
	t.Parallel()
	root, home, tally, trees, ready, release := systemLandingRaceFixture(t)
	configureArtifactLandingFixture(t, root)
	f := testrepo.NewGateFixture(t.TempDir())
	f.Environment = []string{"HOME", "LAND_BASELINE_BENCH", "LAND_GATE_TALLY", "LAND_GATE_TREES", "LAND_RACE_READY", "LAND_RACE_RELEASE"}
	f.Command("go")
	// The hand-off keeps the spelling that the gate's phase-table predicate recognizes, so the
	// gate receives BENCH_KIT and BENCH_RUN_BINARY. The script reads both before the background
	// launch, so an absent value fails the gate and not only the nested runner.
	gate := "set -eu\nroot=${1:-$(" + f.Command("git") + " rev-parse --show-toplevel)}\nkit=${BENCH_KIT:?}\nbench=${BENCH_RUN_BINARY:?}\n" +
		f.Command("env") + " BENCH_KIT=\"$kit\" BENCH_RUN_BINARY=\"$bench\" \"$bench\" gate-phases \"$root\" &\nnested=$!\n" +
		boundedWait(f, artifactBarrierWindow) + f.Command("cat") + " \"$LAND_RACE_READY\" > /dev/null\nkill -KILL \"$nested\"\n"
	f.MustWrite(t, root, gate, gate)
	systemGit(t, root, "add", ".")
	systemGit(t, root, systemCommitArgs("-qm", "nested phase gate")...)
	systemGit(t, root, "update-ref", "refs/bench/green/main", systemGitOutput(t, root, "rev-parse", "HEAD"))
	evidence := filepath.Join(systemGitOutput(t, root, "rev-parse", "--path-format=absolute", "--git-common-dir"), "bench-gate-evidence")
	before := evidenceNames(t, evidence)
	private := t.TempDir()
	base := systemGitOutput(t, root, "rev-parse", "main")

	cmd, out, errOut, done := startArtifactAuthorization(t, root, home, tally, trees, ready, release, private, base, "nested-gate", "winner.txt", "nested phase subject")
	// The released phase lets the surviving group end after the row has its answer.
	t.Cleanup(func() { releaseArtifactBarrier(t, release) })
	// The landing contains the gate's own barrier wait.
	if code := exitedWithin(t, cmd, done, artifactBarrierWindow); code == 0 || !strings.Contains(out.String()+errOut.String(), "stays retained") {
		t.Fatalf("landing with a surviving nested phase = (%d, %q, %q), want the retained-bundle refusal", code, out.String(), errOut.String())
	}
	if after := evidenceNames(t, evidence); !reflect.DeepEqual(after, before) {
		t.Fatalf("gate evidence after the refusal = %v, want the prior %v", after, before)
	}
	bundle := oneProspectiveBundle(t, private)
	checkout := filepath.Join(bundle, prospectiveartifact.CheckoutName)
	if _, err := os.Stat(checkout); err != nil {
		t.Fatalf("checkout with a surviving phase user = %v, want retained", err)
	}
	requireCheckoutRegistered(t, root, checkout, true)
	// The surviving phase runs the selected executable, so the selection must keep it.
	if selected, err := filepath.Glob(filepath.Join(bundle, "bench-run-*", "bench")); err != nil || len(selected) != 1 {
		t.Fatalf("selected executable with a surviving phase user = %v, %v; want the one retained executable", selected, err)
	}
}

// evidenceNames lists the retained evidence records, and none when the store is absent.
func evidenceNames(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names
}

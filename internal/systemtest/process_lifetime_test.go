//go:build system

package systemtest

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/gocache"
	"github.com/gibbonmi/bench/internal/subprocess/processtest"
	"github.com/gibbonmi/bench/internal/subprocess/protection"
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
	f := nestedFixture{repo: landRouteRepo(t, "process lifetime [journey]-"), bin: t.TempDir(), marks: t.TempDir()}
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
	window := bounds.TestDeadline(inner)
	select {
	case err := <-done:
		return systemExitCode(err)
	case <-time.After(window):
		_ = cmd.Process.Kill()
		t.Fatal(bounds.TestTimeoutVerdict("the nested CLI to exit", window))
	}
	return -1
}

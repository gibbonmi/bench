package probe

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/gocache"
)

// gateHolderRootEnv and gateHolderReleaseEnv carry the fixture and the release marker to
// the lock-holding half of the gate row. A live gate holds its lock from another process,
// so the refusal is only proven by another process holding it here too.
const (
	gateHolderRootEnv    = "BENCH_PROBE_GATE_ROOT"
	gateHolderReleaseEnv = "BENCH_PROBE_GATE_RELEASE"
)

// probeArgs answers the accepted invocation when it is called bare, and the argv spelled
// in full otherwise. Every refusal row reads beside the accepted call it departs from, so
// the fault each row names is the only difference a reader has to hold.
func probeArgs(spelled ...string) []string {
	if len(spelled) != 0 {
		return spelled
	}
	return []string{"clamp.go", "--swap", "n < 0", "--with", "n > 0", "--package", "./", "--run", "^TestClampNegative$"}
}

// requireRefused asserts the one structured line, the exit, and the three facts every
// refusal owes: the subject unchanged, no copy under the home, and no run child started.
func requireRefused(t *testing.T, f *fixture, out string, code int, want string, wantCode int) {
	t.Helper()
	if out != want || code != wantCode {
		t.Fatalf("refusal = (%q, %d), want (%q, %d)", out, code, want, wantCode)
	}
	requireSubjectBytes(t, f, clampSource)
	requireHomeEmpty(t, f)
	requireNoRunChild(t, f)
}

// refusalFixture is a fixture whose stub `go` records any start, so PB25's third fact is
// observable on every refusal row.
func refusalFixture(t *testing.T) *fixture {
	t.Helper()
	f := newFixture(t)
	installStubGo(t, f, cannedPass, 0, "")
	return f
}

// PB18: a mutation that matches zero or several sites refuses, because replacing the
// first of several would mutate a site the caller never named.
func TestProbeRefusesAnAmbiguousMutation(t *testing.T) {
	for _, tc := range []struct{ name, old, want string }{
		{"several matches", "return", "error: probe mutation ambiguous — the old string matches 2 times, want exactly 1\n"},
		{"no match", "absent text", "error: probe mutation ambiguous — the old string matches 0 times, want exactly 1\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := refusalFixture(t)
			out, code := runProbe(t, probeArgs("clamp.go", "--swap", tc.old, "--with", "x", "--package", "./", "--run", "^TestClampNegative$")...)
			requireRefused(t, f, out, code, tc.want, 1)
		})
	}
}

// PB19: a replacement equal to the old string changes nothing, so the run would be
// vacuous and report silent over an unchanged file.
func TestProbeRefusesAnEmptyMutation(t *testing.T) {
	f := refusalFixture(t)
	out, code := runProbe(t, probeArgs("clamp.go", "--swap", "n < 0", "--with", "n < 0", "--package", "./", "--run", "^TestClampNegative$")...)
	requireRefused(t, f, out, code, "error: probe mutation empty — --with equals --swap\n", 1)
}

// PB20: the gate grades the bytes it read, so a probe refuses while a gate run holds the
// tree rather than changing the gate's subject under it.
func TestProbeRefusesUnderALiveGateRun(t *testing.T) {
	if root := os.Getenv(gateHolderRootEnv); root != "" {
		holdGateLock(t, root, os.Getenv(gateHolderReleaseEnv))
		return
	}
	f := refusalFixture(t)
	release := startGateHolder(t, f)
	defer release()
	out, code := runProbe(t, probeArgs()...)
	want := "error: gate execution in progress — wait for the gate run to finish before you mutate the tree\n"
	requireRefused(t, f, out, code, want, 1)
}

// startGateHolder runs the lock-holding half in its own process and answers the release.
func startGateHolder(t *testing.T, f *fixture) func() {
	t.Helper()
	releaseMarker := filepath.Join(scratchDir(t), "gate-release")
	holder := exec.Command(os.Args[0], "-test.run=^TestProbeRefusesUnderALiveGateRun$", "-test.timeout="+interruptDeadline().String())
	holder.Env = append(os.Environ(), gateHolderRootEnv+"="+f.root, gateHolderReleaseEnv+"="+releaseMarker)
	if err := holder.Start(); err != nil {
		t.Fatal(err)
	}
	awaitFile(t, gateLockPath(f.root))
	return func() {
		_ = os.WriteFile(releaseMarker, nil, 0o644)
		_ = holder.Wait()
	}
}

// gateLockPath spells the execution lock both halves of the gate row name. internal/gate
// does not export the path, so this file owns the one spelling the test side derives.
func gateLockPath(root string) string {
	return filepath.Join(root, ".git", "bench-gate.lock")
}

// holdGateLock takes the execution lock the gate takes and parks until the release marker
// appears, so the refusal under test observes a genuinely held lock.
func holdGateLock(t *testing.T, root, releaseMarker string) {
	t.Helper()
	lock, err := os.OpenFile(gateLockPath(root), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Close()
	held := gocache.RecordLock(syscall.F_WRLCK)
	if err := syscall.FcntlFlock(lock.Fd(), syscall.F_SETLK, &held); err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(releaseMarker); err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// awaitFile waits for the holder's lock file to appear and then settles, because the
// create and the lock are two steps in the holder's own process.
func awaitFile(t *testing.T, path string) {
	t.Helper()
	expiry := time.Now().Add(interruptDeadline())
	for time.Now().Before(expiry) {
		if info, err := os.Stat(path); err == nil && info.Mode().IsRegular() {
			time.Sleep(50 * time.Millisecond)
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("the gate lock holder did not create %s within %s", path, interruptDeadline())
}

// PB21: only a regular file inside the repository mutates, so each hostile subject names
// what it is instead of the verb following or widening it.
func TestProbeRefusesAHostileSubject(t *testing.T) {
	for _, tc := range []struct{ name, reason string }{
		{"absent.go", "absent"},
		{"cmd", "a directory"},
		{"link.go", "a symlink"},
		{"pipe", "a special file"},
	} {
		t.Run(tc.reason, func(t *testing.T) {
			f := refusalFixture(t)
			switch tc.name {
			case "link.go":
				if err := os.Symlink(f.subject, filepath.Join(f.root, "link.go")); err != nil {
					t.Fatal(err)
				}
			case "pipe":
				if err := syscall.Mkfifo(filepath.Join(f.root, "pipe"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			out, code := runProbe(t, probeArgs(tc.name, "--omit", "n < 0", "--package", "./")...)
			want := "error: probe subject unavailable — " + tc.name + " is " + tc.reason + "\n"
			requireRefused(t, f, out, code, want, 1)
		})
	}

	t.Run("outside the repository", func(t *testing.T) {
		f := refusalFixture(t)
		outside := filepath.Join(scratchDir(t), "outside.go")
		writeFixtureFile(t, outside, "package outside\n", 0o644)
		out, code := runProbe(t, probeArgs(outside, "--omit", "package", "--package", "./")...)
		if !strings.HasSuffix(out, " is outside the repository\n") || code != 1 {
			t.Fatalf("refusal = (%q, %d), want the outside-the-repository line and 1", out, code)
		}
		requireSubjectBytes(t, f, clampSource)
		requireHomeEmpty(t, f)
		requireNoRunChild(t, f)
	})
}

// PB22: the grammar names exactly one mutation form and exactly one selection form, so a
// missing or doubled form is a usage error rather than a defaulted run.
func TestProbeUsageNamesOneSelectionAndOneMutation(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
	}{
		{"no selection", []string{"clamp.go", "--omit", "n < 0"}},
		{"two selections", []string{"clamp.go", "--omit", "n < 0", "--package", "./", "--check", "line-routing"}},
		{"run without package", []string{"clamp.go", "--omit", "n < 0", "--check", "line-routing", "--run", "^X$"}},
		{"two mutations", []string{"clamp.go", "--omit", "n < 0", "--swap", "a", "--with", "b", "--package", "./"}},
		{"no mutation", []string{"clamp.go", "--package", "./"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := refusalFixture(t)
			out, code := runProbe(t, tc.args...)
			requireRefused(t, f, out, code, grammar.Help+"\n", 2)
		})
	}
}

// PB23: the selection is parsed before any write, so an unknown check answers `bench
// test`'s own refusal and leaves no copy behind.
func TestProbeRefusesAnUnknownCheckBeforeAnyWrite(t *testing.T) {
	f := refusalFixture(t)
	out, code := runProbe(t, probeArgs("clamp.go", "--omit", "n < 0", "--check", "no-such-check")...)
	if code != 2 || !strings.HasPrefix(out, "unknown check: no-such-check\n") || !strings.Contains(out, "checks:\n") {
		t.Fatalf("refusal = (%q, %d), want the unknown-check inventory and 2", out, code)
	}
	requireSubjectBytes(t, f, clampSource)
	requireHomeEmpty(t, f)
	requireNoRunChild(t, f)
}

// PB24: neither the prose grade nor the system suite reaches the report a verdict derives
// from, so both refuse as probe targets.
func TestProbeRefusesProseAndSystemChecks(t *testing.T) {
	for _, check := range []string{"prose", "system"} {
		t.Run(check, func(t *testing.T) {
			f := refusalFixture(t)
			out, code := runProbe(t, probeArgs("clamp.go", "--omit", "n < 0", "--check", check)...)
			want := "error: probe focused run unsupported — --check prose and --check system are not probe targets\n"
			requireRefused(t, f, out, code, want, 1)
		})
	}
}

// PB26: the verb answers the shared not-in-repo error outside a repository, so it matches
// every other Bench command there.
func TestProbeOutsideARepository(t *testing.T) {
	t.Chdir(scratchDir(t))
	out, code := runProbe(t, probeArgs()...)
	want := "error: not in a git repository — run inside a Bench-linked repo\n"
	if out != want || code != 1 {
		t.Fatalf("refusal = (%q, %d), want (%q, 1)", out, code, want)
	}
}

// PB44: preservation runs before the mutation, so a home that cannot hold the copy leaves
// the subject exactly as it was.
func TestProbeRefusesWhenPreservationFails(t *testing.T) {
	f := refusalFixture(t)
	writeFixtureFile(t, filepath.Join(f.home, "probe"), "not a directory\n", 0o644)
	out, code := runProbe(t, probeArgs()...)
	if code != 1 || !strings.HasPrefix(out, "error: probe preservation failed — ") {
		t.Fatalf("refusal = (%q, %d), want the preservation failure and 1", out, code)
	}
	requireSubjectBytes(t, f, clampSource)
	requireNoRunChild(t, f)
}

// PB45: a mutation write that fails leaves the subject unchanged, removes the copy, and
// starts no run, so a probe never runs over a file it did not mutate.
func TestProbeRefusesWhenTheMutationWriteFails(t *testing.T) {
	f := refusalFixture(t)
	if err := os.Chmod(f.root, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(f.root, 0o755) })
	out, code := runProbe(t, probeArgs()...)
	if code != 1 || !strings.HasPrefix(out, "error: probe mutation failed — ") {
		t.Fatalf("refusal = (%q, %d), want the mutation failure and 1", out, code)
	}
	requireSubjectBytes(t, f, clampSource)
	requireHomeEmpty(t, f)
	requireNoRunChild(t, f)
}

// PB30: every help spelling answers the grammar's usage line on stdout with exit 0.
func TestProbeHelpSpellings(t *testing.T) {
	for _, spelling := range []string{"--help", "-h", "help"} {
		t.Run(spelling, func(t *testing.T) {
			out, code := Command([]string{spelling})
			if code != 0 || out != grammar.Help+"\n" {
				t.Fatalf("%s = (%q, %d), want the usage line and 0", spelling, out, code)
			}
			if !strings.HasPrefix(out, "usage: bench probe <file> (--swap") {
				t.Fatalf("usage line = %q", out)
			}
		})
	}
}

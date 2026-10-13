package probe

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// gateRefusal is the refusal a probe answers while a gate run holds the tree.
const gateRefusal = "error: gate execution in progress — wait for the gate run to finish before you mutate the tree\n"

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
	miss := probeRow(t, "invalid", "clamp.go", "swap", "substring-miss", 0, nil, "untouched") +
		selectionRow(t, "package", "./", "^TestClampNegative$", "", 0) + executionRow(t, 1, nil, 0)
	for _, tc := range []struct{ name, old, want string }{
		{"no match", "absent text", miss},
		{"several matches", "return", "error: probe mutation ambiguous — the old string matches 2 times, want exactly 1\n"},
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
// tree rather than changing the gate's subject under it. The test is also the lock holder's
// entry: a process that startLockHolder starts holds the lock its environment names.
func TestProbeRefusesUnderALiveGateRun(t *testing.T) {
	if path := os.Getenv(lockHolderPathEnv); path != "" {
		holdRecordLock(t, path, os.Getenv(lockHolderReleaseEnv))
		return
	}
	f := refusalFixture(t)
	release := startLockHolder(t, gateLockPath(f.root))
	defer release()
	out, code := runProbe(t, probeArgs()...)
	requireRefused(t, f, out, code, gateRefusal, 1)
}

// gateLockPath spells the execution lock both halves of the gate row name. internal/gate
// does not export the path, so this file owns the one spelling the test side derives.
func gateLockPath(root string) string {
	return filepath.Join(root, ".git", "bench-gate.lock")
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
		{"run with a check other than system", []string{"clamp.go", "--omit", "n < 0", "--check", "line-routing", "--run", "^X$"}},
		{"delete with omit", []string{"clamp.go", "--delete", "--omit", "n < 0", "--package", "./"}},
		{"delete with unwrap", []string{"clamp.go", "--delete", "--unwrap", "ensureNonNegative(n)", "--package", "./"}},
		{"delete with omit-file", []string{"clamp.go", "--delete", "--omit-file", "x", "--package", "./"}},
		{"insert-before with swap", []string{"clamp.go", "--insert-before", "a", "--with", "b", "--swap", "c", "--package", "./"}},
		{"insert-before with omit", []string{"clamp.go", "--insert-before", "a", "--with", "b", "--omit", "n < 0", "--package", "./"}},
		{"insert-before with delete", []string{"clamp.go", "--insert-before", "a", "--with", "b", "--delete", "--package", "./"}},
		{"omit-file with swap", []string{"clamp.go", "--omit-file", "x", "--swap", "a", "--with", "b", "--package", "./"}},
		{"omit-file with omit", []string{"clamp.go", "--omit-file", "x", "--omit", "n < 0", "--package", "./"}},
		{"omit-file with unwrap", []string{"clamp.go", "--omit-file", "x", "--unwrap", "ensureNonNegative(n)", "--package", "./"}},
		{"with only swap", []string{"clamp.go", "--with", "x", "--package", "./"}},
		{"omit-file with replacement", []string{"clamp.go", "--omit-file", "x", "--with", "n", "--package", "./"}},
		{"unwrap with replacement", []string{"clamp.go", "--unwrap", "ensureNonNegative(n)", "--with", "n", "--package", "./"}},
		{"omit with replacement", []string{"clamp.go", "--omit", "n < 0", "--with", "n", "--package", "./"}},
		{"no mutation", []string{"clamp.go", "--package", "./"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := refusalFixture(t)
			out, code := runProbe(t, tc.args...)
			requireRefused(t, f, out, code, grammar.Help+"\n", 2)
		})
	}
}

// A later --package with an empty value refuses as the empty flag value, because an empty
// expression names no package.
func TestProbeRefusesAnEmptyLaterPackage(t *testing.T) {
	f := refusalFixture(t)
	out, code := runProbe(t, probeArgs("clamp.go", "--swap", "n < 0", "--with", "n > 0", "--package", "./", "--package", "", "--run", "^TestClampNegative$")...)
	requireRefused(t, f, out, code, toon.Usage(grammar.Cmd, usage.EmptyFlagValue("--package"))+"\n", 2)
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
// the subject exactly as it was. The baseline is the one run this refusal allows, because
// it runs before preservation; the refusal itself starts no mutated run.
func TestProbeRefusesWhenPreservationFails(t *testing.T) {
	f := refusalFixture(t)
	writeFixtureFile(t, filepath.Join(f.home, "probe"), "not a directory\n", 0o644)
	out, code := runProbe(t, probeArgs()...)
	if code != 1 || !strings.HasPrefix(out, "error: probe preservation failed — ") {
		t.Fatalf("refusal = (%q, %d), want the preservation failure and 1", out, code)
	}
	requireSubjectBytes(t, f, clampSource)
	requireBaselineOnly(t, f)
}

// PB45: a mutation write that fails leaves the subject unchanged and removes the copy, so a
// probe never runs over a file it did not mutate. The baseline is the one run that already
// happened, and the refusal starts no mutated run after it.
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
	requireBaselineOnly(t, f)
}

// probeHelpNotes is the two facts the help owes, spelled independently of the producer.
// A dropped, reworded, or extra note reds this row rather than passing as a fresh
// expectation.
var probeHelpNotes = []string{
	"--package <expr> takes one Go package expression as one argument, so a space-separated list is not split; repeat --package to run several expressions in one call, or use a parent ./... pattern.",
	"--check <name> names a check from the bench test --help inventory.",
}

// PB30, DG4, PS35: every help spelling answers the usage line and then the notes block with
// exit 0, and the notes name no copy-aside route and no not-a-target sentence.
func TestProbeHelpSpellings(t *testing.T) {
	for _, spelling := range []string{"--help", "-h", "help"} {
		t.Run(spelling, func(t *testing.T) {
			out, code := Command([]string{spelling})
			if !strings.HasPrefix(out, "usage: bench probe <file> (--swap") {
				t.Fatalf("usage line = %q", out)
			}
			want := grammar.Help + "\nnotes:\n  " + strings.Join(probeHelpNotes, "\n  ") + "\n"
			if code != 0 || out != want {
				t.Fatalf("%s = (%q, %d), want (%q, 0)", spelling, out, code, want)
			}
		})
	}
}

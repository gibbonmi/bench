package runtime

import (
	"os"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/contract"
	"github.com/gibbonmi/bench/internal/status"
)

// handoffPhaseFixture parks an idea so the board's leading action is a phase invocation.
// Without one there is nothing for the harness table to translate and the paired-absence
// rows below would pass vacuously.
func handoffPhaseFixture(t *testing.T) contract.Fixture {
	t.Helper()
	f := handoffFixtureOnMain(t)
	f.WriteFile("IDEAS.md", "- 2026-07-05  parked idea\n")
	f.CommitAll("parked")
	return f
}

func TestHandoffHarnessCodex(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	f := handoffPhaseFixture(t)

	out := f.Bench("handoff", "--harness", "codex")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "$bench-what-next")
	contract.RequireNotContains(t, out.Stdout, "/bench-")
	contract.RequireNotContains(t, f.ReadFile(status.HandoffFile), "/bench-")
}

func TestHandoffHarnessDefault(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	for _, args := range [][]string{{"handoff"}, {"handoff", "--harness", "claude"}} {
		f := handoffPhaseFixture(t)
		out := f.Bench(args...)
		out.RequireExit(0)
		contract.RequireContains(t, out.Stdout, "/bench-what-next")
		contract.RequireNotContains(t, out.Stdout, "$bench-")
		contract.RequireNotContains(t, f.ReadFile(status.HandoffFile), "$bench-")
	}
}

func TestHandoffNextOverride(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	f := handoffPhaseFixture(t)

	out := f.Bench("handoff", "--next", "bench debug the canary hang")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "`bench debug the canary hang`")
	// The derived action must be replaced, not accompanied. The signal clause is the tell:
	// an override that appended rather than replaced would still name the board row it came
	// from, so its absence is what distinguishes the two. The needle is the renderer's own
	// wording — an approximation of it can never fail, which is what this assertion did
	// while reading "the board's leading signal" against a renderer that emits "invocable".
	contract.RequireNotContains(t, out.Stdout, "/bench-what-next")
	contract.RequireNotContains(t, out.Stdout, "the board's leading invocable signal")
	contract.RequireNotContains(t, f.ReadFile(status.HandoffFile), "/bench-what-next")
}

func TestHandoffRejectsUnknownHarness(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	f := handoffPhaseFixture(t)

	out := f.Bench("handoff", "--harness", "gemini")
	out.RequireExit(2)
	contract.RequireContains(t, out.Stdout, "usage: bench handoff")
	contract.RequireContains(t, out.Stdout, "--harness gemini")
	if f.Exists(status.HandoffFile) {
		t.Fatal("a usage error wrote " + status.HandoffFile)
	}
}

func TestHandoffArgGrammar(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	cases := []struct {
		args []string
		want string
	}{
		{[]string{"handoff", "--harness", "claude", "--harness", "claude"}, "unknown argument: --harness"},
		{[]string{"handoff", "--next", "a", "--next", "b"}, "unknown argument: --next"},
		{[]string{"handoff", "--write"}, "unknown argument: --write"},
		{[]string{"handoff", "-x"}, "unknown argument: -x"},
		{[]string{"handoff", "extra"}, "unknown argument: extra"},
		{[]string{"handoff", "--harness"}, "missing argument: --harness"},
		// A value-taking flag given nothing is a mistyped invocation, not a request for
		// the derived line: the shared grammar cannot see it, so the command names it.
		{[]string{"handoff", "--next", ""}, `unknown argument: --next ""`},
	}
	for _, tc := range cases {
		f := handoffFixtureOnMain(t)
		f.WriteFile("tracked.txt", "base\n")
		f.CommitAll("base")
		out := f.Bench(tc.args...)
		out.RequireExit(2)
		contract.RequireContains(t, out.Stdout, "usage: bench handoff")
		contract.RequireContains(t, out.Stdout, tc.want)
		if f.Exists(status.HandoffFile) {
			t.Fatalf("args %q wrote %s", tc.args, status.HandoffFile)
		}
	}
}

func TestHandoffDegenerateGit(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()

	detached := handoffFixtureOnMain(t)
	detached.WriteFile("tracked.txt", "base\n")
	detached.CommitAll("base")
	detached.Git("checkout", "-q", "--detach")
	out := detached.Bench("handoff")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "Branch: detached HEAD (no branch)")
	contract.RequireContains(t, out.Stdout, "HEAD `"+handoffShortSHA(t, detached)+"`")

	// No commits: HEAD names nothing, and the landed-state query cannot resolve a default
	// branch either, so three separate values degrade rather than the command failing.
	empty := handoffFixtureOnMain(t)
	out = empty.Bench("handoff")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "HEAD unknown (no commits yet)")
	contract.RequireContains(t, out.Stdout, "dirty state unknown")
	contract.RequireContains(t, out.Stdout, "unpushed count unknown")

	noOrigin := handoffFixtureOnMain(t)
	noOrigin.WriteFile("tracked.txt", "base\n")
	noOrigin.CommitAll("base")
	out = noOrigin.Bench("handoff")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "(origin unknown)")

	// An absent fact reads as stated, never as an empty inline-code span.
	for _, o := range []string{detached.Bench("handoff").Stdout, empty.Bench("handoff").Stdout, noOrigin.Bench("handoff").Stdout} {
		contract.RequireNotContains(t, o, "``")
	}
}

func TestHandoffRefusesControlBytes(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	f := handoffFixtureOnMain(t)
	f.WriteFile("tracked.txt", "base\n")
	f.CommitAll("base")
	// git rejects a control byte in a ref name, so the reachable carrier is a spec's
	// Status value, which reaches the rendered block through the same field validation.
	f.WriteFile("specs/hostile.md", "# Hostile\n\nStatus: sta\x1bged\n")

	out := f.Bench("handoff")
	out.RequireExit(1)
	contract.RequireContains(t, out.Stdout, "error: ")
	contract.RequireNotContains(t, out.Stdout, "\x1b")
	if f.Exists(status.HandoffFile) {
		t.Fatal("a refused render still wrote " + status.HandoffFile)
	}
}

// TestHandoffRefusesNewlineFields covers the control bytes the TOON predicate permits.
// TOON escapes tab, newline, and return; a line-structured markdown document cannot, so a
// value carrying one splits its own field. The `--next` carrier is the sharp one: a newline
// there can write a second `## State` heading, after which every later run refuses the
// document as ambiguous and the reviewer's real State is unreachable to the command.
func TestHandoffRefusesNewlineFields(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	f := handoffFixtureOnMain(t)
	f.WriteFile("tracked.txt", "base\n")
	f.CommitAll("base")

	out := f.Bench("handoff", "--next", "run this\n\n## State\n\nhijacked")
	out.RequireExit(1)
	contract.RequireContains(t, out.Stdout, "error: ")
	if f.Exists(status.HandoffFile) {
		t.Fatal("a refused render still wrote " + status.HandoffFile)
	}

	// A spec Status carrying a tab reaches the same guard from the derived side.
	f.WriteFile("specs/hostile.md", "# Hostile\n\nStatus: sta\tged\n")
	tabbed := f.Bench("handoff")
	tabbed.RequireExit(1)
	contract.RequireContains(t, tabbed.Stdout, "error: ")
}

func TestHandoffUnwritableTarget(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	if os.Geteuid() == 0 {
		capability.Capability(t, capability.Privilege, "root bypasses directory permissions; cannot simulate an unwritable target")
	}
	f := handoffFixtureOnMain(t)
	f.WriteFile("tracked.txt", "base\n")
	f.CommitAll("base")
	if err := os.Chmod(f.Root, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(f.Root, 0o755) })

	out := f.Bench("handoff")
	out.RequireExit(1)
	// A structured line naming the path, not the bare permission trace git or Go would
	// otherwise surface.
	contract.RequireContains(t, out.Stdout, "error: ")
	contract.RequireContains(t, out.Stdout, status.HandoffFile)
	if !strings.HasPrefix(out.Stdout, "error: ") {
		t.Fatalf("write failure did not take the structured error shape:\n%s", out.Stdout)
	}
}

func TestHandoffNotInRepo(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	f := contract.NewFixture(t, contract.WithNoRepo())

	out := f.Bench("handoff")
	out.RequireExit(1)
	contract.RequireContains(t, out.Stdout, "error: not in a git repository — run inside a Bench-linked repo")
	if f.Exists(status.HandoffFile) {
		t.Fatal("the not-in-a-repo path wrote " + status.HandoffFile)
	}
}

func TestHandoffRouted(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Parallel()
	f := handoffFixtureOnMain(t)
	f.WriteFile("tracked.txt", "base\n")
	f.CommitAll("base")

	// Through the wrapper, not the binary: the porcelain route and its advertisement in
	// the wrapper's own help are what make the command discoverable at the shell.
	out := f.Run("bash", benchPath(t), "handoff")
	out.RequireExit(0)
	contract.RequireContains(t, out.Stdout, "# Session handoff")
	help := f.Run("bash", benchPath(t), "help")
	help.RequireExit(0)
	contract.RequireContains(t, help.Stdout, "bench handoff")
}

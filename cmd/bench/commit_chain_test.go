package main

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/commit"
	"github.com/gibbonmi/bench/internal/worktree"
)

const (
	chainSlug      = "chain-fixture"
	chainWorktree  = "/fixture/worktree"
	chainPublished = "0123456789abcdef0123456789abcdef01234567"
)

// chainCalls records the root, the home, and the argv each injected build received, and
// the argv each injected preflight received.
type chainCalls struct {
	builds     [][]string
	preflights [][]string
}

// runCommitChain replaces the three chain steps and runs `bench commit` through the
// production dispatcher, with --preflight-build naming slug when slug is not empty. The
// commit step prints one line and answers slug, published, and commitExit; the build and
// preflight steps print one line each and answer their own exits.
func runCommitChain(t *testing.T, slug, published string, commitExit, buildExit, preflightExit int) (stdout string, code int, calls *chainCalls) {
	t.Helper()
	t.Setenv(benchhome.Env, t.TempDir())
	old := commitChain
	t.Cleanup(func() { commitChain = old })
	calls = &chainCalls{}
	commitChain = chainSteps{
		Commit: func(_ []string, stdout, _ io.Writer) (commit.Outcome, int) {
			fmt.Fprintln(stdout, "commit step")
			return commit.Outcome{Root: chainWorktree, Published: published, PreflightBuild: slug}, commitExit
		},
		Build: func(root, home string, args []string, stdout, _ io.Writer) int {
			calls.builds = append(calls.builds, append([]string{root, home}, args...))
			fmt.Fprintln(stdout, "build step")
			return buildExit
		},
		Preflight: func(args []string) (string, int) {
			calls.preflights = append(calls.preflights, args)
			return "preflight step\n", preflightExit
		},
	}
	args := []string{"commit", "-m", "m"}
	if slug != "" {
		args = append(args, commit.PreflightBuildFlag, slug)
	}
	var out, errOut bytes.Buffer
	code = Command{Stdout: &out, Stderr: &errOut, Executable: "bench"}.Run(append(args, "--", "a.txt"))
	return out.String(), code, calls
}

// wantChainLine is the public spelling of the chain line, written apart from its renderer,
// so a renamed field or a changed step state reds every chain row.
func wantChainLine(sha, build, preflight string) string {
	return "commit-chain{commit=" + sha + ",build=" + build + ",preflight=" + preflight + "}"
}

func lastLine(stdout string) string {
	lines := strings.Split(strings.TrimSuffix(stdout, "\n"), "\n")
	return lines[len(lines)-1]
}

// BO51: a green commit, build, and preflight run in order, each prints its own response,
// and the chain line names the published commit at exit 0.
func TestCommitChainRunsThreeSteps(t *testing.T) {
	stdout, code, calls := runCommitChain(t, chainSlug, chainPublished, 0, 0, 0)
	want := "commit step\nbuild step\npreflight step\n" + wantChainLine(chainPublished, "green", "green") + "\n"
	if code != 0 || stdout != want {
		t.Fatalf("chain = (%d, %q), want (0, %q)", code, stdout, want)
	}
	if want := [][]string{{chainWorktree, worktree.Home(), chainWorktree}}; !reflect.DeepEqual(calls.builds, want) {
		t.Fatalf("builds = %q, want %q: the commit's own worktree as root and target under the Bench home", calls.builds, want)
	}
	if want := [][]string{{"build", chainSlug, "--source-tip", chainPublished}}; !reflect.DeepEqual(calls.preflights, want) {
		t.Fatalf("preflight argv = %q, want %q", calls.preflights, want)
	}
}

// BO52: a lane refusal publishes nothing, so the chain calls no build and exits 1.
func TestCommitChainSkipsAfterRefusal(t *testing.T) {
	stdout, code, calls := runCommitChain(t, chainSlug, "", 1, 0, 0)
	if want := wantChainLine("none", "skipped", "skipped"); code != 1 || lastLine(stdout) != want {
		t.Fatalf("chain = (%d, %q), want exit 1 ending in %q", code, stdout, want)
	}
	if len(calls.builds) != 0 || len(calls.preflights) != 0 {
		t.Fatalf("a refused commit ran builds %q and preflights %q", calls.builds, calls.preflights)
	}
}

// BO53: a red build returns its own exit, calls no preflight, and still names the commit.
func TestCommitChainStopsAtRedBuild(t *testing.T) {
	stdout, code, calls := runCommitChain(t, chainSlug, chainPublished, 0, 130, 0)
	if want := wantChainLine(chainPublished, "red", "skipped"); code != 130 || lastLine(stdout) != want {
		t.Fatalf("chain = (%d, %q), want exit 130 ending in %q", code, stdout, want)
	}
	if len(calls.preflights) != 0 {
		t.Fatalf("a red build ran preflights %q", calls.preflights)
	}
}

// BO54: a red preflight after a green build is the chain's exit.
func TestCommitChainReportsRedPreflight(t *testing.T) {
	stdout, code, _ := runCommitChain(t, chainSlug, chainPublished, 0, 0, 1)
	if want := wantChainLine(chainPublished, "green", "red"); code != 1 || lastLine(stdout) != want {
		t.Fatalf("chain = (%d, %q), want exit 1 ending in %q", code, stdout, want)
	}
}

// BO71: a commit that published without reconciling its checkout stops the chain at exit
// 3, and the chain line names the published commit.
func TestCommitChainStopsAtRemainder(t *testing.T) {
	stdout, code, calls := runCommitChain(t, chainSlug, chainPublished, 3, 0, 0)
	if want := wantChainLine(chainPublished, "skipped", "skipped"); code != 3 || lastLine(stdout) != want {
		t.Fatalf("chain = (%d, %q), want exit 3 ending in %q", code, stdout, want)
	}
	if len(calls.builds) != 0 || len(calls.preflights) != 0 {
		t.Fatalf("an unreconciled commit ran builds %q and preflights %q", calls.builds, calls.preflights)
	}
}

// BO56: a commit without --preflight-build runs the commit step alone through the
// dispatcher: no build, no preflight, no chain line, and the commit step's exit.
func TestCommitWithoutPreflightBuildRunsCommitAlone(t *testing.T) {
	stdout, code, calls := runCommitChain(t, "", chainPublished, 0, 130, 1)
	if code != 0 || stdout != "commit step\n" {
		t.Fatalf("commit = (%d, %q), want (0, %q)", code, stdout, "commit step\n")
	}
	if len(calls.builds) != 0 || len(calls.preflights) != 0 {
		t.Fatalf("a commit without the flag ran builds %q and preflights %q", calls.builds, calls.preflights)
	}
}

// A help request through the production binding answers the commit's help, then the
// chain's help line, at exit 0.
func TestCommitHelpEndsWithChainHelp(t *testing.T) {
	t.Setenv(benchhome.Env, t.TempDir())
	var out, errOut bytes.Buffer
	code := Command{Stdout: &out, Stderr: &errOut, Executable: "bench"}.Run([]string{"commit", "--help"})
	if code != 0 || !strings.HasPrefix(out.String(), "usage: bench commit") || !strings.HasSuffix(out.String(), "\n"+chainHelp+"\n") {
		t.Fatalf("help = (%d, %q), want the commit help ending in the chain help line", code, out.String())
	}
}

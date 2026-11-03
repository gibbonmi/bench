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
	chainCandidate = "fedcba9876543210fedcba9876543210fedcba98"
	chainPublished = "0123456789abcdef0123456789abcdef01234567"
)

// chainCalls records the root, the home, and the argv each injected build received, and
// the root, the slug, the base, and the candidate each injected preflight received.
type chainCalls struct {
	builds     [][]string
	preflights [][]string
}

// runCommitChain replaces the three chain steps and runs `bench commit` through the
// production dispatcher, with --preflight-build naming slug when slug is not empty and
// --base naming base when base is not empty. The commit step stands in for RunChecked: a
// commitExit of 1 is a refusal before the check, so it calls no check; otherwise it calls a
// supplied check with chainWorktree and chainCandidate, and a check error refuses at exit 1
// with nothing published. It then prints one line and answers published and commitExit. The
// build and preflight steps print one line each and answer their own exits.
func runCommitChain(t *testing.T, slug, base, published string, commitExit, buildExit, preflightExit int) (stdout string, code int, calls *chainCalls) {
	t.Helper()
	t.Setenv(benchhome.Env, t.TempDir())
	old := commitChain
	t.Cleanup(func() { commitChain = old })
	calls = &chainCalls{}
	commitChain = chainSteps{
		Commit: func(_ []string, stdout, _ io.Writer, check func(root, candidate string) error) (commit.Outcome, int) {
			outcome := commit.Outcome{Root: chainWorktree, PreflightBuild: slug, Base: base}
			if commitExit != 1 && check != nil && check(chainWorktree, chainCandidate) != nil {
				commitExit = 1
			}
			fmt.Fprintln(stdout, "commit step")
			if commitExit != 1 {
				outcome.Published = published
			}
			return outcome, commitExit
		},
		Build: func(root, home string, args []string, stdout, _ io.Writer) int {
			calls.builds = append(calls.builds, append([]string{root, home}, args...))
			fmt.Fprintln(stdout, "build step")
			return buildExit
		},
		Preflight: func(root, slug, base, candidate string) (string, int) {
			calls.preflights = append(calls.preflights, []string{root, slug, base, candidate})
			return "preflight step\n", preflightExit
		},
	}
	args := []string{"commit", "-m", "m"}
	if slug != "" {
		args = append(args, commit.PreflightBuildFlag, slug)
	}
	if base != "" {
		args = append(args, commit.BaseFlag, base)
	}
	var out, errOut bytes.Buffer
	code = Command{Stdout: &out, Stderr: &errOut, Executable: "bench"}.Run(append(args, "--", "a.txt"))
	return withoutTreeRow(out.String()), code, calls
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

// BO51, PP17, PP18, PP23: before publication, the chain runs the build and then the
// preflight at the candidate commit, each prints its own response, then the commit publishes
// and the chain line names the published commit at exit 0.
func TestCommitChainRunsThreeSteps(t *testing.T) {
	stdout, code, calls := runCommitChain(t, chainSlug, "", chainPublished, 0, 0, 0)
	want := "build step\npreflight step\ncommit step\n" + wantChainLine(chainPublished, "green", "green") + "\n"
	if code != 0 || stdout != want {
		t.Fatalf("chain = (%d, %q), want (0, %q)", code, stdout, want)
	}
	if want := [][]string{{chainWorktree, worktree.Home(), chainWorktree}}; !reflect.DeepEqual(calls.builds, want) {
		t.Fatalf("builds = %q, want %q: the commit's own worktree as root and target under the Bench home", calls.builds, want)
	}
	if want := [][]string{{chainWorktree, chainSlug, "", chainCandidate}}; !reflect.DeepEqual(calls.preflights, want) {
		t.Fatalf("preflight calls = %q, want %q: the candidate commit, not the published one", calls.preflights, want)
	}
}

// BC38, BC40, PP19: a chain that names a base gives it to the preflight with the candidate
// commit, and the commit help names --base <commit> in the --preflight-build group of the
// usage line and of the chain line.
func TestCommitChainPassesTheBase(t *testing.T) {
	t.Run("preflight call", func(t *testing.T) {
		_, code, calls := runCommitChain(t, chainSlug, "chunk-base", chainPublished, 0, 0, 0)
		if want := [][]string{{chainWorktree, chainSlug, "chunk-base", chainCandidate}}; code != 0 || !reflect.DeepEqual(calls.preflights, want) {
			t.Fatalf("chain = (%d, %q), want exit 0 and preflight calls %q", code, calls.preflights, want)
		}
	})
	t.Run("help", func(t *testing.T) {
		const group = "--preflight-build <slug> [--base <commit>]"
		t.Setenv(benchhome.Env, t.TempDir())
		var out, errOut bytes.Buffer
		code := Command{Stdout: &out, Stderr: &errOut, Executable: "bench"}.Run([]string{"commit", "--help"})
		usageLine, chainLine := strings.SplitN(out.String(), "\n", 2)[0], lastLine(withoutTreeTargetLine("commit", out.String()))
		if code != 0 || !strings.Contains(usageLine, "[--dry-run | "+group+"]") || !strings.HasPrefix(chainLine, group+": ") {
			t.Fatalf("help = (%d, %q), want %q in the usage line and at the start of the chain line", code, out.String(), group)
		}
	})
}

// BO52, PP22: a lane refusal before the check publishes nothing, so the chain calls no build
// and exits 1 with every step skipped.
func TestCommitChainSkipsAfterRefusal(t *testing.T) {
	stdout, code, calls := runCommitChain(t, chainSlug, "", "", 1, 0, 0)
	if want := wantChainLine("none", "skipped", "skipped"); code != 1 || lastLine(stdout) != want {
		t.Fatalf("chain = (%d, %q), want exit 1 ending in %q", code, stdout, want)
	}
	if len(calls.builds) != 0 || len(calls.preflights) != 0 {
		t.Fatalf("a refused commit ran builds %q and preflights %q", calls.builds, calls.preflights)
	}
}

// BO53, PP20: a red build refuses the commit at exit 1, calls no preflight, and publishes
// nothing.
func TestCommitChainStopsAtRedBuild(t *testing.T) {
	stdout, code, calls := runCommitChain(t, chainSlug, "", chainPublished, 0, 130, 0)
	if want := wantChainLine("none", "red", "skipped"); code != 1 || lastLine(stdout) != want {
		t.Fatalf("chain = (%d, %q), want exit 1 ending in %q", code, stdout, want)
	}
	if len(calls.preflights) != 0 {
		t.Fatalf("a red build ran preflights %q", calls.preflights)
	}
}

// BO54, PP21: a red preflight after a green build refuses the commit at exit 1, which means
// that nothing was committed, whatever the preflight's own exit.
func TestCommitChainReportsRedPreflight(t *testing.T) {
	for _, preflightExit := range []int{1, 2, 3} {
		stdout, code, _ := runCommitChain(t, chainSlug, "", chainPublished, 0, 0, preflightExit)
		if want := wantChainLine("none", "green", "red"); code != 1 || lastLine(stdout) != want {
			t.Fatalf("preflight exit %d: chain = (%d, %q), want exit 1 ending in %q", preflightExit, code, stdout, want)
		}
	}
}

// BO71, PP24: a commit that published without reconciling its checkout keeps exit 3, and the
// chain line names the published commit with both steps green.
func TestCommitChainStopsAtRemainder(t *testing.T) {
	stdout, code, _ := runCommitChain(t, chainSlug, "", chainPublished, 3, 0, 0)
	if want := wantChainLine(chainPublished, "green", "green"); code != 3 || lastLine(stdout) != want {
		t.Fatalf("chain = (%d, %q), want exit 3 ending in %q", code, stdout, want)
	}
}

// BO56: a commit without --preflight-build runs the commit step alone through the
// dispatcher: no build, no preflight, no chain line, and the commit step's exit.
func TestCommitWithoutPreflightBuildRunsCommitAlone(t *testing.T) {
	stdout, code, calls := runCommitChain(t, "", "", chainPublished, 0, 130, 1)
	if code != 0 || stdout != "commit step\n" {
		t.Fatalf("commit = (%d, %q), want (0, %q)", code, stdout, "commit step\n")
	}
	if len(calls.builds) != 0 || len(calls.preflights) != 0 {
		t.Fatalf("a commit without the flag ran builds %q and preflights %q", calls.builds, calls.preflights)
	}
}

// PP25, PP26: a help request through the production binding answers the commit's help, then
// the chain's help line, at exit 0. The chain line states that its steps run before
// publication, and no help line names the retired exit 4.
func TestCommitHelpEndsWithChainHelp(t *testing.T) {
	t.Setenv(benchhome.Env, t.TempDir())
	var out, errOut bytes.Buffer
	code := Command{Stdout: &out, Stderr: &errOut, Executable: "bench"}.Run([]string{"commit", "--help"})
	if code != 0 || !strings.HasPrefix(out.String(), "usage: bench commit") || !strings.HasSuffix(withoutTreeTargetLine("commit", out.String()), "\n"+chainHelp+"\n") {
		t.Fatalf("help = (%d, %q), want the commit help ending in the chain help line", code, out.String())
	}
	if want := "--preflight-build <slug> [--base <commit>]: before publication,"; !strings.HasPrefix(chainHelp, want) || !strings.Contains(chainHelp, "a red step publishes nothing") {
		t.Fatalf("chain help = %q, want it to start with %q and state that a red step publishes nothing", chainHelp, want)
	}
	if strings.Contains(out.String(), "exit 4") {
		t.Fatalf("help = %q, want no exit 4", out.String())
	}
}

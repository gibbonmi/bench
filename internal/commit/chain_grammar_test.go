package commit

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// parseArgs is the positional view of parseRequest that the parser table in
// commit_test.go reads.
func parseArgs(args []string) (msg string, paths []string, dryRun bool, help string, usageErr string) {
	req, help, usageErr := parseRequest(args)
	return req.msg, req.paths, req.dryRun, help, usageErr
}

// Command is the exit-only view of Run that the package tests read.
func Command(args []string, stdout, stderr io.Writer) int {
	_, exit := Run(args, stdout, stderr)
	return exit
}

// An empty --preflight-build or --base value is what an unset shell variable expands to, so
// the grammar refuses it before any repository read.
func TestCommitChainRefusesEmptyValue(t *testing.T) {
	for _, tc := range []struct {
		name, flag string
		args       []string
	}{
		{name: "slug", flag: PreflightBuildFlag, args: []string{"-m", "m", PreflightBuildFlag, "", "a.txt"}},
		{name: "base", flag: BaseFlag, args: []string{PreflightBuildFlag, "s", BaseFlag, "", "-m", "m", "p"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Chdir(t.TempDir())
			var stdout, stderr bytes.Buffer
			code := Command(tc.args, &stdout, &stderr)
			if want := toon.Usage(grammar.Cmd, usage.EmptyFlagValue(tc.flag)) + "\n"; code != 2 || stdout.Len() != 0 || stderr.String() != want {
				t.Fatalf("exit = %d, stdout = %q, stderr = %q; want exit 2, no stdout, and stderr %q", code, stdout.String(), stderr.String(), want)
			}
		})
	}
}

// BO55: a dry run grades without publishing, so it has no commit for the chain to build
// on. The combination is a usage refusal before any repository read.
func TestCommitChainRefusesDryRun(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Command([]string{"--dry-run", "--preflight-build", "slug", "-m", "m", "a.txt"}, &stdout, &stderr)
	if code != 2 || stdout.Len() != 0 {
		t.Fatalf("exit = %d, stdout = %q; want exit 2 and no stdout", code, stdout.String())
	}
	if line := strings.TrimSuffix(stderr.String(), "\n"); !strings.HasPrefix(line, grammar.Help+" (") || strings.Count(line, "\n") != 0 {
		t.Fatalf("stderr = %q, want the one commit usage line with its reason", stderr.String())
	}
}

// A refused chained commit re-runs with its base, so the pasted route passes the same base
// to the build preflight as the caller's own commit.
func TestRefusalReRunKeepsTheBase(t *testing.T) {
	quoted := sanitize.ShellQuote
	f := assignedCommitSet(t, "rerun", gateScript("exit 0"))
	code, stdout, stderr := runCommand(t, f.checkout, PreflightBuildFlag, "commit-fixture", "--base", "chunk-base", "-m", "m", "--", "missing.txt")
	next, printed := printedNext(stderr)
	steps := refusalroute.Steps(next)
	want := callerCommit(quoted("rerun"), PreflightBuildFlag, quoted("commit-fixture"), "--base", quoted("chunk-base"), "-m", quoted("m"), "--", quoted("missing.txt"))
	if code != 1 || !printed || steps[len(steps)-1] != want {
		t.Fatalf("commit = (%d, %q, %q), want exit 1 and a next= route that ends with %q", code, stdout, stderr, want)
	}
}

// Run hands the command layer the checkout it ran in, the slug --preflight-build named,
// the base --base named, and the commit it published: HEAD after a green commit or a
// publication remainder, and nothing after a refusal.
func TestRunOutcomeNamesThePublishedCommit(t *testing.T) {
	for _, tc := range []struct {
		name      string
		gateExit  int
		write     func(t *testing.T, root string)
		path      string
		wantExit  int
		published bool
	}{
		{name: "green", write: noWrite, path: "tracked.txt", published: true},
		{name: "remainder", write: unreconcilableGate("named"), path: "named", wantExit: 3, published: true},
		{name: "refusal", gateExit: 1, write: noWrite, path: "tracked.txt", wantExit: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, before := landingRepo(t, tc.gateExit, tc.write)
			if tc.path == "named" {
				gittest.MustRun(t, root, "reset", "-q", "--hard", "HEAD")
				namedDir(t, root, "named")
			}
			t.Chdir(root)
			outcome, code := Run([]string{"-m", "m", "--preflight-build", "slug", "--base", "chunk-base", tc.path}, io.Discard, io.Discard)
			head := gittest.Output(t, root, "rev-parse", "HEAD")
			want := Outcome{
				Root:           gittest.Output(t, root, "rev-parse", "--show-toplevel"),
				PreflightBuild: "slug",
				Base:           "chunk-base",
			}
			if tc.published {
				want.Published = head
			}
			if code != tc.wantExit || outcome != want || tc.published == (head == before) {
				t.Fatalf("Run = (%+v, %d), want (%+v, %d); HEAD %s -> %s", outcome, code, want, tc.wantExit, before, head)
			}
		})
	}
}

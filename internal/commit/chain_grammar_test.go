package commit

import (
	"bytes"
	"io"
	"strings"
	"testing"

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

// An empty --preflight-build value is what an unset shell variable expands to, so the
// grammar refuses it before any repository read.
func TestCommitChainRefusesEmptySlug(t *testing.T) {
	t.Chdir(t.TempDir())
	var stdout, stderr bytes.Buffer
	code := Command([]string{"-m", "m", "--preflight-build", "", "a.txt"}, &stdout, &stderr)
	if want := toon.Usage(grammar.Cmd, usage.EmptyFlagValue(PreflightBuildFlag)) + "\n"; code != 2 || stdout.Len() != 0 || stderr.String() != want {
		t.Fatalf("exit = %d, stdout = %q, stderr = %q; want exit 2, no stdout, and stderr %q", code, stdout.String(), stderr.String(), want)
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

// Run hands the command layer the checkout it ran in, the slug --preflight-build named,
// and the commit it published: HEAD after a green commit or a publication remainder,
// and nothing after a refusal.
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
				runGit(t, root, "reset", "-q", "--hard", "HEAD")
				namedDir(t, root, "named")
			}
			t.Chdir(root)
			outcome, code := Run([]string{"-m", "m", "--preflight-build", "slug", tc.path}, io.Discard, io.Discard)
			head := strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD")))
			want := Outcome{
				Root:           strings.TrimSpace(string(runGit(t, root, "rev-parse", "--show-toplevel"))),
				PreflightBuild: "slug",
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

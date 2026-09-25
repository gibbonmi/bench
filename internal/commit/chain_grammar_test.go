package commit

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

// parseArgs is the positional view of parseRequest that the parser table in
// commit_test.go reads.
func parseArgs(args []string) (msg string, paths []string, dryRun bool, help string, usageErr string) {
	req, help, usageErr := parseRequest(args)
	return req.msg, req.paths, req.dryRun, help, usageErr
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

// chainRun runs the chain with the real commit step in root and records the build target
// and the preflight argv the chain hands on. Both chained steps answer green.
func chainRun(t *testing.T, root string, args ...string) (code int, stdout string, builds []string, preflights [][]string) {
	t.Helper()
	t.Chdir(root)
	steps := ChainSteps{
		Commit: Run,
		Build: func(_, _ string, args []string, _, _ io.Writer) int {
			builds = append(builds, args...)
			return 0
		},
		Home: t.TempDir,
		Preflight: func(args []string) (string, int) {
			preflights = append(preflights, args)
			return "", 0
		},
	}
	var out, errOut bytes.Buffer
	code = Chain(steps, args, &out, &errOut)
	return code, out.String(), builds, preflights
}

// The real commit step hands the chain the commit it published and the checkout it ran
// in, so the chain line and the preflight tip name the new HEAD.
func TestChainReadsThePublishedCommitFromTheCommitStep(t *testing.T) {
	root, before := landingRepo(t, 0, noWrite)
	code, stdout, builds, preflights := chainRun(t, root, "-m", "m", "--preflight-build", "slug", "tracked.txt")
	head := strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD")))
	if code != 0 || head == before {
		t.Fatalf("exit = %d, HEAD %s -> %s; want a published commit at exit 0; stdout=%q", code, before, head, stdout)
	}
	if !strings.HasSuffix(stdout, "commit-chain{commit="+head+",build=green,preflight=green}\n") {
		t.Fatalf("stdout = %q, want the chain line naming the new HEAD %s", stdout, head)
	}
	toplevel := strings.TrimSpace(string(runGit(t, root, "rev-parse", "--show-toplevel")))
	if len(builds) != 1 || builds[0] != toplevel {
		t.Fatalf("build targets = %q, want the commit's checkout %q", builds, toplevel)
	}
	if len(preflights) != 1 || strings.Join(preflights[0], " ") != "build slug --source-tip "+head {
		t.Fatalf("preflight argv = %q, want the build preflight at the new HEAD", preflights)
	}
}

// A commit that published without reconciling hands the chain its published commit, and
// a refused commit hands it none.
func TestChainNamesThePublicationOfEveryCommitExit(t *testing.T) {
	t.Run("remainder", func(t *testing.T) {
		root, _ := landingRepo(t, 0, unreconcilableGate("named"))
		runGit(t, root, "reset", "-q", "--hard", "HEAD")
		namedDir(t, root, "named")
		code, stdout, builds, _ := chainRun(t, root, "-m", "m", "--preflight-build", "slug", "named")
		head := strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD")))
		if code != 3 || len(builds) != 0 || !strings.HasSuffix(stdout, "commit-chain{commit="+head+",build=skipped,preflight=skipped}\n") {
			t.Fatalf("chain = (%d, %q, builds %q), want exit 3 naming HEAD %s with no build", code, stdout, builds, head)
		}
	})
	t.Run("refusal", func(t *testing.T) {
		root, _ := landingRepo(t, 1, noWrite)
		code, stdout, builds, _ := chainRun(t, root, "-m", "m", "--preflight-build", "slug", "tracked.txt")
		if code != 1 || len(builds) != 0 || !strings.HasSuffix(stdout, "commit-chain{commit=none,build=skipped,preflight=skipped}\n") {
			t.Fatalf("chain = (%d, %q, builds %q), want exit 1 naming no commit and no build", code, stdout, builds)
		}
	})
}

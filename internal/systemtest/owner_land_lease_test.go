//go:build system

package systemtest

import (
	"os"
	"strings"
	"testing"
)

// TestTwoLandingProcessesShareOneLease is BW10. Each landing runs in its own process, so
// only the kernel record lock tells the contender that the holder is alive. The fixture's
// gate waits at its barrier for a tree with loser.txt, so the holder writes that file.
func TestTwoLandingProcessesShareOneLease(t *testing.T) {
	root, home, tally, trees, ready, release := systemLandingRaceFixture(t)
	base := systemGitOutput(t, root, "rev-parse", "main")
	env := systemLandEnv(root, home, tally, trees, ready, release)
	holder := systemReviewedLandingWorktree(t, root, home, env, "lease-holder", "lease holder", "loser.txt", base)
	contender := systemReviewedLandingWorktree(t, root, home, env, "lease-contender", "lease contender", "winner.txt", base)
	holderDone, holderOut, holderErr := systemLandAtBarrier(t, root, home, tally, trees, ready, release, holder, base)

	contended := systemLand(t, root, home, tally, trees, ready, release, contender, base)
	detail := "refused{detail=landing destination refs/heads/main is leased by assignment " + holder.assignment + ","
	if contended.code != 1 || !strings.HasPrefix(contended.stdout, detail) || strings.Contains(contended.stdout, "bench doctor") {
		t.Fatalf("contender land = (%d, %q, %q), want the lease refusal naming %s", contended.code, contended.stdout, contended.stderr, holder.assignment)
	}
	releaseArtifactBarrier(t, release)
	if code := systemExitCode(<-holderDone); code != 0 || !strings.Contains(holderOut.String(), ",worktree=released,") {
		t.Fatalf("holder land = (%d, %q, %q), want a released publication", code, holderOut.String(), holderErr.String())
	}
	if published := systemGitOutput(t, root, "ls-tree", "--name-only", "main", "loser.txt"); published != "loser.txt" {
		t.Fatalf("main after the holder's publication lists %q, want loser.txt", published)
	}
	if got, err := os.ReadFile(tally); err != nil || string(got) != "l" {
		t.Fatalf("gate tally = %q, %v; want the holder's one run", got, err)
	}
	if got := systemGitOutput(t, contender.path, "rev-parse", "HEAD"); got != contender.tip {
		t.Fatalf("contender source tip moved: got %s want %s", got, contender.tip)
	}
}

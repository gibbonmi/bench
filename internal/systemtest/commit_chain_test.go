//go:build system

package systemtest

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// TestCommitChainRefusesRedCandidate is PP28. In a linked-project worktree, the real
// `bench commit --preflight-build` grades a plan commit before publication. A red plan commit
// exits 1 with the branch ref unchanged, and its repaired retry publishes, with the chain line
// naming the new branch tip.
func TestCommitChainRefusesRedCandidate(t *testing.T) {
	t.Parallel()
	const (
		slug   = "record"
		spec   = "specs/" + slug + "/spec.md"
		ticket = "specs/" + slug + "/tickets/one.md"
	)
	fixture := scaffoldRecordedPublicationRepo(t)
	commitmenttest.Write(t, fixture.root, spec, preflighttest.SpecBody(slug))
	commitmenttest.Write(t, fixture.root, ticket, preflighttest.TicketDoc("One", "PF1", "PF2"))
	commitmenttest.SeedAdmission(t, fixture.root, spec)
	// The chain's build step runs the project's build script. A linked project publishes no
	// executable, so the planted script builds nothing and the build preflight's binary-seal
	// row stays not applicable.
	commitmenttest.Write(t, fixture.root, "scripts/go-build.sh", "#!/bin/sh\nexit 0\n")
	systemGit(t, fixture.root, "add", ".")
	systemGit(t, fixture.root, systemCommitArgs("-qm", "chain fixture")...)
	systemGit(t, fixture.root, "update-ref", "refs/bench/green/main", systemGitOutput(t, fixture.root, "rev-parse", "HEAD"))

	source := systemCreateLandingWorktree(t, fixture.root, fixture.home, "commit-chain", "commit chain")
	commitChain := func(ticketBody string) processResult {
		t.Helper()
		commitmenttest.Write(t, source.path, ticket, ticketBody)
		return systemSelected(t, source.path, fixture.environment(fixture.home), "commit", "--preflight-build", slug, "-m", "plan", "--", ticket)
	}
	before := systemGitOutput(t, source.path, "rev-parse", source.branch)

	// A ticket that covers a row the coverage map does not hold is a red plan.
	red := commitChain(preflighttest.TicketDoc("One", "PF1", "PF9"))
	if red.code != 1 || !strings.HasSuffix(strings.TrimSpace(red.stdout), "commit-chain{commit=none,build=green,preflight=red}") {
		t.Fatalf("red plan commit = (%d, %q, %q), want exit 1 and a red preflight with nothing committed", red.code, red.stdout, red.stderr)
	}
	if after := systemGitOutput(t, source.path, "rev-parse", source.branch); after != before {
		t.Fatalf("a red plan commit moved %s from %s to %s", source.branch, before, after)
	}

	green := commitChain(preflighttest.TicketDoc("One", "PF1", "PF2") + "\nA plan note.\n")
	tip := systemGitOutput(t, source.path, "rev-parse", source.branch)
	if want := "commit-chain{commit=" + tip + ",build=green,preflight=green}"; green.code != 0 || tip == before || !strings.HasSuffix(strings.TrimSpace(green.stdout), want) {
		t.Fatalf("repaired plan commit = (%d, %q, %q), want exit 0 ending in %q", green.code, green.stdout, green.stderr, want)
	}
}

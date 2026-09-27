// The retire listing: after `bench spec retire` removes a merged spec, the dispatch names the
// assignments that spec's work superseded and the count of unique unclaimed refs. It lives
// here, not in internal/spec, because internal/worktree already imports internal/spec.

package main

import (
	"fmt"
	"strings"

	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/worktree"
)

// retireListingCommand runs the spec family and prints its response with the retire listing
// inserted before the retire's next: line. The exit code is always the family's own.
func retireListingCommand(fn func([]string) (string, int)) commandHandler {
	return func(c Command, args []string) int {
		out, code := fn(args)
		fmt.Fprint(c.Stdout, withRetireListing(args, out, code))
		return code
	}
}

// withRetireListing inserts the listing only after a code-0 retire that printed a next: line.
// Help and every refusal print no next: line, so they pass through unchanged.
func withRetireListing(args []string, out string, code int) string {
	if code != 0 || len(args) < 2 || args[0] != "retire" {
		return out
	}
	at := strings.LastIndex("\n"+out, "\n"+spec.RetireNextPrefix)
	if at < 0 {
		return out
	}
	return out[:at] + retireListing(spec.RepoBase(), spec.SlugOf(args[len(args)-1])) + out[at:]
}

// retireListing is one candidate line per active or cleanup-pending assignment whose label or
// request token contains slug, then the unique count line. It reads and discards nothing.
func retireListing(root, slug string) string {
	assignments, err := intent.Assignments(root)
	if err != nil {
		return uniqueRefsLine("unavailable — " + err.Error())
	}
	var b strings.Builder
	for _, a := range assignments {
		if supersededCandidate(a, slug) {
			fmt.Fprintf(&b, "superseded candidate: %s %s — %s\n", a.ID, a.Label, worktree.DiscardTargetCommand(a.ID))
		}
	}
	counts, err := worktree.CountUnclaimedRefs(root)
	if err != nil {
		return b.String() + uniqueRefsLine("unavailable — "+err.Error())
	}
	return b.String() + uniqueRefsLine(fmt.Sprintf("%d — %s", counts.Unique(), worktree.UnclaimedPlanCommand()))
}

// supersededCandidate holds a live record whose label or request token names the slug. A
// complete record is retired work; the calling worktree's own record is listed like any other.
func supersededCandidate(a intent.Assignment, slug string) bool {
	if a.State != intent.StateActive && a.State != intent.StateCleanupPending {
		return false
	}
	return strings.Contains(a.Label, slug) || strings.Contains(a.RequestToken, slug)
}

// uniqueRefsLine is the count line with its value: a count and the plan command, or the
// unavailable reason.
func uniqueRefsLine(value string) string { return "unique refs: " + value + "\n" }

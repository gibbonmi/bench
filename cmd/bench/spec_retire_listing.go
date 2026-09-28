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
func retireListingCommand(fn func([]string) (string, int, string)) commandHandler {
	return func(c Command, args []string) int {
		out, code, operand := fn(args)
		fmt.Fprint(c.Stdout, withRetireListing(operand, out, code))
		return code
	}
}

// withRetireListing inserts the listing only after a code-0 retire that printed a next: line.
// A refusal can reflect an operand that carries a next: line, so the exit code, not the
// marker, keeps every refusal unchanged.
func withRetireListing(operand, out string, code int) string {
	if code != 0 || operand == "" {
		return out
	}
	at := strings.LastIndex("\n"+out, "\n"+spec.RetireNextPrefix)
	if at < 0 {
		return out
	}
	return out[:at] + retireListing(spec.RepoBase(), spec.SlugOf(operand)) + out[at:]
}

// retireListing is one candidate line per active or cleanup-pending assignment whose label or
// request token contains slug, then the count line. A faulted row is not unique, so it adds
// its own count rather than a unique one. It reads and discards nothing.
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
	count := fmt.Sprint(counts.Unique())
	if counts.Faulted > 0 {
		count += fmt.Sprintf(", %d faulted", counts.Faulted)
	}
	return b.String() + uniqueRefsLine(count+" — "+worktree.UnclaimedPlanCommand())
}

// supersededCandidate holds a live record whose label or request token names the slug. A
// complete record is retired work; the calling worktree's own record is listed like any other.
func supersededCandidate(a intent.Assignment, slug string) bool {
	if a.State != intent.StateActive && a.State != intent.StateCleanupPending {
		return false
	}
	return strings.Contains(a.Label, slug) || strings.Contains(a.RequestToken, slug)
}

// uniqueRefsLine is the count line with its value: the unique count, a faulted count when one
// is nonzero, and the plan command, or the unavailable reason.
func uniqueRefsLine(value string) string { return "unique refs: " + value + "\n" }

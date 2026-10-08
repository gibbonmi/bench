// Refusal route fixtures for the commitment faces that the commit prints: the commitment
// policy refuses the commit's candidate and raises its face, and the commit prints that
// face's route. The fixture checkouts' delivery admission lives here too.
package commit

import (
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/refusalroute/routetest"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/worktree"
)

// commitFixtureSpec is the delivery the fixture repository's policy approves.
const commitFixtureSpec = "specs/commit-fixture/spec.md"

// assignedCommitSet is the commit fixture whose checkout is an assignment that the create
// verb made and the delivery admits.
func assignedCommitSet(t *testing.T, label string, gate func(t *testing.T, root string)) commitSet {
	t.Helper()
	f := unboundCommitSet(t, label, gate)
	commitmenttest.Admit(t, f.checkout, label, commitFixtureSpec)
	return f
}

// unboundCommitSet is the commit fixture whose checkout is an assignment that the create
// verb made and no delivery admits, so a printed tree target resolves as the CLI resolves
// it. gate writes the checkout's gate, which its base commit holds. The caller's change to
// a.txt waits uncommitted.
func unboundCommitSet(t *testing.T, label string, gate func(t *testing.T, root string)) commitSet {
	t.Helper()
	f := primaryCommitSet(t)
	runRouteStep(t, f, "bench worktree create --request "+label+" --label "+label)
	f.checkout = createdCheckout(t, f, label)
	prepareLandingCheckout(t, f.checkout, gate)
	runGit(t, f.checkout, "reset", "-q", "--hard", "HEAD")
	mustWrite(t, filepath.Join(f.checkout, "a.txt"), "changed\n", 0o644)
	return f
}

// createdCheckout is the checkout of the assignment that a create of label made.
func createdCheckout(t *testing.T, f commitSet, label string) string {
	t.Helper()
	checkout, err := worktree.TreeTarget(f.primary, label)
	if err != nil {
		t.Fatalf("created assignment %q: %v", label, err)
	}
	return checkout
}

// admitCreated admits the assignment that a create of label made, as the operator's start
// of the delivery would, and returns its checkout.
func admitCreated(t *testing.T, f commitSet, label string) string {
	t.Helper()
	checkout := createdCheckout(t, f, label)
	commitmenttest.Admit(t, checkout, label, commitFixtureSpec)
	return checkout
}

// commitmentFaceFixtures are the walk's fixtures of the commitment faces that the commit
// prints. rerun is the fixtures' own commit at the worktree of a label.
func commitmentFaceFixtures(rerun func(label string, flags ...string) string) []commitFaceFixture {
	const spaced, unsafe = "specs/light path/tickets/one.md", "specs/light\x1bpath/tickets/one.md"
	return []commitFaceFixture{
		{
			// No delivery admits the assignment, and no light-path ticket covers a.txt. The
			// operator names the outcome, its own request, and the deliverable that the
			// fixture's policy approves.
			face: refusalroute.CommitmentUnbound,
			exit: 1,
			build: func(t *testing.T) commitSet {
				return unboundCommitSet(t, refusalroute.CommitmentUnbound, gateScript("exit 0"))
			},
			contains: "bench commitment start --outcome ",
			slots:    []string{"<id>", commitmenttest.DeliveryOutcome, "<request>", refusalroute.CommitmentUnbound, "<path>", commitFixtureSpec},
		},
		{
			// The branch adds a recommended sequence that the active policy does not project.
			// The reviewer withdraws the protected change.
			face: refusalroute.CommitmentDecision,
			exit: 1,
			build: func(t *testing.T) commitSet {
				return assignedCommitSet(t, refusalroute.CommitmentDecision, func(t *testing.T, root string) {
					gateScript("exit 0")(t, root)
					mustWrite(t, filepath.Join(root, "ROADMAP.md"), "# Roadmap\n\n## Recommended sequence\n\n1. "+commitmenttest.DeliveryOutcome+"\n", 0o644)
				})
			},
			prefix:   routetest.ReviewerMarker,
			contains: "bench commitment plan --input ",
			carry: map[int]func(*testing.T, commitSet){0: func(t *testing.T, f commitSet) {
				runGit(t, f.checkout, "rm", "-q", "ROADMAP.md")
				runGit(t, f.checkout, "commit", "-qm", "withdraw the protected change")
			}},
		},
		lightPathFixture("", spaced, sanitize.ShellQuote(spaced), rerun),
		lightPathFixture("ticket that is not line-safe", unsafe, "<"+refusalroute.FactTicket+">", rerun),
		lightPathSpanFixture(),
	}
}

// lightPathFixture is a commitment-light-path-outside fixture: the checkout's one committed
// light-path ticket, at ticket, lists b.txt and not the caller's a.txt. The route names the
// ticket as named. The walk widens the Writes line, the route commits the ticket, and the
// caller's own re-run then publishes. A ticket that prints its placeholder is the value the
// operator holds.
func lightPathFixture(cause, ticket, named string, rerun func(label string, flags ...string) string) commitFaceFixture {
	label := refusalroute.CommitmentLightPathOutside
	return commitFaceFixture{
		face: label, cause: cause, exit: 1,
		build: func(t *testing.T) commitSet {
			return unboundCommitSet(t, label, func(t *testing.T, root string) {
				gateScript("exit 0")(t, root)
				commitmenttest.WriteLightTicket(t, root, ticket, "b.txt")
			})
		},
		contains: " line of " + named + "; then " + callerCommit(sanitize.ShellQuote(label), "-m", "<msg>", "--", named) + "; then ",
		suffix:   rerun(label),
		absent:   "bench commitment start",
		carry: map[int]func(*testing.T, commitSet){0: func(t *testing.T, f commitSet) {
			commitmenttest.WriteLightTicket(t, f.checkout, ticket, "b.txt", "a.txt")
		}},
		slots: []string{named, sanitize.ShellQuote(ticket), "<msg>", "'widen the ticket'"},
	}
}

// lightPathSpanFixture is the commitment-light-path-span fixture: the caller commits a.txt
// and b.txt, and each has its own light-path ticket. The operator commits the path of one
// ticket, a.txt, which the printed commit carries out, so the instruction needs no step of
// its own.
func lightPathSpanFixture() commitFaceFixture {
	label := refusalroute.CommitmentLightPathSpan
	return commitFaceFixture{
		face: label, exit: 1, paths: []string{"b.txt"},
		build: func(t *testing.T) commitSet {
			f := unboundCommitSet(t, label, func(t *testing.T, root string) {
				gateScript("exit 0")(t, root)
				commitmenttest.WriteLightTicket(t, root, "specs/one/tickets/one.md", "a.txt")
				commitmenttest.WriteLightTicket(t, root, "specs/two/tickets/one.md", "b.txt")
			})
			mustWrite(t, filepath.Join(f.checkout, "b.txt"), "b\n", 0o644)
			return f
		},
		suffix: callerCommit(sanitize.ShellQuote(label), "-m", "<msg>", "--", "<path>..."),
		absent: routetest.ReviewerMarker,
		carry:  map[int]func(*testing.T, commitSet){0: func(*testing.T, commitSet) {}},
		slots:  []string{"<msg>", "'one ticket'", "<path>...", "'a.txt'"},
	}
}

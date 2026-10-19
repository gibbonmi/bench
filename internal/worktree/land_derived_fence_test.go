// Derived-fence tests for the landing command: the ticket `Writes:` union plus the implicit
// entries authorize each changed path, and a spec fence section authorizes nothing.
package worktree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestLandCommandAuthorizesFromTicketWrites covers CR55. The source adds a path to the
// ticket's `Writes:` line and commits that path, and no fence section lists it.
func TestLandCommandAuthorizesFromTicketWrites(t *testing.T) {
	t.Parallel()
	request := "land-derived-fence-ticket"
	f := landSurface(t, request)
	ticket := filepath.Join(f.creation.Path, filepath.FromSlash(landingTicket))
	body, err := os.ReadFile(ticket)
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, ticket, withTicketWrite(body, "ticket-only.txt"), 0o644)
	gitRun(t, f.creation.Path, "add", landingTicket)
	commitInWorktree(t, f.creation.Path, "ticket-only.txt", "ticket bytes\n", "ticket-only write")
	if spec := gitOutput(t, f.creation.Path, "show", "HEAD:specs/x/spec.md"); strings.Contains(spec, "ticket-only.txt") {
		t.Fatalf("spec names the ticket-only path:\n%s", spec)
	}
	refreshLandingEvidence(t, f.creation.Path, f.base)
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	r := runVerb(t, verbLand, f.processHomeCall(landArgs(request, f.base, tip, f.creation.Path)...))
	if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released,census=0}") {
		t.Fatalf("ticket-only landing = (%d, %q, %q), want released", r.exit, r.stdout, r.stderr)
	}
	if got := gitOutput(t, f.root, "show", "main:ticket-only.txt"); got != "ticket bytes" {
		t.Fatalf("published ticket-only path = %q", got)
	}
}

// TestLandCommandAuthorizesTheReviewPickup covers CR58. The landing's own review record is
// an implicit entry that no ticket names, and another review record stays outside the fence.
func TestLandCommandAuthorizesTheReviewPickup(t *testing.T) {
	t.Parallel()
	t.Run("pickup", func(t *testing.T) {
		t.Parallel()
		request := "land-derived-fence-pickup"
		f := landSurface(t, request)
		if ticket := gitOutput(t, f.creation.Path, "show", f.tip+":"+landingTicket); strings.Contains(ticket, "reviews/x.md") {
			t.Fatalf("ticket names the review pickup:\n%s", ticket)
		}
		r := runVerb(t, verbLand, f.processHomeCall(landArgs(request, f.base, f.tip, f.creation.Path)...))
		if r.exit != 0 || !strings.Contains(r.stdout, "worktree=released,census=0}") {
			t.Fatalf("pickup landing = (%d, %q, %q), want released", r.exit, r.stdout, r.stderr)
		}
		if got := gitOutput(t, f.root, "show", "main:reviews/x.md"); got == "" {
			t.Fatal("landing published no review pickup")
		}
	})
	t.Run("another review record", func(t *testing.T) {
		t.Parallel()
		request := "land-derived-fence-other-review"
		f := landSurface(t, request)
		commitInWorktree(t, f.creation.Path, "reviews/other.md", "other record\n", "another review record")
		tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
		r := runVerb(t, verbLand, f.processHomeCall(landArgs(request, f.base, tip, f.creation.Path)...))
		if r.exit != 1 || !strings.Contains(r.stdout, refusalPathsTable+"[1]{path}:\n  reviews/other.md\n") {
			t.Fatalf("other review landing = (%d, %q, %q), want reviews/other.md refused", r.exit, r.stdout, r.stderr)
		}
	})
}

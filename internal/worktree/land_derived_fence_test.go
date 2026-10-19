// Derived-fence tests for the landing command: the ticket `Writes:` union plus the implicit
// entries authorize each changed path, and a spec fence section authorizes nothing.
package worktree

import (
	"os"
	"path"
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

// TestLandCommandReadsOnlyTrackedTickets checks that the fence reads only the tickets of the
// source tip. A git-ignored ticket-shaped file in the clean source checkout names a committed
// path, and the landing refuses that path.
func TestLandCommandReadsOnlyTrackedTickets(t *testing.T) {
	t.Parallel()
	request := "land-derived-fence-ignored-ticket"
	f := landSurface(t, request)
	const notes = ".notes.md"
	exclude := gitOutput(t, f.creation.Path, "rev-parse", "--path-format=absolute", "--git-path", "info/exclude")
	mustMkdirAll(t, filepath.Dir(exclude), 0o755)
	mustWrite(t, exclude, []byte(notes+"\n"), 0o644)
	ignored := path.Join(path.Dir(landingTicket), notes)
	mustWrite(t, filepath.Join(f.creation.Path, filepath.FromSlash(ignored)), []byte("Writes: ignored-only.txt\n"), 0o644)
	if status := gitOutput(t, f.creation.Path, "status", "--porcelain", "--ignored", path.Dir(ignored)); status != "!! "+ignored {
		t.Fatalf("ticket-shaped file status = %q, want it ignored", status)
	}
	commitInWorktree(t, f.creation.Path, "ignored-only.txt", "ignored bytes\n", "write named only by an ignored file")
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	r := runVerb(t, verbLand, f.processHomeCall(landArgs(request, f.base, tip, f.creation.Path)...))
	if r.exit != 1 || !strings.Contains(r.stdout, refusalPathsTable+"[1]{path}:\n  ignored-only.txt\n") {
		t.Fatalf("ignored-ticket landing = (%d, %q, %q), want ignored-only.txt refused", r.exit, r.stdout, r.stderr)
	}
}

// TestLandCommandRefusesAnUnreadableTicketSet checks that a tracked ticket read that fails
// refuses the landing, and no checkout ticket replaces it. The tip commits a ticket as a
// symbolic link to a ticket-shaped file whose `Writes:` line names a committed path.
func TestLandCommandRefusesAnUnreadableTicketSet(t *testing.T) {
	t.Parallel()
	request := "land-derived-fence-linked-ticket"
	f := landSurface(t, request)
	linked := path.Join(path.Dir(landingTicket), "x.md")
	mustWrite(t, filepath.Join(f.creation.Path, "specs", "x", "linked.md"), []byte("Writes: linked-only.txt\n"), 0o644)
	if err := os.Symlink("../linked.md", filepath.Join(f.creation.Path, filepath.FromSlash(linked))); err != nil {
		t.Fatal(err)
	}
	gitRun(t, f.creation.Path, "add", "specs/x")
	commitInWorktree(t, f.creation.Path, "linked-only.txt", "linked bytes\n", "write named only by a linked ticket")
	refreshLandingEvidence(t, f.creation.Path, f.base)
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")
	r := runVerb(t, verbLand, f.processHomeCall(landArgs(request, f.base, tip, f.creation.Path)...))
	if r.exit != 1 || !strings.Contains(r.stdout, "special source "+linked) {
		t.Fatalf("linked-ticket landing = (%d, %q, %q), want the link refused", r.exit, r.stdout, r.stderr)
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

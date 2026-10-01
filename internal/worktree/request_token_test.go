package worktree

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

// The create verb persists the plain request token beside its digest, and list prints
// it, so a resumed landing reuses the exact token without a refusal round-trip. The
// digest stays the authorization identity; the token column is recall, not proof.
func TestListPrintsThePersistedRequestToken(t *testing.T) {
	f := newOwnedAssignment(t, "token-recall")
	if f.creation.Assignment.RequestToken == "" {
		t.Fatal("create persisted no request token")
	}
	chdir(t, f.root)
	listed := runVerb(t, verbList, f.call())
	if listed.exit != 0 {
		t.Fatalf("list code=%d out=%q", listed.exit, listed.stdout)
	}
	if !strings.Contains(listed.stdout, "request") || !strings.Contains(listed.stdout, f.creation.Assignment.RequestToken) {
		t.Fatalf("list omitted the request token %q:\n%s", f.creation.Assignment.RequestToken, listed.stdout)
	}
}

// A record written before the field existed carries none: it still loads, and its row
// prints an empty token cell rather than failing the whole listing.
func TestListToleratesAPreTokenRecord(t *testing.T) {
	f := newOwnedAssignment(t, "token-absent")
	stripped := f.creation.Assignment
	stripped.RequestToken = ""
	mustNoError(t, intent.PutAssignment(f.root, stripped))
	chdir(t, f.root)
	if path := runVerb(t, verbPath, f.call(f.creation.Assignment.ID)); path.exit != 0 {
		t.Fatalf("pre-token record no longer resolves: %s", path.stderr)
	}
	listed := runVerb(t, verbList, f.call())
	if listed.exit != 0 || !strings.Contains(listed.stdout, f.creation.Assignment.ID) {
		t.Fatalf("list dropped the pre-token record (code=%d):\n%s", listed.exit, listed.stdout)
	}
	if !strings.Contains(listed.stdout, f.creation.Assignment.Label+",\"\",") {
		t.Fatalf("pre-token row does not show an empty request cell:\n%s", listed.stdout)
	}
}

package worktree

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

// The create verb persists the plain request token beside its digest, and list prints
// it, so a resumed landing reuses the exact token without a refusal round-trip. The
// digest stays the authorization identity; the token column is recall, not proof.
func TestListPrintsThePersistedRequestToken(t *testing.T) {
	root, creation := newOwnedAssignment(t, "token-recall")
	if creation.Assignment.RequestToken == "" {
		t.Fatal("create persisted no request token")
	}
	chdir(t, root)
	out, code := ListCommand(nil)
	if code != 0 {
		t.Fatalf("list code=%d out=%q", code, out)
	}
	if !strings.Contains(out, "request") || !strings.Contains(out, creation.Assignment.RequestToken) {
		t.Fatalf("list omitted the request token %q:\n%s", creation.Assignment.RequestToken, out)
	}
}

// A record written before the field existed carries none: it still loads, and its row
// prints an empty token cell rather than failing the whole listing.
func TestListToleratesAPreTokenRecord(t *testing.T) {
	root, creation := newOwnedAssignment(t, "token-absent")
	stripped := creation.Assignment
	stripped.RequestToken = ""
	mustNoError(t, intent.PutAssignment(root, stripped))
	chdir(t, root)
	var stdout, stderr bytes.Buffer
	if code := PathCommand(root, []string{creation.Assignment.ID}, &stdout, &stderr); code != 0 {
		t.Fatalf("pre-token record no longer resolves: %s", stderr.String())
	}
	out, code := ListCommand(nil)
	if code != 0 || !strings.Contains(out, creation.Assignment.ID) {
		t.Fatalf("list dropped the pre-token record (code=%d):\n%s", code, out)
	}
	if !strings.Contains(out, creation.Assignment.Label+",\"\",") {
		t.Fatalf("pre-token row does not show an empty request cell:\n%s", out)
	}
}

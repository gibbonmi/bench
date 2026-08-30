// The gate authority: attribution of an unpublished tree, beside the lane authority.
package authorization

import (
	"bytes"
	"context"
	"testing"
)

// An infrastructure attribution carries the gate's own open reason, so the refusal that
// reads this result can name the cause instead of the kind alone. A fixture with no gate
// script leaves the subject open, which is exactly the operational outcome the kind
// covers.
func TestAuthorizeCarriesTheInfrastructureReason(t *testing.T) {
	root, tree := laneFixture(t)
	var stdout, stderr bytes.Buffer
	result := AuthorizeWithWriters(context.Background(), root, tree, &stdout, &stderr)
	if result.Kind != Infrastructure {
		t.Fatalf("kind = %q, want %q", result.Kind, Infrastructure)
	}
	if result.Reason == "" {
		t.Fatalf("reason = %q, want the gate's open reason", result.Reason)
	}
	t.Logf("infrastructure reason = %q", result.Reason)
}

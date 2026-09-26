package evidencecmd_test

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// TestReviewRecordChargeOrder is RE12. The author record commits first, and a review charge
// prepared from that record commit binds to the current action. A later record commit
// moves the source tip, so the same charge then refuses as stale. The review record is no
// prepared repository source, so only the source-tip comparison can refuse the stale charge.
func TestReviewRecordChargeOrder(t *testing.T) {
	_, slug, _ := preflighttest.SeedReviewEvidence(t, false)
	commitRecord := func(body, message string) string {
		t.Helper()
		preflighttest.MustWriteFile(t, "reviews/"+slug+".md", body)
		preflighttest.RunGit(t, "add", "-A")
		preflighttest.RunGit(t, "commit", "-q", "-m", message)
		return preflighttest.RunGit(t, "rev-parse", "HEAD")
	}

	recordTip := commitRecord("# Review record\n\nAuthor verification and probe record.\n", "author record")
	identity, prepared, _ := prepareEvidence(t, preflighttest.ReviewArgs(t, slug))
	if prepared["source_tip"] != recordTip {
		t.Fatalf("prepared source tip = %v, want the record commit %s", prepared["source_tip"], recordTip)
	}
	out, code := checkCurrent(t, identity)
	if code != 0 || !strings.Contains(out, "current[1]") {
		t.Fatalf("charge prepared after the record commit = (%d):\n%s", code, out)
	}

	commitRecord("# Review record\n\nAuthor verification and probe record.\n\nA later record edit.\n", "later record")
	out, code = checkCurrent(t, identity)
	if code != 1 || !strings.Contains(out, "is not the prepared source tip "+recordTip) || strings.Contains(out, "current[1]") {
		t.Fatalf("charge after a later record commit = (%d):\n%s", code, out)
	}
}

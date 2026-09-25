package gate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// TestReviewCheckpointUncoveredRangeNamesItsCommits grades the uncovered-source face.
// A repair committed after the last chunk review leaves a range no review covers, and the
// refusal names that range, so the caller does not find the commits by hand.
func TestReviewCheckpointUncoveredRangeNamesItsCommits(t *testing.T) {
	f := checkpointFixture(t)
	reviewed := f.Record.Chunks[0].Tip
	f.Write("source.txt", "later repair\n")
	f.Commit("unreviewed repair")
	code, out := runCheckpoint(t, f)
	if want := reviewed + ".." + f.Tip(); code == 0 || !strings.Contains(out, want) {
		t.Fatalf("uncovered-source refusal = (%d, %q), want the uncovered range %s", code, out, want)
	}
}

// TestReviewCheckpointChainGapNamesTheExpectedBase grades the chain-gap face. The
// refusal names the predecessor tip the chunk base must be, and the rule that places
// plan commits and default-branch merges inside the chunk delta.
func TestReviewCheckpointChainGapNamesTheExpectedBase(t *testing.T) {
	f := recordtest.Attach(t, outcomeFixture(t), 2)
	f.AddChunk()
	f.Save()
	f.Commit("retain first results")
	f.AddChunk()
	f.Save()
	f.Commit("retain second results")
	expected := f.Record.Chunks[0].Tip
	f.Record.Chunks[1].Base = f.Record.Chunks[0].Base
	f.Save()
	var out, diagnostic bytes.Buffer
	code := RunCommand([]string{f.Root, "--checkpoint", recordtest.Spec, "--chunk", "2"}, &out, &diagnostic)
	text := out.String() + diagnostic.String()
	if code == 0 || !strings.Contains(text, "expected base "+expected) {
		t.Fatalf("chain-gap refusal = (%d, %q), want expected base %s", code, text, expected)
	}
	if !strings.Contains(text, "plan commits land before the ticket merge, a default-branch merge lands only before the first chunk") {
		t.Fatalf("chain-gap refusal = %q, want the chain rule", text)
	}
}

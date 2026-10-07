package gate

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	rr "github.com/gibbonmi/bench/internal/reviewrecord"
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
	if !strings.Contains(text, "plan commits land before the ticket merge, a default-branch merge lands before the first chunk or as a fold in the last chunk delta, and only record commits and comment-only corrections follow a chunk tip") {
		t.Fatalf("chain-gap refusal = %q, want the chain rule", text)
	}
}

func seedCommentGap(f *recordtest.Fixture) {
	f.Write("comment_gap_test.go", "package fixture\nvar value = 1 // before\n")
}

func correctCommentGap(f *recordtest.Fixture) {
	f.Write("comment_gap_test.go", "package fixture\nvar value = 1 // after\n")
	f.Commit("correct Go comment")
}

func TestReviewCheckpointCommentOnlyGap(t *testing.T) {
	for _, kind := range []string{"chunk", "complete", "reaffirm Standards"} {
		t.Run(kind, func(t *testing.T) {
			f := attachedCheckpointFixture(t, recordtest.Attach, seedCommentGap)
			frozen := f.Record.Chunks[0]
			if kind == "reaffirm Standards" {
				f.Record.Chunks[0].Reviews[0].Outcome = "fail"
				f.Record.Chunks[0].Reviews[0].FindingIDs = []string{"comment"}
				f.Save()
				f.Commit("retain comment finding")
				if code, out := runCheckpoint(t, f); code == 0 || !strings.Contains(out, "Standards") {
					t.Fatalf("unresolved finding accepted: %d %s", code, out)
				}
				if _, err := os.Stat(filepath.Join(f.Root, ".gate-run-count")); !os.IsNotExist(err) {
					t.Fatalf("oracle ran before finding refusal: %v", err)
				}
			}
			correctCommentGap(f)
			if kind == "reaffirm Standards" {
				entry, err := rr.RecordReview(f.Root, recordtest.Spec, rr.ReviewCall{Chunk: "1", Axis: "Standards", Evidence: f.Evidence("corrected-Standards", frozen.SourceDigest, "independent-review")})
				if err != nil {
					t.Fatal(err)
				}
				if entry.Base != frozen.Base || entry.Tip != frozen.Tip || entry.SourceDigest != frozen.SourceDigest {
					t.Fatalf("superseding review moved its frozen source: %+v", entry)
				}
			}
			args := []string{f.Root, "--checkpoint", recordtest.Spec, "--chunk", "1"}
			if kind == "complete" {
				f.Complete()
				digest, err := rr.SourceDigest(f.Root, f.Tree(), recordtest.Spec)
				if err != nil {
					t.Fatal(err)
				}
				f.Record.Completion.SourceDigest = digest
				f.Record.Completion.Verification = f.Verification("corrected-final", digest, f.Plan.FinalVerification)
				f.Save()
				args = []string{f.Root, "--checkpoint", recordtest.Spec, "--complete"}
			}
			var out, diagnostic bytes.Buffer
			if code := RunCommand(args, &out, &diagnostic); code != 0 {
				t.Fatalf("comment-only checkpoint: %d %s %s", code, &out, &diagnostic)
			}
			if got := outcomeRuns(t, f.Root); got != 1 {
				t.Fatalf("oracle runs = %d, want 1", got)
			}
			retained, err := rr.Read(f.Root, recordtest.Spec)
			if err != nil {
				t.Fatal(err)
			}
			got := retained.Chunks[0]
			if got.Base != frozen.Base || got.Tip != frozen.Tip || got.SourceDigest != frozen.SourceDigest || !reflect.DeepEqual(got.Reviews[1:3], frozen.Reviews[1:3]) {
				t.Fatalf("accepted gap changed frozen evidence: %+v", got)
			}
			wantReviews := 3
			if kind == "reaffirm Standards" {
				wantReviews++
			}
			if len(got.Reviews) != wantReviews {
				t.Fatalf("reviews = %d, want %d", len(got.Reviews), wantReviews)
			}
		})
	}
}

func TestReviewCheckpointKeepsStrictEvidence(t *testing.T) {
	for _, tc := range []struct{ name, reason string }{
		{"statement", "chunk 1: stale reviewed source: no chunk review covers"},
		{"dirty statement", "chunk 1: stale reviewed source"},
		{"forged digest", "chunk 1: stale source digest"},
		{"stale completion", "completion is incomplete or stale"},
		{"chain base", "expected base"},
		{"moved chunk", "chunk 1: stale Standards"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var f *recordtest.Fixture
			if tc.name == "chain base" {
				f = recordtest.Attach(t, outcomeFixture(t), 2)
				seedCommentGap(f)
				f.AddChunk()
				f.Save()
				f.Commit("retain first chunk")
			} else {
				f = attachedCheckpointFixture(t, recordtest.Attach, seedCommentGap)
			}
			if tc.name == "stale completion" {
				f.Complete()
				f.Save()
				f.Commit("retain completion")
			}
			correctCommentGap(f)
			args := []string{f.Root, "--checkpoint", recordtest.Spec, "--chunk", "1"}
			switch tc.name {
			case "statement", "dirty statement":
				f.Write("comment_gap_test.go", "package fixture\nvar value = 2 // after\n")
				if tc.name == "statement" {
					f.Commit("change Go statement")
				}
			case "forged digest", "moved chunk":
				digest, err := rr.SourceDigest(f.Root, f.Tree(), recordtest.Spec)
				if err != nil {
					t.Fatal(err)
				}
				f.Record.Chunks[0].SourceDigest = digest
				if tc.name == "moved chunk" {
					f.Record.Chunks[0].Tip = f.Tip()
					f.Record.Chunks[0].Verification = f.Verification("corrected-chunk", digest, f.Plan.Chunks[0].Verification)
				}
				f.Save()
			case "stale completion":
				args = []string{f.Root, "--checkpoint", recordtest.Spec, "--complete"}
			case "chain base":
				f.AddChunk()
				f.Save()
				f.Commit("retain second chunk")
				args = []string{f.Root, "--checkpoint", recordtest.Spec, "--chunk", "2"}
			}
			var out, diagnostic bytes.Buffer
			code := RunCommand(args, &out, &diagnostic)
			text := out.String() + diagnostic.String()
			if code == 0 || !strings.Contains(text, tc.reason) {
				t.Fatalf("strict evidence = %d %s, want %s", code, text, tc.reason)
			}
			if _, err := os.Stat(filepath.Join(f.Root, ".gate-run-count")); !os.IsNotExist(err) {
				t.Fatalf("oracle ran before evidence refusal: %v", err)
			}
		})
	}
}

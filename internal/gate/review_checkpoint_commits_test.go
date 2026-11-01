package gate

import (
	"bytes"
	"os"
	"reflect"
	"strings"

	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"testing"

	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// TestReviewCheckpointLaterSource refuses a source change after the reviewed tip. An
// uncommitted change refuses at the checkout, and a committed change refuses as stale.
func TestReviewCheckpointLaterSource(t *testing.T) {
	for _, committed := range []bool{false, true} {
		t.Run(map[bool]string{false: "dirty source", true: "unreviewed repair"}[committed], func(t *testing.T) {
			f := checkpointFixture(t)
			f.Write("source.txt", "later repair\n")
			reason := cleanCheckoutRefusal
			if committed {
				f.Commit("unreviewed repair")
				reason = "chunk 1: stale reviewed source"
			}
			if code, out := runCheckpoint(t, f); code == 0 || !strings.Contains(out, reason) {
				t.Fatalf("uncovered delta accepted: %d %s", code, out)
			}
		})
	}
}

// laterTicketFixture is one recorded chunk whose verification sits at the commit of its
// first ticket. A later ticket's commit, which also writes each path of plants, then
// holds the chunk tip and its reviews.
func laterTicketFixture(t *testing.T, plants ...string) *recordtest.Fixture {
	t.Helper()
	f := checkpointFixture(t)
	f.LaterTicket(plants...)
	f.Save()
	f.Commit("retain the later chunk tip")
	return f
}

// bindTo answers the entry change that binds an entry to commit through f.
func bindTo(f *recordtest.Fixture, commit string) func(*rr.Verification) {
	return func(entry *rr.Verification) { f.Bind(entry, commit) }
}

// refusesStaleEntries applies change to each chunk 1 entry of f and requires that the
// checkpoint refuses the stale entry before the oracle runs.
func refusesStaleEntries(t *testing.T, f *recordtest.Fixture, change func(*rr.Verification)) {
	t.Helper()
	for i := range f.Record.Chunks[0].Verification {
		change(&f.Record.Chunks[0].Verification[i])
	}
	resave(t, f, "retain the stale entries")
	if code, out := runCheckpoint(t, f); code == 0 || !strings.Contains(out, "chunk 1: verification tests is incomplete, failed, or stale") {
		t.Fatalf("checkpoint = %d %s; want the stale entry refused", code, out)
	}
	assertNoOracleRun(t, f)
}

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
// plan commits in the chunk delta and a default-branch merge in a recorded fold.
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
	resave(t, f, "move the second chunk base")
	var out, diagnostic bytes.Buffer
	code := RunCommand([]string{f.Root, "--checkpoint", recordtest.Spec, "--chunk", "2"}, &out, &diagnostic)
	text := out.String() + diagnostic.String()
	if code == 0 || !strings.Contains(text, "expected base "+expected) {
		t.Fatalf("chain-gap refusal = (%d, %q), want expected base %s", code, text, expected)
	}
	if !strings.Contains(text, "plan commits land before the ticket merge, a default-branch merge lands before the first chunk or as a recorded fold, and only record commits, recorded folds, and proved comment-only Go corrections follow a chunk tip until its checkpoint") {
		t.Fatalf("chain-gap refusal = %q, want the chain rule", text)
	}
}

// chainedOverCommentGap is two recorded chunks whose chunk 2 base follows the chunk 1 tip by a
// Go comment edit. It also answers the chunk 1 entry as the record first retained it.
func chainedOverCommentGap(t *testing.T) (*recordtest.Fixture, rr.Chunk) {
	t.Helper()
	f := recordtest.Attach(t, outcomeFixture(t), 2)
	seedCommentGap(f)
	f.AddChunk()
	f.Save()
	f.Commit("retain first chunk")
	frozen, err := rr.Read(f.Root, recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	correctCommentGap(f)
	f.AddChunk()
	f.Save()
	f.Commit("retain second chunk")
	return f, frozen.Chunks[0]
}

// TestReviewCheckpointChainAcceptsACommentGap covers CG01 to CG03: the chunk 2 checkpoint and
// the complete checkpoint exit 0 over a comment gap before the chunk 2 base, and the chunk 1
// entry keeps its frozen evidence.
func TestReviewCheckpointChainAcceptsACommentGap(t *testing.T) {
	for _, kind := range []string{"chunk", "complete"} {
		t.Run(kind, func(t *testing.T) {
			f, frozen := chainedOverCommentGap(t)
			args := []string{f.Root, "--checkpoint", recordtest.Spec, "--chunk", "2"}
			if kind == "complete" {
				f.Complete()
				f.Save()
				f.Commit("retain completion")
				args = []string{f.Root, "--checkpoint", recordtest.Spec, "--complete"}
			}
			var out, diagnostic bytes.Buffer
			if code := RunCommand(args, &out, &diagnostic); code != 0 {
				t.Fatalf("%s checkpoint over a comment gap = %d %s %s, want 0", kind, code, &out, &diagnostic)
			}
			retained, err := rr.Read(f.Root, recordtest.Spec)
			if err != nil {
				t.Fatal(err)
			}
			got := retained.Chunks[0]
			if got.Base != frozen.Base || got.Tip != frozen.Tip || got.SourceDigest != frozen.SourceDigest || !reflect.DeepEqual(got.Reviews, frozen.Reviews) || !reflect.DeepEqual(got.Verification, frozen.Verification) {
				t.Fatalf("chunk 1 after the comment gap = %+v, want the frozen entry %+v", got, frozen)
			}
		})
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
				if _, err := os.Stat(runCountWitness.path(t, f.Root)); !os.IsNotExist(err) {
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
				f.Commit("retain the superseding review")
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
				f.Commit("retain corrected completion")
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

// TestReviewCheckpointKeepsStrictEvidence refuses each case with its reason. The moved chunk
// holds its axis results at the earlier pair across a test-only delta, so it also routes to
// the confirmation form.
func TestReviewCheckpointKeepsStrictEvidence(t *testing.T) {
	for _, tc := range []struct{ name, reason, next string }{
		{"statement", "chunk 1: stale reviewed source: no chunk review covers", ""},
		{"dirty statement", cleanCheckoutRefusal, ""},
		{"forged digest", "chunk 1: stale source digest", ""},
		{"stale completion", "completion is incomplete or stale", ""},
		{"chain base", "expected base", ""},
		{"moved chunk", "chunk 1: stale Standards", "next=bench record 'confirmation' --help\n"},
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
				resave(t, f, "retain the "+tc.name)
			case "stale completion":
				args = []string{f.Root, "--checkpoint", recordtest.Spec, "--complete"}
			case "chain base":
				f.Write("comment_gap_test.go", "package fixture\nvar value = 2 // after\n")
				f.Commit("change Go statement")
				f.AddChunk()
				f.Save()
				f.Commit("retain second chunk")
				args = []string{f.Root, "--checkpoint", recordtest.Spec, "--chunk", "2"}
			}
			var out, diagnostic bytes.Buffer
			code := RunCommand(args, &out, &diagnostic)
			text := out.String() + diagnostic.String()
			if code == 0 || !strings.Contains(text, tc.reason) || !strings.Contains(text, tc.next) {
				t.Fatalf("strict evidence = %d %s, want %s and the route %q", code, text, tc.reason, tc.next)
			}
			if _, err := os.Stat(runCountWitness.path(t, f.Root)); !os.IsNotExist(err) {
				t.Fatalf("oracle ran before evidence refusal: %v", err)
			}
		})
	}
}

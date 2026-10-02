package recordcmd_test

import (
	"reflect"
	"strings"
	"testing"

	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordcmd"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/toon"
)

// reviewArgs is the review argv for one axis result on chunk 1 with no finding. Each
// finding follows as one --finding flag, in order.
func reviewArgs(t *testing.T, axis, id string, findings []string, edits ...map[string]string) []string {
	t.Helper()
	args := formArgs(t, "review", map[string]string{"--chunk": "1", "--axis": axis, "--id": id, "--performer": "fixture-reviewer-" + id,
		"--model": "unknown", "--effort": "unknown", "--ref": "fixture:terminal-result", "--excerpt": excerptFile(t, "excerpt.txt", "no findings\n")}, edits...)
	for _, finding := range findings {
		args = append(args, "--finding", finding)
	}
	return args
}

// review runs the review form and requires exit 0.
func review(t *testing.T, f *recordtest.Fixture, axis, id string, findings ...string) string {
	t.Helper()
	return succeed(t, f, reviewArgs(t, axis, id, findings))
}

// reviewed is the review result id in the chunk 1 list.
func reviewed(t *testing.T, f *recordtest.Fixture, id string) rr.Review {
	t.Helper()
	for _, result := range onlyChunk(t, f).Reviews {
		if result.ID == id {
			return result
		}
	}
	t.Fatalf("chunk 1 holds no review result %s", id)
	return rr.Review{}
}

func TestRecordReviewCopiesTheChunkPair(t *testing.T) {
	f := recorded(t, 1)
	chunk := onlyChunk(t, f)
	advance(f, "later")
	review(t, f, "Standards", "standards-1")
	got := reviewed(t, f, "standards-1")
	if got.Base != chunk.Base || got.Tip != chunk.Tip || got.SourceDigest != chunk.SourceDigest || got.Tip == f.Tip() || got.SourceDigest == sourceDigest(t, f) {
		t.Fatalf("review pair = %s..%s at %s, want the chunk pair %s..%s at %s, not HEAD", got.Base, got.Tip, got.SourceDigest, chunk.Base, chunk.Tip, chunk.SourceDigest)
	}
}

func TestRecordReviewFirstResultSupersedesNothing(t *testing.T) {
	f := recorded(t, 1)
	review(t, f, "Spec", "spec-1")
	review(t, f, "Standards", "standards-1")
	if got := reviewed(t, f, "standards-1").Supersedes; len(got) != 0 {
		t.Fatalf("supersedes = %q, want none", got)
	}
}

func TestRecordReviewSupersedesTheSameAxis(t *testing.T) {
	f := recorded(t, 1)
	review(t, f, "Standards", "standards-1", "R1")
	review(t, f, "Spec", "spec-1")
	review(t, f, "Standards", "standards-2")
	review(t, f, "Standards", "standards-3")
	for id, want := range map[string]string{"standards-2": "standards-1", "standards-3": "standards-2"} {
		if got := reviewed(t, f, id).Supersedes; !reflect.DeepEqual(got, []string{want}) {
			t.Fatalf("%s supersedes = %q, want [%s]", id, got, want)
		}
	}
}

func TestRecordReviewWithoutFindingsPasses(t *testing.T) {
	f := recorded(t, 1)
	review(t, f, "Standards", "standards-1")
	if got := reviewed(t, f, "standards-1"); got.Outcome != "pass" || got.State != "completed" || len(got.FindingIDs) != 0 {
		t.Fatalf("result = %s %s with findings %q, want completed pass with none", got.State, got.Outcome, got.FindingIDs)
	}
}

func TestRecordReviewWithFindingsFails(t *testing.T) {
	f := recorded(t, 1)
	for id, findings := range map[string][]string{"standards-1": {"R1", "R2"}, "spec-1": {"R2", "R1"}} {
		review(t, f, map[string]string{"standards-1": "Standards", "spec-1": "Spec"}[id], id, findings...)
		if got := reviewed(t, f, id); got.Outcome != "fail" || got.State != "completed" || !reflect.DeepEqual(got.FindingIDs, findings) {
			t.Fatalf("result = %s %s with findings %q, want completed fail with %q", got.State, got.Outcome, got.FindingIDs, findings)
		}
	}
}

func TestRecordReviewLandsInTheReviewList(t *testing.T) {
	f := recorded(t, 1)
	review(t, f, "Coverage", "coverage-1")
	chunk := onlyChunk(t, f)
	if len(chunk.Reviews) != 1 || len(chunk.Verification) != 0 || chunk.Reviews[0].Role != "independent-review" {
		t.Fatalf("chunk = %d verification and %d review results, want the one independent-review result coverage-1", len(chunk.Verification), len(chunk.Reviews))
	}
}

func TestRecordReviewNamesTheChunkForm(t *testing.T) {
	f := recorded(t, 2)
	refuseArgs(t, f, "bench record chunk", reviewArgs(t, "Standards", "standards-1", nil, map[string]string{"--chunk": "2"}))
}

func TestRecordReviewRefusesTheImplementationSession(t *testing.T) {
	f, _ := recordtest.NewLinked(t, 1)
	base, _ := advance(f, "chunk 1")
	f.RecordChunk(base)
	f.Save()
	refuseArgs(t, f, "invalid independent review axis or performer",
		reviewArgs(t, "Standards", "standards-1", nil, map[string]string{"--performer": f.Record.ImplementationSession}))
}

func TestRecordReviewRefusesADuplicateID(t *testing.T) {
	f := recorded(t, 1)
	review(t, f, "Standards", "standards-1")
	refuseArgs(t, f, "standards-1", reviewArgs(t, "Spec", "standards-1", nil))
}

func TestRecordReviewReportsItsRow(t *testing.T) {
	f := recorded(t, 1)
	review(t, f, "Standards", "standards-1")
	out := review(t, f, "Standards", "standards-2", "R1")
	got := reviewed(t, f, "standards-2")
	want, err := toon.Table("review", []string{"chunk", "id", "axis", "outcome", "supersedes", "source_digest", "excerpt_digest"},
		[][]string{{"1", got.ID, got.Axis, got.Outcome, "standards-1", got.SourceDigest, got.NativeRef.Digest}})
	if err != nil {
		t.Fatal(err)
	}
	if out != want {
		t.Fatalf("output = %q, want %q", out, want)
	}
}

func TestRecordedChunkPassesTheCheckpoint(t *testing.T) {
	f := recorded(t, 1)
	verify(t, f)
	verify(t, f, unprobed)
	for _, axis := range rr.Axes() {
		review(t, f, axis, "review-"+axis)
	}
	f.Commit("record chunk 1")
	if err := rr.Check(f.Root, f.Tree(), f.Tip(), recordtest.Spec, "1", false); err != nil {
		t.Fatalf("checkpoint for chunk 1: %v", err)
	}
}

func TestRecordReviewGrammarRefusals(t *testing.T) {
	f := recorded(t, 1)
	for name, edit := range map[string]map[string]string{
		"axis outside":    {"--axis": "Security"},
		"missing axis":    {"--axis": drop},
		"missing excerpt": {"--excerpt": drop},
	} {
		out, code := recordcmd.Command(f.Root, reviewArgs(t, "Standards", "standards-1", nil, edit))
		if code != 2 || !strings.HasPrefix(out, "usage: bench record review") {
			t.Errorf("%s = exit %d, output %q; want exit 2 with a usage: bench record review line", name, code, out)
		}
	}
}

func TestRecordReviewRefusesAControlCharacterInAFlag(t *testing.T) {
	f := recorded(t, 1)
	for _, name := range []string{"--chunk", "--id", "--performer", "--model", "--effort", "--ref", "--finding", "--excerpt"} {
		edit, findings := map[string]string{name: "value\x1b"}, []string(nil)
		if name == "--finding" {
			edit, findings = nil, []string{"R1\x1b", "R2"}
		}
		out := refuseArgs(t, f, name+" holds a control character", reviewArgs(t, "Standards", "standards-1", findings, edit))
		if strings.Contains(out, "\x1b") {
			t.Errorf("%s refusal = %q, want no control character", name, out)
		}
	}
}

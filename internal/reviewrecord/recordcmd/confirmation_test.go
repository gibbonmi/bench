package recordcmd_test

import (
	"reflect"
	"strings"
	"testing"

	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/toon"
)

// mutation is the caller's mutation text of each confirmation call below.
const mutation = "omit the repaired assertion"

// confirmArgs is the confirmation argv for chunk 1 that imports one biting probe answer. Each
// finding follows as one --finding flag, in order.
func confirmArgs(t *testing.T, id string, findings ...string) []string {
	t.Helper()
	args := formArgs(t, "confirmation", map[string]string{"--chunk": "1", "--id": id, "--performer": recordtest.Orchestrator, "--model": "unknown", "--effort": "unknown",
		"--exit-code": "0", "--ref": "fixture:terminal-result", "--probe-result": excerptFile(t, "probe.out", answer("bit", "yes", 0, "1", report)), "--mutation": mutation})
	for _, finding := range findings {
		args = append(args, "--finding", finding)
	}
	return args
}

// repaired is a recorded chunk 1 whose Standards result holds F-1 and whose Spec and Coverage
// results pass at the reviewed pair. It then commits a repair that writes path, and the chunk
// entry moves to that repair. It answers the fixture and the reviewed tip.
func repaired(t *testing.T, path string) (*recordtest.Fixture, string) {
	t.Helper()
	f := recorded(t, 1)
	chunk := onlyChunk(t, f)
	review(t, f, "Standards", "standards-1", "F-1")
	review(t, f, "Spec", "spec-1")
	review(t, f, "Coverage", "coverage-1")
	repairAt(t, f, path, chunk.Base)
	return f, chunk.Tip
}

// repairAt commits a new version of path and moves the chunk 1 entry from base to that commit.
func repairAt(t *testing.T, f *recordtest.Fixture, path, base string) {
	t.Helper()
	f.Write(path, "package fixture\n\n// "+f.Tip()+"\n")
	f.Commit("repair " + path)
	record(t, f, "1", base, f.Tip())
}

// lastConfirmation is the last confirmation of chunk 1.
func lastConfirmation(t *testing.T, f *recordtest.Fixture) rr.Confirmation {
	t.Helper()
	confirmations := onlyChunk(t, f).Confirmations
	if len(confirmations) == 0 {
		t.Fatal("chunk 1 holds no confirmation")
	}
	return confirmations[len(confirmations)-1]
}

// TestRecordConfirmationWritesTheDerivedPair covers RC38: the entry holds the reviewed tip of
// the axis results, the chunk pair as its repaired pair, the imported probe, and the caller's
// mutation text, and the answer reports the digest of the reviewed tip's tree.
func TestRecordConfirmationWritesTheDerivedPair(t *testing.T) {
	f, reviewedTip := repaired(t, "repair_test.go")
	chunk := onlyChunk(t, f)
	out := succeed(t, f, confirmArgs(t, "confirm-1", "F-1"))
	got := lastConfirmation(t, f)
	data := answer("bit", "yes", 0, "1", report)
	want := rr.Confirmation{Evidence: rr.Evidence{ID: "confirm-1", Performer: recordtest.Orchestrator, Role: rr.ConfirmationRole, Model: "unknown", Effort: "unknown",
		SourceDigest: chunk.SourceDigest, State: "completed", Outcome: "pass", NativeRef: rr.NativeRef{Ref: "fixture:terminal-result", Digest: rr.Digest([]byte(data)), Excerpt: data}},
		Base: chunk.Base, ReviewedTip: reviewedTip, RepairedTip: chunk.Tip, FindingIDs: []string{"F-1"}, Supersedes: []string{}}
	want.Probe = rr.Probe{Mutation: mutation, Outcome: "bit", ExitCode: 1, Restore: "pass", NativeRef: want.NativeRef}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("confirmation = %+v, want %+v", got, want)
	}
	reviewedDigest, err := rr.SourceDigest(f.Root, f.Git("rev-parse", reviewedTip+"^{tree}"), recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	row, err := toon.Table("confirmation", []string{"chunk", "id", "reviewed_tip", "reviewed_digest", "repaired_tip", "source_digest", "supersedes", "excerpt_digest"},
		[][]string{{"1", "confirm-1", reviewedTip, reviewedDigest, chunk.Tip, chunk.SourceDigest, "", want.NativeRef.Digest}})
	if err != nil {
		t.Fatal(err)
	}
	if out != row {
		t.Fatalf("output = %q, want %q", out, row)
	}
}

// TestRecordConfirmationRefusals covers RC39, RC40, and RC41: each refusal leaves the record
// bytes unchanged.
func TestRecordConfirmationRefusals(t *testing.T) {
	for _, tc := range []struct {
		name, path, want string
		findings         []string
		change           func(*testing.T, *recordtest.Fixture)
	}{
		{"different pairs", "repair_test.go", "chunk 1: mixed review pairs", []string{"F-1"}, func(t *testing.T, f *recordtest.Fixture) { review(t, f, "Coverage", "coverage-2") }},
		{"chunk pair", "repair_test.go", "chunk 1: the axis results sit at the chunk pair", nil, func(t *testing.T, f *recordtest.Fixture) {
			review(t, f, "Standards", "standards-2")
			review(t, f, "Spec", "spec-2")
			review(t, f, "Coverage", "coverage-2")
		}},
		{"non-test delta", "repair.go", "the delta to the chunk tip is not test-only: repair.go: path is not a Go test file", []string{"F-1"}, nil},
		{"omitted finding", "repair_test.go", "chunk 1: confirmation confirm-1 closes the findings [], and the open findings at the reviewed pair are [F-1]", nil, nil},
		{"extra finding", "repair_test.go", "closes the findings [F-1 F-2], and the open findings at the reviewed pair are [F-1]", []string{"F-1", "F-2"}, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, _ := repaired(t, tc.path)
			if tc.change != nil {
				tc.change(t, f)
			}
			refuseArgs(t, f, tc.want, confirmArgs(t, "confirm-1", tc.findings...))
		})
	}
}

// TestRecordConfirmationSupersedesThePrior covers RC42: a confirmation after a second
// test-only repair names the first confirmation in its supersession link.
func TestRecordConfirmationSupersedesThePrior(t *testing.T) {
	f, _ := repaired(t, "repair_test.go")
	succeed(t, f, confirmArgs(t, "confirm-1", "F-1"))
	repairAt(t, f, "repair_test.go", onlyChunk(t, f).Base)
	out := succeed(t, f, confirmArgs(t, "confirm-2", "F-1"))
	if got := lastConfirmation(t, f); got.ID != "confirm-2" || !reflect.DeepEqual(got.Supersedes, []string{"confirm-1"}) || !strings.Contains(out, ",confirm-1,") {
		t.Fatalf("second confirmation %s supersedes %q with output %q, want confirm-2 superseding confirm-1", got.ID, got.Supersedes, out)
	}
}

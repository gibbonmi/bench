package reviewrecord

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReviewRecord(t *testing.T) {
	t.Run("missing clean axis", func(t *testing.T) {
		if err := CheckReviews(Chunk{ID: "1"}, []string{"author"}, false); err == nil || !strings.Contains(err.Error(), "missing Standards") {
			t.Fatalf("missing clean review result: %v", err)
		}
	})
}

// TestCarriedReviewKeepsTheRoleCondition grades a carried Standards result whose role is not
// the independent review role. The parser refuses that role in every record, so the
// checkpoint never reaches this condition of the carry branch, and this test calls it.
func TestCarriedReviewKeepsTheRoleCondition(t *testing.T) {
	chunk := Chunk{ID: "1", Base: "base", Tip: "reviewed", SourceDigest: "reviewed-digest"}
	for _, axis := range Axes() {
		chunk.Reviews = append(chunk.Reviews, Review{Evidence: Evidence{ID: axis, Performer: "reviewer-" + axis, Role: "independent-review", SourceDigest: chunk.SourceDigest, State: "completed", Outcome: "pass"}, Axis: axis, Base: chunk.Base, Tip: chunk.Tip})
	}
	standards := &chunk.Reviews[0]
	standards.Role, standards.Outcome, standards.FindingIDs = "author-verification", "fail", []string{"F-1"}
	pair := chunk.reviewed()
	pair.carried = true
	want := `chunk 1: invalid Standards performer "reviewer-Standards" or role "author-verification"`
	if err := checkReviews(pair, []string{"author"}, 3); err == nil || !strings.Contains(err.Error(), want) {
		t.Fatalf("carried result with a wrong role: %v; want %q", err, want)
	}
}

// TestProbeChecksAgree grades the planned-probe check and the confirmation check on each
// probe: the two checks give one verdict, because they call one probe-bite predicate.
func TestProbeChecksAgree(t *testing.T) {
	zero := 0
	requirement := Requirement{ID: "tests", Command: "go test ./...", Probe: "omit the check"}
	record := Record{Version: 1, ImplementationSession: "author"}
	chunk := Chunk{ID: "1", Base: "base", Tip: "repaired", SourceDigest: "repaired-digest"}
	for i, probe := range []Probe{
		{Outcome: "bit", ExitCode: 1, Restore: "pass"},
		{Outcome: "silent", ExitCode: 1, Restore: "pass"},
		{Outcome: "bit", ExitCode: 0, Restore: "pass"},
		{Outcome: "bit", ExitCode: 1, Restore: "fail"},
	} {
		probe.Mutation = requirement.Probe
		verification := Verification{Evidence: Evidence{Performer: "author", Role: "author-verification", SourceDigest: chunk.SourceDigest, State: "completed", Outcome: "pass"}, Requirement: requirement.ID, Command: requirement.Command, ExitCode: &zero, Probe: &probe}
		planned := checkVerification([]Verification{verification}, []Requirement{requirement}, Plan{}, record, chunk.SourceDigest, Repair{Chunk: chunk.ID})
		chunk.Confirmations = []Confirmation{{Evidence: Evidence{Performer: "reconciler", SourceDigest: chunk.SourceDigest, State: "completed"}, Base: chunk.Base, ReviewedTip: "reviewed", RepairedTip: chunk.Tip, Probe: probe}}
		confirmed := confirms(chunk, ReviewedPair{Tip: "reviewed"}, "reconciler")
		if (planned == nil) != (confirmed == nil) || (i == 0 && planned != nil) {
			t.Errorf("probe %+v: planned-probe check %v, confirmation check %v; want one verdict, and a pass for the probe that bit", probe, planned, confirmed)
		}
	}
}

func TestReviewRecordSchema(t *testing.T) {
	data, _ := json.Marshal(Record{Version: 99})
	if _, err := Parse(data); err == nil || !strings.Contains(err.Error(), "unsupported version") {
		t.Fatalf("unsupported version reached evidence validation: %v", err)
	}
}

const renderSpec = "specs/example/spec.md"

// renderRecord is a valid record whose one completed verification embeds
// excerpt, so a test can read the rendered document back through Read.
func renderRecord(excerpt string) Record {
	exit := 0
	evidence := Evidence{ID: "v1", Performer: "author", Role: "author-verification", Model: "unknown", Effort: "unknown", SourceDigest: strings.Repeat("a", 40), State: "completed", Outcome: "pass", NativeRef: NativeRef{Ref: "fixture:result", Excerpt: excerpt, Digest: Digest([]byte(excerpt))}}
	return Record{Version: 1, Spec: renderSpec, PlanDigest: "sha256:plan", ImplementationSession: "author", Completion: Completion{Verification: []Verification{{Evidence: evidence, Requirement: "focused", Command: "bench test", ExitCode: &exit}}}}
}

// render renders record over document, then reads the result back through
// Read, so each test grades the document that a reader sees.
func render(t *testing.T, document []byte, record Record) (string, Record) {
	t.Helper()
	out, err := Render(document, record)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	path, err := RecordPath(renderSpec)
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Dir(filepath.Join(root, path)), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, path), out, 0o644); err != nil {
		t.Fatal(err)
	}
	loaded, err := Read(root, renderSpec)
	if err != nil {
		t.Fatalf("Read of the rendered document: %v\n%s", err, out)
	}
	return string(out), loaded
}

func TestRenderKeepsProseAroundTheFence(t *testing.T) {
	before := "# Review outcomes\n\nThe author notes.\n\n```bench-review-record\n"
	after := "```\n\n## Later notes\n\nKept.\n"
	out, _ := render(t, []byte(before+"{\"version\": 0}\n"+after), renderRecord("ok"))
	if !strings.HasPrefix(out, before) || !strings.HasSuffix(out, after) || strings.Contains(out, `"version": 0`) {
		t.Fatalf("Render changed bytes outside the payload lines:\n%s", out)
	}
}

func TestRenderAppendsAFenceAfterUnterminatedProse(t *testing.T) {
	prose := "# Review outcomes\n\nNo fence yet."
	out, _ := render(t, []byte(prose), renderRecord("ok"))
	if !strings.HasPrefix(out, prose+"\n\n```bench-review-record\n") {
		t.Fatalf("Render did not put the fence on its own line after the prose:\n%s", out)
	}
}

func TestRenderStartsANewDocument(t *testing.T) {
	out, _ := render(t, nil, renderRecord("ok"))
	if !strings.HasPrefix(out, "# Review outcomes\n\n```bench-review-record\n") {
		t.Fatalf("Render of a nil document lacks the heading:\n%s", out)
	}
}

func TestRenderRoundTripsAHostileExcerpt(t *testing.T) {
	excerpt := "a\tb\rc\x1bd\ne <&> f"
	_, loaded := render(t, nil, renderRecord(excerpt))
	if got := loaded.Completion.Verification[0].NativeRef.Excerpt; got != excerpt {
		t.Fatalf("excerpt = %q, want %q", got, excerpt)
	}
}

func TestRenderWritesAReadablePayload(t *testing.T) {
	out, _ := render(t, nil, renderRecord("fmt & vet"))
	if !strings.Contains(out, "\n  \"version\": 1,\n") || !strings.Contains(out, `"excerpt": "fmt & vet"`) {
		t.Fatalf("payload is not two-space indented with a raw &:\n%s", out)
	}
}

// The malformed documents below are built by hand on purpose: no encoder can
// produce them.
func TestRenderRefusesAnUnterminatedFence(t *testing.T) {
	if out, err := Render([]byte("# Review outcomes\n\n```bench-review-record\n{\n"), renderRecord("ok")); err == nil {
		t.Fatalf("Render accepted an unterminated fence:\n%s", out)
	}
}

func TestRenderRefusesADuplicateFence(t *testing.T) {
	fence := "```bench-review-record\n{}\n```\n"
	if out, err := Render([]byte("# Review outcomes\n\n"+fence+"\n"+fence), renderRecord("ok")); err == nil {
		t.Fatalf("Render accepted a duplicate fence:\n%s", out)
	}
}

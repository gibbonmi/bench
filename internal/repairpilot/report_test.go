package repairpilot

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/assessment"
)

func TestRepairPilotReport(t *testing.T) {
	t.Run("stalled", func(t *testing.T) {
		requireFixtureCase(t)
		document := reportFixture()
		document.Observations = []observation{reportRepair("repair-stalled", "stalled")}
		document.Audits = []audit{reportAudit("audit-stalled", "repair-stalled", "stalled")}
		out := runReport(t, document, false)
		assertReportClass(t, out, "stalled", "present")
		assertReportClass(t, out, "productive", "missing")
		unaudited := reportFixture()
		unaudited.Observations = []observation{reportRepair("unaudited-stalled", "stalled")}
		assertReportClass(t, runReport(t, unaudited, false), "stalled", "missing")
	})

	t.Run("productive", func(t *testing.T) {
		requireFixtureCase(t)
		document := reportFixture()
		key := sequenceKey("source-productive", "spec-a", "chunk-a")
		first := reportFailure("productive-first", "assignment-a", "", "")
		first.Sequence = key
		repeated := reportFailure("productive-repeat", "assignment-a", "", "")
		repeated.Sequence = key
		repeated.ObservedAt = parseTime("2026-09-03T11:00:00Z")
		repair := reportRepair("repair-productive", "productive")
		repair.Sequence = key
		document.Observations = []observation{first, repeated, repair}
		document.Audits = []audit{reportAudit("audit-productive", "repair-productive", "productive")}
		out := runReport(t, document, false)
		assertReportClass(t, out, "productive", "present")
		assertReportValue(t, out, "comparable_repeats", "1")
	})

	t.Run("unchanged", func(t *testing.T) {
		requireFixtureCase(t)
		document := reportFixture()
		document.Observations = []observation{reportRerun("rerun-unchanged")}
		out := runReport(t, document, false)
		assertReportClass(t, out, "unchanged", "present")
		assertReportClass(t, out, "stalled", "missing")
		assertReportClass(t, out, "productive", "missing")
		invalid := reportRerun("rerun-changed")
		invalid.Rerun.UnchangedContent = false
		document.Observations = []observation{invalid}
		assertReportClass(t, runReport(t, document, false), "unchanged", "missing")
	})

	t.Run("overlap", func(t *testing.T) {
		requireFixtureCase(t)
		document := reportFixture()
		first := reportFailure("overlap-a", "assignment-a", "2026-09-02T09:00:00Z", "2026-09-02T11:00:00Z")
		second := reportFailure("overlap-b", "assignment-b", "2026-09-02T10:00:00Z", "2026-09-02T12:00:00Z")
		second.Sequence = first.Sequence
		document.Observations = []observation{first, second}
		out := runReport(t, document, false)
		assertReportClass(t, out, "overlap", "present")
		second.AssignmentID = first.AssignmentID
		document.Observations = []observation{first, second}
		assertReportClass(t, runReport(t, document, false), "overlap", "missing")
		second.AssignmentID = "assignment-b"
		second.Sequence = sequenceKey("other-source", "spec-a", "chunk-a")
		document.Observations = []observation{first, second}
		assertReportClass(t, runReport(t, document, false), "overlap", "missing")
	})

	t.Run("incomplete", func(t *testing.T) {
		requireFixtureCase(t)
		document := reportFixture()
		complete := reportFailure("complete", "assignment-a", "", "")
		complete.Endpoint = &endpoint{Kind: "reviewer-handoff", Reference: *nativeReference("native:endpoint")}
		incomplete := reportFailure("incomplete", "assignment-b", "", "")
		incomplete.Sequence = sequenceKey("source-b", "spec-a", "chunk-b")
		document.Observations = []observation{complete, incomplete}
		out := runReport(t, document, true)
		assertReportValue(t, out, "completed_sequences", "1")
		assertReportValue(t, out, "incomplete_sequences", "1")
		rows := reportRows(t, out, "sequences")
		if !reflect.DeepEqual(rowFields(rows, "source", "state"), [][]string{{"source-b", "incomplete"}, {"source-complete", "complete"}}) {
			t.Fatalf("sequence rows = %#v", rows)
		}
	})

	t.Run("missing-class", func(t *testing.T) {
		requireFixtureCase(t)
		for _, class := range requiredReportClasses {
			document := fixtureMissingClass(class)
			out := runReport(t, document, false)
			assertReportValue(t, out, "sample", "inconclusive")
			assertReportClass(t, out, class, "missing")
		}
		complete := runReport(t, completeReportFixture(), false)
		assertReportValue(t, complete, "sample", "coverage-complete")
	})

	t.Run("full", func(t *testing.T) {
		requireFixtureCase(t)
		document := completeReportFixture()
		unknown := reportFailure("unknown-evidence", "assignment-c", "2026-09-02T09:30:00Z", "")
		document.Observations = append(document.Observations, unknown)
		document.Observations[2].Endpoint = &endpoint{
			Kind: "reviewer-handoff", Reference: *nativeReference("native:full-endpoint"),
			CoveredBlockerRefs: []assessment.Reference{*nativeReference("native:covered-blocker")},
		}
		out := runReport(t, document, true)
		for _, table := range []string{"sequences", "observations", "failures", "observation_refs", "repairs", "reruns", "endpoints", "audits", "audit_refs", "intervals", "interval_comparisons", "evidence_gaps"} {
			reportRows(t, out, table)
		}
		for table, count := range map[string]int{
			"sequences": 5, "observations": 6, "failures": 3, "observation_refs": 6, "repairs": 2, "reruns": 1,
			"endpoints": 1, "endpoint_blocker_refs": 1, "audits": 2, "audit_refs": 2, "intervals": 3, "interval_comparisons": 1, "evidence_gaps": 1,
		} {
			if rows := reportRows(t, out, table); len(rows) != count {
				t.Errorf("%s rows = %d, want %d", table, len(rows), count)
			}
		}
		for _, value := range []string{"the target guard is missing", "add the target guard", "reviewed native evidence", "native:verification-repair-productive", "unknown-evidence", "partial-interval", "reviewer-handoff"} {
			if !strings.Contains(out, value) {
				t.Errorf("full report missing retained value %q", value)
			}
		}
		assertRowsContain(t, out, "sequences", map[string]string{"source": "source-unknown-evidence", "state": "incomplete"})
		assertRowsContain(t, out, "observations", map[string]string{"id": "repair-productive", "proposal_label": "productive", "effective_label": "productive"})
		assertRowsContain(t, out, "audits", map[string]string{"audit_id": "audit-stalled", "conclusion": "stalled", "verification_native": "native:verification-audit-stalled"})
		assertRowsContain(t, out, "observation_refs", map[string]string{"observation_id": "unknown-evidence", "native": "native:observation-unknown-evidence"})
		assertRowsContain(t, out, "endpoint_blocker_refs", map[string]string{"observation_id": "rerun-unchanged", "native": "native:covered-blocker"})
		assertRowsContain(t, out, "evidence_gaps", map[string]string{"observation_id": "unknown-evidence", "kind": "partial-interval"})
		if got := rowField(reportRows(t, out, "observations"), "id"); !reflect.DeepEqual(got, []string{"overlap-a", "overlap-b", "repair-productive", "repair-stalled", "rerun-unchanged", "unknown-evidence"}) {
			t.Fatalf("full observation inventory = %v", got)
		}
	})

	t.Run("default", func(t *testing.T) {
		requireFixtureCase(t)
		document := completeReportFixture()
		document.Observations[2].Endpoint = &endpoint{Kind: "reviewer-handoff", Reference: *nativeReference("native:default-endpoint")}
		partial := reportFailure("default-gap", "assignment-c", "2026-09-02T09:30:00Z", "")
		document.Observations = append(document.Observations, partial)
		out := runReport(t, document, false)
		if strings.Contains(out, "observations[") || !strings.Contains(out, "required_classes[") {
			t.Fatalf("default report has wrong detail shape:\n%s", out)
		}
		for field, value := range map[string]string{
			"state": "stopped", "activated_at": "2026-09-02T12:00:00Z", "collection_ends_at": "2026-09-16T12:00:00Z", "cutoff_at": "2026-09-09T12:00:00Z",
			"sample": "coverage-complete", "completed_sequences": "1", "incomplete_sequences": "4", "unknown_labels": "4", "evidence_gaps": "1",
		} {
			assertReportValue(t, out, field, value)
		}
		for _, class := range requiredReportClasses {
			assertReportClass(t, out, class, "present")
		}
	})

	t.Run("provisional", func(t *testing.T) {
		requireFixtureCase(t)
		document := completeReportFixture()
		document.CutoffAt = nil
		out := runReportAt(t, document, false, parseTime("2026-09-03T12:00:00Z"))
		assertReportValue(t, out, "sample", "provisional")
	})

	t.Run("partial-interval", func(t *testing.T) {
		requireFixtureCase(t)
		document := reportFixture()
		first := reportFailure("partial-a", "assignment-a", "2026-09-02T09:00:00Z", "")
		second := reportFailure("partial-b", "assignment-b", "2026-09-02T10:00:00Z", "2026-09-02T12:00:00Z")
		second.Sequence = first.Sequence
		document.Observations = []observation{first, second}
		out := runReport(t, document, false)
		assertReportClass(t, out, "overlap", "missing")
		assertReportValue(t, out, "evidence_gaps", "1")
	})

	t.Run("interval-edge", func(t *testing.T) {
		requireFixtureCase(t)
		for _, pair := range [][4]string{
			{"2026-09-02T09:00:00Z", "2026-09-02T10:00:00Z", "2026-09-02T10:00:00Z", "2026-09-02T11:00:00Z"},
			{"2026-09-02T09:00:00Z", "2026-09-02T09:00:00Z", "2026-09-02T08:00:00Z", "2026-09-02T10:00:00Z"},
		} {
			document := reportFixture()
			first := reportFailure("edge-a", "assignment-a", pair[0], pair[1])
			second := reportFailure("edge-b", "assignment-b", pair[2], pair[3])
			second.Sequence = first.Sequence
			document.Observations = []observation{first, second}
			assertReportClass(t, runReport(t, document, false), "overlap", "missing")
		}
	})

	t.Run("interval-provenance", func(t *testing.T) {
		requireFixtureCase(t)
		document := reportFixture()
		first := reportFailure("unproven-a", "assignment-a", "2026-09-02T09:00:00Z", "2026-09-02T11:00:00Z")
		first.IntervalReference = nil
		second := reportFailure("unproven-b", "assignment-b", "2026-09-02T10:00:00Z", "2026-09-02T12:00:00Z")
		second.Sequence = first.Sequence
		document.Observations = []observation{first, second}
		out := runReport(t, document, true)
		assertReportClass(t, out, "overlap", "missing")
		assertReportValue(t, out, "evidence_gaps", "1")
		assertRowsContain(t, out, "interval_comparisons", map[string]string{"left_observation_id": "unproven-a", "right_observation_id": "unproven-b", "status": "unknown"})

		unbounded := reportFailure("unbounded", "assignment-c", "", "")
		unbounded.Sequence = first.Sequence
		unbounded.IntervalReference = nativeReference("native:unbounded-interval")
		document.Observations = []observation{unbounded}
		out = runReport(t, document, true)
		assertRowsContain(t, out, "intervals", map[string]string{"observation_id": "unbounded", "status": "unbounded", "native": "native:unbounded-interval"})
		assertRowsContain(t, out, "evidence_gaps", map[string]string{"observation_id": "unbounded", "kind": "unbounded-interval"})
	})

	t.Run("order", func(t *testing.T) {
		requireFixtureCase(t)
		document := orderingFixture()
		out := runReport(t, document, true)
		reordered := document
		reordered.Observations = reverseObservations(document.Observations)
		reordered.Audits = reverseAudits(document.Audits)
		if repeated := runReport(t, reordered, true); repeated != out {
			t.Fatal("full report changed with import order")
		}
		observationRows := reportRows(t, out, "observations")
		if got := rowField(observationRows, "id"); !reflect.DeepEqual(got, []string{"sequence-first", "time-first", "id-a", "id-b"}) {
			t.Fatalf("observation order = %v", got)
		}
		auditRows := reportRows(t, out, "audits")
		if got := rowField(auditRows, "audit_id"); !reflect.DeepEqual(got, []string{"audit-a", "audit-z", "audit-other"}) {
			t.Fatalf("audit order = %v", got)
		}
	})
}

func TestRepairPilotIsolation(t *testing.T) {
	t.Run("read-only", func(t *testing.T) {
		requireFixtureCase(t)
		document := completeReportFixture()
		options := reportOptions(t, document)
		before, err := os.ReadFile(documentPath(options))
		if err != nil {
			t.Fatal(err)
		}
		if out, code := Command(options, []string{"report", "--full"}); code != 0 {
			t.Fatalf("report = output %q, exit %d", out, code)
		}
		after, err := os.ReadFile(documentPath(options))
		if err != nil {
			t.Fatal(err)
		}
		if string(after) != string(before) {
			t.Fatal("report changed the stored document")
		}
	})
}

func reportFixture() Document {
	cutoff := parseTime("2026-09-09T12:00:00Z")
	return Document{Version: 1, RepositoryKey: "bench-123", ActivatedAt: parseTime("2026-09-02T12:00:00Z"), CutoffAt: &cutoff, Observations: []observation{}, Audits: []audit{}}
}

func fixtureMissingClass(class string) Document {
	document := completeReportFixture()
	removeID := map[string]string{"productive": "repair-productive", "stalled": "repair-stalled", "unchanged": "rerun-unchanged", "overlap": "overlap-b"}[class]
	observations := []observation{}
	for _, item := range document.Observations {
		if item.ID != removeID {
			observations = append(observations, item)
		}
	}
	document.Observations = observations
	audits := []audit{}
	for _, item := range document.Audits {
		if item.ObservationID != removeID {
			audits = append(audits, item)
		}
	}
	document.Audits = audits
	return document
}

func completeReportFixture() Document {
	document := reportFixture()
	productive := reportRepair("repair-productive", "productive")
	stalled := reportRepair("repair-stalled", "stalled")
	unchanged := reportRerun("rerun-unchanged")
	overlapA := reportFailure("overlap-a", "assignment-a", "2026-09-02T09:00:00Z", "2026-09-02T11:00:00Z")
	overlapB := reportFailure("overlap-b", "assignment-b", "2026-09-02T10:00:00Z", "2026-09-02T12:00:00Z")
	overlapB.Sequence = overlapA.Sequence
	document.Observations = []observation{productive, stalled, unchanged, overlapA, overlapB}
	document.Audits = []audit{
		reportAudit("audit-productive", productive.ID, "productive"),
		reportAudit("audit-stalled", stalled.ID, "stalled"),
	}
	return document
}

func orderingFixture() Document {
	document := reportFixture()
	firstSequence := sequenceKey("source-a", "spec-a", "chunk-a")
	secondSequence := sequenceKey("source-b", "spec-a", "chunk-a")
	sequenceFirst := reportFailure("sequence-first", "assignment-a", "", "")
	sequenceFirst.Sequence = firstSequence
	sequenceFirst.ObservedAt = parseTime("2026-09-05T12:00:00Z")
	timeFirst := reportFailure("time-first", "assignment-a", "", "")
	timeFirst.Sequence = secondSequence
	timeFirst.ObservedAt = parseTime("2026-09-03T10:00:00Z")
	idB := reportFailure("id-b", "assignment-a", "", "")
	idB.Sequence = secondSequence
	idB.ObservedAt = parseTime("2026-09-03T11:00:00Z")
	idA := reportFailure("id-a", "assignment-a", "", "")
	idA.Sequence = secondSequence
	idA.ObservedAt = idB.ObservedAt
	document.Observations = []observation{idB, sequenceFirst, idA, timeFirst}
	document.Audits = []audit{
		reportAudit("audit-z", "id-a", "stalled"),
		reportAudit("audit-other", "id-b", "stalled"),
		reportAudit("audit-a", "id-a", "productive"),
	}
	return document
}

func reportRepair(id, label string) observation {
	return observation{
		ID: id, ObservedAt: parseTime("2026-09-03T12:00:00Z"), Sequence: sequenceKey("source-"+id, "spec-a", "chunk-a"),
		AssignmentID: "assignment-a", SessionID: "session-a", SourceRevision: "revision-b", Stage: "post-review", Kind: "repair",
		References:       []assessment.Reference{*nativeReference("native:observation-" + id)},
		Repair:           &repair{Hypothesis: "the target guard is missing", IntendedChange: "add the target guard", VerificationReference: *nativeReference("native:verification-" + id)},
		ProgressProposal: &progressProposal{Label: label, Reason: "the operator proposes this label"},
	}
}

func reportRerun(id string) observation {
	return observation{
		ID: id, ObservedAt: parseTime("2026-09-04T12:00:00Z"), Sequence: sequenceKey("source-"+id, "spec-a", "chunk-a"),
		AssignmentID: "assignment-a", SessionID: "session-a", SourceRevision: "revision-b", Stage: "post-review", Kind: "rerun",
		References: []assessment.Reference{*nativeReference("native:observation-" + id)},
		Rerun:      &rerun{VerificationReference: *nativeReference("native:rerun-" + id), UnchangedContent: true},
	}
}

func reportFailure(id, assignment, start, end string) observation {
	item := failureInput(id, sequenceKey("source-"+id, "spec-a", "chunk-a")).Observation
	item.AssignmentID = assignment
	if start != "" {
		item.StartedAt = timestamp(start)
		item.IntervalReference = nativeReference("native:interval-" + id)
	}
	if end != "" {
		item.EndedAt = timestamp(end)
		item.IntervalReference = nativeReference("native:interval-" + id)
	}
	return *item
}

func reportAudit(id, observationID, label string) audit {
	return *auditInput(id, observationID, label).Audit
}

func reportOptions(t *testing.T, document Document) Options {
	t.Helper()
	options := testOptions(t.TempDir())
	if err := os.MkdirAll(filepath.Dir(documentPath(options)), 0o700); err != nil {
		t.Fatal(err)
	}
	writeTestDocument(t, documentPath(options), document)
	options.Now = parseTime("2026-09-10T12:00:00Z")
	return options
}

func runReport(t *testing.T, document Document, full bool) string {
	t.Helper()
	return runReportAt(t, document, full, parseTime("2026-09-10T12:00:00Z"))
}

func runReportAt(t *testing.T, document Document, full bool, now time.Time) string {
	t.Helper()
	options := reportOptions(t, document)
	options.Now = now
	args := []string{"report"}
	if full {
		args = append(args, "--full")
	}
	out, code := Command(options, args)
	if code != 0 {
		t.Fatalf("report = output %q, exit %d", out, code)
	}
	return out
}

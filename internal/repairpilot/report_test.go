package repairpilot

import (
	"encoding/csv"
	"os"
	"path/filepath"
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
	})

	t.Run("productive", func(t *testing.T) {
		requireFixtureCase(t)
		document := reportFixture()
		document.Observations = []observation{reportRepair("repair-productive", "productive")}
		document.Audits = []audit{reportAudit("audit-productive", "repair-productive", "productive")}
		out := runReport(t, document, false)
		assertReportClass(t, out, "productive", "present")
	})

	t.Run("unchanged", func(t *testing.T) {
		requireFixtureCase(t)
		document := reportFixture()
		document.Observations = []observation{reportRerun("rerun-unchanged")}
		out := runReport(t, document, false)
		assertReportClass(t, out, "unchanged", "present")
	})

	t.Run("overlap", func(t *testing.T) {
		requireFixtureCase(t)
		document := reportFixture()
		first := reportFailure("overlap-a", "assignment-a", "2026-09-02T09:00:00Z", "2026-09-02T11:00:00Z")
		second := reportFailure("overlap-b", "assignment-b", "2026-09-02T10:00:00Z", "2026-09-02T12:00:00Z")
		document.Observations = []observation{first, second}
		out := runReport(t, document, false)
		assertReportClass(t, out, "overlap", "present")
	})

	t.Run("incomplete", func(t *testing.T) {
		requireFixtureCase(t)
		document := reportFixture()
		complete := reportFailure("complete", "assignment-a", "", "")
		complete.Endpoint = &endpoint{Kind: "reviewer-handoff", Reference: *nativeReference("native:endpoint")}
		incomplete := reportFailure("incomplete", "assignment-b", "", "")
		incomplete.Sequence = sequenceKey("source-b", "spec-a", "chunk-b")
		document.Observations = []observation{complete, incomplete}
		out := runReport(t, document, false)
		assertReportValue(t, out, "completed_sequences", "1")
		assertReportValue(t, out, "incomplete_sequences", "1")
	})

	t.Run("missing-class", func(t *testing.T) {
		requireFixtureCase(t)
		document := completeReportFixture()
		document.Observations = document.Observations[:len(document.Observations)-1]
		out := runReport(t, document, false)
		assertReportValue(t, out, "sample", "inconclusive")
		assertReportClass(t, out, "overlap", "missing")
	})

	t.Run("full", func(t *testing.T) {
		requireFixtureCase(t)
		document := completeReportFixture()
		unknown := reportFailure("unknown-evidence", "assignment-c", "2026-09-02T09:30:00Z", "")
		document.Observations = append(document.Observations, unknown)
		document.Observations[2].Endpoint = &endpoint{Kind: "reviewer-handoff", Reference: *nativeReference("native:full-endpoint")}
		out := runReport(t, document, true)
		for _, table := range []string{"sequences[", "observations[", "failures[", "observation_refs[", "repairs[", "reruns[", "endpoints[", "audits[", "audit_refs[", "intervals[", "interval_comparisons[", "evidence_gaps["} {
			if !strings.Contains(out, table) {
				t.Errorf("full report missing table %q:\n%s", table, out)
			}
		}
		for _, value := range []string{"the target guard is missing", "add the target guard", "reviewed native evidence", "native:verification-repair-productive", "unknown-evidence", "partial-interval", "reviewer-handoff"} {
			if !strings.Contains(out, value) {
				t.Errorf("full report missing retained value %q", value)
			}
		}
	})

	t.Run("default", func(t *testing.T) {
		requireFixtureCase(t)
		out := runReport(t, completeReportFixture(), false)
		if strings.Contains(out, "observations[") || !strings.Contains(out, "required_classes[") {
			t.Fatalf("default report has wrong detail shape:\n%s", out)
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
			document.Observations = []observation{
				reportFailure("edge-a", "assignment-a", pair[0], pair[1]),
				reportFailure("edge-b", "assignment-b", pair[2], pair[3]),
			}
			assertReportClass(t, runReport(t, document, false), "overlap", "missing")
		}
	})

	t.Run("interval-provenance", func(t *testing.T) {
		requireFixtureCase(t)
		document := reportFixture()
		first := reportFailure("unproven-a", "assignment-a", "2026-09-02T09:00:00Z", "2026-09-02T11:00:00Z")
		first.IntervalReference = nil
		second := reportFailure("unproven-b", "assignment-b", "2026-09-02T10:00:00Z", "2026-09-02T12:00:00Z")
		document.Observations = []observation{first, second}
		out := runReport(t, document, false)
		assertReportClass(t, out, "overlap", "missing")
		assertReportValue(t, out, "evidence_gaps", "1")
	})

	t.Run("order", func(t *testing.T) {
		requireFixtureCase(t)
		document := completeReportFixture()
		document.Observations[0], document.Observations[1] = document.Observations[1], document.Observations[0]
		document.Audits[0], document.Audits[1] = document.Audits[1], document.Audits[0]
		out := runReport(t, document, true)
		if repeated := runReport(t, document, true); repeated != out {
			t.Fatal("repeated full report changed with the same document and clock")
		}
		if strings.Index(out, "repair-productive") > strings.Index(out, "repair-stalled") {
			t.Fatalf("observations are not in sequence/time/id order:\n%s", out)
		}
		if strings.Index(out, "audit-productive") > strings.Index(out, "audit-stalled") {
			t.Fatalf("audits are not in observation/id order:\n%s", out)
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

func completeReportFixture() Document {
	document := reportFixture()
	productive := reportRepair("repair-productive", "productive")
	stalled := reportRepair("repair-stalled", "stalled")
	unchanged := reportRerun("rerun-unchanged")
	overlapA := reportFailure("overlap-a", "assignment-a", "2026-09-02T09:00:00Z", "2026-09-02T11:00:00Z")
	overlapB := reportFailure("overlap-b", "assignment-b", "2026-09-02T10:00:00Z", "2026-09-02T12:00:00Z")
	document.Observations = []observation{productive, stalled, unchanged, overlapA, overlapB}
	document.Audits = []audit{
		reportAudit("audit-productive", productive.ID, "productive"),
		reportAudit("audit-stalled", stalled.ID, "stalled"),
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

func assertReportClass(t *testing.T, out, class, status string) {
	t.Helper()
	if !strings.Contains(out, "  "+class+","+status+"\n") {
		t.Fatalf("class %s = want %s:\n%s", class, status, out)
	}
}

func assertReportValue(t *testing.T, out, column, value string) {
	t.Helper()
	lines := strings.Split(out, "\n")
	open := strings.Index(lines[0], "{")
	close := strings.LastIndex(lines[0], "}")
	if open < 0 || close < open || len(lines) < 2 {
		t.Fatalf("report summary is malformed:\n%s", out)
	}
	columns := strings.Split(lines[0][open+1:close], ",")
	values, err := csv.NewReader(strings.NewReader(strings.TrimSpace(lines[1]))).Read()
	if err != nil {
		t.Fatalf("parse report summary: %v\n%s", err, out)
	}
	for index, candidate := range columns {
		if candidate == column && index < len(values) && values[index] == value {
			return
		}
	}
	t.Fatalf("report %s = want %s:\n%s", column, value, out)
}

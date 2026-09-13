package repairpilot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/assessment"
)

func TestRepairPilotFailures(t *testing.T) {
	t.Run("repeated", func(t *testing.T) {
		h := newRecordHarness(t)
		first := failureInput("failure-first", sequenceKey("source-a", "spec-a", "chunk-a"))
		first.Observation.Failures = append(first.Observation.Failures,
			makeFailure("lint", "REQ-2", "second blocker", "inherited", false, "native:failure-2"))
		second := failureInput("failure-repeat", first.Observation.Sequence)
		second.Observation.Failures = append([]failure(nil), first.Observation.Failures...)
		want := append([]failure(nil), first.Observation.Failures...)
		h.accept(t, first)
		h.accept(t, second)
		document := h.document(t)
		if !reflect.DeepEqual(document.Observations[0].Failures, want) || !reflect.DeepEqual(document.Observations[1].Failures, want) {
			t.Fatalf("repeated blocker changed: %+v", document.Observations)
		}
		assertSummaryValue(t, h.report(t), "comparable_repeats", "2")
	})

	t.Run("partial", func(t *testing.T) {
		h := newRecordHarness(t)
		input := failureInput("partial", sequenceKey("source-a", "spec-a", "chunk-a"))
		input.Observation.FailureCompleteness = "first-only"
		h.accept(t, input)
		if got := h.document(t).Observations[0].FailureCompleteness; got != "first-only" {
			t.Fatalf("failure completeness = %q", got)
		}
	})

	t.Run("generic", func(t *testing.T) {
		h := newRecordHarness(t)
		first := failureInput("generic-first", sequenceKey("source-a", "spec-a", "chunk-a"))
		first.Observation.Failures[0].Identity = ""
		first.Observation.Failures[0].Diagnostic = "generic failure"
		second := failureInput("generic-second", first.Observation.Sequence)
		second.Observation.Failures[0].Identity = ""
		second.Observation.Failures[0].Diagnostic = "generic failure"
		h.accept(t, first)
		h.accept(t, second)
		failures := h.document(t).Observations
		if comparableFailure(failures[0].Failures[0], failures[1].Failures[0]) {
			t.Fatal("generic diagnostics became comparable without defect identity")
		}
		assertSummaryValue(t, h.report(t), "comparable_repeats", "0")
	})

	t.Run("rerun", func(t *testing.T) {
		h := newRecordHarness(t)
		key := sequenceKey("source-a", "spec-a", "chunk-a")
		h.accept(t, failureInput("rerun-first", key))
		input := rerunInput("rerun-second", key)
		h.accept(t, input)
		document := h.document(t)
		if repairCount(document) != 0 || document.Observations[1].Rerun == nil || !document.Observations[1].Rerun.UnchangedContent {
			t.Fatalf("rerun became a repair: %+v", document.Observations[1])
		}
	})

	t.Run("vocabulary", func(t *testing.T) {
		h := newRecordHarness(t)
		ownerships := []string{"diff-owned", "inherited", "spec-predicted", "unknown"}
		completeness := []string{"complete", "first-only", "unknown"}
		for i, ownership := range ownerships {
			input := failureInput("ownership-"+ownership, sequenceKey("source-"+ownership, "spec-a", "chunk-a"))
			input.Observation.Failures[0].Ownership = ownership
			input.Observation.FailureCompleteness = completeness[i%len(completeness)]
			h.accept(t, input)
		}
		document := h.document(t)
		for i, item := range document.Observations {
			if item.Failures[0].Ownership != ownerships[i] || item.FailureCompleteness != completeness[i%len(completeness)] {
				t.Fatalf("producer vocabulary changed: %+v", item)
			}
		}
	})
}

func TestRepairPilotEndpoints(t *testing.T) {
	t.Run("unverified", func(t *testing.T) {
		h := newRecordHarness(t)
		key := sequenceKey("source-a", "spec-a", "chunk-a")
		first := failureInput("closure-first", key)
		first.Observation.Failures = append(first.Observation.Failures, makeFailure("lint", "REQ-2", "second blocker", "inherited", true, "native:failure-2"))
		h.accept(t, first)
		before := h.bytes(t)
		closing := repairInput("closure-partial", key)
		closing.Observation.Endpoint = &endpoint{Kind: "verified-closure", Reference: *nativeReference("native:closure"), CoveredBlockerRefs: []assessment.Reference{*nativeReference("native:failure-1")}}
		h.refuse(t, closing)
		if got := h.bytes(t); !reflect.DeepEqual(got, before) {
			t.Fatal("partial closure changed prior evidence")
		}
	})

	t.Run("closed", func(t *testing.T) {
		h := newRecordHarness(t)
		key := sequenceKey("source-a", "spec-a", "chunk-a")
		first := failureInput("closed-first", key)
		first.Observation.Failures = append(first.Observation.Failures, makeFailure("lint", "REQ-2", "second blocker", "inherited", true, "native:failure-2"))
		h.accept(t, first)
		closing := repairInput("closed-endpoint", key)
		closing.Observation.Endpoint = &endpoint{Kind: "verified-closure", Reference: *nativeReference("native:closure"), CoveredBlockerRefs: []assessment.Reference{*nativeReference("native:failure-1"), *nativeReference("native:failure-2")}}
		h.accept(t, closing)
		document := h.document(t)
		if completedSequenceCount(document) != 1 || document.Observations[1].Endpoint.Kind != "verified-closure" {
			t.Fatalf("closure did not complete the sequence: %+v", document.Observations)
		}
	})

	t.Run("handoff", func(t *testing.T) {
		h := newRecordHarness(t)
		key := sequenceKey("source-a", "spec-a", "chunk-a")
		h.accept(t, failureInput("handoff-first", key))
		handoff := failureInput("handoff-endpoint", key)
		handoff.Observation.Endpoint = &endpoint{Kind: "reviewer-handoff", Reference: *nativeReference("native:review-handoff")}
		h.accept(t, handoff)
		document := h.document(t)
		if completedSequenceCount(document) != 1 || document.Observations[1].Endpoint.Kind != "reviewer-handoff" || len(document.Observations[1].Failures) == 0 {
			t.Fatalf("reviewer handoff erased unresolved blockers: %+v", document.Observations[1])
		}
	})
}

func TestRepairPilotCutoff(t *testing.T) {
	t.Run("count", func(t *testing.T) {
		h := newRecordHarness(t)
		completeSequences(t, h, 10)
		document := h.document(t)
		if completedSequenceCount(document) != 10 || document.CutoffAt == nil || !document.CutoffAt.Equal(h.options.Now) {
			t.Fatalf("count cutoff = completed %d, cutoff %v", completedSequenceCount(document), document.CutoffAt)
		}
		h.refuse(t, endpointInput("endpoint-11", sequenceKey("source-11", "spec-a", "chunk-11")))
	})

	t.Run("deadline", func(t *testing.T) {
		h := newRecordHarness(t)
		h.options.Now = parseTime("2026-09-15T12:00:00Z")
		before := h.bytes(t)
		h.refuse(t, failureInput("deadline", sequenceKey("source-a", "spec-a", "chunk-a")))
		if !reflect.DeepEqual(h.bytes(t), before) {
			t.Fatal("deadline-boundary import changed the document")
		}
	})

	t.Run("earlier", func(t *testing.T) {
		h := newRecordHarness(t)
		completeSequences(t, h, 9)
		h.options.Now = parseTime("2026-09-15T12:00:00Z")
		h.refuse(t, endpointInput("endpoint-10", sequenceKey("source-10", "spec-a", "chunk-10")))
		if got := completedSequenceCount(h.document(t)); got != 9 {
			t.Fatalf("deadline lost the earlier bound: completed %d", got)
		}
	})

	t.Run("late-import", func(t *testing.T) {
		h := newRecordHarness(t)
		completeSequences(t, h, 10)
		late := failureInput("late", sequenceKey("late-source", "spec-a", "late-chunk"))
		late.Observation.ObservedAt = parseTime("2026-09-02T12:00:00Z")
		h.refuse(t, late)
		if got := len(h.document(t).Observations); got != 10 {
			t.Fatalf("late import changed frozen sample: %d", got)
		}
	})

	t.Run("audit", func(t *testing.T) {
		h := newRecordHarness(t)
		completeSequences(t, h, 10)
		input := auditInput("post-cutoff-audit", "endpoint-1", "productive")
		h.accept(t, input)
		document := h.document(t)
		if len(document.Audits) != 1 || len(document.Observations) != 10 {
			t.Fatalf("post-cutoff audit changed observations: %+v", document)
		}
	})

	for _, boundary := range []struct {
		name  string
		stamp string
	}{
		{name: "before-activation", stamp: "2026-08-31T12:00:00Z"},
		{name: "future", stamp: "2026-09-11T12:00:00Z"},
	} {
		for _, field := range []struct {
			name   string
			assign func(*observation, *time.Time)
		}{
			{name: "observed-at", assign: func(value *observation, stamp *time.Time) { value.ObservedAt = *stamp }},
			{name: "started-at", assign: func(value *observation, stamp *time.Time) { value.StartedAt = stamp }},
			{name: "ended-at", assign: func(value *observation, stamp *time.Time) { value.EndedAt = stamp }},
		} {
			t.Run(boundary.name+"/"+field.name, func(t *testing.T) {
				h := newRecordHarness(t)
				input := failureInput(boundary.name+"-"+field.name, sequenceKey("source-a", "spec-a", "chunk-a"))
				stamp := parseTime(boundary.stamp)
				field.assign(input.Observation, &stamp)
				before := h.bytes(t)
				h.refuse(t, input)
				if !reflect.DeepEqual(h.bytes(t), before) {
					t.Fatal("out-of-window observation changed the document")
				}
			})
		}
	}
}

func TestRepairPilotAudit(t *testing.T) {
	t.Run("no-repair", func(t *testing.T) {
		h := newRecordHarness(t)
		h.accept(t, failureInput("audit-target", sequenceKey("source-a", "spec-a", "chunk-a")))
		input := auditInput("no-repair", "audit-target", "productive")
		input.Audit.Conclusion.RepairedTarget = ""
		out := h.accept(t, input)
		if got := effectiveProgressLabel(h.document(t), "audit-target"); got != "unknown" {
			t.Fatalf("unsupported progress label = %q", got)
		}
		assertSummaryValue(t, out, "progress_unknown", "1")
	})

	t.Run("proposal", func(t *testing.T) {
		h := newRecordHarness(t)
		input := failureInput("proposal-target", sequenceKey("source-a", "spec-a", "chunk-a"))
		input.Observation.ProgressProposal = &progressProposal{Label: "productive", Reason: "author believes the repair worked"}
		out := h.accept(t, input)
		document := h.document(t)
		if document.Observations[0].ProgressProposal == nil || effectiveProgressLabel(document, "proposal-target") != "unknown" {
			t.Fatalf("proposal became an accepted label: %+v", document)
		}
		assertSummaryValue(t, out, "progress_unknown", "1")
	})

	t.Run("supported", func(t *testing.T) {
		h := newRecordHarness(t)
		h.accept(t, failureInput("supported-target", sequenceKey("source-a", "spec-a", "chunk-a")))
		input := auditInput("supported-audit", "supported-target", "productive")
		out := h.accept(t, input)
		document := h.document(t)
		if effectiveProgressLabel(document, "supported-target") != "productive" || document.Audits[0].Conclusion.RepairedTarget != "REQ-1" || !assessment.ValidReference(document.Audits[0].Conclusion.VerificationReference) {
			t.Fatalf("supported audit lost proof: %+v", document.Audits[0])
		}
		assertSummaryValue(t, out, "progress_productive", "1")
	})

	t.Run("conflicting", func(t *testing.T) {
		invalid := newRecordHarness(t)
		invalid.accept(t, failureInput("same-label-target", sequenceKey("source-b", "spec-a", "chunk-b")))
		invalid.accept(t, auditInput("same-label-a", "same-label-target", "productive"))
		invalid.accept(t, auditInput("same-label-b", "same-label-target", "productive"))
		nonConflict := auditInput("same-label-resolution", "same-label-target", "productive")
		nonConflict.Audit.ResolvesAuditIDs = []string{"same-label-a", "same-label-b"}
		invalid.refuse(t, nonConflict)

		h := newRecordHarness(t)
		h.accept(t, failureInput("conflict-target", sequenceKey("source-a", "spec-a", "chunk-a")))
		h.accept(t, auditInput("productive-audit", "conflict-target", "productive"))
		h.accept(t, auditInput("stalled-audit", "conflict-target", "stalled"))
		if got := effectiveProgressLabel(h.document(t), "conflict-target"); got != "unknown" {
			t.Fatalf("conflicting label = %q, want unknown", got)
		}
		resolution := auditInput("resolution-audit", "conflict-target", "productive")
		resolution.Audit.ResolvesAuditIDs = []string{"productive-audit", "stalled-audit"}
		h.accept(t, resolution)
		document := h.document(t)
		if effectiveProgressLabel(document, "conflict-target") != "productive" || len(document.Audits) != 3 {
			t.Fatalf("resolution did not retain and resolve both audits: %+v", document.Audits)
		}
		later := auditInput("later-stalled-audit", "conflict-target", "stalled")
		out := h.accept(t, later)
		if got := effectiveProgressLabel(h.document(t), "conflict-target"); got != "unknown" {
			t.Fatalf("later contradiction label = %q, want unknown", got)
		}
		assertSummaryValue(t, out, "progress_unknown", "1")
	})
}

func TestRepairPilotGrammar(t *testing.T) {
	t.Run("usage", func(t *testing.T) {
		options := testOptions(t.TempDir())
		if err := os.MkdirAll(filepath.Dir(documentPath(options)), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink("missing-target", documentPath(options)); err != nil {
			t.Fatal(err)
		}
		out, code := Command(options, []string{"bogus"})
		if code != 2 || !strings.Contains(out, "unknown argument: bogus") {
			t.Fatalf("bogus operation = output %q, exit %d; want usage at exit 2", out, code)
		}
	})
	t.Run("hostile", func(t *testing.T) {
		requireFixtureCase(t)
		for _, fixture := range recordInputHostileFixtures(t) {
			t.Run(fixture.name, func(t *testing.T) {
				options := testOptions(t.TempDir())
				if out, code := Command(options, []string{"activate"}); code != 0 {
					t.Fatalf("activate = output %q, exit %d", out, code)
				}
				before, err := os.ReadFile(documentPath(options))
				if err != nil {
					t.Fatal(err)
				}
				inputPath := filepath.Join(t.TempDir(), "input.json")
				fixture.make(t, inputPath, options)
				out, code := Command(options, []string{"record", "--input", inputPath})
				if code != 1 || !strings.Contains(out, "refused") {
					t.Fatalf("hostile input %s = output %q, exit %d", fixture.name, out, code)
				}
				assertDocumentBytes(t, documentPath(options), before)
			})
		}
	})
	t.Run("path-shape", func(t *testing.T) {
		requireFixtureCase(t)
		options := testOptions(t.TempDir())
		if out, code := Command(options, []string{"activate"}); code != 0 {
			t.Fatalf("activate = output %q, exit %d", out, code)
		}
		options.Now = parseTime("2026-09-10T12:00:00Z")
		path := filepath.Join(t.TempDir(), "input [*].json")
		writeFixture(t, path, recordInputBytes(t, failureInput("path-shape", sequenceKey("source-a", "spec-a", "chunk-a"))))
		if out, code := Command(options, []string{"record", "--input", path}); code != 0 {
			t.Fatalf("record path shape = output %q, exit %d", out, code)
		}
	})
	t.Run("no-final-newline", func(t *testing.T) {
		options := testOptions(t.TempDir())
		if out, code := Command(options, []string{"activate"}); code != 0 {
			t.Fatalf("activate = output %q, exit %d", out, code)
		}
		options.Now = parseTime("2026-09-10T12:00:00Z")
		data, err := json.Marshal(failureInput("no-newline", sequenceKey("source-a", "spec-a", "chunk-a")))
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(t.TempDir(), "input.json")
		writeFixture(t, path, data)
		if out, code := Command(options, []string{"record", "--input", path}); code != 0 {
			t.Fatalf("record without final newline = output %q, exit %d", out, code)
		}
	})
}

func repairInput(id string, key sequence) recordInput {
	return recordInput{Version: 1, Observation: &observation{
		ID: id, ObservedAt: parseTime("2026-09-03T12:00:00Z"), Sequence: key,
		AssignmentID: "assignment-a", SessionID: "session-a", SourceRevision: "revision-b",
		Stage: "post-review", Kind: "repair", References: []assessment.Reference{*nativeReference("native:observation-" + id)},
		Repair: &repair{Hypothesis: "the guard is missing", IntendedChange: "add the guard", VerificationReference: *nativeReference("native:verification-" + id)},
	}}
}

func rerunInput(id string, key sequence) recordInput {
	return recordInput{Version: 1, Observation: &observation{
		ID: id, ObservedAt: parseTime("2026-09-03T12:00:00Z"), Sequence: key,
		AssignmentID: "assignment-a", SessionID: "session-a", SourceRevision: "revision-a",
		Stage: "pre-review", Kind: "rerun", References: []assessment.Reference{*nativeReference("native:observation-" + id)},
		Rerun: &rerun{VerificationReference: *nativeReference("native:verification-" + id), UnchangedContent: true},
	}}
}

func completeSequences(t *testing.T, h *recordHarness, count int) {
	t.Helper()
	for i := 1; i <= count; i++ {
		id := "endpoint-" + string(rune('0'+i))
		if i == 10 {
			id = "endpoint-10"
		}
		h.accept(t, endpointInput(id, sequenceKey("source-"+id, "spec-a", "chunk-"+id)))
	}
}

func auditInput(id, observationID, label string) recordInput {
	return recordInput{Version: 1, Audit: &audit{
		ID: id, ObservationID: observationID, Auditor: "reviewer-session",
		EvidenceReferences: []assessment.Reference{*nativeReference("native:evidence-" + id)}, Reason: "reviewed native evidence",
		Conclusion: &progressConclusion{Label: label, RepairedTarget: "REQ-1", VerificationReference: *nativeReference("native:verification-" + id)},
	}}
}

type shortWriteFile struct {
	*os.File
}

func (file shortWriteFile) Write(data []byte) (int, error) {
	return file.File.Write(data[:len(data)/2])
}

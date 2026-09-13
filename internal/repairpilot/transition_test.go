package repairpilot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
)

func TestRepairPilotFailures(t *testing.T) {
	t.Run("repeated", func(t *testing.T) {
		h := newRecordHarness(t)
		first := failureInput("failure-first", sequenceKey("source-a", "spec-a", "chunk-a"))
		second := failureInput("failure-repeat", first.Observation.Sequence)
		h.accept(t, first)
		h.accept(t, second)
		document := h.document(t)
		if !reflect.DeepEqual(document.Observations[0].Failures[0], document.Observations[1].Failures[0]) {
			t.Fatalf("repeated blocker changed: %+v", document.Observations)
		}
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
		closing.Observation.Endpoint = &endpoint{Kind: "verified-closure", Reference: "native:closure", CoveredBlockerRefs: []string{"native:failure-1"}}
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
		closing.Observation.Endpoint = &endpoint{Kind: "verified-closure", Reference: "native:closure", CoveredBlockerRefs: []string{"native:failure-1", "native:failure-2"}}
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
		handoff.Observation.Endpoint = &endpoint{Kind: "reviewer-handoff", Reference: "native:review-handoff"}
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

	t.Run("before-activation", func(t *testing.T) {
		h := newRecordHarness(t)
		input := failureInput("before-activation", sequenceKey("source-a", "spec-a", "chunk-a"))
		input.Observation.ObservedAt = parseTime("2026-08-31T12:00:00Z")
		before := h.bytes(t)
		h.refuse(t, input)
		if !reflect.DeepEqual(h.bytes(t), before) {
			t.Fatal("pre-activation observation changed the document")
		}
	})

	t.Run("future", func(t *testing.T) {
		h := newRecordHarness(t)
		input := failureInput("future", sequenceKey("source-a", "spec-a", "chunk-a"))
		input.Observation.ObservedAt = parseTime("2026-09-11T12:00:00Z")
		before := h.bytes(t)
		h.refuse(t, input)
		if !reflect.DeepEqual(h.bytes(t), before) {
			t.Fatal("future observation changed the document")
		}
	})
}

func TestRepairPilotAudit(t *testing.T) {
	t.Run("no-repair", func(t *testing.T) {
		h := newRecordHarness(t)
		h.accept(t, failureInput("audit-target", sequenceKey("source-a", "spec-a", "chunk-a")))
		input := auditInput("no-repair", "audit-target", "productive")
		input.Audit.Conclusion.RepairedTarget = ""
		h.accept(t, input)
		if got := effectiveProgressLabel(h.document(t), "audit-target"); got != "unknown" {
			t.Fatalf("unsupported progress label = %q", got)
		}
	})

	t.Run("proposal", func(t *testing.T) {
		h := newRecordHarness(t)
		input := failureInput("proposal-target", sequenceKey("source-a", "spec-a", "chunk-a"))
		input.Observation.ProgressProposal = &progressProposal{Label: "productive", Reason: "author believes the repair worked"}
		h.accept(t, input)
		document := h.document(t)
		if document.Observations[0].ProgressProposal == nil || effectiveProgressLabel(document, "proposal-target") != "unknown" {
			t.Fatalf("proposal became an accepted label: %+v", document)
		}
	})

	t.Run("supported", func(t *testing.T) {
		h := newRecordHarness(t)
		h.accept(t, failureInput("supported-target", sequenceKey("source-a", "spec-a", "chunk-a")))
		input := auditInput("supported-audit", "supported-target", "productive")
		h.accept(t, input)
		document := h.document(t)
		if effectiveProgressLabel(document, "supported-target") != "productive" || document.Audits[0].Conclusion.RepairedTarget != "REQ-1" || document.Audits[0].Conclusion.VerificationReference == "" {
			t.Fatalf("supported audit lost proof: %+v", document.Audits[0])
		}
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
		fixtures := []storedHostileFixture{
			{name: "empty", make: func(t *testing.T, path string, _ Options) { writeFixture(t, path, nil) }},
			{name: "malformed", make: func(t *testing.T, path string, _ Options) { writeFixture(t, path, []byte("{\n")) }},
			{name: "duplicate-key", make: func(t *testing.T, path string, _ Options) {
				writeFixture(t, path, []byte("{\"version\":1,\"version\":1}\n"))
			}},
			{name: "unknown-field", make: func(t *testing.T, path string, _ Options) {
				writeFixture(t, path, []byte("{\"version\":1,\"extra\":true}\n"))
			}},
			{name: "unsupported-version", make: func(t *testing.T, path string, _ Options) { writeFixture(t, path, []byte("{\"version\":2}\n")) }},
			{name: "both-payloads", make: func(t *testing.T, path string, _ Options) {
				writeFixture(t, path, []byte("{\"version\":1,\"observation\":{},\"audit\":{}}\n"))
			}},
			{name: "oversized", make: func(t *testing.T, path string, _ Options) {
				writeFixture(t, path, []byte(strings.Repeat("x", int(bounds.ControlRecordLimit)+1)))
			}},
			{name: "directory", make: func(t *testing.T, path string, _ Options) {
				if err := os.Mkdir(path, 0o700); err != nil {
					t.Fatal(err)
				}
			}},
			{name: "symlink", make: func(t *testing.T, path string, _ Options) {
				if err := os.Symlink("missing", path); err != nil {
					t.Fatal(err)
				}
			}},
		}
		fixtures = append(fixtures, platformStoredHostileFixtures()...)
		for _, fixture := range fixtures {
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
		Stage: "post-review", Kind: "repair", References: []string{"native:observation-" + id},
		Repair: &repair{Hypothesis: "the guard is missing", IntendedChange: "add the guard", VerificationReference: "native:verification-" + id},
	}}
}

func rerunInput(id string, key sequence) recordInput {
	return recordInput{Version: 1, Observation: &observation{
		ID: id, ObservedAt: parseTime("2026-09-03T12:00:00Z"), Sequence: key,
		AssignmentID: "assignment-a", SessionID: "session-a", SourceRevision: "revision-a",
		Stage: "pre-review", Kind: "rerun", References: []string{"native:observation-" + id},
		Rerun: &rerun{VerificationReference: "native:verification-" + id, UnchangedContent: true},
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
		EvidenceReferences: []string{"native:evidence-" + id}, Reason: "reviewed native evidence",
		Conclusion: &progressConclusion{Label: label, RepairedTarget: "REQ-1", VerificationReference: "native:verification-" + id},
	}}
}

type shortWriteFile struct {
	*os.File
}

type storedHostileFixture struct {
	name string
	make func(*testing.T, string, Options)
}

func (file shortWriteFile) Write(data []byte) (int, error) {
	return file.File.Write(data[:len(data)/2])
}

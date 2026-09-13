package repairpilot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestRepairPilotIdentity(t *testing.T) {
	t.Run("session", func(t *testing.T) {
		h := newRecordHarness(t)
		first := failureInput("session-first", sequenceKey("source-a", "spec-a", "chunk-a"))
		second := failureInput("session-second", first.Observation.Sequence)
		second.Observation.SessionID = "session-b"
		h.accept(t, first)
		h.accept(t, second)
		if got := sequenceCount(h.document(t)); got != 1 {
			t.Fatalf("sequence count = %d, want 1", got)
		}
	})

	t.Run("assignment", func(t *testing.T) {
		h := newRecordHarness(t)
		first := failureInput("assignment-first", sequenceKey("source-a", "spec-a", "chunk-a"))
		first.Observation.AssignmentID = "assignment-a"
		first.Observation.StartedAt = timestamp("2026-09-02T09:00:00Z")
		first.Observation.EndedAt = timestamp("2026-09-02T11:00:00Z")
		first.Observation.IntervalReference = "native:interval-a"
		second := failureInput("assignment-second", first.Observation.Sequence)
		second.Observation.AssignmentID = "assignment-b"
		second.Observation.StartedAt = timestamp("2026-09-02T10:00:00Z")
		second.Observation.EndedAt = timestamp("2026-09-02T12:00:00Z")
		second.Observation.IntervalReference = "native:interval-b"
		h.accept(t, first)
		h.accept(t, second)
		document := h.document(t)
		if document.Observations[0].AssignmentID != "assignment-a" || document.Observations[1].AssignmentID != "assignment-b" {
			t.Fatalf("assignment identities changed: %+v", document.Observations)
		}
	})

	t.Run("findings", func(t *testing.T) {
		h := newRecordHarness(t)
		input := failureInput("two-findings", sequenceKey("source-a", "spec-a", "chunk-a"))
		input.Observation.Failures = append(input.Observation.Failures, makeFailure("lint", "REQ-2", "second blocker", "inherited", true, "native:failure-2"))
		h.accept(t, input)
		document := h.document(t)
		if len(document.Observations) != 1 || len(document.Observations[0].Failures) != 2 || sequenceCount(document) != 1 {
			t.Fatalf("findings inflated the sequence: %+v", document.Observations)
		}
	})

	t.Run("sources", func(t *testing.T) {
		h := newRecordHarness(t)
		h.accept(t, failureInput("source-a", sequenceKey("source-a", "spec-a", "chunk-a")))
		h.accept(t, failureInput("source-b", sequenceKey("source-b", "spec-a", "chunk-a")))
		if got := sequenceCount(h.document(t)); got != 2 {
			t.Fatalf("sequence count = %d, want 2", got)
		}
	})

	t.Run("stages", func(t *testing.T) {
		h := newRecordHarness(t)
		first := failureInput("pre-review", sequenceKey("source-a", "spec-a", "chunk-a"))
		second := failureInput("post-review", first.Observation.Sequence)
		second.Observation.Stage = "post-review"
		h.accept(t, first)
		h.accept(t, second)
		document := h.document(t)
		if sequenceCount(document) != 1 || document.Observations[0].Stage != "pre-review" || document.Observations[1].Stage != "post-review" {
			t.Fatalf("review stages changed sequence identity: %+v", document.Observations)
		}
	})

	t.Run("missing", func(t *testing.T) {
		for _, mutate := range []func(*observation){
			func(value *observation) { value.AssignmentID = "" },
			func(value *observation) { value.SourceRevision = "" },
			func(value *observation) { value.Sequence.Source = "" },
		} {
			h := newRecordHarness(t)
			input := failureInput("missing", sequenceKey("source-a", "spec-a", "chunk-a"))
			mutate(input.Observation)
			h.refuse(t, input)
			if len(h.document(t).Observations) != 0 {
				t.Fatal("missing attribution changed the document")
			}
		}
	})

	t.Run("retry", func(t *testing.T) {
		h := newRecordHarness(t)
		input := failureInput("retry", sequenceKey("source-a", "spec-a", "chunk-a"))
		h.accept(t, input)
		h.accept(t, input)
		if got := len(h.document(t).Observations); got != 1 {
			t.Fatalf("observation count = %d, want 1", got)
		}
	})

	t.Run("conflict", func(t *testing.T) {
		h := newRecordHarness(t)
		input := failureInput("conflict", sequenceKey("source-a", "spec-a", "chunk-a"))
		h.accept(t, input)
		before := h.bytes(t)
		input.Observation.Failures[0].Diagnostic = "changed"
		h.refuse(t, input)
		if got := h.bytes(t); !reflect.DeepEqual(got, before) {
			t.Fatal("conflicting import changed prior evidence")
		}
	})

	t.Run("nonblocking-start", func(t *testing.T) {
		h := newRecordHarness(t)
		input := failureInput("nonblocking", sequenceKey("source-a", "spec-a", "chunk-a"))
		input.Observation.Failures[0].Blocking = false
		h.refuse(t, input)
		if sequenceCount(h.document(t)) != 0 {
			t.Fatal("nonblocking observation started a sequence")
		}
	})

	t.Run("first-blocker", func(t *testing.T) {
		invalid := newRecordHarness(t)
		attempt := repairInput("invented-repair", sequenceKey("source-a", "spec-a", "chunk-a"))
		attempt.Observation.Failures = []failure{makeFailure("unit", "REQ-1", "first blocker", "diff-owned", true, "native:failure-1")}
		attempt.Observation.FailureCompleteness = "complete"
		invalid.refuse(t, attempt)
		if len(invalid.document(t).Observations) != 0 {
			t.Fatal("repair-shaped first blocker started a sequence")
		}
		h := newRecordHarness(t)
		input := failureInput("first-blocker", sequenceKey("source-a", "spec-a", "chunk-a"))
		h.accept(t, input)
		document := h.document(t)
		if document.Observations[0].Kind != "failure" || repairCount(document) != 0 {
			t.Fatalf("first blocker invented a repair: %+v", document.Observations[0])
		}
	})

	t.Run("specs", func(t *testing.T) {
		h := newRecordHarness(t)
		h.accept(t, failureInput("spec-a", sequenceKey("source-a", "spec-a", "chunk-a")))
		h.accept(t, failureInput("spec-b", sequenceKey("source-a", "spec-b", "chunk-a")))
		if got := sequenceCount(h.document(t)); got != 2 {
			t.Fatalf("sequence count = %d, want 2", got)
		}
	})

	t.Run("interval", func(t *testing.T) {
		h := newRecordHarness(t)
		input := failureInput("interval", sequenceKey("source-a", "spec-a", "chunk-a"))
		input.Observation.StartedAt = timestamp("2026-09-02T09:00:00Z")
		input.Observation.EndedAt = timestamp("2026-09-02T11:00:00Z")
		input.Observation.IntervalReference = "native:interval"
		h.accept(t, input)
		stored := h.document(t).Observations[0]
		if !stored.StartedAt.Equal(*input.Observation.StartedAt) || !stored.EndedAt.Equal(*input.Observation.EndedAt) || stored.IntervalReference != "native:interval" {
			t.Fatalf("interval changed: %+v", stored)
		}
	})

	t.Run("contributors", func(t *testing.T) {
		h := newRecordHarness(t)
		first := failureInput("contributor-a", sequenceKey("source-a", "spec-a", "chunk-a"))
		second := failureInput("contributor-b", first.Observation.Sequence)
		second.Observation.AssignmentID = "assignment-b"
		second.Observation.SessionID = "session-b"
		h.accept(t, first)
		h.accept(t, second)
		if got := sequenceCount(h.document(t)); got != 1 {
			t.Fatalf("sequence count = %d, want 1", got)
		}
	})

	t.Run("reversed-interval", func(t *testing.T) {
		h := newRecordHarness(t)
		input := failureInput("reversed", sequenceKey("source-a", "spec-a", "chunk-a"))
		input.Observation.StartedAt = timestamp("2026-09-02T11:00:00Z")
		input.Observation.EndedAt = timestamp("2026-09-02T09:00:00Z")
		input.Observation.IntervalReference = "native:interval"
		h.refuse(t, input)
		if len(h.document(t).Observations) != 0 {
			t.Fatal("reversed interval changed the document")
		}
	})
}

type recordHarness struct {
	t       *testing.T
	options Options
}

func newRecordHarness(t *testing.T) *recordHarness {
	t.Helper()
	requireFixtureCase(t)
	options := testOptions(t.TempDir())
	if out, code := Command(options, []string{"activate"}); code != 0 {
		t.Fatalf("activate = output %q, exit %d", out, code)
	}
	options.Now = time.Date(2026, 9, 10, 12, 0, 0, 0, time.UTC)
	return &recordHarness{t: t, options: options}
}

func (h *recordHarness) accept(t *testing.T, input recordInput) {
	t.Helper()
	out, code := h.run(t, input)
	if code != 0 || (!strings.Contains(out, "active") && !strings.Contains(out, "stopped")) {
		t.Fatalf("record = output %q, exit %d", out, code)
	}
}

func (h *recordHarness) refuse(t *testing.T, input recordInput) {
	t.Helper()
	if out, code := h.run(t, input); code != 1 || !strings.Contains(out, "refused") {
		t.Fatalf("record = output %q, exit %d; want refusal", out, code)
	}
}

func (h *recordHarness) run(t *testing.T, input recordInput) (string, int) {
	t.Helper()
	data, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "input.json")
	if err := os.WriteFile(path, append(data, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
	return Command(h.options, []string{"record", "--input", path})
}

func (h *recordHarness) document(t *testing.T) Document {
	t.Helper()
	document, state, err := load(h.options)
	if err != nil || state == "absent" {
		t.Fatalf("load = state %s, err %v", state, err)
	}
	return document
}

func (h *recordHarness) bytes(t *testing.T) []byte {
	t.Helper()
	data, err := os.ReadFile(documentPath(h.options))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func failureInput(id string, sequence sequence) recordInput {
	return recordInput{Version: 1, Observation: &observation{
		ID: id, ObservedAt: parseTime("2026-09-02T12:00:00Z"), Sequence: sequence,
		AssignmentID: "assignment-a", SessionID: "session-a", SourceRevision: "revision-a",
		Stage: "pre-review", Kind: "failure", FailureCompleteness: "complete",
		Failures:   []failure{makeFailure("unit", "REQ-1", "first blocker", "diff-owned", true, "native:failure-1")},
		References: []string{"native:observation-" + id},
	}}
}

func endpointInput(id string, key sequence) recordInput {
	input := failureInput(id, key)
	input.Observation.Endpoint = &endpoint{Kind: "reviewer-handoff", Reference: "native:handoff-" + id}
	return input
}

func makeFailure(check, identity, diagnostic, ownership string, blocking bool, reference string) failure {
	return failure{Check: check, Identity: identity, Diagnostic: diagnostic, Ownership: ownership, Blocking: blocking, Reference: reference}
}

func sequenceKey(source, spec, chunk string) sequence {
	return sequence{Source: source, Spec: spec, Chunk: chunk}
}

func timestamp(value string) *time.Time {
	parsed := parseTime(value)
	return &parsed
}

func parseTime(value string) time.Time {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		panic(err)
	}
	return parsed
}

func sequenceCount(document Document) int {
	values := map[sequence]bool{}
	for _, item := range document.Observations {
		values[item.Sequence] = true
	}
	return len(values)
}

func repairCount(document Document) int {
	count := 0
	for _, item := range document.Observations {
		if item.Kind == "repair" {
			count++
		}
	}
	return count
}

func requireFixtureCase(t *testing.T) {
	t.Helper()
	group, name, ok := strings.Cut(t.Name(), "/")
	if !ok {
		return
	}
	group = strings.TrimPrefix(group, "TestRepairPilot")
	group = strings.ToLower(group)
	data, err := os.ReadFile(filepath.Join("testdata", "cases.json"))
	if err != nil {
		t.Fatal(err)
	}
	var inventory struct {
		Version int                 `json:"version"`
		Groups  map[string][]string `json:"groups"`
	}
	if err := json.Unmarshal(data, &inventory); err != nil || inventory.Version != 1 {
		t.Fatalf("fixture inventory = version %d, err %v", inventory.Version, err)
	}
	for _, candidate := range inventory.Groups[group] {
		if candidate == name {
			return
		}
	}
	t.Fatalf("fixture case %s/%s is not inventoried", group, name)
}

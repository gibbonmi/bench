package assessment

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/census"
)

// orchestrated is a run whose attempts include the orchestrator's own work.
func orchestrated(root string) Run {
	r := fixtureRun(root)
	r.Attempts[0].State = "succeeded"
	r.Attempts = append(r.Attempts, Attempt{
		AttemptID: "attempt-orchestrator", ChunkID: "1", Role: "orchestration",
		SessionID: "session-orchestrator", Model: "synthetic", Effort: "medium",
		State: "succeeded",
	})
	return r
}

// writeCensus seeds one assignment's census file with three raw calls.
func writeCensus(t *testing.T, s Store, assignment string) {
	t.Helper()
	dir := census.Dir(s.Home, s.Root)
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	body := "2026-01-01T00:00:00Z\tgit\n2026-01-01T00:00:01Z\tgo\n2026-01-01T00:00:02Z\trm\n"
	if err := os.WriteFile(filepath.Join(dir, assignment), []byte(body), 0600); err != nil {
		t.Fatal(err)
	}
}

// pick maps one native selector onto a named attempt and role.
func pick(id, attempt, role string) map[string]any {
	return map[string]any{"id": id, "attempt_id": attempt, "chunk_id": "1", "role": role}
}

// runInput writes a record input built from the given run with overlaid keys.
func runInput(t *testing.T, r Run, extra map[string]any) string {
	t.Helper()
	data, _ := json.Marshal(r)
	var object map[string]any
	json.Unmarshal(data, &object)
	for k, v := range extra {
		object[k] = v
	}
	data, _ = json.Marshal(object)
	path := filepath.Join(t.TempDir(), "input.json")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

const (
	assignmentA = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	assignmentB = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
)

// DI14: one import collects explicit evidence from several assignments.
func TestAssessmentAssignmentBatches(t *testing.T) {
	s := Store{Home: t.TempDir(), Root: t.TempDir()}
	writeCensus(t, s, assignmentA)
	writeCensus(t, s, assignmentB)
	input := runInput(t, orchestrated(s.Root), map[string]any{
		"bench_input_batches": []any{
			map[string]any{"assignment_id": assignmentA, "census_event_ids": []any{pick(assignmentA+":1", "attempt-1", "implementation")}},
			map[string]any{"assignment_id": assignmentB, "census_event_ids": []any{pick(assignmentB+":2", "attempt-orchestrator", "orchestration")}},
		},
	})
	if out, code := Command(s, []string{"record", "--input", input}); code != 0 {
		t.Fatalf("two assignments refused in one import: %s", out)
	}
	r, err := s.Read("run-1")
	if err != nil {
		t.Fatal(err)
	}
	// Each batch must reach its own attempt. A singular-only collector would
	// carry at most one of these.
	for _, want := range []struct{ attempt, key string }{
		{"attempt-1", "raw_commands/" + assignmentA + ":1"},
		{"attempt-orchestrator", "raw_commands/" + assignmentB + ":2"},
	} {
		found := false
		for _, a := range r.Attempts {
			if a.AttemptID == want.attempt {
				_, found = a.Measures[want.key]
			}
		}
		if !found {
			t.Fatalf("attempt %s lost its batch measure %s", want.attempt, want.key)
		}
	}
}

// DI15: conflicting ownership across batches refuses, and the refusal leaves
// the stored record untouched.
func TestAssessmentBatchConflicts(t *testing.T) {
	cases := []struct {
		name, want string
		batches    []any
	}{
		{"both input forms", "duplicate bench input forms", nil},
		{"foreign assignment", "foreign census assignment", []any{
			map[string]any{"assignment_id": assignmentA, "census_event_ids": []any{pick(assignmentB+":1", "attempt-1", "implementation")}},
		}},
		{"repeated assignment batch", "duplicate assignment batch", []any{
			map[string]any{"assignment_id": assignmentA, "census_event_ids": []any{pick(assignmentA+":1", "attempt-1", "implementation")}},
			map[string]any{"assignment_id": assignmentA, "census_event_ids": []any{pick(assignmentA+":2", "attempt-1", "implementation")}},
		}},
		{"unknown attempt", "unknown mapped attempt", []any{
			map[string]any{"assignment_id": assignmentA, "census_event_ids": []any{pick(assignmentA+":1", "attempt-absent", "implementation")}},
		}},
		{"invalid assignment", "invalid expected assignment", []any{
			map[string]any{"assignment_id": "not-an-assignment", "census_event_ids": []any{pick("not-an-assignment:1", "attempt-1", "implementation")}},
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := Store{Home: t.TempDir(), Root: t.TempDir()}
			writeCensus(t, s, assignmentA)
			writeCensus(t, s, assignmentB)
			extra := map[string]any{"bench_input_batches": tc.batches}
			if tc.batches == nil {
				extra["bench_inputs"] = map[string]any{"assignment_id": assignmentA, "census_event_ids": []any{pick(assignmentA+":1", "attempt-1", "implementation")}}
				extra["bench_input_batches"] = []any{map[string]any{"assignment_id": assignmentB, "census_event_ids": []any{pick(assignmentB+":1", "attempt-1", "implementation")}}}
			}
			out, code := Command(s, []string{"record", "--input", runInput(t, orchestrated(s.Root), extra)})
			if code == 0 || !strings.Contains(out, tc.want) {
				t.Fatalf("conflicting import accepted or lost its reason (%d): %s", code, out)
			}
			if _, err := s.Read("run-1"); err == nil {
				t.Fatal("a refused import changed the stored record")
			}
		})
	}
}

// An explicit but empty batch list supplies nothing. It is not a second input
// form, so it neither refuses beside the singular form nor collects anything.
func TestAssessmentEmptyBatches(t *testing.T) {
	t.Run("beside the singular form", func(t *testing.T) {
		s := Store{Home: t.TempDir(), Root: t.TempDir()}
		writeCensus(t, s, assignmentA)
		input := runInput(t, orchestrated(s.Root), map[string]any{
			"bench_inputs":        map[string]any{"assignment_id": assignmentA, "census_event_ids": []any{pick(assignmentA+":1", "attempt-1", "implementation")}},
			"bench_input_batches": []any{},
		})
		if out, code := Command(s, []string{"record", "--input", input}); code != 0 {
			t.Fatalf("an empty batch list refused beside the singular form: %s", out)
		}
		r, err := s.Read("run-1")
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := r.Attempts[0].Measures["raw_commands/"+assignmentA+":1"]; !ok {
			t.Fatal("the singular form stopped collecting beside an empty batch list")
		}
	})
	t.Run("alone", func(t *testing.T) {
		s := Store{Home: t.TempDir(), Root: t.TempDir()}
		input := runInput(t, orchestrated(s.Root), map[string]any{"bench_input_batches": []any{}})
		if out, code := Command(s, []string{"record", "--input", input}); code != 0 {
			t.Fatalf("an empty batch list alone refused: %s", out)
		}
		r, err := s.Read("run-1")
		if err != nil {
			t.Fatal(err)
		}
		if len(r.Attempts[0].Measures) != 0 {
			t.Fatalf("an empty batch list collected something: %+v", r.Attempts[0].Measures)
		}
	})
}

// One native trace selected by two batches under two different mappings is
// ambiguous however it arrived. The selection ledger spans the whole import, so
// the second batch refuses rather than quietly claiming the event twice.
func TestAssessmentCrossBatchMapping(t *testing.T) {
	s := Store{Home: t.TempDir(), Root: t.TempDir()}
	nativeSpan(t, s, assignmentA, true)
	input := runInput(t, orchestrated(s.Root), map[string]any{
		"bench_input_batches": []any{
			map[string]any{"assignment_id": assignmentA, "trace_ids": []any{pick("trace-1", "attempt-1", "implementation")}},
			map[string]any{"assignment_id": assignmentB, "trace_ids": []any{pick("trace-1", "attempt-orchestrator", "orchestration")}},
		},
	})
	out, code := Command(s, []string{"record", "--input", input})
	if code == 0 || !strings.Contains(out, "ambiguous selection mapping") {
		t.Fatalf("one trace reached two attempts through two batches (%d): %s", code, out)
	}
	if _, err := s.Read("run-1"); err == nil {
		t.Fatal("a refused cross-batch import changed the stored record")
	}
}

// A selector repeated with one mapping contributes once rather than twice.
func TestAssessmentRepeatedSelectors(t *testing.T) {
	s := Store{Home: t.TempDir(), Root: t.TempDir()}
	writeCensus(t, s, assignmentA)
	input := runInput(t, orchestrated(s.Root), map[string]any{
		"bench_input_batches": []any{
			map[string]any{"assignment_id": assignmentA, "census_event_ids": []any{
				pick(assignmentA+":1", "attempt-1", "implementation"),
				pick(assignmentA+":1", "attempt-1", "implementation"),
			}},
		},
	})
	if out, code := Command(s, []string{"record", "--input", input}); code != 0 {
		t.Fatalf("an identical repeated selector refused: %s", out)
	}
	r, err := s.Read("run-1")
	if err != nil {
		t.Fatal(err)
	}
	value := r.Attempts[0].Measures["raw_commands/"+assignmentA+":1"].Value
	if value == nil || *value != 1 {
		t.Fatalf("a repeated native event contributed twice: %+v", r.Attempts[0].Measures)
	}
}

// DI19: a singular import keeps the outcome it had before batches existed.
// The two forms carry one selector, so they must agree on every collected
// measure and on every other attempt field.
func TestAssessmentSingularParity(t *testing.T) {
	collected := func(t *testing.T, form string) ([]string, string) {
		t.Helper()
		s := Store{Home: t.TempDir(), Root: t.TempDir()}
		writeCensus(t, s, assignmentA)
		body := map[string]any{"assignment_id": assignmentA, "census_event_ids": []any{pick(assignmentA+":1", "attempt-1", "implementation")}}
		extra := map[string]any{form: body}
		if form == "bench_input_batches" {
			extra[form] = []any{body}
		}
		if out, code := Command(s, []string{"record", "--input", runInput(t, orchestrated(s.Root), extra)}); code != 0 {
			t.Fatalf("%s refused: %s", form, out)
		}
		r, err := s.Read("run-1")
		if err != nil {
			t.Fatal(err)
		}
		// Measure and evidence references embed this store's own paths, which
		// differ between the two runs. Every other field does not, so the
		// comparison keeps the measured values and drops only the references.
		keys := []string{}
		for i := range r.Attempts {
			for key, m := range r.Attempts[i].Measures {
				value := "unknown"
				if m.Value != nil {
					value = fmt.Sprint(*m.Value)
				}
				keys = append(keys, r.Attempts[i].AttemptID+" "+key+"="+value)
				m.Reference = Reference{}
				r.Attempts[i].Measures[key] = m
			}
			r.Attempts[i].Evidence = nil
		}
		sort.Strings(keys)
		rest, _ := json.Marshal(r.Attempts)
		return keys, string(rest)
	}
	oldKeys, oldRest := collected(t, "bench_inputs")
	newKeys, newRest := collected(t, "bench_input_batches")
	if len(oldKeys) == 0 {
		t.Fatal("the singular form collected nothing, so parity proves nothing")
	}
	if !reflect.DeepEqual(oldKeys, newKeys) {
		t.Fatalf("the two forms collected different measures:\n%v\n%v", oldKeys, newKeys)
	}
	if oldRest != newRest {
		t.Fatalf("the two forms disagree:\n%s\n%s", oldRest, newRest)
	}
}

package assessment

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

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

// DI13: orchestration is its own role in show and in comparison.
func TestAssessmentOrchestration(t *testing.T) {
	if !contains(Roles(), "orchestration") {
		t.Fatal("the role vocabulary omits orchestration")
	}
	s := Store{Home: t.TempDir(), Root: t.TempDir()}
	if err := s.Record(orchestrated(s.Root)); err != nil {
		t.Fatal(err)
	}
	out, code := Command(s, []string{"show", "run-1"})
	if code != 0 || !strings.Contains(out, "orchestration") {
		t.Fatalf("show omits the orchestrator (%d): %s", code, out)
	}
	if !strings.Contains(out, "attempt-orchestrator") {
		t.Fatalf("show omits the orchestrator's attempt: %s", out)
	}
}

// contains is the local membership helper for the role vocabulary.
func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

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

// DI16: estimates, actual charges, currencies, measured zero, and unknowns
// stay distinct in one summary.
func TestAssessmentDelegatedCosts(t *testing.T) {
	r := orchestrated("")
	r.RepoKey = fixtureRun("").RepoKey
	r.Attempts[0].Cost.Actual = []Charge{
		{Kind: "invoice", Amount: ptr(0.0), Currency: "USD", Reference: Reference{"billing", "fixture:zero"}},
	}
	r.Attempts[1].Cost.Actual = []Charge{
		{Kind: "invoice", Amount: ptr(4.0), Currency: "EUR", Reference: Reference{"billing", "fixture:eur"}},
		{Kind: "invoice", Reference: Reference{"billing", "fixture:unknown"}},
	}
	got, err := Summarize(r)
	if err != nil {
		t.Fatal(err)
	}
	// A measured zero is a known zero, not an unknown.
	if value, ok := got.Cost.Actual.Known["USD"]; !ok || value != 0 {
		t.Fatalf("a measured zero was lost or read as unknown: %+v", got.Cost.Actual)
	}
	if got.Cost.Actual.Known["EUR"] != 4 {
		t.Fatalf("currencies were mixed: %+v", got.Cost.Actual)
	}
	// A charge with no amount leaves the actual total partial.
	if !got.Cost.Actual.Partial {
		t.Fatal("a missing charge amount did not mark the actual total partial")
	}
	// No attempt carries a rate, so every estimate stays partial and separate
	// from the actual charges above.
	if len(got.Cost.Estimated.Known) != 0 || !got.Cost.Estimated.Partial {
		t.Fatalf("actual charges leaked into the estimate: %+v", got.Cost.Estimated)
	}
}

// DI17: an update cannot drop a failed, cancelled, or incomplete attempt.
func TestAssessmentDelegatedHistory(t *testing.T) {
	s := Store{Home: t.TempDir(), Root: t.TempDir()}
	r := orchestrated(s.Root)
	for _, state := range []string{"failed", "cancelled", "incomplete"} {
		r.Attempts = append(r.Attempts, Attempt{
			AttemptID: "replaced-" + state, ChunkID: "1", Role: "implementation",
			SessionID: "session-replaced-" + state, Model: "synthetic", Effort: "low",
			State: state,
		})
	}
	if err := s.Record(r); err != nil {
		t.Fatal(err)
	}
	shorter := r
	shorter.Attempts = r.Attempts[:2]
	if err := s.Record(shorter); err == nil {
		t.Fatal("an update dropped the failed and replaced attempts")
	}
	stored, err := s.Read("run-1")
	if err != nil || len(stored.Attempts) != len(r.Attempts) {
		t.Fatalf("the refused update changed the stored attempts: %d %v", len(stored.Attempts), err)
	}
}

// DI18: concurrent intervals contribute their union to wall time and their
// separate spans to effort time.
func TestAssessmentDelegatedIntervals(t *testing.T) {
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	span := func(from, to int) []ObservedInterval {
		return []ObservedInterval{{
			Start:     base.Add(time.Duration(from) * time.Second),
			End:       base.Add(time.Duration(to) * time.Second),
			Reference: Reference{"Bench OTEL", "fixture:" + time.Duration(from).String()},
		}}
	}
	r := orchestrated("")
	r.RepoKey = fixtureRun("").RepoKey
	// Two authors overlap for twenty seconds of their thirty-second spans.
	r.Attempts[0].Intervals = span(0, 30)
	r.Attempts[1].Intervals = span(10, 40)
	got, err := Summarize(r)
	if err != nil {
		t.Fatal(err)
	}
	if got.EffortSeconds != 60 {
		t.Fatalf("effort time lost a concurrent author's own span: %v", got.EffortSeconds)
	}
	if got.WallSeconds == nil || *got.WallSeconds != 40 {
		t.Fatalf("wall time summed concurrent intervals instead of taking their union: %v", got.WallSeconds)
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
		// Measure references embed this store's own paths, which differ between
		// the two runs. The keys do not, so they carry the comparison.
		keys := []string{}
		for i := range r.Attempts {
			for key := range r.Attempts[i].Measures {
				keys = append(keys, r.Attempts[i].AttemptID+" "+key)
			}
			r.Attempts[i].Measures = nil
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

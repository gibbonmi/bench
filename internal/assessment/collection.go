package assessment

import (
	"fmt"
	"github.com/gibbonmi/bench/internal/census"
	"github.com/gibbonmi/bench/internal/otelrecord"
	"github.com/gibbonmi/bench/internal/poolkey"
	"path/filepath"
	"strconv"
	"strings"
)

type Mapping struct {
	AttemptID string `json:"attempt_id"`
	ChunkID   string `json:"chunk_id"`
	Role      string `json:"role"`
}
type Selection struct {
	ID string `json:"id"`
	Mapping
}
type BenchInputs struct {
	AssignmentID   string      `json:"assignment_id"`
	TraceIDs       []Selection `json:"trace_ids,omitempty"`
	CensusEventIDs []Selection `json:"census_event_ids,omitempty"`
}
type HarnessInput struct {
	Mapping
	Path      string `json:"path"`
	Format    string `json:"format"`
	SessionID string `json:"session_id"`
	Epoch     int    `json:"epoch"`
	Sequence  int    `json:"sequence"`
	EventID   string `json:"event_id"`
	Counter   string `json:"counter"`
	Mode      string `json:"mode"`
}

func mapped(r *Run, m Mapping) (*Attempt, error) {
	for i := range r.Attempts {
		a := &r.Attempts[i]
		if a.AttemptID == m.AttemptID {
			if a.ChunkID != m.ChunkID || a.Role != m.Role {
				return nil, fmt.Errorf("ambiguous attempt mapping")
			}
			return a, nil
		}
	}
	return nil, fmt.Errorf("unknown mapped attempt")
}
func diagnostic(r *Run, producer, path, reason string) {
	r.Diagnostics = appendReference(r.Diagnostics, Reference{producer, path + ": " + reason})
}
func measure(a *Attempt, key string, v float64, ref Reference) error {
	if a.Measures == nil {
		a.Measures = map[string]Measure{}
	}
	if old, ok := a.Measures[key]; ok && (old.Reference != ref || (old.Value != nil && *old.Value != v)) {
		return fmt.Errorf("conflicting collected measure %s", key)
	}
	a.Measures[key] = Measure{Value: &v, Reference: ref}
	return nil
}
func Collect(s Store, r Run) (Run, error) {
	batches, err := benchBatches(r)
	if err != nil {
		return r, err
	}
	// One selection ledger spans the whole import, not one batch. A native
	// event that two batches map to two different attempts is ambiguous
	// however it arrived, so that conflict must surface across assignments.
	traces, events := map[string]Mapping{}, map[string]Mapping{}
	for _, b := range batches {
		if err := collectSpans(s, &r, b, traces); err != nil {
			return r, err
		}
		if err := collectCensus(s, &r, b, events); err != nil {
			return r, err
		}
	}
	for _, input := range r.HarnessInputs {
		if err := collectHarness(&r, input); err != nil {
			return r, err
		}
	}
	return r, nil
}

// benchBatches normalizes both Bench-input forms into one list, so collection
// keeps a single path. An empty batch list supplies nothing and is not a second
// form. The caller records only after this returns, so every refusal here
// leaves the stored record unchanged.
func benchBatches(r Run) ([]BenchInputs, error) {
	batches := r.BenchInputBatches
	if r.BenchInputs != nil {
		if len(batches) != 0 {
			return nil, fmt.Errorf("invalid duplicate bench input forms; supply bench_inputs or bench_input_batches")
		}
		batches = []BenchInputs{*r.BenchInputs}
	}
	seen := map[string]bool{}
	for _, b := range batches {
		if _, ok := poolkey.SplitAssignmentSegment(poolkey.AssignmentSegment(b.AssignmentID, b.AssignmentID)); !ok {
			return nil, fmt.Errorf("invalid expected assignment")
		}
		if seen[b.AssignmentID] {
			return nil, fmt.Errorf("duplicate assignment batch %s", b.AssignmentID)
		}
		seen[b.AssignmentID] = true
	}
	return batches, nil
}

func collectCensus(s Store, r *Run, b BenchInputs, selected map[string]Mapping) error {
	if len(b.CensusEventIDs) == 0 {
		return nil
	}
	path := filepath.Join(census.Dir(s.Home, s.Root), b.AssignmentID)
	if err := regularInput(path); err != nil {
		return err
	}
	events, problems, err := census.ReadEvents(s.Home, s.Root, b.AssignmentID)
	if err != nil {
		diagnostic(r, "Bench census", path, err.Error())
	}
	for _, problem := range problems {
		diagnostic(r, "Bench census", path, problem)
	}

	for _, sel := range b.CensusEventIDs {
		if !strings.HasPrefix(sel.ID, b.AssignmentID+":") {
			return fmt.Errorf("foreign census assignment")
		}
		a, err := mapped(r, sel.Mapping)
		if err != nil {
			return err
		}
		fresh, err := selectMapping(selected, sel)
		if err != nil {
			return fmt.Errorf("census: %w", err)
		}
		if !fresh {
			continue
		}
		found := false
		for _, event := range events {
			if event.ID != sel.ID {
				continue
			}
			found = true
			ref := Reference{"Bench census", path + "#" + event.ID}
			if err := measure(a, "raw_commands/"+event.ID, 1, ref); err != nil {
				return err
			}
			a.Evidence = appendReference(a.Evidence, Reference{"Bench census", path + "#" + event.ID + " head=" + event.Head})
		}
		if !found {
			diagnostic(r, "Bench census", path+"#"+sel.ID, "selected event missing")
		}
	}
	return nil
}
func appendReference(refs []Reference, ref Reference) []Reference {
	for _, old := range refs {
		if old == ref {
			return refs
		}
	}
	return append(refs, ref)
}

func collectSpans(s Store, r *Run, b BenchInputs, selected map[string]Mapping) error {
	if len(b.TraceIDs) == 0 {
		return nil
	}
	path := otelrecord.Path(s.Home, s.Root)
	if err := regularInput(path); err != nil {
		return err
	}
	ids := []string{}
	for _, sel := range b.TraceIDs {
		ids = append(ids, sel.ID)
	}
	spans, problems, err := otelrecord.ReadSelected(s.Home, s.Root, ids)
	if err != nil {
		diagnostic(r, "Bench OTEL", path, err.Error())
	}
	for _, problem := range problems {
		diagnostic(r, "Bench OTEL", path, problem)
	}

	for _, sel := range b.TraceIDs {
		a, err := mapped(r, sel.Mapping)
		if err != nil {
			return err
		}
		fresh, err := selectMapping(selected, sel)
		if err != nil {
			return fmt.Errorf("trace: %w", err)
		}
		if !fresh {
			continue
		}
		var group []otelrecord.Span
		assignments := map[string]bool{}
		for _, span := range spans {
			if span.TraceID == sel.ID {
				group = append(group, span)
				id := span.Attributes[otelrecord.AttrAssignmentID]
				if id == "" && strings.HasPrefix(span.Seam, "worktree.") && span.Seam != otelrecord.SeamLanding {
					id = span.Attributes[otelrecord.AttrSubjectID]
				}
				if id != "" {
					assignments[id] = true
				}
			}
		}
		if len(group) == 0 {
			diagnostic(r, "Bench OTEL", path+"#"+sel.ID, "selected trace missing")
			continue
		}
		if len(assignments) != 1 || !assignments[b.AssignmentID] {
			return fmt.Errorf("foreign or ambiguous trace assignment")
		}
		for _, span := range group {
			ref := Reference{"Bench OTEL", path + "#" + span.TraceID + "/" + span.SpanID}
			a.Evidence = appendReference(a.Evidence, ref)
			if span.End.IsZero() {
				diagnostic(r, "Bench OTEL", ref.Native, "unfinished span")
				continue
			}
			if err := measure(a, "elapsed_seconds/"+span.TraceID+"/"+span.SpanID, span.Elapsed().Seconds(), ref); err != nil {
				return err
			}
			if value := span.Attributes[otelrecord.AttrMeasurePathCount]; value != "" {
				n, err := strconv.ParseFloat(value, 64)
				if err != nil || !finite(n) || n < 0 {
					return fmt.Errorf("invalid native path count")
				}
				if err := measure(a, "diff_paths/"+span.TraceID+"/"+span.SpanID, n, ref); err != nil {
					return err
				}
			}
			interval := ObservedInterval{Start: span.Start, End: span.End, Reference: ref}
			seen := false
			for _, old := range a.Intervals {
				if old.Reference == ref {
					if old != interval {
						return fmt.Errorf("conflicting native interval")
					}
					seen = true
				}
			}
			if !seen {
				a.Intervals = append(a.Intervals, interval)
			}
		}
	}
	return nil
}

func selectMapping(selected map[string]Mapping, sel Selection) (bool, error) {
	if prior, ok := selected[sel.ID]; ok {
		if prior != sel.Mapping {
			return false, fmt.Errorf("ambiguous selection mapping")
		}
		return false, nil
	}
	selected[sel.ID] = sel.Mapping
	return true, nil
}

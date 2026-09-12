package assessment

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"
)

// DI13: orchestration is its own role in show and in comparison.
func TestAssessmentOrchestration(t *testing.T) {
	if !slices.Contains(Roles(), "orchestration") {
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

	// The row claims show and comparison alike, so the comparison half needs
	// its own assertion. A role filter added inside the comparison owner would
	// pass the show check above.
	p := comparisonPlan()
	runs := comparisonRuns(s.Root, p)
	for i := range runs {
		pinned := orchestrated(s.Root).Attempts[1]
		pinned.Effort = "high"
		pinned.SessionID = fmt.Sprint("session-orchestrator-", i)
		runs[i].Attempts = append(runs[i].Attempts, pinned)
	}
	compared, err := Compare(p, runs)
	if err != nil {
		t.Fatal(err)
	}
	report, code := renderComparison(compared)
	if code != 0 {
		t.Fatal(report)
	}
	if !strings.Contains(report, "orchestration") {
		t.Fatalf("comparison omits the orchestrator role: %s", report)
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

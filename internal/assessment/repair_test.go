package assessment

import (
	"testing"
	"time"
)

func TestAssessmentRecordReviewCounters(t *testing.T) {
	t.Run("missing cumulative field P1", func(t *testing.T) {
		events := []Event{usageEvent("a", 10), usageEvent("b", 0), usageEvent("c", 15)}
		for i := range events {
			events[i].Mode = CumulativeMode
			events[i].Sequence = i
		}
		events[1].Usage.InputUncached = nil
		got, err := UsageTotal(events)
		if err != nil || got.InputUncached == nil || *got.InputUncached != 15 || len(got.Unknown) == 0 {
			t.Fatalf("expected partial 15, got %+v %v", got, err)
		}
	})
	t.Run("regressing epoch C1", func(t *testing.T) {
		events := []Event{usageEvent("a", 10), usageEvent("b", 3), usageEvent("c", 15)}
		for i := range events {
			events[i].Mode = CumulativeMode
			events[i].Sequence = i
		}
		events[1].Epoch = 1
		got, err := UsageTotal(events)
		if err == nil && (got.InputUncached != nil || len(got.Unknown) == 0) {
			t.Fatalf("regressed epoch appears complete: %+v", got)
		}
	})
	t.Run("independent attempt P2", func(t *testing.T) {
		r := fixtureRun(t.TempDir())
		a := r.Attempts[0]
		a.Usage = []Event{usageEvent("bad-a", 10), usageEvent("bad-b", 3)}
		for i := range a.Usage {
			a.Usage[i].Mode = CumulativeMode
			a.Usage[i].Sequence = i
		}
		b := a
		b.AttemptID = "review"
		b.SessionID = "independent"
		b.Role = "review"
		b.Usage = []Event{usageEvent("clean", 7)}
		b.Usage[0].SessionID = b.SessionID
		r.Attempts = []Attempt{a, b}
		got, err := RunUsage(r)
		if err != nil || len(got[b.AttemptID].Unknown) != 0 || got[b.AttemptID].InputUncached == nil || *got[b.AttemptID].InputUncached != 7 {
			t.Fatalf("clean review contaminated: %+v %v", got, err)
		}
	})
}

func TestAssessmentRecordMeasureProvenance(t *testing.T) {
	for _, kind := range []string{"quality", "unknown tool", "unknown actual", "run time", "attempt time"} {
		t.Run(kind, func(t *testing.T) {
			r := fixtureRun(t.TempDir())
			switch kind {
			case "quality":
				r.Quality = map[string]Measure{"score": {Value: ptr(1.0)}}
			case "unknown tool":
				r.Attempts[0].Cost.Other = []Charge{{Kind: "tool"}}
			case "unknown actual":
				r.Attempts[0].Cost.Actual = []Charge{{Kind: "invoice"}}
			case "run time":
				r.StartedAt = ptr(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
			case "attempt time":
				r.Attempts[0].StartedAt = ptr(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
			}
			if err := Validate(r); err == nil {
				t.Fatal("measure without provenance accepted")
			}
		})
	}
}

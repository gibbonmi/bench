package assessment

import (
	"bytes"
	"fmt"
	"github.com/gibbonmi/bench/internal/poolkey"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func ptr[T any](v T) *T { return &v }

func TestAssessmentRecord(t *testing.T) {
	t.Run("inclusive cache A4", func(t *testing.T) {
		got, err := NormalizeUsage(Usage{InputTotal: ptr(int64(100)), InputCached: ptr(int64(80)), TotalSemantics: "inclusive"})
		if err != nil || got.InputUncached == nil || *got.InputUncached != 20 {
			t.Fatalf("inclusive total must yield 20 uncached: %+v, %v", got, err)
		}
	})
}

func fixtureRun(root string) Run {
	return Run{Version: 1, RunID: "run-1", RepoKey: poolkey.Key(root), Source: "synthetic", Condition: "bench", TaskID: "held-out-1", State: "running", Attempts: []Attempt{{AttemptID: "attempt-1", ChunkID: "1", Role: "implementation", SessionID: "session-1", Model: "synthetic", Effort: "high", State: "failed"}}}
}
func TestAssessmentRecordStorage(t *testing.T) {
	s := Store{Home: t.TempDir(), Root: t.TempDir()}
	r := fixtureRun(s.Root)
	assignment := filepath.Join(poolkey.Pool(s.Home, s.Root), "fixture-assignment")
	if err := os.MkdirAll(assignment, 0700); err != nil {
		t.Fatal(err)
	}
	if err := s.Record(r); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(assignment); err != nil {
		t.Fatal(err)
	}
	got, err := s.Read(r.RunID)
	if err != nil || got.RunID != r.RunID {
		t.Fatalf("A1 durable record lost after assignment removal: %+v %v", got, err)
	}
}

func TestAssessmentRecordPreservesHistory(t *testing.T) {
	s := Store{Home: t.TempDir(), Root: t.TempDir()}
	r := fixtureRun(s.Root)
	if err := s.Record(r); err != nil {
		t.Fatal(err)
	}
	r.Attempts = nil
	if err := s.Record(r); err == nil {
		t.Fatal("A28 deleting a known failed attempt must refuse")
	}
	got, err := s.Read(r.RunID)
	if err != nil || len(got.Attempts) != 1 {
		t.Fatalf("A2 history lost: %+v %v", got, err)
	}
}

func usageEvent(id string, n int64) Event {
	return Event{EventID: id, SessionID: "session-1", Mode: "delta", Counter: "tokens", Reference: Reference{"synthetic", "fixture:" + id}, Usage: Usage{InputUncached: ptr(n), InputCached: ptr(int64(0)), Output: ptr(int64(0))}}
}
func TestAssessmentRecordDuplicateEvents(t *testing.T) {
	e := usageEvent("event-1", 7)
	got, err := UsageTotal([]Event{e, e})
	if err != nil || got.InputUncached == nil || *got.InputUncached != 7 {
		t.Fatalf("A5 duplicate event must count once: %+v %v", got, err)
	}
}

func TestAssessmentRecordEpochs(t *testing.T) {
	a := usageEvent("a", 10)
	a.Mode = "cumulative"
	a.Sequence = 1
	b := usageEvent("b", 15)
	b.Mode = "cumulative"
	b.Sequence = 2
	c := usageEvent("c", 3)
	c.Mode = "cumulative"
	c.Epoch = 1
	c.Sequence = 1
	got, err := UsageTotal([]Event{a, b, c})
	if err != nil || got.InputUncached == nil || *got.InputUncached != 18 {
		t.Fatalf("A6 epochs must total 18: %+v %v", got, err)
	}
}

func TestAssessmentRecordPrices(t *testing.T) {
	a := fixtureRun(t.TempDir()).Attempts[0]
	e := usageEvent("price", 20)
	e.Usage.InputCached = ptr(int64(80))
	e.Usage.Output = ptr(int64(5))
	a.Usage = []Event{e}
	a.Cost.Estimated = &Rates{InputUncached: ptr(2.0), InputCached: ptr(0.5), Output: ptr(4.0), UnitScale: 1, Currency: "USD", Source: "synthetic", Date: "2026-01-01", Conditions: "fixture"}
	got, err := Estimate(a)
	if err != nil || got.Estimated.Known["USD"] != 100 || got.Estimated.Partial {
		t.Fatalf("A8 three rates must yield 100: %+v %v", got, err)
	}
}

func TestAssessmentRecordProvenance(t *testing.T) {
	s := Store{Home: t.TempDir(), Root: t.TempDir()}
	r := fixtureRun(s.Root)
	e := usageEvent("unattributed", 1)
	e.Reference = Reference{}
	r.Attempts[0].Usage = []Event{e}
	if err := s.Record(r); err == nil {
		t.Fatal("A12 usage without provenance must refuse")
	}
	if _, err := s.Read(r.RunID); !os.IsNotExist(err) {
		t.Fatalf("unsafe import stored: %v", err)
	}
}

func TestAssessmentRecordWall(t *testing.T) {
	r := fixtureRun(t.TempDir())
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(10 * time.Second)
	a := r.Attempts[0]
	a.StartedAt = &start
	a.EndedAt = &end
	b := a
	b.AttemptID = "review-2"
	b.Role = "review"
	r.Attempts = []Attempt{a, b}
	got, err := Summarize(r)
	if err != nil || got.WallSeconds == nil || *got.WallSeconds != 10 || got.EffortSeconds != 20 {
		t.Fatalf("A11 concurrent spans: %+v %v", got, err)
	}
}

func TestAssessmentRecordUnknownAndCharges(t *testing.T) {
	t.Run("A7 unknown versus zero", func(t *testing.T) {
		absent, _ := UsageTotal(nil)
		zero, _ := UsageTotal([]Event{usageEvent("zero", 0)})
		if absent.InputUncached != nil || zero.InputUncached == nil || *zero.InputUncached != 0 {
			t.Fatal("unknown and measured zero conflated")
		}
	})
	t.Run("A9 estimate is not billing", func(t *testing.T) {
		a := fixtureRun(t.TempDir()).Attempts[0]
		got, err := Estimate(a)
		if err != nil || !got.Actual.Partial || len(got.Actual.Known) != 0 {
			t.Fatalf("invented actual charge: %+v %v", got, err)
		}
	})
	t.Run("A10 missing tool charge", func(t *testing.T) {
		a := fixtureRun(t.TempDir()).Attempts[0]
		a.Usage = []Event{usageEvent("priced", 1)}
		a.Cost.Estimated = &Rates{InputUncached: ptr(1.0), InputCached: ptr(1.0), Output: ptr(1.0), UnitScale: 1, Currency: "USD", Source: "fixture", Date: "2026-01-01", Conditions: "synthetic"}
		a.Cost.Other = []Charge{{Kind: "tool", Currency: "USD"}}
		got, err := Estimate(a)
		if err != nil || !got.Estimated.Partial {
			t.Fatalf("missing tool rate appears complete: %+v %v", got, err)
		}
	})
	t.Run("A24 every role and failure counts", func(t *testing.T) {
		r := fixtureRun(t.TempDir())
		r.Attempts = nil
		for i, role := range []string{"implementation", "repair", "verification", "review", "diagnostic"} {
			a := fixtureRun(t.TempDir()).Attempts[0]
			a.AttemptID = fmt.Sprint("attempt-", i)
			a.Role = role
			if i%2 == 1 {
				a.State = "cancelled"
			}
			a.Cost.Actual = []Charge{{Kind: "invoice", Amount: ptr(float64(i + 1)), Currency: "USD", Reference: Reference{"billing", "fixture:invoice" + fmt.Sprint(i)}}}
			r.Attempts = append(r.Attempts, a)
		}
		got, err := Summarize(r)
		if err != nil || got.Cost.Actual.Known["USD"] != 15 || got.Cost.Actual.Partial {
			t.Fatalf("role costs lost: %+v %v", got, err)
		}
	})
}
func TestAssessmentRecordUpdates(t *testing.T) {
	for _, kind := range []string{"A14 identity", "A25 append", "A26 fill", "A27 identical", "A29 evidence", "A17 write", "A17 rename"} {
		t.Run(kind, func(t *testing.T) {
			s := Store{Home: t.TempDir(), Root: t.TempDir()}
			r := fixtureRun(s.Root)
			e := usageEvent("event", 5)
			e.Usage.Output = nil
			r.Attempts[0].Usage = []Event{e}
			if err := s.Record(r); err != nil {
				t.Fatal(err)
			}
			path, _ := s.path(r.RunID)
			before, _ := os.ReadFile(path)
			refusal := false
			switch kind {
			case "A14 identity":
				r.Attempts[0].SessionID = "foreign"
				refusal = true
			case "A25 append":
				a := r.Attempts[0]
				a.AttemptID = "next"
				a.SessionID = "session-2"
				a.Usage = nil
				r.Attempts = append(r.Attempts, a)
				r.State = "succeeded"
			case "A26 fill":
				r.Attempts[0].Usage[0].Usage.Output = ptr(int64(9))
			case "A29 evidence":
				r.Attempts[0].Usage[0].Usage.InputUncached = ptr(int64(6))
				refusal = true
			case "A17 write":
				r.State = "failed"
				s.Files.WriteFile = func(string, []byte, os.FileMode) error { return fmt.Errorf("injected write failure") }
				refusal = true
			case "A17 rename":
				r.State = "failed"
				s.Files.Rename = func(string, string) error { return fmt.Errorf("injected rename failure") }
				refusal = true
			}
			err := s.Record(r)
			if (err != nil) != refusal {
				t.Fatalf("refusal=%v err=%v", refusal, err)
			}
			after, _ := os.ReadFile(path)
			if (refusal || kind == "A27 identical") && !bytes.Equal(before, after) {
				t.Fatal("prior bytes changed")
			}
			got, err := s.Read(r.RunID)
			if err != nil {
				t.Fatal(err)
			}
			if kind == "A25 append" && (len(got.Attempts) != 2 || got.Attempts[0].State != "failed") {
				t.Fatal("append lost failed history")
			}
			if kind == "A26 fill" && (got.Attempts[0].Usage[0].Usage.Output == nil || *got.Attempts[0].Usage[0].Usage.Output != 9) {
				t.Fatal("unknown measure did not fill")
			}
		})
	}
}

func TestAssessmentRecordCrossAttemptCounters(t *testing.T) {
	r := fixtureRun(t.TempDir())
	a := r.Attempts[0]
	a.Usage = []Event{usageEvent("first", 10)}
	a.Usage[0].Mode = "cumulative"
	a.Usage[0].Sequence = 1
	a.Cost.Estimated = &Rates{InputUncached: ptr(1.0), InputCached: ptr(1.0), Output: ptr(1.0), UnitScale: 1, Currency: "USD", Source: "fixture", Date: "2026-01-01", Conditions: "synthetic"}
	b := a
	b.AttemptID = "repair"
	b.Role = "repair"
	b.Usage = []Event{usageEvent("second", 15)}
	b.Usage[0].Mode = "cumulative"
	b.Usage[0].Sequence = 2
	r.Attempts = []Attempt{a, b}
	got, err := Summarize(r)
	if err != nil || got.Cost.Estimated.Known["USD"] != 15 {
		t.Fatalf("cumulative work billed twice: %+v %v", got, err)
	}
}

func TestAssessmentRecordLargeCounterConflict(t *testing.T) {
	s := Store{Home: t.TempDir(), Root: t.TempDir()}
	r := fixtureRun(s.Root)
	r.Attempts[0].Usage = []Event{usageEvent("large", 9007199254740992)}
	if err := s.Record(r); err != nil {
		t.Fatal(err)
	}
	r.Attempts[0].Usage[0].Usage.InputUncached = ptr(int64(9007199254740993))
	if err := s.Record(r); err == nil {
		t.Fatal("A29 distinct integer counters above float precision must conflict")
	}
}

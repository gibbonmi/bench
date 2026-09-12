package assessment

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func comparisonPlan() Plan {
	line := map[string]Line{}
	for _, role := range Roles() {
		line[role] = Line{Model: "synthetic", Effort: "high"}
	}
	base := Condition{ID: "baseline", Revision: "revision-1", Lines: line, Harness: "synthetic-v1", Limits: map[string]string{"iterations": "3"}, Capabilities: map[string]string{"core": "1"}, Acceptance: []string{"gate"}, ReviewAxes: []string{"Standards", "Spec", "Coverage"}}
	data, _ := json.Marshal(base)
	var candidate Condition
	json.Unmarshal(data, &candidate)
	candidate.ID = "candidate"
	candidate.Capabilities = map[string]string{"core": "1", "candidate": "1"}
	return Plan{Version: 1, ID: "synthetic-plan", Purpose: "default-change", Approval: Reference{"synthetic user", "fixture approval, no trial execution"}, Budget: Budget{Amount: ptr(100.0), Currency: "USD"}, Tasks: []PlanTask{{ID: "held-out-1", Holdout: true, Repetitions: 2}}, Conditions: []Condition{base, candidate}, Variable: "capability", QualityTolerance: &Tolerance{MaxFailureRate: ptr(0.0), Measures: map[string]Range{"defects": {Max: ptr(0.0)}}}}
}
func comparisonInput(t *testing.T, p Plan) string {
	t.Helper()
	data, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "plan.json")
	if err = os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	return path
}
func TestAssessmentComparison(t *testing.T) {
	t.Run("empty descriptive comparison", func(t *testing.T) {
		s := Store{Home: t.TempDir(), Root: t.TempDir()}
		p := comparisonPlan()
		out, code := Command(s, []string{"compare", "--plan", comparisonInput(t, p), "--runs", ""})
		if code != 0 || !strings.Contains(out, "runs[0]") || !strings.Contains(out, "comparison[") || !strings.Contains(out, "help[0]") {
			t.Fatalf("comparison unavailable: %d %s", code, out)
		}
	})
}

func comparisonRuns(root string, p Plan) []Run {
	var runs []Run
	for _, c := range p.Conditions {
		for rep := 1; rep <= p.Tasks[0].Repetitions; rep++ {
			r := fixtureRun(root)
			r.RunID = fmt.Sprintf("%s-%d", c.ID, rep)
			r.Condition = c.ID
			r.Source = c.Revision
			r.Holdout = p.Tasks[0].Holdout
			r.TaskID = p.Tasks[0].ID
			r.State = "succeeded"
			ref := Reference{"synthetic", "fixture/" + r.RunID}
			start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
			end := start.Add(time.Duration(rep*10) * time.Second)
			a := &r.Attempts[0]
			a.State = "succeeded"
			a.SessionID = r.RunID + "-session"
			a.StartedAt = &start
			a.EndedAt = &end
			a.TimeReference = &ref
			a.Usage = []Event{{EventID: "event-1", SessionID: a.SessionID, Mode: DeltaMode, Counter: "synthetic", Reference: ref, Usage: Usage{InputUncached: ptr(int64(1)), InputCached: ptr(int64(2)), Output: ptr(int64(3))}}}
			a.Cost.Estimated = &Rates{InputUncached: ptr(1.0), InputCached: ptr(2.0), Output: ptr(3.0), UnitScale: 1, Currency: "USD", Source: "synthetic", Date: "2026-01-01", Conditions: "fixture rates, not prices"}
			r.Trial = &Trial{PlanID: p.ID, Repetition: ptr(rep), Harness: c.Harness, Limits: c.Limits, Capabilities: c.Capabilities, Acceptance: c.Acceptance, ReviewAxes: c.ReviewAxes, Reference: ref}
			r.Quality = map[string]Measure{"defects": {Value: ptr(0.0), Reference: ref}}
			runs = append(runs, r)
		}
	}
	return runs
}
func TestAssessmentComparisonEligibility(t *testing.T) {
	for _, kind := range []string{"eligible", "missing budget", "missing approval", "missing quality", "FT311", "one repetition", "not held out", "missing repetition", "unequal acceptance", "unequal review", "quality failure", "failed run", "incomplete run"} {
		t.Run(kind, func(t *testing.T) {
			p := comparisonPlan()
			root := t.TempDir()
			runs := comparisonRuns(root, p)
			refuse := false
			switch kind {
			case "missing budget":
				p.Budget.Amount = nil
				refuse = true
			case "missing approval":
				p.Approval = Reference{}
				refuse = true
			case "missing quality":
				p.QualityTolerance = nil
				refuse = true
			case "FT311":
				p.Tasks[0].ID = "FT311-authoring"
				for i := range runs {
					runs[i].TaskID = p.Tasks[0].ID
				}
			case "one repetition":
				p.Tasks[0].Repetitions = 1
				runs = comparisonRuns(root, p)
			case "not held out":
				p.Tasks[0].Holdout = false
				for i := range runs {
					runs[i].Holdout = false
				}
			case "missing repetition":
				runs = runs[1:]
			case "unequal acceptance":
				p.Conditions[1].Acceptance = []string{"different gate"}
				for i := range runs {
					if runs[i].Condition == p.Conditions[1].ID {
						runs[i].Trial.Acceptance = p.Conditions[1].Acceptance
					}
				}
			case "unequal review":
				p.Conditions[1].ReviewAxes = []string{"Standards", "Spec"}
				for i := range runs {
					if runs[i].Condition == p.Conditions[1].ID {
						runs[i].Trial.ReviewAxes = p.Conditions[1].ReviewAxes
					}
				}
			case "quality failure":
				m := runs[2].Quality["defects"]
				m.Value = ptr(1.0)
				runs[2].Quality["defects"] = m
			case "failed run":
				runs[2].State = "failed"
			case "incomplete run":
				runs[2].Attempts[0].EndedAt = nil
			}
			result, err := Compare(p, runs)
			if refuse {
				if err == nil {
					t.Fatal("invalid plan accepted")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if result.Eligible != (kind == "eligible") {
				t.Fatalf("eligibility wrong: %+v", result.Reasons)
			}
		})
	}
}
func TestAssessmentComparisonFixedInputs(t *testing.T) {
	for _, field := range []string{"task", "revision", "model", "effort", "harness", "limits"} {
		t.Run(field, func(t *testing.T) {
			p := comparisonPlan()
			runs := comparisonRuns(t.TempDir(), p)
			r := &runs[2]
			switch field {
			case "task":
				r.TaskID = "unplanned"
			case "revision":
				r.Source = "revision-2"
			case "model":
				r.Attempts[0].Model = "other"
			case "effort":
				r.Attempts[0].Effort = "other"
			case "harness":
				r.Trial.Harness = "other"
			case "limits":
				r.Trial.Limits = map[string]string{"iterations": "30"}
			}
			if _, err := Compare(p, runs); err == nil {
				t.Fatalf("undeclared %s difference accepted", field)
			}
		})
	}
}
func TestAssessmentComparisonCausal(t *testing.T) {
	for _, kind := range []string{"valid", "missing no-bench", "missing current-bench", "missing changed-capability", "fourth arm", "two capabilities"} {
		t.Run(kind, func(t *testing.T) {
			p := comparisonPlan()
			p.Purpose = "kit-causal"
			p.Conditions[0].ID = "current-bench"
			p.Conditions[1].ID = "changed-capability"
			none := p.Conditions[0]
			none.ID = "no-bench"
			none.Capabilities = nil
			p.Conditions = append(p.Conditions, none)
			if strings.HasPrefix(kind, "missing ") {
				id := strings.TrimPrefix(kind, "missing ")
				for i, c := range p.Conditions {
					if c.ID == id {
						p.Conditions = append(p.Conditions[:i], p.Conditions[i+1:]...)
						break
					}
				}
			}
			if kind == "fourth arm" {
				other := none
				other.ID = "unplanned"
				p.Conditions = append(p.Conditions, other)
			}
			if kind == "two capabilities" {
				p.Conditions[1].Capabilities = map[string]string{"core": "1", "candidate": "1", "second": "1"}
			}
			result, err := Compare(p, comparisonRuns(t.TempDir(), p))
			if strings.HasPrefix(kind, "missing ") || kind == "fourth arm" {
				if err == nil {
					t.Fatal("invalid causal arms accepted")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if result.Eligible != (kind == "valid") {
				t.Fatalf("causal eligibility wrong: %+v", result.Reasons)
			}
		})
	}
}
func TestAssessmentComparisonReport(t *testing.T) {
	s := Store{Home: t.TempDir(), Root: t.TempDir()}
	p := comparisonPlan()
	runs := comparisonRuns(s.Root, p)
	runs[0].State = "failed"
	runs[1].Attempts[0].EndedAt = nil
	for i, role := range Roles() {
		a := runs[0].Attempts[0]
		a.AttemptID = fmt.Sprintf("extra-%d", i)
		a.Role = role
		a.State = "failed"
		a.Usage = nil
		a.Cost = Cost{}
		runs[0].Attempts = append(runs[0].Attempts, a)
	}
	ids := []string{}
	for _, r := range runs {
		if err := s.Record(r); err != nil {
			t.Fatal(err)
		}
		ids = append(ids, r.RunID)
	}
	out, code := Command(s, []string{"compare", "--plan", comparisonInput(t, p), "--runs", strings.Join(ids, ",")})
	if code != 0 {
		t.Fatal(out)
	}
	for _, want := range []string{"variation", "stddev", "failed", "incomplete", "estimated", "actual", "partial", "quality", "repair", "review", "verification", "diagnostic", "input_uncached", "input_cached", "output"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %s: %s", want, out)
		}
	}
	d, err := distribution([]*float64{ptr(10.0), ptr(20.0), nil})
	if err != nil || d.Known != 2 || d.Unknown != 1 || d.Mean == nil || *d.Mean != 15 || d.StdDev == nil || *d.StdDev != 5 {
		t.Fatalf("variation wrong: %+v %v", d, err)
	}
}

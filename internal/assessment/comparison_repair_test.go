package assessment

import (
	"fmt"
	"strings"
	"testing"
)

func TestAssessmentComparisonReviewEligibility(t *testing.T) {
	for _, kind := range []string{"incomplete state", "incomplete attempt", "running attempt", "no treatment", "pilot", "descriptive"} {
		t.Run(kind, func(t *testing.T) {
			p := comparisonPlan()
			if kind == "no treatment" {
				p.Conditions[1].Capabilities = p.Conditions[0].Capabilities
			}
			if kind == "pilot" || kind == "descriptive" {
				p.Purpose = kind
			}
			runs := comparisonRuns(t.TempDir(), p)
			if kind == "incomplete state" {
				runs[0].State = "incomplete"
			}
			if kind == "incomplete attempt" {
				runs[0].Attempts[0].State = "incomplete"
			}
			if kind == "running attempt" {
				runs[0].Attempts[0].State = "running"
			}
			got, err := Compare(p, runs)
			if err != nil || got.Eligible {
				t.Fatalf("ineligible evidence accepted: %v %v", got.Reasons, err)
			}
		})
	}
}
func TestAssessmentComparisonPlanDifferences(t *testing.T) {
	for _, field := range []string{"revision", "model", "effort", "harness", "limits"} {
		t.Run(field, func(t *testing.T) {
			p := comparisonPlan()
			c := &p.Conditions[1]
			switch field {
			case "revision":
				c.Revision = "other"
			case "model":
				c.Lines["implementation"] = Line{"other", "high"}
			case "effort":
				c.Lines["implementation"] = Line{"synthetic", "low"}
			case "harness":
				c.Harness = "other"
			case "limits":
				c.Limits = map[string]string{"iterations": "9"}
			}
			if err := ValidatePlan(p); err == nil {
				t.Fatal("undeclared plan difference accepted")
			}
		})
	}
}
func TestAssessmentComparisonRequiredQuality(t *testing.T) {
	p := comparisonPlan()
	p.QualityTolerance = nil
	if err := ValidatePlan(p); err == nil {
		t.Fatal("missing quality tolerance accepted")
	}
}
func TestAssessmentComparisonTrialBinding(t *testing.T) {
	for _, kind := range []string{"plan", "reference", "repetition low", "repetition high", "capability", "acceptance", "review"} {
		t.Run(kind, func(t *testing.T) {
			p := comparisonPlan()
			runs := comparisonRuns(t.TempDir(), p)
			tr := runs[0].Trial
			switch kind {
			case "plan":
				tr.PlanID = "other"
			case "reference":
				tr.Reference = Reference{}
			case "repetition low":
				tr.Repetition = ptr(0)
			case "repetition high":
				tr.Repetition = ptr(3)
			case "capability":
				tr.Capabilities = map[string]string{"wrong": "1"}
			case "acceptance":
				tr.Acceptance = []string{"other"}
			case "review":
				tr.ReviewAxes = []string{"Standards"}
			}
			got, err := Compare(p, runs)
			if kind == "acceptance" || kind == "review" {
				if err != nil || got.Eligible {
					t.Fatalf("assurance accepted: %v %v", got.Reasons, err)
				}
			} else if err == nil {
				t.Fatal("invalid trial accepted")
			}
		})
	}
}
func TestAssessmentComparisonQualityBoundaries(t *testing.T) {
	for _, v := range []float64{1, 2, 4, 5} {
		t.Run(fmt.Sprint(v), func(t *testing.T) {
			p := comparisonPlan()
			p.QualityTolerance.Measures = map[string]Range{"score": {Min: ptr(2.0), Max: ptr(4.0)}}
			runs := comparisonRuns(t.TempDir(), p)
			for i := range runs {
				runs[i].Quality = map[string]Measure{"score": {Value: ptr(3.0), Reference: Reference{"synthetic", "score"}}}
			}
			runs[0].Quality["score"] = Measure{Value: &v, Reference: Reference{"synthetic", "score"}}
			got, err := Compare(p, runs)
			if err != nil || got.Eligible != (v >= 2 && v <= 4) {
				t.Fatalf("boundary %v: %v %v", v, got.Reasons, err)
			}
		})
	}
}
func TestAssessmentComparisonSecondTask(t *testing.T) {
	p := comparisonPlan()
	root := t.TempDir()
	runs := comparisonRuns(root, p)
	second := PlanTask{ID: "second-task", Holdout: true, Repetitions: 2}
	p.Tasks = append(p.Tasks, second)
	copyPlan := p
	copyPlan.Tasks = []PlanTask{second}
	more := comparisonRuns(root, copyPlan)
	for i := range more {
		more[i].RunID = "second-" + more[i].RunID
	}
	runs = append(runs, more[1:]...)
	got, err := Compare(p, runs)
	if err != nil || got.Eligible || !strings.Contains(strings.Join(got.Reasons, " "), second.ID) {
		t.Fatalf("missing second task ignored: %v %v", got.Reasons, err)
	}
}
func TestAssessmentComparisonCardinality(t *testing.T) {
	p := comparisonPlan()
	p.Tasks = nil
	p.Conditions = nil
	for i := 0; i < 1000; i++ {
		p.Tasks = append(p.Tasks, PlanTask{ID: fmt.Sprint("task-", i), Holdout: true, Repetitions: 2})
		c := comparisonPlan().Conditions[0]
		c.ID = fmt.Sprint("condition-", i)
		p.Conditions = append(p.Conditions, c)
	}
	got, err := Compare(p, nil)
	if err != nil || got.Eligible || len(got.Reasons) > len(p.Conditions)+1 {
		t.Fatalf("unbounded missing-repetition expansion: %d %v", len(got.Reasons), err)
	}
}

func TestAssessmentComparisonQualityDiagnostic(t *testing.T) {
	for _, kind := range []string{"bounds", "unknown"} {
		t.Run(kind, func(t *testing.T) {
			p := comparisonPlan()
			p.QualityTolerance.Measures = map[string]Range{}
			p.QualityTolerance.Measures["z"] = Range{Max: ptr(0.0)}
			p.QualityTolerance.Measures["a"] = Range{Max: ptr(0.0)}
			runs := comparisonRuns(t.TempDir(), p)
			for i := range runs {
				runs[i].Quality = map[string]Measure{"a": {Value: ptr(0.0), Reference: Reference{"synthetic", "a"}}, "z": {Value: ptr(0.0), Reference: Reference{"synthetic", "z"}}}
			}
			if kind == "unknown" {
				runs[0].Quality = nil
			} else {
				for name, m := range runs[0].Quality {
					m.Value = ptr(1.0)
					runs[0].Quality[name] = m
				}
			}
			for i := 0; i < 32; i++ {
				got, err := Compare(p, runs)
				if err != nil || got.Eligible || len(got.Reasons) != 1 || !strings.HasSuffix(got.Reasons[0], "/a") {
					t.Fatalf("quality diagnostic unstable or unbounded: %v %v", got.Reasons, err)
				}
			}
		})
	}
}

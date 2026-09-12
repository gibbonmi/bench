package assessment

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestAssessmentCommandDefaults(t *testing.T) {
	s := Store{Home: t.TempDir(), Root: t.TempDir()}
	binding, err := os.ReadFile("../../.bench/lines.env")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(s.Root, ".bench", "lines.env")
	if err = os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, binding, 0600); err != nil {
		t.Fatal(err)
	}
	p := comparisonPlan()
	r := comparisonRuns(s.Root, p)[0]
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	input := filepath.Join(t.TempDir(), "run.json")
	if err = os.WriteFile(input, data, 0600); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"record", "--input", input}, {"list"}, {"show", r.RunID}, {"compare", "--plan", comparisonInput(t, p), "--runs", r.RunID}} {
		if out, code := Command(s, args); code != 0 {
			t.Fatalf("%v: %d %s", args, code, out)
		}
		got, err := os.ReadFile(path)
		if err != nil || !bytes.Equal(got, binding) {
			t.Fatalf("%s changed model defaults: %v", args[0], err)
		}
	}
}
func TestAssessmentComparisonLargeRepetitions(t *testing.T) {
	p := comparisonPlan()
	p.Tasks[0].Repetitions = int(^uint(0) >> 1)
	result, err := Compare(p, nil)
	if err != nil || result.Eligible {
		t.Fatalf("large incomplete plan: %+v %v", result, err)
	}
}
func TestAssessmentComparisonDuplicateTrials(t *testing.T) {
	for _, kind := range []string{"run", "slot", "missing provenance", "unknown repetition"} {
		t.Run(kind, func(t *testing.T) {
			p := comparisonPlan()
			runs := comparisonRuns(t.TempDir(), p)
			switch kind {
			case "run":
				runs = append(runs, runs[0])
			case "slot":
				runs[1].Trial.Repetition = ptr(1)
			case "missing provenance":
				runs[0].Trial = nil
			case "unknown repetition":
				runs[0].Trial.Repetition = nil
			}
			result, err := Compare(p, runs)
			if kind == "run" || kind == "slot" {
				if err == nil {
					t.Fatal("duplicate trial accepted")
				}
			} else if err != nil || result.Eligible {
				t.Fatalf("unknown evidence became eligible: %+v %v", result, err)
			}
		})
	}
}
func TestAssessmentComparisonDeclaredVariables(t *testing.T) {
	for _, variable := range []string{"revision", "model", "effort", "harness", "limits", "capability"} {
		t.Run(variable, func(t *testing.T) {
			p := comparisonPlan()
			p.Variable = variable
			p.Conditions[1].Capabilities = p.Conditions[0].Capabilities
			c := &p.Conditions[1]
			switch variable {
			case "revision":
				c.Revision = "revision-2"
			case "model":
				c.Lines["implementation"] = Line{Model: "candidate", Effort: "high"}
			case "effort":
				c.Lines["implementation"] = Line{Model: "synthetic", Effort: "medium"}
			case "harness":
				c.Harness = "candidate-v2"
			case "limits":
				c.Limits = map[string]string{"iterations": "5"}
			case "capability":
				c.Capabilities = map[string]string{"core": "1", "candidate": "1"}
			}
			runs := comparisonRuns(t.TempDir(), p)
			for i := range runs {
				if runs[i].Condition == c.ID {
					runs[i].Attempts[0].Model = c.Lines["implementation"].Model
					runs[i].Attempts[0].Effort = c.Lines["implementation"].Effort
				}
			}
			result, err := Compare(p, runs)
			if err != nil || !result.Eligible {
				t.Fatalf("declared variable rejected: %v %+v", err, result.Reasons)
			}
		})
	}
}

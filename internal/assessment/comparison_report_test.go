package assessment

import (
	"encoding/json"
	"fmt"
	"github.com/gibbonmi/bench/internal/axi/axitest"
	"strings"
	"testing"
)

func TestAssessmentComparisonUnsafePlan(t *testing.T) {
	for _, kind := range []string{"fifo", "symlink", "parent-link", "oversized", "duplicate-key", "unknown-field"} {
		t.Run(kind, func(t *testing.T) {
			p := comparisonPlan()
			data, err := json.Marshal(p)
			if err != nil {
				t.Fatal(err)
			}
			input := unsafeAssessmentInput(t, data, kind)
			s := Store{Home: t.TempDir(), Root: t.TempDir()}
			out, code := Command(s, []string{"compare", "--plan", input, "--runs", ""})
			if code != 1 {
				t.Fatalf("unsafe plan accepted: %d %s", code, out)
			}
		})
	}
}
func TestAssessmentComparisonCostTotals(t *testing.T) {
	p := comparisonPlan()
	s := Store{Home: t.TempDir(), Root: t.TempDir()}
	runs := comparisonRuns(s.Root, p)
	for i := range runs {
		base := runs[i].Attempts[0]
		runs[i].Attempts = nil
		for j, role := range Roles() {
			a := base
			a.AttemptID = fmt.Sprint("role-", j)
			a.Role = role
			a.SessionID = base.SessionID + role
			a.Usage = append([]Event(nil), base.Usage...)
			a.Usage[0].SessionID = a.SessionID
			a.State = []string{"failed", "cancelled", "succeeded", "incomplete", "running"}[j]
			a.Cost.Actual = []Charge{{Kind: "billing", Amount: ptr(float64(j + 1)), Currency: "USD", Reference: Reference{"synthetic billing", "fixture"}}, {Kind: "billing", Amount: ptr(float64(j + 2)), Currency: "EUR", Reference: Reference{"synthetic billing", "fixture EUR"}}}
			a.Cost.Other = []Charge{{Kind: "tool", Amount: ptr(float64(j + 3)), Currency: "EUR", Reference: Reference{"synthetic pricing", "fixture"}}}
			if j == 0 && i%2 == 0 {
				a.Cost.Actual = append(a.Cost.Actual, Charge{Kind: "unknown", Reference: Reference{"synthetic billing", "unknown"}})
				a.Cost.Other = append(a.Cost.Other, Charge{Kind: "unknown", Reference: Reference{"synthetic pricing", "unknown"}})
			}
			runs[i].Attempts = append(runs[i].Attempts, a)
		}
	}
	got, err := Compare(p, runs)
	if err != nil {
		t.Fatal(err)
	}
	out, code := renderComparison(got)
	if code != 0 {
		t.Fatal(out)
	}
	doc, err := axitest.DecodeDocument(out)
	if err != nil {
		t.Fatal(err)
	}
	rows, err := doc.Rows("conditions")
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range rows {
		cells := row.(map[string]any)
		var cost CostSummary
		if err = json.Unmarshal([]byte(cells["cost"].(string)), &cost); err != nil {
			t.Fatal(err)
		}
		if cost.Estimated.Known["USD"] != 140 || cost.Estimated.Known["EUR"] != 50 || cost.Actual.Known["USD"] != 30 || cost.Actual.Known["EUR"] != 40 || !cost.Estimated.Partial || !cost.Actual.Partial {
			t.Fatalf("condition totals omit costs: %+v", cost)
		}
	}
	if strings.Contains(out, "time_reference") || strings.Contains(out, "total_semantics") || strings.Contains(out, "synthetic billing") {
		t.Fatal("default report leaked native evidence body")
	}
	roleRows, err := doc.Rows("roles")
	if err != nil || len(roleRows) != 10 {
		t.Fatalf("role aggregation missing: %d %v", len(roleRows), err)
	}
}

package assessment

import (
	"encoding/json"
	"fmt"
	"github.com/gibbonmi/bench/internal/axi/axitest"
	"math"
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
	rows, err := doc.Rows("costs")
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]string{"estimated/USD": "140", "estimated/EUR": "50", "actual/USD": "30", "actual/EUR": "40"}
	seen := map[string]int{}
	for _, row := range rows {
		cells := row.(map[string]any)
		if cells["role"] != "" {
			continue
		}
		key := cells["kind"].(string) + "/" + cells["currency"].(string)
		if cells["known_value"] != expected[key] || cells["partial"] != "true" {
			t.Fatalf("condition totals omit costs: %+v", cells)
		}
		seen[cells["condition"].(string)]++
	}
	for _, c := range p.Conditions {
		if seen[c.ID] != 4 {
			t.Fatalf("missing condition cost rows: %v", seen)
		}
	}
	for _, block := range doc.Blocks {
		values, err := doc.Rows(block)
		if err != nil {
			t.Fatal(err)
		}
		for _, row := range values {
			for _, v := range row.(map[string]any) {
				if text, ok := v.(string); ok && (strings.HasPrefix(text, "{") || strings.HasPrefix(text, "[")) {
					t.Fatalf("opaque JSON value in %s: %s", block, text)
				}
			}
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

func TestAssessmentComparisonOverflow(t *testing.T) {
	for _, kind := range []string{"same role", "different role"} {
		t.Run(kind, func(t *testing.T) {
			p := comparisonPlan()
			s := Store{Home: t.TempDir(), Root: t.TempDir()}
			runs := comparisonRuns(s.Root, p)
			for i := 0; i < 2; i++ {
				runs[i].Attempts[0].Cost.Actual = []Charge{{Kind: "billing", Amount: ptr(math.MaxFloat64), Currency: "USD", Reference: Reference{"synthetic billing", "finite maximum"}}}
			}
			if kind == "different role" {
				runs[1].Attempts[0].Role = "review"
			}
			ids := []string{}
			for _, r := range runs {
				if err := s.Record(r); err != nil {
					t.Fatal(err)
				}
				ids = append(ids, r.RunID)
			}
			out, code := Command(s, []string{"compare", "--plan", comparisonInput(t, p), "--runs", strings.Join(ids, ",")})
			if code != 1 || !strings.HasPrefix(out, "error:") || !strings.Contains(out, "overflow") {
				t.Fatalf("aggregate overflow not refused: %d %s", code, out)
			}
		})
	}
}

func TestAssessmentComparisonVariation(t *testing.T) {
	p := comparisonPlan()
	runs := comparisonRuns(t.TempDir(), p)
	runs[0].Quality["a"] = Measure{Value: ptr(2.0), Reference: Reference{"synthetic", "a"}}
	runs[1].Quality["b"] = Measure{Value: ptr(4.0), Reference: Reference{"synthetic", "b"}}
	runs[0].Attempts[0].Cost.Actual = []Charge{{Kind: "billing", Amount: ptr(7.0), Currency: "USD", Reference: Reference{"synthetic", "USD"}}}
	runs[1].Attempts[0].Cost.Actual = []Charge{{Kind: "billing", Amount: ptr(11.0), Currency: "EUR", Reference: Reference{"synthetic", "EUR"}}}
	result, err := Compare(p, runs)
	if err != nil {
		t.Fatal(err)
	}
	out, code := renderComparison(result)
	if code != 0 {
		t.Fatal(out)
	}
	doc, err := axitest.DecodeDocument(out)
	if err != nil {
		t.Fatal(err)
	}
	rows, err := doc.Rows("variation")
	if err != nil {
		t.Fatal(err)
	}
	values := map[string]string{}
	for _, row := range rows {
		c := row.(map[string]any)
		if c["condition"] == "baseline" {
			values[c["metric"].(string)+"/"+c["statistic"].(string)] = c["value"].(string)
		}
	}
	for key, want := range map[string]string{"wall_seconds/known": "2", "wall_seconds/unknown": "0", "wall_seconds/mean": "15", "wall_seconds/stddev": "5", "quality/a/known": "1", "quality/a/unknown": "1", "quality/a/mean": "2", "actual_known/USD/known": "1", "actual_known/USD/unknown": "1", "actual_known/USD/mean": "7", "actual_known/EUR/mean": "11"} {
		if values[key] != want {
			t.Fatalf("variation %s = %s, want %s", key, values[key], want)
		}
	}
	rows, err = doc.Rows("usage")
	if err != nil {
		t.Fatal(err)
	}
	seen := 0
	for _, row := range rows {
		c := row.(map[string]any)
		if c["condition"] != "baseline" {
			continue
		}
		want := map[string]string{"input_uncached": "2", "input_cached": "4", "output": "6"}[c["category"].(string)]
		if c["known_value"] != want || c["partial"] != "false" {
			t.Fatalf("typed usage wrong: %v", c)
		}
		seen++
	}
	if seen != 3 {
		t.Fatal("missing usage category")
	}
	rows, err = doc.Rows("costs")
	if err != nil {
		t.Fatal(err)
	}
	unknown := false
	for _, row := range rows {
		c := row.(map[string]any)
		if c["condition"] == "candidate" && c["role"] == "" && c["kind"] == "actual" {
			unknown = c["currency"] == "" && c["known_value"] == "unknown" && c["partial"] == "true"
		}
	}
	if !unknown {
		t.Fatal("unknown actual charge became a known zero")
	}
}

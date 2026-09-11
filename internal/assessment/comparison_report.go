package assessment

import (
	"fmt"
	"github.com/gibbonmi/bench/internal/bounds"
	"math"
	"sort"
	"strings"
)

type Distribution struct {
	Known   int      `json:"known"`
	Unknown int      `json:"unknown"`
	Min     *float64 `json:"min,omitempty"`
	Max     *float64 `json:"max,omitempty"`
	Mean    *float64 `json:"mean,omitempty"`
	StdDev  *float64 `json:"stddev,omitempty"`
}

func distribution(values []*float64) (Distribution, error) {
	d := Distribution{}
	mean := 0.0
	for _, p := range values {
		if p == nil {
			d.Unknown++
			continue
		}
		v := *p
		d.Known++
		mean += (v - mean) / float64(d.Known)
		if d.Min == nil || v < *d.Min {
			n := v
			d.Min = &n
		}
		if d.Max == nil || v > *d.Max {
			n := v
			d.Max = &n
		}
	}
	if d.Known == 0 {
		return d, nil
	}
	d.Mean = &mean
	variance := 0.0
	for _, p := range values {
		if p != nil {
			delta := *p - mean
			variance += delta * delta / float64(d.Known)
		}
	}
	deviation := math.Sqrt(variance)
	if !finite(mean) || !finite(deviation) {
		return d, fmt.Errorf("comparison variation overflow")
	}
	d.StdDev = &deviation
	return d, nil
}
func compareCommand(s Store, path, ids string) (string, int) {
	var p Plan
	if err := readJSON(path, &p); err != nil {
		return fail(err)
	}
	if err := ValidatePlan(p); err != nil {
		return fail(err)
	}
	runs := []Run{}
	retained := int64(0)
	if ids != "" {
		for _, id := range strings.Split(ids, ",") {
			r, err := s.Read(id)
			if err != nil {
				return fail(err)
			}
			retained += int64(len(encoded(r)))
			if retained > bounds.ControlRecordLimit {
				return fail(fmt.Errorf("selected records exceed comparison bound"))
			}
			runs = append(runs, r)
		}
	}
	report, err := Compare(p, runs)
	if err != nil {
		return fail(err)
	}
	return renderComparison(report)
}
func renderComparison(report Comparison) (string, int) {
	roles, err := comparisonRoles(report)
	if err != nil {
		return fail(err)
	}
	conditions, err := comparisonConditions(report)
	if err != nil {
		return fail(err)
	}
	runRows, outcomes, roleRows, usageRows, costRows, conditionRows, variationRows := [][]string{}, [][]string{}, [][]string{}, [][]string{}, [][]string{}, [][]string{}, [][]string{}
	for _, r := range report.Runs {
		runRows = append(runRows, []string{r.Run.RunID, r.Run.Condition, r.Run.TaskID, r.Run.State})
		outcomes = append(outcomes, []string{r.Run.RunID, "incomplete", fmt.Sprint(r.Summary.Incomplete || r.Run.State == "incomplete"), "true"})
		for _, name := range sortedKeys(r.Run.Quality) {
			m := r.Run.Quality[name]
			outcomes = append(outcomes, []string{r.Run.RunID, "quality/" + name, number(m.Value), fmt.Sprint(m.Value != nil)})
		}
	}
	for _, r := range roles {
		for _, state := range sortedKeys(r.states) {
			roleRows = append(roleRows, []string{r.condition, r.role, state, fmt.Sprint(r.states[state])})
		}
		for i, n := range []*int64{r.usage.InputUncached, r.usage.InputCached, r.usage.Output} {
			value := "unknown"
			if n != nil {
				value = fmt.Sprint(*n)
			}
			usageRows = append(usageRows, []string{r.condition, r.role, []string{"input_uncached", "input_cached", "output"}[i], value, fmt.Sprint(n == nil || len(r.usage.Unknown) > 0)})
		}
		costRows = appendCostRows(costRows, r.condition, r.role, r.cost)
	}
	for _, c := range conditions {
		conditionRows = append(conditionRows, []string{c.condition, fmt.Sprint(c.runs)})
		costRows = appendCostRows(costRows, c.condition, "", c.cost)
		for _, metric := range sortedKeys(c.metrics) {
			d, err := distribution(c.metrics[metric])
			if err != nil {
				return fail(err)
			}
			d.Unknown = c.runs - d.Known
			variationRows = append(variationRows, []string{c.condition, metric, "known", fmt.Sprint(d.Known)}, []string{c.condition, metric, "unknown", fmt.Sprint(d.Unknown)})
			for _, stat := range []struct {
				name  string
				value *float64
			}{{"min", d.Min}, {"max", d.Max}, {"mean", d.Mean}, {"stddev", d.StdDev}} {
				variationRows = append(variationRows, []string{c.condition, metric, stat.name, number(stat.value)})
			}
		}
	}
	reasonRows := [][]string{}
	sort.Strings(report.Reasons)
	for _, reason := range report.Reasons {
		reasonRows = append(reasonRows, []string{reason})
	}
	tables := []struct {
		name   string
		fields []string
		rows   [][]string
	}{
		{"comparison", []string{"plan", "purpose", "eligible", "limits"}, [][]string{{report.Plan.ID, report.Plan.Purpose, fmt.Sprint(report.Eligible), "eligibility is not adoption approval; approval reference is unverified; inspect native details with bench assessment show <run-id>"}}},
		{"runs", []string{"run_id", "condition", "task", "state"}, runRows},
		{"outcomes", []string{"run_id", "measure", "value", "known"}, outcomes},
		{"roles", []string{"condition", "role", "state", "attempts"}, roleRows},
		{"usage", []string{"condition", "role", "category", "known_value", "partial"}, usageRows},
		{"costs", []string{"condition", "role", "kind", "currency", "known_value", "partial"}, costRows},
		{"conditions", []string{"condition", "runs"}, conditionRows},
		{"variation", []string{"condition", "metric", "statistic", "value"}, variationRows},
		{"reasons", []string{"reason"}, reasonRows},
	}
	out := ""
	for i, t := range tables {
		part, code := table(t.name, t.fields, t.rows, i == len(tables)-1)
		if code != 0 {
			return part, code
		}
		out += part
	}
	return out, 0
}
func number(v *float64) string {
	if v == nil {
		return "unknown"
	}
	return fmt.Sprint(*v)
}
func sortedKeys[V any](values map[string]V) []string {
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
func appendCostRows(rows [][]string, condition, role string, cost CostSummary) [][]string {
	for _, item := range []struct {
		kind  string
		money Money
	}{{"estimated", cost.Estimated}, {"actual", cost.Actual}} {
		currencies := sortedKeys(item.money.Known)
		if len(currencies) == 0 {
			currencies = []string{""}
		}
		for _, currency := range currencies {
			value := "unknown"
			if n, ok := item.money.Known[currency]; ok {
				value = fmt.Sprint(n)
			}
			rows = append(rows, []string{condition, role, item.kind, currency, value, fmt.Sprint(item.money.Partial)})
		}
	}
	return rows
}

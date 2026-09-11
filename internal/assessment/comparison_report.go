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
	out, code := table("comparison", []string{"plan", "purpose", "eligible", "limits"}, [][]string{{report.Plan.ID, report.Plan.Purpose, fmt.Sprint(report.Eligible), "eligibility is not adoption approval; approval reference is unverified"}}, false)
	if code != 0 {
		return out, code
	}
	rows := [][]string{}
	for _, r := range report.Runs {
		rows = append(rows, []string{r.Run.RunID, r.Run.Condition, r.Run.TaskID, r.Run.State, encoded(r.Summary), encoded(r.Usage), encoded(r.Run.Quality), encoded(r.Run.Attempts)})
	}
	more, code := table("runs", []string{"run_id", "condition", "task", "state", "summary", "usage", "quality", "attempts"}, rows, false)
	if code != 0 {
		return more, code
	}
	out += more
	rows = nil
	for _, c := range report.Plan.Conditions {
		metrics := map[string][]*float64{}
		selected := []ComparedRun{}
		keys := map[string]bool{"wall_seconds": true, "effort_seconds": true}
		totals := CostSummary{Estimated: Money{Known: map[string]float64{}}, Actual: Money{Known: map[string]float64{}}}
		for _, r := range report.Runs {
			if r.Run.Condition != c.ID {
				continue
			}
			selected = append(selected, r)
			for currency := range r.Summary.Cost.Estimated.Known {
				keys["estimated_known/"+currency] = true
			}
			for currency := range r.Summary.Cost.Actual.Known {
				keys["actual_known/"+currency] = true
			}
			for k := range r.Run.Quality {
				keys["quality/"+k] = true
			}
			addMoney(&totals.Estimated, r.Summary.Cost.Estimated)
			addMoney(&totals.Actual, r.Summary.Cost.Actual)
		}
		if len(selected) == 0 {
			totals.Estimated.Partial = true
			totals.Actual.Partial = true
		}
		for _, money := range []Money{totals.Estimated, totals.Actual} {
			for _, value := range money.Known {
				if !finite(value) {
					return fail(fmt.Errorf("comparison cost overflow"))
				}
			}
		}
		for k := range keys {
			for _, r := range selected {
				var value *float64
				switch {
				case k == "wall_seconds":
					value = r.Summary.WallSeconds
				case k == "effort_seconds":
					n := r.Summary.EffortSeconds
					if !r.Summary.Incomplete {
						value = &n
					}
				case strings.HasPrefix(k, "quality/"):
					value = r.Run.Quality[strings.TrimPrefix(k, "quality/")].Value
				default:
					money := r.Summary.Cost.Estimated
					currency := strings.TrimPrefix(k, "estimated_known/")
					if strings.HasPrefix(k, "actual_known/") {
						money = r.Summary.Cost.Actual
						currency = strings.TrimPrefix(k, "actual_known/")
					}
					if n, ok := money.Known[currency]; ok {
						value = &n
					}
				}
				metrics[k] = append(metrics[k], value)
			}
		}
		variation := map[string]Distribution{}
		for k, values := range metrics {
			d, err := distribution(values)
			if err != nil {
				return fail(err)
			}
			variation[k] = d
		}
		rows = append(rows, []string{c.ID, fmt.Sprint(len(selected)), encoded(totals), encoded(variation)})
	}
	more, code = table("conditions", []string{"condition", "runs", "cost", "variation"}, rows, false)
	if code != 0 {
		return more, code
	}
	out += more
	rows = nil
	sort.Strings(report.Reasons)
	for _, reason := range report.Reasons {
		rows = append(rows, []string{reason})
	}
	more, code = table("reasons", []string{"reason"}, rows, true)
	return out + more, code
}

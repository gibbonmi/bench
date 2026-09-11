package assessment

import (
	"fmt"
	"sort"
	"time"
)

type Summary struct {
	Cost          CostSummary `json:"cost"`
	WallSeconds   *float64    `json:"wall_seconds,omitempty"`
	EffortSeconds float64     `json:"effort_seconds"`
	Incomplete    bool        `json:"incomplete"`
}

func Summarize(r Run) (Summary, error) {
	out := Summary{Cost: CostSummary{Estimated: Money{Known: map[string]float64{}}, Actual: Money{Known: map[string]float64{}}}}
	if err := Validate(r); err != nil {
		return out, err
	}
	usages, err := RunUsage(r)
	if err != nil {
		return out, err
	}
	out.Incomplete = len(r.Diagnostics) > 0
	var spans [][2]time.Time
	for _, a := range r.Attempts {
		cost, err := estimateUsage(a, usages[a.AttemptID])
		if err != nil {
			return out, err
		}
		addMoney(&out.Cost.Estimated, cost.Estimated)
		addMoney(&out.Cost.Actual, cost.Actual)
		var observed [][2]time.Time
		for _, span := range a.Intervals {
			observed = append(observed, [2]time.Time{span.Start, span.End})
		}
		if a.StartedAt != nil && a.EndedAt != nil {
			observed = append(observed, [2]time.Time{*a.StartedAt, *a.EndedAt})
		}
		if len(observed) == 0 || (a.StartedAt != nil && a.EndedAt == nil) {
			out.Incomplete = true
		}
		out.EffortSeconds += unionSeconds(observed)
		spans = append(spans, observed...)
	}
	if len(r.Attempts) == 0 {
		out.Cost.Estimated.Partial = true
		out.Cost.Actual.Partial = true
		out.Incomplete = true
	}
	if r.StartedAt != nil && r.EndedAt != nil {
		v := r.EndedAt.Sub(*r.StartedAt).Seconds()
		out.WallSeconds = &v
	} else if len(spans) > 0 {
		total := unionSeconds(spans)
		out.WallSeconds = &total
	} else {
		out.Incomplete = true
	}
	for _, m := range []Money{out.Cost.Estimated, out.Cost.Actual} {
		for _, v := range m.Known {
			if !finite(v) {
				return out, fmt.Errorf("cost overflow")
			}
		}
	}
	return out, nil
}
func addMoney(dst *Money, src Money) {
	dst.Partial = dst.Partial || src.Partial
	for currency, n := range src.Known {
		dst.Known[currency] += n
	}
}

func unionSeconds(spans [][2]time.Time) float64 {
	if len(spans) == 0 {
		return 0
	}
	sort.Slice(spans, func(i, j int) bool { return spans[i][0].Before(spans[j][0]) })
	start, end := spans[0][0], spans[0][1]
	total := 0.0
	for _, p := range spans[1:] {
		if p[0].After(end) {
			total += end.Sub(start).Seconds()
			start, end = p[0], p[1]
		} else if p[1].After(end) {
			end = p[1]
		}
	}
	return total + end.Sub(start).Seconds()
}

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
	var spans [][2]time.Time
	for _, a := range r.Attempts {
		cost, err := estimateUsage(a, usages[a.AttemptID])
		if err != nil {
			return out, err
		}
		addMoney(&out.Cost.Estimated, cost.Estimated)
		addMoney(&out.Cost.Actual, cost.Actual)
		if a.StartedAt == nil || a.EndedAt == nil {
			out.Incomplete = true
			continue
		}
		spans = append(spans, [2]time.Time{*a.StartedAt, *a.EndedAt})
		out.EffortSeconds += a.EndedAt.Sub(*a.StartedAt).Seconds()
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
		sort.Slice(spans, func(i, j int) bool { return spans[i][0].Before(spans[j][0]) })
		start, end := spans[0][0], spans[0][1]
		total := 0.0
		for _, p := range spans[1:] {
			if p[0].After(end) {
				total += end.Sub(start).Seconds()
				start = p[0]
				end = p[1]
			} else if p[1].After(end) {
				end = p[1]
			}
		}
		total += end.Sub(start).Seconds()
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

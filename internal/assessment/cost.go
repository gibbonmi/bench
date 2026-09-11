package assessment

import (
	"fmt"
	"math"
)

type Money struct {
	Known   map[string]float64 `json:"known"`
	Partial bool               `json:"partial"`
}
type CostSummary struct {
	Estimated Money `json:"estimated"`
	Actual    Money `json:"actual"`
}

func Estimate(a Attempt) (CostSummary, error) {
	u, err := UsageTotal(a.Usage)
	if err != nil {
		return CostSummary{}, err
	}
	return estimateUsage(a, u)
}
func estimateUsage(a Attempt, u Usage) (CostSummary, error) {
	out := CostSummary{Estimated: Money{Known: map[string]float64{}}, Actual: Money{Known: map[string]float64{}, Partial: len(a.Cost.Actual) == 0}}
	var err error
	r := a.Cost.Estimated
	if r == nil {
		out.Estimated.Partial = true
	} else {
		if r.Currency == "" || r.Source == "" || r.Date == "" || r.Conditions == "" || !finite(r.UnitScale) || r.UnitScale <= 0 {
			return out, fmt.Errorf("incomplete rate provenance")
		}
		rates := []*float64{r.InputUncached, r.InputCached, r.Output}
		for i, n := range []*int64{u.InputUncached, u.InputCached, u.Output} {
			rate := rates[i]
			if rate != nil && (!finite(*rate) || *rate < 0) {
				return out, fmt.Errorf("invalid rate")
			}
			if n == nil || rate == nil {
				out.Estimated.Partial = true
				continue
			}
			out.Estimated.Known[r.Currency] += float64(*n) * (*rate) / r.UnitScale
		}
		out.Estimated.Partial = out.Estimated.Partial || len(u.Unknown) > 0
	}
	if err = charges(&out.Estimated, a.Cost.Other); err != nil {
		return out, err
	}
	if err = charges(&out.Actual, a.Cost.Actual); err != nil {
		return out, err
	}
	for _, m := range []Money{out.Estimated, out.Actual} {
		for _, v := range m.Known {
			if !finite(v) {
				return out, fmt.Errorf("cost overflow")
			}
		}
	}
	return out, nil
}
func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }
func charges(m *Money, list []Charge) error {
	for _, c := range list {
		if c.Kind == "" {
			return fmt.Errorf("charge kind missing")
		}
		if c.Amount == nil {
			m.Partial = true
			continue
		}
		if c.Currency == "" || !validReference(c.Reference) || !finite(*c.Amount) || *c.Amount < 0 {
			return fmt.Errorf("invalid charge or missing authoritative reference")
		}
		m.Known[c.Currency] += *c.Amount
	}
	return nil
}
func validReference(r Reference) bool { return r.Producer != "" && r.Native != "" }

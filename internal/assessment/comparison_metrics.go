package assessment

type conditionReport struct {
	condition string
	runs      int
	cost      CostSummary
	metrics   map[string][]*float64
}

func comparisonConditions(report Comparison) ([]conditionReport, error) {
	groups := map[string]*conditionReport{}
	for _, c := range report.Plan.Conditions {
		groups[c.ID] = &conditionReport{condition: c.ID, cost: CostSummary{Estimated: Money{Known: map[string]float64{}}, Actual: Money{Known: map[string]float64{}}}, metrics: map[string][]*float64{"wall_seconds": nil, "effort_seconds": nil}}
	}
	for _, r := range report.Runs {
		g := groups[r.Run.Condition]
		g.runs++
		addMoney(&g.cost.Estimated, r.Summary.Cost.Estimated)
		addMoney(&g.cost.Actual, r.Summary.Cost.Actual)
		add := func(name string, v *float64) {
			if _, ok := g.metrics[name]; !ok {
				g.metrics[name] = nil
			}
			if v != nil {
				g.metrics[name] = append(g.metrics[name], v)
			}
		}
		add("wall_seconds", r.Summary.WallSeconds)
		if !r.Summary.Incomplete {
			value := r.Summary.EffortSeconds
			add("effort_seconds", &value)
		}
		for name, m := range r.Run.Quality {
			add("quality/"+name, m.Value)
		}
		for currency, n := range r.Summary.Cost.Estimated.Known {
			add("estimated_known/"+currency, &n)
		}
		for currency, n := range r.Summary.Cost.Actual.Known {
			add("actual_known/"+currency, &n)
		}
	}
	out := []conditionReport{}
	for _, c := range report.Plan.Conditions {
		g := groups[c.ID]
		if g.runs == 0 {
			g.cost.Estimated.Partial = true
			g.cost.Actual.Partial = true
		}
		if err := checkCost(g.cost); err != nil {
			return nil, err
		}
		out = append(out, *g)
	}
	return out, nil
}

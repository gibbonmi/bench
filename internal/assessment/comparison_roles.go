package assessment

type roleAggregate struct {
	states map[string]int
	cost   CostSummary
	events []Event
}

func comparisonRoleRows(report Comparison) ([][]string, error) {
	type key struct{ condition, role string }
	groups := map[key]*roleAggregate{}
	for _, r := range report.Runs {
		for _, a := range r.Run.Attempts {
			k := key{r.Run.Condition, a.Role}
			g := groups[k]
			if g == nil {
				g = &roleAggregate{states: map[string]int{}, cost: CostSummary{Estimated: Money{Known: map[string]float64{}}, Actual: Money{Known: map[string]float64{}}}}
				groups[k] = g
			}
			g.states[a.State]++
			u := r.Usage[a.AttemptID]
			g.events = append(g.events, Event{EventID: r.Run.RunID + "/" + a.AttemptID, Usage: u})
			cost, err := estimateUsage(a, u)
			if err != nil {
				return nil, err
			}
			addMoney(&g.cost.Estimated, cost.Estimated)
			addMoney(&g.cost.Actual, cost.Actual)
		}
	}
	rows := [][]string{}
	for _, c := range report.Plan.Conditions {
		for _, role := range Roles() {
			g := groups[key{c.ID, role}]
			if g == nil {
				continue
			}
			u, err := sumEvents(g.events, nil)
			if err != nil {
				return nil, err
			}
			if len(u.Unknown) > 0 {
				u.Unknown = []string{"partial"}
			}
			if err := checkCost(g.cost); err != nil {
				return nil, err
			}
			metrics := struct {
				Usage Usage       `json:"usage"`
				Cost  CostSummary `json:"cost"`
			}{u, g.cost}
			rows = append(rows, []string{c.ID, role, encoded(g.states), encoded(metrics)})
		}
	}
	return rows, nil
}

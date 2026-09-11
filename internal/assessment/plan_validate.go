package assessment

import (
	"fmt"
	"reflect"
	"slices"
	"sort"
)

func ValidatePlan(p Plan) error {
	if !safeText(reflect.ValueOf(p)) || p.Version != 1 || !safeID.MatchString(p.ID) {
		return fmt.Errorf("invalid comparison plan identity")
	}
	if _, ok := purposes[p.Purpose]; !ok {
		return fmt.Errorf("unknown comparison purpose")
	}
	if !validReference(p.Approval) || p.Budget.Amount == nil || !finite(*p.Budget.Amount) || *p.Budget.Amount < 0 || p.Budget.Currency == "" {
		return fmt.Errorf("plan requires approval reference and budget")
	}
	if _, ok := variableProjections[p.Variable]; !ok {
		return fmt.Errorf("plan requires one declared experimental variable")
	}
	q := p.QualityTolerance
	if q == nil || q.MaxFailureRate == nil || !finite(*q.MaxFailureRate) || *q.MaxFailureRate < 0 || *q.MaxFailureRate > 1 {
		return fmt.Errorf("plan requires a quality tolerance")
	}
	for name, v := range q.Measures {
		if name == "" || (v.Min == nil && v.Max == nil) {
			return fmt.Errorf("invalid quality range")
		}
		for _, bound := range []*float64{v.Min, v.Max} {
			if bound != nil && (!finite(*bound) || *bound < 0) {
				return fmt.Errorf("invalid quality bound")
			}
		}
		if v.Min != nil && v.Max != nil && *v.Min > *v.Max {
			return fmt.Errorf("reversed quality range")
		}
	}
	if len(p.Tasks) == 0 || len(p.Conditions) == 0 {
		return fmt.Errorf("plan requires tasks and conditions")
	}
	tasks := map[string]bool{}
	for _, task := range p.Tasks {
		if task.ID == "" || tasks[task.ID] || task.Repetitions < 1 {
			return fmt.Errorf("invalid or duplicate task")
		}
		tasks[task.ID] = true
	}
	conditions := map[string]bool{}
	for _, c := range p.Conditions {
		if !safeID.MatchString(c.ID) || conditions[c.ID] || c.Revision == "" || c.Harness == "" || len(c.Lines) == 0 || len(c.Limits) == 0 {
			return fmt.Errorf("incomplete or duplicate condition")
		}
		conditions[c.ID] = true
		if _, ok := c.Lines["implementation"]; !ok {
			return fmt.Errorf("condition must pin implementation line")
		}
		for role, line := range c.Lines {
			if !slices.Contains(Roles(), role) || line.Model == "" || line.Effort == "" {
				return fmt.Errorf("invalid pinned line")
			}
		}
		for _, m := range []map[string]string{c.Limits, c.Capabilities} {
			for k, v := range m {
				if k == "" || v == "" {
					return fmt.Errorf("empty condition value")
				}
			}
		}
		if !unique(c.Acceptance) || !unique(c.ReviewAxes) {
			return fmt.Errorf("duplicate or empty assurance obligation")
		}
		if !reflect.DeepEqual(fixedCondition(p.Conditions[0], p.Variable), fixedCondition(c, p.Variable)) {
			return fmt.Errorf("undeclared condition difference")
		}
	}
	if purposes[p.Purpose].causal {
		if _, ok := resolveCausalArms(p.Conditions); !ok || p.Variable != capabilityVariable {
			return fmt.Errorf("kit-causal plan requires exactly three FT231 conditions")
		}
	}
	return nil
}
func unique(values []string) bool {
	seen := map[string]bool{}
	for _, v := range values {
		if v == "" || seen[v] {
			return false
		}
		seen[v] = true
	}
	return true
}
func sameSet(a, b []string) bool {
	left, right := slices.Clone(a), slices.Clone(b)
	sort.Strings(left)
	sort.Strings(right)
	return slices.Equal(left, right)
}
func capabilityChanges(a, b map[string]string) int {
	keys := map[string]bool{}
	for k := range a {
		keys[k] = true
	}
	for k := range b {
		keys[k] = true
	}
	n := 0
	for k := range keys {
		if a[k] != b[k] {
			n++
		}
	}
	return n
}

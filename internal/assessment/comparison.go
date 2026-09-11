package assessment

import (
	"fmt"
	"maps"
	"reflect"
	"sort"
	"strings"
)

type ComparedRun struct {
	Run     Run              `json:"run"`
	Summary Summary          `json:"summary"`
	Usage   map[string]Usage `json:"usage"`
}
type Comparison struct {
	Plan     Plan          `json:"plan"`
	Eligible bool          `json:"eligible"`
	Reasons  []string      `json:"reasons"`
	Runs     []ComparedRun `json:"runs"`
}

func Compare(p Plan, runs []Run) (Comparison, error) {
	out := Comparison{Plan: p, Eligible: true}
	if err := ValidatePlan(p); err != nil {
		return out, err
	}
	reasonSet := map[string]bool{}
	reason := func(s string) {
		if reasonSet[s] {
			return
		}
		reasonSet[s] = true
		out.Reasons = append(out.Reasons, s)
		out.Eligible = false
	}
	if !purposes[p.Purpose].eligible {
		reason("descriptive evidence only; not default-change evidence")
	}
	if len(p.Conditions) < 2 {
		reason("comparison has fewer than two conditions")
	}
	conditions := map[string]Condition{}
	tasks := map[string]PlanTask{}
	different := false
	for _, c := range p.Conditions {
		if !reflect.DeepEqual(fixedCondition(c, ""), fixedCondition(p.Conditions[0], "")) {
			different = true
		}
		conditions[c.ID] = c
		if len(c.Acceptance) == 0 || len(c.ReviewAxes) == 0 || !sameSet(c.Acceptance, p.Conditions[0].Acceptance) || !sameSet(c.ReviewAxes, p.Conditions[0].ReviewAxes) {
			reason("unequal or missing acceptance and review obligations")
		}
	}
	if !different {
		reason("conditions have no experimental difference")
	}
	for _, task := range p.Tasks {
		tasks[task.ID] = task
		if !task.Holdout || strings.Contains(strings.ToLower(task.ID), "ft311") {
			reason("task is not eligible held-out evidence: " + task.ID)
		}
		if task.Repetitions < 2 {
			reason("planned repetitions are only a pilot: " + task.ID)
		}
	}
	if purposes[p.Purpose].causal {
		arms, _ := resolveCausalArms(p.Conditions)
		if len(arms.none.Capabilities) != 0 || len(arms.current.Capabilities) == 0 || capabilityChanges(arms.current.Capabilities, arms.changed.Capabilities) != 1 {
			reason("causal arm must change exactly one capability")
		}
	}
	qualityNames := make([]string, 0, len(p.QualityTolerance.Measures))
	for name := range p.QualityTolerance.Measures {
		qualityNames = append(qualityNames, name)
	}
	sort.Strings(qualityNames)
	seen := map[string]bool{}
	slots := map[trialSlot]bool{}
	filled := map[trialPair]int{}
	counts := map[string]int{}
	failures := map[string]int{}
	for _, r := range runs {
		if seen[r.RunID] {
			return out, fmt.Errorf("duplicate selected run")
		}
		seen[r.RunID] = true
		c, ok := conditions[r.Condition]
		if !ok {
			return out, fmt.Errorf("run has undeclared condition")
		}
		task, ok := tasks[r.TaskID]
		if !ok {
			return out, fmt.Errorf("run has undeclared task")
		}
		summary, err := Summarize(r)
		if err != nil {
			return out, err
		}
		usage, err := RunUsage(r)
		if err != nil {
			return out, err
		}
		if r.Source != c.Revision {
			return out, fmt.Errorf("run revision differs from pinned condition")
		}
		for _, a := range r.Attempts {
			if a.State == "running" || a.State == "incomplete" {
				reason("attempt has incomplete evidence: " + r.RunID + "/" + a.AttemptID)
			}
			line, ok := c.Lines[a.Role]
			if !ok || a.Model != line.Model || a.Effort != line.Effort {
				return out, fmt.Errorf("run model or effort differs from pinned condition")
			}
		}
		if !r.Holdout {
			reason("run is not held out: " + r.RunID)
		}
		if summary.Incomplete || r.State == "running" || r.State == "incomplete" {
			reason("run has incomplete evidence: " + r.RunID)
		}
		if r.Trial == nil {
			reason("run has no trial provenance: " + r.RunID)
		} else {
			tr := r.Trial
			if tr.PlanID != p.ID {
				return out, fmt.Errorf("run belongs to another plan")
			}
			if tr.Harness == "" || tr.Limits == nil {
				reason("run has unknown execution conditions: " + r.RunID)
			} else if tr.Harness != c.Harness || !maps.Equal(tr.Limits, c.Limits) {
				return out, fmt.Errorf("run harness or execution limits differ from plan")
			}
			if !maps.Equal(tr.Capabilities, c.Capabilities) {
				return out, fmt.Errorf("run capabilities differ from plan")
			}
			if !sameSet(tr.Acceptance, c.Acceptance) || !sameSet(tr.ReviewAxes, c.ReviewAxes) {
				reason("run assurance differs from plan: " + r.RunID)
			}
			if tr.Repetition == nil {
				reason("run repetition is unknown: " + r.RunID)
			} else {
				if *tr.Repetition < 1 || *tr.Repetition > task.Repetitions {
					return out, fmt.Errorf("run repetition outside plan")
				}
				key := trialSlot{r.TaskID, r.Condition, *tr.Repetition}
				if slots[key] {
					return out, fmt.Errorf("duplicate trial repetition")
				}
				slots[key] = true
				filled[trialPair{r.TaskID, r.Condition}]++
			}
		}
		counts[r.Condition]++
		if r.State == "failed" || r.State == "cancelled" {
			failures[r.Condition]++
		}
		for _, name := range qualityNames {
			bounds := p.QualityTolerance.Measures[name]
			measure, ok := r.Quality[name]
			if !ok || measure.Value == nil {
				reason("quality is unknown: " + r.RunID + "/" + name)
				break
			}
			if (bounds.Min != nil && *measure.Value < *bounds.Min) || (bounds.Max != nil && *measure.Value > *bounds.Max) {
				reason("quality exceeds tolerance: " + r.RunID + "/" + name)
				break
			}
		}
		out.Runs = append(out.Runs, ComparedRun{Run: r, Summary: summary, Usage: usage})
	}
	completed := map[string]map[string]bool{}
	for pair, n := range filled {
		if n == tasks[pair.task].Repetitions {
			if completed[pair.condition] == nil {
				completed[pair.condition] = map[string]bool{}
			}
			completed[pair.condition][pair.task] = true
		}
	}
	for _, c := range p.Conditions {
		if missing := len(p.Tasks) - len(completed[c.ID]); missing > 0 {
			for _, task := range p.Tasks {
				if !completed[c.ID][task.ID] {
					reason(fmt.Sprintf("missing planned repetitions for %d tasks in %s; first: %s", missing, c.ID, task.ID))
					break
				}
			}
		}
		if counts[c.ID] > 0 && float64(failures[c.ID])/float64(counts[c.ID]) > *p.QualityTolerance.MaxFailureRate {
			reason("failure rate exceeds tolerance: " + c.ID)
		}
	}
	return out, nil
}

type trialPair struct{ task, condition string }
type trialSlot struct {
	task, condition string
	repetition      int
}

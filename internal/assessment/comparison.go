package assessment

import (
	"fmt"
	"maps"
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
	reason := func(s string) {
		for _, old := range out.Reasons {
			if old == s {
				return
			}
		}
		out.Reasons = append(out.Reasons, s)
		out.Eligible = false
	}
	if p.Purpose == "pilot" || p.Purpose == "descriptive" {
		reason("descriptive evidence only; not default-change evidence")
	}
	if len(p.Conditions) < 2 {
		reason("comparison has fewer than two conditions")
	}
	conditions := map[string]Condition{}
	tasks := map[string]PlanTask{}
	for _, c := range p.Conditions {
		conditions[c.ID] = c
		if len(c.Acceptance) == 0 || len(c.ReviewAxes) == 0 || !sameSet(c.Acceptance, p.Conditions[0].Acceptance) || !sameSet(c.ReviewAxes, p.Conditions[0].ReviewAxes) {
			reason("unequal or missing acceptance and review obligations")
		}
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
	if p.Purpose == "kit-causal" {
		if len(conditions["no-bench"].Capabilities) != 0 || len(conditions["current-bench"].Capabilities) == 0 || capabilityChanges(conditions["current-bench"].Capabilities, conditions["changed-capability"].Capabilities) != 1 {
			reason("causal arm must change exactly one capability")
		}
	}
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
			line, ok := c.Lines[a.Role]
			if !ok || a.Model != line.Model || a.Effort != line.Effort {
				return out, fmt.Errorf("run model or effort differs from pinned condition")
			}
		}
		if !r.Holdout {
			reason("run is not held out: " + r.RunID)
		}
		if summary.Incomplete || r.State == "running" {
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
		for name, bounds := range p.QualityTolerance.Measures {
			measure, ok := r.Quality[name]
			if !ok || measure.Value == nil {
				reason("quality is unknown: " + r.RunID + "/" + name)
				continue
			}
			if (bounds.Min != nil && *measure.Value < *bounds.Min) || (bounds.Max != nil && *measure.Value > *bounds.Max) {
				reason("quality exceeds tolerance: " + r.RunID + "/" + name)
			}
		}
		out.Runs = append(out.Runs, ComparedRun{Run: r, Summary: summary, Usage: usage})
	}
	for _, c := range p.Conditions {
		for _, task := range p.Tasks {
			if missing := task.Repetitions - filled[trialPair{task.ID, c.ID}]; missing > 0 {
				reason(fmt.Sprintf("missing %d planned repetitions: %s / %s", missing, task.ID, c.ID))
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

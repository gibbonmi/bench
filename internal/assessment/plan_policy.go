package assessment

type purposePolicy struct{ eligible, causal bool }

var purposes = map[string]purposePolicy{
	"pilot": {}, "descriptive": {}, "default-change": {eligible: true}, "kit-causal": {eligible: true, causal: true},
}

const capabilityVariable = "capability"

var variableProjections = map[string]func(*Condition){
	capabilityVariable: func(c *Condition) { c.Capabilities = nil },
	"revision":         func(c *Condition) { c.Revision = "" },
	"harness":          func(c *Condition) { c.Harness = "" },
	"limits":           func(c *Condition) { c.Limits = nil },
	"model":            func(c *Condition) { projectLines(c, func(l Line) Line { l.Model = ""; return l }) },
	"effort":           func(c *Condition) { projectLines(c, func(l Line) Line { l.Effort = ""; return l }) },
}

func projectLines(c *Condition, project func(Line) Line) {
	lines := map[string]Line{}
	for role, line := range c.Lines {
		lines[role] = project(line)
	}
	c.Lines = lines
}
func fixedCondition(c Condition, variable string) Condition {
	c.ID = ""
	c.Acceptance = nil
	c.ReviewAxes = nil
	if len(c.Capabilities) == 0 {
		c.Capabilities = nil
	}
	if project := variableProjections[variable]; project != nil {
		project(&c)
	}
	return c
}

type causalArms struct{ none, current, changed Condition }

func resolveCausalArms(conditions []Condition) (causalArms, bool) {
	byID := map[string]Condition{}
	for _, c := range conditions {
		byID[c.ID] = c
	}
	none, a := byID["no-bench"]
	current, b := byID["current-bench"]
	changed, c := byID["changed-capability"]
	return causalArms{none, current, changed}, len(conditions) == 3 && a && b && c
}

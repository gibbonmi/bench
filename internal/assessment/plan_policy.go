package assessment

type purposePolicy struct{ eligible, causal bool }

func purpose(name string) (purposePolicy, bool) {
	switch name {
	case "pilot", "descriptive":
		return purposePolicy{}, true
	case "default-change":
		return purposePolicy{eligible: true}, true
	case "kit-causal":
		return purposePolicy{eligible: true, causal: true}, true
	default:
		return purposePolicy{}, false
	}
}

const capabilityVariable = "capability"

func variableProjection(name string) (func(*Condition), bool) {
	switch name {
	case capabilityVariable:
		return func(c *Condition) { c.Capabilities = nil }, true
	case "revision":
		return func(c *Condition) { c.Revision = "" }, true
	case "harness":
		return func(c *Condition) { c.Harness = "" }, true
	case "limits":
		return func(c *Condition) { c.Limits = nil }, true
	case "model":
		return func(c *Condition) { projectLines(c, func(l Line) Line { l.Model = ""; return l }) }, true
	case "effort":
		return func(c *Condition) { projectLines(c, func(l Line) Line { l.Effort = ""; return l }) }, true
	default:
		return nil, false
	}
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
	if project, ok := variableProjection(variable); ok {
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

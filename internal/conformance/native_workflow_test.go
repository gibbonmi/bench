package conformance

import "strings"

type workflowTriggerShape struct {
	pullRequest, pushBranches, mainBranch bool
}

func nativeWorkflowTriggers(text string) workflowTriggerShape {
	var shape workflowTriggerShape
	inOn, inPush, inBranches := false, false, false
	for _, raw := range strings.Split(text, "\n") {
		line := strings.TrimRight(raw, " \t\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " "))
		if indent == 0 {
			inOn = line == "on:"
			inPush, inBranches = false, false
			continue
		}
		if !inOn {
			continue
		}
		if indent == 2 {
			inPush, inBranches = trimmed == "push:", false
			if trimmed == "pull_request:" {
				shape.pullRequest = true
			}
			continue
		}
		if inPush && indent == 4 && trimmed == "branches:" {
			shape.pushBranches, inBranches = true, true
			continue
		}
		if inBranches && indent == 6 && trimmed == "- main" {
			shape.mainBranch = true
		}
	}
	return shape
}

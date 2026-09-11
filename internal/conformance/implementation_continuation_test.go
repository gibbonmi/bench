package conformance

import (
	"fmt"
	"slices"
	"testing"

	"github.com/gibbonmi/bench/internal/anchors"
)

func TestImplementationContinuation(t *testing.T) {
	wantPredicates := map[string]struct {
		file    string
		section string
		needle  string
	}{
		"implementation continuation: craft-line dropped the explicit uncapped policy": {
			file: ".agents/skills/bench-craft-line/SKILL.md", section: "The declaration",
			needle: "The iteration policy is a numeric cap or an explicit `uncapped` policy.",
		},
		"implementation continuation: craft-line dropped verified acceptance progress": {
			file: ".agents/skills/bench-craft-line/SKILL.md", section: "Retained implementation continuation",
			needle: "A verified acceptance improvement is progress.",
		},
		"implementation continuation: craft-line counted expected TDD reds as attempts": {
			file: ".agents/skills/bench-craft-line/SKILL.md", section: "Retained implementation continuation",
			needle: "An expected TDD red is not an attempt.",
		},
		"implementation continuation: craft-line dropped the two-attempt reassessment threshold": {
			file: ".agents/skills/bench-craft-line/SKILL.md", section: "Retained implementation continuation",
			needle: "After two completed attempts with no progress, reassess before the next implementation attempt.",
		},
		"implementation continuation: craft-line dropped the explicit user-budget stop": {
			file: ".agents/skills/bench-craft-line/SKILL.md", section: "Retained implementation continuation",
			needle: "Stop when the run reaches an explicit user budget.",
		},
		"implementation continuation: craft-line dropped the actionable blocked handoff": {
			file: ".agents/skills/bench-craft-line/SKILL.md", section: "Retained implementation continuation",
			needle: "If no useful next check remains, report the unresolved blocker and the smallest decision or evidence needed.",
		},
		"implementation continuation: craft-line dropped useful-evidence progress": {
			file: ".agents/skills/bench-craft-line/SKILL.md", section: "Retained implementation continuation",
			needle: "New useful evidence is progress when it changes the next action.",
		},
		"implementation continuation: craft-line counted individual tool calls as attempts": {
			file: ".agents/skills/bench-craft-line/SKILL.md", section: "Retained implementation continuation",
			needle: "An individual tool call is not an attempt.",
		},
		"implementation continuation: craft-line dropped the discriminating reassessment check": {
			file: ".agents/skills/bench-craft-line/SKILL.md", section: "Retained implementation continuation",
			needle: "State the changed hypothesis and the next discriminating check.",
		},
		"implementation continuation: craft-line dropped numeric-cap exhaustion": {
			file: ".agents/skills/bench-craft-line/SKILL.md", section: "Retained implementation continuation",
			needle: "Stop when the run exhausts a selected numeric cap.",
		},
		"implementation continuation: craft-line dropped the required-decision stop": {
			file: ".agents/skills/bench-craft-line/SKILL.md", section: "Retained implementation continuation",
			needle: "Stop dependent implementation for a required user decision.",
		},
		"implementation continuation: craft-line dropped the terminal external-blocker stop": {
			file: ".agents/skills/bench-craft-line/SKILL.md", section: "Retained implementation continuation",
			needle: "Stop for an external blocker when no independent work remains.",
		},
		"implementation continuation: craft-line dropped the cancellation stop": {
			file: ".agents/skills/bench-craft-line/SKILL.md", section: "Retained implementation continuation",
			needle: "Stop immediately when the user cancels.",
		},
		"implementation continuation: craft-line counted diagnostic-only work as an attempt": {
			file: ".agents/skills/bench-craft-line/SKILL.md", section: "Retained implementation continuation",
			needle: "A diagnostic-only action is not a completed implementation-and-verification attempt.",
		},
		"implementation continuation: operating guide dropped the continuation-policy pointer": {
			file: ".bench/BENCH.md", section: "Workflow",
			needle: "`craft-line` owns the retained implementation continuation policy.",
		},
		"implementation continuation: implementation phase dropped the continuation-policy action": {
			file: ".agents/commands/bench-implement-spec.md", section: "Build",
			needle: "Apply `craft-line`'s retained implementation continuation policy throughout the ticket graph.",
		},
	}
	family := anchorsWithDiagnosticPrefix("implementation continuation: ")
	for diagnostic, predicate := range wantPredicates {
		index := slices.IndexFunc(family, func(anchor anchors.Anchor) bool { return anchor.Diagnostic == diagnostic })
		if index < 0 {
			t.Errorf("implementation-continuation anchor is absent: %s", diagnostic)
			continue
		}
		anchor := family[index]
		if anchor.File != predicate.file || anchor.Section != predicate.section || anchor.Needle != predicate.needle {
			t.Errorf("implementation-continuation predicate drifted for %q: %s", diagnostic, fmt.Sprintf("%#v", anchor))
		}
	}
	runAnchorBites(t, family, func(anchor anchors.Anchor) string { return anchor.Diagnostic })
}

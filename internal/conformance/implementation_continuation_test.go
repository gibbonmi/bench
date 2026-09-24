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
		"implementation continuation: craft-line dropped continuation while progress holds": {
			file: ".agents/skills/bench-craft-line/SKILL.md", section: "Retained implementation continuation",
			needle: "Continue while progress holds inside the approved scope.",
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
		"implementation continuation: craft-line restored an artificial stop for uncapped work": {
			file: ".agents/skills/bench-craft-line/SKILL.md", section: "Retained implementation continuation",
			needle: "An uncapped implementation has no artificial iteration stop within the approved spec.",
		},
		"implementation continuation: craft-line dropped the debug route": {
			file: ".agents/skills/bench-craft-line/SKILL.md", section: "Retained implementation continuation",
			needle: "After reassessment, the ticket author can invoke `$bench-debug`.",
		},
		"implementation continuation: craft-delegate restored its copy of the diagnostic helper boundary that the operating guide owns": {
			file: ".agents/skills/bench-craft-delegate/SKILL.md", section: "",
			needle: "A diagnostic helper can inspect evidence, but it receives no implementation or repair assignment.",
		},
		"implementation continuation: diagnostic escalation changed the author session": {
			file: ".agents/skills/bench-craft-line/SKILL.md", section: "Retained implementation continuation",
			needle: "A diagnostic route does not change the ticket author's session.",
		},
		"implementation continuation: craft-line dropped the author's effort adjustment": {
			file: ".agents/skills/bench-craft-line/SKILL.md", section: "Ticketed-build stage defaults",
			needle: "The ticket author can adjust effort in its own session and reports the change.",
		},
		"implementation continuation: craft-delegate dropped top-tier diagnostic pre-approval": {
			file: ".agents/skills/bench-craft-delegate/SKILL.md", section: "Delegate or retain",
			needle: "During implementation, brief higher-tier diagnostic consultation is pre-approved through the top tier.",
		},
		"implementation continuation: delegation discipline expanded diagnostic request context": {
			file: ".agents/skills/bench-craft-delegate/references/delegation-discipline.md", section: "Read-only returns",
			needle: "A diagnostic request contains only one question, the relevant error, minimal code, and attempted hypotheses.",
		},
		"implementation continuation: delegation discipline dropped the short diagnostic return": {
			file: ".agents/skills/bench-craft-delegate/references/delegation-discipline.md", section: "Read-only returns",
			needle: "A diagnostic return contains a short diagnosis and the next check.",
		},
		"implementation continuation: diagnostic escalation changed the author model": {
			file: ".agents/skills/bench-craft-line/SKILL.md", section: "Retained implementation continuation",
			needle: "A diagnostic route does not change the ticket author's model.",
		},
		"implementation continuation: delegation discipline dropped the actual consultation line": {
			file: ".agents/skills/bench-craft-delegate/references/delegation-discipline.md", section: "Read-only returns",
			needle: "Record the actual consultation line used.",
		},
		"implementation continuation: delegation discipline allowed an undeclared model substitution": {
			file: ".agents/skills/bench-craft-delegate/references/delegation-discipline.md", section: "Read-only returns",
			needle: "If the selected model is unavailable, report the failure and use an available authorized diagnostic route without an undeclared model substitution.",
		},

		"implementation continuation: bounded repair dropped fixed allowance": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Scope and allowance",
			needle: "Each implementation chunk permits at most two repair cycles after its initial review.",
		},
		"implementation continuation: bounded repair restored the retired mode list": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "",
			needle: "This allowance applies to retained, full, delegated, unattended, and light-path implementation runs.",
		},
		"implementation continuation: bounded repair dropped light-path chunk": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Scope and allowance",
			needle: "Light-path work counts as one chunk when it receives review findings.",
		},
		"implementation continuation: bounded repair dropped initial-work exclusion": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Scope and allowance",
			needle: "Initial implementation, pre-review checks, and the first review consume no repair cycles.",
		},
		"implementation continuation: bounded repair dropped cycle unit": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Scope and allowance",
			needle: "A repair cycle is one repair attempt and verification of its affected findings. It can address several findings.",
		},
		"implementation continuation: bounded repair dropped verification-only exclusion": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Scope and allowance",
			needle: "Individual tool calls and unchanged verification reruns consume no additional repair cycles.",
		},
		"implementation continuation: bounded repair dropped progress precedence": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Scope and allowance",
			needle: "The allowance takes precedence over continuation while progress holds after initial review.",
		},
		"implementation continuation: bounded repair dropped progress and review continuity": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Scope and allowance",
			needle: "Progress does not extend the allowance. A fresh review does not reset the same chunk's count.",
		},
		"implementation continuation: bounded repair dropped uncapped-line boundary": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Scope and allowance",
			needle: "An uncapped implementation line does not remove this allowance.",
		},
		"implementation continuation: bounded repair dropped reviewer extension": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Scope and allowance",
			needle: "Only an explicit reviewer decision extends the allowance; record the additional work it permits.",
		},
		"implementation continuation: bounded repair dropped blocker evidence": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Classification and completion",
			needle: "Before accepting a blocker, cite its binding requirement or concrete defect evidence.",
		},
		"implementation continuation: bounded repair dropped exhaustion classification": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Classification and completion",
			needle: "The following unresolved conditions remain blocking, including at exhaustion:",
		},
		"implementation continuation: bounded repair dropped required checks": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Classification and completion",
			needle: "A required check fails.",
		},
		"implementation continuation: bounded repair dropped acceptance failures": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Classification and completion",
			needle: "An approved acceptance requirement fails.",
		},
		"implementation continuation: bounded repair dropped correctness defects": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Classification and completion",
			needle: "A concrete correctness defect remains, even with a green gate.",
		},
		"implementation continuation: bounded repair dropped safety defects": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Classification and completion",
			needle: "A concrete safety defect remains, even with a green gate.",
		},
		"implementation continuation: bounded repair dropped mandatory standards": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Classification and completion",
			needle: "A documented mandatory standard fails, even without an automated check.",
		},
		"implementation continuation: bounded repair dropped optional classification": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Classification and completion",
			needle: "A preference without a binding requirement or concrete defect remains optional advice.",
		},
		"implementation continuation: bounded repair dropped optional census": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Classification and completion",
			needle: "Give optional advice no disposition or finding ID, and exclude it from finding totals.",
		},
		"implementation continuation: bounded repair dropped optional retention": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Classification and completion",
			needle: "Retain optional advice in a separate advice section of the native excerpt and the review pickup.",
		},
		"implementation continuation: bounded repair dropped optional pass": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Classification and completion",
			needle: "A native review can pass with no finding IDs while its prose retains optional advice.",
		},
		"implementation continuation: bounded repair dropped current review evidence": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Classification and completion",
			needle: "After substantive repairs, obtain current results or permitted native reaffirmations from every review axis.",
		},
		"implementation continuation: bounded repair dropped evidence-only definition": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Classification and completion",
			needle: "A review-record prose correction is evidence-only only when it changes no finding, source identity, observation, disposition, or verification claim.",
		},
		"implementation continuation: bounded repair counted evidence-only prose": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Classification and completion",
			needle: "Evidence-only corrections consume no repair cycle. Batch all cited corrections before verification.",
		},
		"implementation continuation: bounded repair widened prose reaffirmation": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Classification and completion",
			needle: "Only the issuing axis reaffirms an evidence-only correction unless it invalidates another axis's evidence.",
		},
		"implementation continuation: bounded repair started gate before evidence freeze": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Classification and completion",
			needle: "The coordinator commits the final review record before starting the whole-project gate. It does not edit the gate subject while the gate runs.",
		},
		"implementation continuation: bounded repair dropped completion safeguards": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Classification and completion",
			needle: "Unresolved findings, stale source identity, missing axes, and incomplete verification still block completion.",
		},
		"implementation continuation: bounded repair dropped pickup state": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "State and handoff",
			needle: "Retain the consumed count and any explicit reviewer extension in the existing review pickup.",
		},
		"implementation continuation: bounded repair dropped light-path state": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "State and handoff",
			needle: "For light-path work without a pickup, retain that state in the existing session handoff.",
		},
		"implementation continuation: bounded repair dropped resume state": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "State and handoff",
			needle: "Before a resumed session starts another repair cycle, read the retained state.",
		},
		"implementation continuation: bounded repair dropped missing-count reconstruction": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "State and handoff",
			needle: "If the count is missing, reconstruct it from available evidence; never assume zero.",
		},
		"implementation continuation: bounded repair dropped unknown-count stop": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "State and handoff",
			needle: "If evidence cannot establish the remaining allowance, return that uncertainty to the reviewer before another repair cycle.",
		},
		"implementation continuation: bounded repair dropped exhaustion report": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "State and handoff",
			needle: "At exhaustion with blockers, report the consumed allowance, completed repairs, remaining blockers, and the needed reviewer decision.",
		},
		"implementation continuation: bounded repair dropped dependent-work stop": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "State and handoff",
			needle: "While that decision is pending, stop further repair cycles and dependent chunks.",
		},
		"implementation continuation: bounded repair dropped early completion": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "State and handoff",
			needle: "The allowance is a maximum, not a target.",
		},
		"implementation continuation: craft-line dropped the bounded-repair owner": {
			file: ".agents/skills/bench-craft-line/SKILL.md", section: "Retained implementation continuation",
			needle: "At initial review or before post-review repairs, including resumed work, read [the bounded repair policy](references/bounded-repair-policy.md).",
		},
		"implementation continuation: implementation phase dropped bounded repairs": {
			file: ".agents/commands/bench-implement-spec.md", section: "Build",
			needle: "At initial review or before post-review repairs, including resumed work, read [the bounded repair policy](../skills/bench-craft-line/references/bounded-repair-policy.md).",
		},
		"implementation continuation: craft-review dropped bounded repairs": {
			file: ".agents/skills/bench-craft-review/SKILL.md", section: "",
			needle: "Before classifying repairs, read [the bounded repair policy](../bench-craft-line/references/bounded-repair-policy.md).",
		},
		"implementation continuation: review phase dropped bounded repairs": {
			file: ".agents/commands/bench-review-implementation.md", section: "Review modes",
			needle: "Before classifying repairs, read [the bounded repair policy](../skills/bench-craft-line/references/bounded-repair-policy.md).",
		},
		"implementation continuation: review progression counted optional advice": {
			file: ".agents/commands/bench-review-implementation.md", section: "Review modes",
			needle: "Here, findings that prevent progression are unresolved blockers; retain optional advice separately under the policy.",
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

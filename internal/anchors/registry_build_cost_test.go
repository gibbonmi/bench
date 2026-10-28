package anchors

import "testing"

// TestBuildCostAnchors keeps the build-cost guidance rows red-capable. Each needle and
// diagnostic is written here independently of the registry, so a rewrite that drops a row
// cannot define itself green.
func TestBuildCostAnchors(t *testing.T) {
	const implement = ".agents/commands/bench-implement-spec.md"
	const plan = ".agents/skills/bench-craft-spec/references/completion-plan.md"
	anchorHarness{group: AfterImplementSpec, rules: []anchorRule{
		{file: implement, needle: "then runs each verification requirement once at that commit", want: "build cost: the implement phase restored the run of each verification requirement", forbidden: true},
		{file: implement, needle: "check-floor", want: "build cost: the implement phase restored the retired check-floor row", forbidden: true},
		{file: plan, needle: "check-floor", want: "build cost: the completion-plan reference restored the retired check-floor row", forbidden: true},
	}}.check(t)
}

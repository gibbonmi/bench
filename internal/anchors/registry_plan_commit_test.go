package anchors

import "testing"

// These independent expectations make removal of either half of the plan-commit route fail.
func TestPlanCommitAnchors(t *testing.T) {
	const implement = ".agents/commands/bench-implement-spec.md"
	anchorHarness{group: AfterImplementSpec, rules: []anchorRule{
		{file: implement, needle: "a dispatch follows uses plain `bench commit`", want: "plan commit: a plan commit that a dispatch follows uses plain bench commit"},
		{file: implement, needle: "a plan commit that no dispatch follows uses `bench commit --preflight-build <slug>`", want: "plan commit: a plan commit that no dispatch follows uses bench commit --preflight-build <slug>"},
	}}.check(t)
}

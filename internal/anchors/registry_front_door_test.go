package anchors

import "testing"

// The independently authored rule proves that removing the production anchor turns the check red.
func TestFrontDoorFailureContinuationAnchorRedOnRemoval(t *testing.T) {
	anchorHarness{
		group: BeforeStructured,
		rules: []anchorRule{
			{
				file:    ".agents/commands/bench.md",
				section: "Exit handoff",
				needle:  "If the routed command fails, continue with `.agents/commands/bench-debug.md` as the active phase.",
				want:    ".agents/commands/bench.md Exit handoff dropped failure continuation through the bench-debug phase",
			},
		},
	}.check(t)
}

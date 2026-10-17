package anchors

import "testing"

// TestCraftResearchAnchorsRedOnRemoval keeps the shaping phase's retired research rules
// out: restoring one raises its own diagnostic. Each needle and diagnostic is written here
// independently of the registry, so a dropped prohibition cannot define itself green.
func TestCraftResearchAnchorsRedOnRemoval(t *testing.T) {
	const shaping = ".agents/commands/bench-shape-idea.md"
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{
				file:      shaping,
				needle:    "and produce a short",
				want:      ".agents/commands/bench-shape-idea.md keeps the Research ticket's own output rule; craft-research owns the durable output",
				forbidden: true,
			},
			{
				file:      shaping,
				needle:    "Include a runnable compatibility probe when the answer",
				want:      ".agents/commands/bench-shape-idea.md keeps the Research ticket's own probe rule; craft-research owns the compatibility evidence rule",
				forbidden: true,
			},
			{
				file:      shaping,
				needle:    "as a read-only delegation. Otherwise resolve it inline.",
				want:      ".agents/commands/bench-shape-idea.md keeps its own research delegation rule; craft-research owns the rounds and the delegates",
				forbidden: true,
			},
		},
	}.check(t)
}

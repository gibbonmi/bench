package anchors

import "testing"

// TestDecisionMapAnchorsRedOnRemoval keeps the retired decision-map forms out of the two
// phase commands. The shaping file must not write the unnamed prose-check term, and neither
// file may keep the flat asset path, so a command that survived the topic-folder move
// unedited turns red. Each needle and diagnostic is written here independently of the
// registry.
func TestDecisionMapAnchorsRedOnRemoval(t *testing.T) {
	const (
		shapeIdea = ".agents/commands/bench-shape-idea.md"
		writeSpec = ".agents/commands/bench-write-spec.md"
	)
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{
				file:      shapeIdea,
				needle:    "prose preflight",
				want:      ".agents/commands/bench-shape-idea.md writes an unnamed prose check; `bench gate-prose` is the one handle",
				forbidden: true,
			},
			{
				file:      shapeIdea,
				needle:    "decisions/assets/",
				want:      ".agents/commands/bench-shape-idea.md keeps the flat decisions/assets/ path; a map-owned asset stays under the topic folder",
				forbidden: true,
			},
			{
				file:      writeSpec,
				needle:    "decisions/assets/",
				want:      ".agents/commands/bench-write-spec.md keeps the flat decisions/assets/ path; a map-owned asset moves with the topic folder",
				forbidden: true,
			},
		},
	}.check(t)
}

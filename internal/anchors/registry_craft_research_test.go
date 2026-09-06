package anchors

import "testing"

// TestCraftResearchAnchorsRedOnRemoval keeps the research skill's anchored sentences and
// the caller pointers red-capable. Each needle and diagnostic is written here
// independently of the registry, so a rewrite that drops a rule cannot define itself green.
func TestCraftResearchAnchorsRedOnRemoval(t *testing.T) {
	const skill = ".agents/skills/bench-craft-research/SKILL.md"
	const shaping = ".agents/commands/bench-shape-idea.md"
	const bootstrap = ".agents/skills/bench-craft-spec/references/bootstrap-authority.md"
	anchorHarness{
		group: AfterImplementSpec,
		rules: []anchorRule{
			{
				file:   skill,
				needle: "The question graph is complete when every load-bearing claim traces to one question node.",
				want:   ".agents/skills/bench-craft-research/SKILL.md dropped the question-graph completeness rule",
			},
			{
				file:   skill,
				needle: "Fan out only when at least two frontier questions are mutually independent.",
				want:   ".agents/skills/bench-craft-research/SKILL.md dropped the mutual-independence threshold for fan-out",
			},
			{
				file:   skill,
				needle: "Synthesize a completed round before opening another.",
				want:   ".agents/skills/bench-craft-research/SKILL.md dropped the synthesize-before-the-next-round rule",
			},
			{
				file:   skill,
				needle: "The artifact under study and first-party upstream documentation or APIs are primary sources.",
				want:   ".agents/skills/bench-craft-research/SKILL.md dropped the primary-source definition",
			},
			{
				file:   skill,
				needle: "The coordinator re-opens every source supporting a load-bearing conclusion and independently checks every join between returns.",
				want:   ".agents/skills/bench-craft-research/SKILL.md dropped the coordinator's source re-open and join check",
			},
			{
				file:   skill,
				needle: "One coordinator-authored durable output per research run, keyed by topic.",
				want:   ".agents/skills/bench-craft-research/SKILL.md dropped the one durable output per run keyed by topic",
			},
			{
				file:   skill,
				needle: "Every load-bearing claim cites an exact local path and line or a primary URL with its retrieval date.",
				want:   ".agents/skills/bench-craft-research/SKILL.md dropped the exact path-and-line or dated primary URL citation rule",
			},
			{
				file:   skill,
				needle: "A byte or wire compatibility claim stays unverified until a separate runnable probe returns.",
				want:   ".agents/skills/bench-craft-research/SKILL.md dropped the runnable-probe rule for a compatibility claim",
			},
			{
				file:   skill,
				needle: "Research never owns a write delegate, a done-claim, a reviewer decision, or a prototype.",
				want:   ".agents/skills/bench-craft-research/SKILL.md dropped the read-side boundary",
			},
			{
				file:   shaping,
				needle: "The Research ticket type charges `craft-research`, and a required compatibility probe becomes a Prototype ticket that the Research ticket names in `Blocked by`.",
				want:   ".agents/commands/bench-shape-idea.md dropped the Research ticket type's pointer at craft-research and the Prototype ticket a compatibility probe needs",
			},
			{
				file:   bootstrap,
				needle: "`craft-research` states the probe rule: a compatibility claim stays unverified\nuntil a separate runnable probe returns.",
				want:   ".agents/skills/bench-craft-spec/references/bootstrap-authority.md dropped the pointer at the craft-research probe rule",
			},
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

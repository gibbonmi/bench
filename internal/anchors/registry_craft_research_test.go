package anchors

import "testing"

// TestCraftResearchAnchorsRedOnRemoval holds the research rules that guidance must keep:
// the question graph completes on one node per load-bearing claim, fan-out waits for two
// mutually independent frontier questions, a round is synthesized before the next opens,
// only the artifact and first-party upstream documentation warrant a finding, the
// coordinator re-opens the sources and checks the joins, one durable output per run
// carries every claim's citation, a compatibility claim waits for a runnable probe, and
// research stays on the read side. Each needle and diagnostic is written here
// independently of the registry, so a rewrite that drops a rule cannot define itself green.
func TestCraftResearchAnchorsRedOnRemoval(t *testing.T) {
	const skill = ".agents/skills/bench-craft-research/SKILL.md"
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
		},
	}.check(t)
}

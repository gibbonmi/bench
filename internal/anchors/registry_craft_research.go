package anchors

// This file holds the research anchor rows — the `craft-research` skill's load-bearing
// sentences and the caller pointers that name the skill — which live here for the line
// budget of registry_data.go. Research reads primary sources, verifies the joins between
// delegate returns, and leaves one cited durable output per run.
var craftResearchAnchors = []Anchor{
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-research/SKILL.md", Kind: Require, Section: "", Needle: "The question graph is complete when every load-bearing claim traces to one question node.", Diagnostic: ".agents/skills/bench-craft-research/SKILL.md dropped the question-graph completeness rule"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-research/SKILL.md", Kind: Require, Section: "", Needle: "Fan out only when at least two frontier questions are mutually independent.", Diagnostic: ".agents/skills/bench-craft-research/SKILL.md dropped the mutual-independence threshold for fan-out"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-research/SKILL.md", Kind: Require, Section: "", Needle: "Synthesize a completed round before opening another.", Diagnostic: ".agents/skills/bench-craft-research/SKILL.md dropped the synthesize-before-the-next-round rule"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-research/SKILL.md", Kind: Require, Section: "", Needle: "The artifact under study and first-party upstream documentation or APIs are primary sources.", Diagnostic: ".agents/skills/bench-craft-research/SKILL.md dropped the primary-source definition"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-research/SKILL.md", Kind: Require, Section: "", Needle: "The coordinator re-opens every source supporting a load-bearing conclusion and independently checks every join between returns.", Diagnostic: ".agents/skills/bench-craft-research/SKILL.md dropped the coordinator's source re-open and join check"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-research/SKILL.md", Kind: Require, Section: "", Needle: "One coordinator-authored durable output per research run, keyed by topic.", Diagnostic: ".agents/skills/bench-craft-research/SKILL.md dropped the one durable output per run keyed by topic"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-research/SKILL.md", Kind: Require, Section: "", Needle: "Every load-bearing claim cites an exact local path and line or a primary URL with its retrieval date.", Diagnostic: ".agents/skills/bench-craft-research/SKILL.md dropped the exact path-and-line or dated primary URL citation rule"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-research/SKILL.md", Kind: Require, Section: "", Needle: "A byte or wire compatibility claim stays unverified until a separate runnable probe returns.", Diagnostic: ".agents/skills/bench-craft-research/SKILL.md dropped the runnable-probe rule for a compatibility claim"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-research/SKILL.md", Kind: Require, Section: "", Needle: "Research never owns a write delegate, a done-claim, a reviewer decision, or a prototype.", Diagnostic: ".agents/skills/bench-craft-research/SKILL.md dropped the read-side boundary"},
}

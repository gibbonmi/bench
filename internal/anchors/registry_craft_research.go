package anchors

// craftResearchAnchors refuse the shaping phase's retired research rules: its own output,
// probe, and delegation rules. `craft-research` owns the durable output, the compatibility
// evidence, and the rounds.
var craftResearchAnchors = []Anchor{
	{Group: AfterImplementSpec, File: ".agents/commands/bench-shape-idea.md", Kind: Forbid, Section: "", Needle: "and produce a short", Diagnostic: ".agents/commands/bench-shape-idea.md keeps the Research ticket's own output rule; craft-research owns the durable output"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-shape-idea.md", Kind: Forbid, Section: "", Needle: "Include a runnable compatibility probe when the answer", Diagnostic: ".agents/commands/bench-shape-idea.md keeps the Research ticket's own probe rule; craft-research owns the compatibility evidence rule"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-shape-idea.md", Kind: Forbid, Section: "", Needle: "as a read-only delegation. Otherwise resolve it inline.", Diagnostic: ".agents/commands/bench-shape-idea.md keeps its own research delegation rule; craft-research owns the rounds and the delegates"},
}

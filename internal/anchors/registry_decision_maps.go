package anchors

// decisionMapAnchors refuse the retired decision-map forms in the two phase commands: the
// unnamed prose check and the flat decisions/assets/ path. A map-owned asset lives under
// its topic folder.
var decisionMapAnchors = []Anchor{
	{Group: AfterImplementSpec, File: ".agents/commands/bench-shape-idea.md", Kind: Forbid, Section: "", Needle: "prose preflight", Diagnostic: ".agents/commands/bench-shape-idea.md writes an unnamed prose check; `bench gate-prose` is the one handle"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-shape-idea.md", Kind: Forbid, Section: "", Needle: "decisions/assets/", Diagnostic: ".agents/commands/bench-shape-idea.md keeps the flat decisions/assets/ path; a map-owned asset stays under the topic folder"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-write-spec.md", Kind: Forbid, Section: "", Needle: "decisions/assets/", Diagnostic: ".agents/commands/bench-write-spec.md keeps the flat decisions/assets/ path; a map-owned asset moves with the topic folder"},
}

package anchors

var frontDoorAnchors = []Anchor{
	{File: ".agents/commands/bench.md", Kind: RequireInSection, Section: "Exit handoff", Needle: "If the routed command fails, continue with `.agents/commands/bench-debug.md` as the active phase.", Diagnostic: ".agents/commands/bench.md Exit handoff dropped failure continuation through the bench-debug phase"},
}

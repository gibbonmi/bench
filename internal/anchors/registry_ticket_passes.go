package anchors

var ticketSlicingAnchors = []Anchor{
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-tickets/SKILL.md", Kind: RequireInSection, Section: "Draft the breakdown", Needle: "Before ticket lock, run a source-clause pass and an executable-route pass.", Diagnostic: "ticket slicing: verification passes must precede ticket lock"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-tickets/SKILL.md", Kind: RequireInSection, Section: "Draft the breakdown", Needle: "Assign each applicable source clause to a ticket or an explicit scope cut.", Diagnostic: "ticket slicing: source clauses need a ticket or explicit scope cut"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-tickets/SKILL.md", Kind: RequireInSection, Section: "Draft the breakdown", Needle: "Trace each claimed operation through its executable owner, using the [map-discipline proof rules](../bench-craft-spec/references/map-discipline.md#before-the-map-locks).", Diagnostic: "ticket slicing: executable routes need their owner and canonical proof rules"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-tickets/SKILL.md", Kind: RequireInSection, Section: "Write one file per ticket", Needle: "A ticket names every relocation destination in `Writes:`, including destinations for snapshots and registry rows.", Diagnostic: "ticket slicing: relocation destinations belong in Writes"},
}

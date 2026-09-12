package anchors

import "testing"

// These independent expectations make removal of a ticket-lock obligation fail.
func TestTicketSlicingPasses(t *testing.T) {
	anchorHarness{group: AfterImplementSpec, rules: []anchorRule{
		{file: ".agents/skills/bench-craft-tickets/SKILL.md", section: "Draft the breakdown", needle: "Before ticket lock, run a source-clause pass and an executable-route pass.", want: "ticket slicing: verification passes must precede ticket lock"},
		{file: ".agents/skills/bench-craft-tickets/SKILL.md", section: "Draft the breakdown", needle: "Assign each applicable source clause to a ticket or an explicit scope cut.", want: "ticket slicing: source clauses need a ticket or explicit scope cut"},
		{file: ".agents/skills/bench-craft-tickets/SKILL.md", section: "Draft the breakdown", needle: "Trace each claimed operation through its executable owner, using the [map-discipline proof rules](../bench-craft-spec/references/map-discipline.md#before-the-map-locks).", want: "ticket slicing: executable routes need their owner and canonical proof rules"},
		{file: ".agents/skills/bench-craft-tickets/SKILL.md", section: "Write one file per ticket", needle: "A ticket names every relocation destination in `Writes:`, including destinations for snapshots and registry rows.", want: "ticket slicing: relocation destinations belong in Writes"},
	}}.check(t)
}

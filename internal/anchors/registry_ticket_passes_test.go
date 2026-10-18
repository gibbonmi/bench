package anchors

import "testing"

// TestTicketSlicingAnchors keeps the ticket-slicing command lines, the map seam cell, and the
// retired-form refusals red-capable. Each needle and diagnostic is written here independently
// of the registry, so a rewrite that drops a row cannot define itself green.
func TestTicketSlicingAnchors(t *testing.T) {
	const writeSpec = ".agents/commands/bench-write-spec.md"
	const issues = "docs/known-issues.md"
	const slicing = ".agents/skills/bench-craft-tickets/references/slicing-checks.md"
	anchorHarness{group: AfterImplementSpec, rules: []anchorRule{
		{file: slicing, section: "Slicing rules", needle: "The slicer runs `bench preflight build <slug>` after each change to a ticket or an acceptance row.", want: "ticket slicing: build preflight reruns after each ticket or acceptance row change"},
		{file: ".agents/skills/bench-craft-spec/references/map-discipline.md", section: "Per row", needle: "The seam cell cites that function as `` `<path>_test.go` (`<Name>`) ``, so `bench coverage --check` and build preflight resolve the name against the tree.", want: "map discipline: an existing test seam cites its function in the resolvable form"},
		{file: writeSpec, needle: "A delegate on `craft-line`'s chunk ticket review line grades the chunk's tickets.", want: "write-spec slicing: the retired review of each chunk returned", forbidden: true},
		{file: writeSpec, section: "Process", needle: "After the ticket review accepts, write the fence section. `bench preflight build <slug> --plan-only` then grades the whole ticket graph.", want: "write-spec slicing: the plan-only preflight grades the whole ticket graph after the ticket review"},
		{file: ".agents/skills/bench-craft-line/SKILL.md", needle: "| Chunk ticket review (`/bench-write-spec`) | cheap + high |", want: "write-spec slicing: craft-line kept the retired chunk ticket review row", forbidden: true},
		{file: issues, needle: "Cap the witnesses in a charge", want: "witness cap: the known-issues row kept a copy of the cap", forbidden: true},
		{file: issues, needle: "| FT405: the record verb reads an excerpt from the last response of the worktree. |", want: "record ticket: the known-issues file kept a row that FT405 retires", forbidden: true},
		{file: slicing, needle: "## Enforced rules", want: "closure removal: the slicing checks kept the enforced-rules copy of the preflight checks", forbidden: true},
		{file: slicing, needle: "--propose-writes", want: "closure removal: the slicing checks kept a writes-proposal step", forbidden: true},
		{file: slicing, needle: "--breakdown", want: "closure removal: the slicing checks kept a breakdown step", forbidden: true},
		{file: writeSpec, needle: "--propose-writes", want: "closure removal: write-spec slicing kept the writes-proposal comparison", forbidden: true},
		{file: ".agents/skills/bench-craft-spec/references/map-only-route.md", needle: "--propose-writes", want: "closure removal: the map-only route kept a writes-proposal step", forbidden: true},
		{file: ".agents/skills/bench-craft-spec/references/map-discipline.md", needle: "breakdown, proposal", want: "closure removal: map discipline kept the breakdown and proposal steps", forbidden: true},
		{file: issues, needle: "writes proposal", want: "closure removal: the known-issues file kept a row that names the writes proposal", forbidden: true},
	}}.check(t)
}

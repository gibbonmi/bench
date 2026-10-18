package anchors

import "testing"

// TestTicketSlicingAnchors keeps the ticket-slicing command lines, the map seam cell, and the
// retired-form refusals red-capable. Each needle and diagnostic is written here independently
// of the registry, so a rewrite that drops a row cannot define itself green.
func TestTicketSlicingAnchors(t *testing.T) {
	const writeSpec = ".agents/commands/bench-write-spec.md"
	const issues = "docs/known-issues.md"
	anchorHarness{group: AfterImplementSpec, rules: []anchorRule{
		{file: ".agents/skills/bench-craft-tickets/references/slicing-checks.md", section: "Slicing rules", needle: "The slicer runs `bench preflight build <slug>` after each change to a ticket or an acceptance row.", want: "ticket slicing: build preflight reruns after each ticket or acceptance row change"},
		{file: ".agents/skills/bench-craft-spec/references/map-discipline.md", section: "Per row", needle: "The seam cell cites that function as `` `<path>_test.go` (`<Name>`) ``, so `bench coverage --check` and build preflight resolve the name against the tree.", want: "map discipline: an existing test seam cites its function in the resolvable form"},
		{file: writeSpec, section: "Process", needle: "Compare each ticket's `Writes:` line with `bench preflight build <slug> --propose-writes` under the slicing checks.", want: "write-spec slicing: each chunk ticket's Writes line meets the propose-writes proposal"},
		{file: writeSpec, needle: "A delegate on `craft-line`'s chunk ticket review line grades the chunk's tickets.", want: "write-spec slicing: the retired review of each chunk returned", forbidden: true},
		{file: writeSpec, section: "Process", needle: "After the ticket review accepts, write the fence section. `bench preflight build <slug> --plan-only` then grades the whole ticket graph.", want: "write-spec slicing: the plan-only preflight grades the whole ticket graph after the ticket review"},
		{file: ".agents/skills/bench-craft-line/SKILL.md", needle: "| Chunk ticket review (`/bench-write-spec`) | cheap + high |", want: "write-spec slicing: craft-line kept the retired chunk ticket review row", forbidden: true},
		{file: issues, needle: "Cap the witnesses in a charge", want: "witness cap: the known-issues row kept a copy of the cap", forbidden: true},
		{file: issues, needle: "| FT405: the record verb reads an excerpt from the last response of the worktree. |", want: "record ticket: the known-issues file kept a row that FT405 retires", forbidden: true},
	}}.check(t)
}

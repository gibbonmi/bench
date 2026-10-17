package anchors

// ticketSlicingAnchors pin the preflight and coverage command lines of ticket slicing and the
// map seam cell, and they forbid the retired chunk ticket review and the known-issues copies.
var ticketSlicingAnchors = []Anchor{
	{Group: AfterImplementSpec, File: slicingChecks, Kind: RequireInSection, Class: Command, Section: "Slicing rules", Needle: "The slicer runs `bench preflight build <slug>` after each change to a ticket or an acceptance row.", Diagnostic: "ticket slicing: build preflight reruns after each ticket or acceptance row change"},
	{Group: AfterImplementSpec, File: mapDiscipline, Kind: RequireInSection, Class: Command, Section: "Per row", Needle: "The seam cell cites that function as `` `<path>_test.go` (`<Name>`) ``, so `bench coverage --check` and build preflight resolve the name against the tree.", Diagnostic: "map discipline: an existing test seam cites its function in the resolvable form"},
	{Group: AfterImplementSpec, File: writeSpecCommand, Kind: RequireInSection, Class: Command, Section: "Process", Needle: "Compare each ticket's `Writes:` line with `bench preflight build <slug> --propose-writes` under the slicing checks.", Diagnostic: "write-spec slicing: each chunk ticket's Writes line meets the propose-writes proposal"},
	{Group: AfterImplementSpec, File: writeSpecCommand, Kind: Forbid, Section: "", Needle: "chunk ticket review line", Diagnostic: "write-spec slicing: the retired review of each chunk returned"},
	{Group: AfterImplementSpec, File: writeSpecCommand, Kind: RequireInSection, Class: Command, Section: "Process", Needle: "After the ticket review accepts, write the fence section. `bench preflight build <slug> --plan-only` then grades the whole ticket graph.", Diagnostic: "write-spec slicing: the plan-only preflight grades the whole ticket graph after the ticket review"},
	{Group: AfterImplementSpec, File: craftLineSkill, Kind: Forbid, Section: "", Needle: "| Chunk ticket review (`/bench-write-spec`)", Diagnostic: "write-spec slicing: craft-line kept the retired chunk ticket review row"},
	{Group: AfterImplementSpec, File: knownIssues, Kind: Forbid, Section: "", Needle: "Cap the witnesses in a charge", Diagnostic: "witness cap: the known-issues row kept a copy of the cap"},
	{Group: AfterImplementSpec, File: knownIssues, Kind: Forbid, Section: "", Needle: "FT405", Diagnostic: "record ticket: the known-issues file kept a row that FT405 retires"},
}

// slicingChecks is the craft-tickets reference that states the slicing rules.
const slicingChecks = ".agents/skills/bench-craft-tickets/references/slicing-checks.md"

// writeSpecCommand is the phase command that reviews the spec and slices its tickets by chunk.
const writeSpecCommand = ".agents/commands/bench-write-spec.md"

// craftLineSkill is the skill whose stage defaults route the write-spec review rounds.
const craftLineSkill = ".agents/skills/bench-craft-line/SKILL.md"

// mapDiscipline is the craft-spec reference that states how a seam cell cites its test.
const mapDiscipline = ".agents/skills/bench-craft-spec/references/map-discipline.md"

// knownIssues is the known-issues file that never carries a second copy of the witness cap
// or a row that FT405 retires.
const knownIssues = "docs/known-issues.md"

package anchors

// ticketSlicingAnchors pin the preflight and coverage command lines of ticket slicing and the
// map seam cell, and they forbid the retired chunk ticket review, the known-issues copies,
// and the retired closure-era forms: each step of the writes proposal and the breakdown,
// the closure and fence-writes rows, the spec fence section, and the delegate sentence
// that called the fence a refusal boundary, in any case or emphasis. They also forbid the
// retired build-entry take rule that took only a completion-plan red.
var ticketSlicingAnchors = []Anchor{
	{Group: AfterImplementSpec, File: slicingChecks, Kind: RequireInSection, Class: Command, Section: "Slicing rules", Needle: "The slicer runs `bench preflight build <slug>` after each change to a ticket or an acceptance row.", Diagnostic: "ticket slicing: build preflight reruns after each ticket or acceptance row change"},
	{Group: AfterImplementSpec, File: mapDiscipline, Kind: RequireInSection, Class: Command, Section: "Per row", Needle: "The seam cell cites that function as `` `<path>_test.go` (`<Name>`) ``, so `bench coverage --check` and build preflight resolve the name against the tree.", Diagnostic: "map discipline: an existing test seam cites its function in the resolvable form"},
	{Group: AfterImplementSpec, File: writeSpecCommand, Kind: Forbid, Section: "", Needle: "chunk ticket review line", Diagnostic: "write-spec slicing: the retired review of each chunk returned"},
	{Group: AfterImplementSpec, File: writeSpecCommand, Kind: RequireInSection, Class: Command, Section: "Process", Needle: "After the ticket review accepts, `bench preflight build <slug> --plan-only` grades the whole ticket graph.", Diagnostic: "write-spec slicing: the plan-only preflight grades the whole ticket graph after the ticket review"},
	{Group: AfterImplementSpec, File: craftLineSkill, Kind: Forbid, Section: "", Needle: "| Chunk ticket review (`/bench-write-spec`)", Diagnostic: "write-spec slicing: craft-line kept the retired chunk ticket review row"},
	{Group: AfterImplementSpec, File: knownIssues, Kind: Forbid, Section: "", Needle: "Cap the witnesses in a charge", Diagnostic: "witness cap: the known-issues row kept a copy of the cap"},
	{Group: AfterImplementSpec, File: knownIssues, Kind: Forbid, Section: "", Needle: "FT405", Diagnostic: "record ticket: the known-issues file kept a row that FT405 retires"},
	{Group: AfterImplementSpec, File: knownIssues, Kind: Forbid, Section: "", Needle: "FT438", Diagnostic: "assignment exemption: the known-issues file kept a row that FT438 retires"},
	{Group: AfterImplementSpec, File: knownIssues, Kind: Forbid, Section: "", Needle: "FT406", Diagnostic: "commit rebuild: the known-issues file kept a row that FT406 retires"},
	{Group: AfterImplementSpec, File: knownIssues, Kind: Forbid, Section: "", Needle: "FT482", Diagnostic: "cache release seal: the known-issues file kept a row that FT482 retires"},
	{Group: AfterImplementSpec, File: slicingChecks, Kind: Forbid, Section: "", Needle: "## Enforced rules", Diagnostic: "closure removal: the slicing checks kept the enforced-rules copy of the preflight checks"},
	{Group: AfterImplementSpec, File: slicingChecks, Kind: Forbid, Section: "", Needle: "--propose-writes", Diagnostic: "closure removal: the slicing checks kept a writes-proposal step"},
	{Group: AfterImplementSpec, File: slicingChecks, Kind: Forbid, Section: "", Needle: "--breakdown", Diagnostic: "closure removal: the slicing checks kept a breakdown step"},
	{Group: AfterImplementSpec, File: writeSpecCommand, Kind: Forbid, Section: "", Needle: "--propose-writes", Diagnostic: "closure removal: write-spec slicing kept the writes-proposal comparison"},
	{Group: AfterImplementSpec, File: mapOnlyRoute, Kind: Forbid, Section: "", Needle: "--propose-writes", Diagnostic: "closure removal: the map-only route kept a writes-proposal step"},
	{Group: AfterImplementSpec, File: mapDiscipline, Kind: Forbid, Section: "", Needle: "breakdown, proposal", Diagnostic: "closure removal: map discipline kept the breakdown and proposal steps"},
	{Group: AfterImplementSpec, File: knownIssues, Kind: Forbid, Section: "", Needle: "writes proposal", Diagnostic: "closure removal: the known-issues file kept a row that names the writes proposal"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Section: "", Needle: "-closure", Diagnostic: "closure removal: the implement-spec staleness route kept a closure row"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Section: "", Needle: "fence-writes", Diagnostic: "closure removal: the implement-spec staleness route kept the fence-writes row"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Section: "", Needle: "If every red row is `completion-plan`, the staleness pass below takes the red.", Diagnostic: "build-entry staleness: the implement-spec take rule dropped the writes-landed and map-sources rows"},
	{Group: AfterImplementSpec, File: craftSpecSkill, Kind: Forbid, Section: "", Needle: "## Ownership fences", Diagnostic: "closure removal: the craft-spec template kept the spec fence section"},
	{Group: AfterImplementSpec, File: delegationDiscipline, Kind: ForbidCaseFoldedEmphasis, Section: "", Needle: "is a refusal boundary", Diagnostic: "closure removal: the delegate guidance called the fence a refusal boundary"},
}

// slicingChecks is the craft-tickets reference that states the slicing rules.
const slicingChecks = ".agents/skills/bench-craft-tickets/references/slicing-checks.md"

// writeSpecCommand is the phase command that reviews the spec and slices its tickets by chunk.
const writeSpecCommand = ".agents/commands/bench-write-spec.md"

// craftLineSkill is the skill whose stage defaults route the write-spec review rounds.
const craftLineSkill = ".agents/skills/bench-craft-line/SKILL.md"

// craftSpecSkill is the craft-spec skill whose template holds no spec fence section.
const craftSpecSkill = ".agents/skills/bench-craft-spec/SKILL.md"

// mapDiscipline is the craft-spec reference that states how a seam cell cites its test.
const mapDiscipline = ".agents/skills/bench-craft-spec/references/map-discipline.md"

// mapOnlyRoute is the craft-spec reference that states the reviewer-selected thin-spec route.
const mapOnlyRoute = ".agents/skills/bench-craft-spec/references/map-only-route.md"

// knownIssues is the known-issues file that never carries a second copy of the witness cap
// or a row that FT405, FT406, FT438, or FT482 retires.
const knownIssues = "docs/known-issues.md"

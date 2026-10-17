package anchors

// CommitmentRuleMarker opens the canonical delivery commitment rule. The shared-rule check
// requires it in the operating guide and refuses a copy in AGENTS.md or README.md.
const CommitmentRuleMarker = "**Deliver only the committed outcome.**"

// CommitmentDiagnosticPrefix opens every diagnostic of the commitment guidance family.
const CommitmentDiagnosticPrefix = "commitment guidance: "

// commitmentAnchors pin the two spec-stage command lines of the commitment family. Each
// prohibition refuses one retired rule. One grant let work start or displace the commitment
// outside the committed outcome. Another grant let a build ship a light-path fix. A
// restriction held light-path work behind the active commitment.
var commitmentAnchors = []Anchor{
	{Group: AfterImplementSpec, File: writeSpecCommand, Kind: Require, Class: Command, Needle: "A staged spec is planning work: `bench preflight build <slug> --plan-only` validates it, and its delivery waits for `bench commitment start`.", Diagnostic: CommitmentDiagnosticPrefix + "spec authoring dropped the planning-only staged spec"},
	{Group: AfterImplementSpec, File: writeSpecCommand, Kind: RequireInSection, Class: Command, Section: "Exit handoff", Needle: "From the primary checkout, run `bench handoff --next` with the start command that `bench status` prints.", Diagnostic: CommitmentDiagnosticPrefix + "spec authoring dropped the start command as the handoff next command"},

	{Group: AfterImplementSpec, File: drainCommand, Kind: Forbid, Needle: "build the item in this session (\"implement now\") by default.", Diagnostic: CommitmentDiagnosticPrefix + "drain restored the default implementation of a light-path item"},
	{Group: AfterImplementSpec, File: drainCommand, Kind: Forbid, Needle: "Open a `ROADMAP.md` row only when the reviewer declines.", Diagnostic: CommitmentDiagnosticPrefix + "drain restored the roadmap row only for a declined item"},
	{Group: AfterImplementSpec, File: drainCommand, Kind: Forbid, Needle: "is implement-now work only after", Diagnostic: CommitmentDiagnosticPrefix + "drain restored the commitment admission of a light-path item"},
	{Group: AfterImplementSpec, File: drainCommand, Kind: Forbid, Needle: "Rewrite the `## Recommended sequence` section", Diagnostic: CommitmentDiagnosticPrefix + "drain restored the sequence rewrite"},
	{Group: AfterImplementSpec, File: drainCommand, Kind: Forbid, Needle: "Rank rows by severity.", Diagnostic: CommitmentDiagnosticPrefix + "drain restored the severity rank"},
	{Group: AfterImplementSpec, File: drainCommand, Kind: Forbid, Needle: "choose actionable work over blocked work", Diagnostic: CommitmentDiagnosticPrefix + "drain restored the actionable-work preference"},
	{Group: AfterImplementSpec, File: drainCommand, Kind: Forbid, Needle: "Only when rows are equally actionable, apply literal dependencies, then explicit reviewer pricing.", Diagnostic: CommitmentDiagnosticPrefix + "drain restored dependencies and pricing as an actionability tiebreaker"},
	{Group: AfterImplementSpec, File: drainCommand, Kind: Forbid, Needle: "rank by descending occurrence count", Diagnostic: CommitmentDiagnosticPrefix + "drain restored the occurrence-count tiebreaker"},
	{Group: AfterImplementSpec, File: drainCommand, Kind: Forbid, Needle: "apply the existing reproduced defect-over-feature rule, then cheapest-first cost rule", Diagnostic: CommitmentDiagnosticPrefix + "drain restored the defect and cost rules as an occurrence tiebreaker"},
	{Group: AfterImplementSpec, File: drainCommand, Kind: Forbid, Needle: "the top line of the refreshed `## Recommended sequence`", Diagnostic: CommitmentDiagnosticPrefix + "drain restored the refreshed sequence as its next command"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-final-check.md", Kind: Forbid, Needle: "Leave the roadmap and capture rows to `/bench-drain`", Diagnostic: CommitmentDiagnosticPrefix + "final check restored roadmap closure as a later drain"},
	{Group: AfterImplementSpec, File: operatingGuide, Kind: Forbid, Needle: "a light-path fix that needs no reviewer decision", Diagnostic: CommitmentDiagnosticPrefix + "operating guide restored the learning fix without admission"},
	{Group: AfterImplementSpec, File: operatingGuide, Kind: Forbid, Needle: "light path and fixes included", Diagnostic: CommitmentDiagnosticPrefix + "operating guide restored the commitment start for light-path work"},
	{Group: AfterImplementSpec, File: operatingGuide, Kind: Forbid, Needle: "a light-path fix that the active committed outcome needs", Diagnostic: CommitmentDiagnosticPrefix + "operating guide restored the learning fix for the active outcome only"},
	{Group: AfterImplementSpec, File: operatingGuide, Kind: Forbid, Needle: "`/bench-implement-spec` included, a `bench learning` entry's light-path fix ships at once", Diagnostic: CommitmentDiagnosticPrefix + "operating guide restored the in-build learning fix"},
	{Group: AfterImplementSpec, File: operatingGuide, Kind: Forbid, Needle: "A small defect you find mid-work is not roadmap work", Diagnostic: CommitmentDiagnosticPrefix + "operating guide restored the fix of every mid-work defect without admission"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench.md", Kind: Forbid, Needle: "take the first `sequence` row", Diagnostic: CommitmentDiagnosticPrefix + "router restored the first sequence row as its work"},
}

// The guidance files that the commitment family names more than once.
const (
	operatingGuide = ".bench/BENCH.md"
	drainCommand   = ".agents/commands/bench-drain.md"
)

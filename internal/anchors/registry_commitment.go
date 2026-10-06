package anchors

// CommitmentRuleMarker opens the canonical delivery commitment rule. The shared-rule check
// requires it in the operating guide and refuses a copy in AGENTS.md or README.md.
const CommitmentRuleMarker = "**Deliver only the committed outcome.**"

// CommitmentDiagnosticPrefix opens every diagnostic of the commitment guidance family.
const CommitmentDiagnosticPrefix = "commitment guidance: "

// commitmentAnchors pin the canonical delivery commitment rule, each phase route to its
// commands, and one prohibition for each retired grant that let a drain, a severity rank,
// or a later reconcile admit or displace work.
var commitmentAnchors = []Anchor{
	{Group: AfterImplementSpec, File: operatingGuide, Kind: RequireInSection, Section: "Workflow", Needle: CommitmentRuleMarker, Diagnostic: CommitmentDiagnosticPrefix + "operating guide dropped the canonical delivery commitment rule"},
	{Group: AfterImplementSpec, File: operatingGuide, Kind: RequireInSection, Section: "Workflow", Needle: "A spec implementation starts only through `bench commitment start` for the eligible outcome that `bench status` names.", Diagnostic: CommitmentDiagnosticPrefix + "operating guide dropped the spec start route"},
	{Group: AfterImplementSpec, File: operatingGuide, Kind: RequireInSection, Section: "Workflow", Needle: "Any other finding, idea, learning, or drained item that needs a spec stays uncommitted intake; minimal support for the active outcome stays in it.", Diagnostic: CommitmentDiagnosticPrefix + "operating guide dropped uncommitted intake"},
	{Group: AfterImplementSpec, File: operatingGuide, Kind: RequireInSection, Section: "Workflow", Needle: "Only my explicit direction changes the commitment, through `bench commitment plan` and then `bench commitment approve`; no drain, label, score, or count displaces it.", Diagnostic: CommitmentDiagnosticPrefix + "operating guide dropped explicit displacement"},
	{Group: AfterImplementSpec, File: operatingGuide, Kind: RequireInSection, Section: "Workflow", Needle: "When you propose work, put confirmed defects first, then refactors, then features, by purpose rather than label; dependencies and the approved order govern execution.", Diagnostic: CommitmentDiagnosticPrefix + "operating guide dropped purpose-based defect and refactor priority"},
	{Group: AfterImplementSpec, File: operatingGuide, Kind: RequireInSection, Section: "Workflow", Needle: "When the active outcome cannot continue, run `bench commitment block` with the reason and tell me; the obligation stays.", Diagnostic: CommitmentDiagnosticPrefix + "operating guide dropped the blocker report"},
	{Group: AfterImplementSpec, File: operatingGuide, Kind: RequireInSection, Section: "Workflow", Needle: "It carries one tickets-only folder with exactly one ticket, and it lands with that folder as `--spec`.", Diagnostic: CommitmentDiagnosticPrefix + "operating guide dropped the light-path observable"},
	{Group: AfterImplementSpec, File: operatingGuide, Kind: RequireInSection, Section: "Workflow", Needle: "Its production paths stay inside that ticket's `Writes:` line.", Diagnostic: CommitmentDiagnosticPrefix + "operating guide dropped the light-path Writes boundary"},
	{Group: AfterImplementSpec, File: operatingGuide, Kind: RequireInSection, Section: "Workflow", Needle: "It may remove a roadmap row that no committed outcome pins as a source; a pinned row changes only through its outcome.", Diagnostic: CommitmentDiagnosticPrefix + "operating guide dropped the light-path row rule"},
	{Group: AfterImplementSpec, File: operatingGuide, Kind: RequireInSection, Section: "Workflow", Needle: "`/bench-drain` dispatches every other light-path fix that its verdicts keep, whatever the active commitment is.", Diagnostic: CommitmentDiagnosticPrefix + "operating guide dropped the drain light-path dispatch"},
	{Group: AfterImplementSpec, File: drainCommand, Kind: Require, Needle: "The commitment rule in `.bench/BENCH.md` decides whether spec intake starts; a drain approval never admits it.", Diagnostic: CommitmentDiagnosticPrefix + "drain dropped its route to the commitment rule"},
	{Group: AfterImplementSpec, File: drainCommand, Kind: Require, Needle: "A drained light-path item needs no commitment: the drain dispatches it under the light-path fix rule in `.bench/BENCH.md`.", Diagnostic: CommitmentDiagnosticPrefix + "drain dropped the delegate route for a light-path item"},
	{Group: AfterImplementSpec, File: drainCommand, Kind: Require, Needle: "The drain does not write the `## Recommended sequence` section.", Diagnostic: CommitmentDiagnosticPrefix + "drain dropped the commitment ownership of the recommended sequence"},
	{Group: AfterImplementSpec, File: drainCommand, Kind: Require, Needle: "It informs a commitment proposal and never reorders committed work.", Diagnostic: CommitmentDiagnosticPrefix + "drain dropped the classification limit"},
	{Group: AfterImplementSpec, File: finalCheckCommand, Kind: RequireInSection, Section: "Exit handoff", Needle: "Verified closure is part of delivery.", Diagnostic: CommitmentDiagnosticPrefix + "final check dropped verified closure from delivery"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Require, Needle: "It declares the line, starts its committed outcome through `bench commitment start`, and works vertical slices at the pre-agreed seams.", Diagnostic: CommitmentDiagnosticPrefix + "implementation dropped its commitment start"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Require, Needle: "A light-path change needs no commitment start; `.bench/BENCH.md` owns that rule.", Diagnostic: CommitmentDiagnosticPrefix + "implementation dropped the light-path start exemption"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-write-spec.md", Kind: Require, Needle: "A staged spec is planning work: `bench preflight build <slug> --plan-only` validates it, and its delivery waits for `bench commitment start`.", Diagnostic: CommitmentDiagnosticPrefix + "spec authoring dropped the planning-only staged spec"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-debug.md", Kind: Require, Needle: "Diagnosis needs no commitment start, but each commit of the repro or the fix follows the commitment rule in `.bench/BENCH.md`.", Diagnostic: CommitmentDiagnosticPrefix + "debug dropped its route to the commitment rule"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench.md", Kind: Require, Needle: "take its `commitment_outlook` row", Diagnostic: CommitmentDiagnosticPrefix + "router dropped the commitment outlook"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-setup-repo.md", Kind: Require, Needle: "Delivery work waits for the initial commitment.", Diagnostic: CommitmentDiagnosticPrefix + "setup dropped the initial commitment"},
	{Group: AfterImplementSpec, File: "DATA_HANDLING.md", Kind: Require, Needle: "The records stay on the local machine.", Diagnostic: CommitmentDiagnosticPrefix + "data inventory dropped the local retention of commitment records"},
	{Group: AfterImplementSpec, File: "DATA_HANDLING.md", Kind: Require, Needle: "They hold no transcript, prompt, objective text, environment value, or credential.", Diagnostic: CommitmentDiagnosticPrefix + "data inventory dropped the record contents boundary"},

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
	{Group: AfterImplementSpec, File: finalCheckCommand, Kind: Forbid, Needle: "Leave the roadmap and capture rows to `/bench-drain`", Diagnostic: CommitmentDiagnosticPrefix + "final check restored roadmap closure as a later drain"},
	{Group: AfterImplementSpec, File: operatingGuide, Kind: Forbid, Needle: "a light-path fix that needs no reviewer decision", Diagnostic: CommitmentDiagnosticPrefix + "operating guide restored the learning fix without admission"},
	{Group: AfterImplementSpec, File: operatingGuide, Kind: Forbid, Needle: "light path and fixes included", Diagnostic: CommitmentDiagnosticPrefix + "operating guide restored the commitment start for light-path work"},
	{Group: AfterImplementSpec, File: operatingGuide, Kind: Forbid, Needle: "a light-path fix that the active committed outcome needs", Diagnostic: CommitmentDiagnosticPrefix + "operating guide restored the learning fix for the active outcome only"},
	{Group: AfterImplementSpec, File: operatingGuide, Kind: Forbid, Needle: "A small defect you find mid-work is not roadmap work", Diagnostic: CommitmentDiagnosticPrefix + "operating guide restored the fix of every mid-work defect without admission"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench.md", Kind: Forbid, Needle: "take the first `sequence` row", Diagnostic: CommitmentDiagnosticPrefix + "router restored the first sequence row as its work"},
}

// The guidance files that the commitment family names more than once.
const (
	operatingGuide    = ".bench/BENCH.md"
	drainCommand      = ".agents/commands/bench-drain.md"
	finalCheckCommand = ".agents/commands/bench-final-check.md"
)

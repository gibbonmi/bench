package anchors

// chargeBindingAnchors pin the charge, probe, claim, retry, and landing rules of the
// delegation discipline, and the transport consumer protocol of the charge evidence.
var chargeBindingAnchors = []Anchor{
	{Group: AfterImplementSpec, File: delegationDiscipline, Kind: RequireInSection, Section: "In the charge", Needle: "Before a write charge dispatches, the coordinator ticks each item of this list against the ticket's `Writes:` line.", Diagnostic: "charge binding: the coordinator ticks the charge list against Writes"},
	{Group: AfterImplementSpec, File: delegationDiscipline, Kind: RequireInSection, Section: "In the charge", Needle: "A test that guards on a root privilege routes through the capability seam, never through a bare `t.Skip`.", Diagnostic: "charge binding: a root-privilege guard routes through the capability seam"},
	{Group: AfterImplementSpec, File: delegationDiscipline, Kind: Forbid, Needle: "A repair fence is the approved chunk union plus the exact paths that the review names.", Diagnostic: "charge binding: the repair fence restored the chunk union plus review-named paths"},
	{Group: AfterImplementSpec, File: delegationDiscipline, Kind: RequireInSection, Section: "In the charge", Needle: "A repair fence is the affected ticket's `Writes:` line.", Diagnostic: "charge binding: the repair fence dropped the affected ticket's Writes line"},
	{Group: AfterImplementSpec, File: delegationDiscipline, Kind: RequireInSection, Section: "In the charge", Needle: "A repair based on a frozen sibling uses an integration assignment from `main`.", Diagnostic: "charge binding: a frozen-sibling repair uses an integration assignment"},
	{Group: AfterImplementSpec, File: delegationDiscipline, Kind: RequireInSection, Section: "Probes", Needle: "A completion-plan probe names the check that independently detects its mutation.", Diagnostic: "charge binding: a completion-plan probe names its detecting check"},
	{Group: AfterImplementSpec, File: delegationDiscipline, Kind: RequireInSection, Section: "Claim schema", Needle: "A stated confidence freezes at return time. The coordinator refuses a later raise.", Diagnostic: "charge binding: a stated confidence freezes at return time"},
	{Group: AfterImplementSpec, File: delegationDiscipline, Kind: RequireInSection, Section: "Retry stops and aggregate readiness", Needle: "Before the final allowed repair attempt, the author checks every review axis for duplicated facts and derived counts.", Diagnostic: "charge binding: the final attempt checks every axis for duplicated facts"},
	{Group: AfterImplementSpec, File: delegationDiscipline, Kind: RequireInSection, Section: "Retry stops and aggregate readiness", Needle: "At the cap, the author preserves a green semantic repair and requests an explicit evidence-scoped extension for the remaining edit.", Diagnostic: "charge binding: the cap preserves the green repair and requests an extension"},
	{Group: AfterImplementSpec, File: delegationDiscipline, Kind: Forbid, Needle: "After a rebase changes the source digest", Diagnostic: "charge binding: the digest rule restored the rebase premise that the workflow rejects"},
	{Group: AfterImplementSpec, File: delegationDiscipline, Kind: RequireInSection, Section: "Before the landing", Needle: "After a `bench worktree merge` changes the source digest, the author repeats verification against the new digest before the coordinator accepts the evidence.", Diagnostic: "charge binding: a worktree merge repeats verification against the new digest"},
	{Group: AfterImplementSpec, File: delegationDiscipline, Kind: RequireInSection, Section: "Before the landing", Needle: "An adoption repair covers the contradiction class with an independent fixture, not only the literal examples.", Diagnostic: "charge binding: an adoption repair covers the contradiction class"},
	{Group: AfterImplementSpec, File: delegationDiscipline, Kind: RequireInSection, Section: "Before the landing", Needle: "Implementation delegation and independent review are separate choices, and the assessment record keeps their costs separate.", Diagnostic: "charge binding: delegation and review keep separate costs"},
	{Group: AfterImplementSpec, File: chargeEvidenceFormat, Kind: RequireInSection, Section: "Consumers", Needle: "A cross-harness consumer retrieves every page through its own shell tool and rebuilds each source with its own decoder.", Diagnostic: "charge binding: a consumer rebuilds each source with its own decoder"},
	{Group: AfterImplementSpec, File: chargeEvidenceFormat, Kind: RequireInSection, Section: "Consumers", Needle: "The follow-on hook refuses a pipeline after a Bench call, so the consumer saves each response and decodes the saved bytes in a separate step.", Diagnostic: "charge binding: the consumer protocol names the hook constraint"},
	{Group: AfterImplementSpec, File: chargeEvidenceFormat, Kind: RequireInSection, Section: "Consumers", Needle: "A read takes the shared lock of the evidence store, so a sandboxed consumer needs write access to that store.", Diagnostic: "charge binding: the consumer protocol names the sandbox constraint"},
}

// delegationDiscipline and chargeEvidenceFormat are the craft-delegate references these
// anchors read.
const (
	delegationDiscipline = ".agents/skills/bench-craft-delegate/references/delegation-discipline.md"
	chargeEvidenceFormat = ".agents/skills/bench-craft-delegate/references/charge-evidence-format.md"
)

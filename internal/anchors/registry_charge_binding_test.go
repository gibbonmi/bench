package anchors

import "testing"

// These independent expectations make removal of a charge-binding rule or a consumer rule fail.
func TestChargeBindingAnchors(t *testing.T) {
	const (
		discipline = ".agents/skills/bench-craft-delegate/references/delegation-discipline.md"
		evidence   = ".agents/skills/bench-craft-delegate/references/charge-evidence-format.md"
	)
	anchorHarness{group: AfterImplementSpec, rules: []anchorRule{
		{file: discipline, section: "In the charge", needle: "Before a write charge dispatches, the coordinator ticks each item of this list against the ticket's `Writes:` line.", want: "charge binding: the coordinator ticks the charge list against Writes"},
		{file: discipline, section: "In the charge", needle: "A test that guards on a root privilege routes through the capability seam, never through a bare `t.Skip`.", want: "charge binding: a root-privilege guard routes through the capability seam"},
		{file: discipline, section: "In the charge", needle: "A repair fence is the approved chunk union plus the exact paths that the review names.", want: "charge binding: the repair fence is the chunk union plus review-named paths"},
		{file: discipline, section: "In the charge", needle: "A repair based on a frozen sibling uses an integration assignment from `main`.", want: "charge binding: a frozen-sibling repair uses an integration assignment"},
		{file: discipline, section: "Probes", needle: "A completion-plan probe names the check that independently detects its mutation.", want: "charge binding: a completion-plan probe names its detecting check"},
		{file: discipline, section: "Claim schema", needle: "A stated confidence freezes at return time. The coordinator refuses a later raise.", want: "charge binding: a stated confidence freezes at return time"},
		{file: discipline, section: "Retry stops and aggregate readiness", needle: "Before the final allowed repair attempt, the author checks every review axis for duplicated facts and derived counts.", want: "charge binding: the final attempt checks every axis for duplicated facts"},
		{file: discipline, section: "Retry stops and aggregate readiness", needle: "At the cap, the author preserves a green semantic repair and requests an explicit evidence-scoped extension for the remaining edit.", want: "charge binding: the cap preserves the green repair and requests an extension"},
		{file: discipline, section: "Before the landing", needle: "After a rebase changes the source digest, the author repeats verification against the new digest before the coordinator accepts the evidence.", want: "charge binding: a rebase repeats verification against the new digest"},
		{file: discipline, section: "Before the landing", needle: "An adoption repair covers the contradiction class with an independent fixture, not only the literal examples.", want: "charge binding: an adoption repair covers the contradiction class"},
		{file: discipline, section: "Before the landing", needle: "Implementation delegation and independent review are separate choices, and the assessment record keeps their costs separate.", want: "charge binding: delegation and review keep separate costs"},
		{file: evidence, section: "Consumers", needle: "A cross-harness consumer retrieves every page through its own shell tool and rebuilds each source with its own decoder.", want: "charge binding: a consumer rebuilds each source with its own decoder"},
		{file: evidence, section: "Consumers", needle: "The follow-on hook refuses a pipeline after a Bench call, so the consumer saves each response and decodes the saved bytes in a separate step.", want: "charge binding: the consumer protocol names the hook constraint"},
		{file: evidence, section: "Consumers", needle: "A read takes the shared lock of the evidence store, so a sandboxed consumer needs write access to that store.", want: "charge binding: the consumer protocol names the sandbox constraint"},
	}}.check(t)
}

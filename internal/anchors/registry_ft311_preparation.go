package anchors

// ft311PreparationAnchors pin the build phase's prepared-charge authority. The rows cover
// the charge the phase consumes before dispatch, the approval a generated artifact never
// supplies, the revalidation a stale source forces, and the bounds that hold repair triage
// short of a coordinator decision. Every row is a Require row, because these rules are
// sentences the canonical readers must carry, not routes they must avoid.
var ft311PreparationAnchors = []Anchor{
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Require, Needle: "Run `bench preflight build <slug> --charge --ticket <ticket> --base <base> --source-tip <tip>` and read the prepared evidence identity it returns. Verify the fence, dependency completion, and expected source tip before the ticket's fresh author starts.", Diagnostic: "prepared build guidance: bench-implement-spec.md dropped charge preparation or fence, dependency, and tip verification"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Require, Needle: "Before each ticket, revalidate prepared source identity and required source bytes. Regenerate a charge when either changes.", Diagnostic: "prepared build guidance: bench-implement-spec.md dropped stale-charge revalidation or regeneration"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-delegate/SKILL.md", Kind: Require, Needle: "Use a read-only mid-tier delegate to interpret ambiguous repair evidence. It distinguishes invalid probes, no executed tests, missing coverage, and scope defects only when evidence supports each classification. Missing or contradictory evidence remains unknown. Resolve the harness's bound mid model at medium effort for one iteration. Use high effort only when existing line rules require it.", Diagnostic: "prepared build guidance: craft-delegate dropped bounded mid-tier repair triage"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-delegate/references/delegation-discipline.md", Kind: Require, Needle: "A triage return cites evidence and proposes a bounded repair charge.\n- The coordinator accepts the diagnosis and routes repairs under existing authoring and verification rules.\n- Triage cannot close a finding, expand a fence, waive a probe, or change default tiers.", Diagnostic: "prepared build guidance: delegation discipline lets triage replace coordinator acceptance or independent verification"},

	// The bounded build action rows below pin the prerequisites the build phase requires
	// together before it acts on prepared evidence. Each Require row states one
	// prerequisite, so a phase that drops delivery, context, the current binding, approval,
	// or the supplement bites alone. Approval and the supplement keep their own rows,
	// because verified evidence supplies neither one. The Forbid row keeps the retired
	// build full charge form out of the same guidance.
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Require, Needle: "Build action requires a narrow author read: act only after this session reads the metadata and ticket pages and targeted sources. An author does not retrieve every evidence page. `bench preflight evidence <id> --verify` is the separate artifact-integrity check.", Diagnostic: "bounded build action: bench-implement-spec.md permits build action without a narrow author read"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Require, Needle: "Build action requires available required context: a retrieval receipt, a terminal cursor, or another consumer's delivery never replaces what this session reads itself.", Diagnostic: "bounded build action: bench-implement-spec.md permits build action without available required context"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Require, Needle: "Build action requires a current-action binding: act only after `bench preflight evidence <id> --check-current` binds the artifact to the current assignment and source pair.", Diagnostic: "bounded build action: bench-implement-spec.md permits build action without a current binding"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Require, Needle: "Build action requires reviewer approval, which a generated charge, a staged artifact, or a verified artifact never supplies.", Diagnostic: "bounded build action: bench-implement-spec.md permits build action without reviewer approval"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Require, Needle: "Build action requires the complete task supplement, which this session writes. The supplement names selectors, model, effort, cap, and mutation, and verified evidence never supplies it.", Diagnostic: "bounded build action: bench-implement-spec.md permits build action without the complete task supplement"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Forbid, Needle: "build <slug> --charge --ticket <ticket> --base <base> --source-tip <tip> --full", Diagnostic: "bounded build action: bench-implement-spec.md re-advertises the retired build full charge form"},

	// The consumer-context rows below keep the build phase's reuse and retrieval rules
	// apart. Reuse answers to the new manifest, so a body digest alone cannot carry a
	// source across a regenerated charge. The narrow author read answers to the fresh author
	// itself, so neither a transferred final cursor nor a peer's receipt delivers a byte in
	// its place.
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Require, Needle: "An author that takes a regenerated charge reuses exact available source bytes only after the new manifest verifies their membership, role, and requiredness. A matching body digest alone authorizes no reuse.", Diagnostic: "consumer context: bench-implement-spec.md permits reuse without verified membership, role, and requiredness"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Require, Needle: "A fresh author makes its own narrow author read from the trusted evidence identity. A transferred final cursor and another consumer's receipt deliver no byte to it.", Diagnostic: "consumer context: bench-implement-spec.md permits another consumer's receipt or a final cursor to replace fresh required context"},

	// The fresh ticket author rows below pin the narrow author charge and the small
	// orchestrator. Each Require row states one bound, and each Forbid row keeps a retired
	// full-retrieval sentence out of the build phase. The handoff-refresh needle starts at
	// the read bound, so a negation of the refresh breaks it.
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Require, Needle: "The author charge carries the ticket, its coverage rows, its `Writes:` fence, the evidence identity, and the declared line.", Diagnostic: "fresh ticket author: bench-implement-spec.md dropped the narrow author charge"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Forbid, Needle: "Retrieve every required source with `bench preflight evidence <id>`, and follow each exact successor command until the stream ends.", Diagnostic: "fresh ticket author: bench-implement-spec.md restored the full evidence retrieval"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Forbid, Needle: "Build action requires verified delivery", Diagnostic: "fresh ticket author: bench-implement-spec.md restored the verified-delivery build rule"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Require, Needle: "The orchestrator reads manifests, returns, and verdicts, not code", Diagnostic: "fresh ticket author: bench-implement-spec.md dropped the orchestrator's read bound"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Require, Needle: "not code, and it refreshes `bench handoff` at each chunk checkpoint.", Diagnostic: "fresh ticket author: bench-implement-spec.md dropped the handoff refresh at each chunk checkpoint"},

	// The unchanged-route rows below hold the migration's declared boundary. The
	// implementation phase keeps its own full-run control, which names no preflight
	// command and is therefore no charge route. Write-spec keeps its existing contract,
	// and the Forbid row keeps charge retrieval out of that phase.
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Require, Needle: "## `--full <spec>`\n\nThis command orchestrates the fresh ticket authors, chunk reviews, the final landing, and `/bench-final-check`.", Diagnostic: "unchanged route: bench-implement-spec.md dropped the phase-level full-run control"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-write-spec.md", Kind: Forbid, Needle: "bench preflight evidence", Diagnostic: "unchanged route: bench-write-spec.md adds charge retrieval to the write-spec phase"},
}

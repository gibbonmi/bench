package anchors

// ft311PreparationAnchors pin the build phase's prepared-charge command lines: the charge
// preparation, the artifact-integrity check, and the current-action binding. The Forbid rows
// keep the retired build full charge form, the full evidence retrieval, the verified-delivery
// rule, and charge retrieval in the write-spec phase out of the guidance.
var ft311PreparationAnchors = []Anchor{
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Require, Needle: "Run `bench preflight build <slug> --charge --ticket <ticket> --base <base> --source-tip <tip>` and read the prepared evidence identity it returns. Verify the fence, dependency completion, and expected source tip before the ticket's fresh author starts.", Diagnostic: "prepared build guidance: bench-implement-spec.md dropped charge preparation or fence, dependency, and tip verification"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Require, Needle: "`bench preflight evidence <id> --verify` is the separate artifact-integrity check.", Diagnostic: "bounded build action: bench-implement-spec.md dropped the separate bench preflight evidence --verify integrity check"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Require, Needle: "Build action requires a current-action binding: act only after `bench preflight evidence <id> --check-current` binds the artifact to the current assignment and source pair.", Diagnostic: "bounded build action: bench-implement-spec.md permits build action without a current binding"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Forbid, Needle: "build <slug> --charge --ticket <ticket> --base <base> --source-tip <tip> --full", Diagnostic: "bounded build action: bench-implement-spec.md re-advertises the retired build full charge form"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Forbid, Needle: "Retrieve every required source with `bench preflight evidence <id>`, and follow each exact successor command until the stream ends.", Diagnostic: "fresh ticket author: bench-implement-spec.md restored the full evidence retrieval"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-implement-spec.md", Kind: Forbid, Needle: "Build action requires verified delivery", Diagnostic: "fresh ticket author: bench-implement-spec.md restored the verified-delivery build rule"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-write-spec.md", Kind: Forbid, Needle: "bench preflight evidence", Diagnostic: "unchanged route: bench-write-spec.md adds charge retrieval to the write-spec phase"},
}

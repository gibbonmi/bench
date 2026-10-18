package anchors

// ft311ReviewDispatchAnchors pin the review phase's prepared-evidence command lines, the
// craft-spec implementation-line labels, and the profile's route to the conditional review
// line. The Forbid rows keep the retired same-family CLI and inline-axis routes, the retained
// author, and the review full charge form out, and they reject retrieval instructions that
// conflict with the narrow-axis owner. One diagnostic prefix lets the conformance owner
// enumerate each family without a second registry.
var ft311ReviewDispatchAnchors = []Anchor{
	{Group: AfterImplementSpec, File: ".agents/commands/bench-review-implementation.md", Kind: Require, Class: Command, Needle: "Prepare the shared evidence once with `bench preflight review <slug> --charge --base <b> --source-tip <t>`.", Diagnostic: "prepared review dispatch: bench-review-implementation.md dropped the one shared-evidence preparation with bench preflight review --charge"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-review-implementation.md", Kind: Require, Class: Command, Needle: "Do not run a second `bench consumers --changed` collection for the same frozen pair.", Diagnostic: "prepared review dispatch: bench-review-implementation.md dropped the refusal of a second bench consumers --changed collection"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-review-implementation.md", Kind: Forbid, Needle: "Retrieve every byte with `bench preflight evidence <id>`", Diagnostic: "prepared review dispatch: bench-review-implementation.md contradicts the narrow-axis contract with full artifact retrieval"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-review-implementation.md", Kind: Forbid, Needle: "Read the base-relative diff facts from that prepared evidence.", Diagnostic: "prepared review dispatch: bench-review-implementation.md contradicts the narrow-axis contract with coordinator diff retrieval"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-review-implementation.md", Kind: Forbid, Needle: "Take the consumer tables and their citation row from the shared evidence's `consumers` identity.", Diagnostic: "prepared review dispatch: bench-review-implementation.md contradicts the narrow-axis contract with coordinator consumer-table retrieval"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-review-implementation.md", Kind: Forbid, Needle: "Walk the `touched=false` rows first, because a consumer outside the diff's file set is the FT210 class an unlisted-consumer miss hides.", Diagnostic: "prepared review dispatch: bench-review-implementation.md leaves the touched-false blast walk in coordinator context"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-review-implementation.md", Kind: Require, Class: Command, Needle: "A narrow axis runs `bench preflight evidence <id> --check-current` one time, to bind the artifact to the assignment and the source pair.", Diagnostic: "narrow review axis: bench-review-implementation.md dropped the single current binding of a narrow axis"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-review-implementation.md", Kind: Forbid, Needle: "run the same axes inline", Diagnostic: "prepared review dispatch: bench-review-implementation.md retains the inline-axis fallback"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-delegate/references/cross-harness-reviewers.md", Kind: Forbid, Needle: "falls back to its own family's CLI", Diagnostic: "prepared review dispatch: cross-harness-reviewers.md retains the own-family CLI fallback"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-spec/SKILL.md", Kind: RequireInSection, Class: Label, Reader: "craft-spec spec template", Section: "Template", Needle: "Line: <resolved implementation model id> / <starting effort>.", Diagnostic: "retained workflow: craft-spec dropped the one implementation-line template"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-spec/SKILL.md", Kind: RequireInSection, Class: Label, Reader: "craft-spec spec template", Section: "Template", Needle: "Implementation-line reason: <hardest material chunk, spec precision, seam uncertainty, and test strength>.", Diagnostic: "retained workflow: craft-spec dropped the implementation-line reason factors"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-spec/SKILL.md", Kind: RequireInSection, Class: Label, Reader: "craft-spec spec template", Section: "Template", Needle: "Harder chunks: <chunk IDs or none>.", Diagnostic: "retained workflow: craft-spec dropped the harder-chunk annotation"},
	{Group: AfterImplementSpec, File: "projects/benchkit.md", Kind: Require, Class: Route, Needle: "conditional review line in `craft-line`", Diagnostic: "retained workflow: project profile dropped the conditional review-line reference"},
	{Group: AfterImplementSpec, File: ".bench/BENCH.md", Kind: Forbid, Needle: "The retained implementation session writes production changes, tests, probes, and repairs.", Diagnostic: "fresh ticket author: operating guide restored the retained implementation session"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-review-implementation.md", Kind: Forbid, Needle: "Accepted findings return to the retained `/bench-implement-spec` session.", Diagnostic: "fresh ticket author: review phase restored the finding return to the retained session"},
	{Group: AfterImplementSpec, File: "docs/field-guide.html", Kind: Forbid, Needle: "the approved implementation line retains one author through production changes, tests, probes, repairs, and chunk reviews.", Diagnostic: "fresh ticket author: field guide restored the one retained author on the approved implementation line"},
	{Group: AfterImplementSpec, File: ".bench/BENCH.md", Kind: Forbid, Needle: "After the last chunk, the retained author reconciles overall acceptance and integration.", Diagnostic: "fresh ticket author: operating guide restored the retained author's final reconciliation"},
	{Group: AfterImplementSpec, File: ".bench/BENCH.md", Kind: Forbid, Needle: "Return findings to the retained author and obtain current repair coverage.", Diagnostic: "fresh ticket author: operating guide restored the finding return to the retained author"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-review-implementation.md", Kind: Require, Class: Command, Needle: "`bench preflight evidence <id> --verify` is the separate artifact-integrity check.", Diagnostic: "bounded review action: bench-review-implementation.md dropped the separate bench preflight evidence --verify integrity check"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-review-implementation.md", Kind: Require, Class: Command, Needle: "Review action requires a current-action binding: act only after `bench preflight evidence <id> --check-current` binds the artifact to the current assignment and source pair.", Diagnostic: "bounded review action: bench-review-implementation.md permits review action without a current binding"},
	{Group: AfterImplementSpec, File: ".agents/commands/bench-review-implementation.md", Kind: Forbid, Needle: "bench preflight review <slug> --charge --base <b> --source-tip <t> --full", Diagnostic: "bounded review action: bench-review-implementation.md re-advertises the retired review full charge form"},

	// The native-handoff row keeps the destination harness able to resume the same review.
	// The trusted identity and the exact retrieval command travel together, because the
	// originating checkout path names no evidence and authorizes no retrieval.
	{Group: AfterImplementSpec, File: ".agents/commands/bench-review-implementation.md", Kind: Require, Class: Command, Needle: "The handoff carries the trusted expected evidence identity and the exact `bench preflight evidence <id>` retrieval command, never the originating checkout path.", Diagnostic: "native handoff: bench-review-implementation.md drops the trusted evidence identity or the exact retrieval command from the capable-harness handoff"},
}

// reviewRepeatRule and finalReconciliationRule are owner sentences in `.bench/BENCH.md`. The
// retained-workflow Forbid rows that keep phase-command copies out read them.
const (
	reviewRepeatRule        = "Repeat delegated review only for a later delta or a cross-chunk concern that invalidates prior evidence."
	finalReconciliationRule = "After the last chunk, the orchestrator reconciles the final acceptance and integration."
)

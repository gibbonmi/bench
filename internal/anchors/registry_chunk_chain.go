package anchors

// chunkChainAnchors pin the chunk chain and the record order that the implement and
// review phases state, and the hardening cap that the bounded repair policy states.
var chunkChainAnchors = []Anchor{
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Build", Needle: "After each ticket commit, run `bench worktree exec <target> -- bench preflight build <slug>`.", Diagnostic: "chunk chain: build preflight runs through the worktree after each ticket commit"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Build", Needle: "The author commits the verification and probe record before the axis dispatch.", Diagnostic: "chunk chain: the probe record commits before the axis dispatch"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Land", Needle: "Plan commits land before the ticket merge, and a `main` merge lands only before the first chunk. Only record commits follow the chunk tip.", Diagnostic: "chunk chain: only record commits follow the chunk tip"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Land", Needle: "The reconciliation commit joins the review delta of the last chunk.", Diagnostic: "chunk chain: the reconciliation commit joins the last chunk delta"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Land", Needle: "Write the ordinary assessment record before the `bench worktree land` step, and append the landing evidence after it.", Diagnostic: "chunk chain: the assessment record precedes the landing"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInStep, Section: "Process", Step: 1, Needle: "A later plan commit is never a chunk base.", Diagnostic: "chunk chain: a plan commit is never a chunk base"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInStep, Section: "Process", Step: 1, Needle: "Merge `main` into the source only before the first chunk. A later chunk base holds the tree of the accepted predecessor tip, so the review chain refuses a `main` merge after the first chunk starts. The landing composes the `main` commits that arrive during the build. The review preflight counts each path that a merge brings in against the ownership fences.", Diagnostic: "chunk chain: a main merge lands only before the first chunk"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInSection, Section: "Review modes", Needle: "A chunk that ends on a repair takes one confirming round of all three axes at its final tip.", Diagnostic: "chunk chain: a repair-ending chunk takes one confirming round"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInStep, Section: "Process", Step: 4, Needle: "On one shared tree, only one axis runs tests or probes while the other axes read.", Diagnostic: "chunk chain: one axis probes a shared tree"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInStep, Section: "Process", Step: 6, Needle: "A review worktree moves to the record commit, and the frozen pair still names the source tip.", Diagnostic: "chunk chain: the frozen pair keeps the source tip after the record commit"},
	{Group: AfterImplementSpec, File: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", Kind: RequireInSection, Section: "Classification and completion", Needle: "After the acceptance rows of a chunk prove, the chunk permits at most one hardening cycle.", Diagnostic: "chunk chain: the hardening cap is one cycle"},
}

// implementPhase and reviewPhase are the phase files that state the chunk chain.
const (
	implementPhase = ".agents/commands/bench-implement-spec.md"
	reviewPhase    = ".agents/commands/bench-review-implementation.md"
)

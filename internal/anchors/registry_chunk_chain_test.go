package anchors

import "testing"

// These independent expectations make removal of a chunk-chain rule fail.
func TestChunkChainAnchors(t *testing.T) {
	const implement = ".agents/commands/bench-implement-spec.md"
	anchorHarness{group: AfterImplementSpec, rules: []anchorRule{
		{file: implement, section: "Build", needle: "After each ticket commit, run `bench worktree exec <target> -- bench preflight build <slug>`.", want: "chunk chain: build preflight runs through the worktree after each ticket commit"},
		{file: implement, section: "Build", needle: "The author commits the verification and probe record before the axis dispatch.", want: "chunk chain: the probe record commits before the axis dispatch"},
		{file: implement, section: "Land", needle: "Plan commits and `main` merges land before the ticket merge. Only record commits follow the chunk tip.", want: "chunk chain: only record commits follow the chunk tip"},
		{file: implement, section: "Land", needle: "The reconciliation commit joins the review delta of the last chunk.", want: "chunk chain: the reconciliation commit joins the last chunk delta"},
		{file: implement, section: "Land", needle: "Write the ordinary assessment record before the `bench worktree land` step, and append the landing evidence after it.", want: "chunk chain: the assessment record precedes the landing"},
		{file: ".agents/commands/bench-review-implementation.md", section: "Process", step: 1, needle: "A later plan commit is never a chunk base.", want: "chunk chain: a plan commit is never a chunk base"},
		{file: ".agents/commands/bench-review-implementation.md", section: "Review modes", needle: "A chunk that ends on a repair takes one confirming round of all three axes at its final tip.", want: "chunk chain: a repair-ending chunk takes one confirming round"},
		{file: ".agents/commands/bench-review-implementation.md", section: "Process", step: 4, needle: "On one shared tree, only one axis runs tests or probes while the other axes read.", want: "chunk chain: one axis probes a shared tree"},
		{file: ".agents/commands/bench-review-implementation.md", section: "Process", step: 6, needle: "A review worktree moves to the record commit, and the frozen pair still names the source tip.", want: "chunk chain: the frozen pair keeps the source tip after the record commit"},
		{file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Classification and completion", needle: "After the acceptance rows of a chunk prove, the chunk permits at most one hardening cycle.", want: "chunk chain: the hardening cap is one cycle"},
	}}.check(t)
}

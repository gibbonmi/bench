package anchors

import "testing"

// These independent expectations make removal of a chunk-chain rule fail.
func TestChunkChainAnchors(t *testing.T) {
	const implement = ".agents/commands/bench-implement-spec.md"
	anchorHarness{group: AfterImplementSpec, rules: []anchorRule{
		{file: implement, section: "Build", needle: "After each ticket commit, run `bench worktree exec <target> -- bench preflight build <slug>`.", want: "chunk chain: build preflight runs through the worktree after each ticket commit"},
		{file: implement, section: "Build", needle: "The author commits the verification and probe record before the axis dispatch.", want: "chunk chain: the probe record commits before the axis dispatch"},
		{file: implement, section: "Build", needle: "Prepare the review charge from that record commit. The sequence is the author record commit, then the review charge, then the axis dispatch.", want: "chunk chain: the review charge follows the author record commit and precedes the axis dispatch"},
		{file: implement, section: "Build", needle: "Prepare the review charge before the author record commit.", want: "chunk chain: the Build section puts the review charge before the author record commit", forbidden: true},
		{file: implement, section: "Land", needle: "Plan commits land before the ticket merge, and a `main` merge lands before the first chunk. When `main` moves during the build, the `main` fold lands before the completion landing and joins the review delta of the last chunk. Only record commits and comment-only corrections follow the chunk tip.", want: "chunk chain: only record and comment-only commits follow the chunk tip"},
		{file: implement, section: "Land", needle: "The reconciliation commit joins the review delta of the last chunk.", want: "chunk chain: the reconciliation commit joins the last chunk delta"},
		{file: implement, section: "Land", needle: "Write the ordinary assessment record before the `bench worktree land` step, and append the landing evidence after it.", want: "chunk chain: the assessment record precedes the landing"},
		{file: implement, section: "Land", needle: "When the orchestrator freezes a chunk after its last ticket, it records the chunk entry with `bench record chunk`.", want: "chunk chain: the orchestrator records the chunk entry at the freeze"},
		{file: implement, section: "Land", needle: "Each ticket author then writes its verification entries at that chunk source with `bench record verification`.", want: "chunk chain: each ticket author records its verification at the chunk source"},
		{file: implement, section: "`--full <spec>`", needle: "A green chunk checkpoint and its handoff refresh are not a phase exit. The orchestrator continues into the successor chunk in the same turn, and it stops only on a `craft-line` stop condition.", want: "chunk chain: a green chunk checkpoint is not a phase exit"},
		{file: ".agents/commands/bench-review-implementation.md", section: "Process", step: 1, needle: "A later plan commit is never a chunk base.", want: "chunk chain: a plan commit is never a chunk base"},
		{file: ".agents/commands/bench-review-implementation.md", section: "Process", step: 1, needle: "Merge `main` into the source before the first chunk. A later chunk base holds the tree of the accepted predecessor tip, so the review chain refuses a `main` merge between two chunks. When `main` moves during the build, the source folds `main` with `bench worktree merge --from main <target>` before the completion landing. The fold joins the review delta of the last chunk. The review preflight authorizes each path that the fold brings in unchanged from the `main` tip. It counts each path that the build changed against the ownership fences.", want: "chunk chain: a main merge lands before the first chunk or as the fold the last chunk reviews"},
		{file: ".agents/commands/bench-review-implementation.md", section: "Review modes", needle: "A chunk that ends on a repair takes one confirming round of all three axes at its final tip.", want: "chunk chain: a repair-ending chunk takes one confirming round"},
		{file: ".agents/commands/bench-review-implementation.md", section: "Process", step: 4, needle: "On one shared tree, only one axis runs tests or probes while the other axes read.", want: "chunk chain: one axis probes a shared tree"},
		{file: ".agents/commands/bench-review-implementation.md", section: "Process", step: 6, needle: "A review worktree moves to the record commit, and the frozen pair still names the source tip.", want: "chunk chain: the frozen pair keeps the source tip after the record commit"},
		{file: ".agents/commands/bench-review-implementation.md", section: "Process", step: 6, needle: "Write each completed chunk, verification, review, and amendment entry with `bench record`, not with hand-built JSON.", want: "chunk chain: step 6 writes each completed entry with bench record"},
		{file: ".agents/commands/bench-review-implementation.md", needle: "Preflight supplies the source and plan digests.", want: "chunk chain: the review phase restored the preflight digest sentence", forbidden: true},
		{file: ".agents/commands/bench-review-implementation.md", needle: "Record the performer, role, model, effort, frozen base and tip, source, state, and native result.", want: "chunk chain: the review phase restored the hand-recording field sentence", forbidden: true},
		{file: ".agents/commands/bench-review-implementation.md", needle: "Embed the minimal native excerpt and its SHA-256 digest; local logs are supplemental evidence.", want: "chunk chain: the review phase restored the hand-embedded excerpt sentence", forbidden: true},
		{file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "Classification and completion", needle: "After the acceptance rows of a chunk prove, the chunk permits at most one hardening cycle.", want: "chunk chain: the hardening cap is one cycle"},
	}}.check(t)
}

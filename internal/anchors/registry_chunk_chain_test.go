package anchors

import "testing"

// These independent expectations make removal of a chunk-chain command line or step rule, or
// the return of a retired chunk-chain sentence, fail.
func TestChunkChainAnchors(t *testing.T) {
	const implement = ".agents/commands/bench-implement-spec.md"
	anchorHarness{group: AfterImplementSpec, rules: []anchorRule{
		{file: implement, section: "Build", needle: "After each ticket commit, run `bench worktree exec <target> -- bench preflight build <slug> --base <chunk-base> --source-tip <tip>`.", want: "chunk chain: build preflight runs through the worktree with the chunk's --base and --source-tip after each ticket commit"},
		{file: implement, section: "Build", needle: "Prepare the review charge before the author record commit.", want: "chunk chain: the Build section puts the review charge before the author record commit", forbidden: true},
		{file: "docs/adr/0023-each-ticket-gets-a-fresh-author.md", needle: "and so does a phase that needs a learning's fix.", want: "chunk chain: ADR 0023 restored the unscoped light-path dispatch of a phase", forbidden: true},
		{file: implement, section: "Land", needle: "A build folds `main` with `bench worktree merge --from main <target>` before the first chunk, or in a chunk gap before the first dispatch of the next chunk.", want: "chunk chain: a build folds main only before the first chunk or in a chunk gap"},
		{file: implement, needle: "A build folds `main` only before the first chunk, at the start of the last chunk, or before the completion landing.", want: "chunk chain: the implement phase restored the main fold at the start of the last chunk", forbidden: true},
		{file: implement, needle: "The fold at the start of the last chunk comes before its first ticket dispatch.", want: "chunk chain: the implement phase restored the last-chunk fold before its first ticket dispatch", forbidden: true},
		{file: implement, needle: "When `main` moves during the build, the `main` fold and the reconciliation commit join the review delta of the last chunk.", want: "chunk chain: the implement phase restored the fold in the review delta of the last chunk", forbidden: true},
		{file: implement, needle: "A fold before the completion landing follows `/bench-review-implementation` step 1, and the last chunk refreezes.", want: "chunk chain: the implement phase restored the refreeze of the last chunk after a fold", forbidden: true},
		{file: ".agents/skills/bench-craft-tickets/SKILL.md", needle: "After each chunk, Review freezes the chunk base and the current tip.", want: "chunk chain: craft-tickets restored the review as the chunk freezer", forbidden: true},
		{file: ".bench/BENCH-reference.md", needle: "Review freezes each chunk's base and current tip.", want: "chunk chain: the reference restored the review as the chunk freezer", forbidden: true},
		{file: implement, needle: "When the orchestrator freezes a chunk after its last ticket", want: "chunk chain: the implement phase restored the orchestrator's chunk record", forbidden: true},
		{file: implement, section: "Land", needle: "Before a successor chunk starts, run `bench gate --checkpoint specs/<slug>/spec.md --chunk <id>`.", want: "chunk chain: the chunk checkpoint gate runs before a successor chunk"},
		{file: implement, needle: "--chunk <id>` before the successor", want: "chunk chain: the implement phase restored the chunk checkpoint gate for every chunk", forbidden: true},
		{file: "README.md", needle: "\"none; chunk checkpoint\"", want: "chunk chain: the README chunk flow restored the chunk checkpoint on the no-successor edge", forbidden: true},
		{file: "README.md", needle: "chunk checkpoint\"| accepted", want: "chunk chain: the README chunk flow put the chunk checkpoint on the edge into accepted", forbidden: true},
		{file: ".agents/commands/bench-review-implementation.md", section: "Process", step: 1, needle: "A later plan commit is never a chunk base.", want: "chunk chain: a plan commit is never a chunk base"},
		{file: ".agents/commands/bench-review-implementation.md", section: "Process", step: 1, needle: "Merge `main` into the source before the first chunk. `/bench-implement-spec` owns each later fold point and its `bench record fold` step. The review chain accepts a recorded fold between two chunks and after the last chunk.", want: "chunk chain: a main merge lands before the first chunk or as a recorded fold"},
		{file: ".agents/commands/bench-review-implementation.md", needle: "A later chunk base holds the tree of the accepted predecessor tip, so the review chain refuses a `main` merge between two chunks.", want: "chunk chain: the review phase restored the chain refusal of a main merge between two chunks", forbidden: true},
		{file: ".agents/commands/bench-review-implementation.md", needle: "The fold joins the review delta of the last chunk.", want: "chunk chain: the review phase restored the fold in the review delta of the last chunk", forbidden: true},
		{file: ".agents/commands/bench-review-implementation.md", needle: "look first for duplicated knowledge in the repair delta", want: "chunk chain: the review phase restored the Standards first look that widens a confirming round", forbidden: true},
		{file: ".agents/skills/bench-craft-delegate/references/delegation-discipline.md", needle: "asks only for findings above the blocking bar", want: "chunk chain: delegation-discipline restored the blocking-bar scope of a re-review", forbidden: true},
		{file: ".agents/commands/bench-review-implementation.md", section: "Process", step: 4, needle: "In the default mode, spawn the axes in parallel sub-agents.", want: "chunk chain: review step 4 dropped the parallel dispatch of the axes"},
		{file: ".agents/commands/bench-review-implementation.md", section: "Process", step: 4, needle: "On one shared tree, only one axis runs tests or probes while the other axes read.", want: "chunk chain: one axis probes a shared tree"},
		{file: ".agents/commands/bench-review-implementation.md", section: "Process", step: 6, needle: "A review worktree moves to the record commit, and the frozen pair still names the source tip.", want: "chunk chain: the frozen pair keeps the source tip after the record commit"},
		{file: ".agents/commands/bench-review-implementation.md", section: "Process", step: 6, needle: "Write each completed chunk, verification, review, and amendment entry with `bench record`, not with hand-built JSON.", want: "chunk chain: step 6 writes each completed entry with bench record"},
		{file: ".agents/commands/bench-review-implementation.md", section: "Process", step: 6, needle: "The last ticket's author writes the frozen pair and its digests with `bench record ticket`, and `bench record review` appends one axis result.", want: "chunk chain: step 6 dropped the last ticket's author as the writer of the frozen pair"},
		{file: ".agents/commands/bench-review-implementation.md", section: "Process", step: 6, needle: "`bench record chunk` writes one chunk entry, `bench record excerpt` stages the last retained response as one excerpt, and `bench record ticket` writes one ticket's entries.", want: "chunk chain: step 6 dropped the chunk, excerpt, or ticket record form"},
		{file: ".agents/commands/bench-review-implementation.md", needle: "Preflight supplies the source and plan digests.", want: "chunk chain: the review phase restored the preflight digest sentence", forbidden: true},
		{file: ".agents/commands/bench-review-implementation.md", needle: "Record the performer, role, model, effort, frozen base and tip, source, state, and native result.", want: "chunk chain: the review phase restored the hand-recording field sentence", forbidden: true},
		{file: ".agents/commands/bench-review-implementation.md", needle: "Embed the minimal native excerpt and its SHA-256 digest; local logs are supplemental evidence.", want: "chunk chain: the review phase restored the hand-embedded excerpt sentence", forbidden: true},
	}}.check(t)
}

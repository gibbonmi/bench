package anchors

import "testing"

// These independent expectations make removal of a chunk-chain command line, or
// the return of a retired chunk-chain sentence, fail.
func TestChunkChainAnchors(t *testing.T) {
	const (
		implement    = ".agents/commands/bench-implement-spec.md"
		repairPolicy = ".agents/skills/bench-craft-line/references/bounded-repair-policy.md"
		reference    = ".bench/BENCH-reference.md"
	)
	anchorHarness{group: AfterImplementSpec, rules: []anchorRule{
		{file: implement, section: "Build", needle: "After each ticket commit, run `bench worktree exec <target> -- bench preflight build <slug> --base <chunk-base> --source-tip <tip>`.", want: "chunk chain: the implement phase restored the exec build preflight after a ticket commit", forbidden: true},
		{file: implement, section: "Build", needle: "Prepare the review charge before the author record commit.", want: "chunk chain: the Build section puts the review charge before the author record commit", forbidden: true},
		{file: "docs/adr/0023-each-ticket-gets-a-fresh-author.md", needle: "and so does a phase that needs a learning's fix.", want: "chunk chain: ADR 0023 restored the unscoped light-path dispatch of a phase", forbidden: true},
		{file: implement, section: "Land", needle: "A build folds `main` with `bench worktree merge --from main <target>` before the first chunk, or in a chunk gap before the first dispatch of the next chunk.", want: "chunk chain: a build folds main only before the first chunk or in a chunk gap"},
		{file: implement, needle: "A build folds `main` only before the first chunk, at the start of the last chunk, or before the completion landing.", want: "chunk chain: the implement phase restored the main fold at the start of the last chunk", forbidden: true},
		{file: implement, needle: "The fold at the start of the last chunk comes before its first ticket dispatch.", want: "chunk chain: the implement phase restored the last-chunk fold before its first ticket dispatch", forbidden: true},
		{file: implement, needle: "When `main` moves during the build, the `main` fold and the reconciliation commit join the review delta of the last chunk.", want: "chunk chain: the implement phase restored the fold in the review delta of the last chunk", forbidden: true},
		{file: implement, needle: "A fold before the completion landing follows `/bench-review-implementation` step 1, and the last chunk refreezes.", want: "chunk chain: the implement phase restored the refreeze of the last chunk after a fold", forbidden: true},
		{file: implement, needle: "Only record commits, recorded folds, and comment-only corrections follow the chunk tip.", want: "chunk chain: the implement phase restored the comment-only correction after the chunk tip", forbidden: true},
		{file: implement, needle: "Each source change except a record commit or a recorded fold moves the chunk tip and takes a refreeze of that chunk.", want: "chunk chain: the implement phase restored the refreeze of a comment-only correction", forbidden: true},
		{file: repairPolicy, needle: "A comment-only correction before the green checkpoint stays evidence-only and still takes a refreeze.", want: "chunk chain: the bounded repair policy restored the refreeze of a comment-only correction", forbidden: true},
		{file: ".agents/skills/bench-craft-tickets/SKILL.md", needle: "After each chunk, Review freezes the chunk base and the current tip.", want: "chunk chain: craft-tickets restored the review as the chunk freezer", forbidden: true},
		{file: ".bench/BENCH-reference.md", needle: "Review freezes each chunk's base and current tip.", want: "chunk chain: the reference restored the review as the chunk freezer", forbidden: true},
		{file: implement, needle: "When the orchestrator freezes a chunk after its last ticket", want: "chunk chain: the implement phase restored the orchestrator's chunk record", forbidden: true},
		{file: implement, section: "Land", needle: "Before a successor chunk starts, run `bench gate --checkpoint specs/<slug>/spec.md --chunk <id>`.", want: "chunk chain: the chunk checkpoint gate runs before a successor chunk"},
		{file: implement, needle: "--chunk <id>` before the successor", want: "chunk chain: the implement phase restored the chunk checkpoint gate for every chunk", forbidden: true},
		{file: "README.md", needle: "\"none; chunk checkpoint\"", want: "chunk chain: the README chunk flow restored the chunk checkpoint on the no-successor edge", forbidden: true},
		{file: "README.md", needle: "chunk checkpoint\"| accepted", want: "chunk chain: the README chunk flow put the chunk checkpoint on the edge into accepted", forbidden: true},
		{file: ".agents/commands/bench-review-implementation.md", needle: "A later chunk base holds the tree of the accepted predecessor tip, so the review chain refuses a `main` merge between two chunks.", want: "chunk chain: the review phase restored the chain refusal of a main merge between two chunks", forbidden: true},
		{file: ".agents/commands/bench-review-implementation.md", needle: "The fold joins the review delta of the last chunk.", want: "chunk chain: the review phase restored the fold in the review delta of the last chunk", forbidden: true},
		{file: ".agents/commands/bench-review-implementation.md", needle: "look first for duplicated knowledge in the repair delta", want: "chunk chain: the review phase restored the Standards first look that widens a confirming round", forbidden: true},
		{file: ".agents/skills/bench-craft-delegate/references/delegation-discipline.md", needle: "asks only for findings above the blocking bar", want: "chunk chain: delegation-discipline restored the blocking-bar scope of a re-review", forbidden: true},
		{file: ".agents/commands/bench-review-implementation.md", needle: "Preflight supplies the source and plan digests.", want: "chunk chain: the review phase restored the preflight digest sentence", forbidden: true},
		{file: ".agents/commands/bench-review-implementation.md", needle: "Record the performer, role, model, effort, frozen base and tip, source, state, and native result.", want: "chunk chain: the review phase restored the hand-recording field sentence", forbidden: true},
		{file: ".agents/commands/bench-review-implementation.md", needle: "Embed the minimal native excerpt and its SHA-256 digest; local logs are supplemental evidence.", want: "chunk chain: the review phase restored the hand-embedded excerpt sentence", forbidden: true},
		{file: implement, needle: "A repair session runs each check that its ticket owes once at its repair commit, and its `bench record ticket` call refreezes the chunk at that commit.", want: "chunk chain: the implement phase restored the full rerun of a repair session", forbidden: true},
		{file: ".bench/BENCH.md", needle: "A post-review repair goes to a fresh repair session for each affected ticket, and that session reruns the ticket's verification with current repair coverage.", want: "chunk chain: the operating guide restored the full rerun of a repair session", forbidden: true},
		{file: implement, section: "Land", needle: "When a chunk freezes, start `bench gate --checkpoint specs/<slug>/spec.md --chunk <id> --prepare` beside the review axes, but only when each axis probes in its own worktree.", want: "chunk chain: the Land section restored the --prepare flag of the retired prepared run", forbidden: true},
		{file: implement, section: "Land", needle: "When a repair commits, run the coordinator probe first, then start that prepared run beside the confirmation.", want: "chunk chain: the Land section restored the prepared run after a repair commit", forbidden: true},
		{file: implement, section: "Land", needle: "Wait on the PID of the prepared run before this checkpoint, because a checkpoint refuses while that run holds the gate lock.", want: "chunk chain: the Land section restored the wait on the PID of a prepared run", forbidden: true},
		{file: repairPolicy, needle: "A prepared chunk run excludes the review record from its subject.", want: "chunk chain: the bounded repair policy restored the subject of a prepared chunk run", forbidden: true},
		{file: repairPolicy, needle: "The Land section of `/bench-implement-spec` orders the coordinator probe, the prepared run, and the chunk checkpoint.", want: "chunk chain: the bounded repair policy restored the order of the prepared run", forbidden: true},
		{file: reference, needle: "Each call without `--prepare` still grades the record evidence and the record prose.", want: "chunk chain: the reference restored the --prepare flag of the retired prepared run", forbidden: true},
		{file: reference, needle: "The Land section of `/bench-implement-spec` says when the orchestrator starts the prepared run of the chunk checkpoint.", want: "chunk chain: the reference restored the retired prepared run", forbidden: true},
		{file: reference, needle: "It skips the evidence pass, the record prose grade, and the tip-moved guard.", want: "chunk chain: the reference restored the list of what a prepared run skips", forbidden: true},
		{file: "projects/benchkit.md", needle: "A prepared chunk run writes a green under that key, or a delta green, before the record commit.", want: "chunk chain: the profile restored the retired prepared chunk run", forbidden: true},
		{file: "docs/adr/0002-accepted-enforcement-postures.md", needle: "A prepared chunk run grades the same subject and writes its green under the same key.", want: "chunk chain: ADR 0002 restored the retired prepared chunk run", forbidden: true},
	}}.check(t)
}

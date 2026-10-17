package anchors

// chunkChainAnchors pin the chunk chain command lines of the implement phase: the post-ticket
// build preflight, the main fold, and the chunk checkpoint gate. They forbid the retired
// sentences of the chain: the review charge before the author record, the unscoped light-path
// dispatch in ADR 0023, the main fold at the start of the last chunk and its refreeze, the
// review as the chunk freezer, the orchestrator's chunk record, the checkpoint gate for every
// chunk, the README checkpoint edges, the review chain's refusal of a fold, the widened
// confirming round, and the hand-recording sentences. The step-scoped rows pin review steps
// 1, 4, and 6.
var chunkChainAnchors = []Anchor{
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Build", Needle: "After each ticket commit, run `bench worktree exec <target> -- bench preflight build <slug> --base <chunk-base> --source-tip <tip>`.", Diagnostic: "chunk chain: build preflight runs through the worktree with the chunk's --base and --source-tip after each ticket commit"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: ForbidInSection, Section: "Build", Needle: "review charge before the author record commit", Diagnostic: "chunk chain: the Build section puts the review charge before the author record commit"},
	{Group: AfterImplementSpec, File: freshAuthorADR, Kind: Forbid, Needle: "and so does a phase that needs a learning's fix.", Diagnostic: "chunk chain: ADR 0023 restored the unscoped light-path dispatch of a phase"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Land", Needle: "A build folds `main` with `bench worktree merge --from main <target>` before the first chunk, or in a chunk gap before the first dispatch of the next chunk.", Diagnostic: "chunk chain: a build folds main only before the first chunk or in a chunk gap"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "folds `main` only before the first chunk, at the start of the last chunk", Diagnostic: "chunk chain: the implement phase restored the main fold at the start of the last chunk"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "The fold at the start of the last chunk comes before its first ticket dispatch.", Diagnostic: "chunk chain: the implement phase restored the last-chunk fold before its first ticket dispatch"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "review delta of the last chunk", Diagnostic: "chunk chain: the implement phase restored the fold in the review delta of the last chunk"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "the last chunk refreezes", Diagnostic: "chunk chain: the implement phase restored the refreeze of the last chunk after a fold"},
	{Group: AfterImplementSpec, File: craftTickets, Kind: ForbidCaseFoldedEmphasis, Needle: "review freezes", Diagnostic: "chunk chain: craft-tickets restored the review as the chunk freezer"},
	{Group: AfterImplementSpec, File: benchReference, Kind: ForbidCaseFoldedEmphasis, Needle: "review freezes", Diagnostic: "chunk chain: the reference restored the review as the chunk freezer"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "When the orchestrator freezes a chunk after its last ticket", Diagnostic: "chunk chain: the implement phase restored the orchestrator's chunk record"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Land", Needle: "Before a successor chunk starts, run `bench gate --checkpoint specs/<slug>/spec.md --chunk <id>`.", Diagnostic: "chunk chain: the chunk checkpoint gate runs before a successor chunk"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "--chunk <id>` before the successor", Diagnostic: "chunk chain: the implement phase restored the chunk checkpoint gate for every chunk"},
	{Group: AfterImplementSpec, File: "README.md", Kind: Forbid, Needle: "\"none; chunk checkpoint\"", Diagnostic: "chunk chain: the README chunk flow restored the chunk checkpoint on the no-successor edge"},
	{Group: AfterImplementSpec, File: "README.md", Kind: Forbid, Needle: "chunk checkpoint\"| accepted", Diagnostic: "chunk chain: the README chunk flow put the chunk checkpoint on the edge into accepted"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInStep, Section: "Process", Step: 1, Needle: "A later plan commit is never a chunk base.", Diagnostic: "chunk chain: a plan commit is never a chunk base"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInStep, Section: "Process", Step: 1, Needle: "Merge `main` into the source before the first chunk. `/bench-implement-spec` owns each later fold point and its `bench record fold` step. The review chain accepts a recorded fold between two chunks and after the last chunk.", Diagnostic: "chunk chain: a main merge lands before the first chunk or as a recorded fold"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: Forbid, Needle: "so the review chain refuses a `main` merge between two chunks", Diagnostic: "chunk chain: the review phase restored the chain refusal of a main merge between two chunks"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: Forbid, Needle: "review delta of the last chunk", Diagnostic: "chunk chain: the review phase restored the fold in the review delta of the last chunk"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: Forbid, Needle: "look first for duplicated knowledge in the repair delta", Diagnostic: "chunk chain: the review phase restored the Standards first look that widens a confirming round"},
	{Group: AfterImplementSpec, File: delegationDiscipline, Kind: Forbid, Needle: "asks only for findings above the blocking bar", Diagnostic: "chunk chain: delegation-discipline restored the blocking-bar scope of a re-review"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInStep, Section: "Process", Step: 4, Needle: "In the default mode, spawn the axes in parallel sub-agents.", Diagnostic: "chunk chain: review step 4 dropped the parallel dispatch of the axes"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInStep, Section: "Process", Step: 4, Needle: "On one shared tree, only one axis runs tests or probes while the other axes read.", Diagnostic: "chunk chain: one axis probes a shared tree"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInStep, Section: "Process", Step: 6, Needle: "A review worktree moves to the record commit, and the frozen pair still names the source tip.", Diagnostic: "chunk chain: the frozen pair keeps the source tip after the record commit"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInStep, Section: "Process", Step: 6, Needle: "Write each completed chunk, verification, review, and amendment entry with `bench record`, not with hand-built JSON.", Diagnostic: "chunk chain: step 6 writes each completed entry with bench record"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInStep, Section: "Process", Step: 6, Needle: "The last ticket's author writes the frozen pair and its digests with `bench record ticket`, and `bench record review` appends one axis result.", Diagnostic: "chunk chain: step 6 dropped the last ticket's author as the writer of the frozen pair"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInStep, Section: "Process", Step: 6, Needle: "`bench record chunk` writes one chunk entry, `bench record excerpt` stages the last retained response as one excerpt, and `bench record ticket` writes one ticket's entries.", Diagnostic: "chunk chain: step 6 dropped the chunk, excerpt, or ticket record form"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: Forbid, Needle: "Preflight supplies the source and plan digests.", Diagnostic: "chunk chain: the review phase restored the preflight digest sentence"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: Forbid, Needle: "Record the performer, role, model, effort, frozen base and tip, source, state, and native result.", Diagnostic: "chunk chain: the review phase restored the hand-recording field sentence"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: Forbid, Needle: "Embed the minimal native excerpt and its SHA-256 digest; local logs are supplemental evidence.", Diagnostic: "chunk chain: the review phase restored the hand-embedded excerpt sentence"},
}

// implementPhase and reviewPhase are the phase files that state the chunk chain.
// freshAuthorADR, craftTickets, and benchReference carry retired sentences that the Forbid rows refuse.
const (
	implementPhase = ".agents/commands/bench-implement-spec.md"
	reviewPhase    = ".agents/commands/bench-review-implementation.md"
	craftTickets   = ".agents/skills/bench-craft-tickets/SKILL.md"
	benchReference = ".bench/BENCH-reference.md"
	freshAuthorADR = "docs/adr/0023-each-ticket-gets-a-fresh-author.md"
)

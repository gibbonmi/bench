package anchors

// chunkChainAnchors pin the chunk chain and the record order that the implement and
// review phases state, and the hardening cap that the bounded repair policy states.
// The first row pins the build-entry route that sends a staleness-class preflight red
// to the staleness pass before the chain starts. The next two rows pin the implement
// phase's pointer to the staleness-pass reference and the cheap-tier fan-out charge that
// the reference states. The fold rows pin the build's fold points and its gap fold in the
// implement phase. They also pin that phase's pointer to the chunk gap notes and its fold
// record. They also pin the recorded
// fold in review step 1. They forbid the retired sentences that fold main at the start of the
// last chunk or refuse a fold between two chunks. They also forbid the retired sentences that put
// a fold in the review delta of the last chunk or refreeze that chunk. A fold row also pins the
// craft-delegate sentence that leaves the build's fold timing to the implement phase.
// The freezer rows pin the last ticket's author as the chunk freezer in craft-tickets and the
// reference, and forbid the retired review freeze in both.
// The light-path rows pin the build's light-path hold in the implement phase and in ADR 0023. They also forbid the ADR's
// retired unscoped clause. The chunk-record rows pin the last ticket author's freeze and chunk
// record. They pin the chunk checkpoint gate before a successor chunk and the last chunk's close
// through the complete checkpoint. They also forbid each retired sentence. The build-log rows pin
// the orchestrator's one build log, its line for each saved step, and its retro reader.
// The `bench record` rows route each completed review-record entry to the verb, name the
// excerpt and ticket forms, and forbid the retired hand-recording sentences.
// The review-round rows pin the parallel axis dispatch, the combined repair charge, and the
// named-findings scope of a confirming round. They also forbid the two retired sentences that widen that scope.
// The confirmation rows pin the bounded repair policy's test-only exception and its finding-list
// disposition, and the review phase's matching exception to the confirming round. The policy rows
// keep the implementation-continuation prefix, so that family's twins select them.
var chunkChainAnchors = []Anchor{
	{Group: AfterImplementSpec, File: implementPhase, Kind: Require, Needle: "If every red row is a `*-closure` row, `fence-writes`, `completion-plan`, or `check-floor`, the staleness pass below takes the red. Any other red stops the phase.", Diagnostic: ".agents/commands/bench-implement-spec.md dropped or widened the staleness route for a red build preflight"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Declare the line, validate the tickets, route the venue", Needle: "charge the staleness pass that `" + stalenessPass + "` states.", Diagnostic: ".agents/commands/bench-implement-spec.md dropped the staleness-pass charge pointer to its reference"},
	{Group: AfterImplementSpec, File: stalenessPass, Kind: RequireInSection, Section: "The charge", Needle: "- Line: the cheap tier, default effort, one iteration for each delegate.\n- Access: read-only. A delegate writes no file.", Diagnostic: stalenessPass + " dropped the cheap-tier, read-only staleness fan-out charge"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Build", Needle: "After each ticket commit, run `bench worktree exec <target> -- bench preflight build <slug> --base <chunk-base> --source-tip <tip>`.", Diagnostic: "chunk chain: build preflight runs through the worktree with the chunk's --base and --source-tip after each ticket commit"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Build", Needle: "The author commits the verification and probe record before the axis dispatch.", Diagnostic: "chunk chain: the probe record commits before the axis dispatch"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Build", Needle: "Prepare the review charge from that record commit. The sequence is the author record commit, then the review charge, then the axis dispatch.", Diagnostic: "chunk chain: the review charge follows the author record commit and precedes the axis dispatch"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: ForbidInSection, Section: "Build", Needle: "review charge before the author record commit", Diagnostic: "chunk chain: the Build section puts the review charge before the author record commit"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Build", Needle: "A build dispatches no light-path fix while it runs.", Diagnostic: "chunk chain: the build dispatches no light-path fix"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Build", Needle: "A fix that the build's own outcome needs lands as the build's own commit under `.bench/BENCH.md`'s fix-don't-park rule.", Diagnostic: "chunk chain: a fix the build needs lands as the build's own commit"},
	{Group: AfterImplementSpec, File: freshAuthorADR, Kind: Require, Needle: "A build dispatches no light-path fix while it runs.", Diagnostic: "chunk chain: ADR 0023 dropped the build's light-path hold"},
	{Group: AfterImplementSpec, File: freshAuthorADR, Kind: Forbid, Needle: "and so does a phase that needs a learning's fix.", Diagnostic: "chunk chain: ADR 0023 restored the unscoped light-path dispatch of a phase"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Land", Needle: "Plan commits land before the ticket merge, and a `main` merge lands before the first chunk or as a recorded fold. Only record commits, recorded folds, and comment-only corrections follow the chunk tip.", Diagnostic: "chunk chain: only record commits, recorded folds, and comment-only commits follow the chunk tip"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Land", Needle: "A build folds `main` with `bench worktree merge --from main <target>` before the first chunk, or in a chunk gap before the first dispatch of the next chunk.", Diagnostic: "chunk chain: a build folds main only before the first chunk or in a chunk gap"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Land", Needle: "The Command Notes of `.bench/BENCH-reference.md` state how a chunk gap opens and closes.", Diagnostic: "chunk chain: the implement phase dropped its pointer to the chunk gap notes"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Land", Needle: "Each fold after the first chunk takes `bench record fold` before the next commit.", Diagnostic: "chunk chain: each fold after the first chunk takes bench record fold"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Land", Needle: "After the last freeze, `main` moves only through the reviewer's landing override, and a recorded fold covers the completion source without a refreeze.", Diagnostic: "chunk chain: after the last freeze main moves only through the landing override"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Land", Needle: "A build never folds `main` between two tickets of a chunk.", Diagnostic: "chunk chain: no main fold between two tickets of a chunk"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Land", Needle: "Inside a chunk, the post-ticket preflight passes the chunk's `--base` and the current `--source-tip`, so a moved `main` does not route a fold.", Diagnostic: "chunk chain: the in-chunk preflight passes the chunk base and the current source tip"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "folds `main` only before the first chunk, at the start of the last chunk", Diagnostic: "chunk chain: the implement phase restored the main fold at the start of the last chunk"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "The fold at the start of the last chunk comes before its first ticket dispatch.", Diagnostic: "chunk chain: the implement phase restored the last-chunk fold before its first ticket dispatch"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "review delta of the last chunk", Diagnostic: "chunk chain: the implement phase restored the fold in the review delta of the last chunk"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "the last chunk refreezes", Diagnostic: "chunk chain: the implement phase restored the refreeze of the last chunk after a fold"},
	{Group: AfterImplementSpec, File: craftDelegate, Kind: Require, Needle: "During a build, `/bench-implement-spec` owns the fold timing, and a ticket charge starts without this fold.", Diagnostic: "chunk chain: craft-delegate leaves the build fold timing to the implement phase"},
	{Group: AfterImplementSpec, File: craftTickets, Kind: Forbid, Needle: "review freezes", Diagnostic: "chunk chain: craft-tickets restored the review as the chunk freezer"},
	{Group: AfterImplementSpec, File: benchReference, Kind: Forbid, Needle: "review freezes", Diagnostic: "chunk chain: the reference restored the review as the chunk freezer"},
	{Group: AfterImplementSpec, File: craftTickets, Kind: Require, Needle: "In each chunk, the last ticket's author freezes the chunk base and the current tip.", Diagnostic: "chunk chain: craft-tickets dropped the last ticket's author as the chunk freezer"},
	{Group: AfterImplementSpec, File: benchReference, Kind: Require, Needle: "The last ticket's author freezes each chunk's base and current tip.", Diagnostic: "chunk chain: the reference dropped the last ticket's author as the chunk freezer"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Land", Needle: "After the last ticket of a chunk, that ticket's author freezes the chunk and records the chunk entry with `bench record ticket`.", Diagnostic: "chunk chain: the last ticket's author records the chunk entry at the freeze"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Build", Needle: "The last ticket's author performs that freeze, as the Land section states.", Diagnostic: "chunk chain: the Build section dropped the last ticket's author as the freeze owner"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "When the orchestrator freezes a chunk after its last ticket", Diagnostic: "chunk chain: the implement phase restored the orchestrator's chunk record"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Land", Needle: "Before a successor chunk starts, run `bench gate --checkpoint specs/<slug>/spec.md --chunk <id>`.", Diagnostic: "chunk chain: the chunk checkpoint gate runs before a successor chunk"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "--chunk <id>` before the successor", Diagnostic: "chunk chain: the implement phase restored the chunk checkpoint gate for every chunk"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Land", Needle: "The last chunk runs no chunk checkpoint gate, because the complete checkpoint that `/bench-final-check` owns grades every chunk again.", Diagnostic: "chunk chain: the last chunk closes through the complete checkpoint alone"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Land", Needle: "The last chunk keeps its chunk entry, its verification entries, and its three review entries.", Diagnostic: "chunk chain: the last chunk keeps its chunk, verification, and review entries"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Land", Needle: "The orchestrator keeps one build log at the primary checkout's `capture/build-log/<slug>.log`.", Diagnostic: "chunk chain: the orchestrator keeps one build log in the primary checkout's capture folder"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Land", Needle: "For each reviewer direction that saves a step, it appends one line that names the direction and the step that it removed.", Diagnostic: "chunk chain: the build log takes one line for each reviewer direction that saves a step"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Land", Needle: "The retro author reads that log.", Diagnostic: "chunk chain: the retro author reads the build log"},
	{Group: AfterImplementSpec, File: freshAuthorADR, Kind: Require, Needle: "It refreshes the session handoff at each chunk checkpoint and at the final reconciliation", Diagnostic: "chunk chain: ADR 0023 dropped the handoff refresh at the final reconciliation"},
	{Group: AfterImplementSpec, File: "README.md", Kind: Require, Needle: "more -->|\"yes; chunk checkpoint\"| chunk", Diagnostic: "chunk chain: the README chunk flow dropped the chunk checkpoint on the successor edge"},
	{Group: AfterImplementSpec, File: "README.md", Kind: Forbid, Needle: "\"none; chunk checkpoint\"", Diagnostic: "chunk chain: the README chunk flow restored the chunk checkpoint on the no-successor edge"},
	{Group: AfterImplementSpec, File: "README.md", Kind: Forbid, Needle: "chunk checkpoint\"| accepted", Diagnostic: "chunk chain: the README chunk flow put the chunk checkpoint on the edge into accepted"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Land", Needle: "Each ticket author then writes its verification entries at that chunk source with `bench record ticket`.", Diagnostic: "chunk chain: each ticket author records its verification at the chunk source"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "`--full <spec>`", Needle: "A green chunk checkpoint and its handoff refresh are not a phase exit. The orchestrator continues into the successor chunk in the same turn, and it stops only on a `craft-line` stop condition.", Diagnostic: "chunk chain: a green chunk checkpoint is not a phase exit"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Section: "Land", Needle: "Write the ordinary assessment record before the `bench worktree land` step, and append the landing evidence after it.", Diagnostic: "chunk chain: the assessment record precedes the landing"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInStep, Section: "Process", Step: 1, Needle: "A later plan commit is never a chunk base.", Diagnostic: "chunk chain: a plan commit is never a chunk base"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInStep, Section: "Process", Step: 1, Needle: "Merge `main` into the source before the first chunk. `/bench-implement-spec` owns each later fold point and its `bench record fold` step. The review chain accepts a recorded fold between two chunks and after the last chunk.", Diagnostic: "chunk chain: a main merge lands before the first chunk or as a recorded fold"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: Forbid, Needle: "so the review chain refuses a `main` merge between two chunks", Diagnostic: "chunk chain: the review phase restored the chain refusal of a main merge between two chunks"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: Forbid, Needle: "review delta of the last chunk", Diagnostic: "chunk chain: the review phase restored the fold in the review delta of the last chunk"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInSection, Section: "Review modes", Needle: "A chunk that ends on a repair takes one confirming round of all three axes at its final tip.", Diagnostic: "chunk chain: a repair-ending chunk takes one confirming round"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInSection, Section: "Review modes", Needle: "Its charge names those findings and the folds to confirm.", Diagnostic: "chunk chain: a confirming round's charge dropped its named findings or its folds"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: Forbid, Needle: "look first for duplicated knowledge in the repair delta", Diagnostic: "chunk chain: the review phase restored the Standards first look that widens a confirming round"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInSection, Section: "Review modes", Needle: "When every accepted finding of a round falls on one ticket, one combined repair charge carries them all.", Diagnostic: "chunk chain: a one-ticket round dropped its one combined repair charge"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInSection, Section: "Review modes", Needle: "That charge holds the findings of every axis and goes to that ticket's repair session.", Diagnostic: "chunk chain: the combined repair charge dropped the findings of every axis or the ticket's repair session"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInSection, Section: "Review modes", Needle: "When the accepted findings fall on several tickets, each affected ticket keeps its own repair session under `.bench/BENCH.md`'s repair rule.", Diagnostic: "chunk chain: a multi-ticket round dropped the repair session of each affected ticket"},
	{Group: AfterImplementSpec, File: delegationDiscipline, Kind: Require, Needle: "A re-review charge names the findings and the folds to confirm, and the round examines only those findings.", Diagnostic: "chunk chain: delegation-discipline dropped the named findings and folds of a re-review charge"},
	{Group: AfterImplementSpec, File: delegationDiscipline, Kind: Forbid, Needle: "asks only for findings above the blocking bar", Diagnostic: "chunk chain: delegation-discipline restored the blocking-bar scope of a re-review"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInStep, Section: "Process", Step: 4, Needle: "In the default mode, spawn the axes in parallel sub-agents.", Diagnostic: "chunk chain: review step 4 dropped the parallel dispatch of the axes"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInStep, Section: "Process", Step: 4, Needle: "On one shared tree, only one axis runs tests or probes while the other axes read.", Diagnostic: "chunk chain: one axis probes a shared tree"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInStep, Section: "Process", Step: 6, Needle: "A review worktree moves to the record commit, and the frozen pair still names the source tip.", Diagnostic: "chunk chain: the frozen pair keeps the source tip after the record commit"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInStep, Section: "Process", Step: 6, Needle: "Write each completed chunk, verification, review, and amendment entry with `bench record`, not with hand-built JSON.", Diagnostic: "chunk chain: step 6 writes each completed entry with bench record"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInStep, Section: "Process", Step: 6, Needle: "The last ticket's author writes the frozen pair and its digests with `bench record ticket`, and `bench record review` appends one axis result.", Diagnostic: "chunk chain: step 6 dropped the last ticket's author as the writer of the frozen pair"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInStep, Section: "Process", Step: 6, Needle: "`bench record excerpt` stages the last retained response as one excerpt, and `bench record ticket` writes one ticket's entries.", Diagnostic: "chunk chain: step 6 dropped the excerpt or the ticket record form"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: Forbid, Needle: "Preflight supplies the source and plan digests.", Diagnostic: "chunk chain: the review phase restored the preflight digest sentence"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: Forbid, Needle: "Record the performer, role, model, effort, frozen base and tip, source, state, and native result.", Diagnostic: "chunk chain: the review phase restored the hand-recording field sentence"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: Forbid, Needle: "Embed the minimal native excerpt and its SHA-256 digest; local logs are supplemental evidence.", Diagnostic: "chunk chain: the review phase restored the hand-embedded excerpt sentence"},
	{Group: AfterImplementSpec, File: boundedRepairPolicy, Kind: RequireInSection, Section: "Classification and completion", Needle: "After the acceptance rows of a chunk prove, the chunk permits at most one hardening cycle.", Diagnostic: "chunk chain: the hardening cap is one cycle"},
	{Group: AfterImplementSpec, File: boundedRepairPolicy, Kind: RequireInSection, Section: "Classification and completion", Needle: "A test-only repair whose coordinator probe bites takes no new axis results; the coordinator writes a confirmation with `bench record confirmation` instead.", Diagnostic: "implementation continuation: bounded repair dropped the test-only confirmation exception"},
	{Group: AfterImplementSpec, File: boundedRepairPolicy, Kind: RequireInSection, Section: "Classification and completion", Needle: "The confirmation's finding list is the coordinator's disposition of the findings that the repair closes.", Diagnostic: "implementation continuation: bounded repair dropped the confirmation's finding-list disposition"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: RequireInSection, Section: "Review modes", Needle: "A test-only repair whose coordinator probe bites closes through `bench record confirmation` and takes no confirming round.", Diagnostic: "chunk chain: a test-only repair with a biting coordinator probe took a confirming round"},
}

// implementPhase and reviewPhase are the phase files that state the chunk chain.
// stalenessPass is the reference that owns the build-entry staleness charge.
// freshAuthorADR is the decision record that states the build's light-path hold and the
// orchestrator's handoff refresh. boundedRepairPolicy is the reference that owns the
// repair allowance, the hardening cap, and the test-only exception. craftTickets and
// benchReference also name the chunk freezer.
const (
	implementPhase      = ".agents/commands/bench-implement-spec.md"
	reviewPhase         = ".agents/commands/bench-review-implementation.md"
	craftTickets        = ".agents/skills/bench-craft-tickets/SKILL.md"
	benchReference      = ".bench/BENCH-reference.md"
	stalenessPass       = ".agents/skills/bench-implement-spec/references/staleness-pass.md"
	freshAuthorADR      = "docs/adr/0023-each-ticket-gets-a-fresh-author.md"
	boundedRepairPolicy = ".agents/skills/bench-craft-line/references/bounded-repair-policy.md"
)

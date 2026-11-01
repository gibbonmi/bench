package anchors

// chunkChainAnchors pin the chunk chain command lines of the implement phase: the main fold and
// the chunk checkpoint gate. They forbid the retired sentences of the chain: the exec build
// preflight after a ticket commit, the review charge before the author record, the unscoped light-path
// dispatch in ADR 0023, the main fold at the start of the last chunk and its refreeze, the
// comment-only correction after a chunk tip, the review as the chunk freezer, the orchestrator's
// chunk record, the checkpoint gate for every chunk, the README checkpoint edges, the review
// chain's refusal of a fold, the widened confirming round, the hand-recording sentences, the
// full rerun of a repair session, and the retired prepared run of a chunk checkpoint.
// They also forbid the retired refreeze of a comment-only correction.
var chunkChainAnchors = []Anchor{
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "bench worktree exec <target> -- bench preflight build", Diagnostic: "chunk chain: the implement phase restored the exec build preflight after a ticket commit"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: ForbidInSection, Section: "Build", Needle: "review charge before the author record commit", Diagnostic: "chunk chain: the Build section puts the review charge before the author record commit"},
	{Group: AfterImplementSpec, File: freshAuthorADR, Kind: Forbid, Needle: "and so does a phase that needs a learning's fix.", Diagnostic: "chunk chain: ADR 0023 restored the unscoped light-path dispatch of a phase"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Class: Command, Section: "Land", Needle: "A build folds `main` with `bench worktree merge --from main <target>` before the first chunk, or in a chunk gap before the first dispatch of the next chunk.", Diagnostic: "chunk chain: a build folds main only before the first chunk or in a chunk gap"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "folds `main` only before the first chunk, at the start of the last chunk", Diagnostic: "chunk chain: the implement phase restored the main fold at the start of the last chunk"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "The fold at the start of the last chunk comes before its first ticket dispatch.", Diagnostic: "chunk chain: the implement phase restored the last-chunk fold before its first ticket dispatch"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "review delta of the last chunk", Diagnostic: "chunk chain: the implement phase restored the fold in the review delta of the last chunk"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "the last chunk refreezes", Diagnostic: "chunk chain: the implement phase restored the refreeze of the last chunk after a fold"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "comment-only corrections follow the chunk tip", Diagnostic: "chunk chain: the implement phase restored the comment-only correction after the chunk tip"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "Each source change except a record commit or a recorded fold moves the chunk tip", Diagnostic: "chunk chain: the implement phase restored the refreeze of a comment-only correction"},
	{Group: AfterImplementSpec, File: boundedRepairPolicy, Kind: Forbid, Needle: "takes a refreeze", Diagnostic: "chunk chain: the bounded repair policy restored the refreeze of a comment-only correction"},
	{Group: AfterImplementSpec, File: craftTickets, Kind: ForbidCaseFoldedEmphasis, Needle: "review freezes", Diagnostic: "chunk chain: craft-tickets restored the review as the chunk freezer"},
	{Group: AfterImplementSpec, File: benchReference, Kind: ForbidCaseFoldedEmphasis, Needle: "review freezes", Diagnostic: "chunk chain: the reference restored the review as the chunk freezer"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "When the orchestrator freezes a chunk after its last ticket", Diagnostic: "chunk chain: the implement phase restored the orchestrator's chunk record"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: RequireInSection, Class: Command, Section: "Land", Needle: "Before a successor chunk starts, run `bench gate --checkpoint specs/<slug>/spec.md --chunk <id>`.", Diagnostic: "chunk chain: the chunk checkpoint gate runs before a successor chunk"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "--chunk <id>` before the successor", Diagnostic: "chunk chain: the implement phase restored the chunk checkpoint gate for every chunk"},
	{Group: AfterImplementSpec, File: "README.md", Kind: Forbid, Needle: "\"none; chunk checkpoint\"", Diagnostic: "chunk chain: the README chunk flow restored the chunk checkpoint on the no-successor edge"},
	{Group: AfterImplementSpec, File: "README.md", Kind: Forbid, Needle: "chunk checkpoint\"| accepted", Diagnostic: "chunk chain: the README chunk flow put the chunk checkpoint on the edge into accepted"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: Forbid, Needle: "so the review chain refuses a `main` merge between two chunks", Diagnostic: "chunk chain: the review phase restored the chain refusal of a main merge between two chunks"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: Forbid, Needle: "review delta of the last chunk", Diagnostic: "chunk chain: the review phase restored the fold in the review delta of the last chunk"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: Forbid, Needle: "look first for duplicated knowledge in the repair delta", Diagnostic: "chunk chain: the review phase restored the Standards first look that widens a confirming round"},
	{Group: AfterImplementSpec, File: delegationDiscipline, Kind: Forbid, Needle: "asks only for findings above the blocking bar", Diagnostic: "chunk chain: delegation-discipline restored the blocking-bar scope of a re-review"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: Forbid, Needle: "Preflight supplies the source and plan digests.", Diagnostic: "chunk chain: the review phase restored the preflight digest sentence"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: Forbid, Needle: "Record the performer, role, model, effort, frozen base and tip, source, state, and native result.", Diagnostic: "chunk chain: the review phase restored the hand-recording field sentence"},
	{Group: AfterImplementSpec, File: reviewPhase, Kind: Forbid, Needle: "Embed the minimal native excerpt and its SHA-256 digest; local logs are supplemental evidence.", Diagnostic: "chunk chain: the review phase restored the hand-embedded excerpt sentence"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: Forbid, Needle: "runs each check that its ticket owes once at its repair commit", Diagnostic: "chunk chain: the implement phase restored the full rerun of a repair session"},
	{Group: AfterImplementSpec, File: operatingGuide, Kind: Forbid, Needle: "reruns the ticket's verification with current repair coverage", Diagnostic: "chunk chain: the operating guide restored the full rerun of a repair session"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: ForbidInSection, Section: "Land", Needle: "--prepare", Diagnostic: "chunk chain: the Land section restored the --prepare flag of the retired prepared run"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: ForbidInSection, Section: "Land", Needle: "then start that prepared run beside the confirmation", Diagnostic: "chunk chain: the Land section restored the prepared run after a repair commit"},
	{Group: AfterImplementSpec, File: implementPhase, Kind: ForbidInSection, Section: "Land", Needle: "Wait on the PID of the prepared run", Diagnostic: "chunk chain: the Land section restored the wait on the PID of a prepared run"},
	{Group: AfterImplementSpec, File: boundedRepairPolicy, Kind: Forbid, Needle: "A prepared chunk run excludes the review record from its subject.", Diagnostic: "chunk chain: the bounded repair policy restored the subject of a prepared chunk run"},
	{Group: AfterImplementSpec, File: boundedRepairPolicy, Kind: Forbid, Needle: "orders the coordinator probe, the prepared run, and the chunk checkpoint", Diagnostic: "chunk chain: the bounded repair policy restored the order of the prepared run"},
	{Group: AfterImplementSpec, File: benchReference, Kind: Forbid, Needle: "--prepare", Diagnostic: "chunk chain: the reference restored the --prepare flag of the retired prepared run"},
	{Group: AfterImplementSpec, File: benchReference, Kind: Forbid, Needle: "prepared run", Diagnostic: "chunk chain: the reference restored the retired prepared run"},
	{Group: AfterImplementSpec, File: benchReference, Kind: Forbid, Needle: "It skips the evidence pass, the record prose grade, and the tip-moved guard.", Diagnostic: "chunk chain: the reference restored the list of what a prepared run skips"},
	{Group: AfterImplementSpec, File: benchkitProfile, Kind: Forbid, Needle: "prepared chunk run", Diagnostic: "chunk chain: the profile restored the retired prepared chunk run"},
	{Group: AfterImplementSpec, File: enforcementPosturesADR, Kind: Forbid, Needle: "prepared chunk run", Diagnostic: "chunk chain: ADR 0002 restored the retired prepared chunk run"},
}

// implementPhase is the phase file that states the chunk chain. reviewPhase, freshAuthorADR,
// craftTickets, benchReference, boundedRepairPolicy, benchkitProfile, and enforcementPosturesADR
// carry retired sentences that the Forbid rows refuse.
const (
	implementPhase         = ".agents/commands/bench-implement-spec.md"
	reviewPhase            = ".agents/commands/bench-review-implementation.md"
	craftTickets           = ".agents/skills/bench-craft-tickets/SKILL.md"
	benchReference         = ".bench/BENCH-reference.md"
	freshAuthorADR         = "docs/adr/0023-each-ticket-gets-a-fresh-author.md"
	boundedRepairPolicy    = ".agents/skills/bench-craft-line/references/bounded-repair-policy.md"
	benchkitProfile        = "projects/benchkit.md"
	enforcementPosturesADR = "docs/adr/0002-accepted-enforcement-postures.md"
)

# FT311 landing completion

Status: staged

Roadmap: FT311

Decision source: specs/ft311-landing-completion/decisions/ft311-coordinator-work.md

Verification log: 2 iteration(s) to accept — Terra/medium reviewed through the native agent surface. Iteration one returned 15 findings with 5 blocking, and the author folded 14 and refuted 1 on the tree. Iteration two returned 1 blocking fence gap and 3 stale sentences, which the author folded after the round.

## Problem

A landing publishes the commit and then stops, so the coordinator runs five more calls to finish the boundary.
Two of those calls rebuild the promotion broker and republish its manifest by hand.
Two more plan and apply the landed-sibling cleanup.
A handoff State still needs a hand edit in the primary checkout's file.
A retrospective still needs the gate-stage timings and the repair table typed from memory.
A hand chunk of this size drops a step whenever a session ends early.

## Solution

`bench worktree land` completes the landing after publication.
It refreshes the promotion broker, then it cleans the sibling worktrees whose work this landing carries.
Each operation is one landing effect with its own recorded result and its own resume path.
`bench handoff --state-file <file>` reads the State body from a file.
`bench retro <slug> --scaffold` prints a retrospective body that carries the landing's gate-stage timings and its repair-attribution table.

This is the third of five FT311 specs.
It serves this kit and every repository that links it.
It changes neither the fold-lane behavior, nor the landing gate, nor the composition.
The publication decision and the published bytes stay free of candidate code.

## User stories

### Refresh the promotion broker

Line: gpt-5.6-terra / medium.
The effect is a ticket-sized write charge under a coordinator mutation probe, which the scorecard routes at medium.

1. As a coordinator, I want `bench worktree land` to refresh the promotion broker after publication, so that one call replaces the hand rebuild and the hand repair.
2. As a coordinator, I want the refresh reported as its own landing effect, so that publication and refresh carry separate results.
3. As a coordinator, I want a destination with no build inputs to report `not-applicable`, so that a linked installation stays untouched.
4. As a coordinator, I want a verified published executable to report the refresh complete, so that a landing costs no needless compile.
5. As a coordinator, I want a failed refresh to exit 3 with the resume command, so that I repair the installation and then continue.
6. As a coordinator, I want a failed refresh to leave the published commit and the marker alone, so that it unpublishes no work.
7. As a reviewer, I want the publication decision and the published bytes free of candidate code, so that the stable-owner rule survives this fold.
8. As a coordinator, I want the executable, its seal, and the manifest installed as one outcome, so that a partial installation cannot authenticate.

### Clean the landing's eligible siblings

Line: gpt-5.6-terra / medium.
The effect is a ticket-sized write charge under a coordinator mutation probe, and it matches its ticket's line.

9. As a coordinator, I want the landing to remove the sibling worktrees whose work it carries, so that two calls leave the close.
10. As a coordinator, I want an assignment that was already landed at the destination base left alone, so that one landing does no repository-wide maintenance.
11. As a coordinator, I want a worktree whose eligibility is not proved retained, so that no unproven removal occurs.
12. As a coordinator, I want each retained row and its reason code printed, so that I know which worktrees the explicit cleanup still owns.
13. As a coordinator, I want the cleanup left unstarted after a failed refresh, so that the worktrees stay while the installation needs repair.
14. As a coordinator, I want `bench worktree clean --landed` unchanged, so that the explicit route still covers every other landed assignment.
15. As a coordinator, I want a cleanup fault reported as the failed effect, so that the landing does not read as complete.

### Resume the unfinished effects

Line: gpt-5.6-terra / medium.
The resume path lands inside the two effect tickets, so it takes their line.

16. As a coordinator, I want `bench worktree land --resume` to continue the unfinished effects, so that a retry costs one call.
17. As a coordinator, I want a resume to leave a complete effect alone, so that a retry executes only unfinished work.
18. As a coordinator, I want a resume to publish nothing a second time, so that the published commit stays the one the review bound.
19. As a coordinator, I want a resume of a released assignment to run the effects, so that the terminal path completes.
20. As a maintainer, I want each effect's completion read from observable state, so that a resume needs no journal of its own.
21. As a coordinator, I want an interrupted refresh continued by the next resume, so that a killed process costs no repair by hand.

### Read the handoff State from a file

Line: gpt-5.6-terra / medium.
The verb composes an existing scan with one new input, and the file shapes are hostile.

22. As a coordinator, I want `bench handoff --state-file <file>` to write the file's bytes as the owned section's State, so that a draft needs no hand edit.
23. As a coordinator, I want the State scan to grade the file's pins, so that a draft cannot name a lost commit.
24. As a coordinator, I want a pin outside the section tip's ancestry refused with nothing written, so that the document keeps the bytes it had.
25. As a coordinator, I want a pin that is not a commit refused by its own fault, so that the repair is exact.
26. As a coordinator, I want an ambiguous abbreviated sha refused by its own fault, so that it does not read as prose.
27. As a coordinator, I want an absent or unreadable state file refused, so that a missing draft never reads as an empty State.
28. As a coordinator, I want a state file with an unrepresentable byte refused, so that a control byte cannot ride into the document.
29. As a coordinator, I want a state file that opens a section heading refused, so that the next run can still parse the document.
30. As a coordinator, I want an empty state file to write an empty State, so that a deliberate reset returns the scaffold guidance.
31. As an existing caller, I want a run without the flag to pass the document's State through unchanged, so that the hand-edited form still works.

### Scaffold the retrospective

Line: gpt-5.6-terra / medium.
The scaffold composes the retrospective parser with a record reader, and both owners exist.

32. As a coordinator, I want `bench retro <slug> --scaffold` to print a retrospective body, so that a delegate starts from derived facts.
33. As a maintainer, I want the scaffold's headings derived from the retrospective parser, so that the writer accepts every scaffold.
34. As a coordinator, I want the scaffold to name each gate stage with its elapsed time, so that the timings are measured.
35. As a coordinator, I want missing timing data written as unknown, so that an absent record does not block the draft.
36. As a coordinator, I want one repair-attribution row per ticket of the slug, so that no ticket is forgotten.
37. As a coordinator, I want the rounds cell and the cause cell written as unknown, so that the scaffold invents no attribution.
38. As a coordinator, I want an absent tickets directory to give one unknown row, so that the table survives a retired spec folder.
39. As a coordinator, I want the scaffold to write no file, so that the later capture call still creates the retrospective.
40. As a coordinator, I want the scaffold form and the body form refused together, so that the two content sources cannot mix.
41. As a maintainer, I want the landing's gate phases recorded under the landing's own trace, so that the scaffold reads one landing's stages.

### Keep the guidance current

Line: gpt-5.6-terra / high.
Guidance prose compounds through every session, so the leverage override applies.

42. As a teammate, I want the reference and the phase guidance to describe the folded chunk, so that a cold session runs the shorter close.
43. As a reviewer, I want the other two FT311 capabilities kept out of this build, so that each approved spec has one outcome.

## Implementation decisions

### The landing effect record

The landing gains two effects after publication: `refresh` and `cleanup`.
Both run after the existing release step, so the marker, the reconcile, the prune, and the release keep their order and their meaning.
The effects run in that order, and a failed effect stops every later effect.

The landing prints one row, `effects[2]{effect,result}`, on stdout.
The row precedes the `landed{...}` record, which stays the last stdout line for every existing reader.
A `result` cell is `complete`, `failed`, `pending`, or `not-applicable`.
`pending` means the effect never started, because an earlier effect failed.

A failed effect reaches the existing incomplete render with its own step name.
The `worktree` cell then reads `incomplete:refresh` or `incomplete:cleanup`, and the exit is 3.
The record carries the same `next=` resume command every other incomplete step carries.
The terminal policy needs no edit, because it composes the step name.

### The broker refresh

The refresh applies only when the landing destination declares Bench build inputs.
The predicate is `freshness.DeclaresBuildInputs`, the same function the landing's broker notice reads.
The notice reads it over the source worktree, and the refresh reads it at the destination root.
A destination that declares none reports `not-applicable`.

The refresh is complete when the destination's published executable verifies against the destination's own sources.
The publication owner installs the executable, its seal, and the broker manifest as one outcome, and it promotes the seal last.
So a verified seal is also the proof that the manifest of that same transaction landed, and the refresh needs no second manifest read.

The refresh first reads that predicate.
A destination that already verifies reports `complete` and builds nothing.
A destination that does not verify runs the destination's own build entry point once, then reads the predicate again.
A second read that still fails reports `failed`.

The build runs through the landing's seam set, so a test replaces it and no fixture waits on a real compile.
The seam set already carries a build field whose default is the run-binary builder, and the worktree build verb consumes it.
The refresh binds a sibling field with that same signature, because the two callers need different manifest directories.
So the existing field, its default, and the worktree build verb keep their bytes.
The refresh's default runs the build owner's subject mode, whose manifest directory is the one that owner already defaults to.
The landing states no directory of its own.

The refresh never executes the published executable.
The build owner's subject mode runs the staged candidate build once, and that run publishes the seal and the manifest.
That staged step runs under the invoked owner, after the release step.
The landing's own remaining work runs under the invoked owner, which the wrapper authenticated before the landing started.
The next landing authenticates the new executable against the manifest digest, in the wrapper, before it launches anything.

### The stable-owner boundary

The candidate tree's build entry point runs inside the refresh and nowhere else.
It runs after the release step, so it cannot reach the gate, the composition, the destination compare-and-swap, the marker, the reconcile, or the release.
The published bytes and the publication decision therefore stay free of candidate code, which is the stable-owner rule this repository already proves.

Four existing fixtures commit the Go build-input manifest, so the refresh applies to each and changes its posture.
This spec re-cuts each one and deletes none.
SOL01 asserts that the candidate script never runs; the re-cut reads its marker at the publication boundary and names the branch the script observed.
SOL04 asserts one owner process and no rebuild disclosure; the re-cut keeps both and adds the refresh result.

SOL17 asserts that a forged primary executable never runs while the landing exits 0.
The re-cut keeps the marker assertion and expects `refresh` `failed` at exit 3.
LF4 resumes a destination that declared build inputs after publication, and it asserts exit 0.
The re-cut expects `refresh` `failed` at exit 3, with the published commit unmoved.
Rows LC6 and LC7 carry the first two re-cuts, and row LC4 carries the result the last two now assert.
A script that observes the published commit proves it ran after publication.

### The landed-sibling cleanup

The cleanup covers only the assignments whose work this landing carries.
An assignment qualifies when the landed proof holds against the published commit and fails against the destination base.
The existing landed proof owner answers both questions, so the selector gains one call and no second proof.

The narrowed scope reaches the selector as one value beside the cleanup options.
An empty value is the whole-repository scope the explicit route already uses, so `bench worktree clean --landed` keeps its rows, its plan, and its fingerprint.
The narrowed scope is reachable only from the landing, which plans and applies in one process and prints no fingerprint.

Every selected row still passes the existing eligibility verdict, and a row whose eligibility is not proved is retained.
The apply re-plans each row before it touches that checkout, and it stops the set on the first drift, exactly as the explicit route does.
The landing prints the plan rows on stderr, beside its other evidence, so stdout keeps its two rows.
The cleanup reports `complete` when the apply returns no error, whatever the plan retained.

### The resume of the effects

The resume runs the same two effects after its existing steps.
It runs them for a released assignment too, so the terminal-receipt path completes rather than returning early.

Each effect's completion is a predicate over observable state, never a recorded step.
The refresh reads the published executable's seal.
The cleanup re-plans the narrowed set and finds it empty.
So a resume needs no journal, an interrupted effect resumes, and a complete effect is skipped.

The resume publishes nothing a second time.
Its publication proofs are unchanged, and the effects run only after those proofs accept the named commit.

### The handoff state file

The grammar becomes `bench handoff [--harness <name>] [--next <command>] [--state-file <path>]`.
Without the flag the command reads the document's State and puts it back unchanged, as today.
With the flag the State body is the file's bytes, and every other field stays derived.

The State scan grades the new body under the same three faults and the same repair instruction.
The three existing faults each refuse as they do today: an off-ancestry pin, a non-commit pin, and an ambiguous sha.
A refusal is raised before the document renders, so the file keeps the bytes it had.

The file is classified before it is read, through the shared control-record classifier.
An absent file, a link, a special file, and an unreadable file each answer a structured refusal at exit 1.
A file over the record limit answers the same refusal.
A byte the document cannot carry refuses through the same predicate the derived fields use, widened to the State body.
A line that would open a section heading refuses, because the next run must still parse the document.

An empty file writes an empty State.
The header then carries the scaffold guidance again, which is the existing answer for a section with no State.
The trailing newline of the file is trimmed once, as the document's own writer trims it.

### The retrospective scaffold

The grammar becomes `bench retro <slug> (--body <markdown> | --scaffold)`.
Exactly one form is required, and either both or neither is a usage error at exit 2.
`--scaffold` prints the body on stdout at exit 0 and writes no file.

The scaffold's headings come from the retrospective parser's own heading list, so the writer accepts every scaffold.
The scaffold emits no improvement item, so it proposes nothing and the destination-marker check stays green.
The repair-cause vocabulary stays in the phase guidance, and the scaffold prints no copy of it.

Under the timings heading the scaffold lists each gate stage of the landing's own authorization run with its elapsed milliseconds.
The landing is the newest completed landing span in this repository's seam record, and the stages are the phase spans of that span's trace.
The scaffold names the published commit and the trace identity, so a reader can verify the selection.
An absent record, an unreadable record, and a record with no landing span each write `unknown` under that heading at exit 0.

Under the repair-attribution heading the scaffold prints the existing three-column table.
It writes one row per ticket file of the slug, in name order, with `unknown` in the rounds cell and the cause cell.
An absent or unreadable tickets directory gives one row whose ticket cell reads `unknown`.

The record reader lives beside the record writer, so the encoded shape has one owner.
The gate scheduler `schedule` opens one span per phase, for the full gate and for the lane alike.
`withGateSpanEnv` hands the phases child the record root and the trace parent from its context.
The landing today passes a background context to the authorization call, so those spans start a trace of their own.
The landing therefore threads its own span context into that call.
The phase spans then sit under the landing span.

### Guidance

The reference states the two landing effects, the state-file form, and the scaffold form.
The final-check guidance replaces one anchored sentence, and the anchors registry needle moves with it.

The old needle is `leftover worktrees are retired by `bench worktree clean --landed`: run the plan, apply it, and carry the plan and apply result in the landing report`.
The new needle is `the landing retires the sibling worktrees whose work it carries, and `bench worktree clean --landed` retires every other landed assignment`.
The mutation fixture that quotes the old bytes moves to the new bytes in the same change.
The changelog records the three verb changes under the unreleased heading.
The glossary keeps `landing effect` as written, because that term already names this behavior.

## Testing decisions

The primary seam for the effects is the package-internal landing entry point with its seam set replaced.
That seam is how every existing follow-up-step failure is proved, and it lets a fixture answer a refresh without a compile.
The two stable-owner proofs keep the real process seam, because their claim is about what the invoked owner executes.

The handoff and the retrospective attach at their public command functions over temporary repositories.
The seam-record reader and the scaffold render attach in their own packages.
The gate stages of a landing are proved through a canned seam record, because two real runs differ in their elapsed times.

Every new test runs in the ordinary Go phase of the project gate.
One anchor pins a sentence this spec changes: the final-check landed-worktree sweep, in the after-implement-spec group.
Ticket 5 moves that needle, its diagnostic, and the `final-check-landed-worktree-sweep` mutation fixture together.
Every other guidance edit is review-owned, and this spec adds no anchor.
The full landing gate remains the code oracle.

### Seam diagram

    trigger: coordinator command from the primary checkout
        |
        v
    argv --> [ landWith / resumeLandWith ] --> effects row, landed record
                  ^ tests attach here: replaced seam set, real git fixture, fault boundary
        |
        v
    plan --> [ planLandedSet / applyLandedSet ] --> cleanup plan rows
                  ^ tests attach here: two-assignment fixture, landed proof against two commits
        |
        v
    argv --> [ handoff.Command ] --> pin block, rewritten section
                  ^ tests attach here: temp repository, state file, document read-back
        |
        v
    argv --> [ roadmap.RetroCommand ] --> scaffold body or captured file
                  ^ tests attach here: canned seam record, tickets directory, retros.Parse

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| LC1 | 1, 2 | A landing whose destination declares build inputs prints `effects[2]{effect,result}` with `refresh` `complete` and exits 0. | New TestLandRefreshesTheBrokerAfterPublication through landWith | A landing that never calls the refresh seam leaves the stub's count at 0. |
| LC2 | 3 | A destination that declares no build inputs prints `refresh` `not-applicable` and calls the build seam zero times. | New TestLandSkipsTheRefreshWithoutBuildInputs through landWith | A refresh that runs unconditionally records a call. |
| LC3 | 4 | A destination whose published executable already verifies prints `refresh` `complete` and calls the build seam zero times. | New TestLandSkipsAFreshBroker through landWith | A refresh that always builds records a call. |
| LC4 | 5 | A refresh that leaves the executable unverified prints `refresh` `failed`, `worktree=incomplete:refresh`, a `next=` resume, and exits 3. | New TestLandReportsAFailedRefresh through landWith, and the SOL17 and LF4 re-cuts through their own entry points | An exit 0, an absent `next=`, or a `released` cell fails the exact comparison. |
| LC5 | 6 | A failed refresh leaves the branch, the green marker, and the checkout at the published commit. | New TestLandReportsAFailedRefresh through landWith | Three git reads red on any unpublish or reset. |
| LC6 | 7 | The candidate build entry point runs once, after the release, and it observes the published commit on the destination branch. | Existing TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion re-cut through LandCommand | A refresh placed before the compare-and-swap makes the script observe the destination base. |
| LC7 | 7 | The invoked owner keeps one process identity through publication, release, and both effects. | Existing TestLandCommandKeepsOneOwnerProcessThroughPublicationAndRelease re-cut through LandCommand | A rebuild-and-re-execute path changes the process or prints its rebuild disclosure. |
| LC8 | 8 | A refresh that installs the executable without its seal reports `failed`. | New TestLandRefreshReadsTheSeal through landWith | A predicate that only stats the executable passes on the half-installed set. |
| LC9 | 9 | A landing that folds a sibling assignment removes that sibling's worktree and prints `cleanup` `complete`. | New TestLandCleansTheFoldedSibling through landWith | The sibling checkout must be gone and its assignment record retired. |
| LC10 | 10 | An assignment already landed at the destination base keeps its worktree and joins no cleanup row. | New TestLandLeavesAPriorLandedAssignment through landWith | A selector that reads only the default branch removes it. |
| LC11 | 11 | A sibling holding uncommitted tracked work is retained, and its checkout survives. | New TestLandRetainsAnUnprovenSibling through landWith | A removal loses the work the preservation verdict names. |
| LC12 | 12 | A retained sibling prints action `retain` and reason code `dirty` on stderr, which the preservation verdict returns. | New TestLandRetainsAnUnprovenSibling through landWith | A silent retention, or an `uncertain` code, fails the exact comparison. |
| LC13 | 13 | A failed refresh prints `cleanup` `pending` and removes no sibling checkout. | New TestLandReportsAFailedRefresh through landWith | A cleanup that runs anyway removes the sibling the operator still needs. |
| LC14 | 15 | A cleanup apply fault prints `cleanup` `failed`, `worktree=incomplete:cleanup`, and exits 3. | New TestLandReportsAFailedCleanup through landWith with the cleanup fault boundary | An exit 0 reports a landing that did not finish. |
| LC15 | 14 | `bench worktree clean --landed` plans and applies every landed assignment as today. | Existing TestCleanLandedApplyRemovesAndSettles and TestCleanLandedApplyReplansEachRowBeforeMutation | A narrowed selector reaching the explicit route drops rows the operator asked for. |
| LC16 | 16, 17 | A resume after a failed refresh runs the refresh and then the cleanup, and it calls a complete effect zero times. | New TestResumeLandCompletesTheUnfinishedEffects through resumeLandWith | A resume that repeats a complete effect records a second call. |
| LC17 | 18 | A resume moves the destination branch to no new commit. | New TestResumeLandCompletesTheUnfinishedEffects through resumeLandWith | A second publication changes the branch tip the review bound. |
| LC18 | 19 | A resume whose assignment is already released still runs both effects. | New TestResumeLandRunsTheEffectsAfterRelease through resumeLandWith | The early terminal return skips the effects and reports a complete landing. |
| LC19 | 20, 21 | A resume with no prior run's output reports each effect from the tree alone, and it finishes an interrupted refresh. | New TestResumeReadsEffectStateFromTheTree through resumeLandWith | A journal-backed decision has nothing to read and reports every effect unfinished. |
| LC20 | 2 | The effects row is the stdout line before the `landed{...}` record. | New TestLandEffectsRowPrecedesTheLandedRecord through landWith | An appended row breaks every reader of the last stdout line. |
| LC21 | 22 | `bench handoff --state-file <file>` writes the file's bytes as the owned section's State. | New TestHandoffStateFileWritesTheState through handoff.Command | A form that keeps the document's State drops the drafted body. |
| LC22 | 23, 24 | A state file pinning a commit outside the section tip's ancestry refuses, and the document keeps its bytes. | New TestHandoffStateFileRefusesAnOffAncestryPin through handoff.Command | A written document proves the refusal came after the render. |
| LC23 | 25 | A state file pinning an object that is not a commit answers the not-a-commit fault. | New TestHandoffStateFileRefusesANonCommitPin through handoff.Command | A shared headline hides which repair the writer owes. |
| LC24 | 26 | A state file abbreviating a sha that names two objects answers the ambiguous fault. | New TestHandoffStateFileRefusesAnAmbiguousPin through handoff.Command | A collapsed answer reads an ambiguous prefix as prose and writes it. |
| LC25 | 27 | An absent state file, a symbolic link, a FIFO, and an unreadable file each refuse at exit 1 with nothing written. | New TestHandoffStateFileRefusesAnUnreadableFile through handoff.Command | A plain read treats an absent file as an empty State and erases the section's State. |
| LC26 | 28 | A state file carrying an escape byte refuses at exit 1. | New TestHandoffStateFileRefusesAControlByte through handoff.Command | The byte otherwise reaches every downstream reader of the document. |
| LC27 | 29 | A state file whose line opens a section heading refuses at exit 1. | New TestHandoffStateFileRefusesASectionHeading through handoff.Command | A written heading makes the next run fail to parse the document. |
| LC28 | 30 | An empty state file writes an empty State, and the header carries the scaffold guidance. | New TestHandoffStateFileAcceptsAnEmptyFile through handoff.Command | An empty file read as absent keeps the old State instead of resetting it. |
| LC29 | 31 | A run without the flag keeps the document's State byte for byte. | New TestCommandKeepsTheOwnedStateWithoutTheFlag in the handoff package | A flag that always rewrites State loses a reviewer's own words. |
| LC30 | 32, 33 | `bench retro <slug> --scaffold` prints a body the retrospective parser accepts at exit 0. | New TestRetroScaffoldParses through RetroCommand | The parser checks order only, so the coordinator probe that swaps one heading in the exported list is the catch. |
| LC31 | 34, 41 | The scaffold lists each gate stage of the landing's own trace with its elapsed milliseconds. | New TestRetroScaffoldNamesTheLandingStages over a canned record of `schedule` phase spans under the landing span | A scaffold that reads the newest gate run names the operator's own gate instead. |
| LC32 | 35 | An absent record, an unreadable record, and a record with no landing span each write `unknown` under the timings heading at exit 0. | New TestRetroScaffoldReportsUnknownTimings through RetroCommand | A refusal blocks the draft the decision source requires. |
| LC33 | 36, 37 | The repair-attribution table carries one row per ticket file, in name order, with `unknown` in the rounds cell and the cause cell. | New TestRetroScaffoldListsTheTickets through RetroCommand | A missing ticket row or an invented cause fails the exact comparison. |
| LC34 | 38 | An absent tickets directory gives one row whose ticket cell reads `unknown`. | New TestRetroScaffoldListsTheTickets through RetroCommand | A scaffold that omits the table leaves the retired-folder case with no shape. |
| LC35 | 39 | `--scaffold` creates no file, and a later `--body` for the same slug still creates one. | New TestRetroScaffoldWritesNoFile through RetroCommand | A scaffold that writes makes the capture call refuse on an existing file. |
| LC36 | 40 | `bench retro <slug> --scaffold --body <markdown>` is a usage error at exit 2 with no file written. | New TestRetroRefusesBothForms through RetroCommand | A form that prefers one source silently discards the other. |
| LC37 | 40 | `bench retro <slug>` with neither form is a usage error at exit 2. | New TestRetroRefusesNeitherForm through RetroCommand | A form made optional writes an empty retrospective. |
| LC38 | 32 | The scaffold emits no improvement item, so the destination-marker check reports no diagnostic over the captured body. | New TestRetroScaffoldProposesNoImprovement through retros.ValidateImprovementMarkers | A scaffolded item without a `Feeds:` line reds the gate for every later capture. |
| LC39 | 41 | The lane's phase spans carry the landing span's trace identity. | New TestLandRecordsOneTraceForThePhases in the worktree package over a bound Bench home | A fresh trace at the authorization call leaves the phases unreachable from the landing. |
| LC40 | 42 | The reference, the final-check guidance, and the changelog state the two effects, the state-file form, and the scaffold form. | Review-owned Spec axis over the four guidance files | A missing paragraph leaves a cold session running the old eight-call close. |
| LC41 | 43 | No reset, tier-default, commit change-set, or test-table change enters the build. | Review-owned scope audit | An added grammar or a reset edit belongs to another spec's fence. |
| LC42 | 9 | The landing's own cleanup prints no fingerprint and asks for no second call. | New TestLandCleansTheFoldedSibling through landWith | A printed apply action reintroduces the two-call chunk this spec removes. |
| LC43 | 12 | A cleanup row whose worktree path is not line safe names its assignment pointer rather than its path. | New TestLandRetainsAnUnprovenSibling through landWith | A raw hostile path in the stderr evidence splits its own line. |

### Edge inventory

The canonical walk covers empty input, boundaries, errors, repetition, process boundaries, and hostile environments.
The attached profile is the shell CLI checklist in projects/benchkit.md.
Every behavior serves this repository and every repository that links the kit, except the refresh, which applies only where the destination declares Bench build inputs.

| Class | Concrete disposition | Rows |
|---|---|---|
| Absent versus empty | A destination with no build-input manifest is not applicable, an empty narrowed cleanup set is complete, and an empty state file writes an empty State. | LC2, LC9, LC28 |
| Already satisfied | A verified executable and an empty eligible set both report complete without a mutation. | LC3, LC9 |
| Errors | A build that leaves the executable unverified, a half-installed set, and a cleanup apply fault each report the failing effect. | LC4, LC8, LC14 |
| Ordering | The candidate build entry point runs after the release, and the cleanup never starts after a failed refresh. | LC6, LC13 |
| Repetition | A resume repeats no complete effect and publishes nothing a second time. | LC16, LC17 |
| Interruption | An interrupted refresh leaves the published commit, and the next resume finishes it from the tree. | LC19 |
| Process boundary | The invoked owner keeps one process identity through publication, release, and both effects. | LC7 |
| Output shape | The effects row precedes the landed record, and the landed record stays the last stdout line. | LC20 |
| File kinds | An absent, linked, special, or unreadable state file refuses before the read. | LC25 |
| Text sinks | An escape byte and a section heading in a state file each refuse before the document renders. | LC26, LC27 |
| Identity | A pin outside the ancestry, a pin that is not a commit, and an ambiguous prefix each answer their own fault. | LC22, LC23, LC24 |
| Paths | A retained worktree whose path is not line safe is named by its assignment pointer. | LC43 |
| Missing evidence | An absent record, an unreadable record, a record with no landing span, and an absent tickets directory all write unknown. | LC32, LC34 |
| Grammar | Both retrospective forms together, and neither form, are usage errors at exit 2. | LC36, LC37 |
| No mutation | The scaffold writes no file, and a refused handoff writes no document. | LC22, LC35 |

**Won't handle:** A landing that refreshes the broker of an installation outside its own destination is a separate capability. Every kit landing publishes into the checkout the broker was built from, so no caller loses the refresh.

**Won't handle:** A shim repair and a pre-push repair inside the landing belong to `bench doctor --fix`. The refresh publishes the executable and its manifest, and the doctor keeps the installation rows.

**Won't handle:** A repository-wide cleanup inside the landing is refused by the decision source. `bench worktree clean --landed` keeps every other landed assignment.

**Won't handle:** A retrospective scaffold that writes the file belongs to the capture call. `--body` keeps its exclusive create, so the delegate still returns a reviewed body.

**Won't handle:** A repair-cause the scaffold infers is refused by the decision source. The delegate reads the table and fills each cause.

**Won't handle:** A resume that recovers a missing terminal receipt is FT311's existing occurrence, not this spec's promise. The receipt refusal keeps its current message.

No new Git flag beyond the existing landed proof, no package-variable substitution, and no deletion of a tree file is required by this design.
If implementation needs one, its ticket must attach the corresponding hostile case before it claims its row.

## Ownership fences

These entries are the union of the ticket Writes.
They authorize implementation only after the reviewer approves this spec and its ticket graph.

- `internal/worktree/land_effects.go`
- `internal/worktree/land_effects_test.go`
- `internal/worktree/land_effects_cleanup_test.go`
- `internal/worktree/land.go`
- `internal/worktree/land_resume.go`
- `internal/worktree/land_freshness_test.go`
- `internal/worktree/land_resume_test.go`
- `internal/worktree/land_specless_test.go`
- `internal/worktree/land_journey_test.go`
- `internal/worktree/land_trace_test.go`
- `internal/worktree/clean_landed.go`
- `internal/worktree/clean_landed_apply_test.go`
- `internal/worktree/worktree.go`
- `internal/runbinary/runbinary.go`
- `internal/runbinary/runbinary_test.go`
- `internal/handoff/handoff.go`
- `internal/handoff/state_file.go`
- `internal/handoff/state_file_test.go`
- `internal/handoff/state_scan_test.go`
- `internal/otelrecord/reader.go`
- `internal/otelrecord/reader_test.go`
- `internal/roadmap/retro.go`
- `internal/roadmap/retro_scaffold.go`
- `internal/roadmap/retro_scaffold_test.go`
- `internal/roadmap/retro_test.go`
- `internal/retros/retros.go`
- `internal/retros/retros_test.go`
- `cmd/bench/main.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `.bench/BENCH-reference.md`
- `.agents/commands/bench-final-check.md`
- `AGENTS.md`
- `CHANGELOG.md`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `tests/canary/workflow-guidance-anchors/final-check-landed-worktree-sweep`
- `tests/canary/package-core-guard/unrouted-subcommand`
- `tests/canary/skills-index-command-adapters`
- `tests/canary/load-validity-metadata`
- `tests/canary/docs-currency-token-diet`
- `tests/canary/workflow-guidance-anchors`
- `reviews/ft311-landing-completion.md`
- `ROADMAP.md`
- `roadmap/FT311.md`
- `specs/ft311-diagnostics/spec.md`
- `specs/ft311-diagnostics/decisions`
- `specs/ft311-landing-completion/decisions`

The last five entries are the spec phase's own change, which the landed range carries.
They are the decision map's source and destination paths, the two roadmap references to the old path, and the implemented spec's own reference.
No ticket writes them.

The five command registries and the five canary entries are closure declarations from the preflight proposal.
The registries are named because the worktree, roadmap, and command packages are bound, and the canary entries pin the guidance and help surfaces.
The expected edit inside them is the two grammar suffixes.
One guidance fixture changes: ticket 5 moves `final-check-landed-worktree-sweep` onto the new needle bytes.

The coordinator owns the conditional review pickup, so no ticket names it.

The build cannot edit this spec, its acceptance rows, or its ticket graph without the existing spec-change authority.
Five fenced files are over the line budget: `cmd/bench/main.go`, `cmd/bench/command_registry_test.go`, `internal/worktree/worktree.go`, `internal/worktree/land_freshness_test.go`, and `internal/handoff/state_scan_test.go`.
`cmd/bench/main.go`, `cmd/bench/command_registry_test.go`, and `internal/worktree/worktree.go` change in place and gain no line.
`internal/worktree/land_freshness_test.go` and `internal/handoff/state_scan_test.go` shed their moved cases in the same commit that adds to them.

## Out of scope

These are planning estimates of file edits, not measured implementation costs.
Each future capability gets its own specification and one whole-project landing gate.

| Capability | Estimated price | Derivation |
|---|---|---|
| Recoverable reset | 10 edits, 1 gate runs | Worktree grammar, checkpoint preservation, plan validation, restore, and closure tests |
| Lower-tier trials | 6 edits, 1 gate runs | Task definitions, twelve matched runs' evidence projection, comparison, and reviewer decision record |
| Broker refresh for a foreign installation | 5 edits, 1 gate runs | An install-root resolver, its refusal, a manifest-directory operand, and two tests |
| Landing-effect journal and a status row | 6 edits, 1 gate runs | A record shape, its writer, a reader, a status row, and its tests |

The remaining order is recoverable reset, then trials.
FT120, FT254, FT258, and FT290 retain their neighboring subjects.
The existing merge lane, the composition, and the landing gate remain unchanged.

## Further notes

### Review round and dogfood runs

The review round ran over the uncommitted draft on base `c8256a7d5749c621f2db1eb60ad9d05116dd8a8f`, through the native agent surface at Terra/medium.
Iteration one returned 15 findings across the six rubric questions and the ticket quiz, with 5 blocking and 14 repair targets.
The tree refuted the phase-span finding.
Iteration two returned 1 blocking fence gap and 3 stale sentences, and the author folded them after the round.
No dogfood run applies, because the staged spec changes no executable.

### Decision provenance

The whole ready map and its topic folder moved together into this spec's decisions directory, from the diagnostics spec that carried them.
Three references to the old path moved with it in the same change.
They are the roadmap sequence line, the FT311 detail file, and the diagnostics spec's source line.
The diagnostics spec's decision-provenance paragraph named a return to the top-level decisions folder.
The FT311 detail file superseded that instruction, so this spec repaired the paragraph and flags the repair below.
The remaining FT311 spec consumes this same source in place, and the last of them returns it or retires it.

The three structured map sources were reread on 2026-09-09.
The FT311 detail file's landing-completion group is the fold list this spec settles, and its residual list is not.
The research report's Q1 landing chunk and its Q3 implementation facts were checked against the current landing, cleanup, broker, handoff, and retrospective owners.
The report's fact T5, that the landing already prunes unclaimed landed branches, holds and is narrower than this fold.
The kit profile's tier binding, hostile-input checklist, and lane table were reread, and none changes here.

The reviewer closed two decisions on 2026-09-09.
The landing performs the broker refresh.
This spec re-cuts the two stable-owner proofs, and it deletes neither.
Each re-cut proof reads its marker at the publication boundary and names the destination branch its script observed.
A failed refresh or a failed cleanup after publication exits 3 with the `effects` row and the `next=` resume.
Both decisions stay closed.

### Source-sentence-to-row accounting

References below identify resolved tickets in the single decision source and the FT311 detail file's landing-completion group.
Each table entry paraphrases its assigned clause without creating another decision source.

| Source clause | Disposition |
|---|---|
| FT311: `bench worktree land` republishes the broker and applies the landed-sibling plan | LC1 through LC15 |
| FT311: `bench handoff --state-file <file>` reads the state from a file | LC21 through LC29 |
| FT311: `bench retro <slug>` scaffolds the retro with the landing's gate timings and the repair table | LC30 through LC39 |
| FT311: the fold-lane behavior is a fact, not an open decision | LC41 |
| FT311: the two diagnostics residuals wait for a reviewer decision | Out of scope, and no row |
| Ticket 1: landing refreshes the broker and cleans eligible siblings | LC1, LC9 |
| Ticket 1: handoff accepts state from a file | LC21 |
| Ticket 1: retrospective preparation derives available timing and repair facts | LC31, LC33 |
| Ticket 2: each effect has a recorded result and an independent resume path | LC1, LC16, LC19 |
| Ticket 2: publication succeeded and refresh remains incomplete | LC4, LC5 |
| Ticket 2: a retry executes only unfinished effects | LC16 |
| Ticket 2: cleanup preserves every worktree whose eligibility cannot be proved | LC11, LC12 |
| Ticket 3: the CLI derives revisions, timings, and evidence handles from canonical sources | LC31, LC33, LC39 |
| Ticket 3: delegates perform work that requires interpretation | LC33, and the cause cells stay unknown |
| Ticket 7: cleanup covers only assignments this landing includes | LC9, LC10 |
| Ticket 7: one landing does no repository-wide maintenance | LC10, LC15 |
| Ticket 8: mutations stop after a post-publication failure | LC13 |
| Ticket 8: the result reports each effect as complete, failed, or pending | LC1, LC4, LC13, LC14 |
| Ticket 8: resume verifies publication and continues without publishing again | LC17, LC18 |
| Ticket 8: the worktrees remain while the installation requires repair | LC13 |
| Ticket 12: drafts mark missing facts unknown | LC32, LC33, LC34 |
| Ticket 12: missing timing data does not block a draft | LC32 |
| Ticket 12: a missing or conflicting revision identity blocks a dependent continuation instruction | LC22, LC23, LC24 |
| Ticket 12: drafts invent no attribution | LC33 |
| Ticket 16: a derived fence addition needs approval only when it grants new authority | The fence list above, and the ticket extension note |
| Tickets 10 and 18: five specs, one source, fixed order | LC41 and Out of scope |

### Reader sweep and enforcement reads

The repository sweep used hidden-file coverage over Go, Markdown, fixtures, scripts, and workflow files.
The `landed{` record has five readers, all inside this repository's own tests, and the effects row precedes it rather than joining it.
The moved decision map has four references outside the moved folder, and the change repairs three of them.

| Reader or owner | Read evidence and disposition |
|---|---|
| worktree landAttributed and landWith | The follow-up step order and the incomplete render. LC1 through LC14 add two steps after the release. |
| worktree resumeLandWith and resumeAssignment | The publication proofs, the active and released branches, and the terminal receipt. LC16 through LC19 extend the released branch. |
| worktree landedComplete, landedIncomplete, and landingResumeNext | The terminal record and the resume command. Neither changes; the new step names compose the existing policy. |
| landingpolicy Terminal, Residue, Publication, and ResumeMarker | The typed policies. None changes, because the step name is an input. |
| worktree brokerChangeNotice and brokerInstallStep | The build-input predicate and the install-step sentence. The refresh reuses the predicate and the notice stays. |
| freshness DeclaresBuildInputs, BuildInputs, Verify, and PublishedExecutable | The applicability and completion predicates. The refresh composes them and states no path of its own. |
| freshness Publish and its transaction | The executable, seal, and manifest install as one outcome, seal last. LC8 rests on that order. |
| runbinary Build and canonicalBuild | The build seam and its manifest-directory flag. The refresh adds the subject-mode form beside it. |
| adopt doctorFix, publishBrokerManifest, and WriteBrokerManifest | The installation repair and its wrapper-relative manifest. Neither changes, and the landing calls neither. |
| adopt evalBrokerManifestRow | The five broker predicates the doctor row applies. It stays the operator's health report. |
| bin/bench.sh land_route, land_read_manifest, and land_rebuild_broker | The wrapper's authentication of the broker before execution. It does not change. |
| worktree planLandedSet, selectLandedCleanupRow, and applyLandedSet | The landed selector, its fingerprint, and its per-row re-plan. LC9 through LC14 add one scope value. |
| worktree retainForLandedPreservation and lifecyclepolicy AutomaticPreservation | The eligibility verdict for an automatic removal. It does not change. |
| worktree cleanCommand `--landed` branch | The explicit plan and apply route. LC15 pins its rows and its fingerprint. |
| git LandedInDefault | The landed proof by ancestry and by content. The scope value calls it a second time against the destination base. |
| handoff Command, writeSection, and scanState | The section rewrite and the three State faults. LC21 through LC29 change the State input only. |
| handoffdoc Update, Parse, Render, and UnfencedLines | The document store and its section split. LC27 keeps the parse valid. |
| handoff validate and toon.Representable | The unrepresentable-field refusal. LC26 widens the same predicate to the State body. |
| status HandoffFile and conformance checkHandoffShape | The document path and the single-source Shape text. Neither changes. |
| roadmap RetroCommand and inboxRoot | The exclusive create and the primary-local destination. LC35 keeps both for `--body`. |
| retros Parse, requiredHeadings, Recommendations, and ValidateImprovementMarkers | The heading list and the destination-marker rule. LC30 and LC38 render from the same list. |
| bench-final-check retro template and its anchor needles | The guidance source for the headings and the cause vocabulary. Ticket 5 keeps every needle. |
| otelrecord Begin, BeginIn, Encode, and Writer | The span protocol and the record file. LC39 threads the landing context and LC31 reads the same shape. |
| gate beginLaneSpan, beginGateSpan, and schedule | The lane span, the gate span, and the phase spans below them. Neither changes. |
| tickets Enumerate and Entry | The ticket enumeration seam. LC33 reads it and adds no second walk. |
| cmd/bench registry rows for handoff and retro | The two help suffixes. Tickets 3 and 4 edit them in place. |
| help_inventory_test and command_registry_test | The pinned help rows and the AXI matrix. Tickets 3 and 4 update both pins. |
| axi_query_registry_test and subcommand_routing_table_test | Membership and route ownership only. No new command. |
| CONTEXT.md `landing effect` and `eligibility verdict` | Both terms already name this behavior. Neither changes. |
| projects/benchkit.md lane table and tier binding | The lane argv and the model binding. Neither changes. |

### Pre-review proof checklist

- Cited symbols: every symbol in the reader sweep table above was read in the current tree on 2026-09-09. The review round found two cited test functions that do not exist, and LC29 and LC37 now name new tests instead.
- Import edges: worktree already imports freshness, runbinary, landing, gate, and git. runbinary already imports freshness. roadmap already imports retros. The otelrecord reader adds no import to its consumers beyond roadmap, and otelrecord imports neither worktree nor gate, which this design keeps.
- Source-row clauses and occurrences: the accounting table enumerates the three fold sentences of the FT311 detail file and the nine binding decision tickets. No source clause is rewritten by this spec.
- Promised field labels: `effects`, `effect`, `result`, `complete`, `failed`, `pending`, `not-applicable`, `refresh`, `cleanup`, `incomplete:refresh`, `incomplete:cleanup`, `--state-file`, `--scaffold`, `unknown`.
- Changed-function callers: landAttributed and resumeLandWith are called by their own entry points only. planLandedSet, applyLandedSet, and selectLandedCleanupRow are called by the clean route, the re-plan, and their tests. handoff writeSection is called by Command only. RetroCommand is called by the registry and its tests. retros.Parse is called by RetroCommand and ValidateImprovementMarkers. runbinary.Build is called by the worktree build verb and the landing seam set.
- Copy survival: none. No copy is replaced by a new owner. The scaffold renders from the parser's heading list rather than a second list, and the coordinator's swapped-heading probe is the catch.
- Git flags: the design adds no Git flag. The landed proof calls `merge-base --is-ancestor`, `rev-list --merges`, and `cherry`, which the existing owner already runs, with the destination base as the second subject.

### Flagged additions and engineering choices

No new product promise beyond the reviewed map is intended.
The `effects` row and its four result words, the two new step names, the narrowed cleanup scope value, and the `unknown` cells are engineering choices.

Story 41 and row LC39 are the enabling engineering choice for the fold's own words, the landing's gate timings.
The gate scheduler `schedule` already opens one phase span for the full gate and for the lane.
`withGateSpanEnv` already carries the record root and the trace parent to the phases child.
The landing alone passes a background context, so no reader can group those spans by landing.
The threading is therefore the smallest change that makes the promise readable, not a promise of its own.

The stable-owner boundary was the one contestable call, and the reviewer closed it on 2026-09-09.
The decision source folds the broker refresh into the landing, and the shipped design proves that the landing never executes candidate build code.
This spec keeps the invariant that matters: no candidate code reaches the publication decision or the published bytes.
It re-cuts the four proofs, so each keeps its own assertion beside the refresh result.
The refuse-the-fold option is closed, and stories 1 through 8 stay in this spec.

The reviewer also closed the exit posture on 2026-09-09.
A failed refresh or a failed cleanup exits 3, like every other incomplete follow-up step.

Three smaller calls are flagged.
The refresh applies on the build-input predicate alone, rather than on a second kit-checkout read, because one predicate cannot disagree with itself.
The scaffold prints to stdout and writes no file, because `bench retro --body` owns the exclusive create.
A scaffold that wrote the file would make that create refuse.
The repaired paragraph inside the implemented diagnostics spec is a prose repair, not a change to that spec's rows or fences.

### Implementation ticket approval table

| Ticket | Blocked by | Delivered outcome |
|---|---|---|
| 1. Refresh the promotion broker after publication | none | The refresh effect, a one-cell effects row, its resume, and the four re-cut fixtures |
| 2. Clean the landing's eligible siblings | 1.md | The two-cell effects row, the narrowed landed scope, the cleanup effect, and its resume |
| 3. Read the handoff State from a file | none | `bench handoff --state-file <file>` with the existing scan and the hostile-file refusals |
| 4. Scaffold the retrospective from the landing's record | 2.md, 3.md | `bench retro <slug> --scaffold`, the seam-record reader, and the landing's single trace |
| 5. Fold the guidance and the changelog | 3.md, 4.md | The reference, the phase guidance, and the changelog state the landed behavior |

Tickets 1 and 3 form the first frontier and run in parallel, one worktree each.
Ticket 2 follows ticket 1 because both write the effects owner and the two landing entry points.
Ticket 2 also widens the effects row that ticket 1 prints with one cell.
Ticket 4 follows tickets 2 and 3, because it writes `internal/worktree/land.go` with the first and the command help with the second.

Ticket 2 carries the cleanup invariant for the worktree package, and ticket 4 carries the record invariant for the otelrecord package.
Ticket 5 is the last writer of the shared guidance files.
The review pickup is created only for actionable findings.

### Approval surface

| Subject | Disposition requested |
|---|---|
| Stories and lines | Approve the six outcome groups and their bound model efforts. |
| Seams | Approve the two landing seams, the cleanup selector, and the two public command seams. |
| Acceptance and edges | Approve LC1 through LC43 and the six explicit exclusions. |
| Ownership fences | Approve the exact union above for implementation. |
| Scope and tickets | Approve the five-ticket graph within the third FT311 capability. |
| Stable-owner boundary | Closed on 2026-09-09; note the four re-cut proofs before ticket 1 dispatches. |

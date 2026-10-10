# Build directions

Status: staged

Roadmap: FT410

Decision source: named reviewed artifact .bench/commitment.json outcome build-directions, approved 2026-10-10 in commit 44ac5793.

The outcome cites roadmap/FT410.md and the reviewer directions at capture/ft217-build-timing-2026-10-09.log:32.
Two late decisions on 2026-10-10 close its open forks.
A combined repair charge for findings on several tickets waits for FT405.
The build log is a git-ignored capture file.

Verification log: 4 iteration(s) to accept — the spec round took 3 iterations, for event-order defects and then a fold-point contradiction; the ticket round took 1.

## Problem

The FT217 build paid for steps that the guidance asked for and that gave no new evidence.
The timing log records each cost (capture/ft217-build-timing-2026-10-09.log).
Four `main` folds each ran the whole gate, and three of them fell inside a chunk.
One fold composed a light-path landing that the build itself dispatched.

Two rules force a fold inside a chunk today.
The bare post-ticket preflight at .agents/commands/bench-implement-spec.md:50 grades `base-current` red when `main` moves.
Its next action is `bench worktree merge --from main <target>` (internal/preflight/decision.go:268-288).
The stale-base rule at .agents/skills/bench-craft-delegate/SKILL.md:67-68 folds before each new charge.

The operating guide lets a build ship a light-path fix at once (.bench/BENCH.md:153).
The orchestrator records each chunk entry (.agents/commands/bench-implement-spec.md:58).
That rule cost one freeze and record round trip for each chunk until the reviewer moved it to the author.
The implement phase runs a chunk checkpoint gate for every chunk, the last chunk included (:58).
The complete checkpoint then grades the last chunk again (internal/reviewrecord/coverage.go:44-144).

The first author ran 54 serial probes after its code was done (log:5).
The witness cap lives only in a known-issues workaround (docs/known-issues.md:12).
A confirming round raised a new finding and started a second repair cycle (log:15).
No guidance names a build log, so the saved steps live in an ad hoc untracked file.

## Solution

The guidance states nine reviewer directions, and the tree keeps one owner for each rule.
A build folds `main` only at three sanctioned points, never between two tickets of a chunk.
A build dispatches no light-path fix while it runs.
The ticket author freezes the chunk, records the chunk entry, and records its own verification.
The chunk checkpoint gate runs only before a successor chunk.
The last chunk goes straight to final reconciliation and closes through the complete checkpoint alone.

The author charge states the witness cap and the batched `bench test` call.
It also says that moved or merged code keeps its error handling and validation.
The axes of a round run in parallel, and a one-ticket round sends all findings in one combined repair charge.
A confirming round examines only the findings that its charge names.
The orchestrator keeps one git-ignored build log and records each reviewer direction that saves a step.
FT379 and FT369 stay open, and the landing retires FT410 only.

## User stories

Line: opus / high.
Implementation-line reason: leverage, because the edits are guidance that steers every later build.
The spec is exact, the seams are known anchor rows, and anchor and conformance tests cover them.
Harder chunks: BD-C4, because it holds the one Go production change and it touches the landing predicate.

### Fold timing

1. As the orchestrator, I want no `main` fold between two tickets of a chunk, so that a moved `main` adds no gate run there.
2. As the orchestrator, I want the three sanctioned fold points named, so that I know where a fold can go.
3. As the orchestrator, I want the post-ticket preflight inside a chunk to pass an explicit range, so that a moved `main` routes no fold.
4. As a build coordinator, I want the craft-delegate stale-base fold scoped out of a build, so that two rules agree.
5. As a maintainer, I want the anchored fold sentences of both phases to stay, so that the chunk chain keeps its pins.

### Light paths during a build

6. As the orchestrator, I want implement-spec to forbid a light-path dispatch while a build runs, so that no side landing moves `main`.
7. As the orchestrator, I want a fix that the build needs to land as the build's own commit, so that the fix still ships.
8. As a reader of the operating guide, I want its light-path rule to stop granting an in-build fix, so that the guide keeps one source.
9. As a reader of the operating guide, I want its light-path rule to point to implement-spec for a build, so that I find the owner.
10. As a reader of ADR 0023, I want its light-path consequence to agree with the build rule, so that the ADR stays current.

### Chunk records

11. As a ticket author, I want to freeze the chunk and record its entry after the last ticket, so that a round trip goes away.
12. As a ticket author, I want to record my own verification entries at the chunk source, so that the recorder rule stays with me.
13. As a maintainer, I want the retired orchestrator-records sentence forbidden, so that it cannot return.

### Checkpoints and the last chunk

14. As the orchestrator, I want the chunk checkpoint gate only before a successor chunk, so that the last chunk repeats no gate run.
15. As the orchestrator, I want the last chunk to go straight to final reconciliation, so that one checkpoint grades it.
16. As the orchestrator, I want the last chunk to close through the complete checkpoint of `/bench-final-check`, so that one phase owns `--complete`.
17. As the orchestrator, I want the last chunk to keep its record, verification, and review entries, so that the complete checkpoint accepts it.
18. As the orchestrator, I want a handoff refresh at the final reconciliation, so that the last chunk keeps a refresh point.
19. As a README reader, I want the flow diagram to show the chunk checkpoint on the successor edge only, so that it matches the rule.
20. As a reader of ADR 0023, I want its handoff consequence to name the final reconciliation, so that the ADR stays current.
21. As a maintainer, I want the retired every-chunk checkpoint sentence forbidden, so that it cannot return.

### The author charge

22. As a ticket author, I want the charge to keep the guards of moved or merged code, so that a move drops no validation.
23. As a ticket author, I want the charge to define a witness, so that the cap counts one thing.
24. As a ticket author, I want the charge to cap the witnesses by rows and new branches, so that I run no long probe loop.
25. As a ticket author, I want to meet the cap with fewer independent expectations, so that no demonstrated red goes missing.
26. As a ticket author, I want one `bench test` call for several packages, so that I pay one build for the floor checks.
27. As a ticket author, I want the charge to keep its reach list of callers, tests, fixtures, and exports, so that the self-check stays whole.
28. As a maintainer, I want the known-issues row to point to implement-spec for the cap, so that the cap has one source.
29. As a maintainer, I want the one-red rule for each independent expectation to stay, so that the cap keeps the one-source rule.

### Review rounds

30. As the review coordinator, I want the parallel dispatch of the axes pinned, so that the default cannot drift.
31. As the review coordinator, I want one combined repair charge for a one-ticket round, so that one session fixes every finding.
32. As the review coordinator, I want one repair session for each affected ticket in a multi-ticket round, so that the record rule holds.
33. As a confirming-round axis, I want to examine only the findings that my charge names, so that no new repair cycle starts.
34. As a confirming-round axis, I want the round to stay three axes at the final tip, so that the checkpoint finds three entries.
35. As a maintainer, I want the Standards first-look sentence removed and forbidden, so that no sentence widens a confirming round.
36. As a maintainer, I want the delegation-discipline re-review sentence to agree with the named-findings scope, so that two files agree.

### The build log

37. As the orchestrator, I want one build log in the primary checkout's capture folder, so that it survives the worktree's release.
38. As the orchestrator, I want one log line for each reviewer direction that saves a step, so that the log names each saving.
39. As the retro author, I want implement-spec to name me as the log's reader, so that the saved steps reach the retro.
40. As a kit maintainer, I want the kit `.gitignore` to ignore `capture/build-log/`, so that the log never dirties the tree.
41. As a linked-repository maintainer, I want `bench init` to add the build-log ignore line, so that my repository ignores the log.
42. As the landing, I want `landing.LocalCapturePath` to accept the build-log directory and its clean paths, so that the init agreement holds.
43. As the landing, I want an unsafe or collapsed build-log path refused, so that the allowance stays narrow.
44. As the landing, I want a build-log write to leave the checkout fingerprint unchanged, so that a log write during a landing refuses no landing.

### Open rows

45. As the reviewer, I want FT379 and FT369 unedited and open, so that their undelivered parts keep their owners.
46. As the reviewer, I want the landing to retire FT410 only, so that the roadmap closes exactly the delivered row.

## Implementation decisions

### One owner for each rule

Implement-spec owns the build's fold timing, its light-path hold, its record order, its checkpoint timing, the author charge, and the build log.
The operating guide keeps the light-path dispatch rule outside a build and points to implement-spec for the build rule.
The review phase owns the parallel dispatch, the combined repair charge, and the confirming-round scope.
`/bench-final-check` keeps the complete checkpoint command, and implement-spec names that phase instead of the flag.
The operating guide keeps its repair rule, one fresh repair session for each affected ticket, unchanged.
Each edit replaces a sentence in place or adds a sentence, so that no second statement of a rule survives.

### Fold timing

The Land section of implement-spec takes the five sentences below in a new paragraph.

> A build folds `main` only before the first chunk, at the start of the last chunk, or before the completion landing. The fold at the start of the last chunk comes before its first ticket dispatch. It never folds `main` between two tickets of a chunk. A fold before the completion landing follows `/bench-review-implementation` step 1, and the last chunk refreezes. Inside a chunk, the post-ticket preflight passes the chunk's `--base` and the current `--source-tip`, so a moved `main` does not route a fold.

A fold at the start of the last chunk joins that chunk's review delta, so the pinned sentence at implement-spec:56 holds.
The last ticket's author then freezes the chunk after that fold, so no refreeze and no records-only resumption follow.
A `main` move after that start takes the fold before the completion landing, and the last chunk then refreezes.
That refreeze follows the named record action of the checkpoint evidence, which implement-spec:60 already routes.
The existing pinned sentences at implement-spec:56 and review-implementation:98 stay byte-exact.

The explicit range makes `base-current` green (internal/preflight/decision.go:262-266).
The docs/known-issues.md:14 row stays, because it retires only on FT342.

The craft-delegate stale-base rule keeps its two sentences byte-exact, because registry_data.go:69 and fixture_bite_test.go:350 quote them.
Line 68 of that skill takes the sentence below at its end.

> During a build, `/bench-implement-spec` owns the fold timing, and a ticket charge starts without this fold.

The craft-delegate SKILL.md line count therefore does not change.

### Light paths during a build

Implement-spec takes the two sentences below in a new paragraph of its Build section.

> A build dispatches no light-path fix while it runs. A fix that the build's own outcome needs lands as the build's own commit under `.bench/BENCH.md`'s fix-don't-park rule.

The operating guide's light-path sentence at .bench/BENCH.md:153 keeps its `**Delegate each light-path fix.**` label byte-exact, because internal/conformance/commitment_guidance_test.go:35 mutates from it.
The sentence after the label becomes the sentence below.

> Outside a build, a `bench learning` entry's light-path fix ships at once only when the active outcome needs it; `/bench-implement-spec` states the build rule.

Every other sentence of that paragraph stays, so the anchored neighbours at registry_retained_workflow.go:243-244 and registry_commitment.go:24 hold.

ADR 0023 line 19 ends its light-path consequence with "and so does a phase outside a build that needs a learning's fix."
It then adds "A build dispatches no light-path fix while it runs."
The `/bench-final-check` sentence at final-check:52 stays, because it governs a landing from another session.

### Chunk records and checkpoints

The Land paragraph at implement-spec:58 replaces its second sentence with the sentence below.

> After the last ticket of a chunk, that ticket's author freezes the chunk and records the chunk entry with `bench record chunk`.

The anchor row at registry_chunk_chain.go:20 and its mirror at registry_chunk_chain_test.go:16 change in place to that needle.
The next sentence, "Each ticket author then writes its verification entries at that chunk source with `bench record verification`.", stays.
The last sentence becomes the two sentences below.

> Commit the artifact on its ordinary lane. Before a successor chunk starts, run `bench gate --checkpoint specs/<slug>/spec.md --chunk <id>`.

The `bench record chunk` form has no recorder field, so no code change follows (internal/reviewrecord/recordcmd/chunk.go:10-11).

The Build sentence at implement-spec:50 that starts "After the last ticket in a chunk, freeze the chunk delta" stays byte-exact.
docs_workflow_helpers_test.go:109 pins its lowercase form.
That paragraph takes the sentence below after it, and it then holds six sentences.

> The last ticket's author performs that freeze, as the Land section states.

The paragraph at implement-spec:62 keeps its first sentence byte-exact, because docs_workflow_helpers_test.go:111 pins its lowercase form.
It takes the two sentences below after the first sentence.

> The last chunk runs no chunk checkpoint gate, because the complete checkpoint that `/bench-final-check` owns grades every chunk again. The last chunk keeps its chunk entry, its verification entries, and its three review entries.

The complete checkpoint needs no earlier chunk-checkpoint run, and it grades every recorded chunk (internal/reviewrecord/coverage.go:44-144; internal/gate/complete_checkpoint.go:19-43).

The handoff clause at implement-spec:31 becomes "not code, and it refreshes `bench handoff` at each chunk checkpoint and at the final reconciliation."
The anchor row at registry_ft311_preparation.go:43 changes in place to that needle, so no counted family changes.
ADR 0023 line 15 takes the same change.
Its new sentence is "It refreshes the session handoff at each chunk checkpoint and at the final reconciliation."
Its reconciliation clause stays as a separate sentence.
The README chunk flow at README.md:81-82 moves the label: `blockers -->|"none"| more{"More approved chunks?"}` and `more -->|"yes; chunk checkpoint"| chunk`.

### The author charge

The author-charge paragraph at implement-spec:25 takes the sentence below after its reach-list sentence.

> Code that the author moves or merges keeps its error handling and validation.

That paragraph then holds six sentences, the prose bound.

A new paragraph after it states the witness cap in the four sentences below.

> A witness is one demonstrated red of one independent test expectation. The charge caps the witnesses at one for each coverage row plus one for each production branch that the ticket adds. The author meets the cap with fewer independent expectations, such as one table-driven test, and skips no witness. The author runs the floor checks and witnesses of several packages in one `bench test` call that repeats `--package`.

The `bench test` help states the repeated `--package` grammar (internal/testreport/selection_facts.go:89).
The `bench probe` help owns its one-package limit, so the charge does not restate it.

The delegation-discipline rule "The second is one recorded red for each independent test expectation." stays byte-exact (registry_charge_binding.go:18).
The cap bounds the count of independent expectations, so the two rules agree.
The docs/known-issues.md:12 workaround cell becomes the two sentences below.

> Apply the witness cap and the batched `bench test` call that `/bench-implement-spec` states. Run the probes after the floor checks, so that an invalid mutant does not block the commit.

That row keeps its symptom cell and its FT168 retire cell.

### Review rounds

The review phase's step 4 already states "In the default mode, spawn the axes in parallel sub-agents." (review-implementation:118-119).
A new Require row pins that sentence in its step.
The "Review modes" section takes the three sentences below as their own paragraph after its repair-coverage paragraph.
The paragraph at review-implementation:46-49 keeps five sentences, so the prose bound holds.

> When every accepted finding of a round falls on one ticket, one combined repair charge carries them all. That charge holds the findings of every axis and goes to that ticket's repair session. When the accepted findings fall on several tickets, each affected ticket keeps its own repair session under `.bench/BENCH.md`'s repair rule.

The one-session-one-ticket rule stays (internal/reviewrecord/delegated.go:153-154), and FT405 owns the multi-ticket case.

The confirming-round sentence at review-implementation:48 becomes the two sentences below.

> A confirming round reads only the repair delta and examines only the findings that its charge names. Its charge names those findings and the folds to confirm.

The sentence at :49, "That charge directs the Standards axis to look first for duplicated knowledge in the repair delta.", goes away.
The repair author already runs a duplicated-facts sweep of its own delta (delegation-discipline.md:96), so the removal leaves that check in place.
The round stays three axes at the final tip (registry_chunk_chain.go:26), because the checkpoint needs three review entries (internal/reviewrecord/record.go:111-119).
The anchor row at registry_ft311_review_dispatch.go:22 changes in place to the first new sentence.
Its diagnostic drops the clause "or drop its named folds", because BD59 owns that clause.
A new row in registry_chunk_chain.go pins the second sentence, so the folds clause keeps its pin.

The delegation-discipline bullet at :216-217 becomes "- A re-review charge names the findings and the folds to confirm, and the round examines only those findings."
The repeated-review sentence at review-implementation:46 stays, because it governs a review of a later delta, not a repair close.

### The build log

Implement-spec takes the three sentences below in a new paragraph of its Land section.

> The orchestrator keeps one build log at the primary checkout's `capture/build-log/<slug>.log`. For each reviewer direction that saves a step, it appends one line that names the direction and the step that it removed. The retro author reads that log.

The log has no schema and no verb, and no gate reads its content.

The kit `.gitignore` takes the line `capture/build-log/` in its capture block.
The `landing` package owns the directory name as one exported constant, `landing.BuildLogDir`, with the value `capture/build-log/`.
`LocalCapturePath` and that constant move to a new file, internal/landing/local_capture.go, so landing.go gains headroom under its 400-line limit.

`landing.LocalCapturePath` accepts that directory, its slash-less form, and each clean path beneath it.
It refuses a control byte, a `..` segment, a doubled slash, and a sibling name that only shares the prefix.
`CheckoutFingerprint` runs git status with `--untracked-files=all --ignored` (internal/landing/landing.go:326), and the exact git spelling of an ignored directory there is not observed.
The predicate accepts the directory spellings and each clean path beneath them, so either spelling passes.
The `RuntimeIgnoredPath` guards at internal/landing/landing.go:340-354 are the precedent for that shape.

`internal/adopt/init.go` adds the constant to `localCaptureIgnoreLines`, so `bench init` writes the ignore line.
The new import edge from `internal/adopt` to `internal/landing` makes no cycle: `go list -deps ./internal/landing` lists no `internal/adopt`.
A repository that ran `bench init` before this landing gains the line at its next `bench init`.

The landing reads ignored paths through this predicate in two places.
The destination fingerprint filter drops each accepted entry (internal/landing/landing.go:378).
A new log file created during a landing would otherwise change that fingerprint, and the landing would refuse (landing.go:256-257).
The source cleanup plan reads the same predicate (internal/worktree/build_outputs.go:81).
The init agreement test reds once init writes a line that the predicate refuses (internal/adopt/adopt_test.go:579-585).

A conformance test beside TestLocalCaptureFilesAreNotTracked reads the kit `.gitignore` and requires a line equal to `landing.BuildLogDir`.
The constant therefore stays the one source, and internal/conformance/tier_live_tree_test.go classifies the new live-tree test.

### Budgets and shared files

Implement-spec, the operating guide, and the craft-delegate skill sit at their line budgets (projects/benchkit.md:490, :491, :496).
A ticket that adds net lines to one of them moves that budget row in the same ticket, as commit 47eed422 did.
The edits to the operating guide and the craft-delegate skill add no line.

Implement-spec gains four paragraphs of two lines each, so its budget moves from 85 to 93 in four steps.
Ticket 01 moves it to 87, ticket 02 to 89, ticket 04 to 91, and ticket 06 to 93.
Ticket 03 edits sentences in place and adds no line.

`internal/landing/landing_test.go` (384 lines) and `internal/landing/landing.go` (387 lines) sit near the 400-line limit, so ticket 06 adds two new files there.
`internal/anchors/registry_data.go`, `internal/anchors/registry_data_test.go`, and `internal/adopt/adopt_test.go` exceed the structure limit, so they cannot grow.
New anchor rows go into `internal/anchors/registry_chunk_chain.go`, `internal/anchors/registry_ticket_passes.go`, or `internal/anchors/registry_commitment.go`, and no file joins `internal/anchors/`.
Each ticket adds one `CHANGELOG.md` entry.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| BD-C1 / 01-fold-only-at-chunk-boundaries.md, 02-hold-light-paths-during-a-build.md, 03-close-the-last-chunk-through-the-complete-checkpoint.md | The build flow folds at chunk boundaries, holds light paths, gives the chunk record to the author, and runs the chunk gate only before a successor | BD1–BD21, BD50–BD52, BD55, BD57, BD58, BD60, BD61 | TestChunkChainAnchors, TestCommitmentGuidance, TestReviewConvergenceContractCurrentDocs, TestWorkflowCadenceAnchorsRejectDeletionAndSwap, TestEveryRetainedFixtureBitesThroughRegisteredOwner, TestGuidanceProseBudgetsHoldOnTheLiveTree | no |
| BD-C2 / 04-cap-witnesses-and-keep-moved-guards.md | The author charge states the moved-code guard, the witness definition and cap, and the batched `bench test` call | BD22–BD30 | TestWriteSpecChunkedSlicing, TestChargeBindingAnchors, TestGuidanceProseBudgetsHoldOnTheLiveTree | no |
| BD-C3 / 05-scope-review-rounds.md | The review phase pins parallel axes, sends a one-ticket round in one combined repair charge, and scopes a confirming round to its named findings | BD31–BD38, BD54, BD59 | TestChunkChainAnchors, TestPreparedReviewGuidance, TestEveryRetainedFixtureBitesThroughRegisteredOwner | no |
| BD-C4 / 06-keep-the-build-log.md | The orchestrator keeps a git-ignored build log, and the landing and `bench init` treat it as local capture | BD39–BD49, BD53 | TestChunkChainAnchors, TestLocalCapturePathRecognizesTheBuildLog, TestCheckoutFingerprintIgnoresTheBuildLog, TestKitGitignoreIgnoresTheBuildLog, TestInitIgnoresLocalCaptureFiles, TestLocalCaptureIgnoreLinesMatchTheLandingPredicate | yes |

The ticket basenames are planned names, and the slicer may rename them.
Tickets 01, 02, 03, 04, and 06 all edit implement-spec, so they run serially in that order.
Each edit to implement-spec, the craft-delegate skill, or delegation-discipline changes the charge evidence of later tickets (internal/preflight/charge_pack.go:34-42).
The orchestrator therefore regenerates each later charge after such a commit.
This build follows the guidance on `main` at its build entry.
The new rules govern each build that starts after this landing.

## Testing decisions

- A good test reads the live guidance file and fails when the exact rule sentence is absent or a retired sentence returns.
- The anchor registry is the gate seam for each guidance rule, and the `docs-currency-workflow` check evaluates every row.
- Each new or rewritten row takes an independent mirror rule in an anchor harness test.
- registry_ft311_preparation.go and registry_ft311_review_dispatch.go have no harness test, so the BD35 and BD50 mirrors go to TestChunkChainAnchors.
- A removal of a needle then turns that harness red.
- A retired grant sentence takes its own Forbid row, one row for each sentence.
- The prose budget check and the canary check are regression seams, and each edit keeps every canary `old` string byte-exact.
- The build-log predicate takes a table test in a new landing test file, modelled on TestRuntimeIgnoredPathRejectsUnsafeOrCollapsedNames.
- A fingerprint test in that file writes an ignored build-log file into a real checkout and compares `CheckoutFingerprint` before and after.

### Seam diagram

    trigger: the gate's docs-currency-workflow check, and go test on each harness
        │
        ▼
    guidance file bytes  ──▶  [ anchors.EvaluateGroup over the registry rows ]  ──▶  diagnostics
                      ◀ tests attach here: each anchor harness removes or restores one needle and expects its diagnostic

    trigger: CheckoutFingerprint over a checkout that holds an ignored capture/build-log/<slug>.log
        │
        ▼
    git status --ignored entries  ──▶  [ landing.LocalCapturePath ]  ──▶  filtered fingerprint status
                      ◀ tests attach here: the predicate table test and the fingerprint test

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| BD1 | 1 | Implement-spec Land holds "It never folds `main` between two tickets of a chunk." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A missing no-fold sentence raises the new Require diagnostic. |
| BD2 | 2 | Implement-spec Land holds "A build folds `main` only before the first chunk, at the start of the last chunk, or before the completion landing." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A missing fold-point sentence raises the new Require diagnostic. |
| BD60 | 2 | Implement-spec Land holds "The fold at the start of the last chunk comes before its first ticket dispatch." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A missing dispatch-order sentence raises the new Require diagnostic. |
| BD61 | 2 | Implement-spec Land holds "A fold before the completion landing follows `/bench-review-implementation` step 1, and the last chunk refreezes." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A missing late-fold sentence raises the new Require diagnostic. |
| BD3 | 3 | Implement-spec Land holds "Inside a chunk, the post-ticket preflight passes the chunk's `--base` and the current `--source-tip`, so a moved `main` does not route a fold." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A missing explicit-range sentence raises the new Require diagnostic. |
| BD4 | 4 | The craft-delegate skill holds "During a build, `/bench-implement-spec` owns the fold timing, and a ticket charge starts without this fold." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A missing scope sentence raises the new Require diagnostic. |
| BD5 | 4 | The two craft-delegate stale-base sentences stay byte-exact. | `internal/conformance/fixture_bite_test.go` (`TestWorkflowCadenceAnchorsRejectDeletionAndSwap`) | A rewrap or rewrite of either sentence breaks the byte-exact mutation row. |
| BD6 | 5 | The implement-spec:56 fold sentences and the review-implementation:98 step 1 sentence stay pinned. | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | An edit that drops either pinned sentence raises its existing diagnostic. |
| BD7 | 6 | Implement-spec holds "A build dispatches no light-path fix while it runs." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A missing hold sentence raises the new Require diagnostic. |
| BD8 | 7 | Implement-spec holds "A fix that the build's own outcome needs lands as the build's own commit under `.bench/BENCH.md`'s fix-don't-park rule." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A missing own-commit sentence raises the new Require diagnostic. |
| BD9 | 8 | The operating guide no longer holds "`/bench-implement-spec` included, a `bench learning` entry's light-path fix ships at once". | `internal/conformance/commitment_guidance_test.go` (`TestCommitmentGuidance`) | A restored in-build grant raises the new Forbid diagnostic. |
| BD10 | 8 | The operating guide holds "Outside a build, a `bench learning` entry's light-path fix ships at once only when the active outcome needs it". | `internal/conformance/commitment_guidance_test.go` (`TestCommitmentGuidance`) | A removed outside-a-build clause raises the new Require diagnostic. |
| BD57 | 9 | The operating guide holds "`/bench-implement-spec` states the build rule." | `internal/conformance/commitment_guidance_test.go` (`TestCommitmentGuidance`) | A removed pointer clause raises the new Require diagnostic. |
| BD11 | 8 | The operating guide keeps the `**Delegate each light-path fix.**` label and its pinned neighbour sentences. | `internal/conformance/commitment_guidance_test.go` (`TestCommitmentGuidance`) | A changed label leaves the mutation anchor without a match, and the test fails. |
| BD12 | 10 | ADR 0023 holds "A build dispatches no light-path fix while it runs." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A missing ADR sentence raises the new Require diagnostic. |
| BD13 | 10 | ADR 0023 no longer holds "and so does a phase that needs a learning's fix." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A restored unscoped clause raises the new Forbid diagnostic. |
| BD14 | 11 | Implement-spec Land holds "After the last ticket of a chunk, that ticket's author freezes the chunk and records the chunk entry with `bench record chunk`." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A missing author-records sentence raises the rewritten Require diagnostic. |
| BD58 | 11 | Implement-spec Build holds "The last ticket's author performs that freeze, as the Land section states." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A missing freeze-owner sentence raises the new Require diagnostic. |
| BD15 | 12 | Implement-spec Land keeps "Each ticket author then writes its verification entries at that chunk source with `bench record verification`." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A dropped verification sentence raises its existing diagnostic. |
| BD16 | 13 | Implement-spec no longer holds "When the orchestrator freezes a chunk after its last ticket". | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A restored orchestrator-records sentence raises the new Forbid diagnostic. |
| BD17 | 14 | Implement-spec Land holds "Before a successor chunk starts, run `bench gate --checkpoint specs/<slug>/spec.md --chunk <id>`." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A missing successor-only sentence raises the new Require diagnostic. |
| BD18 | 21 | Implement-spec no longer holds "--chunk <id>` before the successor". | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A restored every-chunk sentence raises the new Forbid diagnostic. |
| BD19 | 15 | Implement-spec holds "The last chunk runs no chunk checkpoint gate, because the complete checkpoint that `/bench-final-check` owns grades every chunk again." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A missing last-chunk sentence raises the new Require diagnostic. |
| BD20 | 16 | Implement-spec keeps "after the last chunk, reconcile every acceptance row and the integrated behavior" in lowercase form. | `internal/conformance/docs_workflow_helpers_test.go` (`TestReviewConvergenceContractCurrentDocs`) | A rewrite of the reconciliation sentence fails the pinned phrase. |
| BD21 | 17 | Implement-spec holds "The last chunk keeps its chunk entry, its verification entries, and its three review entries." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A missing record-keep sentence raises the new Require diagnostic. |
| BD50 | 18 | Implement-spec holds "not code, and it refreshes `bench handoff` at each chunk checkpoint and at the final reconciliation." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A missing final-reconciliation refresh raises the rewritten Require diagnostic. |
| BD51 | 19 | README.md holds `more -->\|"yes; chunk checkpoint"\| chunk`. | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A diagram without the checkpoint on the successor edge raises the new Require diagnostic. |
| BD55 | 19 | README.md no longer holds `"none; chunk checkpoint"`. | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A diagram that keeps the checkpoint on the no-successor edge raises the new Forbid diagnostic. |
| BD52 | 20 | ADR 0023 holds "It refreshes the session handoff at each chunk checkpoint and at the final reconciliation". | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A missing ADR refresh clause raises the new Require diagnostic. |
| BD22 | 22 | Implement-spec holds "Code that the author moves or merges keeps its error handling and validation." | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecChunkedSlicing`) | A missing guard sentence raises the new Require diagnostic. |
| BD23 | 23 | Implement-spec holds "A witness is one demonstrated red of one independent test expectation." | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecChunkedSlicing`) | A missing definition raises the new Require diagnostic. |
| BD24 | 24 | Implement-spec holds "The charge caps the witnesses at one for each coverage row plus one for each production branch that the ticket adds." | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecChunkedSlicing`) | A missing cap sentence raises the new Require diagnostic. |
| BD25 | 25 | Implement-spec holds "The author meets the cap with fewer independent expectations, such as one table-driven test, and skips no witness." | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecChunkedSlicing`) | A missing fewer-expectations sentence raises the new Require diagnostic. |
| BD26 | 26 | Implement-spec holds "The author runs the floor checks and witnesses of several packages in one `bench test` call that repeats `--package`." | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecChunkedSlicing`) | A missing batched-call sentence raises the new Require diagnostic. |
| BD27 | 27 | Implement-spec keeps "Before the first edit, the author lists every place the change must reach: callers, tests, fixtures, and exports." | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecChunkedSlicing`) | A dropped reach list raises its existing diagnostic. |
| BD28 | 28 | docs/known-issues.md no longer holds "Cap the witnesses in a charge". | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecChunkedSlicing`) | A surviving copy of the cap raises the new Forbid diagnostic. |
| BD29 | 28 | docs/known-issues.md holds "Apply the witness cap and the batched `bench test` call that `/bench-implement-spec` states." | `internal/anchors/registry_ticket_passes_test.go` (`TestWriteSpecChunkedSlicing`) | A missing owner pointer raises the new Require diagnostic. |
| BD30 | 29 | Delegation-discipline keeps "The second is one recorded red for each independent test expectation." | `internal/anchors/registry_charge_binding_test.go` (`TestChargeBindingAnchors`) | A dropped one-red rule raises its existing diagnostic. |
| BD31 | 30 | Review-implementation step 4 holds "In the default mode, spawn the axes in parallel sub-agents." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A serial dispatch sentence in place of the parallel one raises the new Require diagnostic. |
| BD32 | 31 | Review-implementation holds "When every accepted finding of a round falls on one ticket, one combined repair charge carries them all." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A missing combined-charge sentence raises the new Require diagnostic. |
| BD33 | 31 | Review-implementation holds "That charge holds the findings of every axis and goes to that ticket's repair session." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A missing every-axis sentence raises the new Require diagnostic. |
| BD34 | 32 | Review-implementation holds "When the accepted findings fall on several tickets, each affected ticket keeps its own repair session under `.bench/BENCH.md`'s repair rule." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A missing multi-ticket sentence raises the new Require diagnostic. |
| BD35 | 33 | Review-implementation holds "A confirming round reads only the repair delta and examines only the findings that its charge names." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A missing named-findings scope raises the rewritten Require diagnostic. |
| BD59 | 33 | Review-implementation holds "Its charge names those findings and the folds to confirm." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A dropped folds clause raises the new Require diagnostic. |
| BD36 | 34 | Review-implementation keeps "A chunk that ends on a repair takes one confirming round of all three axes at its final tip." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A dropped three-axis sentence raises its existing diagnostic. |
| BD37 | 35 | Review-implementation no longer holds "look first for duplicated knowledge in the repair delta". | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A restored widening sentence raises the new Forbid diagnostic. |
| BD38 | 36 | Delegation-discipline holds "the round examines only those findings". | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A missing named-findings clause raises the new Require diagnostic. |
| BD54 | 36 | Delegation-discipline no longer holds "asks only for findings above the blocking bar". | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A restored blocking-bar scope raises the new Forbid diagnostic. |
| BD39 | 37 | Implement-spec holds "The orchestrator keeps one build log at the primary checkout's `capture/build-log/<slug>.log`." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A missing build-log sentence raises the new Require diagnostic. |
| BD40 | 38 | Implement-spec holds "For each reviewer direction that saves a step, it appends one line that names the direction and the step that it removed." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A missing step-line sentence raises the new Require diagnostic. |
| BD41 | 39 | Implement-spec holds "The retro author reads that log." | `internal/anchors/registry_chunk_chain_test.go` (`TestChunkChainAnchors`) | A missing retro reader raises the new Require diagnostic. |
| BD42 | 40 | The kit `.gitignore` holds a line equal to `landing.BuildLogDir`. | planned internal/conformance/local_capture_tracking_test.go (TestKitGitignoreIgnoresTheBuildLog) | A kit `.gitignore` without the line fails the comparison. |
| BD43 | 41 | `bench init` writes `capture/build-log/` after the three existing capture lines. | `internal/adopt/adopt_test.go` (`TestInitIgnoresLocalCaptureFiles`) | An init list without the constant leaves the expected byte string unmatched. |
| BD44 | 41 | Each line that `bench init` ignores satisfies `landing.LocalCapturePath`. | `internal/adopt/adopt_test.go` (`TestLocalCaptureIgnoreLinesMatchTheLandingPredicate`) | An init line that the predicate refuses fails the agreement test. |
| BD45 | 42 | `landing.LocalCapturePath` accepts `capture/build-log/`, `capture/build-log`, and `capture/build-log/<slug>.log`. | planned internal/landing/local_capture_test.go (TestLocalCapturePathRecognizesTheBuildLog) | A predicate that keeps only the three exact files refuses each accepted case. |
| BD46 | 43 | `landing.LocalCapturePath` refuses `capture/build-log-foreign/x.log`. | planned internal/landing/local_capture_test.go (TestLocalCapturePathRecognizesTheBuildLog) | A bare prefix match accepts the sibling name. |
| BD47 | 43 | `landing.LocalCapturePath` refuses a `..` segment, a doubled slash, and a control byte under the build-log name. | planned internal/landing/local_capture_test.go (TestLocalCapturePathRecognizesTheBuildLog) | A predicate without the clean-path guards accepts each unsafe case. |
| BD48 | 44 | `landing.CheckoutFingerprint` returns the same value before and after an ignored `capture/build-log/<slug>.log` appears in the checkout. | planned internal/landing/local_capture_test.go (TestCheckoutFingerprintIgnoresTheBuildLog) | A predicate that refuses the build log keeps its `!!` entry in the fingerprint status, so the two values differ. |
| BD49 | 45 | The landed tree leaves roadmap/FT379.md and roadmap/FT369.md byte-identical to the build base. | review-owned: Spec axis compares both files at the final tip | An edit that narrows either row shows in the final delta. |
| BD53 | 46 | The landing retires roadmap/FT410.md and its `ROADMAP.md` row and no other row. | review-owned: `/bench-final-check` reads the landing commit | A delivery closure that removes a second row shows in the landing diff. |

### Edge inventory

The audience of every guidance rule is every repository that links the kit, because the kit ships these files (.bench/consumer-payload.json:8, :11, :15).
The build-log ignore line serves both audiences: the kit through its own `.gitignore`, and a linked repository through `bench init`.
For an absent `capture/build-log/` directory, the predicate never runs, and a landing proceeds as today.
For an empty directory, Git reports no path, and a landing proceeds as today.

**Won't handle:**

- A refreeze-free route for a late `main` move in the last chunk — FT342 owns it, and the completion-landing fold survives at review-implementation:98.
- A plan commit that `base-current` refuses while `main` has moved — FT380-plan-commit and FT342 own the plan-commit lane, and the explicit-range preflight stays the in-scope caller.
- A light-path fix that another session lands while a build runs — FT342 owns the landing refusal, and `/bench-final-check` keeps its rule at final-check:52.
- A combined repair charge for findings on several tickets — FT405 owns one session that records every ticket, and the per-ticket repair rule at .bench/BENCH.md:127 survives.
- An evidence-only confirmation form for a test-only repair — FT379 owns that review-record form, and the three-axis confirming round survives.
- A schema or a verb for the build log — the reviewer chose a plain capture file, and the retro author stays its one reader.
- A linked repository that never reruns `bench init` — its build log stays untracked, as the FT217 log was, until `bench init` runs.

## Ownership fences

The fence is the union of the ticket `Writes:` lines, plus the review pickup.
The reviewer disposition of this fence is part of the spec sign-off.

- `.agents/commands/bench-implement-spec.md`
- `.agents/commands/bench-review-implementation.md`
- `.agents/skills/bench-craft-delegate/SKILL.md`
- `.agents/skills/bench-craft-delegate/references/delegation-discipline.md`
- `.bench/BENCH-reference.md`
- `.bench/BENCH.md`
- `.gitignore`
- `CHANGELOG.md`
- `README.md`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `docs/adr/0023-each-ticket-gets-a-fresh-author.md`
- `docs/field-guide.html`
- `docs/known-issues.md`
- `internal/adopt/adopt_test.go`
- `internal/adopt/init.go`
- `internal/anchors/anchor_harness_diagnostics_test.go`
- `internal/anchors/registry_calibration.go`
- `internal/anchors/registry_calibration_test.go`
- `internal/anchors/registry_charge_binding.go`
- `internal/anchors/registry_charge_binding_test.go`
- `internal/anchors/registry_check_floor.go`
- `internal/anchors/registry_check_floor_test.go`
- `internal/anchors/registry_chunk_chain.go`
- `internal/anchors/registry_chunk_chain_test.go`
- `internal/anchors/registry_commitment.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_debug_loop.go`
- `internal/anchors/registry_decision_maps.go`
- `internal/anchors/registry_decision_maps_test.go`
- `internal/anchors/registry_ft311_preparation.go`
- `internal/anchors/registry_ft311_review_dispatch.go`
- `internal/anchors/registry_retained_workflow.go`
- `internal/anchors/registry_ticket_passes.go`
- `internal/anchors/registry_ticket_passes_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/charge_evidence_guidance_test.go`
- `internal/conformance/commitment_guidance_test.go`
- `internal/conformance/docs_workflow_helpers_test.go`
- `internal/conformance/fixture_bite_test.go`
- `internal/conformance/ft311_preparation_test.go`
- `internal/conformance/implementation_continuation_test.go`
- `internal/conformance/local_capture_tracking_test.go`
- `internal/conformance/retained_workflow_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/conformance/tier_live_tree_test.go`
- `internal/landing/landing.go`
- `internal/landing/landing_test.go`
- `internal/landing/local_capture.go`
- `internal/landing/local_capture_test.go`
- `internal/landing/state_test.go`
- `internal/worktree/build_outputs.go`
- `internal/worktree/land_local_capture_test.go`
- `projects/benchkit.md`
- `tests/canary/claude-agent-definitions`
- `tests/canary/docs-currency-token-diet`
- `tests/canary/guidance-prose-budgets`
- `tests/canary/line-routing`
- `tests/canary/load-validity-metadata`
- `tests/canary/skill-description-budgets`
- `tests/canary/skills-index-command-adapters`
- `tests/canary/workflow-guidance-anchors`
- `reviews/build-directions.md`

## Out of scope

- An enforced build window that refuses side landings is FT342: about 20 edits and 6 gate runs.
- One repair session that records every ticket of its chunk is FT405: about 8 edits and 3 gate runs.
- An evidence-only confirmation form and a record-only repair cycle are FT379: about 10 edits and 4 gate runs.
- A `bench assessment` stage-timing event that reads the build log is FT231: about 12 edits and 4 gate runs.

## Further notes

Flagged additions for reviewer veto:

- The landing predicate change in `internal/landing` is in-scope support for the build log. A refused build-log entry changes the destination fingerprint, and the init agreement test reds.
- The exported build-log directory constant lives in `internal/landing`, and `internal/adopt` gains an import of that package.
- `LocalCapturePath` moves to a new file, internal/landing/local_capture.go, and its tests take a new file beside it, because both landing files sit near the 400-line limit.
- A new conformance test in internal/conformance/local_capture_tracking_test.go reads the kit `.gitignore` against `landing.BuildLogDir`.
- The last chunk's fold moves to the start of that chunk, before its first ticket dispatch, so the author's freeze follows the fold.
- Implement-spec's line budget rises from 85 to 93 across four tickets.
- The craft-delegate scope is a new sentence, and its two pinned stale-base sentences stay byte-exact.
  The authoring charge asked for a rewrite of registry_data.go:69, and this spec leaves that row unchanged.
- The Standards first-look sentence at review-implementation:49 goes away, because criterion confirming-scope forbids a widened round.
- The handoff refresh gains the final reconciliation, because FT410 removes the last chunk checkpoint.
- FT342's open text says that a build's own light-path dispatches wait for a gap (roadmap/FT342.md:31). That text now contradicts the build rule, and the FT342 spec reconciles it.
- FT369's self-check bullet (roadmap/FT369.md:65-67) becomes delivered, and the drain decides how to narrow that row.
- BENCH.md:135 makes a dependent chunk wait for each prerequisite chunk checkpoint. The last chunk is never a prerequisite, so the rule still holds.

Source-sentence-to-row table:

| source criterion | rows |
|---|---|
| build-directions.no-mid-chunk-fold | BD1–BD6, BD60, BD61 |
| build-directions.no-in-build-light-path | BD7–BD13, BD57 |
| build-directions.author-self-check | BD22, BD27 |
| build-directions.witness-cap | BD23–BD26, BD28–BD30 |
| build-directions.author-records | BD14–BD16, BD58 |
| build-directions.parallel-axes | BD31–BD34 |
| build-directions.confirming-scope | BD35–BD38, BD54, BD59 |
| build-directions.step-log | BD39–BD48 |
| build-directions.FT410-last-chunk | BD17, BD18 |
| build-directions.FT410-complete | BD19–BD21, BD50–BD52, BD55 |
| build-directions.open | BD49, BD53 |

The pre-review proof checklist follows.

- Cited symbols: `LocalCapturePath`, `RuntimeIgnoredPath`, `localCaptureIgnoreLines`, `ignoreLocalCapture`, `checkSource`, `checkReviews`, and `RecordChunk` resolve in their current owners.
- Import edges: `go list -deps ./internal/landing` lists no `internal/adopt`, so the new edge from `internal/adopt` makes no cycle.
- Source-row clauses and occurrences: all eleven criteria, the FT410 body, and the nine directions at log:32 were read; direction 8 belongs to FT379.
- Promised field labels: none.
- Moved symbols: `LocalCapturePath` moves from internal/landing/landing.go to internal/landing/local_capture.go in the same package, so no caller changes.
- Changed-function callers: `LocalCapturePath` has two production callers (internal/landing/landing.go:378; internal/worktree/build_outputs.go:81) and one test caller (internal/adopt/adopt_test.go:581).
- Copy survival: BD28 reds a surviving known-issues copy of the cap, and BD13, BD16, BD18, BD37, BD54, and BD55 red each retired sentence.
- Rendered-shape readers: the README diagram label is the one rendered shape, and BD51 and BD55 grade its new and old labels.
- Pin operators: anchor rows use the normalized substring match of internal/anchors/match.go:29-41; the canary and fixture-bite rows use exact bytes.
- Entry reads: none.
- Derived expectations: each mirror rule restates its needle and diagnostic by hand; no expectation comes from the code under test.
- Consolidated rules: none.
- Quantified obligations: "every chunk" in the complete checkpoint holds by internal/reviewrecord/coverage.go:136-141; "each affected ticket" in the repair rule stays as written.
- Workflow-step writes: the author's chunk record writes `reviews/<slug>.md` at the chunk source, as the orchestrator's record did, so the checkpoint digests do not change.

## Completion plan

This version 1 inventory records planned obligations only.
Before the first dispatch, the plan amendment records a version 2 plan with the run, orchestrator, and assignment identities.
Each chunk names one focused check with a probe for each ticket, one batched package call, and the root conformance check.
The batched call repeats `--package` for each package that the chunk writes, so the floor of each ticket is one `bench test` call.

```bench-completion-plan
{"version":1,"chunks":[{"id":"BD-C1","tickets":["01-fold-only-at-chunk-boundaries.md","02-hold-light-paths-during-a-build.md","03-close-the-last-chunk-through-the-complete-checkpoint.md"],"verification":[{"id":"bd-01-focused","command":"bench test --package ./internal/anchors --run TestChunkChainAnchors","probe":"bd-01-omission"},{"id":"bd-02-focused","command":"bench test --package ./internal/conformance --run TestCommitmentGuidance","probe":"bd-02-omission"},{"id":"bd-03-focused","command":"bench test --package ./internal/anchors --run TestChunkChainAnchors","probe":"bd-03-omission"},{"id":"bd-c1-floor","command":"bench test --package ./internal/anchors --package ./internal/conformance --package ./cmd/bench"},{"id":"bd-c1-conformance","command":"bench test --check root-conformance"}]},{"id":"BD-C2","tickets":["04-cap-witnesses-and-keep-moved-guards.md"],"verification":[{"id":"bd-04-focused","command":"bench test --package ./internal/anchors --run TestWriteSpecChunkedSlicing","probe":"bd-04-omission"},{"id":"bd-c2-floor","command":"bench test --package ./internal/anchors --package ./internal/conformance --package ./cmd/bench"},{"id":"bd-c2-conformance","command":"bench test --check root-conformance"}]},{"id":"BD-C3","tickets":["05-scope-review-rounds.md"],"verification":[{"id":"bd-05-focused","command":"bench test --package ./internal/anchors --run TestChunkChainAnchors","probe":"bd-05-omission"},{"id":"bd-c3-floor","command":"bench test --package ./internal/anchors --package ./internal/conformance --package ./cmd/bench"},{"id":"bd-c3-conformance","command":"bench test --check root-conformance"}]},{"id":"BD-C4","tickets":["06-keep-the-build-log.md"],"verification":[{"id":"bd-06-focused","command":"bench test --package ./internal/landing --run 'TestLocalCapturePathRecognizesTheBuildLog|TestCheckoutFingerprintIgnoresTheBuildLog'","probe":"bd-06-omission"},{"id":"bd-06-init","command":"bench test --package ./internal/adopt --run 'TestInitIgnoresLocalCaptureFiles|TestLocalCaptureIgnoreLinesMatchTheLandingPredicate'"},{"id":"bd-c4-floor","command":"bench test --package ./internal/landing --package ./internal/adopt --package ./internal/worktree --package ./internal/anchors --package ./internal/conformance --package ./cmd/bench"},{"id":"bd-c4-conformance","command":"bench test --check root-conformance"}]}],"final_verification":[{"id":"final-preflight","command":"bench preflight build build-directions"},{"id":"final-conformance","command":"bench test --check root-conformance"},{"id":"final-workflow","command":"bench test --check docs-currency-workflow"}]}
```

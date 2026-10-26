# One grade key lets every whole-gate verb reuse another run's green

Status: staged

Roadmap: FT314

Decision source: specs/gate-reuse/decisions/gate-reuse.md (ready compiled map), decisions #1 to #14, on the map-only route

Verification log: 2 iteration(s) to accept — map-only route; the Fable/medium ticket round took 2 iterations (a GR13 producer in a later ticket, a GR58 seam outside its test package, and absent helper inputs in ticket 01).

## Problem

A spec landing grades the same tree as its final complete checkpoint, and it pays a second whole gate.
In 19 days, 30 spec landings did this, at a median of 288 s for each whole gate.
The evidence store names a green by the tree and an oracle digest that frames the checkout path and the baseline runner identity, so a green never crosses checkouts.

## Solution

One grade key names each whole-gate green by its verdict inputs only.
These are the tree, the toolchain, the baseline phase schedule, the capability posture, the run binary member, and the running executable digest.
Every whole-gate verb looks up a green under that key before it picks a reduced schedule.
In the kit, the complete checkpoint runs under the landing broker, so its green answers the landing.

## User stories

Line: opus / medium.
Implementation-line reason: GR-C1 and GR-C2 change the oracle's reuse key, which the profile routes at mid effort. The seams exist, and each wrong key fails a run-count row.
Harder chunks: GR-C1, GR-C2.

### The grade key (decisions #1, #3, #4, #10, #14)

1. As a landing operator, I want a green of the same verdict inputs from another checkout to answer my run, so that one tree pays one whole gate.
2. As a gate maintainer, I want the grade key to leave out the checkout path, so that two checkouts of one repository share one green.
3. As a gate maintainer, I want the grade key to frame no run purpose, so that a verb's own checks, not the key, decide its purpose.
4. As a gate maintainer, I want the bytes of the baseline phase manifest in the key, with an absent manifest distinct from an empty one, so that a changed schedule runs the gate.
5. As a release operator, I want the capability posture in the key, so that a relaxed green never answers a strict run.
6. As a gate maintainer, I want the run binary member to follow the run binary's source root, so that a spec that changes Go code still reuses its own checkpoint green.
7. As a gate maintainer, I want the running executable digest in the key, so that a green from a candidate build never answers the broker's landing.
8. As a build orchestrator, I want only the chunk checkpoint's review record left out of the tree, so that every gate reader sees the published tree.
9. As a gate maintainer, I want the store to hold whole greens only, so that a delta green or a reader record never answers a lookup.

### Evidence lifetime (decision #12)

10. As a landing operator, I want a reused green to answer only inside the window of the run that recorded it, so that reuse never extends the window.
11. As a gate maintainer, I want a red or timed-out run to remove the green under its key, so that no checkout reuses a contradicted green.
12. As a dashboard reader, I want a reuse to record the checkout's gate cache, so that `bench status` shows the reused green.

### The verbs (decisions #5, #11, #13)

13. As a landing operator, I want each whole-gate verb to look up a green before it picks a reader set or a delta set, so that a known verdict costs no phase.
14. As a build orchestrator, I want each verb to run its own purpose and evidence checks on a hit, so that a reused green never skips a verb's own refusal.
15. As a gate operator, I want `--fresh` and a declared lane to skip the lookup, so that an explicit run and a lane keep their meaning.
16. As a build orchestrator, I want the complete checkpoint in the kit to run under the landing broker, so that its green answers the landing.
17. As a gate operator, I want a complete checkpoint with an unusable broker to refuse with a named face, so that it never runs under an unknown executable.

### Records (decisions #2, #6, #7, #8, #9)

18. As a teammate who walks in, I want an ADR and the current docs to state the grade key, so that I read the decided state without the history.
19. As a reviewer, I want the measured saving and the executable roles in the compiled map, so that the build's premise stays checkable.
20. As a drain owner, I want the merge plan disclosure and the light-path stack left out, so that this spec ships only the key that FT450 and FT451 reuse.

## Testing decisions

A good test drives a real gate run on a fixture repository and counts the runs. A run count tells a reuse from a run without a read of any digest.
The prior art is `outcomeFixture` in `internal/gate/run_outcomes_test.go`. Its gate appends to a count file in the common Git directory, which every linked worktree shares.
The identity rows compare two computed digests and hold no literal digest.

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| GR1 | 2 | With an in-root launcher and tool, two identity roots with equal content and a space in one path give one grade identity | `internal/gate/grade_identity_test.go` (`TestGradeIdentityLeavesOutTheCheckoutPath`) | An identity that frames the identity root gives two identities |
| GR2 | 2 | Two identity roots with equal content still give two oracle digests | `internal/gate/grade_identity_test.go` (`TestGradeIdentityLeavesOutTheCheckoutPath`) | An oracle digest that drops the path changes the gate cache key that FT449 closed |
| GR3 | 4 | An absent baseline phase manifest and a present empty one give two grade identities | `internal/gate/baseline_schedule_identity_test.go` (`TestGradeIdentityCarriesTheBaselineSchedule`) | An identity that reads both states as no schedule gives one identity |
| GR4 | 4 | Two baseline phase manifests with different bytes give two grade identities | `internal/gate/baseline_schedule_identity_test.go` (`TestGradeIdentityCarriesTheBaselineSchedule`) | An identity without the schedule gives one identity |
| GR5 | 4 | A manifest that is a live symbolic link gives two grade identities for two targets with different bytes | `internal/gate/baseline_schedule_identity_test.go` (`TestGradeIdentityFollowsTheLoaderRead`) | An identity that frames the link and not the loaded bytes gives one identity |
| GR6 | 4 | A dangling link and a FIFO at the manifest give one grade identity, and the subject build returns | `internal/gate/baseline_schedule_identity_test.go` (`TestGradeIdentityGivesADefectOneValue`) | An identity that frames the defect path gives two identities, and one that opens the FIFO blocks |
| GR7 | 5 | The capability posture changes the grade identity | `internal/gate/oracle_identity_test.go` (`TestOracleIdentityCarriesTheCapabilityPosture`) | An identity without the posture stays equal across the two postures |
| GR8 | 6, 7 | Each of the tool bytes, `PATH`, an environment value, the declared path bytes, the run binary member, and the running executable digest changes the grade identity | `internal/gate/oracle_identity_test.go` (`TestOracleIdentityAndOracleShareTheClosure`) | An identity that drops one member stays equal across its change |
| GR9 | 2 | The reader plan admits a stored table whose identity a sibling identity root computed over equal content | `internal/gate/grade_identity_test.go` (`TestReaderPlanAdmitsATableOfASiblingCheckout`) | A path-bound identity refuses the reader set |
| GR45 | 2 | With a `PATH` tool inside the checkout, an ordinary subject and a prospective subject of one tree give one grade identity | `internal/gate/grade_identity_test.go` (`TestGradeIdentityKeepsTheOrdinaryAndProspectiveSubjectsEqual`) | A spelling relative to the runtime root alone splits the two identities |
| GR54 | 6 | A prospective subject of tree T gives one grade identity from a worktree whose build inputs equal T and from a primary checkout whose build inputs differ | `internal/gate/grade_identity_test.go` (`TestGradeIdentityFollowsTheRunBinarySource`) | A member over the identity root's build inputs gives two identities |
| GR55 | 6 | Two run binary source roots that declare no build inputs give one grade identity | `internal/gate/grade_identity_test.go` (`TestGradeIdentityGivesAnUndeclaredSourceOneValue`) | A member that digests undeclared content gives two identities |
| GR12 | 1 | Two trees under one subject name two evidence records | `internal/gate/prospective_owner_test.go` (`TestProspectiveEvidenceKeysToTheTreeAndRunBinary`) | A key without the tree names one record |
| GR13 | 6, 7 | A changed run binary member names another evidence record, and a changed running executable digest names another record | `internal/gate/prospective_owner_test.go` (`TestProspectiveEvidenceKeysToTheTreeAndRunBinary`) | A key without either member names one record |
| GR10 | 1 | After a green prospective run of one tree in a linked worktree, a run of that tree from the primary checkout leaves the run count at 1 | `internal/gate/grade_key_test.go` (`TestGradeKeyCrossesCheckouts`) | A key that frames the checkout path runs the gate again |
| GR11 | 1 | That primary run prints exactly `gate: green (fresh verdict reused for this tree)` | `internal/gate/grade_key_test.go` (`TestGradeKeyCrossesCheckouts`) | A changed line breaks each reader of the reuse line |
| GR14 | 5 | After a green under the relaxed posture, a run with `BENCH_REQUIRE_CAPABILITIES=1` runs the gate, and the run count reads 2 | `internal/gate/grade_key_test.go` (`TestGradeKeyBindsThePosture`) | A key without the posture reuses the relaxed green |
| GR15 | 4 | After a green in a linked worktree, a primary run whose baseline phase manifest differs runs the gate | `internal/gate/grade_key_test.go` (`TestGradeKeyBindsTheBaselineSchedule`) | A key without the schedule reuses the green |
| GR16 | 3, 14 | After an unpurposed green of the published tree, the complete checkpoint of that tree leaves the run count at 1, and a completion evidence defect still refuses it | `internal/gate/grade_key_test.go` (`TestCompleteCheckpointReusesAnUnpurposedGreenOfItsTree`) | A key with the purpose runs the gate again, and a lookup that skips the evidence pass accepts the defect |
| GR17 | 10 | A linked-worktree green whose recorded time is older than the window does not answer the primary run | `internal/gate/grade_key_test.go` (`TestGradeKeyKeepsTheSourceWindow`) | A window measured from the reuse reuses the old green |
| GR18 | 11 | After a linked-worktree green, a red primary run of the same key removes it, and the next linked-worktree run runs the gate | `internal/gate/grade_key_test.go` (`TestGradeKeyRedRemovesAnotherCheckoutsGreen`) | A red that removes only its own record leaves the green |
| GR19 | 11 | A timed-out run removes the green under its key | `internal/gate/timeout_recovery_count_test.go` (`TestGateRunTimeoutInvalidatesOldEvidence`) | A timeout that keeps the record lets the next run reuse it |
| GR20 | 1 | An evidence record whose `grade` field differs from its file name does not answer, and the run runs the gate | `internal/gate/grade_key_test.go` (`TestGradeKeyRefusesAMisnamedRecord`) | A reader that trusts the file name reuses a copied record |
| GR21 | 9 | A delta checkpoint after an expired whole green of its base still runs the delta set | `internal/gate/delta_checkpoint_test.go` (`TestDeltaCheckpointReadsTheBaseWithoutExpiry`) | A base read under another key finds no green and runs the whole gate |
| GR65 | 8 | A whole green of the chunk subject without its review record serves the chunk checkpoint and never serves the complete checkpoint of the published tree | `internal/gate/review_checkpoint_test.go` (`TestReviewCheckpointReuse`) | A key that strips the record for every verb lets a chunk green answer the complete checkpoint |
| GR53 | 9 | A reader record never answers as a reusable green | `internal/gate/reader_set_test.go` (`TestReaderRecordIsNeverReusableGreen`) | A store that reads the reader record as a grade reuses a partial run |
| GR57 | 9 | A delta green authorizes no landing | `internal/gate/delta_checkpoint_test.go` (`TestDeltaGreenAuthorizesNoLanding`) | A lookup that reads the delta green lands a partial run |
| GR22 | 12 | A primary run with no gate cache that reuses a linked-worktree green writes the cache with the tree, `green`, and the source run's recorded time | `internal/gate/reused_green_test.go` (`TestReusedGreenRecordsTheGateCache`) | A reuse that writes nothing leaves the cache absent |
| GR47 | 12 | A primary run whose gate cache holds a red verdict of the same tree replaces it with the reused green and the source run's time | `internal/gate/reused_green_test.go` (`TestReusedGreenReplacesARedCacheOfItsTree`) | A write rule that compares only the tree leaves the red cache |
| GR23 | 12 | A reuse that the gate cache already describes leaves the cache bytes unchanged | `internal/gate/run_outcomes_test.go` (`TestGateRunGreenReusesRecordedVerdict`) | A reuse that rewrites the record moves its recorded time |
| GR24 | 12 | A reuse whose gate-cache write fails returns an operational result with no reusable green | `internal/gate/durable_replacement_test.go` (`TestReusedGreenCacheWriteFailureRefuses`) | A reuse that ignores the failed write reports green beside a stale cache |
| GR44 | 12 | A reuse while another run holds the checkout's gate lock, with a cache of another tree, prints the reuse line and leaves the cache bytes unchanged | `internal/gate/reused_green_test.go` (`TestReusedGreenBesideAHeldLockWritesNothing`) | A reuse that waits for the lock turns a green into a refusal |
| GR25 | 13 | In a fixture with no lane, after `bench worktree merge --from main` ran the whole gate, a spec-less landing of the merge tip leaves the gate tally at one run | planned TestLandReusesTheMergeGreen in internal/worktree/land_reuse_test.go | A key that frames the checkout path runs the landing gate again |
| GR26 | 13 | That landing advances the green marker to its published commit | planned TestLandReusesTheMergeGreen in internal/worktree/land_reuse_test.go | A publication that reads the reuse as another kind leaves the marker behind |
| GR27 | 16 | After a green complete checkpoint in the source worktree, the spec landing leaves the gate tally at one run | planned TestLandReusesTheCompleteCheckpointGreen in internal/worktree/land_reuse_test.go | A key that frames the checkout path runs the landing gate again |
| GR28 | 6 | When the source changes a Bench build input, the spec landing after a complete checkpoint under an equal running executable leaves the tally at one run | planned TestLandReusesACheckpointGreenOfAGoChange in internal/worktree/land_reuse_test.go | A key over the identity root's build inputs runs the landing gate again |
| GR66 | 13 | In a fixture with no lane, a merge whose composed tree another checkout graded green adds no whole-gate run | planned TestMergeReusesAGreenOfTheComposedTree in internal/worktree/land_reuse_test.go | A merge that skips the lookup runs the whole gate |
| GR58 | 7 | After a complete checkpoint under another running executable, a prospective run of the same tree from the primary checkout runs the whole gate | planned TestCheckpointUnderAnotherExecutableRunsTheGate in internal/gate/verb_reuse_test.go | A key without the running executable reuses a candidate's green |
| GR48 | 13 | A plain `bench gate` in the primary checkout after a green prospective run of the same tree in a linked worktree prints the reuse line | planned TestGateReusesASiblingGreen in internal/gate/verb_reuse_test.go | A plain gate keyed to its checkout runs the gate again |
| GR49 | 13 | In a repository with no lane, `bench commit` of a composition that another checkout graded green adds no whole-gate run | planned TestCommitReusesAGreenOfTheComposedTree in internal/commit/commit_reuse_test.go | A commit that skips the lookup runs the whole gate |
| GR50 | 15 | `bench gate --fresh` after a green runs the gate again | `internal/gate/run_outcomes_test.go` (`TestGateRunCommandFreshRerunsGreen`) | A fresh run that reuses hides a red |
| GR51 | 15 | A commit with a declared lane runs its lane when a whole green of the composed tree exists | planned TestLaneCommitKeepsItsLane in internal/commit/commit_reuse_test.go | A lane verb that looks up a grade skips its lane checks |
| GR52 | 13 | A kit landing that is eligible for the reader set reuses a whole green of its tree and writes no reader record | planned TestReaderSetLandingReusesAWholeGreen in internal/gate/verb_reuse_test.go | A landing that picks the reader set before the lookup runs the readers |
| GR59 | 16 | In a repository that declares build inputs, a complete checkpoint whose running executable differs from the broker manifest execs the broker with the same arguments and working directory | planned TestCompleteCheckpointExecsTheBroker in internal/gate/complete_checkpoint_broker_test.go | A checkpoint that runs in its own process grades under a candidate |
| GR60 | 16 | A complete checkpoint whose running executable equals the broker runs in its own process | planned TestCompleteCheckpointRunsUnderTheBroker in internal/gate/complete_checkpoint_broker_test.go | A checkpoint that always execs loops or pays an extra start |
| GR61 | 17 | A manifest that does not read refuses the complete checkpoint with `checkpoint-broker-unavailable` and no gate run | planned TestCompleteCheckpointRefusesAnUnreadableManifest in internal/gate/complete_checkpoint_broker_test.go | A checkpoint that falls back runs under an unknown executable |
| GR62 | 17 | A broker whose digest differs from its manifest refuses the complete checkpoint with `checkpoint-broker-unavailable` and no gate run | planned TestCompleteCheckpointRefusesABrokerDigestMismatch in internal/gate/complete_checkpoint_broker_test.go | A checkpoint that execs a changed broker trusts tampered bytes |
| GR67 | 16 | A declared-build complete checkpoint under an executable other than the installed candidate execs the installed candidate and does not refuse with `candidate-not-installed` | planned TestDeclaredBuildCheckpointExecsTheInstalledCandidate in internal/gate/complete_checkpoint_broker_test.go | A comparison after `refuseUninstalled` refuses a run that the broker exec would serve |
| GR63 | 16 | A repository that declares no build inputs runs the complete checkpoint under its running executable and reads no manifest | planned TestCompleteCheckpointKeepsTheExecutableWithoutBuildInputs in internal/gate/complete_checkpoint_broker_test.go | A checkpoint that always reads the manifest refuses every linked repository |
| GR64 | 16 | A chunk checkpoint whose running executable differs from the broker runs in its own process | planned TestChunkCheckpointKeepsTheRunningExecutable in internal/gate/complete_checkpoint_broker_test.go | A chunk checkpoint that execs the broker loses the candidate's own grammar |
| GR29 | 18 | An ADR states the grade key members, reuse across checkouts and verbs, the posture binding, and the broker checkpoint | review-owned at the GR-C4 checkpoint | A missing ADR leaves the decision only in this spec |
| GR30 | 18 | ADR 16's evidence sentence, ADR 2, the profile, the reference, and the glossary state the grade key as the reuse rule | review-owned at the GR-C4 checkpoint | A stale sentence tells a teammate that a green never crosses checkouts |

Not covered: story 19 — the compiled map holds the measurement and the executable roles, and no build behavior proves them.
Not covered: story 20 — FT450, FT451, and the parked merge plan idea own the excluded scope.

## Implementation chunks

GR-C1 creates the grade identity that every later ticket consumes, so it stays one small chunk.
GR-C2 keys the store and the gate cache to it. GR-C3 proves each verb and moves the complete checkpoint under the broker. GR-C4 records the decided state.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| GR-C1 / 01-compute-the-grade-identity.md | One grade identity frames the verdict inputs with no checkout path, no purpose, and no baseline runner | GR1, GR2, GR3, GR4, GR5, GR6, GR7, GR8, GR9, GR45, GR54, GR55 | TestGradeIdentityLeavesOutTheCheckoutPath, TestGradeIdentityCarriesTheBaselineSchedule, TestGradeIdentityFollowsTheLoaderRead, TestGradeIdentityGivesADefectOneValue, TestOracleIdentityCarriesTheCapabilityPosture, TestOracleIdentityAndOracleShareTheClosure, TestReaderPlanAdmitsATableOfASiblingCheckout, TestGradeIdentityKeepsTheOrdinaryAndProspectiveSubjectsEqual, TestGradeIdentityFollowsTheRunBinarySource, TestGradeIdentityGivesAnUndeclaredSourceOneValue | yes |
| GR-C2 / 02-name-evidence-by-the-grade-key.md, 03-record-a-reused-green.md | The store names each green by the grade key across checkouts, and a reuse records the checkout's gate cache | GR12, GR13, GR10, GR11, GR14, GR15, GR16, GR17, GR18, GR19, GR20, GR21, GR65, GR53, GR57, GR22, GR47, GR23, GR24, GR44 | TestProspectiveEvidenceKeysToTheTreeAndRunBinary, TestGradeKeyCrossesCheckouts, TestGradeKeyBindsThePosture, TestGradeKeyBindsTheBaselineSchedule, TestCompleteCheckpointReusesAnUnpurposedGreenOfItsTree, TestGradeKeyKeepsTheSourceWindow, TestGradeKeyRedRemovesAnotherCheckoutsGreen, TestGateRunTimeoutInvalidatesOldEvidence, TestGradeKeyRefusesAMisnamedRecord, TestDeltaCheckpointReadsTheBaseWithoutExpiry, TestReviewCheckpointReuse, TestReaderRecordIsNeverReusableGreen, TestDeltaGreenAuthorizesNoLanding, TestReusedGreenRecordsTheGateCache, TestReusedGreenReplacesARedCacheOfItsTree, TestGateRunGreenReusesRecordedVerdict, TestReusedGreenCacheWriteFailureRefuses, TestReusedGreenBesideAHeldLockWritesNothing | yes |
| GR-C3 / 04-reuse-a-green-in-every-whole-gate-verb.md, 05-run-the-complete-checkpoint-under-the-broker.md | Each whole-gate verb reuses another run's green, and the kit's complete checkpoint runs under the landing broker | GR25, GR26, GR27, GR28, GR66, GR58, GR48, GR49, GR50, GR51, GR52, GR59, GR60, GR61, GR62, GR67, GR63, GR64 | TestLandReusesTheMergeGreen, TestLandReusesTheCompleteCheckpointGreen, TestLandReusesACheckpointGreenOfAGoChange, TestMergeReusesAGreenOfTheComposedTree, TestCheckpointUnderAnotherExecutableRunsTheGate, TestGateReusesASiblingGreen, TestCommitReusesAGreenOfTheComposedTree, TestGateRunCommandFreshRerunsGreen, TestLaneCommitKeepsItsLane, TestReaderSetLandingReusesAWholeGreen, TestCompleteCheckpointExecsTheBroker, TestCompleteCheckpointRunsUnderTheBroker, TestCompleteCheckpointRefusesAnUnreadableManifest, TestCompleteCheckpointRefusesABrokerDigestMismatch, TestDeclaredBuildCheckpointExecsTheInstalledCandidate, TestCompleteCheckpointKeepsTheExecutableWithoutBuildInputs, TestChunkCheckpointKeepsTheRunningExecutable | no |
| GR-C4 / 06-record-the-grade-key-rule.md | An ADR and the current docs state the grade key and the broker checkpoint | GR29, GR30 | review-owned rows, root conformance, prose | no |

### Completion plan

The plan uses version 2 with a delegate execution block. Each verification names its ticket, and each ticket of a multi-ticket chunk owns its own gate and conformance runs. The build records one fresh author for each ticket on the declared line, and review closes each frozen chunk before the next author starts.
A probe names bytes that the chunk keeps or writes. The plan amendment re-cuts a probe whose bytes the author moved.

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"ft314-gate-reuse-full-20261026","orchestrator_session":"claude:gr-orchestrator-20261026","author_limit":3,"assignments":{"01-compute-the-grade-identity.md":[{"session":"claude:gr-t01-author-20261026","assignment":"337ab62820dc0982b002fe23e8ebb850","model":"opus","effort":"medium","source":"5cc6bedca3e101a6be42c298f753fd00d0dcbe81","native_ref":"claude-agent:bench-writer-medium:gr-t01-author-20261026"},{"session":"claude:gr-t01-repair1-20261026","assignment":"337ab62820dc0982b002fe23e8ebb850","model":"opus","effort":"medium","source":"e7a51cc4f1bc5e67763313259d4dd80222a6749f","native_ref":"claude-agent:bench-writer-medium:gr-t01-repair1-20261026","predecessor":"claude:gr-t01-author-20261026","trigger":"user-directed","stopped":"claude-agent:bench-writer-medium:gr-t01-author-20261026","preserved":"e7a51cc4f1bc5e67763313259d4dd80222a6749f"}],"02-name-evidence-by-the-grade-key.md":[{"session":"claude:gr-t02-author-20261026","assignment":"da22b98318a660b72d9bab3333f0faf1","model":"opus","effort":"medium","source":"a5f9c772c7b5fb4812ab5b42c1c1611e84bbec4f","native_ref":"claude-agent:bench-writer-medium:gr-t02-author-20261026"},{"session":"claude:gr-t02-repair1-20261026","assignment":"da22b98318a660b72d9bab3333f0faf1","model":"opus","effort":"medium","source":"1c7f5ce2fb08d5949c29dbe17104fa46aff360b5","native_ref":"claude-agent:bench-writer-medium:gr-t02-repair1-20261026","predecessor":"claude:gr-t02-author-20261026","trigger":"user-directed","stopped":"claude-agent:bench-writer-medium:gr-t02-author-20261026","preserved":"1c7f5ce2fb08d5949c29dbe17104fa46aff360b5"}],"03-record-a-reused-green.md":[{"session":"claude:gr-t03-author-20261026","assignment":"e85f952832c56ff5bf5ed50c3e7defc6","model":"opus","effort":"medium","source":"c702eb07d1869efa03cd602ad305776e80744aa8","native_ref":"claude-agent:bench-writer-medium:gr-t03-author-20261026"}],"04-reuse-a-green-in-every-whole-gate-verb.md":[],"05-run-the-complete-checkpoint-under-the-broker.md":[],"06-record-the-grade-key-rule.md":[]}},"chunks":[{"id":"GR-C1","tickets":["01-compute-the-grade-identity.md"],"verification":[{"id":"gate","command":"bench test --package ./internal/gate","ticket":"01-compute-the-grade-identity.md"},{"id":"authorization","command":"bench test --package ./internal/gate/authorization","ticket":"01-compute-the-grade-identity.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"01-compute-the-grade-identity.md"},{"id":"posture-probe","command":"bench probe internal/gate/subject.go --swap 'strconv.FormatBool(strict())' --with 'strconv.FormatBool(false)' --package ./internal/gate --run TestOracleIdentityCarriesTheCapabilityPosture --expect TestOracleIdentityCarriesTheCapabilityPosture","probe":"Drop the posture from the grade identity. The two postures then give one identity, so the posture row fails.","ticket":"01-compute-the-grade-identity.md"}]},{"id":"GR-C2","tickets":["02-name-evidence-by-the-grade-key.md","03-record-a-reused-green.md"],"verification":[{"id":"gate","command":"bench test --package ./internal/gate","ticket":"02-name-evidence-by-the-grade-key.md"},{"id":"authorization","command":"bench test --package ./internal/gate/authorization","ticket":"02-name-evidence-by-the-grade-key.md"},{"id":"worktree","command":"bench test --package ./internal/worktree","ticket":"02-name-evidence-by-the-grade-key.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"02-name-evidence-by-the-grade-key.md"},{"id":"reuse-probe","command":"bench probe internal/gate/engine.go --swap 'inspection.ReusableGreen, inspection.Reason = true, \"\"' --with 'inspection.ReusableGreen, inspection.Reason = false, \"probe\"' --package ./internal/gate --run TestGradeKeyCrossesCheckouts --expect TestGradeKeyCrossesCheckouts","probe":"Make the store answer no green. The primary run then runs the gate again, so the cross-checkout row fails.","ticket":"02-name-evidence-by-the-grade-key.md"},{"id":"gate-03","command":"bench test --package ./internal/gate","ticket":"03-record-a-reused-green.md"},{"id":"root-conformance-03","command":"bench test --check root-conformance","ticket":"03-record-a-reused-green.md"},{"id":"commands","command":"bench test --package ./cmd/bench","ticket":"02-name-evidence-by-the-grade-key.md"}]},{"id":"GR-C3","tickets":["04-reuse-a-green-in-every-whole-gate-verb.md","05-run-the-complete-checkpoint-under-the-broker.md"],"verification":[{"id":"gate","command":"bench test --package ./internal/gate","ticket":"04-reuse-a-green-in-every-whole-gate-verb.md"},{"id":"worktree","command":"bench test --package ./internal/worktree","ticket":"04-reuse-a-green-in-every-whole-gate-verb.md"},{"id":"commit","command":"bench test --package ./internal/commit","ticket":"04-reuse-a-green-in-every-whole-gate-verb.md"},{"id":"refusalroute","command":"bench test --package ./internal/refusalroute","ticket":"05-run-the-complete-checkpoint-under-the-broker.md"},{"id":"commands","command":"bench test --package ./cmd/bench","ticket":"05-run-the-complete-checkpoint-under-the-broker.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"04-reuse-a-green-in-every-whole-gate-verb.md"},{"id":"land-probe","command":"bench probe internal/gate/engine.go --swap 'inspection.ReusableGreen, inspection.Reason = true, \"\"' --with 'inspection.ReusableGreen, inspection.Reason = false, \"probe\"' --package ./internal/worktree --run TestLandReusesTheCompleteCheckpointGreen --expect TestLandReusesTheCompleteCheckpointGreen","probe":"Make the store answer no green. The landing then runs the gate again, so the checkpoint-then-landing row fails.","ticket":"04-reuse-a-green-in-every-whole-gate-verb.md"},{"id":"gate-05","command":"bench test --package ./internal/gate","ticket":"05-run-the-complete-checkpoint-under-the-broker.md"},{"id":"root-conformance-05","command":"bench test --check root-conformance","ticket":"05-run-the-complete-checkpoint-under-the-broker.md"}]},{"id":"GR-C4","tickets":["06-record-the-grade-key-rule.md"],"verification":[{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"06-record-the-grade-key-rule.md"},{"id":"prose","command":"bench test --check prose","ticket":"06-record-the-grade-key-rule.md"}]}],"final_verification":[{"id":"integration","command":"bench test --changed --base cfbd1d360c2d0048bad07d142468c4bb47f8fc20"},{"id":"system","command":"bench test --check system"},{"id":"whole-project","command":"bench gate"}]}
```

## Further notes

### Gap log

Each entry takes one class: (a) map discretion, (b) a seam or test choice, (c) a missing fact, or (d) a reviewer decision.

- (a) The evidence record names its grade key in a `grade` field. The gate cache record keeps its `oracle` field, because the oracle digest keeps its frames.
- (a) The broker refusal face is `checkpoint-broker-unavailable`. Its route repeats the landing's broker repair advice.
- (b) The run purpose leaves the store key and keeps its frames in the per-checkout oracle digest, which `applyCheckpoint` derives today.
- (b) `EvidenceInspection` carries the grade key in place of its `Oracle` field, so the landing authorization token binds the key.
- (b) `evidencePath`, `retainGreen`, and `invalidateEvidence` keep their signatures. Their test callers stay outside the `Writes:` lines.
- (b) A package-level seam supplies the running executable to the subject build and computes its digest once for each process, so a test can vary it.
- (b) One helper owns the run binary's source-root rule. The selection and the grade identity both call it, and the identity root stands in for the storage root.
- (b) The oracle digest keeps the baseline runner frame, so the per-checkout gate cache keeps its key. Only the grade identity drops it.
- (b) The broker comparison runs before `refuseUninstalled`. A declared build's installed candidate is the broker, so the exec serves that build.
- (b) A replaceable exec function carries the broker exec, so a test observes the arguments and the working directory with no process replacement.
- (b) New gate tests take new files. The growth ratchet grades only files over their line limit (`internal/structure/structure.go:309`), and `refusal_route_test.go` is at its 400-line limit.
- (c) No byte comparison of two run binaries from one tree was made. The validation plan of the measurement asset owns it.

### Ticket sizing

Each estimate takes the byte total of the ticket's `Writes:` files at about four bytes per token. It adds about 50k tokens of fixed context and 20k to 30k tokens for the verification runs.
Each estimate is an upper bound, and each stays under the 150k-token author target.

| ticket | `Writes:` bytes | estimate |
|---|---|---|
| 01-compute-the-grade-identity.md | about 145 KB, 74 KB of it `.bench/commitment.json` | about 115k tokens |
| 02-name-evidence-by-the-grade-key.md | about 170 KB | about 120k tokens |
| 03-record-a-reused-green.md | about 70 KB | about 90k tokens |
| 04-reuse-a-green-in-every-whole-gate-verb.md | about 60 KB | about 95k tokens |
| 05-run-the-complete-checkpoint-under-the-broker.md | about 60 KB | about 90k tokens |
| 06-record-the-grade-key-rule.md | about 245 KB | about 125k tokens |

`.bench/commitment.json` joins ticket 01, because the spec-stage close commits the binding and no author reads that file.
The estimate for ticket 06 counts the 101 KB of `CHANGELOG.md`, which the author appends to and need not read whole.

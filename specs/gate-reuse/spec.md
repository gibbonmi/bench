# One authorization key lets a composition reuse gate evidence across runs

Status: staged

Roadmap: FT314

Decision source: `roadmap/FT314.md`, a named reviewed artifact (reviewer approval 2026-10-14, recorded by drain d-6c97ed54c629), with two late clarifications closed in the reviewer-confirmed conversation on 2026-10-26.

Verification log: pending — the spec review round is open.

## Problem

A composition that one run already graded green pays a second whole gate when a later run grades it again from another checkout.
The FT314 census of 2026-10-07 to 2026-10-14 counted 29 merge-then-land pairs that ran two whole gates on one composition, about two hours in total.

Since commit `3519feb75` (2026-10-19), a kit `bench worktree merge` grades through the target's built-in lane, so that kit pair no longer recurs.
The source is `mergeOwner` in `internal/worktree/merge.go` and `TestMergeFromThePrimaryGradesWithTheTargetLane`. Two duplicates remain:

- A merge in a repository with no declared lane runs the whole gate, and the spec-less landing of that merge tip runs it again.
- A `--complete` checkpoint grades the published tree of its source, and the spec landing of that source grades the same tree again.

FT451 also needs a key that a landing can look up, and FT450 reuses the same key for each composed stack.

The gate cannot reuse a green across checkouts today. The prospective evidence store names a green by the tree and the oracle digest.
The oracle digest frames the caller's checkout path, so each checkout names another record. The sources are `buildSubjectForTree` in `internal/gate/subject.go` and `evidenceName` in `internal/gate/engine.go`.

The same key also misses two members that the reviewer named. The capability posture enters only the tree-free identity, so a relaxed green answers a strict run (`TestOracleIdentityCarriesTheCapabilityPosture`).
The bytes of the baseline phase manifest enter no digest, although that manifest selects the phases of a prospective run (`phasePolicyRoot` and `readManifest` in `internal/gate/manifest.go`).

A coordinator also cannot learn the cost of a merge before it runs. On 2026-09-22 a review-venue update ran a whole gate, and two concurrent follow-on updates refused its active owner.

## Solution

One authorization key names each prospective green. It binds the composed tree, the gate version, the toolchain, the capability posture, and the run purpose. It does not bind the caller's checkout path.
A green that one run records answers any later run whose key is equal, from any checkout of the repository. A change of any key member runs the whole gate.

A reused green answers only inside the freshness window of the run that recorded it. A red or timed-out run removes the green under its key.
A reuse writes the checkout's gate cache with the recorded time of that run, so the dashboard shows the green and the window does not slide.

A spec-less landing after a merge in a repository with no lane then reuses the merge green. A spec landing after its green `--complete` checkpoint reuses that green, when the build changes no Bench build input.

`bench worktree merge --plan --from <commit|target> <target>` prints the merge kind and the projected grade, and it changes no ref, checkout, or record. The grade is `none`, `lane`, `reuse`, or `whole`.
The merge that follows on an unchanged repository runs the projected grade.

## User stories

Line: opus / medium.
Implementation-line reason: GR-C1 changes the oracle's reuse key, which the profile routes at mid effort and above. Its seams exist, its rows run real gate fixtures, and each wrong key fails a run-count row.
Harder chunks: GR-C1.

### The authorization key

1. As a landing operator, I want another checkout's green of the same composition to answer my run, so that one composition pays one whole gate.
2. As a gate maintainer, I want the key to bind the composed tree, so that a green of another tree never answers.
3. As a gate maintainer, I want the key to bind the gate version, so that a green from another policy, launcher, or baseline runner never answers.
4. As a gate maintainer, I want the key to bind the toolchain, so that a green from another `PATH`, tool, environment value, or declared input never answers.
5. As a release operator, I want the key to bind the capability posture, so that a relaxed green never answers a strict run.
6. As a gate maintainer, I want the key to bind the run purpose, so that an unpurposed green never answers a completion.
7. As a gate maintainer, I want the key to leave out the caller's checkout path, so that two checkouts of one repository share one green.
8. As a gate maintainer, I want an executable inside the runtime root or the identity root spelled by its path relative to that root, so that the checkout path cannot enter the key through a launcher or a tool.
9. As a gate maintainer, I want the bytes of the baseline phase manifest in the key, with an absent manifest distinct from an empty one, so that a changed schedule runs the whole gate.
10. As a gate maintainer, I want the identity to read the baseline manifest through the loader's own read, so that the key binds the schedule that the run uses and a special file never blocks.
11. As a gate maintainer, I want a change of any key member to run the whole gate, so that reuse never crosses an environment change.
12. As a dashboard reader, I want the per-checkout oracle digest to keep its frames, so that each gate cache keeps its meaning.
13. As a landing operator, I want an ordinary green and a prospective green of one tree in one checkout to keep one key, so that the inherited-green attribution keeps its answer.

### Evidence lifetime

14. As a landing operator, I want a reused green to answer only inside the freshness window of the run that recorded it, so that reuse never extends the window.
15. As a gate maintainer, I want a red or timed-out run to remove the green under its key, so that no checkout reuses a green that a newer run contradicts.
16. As a gate maintainer, I want each evidence record to carry its key, so that a misnamed record fails closed.
17. As a dashboard reader, I want a reuse to write the checkout's gate cache with the recorded time of the source run, so that `bench status` shows the reused green.
18. As a gate maintainer, I want a reuse that the gate cache already describes to write nothing, so that a repeated reuse never moves the recorded time.
19. As a gate maintainer, I want a failed gate-cache write to refuse the reuse, so that no green claims a cache that the run could not write.
20. As a gate operator, I want a reuse beside another run in the same checkout to answer with no cache write, so that the write never turns a reuse into a refusal.
21. As a gate operator, I want the reuse line unchanged, so that every reader of that line keeps working.
22. As a build orchestrator, I want the delta checkpoint to keep its base read, so that a test-only repair still runs the delta set.

### The landing evidence chain

23. As a landing operator in a repository with no lane, I want a spec-less landing after a merge to reuse the merge green, so that the fold pays the only whole gate.
24. As a landing operator, I want a landing that reuses a green to advance the green marker, so that the next landing plans from that commit.
25. As a build orchestrator, I want a spec landing after its green `--complete` checkpoint to reuse that green, so that the completion pays one whole gate.
26. As a build orchestrator, I want a landing whose source changes a Bench build input to run the whole gate, so that a green from another baseline runner never answers.
27. As a kit maintainer, I want a reader table that another checkout stored under an equal identity to admit the next reader set, so that a reuse does not force the next landing to the whole gate.
28. As a teammate who walks in, I want an ADR and the current docs to state the reuse rule, so that I read the decided state without the history.

### The merge plan

29. As a coordinator, I want `bench worktree merge --plan` to print the merge kind and the projected grade, so that I know the full-gate cost before I run the merge.
30. As a coordinator, I want one grade name for each schedule, so that `none`, `lane`, `reuse`, and `whole` each mean one cost.
31. As a coordinator, I want the plan to change no ref, checkout, gate record, or evidence, so that a plan is safe beside another run.
32. As a coordinator, I want the merge after a plan on an unchanged repository to run the projected grade, so that the projection preserves the gate schedule.
33. As a coordinator, I want a plan refusal to take the merge's face with a rerun that keeps `--plan`, so that its route repeats the plan.
34. As a CLI maintainer, I want the merge help row derived from the one merge grammar, so that help and the refusal grammar cannot drift.
35. As a teammate, I want the reference, the profile, and the glossary to state the plan form, so that I find the disclosure without a source read.

### Reviewed exclusion

36. As a drain owner, I want the landing of several light-path sources under one gate left to FT450's landing queue, so that this spec ships only the key that the queue reuses.

## Implementation decisions

### Terms

- **authorization key** — the digest that names one prospective green in the evidence store. It frames the composed tree, the oracle identity, and the run purpose. Avoid: grade key, evidence key, reuse key.
- **oracle identity** — the digest of the gate version, the toolchain, and the capability posture, with no tree and no checkout path. Avoid: subject identity, environment hash.
- **grade** — the schedule that a merge runs on its composed tree: `none`, `lane`, `reuse`, or `whole`. Avoid: cost class, merge mode.
- **merge plan** — the `--plan` form of the worktree merge. It projects the kind and the grade, and it changes no ref, checkout, or record. Avoid: dry run, merge preview.

### The oracle identity

The identity takes every frame of the oracle digest except the tree and the checkout path. The members fall into three groups:

- The gate version is the policy version, the gate resolution and its launcher closure, the baseline runner identity, and the baseline phase schedule.
- The toolchain is `PATH`, the declared tools, the declared environment values, the declared path bytes, and the gate-input manifest.
- The capability posture is the `BENCH_REQUIRE_CAPABILITIES` value, framed last as today.

An executable inside the runtime root enters the identity by its path relative to the runtime root. Otherwise, an executable inside the identity root enters by its path relative to the identity root.
Every other executable keeps its absolute path. Today `identityPath` spells both cases as one identity-root path, so an ordinary subject and a prospective subject of one checkout share one digest. The relative spelling keeps that equality.

The baseline phase schedule is the answer of `readManifest` at the baseline's `.bench/phases.json`, the loader's own read. Present bytes enter as bytes, and a live link enters as the bytes of its target, which the run then uses.
An absent manifest enters as one distinct value. A defect answer, such as a dangling link, a FIFO, or a directory, enters as one stable defect value with no path in it.
`readManifest` stats the path before it opens it, so a FIFO never blocks the identity. A defect manifest is never a schedule that a run uses: the loader refuses it before any phase, and a plain gate script never reads it.

The oracle digest keeps its frames and its order, which is FT449's closed rule. It keys the per-checkout gate cache that `bench status` and `bench gate` read. Rows GR1 to GR9 and GR45 grade the identity.

### The authorization key

The key frames the tree and the oracle identity. A purposed run then frames its purpose over the key: the checkpoint value, the completion folder, and the source tip.
That purpose framing exists today over the oracle digest in `applyCheckpoint`. One function frames the purpose for both digests, so the two framings cannot drift.

The evidence store under the common Git directory names each green by its key. The record carries the key in an `authorization` field and no `oracle` field.
A reader refuses a record whose key differs from its own, as changed evidence. Every writer and reader of the store routes through one name function. Rows GR10 to GR21 grade the key.

### Reuse and its lifetime

The reuse predicate is today's `reusableEvidence` decision over the key. The store must hold a fresh green under the key, and the checkout's gate cache must not be pending, invalid, or unavailable.
A store record is fresh only inside the freshness window of its own recorded time. The run evaluates the predicate at its two existing sites, before the gate lock and under it.

- Before the lock: when the predicate holds and the checkout's gate cache already describes this subject's fresh green, the run answers with no write, as today.
- At the lock: the run takes the gate lock with no wait, as today. When another run holds the lock and the predicate held before the lock, the run answers the reuse with no write.
- Under the lock: when the predicate holds, the run writes the gate cache unless it already describes this subject's fresh green. The write uses the lock that the run already holds.

The written cache record holds the tree, the checkout's oracle digest, `green`, and the recorded time of the store record. So an absent cache, a cache of another tree, and a red, timed-out, or expired cache of the same tree each take the write.
A failed write returns an operational result, and the run reports no green. Rows GR22, GR23, GR24, GR44, and GR47 grade this step.

A red or timed-out run removes the record under its key, whichever checkout wrote it. Each run in each checkout holds only its own checkout's gate lock.
So the store holds the outcome of the run that finishes last under a key, and each outcome is a real run of an identical subject.

The capability posture now binds every reuse, the same-checkout retry included. This reverses FT449's rule that a green stays reusable under the other posture.
The reviewer approved that reversal on 2026-10-26.

### Trust posture

A reused green authorizes a run with no execution. Its authority is one store record: a `0600` file in a `0700` directory under the common Git directory. Only a gate run transaction writes it.
A record cannot authenticate the run that wrote it. The wrapper's freshness check binds the running executable to its sources.

In a root that ships `.bench/gate-prospective.sh`, as the kit does, the identity frames the baseline runner identity. A kit worktree whose Bench build inputs differ from the destination's then has another identity, so its green never answers the destination's landing.
A root without that script frames no baseline runner. Its gate runs the launcher from the graded tree and the installed kit executable, and the tree member binds the launcher.

ADR 2's threat model covers an accidental weakening, not a determined adversary. Under that model the residual stays accepted, and ADR 2's posture 5 names it.
An ordinary run grades the working tree in place, and a prospective run grades a private checkout of the same tree. Both runs share one key within one checkout today, and after this spec that equivalence crosses checkouts.

### The merge plan

The grammar becomes `bench worktree merge [--plan] --from <commit|target> <target>`. The plan resolves the target, the incoming commit, the previous tip, and the checkout fingerprint as the merge does.
It decides the kind and the composed tree through the merge owner's own composition, with no authorization.

The grade follows from the kind and the target, in this order:

- `none` is the grade of a `current` kind, whatever lane the target declares.
- `lane` is the grade of a target that declares a lane.
- `reuse` is the grade of a target with no lane when the reuse predicate holds for the composed tree.
- `whole` is the grade of every other target with no lane.

The gate exports the reuse predicate for a prospective tree, and the exported form is the `reusableEvidence` decision with its cache-state refusal. The merge run and the plan both read it.
So the projection cannot diverge from the run's decision on an unchanged repository. Row GR46 grades a pending cache.

The plan prints `merge_plan{worktree=<id>,from=<commit>,kind=<kind>,previous_tip=<commit>,tree=<tree>,grade=<grade>}` and exits 0. A refusal takes the merge's face for the same cause, and its rerun keeps `--plan`.
The plan moves no ref, resets no checkout, writes no gate cache or evidence record, and takes no gate lock.

The reuse predicate needs a subject, so the plan opens a prospective artifact bundle as evidence inspection does. The bundle registers a private checkout for the read and removes that registration on close.
The plan can also write unreachable tree objects, as the landing dry run does. The merge help row derives its suffix from the one merge grammar, as its sibling rows do.

### Docs

A new ADR 31 records the decision: one authorization key names each prospective green, and a green crosses runs and checkouts under that key.
ADR 2's posture 5 and ADR 16's evidence sentence name the key. The profile, the reference, the glossary, and the changelog state the current rule and the plan form.

### Consolidation table

| site | old rule | new rule | disposition |
|---|---|---|---|
| oracle identity frames | every oracle frame but the tree, then the posture | every oracle frame but the tree and the checkout path, with relative in-root paths, the baseline schedule, then the posture | GR1, GR3 to GR8, GR45 |
| oracle digest | tree, checkout path, and the other frames | unchanged | GR2 |
| evidence store name | the tree and the oracle digest | the authorization key | GR10, GR12, GR13 |
| evidence record match | tree and oracle equal | tree and authorization key equal | GR20 |
| purpose framing | over the oracle digest | over the oracle digest and over the key, by one function | GR16 |
| posture in reuse | not bound | bound | GR14 |
| gate cache on a reuse | no write | a write with the source time unless the cache describes the subject's fresh green, under the held lock | GR22, GR23, GR24, GR44, GR47 |
| reader table identity | path-bound | path-free | GR9 |
| merge help row | a hand-written suffix | the suffix of the merge grammar | GR42 |

## Implementation chunks

Each chunk follows its predecessor's green checkpoint and its independent review. GR-C1 changes the key at the gate seam and proves the landing chain on it.
GR-C2 adds the merge plan, which reads the reuse predicate that GR-C1 settles. GR-C3 records the decided state in the ADRs and the docs.
The docs take their own chunk because the docs files alone hold about 60k tokens, and each ticket stays near 150k tokens of author context.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| GR-C1 / 01-leave-the-checkout-out-of-the-identity.md, 02-key-the-evidence-store.md, 03-record-a-reused-green.md | One authorization key names every prospective green, a reuse records the checkout's gate cache, and a landing publishes on another run's green | GR1, GR2, GR3, GR4, GR5, GR6, GR7, GR8, GR9, GR45, GR10, GR11, GR12, GR13, GR14, GR15, GR16, GR17, GR18, GR19, GR20, GR21, GR22, GR23, GR24, GR44, GR47, GR25, GR26, GR27, GR28 | TestOracleIdentityLeavesOutTheCheckoutPath, TestOracleIdentityCarriesTheBaselineSchedule, TestOracleIdentityFollowsTheLoaderRead, TestOracleIdentityGivesADefectOneValue, TestOracleIdentityKeepsTheOrdinaryAndProspectiveSubjectsEqual, TestReaderPlanAdmitsATableOfASiblingCheckout, TestAuthorizationKeyCrossesCheckouts, TestAuthorizationKeyBindsThePosture, TestAuthorizationKeyBindsTheBaselineSchedule, TestCompleteCheckpointIgnoresAnUnpurposedGreenOfItsTree, TestAuthorizationKeyKeepsTheSourceWindow, TestAuthorizationKeyRedRemovesAnotherCheckoutsGreen, TestAuthorizationKeyRefusesAMisnamedRecord, TestReusedGreenRecordsTheGateCache, TestReusedGreenReplacesARedCacheOfItsTree, TestReusedGreenCacheWriteFailureRefuses, TestReusedGreenBesideAHeldLockWritesNothing, TestLandReusesTheMergeGreen, TestLandReusesTheCompleteCheckpointGreen, TestLandRunsTheGateForAnotherBaselineRunner | yes |
| GR-C2 / 04-plan-a-merge-grade.md | `bench worktree merge --plan` projects the kind and the grade, and the merge runs that grade | GR31, GR32, GR33, GR34, GR35, GR36, GR37, GR38, GR39, GR40, GR46, GR41, GR42 | TestMergePlanNamesEachGrade, TestMergePlanPrintsOneRecord, TestMergePlanChangesNothing, TestMergeRunsThePlannedGrade, TestMergePlanProjectsAPendingCacheAsWhole, TestMergePlanRefusalKeepsThePlan, TestHelpInventoryIsComplete | no |
| GR-C3 / 05-record-the-reuse-rule.md | ADR 31 and the current docs state the authorization key and the plan form | GR29, GR30, GR43 | review-owned rows, root conformance, prose | no |

### Completion plan

The plan uses version 1. The build records one fresh author for each ticket on the declared line, and review closes each frozen chunk before the next author starts.
A probe names bytes that the chunk keeps or writes. The plan amendment re-cuts a probe whose bytes the author moved.

```bench-completion-plan
{"version":1,"chunks":[{"id":"GR-C1","tickets":["01-leave-the-checkout-out-of-the-identity.md","02-key-the-evidence-store.md","03-record-a-reused-green.md"],"verification":[{"id":"gate","command":"bench test --package ./internal/gate"},{"id":"authorization","command":"bench test --package ./internal/gate/authorization"},{"id":"worktree","command":"bench test --package ./internal/worktree"},{"id":"commands","command":"bench test --package ./cmd/bench"},{"id":"root-conformance","command":"bench test --check root-conformance"},{"id":"census","command":"bench test --package ./internal/conformance --run TestBranchNativeArchitectureCensus"},{"id":"posture-probe","command":"bench probe internal/gate/subject.go --swap 'frame(id, strconv.FormatBool(strict()))' --with 'frame(id, \"false\")' --package ./internal/gate --run TestAuthorizationKeyBindsThePosture --expect TestAuthorizationKeyBindsThePosture","probe":"Drop the posture from the oracle identity. The strict run then reuses the relaxed green, so the posture row fails."},{"id":"reuse-probe","command":"bench probe internal/gate/engine.go --swap 'inspection.ReusableGreen, inspection.Reason = true, \"\"' --with 'inspection.ReusableGreen, inspection.Reason = false, \"probe\"' --package ./internal/worktree --run TestLandReusesTheMergeGreen --expect TestLandReusesTheMergeGreen","probe":"Make the store answer no green. The landing then cannot reuse the merge green, so the merge-reuse row fails."}]},{"id":"GR-C2","tickets":["04-plan-a-merge-grade.md"],"verification":[{"id":"worktree","command":"bench test --package ./internal/worktree"},{"id":"landing","command":"bench test --package ./internal/landing"},{"id":"gate","command":"bench test --package ./internal/gate"},{"id":"usage","command":"bench test --package ./internal/usage"},{"id":"commands","command":"bench test --package ./cmd/bench"},{"id":"root-conformance","command":"bench test --check root-conformance"},{"id":"current-probe","command":"bench probe internal/landing/merge.go --swap 'if contains {' --with 'if false {' --package ./internal/worktree --run TestMergePlanNamesEachGrade --expect TestMergePlanNamesEachGrade","probe":"Skip the current-kind decision. A current target then composes and grades, so the none-grade row fails."}]},{"id":"GR-C3","tickets":["05-record-the-reuse-rule.md"],"verification":[{"id":"root-conformance","command":"bench test --check root-conformance"},{"id":"prose","command":"bench test --check prose"}]}],"final_verification":[{"id":"integration","command":"bench test --changed --base cfbd1d360c2d0048bad07d142468c4bb47f8fc20"},{"id":"system","command":"bench test --check system"},{"id":"whole-project","command":"bench gate"}]}
```

## Testing decisions

A good test drives a real gate run on a fixture repository and counts the runs. A run count tells a reuse from a run without a read of any digest.
The prior art is `outcomeFixture` in `internal/gate/run_outcomes_test.go`. Its gate appends to a count file in the common Git directory, which every linked worktree shares.

The identity rows build subjects through `buildSubjectForTree` with two identity roots, as `TestOracleIdentityAndOracleShareTheClosure` does. They compare two computed digests and hold no literal digest.
The store rows run `Execute` and `ExecuteTree` from a primary checkout and from a linked worktree of one fixture. Each run builds its subject from disk, so no in-memory value crosses the two runs.

The failed cache write uses the replacer seam of `internal/gate/durable_replacement_test.go`. The other gate-cache rows live in `internal/gate/reused_green_test.go`.
The purpose row lives in `internal/gate/authorization_key_test.go`, because `complete_checkpoint_test.go` holds 357 lines. It calls `closureSource` and the published-tree derivation of that file in the same package.
The landing rows use the landing fixture of `internal/worktree/land_fixtures_test.go`, whose prospective gate appends to a tally in the common Git directory. The merge-then-land row composes that fixture with the merge set of `mergeSetAt`.
The merge plan rows use the merge sets of `internal/worktree/merge_target_lane_test.go` and `internal/worktree/refusal_route_follow_test.go`, which declare a lane or none and count whole-gate runs.

`internal/gate/reader_plan_test.go` holds 400 lines, so the reader-table row lives in `internal/gate/oracle_identity_test.go` and calls the reader plan's in-package fixture.
The baseline-schedule rows take a new sibling file, `internal/gate/baseline_schedule_identity_test.go`, so that `oracle_identity_test.go` stays under its 400-line limit.
Each new gate test builds its linked worktree through an existing in-package helper such as `runGitIn`. The architecture census refuses a direct `exec.Command`, or a `Run` or `Output` call with a `worktree` literal, in a gate test.
No test crosses a package seam through another package's private helper. The gate seam observes every row through the `test` and `system` phases of the whole gate.

### Seam diagram

    trigger: a gate run, a checkpoint, a landing, or a merge in any checkout
        │
        ▼
    tree + baseline + toolchain + posture + purpose
        │
        ▼
    [ subject build: oracle digest, oracle identity, authorization key ]
        ◀ tests attach here: two identity roots, two postures, or two manifests, then compare digests
        │
        ▼
    [ evidence store under the common Git directory ]  ──▶  reuse, or a whole run
        ◀ tests attach here: run in a linked worktree, then in the primary, and count runs
        │
        ▼
    [ checkout gate cache ]  ──▶  bench status, bench gate
        ◀ tests attach here: read the cache record after a reuse, with the lock free or held

    trigger: bench worktree merge --plan --from <x> <target>
        │
        ▼
    [ merge owner composition ] ──▶ [ lane, or the exported reuse predicate ] ──▶ merge_plan{...}
        ◀ tests attach here: plan, then merge, and compare the grade with the tally and output

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| GR1 | 7, 8 | With an in-root launcher and tool, two identity roots with equal content and a space in one path give one oracle identity | planned TestOracleIdentityLeavesOutTheCheckoutPath in internal/gate/oracle_identity_test.go | An identity that frames the identity root, or spells an in-root path through it, gives two identities |
| GR2 | 12 | Two identity roots with equal content still give two oracle digests | planned TestOracleIdentityLeavesOutTheCheckoutPath in internal/gate/oracle_identity_test.go | An oracle digest that also drops the path changes the gate cache key that FT449 closed |
| GR3 | 9 | An absent baseline phase manifest and a present empty one give two oracle identities | planned TestOracleIdentityCarriesTheBaselineSchedule in internal/gate/baseline_schedule_identity_test.go | An identity that reads both states as no schedule gives one identity |
| GR4 | 9 | Two baseline phase manifests with different bytes give two oracle identities | planned TestOracleIdentityCarriesTheBaselineSchedule in internal/gate/baseline_schedule_identity_test.go | An identity without the schedule gives one identity |
| GR5 | 10 | A baseline phase manifest that is a live symbolic link gives two oracle identities for two targets with different bytes | planned TestOracleIdentityFollowsTheLoaderRead in internal/gate/baseline_schedule_identity_test.go | An identity that frames the link and not the bytes that `readManifest` returns gives one identity for two schedules |
| GR6 | 10 | A dangling link and a FIFO at the baseline manifest give one oracle identity, and the subject build returns | planned TestOracleIdentityGivesADefectOneValue in internal/gate/baseline_schedule_identity_test.go | An identity that frames the defect text with its path gives two identities, and one that opens the FIFO blocks past the test deadline |
| GR7 | 5 | The capability posture changes the oracle identity | `internal/gate/oracle_identity_test.go` (`TestOracleIdentityCarriesTheCapabilityPosture`) | An identity without the posture stays equal across the two postures |
| GR8 | 3, 4, 11 | Each of tool bytes, policy, `PATH`, baseline runner, environment value, and declared path bytes changes the oracle identity | `internal/gate/oracle_identity_test.go` (`TestOracleIdentityAndOracleShareTheClosure`) | An identity that drops one of these inputs stays equal across its change |
| GR9 | 27 | The reader plan admits a stored table whose identity a sibling identity root computed over equal content | planned TestReaderPlanAdmitsATableOfASiblingCheckout in internal/gate/oracle_identity_test.go | A path-bound identity refuses the reader set with `new oracle identity` |
| GR45 | 8, 13 | With a `PATH` tool inside the checkout, an ordinary subject and a prospective subject of one tree from that checkout give one oracle identity | planned TestOracleIdentityKeepsTheOrdinaryAndProspectiveSubjectsEqual in internal/gate/oracle_identity_test.go | A spelling relative to the runtime root alone spells the tool absolutely for the prospective subject, and the identities split |
| GR10 | 1, 7 | After a green prospective run of one tree in a linked worktree, a prospective run of that tree from the primary checkout leaves the run count at 1 | planned TestAuthorizationKeyCrossesCheckouts in internal/gate/authorization_key_test.go | A key that frames the checkout path runs the gate again, and the count reads 2 |
| GR11 | 21 | That primary run prints exactly `gate: green (fresh verdict reused for this tree)` | planned TestAuthorizationKeyCrossesCheckouts in internal/gate/authorization_key_test.go | A changed or extra line breaks each reader of the reuse line |
| GR12 | 2 | Two trees under one subject name two evidence records | `internal/gate/prospective_owner_test.go` (`TestProspectiveEvidenceKeysToTheTreeAndBaselineRunnerIdentity`) | A key without the tree names one record |
| GR13 | 3 | A changed baseline runner identity names another evidence record | `internal/gate/prospective_owner_test.go` (`TestProspectiveEvidenceKeysToTheTreeAndBaselineRunnerIdentity`) | A key without the baseline runner names one record |
| GR14 | 5, 11 | After a green under the relaxed posture, a run with `BENCH_REQUIRE_CAPABILITIES=1` runs the gate, and the run count reads 2 | planned TestAuthorizationKeyBindsThePosture in internal/gate/authorization_key_test.go | A key without the posture reuses the relaxed green, and the count stays 1 |
| GR15 | 9, 11 | After a green in a linked worktree, a primary run whose baseline phase manifest differs runs the gate | planned TestAuthorizationKeyBindsTheBaselineSchedule in internal/gate/authorization_key_test.go | A key without the schedule reuses the green |
| GR16 | 6 | After an unpurposed green prospective run of the published tree, the complete checkpoint of that tree runs the gate, and the run count reads 2 | planned TestCompleteCheckpointIgnoresAnUnpurposedGreenOfItsTree in internal/gate/authorization_key_test.go | A key without the purpose frames reuses the unpurposed green of the same tree, and the count stays 1 |
| GR17 | 14 | A linked-worktree green whose recorded time is older than the freshness window does not answer the primary run | planned TestAuthorizationKeyKeepsTheSourceWindow in internal/gate/authorization_key_test.go | A window measured from the reuse, or no window, reuses the old green |
| GR18 | 15 | After a linked-worktree green, a red primary run of the same key removes the green, and the next linked-worktree run runs the gate | planned TestAuthorizationKeyRedRemovesAnotherCheckoutsGreen in internal/gate/authorization_key_test.go | A red that removes only its own checkout's record leaves the green, and the count does not grow |
| GR19 | 15 | A timed-out run removes the green under its key | `internal/gate/timeout_recovery_count_test.go` (`TestGateRunTimeoutInvalidatesOldEvidence`) | A timeout that keeps the record lets the next run reuse it |
| GR20 | 16 | An evidence record whose `authorization` field differs from its file name does not answer, and the run runs the gate | planned TestAuthorizationKeyRefusesAMisnamedRecord in internal/gate/authorization_key_test.go | A reader that trusts the file name reuses a copied record |
| GR21 | 22 | A delta checkpoint after an expired whole green of its base still runs the delta set | `internal/gate/delta_checkpoint_test.go` (`TestDeltaCheckpointReadsTheBaseWithoutExpiry`) | A base read under another key finds no green, and the checkpoint runs the whole gate |
| GR22 | 17 | A primary run with no gate cache that reuses a linked-worktree green writes the cache with the tree, `green`, and the recorded time of the worktree run | planned TestReusedGreenRecordsTheGateCache in internal/gate/reused_green_test.go | A reuse that writes nothing, or that takes the held lock again and skips the write, leaves the cache absent |
| GR47 | 17 | A primary run whose gate cache holds a red verdict of the same tree replaces it with the reused green and the worktree run's recorded time | planned TestReusedGreenReplacesARedCacheOfItsTree in internal/gate/reused_green_test.go | A write rule that compares only the tree leaves the red cache in place |
| GR23 | 18 | A reuse that the gate cache already describes leaves the cache bytes unchanged | `internal/gate/run_outcomes_test.go` (`TestGateRunGreenReusesRecordedVerdict`) | A reuse that rewrites the record moves its recorded time |
| GR24 | 19 | A reuse whose gate-cache write fails returns an operational result with no reusable green | planned TestReusedGreenCacheWriteFailureRefuses in internal/gate/durable_replacement_test.go | A reuse that ignores the failed write reports green beside a stale cache |
| GR44 | 20 | A reuse while another run holds the checkout's gate lock, with a cache of another tree, prints the reuse line and leaves the cache bytes unchanged | planned TestReusedGreenBesideAHeldLockWritesNothing in internal/gate/reused_green_test.go | A reuse that waits for or refuses on the lock turns a green into a refusal |
| GR25 | 23 | In a fixture with no lane, after `bench worktree merge --from main` ran the whole gate, a spec-less landing of the merge tip leaves the gate tally at one run | planned TestLandReusesTheMergeGreen in internal/worktree/land_reuse_test.go | A key that frames the checkout path runs the landing gate again, and the tally holds two runs |
| GR26 | 24 | That landing advances the green marker to its published commit | planned TestLandReusesTheMergeGreen in internal/worktree/land_reuse_test.go | A publication that reads the reuse as another kind leaves the marker at the old commit |
| GR27 | 25 | After a green `bench gate --checkpoint specs/x/spec.md --complete` in the source worktree, the spec landing leaves the gate tally at one run | planned TestLandReusesTheCompleteCheckpointGreen in internal/worktree/land_reuse_test.go | A key that frames the checkout path runs the landing gate again, and the tally holds two runs |
| GR28 | 26 | When the source changes a file that the declared build inputs name, the spec landing after the complete checkpoint leaves the gate tally at two runs | planned TestLandRunsTheGateForAnotherBaselineRunner in internal/worktree/land_reuse_test.go | A key without the baseline runner reuses the checkpoint green, and the tally holds one run |
| GR29 | 28 | ADR 31 states the key members, the reuse across checkouts, the posture binding, and the trust posture | review-owned at the GR-C3 checkpoint | A missing ADR leaves the decision only in this spec, which retires at landing |
| GR30 | 28 | ADR 2, ADR 16, the profile, the reference, and the glossary state the authorization key as the reuse rule | review-owned at the GR-C3 checkpoint | A stale sentence tells a teammate that a green never crosses checkouts |
| GR31 | 30 | `bench worktree merge --plan` of a current target that declares a lane prints `grade=none` | planned TestMergePlanNamesEachGrade in internal/worktree/merge_plan_test.go | A plan that tests the lane before the kind names a cost that the merge never pays |
| GR32 | 30 | The plan of a non-current target with a declared lane prints `grade=lane` | planned TestMergePlanNamesEachGrade in internal/worktree/merge_plan_test.go | A plan that ignores the lane names the whole gate |
| GR33 | 30 | The plan of a target with no lane and a fresh green of the composed tree prints `grade=reuse` | planned TestMergePlanNamesEachGrade in internal/worktree/merge_plan_test.go | A plan that skips the reuse predicate names the whole gate |
| GR34 | 30 | The plan of a target with no lane and no such green prints `grade=whole` | planned TestMergePlanNamesEachGrade in internal/worktree/merge_plan_test.go | A plan that assumes a reuse hides the full-gate cost |
| GR35 | 29 | The plan prints exactly one `merge_plan{worktree=...,from=...,kind=...,previous_tip=...,tree=...,grade=...}` line at exit 0 | planned TestMergePlanPrintsOneRecord in internal/worktree/merge_plan_test.go | A plan that prints the `merged{...}` record or another shape breaks the projection's readers |
| GR36 | 31 | The plan leaves the target branch tip, the checkout fingerprint, the `git worktree list` registrations, the gate cache bytes, and the evidence store names unchanged | planned TestMergePlanChangesNothing in internal/worktree/merge_plan_test.go | A plan that runs the merge or its grade, or leaves its private checkout registered, changes one of these facts |
| GR37 | 32 | After a plan with `grade=lane`, the merge prints a lane pass line and writes no whole-gate tally | planned TestMergeRunsThePlannedGrade in internal/worktree/merge_plan_test.go | A plan whose lane decision differs from the merge's names the wrong schedule |
| GR38 | 32 | After a plan with `grade=reuse`, the merge prints the reuse line and adds no whole-gate tally | planned TestMergeRunsThePlannedGrade in internal/worktree/merge_plan_test.go | A plan whose reuse predicate differs from the run's names a reuse that the run does not take |
| GR39 | 32 | After a plan with `grade=whole`, the merge adds one whole-gate tally | planned TestMergeRunsThePlannedGrade in internal/worktree/merge_plan_test.go | A plan that names `whole` while the run reuses misstates the cost |
| GR40 | 32 | After a plan with `grade=none`, the merge reports the `current` kind and adds no whole-gate tally | planned TestMergeRunsThePlannedGrade in internal/worktree/merge_plan_test.go | A plan whose kind decision differs from the merge's names a cost that the merge never pays |
| GR46 | 32 | A no-lane target whose gate cache holds an interrupted pending record, with a fresh green of the composed tree, plans `grade=whole`, and the merge then adds one whole-gate tally | planned TestMergePlanProjectsAPendingCacheAsWhole in internal/worktree/merge_plan_test.go | A plan that reads the store alone through `InspectTreeContext` projects `reuse` while the run refuses the reuse |
| GR41 | 33 | A plan of a conflicted composition refuses with the `merge-conflict` face, and its `next=` route ends with a rerun that holds `--plan` | planned TestMergePlanRefusalKeepsThePlan in internal/worktree/merge_plan_test.go | A route that names the bare merge rerun runs the merge instead of the plan |
| GR42 | 34 | `bench help` prints the merge row with `[--plan]` from the one merge grammar | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`) | A hand-written row without `[--plan]` fails the pinned help line |
| GR43 | 35 | The reference's Command Notes, the profile's merge paragraph, and the glossary state the plan form and its four grades | review-owned at the GR-C3 checkpoint | A stale merge paragraph hides the disclosure from a teammate |

Not covered: story 36 — FT450's landing queue owns the landing of several sources under one gate, and it reuses the key that GR10 grades.

### Edge inventory

The hostile-input walk of the profile checklist gives these dispositions:

- Paths with spaces: GR1 puts a space in one identity root. The plan prints hex identities and fixed words, so no path reaches its record.
- Control bytes in git-sourced text: the plan record holds the assignment id, commit and tree hashes, and fixed words. No free text reaches it.
- A live symbolic link where a file is expected: GR5 frames the bytes that the loader reads through the link, which are the bytes that the run uses.
- A dangling symbolic link and a special file in a discovered path: GR6 frames one stable defect value, and `readManifest` stats before it opens.
- Absent versus present-but-empty: GR3 gives the two manifest states two identities.
- A command whose own write changes a fact it reports: GR36 proves that the plan changes no reported fact. GR22 and GR47 prove that a reuse writes the cache that `bench status` reports.
- State serialized by one process and reloaded by a fresh one: GR10 and GR25 build each subject from disk in a separate run.
- A temporary root with a symbolic-link component: the subject build resolves each root to its canonical spelling first, so the relative spelling holds.
- Concurrent writers: the writer sweep under Further notes names each writer of one key, and GR44 grades a reuse beside a held lock.

**Won't handle:** a landing plan form — the merge plan serves the coordinator before a fold, and each landing still runs or reuses under its own lease.
**Won't handle:** several light-path sources landed under one gate — each single landing still reuses its own source's green under GR10, and FT450 owns the stack.
**Won't handle:** a prune of the evidence store — each reader still reads one `0600` record for its key, and the store already keeps one record for each subject today.
**Won't handle:** two checkouts that grade one key at the same time — each run still holds its own gate lock, and the store keeps the outcome of the real run that finishes last.
**Won't handle:** the run-binary identity of an installed kit in a linked repository with no declared build inputs — the same-checkout retry there keeps its `unbound` runner value, as today.
**Won't handle:** a reuse line that names the source run — every reader still matches the line `gate: green (fresh verdict reused for this tree)` exactly.
**Won't handle:** guidance that tells a coordinator to plan before each fold — each fold still runs as written, and the help row names the plan form.
**Won't handle:** an edit of the baseline phase manifest while a run is in progress — the settle check still compares the tree and the oracle digest, and the profile forbids a repository change during a gate run.

## Out of scope

- A landing queue for several sources under one gate is FT450's capability. Estimate: 14 edits, 4 gate runs, after its own shaping session.
- A prune of the evidence store by age is a separate maintenance capability. Estimate: 3 edits, 1 gate run.
- A plan form of `bench worktree land` is a separate disclosure. Estimate: 6 edits, 2 gate runs.
- Coordinator guidance to plan before a fold is a separate guidance change. Estimate: 2 edits, 1 gate run.
- A run-binary member for an installed kit in a linked repository is a separate key widening. Estimate: 4 edits, 1 gate run.

## Further notes

### Late clarifications

The reviewer closed two forks on 2026-10-26 in this phase's conversation:

- The key design is one key for all runs. The prospective evidence store keys every green by the authorization key, a posture change now reruns, and the reader-table identity becomes path-free.
- The disclosure is a `--plan` form of `bench worktree merge`. It prints the kind and the projected grade, and it changes no ref, checkout, or record.

### Flagged additions

- The gate-cache write on a reuse (GR22, GR23, GR24, GR44, GR47) keeps `bench status` and the green-marker bootstrap on the reused green. The decision source does not name it.
- The baseline phase schedule as a gate-version member (GR3 to GR6, GR15) is the input that selects the phases of a prospective run.
- The spelling of an executable relative to the runtime root or the identity root (GR1, GR45) keeps the ordinary and prospective subjects of one checkout equal.
- The `authorization` field of the evidence record (GR20) makes a misnamed record fail closed.
- The equivalence of an ordinary green and a prospective green now crosses checkouts. It follows from the one-key decision and ADR 2's posture 5.
- The derived merge help row (GR42) removes a hand copy that the grammar change would leave stale.
- The `previous_tip` and `tree` cells of the plan record let a coordinator compare the plan with a later landing.
- ADR 31 and the changelog entries record the decided state for a teammate, as invariant 3 requires.
- The reversal of FT449's posture rule and of the identity-root case of `TestOracleIdentityAndOracleShareTheClosure` is reviewer-approved, and the build edits that case.

### Source-to-row table

| reviewed source clause | occurrence | rows or disposition |
|---|---|---|
| An unchanged landing composition reuses gate evidence only under a complete authorization key | FT314 title | GR10, GR25, GR27 |
| A cross-run reuse would skip a duplicate full gate for an unchanged composition | FT314 paragraph 1 | GR25, GR27 |
| The authorization key binds the composed tree | FT314 paragraph 2 | GR12 |
| The key binds the gate version | FT314 paragraph 2 | GR8, GR13, GR15, GR28 |
| The key binds the toolchain | FT314 paragraph 2 | GR8 |
| The key binds the capability posture | FT314 paragraph 2 | GR7, GR14 |
| The key covers a landing that reuses its `bench worktree merge` verdict | FT314 paragraph 2 | GR25, GR26 |
| A change to any key member runs the full gate | FT314 paragraph 2 | GR8, GR14, GR15, GR16, GR28 |
| The design must preserve the gate as the oracle | FT314 paragraph 3 | GR17, GR18, GR19, GR24 |
| It must not treat an old receipt as authority for a different execution environment | FT314 paragraph 3 | GR14, GR15, GR17, GR28 |
| A review-venue merge discloses its full-gate cost before execution | FT314 paragraph 4 | GR31, GR32, GR33, GR34, GR35, GR36 |
| The projection must preserve the gate schedule | FT314 paragraph 4 | GR37, GR38, GR39, GR40, GR46 |
| A merge and the landing that follows grade the same composition with two full gates | FT314 paragraph 5 | GR25 |
| A drain that lands several light-path sources could land them in order under one gate | FT314 paragraph 5 | Won't handle, story 36 |
| One key for all runs, the posture reversal, and a path-free reader-table identity | late clarification 1 | GR1, GR9, GR14 |
| A `--plan` form of the merge | late clarification 2 | GR31 to GR41, GR46 |

### Reader and writer sweep

The whole hidden-tree sweep excluded only `.git`. Its needles were `bench-gate-evidence`, `evidenceName`, `evidencePath`, `reusableEvidence`, `retainGreen`, `invalidateEvidence`, and `inspectEvidence`.
It also searched the reuse line, `worktree merge --from`, `WorktreeMerge`, `BENCH_REQUIRE_CAPABILITIES`, `readManifest`, and `Identity`.

The writers of the evidence store are `retainGreen` and `invalidateEvidence`. The second serves the red and timeout steps of `executeSubjectOwning` and `retainedFailure`.
The readers are `reusableEvidence` at two sites of `executeSubjectOwning`, `inspectProspective` through `InspectTreeContext` and `deltaPlan`, and `ValidateProjectGreen`.
`authorization.AuthorizeWithWriters` reads `InspectTreeContext` after each run, and `inheritedSubjectGreen` reads `ValidateProjectGreen`. Every site reaches the store through `evidencePath`.

The readers of the oracle identity are `readerTrace.table` and the reader plan's identity compare. No other production code reads it.
The readers of the baseline phase manifest are `phaseTable`, `LaneFor`, and `narrowsKitTable`, each through `readManifest`. The identity becomes a fourth caller of that one read.
The system test `TestProcessLifetimeNestedGate` lists the store names before and after a refusal, so a renamed record keeps it green.

The `.github` workflows set `BENCH_REQUIRE_CAPABILITIES=1` in `native-runtime.yml` and `release.yml`, and they call no merge. A fresh CI checkout holds no green, so the posture binding changes no CI run.
No anchor of the registry pins a sentence that the docs change: `bench anchors` returned no evidence-key needle for the profile, the reference, ADR 2, or the glossary.

### Ticket sizing

Each ticket stays near 150k tokens of author context. The estimate takes the byte total of the ticket's `Writes:` files at about four bytes per token.
It adds about 50k tokens of fixed context for the agreement, the skills, the spec rows, and the ticket. It adds about 20k to 30k tokens for the verification runs of its packages and its probes.
Each estimate is an upper bound that rounds up.

| ticket | `Writes:` bytes | estimate |
|---|---|---|
| 01-leave-the-checkout-out-of-the-identity.md | about 127 KB, 74 KB of it `.bench/commitment.json` | about 105k tokens |
| 02-key-the-evidence-store.md | about 120 KB | about 110k tokens |
| 03-record-a-reused-green.md | about 120 KB | about 110k tokens |
| 04-plan-a-merge-grade.md | about 160 KB | about 120k tokens |
| 05-record-the-reuse-rule.md | about 245 KB | about 130k tokens |

`.bench/commitment.json` joins ticket 01, because the spec-stage close commits the binding and no author reads that file.
Ticket 02 keeps the signature of `evidencePath`. Five gate test files call it for the store directory alone, so they stay out of its `Writes:` line.

### Pre-review proof checklist

- Cited symbols: `buildSubjectForTree`, `hashProspectivePreparation`, `baselineRunnerIdentity`, `identityPath`, `phasePolicyRoot`, `readManifest`, `applyCheckpoint`, `evidenceName`, `evidencePath`, `inspectEvidenceWindowed`, `reusableEvidence`, `retainGreen`, `invalidateEvidence`, `mergeOwner`, `Owner.Merge`, and `worktreeSuffix` resolve in the sources named above.
- Import edges: none new. `internal/worktree` already imports `internal/gate` and `internal/landing`, and `cmd/bench` already imports `internal/usage`.
- Source-row clauses and occurrences: the source-to-row table quotes each FT314 clause and both late clarifications.
- Promised field labels: the plan record cells are `worktree`, `from`, `kind`, `previous_tip`, `tree`, and `grade`. The grades are `none`, `lane`, `reuse`, and `whole`. The evidence record field is `authorization`.
- Changed-function callers: `bench consumers` listed `evidenceName` (2 production sites, 1 test), `reusableEvidence` (2), `retainGreen` (1, 1 test), `invalidateEvidence` (3), and `inspectEvidenceWindowed` (2). It also listed `buildSubjectForTree` (2, 2 tests), `Owner.Merge` (2, 9 tests), and `mergeRerun` (1, 1 test).
- Copy survival: GR42 fails while the hand-written merge help suffix survives. GR5 fails when the identity reads the manifest by a second route.
- Rendered-shape readers: the merge grammar line is pinned in `cmd/bench/help_inventory_test.go` and derived in `internal/worktree/verb_runner_check_test.go`. The reuse line keeps its bytes, and its four readers stay green.
- Pin operators: GR10, GR14, GR16, GR18, GR25, GR27, and GR28 compare run counts with `!=`. GR11 and GR35 compare the exact output. GR36 and GR44 compare bytes and name lists for equality.
- Entry reads: the identity reads the baseline manifest through `buildSubjectForTree`, whose identity-root parameter is the internal form. GR3 to GR6 grade it.
- Derived expectations: the identity rows compare two computed digests. The run-count rows count fixture runs. No expectation comes from the code under test.
- Consolidated rules: the consolidation table maps each changed rule to its rows.
- Quantified obligations: every store reader and writer above routes through one name function, and GR10 with GR20 grades the named record.
- Workflow-step writes: none. The plan is a read-only verb form, and no workflow step changes.

### Bootstrap authority trace

A reused green replaces an execution, so the trust section above traces each hop. The wrapper selects the Bench executable, and the freshness check binds it to its sources.
That executable computes the key and reads the store. The store record carries no proof of its writer. In a root that ships the prospective script, the key's baseline-runner member is the independent binding.
The reviewer-visible trust assumption is ADR 2's accidental-weakening model, and GR28 proves that another baseline runner never answers.

# Cut subsumed tests on probe proof

Status: staged

Roadmap: FT411

Decision source: Named reviewed artifact: the bound outcome criteria `FT411.cut` and `FT411.probe` of outcome FT411 in `.bench/commitment.json` (reviewer decision 2026-10-18, drain d-359accf567c1 ticket triage). Under the reviewer decision of 2026-10-29 that `roadmap/FT437.md` records, these criteria win over the text of `roadmap/FT411.md`.

Verification log: pending — the spec stage writes it at close.

## Problem

The 2026-10-09 test-efficiency reviews listed 27 table rows of tests that add no value.
Fifteen rows are in `capture/test-efficiency-review-2026-10-09.md` and twelve rows are in `capture/test-efficiency-review-worktree-2026-10-09.md`.
Each row claims that a kept test already catches what the row's test catches.
Neither review ran a probe for any claim, so a cut on a review claim alone can remove the only test that catches a fault.

The gate pays for each subsumed test on every run.
The worktree package also pins its exact top-level test count at 880 (`internal/worktree/parallel_census_test.go:926`), so a careless cut turns the gate red.

## Solution

Cut each review row whose probe proves a kept test red on the same mutation.
A probe is a `bench probe` run that applies one production mutation and expects one named test to fail.
The spec stage ran one probe set for each row at commit `7b0a83e44`, and the probe results decide each row.

Where a kept test stayed silent and the review names a move, the build moves the missing predicate into the kept test first.
Where the review names a merge, one test keeps every predicate of the tests it replaces.
A row whose probe stays silent keeps its test.
Each chunk checkpoint grades one bit probe for each cut predicate on the kept test at the chunk tip.

## User stories

Line: sonnet / high.
Implementation-line reason: The hardest material chunk is FT411-C2, because its merges and moves edit shared landing fixtures. Spec precision is exact: each row names its candidate, its kept test, and its mutation. The seams are existing tests, so seam uncertainty is low. Test strength is high, because each cut predicate is a checkpoint probe and the worktree count pin reds a miscount.
Harder chunks: FT411-C2.

### Cuts that keep coverage

1. As a maintainer, I want each subsumed test cut, so that the gate spends no time on a test that adds no failure detection.
2. As a reviewer, I want a probe on the kept test for each cut predicate, so that a cut never removes the only test that catches a fault.
3. As a reviewer, I want each probe to use a mutation that the cut test catches, so that the kept test proves the same coverage and not an easier one.
4. As a reviewer, I want a moved predicate to land in the kept test before its source test goes, so that a kept test that missed a mutation catches it after the move.
5. As a reviewer, I want each merged test to keep every predicate of the tests it replaces, so that a merge removes only repeated arrangement.
6. As a maintainer, I want a row whose probe stays silent to keep its test, so that the cut rule never trades coverage for time.

### Adoption and repair suites

7. As a maintainer, I want the legacy doctor equality test cut, so that no test compares a function with the function that it calls.
8. As a maintainer, I want the undo-of-created test cut, so that one undo test covers the created path and the restored path.
9. As a maintainer, I want the managed repair test cut, so that the undo test alone grades the repaired bytes.
10. As a maintainer, I want the three recovery boundary tests merged into one test with one linked session, so that one repair run proves all three boundaries.
11. As a maintainer, I want the consumer-only scope substring test cut, so that the behavior test of the generated gate holds that scope.
12. As a maintainer, I want the seeded gate inputs test cut after its exact-byte check moves, so that the repeated-setup test also proves the seeded bytes.
13. As a maintainer, I want the immutable inventories test cut after its three unique rows move, so that one decision table holds every legacy refusal.
14. As a maintainer, I want the topic-push silence test cut, so that the legacy-pin test alone proves a silent topic push.
15. As a maintainer, I want seven Upgrade preview e2e rows cut, so that the pure plan-kinds table holds the count partitions.
16. As a maintainer, I want the four remaining Upgrade preview e2e rows kept, so that the real Link path still proves the preview from end to end.
17. As a reviewer, I want each staged spec that cites a cut test to cite its kept test, so that the coverage check stays green.
18. As a maintainer, I want the decision-domain registry to drop a deleted decision test file, so that the census describes the tree.

### Landing suites

19. As a maintainer, I want the registry completeness test cut, so that the follow walk alone proves that each landing face has a fixture.
20. As a maintainer, I want the zero-census test cut, so that the no-heads test alone proves a zero count.
21. As a maintainer, I want the effects-order test cut after its tail-order check moves, so that the refresh test proves the order of the records.
22. As a maintainer, I want the slug row of the spec-form resume test cut, so that only the path row adds a resume.
23. As a reviewer, I want the spec's reviewer face list moved into the follow walk, so that an authority change in the registry still turns a test red.
24. As a maintainer, I want the single-owner test cut after its release check moves, so that the promotion test proves the released checkout.
25. As a maintainer, I want the three census landings merged into one landing, so that one landing proves the count, the heads, the output, and the drop.
26. As a maintainer, I want the PL30 clean row cut, so that the published-identity test alone proves a clean resume over a moved destination.
27. As a maintainer, I want the red-route test cut, so that the follow walk and the lease test hold its three predicates.
28. As a reviewer, I want each whole-test cut in `internal/worktree` to lower the count pin in the same ticket, so that the pin stays exact.

### Rows that stay

29. As a maintainer, I want the detected-project gate system test kept, so that the only test of the detected-branch hygiene line stays.
30. As a maintainer, I want the Link partition e2e rows kept, so that the on-disk outcomes that the pure table cannot see stay graded.
31. As a reviewer, I want review rows that cut no test kept out of this spec, so that fixture sharing gets its own measured spec.

## Implementation decisions

### The cut rule

A review row is a cut when the spec-stage probe set proves a kept test red on each mutation that the cut predicate catches.
A move precedes a cut only where a spec-stage probe found the kept test silent. T1-4 and WT-11 therefore cut without the move or the fold that their review rows propose, because their kept tests already bit.
A merge keeps one test and removes the others. The kept test carries every predicate of the removed tests.
A row whose spec-stage probe found every kept test silent keeps its test. This spec does not add a replacement test for such a row.

A row ID is `T1-<n>` for row n of the first review's Table 1, and `WT-<n>` for row n of the worktree review's Table 1.

### Row dispositions

| row | candidate (spec-stage location) | disposition | kept test |
|---|---|---|---|
| T1-1 | four socket rows | delivered before the review commit by `2b99970ed` and `c34b8644e` | none |
| T1-2 | `TestCompatibilityLegacyDoctor` (`internal/adopt/compatibility_test.go:101`) | cut | `TestDoctorParsesHelpAndUnknownFlags` |
| T1-3 | `TestCompatibilityUndoCreated` (`internal/adopt/repairtest/repair_test.go:76`) | cut | `TestCompatibilityUndo` |
| T1-4 | `TestCompatibilityManagedRepair` (`repair_test.go:18`) | cut | `TestCompatibilityUndo` |
| T1-5 | the Policy, Interruption, and Runtime boundary tests (`repair_test.go:217`, `:243`, `:271`) | merge | `TestCompatibilityRecoveryBoundaries` (new) |
| T1-6 | `TestSetupGateDeclaresConsumerOnlyValidityScope` (`internal/adopt/setup_test.go:159`) | cut | `TestSetupGateRejectsIgnoredDeclaredInputs` |
| T1-7 | `TestSetupGateRejectsIgnoredDeclaredInputs` | no cut: fixture sharing | none |
| T1-8 | `TestDetectedProjectGateRejectsIgnoredDeclaredInput` (`internal/systemtest/adoption_test.go:18`) | stays: probe silent | none |
| T1-9 | `TestSetupSeedsGateInputs` (`setup_test.go:53`) | move, then cut | `TestSetupTwiceLeavesGateInputsIdentical`, `TestSetupSeedsGateInputsInZeroSignalRepo` |
| T1-10 | `TestPlanLifecycleFromImmutableInventories` (`internal/adopt/decision_test.go:8`) | move, then cut | `TestCompleteLifecycleDecisions`, `TestLifecycleDecisionOrder` |
| T1-11 | `TestPrePushHookTopicPushIsSilent` (`internal/adopt/link_hook_test.go:332`) | cut | `TestPrePushHookIgnoresALegacyPin` |
| T1-12 | seven e2e rows of `TestLifecycleUpgradePreviewCompatibility` (`internal/adopt/lifecycle_preview_test.go:93-109`) | cut | its `plan kinds` subtest |
| T1-13 | the Link, Dropped, and Claude partition e2e tables | stays: probe silent | none |
| T1-14 | `TestRootConformanceCheckPrintsFailures` | no cut: run scoping | none |
| T1-15 | pasted git helpers | delivered by `2d44df8b2` and `c3f718143` (FT360) | none |
| WT-1 | `TestLandingRefusalRegistryHasAProducingFixture` (`internal/worktree/identity_component_test.go:217`) | cut | `TestLandingFacesFollowTheirRoutes` |
| WT-2 | `TestLandCommandStatesZeroForAnAssignmentWithNoRecords` (`land_census_test.go:55`) | cut | `TestLandCommandPrintsNoHeadsLineWithoutRecords` |
| WT-3 | `TestLandEffectsRowPrecedesTheLandedRecord` (`land_effects_test.go:337`) | move, then cut | `TestLandSkipsTheRefreshWithoutBuildInputs` |
| WT-4 | the `x` row of `TestResumeLandCommandAcceptsSpecSlugAndPath` (`land_resume_test.go:208`) | cut | `TestResumeLandCommandCompletesAnInterruptedMarker` |
| WT-5 | `TestReviewerLandFacesOpenWithTheMarker` (`refusal_route_test.go:172`) | move, then cut | `TestLandingFacesFollowTheirRoutes` |
| WT-6 | `TestLandCommandKeepsOneOwnerProcessThroughPublicationAndRelease` (`land_freshness_test.go:59`) | move, then cut | `TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion` |
| WT-7 | `deferredMergedDetails` (`identity_component_test.go:256`) | no cut: dead allowance, no probe applies | none |
| WT-8 | the census landings at `land_census_test.go:37`, `:68`, and `land_census_output_test.go:13` | merge | `TestLandCommandStatesTheCensusCountAndDropsTheRecords` |
| WT-9 | the two rows of `TestResumeLandCommandPreauthenticatesCompletedRequestAndPath` | no cut: fixture sharing | none |
| WT-10 | the `PL30/clean` row of `TestResumeLandCommandPublicRefusesDestructiveDestinationState` (`land_resume_refusal_test.go:19`) | cut | `TestResumeLandCommandPublicBindsPublishedLandingIdentity` |
| WT-11 | `TestLandingRedRouteNamesTheRepair` (`refusal_route_test.go:141`) | cut | `TestLandReleasesTheLeaseOnRefusal`, `TestLandingFacesFollowTheirRoutes` |
| WT-12 | pasted identity and landing literals | identity delivered by FT360; the other literals cut no test | none |

The review's system-suite note S1 is not a table row. The twelve copy-based lifecycle tests are not table rows either. The Out of scope section prices both.

### Moves

- T1-9: `TestSetupTwiceLeavesGateInputsIdentical` compares its seeded file with `wantSeededGateInputs`. The comparison runs on its `go.mod` repository.
- T1-10: `TestCompleteLifecycleDecisions` gains three rows over the legacy inventory path: a duplicated current asset and a desired asset with no fingerprint each refuse, and an empty input plans no operation and no refusal.
- WT-3: `TestLandSkipsTheRefreshWithoutBuildInputs` asserts that the effects table is the tail before the landed record, and that the landed record is the last stdout line.
- WT-5: `TestLandingFacesFollowTheirRoutes` holds the spec's list of the seven reviewer land faces as an independent list. Each listed face's printed route opens with the reviewer marker. The list stays independent of the registry, because the walk takes its authority from the same registry that renders the route.
- WT-6: `TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion` asserts that the released checkout is absent after the landing.

The WT-5 move turns the AGENTS.md independent-expectation exception into a recorded red. Its checkpoint probe changes one listed face's authority in the registry, and the follow walk then fails.

### Merges

- T1-5: `TestCompatibilityRecoveryBoundaries` runs one linked session and one `doctor --fix`. It keeps each planted policy file, the supervised process, and the runtime directory check, and it asserts all three decision rows.
- WT-8: `TestLandCommandStatesTheCensusCountAndDropsTheRecords` seeds the head records and the output records before one landing. It asserts the count suffix, the dropped record file, the surviving handoff sections, the head breakdown, and the output breakdown. The count stays the raw-call count, so a count that includes output records turns it red.

### Count pin and readers

- Each ticket that removes a top-level test from `internal/worktree` lowers `worktreeTestCount` by its own removal count, from the value at its base. `TestPackageTestCountPin` grades the result, and its comparison is exact.
- No candidate is a serial test, so `worktreeSerialCeiling` does not change.
- The ticket that cuts T1-3, T1-4, and T1-5 also re-points the seam cells of rows LTE4, LTE9, LTE12, LTE13, and LTE14 in `specs/landing-test-efficiency/spec.md`. LTE4 and LTE9 cite `TestCompatibilityUndo`, and LTE12 to LTE14 cite `TestCompatibilityRecoveryBoundaries`. The `docs-currency-workflow` check runs `bench coverage --check` over every spec, so a stale cell turns the gate red.
- The T1-10 cut empties `internal/adopt/decision_test.go`, so the ticket deletes the file and removes its key from `directArchitectureTests` in `internal/conformance/ordinary_build_census_test.go`. The census reds a key whose file is absent.
- The WT-8 merge empties `internal/worktree/land_census_output_test.go`, so the ticket deletes that file.
- The comment at `internal/worktree/land_refusal_fixture_test.go:36` names the WT-1 candidate. The WT-1 ticket points that comment at the follow walk.
- A helper that only a cut test uses goes with the cut.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| FT411-C1 / 01-cut-subsumed-repair-tests.md, 02-cut-subsumed-adopt-tests.md | The adoption and repair suites lose their subsumed tests, and each kept test bites on every cut predicate | TC1–TC24, TC41, TC42 | TestDoctorParsesHelpAndUnknownFlags, TestCompatibilityUndo, TestCompatibilityRecoveryBoundaries, TestSetupGateRejectsIgnoredDeclaredInputs, TestSetupTwiceLeavesGateInputsIdentical, TestCompleteLifecycleDecisions, TestLifecycleDecisionOrder, TestPrePushHookIgnoresALegacyPin, TestLifecycleUpgradePreviewCompatibility, TestDetectedProjectGateRejectsIgnoredDeclaredInput, TestLifecycleLinkPartitions | no |
| FT411-C2 / 03-cut-subsumed-landing-tests.md, 04-cut-subsumed-refusal-route-tests.md | The landing suites lose their subsumed tests, each kept test bites on every cut predicate, and the count pin stays exact | TC25–TC40 | TestLandingFacesFollowTheirRoutes, TestLandCommandPrintsNoHeadsLineWithoutRecords, TestLandSkipsTheRefreshWithoutBuildInputs, TestResumeLandCommandCompletesAnInterruptedMarker, TestResumeLandCommandAcceptsSpecSlugAndPath, TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion, TestLandCommandStatesTheCensusCountAndDropsTheRecords, TestResumeLandCommandPublicBindsPublishedLandingIdentity, TestLandReleasesTheLeaseOnRefusal, TestPackageTestCountPin | yes |

### Completion plan

```bench-completion-plan
{
  "version": 1,
  "chunks": [
    {
      "id": "FT411-C1",
      "tickets": [
        "01-cut-subsumed-repair-tests.md",
        "02-cut-subsumed-adopt-tests.md"
      ],
      "verification": [
        {"id": "adopt-suites", "command": "bench test --package ./internal/adopt --package ./internal/adopt/repairtest"},
        {"id": "census", "command": "bench test --check ordinary-build-census"},
        {"id": "spec-citations", "command": "bench test --check docs-currency-workflow"},
        {"id": "tc1", "command": "bench test --package ./internal/adopt --run '^TestDoctorParsesHelpAndUnknownFlags$'", "probe": "internal/adopt/compatibility.go: swap `args[0] == \"--compat\"` for `args[0] != \"--compat\"`; expect TestDoctorParsesHelpAndUnknownFlags"},
        {"id": "tc2", "command": "bench test --package ./internal/adopt/repairtest --run '^TestCompatibilityUndo$'", "probe": "internal/adopt/transaction/publication.go: swap `remove = os.Remove` for `remove = func(string) error { return nil }`; expect TestCompatibilityUndo"},
        {"id": "tc3", "command": "bench test --package ./internal/adopt/repairtest --run '^TestCompatibilityUndo$'", "probe": "internal/adopt/transaction/publication.go: swap `_, err = f.Write(value.Data)` for a write of a copy whose first byte has bit 0x20 flipped; expect TestCompatibilityUndo"},
        {"id": "tc4", "command": "bench test --package ./internal/adopt/repairtest --run '^TestCompatibilityRecoveryBoundaries$'", "probe": "internal/compatibility/recovery.go: swap `{\"security-policy\", \"decision\",` for `{\"security-policy\", \"info\",`; expect TestCompatibilityRecoveryBoundaries"},
        {"id": "tc5", "command": "bench test --package ./internal/adopt/repairtest --run '^TestCompatibilityRecoveryBoundaries$'", "probe": "internal/compatibility/recovery.go: swap `{\"process-interruption\", \"decision\",` for `{\"process-interruption\", \"info\",`; expect TestCompatibilityRecoveryBoundaries"},
        {"id": "tc6", "command": "bench test --package ./internal/adopt/repairtest --run '^TestCompatibilityRecoveryBoundaries$'", "probe": "internal/compatibility/recovery.go: swap `{\"private-runtime\", \"decision\",` for `{\"private-runtime\", \"info\",`; expect TestCompatibilityRecoveryBoundaries"},
        {"id": "tc7", "command": "bench test --package ./internal/adopt/repairtest --run '^TestCompatibilityRecoveryBoundaries$'", "probe": "internal/adopt/compatibility.go: insert a write of `x` to `$CODEX_HOME/config.toml` before `namespace := filepath.Join(home, \"compatibility-repairs\")`; expect TestCompatibilityRecoveryBoundaries"},
        {"id": "tc10", "command": "bench test --package ./internal/adopt --run '^TestSetupGateRejectsIgnoredDeclaredInputs$'", "probe": "internal/adopt/setup.go: omit `BENCH_CONFORMANCE_CONSUMER_ONLY=1 `; expect TestSetupGateRejectsIgnoredDeclaredInputs"},
        {"id": "tc11", "command": "bench test --package ./internal/adopt --run '^TestSetupSeedsGateInputsInZeroSignalRepo$'", "probe": "internal/adopt/setup.go: omit `\"BENCH_RUN_BINARY\", `; expect TestSetupSeedsGateInputsInZeroSignalRepo"},
        {"id": "tc12", "command": "bench test --package ./internal/adopt --run '^TestSetupTwiceLeavesGateInputsIdentical$'", "probe": "internal/adopt/setup.go: swap `content: scaffoldGateInputs()` for `content: scaffoldGateInputs() + map[bool]string{false: \" \"}[facts.zeroSignal]`; expect TestSetupTwiceLeavesGateInputsIdentical"},
        {"id": "tc13", "command": "bench test --package ./internal/adopt --run '^TestCompleteLifecycleDecisions$'", "probe": "internal/adopt/decision.go: swap `; exists {` for `; exists && false {`; expect TestCompleteLifecycleDecisions"},
        {"id": "tc14", "command": "bench test --package ./internal/adopt --run '^TestCompleteLifecycleDecisions$'", "probe": "internal/adopt/decision.go: omit ` || asset.Fingerprint == \"\"`; expect TestCompleteLifecycleDecisions"},
        {"id": "tc15", "command": "bench test --package ./internal/adopt --run '^TestCompleteLifecycleDecisions$'", "probe": "internal/adopt/decision.go: insert an `empty` refusal for an input with no current and no desired asset before `decisions := make([]LifecycleDecision, 0, len(current)+len(desired))`; expect TestCompleteLifecycleDecisions"},
        {"id": "tc16", "command": "bench test --package ./internal/adopt --run '^TestCompleteLifecycleDecisions$'", "probe": "internal/adopt/decision.go: swap `case lifecycleAdd, lifecycleChange, lifecycleRemove:` for `case lifecycleAdd, lifecycleChange:`; expect TestCompleteLifecycleDecisions"},
        {"id": "tc17", "command": "bench test --package ./internal/adopt --run '^TestLifecycleDecisionOrder$'", "probe": "internal/adopt/decision.go: swap `return decisions[i].Path < decisions[j].Path` for `return false && decisions[i].Path < decisions[j].Path`; expect TestLifecycleDecisionOrder"},
        {"id": "tc19", "command": "bench test --package ./internal/adopt --run '^TestPrePushHookIgnoresALegacyPin$'", "probe": "internal/adopt/prepush.sh: insert `echo \"topic push\" >&2` before `exit 0`; expect TestPrePushHookIgnoresALegacyPin"},
        {"id": "tc20", "command": "bench test --package ./internal/adopt --run '^TestLifecycleUpgradePreviewCompatibility$/^plan_kinds$'", "probe": "internal/adopt/upgrade.go: swap `if fingerprint == \"\" {` for `if false {`; expect TestLifecycleUpgradePreviewCompatibility"},
        {"id": "tc21", "command": "bench test --package ./internal/adopt --run '^TestLifecycleUpgradePreviewCompatibility$/^plan_kinds$'", "probe": "internal/adopt/upgrade.go: swap `fingerprint := old` for `fingerprint := entry.kind + \":new\" + old[:0]`; expect TestLifecycleUpgradePreviewCompatibility"},
        {"id": "tc22", "command": "bench test --package ./internal/adopt --run '^TestLifecycleUpgradePreviewCompatibility$/^plan_kinds$'", "probe": "internal/adopt/upgrade.go: swap `removed++` for `added++`; expect TestLifecycleUpgradePreviewCompatibility"},
        {"id": "tc23", "command": "bench test --package ./internal/adopt --run '^TestLifecycleUpgradePreviewCompatibility$/^plan_kinds$'", "probe": "internal/adopt/upgrade.go: swap `Preserve: []string{\"CLAUDE.md\"}` for `Preserve: nil`; expect TestLifecycleUpgradePreviewCompatibility"},
        {"id": "tc24", "command": "bench test --package ./internal/adopt --run '^TestLifecycleUpgradePreviewCompatibility$/^hook_refresh$'", "probe": "internal/adopt/upgrade.go: swap `direction != 0 && prePushRefreshEligible(root)` for `false && direction != 0 && prePushRefreshEligible(root)`; expect TestLifecycleUpgradePreviewCompatibility"}
      ]
    },
    {
      "id": "FT411-C2",
      "tickets": [
        "03-cut-subsumed-landing-tests.md",
        "04-cut-subsumed-refusal-route-tests.md"
      ],
      "verification": [
        {"id": "worktree-suite", "command": "bench test --package ./internal/worktree"},
        {"id": "tc25", "command": "bench test --package ./internal/worktree --run '^TestLandCommandPrintsNoHeadsLineWithoutRecords$'", "probe": "internal/worktree/land.go: swap `return counts[assignment]` for `return counts[assignment] + 1`; expect TestLandCommandPrintsNoHeadsLineWithoutRecords"},
        {"id": "tc26", "command": "bench test --package ./internal/worktree --run '^TestLandSkipsTheRefreshWithoutBuildInputs$'", "probe": "internal/worktree/land_effects.go: swap `return landedComplete(stdout, result, active, records)` for a call that keeps the exit code, prints one trailer line after the landed record, and returns the code; expect TestLandSkipsTheRefreshWithoutBuildInputs"},
        {"id": "tc27", "command": "bench test --package ./internal/worktree --run '^TestResumeLandCommandCompletesAnInterruptedMarker$'", "probe": "internal/worktree/land_resume.go: swap `stderr) != 0 {` for `stderr) == 0 {`; expect TestResumeLandCommandCompletesAnInterruptedMarker"},
        {"id": "tc28", "command": "bench test --package ./internal/worktree --run '^TestResumeLandCommandAcceptsSpecSlugAndPath$'", "probe": "internal/worktree/land_resume.go: swap `landingSlug(parsed.Flags[\"--spec\"])` for `parsed.Flags[\"--spec\"]`; expect TestResumeLandCommandAcceptsSpecSlugAndPath"},
        {"id": "tc29", "command": "bench test --package ./internal/worktree --run '^TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion$'", "probe": "internal/worktree/land_effects.go: swap `if j.releaseLandingAssignment(j, a, root, []string{\"--request\", request, path}, io.Discard, &diagnostic) == 0 {` for `if true {`; expect TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion"},
        {"id": "tc30", "command": "bench test --package ./internal/worktree --run '^TestLandCommandStatesTheCensusCountAndDropsTheRecords$'", "probe": "internal/worktree/land.go: omit `fmt.Fprintf(stderr, \"census heads{%s}\\n\", breakdown)`; expect TestLandCommandStatesTheCensusCountAndDropsTheRecords"},
        {"id": "tc31", "command": "bench test --package ./internal/worktree --run '^TestLandCommandStatesTheCensusCountAndDropsTheRecords$'", "probe": "internal/worktree/land.go: omit `fmt.Fprintf(stderr, \"census output{%s}\\n\", breakdown)`; expect TestLandCommandStatesTheCensusCountAndDropsTheRecords"},
        {"id": "tc32", "command": "bench test --package ./internal/worktree --run '^TestLandCommandStatesTheCensusCountAndDropsTheRecords$'", "probe": "internal/worktree/lifecycle.go: swap `_ = census.Drop(a.home, root, id)` for `_ = census.Drop`; expect TestLandCommandStatesTheCensusCountAndDropsTheRecords"},
        {"id": "tc33", "command": "bench test --package ./internal/worktree --run '^TestLandCommandStatesTheCensusCountAndDropsTheRecords$'", "probe": "internal/census/output.go: swap `assignment+outputSuffix` for `assignment` in the RecordOutput write; expect TestLandCommandStatesTheCensusCountAndDropsTheRecords"},
        {"id": "tc34", "command": "bench test --package ./internal/worktree --run '^TestResumeLandCommandPublicBindsPublishedLandingIdentity$'", "probe": "internal/worktree/landingpolicy/landingpolicy.go: swap `if !f.MarkerPresent || !f.MarkerReachesPublished() {` for `if true {`; expect TestResumeLandCommandPublicBindsPublishedLandingIdentity"},
        {"id": "tc36", "command": "bench test --package ./internal/worktree --run '^TestLandingFacesFollowTheirRoutes$'", "probe": "internal/refusalroute/faces_land.go: swap `Name:      \"destination-not-clean\",` for `Name:      \"destination-not-clean-x\",`; expect TestLandingFacesFollowTheirRoutes"},
        {"id": "tc37", "command": "bench test --package ./internal/worktree --run '^TestLandingFacesFollowTheirRoutes$'", "probe": "internal/refusalroute/faces_land.go: swap the `Authority: Reviewer,` line of the face whose sentence is `landing destination is not clean` for `Authority: Agent,`; expect TestLandingFacesFollowTheirRoutes"},
        {"id": "tc38", "command": "bench test --package ./internal/worktree --run '^TestLandReleasesTheLeaseOnRefusal$'", "probe": "internal/landing/landing_refusal.go: swap `\"prospective authorization refused: \"` for `\"prospective authorization declined: \"`; expect TestLandReleasesTheLeaseOnRefusal"},
        {"id": "tc39", "command": "bench test --package ./internal/worktree --run '^TestLandingFacesFollowTheirRoutes$'", "probe": "internal/worktree/land_rerun.go: swap `<repaired-source-tip>` for `<repaired-tip>` in `repairedSourceTipFlag`; expect TestLandingFacesFollowTheirRoutes"},
        {"id": "tc40", "command": "bench test --package ./internal/worktree --run '^TestLandingFacesFollowTheirRoutes$'", "probe": "internal/refusalroute/faces_land.go: omit the `commitAt(FactLabel),` step of the `land-red` route; expect TestLandingFacesFollowTheirRoutes"}
      ]
    }
  ],
  "final_verification": [
    {"id": "adopt-suites", "command": "bench test --package ./internal/adopt --package ./internal/adopt/repairtest"},
    {"id": "worktree-suite", "command": "bench test --package ./internal/worktree"},
    {"id": "census", "command": "bench test --check ordinary-build-census"},
    {"id": "spec-citations", "command": "bench test --check docs-currency-workflow"},
    {"id": "coverage", "command": "bench coverage --check subsumed-test-cuts"}
  ]
}
```

## Testing decisions

- A good test here is the kept test itself. Each acceptance row names one production mutation that the cut predicate caught at the spec stage, and the kept test must fail on it at the chunk tip.
- The seams are the existing kept tests. The two merged tests and the five moved predicates are the only new test code.
- The chunk checkpoint observes each cut through its probe requirement. `TestPackageTestCountPin`, the `ordinary-build-census` check, and the `docs-currency-workflow` check observe the structural edits.
- No row compares two rendered reports, and no row cites a test-only helper across a package boundary.

### Seam diagram

    trigger: bench probe at the chunk tip, recorded for the chunk checkpoint
        │
        ▼
    one named production mutation  ──▶  [ kept test in its own package ]  ──▶  verdict bit or silent
                      ◀ tests attach here: --expect names the kept test, and the cut test no longer exists

    trigger: the gate on the integrated tree
        │
        ▼
    the edited test tree  ──▶  [ TestPackageTestCountPin, ordinary-build-census, docs-currency-workflow ]  ──▶  green or red
                      ◀ tests attach here: the exact pin, the decision-file registry, and the spec seam citations

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| TC1 | 2, 3, 7 | A swap of `args[0] == "--compat"` for `args[0] != "--compat"` in the doctor dispatch turns the kept help test red after the T1-2 cut | `internal/adopt/adopt_test.go` (`TestDoctorParsesHelpAndUnknownFlags`) | The cut test caught only a dispatch change, so a silent kept test would leave dispatch ungraded |
| TC2 | 2, 3, 8 | A no-op file removal in the transaction publication turns the kept undo test red after the T1-3 cut | `internal/adopt/repairtest/repair_test.go` (`TestCompatibilityUndo`) | Undo of a created hook must remove it, and only the kept undo test then grades that removal |
| TC3 | 2, 3, 9 | A first-byte flip in each published file turns the kept undo test red after the T1-4 cut | `internal/adopt/repairtest/repair_test.go` (`TestCompatibilityUndo`) | The cut test compared the repaired hook bytes, and the kept test must still fail on wrong repaired bytes |
| TC4 | 5, 10 | A recovery row that reports `security-policy` as `info` turns the merged boundary test red | planned TestCompatibilityRecoveryBoundaries in internal/adopt/repairtest/repair_test.go | A merge that drops the policy decision row would pass without this probe |
| TC5 | 5, 10 | A recovery row that reports `process-interruption` as `info` turns the merged boundary test red | planned TestCompatibilityRecoveryBoundaries in internal/adopt/repairtest/repair_test.go | A merge that drops the interruption decision row would pass without this probe |
| TC6 | 5, 10 | A recovery row that reports `private-runtime` as `info` turns the merged boundary test red | planned TestCompatibilityRecoveryBoundaries in internal/adopt/repairtest/repair_test.go | A merge that drops the runtime decision row would pass without this probe |
| TC7 | 5, 10 | A repair that writes the CODEX_HOME config file turns the merged boundary test red | planned TestCompatibilityRecoveryBoundaries in internal/adopt/repairtest/repair_test.go | A merge that drops the unchanged-configuration assertions would pass without this probe |
| TC8 | 5, 10 | The merged boundary test keeps the assertion that the supervised process survives the repair | review-owned: the merged test body | No production mutation reaches the process check, so review compares the merged body with the three removed bodies |
| TC9 | 17 | The landing-test-efficiency coverage map cites `TestCompatibilityUndo` in LTE4 and LTE9 and the merged boundary test in LTE12 to LTE14 | `bench test --check docs-currency-workflow` | The check runs `bench coverage --check` over every spec and reds a cited name that the tree does not declare |
| TC10 | 2, 3, 11 | Omission of the `BENCH_CONFORMANCE_CONSUMER_ONLY=1 ` prefix in the seeded gate turns the kept gate-behavior test red after the T1-6 cut | `internal/adopt/setup_test.go` (`TestSetupGateRejectsIgnoredDeclaredInputs`) | The cut test only matched the line text, so the kept behavior test must fail when the scope prefix goes |
| TC11 | 2, 3, 12 | Omission of `"BENCH_RUN_BINARY", ` from the seeded gate inputs turns the zero-signal seed test red after the T1-9 cut | `internal/adopt/setup_test.go` (`TestSetupSeedsGateInputsInZeroSignalRepo`) | A wrong seeded list must still fail one exact-byte test |
| TC12 | 4, 12 | A seeded gate-inputs file that gains one byte only in a repository with a signal turns the repeated-setup test red | `internal/adopt/setup_test.go` (`TestSetupTwiceLeavesGateInputsIdentical`) | Before the move both kept tests stayed silent, so only the moved exact-byte check catches this |
| TC13 | 4, 13 | A legacy inventory that accepts a duplicated current asset turns the decision table red | `internal/adopt/lifecycle_decision_test.go` (`TestCompleteLifecycleDecisions`) | Before the move the kept table stayed silent on this mutation |
| TC14 | 4, 13 | A legacy inventory that accepts a desired asset with no fingerprint turns the decision table red | `internal/adopt/lifecycle_decision_test.go` (`TestCompleteLifecycleDecisions`) | Before the move the kept table stayed silent on this mutation |
| TC15 | 4, 13 | A legacy plan that refuses an input with no current and no desired asset turns the decision table red | `internal/adopt/lifecycle_decision_test.go` (`TestCompleteLifecycleDecisions`) | Before the move the kept table stayed silent when the refusal followed the preserve check |
| TC16 | 2, 3, 13 | A plan that drops remove operations turns the decision table red after the T1-10 cut | `internal/adopt/lifecycle_decision_test.go` (`TestCompleteLifecycleDecisions`) | The cut test's main assertion listed a remove operation, and the kept legacy subtest must still fail without it |
| TC17 | 2, 3, 13 | A plan that leaves its decisions unsorted turns the decision order test red after the T1-10 cut | `internal/adopt/lifecycle_decision_test.go` (`TestLifecycleDecisionOrder`) | The cut test asserted a repeatable order, and the kept order test must hold it |
| TC18 | 18 | The decision-domain registry names no deleted decision test file | `bench test --check ordinary-build-census` | The census reds a registry key whose file is absent |
| TC19 | 2, 3, 14 | A hook that prints one stderr line on a topic push turns the legacy-pin test red after the T1-11 cut | `internal/adopt/link_hook_test.go` (`TestPrePushHookIgnoresALegacyPin`) | The legacy-pin test is the only kept test that asserts empty stderr on a topic push |
| TC20 | 2, 3, 15 | A desired asset with an empty fingerprint turns the plan-kinds table red after the T1-12 cut | `internal/adopt/lifecycle_preview_test.go` (`TestLifecycleUpgradePreviewCompatibility`) | The cut unowned rows caught this mutation, so the pure table must also catch it |
| TC21 | 2, 3, 15 | A non-file asset whose fingerprint ignores its recorded hash turns the plan-kinds table red after the T1-12 cut | `internal/adopt/lifecycle_preview_test.go` (`TestLifecycleUpgradePreviewCompatibility`) | The cut inline and adapter drift rows caught this mutation, so the pure table must also catch it |
| TC22 | 2, 3, 15 | A plan that counts a removal as an addition turns the plan-kinds table red after the T1-12 cut | `internal/adopt/lifecycle_preview_test.go` (`TestLifecycleUpgradePreviewCompatibility`) | The cut dropped-row case caught this mutation, so the pure table must also catch it |
| TC23 | 2, 3, 15 | A plan with no preserved CLAUDE.md path turns the plan-kinds table red after the T1-12 cut | `internal/adopt/lifecycle_preview_test.go` (`TestLifecycleUpgradePreviewCompatibility`) | The cut preserved-CLAUDE.md row caught this mutation, so the pure table must also catch it |
| TC24 | 16 | An Upgrade preview that never counts the hook refresh turns the kept hook-refresh row red | `internal/adopt/lifecycle_preview_test.go` (`TestLifecycleUpgradePreviewCompatibility`) | A cut that also removes the hook rows would leave the refresh count ungraded |
| TC25 | 2, 3, 20 | A census count one above the record count turns the no-heads test red after the WT-2 cut | `internal/worktree/land_census_test.go` (`TestLandCommandPrintsNoHeadsLineWithoutRecords`) | The cut test asserted only the zero count, and the kept test must still fail on a wrong count |
| TC26 | 4, 21 | A landing that prints one line after the landed record turns the refresh test red | `internal/worktree/land_effects_test.go` (`TestLandSkipsTheRefreshWithoutBuildInputs`) | Before the move the refresh test stayed silent, because it read the effects table with a substring check |
| TC27 | 2, 3, 22 | A resume that reads a released source as a failed release turns the interrupted-marker test red after the WT-4 cut | `internal/worktree/land_resume_test.go` (`TestResumeLandCommandCompletesAnInterruptedMarker`) | The slug row asserted the same release outcome on the same arrangement |
| TC28 | 22 | A resume that skips the slug normalization of `--spec` turns the kept path row red | `internal/worktree/land_resume_test.go` (`TestResumeLandCommandAcceptsSpecSlugAndPath`) | A cut that also removes the path row would leave the path form ungraded |
| TC29 | 4, 24 | A landing that skips the release of its landed assignment turns the promotion test red | `internal/worktree/land_freshness_test.go` (`TestLandCommandNeverRunsCandidateLandingCodeDuringItsOwnPromotion`) | Before the move the promotion test stayed silent, because it never read the checkout after the landing |
| TC30 | 5, 25 | A landing that omits the census head line turns the merged census test red | `internal/worktree/land_census_test.go` (`TestLandCommandStatesTheCensusCountAndDropsTheRecords`) | Before the merge only the head-breakdown test caught this omission |
| TC31 | 5, 25 | A landing that omits the census output line turns the merged census test red | `internal/worktree/land_census_test.go` (`TestLandCommandStatesTheCensusCountAndDropsTheRecords`) | Before the merge only the output-breakdown test caught this omission |
| TC32 | 5, 25 | A release that keeps the census records turns the merged census test red | `internal/worktree/land_census_test.go` (`TestLandCommandStatesTheCensusCountAndDropsTheRecords`) | The merged test must keep the dropped-record assertion |
| TC33 | 5, 25 | A census that writes output records into the raw-call file turns the merged census test red | `internal/worktree/land_census_test.go` (`TestLandCommandStatesTheCensusCountAndDropsTheRecords`) | The output-breakdown test held that output records leave the raw-call count unchanged |
| TC34 | 2, 3, 26 | A resume that refuses every move of its destination turns the published-identity test red after the WT-10 cut | `internal/worktree/land_resume_test.go` (`TestResumeLandCommandPublicBindsPublishedLandingIdentity`) | The PL30 clean row asserted the same clean resume after a destination move |
| TC35 | 28 | `TestPackageTestCountPin` passes with `worktreeTestCount` equal to the live top-level test count after each ticket | `internal/worktree/parallel_census_test.go` (`TestPackageTestCountPin`) | The pin comparison is exact, so an unpinned removal or a stale pin turns it red |
| TC36 | 2, 3, 19 | A renamed land face in the registry turns the follow walk red after the WT-1 cut | `internal/worktree/refusal_route_follow_test.go` (`TestLandingFacesFollowTheirRoutes`) | The cut completeness test caught a face that loses its fixture |
| TC37 | 4, 23 | A reviewer land face that the registry gives to the agent turns the follow walk red | `internal/worktree/refusal_route_follow_test.go` (`TestLandingFacesFollowTheirRoutes`) | Before the move the walk stayed silent, because it compared the route with the same registry |
| TC38 | 2, 3, 27 | A changed prospective-refusal detail prefix turns the lease test red after the WT-11 cut | `internal/worktree/land_lease_test.go` (`TestLandReleasesTheLeaseOnRefusal`) | The follow walk does not read the detail prefix, so the lease test is its kept witness |
| TC39 | 2, 3, 27 | A changed repaired-source-tip placeholder turns the follow walk red after the WT-11 cut | `internal/worktree/refusal_route_follow_test.go` (`TestLandingFacesFollowTheirRoutes`) | The cut test pinned the rerun suffix, and the walk must still fail when the rerun cannot fill its slot |
| TC40 | 2, 3, 27 | A land-red route with no commit step turns the follow walk red after the WT-11 cut | `internal/worktree/refusal_route_follow_test.go` (`TestLandingFacesFollowTheirRoutes`) | The cut test required `bench commit --in` in the route |
| TC41 | 6, 29 | A seeded gate with no hygiene line in the detected-project branch turns the detected-project system test red | `internal/systemtest/adoption_test.go` (`TestDetectedProjectGateRejectsIgnoredDeclaredInput`) | The smoke journey stayed silent on this mutation, so a cut of this test would remove the only witness |
| TC42 | 6, 30 | A Link observer that follows a dangling destination symlink turns the dangling-leaf partition row red | `internal/adopt/lifecycle_consumption_test.go` (`TestLifecycleLinkPartitions`) | The kept mechanism rows and the pure table stayed silent on this mutation, so the T1-13 cut set fails its probe |

Not covered: story 1 — the probe rows TC1 to TC42 and the chunk checkpoints carry it, and time saved is a measurement, not a predicate.
Not covered: story 31 — the Out of scope section owns it, and no test can red an exclusion.

### Edge inventory

- A probe that reads `invalid` proves nothing. The author replaces the mutation at the same predicate, keeps the expected witness, and records the replacement and its reason in the ticket return.
- A mutation substring can drift when `main` moves before the build. The build-entry staleness pass re-reads each probe anchor. A `substring-miss` row routes to the same replacement rule.
- A probe can bite through a fixture precondition instead of the predicate. Each requirement names the kept test through `--expect`, and review reads the `first_failure` cell.
- Another landing can change `worktreeTestCount` before a ticket commits. The ticket lowers the pin from the value at its own base, and the exact pin grades the result.
- The drain can retire `specs/landing-test-efficiency` before ticket 01 runs. The ticket then skips the re-point, because no reader of the cut names remains.
- A merged test can leak state between its boundaries. The T1-5 merge keeps each planted file at the path that its source test used, and it starts the supervised process before the repair.
- No hostile-input class applies: the build changes test files, one conformance registry, and one staged spec, and it adds no input surface.
- **Won't handle:** T1-7, WT-9 — fixture sharing removes repeated arrangement but no test, so the probe rule has nothing to grade; `TestSetupGateRejectsIgnoredDeclaredInputs` and `TestResumeLandCommandPreauthenticatesCompletedRequestAndPath` stay as they are.
- **Won't handle:** T1-14 — the review asks for a scoped run, not a cut; `TestRootConformanceCheckPrintsFailures` stays the real-child witness.
- **Won't handle:** WT-7 — no production mutation reaches a branch over an empty slice, so no probe can prove the cut; `TestPackageDetailSentencesNameOneComponent` keeps its live predicate.
- **Won't handle:** T1-13 partial cuts — the review row names no per-row cut set, and its probe found a kept-silent mutation; `TestLifecycleLinkPartitions`, `TestLifecycleDroppedRows`, and `TestLifecycleClaudeOwnership` stay.

## Out of scope

- Fixture sharing for T1-7 and WT-9, and a scoped root-conformance run for T1-14: a separate time-saving capability with its own measurement. About 4 edits, 2 gate runs.
- Removal of the dead `deferredMergedDetails` allowance (WT-7): a light-path cleanup with no probe. About 1 edit, 1 gate run.
- Shared helpers for the remaining `"--spec", "x"` arguments and `worktree=released,census=0}` literals (WT-12): a helper refactor that cuts no test. About 12 edits, 2 gate runs.
- A per-row partition cut of the lifecycle e2e tables (T1-13): it needs a per-row mechanism map and a probe for each row. About 4 edits, 3 gate runs.
- The twelve copy-based lifecycle tests that `roadmap/FT411.md` names: research found about 25 rows over four tests, and `linkedCopy` copies a real Link golden. No set of twelve reproduces, so the reviewer names the set first. 0 edits until then.
- The system-suite owner-cleanup note S1: the code under test is in test files, so no production mutation can prove a cut. 0 edits.

## Further notes

### Spec-stage probe evidence

Every probe ran through `bench probe --in FT411-spec --mutants <set>` at commit `7b0a83e44`, and every row restored `yes`. A pair `bit/bit` means the candidate and the kept test each failed.

| row | mutation subject | candidate | kept test before the build |
|---|---|---|---|
| T1-2 | dispatch swap and dropped arguments in `compatibility.go` | bit | bit |
| T1-3 | no-op removal in `transaction/publication.go` | bit | bit |
| T1-4 | half write and first-byte flip in `transaction/publication.go` | bit | bit |
| T1-5 | each recovery row reported as `info` | bit (own test only) | merged test is new |
| T1-5 | config write before the repair namespace | bit (Policy and Runtime) | merged test is new |
| T1-6 | whole hygiene line and the scope prefix in `setup.go` | bit | bit |
| T1-8 | detected-branch hygiene removal in `setup.go` | bit | silent (`TestAdoptionSmokeJourney`) |
| T1-9 | `"BENCH_RUN_BINARY", ` omission | bit | bit (zero-signal test), silent (repeated-setup test) |
| T1-9 | one extra byte in a repository with a signal | bit | silent (both) |
| T1-10 | missing fingerprint, duplicate, empty refusal after the preserve check | bit | silent |
| T1-10 | remove operation dropped, decisions unsorted | bit | bit |
| T1-11 | stderr line before `exit 0` in `prepush.sh` | bit | bit |
| T1-12 | four count mutations in `upgrade.go` | bit | bit (`plan kinds`) |
| T1-12 | hook refresh never counted | not run | bit (kept hook row) |
| T1-13 | `os.Lstat` to `os.Stat` in the destination observer | bit (dangling row) | silent (kept mechanism rows and pure table) |
| WT-1 | two face renames in `faces_land.go` | bit | bit |
| WT-2 | count plus one, released exit code 2 | bit | bit |
| WT-3 | trailer after the landed record | bit | silent |
| WT-3 | effects line omitted | bit | bit |
| WT-4 | release result inverted | bit | bit |
| WT-4 | slug normalization skipped | not run | bit (kept path row) |
| WT-5 | `destination-not-clean` authority set to `Agent` | bit | silent |
| WT-6 | release skipped | bit | silent |
| WT-6 | freshness verify disabled | bit | bit |
| WT-8 | head line, output line, record drop | bit (own test only) | merged test is new |
| WT-10 | every marker check refused | bit | bit |
| WT-11 | detail prefix | bit | bit (`TestLandReleasesTheLeaseOnRefusal`), silent (follow walk) |
| WT-11 | rerun placeholder, commit step | bit | bit (follow walk) |

### Research sources

Six read-only readers in three pairs verified the rows, the enforcement surface, and the readers at commit `7b0a83e44`. The coordinator re-opened the cited lines and closed the frontier with the probes above.

- Enforcement: the exact worktree pin (`internal/worktree/parallel_census_test.go:926-943`), the exact serial ceiling (`:906`), the spec citation sweep (`internal/conformance/docs_workflow_checks_test.go:332-352`, `internal/coverage/citations.go:189-217`), and the decision-file registry (`internal/conformance/ordinary_build_census_test.go:105-123`).
- No anchor pins any candidate test file (`bench anchors`), no canary pins one (`bench canary`), and no systemtest ledger names one.
- `docs/known-issues.md` holds no row about these packages.
- Not verified: the full body of `internal/conformance/fixture_bite_test.go`, and the per-row mechanism map of the T1-13 tables.

### Authoring venue

The reviewer directed on 2026-10-29 that this session writes the spec and slices the tickets, in place of the fork that `projects/benchkit.md` names.

### Flagged additions

- TC7, TC24, TC28, TC41, and TC42 are guard rows. They red an over-cut or a merge that drops an assertion. The decision source names no such row, but the cut criterion needs each guard to prove that coverage stays.
- The T1-4 and WT-11 cuts drop the review's move and fold, because the spec-stage probes bit without them.
- T1-10 moves a third row (the empty plan) that the review did not name, because its probe found the kept table silent.
- The re-point of five seam cells in `specs/landing-test-efficiency/spec.md` edits another staged spec. It changes no behavior of that spec. Flag for reviewer veto.

### Source-sentence-to-row table

| source sentence | rows |
|---|---|
| `FT411.cut`: each row whose probe proves a kept test red on the same mutation is cut | TC1–TC7, TC10–TC17, TC19–TC23, TC25–TC27, TC29–TC34, TC36–TC40 |
| `FT411.cut`: a whole-test cut in internal/worktree repins its top-level test count in the same ticket | TC35 |
| `FT411.probe`: before each cut, a probe proves that the subsuming test turns red on the same mutation | the completion-plan probe of each row above, and the spec-stage evidence table |
| A row whose probe does not prove a kept test is not cut | TC41, TC42, and the Won't handle lines |

### Pre-review proof checklist

- `Cited symbols`: every cited test function resolved by `rg` at `7b0a83e44`. The merged `TestCompatibilityRecoveryBoundaries` is planned.
- `Import edges`: none. The build adds no import.
- `Source-row clauses and occurrences`: the two criteria are quoted in the table above. The occurrences of `roadmap/FT411.md` are its two 2026-10-09 review lines.
- `Promised field labels`: none.
- `Changed-function callers`: no production function changes.
- `Copy survival`: TC4 to TC7 and TC30 to TC33 red a merged test that drops a predicate of a removed copy.
- `Rendered-shape readers`: none. No rendered output changes.
- `Pin operators`: `TestPackageTestCountPin` fails on `count < worktreeTestCount` and on `count > worktreeTestCount`.
- `Entry reads`: none.
- `Derived expectations`: the WT-5 list stays independent of the registry, and TC37 records its red.
- `Consolidated rules`: none.
- `Quantified obligations`: each review row has one disposition in the table above, and each cut predicate has one checkpoint probe.
- `Workflow-step writes`: none.

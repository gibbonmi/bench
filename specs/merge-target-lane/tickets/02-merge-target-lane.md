# Grade a merge with the target's lane from any checkout

Blocked by: 01-kit-for-root.md
Writes: internal/worktree/merge.go, internal/worktree/merge_refusal.go, internal/worktree/merge_target_lane_test.go (new), internal/worktree/parallel_census_test.go, internal/worktree/journey_test.go, internal/worktree/merge_caller_root_test.go, internal/worktree/refusal_route_follow_test.go, internal/worktree/merge_target_grade_test.go, internal/worktree/merge_route_test.go, internal/worktree/merge_fold_red_test.go, internal/worktree/verb_fixture_test.go, internal/worktree/merge_test.go, internal/worktree/single_read_census_cases_test.go, internal/worktree/effect_census_test.go, docs/research/ft311-session-orchestration.md, CHANGELOG.md
Covers: R10, R11, R12, R13, R14, R15, R16, R19, R20, R21

## What to build

A `bench worktree merge` from any checkout of the kit repository grades the composed tree with the target's built-in selective lane.
An adopted kit, a foreign kit, and an unknown repository keep their current grade.

Ticket 01 supplies `func KitForRoot(kit, root string) string` in the gate package.
It answers `root`, unchanged as the caller spelled it, only when each of these conditions holds:

- `kit` is not empty.
- The canonical spelling of the git top level of `kit` is the canonical spelling of `kit`.
- The canonical spelling of the git common directory of `kit` is the canonical spelling of the git common directory of `root`.

Every other input answers `kit`, unchanged as the caller spelled it.
A git failure or a canonical-path failure at any step answers `kit`.

In `mergeOwner` (`internal/worktree/merge.go`), pass `gate.KitForRoot(kit, target)` to `gate.LaneForCommitAtKit`.
Change only that one line, so the fold (`merge.go:95`) and the target-alone grade (`merge_refusal.go:84`) use one derivation.
Do not resolve the kit at either caller, and do not read `BENCH_KIT` or the working directory below the verb entry.
Replace the claim "the way `bench commit` resolves it" in the `mergeOwner` doc comment.
The new comment states that the merge resolves the kit for the target, as the wrapper resolves the kit for its working directory.
`merge.go` has 22 lines of headroom against its 400-line budget.

Add four parallel top-level tests in the new file `internal/worktree/merge_target_lane_test.go`, with the names that the coverage map gives:

- R10 and R11: the kit is the fixture's primary checkout, as in `land_broker_notice_test.go`.
- R12 and R13: the same kit, and the incoming change makes the lane red while the target alone passes its lane.
- R14: the kit is a second linked worktree of the same repository.
- R15: the kit is `<root>/.bench`, a directory that exists, and the target declares no lane.

The lane fixture commits `go.mod`, one trivial Go package, and a whole gate that writes a tally file, before the assignment exists.
The `go.mod` has a `go` directive older than the host toolchain and no `toolchain` line.
The green incoming change edits only `go.mod`, so the built-in lane selects only the class that runs `go vet` and `go build`.
The red incoming change adds an unknown directive to `go.mod`, so `go vet` exits 1.
In R12 the whole gate writes the tally and exits 1, so a whole-gate grade of the target alone picks `merge-target-red`.

Drive the merge verb through `runVerb` and the `mergeSet.kit` override.
Run git through the package harness (`gitRun` and its peers), because `TestSerialJourneyHarnessCensus` refuses `exec.Command`, `gittest.`, `Setenv`, and `Chdir` in other test files.
Set no environment variable, so the four tests stay parallel and the serial ceiling of 45 stays.
Reuse `commitWholeGate`, `wholeGateMergeFixture`, `foldRedFixture`, and `mergeRefusalFixture` where they fit; add a fixture variant beside its owner when one does not fit.
Put a variant of `mergeFixture` in `merge_target_lane_test.go`, because its owner `merge_test.go` must not grow.
The spec rates the embed-target scan and the build cache lock in this fixture as uncertain.
If the real built-in lane cannot run in the fixture, stop and report a material acceptance shortfall; do not stub the lane resolution.

Raise `worktreeTestCount` in `internal/worktree/parallel_census_test.go` from 798 to 802 in place, with no net line.
That file is 1096 lines, and `merge_test.go` is 901 lines; neither has a grant, so neither grows.
Correct the `mergeSetOver` doc comment in `merge_caller_root_test.go`: the manifest lane applies under a kit apart from the target's repository.
Correct the `mergeSet.kit` doc comment in `verb_fixture_test.go` if it states the same old rule.

Add one `CHANGELOG.md` entry under `[Unreleased]`.
The entry states that a merge from any checkout of the kit repository now grades with the target's lane.
It also states that before, such a merge ran the whole-project gate.
The research note `docs/research/ft311-session-orchestration.md` cites `merge.go:383` in its link table; no gate reads that link, and no edit is planned.
That link was already stale at the base commit, where `merge.go` has 378 lines, so this build does not cause the drift.

## Acceptance

- [ ] `mergeOwner` passes `gate.KitForRoot(kit, target)` to `LaneForCommitAtKit`, and its two production callers pass the verb entry's kit unchanged.
- [ ] A merge from the primary checkout prints a lane line that starts `lane{outcome=pass,` and carries `classes=` (R10).
- [ ] The same merge leaves the whole-gate tally file absent (R11).
- [ ] A fold that the lane grades red refuses with face `merge-fold-red` when the target alone passes its lane (R12).
- [ ] The same refused merge leaves the whole-gate tally file absent (R13).
- [ ] A merge whose kit is another linked worktree prints a lane line that carries `classes=` (R14).
- [ ] A merge under an adopted `<root>/.bench` kit on a target with no declared lane writes the tally and prints no `lane{` line (R15).
- [ ] `TestSingleReadCensusOnTheLiveTree` stays green (R16).
- [ ] `TestMergeGradesIncomingProseFromTheComposedTreeWhateverTheCallerRoot` exits 0 (R19).
- [ ] `TestMergeRunsTheDeclaredLaneOnTheComposedTree` stays green (R20).
- [ ] `TestEffectBoundaryCensus` stays green (R21).
- [ ] `TestPackageTestCountPin` reads 802, and `TestSerialSetStaysBelowTheCeiling` reads 45.
- [ ] A probe that resolves the kit at the fold call site only reds R12; the restore turns it green again.

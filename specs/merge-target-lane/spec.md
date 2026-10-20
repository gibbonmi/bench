# Grade a merge with the target's lane from any checkout

Status: staged
Roadmap: FT335
Decision source: named reviewed artifact — roadmap/FT335.md (drained 2026-10-07 in d-6259cc0d8421) with capture/ft335-maintenance-merge-observation-2026-10-07.md.
Verification log: 4 iteration(s) to accept — the spec round took 2: iteration 1 rejected a second production owner of the kit rule, and reviewer decisions D1 and D2 closed it. The ticket round took 2: iteration 1 rejected one chunk for the seam ticket and its consumer, and the split into C1 and C2 closed it. Slicing research moved R09 to a new file and removed the toolchain environment variable.

## Problem

A reviewer runs `bench worktree merge --from <x> <target>` from the primary checkout.
The target is a worktree of the kit repository, and its kit carries the built-in selective lane.
The merge still runs the whole-project gate on the composed tree.
On 2026-10-07 a Markdown-only fold ran all six gate phases from the primary checkout.
The same merge from inside the target ran only its prose lane.

The cause is one kit and root pair.
The verb entry reads `BENCH_KIT` once, and from the primary checkout the wrapper sets it to the primary checkout ([effects.go](../../internal/worktree/effects.go:35), [bench.sh](../../bin/bench.sh:295)).
[mergeOwner](../../internal/worktree/merge.go:369) passes that kit with the target worktree to [LaneForCommitAtKit](../../internal/gate/kit_source.go:78).
The two directories differ, so [LaneFor](../../internal/gate/manifest.go:171) reads the target's manifest.
The kit repository tracks no `.bench/phases.json`, so the lane is nil.
[NewForLane](../../internal/landing/lane_owner.go:30) then answers the whole-project owner.

Two production callers share `mergeOwner`: the fold at [merge.go](../../internal/worktree/merge.go:95) and the target-alone grade at [merge_refusal.go](../../internal/worktree/merge_refusal.go:84).
`bench consumers worktree.mergeOwner worktree.mergeTargetGrade` lists exactly these call sites and the one `mergeTargetGrade` call at `merge.go:124`.
A caller in another worktree of the same repository meets the same pair, because the wrapper sets the kit to the caller's own worktree.

The caller-root face of FT335 is closed.
Commit `e2a23bac` moved the grading root to the target worktree ([landing merge](../../internal/landing/merge.go:100), [GradeTarget](../../internal/landing/merge.go:142)).
That fix did not change lane resolution, and no test sets the kit to a checkout of the target's repository.
Every merge fixture uses `t.TempDir()` as its kit ([merge_caller_root_test.go](../../internal/worktree/merge_caller_root_test.go:29)).

## Solution

A merge resolves the kit for its target the way the wrapper resolves the kit for its working directory.
When the kit is the top level of its own tree, and the target is a working tree of the same repository, the target is the kit.
The target then grades with the built-in selective lane, from any checkout of the kit repository.
Every other kit stays as the caller gave it, so an adopted repository and a foreign kit keep their current lane or the whole-project gate.

The [benchkit profile](../../projects/benchkit.md:75) already states the contract: "The target's fast lane grades the composed tree".
This spec makes that sentence true for a primary-checkout caller.

## User stories

Line: opus / medium.
Implementation-line reason: the hardest chunk is the verb fixture that runs the real built-in lane with the Go toolchain. The spec gives exact predicates, and planned rows cover each wrong answer. The mid tier comes from the seam uncertainty of that fixture, which no probe has run yet. The medium effort comes from the project route for gate and conformance logic.
Harder chunks: C2.

### The merge grades with the target's lane

1. As a reviewer who merges from the primary checkout, I want the target's built-in selective lane to grade the composed tree, so that a small fold does not pay for the whole-project gate.
2. As a reviewer who merges from the primary checkout, I want the whole-project gate not to run when the target's lane applies, so that one merge runs one authority.
3. As a coordinator who merges from another worktree of the kit repository, I want the target's lane too, so that the caller's checkout does not change the grade.
4. As a reviewer who reads a red fold, I want the target-alone grade to use the same lane as the fold, so that the refusal face names the true cause.
5. As a reviewer, I want a merge from the primary checkout to select the same kit as a merge from inside the target, so that the exec workaround is not necessary.

### The closed fallbacks stay

6. As a maintainer of an adopted repository, I want a merge on a target with no declared lane to keep the whole-project gate, so that an unconfigured repository never gets a weaker grade.
7. As a maintainer of a linked repository, I want a kit from a different repository to stay the kit, so that a foreign kit never grants the built-in lane.
8. As an operator, I want a target outside any readable repository to keep the given kit, so that an unknown answer falls back to the stricter grade.
9. As an operator, I want an empty kit to keep its current meaning, so that the root-as-kit fallback does not change.
10. As an operator whose kit path goes through a symbolic link, I want the same answer as the physical path, so that a path spelling does not change the grade.
11. As an operator whose checkout path holds a space, I want the same answer, so that a legal path does not change the grade.
12. As an operator, I want a kit path that does not exist to stay the kit, so that a missing directory falls back to the stricter grade.

### One rule and no new reads

13. As a kit maintainer, I want the merge's kit rule and the wrapper's kit rule to give the same answer over the same repository shapes, so that the two forms cannot drift.
14. As a kit maintainer, I want the fix to add no second read of `BENCH_KIT` and no read of the working directory, so that the single-read rule stays true.
15. As a kit maintainer, I want the git common-directory query to stay in the git adapter, so that the plumbing owner stays single.
16. As a kit maintainer, I want `bench commit`, the gate phase table, and the other kit readers unchanged, so that this fix stays inside the merge.
17. As a reviewer, I want the manifest lane and the caller-root grade of a merge unchanged, so that the shipped caller-root fix stays closed.

## Implementation decisions

The gate package gains one exported pure form, `KitForRoot(kit, root string) string`.
It applies the two conditions of the wrapper's `kit_dir` rule ([bench.sh](../../bin/bench.sh:67)).
The wrapper answers the top level of the working directory, and `KitForRoot` answers `root` unchanged.
The two answers agree when `root` is a top level, which holds for a merge target and for each R09 shape.
It takes both values as arguments and reads no environment variable and no working directory.
It answers `root` only when each of these conditions holds:

- `kit` is not empty.
- The git top level of `kit` is the same directory as `kit`.
- The git common directory of `kit` is the same directory as the git common directory of `root`.

The answer `root` keeps the caller's spelling of `root`, and every other input answers `kit` as the caller spelled it.
A git failure or a canonical-path failure at any step answers `kit`, so an unknown repository falls back to the caller's kit.
The comparison uses canonical spellings through the existing path owner, as [KitSourceCheckoutAtKit](../../internal/gate/kit_source.go:53) does.
It uses the existing git adapter queries [RootAt](../../internal/git/git.go:96) and [CommonDir](../../internal/git/worktree_admin.go:193).
The gate package already imports both owners, so the spec adds no import edge.

`mergeOwner` passes `gate.KitForRoot(kit, target)` to `LaneForCommitAtKit`.
Both production callers reach that one line, so the fold and the target-alone grade use one derivation.
The `mergeOwner` doc comment names the target-side kit resolution.
Its current claim, "the way `bench commit` resolves it", is false for a primary-checkout caller.

For the kit repository, the resolved kit is the target, so `LaneForCommitAtKit` answers the built-in lane with `Selective` set and an empty `Kit`.
The lane then builds its run binary from the private checkout of the composed tree ([selectLaneRunBinary](../../internal/gate/lane.go:263)), as a target-local commit does.

The fix stays out of `kitRootAt`, `LaneForCommitAtKit`, and `KitSourceCheckoutAtKit`.
Those helpers serve `bench commit`, the phase table, the test report, preflight, coverage, and the adoption doctor.
A change there changes their answers too, and FT421 owns the closest of those questions.

The binary cannot tell an explicit `BENCH_KIT` from the value the wrapper derives, because the wrapper exports both under one name ([bench.sh](../../bin/bench.sh:317)).
Under reviewer decision D2, the merge resolves every kit that is the top level of a checkout of the target's repository to the target.

The wrapper rule stays in Bash, because the wrapper selects the binary before any Go code runs.
Reviewer decision D1 accepts the Bash form and the Go form as an exception to the one-source rule, with row R09 as the drift guard.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| C1 / 01-kit-for-root.md | The gate package answers the kit for a root by the wrapper's rule, and the Go form and the Bash form agree over the wrapper's repository shapes | R01, R02, R03, R04, R05, R06, R07, R08, R09, R17, R18 | KitForRoot unit tests, wrapper differential, BENCH_KIT read check, plumbing owner check, root conformance, gate and cmd package runs | no |
| C2 / 02-merge-target-lane.md | A merge from any checkout of the kit repository grades with the target's selective lane, and the adopted, foreign, and unknown kits keep their current grade | R10, R11, R12, R13, R14, R15, R16, R19, R20, R21 | merge verb lane and fallback witnesses, single-read census, effect census, existing merge lane tests | yes |

Ticket 01 creates the `KitForRoot` seam with its unit rows and the differential row, so it is its own review chunk.
Ticket 02 makes `mergeOwner` consume the seam and adds the merge verb rows.
C2 starts after the C1 chunk review closes.

### Completion plan

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"ft335-merge-target-lane-full-20261019","orchestrator_session":"claude:ft335-orchestrator-20261019","author_limit":1,"assignments":{"01-kit-for-root.md":[{"session":"claude:ft335-t01-author-20261019","assignment":"8ef5866f9492f2fc9d870172e0aac8d2","model":"opus","effort":"medium","source":"04ee12d5dd582d7f6e567c6d3ab156ea7608ff5c","native_ref":"claude-agent:bench-writer-medium:ft335-t01-author-20261019"},{"session":"claude:ft335-t01-repair1-20261019","assignment":"8ef5866f9492f2fc9d870172e0aac8d2","model":"opus","effort":"medium","source":"6cde7864d77526dffdb139d76fc5979e9bfce730","native_ref":"claude-agent:bench-writer-medium:ft335-t01-repair1-20261019","predecessor":"claude:ft335-t01-author-20261019","trigger":"user-directed","stopped":"claude-agent:bench-writer-medium:ft335-t01-author-20261019","preserved":"03bec7089b851acf9865fd84b62493624e3970b7"}],"02-merge-target-lane.md":[]}},"chunks":[{"id":"C1","tickets":["01-kit-for-root.md"],"verification":[{"id":"t1-kit-rule","command":"bench test --package ./internal/gate --run 'KitForRoot|LaneForCommit|LaneFor|KitValue|KitSourceCheckout'","probe":"Remove the top-level condition from KitForRoot and record the R03 red, the restore, and the green.","ticket":"01-kit-for-root.md"},{"id":"t1-wrapper-differential","command":"bench test --package ./cmd/bench --run 'ShellWrapper'","probe":"Remove the top-level condition from KitForRoot and record the R09 red, the restore, and the green.","ticket":"01-kit-for-root.md"},{"id":"t1-plumbing-owner","command":"bench test --check git-plumbing-owner","ticket":"01-kit-for-root.md"},{"id":"t1-root-conformance","command":"bench test --check root-conformance","ticket":"01-kit-for-root.md"},{"id":"t1-gate-package","command":"bench test --package ./internal/gate","ticket":"01-kit-for-root.md"},{"id":"t1-cmd-package","command":"bench test --package ./cmd/bench","ticket":"01-kit-for-root.md"}]},{"id":"C2","tickets":["02-merge-target-lane.md"],"verification":[{"id":"t2-merge-lane","command":"bench test --package ./internal/worktree --run 'Merge|SingleReadCensus|EffectBoundaryCensus'","probe":"Resolve the kit at the fold call site only and record the R12 red, the restore, and the green.","ticket":"02-merge-target-lane.md"},{"id":"t2-root-conformance","command":"bench test --check root-conformance","ticket":"02-merge-target-lane.md"},{"id":"t2-worktree-package","command":"bench test --package ./internal/worktree","ticket":"02-merge-target-lane.md"},{"id":"t2-cmd-package","command":"bench test --package ./cmd/bench","ticket":"02-merge-target-lane.md"}]}],"final_verification":[{"id":"worktree","command":"bench test --package ./internal/worktree"},{"id":"gate","command":"bench test --package ./internal/gate"},{"id":"coverage","command":"bench coverage --check merge-target-lane"}]}
```

## Testing decisions

A good test drives the merge verb or the pure kit form with real git repositories.
It observes the lane line, the whole-gate tally, the refusal face, or the returned kit.
It never stubs the lane resolution, because the defect lives in that resolution.

The unit seam is `gate.KitForRoot` over real repositories built in the test.
Each fixture shape is a separate test function, so each row names one predicate.
The fixtures hold a primary checkout with a linked worktree, a second linked worktree, an adopted `.bench` kit, an unrelated repository, a plain directory, and a symbolic link.
Each adopted `.bench` kit in R03 and R15 is a directory that exists, because a missing directory makes git fail and hides the top-level condition.
The R06 test sets its working directory to the primary checkout of the root's repository with `t.Chdir`, so a form without the empty-kit guard answers `root`.
Do not mark the R06 test parallel, because `t.Chdir` refuses a parallel test.
The gate fixtures run git through the package's `runGitIn` helper, because the ordinary-build census refuses a `worktree` literal in a `gittest` call or an `os/exec` process in a gate test.

The verb seam is the merge verb through the existing `runVerb` harness and the `mergeSet.kit` override ([verb_fixture_test.go](../../internal/worktree/verb_fixture_test.go:58)).
The precedent for a kit that is the fixture's primary checkout is [land_broker_notice_test.go](../../internal/worktree/land_broker_notice_test.go:50).
The fixture commits `go.mod`, one trivial Go package, and a whole gate that writes a tally, before the assignment exists.
The incoming change edits only `go.mod`, so the built-in lane selects only the class that runs `go vet` and `go build`.
Those two checks run no Bench run binary, so the fixture needs no Bench build.
The fixture's `go.mod` has a `go` directive older than the host toolchain and no `toolchain` line, so the lane never fetches a toolchain.
The test sets no environment variable, so it stays parallel and the package's serial ceiling does not change.

A probe on this host showed the Go side of that fixture.
`go vet` and `go build` with `-trimpath` exit zero on a module with one trivial package.
`go vet` exits 1 on a `go.mod` with an unknown directive.
The probe did not run the lane itself; the embed-target scan and the build cache lock in that fixture stay `uncertain` until ticket 02 runs it.

The differential seam is [TestShellWrapperResolvesKitForTheInvocation](../../cmd/bench/main_test.go:680).
Its repository shapes run through the wrapper and through `gate.KitForRoot`, and the two answers must be equal.
The R09 test lives in a new file in `cmd/bench`, because `main_test.go` is over its structure budget and the lane's growth check refuses a line added there.
One wrapper fixture serves both tests, so the existing test's fixture moves out of `main_test.go` into a shared helper, and `main_test.go` only shrinks.

### Seam diagram

    trigger: bench worktree merge --from <x> <target>, from any checkout
        │
        ▼
    ambient kit (BENCH_KIT, read once)  ──▶  [ gate.KitForRoot(kit, target) ]  ──▶  kit for the target
                                                  ◀ unit rows R01–R08: real repositories, returned kit
        │
        ▼
    [ mergeOwner → LaneForCommitAtKit → NewForLane ]  ──▶  lane owner or whole-project owner
        │                          │
        ▼                          ▼
    fold grade (merge.go:95)   target-alone grade (merge_refusal.go:84)
        ◀ verb rows R10–R15: lane line, whole-gate tally, refusal face

    wrapper kit_dir (bin/bench.sh)  ══ differential R09 ══  gate.KitForRoot

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| R01 | 1, 11 | A linked worktree of the kit's repository answers that worktree, under a parent directory whose name holds a space | `internal/gate/kit_root_test.go` (`TestKitForRootNamesALinkedWorktreeOfThePrimaryKit`) | A form that keeps the kit answers the primary checkout |
| R02 | 3 | A kit that is another linked worktree of the same repository answers the target worktree | `internal/gate/kit_root_test.go` (`TestKitForRootNamesALinkedWorktreeOfASiblingKit`) | A form that matches only the primary checkout keeps the sibling kit |
| R03 | 6 | An adopted kit at `<repo>/.bench` stays the kit for a linked worktree of that repository | `internal/gate/kit_root_test.go` (`TestKitForRootKeepsAnAdoptedSubdirectoryKit`) | A form without the top-level condition answers the worktree |
| R04 | 7 | A kit from a different repository stays the kit | `internal/gate/kit_root_test.go` (`TestKitForRootKeepsAKitFromAnotherRepository`) | A form that skips the common-directory comparison answers the root |
| R05 | 8 | A root outside any repository keeps the kit | `internal/gate/kit_root_test.go` (`TestKitForRootKeepsTheKitOutsideARepository`) | A form that treats a git failure as a match answers the root |
| R06 | 9 | An empty kit answers the empty kit when the working directory is the primary checkout of the root's repository | `internal/gate/kit_root_test.go` (`TestKitForRootKeepsAnEmptyKit`) | A form without the empty-kit guard resolves the empty kit to the working directory and answers the root |
| R07 | 10 | A kit spelled through a symbolic link to the primary checkout answers the linked worktree | `internal/gate/kit_root_test.go` (`TestKitForRootMatchesASymlinkedKitSpelling`) | A raw string comparison of the top level keeps the kit |
| R08 | 12 | A kit path that does not exist stays the kit | `internal/gate/kit_root_test.go` (`TestKitForRootKeepsAnAbsentKit`) | A form that treats an unreadable kit as a match answers the root |
| R09 | 5, 13 | For each repository shape of the wrapper test without an environment override, `gate.KitForRoot` gives the kit the wrapper exports | `cmd/bench/wrapper_kit_test.go` (`TestShellWrapperKitMatchesKitForRoot`) | A Go rule that drops either wrapper condition disagrees on the adopted or the project shape |
| R10 | 1 | A merge from the primary checkout prints a lane line that starts `lane{outcome=pass,` and carries `classes=` | planned TestMergeFromThePrimaryGradesWithTheTargetLane in internal/worktree/merge_target_lane_test.go | Today the nil lane runs the whole gate, which prints no lane line |
| R11 | 2 | The same merge leaves the whole-gate tally file absent | planned TestMergeFromThePrimaryGradesWithTheTargetLane in internal/worktree/merge_target_lane_test.go | Today the whole gate runs and writes the tally |
| R12 | 4 | A fold that the lane grades red refuses with face `merge-fold-red` when the target alone passes its lane | planned TestMergeFromThePrimaryGradesTheTargetAloneWithTheFoldLane in internal/worktree/merge_target_lane_test.go | A target-alone grade under the caller's kit runs the whole gate, whose `exit 1` body picks `merge-target-red` |
| R13 | 4 | The same refused merge leaves the whole-gate tally file absent | planned TestMergeFromThePrimaryGradesTheTargetAloneWithTheFoldLane in internal/worktree/merge_target_lane_test.go | A whole-gate grade of either the fold or the target writes the tally |
| R14 | 3 | A merge whose kit is another linked worktree prints a lane line that carries `classes=` | planned TestMergeFromASiblingWorktreeGradesWithTheTargetLane in internal/worktree/merge_target_lane_test.go | A resolution that only matches the primary checkout runs the whole gate |
| R15 | 6 | A merge under an adopted `<root>/.bench` kit on a target with no declared lane writes the whole-gate tally and prints no `lane{` line | planned TestMergeUnderAnAdoptedKitKeepsTheWholeGate in internal/worktree/merge_target_lane_test.go | A resolution without the top-level condition runs the built-in lane |
| R16 | 14 | Each `KitValue` read in the worktree package stays at a verb entry | `internal/worktree/single_read_census_cases_test.go` (`TestSingleReadCensusOnTheLiveTree`) | A `KitValue` call inside `KitForRoot` or below the merge entry turns the census red |
| R17 | 14 | `KitValue` stays the one read of `BENCH_KIT` in the gate package | `internal/gate/kit_value_test.go` (`TestKitValueIsTheOneBenchKitRead`) | A `BENCH_KIT` read inside the new form turns the check red |
| R18 | 15 | No production file outside `internal/git` spells `--git-common-dir` in a `rev-parse` call | conformance check `git-plumbing-owner`, run by `bench test --check git-plumbing-owner` | An inline common-directory query in the gate package turns the check red |
| R19 | 17 | A merge under a kit apart from the target still takes the target's manifest lane from any caller root | `internal/worktree/merge_caller_root_test.go` (`TestMergeGradesIncomingProseFromTheComposedTreeWhateverTheCallerRoot`) | A resolution that takes the caller root's failing lane makes the merge exit nonzero, and the test requires exit 0 |
| R20 | 17 | A merge on a target with a declared manifest lane runs that lane once on the composed tree | `internal/worktree/merge_test.go` (`TestMergeRunsTheDeclaredLaneOnTheComposedTree`) | A resolution that bypasses the manifest lane changes the tally or the lane line |
| R21 | 14 | No production file of the worktree package outside `effects.go` reads the environment or the working directory | `internal/worktree/effect_census_test.go` (`TestEffectBoundaryCensus`) | An `os.Getenv` or `os.Getwd` call in `merge.go` turns the census red |

Not covered: story 16 — review-owned; the gate helpers stay unchanged as a design fence, no row can tell a correct helper from a moved rule, and FT421 owns their behavior.

Planned tests are future evidence.
No new red, green, or mutation result is claimed at spec time.

### Edge inventory

The edge walk covers each hostile-input class of the profile that reaches a path argument.

- A path that holds a space reaches git as one argument; R01 covers it.
- A symbolic link spelling of the kit takes its canonical spelling; R07 covers it.
- A linked worktree's `.git` pointer file names its administration directory; `CommonDir` answers an absolute path, and R01 and R02 use real pointer files.
- An absent kit path and an empty kit are distinct inputs; R08 and R06 cover them.
- A git failure answers the caller's kit, which is the stricter grade; R05 covers it.
- Control bytes, TOON cells, and Markdown sinks do not apply, because the change adds no output.

Won't handle: the coverage census kit-phase predicate — FT421 owns that surface, and the merge owner is the surviving in-scope caller.
Won't handle: `bench commit` under an inherited `BENCH_KIT` that names another checkout — the wrapper resolves the kit for an uninherited caller, commit refuses the primary checkout, and `mergeOwner` stays the in-scope caller.
Won't handle: `bench preflight evidence --check-current` from the primary checkout — that row of `docs/known-issues.md` is a different verb, and the merge verb stays the in-scope caller.
Won't handle: lane selection for `bench worktree land` — the landing grades with the whole-project gate by design, and the merge stays the in-scope caller.
Won't handle: the wrapper test shapes "no git on PATH" and "an explicit BENCH_KIT wins" in R09 — they test the wrapper's environment handling, and the Go form receives the already exported kit.

## Out of scope

- The FT421 coverage census predicate: a separate capability with its own spec, about 3 edits and 1 gate run.
- A Go owner for the wrapper's kit rule that the wrapper consumes: the wrapper selects the binary first, so this needs its own design, about 4 edits and 2 gate runs.

## Further notes

### Reviewer decisions

- D1 (2026-10-19, reviewer sign-off in the spec conversation): AGENTS.md requires one source for each fact, and its test-expectation exception does not cover two production forms. The reviewer accepts the Bash `kit_dir` rule and the Go `gate.KitForRoot` rule as a recorded exception, with R09 as the drift guard. The rejected alternative is a Go owner that the wrapper consumes; Out of scope prices it.
- D2 (2026-10-19, reviewer sign-off in the spec conversation): the merge moves an explicit `BENCH_KIT` to the target when that kit is the top level of a checkout of the target's repository. The target, not the caller, owns the lane that grades the target's composed tree. A kit from another repository and an adopted `.bench` kit stay unchanged.

### Documents that the fix makes true

[projects/benchkit.md](../../projects/benchkit.md:75) states that the target's fast lane grades the composed tree.
The commitment criterion `FT335.lane` in `.bench/commitment.json` states the same obligation.
[docs/research/ft311-session-orchestration.md](../../docs/research/ft311-session-orchestration.md:135) states that `mergeOwner` resolves the target's lane.
Each claim becomes true after the fix, so none takes an edit.

### Flagged additions

- `gate.KitForRoot`, a new exported pure form, because the fix needs one tested owner for the kit rule.
- A standards exception: the kit rule has a Bash form and a Go form under reviewer decision D1.
- Row R09, a wrapper differential, because the rule now has a Bash form and a Go form.
- Story 11 and row R01's path with a space, from the profile's hostile-input walk.
- Rows R02 and R14, a caller in a sibling worktree, because the source names only the primary checkout.
- Row R07, a symbolic link spelling, from the profile's hostile-input walk.

### Source-sentence-to-row table

| Source clause | Coverage |
|---|---|
| FT335: a merge from the primary checkout runs the whole-project gate where the target declares a lane | R01, R10, R11 |
| FT335: the caller kit and the target pair resolve to no lane, and a nil lane takes the whole-project owner | R10, R12, R13 |
| FT335: the whole-project fallback for an unconfigured linked repo stays | R03, R04, R15 |
| FT335: the caller-root face shipped in `e2a23bac`, and the landing merge grades the composed tree from the target worktree | R19, R20 |
| Occurrences 2026-09-22 and 2026-10-07: full prospective gate from the primary checkout | R10, R11 |
| Occurrences 2026-09-25, 2026-09-26, and 2026-09-27: caller-root faces, closed by `e2a23bac` | R19 |
| Observation: the target-local merge preserves the target kit through the wrapper's same-repository detection | R09, R01 |

### Pre-review proof checklist

- Cited symbols: `mergeOwner`, `mergeTargetGrade`, `LaneForCommitAtKit`, `LaneFor`, `NewForLane`, `KitSourceCheckoutAtKit`, `RootAt`, `CommonDir`, `selectLaneRunBinary`, and `kit_dir` were opened at HEAD c2aa51652.
- Import edges: none new; the gate package already imports `internal/git` and `internal/canonicalpath`, and `cmd/bench` already imports the gate package.
- Source-row clauses and occurrences: the source table lists each clause and each occurrence.
- Promised field labels: the lane line prefix `lane{outcome=pass,`, the token `classes=`, and the faces `merge-fold-red` and `merge-target-red`.
- Changed-function callers: `bench consumers worktree.mergeOwner worktree.mergeTargetGrade` lists `merge.go:95`, `merge_refusal.go:84`, and `merge.go:124`.
- Copy survival: none; no copy is replaced.
- Rendered-shape readers: none; no output shape changes.
- Pin operators: R10 and R14 assert a substring of stdout; R11, R13, and R15 assert the tally file's absence or presence; R12 asserts face equality; K rows assert string equality of the returned kit; R16, R17, R18, R20, and R21 are pass or fail results of existing checks; R19 asserts exit 0.
- Entry reads: `KitForRoot` takes its values as arguments, and `mergeOwner` keeps the ambient kit from the verb entry.
- Derived expectations: K rows expect fixture paths; R09 takes the wrapper's answer as its oracle; R10 and R14 assert the selective token, not a check list derived from the class table.
- Consolidated rules: none.
- Quantified obligations: both production callers of `mergeOwner` take R10 through R13.
- Workflow-step writes: none.

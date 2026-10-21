# Cut the landing gate's long poles without losing failure detection

Status: staged

Roadmap: FT402

Decision source: Reviewer-confirmed current conversation, 2026-10-21.

Verification log: pending.

## Problem

A landing spends about five minutes in the whole-project gate. The three newest gates on `main` took 296 to 308 seconds.
The test phase took 173 to 174 seconds, and the system phase took 100 to 108 seconds after it.
The scheduler runs one phase at a time, so the system phase adds its whole wall time to the gate (`internal/gate/runner.go:196`).

The earlier FT402 spec, `landing-test-efficiency`, aimed at adoption and repair fixtures. The tmpfs base for the kit test run removed that cost.
On the tmpfs base, `internal/adopt` takes about 22 seconds, `internal/adopt/repairtest` about 3 seconds, and `internal/conformance` about 37 seconds.
That spec would save about 5 to 10 seconds, which is inside the run-to-run spread. The build stopped before code.

The current long poles are these packages:

| Package | Gate time | Shape of the cost |
|---|---:|---|
| `internal/systemtest` | 100 to 108 s | 105 top-level tests in one serial process; no test calls `t.Parallel` |
| `internal/probe` | 87 s (66 s alone) | 131 tests; each fixture runs `git init` and publishes a seal |
| `internal/worktree` | 70 to 71 s | 45 serial tests; 41.9 percent of git starts repeat a query |
| `internal/testreport` | 65 to 73 s | two tests copy the kit and build a real run binary each |

Sources: the gate logs `.logs/gate-20261021T192218.079145862Z-3443832.jsonl`, `gate-20261021T191248.934370406Z-2945722`, and `gate-20261021T154223.033572571Z-2308306`, with their `.out` files. The S1 to S5 research returns of 2026-10-21 supply the cost shapes.

## Solution

Run the system suite's independent tests in parallel inside its one `go test` process. Keep its one owner, its one identity ledger, and its pinned argv.
Build each expensive probe and testreport fixture once for each package run, and copy it for each test. Keep every witness on its real child.
Cache the stable git location answers for one command run, and remove the direct duplicate queries in worktree commands.
Move the worktree tests that bind the working directory or the child environment into the parallel set, through explicit seams.

Measure each chunk against the recorded baseline. A gain inside the recorded spread counts as no change. No wall-clock threshold enters the gate.

### Decisions closed on 2026-10-21

- The reviewer directed a re-spec of FT402, because the earlier premise failed on the tmpfs base.
- The scope is the system suite, the worktree git queries, the serial worktree tests, and the probe and testreport packages.
- The system suite gains `t.Parallel` in one process. Separate shard processes are out of scope.
- The git cache holds only the stable location answers of a root: the common directory, the top level, and the admin directory.
- The probe and testreport cuts share fixtures and keep every witness. No test is removed.
- The serial worktree tests that move are the working-directory group and the child-environment group.
- The earlier `landing-test-efficiency` spec stays staged and unbound, because `conformance-observation-follow-on` depends on its source-observation chunks. FT402 binds to this spec only.
- The Claude author line for this build is opus at xhigh effort.

## User stories

Line: opus / xhigh.
Implementation-line reason: LGP-C4 changes a serial journey suite with shared repositories, and LGP-C2 adds a cache under hostile-path tests. Both seams hold real uncertainty, and the reviewer set xhigh on 2026-10-21.
Harder chunks: LGP-C2, LGP-C4.

### Shared fixtures in probe and testreport

1. As a maintainer, I want the probe clamp repository built once for each package run, so that each test skips `git init` and the seal publication.
2. As a test author, I want a private copy of the clamp repository for each probe test, so that mutations stay in one test.
3. As a test author, I want each copy to keep the sealed run binary valid, so that the run-binary freshness check still accepts the copy.
4. As a maintainer, I want the PB1 witness to keep its real `go test` child, so that the real toolchain still proves a bite.
5. As a maintainer, I want the PB22 and PS40 to PS43 witnesses to keep their assertions, so that usage and expectation refusals stay pinned.
6. As a maintainer, I want the testreport run binary built once for each package run, so that two kit-copy tests do not each build it.
7. As a maintainer, I want `TestRootConformanceCheckPrintsFailures` to keep its real conformance child, so that the planted defect still reaches the diagnostic.
8. As a maintainer, I want the seam test of `selectRunBinary` to keep its real selection, so that the owned-selection path keeps a real-build test.

### Git location answers for one command run

9. As a developer, I want one command run to query the common directory of a root once, so that repeated queries cost no starts.
10. As a developer, I want the top level and the admin directory cached the same way, so that the stable location facts share one rule.
11. As a reviewer, I want every cached common directory validated again on each read, so that a later symlink or missing directory still refuses.
12. As a reviewer, I want refs, HEAD, registrations, the index, and the default branch read live, so that a write cannot leave a stale answer.
13. As a test author, I want the package-level git functions to stay uncached, so that each hostile-path test still observes its own starts.
14. As a maintainer, I want no process-global or cross-run cache, so that tests that swap PATH or recreate repositories see fresh answers.
15. As a maintainer, I want the four plumbing flags to stay inside `internal/git`, so that the git plumbing owner check stays green.
16. As a developer, I want `planExplicitWith` to stop repeating its common-directory query, so that one explicit plan costs fewer starts.
17. As a developer, I want `validateOwnerMarker` to stop asking for the common directory three times, so that each exec target costs fewer starts.

### Serial worktree tests in the parallel set

18. As a test author, I want a verb call to carry its working directory, so that an operand test needs no process `chdir`.
19. As a test author, I want a verb call to carry its child base environment, so that an exec test binds no process environment.
20. As a maintainer, I want the working-directory group and the child-environment group to call `t.Parallel`, so that they leave the serial head of the package.
21. As a maintainer, I want the serial ceiling to equal the recorded arithmetic of moved and new serial tests, so that the pin stays exact.
22. As a maintainer, I want no existing worktree test removed or merged, so that the count pin rises only by the new tests.
23. As a maintainer, I want the Home-default graders, the PATH-stub tests, and the removed-directory test to stay serial, so that their process-level subject stays graded.
24. As a maintainer, I want the effect boundary census to cover the new environment read, so that the child base environment has one owner.

### System suite tests in parallel

25. As a maintainer, I want independent system tests to call `t.Parallel`, so that the system phase stops adding their wall time one by one.
26. As a maintainer, I want one `TestMain` owner and one identity ledger to remain, so that the profile's statement of the suite stays true.
27. As a reviewer, I want the unfiltered run to verify all four terminal outcomes, so that a parallel run cannot skip a teardown proof.
28. As a reviewer, I want the ledger writes safe under concurrent tests, so that a parallel start cannot lose an observation.
29. As a test author, I want each test that needs a repository to take a private one, so that no parallel test sees a half-made shared repository.
30. As a test author, I want the three TMPDIR bindings replaced with explicit child values, so that no system test binds TMPDIR in the process.
31. As a maintainer, I want a census that classifies each system test as parallel or serial with a reason, so that no test stays serial unnoticed.
32. As a maintainer, I want the census to refuse a parallel test that reaches an owner repository, so that shared state stays serial.
33. As a maintainer, I want the gate's system argv and the phase table unchanged, so that the census pins and the profile table stay true.

### Measured results

34. As a maintainer, I want each chunk's affected package and phase timed against the recorded baseline, so that a speed claim stays reviewable.
35. As a maintainer, I want each reverted attempt recorded, so that a later run does not try it again.
36. As a reviewer, I want the final report to pair the timings with the retained failure proofs, so that speed alone establishes no quality.

## Implementation decisions

### Shared fixtures in probe and testreport

The probe test package gains a package-run golden of the clamp repository. It follows the adopt golden precedent (`internal/adopt/golden_test.go:20-69`): one `sync.Once` build and one `TestMain` cleanup.
`clampFixture` and `refusalFixture` take a private copy of the golden for each test. The copy keeps file bytes and modes, and it keeps a valid git repository.
The golden build calls the fixture publisher once. `probetest.Repository` takes that publisher as an argument (`internal/probe/probetest/fixture.go:69`), so a counting publisher observes the builds.
The seal holds source and executable digests and no root path (`internal/freshness/freshness_publish.go:97-111`). So a byte-identical copy keeps a valid seal, when the executable, the seal, and the manifest directory move together.
The build keeps the environment bindings and the `t.Chdir` of `probetest.Repository` (`internal/probe/probetest/fixture.go:106-112`) for each test.

`TestProbeBitesWhenTheFocusedTestFails` keeps its real `go test` child (PB1). The witnesses PB22 and PS40 to PS43 keep every assertion. The build removes no test from either package.

The testreport package builds its run binary once for each package run, as `testRunSelector` does (`internal/worktree/test_run_test.go:11-52`). An inherited `BENCH_RUN_BINARY` wins, and a present invalid selection refuses.
A counting builder wraps the run-binary build, so the two kit-copy tests together show at most one build.
The probe verb's own run-binary check can still run `go list -deps` through `freshness.VerifyAt` (`internal/runbinary/runbinary.go:155-163`). This spec claims no saving from that check.
`TestRootConformanceCheckPrintsFailures` and `TestNamedCheckRunsOnlyRegisteredDevScope` keep their real children and their kit copies. The build reuses a run binary only where the kit-copy sources match the binary's seal.

### Git location answers for one command run

`internal/git` gains a run value that a command creates at its entry and passes down. The worktree `ambient` value carries it (`internal/worktree/effects.go:21-37`).
The run value caches the raw answers of the common-directory, top-level, and admin-directory queries, keyed by the root spelling that the caller passes.
On every read, the run value runs the existing validation again, such as `validateCommonDir` (`internal/git/worktree_admin.go:245-259`). The cache saves the git start, not the filesystem check.
The package-level functions `CommonDir`, `AdminDir`, `Worktrees`, and `Root` stay uncached. The hostile tests in `internal/git/worktree_admin_hostile_test.go` call them directly and keep their meaning.
The run value has its own `Worktrees` form. That form reads the common directory through the run value, and it keeps the admin scan and the porcelain start live.

The run value holds no ref, HEAD, registration, index, default-branch, or ancestry answer. Those reads stay live (S2 research, finding 21 and 22).
No package variable holds a cache. A direct caller with no run value gets a fresh query, as today.

`validateOwnerMarker` (`internal/worktree/classifier.go:397-426`) asks for the common directory three times today. The starts come from `canonicalRoot`, from `Worktrees`, and from `CommonDir`.
`canonicalRoot` (`internal/worktree/worktree.go:100`) calls `poolkey.Canonical`, which runs `git.CommonDir` (`internal/poolkey/poolkey.go:18-23`). So `internal/poolkey` gains a form that reads through the run value, and `Canonical` keeps its signature for its other callers.
`planExplicitWith` (`internal/worktree/subshell.go:160-262`) starts eight common-directory queries in an owned plan today. Three come from its own `canonicalRoot`, `Worktrees`, and `CommonDir`, and three more from the nested `validateOwnerMarker`.
The other two come from the ledger reads `intent.LifecycleEvidence` and `intent.Assignments`, through `intent.Address` (`internal/intent/intent.go:34-38`).
The run value cuts the six location starts to one. The two ledger reads stay live, because the closed decision does not move `internal/intent` onto the run value. The default-branch reads also stay as they are.
The cache key is the root spelling. `canonicalRoot` returns the canonical root, so a caller that passes the canonical primary root hits the same key on each later read.

### Serial worktree tests in the parallel set

`verbCall` (`internal/worktree/verb_runner_test.go:47-55`) gains a working-directory field and a child base environment field.
`canonicalpath.Resolve` reads the process directory for a relative path (`internal/canonicalpath/canonicalpath.go:46-50`). It joins the parts as text and resolves symlinks before it cleans, on purpose.
A new `canonicalpath` form takes the directory as an argument and keeps exactly those semantics. `canonicalpath.Operand` is not that form, because it cleans lexically and resolves no symlink.
`effects.go` reads the process directory once at the verb boundary. A failed read is carried forward and surfaces only when a relative operand needs it.
So an absolute operand still works in a removed directory, and `TestLandingHandsBackAPathWithNoCanonicalForm` (`internal/worktree/land_resume_refusal_test.go:343`) keeps its behavior.
The child base environment enters `ambient` once at the verb boundary in `effects.go`. Four production sites read `os.Environ()` today, and each one takes the ambient value instead.
The four sites are `execEnv` (`internal/worktree/exec.go:201`), `gitInput` (`internal/worktree/clean.go:218`), the doctor child (`internal/worktree/build.go:331`), and the subshell entry (`internal/worktree/subshell.go:40`).

The working-directory group holds 17 tests, and the child-environment group holds 10 tests (S3 research). The build moves each test that its seam frees and records each test that stays.
Both pins change by recorded arithmetic. Each ticket that touches the package records each term in its verification excerpt.

- `worktreeTestCount` (`internal/worktree/parallel_census_test.go:1061`) equals 805 plus the new top-level tests.
- `worktreeSerialCeiling` (`parallel_census_test.go:1041`) equals 45, minus the moved tests, plus the new serial tests.

LGP-C2 adds serial tests, because its start counts use `journeyStubGit`, which binds PATH (`internal/worktree/journey_test.go:64-68`). So LGP-C2 also writes `parallel_census_test.go`.
The Home-default graders, the PATH-stub tests, and `TestLandingHandsBackAPathWithNoCanonicalForm` stay serial, because the process state is their subject.
`TestEffectBoundaryCensus` (`internal/worktree/effect_census_test.go:20`) gains `os.Environ` in its forbidden set, because `effects.go` becomes the one owner of that read. No other `os.Environ` read stays in the package's production files.

### System suite tests in parallel

The one `go test -tags=system` process stays. `TestMain`, `newSystemOwner`, the three owner repositories, and `verify` keep their shape (`internal/systemtest/owner_test.go:52-113`, `:350-373`).
Every ledger write in the owner code takes the owner mutex today. The one exception is `TestOwnerSelectionChild` (`internal/systemtest/owner_selection_test.go:156-170`).
That test plants the negative controls in a child process, where every other test is skipped. The lock census exempts it by name, and the controls stay unchanged.

About 40 of the 105 tests reach `owner.repos` today. Thirteen tests name it directly, and the rest reach it through `newCompatibilitySession` or `newSessionEnvironmentFixtureWithGoDir`.
An owner repository is only `git init -q -b main` in a temporary directory (`internal/systemtest/owner_test.go:102-112`). Most of those tests use one only because they need a fresh repository.
So each such test takes a private journey repository instead. One helper builds it, and `newSystemOwner` builds its three repositories through the same helper.
Many tests already build a private repository with the same `git init` step, for example `internal/systemtest/land_route_test.go:186`. The profile's "at most three disposable repositories" describes the owner, and it stays true.

Two tests keep `owner.repos`, because the owner's repositories are their subject. `TestSelectedExecutableComposition` runs over all three (`owner_selection_test.go:181`), and `TestOwnerSelectionChild` plants `repos = nil` as a control (`:156`).

A system test calls `t.Parallel` unless it holds a serial reason. A serial reason is one of these:

- The owner's repositories are the test's subject. Only the two tests above hold this reason.
- The test body, or a test-file helper that is not an owner writer, names a ledger field (`starts`, `seen`, or `terminal`).
- The test binds process state, such as `t.Setenv`.

`TestRedTimeoutAndDescendantTeardown` stays whole. Go runs the serial tests first and the parallel tests after them, so a serial mutator never overlaps a parallel test.
The three `t.Setenv("TMPDIR", ...)` calls (`internal/systemtest/owner_artifact_recovery_test.go:24`, `:280`, `:301`) become explicit child environment values. Those children already take an explicit environment, such as `mergeEnvironment(systemBaseEnvironment(), artifactLandEnv(...))` at `:172`.

A new census test in `internal/systemtest` parses the package's test files. It refuses a top-level test that neither calls `t.Parallel` nor appears in a serial table with a reason.
It also refuses a parallel test that reaches `owner.repos` through its body or a test-file helper. The precedent is the worktree census walk (`internal/worktree/parallel_census_test.go:326-360`).
It refuses a parallel test whose body, or a non-owner helper, names a ledger field. The ledger rule is direct, and the walk does not follow owner writers.
An owner writer is a function whose every ledger-field write sits between `owner.mu.Lock()` and `owner.mu.Unlock()`. The census derives this set by that rule, and LGP28 grades the same mutex fact.
Today the rule admits eight writers. Examples are `runWithInput` (`internal/systemtest/owner_test.go:171`), `timeoutProcessGroup` (`:244`), and `startArtifactLand` (`owner_artifact_recovery_test.go:165`).
Every bench child shares one `owner.home` by default (`internal/systemtest/owner_environment_test.go:56-60`). The build runs the parallel set under the default `-parallel` and records any conflict on that home as a serial reason.
The gate argv, the phase table, the `projects/benchkit.md` system sentences, and ADR 0024 stay unchanged.

### Measurement

The baseline is the three newest `main` gates of 2026-10-21, listed under Problem. Each chunk's author records the affected package time from `bench test --full` before and after its change.
The final reconciliation pairs the package and phase timings of the landing gates with each retained red proof. The build records each reverted attempt in the build log.

## Implementation chunks

| Stable chunk ID / tickets | Delivered outcome | Acceptance rows | Tests | Harder chunk |
|---|---|---|---|---|
| LGP-C1 / pending | Probe and testreport fixtures built once for each package run, with every witness kept | LGP1–LGP8, LGP37, LGP34 | probe and testreport package tests | no |
| LGP-C2 / pending | Stable git location answers cached for one command run, and direct duplicates removed | LGP9–LGP17, LGP21, LGP22, LGP34 | git and worktree package tests | yes |
| LGP-C3 / pending | Working-directory and child-environment worktree tests in the parallel set | LGP18–LGP24, LGP34 | worktree census and moved tests | no |
| LGP-C4 / pending | System tests in parallel inside one process, with a serial-reason census | LGP25–LGP36 | system suite and its census | yes |

LGP-C2 and LGP-C3 both edit `internal/worktree/effects.go` and `internal/worktree/parallel_census_test.go`. LGP-C3 therefore follows LGP-C2. The other chunks have no write overlap.
LGP34 belongs to each chunk's author. LGP35 and LGP36 belong to the orchestrator's final reconciliation after the last chunk.
The ticket slicer replaces each `pending` cell with ticket basenames.

### Completion plan

```bench-completion-plan
{"version":1,"chunks":[{"id":"LGP-C1","tickets":["01-share-probe-and-testreport-fixtures.md"],"verification":[{"id":"probe","command":"bench test --package ./internal/probe"},{"id":"testreport","command":"bench test --package ./internal/testreport"},{"id":"root-conformance","command":"bench test --check root-conformance"},{"id":"copy-isolation","command":"bench test --package ./internal/probe","probe":"return the golden root instead of a private copy"},{"id":"golden-use","command":"bench test --package ./internal/probe","probe":"build a fresh repository in clampFixture"}]},{"id":"LGP-C2","tickets":["02-cache-stable-git-locations.md"],"verification":[{"id":"git","command":"bench test --package ./internal/git"},{"id":"worktree","command":"bench test --package ./internal/worktree"},{"id":"root-conformance","command":"bench test --check root-conformance"},{"id":"revalidation","command":"bench test --package ./internal/git","probe":"skip the validation on a cached common directory"}]},{"id":"LGP-C3","tickets":["03-move-cwd-and-child-env-tests.md"],"verification":[{"id":"worktree","command":"bench test --package ./internal/worktree"},{"id":"root-conformance","command":"bench test --check root-conformance"},{"id":"ceiling","command":"bench test --package ./internal/worktree","probe":"leave worktreeSerialCeiling at 45"}]},{"id":"LGP-C4","tickets":["04-run-system-tests-in-parallel.md"],"verification":[{"id":"system","command":"bench test --check system"},{"id":"root-conformance","command":"bench test --check root-conformance"},{"id":"serial-census","command":"bench test --check system","probe":"drop the serial reason of an owner-repository mutator"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/landing-gate-long-poles/spec.md"},{"id":"probe","command":"bench test --package ./internal/probe"},{"id":"testreport","command":"bench test --package ./internal/testreport"},{"id":"git","command":"bench test --package ./internal/git"},{"id":"worktree","command":"bench test --package ./internal/worktree"},{"id":"root-conformance","command":"bench test --check root-conformance"},{"id":"system","command":"bench test --check system"}]}
```

## Testing decisions

A good test drives the real verb or the real `go test` child, and it observes starts, file state, or diagnostics. Elapsed time is evidence for review, not an assertion.
Three structural checks prove the removed duplicate work:

- a start count from the stub git argument log;
- a census of the serial set;
- a fixture build count.

The stub git argument log (`internal/gittest/gittest.go:84`) is the start counter. A line-count assertion over it is new; the hostile tests only assert absence today.
A fixture build counter wraps the fixture publisher or the run-binary builder. It follows the mutex-guarded log of the worktree selector (`internal/worktree/test_run_test.go:33-39`), because `sync.Once` keeps no count.

The system census follows the worktree census walk. It refuses by test name and names the reason table. The ledger negative controls in `TestUnfilteredRunVerifiesLedger` stay unchanged.
The gate observes the result through `bench test --check system`, the package tests, and `bench test --check root-conformance`.

### Seam diagram

    trigger: one probe or testreport test
        │
        ▼
    golden (once per package run) ──▶ [ private copy ] ──▶ real verb or go child
                      ◀ tests attach here: build count, copy isolation, witness reds

    trigger: one worktree command run
        │
        ▼
    root ──▶ [ git run value: common dir, top level, admin dir ] ──▶ validated answer
                      ◀ tests attach here: stub git argument log line counts

    trigger: go test -tags=system (one process)
        │
        ▼
    TestMain owner ──▶ [ serial set, then parallel set ] ──▶ ledger verify
                      ◀ tests attach here: system census, ledger negative controls

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| LGP1 | 1 | Two `clampFixture` calls in one package run invoke the fixture publisher once. | planned TestClampGoldenBuildsOnce in internal/probe/golden_test.go | A fixture that still builds its own repository invokes the publisher a second time. |
| LGP2 | 2 | A write in one fixture copy leaves a second copy's subject bytes unchanged. | planned TestClampCopiesAreIsolated in internal/probe/golden_test.go | A shared root or a hard link shows the write in the second copy. |
| LGP3 | 3 | A fixture copy passes the run-binary freshness check for its own root. | planned TestClampCopyKeepsAFreshSeal in internal/probe/golden_test.go | A stale seal refuses the selection with the freshness diagnostic. |
| LGP4 | 4 | The PB1 witness still runs a real `go test` child and bites. | `internal/probe/probe_test.go` (`TestProbeBitesWhenTheFocusedTestFails`) | A stubbed child or a skipped run loses the packages table. |
| LGP5 | 5 | The PB22 usage witness keeps every refusal case on the shared fixture. | `internal/probe/refusal_test.go` (`TestProbeUsageNamesOneSelectionAndOneMutation`) | A dropped case or a changed usage text fails the comparison. |
| LGP6 | 6 | The two kit-copy tests together invoke the run-binary builder at most once. | planned TestOwnedSelectionBuildsOnce in internal/testreport/selection_once_test.go | A test that still builds its own selection invokes the builder a second time. |
| LGP7 | 7 | The planted row-next-grammar defect still reaches the diagnostic through the real child. | `internal/testreport/root_conformance_test.go` (`TestRootConformanceCheckPrintsFailures`) | A canned child or a reused stale binary loses exit 1 or the diagnostic. |
| LGP8 | 8 | The named check still runs only the registered dev-scope check through the real selection. | `internal/testreport/check_test.go` (`TestNamedCheckRunsOnlyRegisteredDevScope`) | A changed selection path breaks the header row or the timing file. |
| LGP9 | 9 | Two common-directory reads in one run start git once. | planned TestRunCachesTheCommonDir in internal/git/run_cache_test.go | Each extra start adds one line to the stub argument log. |
| LGP10 | 10 | Two top-level reads and two admin-directory reads in one run start git once each. | planned TestRunCachesTopLevelAndAdminDir in internal/git/run_cache_test.go | Each extra start adds one line to the stub argument log. |
| LGP11 | 11 | A common directory replaced by a symlink after the first read refuses on the second read. | planned TestRunRevalidatesACachedCommonDir in internal/git/run_cache_test.go | A cache that skips validation returns the stale answer. |
| LGP12 | 12 | A branch moved by `update-ref` inside one run reads its new commit. | planned TestRunReadsRefsLive in internal/git/run_cache_test.go | A cached ref returns the old commit. |
| LGP13 | 13 | The bad common-directory cases still refuse before the porcelain start. | `internal/git/worktree_admin_hostile_test.go` (`TestWorktreesRejectsBadCommonDirBeforePorcelain`) | A cached package-level answer hides the hostile stub answer. |
| LGP14 | 14 | Two run values over one root each start their own common-directory query. | planned TestRunValuesShareNoCache in internal/git/run_cache_test.go | A process-global cache saves the second start. |
| LGP15 | 15 | A plumbing flag written outside `internal/git` reds the git plumbing owner check. | `internal/conformance/git_plumbing_owner_test.go` (`TestGitPlumbingOwnerRedsARetypedFlag`) | A run value that spells a flag outside `internal/git` fails the registered check. |
| LGP16 | 16 | A direct `planExplicitWith` call for an owned target under a canonical primary root starts the common-directory query exactly three times. | planned TestExplicitPlanReadsEachLocationOnce in internal/worktree/journey_test.go | Each surviving location duplicate adds a line above the one run-value start and the two ledger-read starts. |
| LGP17 | 17 | One `validateOwnerMarker` call under a canonical primary root starts the common-directory query exactly once. | planned TestOwnerMarkerReadsTheCommonDirOnce in internal/worktree/journey_test.go | A `canonicalRoot`, `Worktrees`, or `CommonDir` start outside the run value adds a second line. |
| LGP18 | 18 | A relative operand resolves against the verb call's directory, not the process directory. | planned TestVerbCallDirectoryResolvesOperands in internal/worktree/identifier_operand_test.go | A process-directory read resolves the operand elsewhere and the path differs. |
| LGP19 | 19 | An exec child sees a value from the verb call's base environment and not a process-only value. | planned TestExecChildTakesTheCallEnvironment in internal/worktree/exec_test.go | An `os.Environ()` read shows the process value in the child. |
| LGP20 | 20 | Each moved test calls `t.Parallel` and reaches neither `bindEnv` nor `chdir`. | `internal/worktree/parallel_census_test.go` (`TestSerialSetStaysBelowTheCeiling`) | A moved test that still binds the process lands in the serial set. |
| LGP21 | 21 | The serial set equals 45 minus the moved tests plus the new serial tests. | `internal/worktree/parallel_census_test.go` (`TestSerialSetStaysBelowTheCeiling`) | A ceiling that misses one term reds with the exact-ceiling diagnostic. |
| LGP22 | 22 | The top-level test count equals 805 plus the new top-level tests. | `internal/worktree/parallel_census_test.go` (`TestPackageTestCountPin`) | A removed or merged test leaves the count below the recorded sum. |
| LGP23 | 23 | The removed-directory test stays serial and keeps its real process directory. | planned TestRemovedDirectoryOperandStaysSerial in internal/worktree/parallel_census_test.go | A moved test loses the real `filepath.Abs` failure. |
| LGP24 | 24 | An `os.Environ` read outside `effects.go` reds the effect census. | `internal/worktree/effect_census_test.go` (`TestEffectBoundaryCensus`) | A second environment owner passes an unextended census. |
| LGP25 | 25 | Every system test without a serial reason calls `t.Parallel`. | planned TestSystemParallelCensus in internal/systemtest/parallel_census_test.go | A new test with no call and no reason reds by name. |
| LGP26 | 26 | The package keeps one `TestMain` and three owner repositories. | planned TestSystemParallelCensus in internal/systemtest/parallel_census_test.go | A second owner or a fourth repository changes the parsed facts. |
| LGP27 | 27 | The unfiltered run still names each missing terminal outcome. | `internal/systemtest/owner_selection_test.go` (`TestUnfilteredRunVerifiesLedger`) | A ledger that skips a terminal outcome loses its diagnostic. |
| LGP28 | 28 | Every ledger write outside `TestOwnerSelectionChild` takes the owner mutex. | planned TestLedgerWritesHoldTheOwnerLock in internal/systemtest/parallel_census_test.go | An unguarded write to `starts`, `seen`, or `terminal` reds by name. |
| LGP29 | 29 | A parallel test that reaches `owner.repos`, or names a ledger field outside a mutex-guarded owner writer, reds the census. | planned TestSystemParallelCensus in internal/systemtest/parallel_census_test.go | A parallel test that shares a repository or writes the ledger unguarded passes with no refusal. |
| LGP30 | 30 | No `t.Setenv("TMPDIR"` call remains in the package. | planned TestSystemParallelCensus in internal/systemtest/parallel_census_test.go | Each remaining TMPDIR binding reds by test name. |
| LGP31 | 31 | Each serial table row names a test that exists and a reason. | planned TestSystemSerialTableIsExact in internal/systemtest/parallel_census_test.go | A stale row or a row with no reason reds. |
| LGP32 | 32 | The owner-repository serial reason names exactly `TestSelectedExecutableComposition` and `TestOwnerSelectionChild`, and both reach `owner.repos`. | planned TestSystemSerialTableIsExact in internal/systemtest/parallel_census_test.go | A test kept serial on an owner repository beyond the two subject tests reds by name. |
| LGP33 | 33 | The system phase argv stays the pinned single `go test` call. | `internal/gate/branch_native_phases_test.go` (`TestBenchkitPhasesUseBranchNativeDrivers`) | A shard driver or a changed operand fails the exact argv. |
| LGP34 | 34 | Each chunk records its package time before and after its change. | Review-owned: chunk verification record. | A chunk with no timing cannot support its speed claim. |
| LGP35 | 35 | The build log names each reverted attempt. | Review-owned: build log at the final reconciliation. | A missing entry lets a later run repeat the attempt. |
| LGP36 | 36 | The final report pairs landing-gate timings with each retained red proof. | Review-owned: final implementation reconciliation. | A faster green run alone cannot establish retained quality. |
| LGP37 | 5 | The PS40 to PS43 witness keeps its passing-run guard and its expectation refusals on the shared fixture. | `internal/probe/verdict_test.go` (`TestExpectRequiresTheNamedWitness`) | A dropped subtest or a changed refusal fails the comparison. |

### Edge inventory

- Fixture edges: a write in one copy (LGP2), and a copy's seal (LGP3).
- Git edges: a hostile common directory after a cached read (LGP11, LGP13), a ref write inside the run (LGP12), and two run values in one process (LGP14).
- Worktree edges: a relative operand (LGP18), a removed process directory (LGP23), a process-only environment value (LGP19), and a PATH stub (LGP21 serial-set arithmetic).
- System edges: a test that reaches an owner repository (LGP29, LGP32), a test that binds process state (LGP30), and an unfiltered run that misses one terminal outcome (LGP27).

**Won't handle:** A failed golden build. `clampFixture` fails the test with the build error, as `probetest.Repository` does today, and that remains the surviving caller.

**Won't handle:** A filtered system run. `TestMain` skips the ledger verify under `-run` today (`internal/systemtest/owner_test.go:64`), and that behavior is unchanged.

**Won't handle:** Separate system shard processes. The single-process parallel run keeps one ledger, and `TestUnfilteredRunVerifiesLedger` remains its surviving caller.

**Won't handle:** Phase overlap. ADR 0024 requires a resource census, and the serial scheduler in `internal/gate/runner.go` remains the surviving caller.

**Won't handle:** Cached refs, HEAD, registrations, or default-branch answers. `ResolvedDefault` and the live readers remain the surviving callers.

**Won't handle:** PATH-stub and Home-default worktree tests in the parallel set. `journeyStubGit` and the boundary-default graders remain the surviving callers.

**Won't handle:** Canned children for the witnesses. PB1 and `TestRootConformanceCheckPrintsFailures` remain the surviving callers of the real children.

## Out of scope

| Separate capability | Estimated work | Reason |
|---|---|---|
| System shard processes with a merged ledger file | 9 edits, 3 gate runs | Changes the argv pins, the profile, and the ledger format. |
| Phase overlap of test and system | 6 edits, 3 gate runs | Requires the ADR 0024 resource census. |
| A higher Go package parallelism on this host | 1 edit, 2 gate runs | A reviewer decision on the shared host load. |
| Git cache invalidation at each write | 12 edits, 2 gate runs | Larger stale-answer risk and review surface. |
| Canned children for the witnesses | 4 edits, 1 gate run | Moves detection out of the real child. |
| Earlier conformance source sharing and one-checkout authorization | 20 edits, 4 gate runs | Measured at a few seconds on the tmpfs base. |

## Further notes

### Write collisions with other staged specs

Several staged specs write the same files. This spec does not order them; the commitment order decides.

- `process-lifetime` (FT362) writes `internal/systemtest/owner_artifact_recovery_test.go`, `internal/testreport/command.go`, and several `internal/worktree` production files.
- `shared-test-fixtures` (FT360) writes `internal/systemtest/owner_test.go`, `internal/worktree/parallel_census_test.go`, and `internal/probe/probe_test.go`. Its ticket 08 shares the system fixture identity, which overlaps the journey-repository helper of LGP-C4.
- `cleanup-member-traversal` writes `internal/worktree/clean_*` files and can change the worktree test count.
- `shell-argument-ownership` writes `internal/worktree` operand tests, such as `identifier_operand_test.go`.

### Research inputs

Five read-only slices ran on 2026-10-21, each with two independent readers. The coordinator closed the disagreements with its own reads.
The gate sets `BENCH_RUN_BINARY` and `BENCH_KIT` for every phase (`internal/gate/phases.go:166-180`). `Worktrees` repeats the common-directory query itself (`internal/git/worktree_admin.go:296`).
`canonicalpath.Resolve` reads the process directory for a relative path, and `targetPath` refuses a relative path target (`internal/worktree/path.go:305-306`).
The secondary reviews of 2026-10-09 cite 767 worktree tests. The tree pins 805 (`internal/worktree/parallel_census_test.go:1061`).

### Flagged additions

These promises go beyond the literal decision source. Each one is the structural guard that a closed seam needs to stay red-capable. The reviewer can veto each one.

- The effect census gains `os.Environ` (story 24, LGP24). It guards the one owner of the child base environment.
- The system lock census (story 28, LGP28). It guards the ledger under parallel tests.
- The system serial-reason census and its exact table (stories 31 and 32, LGP25, LGP26, LGP29 to LGP32). It guards the parallel set.
- The removed-directory row (LGP23). It keeps the one process-directory test serial.
- A new `canonicalpath` form that takes the directory as an argument. It keeps the symlink semantics of `Resolve`.
- A `poolkey` form that reads the common directory through the run value. Without it, `canonicalRoot` keeps a duplicate start that the closed decision names.
- One journey-repository helper that `newSystemOwner` and each migrated system test share. It applies the closed decision that parallel tests take private repositories.

### Source sentence to rows

| Closed decision of 2026-10-21 | Rows |
|---|---|
| The system suite gains `t.Parallel` in one process, with one owner, one ledger, and the pinned argv. | LGP25 to LGP33 |
| The git cache holds only the common directory, the top level, and the admin directory for one run. | LGP9 to LGP15 |
| The direct duplicate queries go away. | LGP16, LGP17 |
| The probe and testreport fixtures are built once for each package run, and every witness is kept. | LGP1 to LGP8, LGP37 |
| The working-directory and child-environment worktree groups move into the parallel set. | LGP18 to LGP24 |
| A gain inside the spread counts as no change, and each reverted attempt is recorded. | LGP34 to LGP36 |

### Pre-review proof checklist

- `Cited symbols`: the review of 2026-10-21 resolved every cited symbol and test name at `43d84ec19`.
- `Import edges`: none. The spec adds no cross-package import edge except the run value from `internal/git` into `internal/worktree`, which already imports it.
- `Source-row clauses and occurrences`: none.
- `Promised field labels`: none.
- `Changed-function callers`: `canonicalpath.Resolve` and `poolkey.Canonical` keep their signatures. The callers of the four `os.Environ` sites, of `planExplicitWith`, and of `validateOwnerMarker` are listed under Implementation decisions. The slicer runs `bench consumers` for each.
- `Copy survival`: LGP1 and LGP6 fail when a per-test build survives.
- `Rendered-shape readers`: none. No rendered output changes.
- `Pin operators`: `TestPackageTestCountPin` refuses a count below and a count above the pin, and `TestSerialSetStaysBelowTheCeiling` requires equality (`internal/worktree/parallel_census_test.go:1043-1078`).
- `Entry reads`: the process directory and `os.Environ` move into `effects.go` (LGP18, LGP19, LGP24).
- `Derived expectations`: none. Each count is an independent literal or a recorded sum.
- `Consolidated rules`: none.
- `Quantified obligations`: the coordinator swept each set against `main` at `43d84ec19` on 2026-10-21.
  - "Every ledger write" (LGP28): ten write sites. Eight guarded writers hold eight of them, and `TestOwnerSelectionChild` holds the other two.
  - "No other `os.Environ` read" (LGP24): four production sites, all in top-level files of `internal/worktree`.
  - "Every witness" (LGP4, LGP5, LGP7, LGP8, LGP37): PB1 is the only real-`go` probe test. `installOwnedSelection` has exactly two callers.
  - "At most once" (LGP6): the planted defect is outside the run-binary build inputs (`scripts/go-build.inputs`), so both kit copies share one digest.
  - Common-directory starts (LGP16, LGP17): eight in an owned explicit plan and three in one owner-marker validation, counted by reading.
  - Owner-repository reach (LGP29, LGP32): 13 direct test sites and two helpers, about 40 tests in all. Two tests have the owner repositories as their subject.
  - "No `t.Setenv("TMPDIR"`" (LGP30): three sites.
- `Workflow-step writes`: none.

### What this spec did not verify

- The per-test cost split inside the probe fixture between `git init` and the seal publication.
- Whether each of the 17 operand tests needs only the verb-call directory.
- Whether concurrent bench children that share one `owner.home` conflict. The build records each conflict as a serial reason.
- Whether the probe verb's selection path reaches `freshness.VerifyAt` and its `go list -deps` start.

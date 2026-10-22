# Cut the landing gate's long poles without losing failure detection

Status: staged

Roadmap: FT402

Decision source: Reviewer-confirmed current conversation, 2026-10-21.

Verification log: 6 iteration(s) to accept — the spec round took 4 iterations, because unswept quantified sets produced NB1, NB2, and two sweep corrections; the ticket round took 2, because the slicer skipped written slicing rules and one ticket passed the context limit.

Implementation approval: approved by the reviewer on 2026-10-21, including the 16-ticket breakdown and the declared implementation line.

## Problem

A landing spends about five minutes in the whole-project gate. The three baseline gates on `main` took 296 to 308 seconds.
The test phase took 173 to 174 seconds, and the system phase took 100 to 108 seconds after it.
The scheduler runs one phase at a time, so the system phase adds its whole wall time to the gate (`internal/gate/runner.go:196`).

The earlier FT402 spec, `landing-test-efficiency`, aimed at adoption and repair fixtures. The tmpfs base for the kit test run removed that cost.
On the tmpfs base, `internal/adopt` takes about 22 seconds, `internal/adopt/repairtest` about 3 seconds, and `internal/conformance` about 37 seconds.
That spec would save about 5 to 10 seconds, which is inside the run-to-run spread. The build stopped before code.

The current long poles are these packages:

| Package | Gate time | Shape of the cost |
|---|---:|---|
| `internal/systemtest` | 100 to 108 s | 105 top-level tests in one serial process; no test calls `t.Parallel` |
| `internal/probe` | 88 to 90 s (66 s alone) | 131 tests; each fixture runs `git init` and publishes a seal |
| `internal/worktree` | 70 to 71 s | 45 serial tests; 41.9 percent of git starts repeat a query |
| `internal/testreport` | 65 to 68 s | two tests copy the kit and build a real run binary each |

Sources: the gate logs `.logs/gate-20261021T192218.079145862Z-3443832.jsonl`, `gate-20261021T191248.934370406Z-2945722`, and `gate-20261021T154223.033572571Z-2308306`, with their `.out` files. The S1 to S5 research returns of 2026-10-21 supply the cost shapes.
A fourth gate of 2026-10-21 (`gate-20261021T233107.745850231Z-117051`) took 374 seconds under host load. It shows how wide the run-to-run spread is.

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
- Each parallel system test takes a private repository instead of a shared owner repository.
- The git cache holds only the stable location answers of a root: the common directory, the top level, and the admin directory.
- The probe and testreport cuts share fixtures and keep every witness. No test is removed.
- The serial worktree tests that move are the working-directory group and the child-environment group.
- The earlier `landing-test-efficiency` spec stays staged and unbound, because `conformance-observation-follow-on` depends on its source-observation chunks. FT402 binds to this spec only.
- The Claude author line for this build is opus at xhigh effort.

## User stories

Line: opus / xhigh.
Implementation-line reason: LGP-C4A and LGP-C4B change a serial journey suite with shared repositories, and LGP-C2A to LGP-C2C add a cache under hostile-path tests. Both seams hold real uncertainty, and the reviewer set xhigh on 2026-10-21.
Harder chunks: LGP-C2A, LGP-C2B, LGP-C2C, LGP-C4A, LGP-C4B.

### Shared fixtures in probe and testreport

1. As a maintainer, I want the probe clamp repository built once for each package run, so that each test skips `git init` and the seal publication.
2. As a test author, I want a private copy of the clamp repository for each probe test, so that mutations stay in one test.
3. As a test author, I want each copy to keep the sealed run binary valid, so that the run-binary freshness check still accepts the copy.
4. As a maintainer, I want the PB1 witness to keep its real `go test` child, so that the real toolchain still proves a bite.
5. As a maintainer, I want the PB22 and PS40 to PS43 witnesses to keep their assertions, so that usage and expectation refusals stay pinned.
6. As a maintainer, I want the testreport run binary built once for each package run, so that two kit-copy tests do not each build it.
7. As a maintainer, I want `TestRootConformanceCheckPrintsFailures` to keep its real conformance child, so that the planted defect still reaches the diagnostic.
8. As a maintainer, I want the kit-copy tests to keep the real selection path of `selectRunBinary`, so that an inherited binary is still validated against each kit copy. In the gate the inherited selection wins, so `TestBuildLeavesTheWrapperManifestUntouched` (`internal/runbinary/build_script_test.go:26-57`) keeps the real build under test.

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

The probe test package gains a package-run golden of the clamp repository. It follows the adopt golden precedent (`internal/adopt/golden_test.go:15-76`): one `sync.Once` build and one `TestMain` cleanup.
`clampFixture` and `refusalFixture` take a private copy of the golden for each test. The copy keeps file bytes and modes, and it keeps a valid git repository.
The golden build calls the fixture publisher once. `probetest.Repository` takes that publisher as an argument (`internal/probe/probetest/fixture.go:69`), so a counting publisher observes the builds.
The seal holds source and executable digests and no root path (`internal/freshness/freshness_publish.go:97-111`). So a byte-identical copy keeps a valid seal, when the executable, the seal, and the manifest directory move together.
The publisher also writes a broker manifest whose `path` field binds the golden executable (`internal/brokermanifest/brokermanifest.go:50-55`). No probe or testreport reader reads that manifest, because `runbinary.Validate` checks only the seal (`internal/runbinary/runbinary.go:176-202`). So the stale manifest of a copy is harmless.
The build keeps the environment bindings and the `t.Chdir` of `probetest.Repository` (`internal/probe/probetest/fixture.go:106-112`) for each test.

`TestProbeBitesWhenTheFocusedTestFails` keeps its real `go test` child (PB1). The witnesses PB22 and PS40 to PS43 keep every assertion. The build removes no test from either package.

The testreport package builds its run binary once for each package run through `runbinary.Shared` (`internal/runbinary/shared.go`), as `testRunSelector` does (`internal/worktree/test_run_test.go:19-34`). An inherited `BENCH_RUN_BINARY` wins, and a present invalid selection refuses, by the one rule of `Factory.Select` (`internal/runbinary/runbinary.go:77-85`).
A counting builder wraps the run-binary build, so the two kit-copy tests together show at most one build.
The probe verb's own run-binary check can still run `go list -deps` through `freshness.VerifyAt` (`internal/runbinary/runbinary.go:162-171`). This spec claims no saving from that check.
`TestRootConformanceCheckPrintsFailures` and `TestNamedCheckRunsOnlyRegisteredDevScope` keep their real children and their kit copies. The build reuses a run binary only where the kit-copy sources match the binary's seal.

### Git location answers for one command run

`internal/git` gains a run value that a command creates at its entry and passes down. The worktree `ambient` value carries it (`internal/worktree/effects.go:22-36`).
`newAmbient` and `callAmbient` build a fresh run value for each verb run. No `joins` value holds it, because a test reuses one `joins` value across verb runs (`internal/worktree/verb_runner_test.go:142-146`).
`planExplicitWith` takes the run value as a parameter. Its seven callers pass the value from their `ambient`, or a fresh run value at a public entry.
The run value caches the raw answers of the common-directory, top-level, and admin-directory queries, keyed by the root spelling that the caller passes.
On every read, the run value runs the existing validation again, such as `validateCommonDir` (`internal/git/worktree_admin.go:245-260`). The cache saves the git start, not the filesystem check.
The package-level functions `CommonDir`, `AdminDir`, `Worktrees`, and `Root` stay uncached. The hostile tests in `internal/git/worktree_admin_hostile_test.go` call them directly and keep their meaning.
The run value has its own `Worktrees` form. That form reads the common directory through the run value, and it keeps the admin scan and the porcelain start live.

The run value holds no ref, HEAD, registration, index, default-branch, or ancestry answer. Those reads stay live (S2 research, finding 21 and 22).
No package variable holds a cache. A direct caller with no run value gets a fresh query, as today.

`validateOwnerMarker` (`internal/worktree/classifier.go:396-463`) asks for the common directory three times today. The starts come from `canonicalRoot`, from `Worktrees`, and from `CommonDir`.
`canonicalRoot` (`internal/worktree/worktree.go:100`) calls `poolkey.Canonical`, which runs `git.CommonDir` (`internal/poolkey/poolkey.go:18-23`). So `internal/poolkey` gains a form that reads through the run value, and `Canonical` keeps its signature for its other callers.
`planExplicitWith` (`internal/worktree/subshell.go:160-370`) starts six common-directory queries in an owned plan after ticket 04. Three come from its own `canonicalRoot`, `Worktrees`, and `CommonDir`, and one more from the nested `validateOwnerMarker`, whose fresh run value reads once.
The other two come from the ledger reads `intent.LifecycleEvidence` and `intent.Assignments`, through `intent.Address` (`internal/intent/intent.go:34-38`).
The run value cuts the four location starts to one. The two ledger reads stay live, because the closed decision does not move `internal/intent` onto the run value. The default-branch reads also stay as they are.
The cache key is the root spelling. `canonicalRoot` returns the canonical root, so a caller that passes the canonical primary root hits the same key on each later read.

### Serial worktree tests in the parallel set

`verbCall` (`internal/worktree/verb_runner_test.go:47-55`) gains a working-directory field and a child base environment field.
`canonicalpath.Resolve` reads the process directory for a relative path (`internal/canonicalpath/canonicalpath.go:46-50`). It joins the parts as text and resolves symlinks before it cleans, on purpose.
A new `canonicalpath` form takes the directory as an argument and keeps exactly those semantics. `canonicalpath.Operand` is not that form, because it cleans lexically and resolves no symlink.
`effects.go` reads the process directory once at the verb boundary. A failed read is carried forward and surfaces only when a relative operand needs it.
So an absolute operand still works in a removed directory, and `TestLandingHandsBackAPathWithNoCanonicalForm` (`internal/worktree/land_resume_refusal_test.go:343`) keeps its behavior.

A verb boundary that holds `ambient` makes a relative path operand absolute against the call directory before it calls `planExplicitWith` or `cleanupIdentity`.
So both functions keep their signatures, and `resolveOperand` (`internal/worktree/subshell.go:125-138`) then sees an absolute path.
The boundary leaves an operand that starts with `~` unchanged, because `expandHomeTarget` (`internal/worktree/path.go:276`) owns its refusal.
A ledger-fed caller passes a recorded absolute path, and the public entries `PlanExplicit`, `PlanAutomatic`, and `ApplyExplicit` keep the process directory.

The child base environment enters `ambient` once at the verb boundary in `effects.go`. Four production sites read `os.Environ()` today, and each one takes the ambient value instead.
The four sites are `execEnv` (`internal/worktree/exec.go:201`), `gitInput` (`internal/worktree/clean.go:218`), the doctor child (`internal/worktree/build.go:331`), and the subshell entry (`internal/worktree/subshell.go:40`).

An exported chain that holds no `ambient` takes the base environment from one `effects.go` function, as `currentTime` serves the clock (`internal/worktree/effects.go:23`). That chain is `RunTreeChild` (`exec.go:119`), whose one caller is `internal/treetarget/run.go:66`.
The single-read census refuses an `effects.go` read from an unexported declaration. So `gitInput`, `commitTree`, and `captureLayers` on the reset recovery path take the environment as a parameter, from the `ambient` at `retireCheckout` and `resetWith`, by consultant decision of 2026-10-22.
So `effects.go` stays the one `os.Environ` reader, and no package outside `internal/worktree` changes. The exec entry `ExecCommandResolving` (`exec.go:60`) is a verb boundary, so it builds its own `ambient` value.

The working-directory group holds 17 tests, and the child-environment group holds 10 tests (S3 research). The build moves each test that its seam frees and records each test that stays.
Both pins change by recorded arithmetic. Each ticket that touches the package records each term in its verification excerpt.

- `worktreeTestCount` (`internal/worktree/parallel_census_test.go:1061`) equals 805 plus the new top-level tests.
- `worktreeSerialCeiling` (`parallel_census_test.go:1041`) equals 45, minus the moved tests, plus the new serial tests.

LGP-C2B and LGP-C2C add serial tests, because their start counts use `journeyStubGit`, which binds PATH (`internal/worktree/journey_test.go:64-68`). So both also write `parallel_census_test.go`.
The Home-default graders, the PATH-stub tests, and `TestLandingHandsBackAPathWithNoCanonicalForm` stay serial, because the process state is their subject.
`TestEffectBoundaryCensus` (`internal/worktree/effect_census_test.go:20`) gains `os.Environ` in its forbidden set, because `effects.go` becomes the one owner of that read. No other `os.Environ` read stays in the package's production files.

### System suite tests in parallel

The one `go test -tags=system` process stays. `TestMain`, `newSystemOwner`, the three owner repositories, and `verify` keep their shape (`internal/systemtest/owner_test.go:52-115`, `:350-373`).
Every ledger write in the owner code takes the owner mutex today. The one exception is `TestOwnerSelectionChild` (`internal/systemtest/owner_selection_test.go:156-170`).
That test plants the negative controls in a child process, where every other test is skipped. The lock census exempts it by name, and the controls stay unchanged.

About 40 of the 105 tests reach `owner.repos` today. Fourteen tests name it directly at 15 sites, and the rest reach it through `newCompatibilitySession` or `newSessionEnvironmentFixtureWithGoDir`.
An owner repository is only `git init -q -b main` in a temporary directory (`internal/systemtest/owner_test.go:102-112`). Most of those tests use one only because they need a fresh repository.
So each such test takes a private journey repository instead. One helper builds it, and `newSystemOwner` builds its three repositories through the same helper.
Many tests already build a private repository with the same `git init` step, for example `internal/systemtest/land_route_test.go:186`. The profile's "at most three disposable repositories" describes the owner, and it stays true.

`TestBranchNativeArchitectureCensus` (`internal/conformance/ordinary_build_census_test.go:304`) already pins `func TestMain(`, the `for range 3` loop, and `len(o.repos) != 3` in `owner_test.go` (`:88-99`).
It is the one source of LGP26, so the new census does not parse those facts again.

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
An owner writer is a function that writes at least one ledger field, and whose every ledger-field reference sits between a `Lock` and an `Unlock` of the owner mutex. The owner value is spelled `owner` or `o`, so the lock is `owner.mu.Lock()` or `o.mu.Lock()`. The census derives this set by that rule, and LGP28 grades the same mutex fact.
Today the rule admits eight writers. Examples are `runWithInput` (`internal/systemtest/owner_test.go:171`), `timeoutProcessGroup` (`:244`), and `startArtifactLand` (`owner_artifact_recovery_test.go:165`).
The census matches selectors on the owner value, such as `owner.starts` or `o.starts`, and not bare words. A local `seen` (`internal/systemtest/charge_evidence_test.go:378`) and the strings in `TestUnfilteredRunVerifiesLedger` are not ledger references.
Every bench child shares one `owner.home` by default (`internal/systemtest/owner_environment_test.go:56-60`). The build runs the parallel set under the default `-parallel`.
A test that fails under parallel load keeps a serial row, and its reason quotes the failing run. A shared `owner.home`, a timing bound, or a shared Go build cache are the likely causes.

The census walk has one source. A new non-test package, `internal/testcensus`, holds the test-file parser, the top-level test predicate, and the direct `t.Parallel` predicate.
The worktree census imports it in place of its own copies (`internal/worktree/parallel_census_test.go:52`, `:95`, `:292`), and the system census imports it too.
Each census keeps its own rules, such as the worktree serial helpers and the owner-writer rule.
`internal/gittest` and `internal/testrepo` set the precedent for a non-test package that only tests import.
The gate argv, the phase table, the `projects/benchkit.md` system sentences, and ADR 0024 stay unchanged.

The system migration runs as expand, migrate, and contract. The expand ticket adds the census and lists each unmoved test with the temporary reason "pending migration".
Each migration ticket removes those rows for its files. The contract ticket forbids the temporary reason (LGP38).
`landRouteRepo` (`internal/systemtest/land_route_test.go:180-190`) already builds a private repository with the same `git init` step. It wraps the new journey-repository helper, so the step has one source.
The recovery tests pass the private TMPDIR to every child that writes a bundle. The routes are `artifactLandEnv` and `startArtifactLand` (`owner_artifact_recovery_test.go:161-191`), and `systemLandEnv` for the review children (`:29`, `:59`, `:398`).
The serial table lives in its own new file beside the census, so that neither file passes 400 lines. A new file over its limit reds the growth ratchet, because its base count is zero.

Four system test files are at or over the 400-line limit with no grant. They are `charge_evidence_test.go` (400 lines), `owner_artifact_recovery_test.go` (430), `session_start_test.go` (457), and `land_route_test.go` (508).
`owner_test.go` has 485 lines against its grant of 520 (`.bench/structure.budgets:76`), so 35 lines of headroom remain.
The structure growth ratchet reds a file over its limit that gains lines (`internal/structure/structure.go:278-332`). So each ticket that adds `t.Parallel` lines to these files offsets them or splits the file, which the ratchet never punishes.

### Measurement

The baseline is the three newest `main` gates of 2026-10-21, listed under Problem. Each chunk's author records the affected package time from `bench test --full` before and after its change.
The final reconciliation pairs the package and phase timings of the landing gates with each retained red proof. The build records each reverted attempt in the build log.

## Implementation chunks

| Stable chunk ID / tickets | Delivered outcome | Acceptance rows | Tests | Harder chunk |
|---|---|---|---|---|
| LGP-C1 / 01-share-the-probe-clamp-golden.md, 02-build-the-testreport-selection-once.md | Probe and testreport fixtures built once for each package run, with every witness kept | LGP1–LGP8, LGP37, LGP34 | probe and testreport package tests | no |
| LGP-C2A / 03-add-the-git-run-value.md | A git run value that caches the stable location answers of one command run | LGP9–LGP15, LGP34 | git and poolkey package tests | yes |
| LGP-C2B / 04-carry-the-run-value-and-read-owner-markers-once.md | Each verb run carries one run value, and owner-marker validation reads the common directory once | LGP17, LGP21, LGP22, LGP34 | worktree package tests | yes |
| LGP-C2C / 05-read-the-explicit-plan-locations-once.md | An explicit cleanup plan reads its location answers once | LGP16, LGP21, LGP22, LGP34 | worktree package tests | yes |
| LGP-C3A / 06-carry-the-call-directory.md | A verb run carries its working directory, and relative operands resolve against it | LGP18, LGP22, LGP34 | canonicalpath and worktree package tests | no |
| LGP-C3B / 07-move-the-working-directory-tests.md | The working-directory worktree tests in the parallel set | LGP20–LGP23, LGP34 | worktree census and moved tests | no |
| LGP-C3C / 08-carry-the-child-environment.md | A verb run carries its child base environment, and effects.go owns the environment read | LGP19, LGP22, LGP24, LGP34 | worktree package tests and the effect census | no |
| LGP-C3D / 09-move-the-exec-child-environment-tests.md | The exec child-environment worktree tests in the parallel set | LGP20–LGP22, LGP34 | worktree census and moved tests | no |
| LGP-C4A / 10-add-the-system-parallel-census.md | A system census with a serial-reason table, and one journey-repository helper | LGP25, LGP26, LGP28, LGP29, LGP31, LGP34 | system suite and its census | yes |
| LGP-C4B / 11-run-the-owner-and-recovery-tests-in-parallel.md, 12-run-the-compatibility-and-adoption-tests-in-parallel.md, 13-run-the-session-start-and-guard-tests-in-parallel.md, 14-run-the-land-and-telemetry-tests-in-parallel.md, 15-run-the-exec-and-charge-tests-in-parallel.md, 16-close-the-system-serial-table.md | System tests in parallel inside one process, and a closed serial table | LGP27, LGP30, LGP32–LGP36, LGP38 | system suite and its census | yes |

Each seam-creating ticket is its own chunk, so that its review closes before a consumer starts. LGP-C2A creates the run value, LGP-C2B its carrier, LGP-C3A the call directory, LGP-C3C the child environment, and LGP-C4A the census and the helper.
LGP-C2B through LGP-C3D all edit `internal/worktree/effects.go` and `internal/worktree/parallel_census_test.go`, so they run in order. Each LGP-C4B ticket edits the serial table, so those tickets run in order.
LGP34 belongs to each ticket's author. LGP35 and LGP36 belong to the orchestrator's final reconciliation after the last chunk.
The ticket review of 2026-10-21 split three chunks to meet the context limit. Old-to-new chunk IDs: LGP-C2B became LGP-C2B and LGP-C2C, and LGP-C3 became LGP-C3A, LGP-C3B, LGP-C3C, and LGP-C3D.
Each ticket fits about 150k tokens of author context, by reviewer direction of 2026-10-21.

### Completion plan

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"ft402-landing-gate-long-poles-full-20261021","orchestrator_session":"claude:lgp-orchestrator-20261021","author_limit":1,"assignments":{"01-share-the-probe-clamp-golden.md":[{"session":"claude:lgp-t01-author-20261021","assignment":"f40d417548f09f3896c5ac7a05513606","model":"opus","effort":"high","source":"acbf1803a55dab9674543491d9a416c41955847c","native_ref":"claude-agent:bench-writer-high:lgp-t01-author-20261021"},{"session":"claude:lgp-t01-repair1-20261021","assignment":"f40d417548f09f3896c5ac7a05513606","model":"opus","effort":"high","source":"d55212cfb07b18e0255fdfb07aae7090f2e9dea4","native_ref":"claude-agent:bench-writer-high:lgp-t01-repair1-20261021","predecessor":"claude:lgp-t01-author-20261021","trigger":"user-directed","stopped":"claude-agent:bench-writer-high:lgp-t01-author-20261021","preserved":"c1d58826bf5f83d20d2060411e89d8ba742b7d5b"}],"02-build-the-testreport-selection-once.md":[{"session":"claude:lgp-t02-author-20261021","assignment":"f40d417548f09f3896c5ac7a05513606","model":"opus","effort":"high","source":"dc17fbba8606918a3f68d9a9a307440213750a15","native_ref":"claude-agent:bench-writer-high:lgp-t02-author-20261021"},{"session":"claude:lgp-t02-repair1-20261021","assignment":"f40d417548f09f3896c5ac7a05513606","model":"opus","effort":"high","source":"844bd78914e6125b0dd227571d9310e003374b1a","native_ref":"claude-agent:bench-writer-high:lgp-t02-repair1-20261021","predecessor":"claude:lgp-t02-author-20261021","trigger":"user-directed","stopped":"claude-agent:bench-writer-high:lgp-t02-author-20261021","preserved":"c1d58826bf5f83d20d2060411e89d8ba742b7d5b"}],"03-add-the-git-run-value.md":[{"session":"claude:lgp-t03-author-20261021","assignment":"f40d417548f09f3896c5ac7a05513606","model":"opus","effort":"xhigh","source":"23ff33f091edcf5b67b53a7e71174724fc6973bf","native_ref":"claude-agent:bench-writer-xhigh:lgp-t03-author-20261021"},{"session":"claude:lgp-t03-repair1-20261021","assignment":"f40d417548f09f3896c5ac7a05513606","model":"opus","effort":"xhigh","source":"95552ef9327e30c1a749b3a2a352a819bbfd5619","native_ref":"claude-agent:bench-writer-xhigh:lgp-t03-repair1-20261021","predecessor":"claude:lgp-t03-author-20261021","trigger":"user-directed","stopped":"claude-agent:bench-writer-xhigh:lgp-t03-author-20261021","preserved":"3da1db02832c86be6e8925388dc867a6df1afa97"}],"04-carry-the-run-value-and-read-owner-markers-once.md":[{"session":"claude:lgp-t04-author-20261021","assignment":"f40d417548f09f3896c5ac7a05513606","model":"opus","effort":"xhigh","source":"ce1e87b62d651acdde0db9e1af290a57d7394f28","native_ref":"claude-agent:bench-writer-xhigh:lgp-t04-author-20261021"}],"05-read-the-explicit-plan-locations-once.md":[{"session":"claude:lgp-t05-author-20261021","assignment":"f40d417548f09f3896c5ac7a05513606","model":"opus","effort":"xhigh","source":"2bd9676bccd899790f3ddce780241656446aaa45","native_ref":"claude-agent:bench-writer-xhigh:lgp-t05-author-20261021"},{"session":"claude:lgp-t05-author2-20261021","assignment":"f40d417548f09f3896c5ac7a05513606","model":"opus","effort":"xhigh","source":"83bfa89b7fdc5fadc823cc88d6c400dcee00f814","native_ref":"claude-agent:bench-writer-xhigh:lgp-t05-author2-20261021","predecessor":"claude:lgp-t05-author-20261021","trigger":"terminal-failure","stopped":"claude-agent:bench-writer-xhigh:lgp-t05-author-20261021","preserved":"83bfa89b7fdc5fadc823cc88d6c400dcee00f814"}],"06-carry-the-call-directory.md":[{"session":"claude:lgp-t06-author-20261022","assignment":"f40d417548f09f3896c5ac7a05513606","model":"opus","effort":"high","source":"49b873d483d574b590166bc415f2000afc958cac","native_ref":"claude-agent:bench-writer-high:lgp-t06-author-20261022"},{"session":"claude:lgp-t06-repair1-20261022","assignment":"f40d417548f09f3896c5ac7a05513606","model":"opus","effort":"high","source":"6b3b9a15e9dfbbcb7e295380dc1eaa4d30a85550","native_ref":"claude-agent:bench-writer-high:lgp-t06-repair1-20261022","predecessor":"claude:lgp-t06-author-20261022","trigger":"user-directed","stopped":"claude-agent:bench-writer-high:lgp-t06-author-20261022","preserved":"a09abf90400cc7cfc1327eca0076e577b517f241"},{"session":"claude:lgp-t06-repair2-20261022","assignment":"f40d417548f09f3896c5ac7a05513606","model":"opus","effort":"high","source":"b11587d0e75292e6a0ff1bd1a4f28349c641b86f","native_ref":"claude-agent:bench-writer-high:lgp-t06-repair2-20261022","predecessor":"claude:lgp-t06-repair1-20261022","trigger":"user-directed","stopped":"claude-agent:bench-writer-high:lgp-t06-repair1-20261022","preserved":"53fe8f83b9fc6ecdb08fcf76f122623b96fda0a4"}],"07-move-the-working-directory-tests.md":[{"session":"claude:lgp-t07-author-20261022","assignment":"f40d417548f09f3896c5ac7a05513606","model":"opus","effort":"high","source":"645b0ef6eb912575113c636d131ed92e2c1f1d50","native_ref":"claude-agent:bench-writer-high:lgp-t07-author-20261022"}],"08-carry-the-child-environment.md":[{"session":"claude:lgp-t08-author-20261022","assignment":"f40d417548f09f3896c5ac7a05513606","model":"opus","effort":"high","source":"10e5ba796a8dd686807dbce4ff86748134ad1cff","native_ref":"claude-agent:bench-writer-high:lgp-t08-author-20261022"}],"09-move-the-exec-child-environment-tests.md":[],"10-add-the-system-parallel-census.md":[],"11-run-the-owner-and-recovery-tests-in-parallel.md":[],"12-run-the-compatibility-and-adoption-tests-in-parallel.md":[],"13-run-the-session-start-and-guard-tests-in-parallel.md":[],"14-run-the-land-and-telemetry-tests-in-parallel.md":[],"15-run-the-exec-and-charge-tests-in-parallel.md":[],"16-close-the-system-serial-table.md":[]}},"chunks":[{"id":"LGP-C1","tickets":["01-share-the-probe-clamp-golden.md","02-build-the-testreport-selection-once.md"],"verification":[{"id":"t01-pkg-cmd-bench","command":"bench test --package ./cmd/bench","ticket":"01-share-the-probe-clamp-golden.md"},{"id":"t01-pkg-internal-probe-probetest","command":"bench test --package ./internal/probe/probetest","ticket":"01-share-the-probe-clamp-golden.md"},{"id":"t01-pkg-internal-reviewrecord-recordcmd","command":"bench test --package ./internal/reviewrecord/recordcmd","ticket":"01-share-the-probe-clamp-golden.md"},{"id":"t01-probe","command":"bench test --package ./internal/probe","ticket":"01-share-the-probe-clamp-golden.md"},{"id":"t01-root-conformance","command":"bench test --check root-conformance","ticket":"01-share-the-probe-clamp-golden.md"},{"id":"t01-golden-use","command":"bench test --package ./internal/probe","probe":"build a fresh repository in clampFixture","ticket":"01-share-the-probe-clamp-golden.md"},{"id":"t01-copy-isolation","command":"bench test --package ./internal/probe","probe":"return the golden root instead of a private copy","ticket":"01-share-the-probe-clamp-golden.md"},{"id":"t02-testreport","command":"bench test --package ./internal/testreport","ticket":"02-build-the-testreport-selection-once.md"},{"id":"t02-root-conformance","command":"bench test --check root-conformance","ticket":"02-build-the-testreport-selection-once.md"},{"id":"t02-selection-once","command":"bench test --package ./internal/testreport","probe":"build the selection in each installOwnedSelection call","ticket":"02-build-the-testreport-selection-once.md"},{"id":"t01-testrepo","command":"bench test --package ./internal/testrepo","ticket":"01-share-the-probe-clamp-golden.md"},{"id":"t01-adopt","command":"bench test --package ./internal/adopt","ticket":"01-share-the-probe-clamp-golden.md"},{"id":"t01-copy-modes","command":"bench test --package ./internal/testrepo","probe":"drop the mode restore of the shared tree copy; the copy-mode test must red","ticket":"01-share-the-probe-clamp-golden.md"},{"id":"t02-runbinary","command":"bench test --package ./internal/runbinary","ticket":"02-build-the-testreport-selection-once.md"},{"id":"t02-worktree","command":"bench test --package ./internal/worktree","ticket":"02-build-the-testreport-selection-once.md"},{"id":"t02-select-inherit","command":"bench test --package ./internal/runbinary","probe":"ignore a present inherited value in Factory.Select; the inherited-selection test must red","ticket":"02-build-the-testreport-selection-once.md"},{"id":"t02-pkg-cmd-bench","command":"bench test --package ./cmd/bench","ticket":"02-build-the-testreport-selection-once.md"}]},{"id":"LGP-C2A","tickets":["03-add-the-git-run-value.md"],"verification":[{"id":"t03-git","command":"bench test --package ./internal/git","ticket":"03-add-the-git-run-value.md"},{"id":"t03-poolkey","command":"bench test --package ./internal/poolkey","ticket":"03-add-the-git-run-value.md"},{"id":"t03-root-conformance","command":"bench test --check root-conformance","ticket":"03-add-the-git-run-value.md"},{"id":"t03-revalidation","command":"bench test --package ./internal/git","probe":"skip the validation on a cached common directory","ticket":"03-add-the-git-run-value.md"},{"id":"t03-gittest","command":"bench test --package ./internal/gittest","ticket":"03-add-the-git-run-value.md"},{"id":"t03-admin-key","command":"bench test --package ./internal/git","probe":"key Run.AdminDir under the common-directory location; the linked-worktree test must red","ticket":"03-add-the-git-run-value.md"},{"id":"t03-failure-memo","command":"bench test --package ./internal/git","probe":"store a failed answer in Run.answer; the failed-query test must red","ticket":"03-add-the-git-run-value.md"}]},{"id":"LGP-C2B","tickets":["04-carry-the-run-value-and-read-owner-markers-once.md"],"verification":[{"id":"t04-pkg-cmd-bench","command":"bench test --package ./cmd/bench","ticket":"04-carry-the-run-value-and-read-owner-markers-once.md"},{"id":"t04-worktree","command":"bench test --package ./internal/worktree","ticket":"04-carry-the-run-value-and-read-owner-markers-once.md"},{"id":"t04-root-conformance","command":"bench test --check root-conformance","ticket":"04-carry-the-run-value-and-read-owner-markers-once.md"},{"id":"t04-marker-run-value","command":"bench test --package ./internal/worktree","probe":"skip the run value in validateOwnerMarker","ticket":"04-carry-the-run-value-and-read-owner-markers-once.md"}]},{"id":"LGP-C2C","tickets":["05-read-the-explicit-plan-locations-once.md"],"verification":[{"id":"t05-pkg-cmd-bench","command":"bench test --package ./cmd/bench","ticket":"05-read-the-explicit-plan-locations-once.md"},{"id":"t05-worktree","command":"bench test --package ./internal/worktree","ticket":"05-read-the-explicit-plan-locations-once.md"},{"id":"t05-root-conformance","command":"bench test --check root-conformance","ticket":"05-read-the-explicit-plan-locations-once.md"},{"id":"t05-plan-run-value","command":"bench test --package ./internal/worktree","probe":"read the common directory again in planExplicitWith","ticket":"05-read-the-explicit-plan-locations-once.md"}]},{"id":"LGP-C3A","tickets":["06-carry-the-call-directory.md"],"verification":[{"id":"t06-pkg-cmd-bench","command":"bench test --package ./cmd/bench","ticket":"06-carry-the-call-directory.md"},{"id":"t06-worktree","command":"bench test --package ./internal/worktree","ticket":"06-carry-the-call-directory.md"},{"id":"t06-canonicalpath","command":"bench test --package ./internal/canonicalpath","ticket":"06-carry-the-call-directory.md"},{"id":"t06-root-conformance","command":"bench test --check root-conformance","ticket":"06-carry-the-call-directory.md"},{"id":"t06-call-directory","command":"bench test --package ./internal/worktree","probe":"resolve a relative operand against the process directory","ticket":"06-carry-the-call-directory.md"},{"id":"t06-routed-sites","command":"bench test --package ./internal/worktree","probe":"resolve the land, land resume, and reauthorize operands against the process directory; the LGP18 test must red","ticket":"06-carry-the-call-directory.md"},{"id":"t06-clean-failed-read","command":"bench test --package ./internal/worktree","probe":"drop the carried directory-read failure on the clean path; the clean failed-read test must red","ticket":"06-carry-the-call-directory.md"},{"id":"t06-apply-call-directory","command":"bench test --package ./internal/worktree","probe":"gate the clean call-directory join to the plan posture; the clean --apply case of the LGP18 test must red","ticket":"06-carry-the-call-directory.md"}]},{"id":"LGP-C3B","tickets":["07-move-the-working-directory-tests.md"],"verification":[{"id":"t07-pkg-cmd-bench","command":"bench test --package ./cmd/bench","ticket":"07-move-the-working-directory-tests.md"},{"id":"t07-worktree","command":"bench test --package ./internal/worktree","ticket":"07-move-the-working-directory-tests.md"},{"id":"t07-root-conformance","command":"bench test --check root-conformance","ticket":"07-move-the-working-directory-tests.md"},{"id":"t07-ceiling","command":"bench test --package ./internal/worktree","probe":"leave the ceiling at its ticket-06 value","ticket":"07-move-the-working-directory-tests.md"}]},{"id":"LGP-C3C","tickets":["08-carry-the-child-environment.md"],"verification":[{"id":"t08-pkg-cmd-bench","command":"bench test --package ./cmd/bench","ticket":"08-carry-the-child-environment.md"},{"id":"t08-worktree","command":"bench test --package ./internal/worktree","ticket":"08-carry-the-child-environment.md"},{"id":"t08-root-conformance","command":"bench test --check root-conformance","ticket":"08-carry-the-child-environment.md"},{"id":"t08-environ-owner","command":"bench test --package ./internal/worktree","probe":"restore the os.Environ read in execEnv","ticket":"08-carry-the-child-environment.md"}]},{"id":"LGP-C3D","tickets":["09-move-the-exec-child-environment-tests.md"],"verification":[{"id":"t09-pkg-cmd-bench","command":"bench test --package ./cmd/bench","ticket":"09-move-the-exec-child-environment-tests.md"},{"id":"t09-worktree","command":"bench test --package ./internal/worktree","ticket":"09-move-the-exec-child-environment-tests.md"},{"id":"t09-root-conformance","command":"bench test --check root-conformance","ticket":"09-move-the-exec-child-environment-tests.md"},{"id":"t09-exec-bind","command":"bench test --package ./internal/worktree","probe":"bind the caller environment again in childOutputAtHome","ticket":"09-move-the-exec-child-environment-tests.md"}]},{"id":"LGP-C4A","tickets":["10-add-the-system-parallel-census.md"],"verification":[{"id":"t10-pkg-cmd-bench","command":"bench test --package ./cmd/bench","ticket":"10-add-the-system-parallel-census.md"},{"id":"t10-worktree","command":"bench test --package ./internal/worktree","ticket":"10-add-the-system-parallel-census.md"},{"id":"t10-testcensus","command":"bench test --package ./internal/testcensus","ticket":"10-add-the-system-parallel-census.md"},{"id":"t10-system","command":"bench test --check system","ticket":"10-add-the-system-parallel-census.md"},{"id":"t10-root-conformance","command":"bench test --check root-conformance","ticket":"10-add-the-system-parallel-census.md"},{"id":"t10-serial-census","command":"bench test --check system","probe":"drop the serial reason of an owner-repository mutator","ticket":"10-add-the-system-parallel-census.md"}]},{"id":"LGP-C4B","tickets":["11-run-the-owner-and-recovery-tests-in-parallel.md","12-run-the-compatibility-and-adoption-tests-in-parallel.md","13-run-the-session-start-and-guard-tests-in-parallel.md","14-run-the-land-and-telemetry-tests-in-parallel.md","15-run-the-exec-and-charge-tests-in-parallel.md","16-close-the-system-serial-table.md"],"verification":[{"id":"t11-system","command":"bench test --check system","ticket":"11-run-the-owner-and-recovery-tests-in-parallel.md"},{"id":"t11-root-conformance","command":"bench test --check root-conformance","ticket":"11-run-the-owner-and-recovery-tests-in-parallel.md"},{"id":"t11-tmpdir-child","command":"bench test --check system","probe":"drop TMPDIR from the startArtifactLand child environment","ticket":"11-run-the-owner-and-recovery-tests-in-parallel.md"},{"id":"t12-system","command":"bench test --check system","ticket":"12-run-the-compatibility-and-adoption-tests-in-parallel.md"},{"id":"t12-root-conformance","command":"bench test --check root-conformance","ticket":"12-run-the-compatibility-and-adoption-tests-in-parallel.md"},{"id":"t12-compatibility-repo","command":"bench test --check system","probe":"point newCompatibilitySession back at owner.repos[0]","ticket":"12-run-the-compatibility-and-adoption-tests-in-parallel.md"},{"id":"t13-system","command":"bench test --check system","ticket":"13-run-the-session-start-and-guard-tests-in-parallel.md"},{"id":"t13-root-conformance","command":"bench test --check root-conformance","ticket":"13-run-the-session-start-and-guard-tests-in-parallel.md"},{"id":"t13-session-repo","command":"bench test --check system","probe":"point newSessionEnvironmentFixtureWithGoDir back at owner.repos[2]","ticket":"13-run-the-session-start-and-guard-tests-in-parallel.md"},{"id":"t14-system","command":"bench test --check system","ticket":"14-run-the-land-and-telemetry-tests-in-parallel.md"},{"id":"t14-root-conformance","command":"bench test --check root-conformance","ticket":"14-run-the-land-and-telemetry-tests-in-parallel.md"},{"id":"t14-otel-parallel","command":"bench test --check system","probe":"remove t.Parallel from one otel test","ticket":"14-run-the-land-and-telemetry-tests-in-parallel.md"},{"id":"t15-system","command":"bench test --check system","ticket":"15-run-the-exec-and-charge-tests-in-parallel.md"},{"id":"t15-root-conformance","command":"bench test --check root-conformance","ticket":"15-run-the-exec-and-charge-tests-in-parallel.md"},{"id":"t15-exec-parallel","command":"bench test --check system","probe":"remove t.Parallel from one exec-bound test","ticket":"15-run-the-exec-and-charge-tests-in-parallel.md"},{"id":"t16-system","command":"bench test --check system","ticket":"16-close-the-system-serial-table.md"},{"id":"t16-root-conformance","command":"bench test --check root-conformance","ticket":"16-close-the-system-serial-table.md"},{"id":"t16-pending-reason","command":"bench test --check system","probe":"keep one pending migration row","ticket":"16-close-the-system-serial-table.md"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/landing-gate-long-poles/spec.md"},{"id":"probe","command":"bench test --package ./internal/probe"},{"id":"testreport","command":"bench test --package ./internal/testreport"},{"id":"git","command":"bench test --package ./internal/git"},{"id":"poolkey","command":"bench test --package ./internal/poolkey"},{"id":"worktree","command":"bench test --package ./internal/worktree"},{"id":"canonicalpath","command":"bench test --package ./internal/canonicalpath"},{"id":"root-conformance","command":"bench test --check root-conformance"},{"id":"system","command":"bench test --check system"}]}
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
| LGP1 | 1 | Two `clampFixture` calls in one package run invoke the fixture publisher once. | `internal/probe/golden_test.go` (`TestClampGoldenBuildsOnce`) | A fixture that still builds its own repository invokes the publisher a second time. |
| LGP2 | 2 | A write in one fixture copy leaves a second copy's subject bytes unchanged. | `internal/probe/golden_test.go` (`TestClampCopiesAreIsolated`) | A shared root or a hard link shows the write in the second copy. |
| LGP3 | 3 | A fixture copy passes the run-binary freshness check for its own root. | `internal/probe/golden_test.go` (`TestClampCopyKeepsAFreshSeal`) | A stale seal refuses the selection with the freshness diagnostic. |
| LGP4 | 4 | The PB1 witness still runs a real `go test` child and bites. | `internal/probe/probe_test.go` (`TestProbeBitesWhenTheFocusedTestFails`) | A stubbed child or a skipped run loses the packages table. |
| LGP5 | 5 | The PB22 usage witness keeps every refusal case on the shared fixture. | `internal/probe/refusal_test.go` (`TestProbeUsageNamesOneSelectionAndOneMutation`) | A dropped case or a changed usage text fails the comparison. |
| LGP6 | 6 | The two kit-copy tests together invoke the run-binary builder at most once. | `internal/testreport/selection_once_test.go` (`TestOwnedSelectionBuildsOnce`) | A helper that builds on each `installOwnedSelection` call invokes the builder a second time. |
| LGP7 | 7 | The planted row-next-grammar defect still reaches the diagnostic through the real child. | `internal/testreport/root_conformance_test.go` (`TestRootConformanceCheckPrintsFailures`) | A canned child or a reused stale binary loses exit 1 or the diagnostic. |
| LGP8 | 8 | The named check still runs only the registered dev-scope check through the real selection. | `internal/testreport/check_test.go` (`TestNamedCheckRunsOnlyRegisteredDevScope`) | A changed selection path breaks the header row or the timing file. |
| LGP9 | 9 | Two common-directory reads in one run start git once. | `internal/git/run_cache_test.go` (`TestRunCachesTheCommonDir`) | Each extra start adds one line to the stub argument log. |
| LGP10 | 10 | Two top-level reads and two admin-directory reads in one run start git once each. | `internal/git/run_cache_test.go` (`TestRunCachesTopLevelAndAdminDir`) | Each extra start adds one line to the stub argument log. |
| LGP11 | 11 | A common directory replaced by a symlink after the first read refuses on the second read. | `internal/git/run_cache_test.go` (`TestRunRevalidatesACachedCommonDir`) | A cache that skips validation returns the stale answer. |
| LGP12 | 12 | A branch moved by `update-ref` inside one run reads its new commit. | `internal/git/run_cache_test.go` (`TestRunReadsRefsLive`) | A cached ref returns the old commit. |
| LGP13 | 13 | The bad common-directory cases still refuse before the porcelain start. | `internal/git/worktree_admin_hostile_test.go` (`TestWorktreesRejectsBadCommonDirBeforePorcelain`) | A cached package-level answer hides the hostile stub answer. |
| LGP14 | 14 | Two run values over one root each start their own common-directory query. | `internal/git/run_cache_test.go` (`TestRunValuesShareNoCache`) | A process-global cache saves the second start. |
| LGP15 | 15 | A plumbing flag written outside `internal/git` reds the git plumbing owner check. | `internal/conformance/git_plumbing_owner_test.go` (`TestGitPlumbingOwnerRedsARetypedFlag`) | A run value that spells a flag outside `internal/git` fails the registered check. |
| LGP16 | 16 | A direct `planExplicitWith` call for an owned target under a canonical primary root starts the common-directory query exactly three times. | `internal/worktree/journey_test.go` (`TestExplicitPlanReadsEachLocationOnce`) | Each surviving location duplicate adds a line above the one run-value start and the two ledger-read starts. |
| LGP17 | 17 | One `validateOwnerMarker` call under a canonical primary root starts the common-directory query exactly once. | `internal/worktree/journey_test.go` (`TestOwnerMarkerReadsTheCommonDirOnce`) | A `canonicalRoot`, `Worktrees`, or `CommonDir` start outside the run value adds a second line. |
| LGP18 | 18 | A relative operand resolves against the verb call's directory, not the process directory. | `internal/worktree/identifier_operand_test.go` (`TestVerbCallDirectoryResolvesOperands`) | A process-directory read resolves the operand elsewhere and the path differs. |
| LGP19 | 19 | An exec child sees a value from the verb call's base environment and not a process-only value. | planned TestExecChildTakesTheCallEnvironment in internal/worktree/exec_call_env_test.go | An `os.Environ()` read shows the process value in the child. |
| LGP20 | 20 | Each moved test calls `t.Parallel` and reaches neither `bindEnv` nor `chdir`. | `internal/worktree/parallel_census_test.go` (`TestParallelCensusOnTheLiveTree`) | A moved test that still binds the process, or that lacks `t.Parallel`, fails the live census. |
| LGP21 | 21 | The serial set equals 45 minus the moved tests plus the new serial tests. | `internal/worktree/parallel_census_test.go` (`TestSerialSetStaysBelowTheCeiling`) | A ceiling that misses one term reds with the exact-ceiling diagnostic. |
| LGP22 | 22 | The top-level test count equals 805 plus the new top-level tests. | `internal/worktree/parallel_census_test.go` (`TestPackageTestCountPin`) | A removed or merged test leaves the count below the recorded sum. |
| LGP23 | 23 | The removed-directory test stays serial and keeps its real process directory. | `internal/worktree/serial_ceiling_test.go` (`TestRemovedDirectoryOperandStaysSerial`) | A moved test loses the real `filepath.Abs` failure. |
| LGP24 | 24 | An `os.Environ` read outside `effects.go` reds the effect census. | `internal/worktree/effect_census_test.go` (`TestEffectBoundaryCensus`) | A second environment owner passes an unextended census. |
| LGP25 | 25 | Every system test without a serial reason calls `t.Parallel`. | planned TestSystemParallelCensus in internal/systemtest/parallel_census_test.go | A new test with no call and no reason reds by name. |
| LGP26 | 26 | The package keeps one `TestMain` and three owner repositories. | `internal/conformance/ordinary_build_census_test.go` (`TestBranchNativeArchitectureCensus`) | A removed `TestMain` or a fourth repository drops a pinned literal of `owner_test.go`, and the check reds. |
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
| LGP38 | 31 | No serial table row holds the temporary reason "pending migration". | planned TestSystemSerialTableIsExact in internal/systemtest/parallel_census_test.go | A migration left unfinished keeps its temporary row and reds by name. |

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
- `shared-test-fixtures` (FT360) writes `internal/systemtest/owner_test.go`, `internal/worktree/parallel_census_test.go`, and `internal/probe/probe_test.go`. Its ticket 08 shares the system fixture identity, which overlaps the journey-repository helper of LGP-C4A. That ticket also writes `internal/systemtest/bench_follow_on_test.go`, which ticket 13 writes.
- `cleanup-member-traversal` writes `internal/worktree/clean_*` files and can change the worktree test count.
- `shell-argument-ownership` writes `internal/worktree` operand tests, such as `identifier_operand_test.go`, and `internal/systemtest/tree_target_test.go`.
- `exact-planning-readers` writes `internal/systemtest/tree_target_test.go`.
- `production-test-seam-policy` writes `internal/systemtest/compatibility_test.go`, which holds its PS40 witness.
- `process-lifetime` (FT362) also adds new `internal/systemtest` test files, which then fall under the system census.

### Research inputs

Five read-only slices ran on 2026-10-21, each with two independent readers. The coordinator closed the disagreements with its own reads.
The gate sets `BENCH_RUN_BINARY` and `BENCH_KIT` for every phase (`internal/gate/phases.go:166-180`). `Worktrees` repeats the common-directory query itself (`internal/git/worktree_admin.go:296`).
`canonicalpath.Resolve` reads the process directory for a relative path, and `targetPath` refuses a relative path target (`internal/worktree/path.go:305-306`).
The secondary reviews of 2026-10-09 cite 767 worktree tests. The tree pins 805 (`internal/worktree/parallel_census_test.go:1061`).

### Flagged additions

These promises go beyond the literal decision source. Each one is the structural guard that a closed seam needs to stay red-capable. The reviewer can veto each one.

- The effect census gains `os.Environ` (story 24, LGP24). It guards the one owner of the child base environment.
- The system lock census (story 28, LGP28). It guards the ledger under parallel tests.
- The system serial-reason census and its exact table (stories 31 and 32, LGP25, LGP29 to LGP32). It guards the parallel set.
- One `internal/testcensus` package that both censuses import. It applies the one-source rule to the census walk, by consultant consensus of 2026-10-21.
- One `effects.go` environment function for the exported chain that holds no `ambient`. It keeps `effects.go` the one `os.Environ` reader without a change outside `internal/worktree`.
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
- `Import edges`: the run value from `internal/git` into `internal/worktree`, which already imports it, and the new `internal/testcensus` into the test files of `internal/worktree` and `internal/systemtest`.
- `Source-row clauses and occurrences`: none.
- `Promised field labels`: none.
- `Changed-function callers`: `canonicalpath.Resolve` and `poolkey.Canonical` keep their signatures. `planExplicitWith` gains the run value as a parameter, and its seven callers are `classifier.go:267`, `subshell.go:141`, `subshell.go:155`, `clean_landed.go:145`, `clean_set.go:191`, `worktree.go:344`, and `resume.go:86`. The `validateOwnerMarker` callers are `identity_component.go:70`, `reauthorize.go:83`, and `subshell.go:225`. The `ownerMarkerRefusal` callers are `identity_component.go:99`, `ownership.go:204`, and `reset.go:158`.
  The `os.Environ` sites are reached through `runWorktreeChild` and `RunTreeChild` (`exec.go:111`, `:119`), `commitTree` and `captureLayers` (`clean.go:212`, `layers.go:149`, `:153`), `installAttributed` (`build.go:158`), and the subshell entry (`subshell.go:40`). Each author runs `bench consumers` for each changed function.
- `Copy survival`: LGP1 and LGP6 fail when a per-test build survives.
- `Rendered-shape readers`: none. No rendered output changes.
- `Pin operators`: `TestPackageTestCountPin` refuses a count below and a count above the pin, and `TestSerialSetStaysBelowTheCeiling` requires equality (`internal/worktree/parallel_census_test.go:1043-1078`).
- `Entry reads`: the process directory and `os.Environ` move into `effects.go` (LGP18, LGP19, LGP24).
- `Derived expectations`: none. Each count is an independent literal or a recorded sum.
- `Consolidated rules`: none.
- `Quantified obligations`: the coordinator swept each set against `main` at `43d84ec19` on 2026-10-21.
  - "Every ledger write" (LGP28): ten write sites. Eight guarded writers hold eight of them, and `TestOwnerSelectionChild` holds the other two.
  - "No other `os.Environ` read" (LGP24): four production sites, all in top-level files of `internal/worktree`.
  - "Every witness" (LGP4, LGP5, LGP7, LGP8, LGP37): PB1 is the real-`go` witness of LGP4. `TestProbeOmitsTheMatchOnce`, `TestProbeIsSilentWhenTheFocusedTestPasses`, and `TestProbeIsInvalidWhenTheMutationDoesNotCompile` (`internal/probe/probe_test.go:206-224`) also run a real `go test` child on the fixture. `installOwnedSelection` has exactly two callers.
  - "At most once" (LGP6): the planted defect is outside the run-binary build inputs (`scripts/go-build.inputs`), so both kit copies share one digest.
  - Common-directory starts (LGP16, LGP17): eight in an owned explicit plan and three in one owner-marker validation, counted by reading before ticket 04. After ticket 04, the plan's red test observed six.
  - Owner-repository reach (LGP29, LGP32): 15 direct sites in 14 tests and two helpers, about 40 tests in all. Two tests have the owner repositories as their subject.
  - "No `t.Setenv("TMPDIR"`" (LGP30): three sites.
- `Workflow-step writes`: none.

### What this spec did not verify

- The per-test cost split inside the probe fixture between `git init` and the seal publication.
- Whether each of the 17 operand tests needs only the verb-call directory.
- Whether concurrent bench children that share one `owner.home` conflict. The build records each conflict as a serial reason.
- Whether the probe verb's selection path reaches `freshness.VerifyAt` and its `go list -deps` start.

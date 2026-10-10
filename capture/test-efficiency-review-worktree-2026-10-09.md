## Result

I analysed `internal/worktree` as the largest current chunk. The system suite is a close second, so I also report its measured per-test times and one large measured saving. I wrote no file in the repository (`git status --short` was empty at the end). All scratch files sit under `/tmp/claude-1000/-home-mgibs-workspace-bench/c5ce0343-7710-43c7-ab2e-2a1e309edc1e/scratchpad/eff2/` (118 MB).

**Why `internal/worktree`:** wall time is a tie inside the load noise, so I chose on resource use.

| Chunk | Wall under tmpfs (each run) | CPU (user+sys) | Shape |
|---|---|---|---|
| `internal/worktree` | `bench test`: 78.0 s, 61.2 s. Direct: 84.1 s, 80.1 s | 259-320 s | 1,385 tests, 767 top level, 722 already `t.Parallel` |
| System suite | `bench test`: 75.4 s, 87.2 s. Direct: 76.9 s | 84 s direct, 142-154 s via `bench test` | 167 tests, 91 top level, fully serial |

- The pairs did not agree on order: worktree was larger in the 22:27/22:30 pair, and the system suite was larger in the 22:54/22:55 pair.
- Worktree uses 2 to 4 times the CPU. In a five-package concurrent run, its 259-320 s was about 34 to 42% of the 755 CPU-seconds.
- The system suite is one serial phase, so its wall time adds one-for-one to the gate. A reviewer who ranks by serial gate wall can swap the two. The system results are in the lead list, the system notes, and Measurements.

**Host load:** every number is under load. Load average was 4 to 16, because other sessions and gates ran at the same time. I compare same-run numbers where I can. I label each estimate.

## Ranked list

Savings are against the tmpfs baselines above.

| # | Change | Chunk | Saving | Effort | Conf. |
|---|---|---|---|---|---|
| 1 | Shard the system suite over 4 processes, or add `t.Parallel` after an audit of the shared owner repos | System | **Measured bound:** 25.9 s wall for 4 concurrent shards, against 75-87 s serial. That saves about 50-60 s of a serial gate phase. All 4 shards passed. | M | 7 |
| 2 | Cache pure git queries per command run (location, content-addressed reads). `internal/git/worktree_admin.go:193-199` `CommonDir` is the largest source. | worktree | **Estimate:** 18-28 s of the 61-84 s package, about 25-35% of CPU. The measured facts behind it are in Table 2 row B. | M, production code | 5 |
| 3 | Fold the duplicate landings in Table 1 (rows 1-5, 8, 9) | worktree | **Measured sub-times:** about 11 s of the 272 s parallel sum for rows of confidence 6 or more, and 14.6 s with rows 6, 9 and 10. Estimate of 3-4 s wall at `-parallel=4`. | S-M | 6-8 |
| 4 | Move the 45 serial tests into the parallel phase | worktree | **Measured:** the serial head is 9.4 s (12% of wall at p4). **Estimate:** 6-7 s saved. | M | 6 |
| 5 | `t.Parallel` on the subtests of 75 serial-chain tests (13.4 s, 5.2 s and 4.6 s are the longest) | worktree | **Measured:** the 14-subtest test ended at 35.1 s of 35.6 s at `-parallel=16`. **Estimate:** 3-5 s at p16, about 0-1 s at p4. | S | 5 |
| 6 | Raise `-parallel` above 4 | worktree | **Measured, the package alone:** 81.1 s at p4, 49.4 s at p8, 35.6 s at p16. **Measured, five packages together:** 106.6 s at p4 and 109.6 s at p8, so no gate gain. The host is CPU-saturated. | S | 8 that it gives no gate gain |
| 7 | Build the landing fixture once and copy it | worktree | **Estimate:** 3-8 s of CPU. Worktrees hold absolute paths, so copies break. I do not recommend it. | L | 3 |

The prior review's gate-level `-p` finding stays. Item 6 shows that `-parallel` moves time between packages, not out of the wall, while the host is saturated.

## Table 1: tests that add no value (`internal/worktree`)

The package pins its exact top-level test count at 767 in `parallel_census_test.go:1061-1078` (`TestPackageTestCountPin`, row WF12). Its comment says "no test is removed or merged for wall-clock". It also pins the serial set at 45 (`:1041`). Rows that remove or merge a whole top-level test therefore need a reviewer decision on the pin. Rows that touch only a subtest or an assertion do not.

Times are elapsed under `-parallel=4` and load. "Covered by" is a block-coverage fact from isolated runs: the test's covered blocks are a subset of the other test's. Coverage is evidence only. I confirmed each row by reading the test.

| # | Test (file:line) | Why it adds no value | Covered by | Action | Conf. |
|---|---|---|---|---|---|
| 1 | `identity_component_test.go:217` `TestLandingRefusalRegistryHasAProducingFixture`. Subtests `:235-240` took 4.32 s and 3,292 git starts. | The completeness code at `:218-234` repeats `routetest.Fixtures` (`routetest.go:46-70`). The 13 subtests re-run `produceLandingFace`, which the follow walk already calls (`refusal_route_follow_test.go:54`). Its coverage is a subset of the follow test (0 unique blocks). | `refusal_route_follow_test.go:33-76` | Delete the subtest loop. Delete the whole test only after a pin decision. Re-point coverage rows LRS1, LRS2 and RR14 in its doc comment. | 8 |
| 2 | `land_census_test.go:55` `TestLandCommandStatesZeroForAnAssignmentWithNoRecords` (0.90 s) | Same arrangement as `:85`. Its assertions (exit 0, `,census=0}`) are a subset of `:85`'s. Coverage is identical. | `land_census_test.go:85` | Delete (pin change) | 9 |
| 3 | `land_effects_test.go:230` `TestLandEffectsRowPrecedesTheLandedRecord` (0.88 s) | Same fixture, joins and landing call as `:109`. I measured 588 git starts for both, and coverage is identical. Only the tail-order assertion differs. | `land_effects_test.go:109` | Add the tail-order check to `:109`, then delete `:230` (pin change) | 7 |
| 4 | `land_resume_test.go:209` the `x` row of `TestResumeLandCommandAcceptsSpecSlugAndPath` (about 0.9 s of the 1.83 s) | Same arrangement and coverage as `:226`. It asserts exit 0, released, empty stderr and the gate tally, which `:226` also asserts. Only the `./specs/x/spec.md` row needs the `chdir` at `:213`. | `land_resume_test.go:226` | Drop the `x` row (no pin change) | 7 |
| 5 | `refusal_route_test.go:172` `TestReviewerLandFacesOpenWithTheMarker` (2.19 s, 1,777 git starts) | It produces 7 faces again. Its coverage is a subset of the follow test. Its list is the spec's list, not the registry's (doc `:168-171`). That independence is necessary: the walk takes its authority from the same registry as the route. I found no recorded red. | `refusal_route_follow_test.go:33` | Do not delete. Move the independent list into the follow loop so each face is produced once. A recorded red is needed for the list. | 6 |
| 6 | `land_freshness_test.go:74` `os.Getpid() != ownerPid` | `runVerb` calls the verb in-process (`verb_runner_test.go:122-150`). A re-exec would replace the test process, so this check can never fail. Also, `:59` has the same coverage as `:23` (Jaccard 1.000), and its other assertions repeat `:23`'s. | `land_freshness_test.go:23` | Delete the `:74` check. Move the release-completed check (`:83`) into `:23`, then delete `:59` (1.68 s, pin change). | 8 (dead check), 5 (`:59`) |
| 7 | `identity_component_test.go:252` `deferredMergedDetails` and its uses at `:271` and `:278-282` | The slice is empty, so the `slices.Contains` branch and the final loop never run. This is a dead allowance. | none needed | Delete the dead branches (0 s) | 7 |
| 8 | `land_census_test.go:37`, `:68` and `land_census_output_test.go:13` (0.91, 0.88 and 0.88 s) | Three landings that differ only in the records seeded before the landing. Coverage differs by 13 blocks or fewer. | one combined landing | Seed both record kinds, assert all outputs in one landing. Saves about 1.8 s. Pin change. | 6 |
| 9 | `land_resume_test.go:244` `TestResumeLandCommandPreauthenticatesCompletedRequestAndPath` (1.80 s, 1,222 git starts) | Two rows build the same full landing. They differ only in the forged resume arguments. The forged resume mutates nothing, so one landing can serve both. | itself | Build one landing, run both forged resumes. Saves about 0.9 s. | 5 |
| 10 | `land_resume_refusal_test.go:35` row `PL30/clean` (1.05 s) | A control row. Its outcome is the same as the real-binary active resume at `land_resume_test.go:128`, which also follows a destination move. | `land_resume_test.go:88-130` | Drop `PL30/clean`. Keep `PL25/clean` as the control. | 5 |
| 11 | `refusal_route_test.go:141` `TestLandingRedRouteNamesTheRepair` (0.77 s) | Coverage is a subset of the follow test. The follow test's `names` field already requires `bench commit --in <label>` (`land_refusal_fixture_test.go:214`). Two unique checks remain: the detail prefix and the `; then <rerun>` suffix. | `refusal_route_follow_test.go:33` | Fold those two checks into the follow walk | 4 |
| 12 | 62 pasted `"-c","user.name=bench","-c","user.email=bench@local"` argument runs, 58 `"--spec","x"`, 30 `worktree=released,census=0}` | One fact repeated many times. `RepoOnBranch` already sets the identity in repository config (`gittest.go:146-147`), and linked worktrees share it. A `commitInWorktree` helper exists and is used 174 times. | `commitInWorktree`, `landArgs` | Add one commit helper and one released-record helper. No time effect. | 6 |

**System suite rows (not read in full):**
- S1. `owner_selection_test.go:44` `TestFocusedRunStillCleansOwnerRoot`. The `pass` and `fail` rows run the same child as the first form of `:19`. Each child starts a full test binary, about 0.3 s. Merge `assertOwnerRootRemoved` into `:19` to save about 0.6 s (conf 6).
- S2. `adoption_test.go:18`, carried from the prior review. It took 1.65 s. Conf 5.

**No findings (worktree):**
- Silent skips: only the 2 socket rows skip, and the charge says to ignore them. The symlink, fifo and privilege capability seams all execute here.
- Sleeps and polling: only 4 bounded `time.After` waits and 1 poll helper (`subshell_test.go:134`). No test with a sleep exceeds 0.9 s.
- Static census tests (parallel, serial, single-read, literal): they run in 0.00 s. They have no time value.
- Per-test binary builds: none. `test_run_test.go:43-52` builds the run binary once through `sync.Once`. A warm build took 1.3 s.

## Table 2: efficiency changes

| Pkg | Change | Estimated saving | Effort | Risk to coverage | Evidence |
|---|---|---|---|---|---|
| A. System | Shard the suite over 4 processes, or add `t.Parallel`. The suite has 0 `t.Parallel`, and it shares one `owner` with 3 repos. | 50-60 s wall. | M | The unfiltered ledger check in `TestMain` (`owner_test.go:63-69`) is skipped when `-run` is set. A real design must merge the ledger across shards. | **Measured:** shards 23.0 to 25.9 s, all exit 0. Serial runs: 75.4, 87.2 and 76.9 s. CPU is 1.1 cores. |
| B. worktree (production) | Cache pure git queries per command run. | 18-28 s of the package, estimate. | M | Stale-cache risk. The hostile-path tests need validation after the cache. A reviewer decision. | **Measured:** 182,555 and 182,559 git starts in two runs. 89% come from the test process (162,015 of 182,209 attributed lines), and 11% from the `bench` child. Time and git starts correlate at r = 0.918 over 483 tests (1.88 ms per start). 53% of starts are pure queries, and 41.9% (76,292) repeat an earlier identical query in the same test. Largest group: `rev-parse` location, 49,952 starts with 92% repeats. One landing costs 588 starts (fixture alone 142). Single call cost: 1.5 ms for `rev-parse`, 1.9 ms for `cat-file` in a shell loop. |
| C. worktree | Move serial tests to parallel. Use child env or an env seam in `runVerb`. | 6-7 s at p4, estimate. | M | None expected. Update `worktreeSerialCeiling` at `parallel_census_test.go:1041`. | **Measured:** 45 serial tests take 9.42 s; the top 5 take 6.1 s. Blockers: 29 `bindEnv` and 20 `chdir` sites. 846 `t.Parallel` calls exist. |
| D. worktree | `t.Parallel` on subtests of serial-chain tests. | 3-5 s at p16, estimate. | S per test | None. Each subtest builds its own fixture. | **Measured:** 75 tests run subtests serially (57.3 s of elapsed). `TestResumeLandCommandPublicRefusesDestructiveDestinationState` takes 13.4 s (9,408 git starts, 14 subtests). |
| E. worktree | Fold the duplicate landings (Table 1). | About 11-15 s of CPU, 3-4 s wall at p4, estimate. | S-M | Pin changes for the whole-test rows. | **Measured:** coverage clusters hold 29 clusters of 112 tests (84.8 s). The largest is 16 near-identical landing tests (14.8 s). The greedy minimal block cover keeps 263 of 767 tests (147.6 s) and leaves 504 (133.5 s) with no unique in-process block. |
| F. host config | Raise `-parallel`. | Alone: 31.7 s at p8 and 45.5 s at p16. In a gate-like run: none. | S | None | **Measured:** see Measurements 8-9. A reviewer decision (`GOFLAGS`). |
| G. worktree | Landing fixture template. | 3-8 s CPU, estimate. | L | Medium. Worktree absolute paths. | **Measured:** fixture setup is about 140 starts, about 24% of one landing's 588. |
| H. System | `t.Parallel` inside the suite. | Not measured beyond A. | M-L | The owner repos are mutated by tests, for example `owner_selection_test.go:227-243` runs `link copy` on `owner.repos[0]`. I did not audit the other 90 tests. | Read `owner_test.go:1-200` and `owner_selection_test.go`. |

**Reading traps:**
- `TestParallelJourneysShareTheHarnessSafely` shows 54-60 s. That is the wait for its parallel subtests' token, not cost.
- Parents of parallel subtests show inflated times for the same reason.

## Measurements: exact commands and results

All ran from the worktree via `bench worktree exec test-efficiency-2 -- ...`. `D` is the scratch directory above. `RUN` is `BENCH_RUN_BINARY=$D/run/bench BENCH_KIT=$PWD`. I built that sealed binary with `bash scripts/go-build.sh --manifest-dir $D/run $PWD $D/run/bench` (1.3 s).

1. `env TMPDIR=/dev/shm bash -c 'time bench test --check system'` gave 167 tests, pass, 75,426 ms (wall 77.4 s, load 4.9-6.3, 22:27). Repeat at 22:55 gave 87,155 ms (wall 88.8 s, load 7.1-8.4).
2. `TMPDIR=/dev/shm go test -count=1 -json ./internal/worktree/` gave 84.07 s, 1,385 tests, 1 fail (`TestCommitmentPlanningBootstrap`, which needs `BENCH_KIT`), 2 socket skips. Wall 88.2 s, user 3m44, sys 1m42.
3. `env TMPDIR=/dev/shm bench test --package ./internal/worktree` gave pass, 77,987 ms (wall 83.2 s, load 9.6-11.3, 22:30). Repeat at 22:54 gave 61,187 ms (wall 63.2 s, load 8.6-7.7).
4. System suite as JSON: `TMPDIR=/dev/shm $RUN BENCH_SYSTEM_ROOT=$PWD go test -tags=system -count=1 -json ./internal/systemtest` gave pass, 76.93 s (wall 77.8 s, user 52.5 s, sys 31.3 s). The top-level times sum to 76.5 s, so it is serial. Slowest tests, in seconds:

| Test | Time |
|---|---|
| `TestLifecyclePublicationCompatibility` (48 subtests) | 6.39 |
| `TestProspectiveArtifactRecoveryAfterKilledLanding` | 5.10 |
| `TestAdoptionSmokeJourney` | 5.04 |
| `TestOneAuthorizationRecoversTheDeadBundleOfADeadAndLivePair` | 4.36 |
| `TestUnfilteredRunVerifiesLedger` (40 children) | 3.66 |
| `TestExecReturnsAtChildExit` | 3.14 |

5. Worktree with the sealed binary: `TMPDIR=/dev/shm $RUN go test -count=1 -json ./internal/worktree/` gave pass, 80.13 s (wall 81.1 s, load 8.4-12.4, user 3m23, sys 1m34). The serial head was 9.42 s; the parallel phase summed to 271.7 s.
6. The same command with `-parallel=8` gave wall 49.4 s, exit 0, load 11-13. With `-parallel=16` it gave wall 35.6 s, exit 0, load 12-14. CPU was unchanged in both.
7. Counting shims on `PATH` (`$D/shim/git`, `$D/shim2/git`; each logs the call, then runs `/usr/bin/git`). Run with `go test -count=1 ./internal/worktree/`. Both gave pass with 182,555 and 182,559 git starts. Shim runs took 120-137 s and I used them for counts only.
8. Gate-like run at `-parallel=4`: `env TMPDIR=/dev/shm GOFLAGS="-p=4 -parallel=4" bench test --package ./internal/worktree --package ./internal/testreport --package ./internal/conformance --package ./internal/preflight --package ./internal/gate --package ./internal/evidencecmd`.
   - I gave `evidencecmd` the wrong path, so it failed with 0 tests. The five other packages ran.
   - Wall 106.6 s (load 9-13). Package times: conformance 101.5 s, testreport 99.7 s, worktree 94.8 s, preflight 63.7 s, gate 30.4 s.
9. The same run without `evidencecmd` and with `-parallel=8` gave wall 109.6 s (load 11-16). Package times: conformance 105.6 s, testreport 107.7 s, worktree 71.9 s, preflight 74.5 s, gate 27.8 s. All pass. Worktree gained 23 s and the others lost time.
10. Coverage: I built the test binary with `go test -c -cover -covermode=set -coverpkg=<worktree,landing,intent,refusalroute,freshness,commitment>`. I ran each of the 767 top-level tests alone with `xargs -P 6`. All 767 passed in isolation. The run took 54.0 s wall and 295 s of summed isolated time. This covered 3,654 blocks. The system-suite shard result is in item 11.
11. System shards: 4 balanced groups (19.1 s each by measured time) via `xargs -P 4` and `go test -tags=system -run '^(...)$'`. Wall 25.9 s. Shards finished at 23.0, 24.6, 25.5 and 25.9 s. All exit 0.
12. Git per-call cost: a 500-iteration bash loop took 0.766 s for `rev-parse --path-format=absolute --git-common-dir` and 0.936 s for `cat-file blob`. An empty `/bin/true` loop took 0.386 s.

## What I read and did not read

**Read:**
- Prior report in full.
- `delegation-discipline.md`.
- Worktree harness: `journey_race_test.go`, `test_run_test.go`, `journey_test.go`, `main_test.go`, `verb_runner_test.go:90-170`, `land_fixtures_test.go`, `completion_fixture_test.go`, `land_refusal_fixture_test.go`.
- Worktree landing and refusal tests: `land_resume_refusal_test.go`, `refusal_route_follow_test.go`, `refusal_route_test.go`, `identity_component_test.go`, `land_effects_test.go`, `land_effects_cleanup_test.go`, `land_journey_test.go`, `land_resume_test.go`, `land_freshness_test.go`, `land_census_test.go`, `land_census_output_test.go`, `land_surface_test.go`, `land_broker_notice_test.go`, `commitment_landing_test.go:1-242`.
- Pins and neighbours: `parallel_census_test.go:1000-1097` and a names search of the rest.
- Production seams: `internal/refusalroute/routetest/routetest.go`, `internal/gittest/gittest.go` (parts), `internal/runbinary/runbinary.go`, `internal/git/worktree_admin.go:120-230`, `internal/gate/phases.go:55-125`.
- System suite: `internal/systemtest/owner_test.go:1-200`, `owner_selection_test.go`.

**Not read or not measured:**
- About 100 of the 123 worktree test files: `commitment_light_landing_test.go`, `clean_set_test.go`, `land_reauthorization_test.go`, `delegated_integration_test.go`, `ownership_test.go`, `clean_landed_apply_test.go`, `land_flags_test.go`, and the rest.
- The system tests other than the two files above. Table 1 for the system suite is a sample.
- No probe or recorded red for any "subsumed" claim. Coverage subset is not assertion subset. I confirmed rows 1-11 by reading only.
- The gate-level effect of any item. Runs 8 and 9 use five packages, not the real gate. I did not run `bench gate`.
- Whether the landing broker keeps the caller's `TMPDIR`. The effect of item 2 on production behaviour.
- Other packages' tests beyond the five in runs 8 and 9.

## Residue and caveats

- **Scratch:** `.../scratchpad/eff2/` holds JSON runs, 3 git-call logs, per-test coverage, the sealed binary and the coverage test binary.
- **`/dev/shm`:** two 4 KB directories, `TestOtelCrashKeepsStartedPhaseLine3664604367` (22:34) and `...3031551307` (22:58), come from my direct system runs. The system test leaks its temp directory when run directly.
- **Other `/dev/shm` directories:** the six-letter ones and `bench-cmd-home-*` belong to other sessions. A live system run was still writing into `LYOZV6` (761 MB) at 23:04. I removed nothing.
- **Reviewer decisions:**
  - Relax `worktreeTestCount` (767) and `worktreeSerialCeiling` (45) for any merge or serial-to-parallel change.
  - Approve any production cache in `internal/git`.
  - Decide whether system-suite sharding replaces the single-process ledger check.
  - Leave `GOFLAGS` at `-parallel=4` while the host stays saturated.

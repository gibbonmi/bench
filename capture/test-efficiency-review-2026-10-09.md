## Interim status, then final report

Nothing is still running. The call that showed no result was my last search, which the permission classifier refused. Two earlier background jobs finished. I ran no new measurements after your stop order.

**Residue from my runs:** 22 entries (about 18 MB) sit in /dev/shm: `bench-cmd-home-*` x20, `bench-npm-cache`, `node-compile-cache`. All were created at 20:28-20:44 by my test runs. My `rm -rf` on those exact names was refused, so I left them. They are harmless and a WSL restart clears them. cmd/bench tests leak 10 `bench-cmd-home-*` dirs per run into TMPDIR. I could not locate the source because the search was refused.

---

## Lead: ranked changes (savings are against the gate logs, not a measured gate run)

Baseline: the 23:48 gate took 739 s. The test phase was 478 s, race 4 s, and system 250 s. The phases run strictly one after another (`internal/gate/runner.go:213-266` launches one phase, then waits). Your newer numbers (adopt 351 s and repairtest 127 s at 00:32) put the test phase near 352 s. I did not read the 00:32 log.

1. **Put the kit test run's TMPDIR on tmpfs (S effort, about 300-375 s per gate).**
   - Almost all of the cost in adopt, repairtest and the long system journeys is fsync on this host's ext4 (/dev/sdd, WSL2). /dev/shm is tmpfs (12 GB, mounted `rw,nosuid,nodev,noatime`, so it is not noexec).
   - One owner picks the base: `internal/env/kit_run.go:41-44` reads TMPDIR, else `os.TempDir()`. The gate merges the run's TMPDIR into every phase (`internal/gate/phases.go:293-307`). `bench test` uses the same owner (`internal/testreport/command.go:225`). The prospective artifact reads `os.TempDir()` (`prospectiveartifact.go:152`).
   - Estimate: test phase about 150 s (about -200 s against 352 s; -325 s against 478 s), plus system about -100 to -175 s.
   - The test-phase figure comes from a package scheduling simulation (below). The system figure comes from 3 sampled slices.
   - Evidence is in the measurements section.
2. **Raise `-p` after item 1 (S effort, about 50 s, reviewer decision).** The host config `/home/mgibs/.config/go/env` holds `GOFLAGS=-p=4 -parallel=4` (`projects/benchkit.md:604`). With tmpfs the package time sums to about 470 s, so `-p=4` becomes the bound. My simulation reproduced the real 477 s test phase. Under tmpfs estimates it gives about 152 s at p=4, 121 s at p=6 and 102 s at p=8.

   The risk is load on a host shared by 2 delegates. Coverage is unaffected.
3. **Fix 4 socket rows that skip on every local run (S effort, correctness, 0 s).** See Table 1 row 1. This is a confirmed defect, so it comes first when you order by "defects first".
4. **Cut the remaining per-test hot spots (S-M effort, about 3-10 s of wall total).** After item 1 the tail is flat, and no single test is more than about 5 s of wall.
   - `testreport` `TestRootConformanceCheckPrintsFailures`: 17.7 s on tmpfs.
   - `conformance` `TestEveryRetainedFixtureBitesThroughRegisteredOwner`: 15.6 s. Its fixture subtests can run in parallel (`fixture_bite_test.go:723-735`; the `completed` map write needs a mutex).
   - `TestSetupGateRejectsIgnoredDeclaredInputs`: 4.8 s.
5. **Fallback only if tmpfs is rejected:** inject root and kit into the adopt entry points so tests can use `t.Parallel`. See Table 2 row G.

After item 1 the adopt-only levers are worth very little. The adopt_test.go migration (planned ticket 2) is sized below and is not worth doing.

---

## Table 1: tests that add no value

Line numbers are current main for files the golden landing did not touch. For `lifecycle_*` files they use pre-golden numbering (I read those before 96c2b342). The time column is single-package ext4 unless I say otherwise.

| # | Test (file:line) | Why it adds no value | What covers it instead | Action | Conf. |
|---|---|---|---|---|---|
| 1 | **DEFECT.** `internal/worktree/clean_landed_hostile_test.go:79` (socket row about :96-101); `:137` (socket row about :148-152); `internal/conformance/prose_budget_test.go:351` (guard :381); `skill_description_budgets_test.go:443` (guard :461) | The 4 `/socket` rows SKIP under the default TMPDIR (measured, see item 9). `net.Listen("unix", path)` fails with "bind: invalid argument" because the path is over the 107-byte limit. The paths are 116 bytes or more even under /tmp, so the gate's private TMPDIR is longer still. The gate runs without `BENCH_REQUIRE_CAPABILITIES`, so the skip is silent. | The fifo and symlink rows of the same tests cover the non-regular refusal. The socket-specific branch is unproven. | Create the socket at a short path, then `os.Rename` it to the long path. This is my suggestion and I did not try it. | 9 (local skip measured); 6 (gate behaviour inferred) |
| 2 | `internal/adopt/compatibility_test.go:101` `TestCompatibilityLegacyDoctor` | The expectation comes from the production function under test. `Doctor` forwards every non-`--compat` argument to `legacyDoctor` (`compatibility.go:49-54`), and the test compares `Doctor` with `legacyDoctor` on those same arguments. It can only catch a dispatch change. | `adopt_test.go:362` (help and unknown flags), `:728` and `:763` (doctor `--fix`) | Delete | 7 |
| 3 | `internal/adopt/repairtest/repair_test.go:76` `TestCompatibilityUndoCreated` (1.81 s) | It is a strict subset of `TestCompatibilityUndo` (`:159`): same remove, `--fix`, stat, undo, absence checks and the same `stderr==""` check. | `repair_test.go:159` | Delete | 8 |
| 4 | `repair_test.go:18` `TestCompatibilityManagedRepair` (1.78 s) | Same arrangement as row 3. Its one unique assertion is that the restored bytes equal the original. | `repair_test.go:159` | Move the byte comparison into `TestCompatibilityUndo`, then delete | 6 |
| 5 | `repair_test.go:217`, `:243`, `:271` (Policy, Interruption and Runtime boundary tests; 3.49, 4.31 and 3.77 s) | Same arrangement each time: a linked session plus `doctor --fix`. Only the asserted decision row differs. | Each other | Merge into one test with one linked session, keeping separate CODEX_HOME files. Saves about 8 s. | 6 |
| 6 | `internal/adopt/setup_test.go:159` `TestSetupGateDeclaresConsumerOnlyValidityScope` | It runs a full Setup to check a substring. The literal is copied from `setup.go:351`. By reading, dropping that line turns the `:111` rows red, because only that hygiene check makes those rows red. I ran no probe. | `setup_test.go:111` | Delete | 7 |
| 7 | `setup_test.go:111` `TestSetupGateRejectsIgnoredDeclaredInputs` (3 rows; 4.8 s on tmpfs; about 1.6 s per row is `bash gate` spawning `bench test --check`) | All 3 rows repeat the same arrangement (git init, Setup, retire sentinel, write files). Only the manifest differs. | n/a | Run Setup once, then loop over the manifests. Saves 2 Setups, about 4 s on ext4 (estimate). | 7 |
| 8 | `internal/systemtest/adoption_test.go:18` `TestDetectedProjectGateRejectsIgnoredDeclaredInput` | It repeats the "space" row of row 7 at the wrapper seam. It adds one `git init`, one `setup --yes` and one gate. | `adoption_test.go:70` (smoke journey, same repo and wrapper) | Fold it into the smoke journey as a last leg. Saves about 5-10 s of system time (estimate). | 5 |
| 9 | `setup_test.go:53` `TestSetupSeedsGateInputs` | It runs a Setup on a go.mod repo and compares exact bytes. `:86` runs two Setups on the same repo but only compares first with second. `:98` already does the exact-bytes check on a bare repo. | `:98` | Move the exact-bytes assertion into `:86`, then delete `:53`. Keep the independent `wantSeededGateInputs` literal. I found no recorded red for it. | 5 |
| 10 | `adopt/decision_test.go:8` `TestPlanLifecycleFromImmutableInventories` | Its main assertion is a subset of the "legacy" subtest at `lifecycle_decision_test.go:82-100`. The invalid-path rows duplicate `:190-195`. | `lifecycle_decision_test.go:81` | Move the two unique rows (duplicate, missing fingerprint) into the refusal table, then delete | 6 |
| 11 | `adopt/link_hook_test.go:345` `TestPrePushHookTopicPushIsSilent` | The legacy-pin test at `:362` is the same topic push, also silent, and asserts more. | `link_hook_test.go:362`, `:319` | Delete (about 0.05 s) | 7 |
| 12 | `lifecycle_preview_test.go` about :75 (e2e Upgrade `--check` rows) | 8 of the 11 rows exercise the same `upgradePlanCounts` branches as the pure table at about :142. Each also repeats the "manifest not rewritten" assertion. | `testPreviewPlanKinds`, about :142 | Keep "converged", one drift row and the two hook-refresh rows. About 2-3 s (estimate). | 5 |
| 13 | `lifecycle_consumption_test.go`, `lifecycle_special_test.go` (e2e partition tables for Claude, Dropped, Link) | They re-drive the partitions of the pure `PlanLifecycle` table (`lifecycle_decision_test.go:37-79`, `:81`) through real entries. | The pure table | Keep one e2e row per consumption mechanism (install, keep and report, keep silent, refuse). | 4 |
| 14 | `internal/testreport/root_conformance_test.go:81` `TestRootConformanceCheckPrintsFailures` | It runs the whole root-conformance suite on a kit copy (17.7 s on tmpfs, 45% of the package) to see one diagnostic printed. The gate already runs the same suite on the live tree in `conformance` (`TestRootConformance`, 17.9 s). | `internal/conformance/gate_entry_test.go` | Scope the run to one check, or stub `go` with canned output. It may be a deliberate bite witness, and I do not know whether a red was recorded. | 4 |
| 15 | `adopt_test.go:412` and `:437` (also `:469`), `setup_prompt_test.go:169`; helpers `adopt_test.go:307` vs `link_hook_test.go:206` | The `initRepo` preamble (`adopt_test.go:510`) is pasted, and `runAdoptGit` equals `runHookGit`. Repo-wide there are about 25 `runGit`-style helper copies. | `gittest.Output`, `testrepo` | Collapse into one helper. No time effect. | 8 |

**No findings:**
- Long sleeps: none found. All `time.Sleep` uses are 10 ms or similar.
- A scan for test functions with no assertions turned up only helper-based assertions, so I report nothing there.

**Your question about the twelve copy-based tests** (DroppedRows, ClaudeOwnership, UnlinkPartitions absent cases, UnlinkRefusals):
- Yes, they could start from `consumerRepo` plus a hand-written manifest, as `TestLifecycleOwnerConsumption` does.
- The saving is small. `linkedCopy` does a `CopyFS` of about 127 files plus a chmod walk, with no fsync. I put it at about 0.1-0.3 s per copy, about 1-4 s in total (estimate, unmeasured). I rank it low (confidence 4).
- The `git init` followed by `RemoveAll(.git)` in `linkedCopy` wastes one git subprocess per copy, about 15-40 ms each. That is 1-2 s over the suite (estimate).
- I did not read `consumerRepoAt`. `golden_test.go` already calls it, so the copy path could reuse it without a git init.

**Sizing the adopt_test.go migration (planned ticket 2):** not worth doing.
- In adopt_test.go only `TestLinkRelinkStaysGreenInAConsumerRepo` does real fsync work. It took 1.74 s, and its second Link is converged and cheap.
- The Init tests measured 0.00-0.01 s each, and the doctor tests do not Link.
- The remaining fsync cost is in `setup_test.go` (9 Setup calls) and `setup_prompt_test.go` (2 interactive Setups). It is also in `TestLifecycleSeedOwnership` (5 Setups), the compat and Upgrade tests, and the rows whose subject is Link.

---

## Table 2: efficiency changes

| Pkg | Change | Estimated saving | Effort | Risk to coverage | Evidence |
|---|---|---|---|---|---|
| A. gate and `bench test` (`internal/env/kit_run.go:41-44`) | Prefer a tmpfs base such as /dev/shm for the private run directory, with fallback to `os.TempDir()`. A per-caller `TMPDIR=/dev/shm` may work with no code change, but I did not verify the landing passes TMPDIR through. A package-scoped alternative is a TMPDIR set in `TestMain` of adopt, repairtest and systemtest. | Test phase about -200 s (against the 00:32 baseline), system -100 to -175 s. Adopt and repairtest together were 478 s of CPU-slot time in the 00:32 gate. | S | None seen. Durability cannot be observed by any test, and fault tests inject at `syncDirectory`. Open points: 2 worktree test failures in my raw-TMPDIR runs, one of them unexplained (item 5); peak tmpfs use was not measured (the tracked tree is 16.1 MB, so a KitCopy is small); socket paths are 4 bytes longer on rows already skipping. | Items 1-8 below |
| B. host Go config | Raise `-p` (4 to 6-8). | About -30 to -50 s after A. | S | None. Load on the shared host. Reviewer decision per `projects/benchkit.md:604`. | `sim.py` |
| C. gate runner | Overlap `test` and `system`. | Up to about 100 s after A. | L | None. ADR 0024 requires a resource census first. | `runner.go:213-266`; `docs/adr/0024-*` |
| D. `internal/conformance` | Run fixture-bite subtests in parallel (`fixture_bite_test.go:723-735`). Also parallelise `TestRootConformance` and `TestWorkflowCadenceAnchorsRejectDeletionAndSwap` if they use no Setenv or Chdir. | About -11 s package time (15.6 s to about 5 s with `-parallel=4`), about 3 s wall. | S | Whether `owner.run` is goroutine-safe is unverified. | json per-test run |
| E. `internal/adopt` | Setup golden: share one bare-repo `setup --yes` among the 3 `setup_test.go:111` rows and `:159` (this is Table 1 rows 6 and 7). | About -6 to -10 s ext4. | S | None. | `setup_test.go` |
| F. `internal/adopt` production | Link stages each asset with an fsync (`link_stage.go:173` via `transaction.go:20-30`), then stages it again with an fsync (`transaction/publication.go:56`), then syncs the directory per publish (`:108`). There are about 127 managed assets, so up to about 3 syncs each. Drop the first-stage sync and batch the directory syncs per parent. | Up to about one third of Link fsyncs, on ext4 only. | M | Durability and recovery ordering are a reviewer decision. The journal record holds the image data, but I did not verify recovery. | code read |
| G. `internal/adopt`, `repairtest` (fallback only) | Add a thin `linkAt(root, kit, ...)` so tests need no Setenv or Chdir. `Link`, `Init`, `Upgrade`, `Setup`, `Unlink`, the doctor rows and the compat entries all call `git.Root()` (cwd) and `gate.KitDir()` (BENCH_KIT) at entry: `link.go:37,46`; `init.go:23,28`; `upgrade.go:54,64`; `setup.go:32,37`; `unlink.go:37`; `doctor.go:234,262,363`; `compatibility.go:119,171,225,258`; `broker.go:53`. `transactionalLink(root, kit, ...)` (`link.go:64`) already takes both. Other blockers: package var `syncDirectory` (swapped at `adopt_test.go:156-158`), `BENCH_LINK_FAULT` (`transaction.go:96`), HOME, PATH, CODEX_HOME and BENCH_HOME. Counts: 32 `t.Setenv` and 20 `Chdir` in adopt; `repairtest/session_test.go:23-27` has 4 `t.Setenv` plus `Chdir`. | Under `-parallel=4`, overlapping fsync waits could cut adopt to roughly a third or quarter (estimate, unmeasured). Worth under 5 s wall once A is done. | L | Keep a few tests on the cwd adapter. | code read |
| H. golden review nit (`linkedCopy` via `consumerRepo`) | Skip `git init` then `RemoveAll(.git)`. | 1-2 s over the suite (estimate). | S | None. | `golden_test.go` |

**Package times, last full gate (23:48):**

| Package | Time |
|---|---|
| cmd/bench | 22.9 s |
| adopt | 477 s |
| repairtest | 155 s |
| commit | 22 s |
| conformance | 59.6 s |
| gate | 28.3 s |
| preflight | 44.5 s |
| evidencecmd | 27.4 s |
| probe | 23.3 s |
| runbinary | 21 s |
| shift | 25 s |
| testreport | 65.5 s |
| worktree | 123.6 s |
| systemtest | 249.8 s |

 The other 85 packages were 15 s or less each.

---

## Measurements I ran (exact commands and results)

All ran from `/home/mgibs/workspace/bench`.

1. `go test -count=1 -json -run 'TestUnlink|TestLinkRelink|TestInit' ./internal/adopt/` (default TMPDIR=/tmp).
   - Link-based tests took 1.22-2.16 s each. Init tests took 0-0.01 s each. The slice took 11.3 s.
   - `TestUnlinkRemovesOnlyCleanManifestAssets` 2.16 s; `TestLinkRelinkStaysGreenInAConsumerRepo` 1.74 s.
2. Same with `TMPDIR=/dev/shm` and `-run 'TestUnlink|TestLinkRelink'`: 0.11-0.22 s per test, 0.99 s for the slice.
3. `go test -count=1 -json ./internal/adopt/repairtest/`: 62.7 s on default TMPDIR. With `TMPDIR=/dev/shm` it took 3.3 s, and every test took 0.34 s or less.
   - On ext4: `TestCompatibilityRepairIdempotence` 10.3 s, `UndoConflict` rows 9.1-9.6 s, `ForeignConflict` 3.7 s, the 3 boundary tests 3.5-4.3 s, the rest 1.3-2.4 s.
4. `TMPDIR=/dev/shm go test -count=1 -p=1` over 11 packages (testreport, preflight, evidencecmd, shift, probe, commit, runbinary, cmd/bench, gate, worktree, conformance):

   | Package | tmpfs alone | Gate 23:48 | Note |
   |---|---|---|---|
   | testreport | 46.1 s | 65.5 s | |
   | preflight | 34.7 s | 44.5 s | |
   | evidencecmd | 24.4 s | 27.4 s | |
   | shift | 7.4 s | 25 s | ext4 alone: 10.2 s |
   | probe | 24.3 s | 23.3 s | |
   | commit | 12.1 s | 22 s | ext4 alone: 13.7 s |
   | runbinary | 16.9 s | 21 s | |
   | cmd/bench | 23.6 s | 22.9 s | |
   | gate | 19.1 s | 28.3 s | ext4 alone: 26.6 s |
   | worktree | 69.1 s, FAIL | 123.6 s | See item 5 |
   | conformance | 62.6 s, FAIL | 59.6 s | See item 6 |

   - The ext4 comparison for shift, commit and gate was `go test -count=1 -p=1 ./internal/shift/ ./internal/commit/ ./internal/gate/`.
   - The gate numbers are loaded, the others quiet, so the comparison is only indicative.
5. worktree:
   - The first tmpfs run failed `TestLandingRetiresAnAbsentLandedSibling` (prospective authorization refused, "evidence absent") and `TestCommitmentPlanningBootstrap`.
   - On default TMPDIR (99.5 s) and again on tmpfs (60.1 s), only `TestCommitmentPlanningBootstrap` failed. It needs the gate's environment: "lane Bench executable unavailable ... scripts/go-build.sh: No such file".
   - `TestLandingRetiresAnAbsentLandedSibling` passed alone on both TMPDIR values and in the rerun. It is an unexplained failure in 1 of 2 tmpfs runs, with no evidence of a tmpfs cause. Treat it as unknown. A real gate run with the redirect is the only proof.
6. conformance on tmpfs failed 3 tests: `TestRootConformance`, `TestGuidanceProseBudgetsHoldOnTheLiveTree` and `TestSkillDescriptionBudgetsHoldOnTheLiveTree`.
   - The messages name `.agents/skills/typesafe-ai` and `.logs/sr-line-comparison`. Those are untracked files in the primary checkout, which the gate's clean checkout would not have.
   - I did not rerun these on /tmp, so this attribution rests on the messages.
7. adopt slices on tmpfs, run before the golden landed:
   - `-run 'TestLifecycle'` gave 16 tests in 7.4 s.
   - `-run 'TestSetup|TestUnlink|TestLink|TestInit|TestDoctor|TestCompat|TestUpgrade|TestTransactionalLink'` gave 48 tests in 7.1 s. `TestSetupGateRejectsIgnoredDeclaredInputs` took 4.8 s of that.
   - I did not run the whole adopt package, and I ran no ext4 equivalent of these slices. My estimate for adopt on tmpfs is about 25 s (unmeasured remainder).
8. System lane through `bench test --check system --run ...`:
   - `'TestDetectedProjectGateRejectsIgnoredDeclaredInput|TestAdoptionSmokeJourney'`: 29.8 s default, 16.2 s with `TMPDIR=/dev/shm`.
   - `'TestLifecycleEntryJourney|TestLifecyclePublicationCompatibility|TestCompatibilityConcurrentWriters'` (15 tests including subtests): 176.5 s default, 8.0 s on tmpfs.
   - I did not run the whole 92-test system lane.
9. Socket skips: `go test -count=1 -v -run 'TestCleanLandedSpecialPathsRetainedWithoutOpening|TestLandedConsumersRejectSpecialGitMetadataBeforePlanning' ./internal/worktree/` showed both `/socket` rows SKIP with "bind: invalid argument". A matching conformance run (`-run 'NonRegularSubjects|NonRegularEntries|GuidanceSweep'`) showed the 2 conformance socket rows SKIP. `bounds` and `landing` socket rows ran without skipping.
10. Per-test tmpfs `-json` runs for testreport, preflight, evidencecmd, probe, conformance and cmd/bench:
    - testreport 39.1 s: `TestRootConformanceCheckPrintsFailures` 17.7 s, `TestCommandDrainsBuilderGroupOnCancelSignal` 6.0 s.
    - conformance 55.0 s: `TestRootConformance` 17.9 s, `TestEveryRetainedFixtureBitesThroughRegisteredOwner` 15.6 s, `TestWorkflowCadenceAnchorsRejectDeletionAndSwap` 5.3 s.
    - preflight, evidencecmd, probe and cmd/bench were flat: no test above about 2.3 s, and most 0.3-0.8 s.
11. Scheduling simulation (`scratchpad/sim.py`): it reproduces 477 s from the 23:48 log at p=4. Estimates for unmeasured packages are 0.8x their gate times, and adopt is taken as 25 s. This is an estimate, not a gate run.

---

## What I read and did not read

**Read:**
- Logs and gate code: gate logs `.logs/gate-*.out` and the 23:48 `.jsonl`, `internal/gate/phases.go`, `runner.go` (schedule), `internal/env/kit_run.go` and `env.go`, ADR 0024 and `docs/known-issues.md`.
- Profile and CLI help: `projects/benchkit.md` (phase table and host limits), `bench test --help`, `bench help`, and the `capture/ft217-build-timing` log.
- adopt production: `compatibility.go` (Doctor and legacyDoctor), `transaction.go`, `link.go`, `link_stage.go`, `decision.go` (head), `setup.go` (lines 320-390), `transaction/publication.go`.
- adopt tests, in full: `adopt_test.go`, `link_transaction_test.go`, `link_hook_test.go`, `compatibility_test.go`, `setup_test.go`, `setup_prompt_test.go`, `lifecycle_decision_test.go`, `decision_test.go`, `link_plan_test.go`, `lifecycle_unlink_test.go` (the post-golden version).
- adopt tests, pre-golden only: `lifecycle_consumption_test.go`, `lifecycle_special_test.go`, `lifecycle_preview_test.go`. I read `golden_test.go` after it landed, and only skimmed `lifecycle_unlink_test.go` there.
- repairtest: `session_test.go`, `repair_test.go`, `privacy_test.go`, and the first 70 lines of `kit_test.go`.
- Other packages: `internal/testreport/root_conformance_test.go` and parts of `check_test.go` and `environment_test.go`. `internal/conformance/harness_test.go` and parts of `fixture_bite_test.go`. `internal/systemtest/adoption_test.go`, parts of `owner_test.go` and `owner_landing_fixture_test.go`, and the system test name list. `internal/testrepo/working_tree.go`, `internal/gittest/gittest.go` (KitCopy), `cmd/bench` build-subject tests, `internal/worktree/clean_landed_hostile_test.go`, and `main_test.go`.

**Not read or not measured:**
- Gate log `gate-20261010T003228...` and the 00:32 numbers (used your figures).
- `internal/gitguard`, `intent`, `recordcmd`, `freshness`, `commitment`, `landing`, `status`, `diff`, `roadmap` and the rest of the long tail (each 8-15 s in the gate).
- Most system test bodies, `adopt/broker_test.go`, the remaining 30 or so adopt tests, and the `transaction` package tests beyond names.
- No whole-package adopt run on either TMPDIR, no whole system lane, and no real gate with the TMPDIR redirect.
- No `bench probe` witness for any "subsumed" claim.
- Whether the landing broker preserves the caller's TMPDIR.

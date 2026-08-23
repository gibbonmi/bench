# Worktree fresh-test invocation census

Measured and inspected 2026-08-23. The dynamic package sample began at
`0e17d428` and the static census finished after the shared checkout advanced to
`148f3a68`. FT113 subsequently landed at `d5915574`. Its only production change
under `internal/worktree` is a 13-line edit to `worktree.go`. It did not change
the test constructors counted here.

## Measurement validity

The clean pre-FT113 package run recorded during the initiating diagnosis was
137.86 seconds wall / 130.39 seconds reported by Go. The coordinator's later
clean delegate gate reported `internal/worktree` at 136.319 seconds.

Three sequential clean gates during the prerequisite landing repair and the
worktree-rule enforcement reported `internal/worktree` at 130.013, 125.779,
and 125.790 seconds. Their median is 125.790 seconds. The commits differ by
those small repairs, so this confirms the current magnitude without replacing
the first spec's exact-subject before measurement.

A JSON timing sample reported 167.41 seconds wall, 157.548 seconds from Go, and
65.98 seconds CPU (39% CPU/wall). It overlapped another delegate's Go
verification and the checkout advanced during the run. That result demonstrates
contention sensitivity but is not an exact-subject baseline. A whole-suite JSON
run was stopped when the overlap was identified. `strace` is unavailable, so
child counts below distinguish exact static expansions from source-reference
counts.

## Where the time accumulates

The contaminated JSON sample contained 282 top-level tests and 286 subtests.
Because the package has no `t.Parallel`, top-level elapsed values summed to the
package span. Thirty-eight tests at or above one second contributed 99.45
seconds. The ten slowest contributed 56.71 seconds:

| test | elapsed seconds |
|---|---:|
| `TestResumeLandCommandPublicRefusesDestructiveDestinationState` | 12.27 |
| `TestRecoveryPreservesEveryGitVisibleLayerWithoutMovingBranchOrIndex` | 8.55 |
| `TestLifecycleFaultBoundariesRemainLockedOrAbsent` | 7.00 |
| `TestLandCommandPublicRealGitJourney` | 5.89 |
| `TestExplicitApplyRevalidatesSafetyEvidence` | 5.72 |
| `TestResumeLandCommandRefusesWhenTerminalReceiptWasEvicted` | 4.29 |
| `TestCleanLandedApplyRefusesInitialDriftWithoutMutation` | 3.83 |
| `TestBuildOutputDeclarationFailsClosed` | 3.11 |
| `TestConcurrentCleanupRecordsOneTransaction` | 3.09 |
| `TestCleanLandedApplyRemovesAndSettles` | 2.96 |

No single test owns the long pole. The package is a serial sum of process and
filesystem waits.

## Invocation topology

| class | current topology |
|---|---|
| selected-binary builds | `buildLandingBinary` has 12 call sites. Each invokes `scripts/go-build.sh`. |
| nested Go commands | Those builds expand to 12 `go build` calls and 12 successor `go list -buildvcs=false -json -deps` calls used to publish the freshness seal. |
| nested `go test` | No static call site. The gate profile also forbids a nested ordinary test driver. |
| private Bench commands | The 12 build-bearing tests start 55 private Bench commands, primarily `worktree land`, plus one `worktree list`. |
| real repository construction | `newWorktreeRepo` has 123 static call sites. Each realized call initializes/configures a repository and adds and commits its base. Table loops increase the dynamic total. |
| Git helpers | Tests contain 226 `gitRun` references, 204 `gitOutput` references, and 19 direct `exec.Command("git")` sites. Production worktree code has 27 direct Git spawn sites plus the shared `internal/git` adapter. |
| gates | Landing fixtures install small prospective shell gates. They do not invoke the six-phase Bench gate or a successor test suite. |
| waits | One SIGINT helper briefly loops on `sleep 1`; the multi-second `time.After` and process-grace values are failure backstops, not green-path fixed sleeps. |

The package also contains 106 `t.Setenv`, 15 `t.Chdir`, two direct
`os.Chdir`, and zero `t.Parallel` references. `TestMain` replaces
process-global `BENCH_HOME` for the package. These owners make indiscriminate
parallelism unsafe before explicit inputs replace the globals.

## Historical comparison

The 2026-08-13 census on `a3b599ea` measured `internal/worktree` at 19.49
seconds and the whole fresh test phase at 31.90 seconds. The latter was floored
by `internal/publication`: `TestFixtureRegistryStagedOpsNeverCheckToolVersions`
spent 30.03 seconds calling through `http.DefaultClient` to
`http://127.0.0.1:1`, where the kernel retried the connection instead of
refusing immediately. The same background-context/default-client path remains
in the current tree.

The current worktree package is therefore a separate regression from the old
publication floor. Repeated builds are material—build-bearing tests contributed
34.68 seconds in the contaminated sample—but broad real-Git construction and
nested public Bench journeys own the larger cumulative topology.

## Consequences for shaping

- Remove duplicate selected-binary builds and policy-only repository
  materialization before pricing concurrency.
- Preserve real Git where Git supplies the behavior, while moving typed policy
  partitions behind the existing decision owners.
- Remove process-global inputs before any `t.Parallel` expansion.
- Treat publication's uncontrolled network wait as a separate bug before a
  whole-suite wall-clock budget becomes authoritative.
- Re-run a clean exact-subject package and whole-suite census after the first
  demand-reduction spec lands.

# Demand-reduction evidence (rows DM1, GF1, NP1)

This file is the one evidence owner for the first-spec measurements. It records
the baseline, the after measurements, and the demand counts. Each count names
the command or the method that produced it.

## Subject

- Commit: `03a5f736e0a11b18a23aa1e75c1c5033663c3a12` (retained integration tip;
  tickets 01–06 landed).
- Host: WSL2, `Linux 6.6.87.2-microsoft-standard-WSL2`, repository on Linux
  ext4 (`/dev/sdd`, `df -T .`). Other load can exist on this shared host; the
  `uptime` load averages beside each run record the observed condition.
- Go: `go1.25.0 linux/amd64` (`go version`).
- Cache state: normal Go build and test caches. No cache was cleared.
- Date: 2026-08-23.

## Baseline (before)

The spec is the source for the baseline. This run did not re-run the old tree.

- `specs/worktree-test-latency/spec.md` records three clean-gate
  `internal/worktree` spans: 130.013, 125.779, and 125.790 seconds.
- The same spec records the historical package span: 19.49 seconds
  (2026-08-13 census).
- The historical 31.90-second whole-suite floor included an unrelated
  30.03-second publication connection wait. A separate `$bench-debug` repair
  owns that wait.

## After (measured on this host, commit 03a5f736)

### Package spans, three runs

Command, run three times:

    go test -count=1 ./internal/worktree/...

Raw output:

Run 1 (`uptime`: `17:07:02 up 9:13, 1 user, load average: 0.12, 0.80, 0.80`):

    ok  	github.com/gibbonmi/bench/internal/worktree	58.022s
    ok  	github.com/gibbonmi/bench/internal/worktree/landingpolicy	0.002s
    ok  	github.com/gibbonmi/bench/internal/worktree/lifecyclepolicy	0.003s
    ok  	github.com/gibbonmi/bench/internal/worktree/reclaimpolicy	0.002s
    ok  	github.com/gibbonmi/bench/internal/worktree/refresh	0.025s

Run 2 (`uptime`: `17:08:10 up 9:14, 1 user, load average: 0.86, 0.90, 0.83`):

    ok  	github.com/gibbonmi/bench/internal/worktree	56.779s
    ok  	github.com/gibbonmi/bench/internal/worktree/landingpolicy	0.003s
    ok  	github.com/gibbonmi/bench/internal/worktree/lifecyclepolicy	0.004s
    ok  	github.com/gibbonmi/bench/internal/worktree/reclaimpolicy	0.002s
    ok  	github.com/gibbonmi/bench/internal/worktree/refresh	0.027s

Run 3 (`uptime`: `17:09:17 up 9:15, 1 user, load average: 1.15, 0.98, 0.87`):

    ok  	github.com/gibbonmi/bench/internal/worktree	56.898s
    ok  	github.com/gibbonmi/bench/internal/worktree/landingpolicy	0.002s
    ok  	github.com/gibbonmi/bench/internal/worktree/lifecyclepolicy	0.003s
    ok  	github.com/gibbonmi/bench/internal/worktree/reclaimpolicy	0.003s
    ok  	github.com/gibbonmi/bench/internal/worktree/refresh	0.025s

Parent-package spans: 58.022, 56.779, and 56.898 seconds (median 56.898). The
three pure policy packages and `refresh` each complete below 0.03 seconds.

### Whole-suite span, one run

Command:

    time go test -count=1 ./...

`uptime` before the run: `17:10:31 up 9:16, 1 user, load average: 0.79, 0.91,
0.85`. Result: 61 packages `ok`, zero `FAIL`. Wall time:

    real	1m10.122s
    user	1m21.044s
    sys	0m47.145s

Inside that run, `internal/worktree` reported 58.985 seconds and
`internal/publication` reported 1.855 seconds.

### Publication wait, reported separately

The historical 30.03-second publication connection wait was not visible in this
whole-suite run (`internal/publication`: 1.855 seconds). This evidence makes no
worktree improvement claim from that wait or from its absence. The
`$bench-debug` repair owns that transport defect.

## Demand counts, before and after

Before counts come from
`specs/worktree-test-latency/decisions/assets/worktree-test-invocation-census.md`
(static census at `148f3a68`). After counts come from static counts at
`03a5f736`; each row names its command or method.

| demand | before | after | after method |
|---|---|---|---|
| executable builds per package run | 12 (`buildLandingBinary` call sites; each ran `go build` plus `go list -json -deps`) | 1 direct-build, or 0 when the gate supplies `BENCH_RUN_BINARY` | `internal/worktree/test_run_test.go` owns selection: `testRunSelector` builds under `sync.Once`; `TestDirectRunBuildsAndSealsOneExecutableForAllJourneys` (SB1) proves one build, and `TestInheritedSelectionReachesJourneysUnchangedWithZeroBuilds` (SB2) proves zero. `rg -c buildLandingBinary internal/worktree` finds no call site. |
| real repositories, static constructor call sites | 123 (`newWorktreeRepo`) | 133 total through the harness: 119 `newWorktreeRepo`, 6 `journeyRepoOnBranch`, 5 `journeyRepo`, 3 `journeyFIFOWorktreeAdmin` | `rg -o '<name>\(' internal/worktree/*_test.go` per constructor, excluding the harness definitions in `journey_test.go`. Every constructor records to the harness effect log. Repository-backed policy partitions moved to the pure owners; the remaining sites are retained journeys and adapters, and the count no longer grows with policy partitions. |
| descendant starts outside the harness | 19 direct `exec.Command("git")` sites, plus 226 `gitRun` and 204 `gitOutput` references with per-test private builds | 0. All descendants route through the harness `descendant` constructor: 37 direct `descendant(` sites, 218 `gitRun(` sites, 217 `gitOutput(` sites, 3 `journeyStubGit(` sites | `TestSerialJourneyHarnessCensus` in `internal/worktree/journey_test.go` fails the package when a test file outside the harness starts a process. `rg -n 'exec\.Command' internal/worktree/*_test.go` finds only the two harness sites (`journey_test.go:109`, `test_run_test.go:68`). |
| environment mutations outside the harness | 106 `t.Setenv` sites across test files, plus a `TestMain` global | 0. Direct mutation exists only in the harness: `bindEnv` (`journey_test.go:92`), the package root in `main_test.go:33`. Tests call `bindEnv` at 101 sites, and each call is recorded | `TestSerialJourneyHarnessCensus` (`journey_test.go`) and `TestEffectBoundaryCensus` (`effect_census_test.go`) enforce the boundary; `rg -n 't\.Setenv\(' internal/worktree/*_test.go` confirms the single harness site. |
| directory changes outside the harness | 15 `t.Chdir` plus 2 `os.Chdir` sites | 0. Direct mutation exists only in the harness `chdir` (`journey_test.go:98`); tests call `chdir` at 39 recorded sites | Same census pair; `rg -n 't\.Chdir\(|os\.Chdir\(' internal/worktree/*_test.go` confirms the single harness site. |

## GF1: unchanged gate driver

The gate still runs one ordinary `go test -count=1 ./...` driver with normal
caches. This spec changed no gate phase, `-count=1` argument, package loop, or
cache policy. The gate argv tests
(`internal/gate/branch_native_phases_test.go` and
`internal/conformance/ordinary_build_census_test.go`) pin that argv and stayed
green in the landed gate runs for tickets 01–06. The measurements above also
use `-count=1` and normal caches.

## NP1: no parallelism added

The completed first-spec diff adds no `t.Parallel` call and no scheduler.
`rg -n 't\.Parallel' internal/worktree -g '*.go'` finds only the three
`purity_census_test.go` enforcement strings, which fail a policy package that
adds the call. The parent package census keeps descendants, environment, and
directories serial behind the one harness.

## Reproducibility

Run these commands at commit `03a5f736` on the reference WSL host, with normal
caches and no concurrent gate:

    uptime
    go test -count=1 ./internal/worktree/...   # run three times
    time go test -count=1 ./...                # run once
    go version

Record `uptime` beside each run. A concurrent gate, a cold cache, a visible
publication wait, or a moved source invalidates a reference result. Do not
change the target instead.

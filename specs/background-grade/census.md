# Background census

This record is the resource census that ADR 0024 requires before adoption of a background gate.
It measures the cost of one background gate proxy to a foreground test run on four pinned cores.
The measurement ran on 2026-10-27 at the source commit `84e9b1d372f0931e75db1072013e63cd4f94ed00`.

## Host

| Fact | Value |
|---|---|
| CPU | 12th Gen Intel(R) Core(TM) i5-12600K, 16 logical cores on the host |
| Cores of the measurement | 4, pinned with `taskset -c 0-3` (`taskset -c 0-3 nproc` prints 4) |
| Memory | `MemTotal: 24608948 kB` |
| Kernel | `6.6.87.2-microsoft-standard-WSL2` |
| `go env GOFLAGS` | `-p=4 -parallel=4` |
| Go | `go1.25.14 linux/amd64` |
| Disk scheduler | `none` on every block device, which includes `sdd`, the disk of the worktree |
| Bench | `/home/mgibs/.local/bin/bench` |

## Foreground command

Each foreground run used this command line from the worktree root, at the host configuration:

```
taskset -c 0-3 bench test --package ./internal/gate
```

The foreground run sets no GOFLAGS value of its own, so it uses `-p=4 -parallel=4` from the host.
One warm-up run went before the five alone runs. The warm-up run took 243.23 s, exited 0, and is not in the alone set.

## Proxy commands

The proxy runs the Go test commands of the whole gate, and it adds `-p=2 -parallel=2` to each command.
The source of the commands is the kit phase table in `internal/gate`.
The table has these three Go test commands:

- The whole test phase, from `wholeTestPhase`.
- The race phase, from `raceDriverArgv`.
- The system phase, from `SystemSuite`.

`BaseTestArgv` adds `-trimpath -count=1 -timeout=10m0s` to each of the three commands.

The proxy uses these values:

- `ROOT=/home/mgibs/.bench/worktrees/bench-2826441890/e76595bcb779900d284d6606cf28cf35-c7bb50cb85b54be3caa07403fc532593`
- `BIN=/tmp/claude-1000/-home-mgibs-workspace-bench/0edeb8f0-0137-45f7-b286-9187eb8eb2d8/scratchpad/run/bench-proxy`, built with `go build -trimpath -buildvcs=false -o $BIN ./cmd/bench`

Each proxy command runs in `$ROOT` with `BENCH_RUN_BINARY=$BIN` and `BENCH_KIT=$ROOT`.
These are the exact proxy command lines:

```
BENCH_CONFORMANCE_ROOT=$ROOT BENCH_CONFORMANCE_TIER=dev go test -trimpath -count=1 -timeout=10m0s -p=2 -parallel=2 "-exec='$BIN' gate-go exec" ./...
go test -trimpath -count=1 -timeout=10m0s -p=2 -parallel=2 -race -v ./internal/worktree ./internal/guards ./internal/census ./internal/responsebound -run '^(TestConcurrentCleanupRecordsOneTransaction|TestParallelJourneysShareTheHarnessSafely|TestParallelJourneysRecordEverySelection|TestScanTimeoutPreservesPartialRowsAndHonestCounts|TestScanEnumerationTimeoutUsesUnknownCounts|TestConcurrentRecordsKeepEveryLine|TestOwnerSerializesConcurrentWrites)$'
BENCH_SYSTEM_ROOT=$ROOT go test -trimpath -count=1 -timeout=10m0s -p=2 -parallel=2 -tags=system ./internal/systemtest
```

One proxy pass runs the three commands in sequence.
A restart loop starts each pass with this wrapper:

```
taskset -c 0-3 nice -n 19 ionice -c 3 bash proxy.sh "$ROOT" "$BIN" <log prefix>
```

The loop starts the next pass when a pass ends, until the last beside run ends.
The proxy started 22 s before the first beside run.

| Pass | Command | Wall (s) | Exit |
|---|---|---|---|
| 1 | whole test | 1000.42 | 1 |
| 1 | race | 29.40 | 1 |
| 1 | system | 1.11 | 1 |
| 2 | whole test | stopped after the fifth beside run | none |

The proxy restarted once, during the fourth beside run.
So each beside run had the proxy beside it for the full run.

## Proxy priority

A sampler read the priority of each process in the session of the proxy.
It read the nice value with `ps -o pid,ni,psr,args -s <session>`, and the I/O class with `ionice -p <pid>`.
It took one sample before the first beside run, one sample 15 s into each beside run, and one sample after the last beside run.

| Process | Nice | I/O class | Affinity | Samples |
|---|---|---|---|---|
| `proxy.sh`, `go test`, each test binary, and each test child | 19 | `idle` | `0-3` | 32 nice readings and 31 I/O class readings in 7 samples |
| The restart loop (`bash -c`) | 0 | `none: prio 0` | `0-15` | 7 samples |

The restart loop only waits for each pass, and it does no test work.
Six more I/O class reads found a process that had already exited.
No sampled proxy process ran at a nice value other than 19 or at an I/O class other than `idle`.

## Runs

| Run | Wall (s) | Exit | Package status |
|---|---|---|---|
| alone-1 | 239.45 | 0 | pass, 728 tests |
| alone-2 | 240.74 | 0 | pass, 728 tests |
| alone-3 | 236.88 | 0 | pass, 728 tests |
| alone-4 | 237.73 | 0 | pass, 728 tests |
| alone-5 | 234.61 | 0 | pass, 728 tests |
| beside-1 | 267.24 | 0 | pass, 728 tests |
| beside-2 | 310.58 | 0 | pass, 728 tests |
| beside-3 | 325.88 | 0 | pass, 728 tests |
| beside-4 | 273.81 | 0 | pass, 728 tests |
| beside-5 | 298.98 | 0 | pass, 728 tests |

The wall time is the time from the start to the end of the foreground command.
The beside-3 run also printed a timeout warning: the `internal/gate` package took 304328 ms of its 600000 ms limit.

## Result

| Measure | Value |
|---|---|
| Alone median | 237.73 s |
| Alone spread | 234.61 s to 240.74 s, a range of 6.13 s (2.6% of the alone median) |
| Beside median | 298.98 s |
| Ratio, beside median to alone median | 1.258 |
| Bound, 1.2 times the alone median | 285.28 s |
| Red foreground runs | 0 |

The beside median is more than 1.2 times the alone median, so the bound does not pass.
No foreground run went red.

## Limits of the measurement

The proxy binary had no seal file, so each proxy test that selects a sealed run binary went red at once.
These tests are in `internal/adopt`, `internal/testreport`, `internal/worktree`, and `internal/systemtest`.
The system command exited after 1.11 s for this reason.
Thus the proxy did less work than a whole gate with a sealed binary.
More proxy work adds contention, so the gap makes the measured cost low, not high.

The proxy also had red tests in `internal/gate` and `internal/subprocess/protection`, and this record does not diagnose them.
The foreground runs of `internal/gate` did not go red.

The disk scheduler `none` does not apply I/O priority classes.
So the idle I/O class of the proxy gave the foreground no I/O priority on this host.

## Reviewer waiver

On 2026-10-27, the reviewer accepted a beside-to-alone ratio of 1.258 and directed the build to continue.
The reason comes from the debug loop of the orchestrator.
In that loop, the foreground slowdown stayed the same under nice 19, under SCHED_IDLE, and with one thread pinned to each physical core.
So the cost comes from shared resources other than CPU time, and no unprivileged priority controls these resources on this host.
Only `-p=1 -parallel=1` decreased the slowdown, and it decreased the slowdown only in part.
In that loop, a rerun with a sealed proxy binary showed the same slowdown pattern.

Verdict: stop

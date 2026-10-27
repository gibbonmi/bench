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
The foreground runs and the proxy ran in this one worktree, on the disk `sdd`.
`bench test` gives its Go child the Bench build cache `$HOME/.cache/bench/go-build`.
The proxy commands set no `GOCACHE` value, so they used the default Go build cache `$HOME/.cache/go-build`.
Both caches are on `sdd`.
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

The loop started the next pass when a pass ended, until the last beside run ended.
Pass 1 started at epoch 1793130668.46, 22 s before the first beside run.

| Pass | Command | Wall (s) | Exit |
|---|---|---|---|
| 1 | whole test | 1000.42 | 1 |
| 1 | race | 29.40 | 1 |
| 1 | system | 1.11 | 1 |
| 2 | whole test | stopped after the fifth beside run | none |

The proxy restarted once, during the fourth beside run.
Pass 1 ended and pass 2 started at epoch 1793131699.40, as `proxy.cmds` records.
The 15 s sample of beside-4 was at epoch 1793131611.52, so the restart came about 103 s into that 273.81 s run.
So each beside run had the proxy beside it for the full run.

## Proxy priority

A sampler read the priority of each process in the session of the proxy.
It read the nice value with `ps -o pid,ni,psr,args -s <session>`.
Then it listed the session again with `ps -o pid= -s <session>`.
For each process in that second list, it read the I/O class with `ionice -p <pid>` and the affinity with `taskset -cp <pid>`.
It took one sample before the first beside run, one sample 15 s into each beside run, and one sample after the last beside run.

| Process | Nice | I/O class | Affinity | Samples |
|---|---|---|---|---|
| `proxy.sh`, `go test`, each test binary, and each test child | 19 | `idle` | `0-3` | 7 samples, with the reading counts in the next table |
| The restart loop (`bash -c`) | 0 | `none: prio 0` | `0-15` | 7 samples |

The restart loop only waits for each pass, and it does no test work.
This table gives the readings of the proxy processes in each sample:

| Sample | Nice readings | I/O class and affinity readings | Reads of an exited process |
|---|---|---|---|
| Before beside-1 | 5 | 5 | 0 |
| beside-1 at 15 s | 4 | 4 | 1 |
| beside-2 at 15 s | 6 | 6 | 0 |
| beside-3 at 15 s | 5 | 4 | 0 |
| beside-4 at 15 s | 3 | 3 | 1 |
| beside-5 at 15 s | 4 | 4 | 4 |
| After beside-5 | 5 | 5 | 0 |
| Total | 32 | 31 | 6 |

Processes started and exited between the two lists of a sample.
In the beside-3 sample, the first list showed a defunct `git` child at nice 19, and the second list did not show it.
So that sample has one more nice reading than I/O class readings.
The 6 reads of an exited process each returned `No such process`, and the first list did not show these processes.
No sampled proxy process had a nice value other than 19, an I/O class other than `idle`, or an affinity other than `0-3`.

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
These tests are in these four packages:

- `internal/adopt`
- `internal/testreport`
- `internal/worktree`
- `internal/systemtest`

The system command exited after 1.11 s for this reason.
Thus the proxy did less work than a whole gate with a sealed binary.
This record infers that more proxy work adds contention, so the gap makes the measured cost low, not high.
No run of this census measured that inference.

The proxy also had red tests in `internal/gate` and `internal/subprocess/protection`, and this record does not diagnose them.
The foreground runs of `internal/gate` did not go red.

The disk scheduler `none` does not apply I/O priority classes.
So the idle I/O class of the proxy gave the foreground no I/O priority on this host.

## Reviewer waiver

On 2026-10-27, the reviewer accepted a beside-to-alone ratio of 1.258 and directed the build to continue.
The reason comes from the debug loop of the orchestrator.
The runs of this section are the orchestrator's runs, not runs of the census above.

The loop foreground is a subset of `internal/gate`:

```
taskset -c <cpus> bench test --package ./internal/gate --run '^Test(Checkpoint|Complete|Delta)'
```

A beside variant ran the whole-test proxy in a restart loop on the same CPUs, and it started the foreground after a 60 s warm-up.
The proxy used the sealed primary binary `/home/mgibs/workspace/bench/dist/bench` as `$BIN`, so these runs are the sealed-binary rerun:

```
BENCH_RUN_BINARY=$BIN BENCH_KIT=$ROOT BENCH_CONFORMANCE_ROOT=$ROOT BENCH_CONFORMANCE_TIER=dev taskset -c <cpus> <priority> go test -trimpath -count=1 -timeout=10m0s -p=2 -parallel=2 "-exec=$BIN gate-go exec" ./...
```

Each variant ran the foreground three times, and every run exited 0.
CPUs 0-3 are two physical cores with two SMT threads each (`lscpu -e`).
CPUs 0,2,4,6 are one thread on each of four physical cores.

| Variant | CPUs | Wall (s) | Median (s) | Ratio to its alone median |
|---|---|---|---|---|
| Alone | 0-3 | 55.77, 58.54, 59.46 | 58.54 | 1.000 |
| Beside, nice 19 and idle I/O | 0-3 | 57.54, 62.63, 69.78 | 62.63 | 1.070 |
| Beside, SCHED_IDLE (`chrt -i 0`) and idle I/O | 0-3 | 60.56, 63.57, 68.77 | 63.57 | 1.086 |
| Beside, `-p=1 -parallel=1`, nice 19, and idle I/O | 0-3 | 59.36, 61.27, 65.12 | 61.27 | 1.047 |
| Alone | 0,2,4,6 | 56.24, 58.74, 57.97 | 57.97 | 1.000 |
| Beside, nice 19 and idle I/O | 0,2,4,6 | 57.85, 63.11, 70.34 | 63.11 | 1.089 |

The foreground slowdown stayed at about the same ratio under nice 19, under SCHED_IDLE, and with one thread on each physical core.
So the cost comes from shared resources other than CPU time, and no unprivileged priority controls these resources on this host.
Only `-p=1 -parallel=1` decreased the slowdown, and it decreased the slowdown only in part.
The sealed proxy binary also slowed the foreground under each priority variant.
The loop foreground is a subset of the census foreground, so its ratios do not compare directly with the census ratio of 1.258.

Verdict: stop

# Measure the background census

Blocked by: none
Writes: capture/build-log/background-grade.log (new)
Covers: BG93, BG94, BG95

## What to build

Before any code change, measure what one background gate costs the foreground on a runner-sized host.
This measurement is the fresh resource census that ADR 0024 requires before adoption.

Pin every process of the measurement to four cores, for example with `taskset -c 0-3`, to model a runner with two to four cores.
The foreground workload is `bench test --package ./internal/gate` at the host configuration.
Run it five times alone, and then five times beside one proxy of a whole gate.
The proxy runs the Go test commands of the whole gate at `-p=2 -parallel=2`, at nice 19, and in the idle I/O class.
The Bench width input does not exist yet, so the proxy passes those flags itself.
Restart the proxy when it ends before the foreground run ends, so the foreground run always has the proxy beside it.

Write the build log `capture/build-log/background-grade.log`.
The log records the exact proxy and foreground command lines, each wall time, each exit, the two medians, and the spread of the alone runs.
It also records the host facts: the cores, the memory, the kernel, and `go env GOFLAGS`.
The log shows the priority that the proxy actually ran at, read from the running proxy process.

The bound passes when two conditions hold. The median beside the proxy is at most 1.2 times the median alone. No foreground run goes red.
If the bound does not pass, the build stops, and the orchestrator reports the log to the reviewer.

## Acceptance

- [ ] The log holds five alone wall times and five beside wall times of `bench test --package ./internal/gate` on four pinned cores.
- [ ] The log holds the exact proxy command lines, with `-p=2` and `-parallel=2`.
- [ ] The log shows nice 19 and the idle I/O class, read from the running proxy process.
- [ ] The log holds both medians, the spread of the alone runs, and the host facts.
- [ ] The last log line states `pass` only when the beside median is at most 1.2 times the alone median and no foreground run is red. Otherwise it states `stop`.

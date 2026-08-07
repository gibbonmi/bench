# Gate memory and process profile

Evidence asset for `decisions/gate-budget.md` #1–#5, #8, and #20. Measured
2026-08-06 local time (`2026-08-07T00:01:39Z`–`00:06:21Z`) on commit
`6607236c0ae6591805ecaf5572bf85ecaad20a80`, a 12-online-core host, with the
foreign FT187 retro left untouched. No other gate process was active. The sampled
gate was green in 282 seconds.

## Method

The sealed binary was rebuilt through the canonical builder, then one fresh gate
was launched under a process-tree sampler:

```text
GOMAXPROCS=2 GOFLAGS=-p=2 bash scripts/go-build.sh /home/devuser/workspace/bench /home/devuser/workspace/bench/dist/bench
setsid env GOMAXPROCS=2 GOFLAGS=-p=2 bench gate --fresh
```

The sampler recorded `/proc/meminfo`, cgroup `memory.current`, `memory.peak`,
`memory.stat`, and `io.stat`, per-process RSS and `smaps_rollup` PSS, process
ancestry and command lines, Go/compiler/linker/gate/canary/contract process
counts, Go-cache size, and new temporary-tree sizes. The cgroup's pre-existing
`memory.peak` was not reset, so the run's peak is the maximum sampled
`memory.current`, not the retained 7.514 GB historical peak.

The sampler's first ancestry closure admitted unrelated ancestors into its
aggregate descendant columns. Raw per-PID rows were intact; the process figures
below come from a corrected offline descendant closure. Meminfo, cgroup, I/O,
cache, and temporary-tree figures were unaffected.

## Memory result

| observable | before | sampled peak | after |
|---|---:|---:|---:|
| cgroup current | 2.929 GB | 6.091 GB | 3.397 GB |
| cgroup anonymous | 0.314 GB | 2.002 GB | 0.321 GB |
| cgroup file | 2.334 GB | 3.601 GB | 2.770 GB |
| cgroup kernel | 0.279 GB | 0.467 GB | 0.305 GB |
| `Cached` | 1.472 GB | 2.684 GB | 1.849 GB |
| `SReclaimable` | 0.417 GB | 0.477 GB | 0.447 GB |
| `SwapFree` | 16.982 GB | 16.982 GB | 16.983 GB |

`MemTotal - MemFree`, the cache-inclusive desktop-style reading, reached
7.371 GB. `MemTotal - MemAvailable`, the pressure-oriented reading, reached only
3.072 GB; `MemAvailable` never fell below 5.255 GB and swap did not deplete.
The near-8 GB symptom is therefore not 8 GB of unreclaimable process memory.

The gate added 1.688 GB anonymous, 1.268 GB file-backed, and 0.189 GB kernel
memory over the cgroup baseline. Anonymous memory returned almost to baseline
after exit while 436 MB of additional file cache remained. The absolute peak was
file-dominated, but duplicated process heaps were the larger transient increment.

At the corrected descendant-PSS peak, 29 seconds into the run, the gate had 80
descendants with 3.033 GB RSS, 1.943 GB PSS, 1.720 GB anonymous PSS, and 222 MB
file PSS. The sample carried 17 `go`, seven `compile`, 33 `bench`, eight
canary-related, and 21 contract/artifact-related processes, including six
simultaneous fixture inner gates. The seven compiler processes alone accounted
for roughly 875 MB PSS.

## Storage and I/O result

The cgroup wrote 4.25 GB and read 142 MB over the run. During the first 34
seconds it wrote 2.153 GB, an average of about 63 MB/s, consistent with the
operator-observed 22 MB/s sustained activity and 130 MB/s burst. Live new
temporary trees peaked at 333 MB and the Go build cache grew by 152 MB. Repeated
creation, reading, and deletion therefore creates substantially more I/O and
page-cache population than the live generated-tree footprint suggests.

Clearing the Go cache before the run would have changed the warm-cache workload.
Its 152 MB growth is material storage churn but is not the source of the
cache-inclusive 7.371 GB peak.

## Concurrency finding

Primary and stripped-subject schedules run concurrently, and each scheduler
launches every dependency-ready outer phase without an admission bound. Seven
outer phase roots overlapped in the sampled run. Independently, canary computed
six fixture workers from the host width and launched six sibling inner gates.
Inner mode removes canary before owner selection and schedules its own phases
sequentially, so this is not recursive canary execution.

The exported `GOMAXPROCS=2 GOFLAGS=-p=2` did not constrain that outer fan-out.
The current closed subject environment carries only `PATH` plus manifest-declared
variables, and the manifest declares neither setting. Canary consequently read
the 12-core host width, chose six workers at inner width two, and pinned only the
inner children. The profile therefore confirms the map's process-boundary budget
destination while showing that current ambient width exports are not yet an
effective outer-gate budget.

This run does not price reserve `r`, its span-inflation threshold, or the two Go
grant splits. It predates the remaining decision-seam and conformance-harness
demand reductions, and one run on one host cannot state variance. It is entry
evidence for the mechanism and for #20's future measurement shape, not an answer
to #8 or #20.

# Gate trustworthiness report — runtime-contract flakes

Status: **fix landed and validated under load; diagnosis closed ahead of the
confirming dump.** Repository `bench`, branch `main`.

The body below is the diagnosis as written at baseline commit `7699473`, kept for
its evidence table and traps. Read this header first — it supersedes the body
wherever they disagree.

- **Fix landed** in `380fd00` (marker deadlines 5s→60s in R14 and 2s→60s in the
  losing-racer repair test, R14 fast-fail on early child exit, `[DEBUG-a4f2]`
  diagnostics promoted untagged) and `daf81d1` (data race on the SIGQUIT-timeout
  path). `specs/load-tolerant-marker-deadlines.md` is `Status: implemented`.
  The chosen option was **tolerate slow persistence**, taken deliberately ahead
  of the confirming goroutine dump — see that spec's reviewer-closed note.
- **Acceptance load window ran 2026-07-23: 3/3 green.** Wall clocks ~570s, 395s,
  394s at guest load averages 26→58 — every run slower than the 352-465s band in
  which the old deadlines failed 2/2, and slower than any recorded green. Zero
  hits of `did not start`, `did not reach pending`, or `RemoveAll cleanup`.
- **The fsync hypothesis remains unconfirmed and now untestable by this route.**
  The 60s deadlines absorb the stall, so R14 no longer fails under load and no
  goroutine dump will be produced. The diagnostics stay compiled in: if a stall
  ever exceeds 60s, the failure carries the per-thread table and the dump.
- **Still open:** the `.bench-contract-env` cleanup flake (below, uninvestigated;
  it did not fire in any of the three runs).

The goal was to make `bench gate` give a trustworthy verdict again. Two tests in
`internal/contract/runtime` fail only under full-gate load, which cost three
identical full-gate runs to land one commit. This report records what is proven,
what is disproven, what is still open, and the exact next action.

## The reproducer

**Only the real gate reproduces it.** No cheaper harness works — this is the
single most useful result here, because it rules out the whole class of
approximations a fresh session would otherwise try first.

```
bash bin/bench.sh gate
```

Then grep the output for `did not start`, `did not reach pending`, or
`RemoveAll cleanup`. A run takes ~4-8 minutes.

| load shape | runs | flake hits |
|---|---|---|
| 32 synthetic CPU workers | 3 | 0 |
| parallel contract tree + 2 loaders | 4 | 0 |
| gate's contract phase + 16 CPU workers | 3 | 0 |
| R14 test alone + 128 CPU workers | 6 | 0 |
| 4 GB memory ballast + real gate | 2 | 0 |
| 6 GB memory ballast + real gate | 2 | 0 |
| **real `bench gate`, machine under user load** | **2** | **2** |
| real `bench gate`, clean idle machine | 5 | 0 |
| real `bench gate`, reviewer-spun load (loadavg 20-90) | 4 | 1 |
| R14 tests alone + 4 guest-side fsync hammers (dd conv=fsync, same ext4) | 6 | 0 |
| real `bench gate`, host-side load, **after the 60s-deadline fix** | 3 | 0 |

The reviewer-spun-load row is from 2026-07-22: three runs at load averages the
original failures never reached stayed green, and the one hit (run 3) landed at
comparable load — so load magnitude alone is not the trigger either; some
co-occurring ingredient (likely I/O or scheduling shape, not CPU) still matters.

## What is proven

**1. The child process is wedged, not slow.** *(Superseded in part — see the
honest correction under the hypothesis section: the measurement behind this
claim never exercised fsync, so wedged-vs-stalled is open pending the
goroutine dump.)* This was the load-bearing finding, and it is why raising the
timeout was originally ruled out.

`startStory5GateOwner` spawns `bash bin/bench.sh gate` and waits 5s for the gate
script to write its started-marker. Instrumenting the failure showed:

```
--- FAIL: TestFT78Story5ProofLedger/R14/locked-pending-inspection (6.00s)
    [DEBUG-a4f2] child STILL ALIVE at deadline
    [DEBUG-a4f2] child output:      (empty)
```

The child is alive, has produced no output, and has not reached the gate script.

**2. The 5s deadline has ~160x headroom, so it is not the bug.** Measured
directly: spawn `bash bin/bench.sh gate` in a fixture repo whose `.bench/gate.sh`
timestamps itself, and time until that script executes. At load average 15 the
result was **30-36 ms across 10/10 samples**. For contention alone to blow a 5s
budget the box would need ~160x oversubscription; the gate reaches roughly 8x.

**3. The wedge is inside `dist/bench gate-run` itself, past lock creation,
before `.bench/gate.sh`.** From the 2026-07-22 hit
(`R14/interrupted-pending-inspection`), the first flake caught with the dump
live: the process group contains exactly one process — `dist/bench gate-run`
(bash had already exec'd into it), alive at the deadline, zero output, gate
lock **created** (`lock=true`), main thread parked in `futex_wait_queue`. So
the shell wrapper is exonerated, the child made real progress (took the lock),
and then blocked somewhere between lock acquisition and spawning the gate
script. The per-process WCHAN cannot name the blocked syscall — in a Go binary
the main thread futex-waits while whichever thread runs the blocking syscall
sits elsewhere — which is why the instrumentation now dumps per-thread state
and SIGQUITs the child for a goroutine dump (see below).

**4. The failure message is misleading and should be fixed regardless.**
`gate owner did not reach pending state` is emitted whether the child wedged,
died instantly, or never started — the test set no `cmd.Stdout`, so a failing
child's stderr went to `/dev/null`. Any future debugging is blind without this.

## What is disproven

- **Early exit / exec failure** (child dies against a `dist/bench` being rewritten,
  ETXTBSY, partial binary). Falsified: the child is alive at the deadline.
- **CPU contention as the ingredient.** 16 runs across four CPU-saturating shapes,
  including 128 busy workers against the test alone, produced zero hits.
- **Memory pressure as a sufficient ingredient.** 4 GB and 6 GB of tmpfs ballast
  produced zero hits in 4 gate runs. Caveat: ballast is *inert* memory; the
  real-world trigger (Chrome + Steam resident) also churns CPU, memory and I/O, so
  this weakens the hypothesis without killing it.

## The strongest remaining hypothesis

**fsync stall in `durableReplace` under host-side I/O pressure** (formed
2026-07-22 from the code read that proof item 3 made possible). Between taking
the gate lock and spawning `.bench/gate.sh`, `gate-run` persists the pending
verdict via `durableReplace` (`internal/gate/verdict.go`), which fsyncs the
temp file **and** the containing directory. On WSL2, ext4 lives on a VHDX
behind the Windows host's I/O scheduler: when the host churns disk (Chrome +
Steam do; synthetic CPU workers and inert memory ballast do not), an fsync can
stall for many seconds. This explains every row of the evidence table at once,
including why load *magnitude* alone (loadavg 90, zero hits) doesn't trigger.

Falsifiable prediction: the SIGQUIT goroutine dump on the next hit shows a
goroutine blocked in `os.File.Sync` inside `durableReplaceWithEngine`, and the
per-thread `ps` shows one thread in a writeback/jbd2 wait-channel rather than
futex.

**Honest correction to proof item 2:** the ~160x-headroom measurement was
taken under CPU load, which never stretches fsync. If this hypothesis holds,
the child is not wedged — it is stalled in one legitimately slow syscall that
only I/O pressure stretches, and the "don't raise the deadline" conclusion is
open again. Whether the fix is a longer/deadline-aware wait, moving the
durable write, or tolerating slow persistence is a reviewer decision once the
dump confirms the site.

The earlier framing — whole-machine pressure, correlated with total gate
duration — remains consistent and is subsumed by the fsync hypothesis:

| condition | gate wall clock | hits |
|---|---|---|
| failing runs (Chrome + Steam resident) | 352-465 s | 2/2 |
| clean-machine runs | 231-240 s | 0/5 |

The failing runs were ~1.5-2x slower end to end. Whatever stretches the whole
gate also triggers the wedge.

**The decisive next experiment:** reproduce with Chrome and Steam actually
running, with the process-group instrumentation in place (below). That dump names
the blocking syscall directly and should end the guessing in one hit.

## Instrumentation currently in the tree

`internal/contract/runtime/runtime_gate_action_proof_test.go` carries these
diagnostics. They are now **committed and permanent, with the `[DEBUG-a4f2]`
tags removed** — `rg 'DEBUG-a4f2'` finds only records like this report, never
code. What they do:

- captures the spawned child's stdout/stderr, previously discarded;
- reports whether the child exited early or was alive at the deadline;
- on a miss, dumps the child's whole process group **per-thread** with kernel
  wait-channels (`ps -eLo pid,tid,pgid,stat,wchan:30,etime,args`) and the gate
  lock's existence;
- then sends the process group **SIGQUIT** and waits up to 3s for the reap, so
  the Go runtime's full goroutine dump lands in the captured stderr. On the
  next hit this names the exact blocked line in `gate-run` — no further
  guessing needed.

The per-thread and SIGQUIT steps were added 2026-07-22 after the first
instrumented hit showed the per-process dump was one level too coarse (proof
item 3). They never caught a flake: the 60s deadlines landed first and stopped
the failure from firing. They stay as the permanent diagnostic for a genuine
hang past 60s.

It also moves the single `cmd.Wait()` into `startStory5GateOwner` — `stop()`
previously called `Wait` too, and a double `Wait` breaks the passing path.

**This instrumentation is gate-green** and weakens no assertion — it only
enriches an existing failure message. Promotion to a permanent diagnostic was
the decision taken in `380fd00`.

## Third flake — cross-runtime corroboration of the fsync hypothesis

2026-07-22, run 5 of the load window:
`TestBinaryRepairContracts/repair_losing-racer_cleanup_contract_failed`
(`internal/contract/surface/binary_repair_hardened_test.go:60`,
`repair did not reach synchronization marker`, 2s deadline). Same shape as
R14: a spawned bench child misses a marker deadline under load. The decisive
detail: the repair child is **Node**, not Go (`bin/bench-repair-binary.mjs`),
and its last operation before writing the ready marker is `await fh.sync()` —
an fsync of the downloaded temp binary. Two independent runtimes stalling at
the last pre-marker fsync is strong circumstantial support for the fsync
hypothesis; a heavy-interpreter cold page-in under host I/O remains the
runner-up explanation (node's startup is also inside the missed 2s window).
This flake had never fired before in any recorded run.

## Second flake — not yet investigated

`TestRuntimeGateContracts/bench_gate_rebuilt_self-host_contract` fails with
`TempDir RemoveAll cleanup: unlinkat .../001/.bench-contract-env: directory not
empty` — a child process still writing inside `t.TempDir()` when Go's cleanup
runs. Passes in 0.4s isolated. It did not fire in any run during this session, so
there is no fresh evidence beyond the original report. `.bench-contract-env` is
created by `isolatedEnv` in `internal/contract/command.go:220` and holds the
fixture's `HOME`, XDG dirs, and `BENCH_HOME`.

## Traps that cost time here

Each of these produced a red that looked like a finding and was not. A fresh
session should recognise them immediately.

1. **Never build `dist/bench` with plain `go build`.** Use
   `bash scripts/go-build.sh <root> <out>`. The helper stamps the version from
   `package.json`; without it, `bench version output contract` and
   `upgrade downgrade refusal` fail deterministically 3/3. Verify with
   `./dist/bench version` → `bench 0.2.0 (linux/amd64)`.
2. **Never write inside the repo while a gate is running.** A `bench idea` call
   mid-run mutated `IDEAS.md` and failed an otherwise all-green gate with
   `gate subject changed during execution`.
3. **Never stop a gate by killing only the wrapper.** Doing so orphans the tree:
   `bench canary` reparents to init, keeps spawning nested `gate-phases` children
   indefinitely, and holds the gate flock. Every later `bench gate` then fails in
   0s with the misleading `gate execution already in progress`. Kill the process
   group, then confirm with
   `pgrep -af 'dist/bench (canary|gate-phases)'` before trusting any result.
   This is parked in `IDEAS.md` as a real product defect, not only a local hazard.

## Next actions, in order

Items 1 and 2 (reproduce under load, then fix) are **done** — see the status
header. What remains:

1. Reviewer: veto or accept the four defaulted decisions flagged in
   `specs/load-tolerant-marker-deadlines.md`, then
   `bench spec retire load-tolerant-marker-deadlines` and retire this report
   (fold the traps and the evidence table into the project profile and
   learnings, then delete it).
2. Investigate the `.bench-contract-env` cleanup flake.
3. Close via `/bench-final-check`. That phase flips
   `specs/consumer-payload-and-phase-contract.md` from `Status: staged` using
   `bench spec implemented` — never hand-write the status line — and deletes
   `reviews/consumer-payload-and-phase-contract.md` in the same green commit.

## Work this diagnosis produced

All of it is now committed; the tree is clean. Retained for provenance.

| path | change | state |
|---|---|---|
| `specs/consumer-payload-and-phase-contract.md` | story 2 acceptance-coverage row reworded so its red signal matches the design and the landed test (a file under a kit-only subtree is not written) | done; `bench coverage --check` green, 12 rows valid |
| `internal/contract/runtime/runtime_gate_action_proof_test.go` | diagnostics plus the 60s deadline | committed in `380fd00`/`daf81d1`; permanent |
| `IDEAS.md` | three parked ideas (gate speed, core capping, orphaned-gate/stale-lock) | parked |
| `ROADMAP.md` | FT88 row proposed for this work | awaiting reviewer veto |

The `consumer-payload-and-phase-contract` status flip and the `reviews/` deletion
belong to that spec's own closing green commit, not to this report.

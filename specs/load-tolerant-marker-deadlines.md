# load-tolerant-marker-deadlines

Status: implemented

## Problem

Two contract tests spawn a bench child and wait a short fixed deadline (5s and
2s) for it to write a synchronization marker. On this WSL2 host, whenever the
Windows host contends the VHDX with real I/O, the child's pre-marker durable
write (fsync) can stall for seconds, the deadline misses, and the gate goes red
on work that is fine. Diagnosis is in `GATE-REPORT.md`: three hits across two
independent runtimes (Go `gate-run` in `durableReplace`, Node
`bench-repair-binary.mjs` at `fh.sync()`), zero hits under any CPU or inert
memory load. The leading hypothesis — corroborated across two runtimes but
**not yet confirmed** by the goroutine dump the report's next-action #1 seeks
— is that the child is stalled in a legitimate slow syscall, not wedged. The
false reds cost a full gate re-run (~4-8 min) each and erode trust in the
oracle.

## Solution

Make the two proven marker waits load-tolerant without weakening what they
prove: raise the deadline to 60s, fail fast the moment the child *dies* instead
of burning the deadline, and keep the diagnosis-grade failure output (child
output, per-thread wait-channels, SIGQUIT goroutine dump) as the permanent
failure message so any future miss self-documents. The production durable
writes are untouched — they are designed crash-safety properties, and the
reviewer closed that decision in conversation on 2026-07-22.

**Reviewer-closed, 2026-07-22: fix ahead of the confirming dump.** The
diagnosis's decisive experiment (an R14 hit carrying the goroutine dump) has
not landed; the reviewer approved proceeding on the corroborated hypothesis
with the residual risk named: if the true cause is a wedge rather than a
stall, the 60s deadline converts a fast false red into a slower true red —
never a false green — and the permanent diagnostics name the blocked line on
that first slower failure.

This spec was written same-session on the reviewer's explicit instruction (the
batch-drain override in `/bench-write-spec`); there is no decision map behind
it. Defaulted decisions are flagged below for post-hoc veto.

## User stories

1. As the gate operator, I want `startStory5GateOwner` to wait up to 60s for
   the started-marker while failing immediately if the child exits first, so
   that a host-side I/O stall no longer reds the gate but a dead child still
   fails in milliseconds.
   Line: `gpt-5.6-terra` / medium. The edit reorders a wait loop and adds an
   in-loop exit-channel check in the same code that already carried a
   double-`Wait` defect, so it is concurrency-semantic rather than
   mechanical, and the gate cannot observe the fast-fail branch.
2. As the next debugger, I want the R14 helper's failure message to
   permanently carry the child's captured output, the per-thread wait-channel
   dump, and the SIGQUIT goroutine dump — with the `[DEBUG-a4f2]` tags and
   debug naming removed — so that a future miss names its blocked line without
   a reproduction campaign.
   Line: `gpt-5.6-luna` / medium. The diagnostic code already exists and runs
   gate-green; the story renames and re-comments it into permanent form.
3. As the gate operator, I want `waitForRepairMarker` to wait up to 60s, so
   that the losing-racer repair test survives the same host-side stall that
   hit it on 2026-07-22.
   Line: `gpt-5.6-luna` / medium. One constant and one comment at an existing
   helper.
4. As the next debugger, I want the losing-racer repair test to capture the
   child's stdout/stderr and include it in the marker-timeout failure, so
   that a repair-side miss is not blind the way the R14 miss was before
   instrumentation.
   Line: `gpt-5.6-luna` / medium. Mirrors the capture pattern story 2 makes
   permanent, at a second test site.

## Implementation decisions

- Deadline value is **60s** at both sites (defaulted — veto point). Rationale:
  observed stalls miss 2-5s windows; 60s is far above any plausible fsync
  stall yet still bounds a genuine hang to one minute, and a dead child fails
  fast regardless.
- The R14 helper keeps the single-`Wait` ownership introduced during
  diagnosis: `startStory5GateOwner` owns the `Wait` via its exit channel and
  `stop()` receives from that channel; the pre-diagnosis double-`Wait` was a
  real defect on the passing path.
- Marker check ordering: poll the marker **before** the exit-channel check in
  each loop iteration. This is defensive ordering for the newly added in-loop
  exit check — no concrete cross-subtest kill race exists today (each subtest
  owns its fixture and its owner) — so a marker-then-exit sequence, however it
  arises, resolves as success rather than a false "exited before pending".
- `waitForRepairMarker` gets **no** exit fast-fail (defaulted — veto point):
  the calling tests invoke `cmd.Wait()` themselves later, and adding a watcher
  goroutine would reintroduce the double-`Wait` defect at a second site.
  Named cost: a genuinely crashed repair child now fails in 60s instead of
  2s — a 30x failure-latency regression on a path that has never fired,
  accepted to keep the diff free of new concurrency.
- The sibling repair waits (e.g. the 15s deadline in
  `binary_repair_install_test.go`) are untouched (defaulted — veto point):
  they have never flaked, and invariant 4 keeps the diff minimal.
- Production durable writes (`durableReplace`, `fh.sync()` in the repair
  script) are untouched — reviewer-closed 2026-07-22.
- No production code changes at all; the diff is confined to the two test
  files.

## Testing decisions

- The trigger condition (host-side VHDX I/O contention) cannot be produced
  from inside the guest on demand — guest-side fsync hammering went 0/6
  (`GATE-REPORT.md` evidence table) — so the fix cannot carry an automated
  red-to-green repro. Validation is the project gate (must stay green) plus a
  manual load-window check: up to 3 `bash bin/bench.sh gate` runs while the
  reviewer spins host-side load, matching the 1-in-3 hit rate observed.
- What a good test is here: these stories edit test harness code, so the
  "seam" is the harness itself and the gate is the only external observer.
  Inventing tests-of-tests for the helpers would add a layer the repo has no
  precedent for; the classification below is honest about that.
- Gate command: `bash bin/bench.sh gate`.

### Seam diagram

    trigger: contract test (R14 subtests / losing-racer subtest)
        │
        ▼
    spawn `bash bin/bench.sh <gate|repair>`  ──▶  [ marker-wait helper ]  ──▶  proceed on marker
                                                  [ 60s deadline       ]  ──▶  fast-fail on child exit (R14 only)
                                                  [                    ]  ──▶  diagnostic dump + Fatal on timeout
                      ◀ observed by the gate: suite green under stall, red with
                        self-documenting dump on genuine hang

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | R14 subtests pass while the child's pre-marker fsync stalls up to ~55s | `startStory5GateOwner` | not TDD-able — the stall needs host-side VHDX contention the guest cannot generate (0/6 guest-side attempts); validated by the manual load-window check in Testing decisions | Under the load window the old 5s deadline hit 1-in-3; three green runs under the same load is the agreed acceptance evidence |
| 1 | A child that exits before the marker fails the helper immediately with the exit error and captured output, not after the 60s deadline | `startStory5GateOwner` failure path | not TDD-able — the helper reports via `t.Fatal` and the repo has no test-of-test-helper seam to intercept it; graded by `/bench-review-implementation` against this row | Review confirms the in-loop exit-channel check precedes the deadline exhaustion; without it a dead child silently burns the full 60s |
| 2 | An R14 marker-timeout failure prints child output, per-thread wchan table, and the Go goroutine dump, with no `[DEBUG-a4f2]` tag anywhere in the tree | `startStory5GateOwner` failure path | `rg 'DEBUG-a4f2'` over the tree returns nothing (red today), and the gate stays green with the diagnostics compiled in; the sweep cannot distinguish promotion from deletion, so `/bench-review-implementation` additionally grades that the capture, dump, and SIGQUIT code survives untagged | The sweep proves the tags left; review proves the diagnostics stayed — together they exclude the cheapest wrong implementation (delete everything tagged) |
| 3 | The losing-racer subtest passes while the repair child's `fh.sync()` stalls up to ~55s | `waitForRepairMarker` | not TDD-able — same untriggerable host-side condition; covered by the same manual load-window check | The 2s deadline missed under the load window on 2026-07-22; surviving the same window is the acceptance evidence |
| 4 | A losing-racer marker timeout includes the child's captured output in the failure message | losing-racer test body | not TDD-able — the failure branch fires only under the untriggerable condition; verified by code review that the capture buffer feeds the Fatal | Without capture the next repair-side miss is blind exactly as the pre-instrumentation R14 miss was; review confirms the buffer is wired into the message |

### Edge inventory

- Child exits before marker (crash, bad env) — story 1 fast-fail; R14 only.
- Child alive past 60s (genuine hang) — diagnostic dump + SIGQUIT + Fatal,
  story 2.
- Child ignores SIGQUIT (non-Go child, masked signal) — bounded by the
  existing 3s reap timeout, then hard kill; already in the promoted code.
- Marker written, then child killed by a racing subtest — marker-first check
  ordering, Implementation decisions.
- **Won't handle:** sibling marker waits that have never flaked (15s install
  test and others) — minimal diff; raise them if they ever fire.
- **Won't handle:** making the trigger reproducible in CI (host-side I/O
  injection) — not possible from the guest; recorded in `GATE-REPORT.md`.

## Out of scope

- **Second flake** (`.bench-contract-env` TempDir RemoveAll) — separate defect
  with its own mechanism (child still writing at cleanup), never reproduced
  this session — 2 edits, 2 gate runs.
- **Orphaned-gate / stale-lock product defect** (killed wrapper leaves canary
  spawning and the flock held) — real product bug, parked in `IDEAS.md` —
  4 edits, 3 gate runs.
- **GATE-REPORT.md retirement** — after the fix lands and survives a load
  window, fold the durable findings into the profile/learnings and delete the
  report — 1 edit, 0 gate runs.

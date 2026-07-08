# concurrent-acquire fail-fast

Status: implemented

Compiled from `decisions/canary-cost.md` (FT40). The map's Handoff is closed with
no uncertainty flags; the contract-test edit was approved as a production change
in map ticket #6.

## Problem

The gate's wall is 130s, and 128s of it is the canary sweep. The sweep is not
computing — it is sleeping. The concurrent-acquire runtime contract test waits
for two spawned `bench worktree` runs to record their worktrees, on a fixed
60-second deadline, and never checks whether the runs are still alive. Against
~18 of the 66 canary fixture trees, both children exit within milliseconds, no
record can ever appear, and the test sleeps out the full minute anyway. Eighteen
60-second sleeps across 16 workers form two waves ≈ 128s of wall.

## Solution

Teach the wait loop the fact it already has access to: a `bench worktree` run
that has exited can never record. When both spawned runs have exited while the
record file holds fewer than two lines, fail immediately with a message distinct
from the timeout message. The 60-second deadline stays as the backstop for
live-but-stuck runs, and the passing path is unchanged. Prototype-measured: the
canary sweep drops from ~128s to ~26s.

## User stories

1. As a kit maintainer running `bench gate`, I want the concurrent-acquire wait
   loop to fail immediately once both spawned `bench worktree` runs have exited
   with fewer than two records, so that canary fixture trees whose children die
   at startup stop costing a 60-second sleep each.
   Line: claude-sonnet-5 / medium. The mechanism, numbers, and terminal
   condition are settled by the map's prototype, so this is transcription at a
   fully known seam, at the medium effort the profile assigns to gate and
   conformance logic.
2. As a reviewer reading a red contract phase, I want the early-exit failure
   message distinct from the existing timeout message, so that "children died
   before overlapping" is distinguishable from "children alive but stuck".
   Line: claude-sonnet-5 / medium. The message is a one-line authoring choice
   inside the same settled edit, so it shares story 1's routing.
3. As a kit maintainer, I want the passing path behaviorally unchanged, so that
   no canary EXPECT, gate verdict, or live-subject behavior moves.
   Line: claude-sonnet-5 / medium. This is a constraint the gate itself
   observes on every run, which is exactly the work the routing table sends
   cheap.
4. As a kit maintainer, I want the 60-second deadline retained as the backstop,
   so that a child that records once and then hangs still ends the test loudly
   instead of hanging the contract phase.
   Line: claude-sonnet-5 / medium. Retaining an existing branch is mechanical
   and carries the same settled-seam rationale as story 1.

## Implementation decisions

- One unit changes: the overlap wait loop inside
  `bench_worktree_concurrent-acquire_contract` (runtime contract suite).
  Outside and unchanged: the barrier design (record file → go-file release),
  the two-run spawn, the post-barrier assertions, the canary sweep, all
  production code.
- New terminal condition: **both spawned runs have exited AND the record file
  holds fewer than two lines → immediate fatal**, with a failure message
  distinct from the existing "did not record within a minute" message. The
  condition is deliberately "both exited": if one child crashes while the other
  holds, the record stays at one line and the 60s backstop still catches it —
  that rare slow case is accepted (map watch-out).
- How the loop observes "both exited": non-blocking drain of the existing
  `done` channel each iteration (the Handoff names this as the natural shape;
  it is the spec's choice here). Two ordering constraints make it correct:
  - **Drained results are not lost.** On the passing path, every result the
    wait loop consumes from `done` must still get its exit-code assertion —
    the final receive loop today expects exactly two results, so the
    implementation retains drained results (or counts receives) rather than
    silently eating them.
  - **Re-read the record after observing the second exit, before failing.**
    The second record can land between the loop's record read and its exit
    observation. Since an exited run can never record later, one final read
    after "both exited" is sufficient to close the race: fatal only if the
    count is still under two.
- On the terminal branch, releasing the go-file before failing is unnecessary
  (no straggler is alive) but harmless; keeping it symmetric with the timeout
  branch is acceptable either way.
- The fix is correct independent of the canary: a run that has exited can never
  record, because the fixture shell records before it holds.

## Testing decisions

- A good test here exercises the test's own observable interface — its verdict
  and failure message — not its internals. The unit under change *is* a
  contract test, so the gate itself is the primary harness: the outer contract
  phase runs the passing path on every gate, and the canary sweep drives the
  new terminal branch ~18 times per sweep (broken fixture trees make both
  children exit early).
- No new test is written for the test. The assertables are: gate green
  (passing path unchanged), canary sweep green with all 66 fixtures biting and
  the vacuity baseline clean, and the sweep wall dropping ~128s → ~26s. The
  wall drop has no gate assertion — verify it manually once after the change
  lands (`bench canary <root>` timed, per the Handoff).
- No fixture EXPECT matches either the old or the new failure message, so no
  canary semantics move with the message change.
- Prior art: the runtime worktree contract tests themselves (the FT37 comment
  atop this test records the barrier's design rationale — the edit must not
  reintroduce a timed-poll overlap assumption).
- Gate command: `.bench/gate.sh` (the project gate).
- Runtime note: rebuild `dist/bench` before hand-running the runtime contract
  suite — it exercises the built binary.

### Seam diagram

    trigger: outer contract phase (every gate run, passing path)
             + ~18 canary fixture inner gates per sweep (terminal path)
        │
        ▼
    record file (worktree paths)   ──▶  [ overlap wait loop in the        ]  ──▶  break to barrier release (pass)
    done channel (two run exits)   ──▶  [ concurrent-acquire contract test ]  ──▶  immediate fatal, distinct message
    60s deadline clock             ──▶  [                                  ]  ──▶  timeout fatal (backstop)
                      ◀ tests attach here: the gate runs the test; a broken
                        fixture tree drives the terminal branch, a live tree
                        drives the passing branch

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | both runs exited with <2 records → immediate fatal | the test run against a canary fixture tree (inner contract suite) | manual red: time the inner contract suite in a behavior-owned fixture tree pre- and post-change — pre fails after a 60s sleep with the timeout message, post fails in milliseconds with the new message; not TDD-able as a conventional red because the unit is itself a test | a missing or dead terminal branch leaves the failure at 60s with the old message, which the timing check sees directly |
| 2 | early-exit message distinct from timeout message | failing test output on the terminal branch | already covered — the story-1 timing check also inspects the message text | a reused timeout message makes the two failure modes indistinguishable in the observed output |
| 3 | passing path unchanged: gate green, canary green, all 66 fixtures bite, vacuity baseline clean | `.bench/gate.sh` (outer contract phase + canary sweep) | already covered — the gate runs the passing path and the full sweep on every run; an over-eager terminal check or an eaten `done` result turns the contract phase red | any behavior change on a live subject is exactly what the gate's existing green asserts against |
| 4 | live-but-stuck run still ends at the 60s backstop | the wait loop's deadline branch | not TDD-able cheaply — forcing it needs a child that records once and then hangs; the branch is retained code, and the map accepts the one-crashed-one-holding case landing here | removing the backstop would hang the phase on a stuck child, and retention is verifiable by reading the loop in review |

### Edge inventory

- **Exit/record race** (second record lands between the record read and the exit
  observation) — resolved as an implementation decision: one final record
  re-read after observing both exits, before failing; guarded ambiently by the
  passing path the gate runs every time (a botched ordering flakes the outer
  contract phase red under load).
- **One child exits, one holds** — **Won't handle** beyond the 60s backstop:
  the record stays at one line, the deadline branch catches it, and the map
  accepts that rare slow case.
- **Child exits nonzero vs zero** — the terminal condition keys on "exited",
  not the exit code; the post-barrier `RequireExit(0)` assertions still grade
  codes on the passing path. No distinct handling.
- **Absent vs empty record file** — unchanged behavior: the existing read
  tolerates both and counts zero lines; the edit adds no new file input.
- **Hostile inputs** — n/a per the Handoff: test-only edit consuming only its
  own children's exit signals and its own record file; the profile checklist
  classes attach to CLI surfaces this change does not touch.
- **Amputated caller** — none: no caller changes; the test's only consumer is
  the contract suite runner.

## Out of scope

- **Scope each canary fixture's inner gate to its EXPECT's phase** — a separate
  capability (a fixture phase-declaration mechanism plus a per-phase vacuity
  baseline), rejected under the lenient budget in map ticket #2 and kept as the
  documented first resort if the gate wall creeps — ~15 edits, ~10 gate runs.
- **FT39 deadline fix** (event-keyed or raised deadline for the outer test's
  legit overlap window) — a separate capability parked on the roadmap pending a
  captured red with the exact message from a real `bench gate` run (map ticket
  #6b) — ~3 edits, ~3 gate runs once the repro exists.
- **Skip-canary caching and fixture batching** — the remaining #3 candidates,
  rejected with the lenient budget; they stay on the map/roadmap menu, not
  here — ~10 edits, ~8 gate runs each.

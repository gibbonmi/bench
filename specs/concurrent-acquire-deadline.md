# concurrent-acquire deadline

Status: draft

Compiled from the FT39 roadmap graduation (captured red: a real `bench gate` run
went red with exactly "second acquire did not record within a minute — the runs
never overlapped" at `runtime_worktree_test.go:355`, 60.06s, on a tree identical
to a green-gated one). There is no `decisions/` map behind this row — the repro
was the graduation criterion — so this spec is reviewer-directed. The one
contestable call it makes is **event-keyed vs. a blind raise** of the deadline;
it is decided here (event-keyed) and flagged for veto in Implementation
decisions. Sibling context: `specs/concurrent-acquire-fail-fast.md` (retired,
`bench spec history concurrent-acquire-fail-fast`) and ADR 0003.

## Problem

The concurrent-acquire runtime contract test spawns two `bench worktree` runs and
proves they genuinely overlap: each shell records its worktree path and then holds
on a barrier the test releases only after seeing both records. To bound the wait,
the test arms a fixed 60-second deadline **at spawn** and fails if fewer than two
records land before it. That deadline therefore times the whole spawn → first
record → second record path, and the spawn portion is unbounded under load:
during a full `bench gate`, the concurrent canary sweep spawns ~18 copies of this
same poll storm alongside the outer contract phase's own run, so getting both
child processes launched and holding can itself approach the minute. The captured
red is exactly this — 60.06s elapsed, behavior correct, the runs would have
overlapped, but the wall-clock budget expired first. The `ea966a3` fail-fast
already turned "both children exited early" into an instant failure; the remaining
weak point is the fixed spawn-relative deadline itself.

## Solution

Stop timing the overlap assertion from spawn. Split the wait into the two windows
it always contained and time each against the event that opens it:

- **First-record window** (spawn → first record): a generous hang-guard only. It
  absorbs the load-driven spawn latency that caused the flake and bounds a truly
  wedged acquire; it is not a behavioral assertion, so it can be set well above
  any plausible spawn spike without dulling any bite.
- **Overlap window** (first record → second record): armed the moment the first
  record appears — i.e. once one run is confirmed acquired and holding. This is
  the window that carries the test's bite: under correct behavior a second acquire
  of a distinct worktree lands in well under a second; under a serialization
  regression the second run blocks on the first's lease forever, so the window
  expires and the test goes red. Because it is measured from a real event, spawn
  latency never leaks into it.

The `ea966a3` fail-fast (both runs exited with fewer than two records → instant
fatal) is retained unchanged, and the barrier design (record-then-hold, test drops
the go-file after two records) is untouched. The passing path stays behaviorally
identical; the canary-driven fail-fast path stays millisecond-fast.

## User stories

1. As a kit maintainer running `bench gate` under load, I want the overlap
   deadline armed at the first record rather than at spawn, so that unbounded
   spawn latency during the concurrent canary sweep no longer counts against the
   overlap assertion and the correct-behavior flake disappears.
   Line: claude-sonnet-5 / medium. The event-keyed shape is fully decided in this
   spec and the edit lands at one known seam, which is the
   transcription-at-a-known-seam work the profile routes to the cheap model at
   the medium effort it assigns to gate and conformance logic.
2. As a kit maintainer, I want a genuine serialization regression to still go red,
   so that removing the flake does not remove the test's reason to exist: when one
   run records and holds while the second blocks on its lease forever, the
   first-record-armed overlap window must expire and fail.
   Line: claude-sonnet-5 / medium. Retaining and re-anchoring an existing deadline
   branch is mechanical once the arming event is specified, and it carries story
   1's settled-seam rationale.
3. As a kit maintainer, I want a generous first-record backstop, so that two runs
   that are alive but never record (a wedged acquire that never exits, which the
   fail-fast cannot catch) still end the contract phase loudly instead of hanging
   it.
   Line: claude-sonnet-5 / medium. This is one added guarded branch inside the
   same settled edit, sharing story 1's routing.
4. As a reviewer reading a red contract phase, I want each of the three terminal
   conditions to carry a distinct message, so that "children died before
   overlapping", "one holds but the other never overlapped", and "nothing recorded
   at all" are told apart from the failure text alone.
   Line: claude-sonnet-5 / medium. The messages are one-line authoring choices
   inside the same edit and move no canary EXPECT, so they share story 1's
   routing.
5. As a kit maintainer, I want the passing path and the fail-fast path
   behaviorally unchanged, so that no gate verdict, canary EXPECT, or live-subject
   behavior moves and the fail-fast still fires in milliseconds against broken
   fixture trees.
   Line: claude-sonnet-5 / medium. This is a constraint the gate itself observes
   on every run and on every canary sweep, which is exactly the work the routing
   table sends cheap.

## Implementation decisions

- **One unit changes:** the overlap wait loop inside
  `bench_worktree_concurrent-acquire_contract` (runtime contract suite). Outside
  and unchanged: the barrier (record file → go-file release), the two-run spawn,
  the `ea966a3` fail-fast branch, the post-barrier assertions
  (`RequireExit(0)`, distinct-path check), the canary sweep, all production code.
  The FT37 comment atop the test warns against a timed-poll *overlap mechanism*;
  this edit must not reintroduce one — overlap stays guaranteed by the barrier,
  and the deadlines remain only failure backstops.

- **Contestable call, decided here and flagged for veto — event-keyed over a
  blind raise.** A blind raise (e.g. 60s → 3m at spawn) is one line, but it is a
  guess that can still flake on a heavier machine and makes every genuine
  serialization regression take the full raised wall to fail. Event-keying instead
  bounds the behavior the test actually asserts (second acquire, measured from a
  confirmed-holding run), so it is robust to spawn load at any magnitude and keeps
  the real-regression fail-latency modest. Recommend event-keyed. **Reviewer veto
  surface:** if you would rather ship the one-line raise, say so — this is your
  call, and the two chosen durations below are equally yours to tune.

- **Two windows, each event-anchored:**
  - *Overlap window* — armed when the record file first reaches ≥1 non-empty line;
    on expiry, fatal. **Recommended duration: 60s from the first record.** Measured
    from a run already confirmed holding, a correct second acquire lands in well
    under a second even under gate load, so 60s is ~100× headroom yet still finite,
    so a serialization deadlock fails red. Keeping the familiar 60s magnitude
    minimizes surprise while changing only *what the clock is measured from*.
  - *First-record backstop* — a generous guard bounding spawn → first record; on
    expiry, fatal. **Recommended duration: 2m.** The captured flake shows the
    spawn-to-first-record path can approach a minute under the concurrent canary
    storm; 2m clears any plausible load spike while still bounding a wedged
    acquire. It never fires once any record appears, so generosity here costs
    nothing on the passing path and dulls no bite (the overlap assertion is
    event-relative, not covered by this guard).

- **The bite is intact because the first record always lands in the regression
  case.** In a serialization bug one run wins the pool, acquires, records, and
  holds while the other blocks on the lease — so the first record appears, the
  overlap window arms, and its expiry converts what is otherwise an infinite hang
  into a red. The rarer sub-case where *neither* run records (both block acquiring)
  is caught by the first-record backstop. Dropping the overlap deadline entirely —
  the cheapest wrong "fix" that leans only on the fail-fast — would turn the
  first-records-then-second-hangs regression into a gate hang, not a red; the
  overlap deadline must be retained.

- **Fail-fast retained and checked every iteration.** The `ea966a3` branch (both
  runs exited with <2 records, after one final re-read → fatal) stays and must be
  evaluated on every loop pass regardless of window, so the canary path still fails
  in milliseconds.

- **Exit/record race, unchanged handling.** A second record can land between the
  loop's record read and its observation of an exit; retain the existing final
  re-read before any fatal, and break to the barrier release whenever the count
  reaches two. An exited run can never record later, so one re-read closes the
  race.

- **Natural implementation shape (freedom):** two `time.Time` deadlines — the
  backstop set at spawn, the overlap deadline set lazily on the transition to
  ≥1 record — evaluated each iteration alongside the existing non-blocking `done`
  drain. Any equivalent structure (a single deadline re-anchored on first record)
  is acceptable provided drained `done` results still receive their exit-code
  assertions and the passing path stays byte-identical.

## Testing decisions

- A good test here exercises the unit's observable interface — its verdict and its
  failure message — not its internals. The unit under change **is itself a
  contract test**, so the gate is the primary harness: the outer contract phase
  runs the passing path on every gate, and the canary sweep drives the fail-fast
  branch ~18× per sweep (broken fixture trees make both children exit early). No
  new test is written *for* the test, matching the sibling fail-fast spec's
  posture.
- The flake removal (story 1) is **not conventionally TDD-able**: a load-timing
  flake has no deterministic failing test. Its red signal is the already-captured
  FT39 red (60.06s, the timeout message, on a tree identical to green); the fix is
  verified by (a) review reading that the overlap deadline is armed at the first
  record, and (b) an empirical loaded-gate soak (see the coverage map) confirming
  the flake message no longer appears.
- The bite retention (story 2) is verifiable by an **induced serialization probe**
  where feasible — forcing single-worktree contention so the second run blocks on
  the first's lease and the overlap window fires with the new message — otherwise
  by review that the overlap deadline is retained and armed on first record. The
  probe requires a pool-size / contention knob whose availability is unconfirmed;
  building a permanent deterministic bite test out of it is a separate capability
  (Out of scope), consistent with how the sibling spec verified its own
  test-of-a-test.
- No canary fixture EXPECT matches the old or the new failure strings (verified:
  the only nearby match is a different test's message in
  `internal/canary/canary_concurrency_test.go`), so the message changes move no
  canary semantics.
- Gate command: `bench gate` (the project gate, `.bench/gate.sh`).
- Runtime note: rebuild `dist/bench` before hand-running the runtime contract
  suite — it exercises the built binary.

### Seam diagram

    trigger: outer contract phase (every gate run — passing path)
             + ~18 canary fixture inner gates per sweep (fail-fast path)
        │
        ▼
    record file (1st line = event) ──▶ [ overlap wait loop in the        ] ──▶ ≥2 records → break, drop go-file (PASS)
    done channel (two run exits)   ──▶ [ concurrent-acquire contract test ] ──▶ both exited, <2 → fatal (fail-fast, ea966a3)
    overlap clock (armed @ 1st rec)──▶ [                                  ] ──▶ overlap window expired → fatal (serialization)
    first-record backstop clock    ──▶ [                                  ] ──▶ no record at all → fatal (wedged acquire)
                      ◀ tests attach here: the gate runs the test; a live tree
                        drives the passing branch, a broken tree drives fail-fast,
                        an induced single-worktree contention drives the overlap branch

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | overlap deadline armed at the first record, not at spawn → correct-behavior runs no longer time out under load | the test's passing path in the outer contract phase, under concurrent canary load | already captured (FT39 graduation): a real `bench gate` went red with `second acquire did not record within a minute — the runs never overlapped` at 60.06s on a tree identical to green. Not further TDD-able (load-timing); post-fix verified by a loaded-gate soak (repeat `bench gate` / `bench canary <root>` N× under load) showing no timeout-message red | if the deadline stayed spawn-relative, the same load that produced the captured red would keep tripping it; anchoring it to the first record removes spawn latency from the timed window, and the soak observes the message's absence |
| 2 | genuine serialization (first records and holds, second blocks on the lease forever) still goes red | the overlap-window branch | induced red where feasible: force single-worktree contention so the second run cannot acquire; the loop fires `second acquire did not record within 60s of the first — the runs never overlapped`. Else review-verified: the overlap deadline is retained and armed on first record | a "fix" that dropped the overlap deadline would turn this regression into an infinite gate hang; retaining and arming it on the first record converts the hang into an observable red |
| 3 | both runs alive but never record → loud failure, not a hang | the first-record backstop branch | not cheaply TDD-able — forcing a live-but-non-recording acquire needs a wedged pool; review-verified that the backstop fires `no bench worktree run recorded within 2m — acquire appears wedged` | without a first-record guard no overlap deadline ever arms (it arms on first record) and the fail-fast needs both runs exited, so this case would hang the contract phase forever |
| 4 | three terminal conditions carry distinct messages | the failing-test output on each branch | already covered by rows 1–3, which each key on a different verbatim string; the strings must be pairwise distinct | a reused message makes two failure modes indistinguishable from the observed red, defeating the reviewer's diagnosis |
| 5 | passing path and fail-fast path behaviorally unchanged: gate green, canary green with all fixtures biting and the vacuity baseline clean, fail-fast fires in ms | `bench gate` (outer contract phase + canary sweep) | already covered — the gate runs the passing path and the full sweep every run; an eaten `done` result, an over-eager backstop, or a slowed fail-fast turns the contract phase or the sweep red | any behavior change on a live subject or any regression in fail-fast latency is exactly what the gate's existing green (and ADR 0003's ~26s canary wall) asserts against |

### Edge inventory

Walked per behavior; each lands as a row above or a **Won't handle** line here.

- **Loaded machine / slow spawn** — row 1: the overlap window is event-keyed and the
  first-record backstop (2m) absorbs the spawn spike; neither trips on correct
  behavior under load.
- **Both runs exit early** — row 5 / unchanged: the `ea966a3` fail-fast fires in
  milliseconds; this edit must not slow it.
- **One run exits, one holds and recorded** — record count sits at 1, the overlap
  window (armed at that first record) expires and fails red — the same accepted
  slow case as the fail-fast spec, now bounded by the overlap window from the first
  record rather than by the old spawn-relative minute.
- **Genuine non-overlap / serialization** — row 2: first record lands, overlap
  window arms and fires; the bite is preserved.
- **Neither run records (both block acquiring)** — row 3: the first-record
  backstop fires.
- **Exit/record race** (second record lands between the read and the exit
  observation) — resolved as an implementation decision: one final re-read before
  any fatal, break to release on count ≥ 2; guarded ambiently by the passing path
  the gate runs every time.
- **One child crashes while the other holds, landing at the overlap deadline
  instead of instantly** — **Won't handle** beyond the overlap window: the record
  stays at one line, the overlap branch catches it after its window, and that rare
  slow case is accepted (inherited from the fail-fast spec's watch-out).
- **Distinguishing a slow-but-correct second acquire from serialization inside the
  overlap window** — **Won't handle**: not attempted; the 60s-from-first-record
  headroom makes a false red on correct behavior implausible, and tightening it to
  probe the boundary would reintroduce flake risk for no bite.
- **Hostile inputs** — n/a: test-only edit consuming its own children's exit
  signals and its own record file; the profile checklist classes attach to CLI
  surfaces this change does not touch.
- **Amputated caller** — none: no caller changes; the test's only consumer is the
  contract suite runner.

## Out of scope

- **A permanent deterministic bite test that forces single-worktree serialization**
  — a separate capability requiring a pool-size / contention knob (availability
  unconfirmed) plus a dedicated fixture, distinct from re-anchoring the existing
  deadline; the bite is retained and review-verifiable without it. ~5 edits,
  ~4 gate runs.
- **A blind raise of the spawn-relative deadline** — the rejected alternative to
  this spec's event-keying, not a future capability; recorded here only so the
  reviewer's veto has a named target. ~1 edit, ~1 gate run.
- **Scope-each-canary-fixture-to-its-EXPECT-phase and the other canary latency
  cuts** — ADR 0003's pre-decided first resort if the sweep wall creeps; a separate
  capability on the canary-cost menu, unrelated to this test's deadline. ~15 edits,
  ~10 gate runs.

# Session handoff

Repository: /home/mgibs/workspace/bench — branch `main`. Everything below is
executable from a cold start; no conversation history is needed.

## State

- **FT91's first arm is built, gate-green, and landed on `main`.**
  `specs/gate-concurrency-budget.md` is `Status: implemented`; the code and the
  status flip ride in one commit (the one that last wrote this file). Nothing
  from this arm is uncommitted. `main` is unpushed — pushing is the reviewer's call.
- **What landed.** `bounds.CanaryInnerWidth = 2` is the one source of k.
  `internal/canary/canary.go` consumes it twice: `innerEnv` strips any inherited
  `GOMAXPROCS=` and appends exactly one `GOMAXPROCS=2`, and
  `fixtureWorkers(runtime.GOMAXPROCS(0), len(fixtures))` derives the sweep's
  worker pool — divide by k, floor at 1, cap at the fixture count.
  `checkBoundsPolicy` carries all three registrations (`required`, `owners`,
  `boundLikeName`), each proved by a canary mutation fixture. `runtime.NumCPU()`
  is gone repo-wide. `internal/gate` is untouched.
- **Story 1's payoff, measured 2026-07-24.** Full gate: **336 s wall, peak
  1-minute load 33** on 16 cores, against the 2026-07-22 baseline of 10–15 min at
  load ~123. The map's k = 2 prediction was 332 s — reproduced within ~1%. The
  conformance phase (335 s) is now the long pole, exactly as the map said.
- **Two closed veto points, not open questions:** k = 2, and the outer
  conformance/contract phases stay uncapped. Don't reopen them.
- **Story 7 nested-run contract.** The canary package's own tests run at
  `GOMAXPROCS=2` inside a fixture's inner gate, where the derived bound is 1. The
  two tests needing overlap gate through `capability.CPU`; the bounds test
  releases and asserts at the derived bound. If you touch
  `internal/canary/canary_concurrency_test.go`, re-prove all three:
  `GOMAXPROCS=2 go test -timeout 120s ./internal/canary` green, `-v` showing
  exactly two `bench-skip kind=capability class=cpu` lines, and the same run at
  full width showing none. A deleted assertion and an honest skip both look
  green; only the emitted line tells them apart.
- **The bounds test grades the cap's direction, and only with its settle.** After
  in-flight first reaches the derived bound it waits before releasing, so an
  over-wide pool oversubscribes and is counted. Without that settle an uncapped
  implementation passes. Verified: reverting `runFixtures` to `runtime.NumCPU()`
  fails it 5/5 with `high-water = 16, want == derived bound 8`.
- **A concurrent session was writing this repo on 2026-07-24** and edited this
  file mid-gate, which the working-tree tripwire correctly caught. Its guidance
  paragraph is preserved under Shape below. Check `bench status` and `ps` for
  other live writers before starting a gate.
- **A leaked test gate process** (`TestFT78Story4ProofLedgerR12armed-stop-blocked`)
  was running for ~30 hours under `/tmp` as of 2026-07-24, perturbing load
  measurements. Not investigated; kill it or diagnose the leak before trusting a
  timing run.

## Next command

`/bench-what-next` in a fresh mid-tier session. `bench status` has drain work
pending: parked ideas plus open learnings, one roadmap row naming a retired spec,
and this arm's own spec now eligible for `bench spec retire
gate-concurrency-budget` once it reaches the default branch. FT91's remaining arms
— capping the outer phases, removing the hardcoded `-count=1`, and a shared
hermetic build cache — are recorded under Out of scope in the spec and are
separate capabilities, not follow-ups to this one.

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Operational gotchas are placed by lifetime, not copied here: one that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes, and one scoped to a
build belongs in that spec's coverage rows. This file names at most when you'll hit
one, never the command — a second copy drifts from the source.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.

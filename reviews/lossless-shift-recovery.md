# Review findings — lossless-shift-recovery (FT79)

Semantic review of the uncommitted shape/spec session: `specs/lossless-shift-recovery.md`,
`decisions/lossless-shift-recovery.md`, and the CHANGELOG/ROADMAP reconcile.
Advisory; the gate stays the oracle. Three findings converge on one gap: the
`interrupted`/exit-130 outcome has no acceptance-coverage row (S1, Spec 1, Coverage F4).

## Standards

**5 findings (1 hard violation, 4 judgment calls). Worst: S1.**

- **S1 (hard) — `interrupted`/130 has no coverage-map row.** `bench-craft-spec`:
  "Assertables become rows"; every edge lands in a row or a Won't-handle line. The
  recovery matrix (decision map) treats signal interruption as a first-class row with
  *new* snapshot-on-interrupt behavior, but `specs/lossless-shift-recovery.md` carries
  it only as edge-inventory prose ("the existing interrupt-cleanup contracts extend…"),
  which is neither a row nor a cut.
- **S2 (judgment) — spec glosses the opaque cheap-tier alias as "Sonnet 5"**
  (`specs/lossless-shift-recovery.md` line 33). Invariant 2 keeps tier tokens opaque;
  `.bench/lines.env` names no model version, and ROADMAP FT8 frames "Sonnet 5" as a
  *mid*-tier candidate, so the gloss is likely wrong as well as over-specified.
- **S3 (judgment) — CHANGELOG entry records a no-op learnings run** the file's own
  append rule ("one entry per `/bench-update-kit` run or learnings-sourced promotion")
  doesn't cover. Matches existing local practice (the 2026-07-13 entries), so
  judgment-call, not drift introduced here.
- **S4 (judgment, drift-risk only) — map and spec restate the exit-code taxonomy,
  env defaults/bounds, and recovery matrix.** No divergence today; inherent to the
  map→spec pipeline. Note, not a fix.
- **S5 (judgment, cosmetic) — decision-map title carries "(FT79)"** where sibling maps
  and the spec omit the tag.

## Spec

**2 findings. Worst: Finding 1.**

- **Spec 1 — missing row for the interrupted/exit-130 outcome.** Decision-map matrix
  row "signal interruption → interrupted 130", Handoff gate-attachment lists "SIGTERM
  mid-adapter", story 6 names 130 — the other five taxonomy codes each have rows; 130
  has none and no falsifiable red signal. Close by adding a Seam A row: SIGTERM
  mid-adapter with a mutation → exit 130, `shift_result` outcome `interrupted`,
  recovery ref resolves.
- **Spec 2 — decided refactor-phase discard carve-out omitted, and contradicted.**
  Decision #9: "The refactor phase's red-gate *rollback* stays discard-by-design."
  The spec never encodes the exception and instead asserts "one uniform rule for
  every post-mutation failure" (Release-ordering decision), which directs a builder
  to snapshot exactly the case #9 says must stay discarded.

All named repo surfaces verified real (`envInt`, `stageTouched`, `runAdapter`,
`temporaryIndex`, `commitTree`, `Fault`, `ReasonUnexpectedLock`, `resume-clean`,
`runtime_shift_test.go`); CHANGELOG/ROADMAP tree-facts hold (IDEAS.md empty, zero
open journal entries, decision map complete).

## Coverage

**8 findings (6 primary, 2 medium). Worst: F1.**

- **F1 — same-second branch/ref collision defeats story 18.** Branch name is
  second-resolution (`internal/shift/loop.go:111`); two shifts in the same second
  derive the same branch and recovery ref, so the second dies at `git switch -c` or
  is rejected by the fail-closed `update-ref` — losing the second shift's work, the
  exact failure FT79 exists to prevent. Row 3,18 runs shifts sequentially and never
  forces the same-second case. Fix granularity (suffix vs finer timestamp) is an
  unmade decision — reviewer's call.
- **F2 — `BENCH_MAX_WALL` validation under-enumerated.** Row 11's only wall value is
  `48h`; malformed (`abc`), zero/negative (`0s`, `-5m` — Go parses negatives), and
  the inclusive `24h` accept are unexercised, so a degenerate parser passes.
- **F3 — cap upper bound and `BENCH_REFACTOR_ITERS` never concretely exercised.**
  Row 11 supplies no `=101` over-bound value and no invalid `_REFACTOR_ITERS` value;
  an impl missing either check passes the row.
- **F4 — interrupt-with-mutation unexercised.** `testShiftInterruptCleanup`'s agent
  writes nothing and asserts cleanup, so a build that skips snapshot-on-interrupt
  passes. Needs a row: SIGTERM after a write → exit 130 + recovery ref with the
  mutated file. (Same gap as S1/Spec 1.)
- **F5 — wall deadline firing mid-gate unexercised.** Row 8 sleeps only the adapter;
  a wall that never wires cancel-gate passes yet hangs on a slow gate. Needs a
  slow-gate + tiny `BENCH_MAX_WALL` row asserting cancellation, snapshot, exit 3.
- **F6 — adapter spawn failure (vs nonzero exit) unpinned.** `runAdapter` swallows
  the start error today (`internal/shift/loop.go:290`), reading spawn failure as a
  clean no-op; row 6,7 only exercises `/bin/false`. Needs a spawn-failure case
  asserting `failed`/exit 1, not `no-op`.
- **F7 (medium) — whitespace-only objective slips past `objective == ""`**
  (`internal/shift/shift.go:137`). Row 9 tests only the absent objective; add a
  `bench shift "   "` row asserting exit 2, no lease — or decide it's acceptable.
- **F8 (medium) — dirty nested repo/submodule in the snapshot loses inner work.**
  `add -A` records a gitlink only; uncommitted work inside the nested repo is not
  captured, and the profile checklist names "dirty nested repositories". Row 1,18
  asserts nested *untracked* paths only. Needs a row or an explicit Won't-handle
  line — an edge to *decide*, not drop.

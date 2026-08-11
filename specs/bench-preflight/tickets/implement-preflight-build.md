# Implement bench preflight build mode

Blocked by: implement-preflight-review.md, harden-preflight-bootstrap-errors.md
Ownership fence: `internal/preflight/`
Integration surfaces: verdict core + gatherer + grammar→implement-preflight-review.md (expanded in place inside `internal/preflight/`, exercised by every B row); hardened bootstrap diagnostics→harden-preflight-bootstrap-errors.md; build-mode verb consumed by phase prose→advertise-preflight-kit-prose.md
Contracts: the `build` mode verb and its applicability table (modes: `review`, `build`; per-check states: run or `not-applicable`; exit derives from run rows only) cross `internal/preflight/`→the phase-entry prose in `.agents/commands/bench-implement-spec.md`, asserted by advertise-preflight-kit-prose.md's D2 against this real grammar
Closure: B1/rows-owned-na, B1/rows-membership-na, B1/diff-nonempty-na, B1/green-exit, B2/present-tickets-real-checks, B2/empty-tickets-red, B3/red-exit-despite-na, B4/build-out-of-fence-red

## What to build

Expand the grammar to accept `bench preflight build <slug>` and give the
verdict core mode applicability: build always runs `base-current` and
`paths-authorized`; `rows-owned` and `rows-membership` run for real when
`specs/<slug>/tickets/` exists (present-but-empty is red — declared rows
unowned) and print `not-applicable` when it does not; `diff-nonempty` is
`not-applicable` in build mode. Not-applicable rows are printed, definitive
verdict rows, and never soften a real red's exit. Same TOON table, exit, and
error contracts as review mode; same CLI-contract and verdict-core test seams.

## Acceptance

- [ ] [B1] (covers PF9) `bench preflight build <slug>` with no `tickets/` directory prints a `not-applicable` row for each of `rows-owned`, `rows-membership`, and `diff-nonempty` — all three individually asserted — runs the rest, and exits 0 when they are green.
- [ ] [B2] (covers PF10) with a present `tickets/` directory the row checks run for real; present-but-empty `tickets/` is red — declared rows unowned.
- [ ] [B3] (covers PF11) `base-current` red in build mode exits 1 even while the ticket checks are not-applicable.
- [ ] [B4] (covers PF21) a tracked change outside every fence entry makes `paths-authorized` red in build mode, exit 1.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| B1/rows-owned-na | omit the `rows-owned` not-applicable row from the table | the fresh-build contract test | seed no tickets dir, run build, expect the missing-row failure |
| B1/rows-membership-na | omit the `rows-membership` not-applicable row | the fresh-build contract test | same seed, expect the missing-row failure |
| B1/diff-nonempty-na | run `diff-nonempty` for real in build mode instead of printing not-applicable | the fresh-build contract test | same seed (empty diff), expect the false-red failure |
| B1/green-exit | exit 1 whenever any row is not-applicable | the fresh-build contract test | same seed, expect the exit-code failure |
| B2/present-tickets-real-checks | print not-applicable whenever tickets exist too | the resumed-build contract test | seed a tickets dir citing all rows, run build, expect the not-run failure |
| B2/empty-tickets-red | treat present-but-empty `tickets/` as not-applicable | the empty-tickets contract test | seed an empty tickets dir, run build, expect the missed-red failure |
| B3/red-exit-despite-na | derive the exit from the last row's verdict | the stale-base build contract test | seed a stale base with no tickets dir, run build, expect the exit-0 failure |
| B4/build-out-of-fence-red | skip `paths-authorized` in build mode | the build out-of-fence contract test | seed a change outside every fence, run build, expect the missed-red failure |

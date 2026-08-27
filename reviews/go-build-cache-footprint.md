# Review pickup: go-build-cache-footprint

Frozen base `caeb19fbe668617d01a7b65d550919fe9da27d88`, reviewed tip
`8f11c0571edf1a6faac15985f0369dedd10ace66`. Raw findings: Standards 5, Spec 3,
Coverage 4. Repair targets after collapse: 8, in tickets 09 to 11.

## Standards

Count 5. Worst: the operator-facing cache refusals are derived twice across
the ticket 04 and 05 fence.

- `auto-fix` — `internal/gocache/command.go:56` and `internal/gocache/clean.go:35`
  carry the byte-identical control-byte refusal, and `cache directory not derived`
  is derived at `clean.go:32`, `clean.go:53`, and `command.go:51`. Rule: AGENTS.md
  "one source per fact". Ticket 10.
- `auto-fix` — `internal/preprelease/preprelease_test.go` T07 comment narrates the
  deleted `goTestArgv`. Rule: `craft-comments` "no narration". Ticket 10.
- `auto-fix` (reviewer chose trim) — `internal/gate/cache_footprint_run_test.go`,
  `internal/gate/lane_test.go` T03, `internal/gate/report_test.go` R08 argue the
  test's design to the reviewer. Ticket 10.
- `no-op` — the "disk pressure is a machine fact" rationale appears in four places;
  prose, not enforcement.
- `auto-fix` — `internal/gocache/clean_test.go` `equal` restates `slices.Equal`.
  Ticket 10.

## Spec

Count 3. Worst: rows L01 to L03 are partial, because no test binds the three
production `gocache.Hold` call sites.

- `auto-fix` — L02 and L03 name a lane run and a `bench test` run, but
  `TestCleanRefusesWhileEachHolderRuns` calls `Hold` directly under three labels.
  Enumeration: `rg gocache.Hold` finds three production sites and no test outside
  `internal/gocache`. Ticket 09.
- `auto-fix` — the C12 reseam to `gate runner integration` was a reviewer decision
  on 2026-08-27 in the coordinating session, but ticket 07 does not record it.
  Ticket 10 records it.
- `no-op` — `bench cache clean` rides the parent's operational exemption because
  one row carries one disposition; flagged for veto, not repaired.

## Coverage

Count 4. Worst: the holder call sites are untested (collapses into Spec 1).

- `auto-fix` — same target as Spec 1. Ticket 09.
- `auto-fix` (reviewer chose unbounded wait) — a holder that cannot take the lock
  within `bounds.CacheHoldWait` runs unlocked. Holders now wait with `F_SETLKW`,
  and the bound goes away. Ticket 09.
- `auto-fix` (reviewer chose the fallback) — a phase runner whose env carries no
  `GOCACHE` prints an empty directory and logs an empty path. `FromEnv` falls back
  to the `HOME` derivation, pinned by new row R18. Ticket 11.
- `no-op` — an out-of-range integer in `trim.txt` renders a bogus timestamp;
  cosmetic, decided no-op.

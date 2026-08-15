# Characterize lock contention, timeout, and persistence failures at the public seam

Blocked by: characterize-gate-run-outcomes.md
Writes: internal/gate (new test files only)

## What to build

Tests only, extending the fixture the first ticket landed. Lock held by an
in-flight run refuses with the owner diagnostic (GC4); lock held by a separate
holder process — the test binary re-exec'd (prior art
`internal/freshness/freshness_test.go`) that opens `<git-dir>/bench-gate.lock`,
takes the fcntl write lock, writes an owner record, and waits — demotes a
retained reusable green to pending on a refused `RunCommand([--fresh, root])`
(GC5; ordinary `Execute` reuses before lock acquisition per GC1); with
`gateTimeout` overridden and restored, a `--fresh` run of a sleeping script
records `timeout` and the earlier green is not reused afterwards (GC6);
sentinel-switched `--fresh` runs whose script strips the evidence directory
mode (GT3), makes it unwritable before a red exit (GT4), or makes the git-dir
unwritable before a green exit (GT5) produce the operational exits and leave
the record states the rows name. GT3–GT5 skip under `internal/capability`'s
privilege check with the fixture restoring modes in cleanup. Each row lands
green today. Grouping: as in the first ticket, no per-row red is stranded; the
rows share the holder-process and permission-flip fixture capabilities.
Consumes the first ticket's fixture unchanged and extends it with the holder
process and the mode-flipping script variants.

Return note (not acceptance): for every row, the mutation applied and the
observed red.

## Acceptance

- [ ] In-flight run holding the lock: concurrent `Execute` `ActionExit 1`, stderr `gate execution already in progress` and `gate owner: pid <n> (alive)`, script not re-run (covers GC4)
- [ ] Retained green + external holder: `Inspect` `Ready`/`green` before, refused `RunCommand([--fresh, root])`, `Pending` after the refused call, script not run (covers GC5)
- [ ] Sleeping script under sub-second `gateTimeout` via `--fresh`: 124, `gate: timeout`, `Inspect` `Ready`/`timeout`; sentinel removed, ordinary run re-runs (counter 3) (covers GC6)
- [ ] Evidence dir stripped to 0500 before green exit: `ActionExit 1`, `gate evidence persistence failed`, not reusable (covers GT3)
- [ ] Evidence dir unwritable before red exit with a retained file present: `ActionExit 1`, `gate evidence invalidation failed`, not reusable (covers GT4)
- [ ] Git-dir unwritable before green exit: `ActionExit 1`, `gate final persistence failed`, `Inspect` `Pending` (covers GT5)

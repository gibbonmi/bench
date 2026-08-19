# Apply the reclaim plan

Blocked by: 01-plan-reclaimable-pool-keys.md
Writes: internal/worktree/pool_reclaim.go, internal/worktree/pool_reclaim_test.go, internal/usage/worktree.go

## What to build

`bench worktree reclaim --apply <fingerprint>`: the destructive half of the
handshake ticket 01 opened. It re-plans, refuses on stdout naming the re-plan
command and exits 1 when the recomputed fingerprint does not match the one
supplied, and otherwise removes each planned key after re-checking that key
against ticket 01's predicate immediately before removal — the fingerprint proves
the plan as a whole is current, the re-check proves each individual key still
qualifies at the instant of removal.

Removal is `os.RemoveAll` on a path whose parent is asserted to be exactly
`benchHome()/worktrees`. It prints per-key verdicts and removed/retained counts
in the same TOON shape the plan uses. An absent, empty, or malformed
`--apply` value is exit 2 with usage and removes nothing; an apply over a pool
with no reclaimable key exits 0 having removed nothing.

No receipt or transaction record is written: one `RemoveAll` is the whole
per-key operation, so an interrupted run leaves the keys it already removed
removed and the rest untouched, which the next plan simply reports as a smaller
target set.

## Acceptance

- [ ] an apply carrying the plan's fingerprint removes exactly the planned keys, leaves every other key present, and prints per-key verdicts with removed and retained counts (AP1, PL3).
- [ ] an apply whose fingerprint no longer matches the pool refuses on stdout naming the re-plan command, exits 1, and removes nothing (AP2).
- [ ] a key that becomes non-reclaimable between plan and apply — a live `.git` pointer appears in it — survives an apply whose fingerprint still matches (AP3).
- [ ] an apply over a pool with no reclaimable key exits 0 having removed nothing (AP4).
- [ ] an absent, empty, and malformed `--apply` value each exit 2 with usage and remove nothing (AP5).
- [ ] every removed path's parent is exactly `$BENCH_HOME/worktrees` (SH7).
- [ ] a key whose source repository was created, keyed, and then deleted is planned and then removed by the real command (RP1).

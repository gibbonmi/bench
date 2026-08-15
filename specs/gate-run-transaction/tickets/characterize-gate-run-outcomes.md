# Characterize the gate run's ordinary outcomes at the public seam

Blocked by: none
Writes: internal/gate (new test files only)

## What to build

Tests only, no production edit. A fixture built on `internal/gittest.RepoOnBranch`
with a committed `.gitignore` (ignoring the fixture's sentinels and the script's
run counter) and an executable `.bench/gate.sh` that increments a counter file,
copies `<git-dir>/bench-last-gate` aside while it runs, and switches behavior on
`.gitignore`d sentinels so the subject never moves between runs. Through
`gate.Execute`, `gate.RunCommand`, and `gate.Inspect` the tests pin: green
records and reuses without re-running (GC1), red records and invalidates (GC2),
`--fresh` re-runs (GC3), the oracle sees the pending record mid-run (GC7),
drift during execution refuses (GC8), no gate exits 3 (GC9), and an abandoned
pending record reloads as pending and does not block (GC11). GC10 is already
covered by the landing tests and is cited, not re-written. Each row lands
green against today's code; the ticket records, for each row, the named
mutation applied to a scratch copy of the tree and the observed red, so the
characterization is shown to bite. Grouping: no per-row red is stranded by a
thinner cut; the rows are grouped by the fixture capability they share (base
script + sentinels + record copy), and per-row tickets would cost one full-gate
commit each for tests-only landings — a reviewer call recorded in the approval
table. Shared contract this ticket publishes for its siblings: the fixture
constructor, the counter/sentinel/record-copy script conventions, and the
`Inspect`-shape assertions, reused unchanged by the contention ticket.
GC10's existing coverage (`internal/worktree/land_test.go:65-73`,
`internal/landing/landing_test.go`) is cited in the test file's header, not
re-written.

Return note (not acceptance): for every row, the mutation applied to a
scratch copy and the observed red.

## Acceptance

- [ ] Green run: exit 0/0, `Inspect` `Ready`/`green` reusable; second run prints the reuse line, counter stays 1, `recorded_at` bytes identical (covers GC1)
- [ ] Red run (exit 7): exit 7/7, `Inspect` `Ready`/`red`; a prior green for the same subject is not reused afterwards (covers GC2)
- [ ] `RunCommand([--fresh, root])` after a fresh green re-runs the script, counter 2 (covers GC3)
- [ ] The script's copy of the record mid-run parses as pending with the test pid and the inspected tree (covers GC7)
- [ ] Script appends to a tracked file: `ActionExit 1`, `GateExit` = script exit, stderr `gate subject changed during execution`, not reusable (covers GC8)
- [ ] No `.bench/gate.sh`, no `BENCH_GATE`: 3/3 and `no gate found` (covers GC9)
- [ ] Hand-written pending record with a dead pid: `Inspect` `Pending`; the next `Execute` runs the script and ends `Ready`/`green` (covers GC11)
- [ ] The prospective path's existing green/reusable coverage is cited from the test file and left in place (covers GC10)

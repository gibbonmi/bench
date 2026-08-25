# Review pickup: one-change-one-grade

Base `91443a2a`, reviewed tip `45581e89`. Three axes ran as read-only delegates.

## Standards

Count: 3 raw, 0 hard violations, 0 actionable. Worst issue: the lane's check
names appear in `CONTEXT.md` and `CHANGELOG.md` beside the one checked source,
the profile table. Rows OG28 and OG31 ask for that prose, so the finding is
`no-op`; the reviewer may veto.

## Spec

Count: 1 raw, 1 actionable. Worst issue: rows OG14 and OG23 say the test reads
the literal `gate: green`, but the landing fixture's gate-script route prints
no Bench-owned verdict line. Both tests read the script's tally instead, which
is the file's existing convention.

- OG14, OG23 — `internal/worktree/land_journey_test.go`,
  `internal/commit/lane_test.go` — non-behavioral seam substitution —
  `auto-fix`: reword the two rows to name the tally the tests read.

## Coverage

Count: 6 raw, 2 actionable. Worst issue: the timed-out lane path in
`gate.RunLane` has no test, and the edge inventory decides it writes nothing.

- A lane fail under `--dry-run` exits 1 and names the check — no test —
  `auto-fix`: add the test in `internal/commit/lane_test.go`.
- A timed-out lane writes no record and publishes nothing — no test —
  `auto-fix`: add the test in `internal/gate/lane_record_test.go` through the
  `gateTimeout` seam.

# Ignore the capture files in linked repos

Blocked by: none
Writes: internal/adopt/init.go, internal/adopt/adopt_test.go
Covers: HS24

## What to build

Verify the premise first: `Init` in internal/adopt/init.go scaffolds the gate,
the canary runner, and the learnings file, and touches no ignore file. Then
append `capture/session-handoff.md`, `capture/IDEAS.md`, and
`capture/learnings.md` to the root `.gitignore` when each line is absent,
creating the file when it does not exist. A second run adds nothing.

## Acceptance

- [ ] A fresh repo after `bench init` has the three entries in `.gitignore`.
- [ ] A second `bench init` leaves the file byte-identical.
- [ ] A repo whose `.gitignore` already holds one entry gains only the other two.
- [ ] Self-probe: append unconditionally, and report the second-run test red.

# Report a reduced verdict as its own board row

Blocked by: none

Ownership fence: `internal/status`
Assumptions: `gate.Inspection` already carries `Reduced`; the reduced record already records the current tree

## What to build

After a reduced run, `bench status` prints `stale (gated tree abc1234, work tree abc1234)
/ re-run the gate` — two identical hashes described as drift, and an action that cannot
clear, because re-running produces another reduced record and the same row. It is the
first thing an operator sees after every capture-only commit, which is the workflow this
feature exists to serve.

`GateInfo` never reads `Inspection.Reduced`, so a reduced green falls into the stale
branch. Carry it, and give a reduced verdict its own row that says what was graded and
what would widen it.

## Acceptance

- [ ] [RA1] A fresh reduced green renders a reduced row, not a stale row, and never prints two identical tree hashes as drift.
- [ ] [RA2] The row's action names the escape that widens the verdict rather than one that reproduces it.
- [ ] [RA3] A genuinely stale verdict, reduced or not, still renders the strong stale row.

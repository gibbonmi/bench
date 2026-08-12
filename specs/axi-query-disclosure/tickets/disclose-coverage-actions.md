# Disclose coverage check and retry actions

Blocked by: record-ft173-log-leverage.md
Writes: `internal/coverage/`

## What to build

Make the public `coverage` query derive `bench coverage --check <spec>` with the resolved spec path for unchecked maps and the exact retry after a named repair for malformed coverage. Derive per row, preserve row values and order, and keep complete and empty results honestly action-free.

## Acceptance

- [ ] [CH1] (covers QD1) every unchecked or repairable coverage row yields its own correctly valued check or retry action, including many-row and path-bearing fixtures that reject sampling and guessed values.
- [ ] [CH2] (covers QD1) complete and empty coverage results append the honest zero-row help block without a redundant check.
- [ ] [CH3] (covers QD6) old-to-new public-command fixtures prove that each named coverage state changes only by its appended help block; primary bytes, streams, exits, and argv behavior remain byte-equal.

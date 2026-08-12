# Disclose coverage check and retry actions

Blocked by: record-ft173-log-leverage.md, disclose-map-actions.md
Writes: `internal/coverage/`

## What to build

Make the successful default `coverage` extraction derive `bench coverage --check <spec>` with the resolved spec path for each mapped row and the exact repair-then-check retry for repairable mapped rows. Derive per row, preserve row values and order, and keep mapped zero-row, historical, and no-map results honestly action-free. `coverage --check` and its refusals retain their existing contracts and do not gain disclosure.

## Acceptance

- [ ] [CH1] (covers QD1) every mapped or repairable mapped row in a successful default extraction yields its own correctly valued check or repair-then-check retry action, including many-row and path-bearing fixtures that reject sampling and guessed values.
- [ ] [CH2] (covers QD1) complete and empty coverage results append the honest zero-row help block without a redundant check.
- [ ] [CH3] (covers QD6) old-to-new public-command fixtures separately cover mapped rows, repairable mapped rows, mapped zero-row, historical, and no-map default extraction; each changes only by its appended help block while primary bytes, streams, exits, and argv behavior remain byte-equal.

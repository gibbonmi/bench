# Disclose guard repair actions

Blocked by: record-ft173-log-leverage.md, disclose-map-actions.md
Writes: `internal/guards/`

## What to build

Make the public `guards` query derive `bench link` only for stale or unwired rows. Emit one action per matching row with stable-order deduplication, leave timeout or incomplete-scan states honest, and preserve the existing aggregate, primary table, streams, exits, and accepted argv.

## Acceptance

- [ ] [GH1] (covers QD1) every stale or unwired row contributes the `bench link` repair action, duplicate templates dedupe stably, and many-row fixtures reject sampled or first-match-only derivation.
- [ ] [GH2] (covers QD1) clean, empty, timeout, and incomplete-scan states append the honest zero-row help block rather than claiming `bench link` can repair them.
- [ ] [GH3] (covers QD6) old-to-new public-command fixtures prove that each named guards state changes only by its appended help block; primary bytes, streams, exits, aggregate facts, and argv behavior remain byte-equal.

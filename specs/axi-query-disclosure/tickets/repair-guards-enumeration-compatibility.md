# Repair guards enumeration compatibility

Blocked by: none
Writes: `internal/guards/guards_test.go`, `internal/guards/testdata`

## What to build

Add the missing public old/new compatibility evidence for an enumeration-timeout incomplete scan. Reuse production timeout seams; do not change guards behavior or infer repair work from unknown rows.

## Acceptance

- [ ] [GE1] (covers QD6) a checked-in pre-disclosure fixture captures the enumeration-timeout `Command` response with `total: unknown`, `omitted: unknown`, timeout reason, empty stderr, and exit 0.
- [ ] [GE2] (covers QD6) the candidate public response equals the checked primary bytes plus exactly `help[0]{cmd,why}:\n`; it does not borrow the partial-timeout fixture without an explicit alias.
- [ ] [GE3] (covers QD6) changing an enumeration-timeout primary byte, its exit, or its help cardinality makes the exact public fixture test red while the existing internal `Scan` metadata test remains green.


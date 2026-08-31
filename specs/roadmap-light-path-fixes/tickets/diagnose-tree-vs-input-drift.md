# Diagnose tree and declared-input drift

Blocked by: none
Writes: internal/gate/run_transaction.go, internal/gate/run_outcomes_test.go
Covers: LF1

## What to build

Make the gate transaction refusal identify the moved subject. Distinguish the
tree hash from each declared manifest path using transaction-owned facts.

## Acceptance

- [ ] Tree movement names the tree hash.
- [ ] Declared-input movement names the exact manifest path.
- [ ] Stable subjects keep the current gate outcome.


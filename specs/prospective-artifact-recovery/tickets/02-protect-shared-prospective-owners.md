# Protect shared prospective owners

Blocked by: 01-recover-dead-prospective-bundles.md
Writes: internal/gate/prospectiveartifact/, internal/gate/engine.go, internal/gate/lane.go, internal/gate/prospective.go, internal/gate/lane_test.go, internal/systemtest/owner_artifact_recovery_test.go

## What to build

Make the owner refuse destructive authority from live, invalid, and foreign state.
Apply the same owner convention to evidence inspection and fast-lane checkouts.
Permission-refused and unknown process probes retain their candidates.

The bundle policy classifies concurrent candidates independently.
The full producer table uses one record and cleanup contract.

## Acceptance

- [ ] PAR11: a concurrent authorization does not remove a blocked live owner's bundle.
- [ ] PAR12: an unsupported record schema retains its bundle root.
- [ ] PAR13: a missing, empty, or malformed record retains its bundle root.
- [ ] PAR14: a symbolic-link or special-file record retains its bundle root.
- [ ] PAR15: a symbolic-link or special-file bundle-root candidate remains unchanged.
- [ ] PAR16: a foreign same-prefix directory remains byte-for-byte unchanged.
- [ ] PAR17: a permission-refused process probe retains its bundle.
- [ ] PAR18: a non-ESRCH process-probe error retains its bundle.
- [ ] PAR19: an answering PID retains its bundle regardless of record age.
- [ ] PAR20: one sweep removes only the dead bundle from a dead-and-live pair.
- [ ] PAR29: the evidence-inspection producer publishes the valid owner record.
- [ ] PAR30: the fast-lane producer publishes the valid owner record.

# Recover dead prospective bundles

Blocked by: none
Writes: internal/gate/prospectiveartifact/ (new), internal/gate/engine.go, internal/gate/prospective.go, internal/gate/prospective_owner_test.go, internal/worktree/land_freshness_test.go, internal/systemtest/owner_artifact_recovery_test.go (new), projects/benchkit.md, CHANGELOG.md

## What to build

Create the prospective artifact owner and route full-gate execution through it.
The owner records the PID and repository before it creates the checkout.
It places every owner-authored run binary under the same bundle root.
It publishes the complete owner record before checkout creation starts.

A fresh authorization removes one recognized dead-owner bundle before it creates its own bundle.
Normal terminal outcomes close the current bundle without changing gate evidence.
Recovery handles partial bundle forms and executes no recovered candidate bytes.

## Acceptance

- [ ] PAR01: a fresh authorization removes the dead owner's checkout and run binary before execution.
- [ ] PAR02: another repository's bundle record cannot authorize deletion.
- [ ] PAR03: cleanup removes the exact Git worktree registration before the bundle root.
- [ ] PAR04: a second fresh-process sweep succeeds without changing another path.
- [ ] PAR05: recovery completes before the new prospective checkout exists.
- [ ] PAR06: the owner record is visible before checkout creation starts.
- [ ] PAR07: the published record has the exact fields, type, and mode.
- [ ] PAR08: recovery removes a dead record-only bundle.
- [ ] PAR09: recovery removes a dead checkout that has no run binary.
- [ ] PAR10: recovery removes a dead stale Git registration.
- [ ] PAR21: hostile path text stays confined to the exact bundle root.
- [ ] PAR22: a green prospective execution leaves no current bundle.
- [ ] PAR23: a red prospective execution leaves no current bundle.
- [ ] PAR24: a timed-out prospective execution leaves no current bundle.
- [ ] PAR25: a cancelled prospective execution leaves no current bundle.
- [ ] PAR26: a run-binary build refusal leaves no current bundle.
- [ ] PAR27: cleanup never removes an inherited baseline run binary.
- [ ] PAR28: the full-gate producer publishes the valid owner record.
- [ ] PAR31: dead-bundle recovery preserves exact green evidence reuse.
- [ ] PAR32: recovery never runs the planted old candidate executable.
- [ ] PAR33: a Git registration-removal failure refuses and retains the dead bundle root.

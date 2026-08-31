# Single-source the prospective checkout name

Blocked by: none
Writes: internal/gate/prospectiveartifact/prospectiveartifact.go, internal/gate/prospectiveartifact/prospectiveartifact_test.go, internal/systemtest/owner_artifact_recovery_test.go
Covers: LF21

## What to build

Expose one checkout child-name owner from prospectiveartifact. Make the system
journey and package tests consume it instead of retyping the value.

## Acceptance

- [ ] Production, package tests, and the system journey use one name owner.
- [ ] The system journey no longer carries three copied literals.
- [ ] Existing checkout layout remains compatible.


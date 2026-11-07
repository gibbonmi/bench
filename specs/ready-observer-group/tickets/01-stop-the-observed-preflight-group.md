# Stop the observed preflight group

Blocked by: none
Writes: internal/bounds/bounds.go, internal/conformance/injected_ports_test.go, internal/conformance/native_workflow_test.go
Covers: none

## What to build

The release-evidence readiness observer runs its preflight under the process-group owner with a limit below the test binary timeout.
A run that never reaches the readiness point makes the observer refuse, and no member of the run's group outlives that refusal.

## Acceptance

- [ ] A run that never reaches readiness makes the observer refuse within its limit.
- [ ] No member of the observed run's process group survives the refusal.
- [ ] The native release-evidence probe still observes the held run and releases it.

# Bound the release-evidence readiness hold

Blocked by: none
Writes: internal/bounds/bounds.go, internal/releaseevidence/release_evidence.go, internal/releaseevidence/protection_test.go
Covers: none

## What to build

A release-preflight run that holds at the evidence readiness point refuses with a named error when no observer releases the hold within a fixed window.
Before this fix, a run whose observer never came waited without bound, and one orphan run held a worktree for 8 hours.

## Acceptance

- [ ] A readiness hold that no observer releases returns a named readiness refusal within its window, also when a test run removes verdict deadlines.
- [ ] A held run that an observer releases still continues, and parent cancellation still ends the hold.
